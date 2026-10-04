<script setup>
import BaseImage from '../base/BaseImage.vue'

defineProps({
  slots: { type: Array, required: true },
  total: { type: Number, required: true },
  tone: { type: String, default: 'accent' }, // 'accent' | 'force'
})
</script>

<template>
  <div class="grid grid-cols-2 gap-1.5 content-start w-24 shrink-0">
    <div
      v-for="i in total"
      :key="i"
      class="w-10 h-10 rounded-card overflow-hidden border"
      :class="i <= slots.length ? (tone === 'accent' ? 'border-accent' : 'border-force') : 'border-dashed border-border'"
    >
      <BaseImage
        v-if="i <= slots.length"
        :src="slots[i - 1].entry.character.image.large"
        :alt="slots[i - 1].entry.character.name.full"
      />
    </div>
  </div>
</template>
