import { computed, ref } from "vue";
import { defineStore } from "pinia";
import { message } from "antdv-next";

import {
  api,
  ApiError,
  clearStoredSession,
  getStoredSession,
  setStoredSession,
  type ProjectRecord,
  type StoredSession,
} from "../api/api";
import { getErrorMessage } from "../utils/format";

interface LoginPayload {
  username: string;
  password: string;
}

export const useAuthStore = defineStore("auth", () => {
  const session = ref<StoredSession | null>(getStoredSession());
  const loginLoading = ref(false);
  const loginError = ref("");
  const loadingProjects = ref(false);
  const projects = ref<ProjectRecord[]>([]);
  const selectedProjectId = ref<number | undefined>(undefined);

  const isAuthenticated = computed(() => Boolean(session.value));
  const userName = computed(() => session.value?.user.username || "admin");

  const clearInvalidProjectSelection = () => {
    if (
      selectedProjectId.value &&
      !projects.value.some((project) => project.id === selectedProjectId.value)
    ) {
      selectedProjectId.value = undefined;
    }
  };

  const logout = (showMessage = true) => {
    clearStoredSession();
    session.value = null;
    projects.value = [];
    selectedProjectId.value = undefined;
    if (showMessage) message.success("已退出登录");
  };

  const handleApiError = (error: unknown, fallback: string) => {
    if (error instanceof ApiError && error.status === 401) {
      logout(false);
      message.warning("登录已过期，请重新登录");
      return;
    }
    message.error(getErrorMessage(error, fallback));
  };

  const loadProjects = async () => {
    if (!session.value) return;
    loadingProjects.value = true;
    try {
      projects.value = await api.listProjects();
      clearInvalidProjectSelection();
    } catch (error) {
      handleApiError(error, "加载项目列表失败");
    } finally {
      loadingProjects.value = false;
    }
  };

  const login = async (payload: LoginPayload) => {
    loginLoading.value = true;
    loginError.value = "";
    try {
      const loginResponse = await api.login(payload);
      session.value = setStoredSession(loginResponse);
      message.success("登录成功");
      await loadProjects();
      return true;
    } catch (error) {
      loginError.value = getErrorMessage(error, "登录失败，请检查账号密码");
      return false;
    } finally {
      loginLoading.value = false;
    }
  };

  const selectProject = (projectId: number) => {
    selectedProjectId.value = projectId;
  };

  return {
    session,
    isAuthenticated,
    userName,
    loginLoading,
    loginError,
    loadingProjects,
    projects,
    selectedProjectId,
    loadProjects,
    login,
    logout,
    selectProject,
  };
});
