<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string
    cardSize?: 'default' | 'small'
  }>(),
  {
    title: '',
    cardSize: 'small',
  },
)
</script>

<template>
  <div class="flex flex-col gap-4">
    <slot name="search" />

    <a-card :size="cardSize" class="overflow-hidden rounded-2 shadow-none [&_.ant-card-body]:p-4">
      <div
        v-if="title || $slots.title || $slots.toolbar"
        class="mb-4 flex flex-wrap items-center justify-between gap-4 max-xl:items-stretch"
      >
        <div class="text-[16px] text-[color:var(--text-primary)] font-600 leading-6">
          <slot name="title">
            {{ title }}
          </slot>
        </div>
        <a-space v-if="$slots.toolbar" class="[&_.ant-btn]:h-9 [&_.ant-select]:min-w-[240px]" wrap>
          <slot name="toolbar" />
        </a-space>
      </div>

      <slot />
    </a-card>

    <slot name="after" />
  </div>
</template>
