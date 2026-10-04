<script setup>
defineProps({
  name: { type: String, required: true },
})
</script>

<template>
  <p class="px-2 pt-2 text-sm font-medium text-text-primary truncate">{{ name }}</p>
</template>
