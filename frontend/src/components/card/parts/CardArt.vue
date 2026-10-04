<script setup>
import BaseImage from '../../base/BaseImage.vue'

defineProps({
  src: { type: String, default: '' },
  alt: { type: String, default: '' },
})
</script>

<template>
  <div class="h-28">
    <BaseImage :src="src" :alt="alt" />
  </div>
</template>
