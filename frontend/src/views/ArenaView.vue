<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import CardRenderer from '../components/card/CardRenderer.vue'
import CardBack from '../components/card/CardBack.vue'
import AttackEffect from '../components/card/AttackEffect.vue'
import BaseButton from '../components/base/BaseButton.vue'
import BaseModal from '../components/base/BaseModal.vue'
import { resolveRound, canUseAttack, regenMana, isTeamDefeated } from '../game/combat.js'
import { locations } from '../data/locations.js'
import { appImage } from '../data/appImages.js'
import { BroadcastGameState, GetCharacterAttacks } from '../../wailsjs/go/main/App.js'
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

// Copies locales et mutables — une carte quitte la main quand elle est
// posée, et n'y revient que si elle survit à la manche.
const teamAHand = ref([...props.handBlue])
const teamBHand = ref([...props.handRed])
const queueA = ref([])
const queueB = ref([])
const slotA = ref(null)
const slotB = ref(null)
const revealed = ref(false)
const resultText = ref('')
const detailsEntry = ref(null)
const contextMenu = ref({ visible: false, x: 0, y: 0, entry: null, player: null })
const winner = ref(null)

// Mana par personnage, pour toute la durée de la partie — indépendant du
// fait d'être en main, en file, ou sur l'emplacement de combat.
const manaState = reactive({})

function initMana(entry) {
  const id = entry.character.id
  if (!(id in manaState)) manaState[id] = entry.stats.maxMana ?? 0
  return manaState[id]
}

function manaOf(entry) {
  return entry ? initMana(entry) : 0
}

// Catalogue des attaques par personnage, chargé une fois.
const attacksCatalog = ref({})
onMounted(async () => {
  const rows = await GetCharacterAttacks()
  const grouped = {}
  for (const a of rows) {
    ;(grouped[a.characterId] ??= []).push(a)
  }
  attacksCatalog.value = grouped
  for (const entry of [...teamAHand.value, ...teamBHand.value]) initMana(entry)
})

function attacksFor(entry) {
  return attacksCatalog.value[entry.character.id] ?? []
}

function queueFor(player) {
  return player === 'a' ? queueA : queueB
}
function handFor(player) {
  return player === 'a' ? teamAHand : teamBHand
}

// L'invité (équipe Rouge) ne reçoit que ce qu'il a le droit de voir : sa
// propre main, son propre emplacement, et l'emplacement adverse caché tant
// que la manche n'est pas révélée.
function slotForBroadcast(slot) {
  if (!slot) return null
  return {
    id: slot.entry.character.id,
    name: slot.entry.character.name.full,
    image: slot.entry.character.image.remote || slot.entry.character.image.large,
    accentColor: slot.entry.accentColor,
    force: slot.stats.force,
    pv: slot.stats.pv,
    attackColor: slot.attack?.color ?? null,
  }
}

function broadcastToGuest() {
  BroadcastGameState(
    JSON.stringify({
      type: 'arenaState',
      location: props.location,
      hand: teamBHand.value.map((entry) => ({
        id: entry.character.id,
        name: entry.character.name.full,
        image: entry.character.image.remote || entry.character.image.large,
        accentColor: entry.accentColor,
        force: entry.stats.force,
        pv: entry.stats.pv,
        mana: manaOf(entry),
        attacks: attacksFor(entry),
      })),
      slotSelf: slotForBroadcast(slotB.value),
      slotOpponent: !slotA.value ? null : revealed.value ? slotForBroadcast(slotA.value) : 'hidden',
      revealed: revealed.value,
      resultText: resultText.value,
      winner: winner.value,
    })
  )
}

onMounted(() => {
  broadcastToGuest()
  EventsOn('guest:connected', () => broadcastToGuest())
  EventsOn('guest:action', (payload) => {
    const action = JSON.parse(payload)
    if (action.type === 'place') {
      const entry = teamBHand.value.find((e) => e.character.id === action.characterId)
      if (entry) placeInSlot('b', entry, action.attackId)
    } else if (action.type === 'resolve') {
      resolveCombat()
    } else if (action.type === 'nextRound') {
      nextRound()
    }
  })
})

function placeInSlot(player, entry, attackId) {
  const hand = handFor(player)
  const attack = attackId ? attacksFor(entry).find((a) => a.id === attackId) : null
  if (attack && !canUseAttack(manaOf(entry), attack)) return

  if (attack) manaState[entry.character.id] -= attack.manaCost

  const slot = { entry, stats: { ...entry.stats }, attack: attack ?? null }
  hand.value = hand.value.filter((e) => e !== entry)
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

function placeFromMenu(attackId) {
  placeInSlot(contextMenu.value.player, contextMenu.value.entry, attackId)
}

function queueNextFromMenu() {
  handFor(contextMenu.value.player).value = handFor(contextMenu.value.player).value.filter(
    (e) => e !== contextMenu.value.entry
  )
  queueFor(contextMenu.value.player).value.push(contextMenu.value.entry)
  closeMenu()
}

function queueFirstFromMenu() {
  handFor(contextMenu.value.player).value = handFor(contextMenu.value.player).value.filter(
    (e) => e !== contextMenu.value.entry
  )
  queueFor(contextMenu.value.player).value.unshift(contextMenu.value.entry)
  closeMenu()
}

function showDetailsFromMenu() {
  detailsEntry.value = contextMenu.value.entry
  closeMenu()
}

function regenRoster(player) {
  const all = [...handFor(player).value, ...queueFor(player).value]
  const slot = player === 'a' ? slotA.value : slotB.value
  if (slot) all.push(slot.entry)
  for (const entry of all) {
    const max = entry.stats.maxMana ?? 0
    manaState[entry.character.id] = regenMana(manaOf(entry), entry.stats.manaRegen ?? 0, max)
  }
}

function resolveCombat() {
  if (!slotA.value || !slotB.value || revealed.value) return

  const outcome = resolveRound(
    { force: slotA.value.stats.force, pv: slotA.value.stats.pv, attackBonus: slotA.value.attack?.forceBonus ?? 0 },
    { force: slotB.value.stats.force, pv: slotB.value.stats.pv, attackBonus: slotB.value.attack?.forceBonus ?? 0 }
  )
  revealed.value = true
  slotA.value.stats.pv = outcome.pvAAfter
  slotB.value.stats.pv = outcome.pvBAfter
  slotA.value.eliminated = outcome.eliminatedA
  slotB.value.eliminated = outcome.eliminatedB

  if (outcome.winner === 'draw') {
    resultText.value = 'Égalité parfaite — aucun dégât'
  } else {
    const winnerSlot = outcome.winner === 'a' ? slotA.value : slotB.value
    const loserSlot = outcome.winner === 'a' ? slotB.value : slotA.value
    resultText.value = loserSlot.eliminated
      ? `${winnerSlot.entry.character.name.full} élimine ${loserSlot.entry.character.name.full} !`
      : `${winnerSlot.entry.character.name.full} gagne la manche`
  }

  regenRoster('a')
  regenRoster('b')
  broadcastToGuest()
}

function nextRound() {
  settleSlot('a')
  settleSlot('b')
  revealed.value = false
  resultText.value = ''
  broadcastToGuest()
  checkVictory()
}

function settleSlot(player) {
  const slotRef = player === 'a' ? slotA : slotB
  const hand = handFor(player)
  const queue = queueFor(player)

  if (slotRef.value && !slotRef.value.eliminated) {
    hand.value.push(slotRef.value.entry)
  }
  slotRef.value = queue.value.length ? { entry: queue.value.shift(), stats: null, attack: null } : null
  if (slotRef.value && !slotRef.value.stats) slotRef.value.stats = { ...slotRef.value.entry.stats }
}

function checkVictory() {
  if (winner.value) return
  const aDefeated = isTeamDefeated(teamAHand.value, queueA.value, slotA.value)
  const bDefeated = isTeamDefeated(teamBHand.value, queueB.value, slotB.value)
  if (aDefeated && bDefeated) winner.value = 'draw'
  else if (aDefeated) winner.value = 'b'
  else if (bDefeated) winner.value = 'a'
  if (winner.value) broadcastToGuest()
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

    <div v-if="winner" class="mb-6 p-4 rounded-card bg-surface-1 border border-accent text-center">
      <p class="text-lg font-medium text-text-primary">
        <template v-if="winner === 'draw'">Match nul — plus personne ne tient debout !</template>
        <template v-else>{{ winner === 'a' ? 'Équipe Bleue' : 'Équipe Rouge' }} remporte la partie !</template>
      </p>
    </div>

    <div class="flex items-center justify-center gap-6 mb-2">
      <div class="relative w-40 h-64 flex items-center justify-center border border-dashed border-border rounded-card">
        <AttackEffect :color="slotA?.attack?.color" :active="revealed && !!slotA?.attack" />
        <CardRenderer
          v-if="slotA && revealed"
          :character="slotA.entry.character"
          :stats="slotA.stats"
          :accent-color="slotA.entry.accentColor"
        />
        <CardBack v-else-if="slotA" />
        <span v-else class="text-text-muted text-sm">Emplacement vide</span>
        <span v-if="slotA?.attack && revealed" class="absolute -bottom-2 text-xs px-2 py-0.5 rounded-full bg-surface-1 border border-border" :style="{ color: slotA.attack.color }">
          {{ slotA.attack.name }}
        </span>
      </div>
      <span class="text-text-muted font-medium">VS</span>
      <div class="relative w-40 h-64 flex items-center justify-center border border-dashed border-border rounded-card">
        <AttackEffect :color="slotB?.attack?.color" :active="revealed && !!slotB?.attack" />
        <CardRenderer
          v-if="slotB && revealed"
          :character="slotB.entry.character"
          :stats="slotB.stats"
          :accent-color="slotB.entry.accentColor"
        />
        <CardBack v-else-if="slotB" />
        <span v-else class="text-text-muted text-sm">Emplacement vide</span>
        <span v-if="slotB?.attack && revealed" class="absolute -bottom-2 text-xs px-2 py-0.5 rounded-full bg-surface-1 border border-border" :style="{ color: slotB.attack.color }">
          {{ slotB.attack.name }}
        </span>
      </div>
    </div>

    <p v-if="resultText" class="text-center text-text-primary font-medium mb-2">{{ resultText }}</p>
    <p v-else-if="slotA && slotB" class="text-center text-text-muted text-sm mb-2">
      Les deux cartes sont posées, face cachée — prêtes pour le combat
    </p>

    <div class="flex justify-center gap-3 mb-6">
      <BaseButton v-if="slotA && slotB && !revealed && !winner" @click="resolveCombat">Résoudre le combat</BaseButton>
      <BaseButton v-if="revealed && !winner" variant="secondary" @click="nextRound">Manche suivante</BaseButton>
    </div>

    <div class="mb-6">
      <p class="text-sm text-accent font-medium mb-2">Équipe Bleue (clic pour placer, clic droit pour lier une attaque)</p>
      <div class="flex gap-3 overflow-x-auto pb-2">
        <div
          v-for="entry in teamAHand"
          :key="'a-' + entry.character.id"
          class="cursor-pointer shrink-0 flex flex-col items-center gap-1"
          @click="placeInSlot('a', entry)"
          @contextmenu.prevent="openMenu($event, 'a', entry)"
        >
          <CardRenderer :character="entry.character" :stats="entry.stats" :accent-color="entry.accentColor" />
          <span class="text-xs text-text-muted">⚡ {{ manaOf(entry) }}/{{ entry.stats.maxMana }}</span>
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
      <p class="text-sm text-force font-medium mb-2">Équipe Rouge (clic pour placer, clic droit pour lier une attaque)</p>
      <div class="flex gap-3 overflow-x-auto pb-2">
        <div
          v-for="entry in teamBHand"
          :key="'b-' + entry.character.id"
          class="cursor-pointer shrink-0 flex flex-col items-center gap-1"
          @click="placeInSlot('b', entry)"
          @contextmenu.prevent="openMenu($event, 'b', entry)"
        >
          <CardRenderer :character="entry.character" :stats="entry.stats" :accent-color="entry.accentColor" />
          <span class="text-xs text-text-muted">⚡ {{ manaOf(entry) }}/{{ entry.stats.maxMana }}</span>
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
      class="fixed z-50 bg-surface-1 border border-border rounded-card shadow-lg py-1 text-sm min-w-[220px]"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
      @click.stop
    >
      <button class="block w-full text-left px-3 py-1.5 hover:bg-surface-2" @click="placeFromMenu()">
        Placer sans attaque
      </button>
      <button
        v-for="attack in attacksFor(contextMenu.entry ?? { character: { id: null } })"
        :key="attack.id"
        class="block w-full text-left px-3 py-1.5 hover:bg-surface-2 disabled:opacity-40 disabled:cursor-not-allowed"
        :disabled="!canUseAttack(manaOf(contextMenu.entry), attack)"
        @click="placeFromMenu(attack.id)"
      >
        <span class="inline-block w-2 h-2 rounded-full mr-1" :style="{ backgroundColor: attack.color }"></span>
        Lier {{ attack.name }} (+{{ attack.forceBonus }} FOR, -{{ attack.manaCost }} mana)
      </button>
      <div class="border-t border-border my-1"></div>
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
        <p class="text-sm text-text-secondary">Mana {{ manaOf(detailsEntry) }}/{{ detailsEntry.stats.maxMana }}</p>
      </div>
    </BaseModal>
  </div>
  </div>
</template>
