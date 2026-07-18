<script setup lang="ts">
import type { DeviceStatus } from '~/types/device'

const props = withDefaults(defineProps<{
  status: DeviceStatus
  size?: 'sm' | 'md'
}>(), {
  size: 'md',
})

const isOnline = computed(() => props.status === 'ONLINE')

const badgeClass = computed(() => [
  'inline-flex items-center gap-2 rounded-full border-2 border-ink font-bold uppercase tracking-wide shadow-brutal-sm',
  props.size === 'sm' ? 'px-2.5 py-0.5 text-[10px]' : 'px-3 py-1 text-xs',
  isOnline.value ? 'bg-ggreen text-white' : 'bg-gred text-white',
])
</script>

<template>
  <span
    :class="badgeClass"
    role="status"
    :aria-label="`Device status: ${status}`"
  >
    <span
      class="size-2 shrink-0 rounded-full bg-white"
      :class="isOnline ? 'animate-pulse-dot' : 'opacity-70'"
      aria-hidden="true"
    />
    {{ status }}
  </span>
</template>
