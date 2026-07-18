<script setup lang="ts">
import type { ConnectionStatus } from '~/types/device'

defineProps<{
  connectionStatus: ConnectionStatus
  lastFetchedAt: string | null
  onlineCount: number
  offlineCount: number
}>()

const emit = defineEmits<{
  toggleMenu: []
}>()

const statusLabel = computed(() => ({
  connected: 'Live',
  disconnected: 'Down',
  checking: 'Checking',
} as const))

const statusTone = computed(() => ({
  connected: 'bg-ggreen text-white',
  disconnected: 'bg-gred text-white',
  checking: 'bg-gyellow text-ink',
} as const))
</script>

<template>
  <header class="sticky top-0 z-20 flex h-[4.75rem] items-center justify-between gap-3 border-b-3 border-ink bg-white px-4 sm:px-8">
    <div class="flex min-w-0 items-center gap-3">
      <button
        type="button"
        class="inline-flex size-10 items-center justify-center rounded-brutal border-3 border-ink bg-white text-lg font-bold shadow-brutal-sm lg:hidden"
        aria-label="Open navigation menu"
        @click="emit('toggleMenu')"
      >
        ☰
      </button>
      <div class="min-w-0">
        <p class="truncate text-xl font-extrabold tracking-tight">
          <slot name="title">
            Dashboard
          </slot>
        </p>
        <p v-if="lastFetchedAt" class="text-xs font-medium text-ink/45">
          Synced
          <time :datetime="lastFetchedAt">{{ new Date(lastFetchedAt).toLocaleTimeString() }}</time>
        </p>
      </div>
    </div>

    <div class="flex shrink-0 items-center gap-2 sm:gap-3">
      <div class="hidden items-center gap-2 text-xs font-bold sm:flex" aria-label="Device summary">
        <span class="rounded-full border-2 border-ink bg-ggreen px-3 py-1 text-white shadow-brutal-sm">
          {{ onlineCount }} online
        </span>
        <span class="rounded-full border-2 border-ink bg-gred px-3 py-1 text-white shadow-brutal-sm">
          {{ offlineCount }} offline
        </span>
      </div>

      <div
        class="inline-flex items-center gap-2 rounded-full border-2 border-ink px-3 py-1 text-xs font-bold shadow-brutal-sm"
        :class="statusTone[connectionStatus]"
        role="status"
        :aria-label="statusLabel[connectionStatus]"
      >
        <span
          class="size-2 rounded-full bg-current"
          :class="connectionStatus === 'connected' ? 'animate-pulse-dot' : ''"
          aria-hidden="true"
        />
        <span class="hidden xs:inline sm:inline">API</span>
        {{ statusLabel[connectionStatus] }}
      </div>
    </div>
  </header>
</template>
