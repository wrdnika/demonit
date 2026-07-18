<script setup lang="ts">
const store = useDeviceStore()
const config = useRuntimeConfig()

let pollTimer: ReturnType<typeof setInterval> | null = null

onMounted(async () => {
  await store.refreshConnection()
  await store.fetchDevices()

  const interval = Number(config.public.pollIntervalMs) || 15000
  pollTimer = setInterval(() => {
    store.fetchDevices()
  }, interval)
})

onBeforeUnmount(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
})
</script>

<template>
  <div class="flex min-h-screen bg-surface-50">
    <AppSidebar />

    <div class="flex min-w-0 flex-1 flex-col">
      <AppNavbar
        :connection-status="store.connectionStatus"
        :last-fetched-at="store.lastFetchedAt"
        :online-count="store.onlineCount"
        :offline-count="store.offlineCount"
      >
        <template #title>
          <slot name="title">
            Monitoring
          </slot>
        </template>
      </AppNavbar>

      <main id="main-content" class="flex-1 overflow-auto p-6">
        <slot />
      </main>
    </div>

    <AlertBanner
      :alerts="store.activeAlerts"
      @dismiss="store.dismissAlert"
      @dismiss-all="store.dismissAllAlerts"
    />
  </div>
</template>
