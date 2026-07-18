<script setup lang="ts">
const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const route = useRoute()

const links = [
  { to: '/', label: 'Overview', color: 'bg-gblue', exact: true },
  { to: '/devices', label: 'Devices', color: 'bg-ggreen', exact: false },
] as const

function isActive(to: string, exact: boolean) {
  if (exact) {
    return route.path === to
  }
  return route.path === to || route.path.startsWith(`${to}/`)
}

function onNavigate() {
  emit('close')
}

watch(() => route.fullPath, () => {
  if (props.open) {
    emit('close')
  }
})
</script>

<template>
  <!-- Mobile backdrop -->
  <Transition
    enter-active-class="transition duration-200"
    enter-from-class="opacity-0"
    leave-active-class="transition duration-150"
    leave-to-class="opacity-0"
  >
    <button
      v-if="open"
      type="button"
      class="fixed inset-0 z-40 bg-ink/40 lg:hidden"
      aria-label="Close navigation"
      @click="emit('close')"
    />
  </Transition>

  <aside
    class="fixed inset-y-0 left-0 z-50 flex h-screen w-64 flex-col border-r-3 border-ink bg-white transition-transform duration-200 lg:sticky lg:translate-x-0"
    :class="open ? 'translate-x-0' : '-translate-x-full'"
    aria-label="Primary"
  >
    <div class="flex items-center justify-between border-b-3 border-ink p-5">
      <div class="flex items-center gap-3">
        <span
          class="flex size-11 shrink-0 items-center justify-center rounded-2xl border-3 border-ink bg-gblue font-extrabold text-white shadow-brutal-sm"
          aria-hidden="true"
        >
          D
        </span>
        <div class="min-w-0">
          <p class="text-[11px] font-semibold uppercase tracking-[0.14em] text-ink/40">
            Device Monitor
          </p>
          <h1 class="text-2xl font-extrabold leading-none tracking-tight">
            Demonit
          </h1>
        </div>
      </div>
      <button
        type="button"
        class="rounded-full border-2 border-ink px-2.5 py-1 text-sm font-bold lg:hidden"
        aria-label="Close menu"
        @click="emit('close')"
      >
        ✕
      </button>
    </div>

    <nav class="flex-1 space-y-2 p-4" aria-label="Main navigation">
      <NuxtLink
        v-for="link in links"
        :key="link.to"
        :to="link.to"
        class="flex items-center gap-3 rounded-brutal border-3 border-ink px-3 py-3 text-sm font-bold transition"
        :class="isActive(link.to, link.exact)
          ? `${link.color} text-white shadow-brutal-sm`
          : 'bg-white text-ink hover:bg-paper'"
        :aria-current="isActive(link.to, link.exact) ? 'page' : undefined"
        @click="onNavigate"
      >
        <span
          class="size-2.5 rounded-full border-2 border-ink"
          :class="isActive(link.to, link.exact) ? 'bg-white' : link.color"
          aria-hidden="true"
        />
        {{ link.label }}
      </NuxtLink>
    </nav>

    <div class="m-4 rounded-brutal border-3 border-ink bg-paper p-4 shadow-brutal-yellow">
      <p class="text-[11px] font-bold uppercase tracking-wide text-ink/45">
        Fleet API
      </p>
      <p class="mt-1 text-sm font-semibold leading-snug">
        Go · Nuxt · Postgres
      </p>
    </div>
  </aside>
</template>
