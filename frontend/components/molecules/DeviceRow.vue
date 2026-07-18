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
  <tr class="border-b-2 border-ink/10 transition-colors hover:bg-accent-soft/40">
    <td class="px-4 py-3">
      <NuxtLink :to="`/devices/${device.id}`" class="block min-w-0 hover:underline">
        <p class="truncate text-sm font-bold">
          {{ device.name }}
        </p>
        <p class="truncate font-mono text-xs text-ink/40" :title="device.id">
          {{ device.id }}
        </p>
      </NuxtLink>
    </td>
    <td class="px-4 py-3">
      <span class="rounded-full border-2 border-ink bg-gyellow px-2.5 py-0.5 text-xs font-bold">
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
    <td class="px-4 py-3 font-mono text-xs font-bold">
      {{ metricLabel(device.ram_usage) }}
    </td>
    <td class="px-4 py-3 font-mono text-xs">
      <time :datetime="device.last_seen">{{ lastSeenLabel }}</time>
    </td>
  </tr>
</template>
