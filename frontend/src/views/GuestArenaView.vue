<script setup>
import { computed } from 'vue'
import CardRenderer from '../components/card/CardRenderer.vue'
import CardBack from '../components/card/CardBack.vue'
import BaseButton from '../components/base/BaseButton.vue'
import { locations } from '../data/locations.js'
import { appImage } from '../data/appImages.js'

const props = defineProps({
  state: { type: Object, required: true },
})
const emit = defineEmits(['action'])

const locationImage = computed(() => {
  const loc = locations.find((l) => l.id === props.state.location)
  return loc ? appImage(loc.image) : ''
})

function toCharacter(c) {
  return { id: c.id, name: { full: c.name }, image: { large: c.image } }
}
function toStats(c) {
  return { force: c.force, pv: c.pv }
}

function place(id) {
  emit('action', { type: 'place', characterId: id })
}
function resolve() {
  emit('action', { type: 'resolve' })
}
function nextRound() {
  emit('action', { type: 'nextRound' })
}
</script>

<template>
  <div
    class="min-h-screen bg-cover bg-center"
    :style="locationImage ? { backgroundImage: `linear-gradient(rgba(244,245,248,0.5), rgba(244,245,248,0.5)), url(${locationImage})` } : {}"
  >
  <div class="p-6">
    <h1 class="text-lg font-medium text-text-primary mb-4">Arène — Équipe Rouge</h1>

    <div class="flex items-center justify-center gap-6 mb-2">
      <div class="w-40 h-64 flex items-center justify-center border border-dashed border-border rounded-card">
        <CardRenderer
          v-if="state.slotOpponent && state.slotOpponent !== 'hidden'"
          :character="toCharacter(state.slotOpponent)"
          :stats="toStats(state.slotOpponent)"
          :accent-color="state.slotOpponent.accentColor"
        />
        <CardBack v-else-if="state.slotOpponent === 'hidden'" />
        <span v-else class="text-text-muted text-sm">Emplacement vide</span>
      </div>
      <span class="text-text-muted font-medium">VS</span>
      <div class="w-40 h-64 flex items-center justify-center border border-dashed border-border rounded-card">
        <CardRenderer
          v-if="state.slotSelf"
          :character="toCharacter(state.slotSelf)"
          :stats="toStats(state.slotSelf)"
          :accent-color="state.slotSelf.accentColor"
        />
        <span v-else class="text-text-muted text-sm">Emplacement vide</span>
      </div>
    </div>

    <p v-if="state.resultText" class="text-center text-text-primary font-medium mb-2">{{ state.resultText }}</p>
    <p v-else-if="state.slotSelf && state.slotOpponent" class="text-center text-text-muted text-sm mb-2">
      Les deux cartes sont posées, face cachée — prêtes pour le combat
    </p>

    <div class="flex justify-center gap-3 mb-6">
      <BaseButton v-if="state.slotSelf && state.slotOpponent && !state.revealed" @click="resolve">
        Résoudre le combat
      </BaseButton>
      <BaseButton v-if="state.revealed" variant="secondary" @click="nextRound">Manche suivante</BaseButton>
    </div>

    <p class="text-sm text-force font-medium mb-2">Ta main (clic pour placer)</p>
    <div class="flex gap-3 overflow-x-auto pb-2">
      <div
        v-for="c in state.hand"
        :key="c.id"
        class="cursor-pointer shrink-0"
        @click="place(c.id)"
      >
        <CardRenderer :character="toCharacter(c)" :stats="toStats(c)" :accent-color="c.accentColor" />
      </div>
    </div>
  </div>
  </div>
</template>
