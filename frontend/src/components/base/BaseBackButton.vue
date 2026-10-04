<script setup>
defineEmits(['click'])
</script>

<template>
  <button
    class="m-4 text-sm text-text-secondary hover:text-text-primary cursor-pointer"
    @click="$emit('click')"
  >
    ← Retour
  </button>
</template>
