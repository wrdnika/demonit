<script setup lang="ts">
const props = withDefaults(defineProps<{
  label: string
  value: number | null | undefined
  max?: number
}>(), {
  max: 100,
})

const pct = computed(() => {
  if (props.value == null || Number.isNaN(props.value)) {
    return null
  }
  return Math.min(Math.max(props.value, 0), props.max)
})

const fillClass = computed(() => {
  if (pct.value == null) {
    return 'bg-surface-200'
  }
  if (pct.value >= 85) {
    return 'bg-alert'
  }
  if (pct.value >= 65) {
    return 'bg-amber-500'
  }
  return 'bg-accent'
})

const display = computed(() =>
  pct.value == null ? '—' : `${pct.value.toFixed(1)}%`,
)
</script>

<template>
  <div class="space-y-1.5">
    <div class="flex items-center justify-between text-xs">
      <span class="font-medium text-surface-800/70">{{ label }}</span>
      <span class="font-mono text-surface-900">{{ display }}</span>
    </div>
    <div
      class="metric-track"
      role="progressbar"
      :aria-valuenow="pct ?? undefined"
      :aria-valuemin="0"
      :aria-valuemax="max"
      :aria-label="`${label} usage`"
    >
      <div
        class="metric-fill"
        :class="fillClass"
        :style="{ width: pct == null ? '0%' : `${pct}%` }"
      />
    </div>
  </div>
</template>
