<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import CardRenderer from '../components/card/CardRenderer.vue'
import BaseButton from '../components/base/BaseButton.vue'
import BaseAspectImage from '../components/base/BaseAspectImage.vue'
import TeamSlotsColumn from '../components/draft/TeamSlotsColumn.vue'
import { rollModifier } from '../game/draftModifiers.js'
import { localCharacterImage } from '../data/images.js'
import { appImage } from '../data/appImages.js'
import { GetCharacters, GetDraftModifiers, BroadcastGameState } from '../../wailsjs/go/main/App.js'
import { EventsOn } from '../../wailsjs/runtime/runtime.js'

const props = defineProps({
  config: { type: Object, required: true },
})
const emit = defineEmits(['done'])

const BID_PRESETS = [10000, 20000, 50000]
const POT_WIDTH = '160px'

const pool = ref([])
const modifiersCatalog = ref([])
const teams = reactive({
  blue: { id: 'blue', name: 'Équipe Bleue', budget: props.config.startingBudget, slots: [], slotBonus: 0 },
  red: { id: 'red', name: 'Équipe Rouge', budget: props.config.startingBudget, slots: [], slotBonus: 0 },
})
const currentLot = ref(null)
const currentModifier = ref(null)
const currentBid = ref(0)
const currentBidder = ref(null)
const customBidBlue = ref(0)
const customBidRed = ref(0)
const bidErrors = reactive({ blue: '', red: '' })
const expandedTeam = ref(null)

onMounted(async () => {
  const rows = await GetCharacters()
  pool.value = rows.map((c) => ({
    character: {
      id: c.id,
      name: { full: c.nameFull, native: c.nameNative },
      image: { large: localCharacterImage(c.imageLocal) || c.imageLarge, remote: c.imageLarge },
      favourites: c.favourites,
    },
    stats: { force: c.force, pv: c.pv, maxMana: c.maxMana, manaRegen: c.manaRegen },
    accentColor: c.seriesColor,
  }))
  modifiersCatalog.value = props.config.modifiersEnabled ? await GetDraftModifiers() : []

  if (props.config.connection === 'lan') {
    broadcastDraftState()
    EventsOn('guest:connected', () => broadcastDraftState())
    EventsOn('guest:action', (payload) => {
      const action = JSON.parse(payload)
      if (action.type === 'bid') placeBid('red', action.amount)
      else if (action.type === 'pass') pass('red')
    })
  }
})

// L'invité (équipe Rouge) voit le lot courant et peut miser depuis son
// téléphone — le Bocal n'a pas d'information cachée entre équipes, donc on
// lui envoie l'état tel quel, pas besoin de filtrer comme en Arène.
function broadcastDraftState() {
  if (props.config.connection !== 'lan') return
  BroadcastGameState(
    JSON.stringify({
      type: 'draftState',
      currentLot: currentLot.value
        ? {
            character: {
              id: currentLot.value.character.id,
              name: currentLot.value.character.name,
              image: { large: currentLot.value.character.image.remote || currentLot.value.character.image.large },
              favourites: currentLot.value.character.favourites,
            },
            stats: currentLot.value.stats,
            accentColor: currentLot.value.accentColor,
          }
        : null,
      currentModifier: currentModifier.value,
      currentBid: currentBid.value,
      currentBidder: currentBidder.value,
      bidError: bidErrors.red,
      red: {
        budget: teams.red.budget,
        slotsUsed: teams.red.slots.length,
        slotsTotal: effectiveSlots('red'),
      },
      blue: {
        budget: teams.blue.budget,
        slotsUsed: teams.blue.slots.length,
        slotsTotal: effectiveSlots('blue'),
      },
      draftDone: draftDone.value,
    })
  )
}

function customBidRefFor(teamId) {
  return teamId === 'blue' ? customBidBlue : customBidRed
}

function effectiveSlots(teamId) {
  return props.config.slots + teams[teamId].slotBonus
}

function totalForce(teamId) {
  return teams[teamId].slots.reduce((sum, s) => sum + s.entry.stats.force, 0)
}

function totalPV(teamId) {
  return teams[teamId].slots.reduce((sum, s) => sum + s.entry.stats.pv, 0)
}

function canBidAmount(teamId, amount) {
  const team = teams[teamId]
  return amount > 0 && team.budget >= amount && team.slots.length < effectiveSlots(teamId)
}

// Mise libre : aucun prix suggéré ni plancher — un personnage peut partir
// pour 1 $ comme pour tout le budget d'une équipe, à leur discrétion.
function canAffordAnything(teamId) {
  return teams[teamId].budget > 0 && teams[teamId].slots.length < effectiveSlots(teamId)
}

const notFull = computed(() => ['blue', 'red'].filter((id) => teams[id].slots.length < effectiveSlots(id)))
const activeTeams = computed(() => notFull.value.filter((id) => canAffordAnything(id)))
const strandedTeams = computed(() => notFull.value.filter((id) => !canAffordAnything(id)))

const draftDone = computed(
  () => notFull.value.length === 0 || (pool.value.length === 0 && !currentLot.value)
)

function toggleExpanded(teamId) {
  expandedTeam.value = expandedTeam.value === teamId ? null : teamId
}

function drawCard() {
  if (!pool.value.length || currentLot.value) return
  const idx = Math.floor(Math.random() * pool.value.length)
  currentLot.value = pool.value[idx]
  pool.value = pool.value.filter((_, i) => i !== idx)
  currentBid.value = 0
  currentBidder.value = null
  bidErrors.blue = ''
  bidErrors.red = ''
  currentModifier.value = props.config.modifiersEnabled
    ? rollModifier(modifiersCatalog.value, props.config.modifiersIntensity)
    : null
  broadcastDraftState()
}

function skipCard() {
  pool.value.push(currentLot.value)
  currentLot.value = null
  currentModifier.value = null
  broadcastDraftState()
}

// Mise libre : n'importe quel montant au-dessus de la mise actuelle et dans
// le budget de l'équipe — préréglages (+10K/+20K/+50K) ou montant au choix,
// même règle des deux côtés.
function placeBid(teamId, amount) {
  bidErrors[teamId] = ''
  const n = Math.round(Number(amount))
  if (!n || n <= 0) {
    bidErrors[teamId] = 'Montant invalide.'
    return
  }
  if (n <= currentBid.value) {
    bidErrors[teamId] = `Doit dépasser la mise actuelle (${currentBid.value.toLocaleString()} $).`
    return
  }
  if (n > teams[teamId].budget) {
    bidErrors[teamId] = 'Dépasse le budget restant.'
    return
  }
  currentBid.value = n
  currentBidder.value = teamId
  broadcastDraftState()
}

function raisePreset(teamId, preset) {
  placeBid(teamId, currentBid.value + preset)
}

function raiseCustom(teamId) {
  placeBid(teamId, customBidRefFor(teamId).value)
  customBidRefFor(teamId).value = 0
}

function pass(teamId) {
  const other = teamId === 'blue' ? 'red' : 'blue'
  if (currentBidder.value === other) {
    finalizeLot(other, currentBid.value)
  } else {
    skipCard()
  }
}

function finalizeLot(winnerId, price) {
  teams[winnerId].budget -= price
  const modifier = currentModifier.value
  if (modifier?.kind === 'money') {
    teams[winnerId].budget += modifier.type === 'bonus' ? modifier.amount : -modifier.amount
  }
  if (modifier?.kind === 'slot') {
    teams[winnerId].slotBonus += modifier.type === 'bonus' ? modifier.amount : -modifier.amount
  }
  teams[winnerId].slots.push({ entry: currentLot.value, price, modifier })
  currentLot.value = null
  currentBid.value = 0
  currentBidder.value = null
  currentModifier.value = null
  broadcastDraftState()
}

function fillRandomly(teamId) {
  while (teams[teamId].slots.length < effectiveSlots(teamId) && pool.value.length) {
    const idx = Math.floor(Math.random() * pool.value.length)
    const entry = pool.value[idx]
    pool.value.splice(idx, 1)
    teams[teamId].slots.push({ entry, price: 0, modifier: null })
  }
  broadcastDraftState()
}

function finishDraft() {
  emit('done', {
    hands: {
      blue: teams.blue.slots.map((s) => s.entry),
      red: teams.red.slots.map((s) => s.entry),
    },
  })
}
</script>

<template>
  <div class="p-6">
    <h1 class="text-lg font-medium text-text-primary mb-4">Le Bocal — Enchères</h1>

    <div class="flex gap-4 mb-2">
      <button
        class="flex-1 text-left p-3 rounded-card border cursor-pointer"
        :class="teams.blue.slots.length >= effectiveSlots('blue') ? 'border-border' : 'border-accent'"
        @click="toggleExpanded('blue')"
      >
        <p class="font-medium text-accent">{{ teams.blue.name }}</p>
        <p class="text-sm text-text-secondary">
          {{ teams.blue.budget.toLocaleString() }} $ · {{ teams.blue.slots.length }}/{{ effectiveSlots('blue') }} slots ·
          {{ totalForce('blue') }} FOR / {{ totalPV('blue') }} PV
        </p>
      </button>
      <button
        class="flex-1 text-left p-3 rounded-card border cursor-pointer"
        :class="teams.red.slots.length >= effectiveSlots('red') ? 'border-border' : 'border-force'"
        @click="toggleExpanded('red')"
      >
        <p class="font-medium text-force">{{ teams.red.name }}</p>
        <p class="text-sm text-text-secondary">
          {{ teams.red.budget.toLocaleString() }} $ · {{ teams.red.slots.length }}/{{ effectiveSlots('red') }} slots ·
          {{ totalForce('red') }} FOR / {{ totalPV('red') }} PV
        </p>
      </button>
    </div>

    <div v-if="expandedTeam" class="mb-6 p-3 bg-surface-1 rounded-card">
      <p class="text-sm font-medium mb-2" :class="expandedTeam === 'blue' ? 'text-accent' : 'text-force'">
        Équipe achetée — {{ teams[expandedTeam].name }}
      </p>
      <div v-if="!teams[expandedTeam].slots.length" class="text-sm text-text-muted">Aucun achat pour l'instant.</div>
      <div class="flex flex-col gap-1">
        <div
          v-for="s in teams[expandedTeam].slots"
          :key="s.entry.character.id"
          class="flex items-center justify-between text-sm py-1 border-b border-border last:border-0"
        >
          <span>
            {{ s.entry.character.name.full }}
            <span v-if="s.modifier" class="text-xs text-text-muted">— {{ s.modifier.label }}</span>
          </span>
          <span class="text-text-secondary">{{ s.price.toLocaleString() }} $</span>
        </div>
      </div>
    </div>

    <div v-if="draftDone" class="mb-6">
      <p class="text-text-primary font-medium mb-3">Draft terminé !</p>
      <div class="flex gap-4 mb-4">
        <div class="flex-1 p-3 rounded-card border border-accent">
          <p class="font-medium text-accent">{{ teams.blue.name }}</p>
          <p class="text-sm text-text-secondary">
            {{ teams.blue.slots.length }} personnages · {{ totalForce('blue') }} FOR / {{ totalPV('blue') }} PV
          </p>
        </div>
        <div class="flex-1 p-3 rounded-card border border-force">
          <p class="font-medium text-force">{{ teams.red.name }}</p>
          <p class="text-sm text-text-secondary">
            {{ teams.red.slots.length }} personnages · {{ totalForce('red') }} FOR / {{ totalPV('red') }} PV
          </p>
        </div>
      </div>
      <BaseButton @click="finishDraft">Continuer</BaseButton>
    </div>

    <template v-else>
      <div class="flex items-start justify-center gap-6 my-6">
        <TeamSlotsColumn :slots="teams.blue.slots" :total="effectiveSlots('blue')" tone="accent" />

        <div class="flex flex-col items-center">
          <BaseAspectImage
            class="relative z-0"
            :src="appImage('bocal.png')"
            alt="Le Bocal"
            ratio="498 / 949"
            :width="POT_WIDTH"
          />

          <div class="relative z-10 -mt-10 min-h-[220px] flex items-start justify-center">
            <Transition name="pop">
              <div v-if="currentLot" :key="currentLot.character.id" class="flex flex-col items-center gap-2">
                <CardRenderer
                  :character="currentLot.character"
                  :stats="currentLot.stats"
                  :accent-color="currentLot.accentColor"
                />
                <span
                  v-if="currentModifier"
                  class="text-xs px-2 py-0.5 rounded-full"
                  :class="currentModifier.type === 'bonus' ? 'bg-pv-soft text-pv' : 'bg-force-soft text-force'"
                >
                  {{ currentModifier.type === 'bonus' ? '✦' : '⚠' }} {{ currentModifier.label }}
                </span>
              </div>
            </Transition>
          </div>
        </div>

        <TeamSlotsColumn :slots="teams.red.slots" :total="effectiveSlots('red')" tone="force" />
      </div>

      <div v-if="currentLot" class="flex flex-col items-center gap-3 mb-6">
        <p class="text-lg font-medium text-text-primary">
          Mise actuelle : {{ currentBid.toLocaleString() }} $
          <span v-if="currentBidder" class="text-text-secondary text-sm">({{ teams[currentBidder].name }})</span>
        </p>

        <div v-if="currentBidder !== 'blue'" class="flex flex-col items-center gap-1">
          <div class="flex gap-2">
            <BaseButton
              v-for="preset in BID_PRESETS"
              :key="'blue-' + preset"
              :disabled="!canBidAmount('blue', currentBid + preset)"
              @click="raisePreset('blue', preset)"
            >
              Bleue +{{ preset.toLocaleString() }} $
            </BaseButton>
          </div>
          <div class="flex gap-1 items-center">
            <input
              v-model.number="customBidBlue"
              type="number"
              :min="currentBid + 1"
              :max="teams.blue.budget"
              placeholder="Montant bleu au choix"
              class="w-40 border border-border rounded-card px-2 py-1 text-sm bg-surface-2 text-text-primary"
            />
            <BaseButton variant="secondary" @click="raiseCustom('blue')">Miser ce montant</BaseButton>
          </div>
          <p v-if="bidErrors.blue" class="text-xs text-force">{{ bidErrors.blue }}</p>
        </div>
        <BaseButton v-if="currentBidder === 'blue'" variant="secondary" @click="pass('red')">Rouge passe</BaseButton>

        <template v-if="config.connection !== 'lan'">
          <div v-if="currentBidder !== 'red'" class="flex flex-col items-center gap-1">
            <div class="flex gap-2">
              <BaseButton
                v-for="preset in BID_PRESETS"
                :key="'red-' + preset"
                :disabled="!canBidAmount('red', currentBid + preset)"
                @click="raisePreset('red', preset)"
              >
                Rouge +{{ preset.toLocaleString() }} $
              </BaseButton>
            </div>
            <div class="flex gap-1 items-center">
              <input
                v-model.number="customBidRed"
                type="number"
                :min="currentBid + 1"
                :max="teams.red.budget"
                placeholder="Montant rouge au choix"
                class="w-40 border border-border rounded-card px-2 py-1 text-sm bg-surface-2 text-text-primary"
              />
              <BaseButton variant="secondary" @click="raiseCustom('red')">Miser ce montant</BaseButton>
            </div>
            <p v-if="bidErrors.red" class="text-xs text-force">{{ bidErrors.red }}</p>
          </div>
        </template>
        <p v-else-if="currentBidder !== 'red'" class="text-sm text-text-muted">Équipe Rouge — mise depuis le téléphone de l'invité.</p>
        <BaseButton v-if="currentBidder === 'red'" variant="secondary" @click="pass('blue')">Bleue passe</BaseButton>

        <BaseButton v-if="!currentBidder" variant="ghost" @click="skipCard">Personne ne mise — remettre dans le pot</BaseButton>
      </div>

      <div v-else class="flex justify-center mb-6">
        <BaseButton :disabled="!activeTeams.length || !pool.length" @click="drawCard">
          Piocher une carte du Bocal
        </BaseButton>
      </div>

      <div v-if="strandedTeams.length" class="flex flex-col gap-2">
        <div
          v-for="teamId in strandedTeams"
          :key="teamId"
          class="flex items-center justify-between p-3 bg-surface-1 rounded-card"
        >
          <p class="text-sm" :class="teamId === 'blue' ? 'text-accent' : 'text-force'">
            {{ teams[teamId].name }} n'a plus les moyens de miser —
            {{ effectiveSlots(teamId) - teams[teamId].slots.length }} slot(s) restant(s)
          </p>
          <BaseButton variant="secondary" @click="fillRandomly(teamId)">Remplir au hasard</BaseButton>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.pop-enter-active {
  animation: draw-out 0.75s cubic-bezier(0.3, 0.6, 0.35, 1);
  transform-origin: center;
  position: relative;
  z-index: 20;
}
@keyframes draw-out {
  0% {
    transform: translate(10px, -60px) scale(0.05) rotate(-30deg);
    opacity: 0;
  }
  35% {
    transform: translate(65px, -10px) scale(0.5) rotate(220deg);
    opacity: 1;
  }
  65% {
    transform: translate(-30px, 25px) scale(0.85) rotate(370deg);
  }
  100% {
    transform: translate(0, 0) scale(1) rotate(360deg);
    opacity: 1;
  }
}
.pop-leave-active {
  position: absolute;
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.pop-leave-to {
  transform: scale(0.85);
  opacity: 0;
}
</style>
