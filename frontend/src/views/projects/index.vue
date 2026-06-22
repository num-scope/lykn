<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useBoolean } from "@sa/hooks";
import { useAntdvForm } from "@/hooks/common/form";
import { useAntdvPaginatedTable, useTableOperate } from "@/hooks/common/table";
import {
  createProject,
  deleteProject,
  downloadProjectPublicKey,
  fetchProject,
  fetchProjects,
  updateProject,
} from "@/service/api";
import { useAuthStore } from "@/store/modules/auth";
import { formatDateTime, getServiceErrorMessage, saveBlobFile, unwrapFlatData } from "@/utils/lykn";

const router = useRouter();
const authStore = useAuthStore();

const searchModel = reactive({
  keyword: "",
  key_bits: undefined as number | undefined,
});

const { formRef, validate, restoreValidation } = useAntdvForm();
const paginationParams = reactive({
  current: 1,
  size: 10,
});
const tableScrollX = 1304;

const modalOpen = ref(false);
const submitting = ref(false);
const operateType = ref<AntdvUI.TableOperateType>("add");
const editingId = ref<number>();
const publicKeyOpen = ref(false);
const publicKeyLoading = ref(false);
const publicKeyProject = ref<Api.Lykn.ProjectRecord>();
const publicKey = ref("");
const downloadingId = ref<number>();

const formModel = reactive<Api.Lykn.CreateProjectPayload>({
  name: "",
  description: "",
  key_bits: 2048,
});

const keyBitOptions = [
  { label: "2048 bit", value: 2048 },
  { label: "3072 bit", value: 3072 },
  { label: "4096 bit", value: 4096 },
];

const searchItems = computed(() => [
  { key: "keyword", label: "关键词", placeholder: "项目名称 / 描述" },
  { key: "key_bits", label: "密钥位数", type: "select" as const, placeholder: "请选择密钥位数" },
]);

const formRules = {
  name: [{ required: true, message: "请输入项目名称" }],
  key_bits: [{ required: true, message: "请选择密钥位数" }],
};

const { loading, data, columns, columnChecks, getData, getDataByPage, mobilePagination } =
  useAntdvPaginatedTable({
    api: fetchProjects,
    transform: (response) => {
      const records = unwrapFlatData(response, "加载项目列表失败") as Api.Lykn.ProjectRecord[];
      const keyword = searchModel.keyword.trim().toLowerCase();

      const filteredRecords = records.filter((item) => {
        const matchKeyword =
          !keyword ||
          item.name.toLowerCase().includes(keyword) ||
          item.description.toLowerCase().includes(keyword);
        const matchKeyBits = !searchModel.key_bits || item.key_bits === searchModel.key_bits;

        return matchKeyword && matchKeyBits;
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
      { key: "name", dataIndex: "name", title: "项目", align: "left", width: 220, fixed: "left" },
      { key: "description", dataIndex: "description", title: "描述", align: "left", width: 280 },
      { key: "key_bits", dataIndex: "key_bits", title: "密钥", align: "center", width: 120 },
      { key: "updated_at", dataIndex: "updated_at", title: "更新时间", align: "center", width: 180 },
      { key: "created_at", dataIndex: "created_at", title: "创建时间", align: "center", width: 180 },
      { key: "operate", title: "操作", align: "center", width: 260, fixed: "right" },
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

const modalTitle = computed(() => (operateType.value === "add" ? "新建项目" : "编辑项目"));

function resetForm() {
  editingId.value = undefined;
  formModel.name = "";
  formModel.description = "";
  formModel.key_bits = 2048;
  void restoreValidation();
}

function handleSearch() {
  getDataByPage();
}

function resetSearch() {
  searchModel.keyword = "";
  searchModel.key_bits = undefined;
  getDataByPage();
}

function openCreateModal() {
  resetForm();
  operateType.value = "add";
  modalOpen.value = true;
}

function openEditModal(record: Api.Lykn.ProjectRecord) {
  operateType.value = "edit";
  editingId.value = record.id;
  formModel.name = record.name;
  formModel.description = record.description;
  formModel.key_bits = record.key_bits;
  modalOpen.value = true;
}

async function submitForm() {
  await validate();

  submitting.value = true;
  try {
    if (operateType.value === "add") {
      unwrapFlatData(await createProject({ ...formModel }), "创建项目失败");
      window.$message?.success("项目已创建");
    } else if (editingId.value) {
      unwrapFlatData(
        await updateProject(editingId.value, {
          name: formModel.name,
          description: formModel.description,
        }),
        "更新项目失败",
      );
      window.$message?.success("项目已更新");
    }

    modalOpen.value = false;
    await getDataByPage();
    await authStore.loadProjects();
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  } finally {
    submitting.value = false;
  }
}

async function handleBatchDelete() {
  setBatchDeleting(true);
  try {
    await Promise.all(checkedRowKeys.value.map((id) => deleteProject(Number(id))));
    await onBatchDeleted();
    await authStore.loadProjects();
  } catch (error) {
    checkedRowKeys.value = [];
    window.$message?.error(getServiceErrorMessage(error, "批量删除项目失败"));
    await getDataByPage();
  } finally {
    setBatchDeleting(false);
  }
}

async function removeProject(record: Api.Lykn.ProjectRecord) {
  try {
    unwrapFlatData(await deleteProject(record.id), "删除项目失败");
    window.$message?.success("项目已删除");
    await getDataByPage();
    await authStore.loadProjects();
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  }
}

async function showPublicKey(record: Api.Lykn.ProjectRecord) {
  publicKeyOpen.value = true;
  publicKeyProject.value = record;
  publicKey.value = "";
  publicKeyLoading.value = true;

  try {
    const detail = unwrapFlatData(await fetchProject(record.id), "获取公钥失败");
    publicKeyProject.value = detail;
    publicKey.value = detail.public_key || "";
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error));
  } finally {
    publicKeyLoading.value = false;
  }
}

async function copyPublicKey() {
  if (!publicKey.value) return;

  await navigator.clipboard.writeText(publicKey.value);
  window.$message?.success("公钥已复制");
}

async function downloadPublicKey(record: Api.Lykn.ProjectRecord) {
  downloadingId.value = record.id;
  try {
    saveBlobFile(await downloadProjectPublicKey(record.id));
  } catch (error) {
    window.$message?.error(getServiceErrorMessage(error, "下载公钥失败"));
  } finally {
    downloadingId.value = undefined;
  }
}

function toProject(record: unknown) {
  return record as Api.Lykn.ProjectRecord;
}

function goLicenses(record: Api.Lykn.ProjectRecord) {
  authStore.selectProject(record.id);
  router.push({ name: "licenses" });
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
    >
      <template #key_bits="{ model }">
        <ASelect
          v-model:value="model.key_bits"
          allow-clear
          class="w-full"
          :options="keyBitOptions"
          placeholder="请选择密钥位数"
        />
      </template>
    </QuerySearchForm>

    <ACard
      title="项目管理"
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
              新建项目
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
            <template v-else-if="column.key === 'name'">
              <ATypographyText strong>{{ toProject(record).name }}</ATypographyText>
            </template>
            <template v-else-if="column.key === 'description'">
              <ATypographyText type="secondary">{{
                toProject(record).description || "-"
              }}</ATypographyText>
            </template>
            <template v-else-if="column.key === 'key_bits'">
              <ATag color="blue">{{ toProject(record).key_bits }} bit</ATag>
            </template>
            <template v-else-if="column.key === 'updated_at'">
              {{ formatDateTime(toProject(record).updated_at) }}
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatDateTime(toProject(record).created_at) }}
            </template>
            <template v-else-if="column.key === 'operate'">
              <div class="flex-center justify-end gap-6px">
                <AButton type="primary" ghost size="small" @click="goLicenses(toProject(record))"
                  >License</AButton
                >
                <AButton size="small" @click="showPublicKey(toProject(record))">公钥</AButton>
                <AButton
                  size="small"
                  :loading="downloadingId === toProject(record).id"
                  @click="downloadPublicKey(toProject(record))"
                  >下载</AButton
                >
                <AButton size="small" @click="openEditModal(toProject(record))">编辑</AButton>
                <APopconfirm
                  title="确定删除该项目？相关 License 会一并删除"
                  @confirm="removeProject(toProject(record))"
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
        <AFormItem label="项目名称" name="name">
          <AInput v-model:value="formModel.name" placeholder="请输入项目名称" />
        </AFormItem>
        <AFormItem label="项目描述" name="description">
          <ATextarea
            v-model:value="formModel.description"
            :auto-size="{ minRows: 3, maxRows: 5 }"
            placeholder="请输入项目描述"
          />
        </AFormItem>
        <AFormItem label="RSA 密钥位数" name="key_bits">
          <ASelect
            v-model:value="formModel.key_bits"
            :disabled="operateType === 'edit'"
            :options="keyBitOptions"
          />
        </AFormItem>
      </AForm>
    </AModal>

    <AModal v-model:open="publicKeyOpen" width="760px" :footer="null" destroy-on-hidden>
      <template #title>{{ publicKeyProject?.name || "项目" }} 公钥</template>
      <ASpin :spinning="publicKeyLoading">
        <ASpace vertical class="w-full">
          <AAlert type="info" show-icon message="客户端 SDK 或 CLI 可使用该公钥验证离线 License" />
          <pre class="max-h-420px overflow-auto rd-6px bg-layout p-16px text-12px leading-6">{{
            publicKey || "暂无公钥"
          }}</pre>
          <ASpace>
            <AButton type="primary" :disabled="!publicKey" @click="copyPublicKey">复制公钥</AButton>
            <AButton
              v-if="publicKeyProject"
              :loading="downloadingId === publicKeyProject.id"
              @click="downloadPublicKey(publicKeyProject)"
            >
              下载 public.pem
            </AButton>
          </ASpace>
        </ASpace>
      </ASpin>
    </AModal>
  </div>
</template>
