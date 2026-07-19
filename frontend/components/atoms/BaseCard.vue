<script setup lang="ts">
withDefaults(defineProps<{
  padding?: boolean
  tone?: 'white' | 'lime' | 'cyan' | 'sun' | 'pink'
}>(), {
  padding: true,
  tone: 'white',
})

const toneClass: Record<string, string> = {
  white: 'bg-surface shadow-brutal',
  lime: 'bg-ggreen text-white shadow-brutal-green',
  cyan: 'bg-gblue text-white shadow-brutal-blue',
  sun: 'bg-gyellow text-[#202124] shadow-brutal-yellow',
  pink: 'bg-gred text-white shadow-brutal-red',
}
</script>

<template>
  <article class="rounded-brutal border-3 border-ink" :class="toneClass[tone]">
    <header
      v-if="$slots.header"
      class="flex items-center justify-between gap-3 border-b-3 border-ink px-5 py-3.5"
    >
      <slot name="header" />
    </header>

    <div :class="padding ? 'p-5' : ''">
      <slot />
    </div>

    <footer
      v-if="$slots.footer"
      class="border-t-3 border-ink px-5 py-3.5"
    >
      <slot name="footer" />
    </footer>
  </article>
</template>
