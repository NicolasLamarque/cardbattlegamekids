<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: Number, required: true },
  max: { type: Number, default: 100 },
  tone: { type: String, default: 'accent' },
})

const percent = computed(() =>
  Math.max(0, Math.min(100, Math.round((props.value / props.max) * 100)))
)
</script>

<template>
  <div class="flex items-center gap-2">
    <span class="text-xs text-text-secondary w-8">{{ label }}</span>
    <div class="flex-1 h-1.5 rounded-full bg-surface-2 overflow-hidden">
      <div
        class="h-full rounded-full"
        :class="{
          'bg-force': tone === 'force',
          'bg-pv': tone === 'pv',
          'bg-accent': tone === 'accent',
          'bg-gold': tone === 'gold',
        }"
        :style="{ width: percent + '%' }"
      />
    </div>
    <span class="text-xs text-text-secondary w-6 text-right">{{ value }}</span>
  </div>
</template>
