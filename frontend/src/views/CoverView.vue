<script setup>
import BaseButton from '../components/base/BaseButton.vue'
import logo from '../assets/images/app/logo.jpg'

defineEmits(['enter'])
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center gap-6 p-6 text-center">
    <img :src="logo" alt="Phoenix Shuffle" class="w-full max-w-2xl rounded-card" />
    <p class="text-text-secondary max-w-md">
      Bâtis ton équipe, mise gros aux enchères, et affronte tes amis en arène.
    </p>
    <BaseButton @click="$emit('enter')">Entrer dans l'arène</BaseButton>
  </div>
</template>
