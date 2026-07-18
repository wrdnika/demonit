<script setup lang="ts">
withDefaults(defineProps<{
  padding?: boolean
}>(), {
  padding: true,
})
</script>

<template>
  <article class="rounded-lg border border-surface-200 bg-white shadow-panel">
    <header
      v-if="$slots.header"
      class="flex items-center justify-between gap-3 border-b border-surface-100 px-5 py-3"
    >
      <slot name="header" />
    </header>

    <div :class="padding ? 'p-5' : ''">
      <slot />
    </div>

    <footer
      v-if="$slots.footer"
      class="border-t border-surface-100 px-5 py-3"
    >
      <slot name="footer" />
    </footer>
  </article>
</template>
