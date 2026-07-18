<script setup lang="ts">
import type { ConnectionStatus } from '~/types/device'

defineProps<{
  connectionStatus: ConnectionStatus
  lastFetchedAt: string | null
  onlineCount: number
  offlineCount: number
}>()

const statusLabel = computed(() => ({
  connected: 'API connected',
  disconnected: 'API disconnected',
  checking: 'Checking connection…',
} as const))

const statusDotClass = computed(() => ({
  connected: 'bg-online',
  disconnected: 'bg-offline',
  checking: 'bg-amber-400 animate-pulse-dot',
} as const))
</script>

<template>
  <header class="flex h-14 items-center justify-between gap-4 border-b border-surface-200 bg-white px-6">
    <div>
      <p class="text-sm font-semibold text-surface-900">
        <slot name="title">
          Dashboard
        </slot>
      </p>
      <p v-if="lastFetchedAt" class="text-xs text-surface-800/50">
        Updated
        <time :datetime="lastFetchedAt">{{ new Date(lastFetchedAt).toLocaleTimeString() }}</time>
      </p>
    </div>

    <div class="flex items-center gap-4">
      <div class="hidden items-center gap-3 text-xs sm:flex" aria-label="Device summary">
        <span class="rounded bg-online-soft px-2 py-1 font-medium text-online">
          {{ onlineCount }} online
        </span>
        <span class="rounded bg-offline-soft px-2 py-1 font-medium text-offline">
          {{ offlineCount }} offline
        </span>
      </div>

      <div
        class="inline-flex items-center gap-2 rounded-md border border-surface-200 bg-surface-50 px-3 py-1.5 text-xs font-medium"
        role="status"
        :aria-label="statusLabel[connectionStatus]"
      >
        <span
          class="size-2 rounded-full"
          :class="statusDotClass[connectionStatus]"
          aria-hidden="true"
        />
        <span class="text-surface-800">{{ statusLabel[connectionStatus] }}</span>
      </div>
    </div>
  </header>
</template>
