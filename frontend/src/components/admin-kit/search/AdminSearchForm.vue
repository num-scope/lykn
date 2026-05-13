<script setup lang="ts">
import { computed } from 'vue'

import AdminSearchPanel from './AdminSearchPanel.vue'
import type { AdminSearchField } from '../types'

const props = withDefaults(
  defineProps<{
    model?: Record<string, any>
    fields?: AdminSearchField[]
    fieldCount?: number
    collapsedCount?: number
    submitText?: string
    resetText?: string
    submitDisabled?: boolean
    resetDisabled?: boolean
  }>(),
  {
    model: () => ({}),
    fields: () => [],
    collapsedCount: 2,
    submitText: '查询',
    resetText: '重置',
    submitDisabled: false,
    resetDisabled: false,
  },
)

const emit = defineEmits<{
  search: []
  reset: []
  fieldChange: [key: string, value: any]
}>()

const resolvedFieldCount = computed(() => props.fieldCount ?? props.fields.length)

const getVisibleFields = (visibleCount: number) => props.fields.slice(0, visibleCount)

const getFieldValue = (field: AdminSearchField) => props.model[field.key]

const setFieldValue = (field: AdminSearchField, value: any) => {
  props.model[field.key] = value
  emit('fieldChange', field.key, value)
}

const onSearch = () => {
  emit('search')
}

const onReset = () => {
  emit('reset')
}
</script>

<template>
  <AdminSearchPanel
    :field-count="resolvedFieldCount"
    :collapsed-count="collapsedCount"
    :submit-text="submitText"
    :reset-text="resetText"
    :submit-disabled="submitDisabled"
    :reset-disabled="resetDisabled"
    @finish="onSearch"
    @reset="onReset"
  >
    <template #default="{ visibleCount }">
      <template v-if="fields.length">
        <a-form-item
          v-for="field in getVisibleFields(visibleCount)"
          :key="field.key"
          :label="field.label"
        >
          <slot
            v-if="field.type === 'custom'"
            :name="`field-${field.key}`"
            :field="field"
            :value="getFieldValue(field)"
            :set-value="(value: any) => setFieldValue(field, value)"
          />
          <a-select
            v-else-if="field.type === 'select'"
            :value="getFieldValue(field)"
            :allow-clear="field.clearable ?? true"
            :disabled="field.disabled"
            class="w-full"
            :placeholder="field.placeholder"
            :options="field.options"
            @update:value="(value: any) => setFieldValue(field, value)"
          />
          <a-range-picker
            v-else-if="field.type === 'dateRange'"
            :value="getFieldValue(field)"
            :allow-clear="field.clearable ?? true"
            :disabled="field.disabled"
            class="w-full"
            :placeholder="field.placeholder ? [field.placeholder, field.placeholder] : undefined"
            @update:value="(value: any) => setFieldValue(field, value)"
          />
          <a-input-number
            v-else-if="field.type === 'number'"
            :value="getFieldValue(field)"
            :disabled="field.disabled"
            class="w-full"
            :placeholder="field.placeholder"
            @update:value="(value: any) => setFieldValue(field, value)"
            @pressEnter="onSearch"
          />
          <a-input
            v-else
            :value="getFieldValue(field)"
            :allow-clear="field.clearable ?? true"
            :disabled="field.disabled"
            class="w-full"
            :placeholder="field.placeholder"
            @update:value="(value: any) => setFieldValue(field, value)"
            @pressEnter="onSearch"
          />
        </a-form-item>
      </template>
      <slot v-else />
    </template>
  </AdminSearchPanel>
</template>
