<script setup lang="ts">
const store = useDeviceStore()
const config = useRuntimeConfig()

const sidebarOpen = ref(false)

let pollTimer: ReturnType<typeof setInterval> | null = null
const { connect, disconnect } = useAlertStream((payload) => {
  store.applyOfflineEvent(payload)
})

onMounted(async () => {
  await store.refreshConnection()
  await store.fetchDevices()
  connect()

  // Keep a slower poll for metric freshness; offline alerts arrive via SSE.
  const interval = Number(config.public.pollIntervalMs) || 15000
  pollTimer = setInterval(() => {
    store.fetchDevices()
  }, interval)
})

onBeforeUnmount(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
  disconnect()
})
</script>

<template>
  <div class="flex min-h-screen bg-paper">
    <AppSidebar :open="sidebarOpen" @close="sidebarOpen = false" />

    <div class="flex min-w-0 flex-1 flex-col">
      <AppNavbar
        :connection-status="store.connectionStatus"
        :last-fetched-at="store.lastFetchedAt"
        :online-count="store.onlineCount"
        :offline-count="store.offlineCount"
        @toggle-menu="sidebarOpen = !sidebarOpen"
      >
        <template #title>
          <slot name="title">
            Monitoring
          </slot>
        </template>
      </AppNavbar>

      <main id="main-content" class="flex-1 overflow-auto p-5 sm:p-8">
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
