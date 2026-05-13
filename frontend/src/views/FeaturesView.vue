<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import { message } from "antdv-next";
import { PlusOutlined } from "@antdv-next/icons";

import {
  AdminDataTable,
  AdminResourcePage,
  AdminRowActions,
  AdminSearchForm,
  AdminStatusTag,
  type AdminPaginationState,
  type AdminRowAction,
  type AdminSearchField,
} from "../components/admin-kit";
import { api, type CreateFeaturePayload, type FeatureRecord } from "../api/api";
import { formatDateTime, getErrorMessage } from "../utils/format";

interface FeatureSearchState {
  keyword: string;
  enabled?: boolean;
}

interface FeatureFormState {
  id?: number;
  code: string;
  name: string;
  description: string;
  enabled: boolean;
}

const features = ref<FeatureRecord[]>([]);
const loading = ref(false);
const saving = ref(false);
const modalOpen = ref(false);
const modalMode = ref<"create" | "edit">("create");
const searchState = reactive<FeatureSearchState>({
  keyword: "",
  enabled: undefined,
});
const pagination = reactive<AdminPaginationState>({
  current: 1,
  pageSize: 10,
  total: 0,
});
const featureForm = reactive<FeatureFormState>({
  code: "",
  name: "",
  description: "",
  enabled: true,
});

const searchFields: AdminSearchField[] = [
  {
    key: "keyword",
    label: "关键词",
    type: "input",
    placeholder: "功能编码 / 名称 / 描述",
  },
  {
    key: "enabled",
    label: "状态",
    type: "select",
    placeholder: "全部状态",
    options: [
      { label: "启用", value: true },
      { label: "停用", value: false },
    ],
  },
];

const featureColumns = [
  { title: "功能", key: "feature", dataIndex: "name" },
  { title: "编码", key: "code", dataIndex: "code" },
  { title: "状态", key: "enabled", dataIndex: "enabled" },
  { title: "描述", key: "description", dataIndex: "description" },
  { title: "更新时间", key: "updated_at", dataIndex: "updated_at" },
  { title: "创建时间", key: "created_at", dataIndex: "created_at" },
  { title: "操作", key: "actions" },
];

const filteredFeatures = computed(() => {
  const keyword = searchState.keyword.trim().toLowerCase();
  return features.value.filter((feature) => {
    const matchesKeyword =
      !keyword ||
      feature.code.toLowerCase().includes(keyword) ||
      feature.name.toLowerCase().includes(keyword) ||
      feature.description.toLowerCase().includes(keyword);
    const matchesEnabled =
      searchState.enabled === undefined || feature.enabled === searchState.enabled;
    return matchesKeyword && matchesEnabled;
  });
});

const pagedFeatures = computed(() => {
  const start = (pagination.current - 1) * pagination.pageSize;
  return filteredFeatures.value.slice(start, start + pagination.pageSize);
});

const modalTitle = computed(() => (modalMode.value === "create" ? "新建功能" : "编辑功能"));

watch(
  filteredFeatures,
  (items) => {
    pagination.total = items.length;
    const maxPage = Math.max(1, Math.ceil(items.length / pagination.pageSize));
    if (pagination.current > maxPage) pagination.current = maxPage;
  },
  { immediate: true },
);

watch(
  () => [searchState.keyword, searchState.enabled],
  () => {
    pagination.current = 1;
  },
);

const loadFeatures = async () => {
  loading.value = true;
  try {
    features.value = await api.listFeatures();
  } catch (error) {
    message.error(getErrorMessage(error, "加载功能列表失败"));
  } finally {
    loading.value = false;
  }
};

const onSearchReset = () => {
  searchState.keyword = "";
  searchState.enabled = undefined;
  pagination.current = 1;
};

const onPageChange = (page: number, pageSize: number) => {
  pagination.current = page;
  pagination.pageSize = pageSize;
};

const resetFeatureForm = () => {
  featureForm.id = undefined;
  featureForm.code = "";
  featureForm.name = "";
  featureForm.description = "";
  featureForm.enabled = true;
};

const openCreateModal = () => {
  resetFeatureForm();
  modalMode.value = "create";
  modalOpen.value = true;
};

const openEditModal = (feature: FeatureRecord) => {
  featureForm.id = feature.id;
  featureForm.code = feature.code;
  featureForm.name = feature.name;
  featureForm.description = feature.description;
  featureForm.enabled = feature.enabled;
  modalMode.value = "edit";
  modalOpen.value = true;
};

const submitFeature = async () => {
  const payload: CreateFeaturePayload = {
    code: featureForm.code.trim(),
    name: featureForm.name.trim(),
    description: featureForm.description.trim(),
    enabled: featureForm.enabled,
  };
  if (!payload.code) {
    message.warning("请输入功能编码");
    return;
  }
  if (!payload.name) {
    message.warning("请输入功能名称");
    return;
  }
  saving.value = true;
  try {
    if (modalMode.value === "create") {
      await api.createFeature(payload);
      message.success("功能已创建");
    } else if (featureForm.id) {
      await api.updateFeature(featureForm.id, payload);
      message.success("功能已更新");
    }
    modalOpen.value = false;
    await loadFeatures();
  } catch (error) {
    message.error(getErrorMessage(error));
  } finally {
    saving.value = false;
  }
};

const deleteFeature = async (feature: FeatureRecord) => {
  try {
    await api.deleteFeature(feature.id);
    message.success(`已删除功能「${feature.name}」`);
    await loadFeatures();
  } catch (error) {
    message.error(getErrorMessage(error));
  }
};

const featureActions = (feature: FeatureRecord): AdminRowAction<FeatureRecord>[] => [
  {
    key: "edit",
    label: "编辑",
    onClick: () => openEditModal(feature),
  },
  {
    key: "delete",
    label: "删除",
    danger: true,
    confirm: `确定删除功能「${feature.name}」？已被套餐使用的功能不能删除。`,
    confirmOkText: "删除",
    onClick: () => deleteFeature(feature),
  },
];

loadFeatures();
</script>

<template>
  <div class="flex flex-col gap-4">
    <AdminResourcePage title="功能管理">
      <template #search>
        <AdminSearchForm
          :model="searchState"
          :fields="searchFields"
          :field-count="searchFields.length"
          @search="pagination.current = 1"
          @reset="onSearchReset"
        />
      </template>

      <template #toolbar>
        <a-button type="primary" @click="openCreateModal">
          <template #icon>
            <PlusOutlined />
          </template>
          新建功能
        </a-button>
        <a-button @click="loadFeatures">刷新</a-button>
      </template>

      <AdminDataTable
        :columns="featureColumns"
        :data-source="pagedFeatures"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="onPageChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'feature'">
            <a-typography-text strong>{{ record.name }}</a-typography-text>
          </template>
          <template v-else-if="column.key === 'code'">
            <a-typography-text copyable class="font-mono text-xs">{{
              record.code
            }}</a-typography-text>
          </template>
          <template v-else-if="column.key === 'enabled'">
            <AdminStatusTag
              :value="record.enabled"
              :active-value="true"
              active-text="启用"
              inactive-text="停用"
              inactive-color="default"
            />
          </template>
          <template v-else-if="column.key === 'description'">
            {{ record.description || "暂无描述" }}
          </template>
          <template v-else-if="column.key === 'updated_at'">
            {{ formatDateTime(record.updated_at) }}
          </template>
          <template v-else-if="column.key === 'created_at'">
            {{ formatDateTime(record.created_at) }}
          </template>
          <template v-else-if="column.key === 'actions'">
            <AdminRowActions :record="record" :actions="featureActions(record)" />
          </template>
        </template>
      </AdminDataTable>
    </AdminResourcePage>

    <a-modal
      v-model:open="modalOpen"
      :title="modalTitle"
      :confirm-loading="saving"
      @ok="submitFeature"
    >
      <a-form layout="vertical" :model="featureForm">
        <a-form-item label="功能编码" required>
          <a-input v-model:value="featureForm.code" placeholder="例如：reports" />
        </a-form-item>
        <a-form-item label="功能名称" required>
          <a-input v-model:value="featureForm.name" placeholder="例如：报表功能" />
        </a-form-item>
        <a-form-item label="功能描述">
          <a-textarea
            v-model:value="featureForm.description"
            :auto-size="{ minRows: 3, maxRows: 5 }"
            placeholder="说明该功能控制的业务能力"
          />
        </a-form-item>
        <a-form-item label="状态">
          <a-switch
            v-model:checked="featureForm.enabled"
            checked-children="启用"
            un-checked-children="停用"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
