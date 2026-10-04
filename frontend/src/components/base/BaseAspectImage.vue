<script setup>
import BaseImage from './BaseImage.vue'

defineProps({
  src: { type: String, default: '' },
  alt: { type: String, default: '' },
  ratio: { type: String, default: '1 / 1' }, // ex. "498 / 949"
  width: { type: String, default: '160px' }, // n'importe quelle valeur CSS
  fit: { type: String, default: 'contain' },
})
</script>

<template>
  <div :style="{ width, aspectRatio: ratio }">
    <BaseImage :src="src" :alt="alt" :fit="fit" />
  </div>
</template>
