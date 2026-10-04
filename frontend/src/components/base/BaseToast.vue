<script setup>
defineProps({
  message: { type: String, default: '' },
  tone: { type: String, default: 'accent' }, // 'accent' | 'force'
})
</script>

<template>
  <Transition name="toast">
    <div
      v-if="message"
      class="fixed top-4 left-1/2 -translate-x-1/2 z-50 px-4 py-2 rounded-card text-sm font-medium shadow-lg"
      :class="tone === 'accent' ? 'bg-accent text-white' : 'bg-force text-white'"
    >
      {{ message }}
    </div>
  </Transition>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, -12px);
}
</style>
