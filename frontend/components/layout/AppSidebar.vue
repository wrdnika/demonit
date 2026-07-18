<script setup lang="ts">
const route = useRoute()

const links = [
  { to: '/', label: 'Overview', exact: true },
  { to: '/devices', label: 'Devices', exact: false },
] as const

function isActive(to: string, exact: boolean) {
  if (exact) {
    return route.path === to
  }
  return route.path === to || route.path.startsWith(`${to}/`)
}
</script>

<template>
  <aside
    class="flex h-full w-60 shrink-0 flex-col border-r border-surface-200 bg-surface-950 text-surface-50"
    aria-label="Primary"
  >
    <div class="border-b border-white/10 px-5 py-5">
      <p class="font-mono text-xs uppercase tracking-widest text-accent">
        Demonit
      </p>
      <h1 class="mt-1 text-lg font-semibold tracking-tight">
        IoT Monitor
      </h1>
    </div>

    <nav class="flex-1 px-3 py-4" aria-label="Main navigation">
      <ul class="space-y-1">
        <li v-for="link in links" :key="link.to">
          <NuxtLink
            :to="link.to"
            class="block rounded-md px-3 py-2 text-sm font-medium transition"
            :class="isActive(link.to, link.exact)
              ? 'bg-white/10 text-white'
              : 'text-white/65 hover:bg-white/5 hover:text-white'"
            :aria-current="isActive(link.to, link.exact) ? 'page' : undefined"
          >
            {{ link.label }}
          </NuxtLink>
        </li>
      </ul>
    </nav>

    <div class="border-t border-white/10 px-5 py-4 text-xs text-white/40">
      Hexagonal Go API · Vue 3
    </div>
  </aside>
</template>
