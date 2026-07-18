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
        class="flex items-start gap-3 rounded-brutal border-3 border-ink bg-gred px-4 py-3 text-white shadow-brutal-red"
        role="alert"
      >
        <span
          class="mt-0.5 inline-flex size-7 shrink-0 items-center justify-center rounded-full border-2 border-ink bg-white text-xs font-bold text-gred"
          aria-hidden="true"
        >
          !
        </span>

        <div class="min-w-0 flex-1">
          <p class="text-sm font-bold">
            Device offline
          </p>
          <p class="truncate text-sm opacity-90">
            {{ alert.message }}
          </p>
        </div>

        <button
          type="button"
          class="rounded-full border-2 border-ink bg-white px-2.5 py-1 text-xs font-bold text-ink"
          :aria-label="`Dismiss alert for ${alert.deviceName}`"
          @click="emit('dismiss', alert.id)"
        >
          ✕
        </button>
      </aside>
    </TransitionGroup>

    <button
      v-if="alerts.length > 1"
      type="button"
      class="brutal-btn-ghost pointer-events-auto text-xs"
      @click="emit('dismissAll')"
    >
      Dismiss all
    </button>
  </div>
</template>
