<script setup lang="ts">
import type { MetricLog } from '~/types/device'

const props = defineProps<{
  metrics: MetricLog[]
}>()

const width = 640
const height = 220
const pad = { top: 16, right: 16, bottom: 28, left: 36 }

const series = computed(() => {
  // API returns newest-first; chart left→right should be oldest→newest.
  const rows = [...props.metrics].reverse().slice(-30)
  return rows.map(r => ({
    t: new Date(r.timestamp).getTime(),
    cpu: r.cpu_usage,
    ram: r.ram_usage,
    label: new Date(r.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
  }))
})

const plotW = width - pad.left - pad.right
const plotH = height - pad.top - pad.bottom

function xAt(i: number, n: number) {
  if (n <= 1) {
    return pad.left + plotW / 2
  }
  return pad.left + (i / (n - 1)) * plotW
}

function yAt(value: number) {
  const clamped = Math.min(100, Math.max(0, value))
  return pad.top + plotH - (clamped / 100) * plotH
}

function linePath(key: 'cpu' | 'ram') {
  const rows = series.value
  if (!rows.length) {
    return ''
  }
  return rows
    .map((row, i) => {
      const cmd = i === 0 ? 'M' : 'L'
      return `${cmd}${xAt(i, rows.length).toFixed(1)} ${yAt(row[key]).toFixed(1)}`
    })
    .join(' ')
}

const yTicks = [0, 25, 50, 75, 100]
</script>

<template>
  <div v-if="!series.length" class="px-2 py-8 text-center text-sm font-bold text-ink/50">
    Not enough samples for a chart yet.
  </div>
  <div v-else class="w-full overflow-x-auto">
    <svg
      :viewBox="`0 0 ${width} ${height}`"
      class="h-auto w-full min-w-[320px]"
      role="img"
      aria-label="CPU and RAM usage over time"
    >
      <rect
        :x="pad.left"
        :y="pad.top"
        :width="plotW"
        :height="plotH"
        class="fill-paper stroke-ink"
        stroke-width="2"
      />

      <g v-for="tick in yTicks" :key="tick">
        <line
          :x1="pad.left"
          :x2="pad.left + plotW"
          :y1="yAt(tick)"
          :y2="yAt(tick)"
          class="stroke-ink/15"
          stroke-width="1"
          stroke-dasharray="4 4"
        />
        <text
          :x="pad.left - 8"
          :y="yAt(tick) + 4"
          text-anchor="end"
          class="fill-ink/50 font-mono text-[10px]"
        >
          {{ tick }}
        </text>
      </g>

      <path
        :d="linePath('cpu')"
        fill="none"
        class="stroke-ggreen"
        stroke-width="3"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
      <path
        :d="linePath('ram')"
        fill="none"
        class="stroke-gblue"
        stroke-width="3"
        stroke-linecap="round"
        stroke-linejoin="round"
      />

      <g v-for="(row, i) in series" :key="i">
        <circle :cx="xAt(i, series.length)" :cy="yAt(row.cpu)" r="3.5" class="fill-ggreen stroke-ink" stroke-width="1.5" />
        <circle :cx="xAt(i, series.length)" :cy="yAt(row.ram)" r="3.5" class="fill-gblue stroke-ink" stroke-width="1.5" />
      </g>

      <text
        v-if="series[0]"
        :x="pad.left"
        :y="height - 8"
        class="fill-ink/45 font-mono text-[10px]"
      >
        {{ series[0].label }}
      </text>
      <text
        v-if="series.length > 1"
        :x="pad.left + plotW"
        :y="height - 8"
        text-anchor="end"
        class="fill-ink/45 font-mono text-[10px]"
      >
        {{ series[series.length - 1].label }}
      </text>
    </svg>

    <div class="mt-3 flex flex-wrap gap-3 text-xs font-bold">
      <span class="inline-flex items-center gap-2">
        <span class="size-3 rounded-full border-2 border-ink bg-ggreen" />
        CPU
      </span>
      <span class="inline-flex items-center gap-2">
        <span class="size-3 rounded-full border-2 border-ink bg-gblue" />
        RAM
      </span>
    </div>
  </div>
</template>
