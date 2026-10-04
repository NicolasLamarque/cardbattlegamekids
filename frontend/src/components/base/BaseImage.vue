<script setup>
import { ref, watch } from 'vue'
import { ImageOff } from 'lucide-vue-next'

const props = defineProps({
  src: { type: String, default: '' },
  alt: { type: String, default: '' },
  fit: { type: String, default: 'cover' }, // 'cover' | 'contain'
})

const failed = ref(false)
watch(() => props.src, () => { failed.value = false })
</script>

<template>
  <div class="w-full h-full flex items-center justify-center bg-accent-soft overflow-hidden">
    <img
      v-if="src && !failed"
      :src="src"
      :alt="alt"
      class="w-full h-full"
      :class="fit === 'contain' ? 'object-contain' : 'object-cover'"
      @error="failed = true"
    />
    <ImageOff v-else class="w-8 h-8 text-accent" />
  </div>
</template>
