<script setup lang="ts">
import { computed, ref } from 'vue'
import { BarsOutlined, BulbOutlined, DatabaseOutlined, UserOutlined } from '@antdv-next/icons'

import type { AdminMenuItem, AdminMenuTheme } from '../types'

const props = withDefaults(
  defineProps<{
    brand?: string
    currentTitle?: string
    menuItems: AdminMenuItem[]
    selectedKeys: string[]
    menuTheme?: AdminMenuTheme
    rootThemeClass?: string
    themeMenuItems: AdminMenuItem[]
    themeSelectedKeys: string[]
    userMenuItems: AdminMenuItem[]
    userName?: string
  }>(),
  {
    brand: 'Admin',
    currentTitle: '',
    menuTheme: 'light',
    rootThemeClass: '',
    userName: 'admin',
  },
)

const emit = defineEmits<{
  navigate: [key: string]
  themeClick: [key: string]
  userClick: [key: string]
}>()

const mobileMenuOpen = ref(false)

const pageTitle = computed(() => props.currentTitle || props.menuItems[0]?.label || '')

const onNavigate = ({ key }: { key: string }) => {
  emit('navigate', key)
  mobileMenuOpen.value = false
}

const onThemeMenuClick = ({ key }: { key: string }) => {
  emit('themeClick', key)
}

const onUserMenuClick = ({ key }: { key: string }) => {
  emit('userClick', key)
}
</script>

<template>
  <a-layout
    :class="[
      'min-h-screen bg-transparent text-[color:var(--text-primary)] [&_.ant-layout]:min-h-screen [&_.ant-layout]:bg-transparent',
      rootThemeClass,
    ]"
  >
    <a-layout-sider
      :width="248"
      :theme="menuTheme"
      class="hidden border-r border-[color:var(--surface-border)] !bg-[var(--shell-surface)] xl:!block"
    >
      <a-flex vertical class="h-full">
        <div class="px-2 py-3">
          <a-space align="center" size="middle" class="admin-brand px-2 pb-3 pt-1">
            <a-avatar shape="square" :size="36" class="admin-brand-mark">
              <template #icon>
                <slot name="brand-icon">
                  <DatabaseOutlined />
                </slot>
              </template>
            </a-avatar>
            <a-typography-title :level="4" class="!m-0 admin-brand-text">
              {{ brand }}
            </a-typography-title>
          </a-space>
          <a-menu
            mode="inline"
            :theme="menuTheme"
            :selected-keys="selectedKeys"
            :items="menuItems"
            class="border-r-0! bg-transparent! [&_.ant-menu-item]:!mx-0 [&_.ant-menu-item]:!my-1 [&_.ant-menu-item]:!w-auto [&_.ant-menu-item]:rounded-[var(--app-radius-sm,4px)] [&_.ant-menu-submenu-title]:!mx-0 [&_.ant-menu-submenu-title]:!my-1 [&_.ant-menu-submenu-title]:!w-auto [&_.ant-menu-submenu-title]:rounded-[var(--app-radius-sm,4px)]"
            @click="onNavigate"
          />
        </div>
      </a-flex>
    </a-layout-sider>

    <a-drawer
      :open="mobileMenuOpen"
      placement="left"
      :size="320"
      :footer="null"
      class="xl:hidden [&_.ant-drawer-content]:!border-[color:var(--surface-border)] [&_.ant-drawer-content]:!bg-[var(--shell-surface)]"
      @close="mobileMenuOpen = false"
    >
      <a-space orientation="vertical" size="middle" class="flex">
        <a-space align="center" size="middle" class="admin-brand">
          <a-avatar shape="square" :size="40" class="admin-brand-mark">
            <template #icon>
              <slot name="brand-icon">
                <DatabaseOutlined />
              </slot>
            </template>
          </a-avatar>
          <a-typography-title :level="4" class="!m-0 admin-brand-text">
            {{ brand }}
          </a-typography-title>
        </a-space>
        <a-menu
          mode="inline"
          :theme="menuTheme"
          :selected-keys="selectedKeys"
          :items="menuItems"
          class="border-r-0! bg-transparent! [&_.ant-menu-item]:!mx-0 [&_.ant-menu-item]:!my-1 [&_.ant-menu-item]:!w-auto [&_.ant-menu-item]:rounded-[var(--app-radius-sm,4px)] [&_.ant-menu-submenu-title]:!mx-0 [&_.ant-menu-submenu-title]:!my-1 [&_.ant-menu-submenu-title]:!w-auto [&_.ant-menu-submenu-title]:rounded-[var(--app-radius-sm,4px)]"
          @click="onNavigate"
        />
      </a-space>
    </a-drawer>

    <a-layout class="min-h-screen bg-transparent">
      <a-layout-header
        class="sticky top-0 z-20 h-16 border-b border-[color:var(--surface-border)] !bg-[var(--shell-surface)] px-4 md:px-6 xl:px-8"
      >
        <a-flex
          justify="space-between"
          align="center"
          gap="middle"
          class="app-header-inner mx-auto h-full w-full max-w-[1600px]"
        >
          <a-space align="center" size="middle">
            <a-button
              type="text"
              class="inline-flex h-9 shadow-none hover:shadow-none xl:hidden"
              aria-label="打开菜单"
              @click="mobileMenuOpen = true"
            >
              <template #icon>
                <BarsOutlined />
              </template>
            </a-button>
            <a-typography-title :level="5" class="!m-0 admin-page-title">
              {{ pageTitle }}
            </a-typography-title>
          </a-space>

          <a-space size="small">
            <slot name="topbar-extra" />
            <a-dropdown
              :trigger="['click']"
              placement="bottomRight"
              :menu="{
                items: themeMenuItems,
                selectable: true,
                selectedKeys: themeSelectedKeys,
                onClick: onThemeMenuClick,
              }"
            >
              <a-button
                type="text"
                class="inline-flex h-9 shadow-none hover:shadow-none"
                aria-label="主题设置"
              >
                <template #icon>
                  <BulbOutlined />
                </template>
              </a-button>
            </a-dropdown>

            <a-dropdown
              :trigger="['click']"
              placement="bottomRight"
              :menu="{ items: userMenuItems, onClick: onUserMenuClick }"
            >
              <a-button type="text" class="h-9 px-[14px] shadow-none hover:shadow-none">
                <a-space size="small">
                  <UserOutlined />
                  <span>{{ userName }}</span>
                </a-space>
              </a-button>
            </a-dropdown>
          </a-space>
        </a-flex>
      </a-layout-header>

      <a-layout-content class="app-content p-4 md:px-6 xl:px-8">
        <div class="app-content-inner mx-auto w-full max-w-[1600px]">
          <slot />
        </div>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>
