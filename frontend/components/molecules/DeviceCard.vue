<script setup lang="ts">
import type { Device } from '~/types/device'

const props = defineProps<{
  device: Device
}>()

const lastSeenLabel = computed(() => formatRelativeTime(props.device.last_seen))

function formatRelativeTime(iso: string): string {
  const then = new Date(iso).getTime()
  if (Number.isNaN(then)) {
    return 'unknown'
  }

  const deltaSec = Math.round((Date.now() - then) / 1000)
  if (deltaSec < 5) {
    return 'just now'
  }
  if (deltaSec < 60) {
    return `${deltaSec}s ago`
  }
  if (deltaSec < 3600) {
    return `${Math.floor(deltaSec / 60)}m ago`
  }
  if (deltaSec < 86400) {
    return `${Math.floor(deltaSec / 3600)}h ago`
  }
  return new Date(iso).toLocaleString()
}
</script>

<template>
  <BaseCard class="h-full transition duration-300 hover:border-accent/40">
    <template #header>
      <div class="min-w-0">
        <h3 class="truncate text-sm font-semibold text-surface-900">
          {{ device.name }}
        </h3>
        <p class="font-mono text-xs text-surface-800/50">
          {{ device.type }}
        </p>
      </div>
      <BaseStatusBadge :status="device.status" />
    </template>

    <div class="space-y-4">
      <MetricBar label="CPU" :value="device.cpu_usage" />
      <MetricBar label="RAM" :value="device.ram_usage" />

      <dl class="grid grid-cols-1 gap-1 text-xs text-surface-800/60">
        <div class="flex justify-between gap-2">
          <dt>Last seen</dt>
          <dd class="font-mono text-surface-900">
            {{ lastSeenLabel }}
          </dd>
        </div>
        <div class="flex justify-between gap-2">
          <dt>Device ID</dt>
          <dd class="truncate font-mono" :title="device.id">
            {{ device.id.slice(0, 8) }}…
          </dd>
        </div>
      </dl>
    </div>
  </BaseCard>
</template>
