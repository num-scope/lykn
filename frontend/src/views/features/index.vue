<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { useBoolean } from "@sa/hooks";
import { useAntdvForm } from "@/hooks/common/form";
import { useAntdvPaginatedTable, useTableOperate } from "@/hooks/common/table";
import { createFeature, deleteFeature, fetchFeatures, updateFeature } from "@/service/api";
import { formatDateTime, getServiceErrorMessage, unwrapFlatData } from "@/utils/lykn";

const searchModel = reactive({
  keyword: "",
  enabled: undefined as "true" | "false" | undefined,
});

const { formRef, validate, restoreValidation } = useAntdvForm();
const paginationParams = reactive({
  current: 1,
  size: 10,
});
const tableScrollX = 1284;

const modalOpen = ref(false);
const submitting = ref(false);
const operateType = ref<AntdvUI.TableOperateType>("add");
const editingId = ref<number>();

const formModel = reactive<Api.Lykn.FeaturePayload>({
  code: "",
  name: "",
  description: "",
  enabled: true,
});

const enabledOptions = [
  { label: "启用", value: "true" },
  { label: "停用", value: "false" },
];

const searchItems = computed(() => [
  { key: "keyword", label: "关键词", placeholder: "功能编码 / 名称 / 描述" },
  {
    key: "enabled",
    label: "状态",
    type: "select" as const,
    placeholder: "请选择状态",
    options: enabledOptions,
  },
]);

const formRules = {
  code: [{ required: true, message: "请输入功能编码" }],
  name: [{ required: true, message: "请输入功能名称" }],
};

const { loading, data, columns, columnChecks, getData, getDataByPage, mobilePagination } =
  useAntdvPaginatedTable({
    api: fetchFeatures,
    transform: (response) => {
      const records = unwrapFlatData(response, "加载功能列表失败") as Api.Lykn.FeatureRecord[];
      const keyword = searchModel.keyword.trim().toLowerCase();

      const filteredRecords = records.filter((item) => {
        const matchKeyword =
          !keyword ||
          item.code.toLowerCase().includes(keyword) ||
          item.name.toLowerCase().includes(keyword) ||
          item.description.toLowerCase().includes(keyword);
        const matchEnabled =
          searchModel.enabled === undefined || String(item.enabled) === searchModel.enabled;

        return matchKeyword && matchEnabled;
      });

      return {
        data: filteredRecords,
        pageNum: paginationParams.current,
        pageSize: paginationParams.size,
        total: filteredRecords.length,
      };
    },
    onPaginationParamsChange: (params) => {
      paginationParams.current = params.current || 1;
      paginationParams.size = params.pageSize || 10;
    },
    columns: () => [
      { key: "index", title: "序号", align: "center", width: 64, fixed: "left" },
      { key: "name", dataIndex: "name", title: "功能", align: "left", width: 180, fixed: "left" },
      { key: "code", dataIndex: "code", title: "编码", align: "left", width: 160 },
      { key: "enabled", dataIndex: "enabled", title: "状态", align: "center", width: 100 },
      { key: "description", dataIndex: "description", title: "描述", align: "left", width: 280 },
      {
        key: "updated_at",
        dataIndex: "updated_at",
        title: "更新时间",
        align: "center",
        width: 180,
      },
      {
        key: "created_at",
        dataIndex: "created_at",
        title: "创建时间",
        align: "center",
        width: 180,
      },
      { key: "operate", title: "操作", align: "center", width: 140, fixed: "right" },
    ],
  });

const { checkedRowKeys, onBatchDeleted } = useTableOperate(data, "id", getData);
const { bool: batchDeleting, setBool: setBatchDeleting } = useBoolean();
const rowSelection = computed(() => ({
  fixed: true,
  selectedRowKeys: checkedRowKeys.value,
  onChange: (keys: (string | number)[]) => {
    checkedRowKeys.value = keys as string[];
  },
}));

const modalTitle = computed(() => (operateType.value === "add" ? "新建功能" : "编辑功能"));

function handleSearch() {
  getDataByPage();
}

function resetSearch() {
  searchModel.keyword = "";
  searchModel.enabled = undefined;
  getDataByPage();
}

function resetForm() {
  editingId.value = undefined;
  formModel.code = "";
  formModel.name = "";
  formModel.description = "";
  formModel.enabled = true;
  void restoreValidation();
}

function openCreateModal() {
  resetForm();
  operateType.value = "add";
  modalOpen.value = true;
}

function openEditModal(record: Api.Lykn.FeatureRecord) {
  operateType.value = "edit";
  editingId.value = record.id;
  formModel.code = record.code;
  formModel.name = record.name;
  formModel.description = toFeature(record).description;
  formModel.enabled = toFeature(record).enabled;
  modalOpen.value = true;
}

async function submitForm() {
  await validate();

  submitting.value = true;
  try {
    if (operateType.value === "add") {
      unwrapFlatData(await createFeature({ ...formModel }), "创建功能失败");
      window.$message?.success("功能已创建");
    } else if (editingId.value) {
      unwrapFlatData(await updateFeature(editingId.value, { ...formModel }), "更新功能失败");
      window.$message?.success("功能已更新");
    }

    modalOpen.value = false;
    await getDataByPage();
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  } finally {
    submitting.value = false;
  }
}

function toFeature(record: unknown) {
  return record as Api.Lykn.FeatureRecord;
}

async function handleBatchDelete() {
  setBatchDeleting(true);
  try {
    await Promise.all(checkedRowKeys.value.map((id) => deleteFeature(Number(id))));
    await onBatchDeleted();
  } catch (error) {
    checkedRowKeys.value = [];
    window.$message?.error(getServiceErrorMessage(error, "批量删除功能失败"));
    await getDataByPage();
  } finally {
    setBatchDeleting(false);
  }
}

async function removeFeature(record: Api.Lykn.FeatureRecord) {
  try {
    unwrapFlatData(await deleteFeature(record.id), "删除功能失败");
    window.$message?.success("功能已删除");
    await getDataByPage();
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  }
}
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <QuerySearchForm
      v-model:model="searchModel"
      :items="searchItems"
      :loading="loading"
      @search="handleSearch"
      @reset="resetSearch"
    />

    <ACard
      title="功能管理"
      variant="borderless"
      :body-style="{ flex: 1, overflow: 'hidden' }"
      class="flex-col-stretch sm:flex-1-hidden card-wrapper"
    >
      <template #extra>
        <TableHeaderOperation
          v-model:columns="columnChecks"
          :disabled-delete="checkedRowKeys.length === 0"
          :loading="loading"
          @refresh="getData"
        >
          <template #default>
            <AButton size="small" ghost type="primary" @click="openCreateModal">
              <template #icon>
                <icon-ic-round-plus class="text-icon" />
              </template>
              新建功能
            </AButton>
            <APopconfirm title="确认删除吗？" @confirm="handleBatchDelete">
              <AButton
                size="small"
                ghost
                danger
                :disabled="checkedRowKeys.length === 0"
                :loading="batchDeleting"
              >
                <template #icon>
                  <icon-ic-round-delete class="text-icon" />
                </template>
                批量删除
              </AButton>
            </APopconfirm>
          </template>
        </TableHeaderOperation>
      </template>

      <ATable
        row-key="id"
        :row-selection="rowSelection"
        size="small"
        :columns="columns"
        :data-source="data"
        :loading="loading"
        :scroll="{ x: tableScrollX }"
        :pagination="mobilePagination"
        class="h-full"
      >
          <template #bodyCell="{ column, record, index }">
            <template v-if="column.key === 'index'">
              {{ index + 1 }}
            </template>
            <template v-else-if="column.key === 'enabled'">
              <ATag :color="toFeature(record).enabled ? 'success' : 'default'">{{
                toFeature(record).enabled ? "启用" : "停用"
              }}</ATag>
            </template>
            <template v-else-if="column.key === 'description'">
              <ATypographyText type="secondary">{{
                toFeature(record).description || "-"
              }}</ATypographyText>
            </template>
            <template v-else-if="column.key === 'updated_at'">
              {{ formatDateTime(toFeature(record).updated_at) }}
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatDateTime(toFeature(record).created_at) }}
            </template>
            <template v-else-if="column.key === 'operate'">
              <div class="flex-center justify-end gap-6px">
                <AButton type="primary" ghost size="small" @click="openEditModal(toFeature(record))"
                  >编辑</AButton
                >
                <APopconfirm
                  title="确定删除该功能？已被套餐使用的功能不能删除"
                  @confirm="removeFeature(toFeature(record))"
                >
                  <AButton danger ghost size="small">删除</AButton>
                </APopconfirm>
              </div>
            </template>
          </template>
      </ATable>
    </ACard>

    <AModal
      v-model:open="modalOpen"
      :title="modalTitle"
      :confirm-loading="submitting"
      destroy-on-hidden
      @ok="submitForm"
    >
      <AForm ref="formRef" layout="vertical" :model="formModel" :rules="formRules">
        <AFormItem label="功能编码" name="code">
          <AInput v-model:value="formModel.code" placeholder="例如：advanced-export" />
        </AFormItem>
        <AFormItem label="功能名称" name="name">
          <AInput v-model:value="formModel.name" placeholder="请输入功能名称" />
        </AFormItem>
        <AFormItem label="功能描述" name="description">
          <ATextarea
            v-model:value="formModel.description"
            :auto-size="{ minRows: 3, maxRows: 5 }"
            placeholder="请输入功能描述"
          />
        </AFormItem>
        <AFormItem label="状态" name="enabled">
          <ASwitch
            v-model:checked="formModel.enabled"
            checked-children="启用"
            un-checked-children="停用"
          />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>
