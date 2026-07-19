<script setup lang="ts">
import type { Device } from '~/types/device'

definePageMeta({
  layout: 'default',
})

const store = useDeviceStore()

useHead({
  title: 'Overview · Demonit',
})

const formOpen = ref(false)
const deleteOpen = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const activeDevice = ref<Device | null>(null)

const stats = computed(() => [
  { label: 'Total', value: store.totalCount, tone: 'white' as const },
  { label: 'Online', value: store.onlineCount, tone: 'lime' as const },
  { label: 'Offline', value: store.offlineCount, tone: 'pink' as const },
  { label: 'Alerts', value: store.activeAlerts.length, tone: 'sun' as const },
])

function openCreate() {
  formMode.value = 'create'
  activeDevice.value = null
  formOpen.value = true
}

function openEdit(device: Device) {
  formMode.value = 'edit'
  activeDevice.value = device
  formOpen.value = true
}

function openDelete(device: Device) {
  activeDevice.value = device
  deleteOpen.value = true
}
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-6 animate-fade-in">
    <header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
      <div class="space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">
          Dashboard overview
        </h2>
        <p class="max-w-2xl text-sm font-medium text-ink/70">
          Live status from the Go API. Add devices, then push heartbeats — stale nodes flip OFFLINE in ~30s.
        </p>
      </div>
      <button type="button" class="brutal-btn self-start" @click="openCreate">
        Add device
      </button>
    </header>

    <section
      v-if="store.error"
      class="rounded-brutal border-3 border-ink bg-gred px-4 py-3 text-sm font-bold text-white shadow-brutal-red"
      role="alert"
    >
      {{ store.error }}
    </section>

    <section aria-label="Fleet statistics" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <BaseCard v-for="stat in stats" :key="stat.label" :tone="stat.tone">
        <p class="text-xs font-bold uppercase tracking-wide">
          {{ stat.label }}
        </p>
        <p class="mt-2 font-mono text-4xl font-bold">
          {{ store.loading && store.totalCount === 0 ? '—' : stat.value }}
        </p>
      </BaseCard>
    </section>

    <section aria-labelledby="recent-devices-heading">
      <BaseCard :padding="false">
        <template #header>
          <h3 id="recent-devices-heading" class="text-sm font-bold uppercase tracking-wide">
            Fleet snapshot
          </h3>
          <NuxtLink to="/devices" class="text-xs font-bold uppercase underline">
            View all
          </NuxtLink>
        </template>

        <div
          v-if="store.loading && store.devices.length === 0"
          class="px-5 py-10 text-center text-sm font-bold"
        >
          Loading devices…
        </div>

        <div
          v-else-if="store.devices.length === 0"
          class="px-5 py-10 text-center text-sm font-bold"
        >
          No devices yet —
          <button type="button" class="underline" @click="openCreate">
            add one
          </button>
        </div>

        <div
          v-else
          class="grid gap-4 p-4 sm:grid-cols-2"
        >
          <DeviceCard
            v-for="device in store.devices.slice(0, 4)"
            :key="device.id"
            :device="device"
            @edit="openEdit"
            @delete="openDelete"
          />
        </div>
      </BaseCard>
    </section>

    <DeviceFormModal
      :open="formOpen"
      :mode="formMode"
      :device="activeDevice"
      @close="formOpen = false"
    />
    <DeleteDeviceModal
      :open="deleteOpen"
      :device="activeDevice"
      @close="deleteOpen = false"
    />
  </div>
</template>
