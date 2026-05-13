<script setup lang="ts" generic="T extends Record<string, any>">
import { computed } from 'vue'

import type { AdminRowAction } from '../types'

const props = defineProps<{
  record: T
  actions: AdminRowAction<T>[]
}>()

const visibleActions = computed(() =>
  props.actions.filter((action) =>
    typeof action.visible === 'function' ? action.visible(props.record) : action.visible !== false,
  ),
)

const isDisabled = (action: AdminRowAction<T>) =>
  typeof action.disabled === 'function' ? action.disabled(props.record) : action.disabled === true

const getConfirm = (action: AdminRowAction<T>) =>
  typeof action.confirm === 'function' ? action.confirm(props.record) : action.confirm

const runAction = (action: AdminRowAction<T>) => {
  if (isDisabled(action)) return
  return action.onClick(props.record)
}
</script>

<template>
  <a-space class="[&_.ant-btn-link]:px-0">
    <template v-for="action in visibleActions" :key="action.key">
      <a-popconfirm
        v-if="getConfirm(action)"
        :title="getConfirm(action)"
        :ok-text="action.confirmOkText || '确定'"
        :cancel-text="action.confirmCancelText || '取消'"
        :disabled="isDisabled(action)"
        @confirm="runAction(action)"
      >
        <a-button type="link" :danger="action.danger" :disabled="isDisabled(action)">
          <template v-if="action.icon" #icon>
            <component :is="action.icon" />
          </template>
          {{ action.label }}
        </a-button>
      </a-popconfirm>
      <a-button
        v-else
        type="link"
        :danger="action.danger"
        :disabled="isDisabled(action)"
        @click="runAction(action)"
      >
        <template v-if="action.icon" #icon>
          <component :is="action.icon" />
        </template>
        {{ action.label }}
      </a-button>
    </template>
  </a-space>
</template>
