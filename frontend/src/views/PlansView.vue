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
import { api, type CreatePlanPayload, type FeatureRecord, type PlanRecord } from "../api/api";
import { formatDateTime, getErrorMessage } from "../utils/format";

interface PlanSearchState {
  keyword: string;
  enabled?: boolean;
}

interface PlanFormState {
  id?: number;
  code: string;
  name: string;
  description: string;
  featureIds: number[];
  maxUsers: number;
  maxDevices: number;
  enabled: boolean;
}

const plans = ref<PlanRecord[]>([]);
const features = ref<FeatureRecord[]>([]);
const loading = ref(false);
const loadingFeatures = ref(false);
const saving = ref(false);
const drawerOpen = ref(false);
const drawerMode = ref<"create" | "edit">("create");
const searchState = reactive<PlanSearchState>({
  keyword: "",
  enabled: undefined,
});
const pagination = reactive<AdminPaginationState>({
  current: 1,
  pageSize: 10,
  total: 0,
});
const planForm = reactive<PlanFormState>({
  code: "",
  name: "",
  description: "",
  featureIds: [],
  maxUsers: 0,
  maxDevices: 1,
  enabled: true,
});

const searchFields: AdminSearchField[] = [
  {
    key: "keyword",
    label: "关键词",
    type: "input",
    placeholder: "套餐编码 / 名称 / 描述",
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

const planColumns = [
  { title: "套餐", key: "plan", dataIndex: "name" },
  { title: "描述", key: "description", dataIndex: "description" },
  { title: "编码", key: "code", dataIndex: "code" },
  { title: "状态", key: "enabled", dataIndex: "enabled" },
  { title: "功能", key: "features" },
  { title: "额度", key: "limits" },
  { title: "更新时间", key: "updated_at", dataIndex: "updated_at" },
  { title: "操作", key: "actions" },
];

const featureOptions = computed(() =>
  features.value.map((feature) => ({
    label: `${feature.name} · ${feature.code}`,
    value: feature.id,
    disabled: !feature.enabled,
  })),
);

const selectedFeatures = computed(() =>
  features.value.filter((feature) => planForm.featureIds.includes(feature.id)),
);

const filteredPlans = computed(() => {
  const keyword = searchState.keyword.trim().toLowerCase();
  return plans.value.filter((plan) => {
    const matchesKeyword =
      !keyword ||
      plan.code.toLowerCase().includes(keyword) ||
      plan.name.toLowerCase().includes(keyword) ||
      plan.description.toLowerCase().includes(keyword);
    const matchesEnabled =
      searchState.enabled === undefined || plan.enabled === searchState.enabled;
    return matchesKeyword && matchesEnabled;
  });
});

const pagedPlans = computed(() => {
  const start = (pagination.current - 1) * pagination.pageSize;
  return filteredPlans.value.slice(start, start + pagination.pageSize);
});

const drawerTitle = computed(() => (drawerMode.value === "create" ? "新建套餐" : "编辑套餐"));

watch(
  filteredPlans,
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

const loadPlans = async () => {
  loading.value = true;
  try {
    plans.value = await api.listPlans();
  } catch (error) {
    message.error(getErrorMessage(error, "加载套餐列表失败"));
  } finally {
    loading.value = false;
  }
};

const loadFeatures = async () => {
  loadingFeatures.value = true;
  try {
    features.value = await api.listFeatures();
  } catch (error) {
    message.error(getErrorMessage(error, "加载功能列表失败"));
  } finally {
    loadingFeatures.value = false;
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

const resetPlanForm = () => {
  planForm.id = undefined;
  planForm.code = "";
  planForm.name = "";
  planForm.description = "";
  planForm.featureIds = [];
  planForm.maxUsers = 0;
  planForm.maxDevices = 1;
  planForm.enabled = true;
};

const openCreateDrawer = () => {
  resetPlanForm();
  drawerMode.value = "create";
  drawerOpen.value = true;
};

const openEditDrawer = (plan: PlanRecord) => {
  planForm.id = plan.id;
  planForm.code = plan.code;
  planForm.name = plan.name;
  planForm.description = plan.description;
  planForm.featureIds = plan.features.map((feature) => feature.id);
  planForm.maxUsers = plan.max_users;
  planForm.maxDevices = plan.max_devices;
  planForm.enabled = plan.enabled;
  drawerMode.value = "edit";
  drawerOpen.value = true;
};

const submitPlan = async () => {
  const payload: CreatePlanPayload = {
    code: planForm.code.trim(),
    name: planForm.name.trim(),
    description: planForm.description.trim(),
    feature_ids: planForm.featureIds,
    max_users: Number(planForm.maxUsers) || 0,
    max_devices: Number(planForm.maxDevices) || 0,
    enabled: planForm.enabled,
  };
  if (!payload.code) {
    message.warning("请输入套餐编码");
    return;
  }
  if (!payload.name) {
    message.warning("请输入套餐名称");
    return;
  }
  if (payload.max_devices < 1) {
    message.warning("最大设备数至少为 1");
    return;
  }
  saving.value = true;
  try {
    if (drawerMode.value === "create") {
      await api.createPlan(payload);
      message.success("套餐已创建");
    } else if (planForm.id) {
      await api.updatePlan(planForm.id, payload);
      message.success("套餐已更新");
    }
    drawerOpen.value = false;
    await loadPlans();
  } catch (error) {
    message.error(getErrorMessage(error));
  } finally {
    saving.value = false;
  }
};

const deletePlan = async (plan: PlanRecord) => {
  try {
    await api.deletePlan(plan.id);
    message.success(`已删除套餐「${plan.name}」`);
    await loadPlans();
  } catch (error) {
    message.error(getErrorMessage(error));
  }
};

const planActions = (plan: PlanRecord): AdminRowAction<PlanRecord>[] => [
  {
    key: "edit",
    label: "编辑",
    onClick: () => openEditDrawer(plan),
  },
  {
    key: "delete",
    label: "删除",
    danger: true,
    confirm: `确定删除套餐「${plan.name}」？已签发 License 会保留套餐快照。`,
    confirmOkText: "删除",
    onClick: () => deletePlan(plan),
  },
];

loadFeatures();
loadPlans();
</script>

<template>
  <div class="flex flex-col gap-4">
    <AdminResourcePage title="套餐管理">
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
        <a-button type="primary" @click="openCreateDrawer">
          <template #icon>
            <PlusOutlined />
          </template>
          新建套餐
        </a-button>
        <a-button @click="loadPlans">刷新</a-button>
      </template>

      <AdminDataTable
        :columns="planColumns"
        :data-source="pagedPlans"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="onPageChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'plan'">
            <a-typography-text strong>{{ record.name }}</a-typography-text>
          </template>
          <template v-else-if="column.key === 'description'">
            {{ record.description || "暂无描述" }}
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
          <template v-else-if="column.key === 'features'">
            <a-space class="whitespace-nowrap">
              <a-tag v-for="feature in record.features" :key="feature.id">{{ feature.name }}</a-tag>
              <a-typography-text v-if="!record.features.length" type="secondary"
                >无</a-typography-text
              >
            </a-space>
          </template>
          <template v-else-if="column.key === 'limits'">
            <a-space class="whitespace-nowrap">
              <a-tag color="blue">用户 {{ record.max_users || "不限" }}</a-tag>
              <a-tag color="green">设备 {{ record.max_devices }}</a-tag>
            </a-space>
          </template>
          <template v-else-if="column.key === 'updated_at'">
            {{ formatDateTime(record.updated_at) }}
          </template>
          <template v-else-if="column.key === 'actions'">
            <AdminRowActions :record="record" :actions="planActions(record)" />
          </template>
        </template>
      </AdminDataTable>
    </AdminResourcePage>

    <a-drawer v-model:open="drawerOpen" :size="760" :title="drawerTitle" :mask-closable="!saving">
      <template #extra>
        <a-space>
          <a-button :disabled="saving" @click="drawerOpen = false">取消</a-button>
          <a-button type="primary" :loading="saving" @click="submitPlan">保存</a-button>
        </a-space>
      </template>
      <a-form layout="vertical" :model="planForm">
        <a-row :gutter="16">
          <a-col :xs="24" :md="12">
            <a-form-item label="套餐编码" required>
              <a-input v-model:value="planForm.code" placeholder="例如：enterprise" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :md="12">
            <a-form-item label="套餐名称" required>
              <a-input v-model:value="planForm.name" placeholder="例如：企业版" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="套餐描述">
          <a-textarea
            v-model:value="planForm.description"
            :auto-size="{ minRows: 3, maxRows: 5 }"
            placeholder="说明该套餐适用的客户或授权边界"
          />
        </a-form-item>
        <a-form-item label="绑定功能">
          <a-select
            v-model:value="planForm.featureIds"
            mode="multiple"
            :options="featureOptions"
            :loading="loadingFeatures"
            placeholder="选择启用的功能点"
          />
        </a-form-item>
        <a-form-item label="已选功能预览">
          <a-space wrap>
            <a-tag v-for="feature in selectedFeatures" :key="feature.id">{{ feature.name }}</a-tag>
            <a-typography-text v-if="!selectedFeatures.length" type="secondary"
              >暂无功能</a-typography-text
            >
          </a-space>
        </a-form-item>
        <a-row :gutter="16">
          <a-col :xs="24" :md="12">
            <a-form-item label="最大用户数" extra="0 表示不限制用户数">
              <a-input-number v-model:value="planForm.maxUsers" class="w-full" :min="0" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :md="12">
            <a-form-item label="最大设备数" required>
              <a-input-number v-model:value="planForm.maxDevices" class="w-full" :min="1" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="状态">
          <a-switch
            v-model:checked="planForm.enabled"
            checked-children="启用"
            un-checked-children="停用"
          />
        </a-form-item>
      </a-form>
    </a-drawer>
  </div>
</template>
