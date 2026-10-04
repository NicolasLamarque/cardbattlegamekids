<script setup>
import { ref, computed, onMounted } from 'vue'
import CardRenderer from '../components/card/CardRenderer.vue'
import CardBack from '../components/card/CardBack.vue'
import BaseButton from '../components/base/BaseButton.vue'
import BaseModal from '../components/base/BaseModal.vue'
import { resolveRound } from '../game/combat.js'
import { locations } from '../data/locations.js'
import { appImage } from '../data/appImages.js'
import { BroadcastGameState } from '../../wailsjs/go/main/App.js'
import { EventsOn } from '../../wailsjs/runtime/runtime.js'

const props = defineProps({
  specialCardsEnabled: { type: Boolean, default: true },
  handBlue: { type: Array, default: () => [] },
  handRed: { type: Array, default: () => [] },
  location: { type: String, default: 'tokyo' },
  connection: { type: String, default: 'solo' },
})

const locationImage = computed(() => {
  const loc = locations.find((l) => l.id === props.location)
  return loc ? appImage(loc.image) : ''
})

const slotA = ref(null)
const slotB = ref(null)
const queueA = ref([])
const queueB = ref([])
const revealed = ref(false)
const resultText = ref('')
const detailsEntry = ref(null)
const contextMenu = ref({ visible: false, x: 0, y: 0, entry: null, player: null })

function queueFor(player) {
  return player === 'a' ? queueA : queueB
}

// L'invité (équipe Rouge) ne reçoit que ce qu'il a le droit de voir : sa
// propre main, son propre emplacement, et l'emplacement adverse caché tant
// que la manche n'est pas révélée.
function slotForBroadcast(slot) {
  if (!slot) return null
  return {
    id: slot.entry.character.id,
    name: slot.entry.character.name.full,
    // L'invité est sur un autre appareil — le chemin d'image local (résolu
    // par Vite côté hôte) ne lui est d'aucune utilité, on lui envoie
    // toujours l'URL distante AniList, qu'il peut charger lui-même.
    image: slot.entry.character.image.remote || slot.entry.character.image.large,
    accentColor: slot.entry.accentColor,
    force: slot.stats.force,
    pv: slot.stats.pv,
  }
}

function broadcastToGuest() {
  BroadcastGameState(
    JSON.stringify({
      type: 'arenaState',
      location: props.location,
      hand: props.handRed.map((entry) => ({
        id: entry.character.id,
        name: entry.character.name.full,
        image: entry.character.image.remote || entry.character.image.large,
        accentColor: entry.accentColor,
        force: entry.stats.force,
        pv: entry.stats.pv,
      })),
      slotSelf: slotForBroadcast(slotB.value),
      slotOpponent: !slotA.value ? null : revealed.value ? slotForBroadcast(slotA.value) : 'hidden',
      revealed: revealed.value,
      resultText: resultText.value,
    })
  )
}

onMounted(() => {
  broadcastToGuest()
  EventsOn('guest:connected', () => broadcastToGuest())
  EventsOn('guest:action', (payload) => {
    const action = JSON.parse(payload)
    if (action.type === 'place') {
      const entry = props.handRed.find((e) => e.character.id === action.characterId)
      if (entry) placeInSlot('b', entry)
    } else if (action.type === 'resolve') {
      resolveCombat()
    } else if (action.type === 'nextRound') {
      nextRound()
    }
  })
})

function placeInSlot(player, entry) {
  const slot = { entry, stats: { ...entry.stats } }
  if (player === 'a') slotA.value = slot
  else slotB.value = slot
  revealed.value = false
  resultText.value = ''
  closeMenu()
  broadcastToGuest()
}

function openMenu(event, player, entry) {
  contextMenu.value = { visible: true, x: event.clientX, y: event.clientY, entry, player }
}

function closeMenu() {
  contextMenu.value.visible = false
}

function placeFromMenu() {
  placeInSlot(contextMenu.value.player, contextMenu.value.entry)
}

function queueNextFromMenu() {
  queueFor(contextMenu.value.player).value.push(contextMenu.value.entry)
  closeMenu()
}

function queueFirstFromMenu() {
  queueFor(contextMenu.value.player).value.unshift(contextMenu.value.entry)
  closeMenu()
}

function showDetailsFromMenu() {
  detailsEntry.value = contextMenu.value.entry
  closeMenu()
}

function resolveCombat() {
  if (!slotA.value || !slotB.value || revealed.value) return
  const outcome = resolveRound(slotA.value.stats, slotB.value.stats)
  revealed.value = true

  if (outcome.winner === 'draw') {
    resultText.value = 'Égalité parfaite — aucun dégât'
  } else {
    const winnerSlot = outcome.winner === 'a' ? slotA.value : slotB.value
    winnerSlot.stats = outcome.winnerStatsAfter
    resultText.value = `${winnerSlot.entry.character.name.full} gagne, mais encaisse les dégâts du combat`
  }

  broadcastToGuest()
}

function nextRound() {
  slotA.value = queueA.value.length ? { entry: queueA.value.shift(), stats: null } : slotA.value
  slotB.value = queueB.value.length ? { entry: queueB.value.shift(), stats: null } : slotB.value
  if (slotA.value && !slotA.value.stats) slotA.value.stats = { ...slotA.value.entry.stats }
  if (slotB.value && !slotB.value.stats) slotB.value.stats = { ...slotB.value.entry.stats }
  revealed.value = false
  resultText.value = ''
  broadcastToGuest()
}
</script>

<template>
  <div
    class="min-h-screen bg-cover bg-center"
    :style="locationImage ? { backgroundImage: `linear-gradient(rgba(244,245,248,0.5), rgba(244,245,248,0.5)), url(${locationImage})` } : {}"
  >
  <div class="p-6" @click="closeMenu">
    <div class="flex items-center gap-2 mb-4">
      <h1 class="text-lg font-medium text-text-primary">Arène</h1>
      <span
        class="text-xs px-2 py-0.5 rounded-full"
        :class="specialCardsEnabled ? 'bg-accent-soft text-accent' : 'bg-surface-2 text-text-muted'"
      >
        Cartes spéciales {{ specialCardsEnabled ? 'activées' : 'désactivées' }}
      </span>
    </div>

    <div class="flex items-center justify-center gap-6 mb-2">
      <div class="w-40 h-64 flex items-center justify-center border border-dashed border-border rounded-card">
        <CardRenderer
          v-if="slotA && revealed"
          :character="slotA.entry.character"
          :stats="slotA.stats"
          :accent-color="slotA.entry.accentColor"
        />
        <CardBack v-else-if="slotA" />
        <span v-else class="text-text-muted text-sm">Emplacement vide</span>
      </div>
      <span class="text-text-muted font-medium">VS</span>
      <div class="w-40 h-64 flex items-center justify-center border border-dashed border-border rounded-card">
        <CardRenderer
          v-if="slotB && revealed"
          :character="slotB.entry.character"
          :stats="slotB.stats"
          :accent-color="slotB.entry.accentColor"
        />
        <CardBack v-else-if="slotB" />
        <span v-else class="text-text-muted text-sm">Emplacement vide</span>
      </div>
    </div>

    <p v-if="resultText" class="text-center text-text-primary font-medium mb-2">{{ resultText }}</p>
    <p v-else-if="slotA && slotB" class="text-center text-text-muted text-sm mb-2">
      Les deux cartes sont posées, face cachée — prêtes pour le combat
    </p>

    <div class="flex justify-center gap-3 mb-6">
      <BaseButton v-if="slotA && slotB && !revealed" @click="resolveCombat">Résoudre le combat</BaseButton>
      <BaseButton v-if="revealed" variant="secondary" @click="nextRound">Manche suivante</BaseButton>
    </div>

    <div class="mb-6">
      <p class="text-sm text-accent font-medium mb-2">Équipe Bleue (clic pour placer, clic droit pour les options)</p>
      <div class="flex gap-3 overflow-x-auto pb-2">
        <div
          v-for="entry in handBlue"
          :key="'a-' + entry.character.id"
          class="cursor-pointer shrink-0"
          @click="placeInSlot('a', entry)"
          @contextmenu.prevent="openMenu($event, 'a', entry)"
        >
          <CardRenderer :character="entry.character" :stats="entry.stats" :accent-color="entry.accentColor" />
        </div>
      </div>
      <div v-if="queueA.length" class="flex items-center gap-2 mt-2 text-xs text-text-secondary">
        <span>File d'attente :</span>
        <span v-for="(q, i) in queueA" :key="q.character.id" class="px-2 py-0.5 rounded-full bg-accent-soft text-accent">
          {{ i + 1 }}. {{ q.character.name.full }}
        </span>
      </div>
    </div>

    <div v-if="connection === 'solo'">
      <p class="text-sm text-force font-medium mb-2">Équipe Rouge (clic pour placer, clic droit pour les options)</p>
      <div class="flex gap-3 overflow-x-auto pb-2">
        <div
          v-for="entry in handRed"
          :key="'b-' + entry.character.id"
          class="cursor-pointer shrink-0"
          @click="placeInSlot('b', entry)"
          @contextmenu.prevent="openMenu($event, 'b', entry)"
        >
          <CardRenderer :character="entry.character" :stats="entry.stats" :accent-color="entry.accentColor" />
        </div>
      </div>
      <div v-if="queueB.length" class="flex items-center gap-2 mt-2 text-xs text-text-secondary">
        <span>File d'attente :</span>
        <span v-for="(q, i) in queueB" :key="q.character.id" class="px-2 py-0.5 rounded-full bg-force-soft text-force">
          {{ i + 1 }}. {{ q.character.name.full }}
        </span>
      </div>
    </div>

    <div v-else class="p-3 bg-surface-1 rounded-card text-sm text-text-muted">
      Équipe Rouge — contrôlée par le joueur invité sur son propre appareil. Sa main reste privée.
    </div>

    <div
      v-if="contextMenu.visible"
      class="fixed z-50 bg-surface-1 border border-border rounded-card shadow-lg py-1 text-sm"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
      @click.stop
    >
      <button class="block w-full text-left px-3 py-1.5 hover:bg-surface-2" @click="placeFromMenu">Placer en arène</button>
      <button class="block w-full text-left px-3 py-1.5 hover:bg-surface-2" @click="queueNextFromMenu">Mettre au prochain tour</button>
      <button class="block w-full text-left px-3 py-1.5 hover:bg-surface-2" @click="queueFirstFromMenu">Mettre en premier à la prochaine joute</button>
      <button class="block w-full text-left px-3 py-1.5 hover:bg-surface-2" @click="showDetailsFromMenu">Voir les détails</button>
    </div>

    <BaseModal :open="!!detailsEntry" @close="detailsEntry = null">
      <div v-if="detailsEntry" class="flex flex-col items-center gap-3">
        <CardRenderer
          :character="detailsEntry.character"
          :stats="detailsEntry.stats"
          :accent-color="detailsEntry.accentColor"
        />
        <p class="text-sm text-text-secondary">{{ detailsEntry.character.favourites }} favoris AniList</p>
      </div>
    </BaseModal>
  </div>
  </div>
</template>
