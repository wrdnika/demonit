<script setup lang="ts">
import type { Device, DeviceStatus, DeviceType } from '~/types/device'

definePageMeta({
  layout: 'default',
})

const store = useDeviceStore()

useHead({
  title: 'Devices · Demonit',
})

const statusOptions: Array<DeviceStatus | 'ALL'> = ['ALL', 'ONLINE', 'OFFLINE']
const typeOptions: Array<DeviceType | 'ALL'> = ['ALL', 'ATM', 'SERVER', 'LAPTOP']

const search = computed({
  get: () => store.searchQuery,
  set: (value: string) => store.setSearchQuery(value),
})

const formOpen = ref(false)
const deleteOpen = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const activeDevice = ref<Device | null>(null)

function openCreate() {
  formMode.value = 'create'
  activeDevice.value = null
  formOpen.value = true
}

function openEdit(device: Device) {
  formMode.value = 'edit'
  activeDevice.value = device
  formOpen.value = true
}

function openDelete(device: Device) {
  activeDevice.value = device
  deleteOpen.value = true
}
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-6 animate-fade-in">
    <header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
      <div class="space-y-1">
        <h2 class="text-3xl font-bold tracking-tight">
          Devices
        </h2>
        <p class="text-sm font-medium text-ink/70">
          Filter the fleet. Click a row for metrics history.
        </p>
      </div>
      <div class="flex flex-wrap gap-2 self-start">
        <button
          type="button"
          class="brutal-btn-ghost"
          :disabled="store.loading"
          @click="store.fetchDevices()"
        >
          {{ store.loading ? 'Refreshing…' : 'Refresh' }}
        </button>
        <button type="button" class="brutal-btn" @click="openCreate">
          Add device
        </button>
      </div>
    </header>

    <BaseCard>
      <form
        class="grid gap-3 md:grid-cols-4"
        role="search"
        aria-label="Filter devices"
        @submit.prevent
      >
        <label class="block space-y-1 md:col-span-2">
          <span class="text-xs font-bold uppercase tracking-wide">Search</span>
          <input
            v-model="search"
            type="search"
            placeholder="Name or UUID"
            class="brutal-input"
          >
        </label>

        <label class="block space-y-1">
          <span class="text-xs font-bold uppercase tracking-wide">Status</span>
          <select
            class="brutal-input"
            :value="store.filterStatus"
            @change="store.setFilterStatus(($event.target as HTMLSelectElement).value as DeviceStatus | 'ALL')"
          >
            <option v-for="opt in statusOptions" :key="opt" :value="opt">
              {{ opt }}
            </option>
          </select>
        </label>

        <label class="block space-y-1">
          <span class="text-xs font-bold uppercase tracking-wide">Type</span>
          <select
            class="brutal-input"
            :value="store.filterType"
            @change="store.setFilterType(($event.target as HTMLSelectElement).value as DeviceType | 'ALL')"
          >
            <option v-for="opt in typeOptions" :key="opt" :value="opt">
              {{ opt }}
            </option>
          </select>
        </label>
      </form>

      <div class="mt-3 flex items-center justify-between text-xs font-bold uppercase tracking-wide">
        <p>
          Showing {{ store.filteredDevices.length }} / {{ store.totalCount }}
        </p>
        <button type="button" class="underline" @click="store.clearFilters()">
          Clear
        </button>
      </div>
    </BaseCard>

    <section aria-label="Device list">
      <div
        v-if="store.filteredDevices.length === 0"
        class="rounded-brutal border-3 border-ink bg-surface px-5 py-10 text-center text-sm font-bold shadow-brutal lg:hidden"
      >
        No devices match the current filters.
      </div>

      <div class="grid gap-4 lg:hidden">
        <DeviceCard
          v-for="device in store.filteredDevices"
          :key="device.id"
          :device="device"
          @edit="openEdit"
          @delete="openDelete"
        />
      </div>

      <BaseCard :padding="false" class="hidden overflow-hidden lg:block">
        <div class="overflow-x-auto">
          <table class="min-w-full text-left">
            <caption class="sr-only">
              Monitored IoT devices
            </caption>
            <thead class="border-b-3 border-ink bg-gyellow text-xs font-bold uppercase tracking-wide text-[#202124]">
              <tr>
                <th scope="col" class="px-4 py-3">
                  Device
                </th>
                <th scope="col" class="px-4 py-3">
                  Type
                </th>
                <th scope="col" class="px-4 py-3">
                  Status
                </th>
                <th scope="col" class="px-4 py-3">
                  CPU
                </th>
                <th scope="col" class="px-4 py-3">
                  RAM
                </th>
                <th scope="col" class="px-4 py-3">
                  Last seen
                </th>
                <th scope="col" class="px-4 py-3">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              <DeviceRow
                v-for="device in store.filteredDevices"
                :key="device.id"
                :device="device"
                @edit="openEdit"
                @delete="openDelete"
              />
            </tbody>
          </table>
        </div>

        <p
          v-if="store.filteredDevices.length === 0"
          class="px-5 py-10 text-center text-sm font-bold"
        >
          No devices match the current filters.
        </p>
      </BaseCard>
    </section>

    <DeviceFormModal
      :open="formOpen"
      :mode="formMode"
      :device="activeDevice"
      @close="formOpen = false"
    />
    <DeleteDeviceModal
      :open="deleteOpen"
      :device="activeDevice"
      @close="deleteOpen = false"
    />
  </div>
</template>
