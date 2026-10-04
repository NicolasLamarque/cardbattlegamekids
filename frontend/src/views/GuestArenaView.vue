<script setup>
import { ref, computed } from 'vue'
import CardRenderer from '../components/card/CardRenderer.vue'
import CardBack from '../components/card/CardBack.vue'
import AttackEffect from '../components/card/AttackEffect.vue'
import BaseButton from '../components/base/BaseButton.vue'
import { locations } from '../data/locations.js'
import { appImage } from '../data/appImages.js'

const props = defineProps({
  state: { type: Object, required: true },
})
const emit = defineEmits(['action'])

const attackMenuFor = ref(null)

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

function place(id, attackId) {
  attackMenuFor.value = null
  emit('action', { type: 'place', characterId: id, attackId })
}
function toggleAttackMenu(id) {
  attackMenuFor.value = attackMenuFor.value === id ? null : id
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

    <div v-if="state.winner" class="mb-6 p-4 rounded-card bg-surface-1 border border-accent text-center">
      <p class="text-lg font-medium text-text-primary">
        <template v-if="state.winner === 'draw'">Match nul — plus personne ne tient debout !</template>
        <template v-else-if="state.winner === 'b'">Ton équipe remporte la partie !</template>
        <template v-else>L'équipe adverse remporte la partie…</template>
      </p>
    </div>

    <div class="flex items-center justify-center gap-6 mb-2">
      <div class="relative w-40 h-64 flex items-center justify-center border border-dashed border-border rounded-card">
        <AttackEffect :color="state.slotOpponent?.attackColor" :active="state.revealed && !!state.slotOpponent?.attackColor" />
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
      <div class="relative w-40 h-64 flex items-center justify-center border border-dashed border-border rounded-card">
        <AttackEffect :color="state.slotSelf?.attackColor" :active="state.revealed && !!state.slotSelf?.attackColor" />
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

    <p class="text-sm text-force font-medium mb-2">Ta main (clic pour placer, re-clic pour lier une attaque)</p>
    <div class="flex gap-3 overflow-x-auto pb-2">
      <div v-for="c in state.hand" :key="c.id" class="shrink-0 flex flex-col items-center gap-1">
        <div class="relative cursor-pointer" @click="c.attacks?.length ? toggleAttackMenu(c.id) : place(c.id)">
          <CardRenderer :character="toCharacter(c)" :stats="toStats(c)" :accent-color="c.accentColor" />
        </div>
        <span class="text-xs text-text-muted">⚡ {{ c.mana }}</span>
        <div v-if="attackMenuFor === c.id" class="flex flex-col gap-1 bg-surface-1 border border-border rounded-card p-2 text-xs">
          <button class="text-left px-2 py-1 rounded hover:bg-surface-2" @click="place(c.id)">Sans attaque</button>
          <button
            v-for="attack in c.attacks"
            :key="attack.id"
            class="text-left px-2 py-1 rounded hover:bg-surface-2 disabled:opacity-40"
            :disabled="c.mana < attack.manaCost"
            @click="place(c.id, attack.id)"
          >
            <span class="inline-block w-2 h-2 rounded-full mr-1" :style="{ backgroundColor: attack.color }"></span>
            {{ attack.name }} (+{{ attack.forceBonus }} FOR, -{{ attack.manaCost }} mana)
          </button>
        </div>
      </div>
    </div>
  </div>
  </div>
</template>
