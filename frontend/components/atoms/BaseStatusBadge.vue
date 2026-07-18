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
  'inline-flex items-center gap-1.5 rounded font-medium uppercase tracking-wide',
  props.size === 'sm' ? 'px-2 py-0.5 text-xs' : 'px-2.5 py-1 text-xs',
  isOnline.value
    ? 'bg-online-soft text-online'
    : 'bg-offline-soft text-offline',
])
</script>

<template>
  <span
    :class="badgeClass"
    role="status"
    :aria-label="`Device status: ${status}`"
  >
    <span
      class="size-1.5 shrink-0 rounded-full"
      :class="isOnline ? 'bg-online animate-pulse-dot' : 'bg-offline'"
      aria-hidden="true"
    />
    {{ status }}
  </span>
</template>
