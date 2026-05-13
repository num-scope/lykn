<script setup lang="ts">
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { DatabaseOutlined, SafetyCertificateOutlined, UserOutlined } from "@antdv-next/icons";

import {
  AdminLoginPage,
  type AdminLoginHighlight,
  type AdminLoginPayload,
} from "../components/admin-kit";
import { useAuthStore } from "../stores/auth";

const router = useRouter();
const authStore = useAuthStore();
const { loginLoading, loginError } = storeToRefs(authStore);

const loginHighlights: AdminLoginHighlight[] = [
  {
    icon: SafetyCertificateOutlined,
    title: "签发可验证的 license",
    note: "基于项目私钥签名，客户端通过公钥完成离线验证。",
  },
  {
    icon: DatabaseOutlined,
    title: "按项目隔离授权资产",
    note: "每个项目拥有独立 RSA 密钥、license 列表与下载入口。",
  },
  {
    icon: UserOutlined,
    title: "面向用户的最小闭环",
    note: "登录、项目管理、license 签发、公钥和 license 下载都在一个控制台内完成。",
  },
];

const handleLogin = async (payload: AdminLoginPayload) => {
  if (await authStore.login(payload)) {
    await router.replace({ name: "projects" });
  }
};
</script>

<template>
  <AdminLoginPage
    :loading="loginLoading"
    :error-message="loginError"
    :highlights="loginHighlights"
    eyebrow="Lykn Console"
    brand-title="Lykn 授权平台"
    headline="管理 Python 项目的"
    headline-strong="License 签发与验证资产"
    description="为项目生成独立密钥，签发可离线验证的 license，并将公钥交付给 Python SDK 或 CLI 使用。"
    footnote="开发环境默认用户为 admin / admin123。首次部署后请立即修改默认凭据。"
    form-title="用户登录"
    form-description="登录后即可管理项目、公钥与 license。"
    submit-text="进入 Lykn 控制台"
    @submit="handleLogin"
  />
</template>
