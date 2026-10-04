<script setup>
import BaseButton from '../components/base/BaseButton.vue'

defineProps({
  hands: { type: Object, required: true },
})
defineEmits(['start'])
</script>

<template>
  <div class="p-6">
    <h1 class="text-lg font-medium text-text-primary mb-4">Distribution des cartes spéciales</h1>

    <div class="flex gap-6 mb-6">
      <div class="flex-1">
        <p class="font-medium text-accent mb-2">Équipe Bleue</p>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="(card, i) in hands.blue"
            :key="i"
            class="text-sm px-3 py-1.5 rounded-card bg-accent-soft text-accent"
          >
            {{ card.name }}
          </span>
        </div>
      </div>
      <div class="flex-1">
        <p class="font-medium text-force mb-2">Équipe Rouge</p>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="(card, i) in hands.red"
            :key="i"
            class="text-sm px-3 py-1.5 rounded-card bg-force-soft text-force"
          >
            {{ card.name }}
          </span>
        </div>
      </div>
    </div>

    <p class="text-xs text-text-muted mb-4">
      Pour l'instant les deux mains s'affichent sur le même écran — la vue privée par appareil
      (téléphone + code QR) arrivera avec le serveur local.
    </p>

    <BaseButton @click="$emit('start')">Commencer la partie</BaseButton>
  </div>
</template>
