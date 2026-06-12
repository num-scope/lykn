<script setup lang="ts">
import { computed, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { message } from "antdv-next";
import { storeToRefs } from "pinia";
import { PlusOutlined } from "@antdv-next/icons";

import {
  AdminDataTable,
  AdminResourcePage,
  AdminRowActions,
  AdminSearchForm,
  type AdminPaginationState,
  type AdminRowAction,
  type AdminSearchField,
} from "../components/admin-kit";
import { api, type CreateProjectPayload, type ProjectRecord } from "../api/api";
import { useAuthStore } from "../stores/auth";
import { downloadFile, formatDateTime, getErrorMessage } from "../utils/format";

const router = useRouter();
const authStore = useAuthStore();
const { projects, loadingProjects } = storeToRefs(authStore);

interface ProjectSearchState {
  keyword: string;
  keyBits?: number;
}

interface ProjectFormState {
  id?: number;
  name: string;
  description: string;
  key_bits: number;
}

const searchState = reactive<ProjectSearchState>({
  keyword: "",
  keyBits: undefined,
});
const pagination = reactive<AdminPaginationState>({
  current: 1,
  pageSize: 10,
  total: 0,
});
const modalOpen = ref(false);
const modalMode = ref<"create" | "edit">("create");
const saving = ref(false);
const publicKeyOpen = ref(false);
const publicKeyLoading = ref(false);
const publicKeyProject = ref<ProjectRecord | null>(null);
const publicKey = ref("");
const downloadingId = ref<number | null>(null);

const projectForm = reactive<ProjectFormState>({
  name: "",
  description: "",
  key_bits: 2048,
});

const searchFields: AdminSearchField[] = [
  {
    key: "keyword",
    label: "关键词",
    type: "input",
    placeholder: "项目名称 / 描述",
  },
  {
    key: "keyBits",
    label: "密钥位数",
    type: "select",
    placeholder: "全部位数",
    options: [
      { label: "2048 bit", value: 2048 },
      { label: "3072 bit", value: 3072 },
      { label: "4096 bit", value: 4096 },
    ],
  },
];

const keyBitOptions = [
  { label: "2048 bit（兼容优先）", value: 2048 },
  { label: "3072 bit（平衡安全与性能）", value: 3072 },
  { label: "4096 bit（高安全）", value: 4096 },
];

const projectColumns = [
  { title: "项目", key: "project", dataIndex: "name" },
  { title: "描述", key: "description", dataIndex: "description" },
  { title: "密钥", key: "key_bits", dataIndex: "key_bits" },
  { title: "更新时间", key: "updated_at", dataIndex: "updated_at" },
  { title: "创建时间", key: "created_at", dataIndex: "created_at" },
  { title: "操作", key: "actions" },
];

const filteredProjects = computed(() => {
  const keyword = searchState.keyword.trim().toLowerCase();
  return projects.value.filter((project) => {
    const matchesKeyword =
      !keyword ||
      project.name.toLowerCase().includes(keyword) ||
      project.description.toLowerCase().includes(keyword);
    const matchesKeyBits = !searchState.keyBits || project.key_bits === searchState.keyBits;
    return matchesKeyword && matchesKeyBits;
  });
});

const pagedProjects = computed(() => {
  const start = (pagination.current - 1) * pagination.pageSize;
  return filteredProjects.value.slice(start, start + pagination.pageSize);
});

const modalTitle = computed(() => (modalMode.value === "create" ? "新建项目" : "编辑项目"));

watch(
  filteredProjects,
  (projects) => {
    pagination.total = projects.length;
    const maxPage = Math.max(1, Math.ceil(projects.length / pagination.pageSize));
    if (pagination.current > maxPage) pagination.current = maxPage;
  },
  { immediate: true },
);

watch(
  () => [searchState.keyword, searchState.keyBits],
  () => {
    pagination.current = 1;
  },
);

const onSearchReset = () => {
  searchState.keyword = "";
  searchState.keyBits = undefined;
  pagination.current = 1;
};

const onPageChange = (page: number, pageSize: number) => {
  pagination.current = page;
  pagination.pageSize = pageSize;
};

const resetProjectForm = () => {
  projectForm.id = undefined;
  projectForm.name = "";
  projectForm.description = "";
  projectForm.key_bits = 2048;
};

const openCreateModal = () => {
  resetProjectForm();
  modalMode.value = "create";
  modalOpen.value = true;
};

const openEditModal = (project: ProjectRecord) => {
  projectForm.id = project.id;
  projectForm.name = project.name;
  projectForm.description = project.description;
  projectForm.key_bits = project.key_bits;
  modalMode.value = "edit";
  modalOpen.value = true;
};

const submitProject = async () => {
  const name = projectForm.name.trim();
  if (!name) {
    message.warning("请输入项目名称");
    return;
  }
  saving.value = true;
  try {
    if (modalMode.value === "create") {
      const payload: CreateProjectPayload = {
        name,
        description: projectForm.description.trim(),
        key_bits: projectForm.key_bits,
      };
      await api.createProject(payload);
      message.success("项目已创建");
    } else if (projectForm.id) {
      await api.updateProject(projectForm.id, {
        name,
        description: projectForm.description.trim(),
      });
      message.success("项目已更新");
    }
    modalOpen.value = false;
    await authStore.loadProjects();
  } catch (error) {
    message.error(getErrorMessage(error));
  } finally {
    saving.value = false;
  }
};

const deleteProject = async (project: ProjectRecord) => {
  try {
    await api.deleteProject(project.id);
    message.success(`已删除项目「${project.name}」`);
    await authStore.loadProjects();
  } catch (error) {
    message.error(getErrorMessage(error));
  }
};

const showPublicKey = async (project: ProjectRecord) => {
  publicKeyOpen.value = true;
  publicKeyLoading.value = true;
  publicKeyProject.value = project;
  publicKey.value = "";
  try {
    const detail = await api.getProject(project.id);
    publicKeyProject.value = detail;
    publicKey.value = detail.public_key || "";
  } catch (error) {
    message.error(getErrorMessage(error, "获取公钥失败"));
  } finally {
    publicKeyLoading.value = false;
  }
};

const copyPublicKey = async () => {
  if (!publicKey.value) return;
  try {
    await navigator.clipboard.writeText(publicKey.value);
    message.success("公钥已复制");
  } catch {
    message.warning("复制失败，请手动复制");
  }
};

const downloadPublicKey = async (project: ProjectRecord) => {
  downloadingId.value = project.id;
  try {
    const file = await api.downloadProjectPublicKey(project.id);
    downloadFile({ ...file, filename: `${project.name || "project"}-public.pem` });
    message.success("公钥下载已开始");
  } catch (error) {
    message.error(getErrorMessage(error, "下载公钥失败"));
  } finally {
    downloadingId.value = null;
  }
};

const projectActions = (project: ProjectRecord): AdminRowAction<ProjectRecord>[] => [
  {
    key: "licenses",
    label: "授权",
    onClick: () => {
      authStore.selectProject(project.id);
      router.push({ name: "licenses" });
    },
  },
  {
    key: "public-key",
    label: "公钥",
    onClick: () => showPublicKey(project),
  },
  {
    key: "download-key",
    label: downloadingId.value === project.id ? "下载中" : "下载",
    disabled: downloadingId.value === project.id,
    onClick: () => downloadPublicKey(project),
  },
  {
    key: "edit",
    label: "编辑",
    onClick: () => openEditModal(project),
  },
  {
    key: "delete",
    label: "删除",
    danger: true,
    confirm: `确定删除项目「${project.name}」？相关 license 会一并删除。`,
    confirmOkText: "删除",
    onClick: () => deleteProject(project),
  },
];
</script>

<template>
  <div class="flex flex-col gap-4">
    <AdminResourcePage title="项目管理">
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
          新建项目
        </a-button>
        <a-button @click="authStore.loadProjects">刷新</a-button>
      </template>

      <AdminDataTable
        :columns="projectColumns"
        :data-source="pagedProjects"
        :loading="loadingProjects"
        :pagination="pagination"
        row-key="id"
        @change="onPageChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'project'">
            {{ record.name }}
          </template>
          <template v-else-if="column.key === 'description'">
            {{ record.description || "暂无描述" }}
          </template>
          <template v-else-if="column.key === 'key_bits'">
            <a-tag color="blue">{{ record.key_bits }} bit</a-tag>
          </template>
          <template v-else-if="column.key === 'updated_at'">
            {{ formatDateTime(record.updated_at) }}
          </template>
          <template v-else-if="column.key === 'created_at'">
            {{ formatDateTime(record.created_at) }}
          </template>
          <template v-else-if="column.key === 'actions'">
            <AdminRowActions :record="record" :actions="projectActions(record)" />
          </template>
        </template>
      </AdminDataTable>
    </AdminResourcePage>

    <a-modal
      v-model:open="modalOpen"
      :title="modalTitle"
      :confirm-loading="saving"
      @ok="submitProject"
    >
      <a-form layout="vertical" :model="projectForm">
        <a-form-item label="项目名称" required>
          <a-input v-model:value="projectForm.name" placeholder="例如：Acme Desktop Pro" />
        </a-form-item>
        <a-form-item label="项目描述">
          <a-textarea
            v-model:value="projectForm.description"
            :auto-size="{ minRows: 3, maxRows: 5 }"
            placeholder="记录产品线、客户范围或签发规则"
          />
        </a-form-item>
        <a-form-item
          label="RSA 密钥位数"
          :extra="modalMode === 'edit' ? '密钥位数创建后不可修改。' : ''"
        >
          <a-select
            v-model:value="projectForm.key_bits"
            :options="keyBitOptions"
            :disabled="modalMode === 'edit'"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="publicKeyOpen" width="760px" :footer="null">
      <template #title> {{ publicKeyProject?.name || "项目" }} 公钥 </template>
      <a-spin :spinning="publicKeyLoading">
        <a-space orientation="vertical" size="middle" class="flex">
          <a-alert
            type="info"
            show-icon
            message="将此公钥集成到客户端 SDK、服务端组件或 CLI，用于验证下载的 license 文件。"
          />
          <pre
            class="max-h-[420px] overflow-auto rounded-2 bg-[var(--surface-2)] p-4 text-left text-xs leading-6"
            >{{ publicKey || "暂无公钥" }}</pre
          >
          <a-space wrap>
            <a-button type="primary" :disabled="!publicKey" @click="copyPublicKey"
              >复制公钥</a-button
            >
            <a-button
              :disabled="!publicKeyProject"
              :loading="downloadingId === publicKeyProject?.id"
              @click="publicKeyProject && downloadPublicKey(publicKeyProject)"
            >
              下载 public.pem
            </a-button>
          </a-space>
        </a-space>
      </a-spin>
    </a-modal>
  </div>
</template>
