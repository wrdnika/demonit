import { defineStore } from 'pinia'
import type { ConnectionStatus, Device, DeviceAlert, DeviceStatus, DeviceType } from '~/types/device'
import { DeviceApiError, useDeviceApi } from '~/composables/useDeviceApi'

interface DevicesState {
  devices: Device[]
  alerts: DeviceAlert[]
  loading: boolean
  error: string | null
  connectionStatus: ConnectionStatus
  lastFetchedAt: string | null
  filterStatus: DeviceStatus | 'ALL'
  filterType: DeviceType | 'ALL'
  searchQuery: string
}

export const useDeviceStore = defineStore('devices', {
  state: (): DevicesState => ({
    devices: [],
    alerts: [],
    loading: false,
    error: null,
    connectionStatus: 'checking',
    lastFetchedAt: null,
    filterStatus: 'ALL',
    filterType: 'ALL',
    searchQuery: '',
  }),

  getters: {
    onlineCount: (state): number =>
      state.devices.filter(d => d.status === 'ONLINE').length,

    offlineCount: (state): number =>
      state.devices.filter(d => d.status === 'OFFLINE').length,

    totalCount: (state): number => state.devices.length,

    activeAlerts: (state): DeviceAlert[] =>
      state.alerts.filter(a => !a.dismissed),

    filteredDevices: (state): Device[] => {
      const query = state.searchQuery.trim().toLowerCase()

      return state.devices.filter((device) => {
        const statusOk = state.filterStatus === 'ALL' || device.status === state.filterStatus
        const typeOk = state.filterType === 'ALL' || device.type === state.filterType
        const searchOk = !query
          || device.name.toLowerCase().includes(query)
          || device.id.toLowerCase().includes(query)

        return statusOk && typeOk && searchOk
      })
    },

    isConnected: (state): boolean => state.connectionStatus === 'connected',
  },

  actions: {
    async fetchDevices() {
      const api = useDeviceApi()
      this.loading = true
      this.error = null

      try {
        const previous = new Map(this.devices.map(d => [d.id, d.status]))
        const next = await api.listDevices()

        this.detectOfflineTransitions(previous, next)
        this.devices = next
        this.lastFetchedAt = new Date().toISOString()
        this.connectionStatus = 'connected'
      }
      catch (err) {
        this.connectionStatus = 'disconnected'
        this.error = err instanceof DeviceApiError
          ? err.message
          : 'Failed to load devices'
      }
      finally {
        this.loading = false
      }
    },

    async refreshConnection() {
      const api = useDeviceApi()
      this.connectionStatus = 'checking'
      const ok = await api.checkHealth()
      this.connectionStatus = ok ? 'connected' : 'disconnected'
      return ok
    },

    detectOfflineTransitions(previous: Map<string, DeviceStatus>, next: Device[]) {
      for (const device of next) {
        const prior = previous.get(device.id)
        if (prior === 'ONLINE' && device.status === 'OFFLINE') {
          this.pushAlert(device)
        }
      }
    },

    pushAlert(device: Device) {
      const alert: DeviceAlert = {
        id: `${device.id}-${Date.now()}`,
        deviceId: device.id,
        deviceName: device.name,
        message: `${device.name} went OFFLINE`,
        createdAt: new Date().toISOString(),
        dismissed: false,
      }
      this.alerts = [alert, ...this.alerts].slice(0, 20)
    },

    dismissAlert(alertId: string) {
      const target = this.alerts.find(a => a.id === alertId)
      if (target) {
        target.dismissed = true
      }
    },

    dismissAllAlerts() {
      this.alerts.forEach((a) => {
        a.dismissed = true
      })
    },

    setFilterStatus(status: DeviceStatus | 'ALL') {
      this.filterStatus = status
    },

    setFilterType(type: DeviceType | 'ALL') {
      this.filterType = type
    },

    setSearchQuery(query: string) {
      this.searchQuery = query
    },

    clearFilters() {
      this.filterStatus = 'ALL'
      this.filterType = 'ALL'
      this.searchQuery = ''
    },
  },
})
