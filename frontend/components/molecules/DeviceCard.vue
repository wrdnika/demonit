<script setup lang="ts">
import type { Device } from '~/types/device'

const props = defineProps<{
  device: Device
}>()

const emit = defineEmits<{
  edit: [device: Device]
  delete: [device: Device]
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
  <BaseCard class="h-full transition hover:shadow-brutal-blue">
    <template #header>
      <div class="min-w-0">
        <h3 class="truncate text-base font-extrabold text-ink">
          {{ device.name }}
        </h3>
        <p class="text-xs font-semibold uppercase tracking-wide text-ink/45">
          {{ device.type }}
        </p>
      </div>
      <BaseStatusBadge :status="device.status" />
    </template>

    <div class="space-y-4">
      <MetricBar label="CPU" :value="device.cpu_usage" />
      <MetricBar label="RAM" :value="device.ram_usage" />
      <StatusPayloadList :payload="device.status_payload" :limit="4" />

      <dl class="grid grid-cols-1 gap-1 border-t-2 border-ink/10 pt-3 text-xs font-medium">
        <div class="flex justify-between gap-2">
          <dt class="text-ink/45">
            Last seen
          </dt>
          <dd class="font-mono font-semibold">
            {{ lastSeenLabel }}
          </dd>
        </div>
        <div class="flex justify-between gap-2">
          <dt class="text-ink/45">
            ID
          </dt>
          <dd class="truncate font-mono" :title="device.id">
            {{ device.id.slice(0, 8) }}…
          </dd>
        </div>
      </dl>
    </div>

    <template #footer>
      <div class="flex flex-wrap items-center gap-2">
        <NuxtLink
          :to="`/devices/${device.id}`"
          class="text-xs font-bold uppercase underline"
        >
          Open
        </NuxtLink>
        <button
          type="button"
          class="text-xs font-bold uppercase underline"
          @click="emit('edit', device)"
        >
          Edit
        </button>
        <button
          type="button"
          class="text-xs font-bold uppercase text-gred underline"
          @click="emit('delete', device)"
        >
          Delete
        </button>
      </div>
    </template>
  </BaseCard>
</template>
