<script setup lang="ts">
import type { DeviceType } from '~/types/device'
import { DeviceApiError } from '~/composables/useDeviceApi'

const store = useDeviceStore()

const name = ref('')
const type = ref<DeviceType>('ATM')
const submitting = ref(false)
const formError = ref<string | null>(null)
const successId = ref<string | null>(null)

const typeOptions: DeviceType[] = ['ATM', 'SERVER', 'LAPTOP']

async function onSubmit() {
  formError.value = null
  successId.value = null
  submitting.value = true

  try {
    const device = await store.registerDevice({
      name: name.value.trim(),
      type: type.value,
    })
    successId.value = device.id
    name.value = ''
    type.value = 'ATM'
  }
  catch (err) {
    formError.value = err instanceof DeviceApiError
      ? err.message
      : 'Failed to register device'
  }
  finally {
    submitting.value = false
  }
}
</script>

<template>
  <BaseCard tone="cyan">
    <template #header>
      <h3 class="text-sm font-bold">
        Register device
      </h3>
    </template>

    <form class="space-y-3" @submit.prevent="onSubmit">
      <label class="block space-y-1.5">
        <span class="text-xs font-bold uppercase tracking-wide opacity-80">Name</span>
        <input
          v-model="name"
          type="text"
          required
          minlength="2"
          maxlength="255"
          placeholder="ATM Cabang 02"
          class="brutal-input text-ink"
        >
      </label>

      <label class="block space-y-1.5">
        <span class="text-xs font-bold uppercase tracking-wide opacity-80">Type</span>
        <select v-model="type" class="brutal-input text-ink" required>
          <option v-for="opt in typeOptions" :key="opt" :value="opt">
            {{ opt }}
          </option>
        </select>
      </label>

      <p v-if="formError" class="rounded-brutal border-2 border-ink bg-gred px-3 py-2 text-xs font-bold text-white" role="alert">
        {{ formError }}
      </p>
      <p v-if="successId" class="rounded-brutal border-2 border-ink bg-ggreen px-3 py-2 text-xs font-bold text-white" role="status">
        Registered · {{ successId.slice(0, 8) }}…
      </p>

      <button type="submit" class="brutal-btn-ghost w-full bg-gyellow" :disabled="submitting">
        {{ submitting ? 'Saving…' : 'Add device' }}
      </button>
    </form>
  </BaseCard>
</template>
