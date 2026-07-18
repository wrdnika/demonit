/** Shared TypeScript contracts mirroring the Go backend domain. */

export type DeviceType = 'ATM' | 'SERVER' | 'LAPTOP'
export type DeviceStatus = 'ONLINE' | 'OFFLINE'

export interface Device {
  id: string
  name: string
  type: DeviceType
  status: DeviceStatus
  last_seen: string
  created_at: string
  updated_at: string
  /** Present when the API joins the latest metric sample. */
  cpu_usage?: number
  ram_usage?: number
  /** Dynamic hardware fields from latest heartbeat (cash_remaining, temperature, …). */
  status_payload?: Record<string, unknown> | string
}

export interface MetricLog {
  id: string
  device_id: string
  cpu_usage: number
  ram_usage: number
  status_payload: Record<string, unknown> | string
  timestamp: string
}

export interface HeartbeatPayload {
  device_id: string
  cpu_usage: number
  ram_usage: number
  status_payload: Record<string, unknown>
}

export interface ApiSuccessResponse<T> {
  success: true
  message?: string
  data: T
}

export interface ApiErrorBody {
  code: string
  message: string
  details?: Record<string, string>
}

export interface ApiErrorResponse {
  success: false
  error: ApiErrorBody
}

export type ApiResponse<T> = ApiSuccessResponse<T> | ApiErrorResponse

export interface DeviceAlert {
  id: string
  deviceId: string
  deviceName: string
  message: string
  createdAt: string
  dismissed: boolean
}

export type ConnectionStatus = 'connected' | 'disconnected' | 'checking'
