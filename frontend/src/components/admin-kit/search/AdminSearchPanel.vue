<script setup lang="ts">
import { DownOutlined, UpOutlined } from '@antdv-next/icons'
import { computed, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    fieldCount: number
    collapsedCount?: number
    submitText?: string
    resetText?: string
    submitDisabled?: boolean
    resetDisabled?: boolean
  }>(),
  {
    collapsedCount: 2,
    submitText: '查询',
    resetText: '重置',
    submitDisabled: false,
    resetDisabled: false,
  },
)

const emit = defineEmits<{
  finish: []
  reset: []
}>()

const expanded = ref(false)

const collapsible = computed(() => props.fieldCount > props.collapsedCount)
const isCollapsed = computed(() => collapsible.value && !expanded.value)

const submit = () => {
  emit('finish')
}

const reset = () => {
  emit('reset')
}

const toggleExpanded = () => {
  expanded.value = !expanded.value
}
</script>

<template>
  <a-form
    layout="horizontal"
    class="border border-[color:var(--surface-border)] rounded-2 bg-[var(--surface-1)] px-4 py-4 [&_.ant-form-item]:mb-0 [&_.ant-form-item-control]:min-w-0 [&_.ant-form-item-control]:w-full [&_.ant-form-item-label]:p-0 [&_.ant-form-item-label]:text-right [&_.ant-form-item-label]:text-[color:var(--text-secondary)] [&_.ant-form-item-label]:font-600 [&_.ant-form-item-label]:leading-none [&_.ant-form-item-row]:grid [&_.ant-form-item-row]:grid-cols-[88px_minmax(0,1fr)] [&_.ant-form-item-row]:items-center [&_.ant-form-item-row]:gap-x-3 [&_.ant-picker]:w-full [&_.ant-picker]:rounded-2 [&_.ant-select]:w-full [&_.ant-select-selector]:rounded-2"
    @finish="submit"
    @submit.prevent
  >
    <div class="grid grid-cols-1 items-end gap-x-4 gap-y-3 md:grid-cols-2 xl:grid-cols-3">
      <slot
        :expanded="expanded"
        :is-collapsed="isCollapsed"
        :visible-count="isCollapsed ? collapsedCount : fieldCount"
      />
      <div
        :class="[
          'flex flex-wrap items-center justify-end gap-2',
          collapsible && expanded && 'col-span-full',
        ]"
      >
        <a-button html-type="button" :disabled="resetDisabled" @click="reset">
          {{ resetText }}
        </a-button>
        <a-button html-type="button" type="primary" :disabled="submitDisabled" @click="submit">
          {{ submitText }}
        </a-button>
        <a-button
          v-if="collapsible"
          html-type="button"
          type="link"
          class="px-0"
          @click="toggleExpanded"
        >
          {{ expanded ? '收起' : '展开' }}
          <component :is="expanded ? UpOutlined : DownOutlined" class="ml-1 text-xs" />
        </a-button>
      </div>
    </div>
  </a-form>
</template>
