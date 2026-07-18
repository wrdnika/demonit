/**
 * Ensures the Pinia store connection probe runs early on client navigations.
 * Data polling itself lives in the default layout to avoid duplicate timers.
 */
export default defineNuxtRouteMiddleware(async () => {
  if (import.meta.server) {
    return
  }

  const store = useDeviceStore()
  if (store.connectionStatus === 'checking' && store.devices.length === 0) {
    await store.refreshConnection()
  }
})
