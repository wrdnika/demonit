import type { Device } from '~/types/device'
import { proxyUpdateDevice } from '../../utils/deviceProxy'

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, 'id')
  if (!id) {
    throw createError({ statusCode: 400, statusMessage: 'id is required' })
  }
  const body = await readBody<{ name: string, type: Device['type'] }>(event)
  if (!body?.name || !body?.type) {
    throw createError({ statusCode: 400, statusMessage: 'name and type are required' })
  }
  return proxyUpdateDevice(id, body)
})
