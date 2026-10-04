<script setup>
defineProps({
  open: { type: Boolean, default: false },
})
defineEmits(['close'])
</script>

<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 bg-black/40 flex items-center justify-center z-50"
      @click.self="$emit('close')"
    >
      <div class="bg-surface-1 rounded-card p-6 max-w-md w-full mx-4">
        <slot />
      </div>
    </div>
  </Teleport>
</template>
