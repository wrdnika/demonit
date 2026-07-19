<script setup lang="ts">
import type { Device, DeviceType } from '~/types/device'
import { DeviceApiError } from '~/composables/useDeviceApi'

const props = defineProps<{
  open: boolean
  mode: 'create' | 'edit'
  device?: Device | null
}>()

const emit = defineEmits<{
  close: []
  saved: [device: Device]
}>()

const store = useDeviceStore()

const name = ref('')
const type = ref<DeviceType>('ATM')
const submitting = ref(false)
const formError = ref<string | null>(null)

const typeOptions: DeviceType[] = ['ATM', 'SERVER', 'LAPTOP']

const title = computed(() => (props.mode === 'edit' ? 'Edit device' : 'Add device'))
const submitLabel = computed(() => {
  if (submitting.value) {
    return props.mode === 'edit' ? 'Saving…' : 'Adding…'
  }
  return props.mode === 'edit' ? 'Save changes' : 'Add device'
})

watch(
  () => [props.open, props.mode, props.device] as const,
  ([open]) => {
    if (!open) {
      return
    }
    formError.value = null
    if (props.mode === 'edit' && props.device) {
      name.value = props.device.name
      type.value = props.device.type
    }
    else {
      name.value = ''
      type.value = 'ATM'
    }
  },
)

function onClose() {
  if (submitting.value) {
    return
  }
  emit('close')
}

async function onSubmit() {
  formError.value = null
  submitting.value = true

  try {
    const payload = {
      name: name.value.trim(),
      type: type.value,
    }
    const device = props.mode === 'edit' && props.device
      ? await store.updateDevice(props.device.id, payload)
      : await store.registerDevice(payload)
    emit('saved', device)
    emit('close')
  }
  catch (err) {
    formError.value = err instanceof DeviceApiError
      ? err.message
      : props.mode === 'edit'
        ? 'Failed to update device'
        : 'Failed to register device'
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
        :aria-label="title"
        @click.self="onClose"
      >
        <div
          class="w-full max-w-md animate-slide-down rounded-brutal border-3 border-ink bg-surface shadow-brutal"
          @click.stop
        >
          <header class="flex items-center justify-between gap-3 border-b-3 border-ink px-5 py-4">
            <h2 class="text-lg font-extrabold tracking-tight">
              {{ title }}
            </h2>
            <button
              type="button"
              class="rounded-full border-2 border-ink px-2.5 py-1 text-sm font-bold"
              aria-label="Close"
              :disabled="submitting"
              @click="onClose"
            >
              ✕
            </button>
          </header>

          <form class="space-y-4 p-5" @submit.prevent="onSubmit">
            <label class="block space-y-1.5">
              <span class="text-xs font-bold uppercase tracking-wide text-ink/70">Name</span>
              <input
                v-model="name"
                type="text"
                required
                minlength="2"
                maxlength="255"
                placeholder="ATM Cabang 02"
                class="brutal-input"
                autofocus
              >
            </label>

            <label class="block space-y-1.5">
              <span class="text-xs font-bold uppercase tracking-wide text-ink/70">Type</span>
              <select v-model="type" class="brutal-input" required>
                <option v-for="opt in typeOptions" :key="opt" :value="opt">
                  {{ opt }}
                </option>
              </select>
            </label>

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
              <button type="submit" class="brutal-btn" :disabled="submitting">
                {{ submitLabel }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
