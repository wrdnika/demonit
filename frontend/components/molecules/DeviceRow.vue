<script setup lang="ts">
import type { Device } from '~/types/device'

const props = defineProps<{
  device: Device
}>()

const lastSeenLabel = computed(() => {
  const then = new Date(props.device.last_seen).getTime()
  if (Number.isNaN(then)) {
    return '—'
  }
  return new Date(props.device.last_seen).toLocaleString()
})

function metricLabel(value?: number) {
  return value == null || Number.isNaN(value) ? '—' : `${value.toFixed(1)}%`
}
</script>

<template>
  <tr class="border-b border-surface-100 transition-colors duration-200 hover:bg-surface-50">
    <td class="px-4 py-3">
      <div class="min-w-0">
        <p class="truncate text-sm font-semibold text-surface-900">
          {{ device.name }}
        </p>
        <p class="truncate font-mono text-xs text-surface-800/45" :title="device.id">
          {{ device.id }}
        </p>
      </div>
    </td>
    <td class="px-4 py-3">
      <span class="rounded bg-surface-100 px-2 py-0.5 font-mono text-xs text-surface-800">
        {{ device.type }}
      </span>
    </td>
    <td class="px-4 py-3">
      <BaseStatusBadge :status="device.status" size="sm" />
    </td>
    <td class="px-4 py-3">
      <div class="w-28">
        <MetricBar label="CPU" :value="device.cpu_usage" />
      </div>
    </td>
    <td class="px-4 py-3 font-mono text-xs text-surface-800">
      {{ metricLabel(device.ram_usage) }}
    </td>
    <td class="px-4 py-3 font-mono text-xs text-surface-800/70">
      <time :datetime="device.last_seen">{{ lastSeenLabel }}</time>
    </td>
  </tr>
</template>
