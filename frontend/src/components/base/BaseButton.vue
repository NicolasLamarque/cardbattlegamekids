<script setup>
defineProps({
  variant: { type: String, default: 'primary' },
})
</script>

<template>
  <button
    class="px-4 py-2 rounded-card font-medium transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
    :class="{
      'bg-accent text-white hover:opacity-90': variant === 'primary',
      'bg-surface-2 text-text-primary hover:bg-border': variant === 'secondary',
      'bg-transparent text-text-secondary hover:bg-surface-2': variant === 'ghost',
    }"
  >
    <slot />
  </button>
</template>
