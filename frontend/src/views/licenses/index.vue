<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import type { Dayjs } from "dayjs";
import { useAntdvForm } from "@/hooks/common/form";
import { useAntdvPaginatedTable } from "@/hooks/common/table";
import { fetchPlans, fetchProjectLicenses, issueLicense, downloadLicense } from "@/service/api";
import { useAuthStore } from "@/store/modules/auth";
import {
  formatDateTime,
  formatPeriod,
  getLicenseState,
  getServiceErrorMessage,
  saveBlobFile,
  unwrapFlatData,
} from "@/utils/lykn";

const authStore = useAuthStore();

const searchModel = reactive({
  projectId: null as number | null,
  keyword: "",
  state: undefined as "active" | "upcoming" | "expired" | undefined,
});

const { formRef, validate, restoreValidation } = useAntdvForm();
const paginationParams = reactive({
  current: 1,
  size: 10,
});
const tableScrollX = 1734;
const plans = ref<Api.Lykn.PlanRecord[]>([]);
const planLoading = ref(false);
const issuing = ref(false);
const drawerOpen = ref(false);
const detailOpen = ref(false);
const detailRecord = ref<Api.Lykn.LicenseRecord>();
const downloadingId = ref<number>();

const formModel = reactive({
  subjectName: "",
  subjectEmail: "",
  subjectOrg: "",
  planId: undefined as number | undefined,
  validRange: null as [Dayjs, Dayjs] | null,
  hardwareKeys: [] as Array<keyof Api.Lykn.LicenseHardwarePayload>,
  hostname: "",
  cpu_id: "",
  disk_serial: "",
  mac_addresses: "",
});

const stateOptions = [
  { label: "生效中", value: "active" },
  { label: "未开始", value: "upcoming" },
  { label: "已过期", value: "expired" },
];

const searchItems = computed(() => [
  { key: "projectId", label: "项目", type: "select" as const, placeholder: "请选择项目" },
  { key: "keyword", label: "关键词", placeholder: "客户 / 邮箱 / UUID / 套餐" },
  {
    key: "state",
    label: "有效性",
    type: "select" as const,
    placeholder: "请选择有效性",
    options: stateOptions,
  },
]);

const hardwareOptions = [
  { label: "主机名", value: "hostname" },
  { label: "CPU ID", value: "cpu_id" },
  { label: "磁盘序列号", value: "disk_serial" },
  { label: "MAC 地址", value: "mac_addresses" },
];

const formRules = {
  subjectName: [{ required: true, message: "请输入授权对象" }],
  planId: [{ required: true, message: "请选择套餐" }],
  validRange: [{ required: true, message: "请选择有效期" }],
};

const projectOptions = computed(() =>
  authStore.projects.map((project) => ({
    label: `${project.name} · ${project.key_bits} bit`,
    value: project.id,
  })),
);

const enabledPlans = computed(() => plans.value.filter((item) => item.enabled));
const planOptions = computed(() =>
  enabledPlans.value.map((item) => ({ label: `${item.name} · ${item.code}`, value: item.id })),
);
const selectedPlan = computed(() => plans.value.find((item) => item.id === formModel.planId));

const { loading, data, columns, columnChecks, getData, getDataByPage, mobilePagination } =
  useAntdvPaginatedTable({
    api: () => {
      if (!authStore.selectedProjectId) {
        return Promise.resolve({ data: [], error: null, response: {} as any });
      }

      return fetchProjectLicenses(authStore.selectedProjectId);
    },
    transform: (response) => {
      const records = !response.error && response.data ? response.data : [];
      const keyword = searchModel.keyword.trim().toLowerCase();

      const filteredRecords = records.filter((item) => {
        const matchKeyword =
          !keyword ||
          item.uuid.toLowerCase().includes(keyword) ||
          item.subject_name.toLowerCase().includes(keyword) ||
          item.subject_email.toLowerCase().includes(keyword) ||
          item.subject_org.toLowerCase().includes(keyword) ||
          item.plan.toLowerCase().includes(keyword) ||
          item.plan_name.toLowerCase().includes(keyword);
        const matchState = !searchModel.state || getLicenseState(item) === searchModel.state;

        return matchKeyword && matchState;
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
      { key: "subject", title: "授权对象", align: "left", width: 220, fixed: "left" },
      { key: "uuid", dataIndex: "uuid", title: "UUID", align: "left", width: 280 },
      { key: "plan", title: "套餐", align: "left", width: 180 },
      { key: "state", title: "状态", align: "center", width: 100 },
      { key: "period", title: "有效期", align: "center", width: 300 },
      { key: "features", title: "功能", align: "left", width: 260 },
      {
        key: "created_at",
        dataIndex: "created_at",
        title: "签发时间",
        align: "center",
        width: 180,
      },
      { key: "operate", title: "操作", align: "center", width: 150, fixed: "right" },
    ],
  });

function stateText(record: Api.Lykn.LicenseRecord) {
  const state = getLicenseState(record);
  if (state === "active") return "生效中";
  if (state === "upcoming") return "未开始";
  if (state === "expired") return "已过期";
  return "未知";
}

function stateColor(record: Api.Lykn.LicenseRecord) {
  const state = getLicenseState(record);
  if (state === "active") return "success";
  if (state === "upcoming") return "processing";
  if (state === "expired") return "error";
  return "default";
}

async function loadPlans() {
  planLoading.value = true;
  try {
    plans.value = unwrapFlatData(await fetchPlans(), "加载套餐列表失败");
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  } finally {
    planLoading.value = false;
  }
}

function handleSearch() {
  getDataByPage();
}

function resetSearch() {
  searchModel.projectId = authStore.selectedProjectId || null;
  searchModel.keyword = "";
  searchModel.state = undefined;
  getDataByPage();
}

function resetForm() {
  formModel.subjectName = "";
  formModel.subjectEmail = "";
  formModel.subjectOrg = "";
  formModel.planId = enabledPlans.value[0]?.id;
  formModel.validRange = null;
  formModel.hardwareKeys = [];
  formModel.hostname = "";
  formModel.cpu_id = "";
  formModel.disk_serial = "";
  formModel.mac_addresses = "";
  void restoreValidation();
}

function openIssueDrawer() {
  if (!authStore.selectedProjectId) {
    window.$message?.warning("请先选择项目");
    return;
  }

  if (!enabledPlans.value.length) {
    window.$message?.warning("请先创建并启用套餐");
    return;
  }

  resetForm();
  drawerOpen.value = true;
}

function hasHardwareKey(key: keyof Api.Lykn.LicenseHardwarePayload) {
  return formModel.hardwareKeys.includes(key);
}

function getHardwarePayload() {
  const payload: Api.Lykn.LicenseHardwarePayload = {};

  if (hasHardwareKey("hostname")) payload.hostname = formModel.hostname.trim();
  if (hasHardwareKey("cpu_id")) payload.cpu_id = formModel.cpu_id.trim();
  if (hasHardwareKey("disk_serial")) payload.disk_serial = formModel.disk_serial.trim();
  if (hasHardwareKey("mac_addresses")) {
    payload.mac_addresses = formModel.mac_addresses
      .split(",")
      .map((item) => item.trim())
      .filter(Boolean);
  }

  return payload;
}

async function submitIssue() {
  await validate();

  if (!authStore.selectedProjectId || !formModel.planId || !formModel.validRange) return;

  issuing.value = true;
  try {
    unwrapFlatData(
      await issueLicense(authStore.selectedProjectId, {
        subject: {
          name: formModel.subjectName.trim(),
          email: formModel.subjectEmail.trim(),
          organization: formModel.subjectOrg.trim(),
        },
        plan_id: formModel.planId,
        not_before: formModel.validRange[0].toISOString(),
        not_after: formModel.validRange[1].toISOString(),
        hardware: getHardwarePayload(),
      }),
      "签发 License 失败",
    );

    window.$message?.success("License 已签发");
    drawerOpen.value = false;
    await getDataByPage();
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  } finally {
    issuing.value = false;
  }
}

async function download(record: Api.Lykn.LicenseRecord) {
  downloadingId.value = record.id;
  try {
    saveBlobFile(await downloadLicense(record.id));
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error, "下载 License 失败"));
  } finally {
    downloadingId.value = undefined;
  }
}

function showDetail(record: Api.Lykn.LicenseRecord) {
  detailRecord.value = record;
  detailOpen.value = true;
}

watch(
  () => authStore.selectedProjectId,
  (projectId) => {
    searchModel.projectId = projectId || null;
    getDataByPage();
  },
  { immediate: true },
);

authStore.loadProjects();
loadPlans();
</script>

<template>
  <div class="min-h-500px flex-col-stretch gap-16px overflow-hidden lt-sm:overflow-auto">
    <QuerySearchForm
      v-model:model="searchModel"
      :items="searchItems"
      :loading="loading"
      @search="handleSearch"
      @reset="resetSearch"
    >
      <template #projectId>
        <ASelect
          v-model:value="authStore.selectedProjectId"
          allow-clear
          class="w-full"
          :loading="authStore.loadingProjects"
          :options="projectOptions"
          placeholder="请选择项目"
        />
      </template>
    </QuerySearchForm>

    <ACard
      title="License 管理"
      variant="borderless"
      :body-style="{ flex: 1, overflow: 'hidden' }"
      class="flex-col-stretch sm:flex-1-hidden card-wrapper"
    >
      <template #extra>
        <TableHeaderOperation v-model:columns="columnChecks" :loading="loading" @refresh="getData">
          <template #default>
            <AButton size="small" ghost type="primary" @click="openIssueDrawer">
              <template #icon>
                <icon-ic-round-plus class="text-icon" />
              </template>
              签发 License
            </AButton>
          </template>
        </TableHeaderOperation>
      </template>

      <ATable
        row-key="id"
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
            <template v-else-if="column.key === 'subject'">
              <div class="flex-col gap-4px">
                <ATypographyText strong>{{ record.subject_name }}</ATypographyText>
                <ATypographyText type="secondary" class="text-12px">{{
                  record.subject_email || "-"
                }}</ATypographyText>
              </div>
            </template>
            <template v-else-if="column.key === 'uuid'">
              <ATypographyText copyable class="font-mono text-12px">{{
                record.uuid
              }}</ATypographyText>
            </template>
            <template v-else-if="column.key === 'plan'">
              {{ record.plan_name || record.plan || "-" }}
            </template>
            <template v-else-if="column.key === 'state'">
              <ATag :color="stateColor(record)">{{ stateText(record) }}</ATag>
            </template>
            <template v-else-if="column.key === 'period'">
              {{ formatPeriod(record.not_before, record.not_after) }}
            </template>
            <template v-else-if="column.key === 'features'">
              <ASpace wrap>
                <ATag v-for="item in record.features" :key="item">{{ item }}</ATag>
                <ATypographyText v-if="!record.features.length" type="secondary"
                  >未配置</ATypographyText
                >
              </ASpace>
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatDateTime(record.created_at) }}
            </template>
            <template v-else-if="column.key === 'operate'">
              <div class="flex-center justify-end gap-6px">
                <AButton type="primary" ghost size="small" @click="showDetail(record)"
                  >详情</AButton
                >
                <AButton
                  size="small"
                  :loading="downloadingId === record.id"
                  @click="download(record)"
                  >下载</AButton
                >
              </div>
            </template>
          </template>
      </ATable>
    </ACard>

    <ADrawer v-model:open="drawerOpen" title="签发 License" :size="640" destroy-on-hidden>
      <template #extra>
        <ASpace>
          <AButton @click="drawerOpen = false">取消</AButton>
          <AButton type="primary" :loading="issuing" @click="submitIssue">签发</AButton>
        </ASpace>
      </template>

      <AForm ref="formRef" layout="vertical" :model="formModel" :rules="formRules">
        <AFormItem label="授权对象" name="subjectName">
          <AInput v-model:value="formModel.subjectName" placeholder="请输入客户或组织名称" />
        </AFormItem>
        <AFormItem label="邮箱" name="subjectEmail">
          <AInput v-model:value="formModel.subjectEmail" placeholder="请输入邮箱" />
        </AFormItem>
        <AFormItem label="组织" name="subjectOrg">
          <AInput v-model:value="formModel.subjectOrg" placeholder="请输入组织" />
        </AFormItem>
        <AFormItem label="套餐" name="planId">
          <ASelect
            v-model:value="formModel.planId"
            :loading="planLoading"
            :options="planOptions"
            placeholder="请选择套餐"
          />
        </AFormItem>
        <AAlert v-if="selectedPlan" type="info" show-icon class="mb-16px">
          <template #message>
            功能 {{ selectedPlan.features.length }} 个，用户
            {{ selectedPlan.max_users || "不限" }}，设备 {{ selectedPlan.max_devices || "-" }}
          </template>
        </AAlert>
        <AFormItem label="有效期" name="validRange">
          <ARangePicker v-model:value="formModel.validRange" class="w-full" show-time />
        </AFormItem>
        <AFormItem label="硬件绑定">
          <ACheckboxGroup v-model:value="formModel.hardwareKeys" :options="hardwareOptions" />
        </AFormItem>
        <AFormItem v-if="hasHardwareKey('hostname')" label="主机名">
          <AInput v-model:value="formModel.hostname" placeholder="例如：office-mac" />
        </AFormItem>
        <AFormItem v-if="hasHardwareKey('cpu_id')" label="CPU ID">
          <AInput v-model:value="formModel.cpu_id" placeholder="例如：BFEBFBFF000906EA" />
        </AFormItem>
        <AFormItem v-if="hasHardwareKey('disk_serial')" label="磁盘序列号">
          <AInput v-model:value="formModel.disk_serial" placeholder="例如：DISK-SERIAL-001" />
        </AFormItem>
        <AFormItem v-if="hasHardwareKey('mac_addresses')" label="MAC 地址">
          <AInput v-model:value="formModel.mac_addresses" placeholder="多个地址用英文逗号分隔" />
        </AFormItem>
      </AForm>
    </ADrawer>

    <ADrawer v-model:open="detailOpen" title="License 详情" :size="560">
      <AEmpty v-if="!detailRecord" description="暂无详情" />
      <ADescriptions v-else :column="1" bordered size="small">
        <ADescriptionsItem label="UUID">
          <ATypographyText copyable class="font-mono text-12px">{{
            detailRecord.uuid
          }}</ATypographyText>
        </ADescriptionsItem>
        <ADescriptionsItem label="授权对象">{{ detailRecord.subject_name }}</ADescriptionsItem>
        <ADescriptionsItem label="邮箱">{{ detailRecord.subject_email || "-" }}</ADescriptionsItem>
        <ADescriptionsItem label="组织">{{ detailRecord.subject_org || "-" }}</ADescriptionsItem>
        <ADescriptionsItem label="套餐">{{
          detailRecord.plan_name || detailRecord.plan || "-"
        }}</ADescriptionsItem>
        <ADescriptionsItem label="有效期">{{
          formatPeriod(detailRecord.not_before, detailRecord.not_after)
        }}</ADescriptionsItem>
        <ADescriptionsItem label="额度">
          用户 {{ detailRecord.limits?.max_users || "不限" }} / 设备
          {{ detailRecord.limits?.max_devices || "-" }}
        </ADescriptionsItem>
        <ADescriptionsItem label="功能">
          <ASpace wrap>
            <ATag v-for="item in detailRecord.features" :key="item">{{ item }}</ATag>
            <ATypographyText v-if="!detailRecord.features.length" type="secondary"
              >未配置</ATypographyText
            >
          </ASpace>
        </ADescriptionsItem>
      </ADescriptions>
    </ADrawer>
  </div>
</template>
