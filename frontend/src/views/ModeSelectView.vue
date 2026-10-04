<script setup>
import BaseImage from '../components/base/BaseImage.vue'
import { appImage } from '../data/appImages.js'

defineEmits(['select', 'admin'])
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center gap-8 p-6 relative">
    <button
      class="absolute top-4 right-4 text-xs text-text-muted hover:text-text-secondary cursor-pointer"
      @click="$emit('admin')"
    >
      Gestion des stats
    </button>
    <h1 class="text-lg font-medium text-text-primary">Choisis ton mode</h1>
    <div class="flex gap-6 flex-wrap justify-center">
      <button
        class="w-64 h-80 rounded-card overflow-hidden border border-border hover:border-accent transition-colors flex flex-col cursor-pointer"
        @click="$emit('select', 'career')"
      >
        <div class="flex-1"><BaseImage :src="appImage('mode-carriere.png')" alt="Mode Carrière" /></div>
       
      </button>
      <button
        class="w-64 h-80 rounded-card overflow-hidden border border-border hover:border-accent transition-colors flex flex-col cursor-pointer"
        @click="$emit('select', 'skirmish')"
      >
        <div class="flex-1"><BaseImage :src="appImage('mode-escarmouche.jpg')" alt="Mode Escarmouche" /></div>
      
      </button>
    </div>
  </div>
</template>
