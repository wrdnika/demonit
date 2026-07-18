<script setup lang="ts">
const props = withDefaults(defineProps<{
  payload?: Record<string, unknown> | string | null
  limit?: number
}>(), {
  limit: 6,
})

const entries = computed(() => {
  const raw = props.payload
  if (!raw) {
    return [] as Array<{ key: string, value: string }>
  }

  let obj: Record<string, unknown> = {}
  if (typeof raw === 'string') {
    try {
      obj = JSON.parse(raw) as Record<string, unknown>
    }
    catch {
      return []
    }
  }
  else {
    obj = raw
  }

  return Object.entries(obj)
    .slice(0, props.limit)
    .map(([key, value]) => ({
      key: key.replaceAll('_', ' '),
      value: formatValue(value),
    }))
})

function formatValue(value: unknown): string {
  if (value == null) {
    return '—'
  }
  if (typeof value === 'number') {
    return Number.isInteger(value) ? String(value) : value.toFixed(1)
  }
  if (typeof value === 'boolean') {
    return value ? 'yes' : 'no'
  }
  if (typeof value === 'object') {
    return JSON.stringify(value)
  }
  return String(value)
}
</script>

<template>
  <dl v-if="entries.length" class="grid gap-2">
    <div
      v-for="item in entries"
      :key="item.key"
      class="flex items-center justify-between gap-2 rounded-xl border-2 border-ink/15 bg-paper px-2.5 py-1.5 text-xs"
    >
      <dt class="font-bold uppercase tracking-wide text-ink/45">
        {{ item.key }}
      </dt>
      <dd class="truncate font-mono font-semibold" :title="item.value">
        {{ item.value }}
      </dd>
    </div>
  </dl>
</template>
