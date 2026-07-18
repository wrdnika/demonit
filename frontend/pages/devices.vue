<script setup lang="ts">
import type { DeviceStatus, DeviceType } from '~/types/device'

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
</script>

<template>
  <div class="mx-auto max-w-6xl space-y-6 animate-fade-in">
    <header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
      <div class="space-y-1">
        <h2 class="text-2xl font-semibold tracking-tight text-surface-900">
          Devices
        </h2>
        <p class="text-sm text-surface-800/60">
          Filter the fleet by status, type, or name. Metrics refresh on the polling interval.
        </p>
      </div>
      <button
        type="button"
        class="self-start rounded-md bg-surface-900 px-3 py-2 text-xs font-medium text-white transition hover:bg-surface-800 disabled:opacity-50"
        :disabled="store.loading"
        @click="store.fetchDevices()"
      >
        {{ store.loading ? 'Refreshing…' : 'Refresh now' }}
      </button>
    </header>

    <BaseCard>
      <form
        class="grid gap-3 md:grid-cols-4"
        role="search"
        aria-label="Filter devices"
        @submit.prevent
      >
        <label class="block space-y-1 md:col-span-2">
          <span class="text-xs font-medium text-surface-800/60">Search</span>
          <input
            v-model="search"
            type="search"
            placeholder="Name or UUID"
            class="w-full rounded-md border border-surface-200 bg-white px-3 py-2 text-sm outline-none transition focus:border-accent"
          >
        </label>

        <label class="block space-y-1">
          <span class="text-xs font-medium text-surface-800/60">Status</span>
          <select
            class="w-full rounded-md border border-surface-200 bg-white px-3 py-2 text-sm outline-none focus:border-accent"
            :value="store.filterStatus"
            @change="store.setFilterStatus(($event.target as HTMLSelectElement).value as DeviceStatus | 'ALL')"
          >
            <option v-for="opt in statusOptions" :key="opt" :value="opt">
              {{ opt }}
            </option>
          </select>
        </label>

        <label class="block space-y-1">
          <span class="text-xs font-medium text-surface-800/60">Type</span>
          <select
            class="w-full rounded-md border border-surface-200 bg-white px-3 py-2 text-sm outline-none focus:border-accent"
            :value="store.filterType"
            @change="store.setFilterType(($event.target as HTMLSelectElement).value as DeviceType | 'ALL')"
          >
            <option v-for="opt in typeOptions" :key="opt" :value="opt">
              {{ opt }}
            </option>
          </select>
        </label>
      </form>

      <div class="mt-3 flex items-center justify-between text-xs text-surface-800/50">
        <p>
          Showing {{ store.filteredDevices.length }} of {{ store.totalCount }}
        </p>
        <button
          type="button"
          class="font-medium text-accent-strong hover:text-accent"
          @click="store.clearFilters()"
        >
          Clear filters
        </button>
      </div>
    </BaseCard>

    <section aria-label="Device list">
      <!-- Mobile / tablet cards -->
      <div class="grid gap-4 lg:hidden">
        <TransitionGroup
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="opacity-0 translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
        >
          <DeviceCard
            v-for="device in store.filteredDevices"
            :key="device.id"
            :device="device"
          />
        </TransitionGroup>
      </div>

      <!-- Desktop table -->
      <BaseCard :padding="false" class="hidden overflow-hidden lg:block">
        <div class="overflow-x-auto">
          <table class="min-w-full text-left">
            <caption class="sr-only">
              Monitored IoT devices
            </caption>
            <thead class="bg-surface-50 text-xs uppercase tracking-wide text-surface-800/50">
              <tr>
                <th scope="col" class="px-4 py-3 font-medium">
                  Device
                </th>
                <th scope="col" class="px-4 py-3 font-medium">
                  Type
                </th>
                <th scope="col" class="px-4 py-3 font-medium">
                  Status
                </th>
                <th scope="col" class="px-4 py-3 font-medium">
                  CPU
                </th>
                <th scope="col" class="px-4 py-3 font-medium">
                  RAM
                </th>
                <th scope="col" class="px-4 py-3 font-medium">
                  Last seen
                </th>
              </tr>
            </thead>
            <tbody>
              <DeviceRow
                v-for="device in store.filteredDevices"
                :key="device.id"
                :device="device"
              />
            </tbody>
          </table>
        </div>

        <p
          v-if="store.filteredDevices.length === 0"
          class="px-5 py-10 text-center text-sm text-surface-800/50"
        >
          No devices match the current filters.
        </p>
      </BaseCard>
    </section>
  </div>
</template>
