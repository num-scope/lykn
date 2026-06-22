<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { useBoolean } from "@sa/hooks";
import { useAntdvForm } from "@/hooks/common/form";
import { useAntdvPaginatedTable, useTableOperate } from "@/hooks/common/table";
import { createPlan, deletePlan, fetchFeatures, fetchPlans, updatePlan } from "@/service/api";
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
const tableScrollX = 1504;

const drawerOpen = ref(false);
const submitting = ref(false);
const operateType = ref<AntdvUI.TableOperateType>("add");
const editingId = ref<number>();
const features = ref<Api.Lykn.FeatureRecord[]>([]);
const featureLoading = ref(false);

const formModel = reactive<Api.Lykn.PlanPayload>({
  code: "",
  name: "",
  description: "",
  feature_ids: [],
  max_users: 0,
  max_devices: 1,
  enabled: true,
});

const enabledOptions = [
  { label: "启用", value: "true" },
  { label: "停用", value: "false" },
];

const searchItems = computed(() => [
  { key: "keyword", label: "关键词", placeholder: "套餐编码 / 名称 / 描述" },
  {
    key: "enabled",
    label: "状态",
    type: "select" as const,
    placeholder: "请选择状态",
    options: enabledOptions,
  },
]);

const featureOptions = computed(() =>
  features.value.map((item) => ({
    label: `${item.name} · ${item.code}`,
    value: item.id,
    disabled: !item.enabled,
  })),
);

const selectedFeatures = computed(() =>
  features.value.filter((item) => formModel.feature_ids.includes(item.id)),
);

const formRules = {
  code: [{ required: true, message: "请输入套餐编码" }],
  name: [{ required: true, message: "请输入套餐名称" }],
  max_devices: [{ required: true, message: "请输入设备额度" }],
};

const { loading, data, columns, columnChecks, getData, getDataByPage, mobilePagination } =
  useAntdvPaginatedTable({
    api: fetchPlans,
    transform: (response) => {
      const records = unwrapFlatData(response, "加载套餐列表失败") as Api.Lykn.PlanRecord[];
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
      { key: "name", dataIndex: "name", title: "套餐", align: "left", width: 180, fixed: "left" },
      { key: "code", dataIndex: "code", title: "编码", align: "left", width: 160 },
      { key: "enabled", dataIndex: "enabled", title: "状态", align: "center", width: 100 },
      { key: "features", dataIndex: "features", title: "功能", align: "left", width: 260 },
      { key: "limits", title: "额度", align: "center", width: 160 },
      { key: "description", dataIndex: "description", title: "描述", align: "left", width: 260 },
      {
        key: "updated_at",
        dataIndex: "updated_at",
        title: "更新时间",
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

const drawerTitle = computed(() => (operateType.value === "add" ? "新建套餐" : "编辑套餐"));

async function loadFeatures() {
  featureLoading.value = true;
  try {
    features.value = unwrapFlatData(await fetchFeatures(), "加载功能列表失败");
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  } finally {
    featureLoading.value = false;
  }
}

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
  formModel.feature_ids = [];
  formModel.max_users = 0;
  formModel.max_devices = 1;
  formModel.enabled = true;
  void restoreValidation();
}

function openCreateDrawer() {
  resetForm();
  operateType.value = "add";
  drawerOpen.value = true;
}

function openEditDrawer(record: Api.Lykn.PlanRecord) {
  operateType.value = "edit";
  editingId.value = record.id;
  formModel.code = record.code;
  formModel.name = record.name;
  formModel.description = toPlan(record).description;
  formModel.feature_ids = toPlan(record).features.map((item) => item.id);
  formModel.max_users = toPlan(record).max_users;
  formModel.max_devices = toPlan(record).max_devices;
  formModel.enabled = toPlan(record).enabled;
  drawerOpen.value = true;
}

async function submitForm() {
  await validate();

  submitting.value = true;
  try {
    const payload: Api.Lykn.PlanPayload = {
      ...formModel,
      max_users: Number(formModel.max_users) || 0,
      max_devices: Number(formModel.max_devices) || 0,
    };

    if (operateType.value === "add") {
      unwrapFlatData(await createPlan(payload), "创建套餐失败");
      window.$message?.success("套餐已创建");
    } else if (editingId.value) {
      unwrapFlatData(await updatePlan(editingId.value, payload), "更新套餐失败");
      window.$message?.success("套餐已更新");
    }

    drawerOpen.value = false;
    await getDataByPage();
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  } finally {
    submitting.value = false;
  }
}

function toPlan(record: unknown) {
  return record as Api.Lykn.PlanRecord;
}

async function handleBatchDelete() {
  setBatchDeleting(true);
  try {
    await Promise.all(checkedRowKeys.value.map((id) => deletePlan(Number(id))));
    await onBatchDeleted();
  } catch (error) {
    checkedRowKeys.value = [];
    window.$message?.error(getServiceErrorMessage(error, "批量删除套餐失败"));
    await getDataByPage();
  } finally {
    setBatchDeleting(false);
  }
}

async function removePlan(record: Api.Lykn.PlanRecord) {
  try {
    unwrapFlatData(await deletePlan(record.id), "删除套餐失败");
    window.$message?.success("套餐已删除");
    await getDataByPage();
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  }
}

loadFeatures();
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
      title="套餐管理"
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
            <AButton size="small" ghost type="primary" @click="openCreateDrawer">
              <template #icon>
                <icon-ic-round-plus class="text-icon" />
              </template>
              新建套餐
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
              <ATag :color="toPlan(record).enabled ? 'success' : 'default'">{{
                toPlan(record).enabled ? "启用" : "停用"
              }}</ATag>
            </template>
            <template v-else-if="column.key === 'features'">
              <ASpace wrap>
                <ATag v-for="item in toPlan(record).features" :key="item.id" color="processing">{{
                  item.name
                }}</ATag>
                <ATypographyText v-if="!toPlan(record).features.length" type="secondary"
                  >未绑定</ATypographyText
                >
              </ASpace>
            </template>
            <template v-else-if="column.key === 'limits'">
              用户 {{ toPlan(record).max_users || "不限" }} / 设备
              {{ toPlan(record).max_devices || "-" }}
            </template>
            <template v-else-if="column.key === 'description'">
              <ATypographyText type="secondary">{{
                toPlan(record).description || "-"
              }}</ATypographyText>
            </template>
            <template v-else-if="column.key === 'updated_at'">
              {{ formatDateTime(toPlan(record).updated_at) }}
            </template>
            <template v-else-if="column.key === 'operate'">
              <div class="flex-center justify-end gap-6px">
                <AButton type="primary" ghost size="small" @click="openEditDrawer(toPlan(record))"
                  >编辑</AButton
                >
                <APopconfirm title="确定删除该套餐？" @confirm="removePlan(toPlan(record))">
                  <AButton danger ghost size="small">删除</AButton>
                </APopconfirm>
              </div>
            </template>
          </template>
      </ATable>
    </ACard>

    <ADrawer v-model:open="drawerOpen" :title="drawerTitle" :size="560" destroy-on-hidden>
      <template #extra>
        <ASpace>
          <AButton @click="drawerOpen = false">取消</AButton>
          <AButton type="primary" :loading="submitting" @click="submitForm">保存</AButton>
        </ASpace>
      </template>

      <AForm ref="formRef" layout="vertical" :model="formModel" :rules="formRules">
        <AFormItem label="套餐编码" name="code">
          <AInput v-model:value="formModel.code" placeholder="例如：pro" />
        </AFormItem>
        <AFormItem label="套餐名称" name="name">
          <AInput v-model:value="formModel.name" placeholder="请输入套餐名称" />
        </AFormItem>
        <AFormItem label="套餐描述" name="description">
          <ATextarea
            v-model:value="formModel.description"
            :auto-size="{ minRows: 3, maxRows: 5 }"
            placeholder="请输入套餐描述"
          />
        </AFormItem>
        <AFormItem label="包含功能" name="feature_ids">
          <ASelect
            v-model:value="formModel.feature_ids"
            mode="multiple"
            allow-clear
            :loading="featureLoading"
            :options="featureOptions"
            placeholder="请选择功能"
          />
        </AFormItem>
        <AFormItem v-if="selectedFeatures.length" label="已选功能">
          <ASpace wrap>
            <ATag v-for="item in selectedFeatures" :key="item.id" color="processing">{{
              item.name
            }}</ATag>
          </ASpace>
        </AFormItem>
        <AFormItem label="用户额度" name="max_users" extra="0 表示不限制用户数">
          <AInputNumber v-model:value="formModel.max_users" class="w-full" :min="0" />
        </AFormItem>
        <AFormItem label="设备额度" name="max_devices">
          <AInputNumber v-model:value="formModel.max_devices" class="w-full" :min="1" />
        </AFormItem>
        <AFormItem label="状态" name="enabled">
          <ASwitch
            v-model:checked="formModel.enabled"
            checked-children="启用"
            un-checked-children="停用"
          />
        </AFormItem>
      </AForm>
    </ADrawer>
  </div>
</template>
