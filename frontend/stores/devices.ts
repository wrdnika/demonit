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
  /** Monotonic counter so overlapping polls cannot apply stale results. */
  fetchSeq: number
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
    fetchSeq: 0,
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
      const seq = ++this.fetchSeq
      this.loading = true
      this.error = null

      try {
        const previous = new Map(this.devices.map(d => [d.id, d.status]))
        const next = await api.listDevices()

        // Drop stale responses if a newer poll already started.
        if (seq !== this.fetchSeq) {
          return
        }

        this.detectOfflineTransitions(previous, next)
        this.devices = next
        this.lastFetchedAt = new Date().toISOString()
        this.connectionStatus = 'connected'
      }
      catch (err) {
        if (seq !== this.fetchSeq) {
          return
        }
        this.connectionStatus = 'disconnected'
        this.error = err instanceof DeviceApiError
          ? err.message
          : 'Failed to load devices'
      }
      finally {
        if (seq === this.fetchSeq) {
          this.loading = false
        }
      }
    },

    async registerDevice(payload: { name: string, type: DeviceType }) {
      const api = useDeviceApi()
      const device = await api.registerDevice(payload)
      this.devices = [device, ...this.devices]
      this.connectionStatus = 'connected'
      return device
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

    /** Apply realtime SSE offline event without waiting for the next poll. */
    applyOfflineEvent(payload: {
      device_id: string
      device_name: string
      last_seen?: string
    }) {
      const existing = this.devices.find(d => d.id === payload.device_id)
      if (existing) {
        if (existing.status === 'OFFLINE') {
          return
        }
        existing.status = 'OFFLINE'
        if (payload.last_seen) {
          existing.last_seen = payload.last_seen
        }
        this.pushAlert(existing)
        return
      }

      this.pushAlert({
        id: payload.device_id,
        name: payload.device_name,
        type: 'ATM',
        status: 'OFFLINE',
        last_seen: payload.last_seen || new Date().toISOString(),
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      })
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
