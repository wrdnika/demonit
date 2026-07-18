<script setup lang="ts">
definePageMeta({
  layout: 'default',
})

const store = useDeviceStore()

useHead({
  title: 'Overview · Demonit',
})

const stats = computed(() => [
  { label: 'Total devices', value: store.totalCount, tone: 'text-surface-900' },
  { label: 'Online', value: store.onlineCount, tone: 'text-online' },
  { label: 'Offline', value: store.offlineCount, tone: 'text-offline' },
  { label: 'Active alerts', value: store.activeAlerts.length, tone: 'text-alert' },
])
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-6 animate-fade-in">
    <header class="space-y-1">
      <h2 class="text-2xl font-semibold tracking-tight text-surface-900">
        Dashboard overview
      </h2>
      <p class="text-sm text-surface-800/60">
        Live fleet health from the monitoring API. Offline transitions raise alerts automatically.
      </p>
    </header>

    <section
      v-if="store.error"
      class="rounded-lg border border-alert/20 bg-alert-soft px-4 py-3 text-sm text-alert"
      role="alert"
    >
      {{ store.error }}
    </section>

    <section aria-label="Fleet statistics" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <BaseCard v-for="stat in stats" :key="stat.label">
        <p class="text-xs font-medium uppercase tracking-wide text-surface-800/50">
          {{ stat.label }}
        </p>
        <p class="mt-2 font-mono text-3xl font-semibold" :class="stat.tone">
          {{ store.loading && store.totalCount === 0 ? '—' : stat.value }}
        </p>
      </BaseCard>
    </section>

    <section aria-labelledby="recent-devices-heading">
      <BaseCard :padding="false">
        <template #header>
          <h3 id="recent-devices-heading" class="text-sm font-semibold text-surface-900">
            Fleet snapshot
          </h3>
          <NuxtLink
            to="/devices"
            class="text-xs font-medium text-accent-strong transition hover:text-accent"
          >
            View all devices
          </NuxtLink>
        </template>

        <div
          v-if="store.loading && store.devices.length === 0"
          class="px-5 py-10 text-center text-sm text-surface-800/50"
        >
          Loading devices…
        </div>

        <div
          v-else-if="store.devices.length === 0"
          class="px-5 py-10 text-center text-sm text-surface-800/50"
        >
          No devices registered yet.
        </div>

        <div
          v-else
          class="grid gap-4 p-5 sm:grid-cols-2 lg:grid-cols-3"
        >
          <TransitionGroup
            enter-active-class="transition duration-300 ease-out"
            enter-from-class="opacity-0 translate-y-2"
            enter-to-class="opacity-100 translate-y-0"
            leave-active-class="transition duration-200 ease-in"
            leave-from-class="opacity-100"
            leave-to-class="opacity-0"
          >
            <DeviceCard
              v-for="device in store.devices.slice(0, 6)"
              :key="device.id"
              :device="device"
            />
          </TransitionGroup>
        </div>
      </BaseCard>
    </section>
  </div>
</template>
