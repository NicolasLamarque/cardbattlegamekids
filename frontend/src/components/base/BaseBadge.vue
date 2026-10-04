<script setup>
defineProps({
  tone: { type: String, default: 'accent' },
})
</script>

<template>
  <span
    class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
    :class="{
      'bg-accent-soft text-accent': tone === 'accent',
      'bg-force-soft text-force': tone === 'force',
      'bg-pv-soft text-pv': tone === 'pv',
      'bg-gold-soft text-gold': tone === 'gold',
    }"
  >
    <slot />
  </span>
</template>
