import type { ApiErrorResponse, ApiResponse, ApiSuccessResponse, Device } from '~/types/device'

interface DeviceBody {
  name: string
  type: Device['type']
}

function adminHeaders(adminApiKey: string) {
  return {
    Accept: 'application/json',
    'Content-Type': 'application/json',
    'X-Admin-API-Key': adminApiKey,
  }
}

function mapUpstreamError(error: unknown, fallback: string) {
  const fetchError = error as {
    statusCode?: number
    statusMessage?: string
    data?: ApiErrorResponse
    message?: string
  }

  if (fetchError.data?.error) {
    throw createError({
      statusCode: fetchError.statusCode || 401,
      statusMessage: fetchError.data.error.message,
      data: fetchError.data.error,
    })
  }

  throw createError({
    statusCode: fetchError.statusCode || 502,
    statusMessage: fetchError.statusMessage || fetchError.message || fallback,
  })
}

export async function proxyCreateDevice(body: DeviceBody) {
  const config = useRuntimeConfig()
  try {
    const response = await $fetch<ApiResponse<Device>>('/api/v1/devices', {
      baseURL: config.apiBaseServer as string,
      method: 'POST',
      body,
      headers: adminHeaders(config.adminApiKey as string),
    })
    if (!response || typeof response !== 'object') {
      throw createError({ statusCode: 502, statusMessage: 'Empty upstream response' })
    }
    if ('success' in response && response.success === false) {
      const err = response as ApiErrorResponse
      throw createError({ statusCode: 400, statusMessage: err.error.message, data: err.error })
    }
    return response as ApiSuccessResponse<Device>
  }
  catch (error: unknown) {
    mapUpstreamError(error, 'Upstream register failed')
  }
}

export async function proxyUpdateDevice(id: string, body: DeviceBody) {
  const config = useRuntimeConfig()
  try {
    const response = await $fetch<ApiResponse<Device>>(`/api/v1/devices/${id}`, {
      baseURL: config.apiBaseServer as string,
      method: 'PUT',
      body,
      headers: adminHeaders(config.adminApiKey as string),
    })
    if (!response || typeof response !== 'object') {
      throw createError({ statusCode: 502, statusMessage: 'Empty upstream response' })
    }
    if ('success' in response && response.success === false) {
      const err = response as ApiErrorResponse
      throw createError({ statusCode: 400, statusMessage: err.error.message, data: err.error })
    }
    return response as ApiSuccessResponse<Device>
  }
  catch (error: unknown) {
    mapUpstreamError(error, 'Upstream update failed')
  }
}

export async function proxyDeleteDevice(id: string) {
  const config = useRuntimeConfig()
  try {
    const response = await $fetch<ApiResponse<{ id: string }>>(`/api/v1/devices/${id}`, {
      baseURL: config.apiBaseServer as string,
      method: 'DELETE',
      headers: adminHeaders(config.adminApiKey as string),
    })
    if (!response || typeof response !== 'object') {
      throw createError({ statusCode: 502, statusMessage: 'Empty upstream response' })
    }
    if ('success' in response && response.success === false) {
      const err = response as ApiErrorResponse
      throw createError({ statusCode: 400, statusMessage: err.error.message, data: err.error })
    }
    return response as ApiSuccessResponse<{ id: string }>
  }
  catch (error: unknown) {
    mapUpstreamError(error, 'Upstream delete failed')
  }
}
