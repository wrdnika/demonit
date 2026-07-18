<script setup lang="ts">
import type { DeviceAlert } from '~/types/device'

defineProps<{
  alerts: DeviceAlert[]
}>()

const emit = defineEmits<{
  dismiss: [id: string]
  dismissAll: []
}>()
</script>

<template>
  <div
    v-if="alerts.length"
    class="pointer-events-none fixed inset-x-0 top-0 z-50 flex flex-col items-center gap-2 p-4 sm:items-end"
    aria-live="assertive"
    aria-relevant="additions"
  >
    <TransitionGroup
      enter-active-class="animate-slide-down"
      leave-active-class="transition duration-200 ease-in opacity-0 -translate-y-2"
      move-class="transition duration-200"
      tag="div"
      class="pointer-events-auto flex w-full max-w-md flex-col gap-2"
    >
      <aside
        v-for="alert in alerts"
        :key="alert.id"
        class="flex items-start gap-3 rounded-lg border border-alert/20 bg-alert-soft px-4 py-3 shadow-panel"
        role="alert"
      >
        <span
          class="mt-0.5 inline-flex size-6 shrink-0 items-center justify-center rounded bg-alert text-xs font-bold text-white"
          aria-hidden="true"
        >
          !
        </span>

        <div class="min-w-0 flex-1">
          <p class="text-sm font-semibold text-alert">
            Device offline
          </p>
          <p class="truncate text-sm text-surface-900">
            {{ alert.message }}
          </p>
        </div>

        <button
          type="button"
          class="rounded px-2 py-1 text-xs font-medium text-alert transition hover:bg-white/60"
          :aria-label="`Dismiss alert for ${alert.deviceName}`"
          @click="emit('dismiss', alert.id)"
        >
          Dismiss
        </button>
      </aside>
    </TransitionGroup>

    <button
      v-if="alerts.length > 1"
      type="button"
      class="pointer-events-auto rounded bg-white px-3 py-1.5 text-xs font-medium text-surface-800 shadow-panel transition hover:bg-surface-100"
      @click="emit('dismissAll')"
    >
      Dismiss all
    </button>
  </div>
</template>
