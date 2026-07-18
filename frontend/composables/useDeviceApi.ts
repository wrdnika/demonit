import type {
  ApiErrorResponse,
  ApiResponse,
  ApiSuccessResponse,
  Device,
  HeartbeatPayload,
  MetricLog,
} from '~/types/device'

export class DeviceApiError extends Error {
  readonly code: string
  readonly statusCode: number
  readonly details?: Record<string, string>

  constructor(message: string, code: string, statusCode: number, details?: Record<string, string>) {
    super(message)
    this.name = 'DeviceApiError'
    this.code = code
    this.statusCode = statusCode
    this.details = details
  }
}

export interface RegisterDevicePayload {
  name: string
  type: Device['type']
}

/**
 * Reusable API client for the Go monitoring backend.
 * Uses runtimeConfig for the base URL — never hardcode hosts.
 */
export function useDeviceApi() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string

  async function request<T>(
    path: string,
    options: Parameters<typeof $fetch<ApiResponse<T>>>[1] = {},
  ): Promise<T> {
    try {
      const response = await $fetch<ApiResponse<T>>(path, {
        baseURL,
        ...options,
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
          ...(options.headers || {}),
        },
        onResponseError({ response }) {
          const body = response._data as ApiErrorResponse | undefined
          const message = body?.error?.message || response.statusText || 'Request failed'
          const code = body?.error?.code || 'HTTP_ERROR'
          throw new DeviceApiError(message, code, response.status, body?.error?.details)
        },
      })

      if (!response || typeof response !== 'object') {
        throw new DeviceApiError('Empty API response', 'EMPTY_RESPONSE', 500)
      }

      if ('success' in response && response.success === false) {
        const err = response as ApiErrorResponse
        throw new DeviceApiError(
          err.error.message,
          err.error.code,
          400,
          err.error.details,
        )
      }

      return (response as ApiSuccessResponse<T>).data
    }
    catch (error) {
      if (error instanceof DeviceApiError) {
        throw error
      }

      const fetchError = error as { statusCode?: number, statusMessage?: string, message?: string }
      throw new DeviceApiError(
        fetchError.statusMessage || fetchError.message || 'Network request failed',
        'NETWORK_ERROR',
        fetchError.statusCode || 0,
      )
    }
  }

  function listDevices() {
    return request<Device[]>('/api/v1/devices', { method: 'GET' })
  }

  function getDevice(id: string) {
    return request<Device>(`/api/v1/devices/${id}`, { method: 'GET' })
  }

  function listDeviceMetrics(id: string, limit = 50) {
    return request<MetricLog[]>(`/api/v1/devices/${id}/metrics?limit=${limit}`, { method: 'GET' })
  }

  async function registerDevice(payload: RegisterDevicePayload) {
    // Nuxt BFF keeps ADMIN_API_KEY on the server — do not call Go directly.
    try {
      const response = await $fetch<ApiResponse<Device>>('/api/devices', {
        method: 'POST',
        body: payload,
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
      })

      if (!response || typeof response !== 'object') {
        throw new DeviceApiError('Empty API response', 'EMPTY_RESPONSE', 500)
      }
      if ('success' in response && response.success === false) {
        const err = response as ApiErrorResponse
        throw new DeviceApiError(err.error.message, err.error.code, 400, err.error.details)
      }
      return (response as ApiSuccessResponse<Device>).data
    }
    catch (error) {
      if (error instanceof DeviceApiError) {
        throw error
      }
      const fetchError = error as {
        statusCode?: number
        statusMessage?: string
        message?: string
        data?: { message?: string, code?: string }
      }
      throw new DeviceApiError(
        fetchError.data?.message || fetchError.statusMessage || fetchError.message || 'Register failed',
        fetchError.data?.code || 'NETWORK_ERROR',
        fetchError.statusCode || 0,
      )
    }
  }

  function sendHeartbeat(payload: HeartbeatPayload) {
    return request<{ device_id: string, status: string }>('/api/v1/heartbeat', {
      method: 'POST',
      body: payload,
    })
  }

  async function checkHealth(): Promise<boolean> {
    try {
      await $fetch('/healthz', {
        baseURL,
        method: 'GET',
        timeout: 4000,
      })
      return true
    }
    catch {
      return false
    }
  }

  return {
    baseURL,
    listDevices,
    getDevice,
    listDeviceMetrics,
    registerDevice,
    sendHeartbeat,
    checkHealth,
  }
}
