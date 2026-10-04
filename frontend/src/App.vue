<script setup>
import { ref, onMounted } from 'vue'
import CoverView from './views/CoverView.vue'
import ModeSelectView from './views/ModeSelectView.vue'
import SetupView from './views/SetupView.vue'
import DraftView from './views/DraftView.vue'
import CardsDealtView from './views/CardsDealtView.vue'
import CollectionView from './views/CollectionView.vue'
import ArenaView from './views/ArenaView.vue'
import StatsAdminView from './views/StatsAdminView.vue'
import ModifiersAdminView from './views/ModifiersAdminView.vue'
import AttacksAdminView from './views/AttacksAdminView.vue'
import BaseBackButton from './components/base/BaseBackButton.vue'
import BaseToast from './components/base/BaseToast.vue'
import { specialCards, dealHands } from './game/specialCards.js'
import { EventsOn } from '../wailsjs/runtime/runtime.js'

const stage = ref('cover')
const selectedMode = ref('')
const gameConfig = ref(null)
const draftHands = ref(null)
const specialHands = ref(null)
const view = ref('collection')
const toastMessage = ref('')
const toastTone = ref('accent')
let toastTimer = null

function showToast(message, tone) {
  toastMessage.value = message
  toastTone.value = tone
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    toastMessage.value = ''
  }, 3000)
}

onMounted(() => {
  EventsOn('guest:connected', () => showToast('Joueur invité connecté !', 'accent'))
  EventsOn('guest:disconnected', () => showToast('Joueur invité déconnecté', 'force'))
})

const previousStage = {
  modeSelect: 'cover',
  setup: 'modeSelect',
  draft: 'setup',
  cardsDealt: 'draft',
  adminStandalone: 'modeSelect',
}

function goBack() {
  stage.value = previousStage[stage.value] ?? 'cover'
}

function selectMode(mode) {
  selectedMode.value = mode
  stage.value = 'setup'
}

function launchSetup(config) {
  gameConfig.value = config
  stage.value = 'draft'
}

function finishDraft({ hands }) {
  draftHands.value = hands
  if (gameConfig.value.specialCardsEnabled) {
    specialHands.value = dealHands(specialCards, ['blue', 'red'], gameConfig.value.specialCardsPerPlayer)
    stage.value = 'cardsDealt'
  } else {
    stage.value = 'game'
  }
}

function startGame() {
  stage.value = 'game'
}
</script>

<template>
  <BaseToast :message="toastMessage" :tone="toastTone" />

  <CoverView v-if="stage === 'cover'" @enter="stage = 'modeSelect'" />

  <div v-else-if="stage === 'modeSelect'">
    <BaseBackButton @click="goBack" />
    <ModeSelectView @select="selectMode" @admin="stage = 'adminStandalone'; view = 'stats'" />
  </div>

  <div v-else-if="stage === 'setup'">
    <BaseBackButton @click="goBack" />
    <SetupView :mode="selectedMode" @launch="launchSetup" />
  </div>

  <div v-else-if="stage === 'draft'">
    <BaseBackButton @click="goBack" />
    <DraftView :config="gameConfig" @done="finishDraft" />
  </div>

  <div v-else-if="stage === 'cardsDealt'">
    <BaseBackButton @click="goBack" />
    <CardsDealtView :hands="specialHands" @start="startGame" />
  </div>

  <div v-else-if="stage === 'adminStandalone'">
    <BaseBackButton @click="goBack" />
    <nav class="flex gap-2 px-4 pb-4">
      <button
        class="px-3 py-1.5 rounded-card text-sm font-medium"
        :class="view === 'stats' ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-2'"
        @click="view = 'stats'"
      >
        Stats des personnages
      </button>
      <button
        class="px-3 py-1.5 rounded-card text-sm font-medium"
        :class="view === 'modifiers' ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-2'"
        @click="view = 'modifiers'"
      >
        Malus / Bonus du Bocal
      </button>
      <button
        class="px-3 py-1.5 rounded-card text-sm font-medium"
        :class="view === 'attacks' ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-2'"
        @click="view = 'attacks'"
      >
        Table des attaques
      </button>
    </nav>
    <StatsAdminView v-if="view === 'stats'" />
    <AttacksAdminView v-else-if="view === 'attacks'" />
    <ModifiersAdminView v-else />
  </div>

  <div v-else>
    <nav class="flex items-center gap-2 p-4 border-b border-border">
      <button
        class="px-3 py-1.5 rounded-card text-sm font-medium"
        :class="view === 'collection' ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-2'"
        @click="view = 'collection'"
      >
        Collection
      </button>
      <button
        class="px-3 py-1.5 rounded-card text-sm font-medium"
        :class="view === 'arena' ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-2'"
        @click="view = 'arena'"
      >
        Arène
      </button>
      <button
        class="px-3 py-1.5 rounded-card text-sm font-medium"
        :class="view === 'stats' ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-2'"
        @click="view = 'stats'"
      >
        Stats
      </button>
      <button
        class="px-3 py-1.5 rounded-card text-sm font-medium"
        :class="view === 'modifiers' ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-2'"
        @click="view = 'modifiers'"
      >
        Malus / Bonus
      </button>
      <button
        class="px-3 py-1.5 rounded-card text-sm font-medium"
        :class="view === 'attacks' ? 'bg-accent text-white' : 'text-text-secondary hover:bg-surface-2'"
        @click="view = 'attacks'"
      >
        Table des attaques
      </button>
    </nav>
    <CollectionView v-if="view === 'collection'" />
    <StatsAdminView v-else-if="view === 'stats'" />
    <ModifiersAdminView v-else-if="view === 'modifiers'" />
    <AttacksAdminView v-else-if="view === 'attacks'" />
    <ArenaView
      v-else
      :special-cards-enabled="gameConfig?.specialCardsEnabled ?? true"
      :hand-blue="draftHands?.blue ?? []"
      :hand-red="draftHands?.red ?? []"
      :location="gameConfig?.location ?? 'tokyo'"
      :connection="gameConfig?.connection ?? 'solo'"
    />
  </div>
</template>
