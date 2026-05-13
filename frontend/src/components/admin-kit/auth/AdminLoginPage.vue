<script setup lang="ts">
import type { FormInstance } from 'antdv-next'
import { h, reactive, ref } from 'vue'
import { LockOutlined, SafetyCertificateOutlined, UserOutlined } from '@antdv-next/icons'

import type { AdminLoginHighlight, AdminLoginPayload } from '../types'

withDefaults(
  defineProps<{
    loading?: boolean
    errorMessage?: string
    highlights?: AdminLoginHighlight[]
    eyebrow?: string
    brandTitle?: string
    headline?: string
    headlineStrong?: string
    description?: string
    footnote?: string
    formTitle?: string
    formDescription?: string
    submitText?: string
  }>(),
  {
    loading: false,
    errorMessage: '',
    highlights: () => [],
    eyebrow: 'Admin Console',
    brandTitle: '管理后台',
    headline: '统一、清晰地管理',
    headlineStrong: '后台系统的关键数据与操作',
    description: '使用标准化登录、导航、搜索、表格和操作反馈，让后台页面保持一致、清楚、容易维护。',
    footnote: '适合需要快速搭建后台管理入口的业务系统。',
    formTitle: '欢迎登录',
    formDescription: '请输入管理员账号和密码。',
    submitText: '进入控制台',
  },
)

const emit = defineEmits<{
  submit: [payload: AdminLoginPayload]
}>()

const formState = reactive<AdminLoginPayload>({
  username: '',
  password: '',
})
const formRef = ref<FormInstance>()

const rules: Record<string, any[]> = {
  username: [
    {
      required: true,
      message: '请输入用户名',
      trigger: ['blur', 'change'],
    },
    { min: 3, message: '用户名至少 3 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: ['blur', 'change'] },
    { min: 6, message: '密码至少 6 个字符', trigger: 'blur' },
  ],
}

const submit = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }
  emit('submit', {
    username: formState.username.trim(),
    password: formState.password,
  })
}
</script>

<template>
  <a-layout
    class="box-border min-h-dvh overflow-x-hidden bg-[var(--page-bg)] px-4 py-5 sm:px-6 lg:overflow-hidden lg:px-10 lg:py-0"
  >
    <div
      class="mx-auto flex min-h-[calc(100dvh-2.5rem)] w-full max-w-6xl items-center lg:min-h-dvh"
    >
      <a-row :gutter="[56, 40]" align="middle" class="w-full">
        <a-col :xs="24" :lg="13" class="order-2 lg:order-1">
          <div class="mx-auto max-w-[640px] lg:pr-8">
            <a-space orientation="vertical" size="large" class="flex">
              <a-space align="center" size="middle">
                <a-avatar :size="56">
                  <template #icon>
                    <SafetyCertificateOutlined class="text-2xl" />
                  </template>
                </a-avatar>
                <a-space orientation="vertical" size="small" class="flex">
                  <a-typography-text
                    class="text-[12px] tracking-[0.24em] text-[color:var(--text-secondary)] uppercase"
                  >
                    {{ eyebrow }}
                  </a-typography-text>
                  <a-typography-title :level="4" class="!m-0">
                    {{ brandTitle }}
                  </a-typography-title>
                </a-space>
              </a-space>

              <a-space orientation="vertical" size="middle" class="flex">
                <a-typography-title
                  :level="1"
                  class="!m-0 max-w-[620px] !text-4xl !leading-[1.18] md:!text-5xl xl:!text-[52px]"
                >
                  {{ headline }}
                  <a-typography-text strong class="block">
                    {{ headlineStrong }}
                  </a-typography-text>
                </a-typography-title>
                <a-typography-paragraph
                  class="!m-0 max-w-[560px] text-[15px] text-[color:var(--text-secondary)] !leading-[1.9]"
                >
                  {{ description }}
                </a-typography-paragraph>
              </a-space>

              <div v-if="highlights.length" class="border-y border-[color:var(--surface-border)]">
                <div
                  v-for="(item, index) in highlights"
                  :key="item.title"
                  class="py-5"
                  :class="
                    index !== highlights.length - 1
                      ? 'border-b border-[color:var(--surface-border)]'
                      : ''
                  "
                >
                  <a-space orientation="vertical" size="small" class="flex">
                    <a-space align="center" size="small">
                      <component
                        :is="item.icon"
                        v-if="item.icon"
                        class="text-base text-[color:var(--text-secondary)]"
                      />
                      <a-typography-title :level="5" class="!m-0">
                        {{ item.title }}
                      </a-typography-title>
                    </a-space>
                    <a-typography-paragraph
                      class="!m-0 max-w-[480px] pl-6 text-[color:var(--text-secondary)] !leading-[1.85]"
                    >
                      {{ item.note }}
                    </a-typography-paragraph>
                  </a-space>
                </div>
              </div>

              <a-typography-paragraph
                class="!m-0 max-w-[560px] text-sm text-[color:var(--text-secondary)] !leading-[1.85]"
              >
                {{ footnote }}
              </a-typography-paragraph>
            </a-space>
          </div>
        </a-col>

        <a-col :xs="24" :lg="11" class="order-1 lg:order-2">
          <div
            class="mx-auto w-full max-w-[400px] rounded-2xl border border-[color:var(--surface-border)] bg-[var(--surface-1)] p-6 shadow-sm sm:p-8"
          >
            <a-space orientation="vertical" size="large" class="flex">
              <a-space orientation="vertical" size="small" class="flex">
                <a-typography-title :level="2" class="!m-0 !leading-[1.18]">
                  {{ formTitle }}
                </a-typography-title>
                <a-typography-paragraph
                  class="!m-0 text-sm text-[color:var(--text-secondary)] !leading-[1.75]"
                >
                  {{ formDescription }}
                </a-typography-paragraph>
              </a-space>

              <a-form
                ref="formRef"
                :model="formState"
                :rules="rules"
                layout="vertical"
                scroll-to-first-error
                @finish="submit"
              >
                <a-form-item label="用户名" name="username" has-feedback html-for="login-username">
                  <a-input
                    id="login-username"
                    v-model:value="formState.username"
                    name="username"
                    size="large"
                    :prefix="h(UserOutlined)"
                    placeholder="请输入用户名"
                    autocomplete="username"
                  />
                </a-form-item>
                <a-form-item label="密码" name="password" has-feedback html-for="login-password">
                  <a-input-password
                    id="login-password"
                    v-model:value="formState.password"
                    name="password"
                    size="large"
                    :prefix="h(LockOutlined)"
                    placeholder="请输入密码"
                    autocomplete="current-password"
                  />
                </a-form-item>
                <a-alert
                  v-if="errorMessage"
                  class="mb-4"
                  type="error"
                  :message="errorMessage"
                  show-icon
                />
                <a-button
                  type="primary"
                  size="large"
                  block
                  :loading="loading"
                  html-type="submit"
                  class="h-11"
                >
                  {{ submitText }}
                </a-button>
              </a-form>
            </a-space>
          </div>
        </a-col>
      </a-row>
    </div>
  </a-layout>
</template>
