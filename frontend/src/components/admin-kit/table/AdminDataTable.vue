<script setup lang="ts">
import { computed, useSlots } from 'vue'

import type { AdminPaginationState } from '../types'

defineOptions({ inheritAttrs: false })

const props = withDefaults(
  defineProps<{
    columns: any[]
    dataSource: any[]
    rowKey?: string | ((record: any, index?: number) => string | number)
    loading?: boolean
    pagination: AdminPaginationState
    pageSizeOptions?: Array<string | number>
    showQuickJumper?: boolean
    showSizeChanger?: boolean
    showTotal?: (total: number, range: [number, number]) => string
  }>(),
  {
    rowKey: 'id',
    loading: false,
    pageSizeOptions: () => [10, 20, 50, 100],
    showQuickJumper: true,
    showSizeChanger: true,
    showTotal: undefined,
  },
)

const emit = defineEmits<{
  change: [page: number, pageSize: number]
}>()

const slots = useSlots()
const tableSlots = computed(() => slots)
const resolvedShowTotal = computed(
  () =>
    props.showTotal ??
    ((total: number, range: [number, number]) => `第 ${range[0]}-${range[1]} 条 / 共 ${total} 条`),
)

const onPageChange = (page: number, pageSize: number) => {
  emit('change', page, pageSize)
}
</script>

<template>
  <a-table
    v-bind="$attrs"
    class="bg-transparent [&_.ant-btn-link]:px-0 [&_.ant-table-cell]:align-middle [&_.ant-table-cell]:whitespace-nowrap [&_.ant-table-container]:!border-0 [&_.ant-table-container]:!rounded-none [&_.ant-table-container]:bg-transparent [&_.ant-table-content]:overflow-x-auto [&_.ant-table]:bg-transparent"
    :columns="columns"
    :data-source="dataSource"
    :row-key="rowKey"
    :loading="loading"
    :pagination="false"
    :scroll="{ x: 'max-content' }"
  >
    <template v-for="(_, name) in tableSlots" #[name]="slotProps">
      <slot :name="name" v-bind="slotProps || {}" />
    </template>
  </a-table>

  <a-pagination
    class="mt-4 flex flex-wrap justify-end gap-y-2.5"
    align="end"
    :current="pagination.current"
    :page-size="pagination.pageSize"
    :total="pagination.total"
    :page-size-options="pageSizeOptions"
    :show-quick-jumper="showQuickJumper"
    :show-size-changer="showSizeChanger"
    :show-total="resolvedShowTotal"
    @change="onPageChange"
  />
</template>
