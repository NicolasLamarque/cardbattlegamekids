<script setup>
defineProps({
  accentColor: { type: String, default: '' },
})
</script>

<template>
  <div
    class="w-40 bg-surface-1 border rounded-card overflow-hidden flex flex-col"
    :style="{ borderColor: accentColor || 'var(--color-border-strong)' }"
  >
    <div class="h-1" :style="{ backgroundColor: accentColor || 'var(--color-accent)' }" />
    <slot />
  </div>
</template>
