<script setup lang="ts">
import { App as AntApp, ConfigProvider } from "antdv-next";
import { storeToRefs } from "pinia";
import { computed, h, onMounted, watch, watchEffect } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  AppstoreOutlined,
  DatabaseOutlined,
  SafetyCertificateOutlined,
  TagsOutlined,
} from "@antdv-next/icons";

import { AdminShell, type AdminMenuItem } from "../components/admin-kit";
import { useAuthStore } from "../stores/auth";
import { useThemeStore } from "../stores/theme";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const themeStore = useThemeStore();
const { userName } = storeToRefs(authStore);
const { antdTheme, theme, rootThemeClass } = storeToRefs(themeStore);

const menuItems: AdminMenuItem[] = [
  {
    key: "projects",
    label: "项目管理",
    icon: DatabaseOutlined,
    hint: "签名密钥与产品线",
  },
  {
    key: "features",
    label: "功能管理",
    icon: TagsOutlined,
    hint: "全局功能点",
  },
  {
    key: "plans",
    label: "套餐管理",
    icon: AppstoreOutlined,
    hint: "功能与额度快照",
  },
  {
    key: "licenses",
    label: "License 管理",
    icon: SafetyCertificateOutlined,
    hint: "签发、下载与追踪",
  },
];

const consoleRouteNames = new Set(["projects", "licenses", "features", "plans"]);
const themeMenuItems: AdminMenuItem[] = [
  { key: "light", label: "明亮主题" },
  { key: "dark", label: "深色主题" },
];
const userMenuItems: AdminMenuItem[] = [{ key: "logout", label: "退出登录", danger: true }];

const isLoginRoute = computed(() => route.name === "login");
const activeKey = computed(() =>
  typeof route.name === "string" && consoleRouteNames.has(route.name) ? route.name : "projects",
);
const selectedKeys = computed(() => [activeKey.value]);
const themeSelectedKeys = computed(() => [theme.value]);
const currentTitle = computed(
  () => menuItems.find((item) => item.key === activeKey.value)?.label || "控制台",
);

const handleNavigate = async (key: string) => {
  if (consoleRouteNames.has(key)) {
    await router.push({ name: key });
  }
};

const handleThemeClick = (key: string) => {
  if (key === "light" || key === "dark") {
    themeStore.setTheme(key);
  }
};

const handleUserClick = async (key: string) => {
  if (key === "logout") {
    authStore.logout();
    await router.replace({ name: "login" });
  }
};

onMounted(() => {
  authStore.loadProjects();
});

watch(
  () => authStore.isAuthenticated,
  (isAuthenticated) => {
    if (!isAuthenticated && !isLoginRoute.value) router.replace({ name: "login" });
  },
);

watchEffect(() => {
  ConfigProvider.config({
    holderRender: (children) =>
      h(ConfigProvider, { theme: antdTheme.value }, () =>
        h(AntApp, { class: ["app-shell", rootThemeClass.value] }, () => children),
      ),
  });
});
</script>

<template>
  <a-config-provider :theme="antdTheme">
    <a-app :class="['app-shell', rootThemeClass]">
      <router-view v-if="isLoginRoute" />
      <AdminShell
        v-else
        brand="Lykn"
        :current-title="currentTitle"
        :menu-items="menuItems"
        :selected-keys="selectedKeys"
        :menu-theme="theme"
        :root-theme-class="rootThemeClass"
        :theme-menu-items="themeMenuItems"
        :theme-selected-keys="themeSelectedKeys"
        :user-menu-items="userMenuItems"
        :user-name="userName"
        @navigate="handleNavigate"
        @theme-click="handleThemeClick"
        @user-click="handleUserClick"
      >
        <template #brand-icon>
          <SafetyCertificateOutlined />
        </template>
        <router-view />
      </AdminShell>
    </a-app>
  </a-config-provider>
</template>

<style>
.app-shell {
  min-height: 100vh;
  color: var(--text-primary);
  background: var(--page-bg);
}
</style>
