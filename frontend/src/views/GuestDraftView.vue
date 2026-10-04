<script setup>
import { ref } from 'vue'
import CardRenderer from '../components/card/CardRenderer.vue'
import BaseButton from '../components/base/BaseButton.vue'

const props = defineProps({
  state: { type: Object, required: true },
})
const emit = defineEmits(['action'])

const customAmount = ref(0)

const BID_PRESETS = [10000, 20000, 50000]

function bid(amount) {
  emit('action', { type: 'bid', amount })
}

function bidPreset(preset) {
  bid(props.state.currentBid + preset)
}

function bidCustom() {
  bid(Number(customAmount.value))
  customAmount.value = 0
}

function pass() {
  emit('action', { type: 'pass' })
}
</script>

<template>
  <div class="p-6 min-h-screen flex flex-col items-center gap-4">
    <h1 class="text-lg font-medium text-text-primary">Le Bocal — ta mise</h1>

    <p class="text-sm text-force font-medium">
      {{ state.red.budget.toLocaleString() }} $ · {{ state.red.slotsUsed }}/{{ state.red.slotsTotal }} slots
    </p>

    <div v-if="state.draftDone" class="text-center">
      <p class="text-text-primary font-medium">Draft terminé — reviens voir l'hôte pour la suite.</p>
    </div>

    <template v-else-if="state.currentLot">
      <CardRenderer
        :character="state.currentLot.character"
        :stats="state.currentLot.stats"
        :accent-color="state.currentLot.accentColor"
      />
      <span
        v-if="state.currentModifier"
        class="text-xs px-2 py-0.5 rounded-full"
        :class="state.currentModifier.type === 'bonus' ? 'bg-pv-soft text-pv' : 'bg-force-soft text-force'"
      >
        {{ state.currentModifier.type === 'bonus' ? '✦' : '⚠' }} {{ state.currentModifier.label }}
      </span>

      <p class="text-base text-text-primary">
        Mise actuelle : {{ state.currentBid.toLocaleString() }} $
        <span v-if="state.currentBidder" class="text-text-secondary text-sm">
          ({{ state.currentBidder === 'red' ? 'toi' : 'Équipe Bleue' }})
        </span>
      </p>

      <template v-if="state.currentBidder !== 'red'">
        <div class="flex gap-2">
          <BaseButton v-for="preset in BID_PRESETS" :key="preset" @click="bidPreset(preset)">
            +{{ preset.toLocaleString() }} $
          </BaseButton>
        </div>
        <div class="flex gap-1 items-center">
          <input
            v-model.number="customAmount"
            type="number"
            placeholder="Montant au choix"
            class="w-40 border border-border rounded-card px-2 py-1 text-sm bg-surface-2 text-text-primary"
          />
          <BaseButton variant="secondary" @click="bidCustom">Miser</BaseButton>
        </div>
        <p v-if="state.bidError" class="text-xs text-force">{{ state.bidError }}</p>
      </template>
      <BaseButton v-else variant="secondary" @click="pass">Passer</BaseButton>
    </template>

    <p v-else class="text-text-secondary text-sm text-center">
      L'hôte pioche la prochaine carte du Bocal…
    </p>
  </div>
</template>
