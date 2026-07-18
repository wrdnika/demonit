import type { DeviceStatus } from '~/types/device'

export interface DeviceOfflineSSEPayload {
  type: 'device_offline'
  data: {
    device_id: string
    device_name: string
    device_type: string
    last_seen: string
    occurred_at: string
  }
}

/**
 * Subscribe to backend SSE offline alerts.
 * Returns a cleanup function that closes the EventSource.
 */
export function useAlertStream(onOffline: (payload: DeviceOfflineSSEPayload['data']) => void) {
  const config = useRuntimeConfig()
  let source: EventSource | null = null

  function connect() {
    if (!import.meta.client) {
      return
    }

    const url = `${config.public.apiBase}/api/v1/events`
    source = new EventSource(url)

    source.addEventListener('device_offline', (evt) => {
      try {
        const parsed = JSON.parse((evt as MessageEvent).data) as DeviceOfflineSSEPayload
        if (parsed?.data?.device_id) {
          onOffline(parsed.data)
        }
      }
      catch {
        // ignore malformed frames
      }
    })

    source.onerror = () => {
      // Browser auto-reconnects; leave EventSource open.
    }
  }

  function disconnect() {
    source?.close()
    source = null
  }

  return { connect, disconnect }
}

export type OfflineStatusPatch = {
  id: string
  status: DeviceStatus
}
