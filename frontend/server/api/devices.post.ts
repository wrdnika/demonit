import type { ApiErrorResponse, ApiResponse, ApiSuccessResponse, Device } from '~/types/device'

interface RegisterBody {
  name: string
  type: Device['type']
}

/**
 * BFF proxy: browser → Nuxt server → Go API.
 * Admin API key stays server-side (runtimeConfig.adminApiKey).
 */
export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const body = await readBody<RegisterBody>(event)

  if (!body?.name || !body?.type) {
    throw createError({
      statusCode: 400,
      statusMessage: 'name and type are required',
    })
  }

  try {
    const response = await $fetch<ApiResponse<Device>>('/api/v1/devices', {
      baseURL: config.apiBaseServer as string,
      method: 'POST',
      body,
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
        'X-Admin-API-Key': config.adminApiKey as string,
      },
    })

    if (!response || typeof response !== 'object') {
      throw createError({ statusCode: 502, statusMessage: 'Empty upstream response' })
    }

    if ('success' in response && response.success === false) {
      const err = response as ApiErrorResponse
      throw createError({
        statusCode: 400,
        statusMessage: err.error.message,
        data: err.error,
      })
    }

    return response as ApiSuccessResponse<Device>
  }
  catch (error: unknown) {
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
      statusMessage: fetchError.statusMessage || fetchError.message || 'Upstream register failed',
    })
  }
})
