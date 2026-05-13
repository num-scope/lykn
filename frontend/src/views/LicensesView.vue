<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import { message, type DescriptionsItemType, type RangePickerProps } from "antdv-next";
import { storeToRefs } from "pinia";

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
import {
  api,
  type IssueLicensePayload,
  type LicenseHardwarePayload,
  type LicenseRecord,
  type PlanRecord,
} from "../api/api";
import { useAuthStore } from "../stores/auth";
import {
  downloadFile,
  formatDateTime,
  formatPeriod,
  getErrorMessage,
  getLicenseState,
} from "../utils/format";

const authStore = useAuthStore();
const { projects, selectedProjectId, loadingProjects } = storeToRefs(authStore);

interface LicenseSearchState {
  keyword: string;
  state?: "active" | "expired" | "upcoming";
}

type LicenseRangeValue = NonNullable<RangePickerProps["value"]>;
type LicenseDateValue = NonNullable<LicenseRangeValue[number]>;
type HardwareBindingKey = "hostname" | "cpu_id" | "disk_serial" | "mac_addresses";

interface LicenseFormState {
  subjectName: string;
  subjectEmail: string;
  subjectOrg: string;
  planId?: number;
  validRange?: LicenseRangeValue;
  hardwareKeys: HardwareBindingKey[];
  hostname: string;
  cpuId: string;
  diskSerial: string;
  macAddress: string;
}

const licenses = ref<LicenseRecord[]>([]);
const plans = ref<PlanRecord[]>([]);
const loadingLicenses = ref(false);
const loadingPlans = ref(false);
const issuing = ref(false);
const issueDrawerOpen = ref(false);
const detailOpen = ref(false);
const detailLicense = ref<LicenseRecord | null>(null);
const downloadingId = ref<number | null>(null);
const searchState = reactive<LicenseSearchState>({
  keyword: "",
  state: undefined,
});
const pagination = reactive<AdminPaginationState>({
  current: 1,
  pageSize: 10,
  total: 0,
});
const licenseForm = reactive<LicenseFormState>({
  subjectName: "",
  subjectEmail: "",
  subjectOrg: "",
  planId: undefined,
  validRange: undefined,
  hardwareKeys: [],
  hostname: "",
  cpuId: "",
  diskSerial: "",
  macAddress: "",
});

const searchFields: AdminSearchField[] = [
  {
    key: "keyword",
    label: "关键词",
    type: "input",
    placeholder: "客户 / 邮箱 / UUID / plan",
  },
  {
    key: "state",
    label: "有效性",
    type: "select",
    placeholder: "全部状态",
    options: [
      { label: "生效中", value: "active" },
      { label: "未开始", value: "upcoming" },
      { label: "已过期", value: "expired" },
    ],
  },
];

const hardwareBindingOptions = [
  { label: "主机名", value: "hostname" },
  { label: "CPU ID", value: "cpu_id" },
  { label: "磁盘序列号", value: "disk_serial" },
  { label: "MAC 地址", value: "mac_addresses" },
];

const licenseColumns = [
  { title: "授权对象", key: "subject", dataIndex: "subject_name" },
  { title: "License UUID", key: "uuid", dataIndex: "uuid" },
  { title: "计划", key: "plan", dataIndex: "plan" },
  { title: "状态", key: "state" },
  { title: "有效期", key: "period" },
  { title: "功能", key: "features" },
  { title: "签发时间", key: "created_at", dataIndex: "created_at" },
  { title: "操作", key: "actions" },
];

const projectOptions = computed(() =>
  projects.value.map((project) => ({
    label: `${project.name} · ${project.key_bits} bit`,
    value: project.id,
  })),
);

const enabledPlans = computed(() => plans.value.filter((plan) => plan.enabled));
const planOptions = computed(() =>
  enabledPlans.value.map((plan) => ({
    label: `${plan.name} · ${plan.code}`,
    value: plan.id,
  })),
);
const selectedPlan = computed(() => plans.value.find((plan) => plan.id === licenseForm.planId));

const detailItems = computed<DescriptionsItemType[]>(() => {
  const license = detailLicense.value;
  if (!license) return [];
  return [
    { key: "uuid", label: "UUID", content: license.uuid },
    { key: "subject", label: "授权对象", content: license.subject_name },
    { key: "email", label: "邮箱", content: license.subject_email || "-" },
    { key: "org", label: "组织", content: license.subject_org || "-" },
    {
      key: "plan",
      label: "套餐",
      content: license.plan_name
        ? `${license.plan_name} · ${license.plan}`
        : license.plan || "default",
    },
    {
      key: "period",
      label: "有效期",
      content: formatPeriod(license.not_before, license.not_after),
    },
    { key: "features", label: "功能", content: license.features.join(", ") || "-" },
    {
      key: "limits",
      label: "额度",
      content: `用户 ${license.limits?.max_users || "不限"} / 设备 ${
        license.limits?.max_devices || "-"
      }`,
    },
  ];
});

const filteredLicenses = computed(() => {
  const keyword = searchState.keyword.trim().toLowerCase();
  return licenses.value.filter((license) => {
    const matchesKeyword =
      !keyword ||
      license.uuid.toLowerCase().includes(keyword) ||
      license.subject_name.toLowerCase().includes(keyword) ||
      license.subject_email.toLowerCase().includes(keyword) ||
      license.subject_org.toLowerCase().includes(keyword) ||
      (license.plan_name || "").toLowerCase().includes(keyword) ||
      license.plan.toLowerCase().includes(keyword);
    const matchesState = !searchState.state || getLicenseState(license) === searchState.state;
    return matchesKeyword && matchesState;
  });
});

const pagedLicenses = computed(() => {
  const start = (pagination.current - 1) * pagination.pageSize;
  return filteredLicenses.value.slice(start, start + pagination.pageSize);
});

watch(
  selectedProjectId,
  (projectId) => {
    if (projectId) {
      loadLicenses(projectId);
    } else {
      licenses.value = [];
    }
  },
  { immediate: true },
);

watch(
  filteredLicenses,
  (items) => {
    pagination.total = items.length;
    const maxPage = Math.max(1, Math.ceil(items.length / pagination.pageSize));
    if (pagination.current > maxPage) pagination.current = maxPage;
  },
  { immediate: true },
);

watch(
  () => [searchState.keyword, searchState.state],
  () => {
    pagination.current = 1;
  },
);

async function loadLicenses(projectId = selectedProjectId.value) {
  if (!projectId) return;
  loadingLicenses.value = true;
  try {
    licenses.value = await api.listProjectLicenses(projectId);
  } catch (error) {
    message.error(getErrorMessage(error, "加载授权列表失败"));
  } finally {
    loadingLicenses.value = false;
  }
}

const loadPlans = async () => {
  loadingPlans.value = true;
  try {
    plans.value = await api.listPlans();
    if (!licenseForm.planId && enabledPlans.value.length) {
      licenseForm.planId = enabledPlans.value[0].id;
    }
  } catch (error) {
    message.error(getErrorMessage(error, "加载套餐列表失败"));
  } finally {
    loadingPlans.value = false;
  }
};

const onSearchReset = () => {
  searchState.keyword = "";
  searchState.state = undefined;
  pagination.current = 1;
};

const onPageChange = (page: number, pageSize: number) => {
  pagination.current = page;
  pagination.pageSize = pageSize;
};

const resetLicenseForm = () => {
  licenseForm.subjectName = "";
  licenseForm.subjectEmail = "";
  licenseForm.subjectOrg = "";
  licenseForm.planId = enabledPlans.value[0]?.id;
  licenseForm.validRange = undefined;
  licenseForm.hardwareKeys = [];
  licenseForm.hostname = "";
  licenseForm.cpuId = "";
  licenseForm.diskSerial = "";
  licenseForm.macAddress = "";
};

const openIssueDrawer = () => {
  if (!selectedProjectId.value) {
    message.warning("请先选择项目");
    return;
  }
  if (!enabledPlans.value.length) {
    message.warning("请先创建并启用套餐");
    return;
  }
  resetLicenseForm();
  issueDrawerOpen.value = true;
};

const toIsoString = (value: LicenseDateValue | undefined) => {
  if (!value) return "";
  return value.toISOString();
};

const hasHardwareKey = (key: HardwareBindingKey) => licenseForm.hardwareKeys.includes(key);

const buildHardwarePayload = (): LicenseHardwarePayload | null => {
  const payload: LicenseHardwarePayload = {};
  if (hasHardwareKey("hostname")) {
    const hostname = licenseForm.hostname.trim();
    if (!hostname) {
      message.warning("请输入主机名");
      return null;
    }
    payload.hostname = hostname;
  }
  if (hasHardwareKey("cpu_id")) {
    const cpuId = licenseForm.cpuId.trim();
    if (!cpuId) {
      message.warning("请输入 CPU ID");
      return null;
    }
    payload.cpu_id = cpuId;
  }
  if (hasHardwareKey("disk_serial")) {
    const diskSerial = licenseForm.diskSerial.trim();
    if (!diskSerial) {
      message.warning("请输入磁盘序列号");
      return null;
    }
    payload.disk_serial = diskSerial;
  }
  if (hasHardwareKey("mac_addresses")) {
    const macAddress = licenseForm.macAddress.trim();
    if (!macAddress) {
      message.warning("请输入 MAC 地址");
      return null;
    }
    payload.mac_addresses = [macAddress];
  }
  return payload;
};

const submitLicense = async () => {
  const projectId = selectedProjectId.value;
  if (!projectId) return;
  if (!licenseForm.subjectName.trim()) {
    message.warning("请输入授权对象名称");
    return;
  }
  if (!licenseForm.planId) {
    message.warning("请选择套餐");
    return;
  }
  if (!licenseForm.validRange?.[0] || !licenseForm.validRange?.[1]) {
    message.warning("请选择有效期");
    return;
  }

  const hardware = buildHardwarePayload();
  if (!hardware) {
    return;
  }

  const payload: IssueLicensePayload = {
    subject: {
      name: licenseForm.subjectName.trim(),
      email: licenseForm.subjectEmail.trim(),
      organization: licenseForm.subjectOrg.trim(),
    },
    plan_id: licenseForm.planId,
    not_before: toIsoString(licenseForm.validRange[0]),
    not_after: toIsoString(licenseForm.validRange[1]),
    hardware,
  };

  if (!payload.not_before || !payload.not_after) {
    message.error("有效期格式不正确");
    return;
  }

  issuing.value = true;
  try {
    await api.issueLicense(projectId, payload);
    message.success("License 已签发");
    issueDrawerOpen.value = false;
    await loadLicenses(projectId);
  } catch (error) {
    message.error(getErrorMessage(error, "签发失败"));
  } finally {
    issuing.value = false;
  }
};

const downloadLicense = async (license: LicenseRecord) => {
  downloadingId.value = license.id;
  try {
    const file = await api.downloadLicense(license.id);
    downloadFile(file);
    message.success("License 下载已开始");
  } catch (error) {
    message.error(getErrorMessage(error, "下载 license 失败"));
  } finally {
    downloadingId.value = null;
  }
};

const openDetail = async (license: LicenseRecord) => {
  detailOpen.value = true;
  detailLicense.value = license;
  try {
    detailLicense.value = await api.getLicense(license.id);
  } catch (error) {
    message.error(getErrorMessage(error, "加载详情失败"));
  }
};

const licenseActions = (license: LicenseRecord): AdminRowAction<LicenseRecord>[] => [
  {
    key: "detail",
    label: "详情",
    onClick: () => openDetail(license),
  },
  {
    key: "download",
    label: downloadingId.value === license.id ? "下载中" : "下载",
    disabled: downloadingId.value === license.id,
    onClick: () => downloadLicense(license),
  },
];

const stateText = (license: LicenseRecord) => {
  const state = getLicenseState(license);
  if (state === "active") return "生效中";
  if (state === "upcoming") return "未开始";
  if (state === "expired") return "已过期";
  return "未知";
};

const stateColor = (license: LicenseRecord) => {
  const state = getLicenseState(license);
  if (state === "active") return "success";
  if (state === "upcoming") return "processing";
  if (state === "expired") return "error";
  return "default";
};

loadPlans();
</script>

<template>
  <div class="flex flex-col gap-4">
    <AdminResourcePage title="License 管理">
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
        <a-select
          :value="selectedProjectId"
          :options="projectOptions"
          :loading="loadingProjects"
          placeholder="选择项目"
          @change="(value: number) => authStore.selectProject(value)"
        />
        <a-button type="primary" :disabled="!selectedProjectId" @click="openIssueDrawer">
          签发 License
        </a-button>
        <a-button @click="loadLicenses()">刷新</a-button>
      </template>

      <a-alert
        v-if="!selectedProjectId"
        class="mb-4"
        type="warning"
        show-icon
        message="请选择项目后查看或签发 license。"
      />

      <AdminDataTable
        :columns="licenseColumns"
        :data-source="pagedLicenses"
        :loading="loadingLicenses"
        :pagination="pagination"
        row-key="id"
        @change="onPageChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'subject'">
            <a-space size="small" class="whitespace-nowrap">
              <a-typography-text strong>{{ record.subject_name }}</a-typography-text>
              <a-typography-text type="secondary">
                / {{ record.subject_email || record.subject_org || "未填写邮箱 / 组织" }}
              </a-typography-text>
            </a-space>
          </template>
          <template v-else-if="column.key === 'uuid'">
            <a-typography-text copyable class="font-mono text-xs">{{
              record.uuid
            }}</a-typography-text>
          </template>
          <template v-else-if="column.key === 'plan'">
            <a-space size="small" class="whitespace-nowrap">
              <a-tag color="blue">{{ record.plan_name || record.plan || "default" }}</a-tag>
              <a-typography-text
                v-if="record.plan_name && record.plan"
                type="secondary"
                class="text-xs"
              >
                · {{ record.plan }}
              </a-typography-text>
            </a-space>
          </template>
          <template v-else-if="column.key === 'state'">
            <AdminStatusTag
              v-if="getLicenseState(record) === 'active' || getLicenseState(record) === 'expired'"
              :value="getLicenseState(record) === 'active'"
              :active-value="true"
              active-text="生效中"
              inactive-text="已过期"
              inactive-color="error"
            />
            <a-tag v-else :color="stateColor(record)">{{ stateText(record) }}</a-tag>
          </template>
          <template v-else-if="column.key === 'period'">
            {{ formatPeriod(record.not_before, record.not_after) }}
          </template>
          <template v-else-if="column.key === 'features'">
            <a-space class="whitespace-nowrap">
              <a-tag v-for="feature in record.features" :key="feature">{{ feature }}</a-tag>
              <a-typography-text v-if="!record.features.length" type="secondary"
                >无</a-typography-text
              >
            </a-space>
          </template>
          <template v-else-if="column.key === 'created_at'">
            {{ formatDateTime(record.created_at) }}
          </template>
          <template v-else-if="column.key === 'actions'">
            <AdminRowActions :record="record" :actions="licenseActions(record)" />
          </template>
        </template>
      </AdminDataTable>
    </AdminResourcePage>

    <a-drawer
      v-model:open="issueDrawerOpen"
      :size="760"
      title="签发 License"
      :mask-closable="!issuing"
      :keyboard="!issuing"
    >
      <template #extra>
        <a-space>
          <a-button :disabled="issuing" @click="issueDrawerOpen = false">取消</a-button>
          <a-button type="primary" :loading="issuing" @click="submitLicense">签发</a-button>
        </a-space>
      </template>
      <a-form layout="vertical" :model="licenseForm">
        <a-row :gutter="16">
          <a-col :xs="24" :md="12">
            <a-form-item label="授权对象" required>
              <a-input v-model:value="licenseForm.subjectName" placeholder="客户或设备名称" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :md="12">
            <a-form-item label="邮箱">
              <a-input v-model:value="licenseForm.subjectEmail" placeholder="license@example.com" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :xs="24" :md="12">
            <a-form-item label="组织">
              <a-input v-model:value="licenseForm.subjectOrg" placeholder="公司或团队名称" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :md="12">
            <a-form-item label="套餐" required>
              <a-select
                v-model:value="licenseForm.planId"
                :options="planOptions"
                :loading="loadingPlans"
                show-search
                placeholder="选择套餐"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="有效期" required>
          <a-range-picker v-model:value="licenseForm.validRange" class="w-full" show-time />
        </a-form-item>
        <a-form-item label="套餐功能预览">
          <a-space wrap>
            <a-tag v-for="feature in selectedPlan?.features || []" :key="feature.id">
              {{ feature.name }}
            </a-tag>
            <a-typography-text v-if="!selectedPlan?.features.length" type="secondary">
              暂无功能
            </a-typography-text>
          </a-space>
        </a-form-item>
        <a-form-item label="套餐额度预览">
          <a-space wrap>
            <a-tag color="blue">用户 {{ selectedPlan?.max_users || "不限" }}</a-tag>
            <a-tag color="green">设备 {{ selectedPlan?.max_devices || "-" }}</a-tag>
          </a-space>
        </a-form-item>
        <a-form-item label="硬件绑定条件">
          <a-checkbox-group
            v-model:value="licenseForm.hardwareKeys"
            :options="hardwareBindingOptions"
          />
        </a-form-item>
        <a-row v-if="licenseForm.hardwareKeys.length" :gutter="16">
          <a-col v-if="hasHardwareKey('hostname')" :xs="24" :md="12">
            <a-form-item label="主机名" required>
              <a-input v-model:value="licenseForm.hostname" placeholder="例如：office-mac" />
            </a-form-item>
          </a-col>
          <a-col v-if="hasHardwareKey('cpu_id')" :xs="24" :md="12">
            <a-form-item label="CPU ID" required>
              <a-input v-model:value="licenseForm.cpuId" placeholder="例如：BFEBFBFF000906EA" />
            </a-form-item>
          </a-col>
          <a-col v-if="hasHardwareKey('disk_serial')" :xs="24" :md="12">
            <a-form-item label="磁盘序列号" required>
              <a-input v-model:value="licenseForm.diskSerial" placeholder="例如：DISK-SERIAL-001" />
            </a-form-item>
          </a-col>
          <a-col v-if="hasHardwareKey('mac_addresses')" :xs="24" :md="12">
            <a-form-item label="绑定 MAC 地址" required>
              <a-input
                v-model:value="licenseForm.macAddress"
                placeholder="例如：AA:BB:CC:DD:EE:FF"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-alert
          v-else
          type="info"
          show-icon
          message="未选择硬件绑定条件时，License 不校验硬件。"
        />
      </a-form>
    </a-drawer>

    <a-drawer v-model:open="detailOpen" :size="560" title="License 详情">
      <a-empty v-if="!detailLicense" description="暂无详情" />
      <a-descriptions v-else :column="1" bordered size="small" :items="detailItems">
        <template #contentRender="{ item }">
          <a-typography-text v-if="item.key === 'uuid'" copyable class="font-mono text-xs">
            {{ item.content }}
          </a-typography-text>
          <a-space v-else-if="item.key === 'features'" wrap>
            <a-tag v-for="feature in detailLicense?.features || []" :key="feature">
              {{ feature }}
            </a-tag>
            <span v-if="!detailLicense?.features.length">-</span>
          </a-space>
          <template v-else>{{ item.content }}</template>
        </template>
      </a-descriptions>
    </a-drawer>
  </div>
</template>
