<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    value: string | number | boolean | null | undefined
    activeValue?: string | number | boolean
    activeText?: string
    inactiveText?: string
    activeColor?: string
    inactiveColor?: string
  }>(),
  {
    activeValue: 1,
    activeText: '启用',
    inactiveText: '停用',
    activeColor: 'success',
    inactiveColor: 'default',
  },
)

const active = computed(() => props.value === props.activeValue)
</script>

<template>
  <a-tag :color="active ? activeColor : inactiveColor">
    <slot :active="active">
      {{ active ? activeText : inactiveText }}
    </slot>
  </a-tag>
</template>
