<script setup lang="ts">
import type { Device, MetricLog } from '~/types/device'
import { DeviceApiError, useDeviceApi } from '~/composables/useDeviceApi'

definePageMeta({
  layout: 'default',
})

const route = useRoute()
const api = useDeviceApi()

const deviceId = computed(() => String(route.params.id || ''))

const device = ref<Device | null>(null)
const metrics = ref<MetricLog[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

useHead({
  title: computed(() => device.value ? `${device.value.name} · Demonit` : 'Device · Demonit'),
})

async function load() {
  loading.value = true
  error.value = null
  try {
    const [d, m] = await Promise.all([
      api.getDevice(deviceId.value),
      api.listDeviceMetrics(deviceId.value, 50),
    ])
    device.value = d
    metrics.value = m
  }
  catch (err) {
    error.value = err instanceof DeviceApiError ? err.message : 'Failed to load device'
    device.value = null
    metrics.value = []
  }
  finally {
    loading.value = false
  }
}

onMounted(load)
watch(deviceId, load)

function formatTime(iso: string) {
  return new Date(iso).toLocaleString()
}

function payloadPreview(payload: Record<string, unknown> | string) {
  if (typeof payload === 'string') {
    return payload
  }
  try {
    return JSON.stringify(payload)
  }
  catch {
    return '{}'
  }
}
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-6 animate-fade-in">
    <div class="flex items-center gap-3">
      <NuxtLink to="/devices" class="brutal-btn-ghost text-xs">
        ← Back
      </NuxtLink>
      <button type="button" class="brutal-btn text-xs" :disabled="loading" @click="load">
        Refresh
      </button>
    </div>

    <section
      v-if="error"
      class="rounded-brutal border-3 border-ink bg-gred px-4 py-3 text-sm font-bold text-white shadow-brutal-red"
      role="alert"
    >
      {{ error }}
    </section>

    <div v-else-if="loading && !device" class="text-sm font-bold">
      Loading device…
    </div>

    <template v-else-if="device">
      <header class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div class="space-y-2">
          <p class="font-mono text-xs font-bold uppercase tracking-widest text-ink/50">
            {{ device.type }}
          </p>
          <h2 class="text-3xl font-bold tracking-tight">
            {{ device.name }}
          </h2>
          <p class="font-mono text-xs break-all">
            {{ device.id }}
          </p>
        </div>
        <BaseStatusBadge :status="device.status" />
      </header>

      <section class="grid gap-4 sm:grid-cols-3">
        <BaseCard tone="lime">
          <p class="text-xs font-bold uppercase tracking-wide">
            CPU
          </p>
          <p class="mt-2 font-mono text-3xl font-bold">
            {{ device.cpu_usage == null ? '—' : `${device.cpu_usage.toFixed(1)}%` }}
          </p>
        </BaseCard>
        <BaseCard tone="cyan">
          <p class="text-xs font-bold uppercase tracking-wide">
            RAM
          </p>
          <p class="mt-2 font-mono text-3xl font-bold">
            {{ device.ram_usage == null ? '—' : `${device.ram_usage.toFixed(1)}%` }}
          </p>
        </BaseCard>
        <BaseCard tone="sun">
          <p class="text-xs font-bold uppercase tracking-wide">
            Last seen
          </p>
          <p class="mt-2 font-mono text-sm font-bold">
            {{ formatTime(device.last_seen) }}
          </p>
        </BaseCard>
      </section>

      <section v-if="metrics.length" aria-labelledby="chart-heading">
        <BaseCard>
          <template #header>
            <h3 id="chart-heading" class="text-sm font-bold uppercase tracking-wide">
              CPU / RAM trend
            </h3>
            <span class="font-mono text-xs font-bold">
              last {{ Math.min(metrics.length, 30) }}
            </span>
          </template>
          <MetricsChart :metrics="metrics" />
        </BaseCard>
      </section>

      <section v-if="device.status_payload" aria-labelledby="payload-heading">
        <BaseCard>
          <template #header>
            <h3 id="payload-heading" class="text-sm font-bold uppercase tracking-wide">
              Hardware payload
            </h3>
          </template>
          <StatusPayloadList :payload="device.status_payload" :limit="20" />
        </BaseCard>
      </section>

      <section aria-labelledby="metrics-heading">
        <BaseCard :padding="false">
          <template #header>
            <h3 id="metrics-heading" class="text-sm font-bold uppercase tracking-wide">
              Metric history
            </h3>
            <span class="font-mono text-xs font-bold">
              {{ metrics.length }} samples
            </span>
          </template>

          <div v-if="metrics.length === 0" class="px-4 py-10 text-center text-sm font-bold">
            No heartbeats yet. POST /api/v1/heartbeat to start logging.
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full text-left text-sm">
              <thead class="border-b-3 border-ink bg-paper text-xs font-bold uppercase tracking-wide">
                <tr>
                  <th class="px-4 py-3">
                    Time
                  </th>
                  <th class="px-4 py-3">
                    CPU
                  </th>
                  <th class="px-4 py-3">
                    RAM
                  </th>
                  <th class="px-4 py-3">
                    Payload
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="row in metrics"
                  :key="row.id"
                  class="border-b-2 border-ink/20"
                >
                  <td class="px-4 py-3 font-mono text-xs whitespace-nowrap">
                    <time :datetime="row.timestamp">{{ formatTime(row.timestamp) }}</time>
                  </td>
                  <td class="px-4 py-3 font-mono font-bold">
                    {{ row.cpu_usage.toFixed(1) }}%
                  </td>
                  <td class="px-4 py-3 font-mono font-bold">
                    {{ row.ram_usage.toFixed(1) }}%
                  </td>
                  <td class="max-w-xs truncate px-4 py-3 font-mono text-xs" :title="payloadPreview(row.status_payload)">
                    {{ payloadPreview(row.status_payload) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </BaseCard>
      </section>
    </template>
  </div>
</template>
