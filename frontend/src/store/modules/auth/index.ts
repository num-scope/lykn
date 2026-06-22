import { computed, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { defineStore } from 'pinia';
import { useLoading } from '@sa/hooks';
import { fetchGetUserInfo, fetchLogin, fetchProjects } from '@/service/api';
import { useRouterPush } from '@/hooks/common/router';
import { localStg } from '@/utils/storage';
import { getServiceErrorMessage } from '@/utils/lykn';
import { SetupStoreId } from '@/enum';
import { $t } from '@/locales';
import { useRouteStore } from '../route';
import { useTabStore } from '../tab';
import { clearAuthStorage, getToken } from './shared';

const defaultRole = import.meta.env.VITE_STATIC_SUPER_ROLE || 'super';

function getUserInfoFromAccount(account?: Api.Auth.UserAccount): Api.Auth.UserInfo {
  return {
    userId: account ? String(account.id) : '',
    userName: account?.username || '',
    roles: account ? [defaultRole] : [],
    buttons: []
  };
}

export const useAuthStore = defineStore(SetupStoreId.Auth, () => {
  const route = useRoute();
  const routeStore = useRouteStore();
  const tabStore = useTabStore();
  const { toLogin, redirectFromLogin } = useRouterPush(false);
  const { loading: loginLoading, startLoading, endLoading } = useLoading();

  const token = ref(getToken());
  const loginError = ref('');
  const loadingProjects = ref(false);
  const projects = ref<Api.Lykn.ProjectRecord[]>([]);
  const selectedProjectId = ref<number>();

  const userInfo = reactive<Api.Auth.UserInfo>(getUserInfoFromAccount());

  const isStaticSuper = computed(() => {
    return import.meta.env.VITE_AUTH_ROUTE_MODE === 'static' && userInfo.roles.includes(defaultRole);
  });

  const isLogin = computed(() => Boolean(token.value));
  const isAuthenticated = computed(() => isLogin.value);
  const userName = computed(() => userInfo.userName || 'admin');

  function syncUserInfo(account: Api.Auth.UserAccount) {
    Object.assign(userInfo, getUserInfoFromAccount(account));
  }

  function ensureSelectedProject() {
    if (selectedProjectId.value && projects.value.some(project => project.id === selectedProjectId.value)) return;

    selectedProjectId.value = projects.value[0]?.id;
  }

  async function loadProjects() {
    if (!token.value) return;

    loadingProjects.value = true;
    try {
      const { data, error } = await fetchProjects();
      if (error || !data) throw error;

      projects.value = data;
      ensureSelectedProject();
    } catch (error) {
      window.$message?.error(getServiceErrorMessage(error, '加载项目列表失败'));
    } finally {
      loadingProjects.value = false;
    }
  }

  function selectProject(projectId: number) {
    selectedProjectId.value = projectId;
  }

  async function resetStore(redirectToLogin = true) {
    recordUserId();
    clearAuthStorage();

    token.value = '';
    loginError.value = '';
    projects.value = [];
    selectedProjectId.value = undefined;
    Object.assign(userInfo, getUserInfoFromAccount());

    if (redirectToLogin && !route.meta.constant) {
      await toLogin();
    }

    tabStore.cacheTabs();
    routeStore.resetStore();
  }

  function recordUserId() {
    if (userInfo.userId) {
      localStg.set('lastLoginUserId', userInfo.userId);
    }
  }

  function checkTabClear() {
    if (!userInfo.userId) return false;

    const lastLoginUserId = localStg.get('lastLoginUserId');
    if (!lastLoginUserId || lastLoginUserId !== userInfo.userId) {
      localStg.remove('globalTabs');
      tabStore.clearTabs();
      localStg.remove('lastLoginUserId');
      return true;
    }

    localStg.remove('lastLoginUserId');
    return false;
  }

  async function login(userNameValue: string, password: string, redirect = true) {
    startLoading();
    loginError.value = '';

    try {
      const { data, error } = await fetchLogin(userNameValue, password);
      if (error || !data) throw error;

      token.value = data.access_token;
      localStg.set('token', data.access_token);
      localStg.set('refreshToken', data.access_token);
      syncUserInfo(data.user);

      await loadProjects();
      await redirectFromLogin(checkTabClear() ? false : redirect);

      window.$notification?.success({
        title: $t('page.login.common.loginSuccess'),
        description: $t('page.login.common.welcomeBack', { userName: userInfo.userName }),
        duration: 4500
      });

      return true;
    } catch (error) {
      loginError.value = getServiceErrorMessage(error, '登录失败，请检查账号密码');
      window.$message?.error(loginError.value);
      await resetStore(false);
      return false;
    } finally {
      endLoading();
    }
  }

  async function loginByToken(loginToken: Api.Auth.LoginToken) {
    token.value = loginToken.token;
    localStg.set('token', loginToken.token);
    localStg.set('refreshToken', loginToken.refreshToken || loginToken.token);

    return getUserInfo();
  }

  async function getUserInfo() {
    const { data, error } = await fetchGetUserInfo();
    if (error || !data) {
      return false;
    }

    syncUserInfo(data);
    await loadProjects();
    return true;
  }

  async function initUserInfo() {
    const maybeToken = getToken();
    if (!maybeToken) return;

    token.value = maybeToken;
    const pass = await getUserInfo();

    if (!pass) {
      await resetStore(false);
    }
  }

  async function logout(showMessage = true) {
    await resetStore(false);
    if (showMessage) window.$message?.success($t('common.logout'));
  }

  return {
    token,
    userInfo,
    isStaticSuper,
    isLogin,
    isAuthenticated,
    userName,
    loginLoading,
    loginError,
    loadingProjects,
    projects,
    selectedProjectId,
    resetStore,
    login,
    loginByToken,
    getUserInfo,
    initUserInfo,
    loadProjects,
    logout,
    selectProject
  };
});
