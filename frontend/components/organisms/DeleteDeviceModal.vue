<script setup lang="ts">
import type { Device } from '~/types/device'
import { DeviceApiError } from '~/composables/useDeviceApi'

const props = defineProps<{
  open: boolean
  device: Device | null
}>()

const emit = defineEmits<{
  close: []
  deleted: [id: string]
}>()

const store = useDeviceStore()
const submitting = ref(false)
const formError = ref<string | null>(null)

watch(
  () => props.open,
  (open) => {
    if (open) {
      formError.value = null
    }
  },
)

function onClose() {
  if (submitting.value) {
    return
  }
  emit('close')
}

async function onConfirm() {
  if (!props.device) {
    return
  }
  formError.value = null
  submitting.value = true
  try {
    await store.deleteDevice(props.device.id)
    emit('deleted', props.device.id)
    emit('close')
  }
  catch (err) {
    formError.value = err instanceof DeviceApiError
      ? err.message
      : 'Failed to delete device'
  }
  finally {
    submitting.value = false
  }
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.open) {
    onClose()
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200"
      enter-from-class="opacity-0"
      leave-active-class="transition duration-150"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        class="fixed inset-0 z-[60] flex items-end justify-center bg-ink/40 p-4 sm:items-center"
        role="dialog"
        aria-modal="true"
        aria-label="Delete device"
        @click.self="onClose"
      >
        <div
          class="w-full max-w-md animate-slide-down rounded-brutal border-3 border-ink bg-surface shadow-brutal"
          @click.stop
        >
          <header class="border-b-3 border-ink px-5 py-4">
            <h2 class="text-lg font-extrabold tracking-tight">
              Delete device?
            </h2>
          </header>

          <div class="space-y-4 p-5">
            <p class="text-sm font-medium text-ink/80">
              This removes
              <span class="font-bold text-ink">{{ device?.name }}</span>
              and its metric history. Agents using this ID will stop matching.
            </p>

            <p
              v-if="formError"
              class="rounded-brutal border-2 border-ink bg-gred px-3 py-2 text-xs font-bold text-white"
              role="alert"
            >
              {{ formError }}
            </p>

            <div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
              <button
                type="button"
                class="brutal-btn-ghost"
                :disabled="submitting"
                @click="onClose"
              >
                Cancel
              </button>
              <button
                type="button"
                class="brutal-btn-danger"
                :disabled="submitting"
                @click="onConfirm"
              >
                {{ submitting ? 'Deleting…' : 'Delete' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
