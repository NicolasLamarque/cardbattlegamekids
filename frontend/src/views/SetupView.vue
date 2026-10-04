<script setup>
import { ref, computed } from 'vue'
import BaseButton from '../components/base/BaseButton.vue'
import BasePanel from '../components/base/BasePanel.vue'
import BaseImage from '../components/base/BaseImage.vue'
import { appImage } from '../data/appImages.js'
import { locations } from '../data/locations.js'
import { StartLanServer } from '../../wailsjs/go/main/App.js'

const props = defineProps({
  mode: { type: String, required: true },
})
const emit = defineEmits(['launch'])

const slots = ref(10)
const budgetChoice = ref(500000)
const customBudget = ref(500000)
const location = ref('tokyo')
const combatMode = ref('blitz')
const specialCardsEnabled = ref(true)
const specialCardsPerPlayer = ref(5)
const modifiersEnabled = ref(true)
const modifiersIntensity = ref('rare')
const connection = ref('solo')
const lanInfo = ref(null)
const lanStarting = ref(false)

async function startLan() {
  lanStarting.value = true
  try {
    lanInfo.value = await StartLanServer()
  } finally {
    lanStarting.value = false
  }
}

const startingBudget = computed(() =>
  budgetChoice.value === 'custom' ? Number(customBudget.value) || 0 : budgetChoice.value
)

function launch() {
  emit('launch', {
    mode: props.mode,
    players: 2,
    slots: slots.value,
    startingBudget: startingBudget.value,
    location: location.value,
    combatMode: combatMode.value,
    specialCardsEnabled: specialCardsEnabled.value,
    specialCardsPerPlayer: specialCardsPerPlayer.value,
    modifiersEnabled: modifiersEnabled.value,
    modifiersIntensity: modifiersIntensity.value,
    connection: connection.value,
  })
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-6">
    <BasePanel class="w-full max-w-md flex flex-col gap-4">
      <h1 class="text-lg font-medium text-text-primary">
        Paramétrage — {{ mode === 'career' ? 'Carrière' : 'Escarmouche' }}
      </h1>

      <div>
        <label class="text-sm text-text-secondary block mb-1">Joueurs</label>
        <select class="w-full border border-border rounded-card px-2 py-1.5 bg-surface-2 text-text-primary">
          <option>2 joueurs</option>
          <option disabled>3-4 joueurs (bientôt)</option>
        </select>
      </div>

      <div>
        <label class="text-sm text-text-secondary block mb-1">Slots par équipe (personnages à acheter)</label>
        <select
          v-model.number="slots"
          class="w-full border border-border rounded-card px-2 py-1.5 bg-surface-2 text-text-primary"
        >
          <option :value="5">5</option>
          <option :value="10">10</option>
          <option :value="15">15</option>
          <option :value="20">20</option>
        </select>
      </div>

      <div>
        <label class="text-sm text-text-secondary block mb-1">Montant de départ</label>
        <select
          v-model="budgetChoice"
          class="w-full border border-border rounded-card px-2 py-1.5 bg-surface-2 text-text-primary mb-2"
        >
          <option :value="100000">100 000 $</option>
          <option :value="500000">500 000 $</option>
          <option :value="1000000">1 000 000 $</option>
          <option value="custom">Montant personnalisé</option>
        </select>
        <input
          v-if="budgetChoice === 'custom'"
          v-model.number="customBudget"
          type="number"
          min="0"
          step="10000"
          class="w-full border border-border rounded-card px-2 py-1.5 bg-surface-2 text-text-primary"
        />
      </div>

      <div>
        <label class="text-sm text-text-secondary block mb-1">Lieu du combat</label>
        <div class="flex gap-2">
          <div v-for="loc in locations" :key="loc.id" class="flex-1">
            <button
              class="w-full block rounded-card overflow-hidden border"
              :class="location === loc.id ? 'border-accent' : 'border-border'"
              @click="location = loc.id"
            >
              <div class="h-32"><BaseImage :src="appImage(loc.image)" :alt="loc.name" /></div>
            </button>
            <p class="text-xs py-1 text-center text-text-secondary">{{ loc.name }}</p>
          </div>
        </div>
      </div>

      <div>
        <label class="text-sm text-text-secondary block mb-1">Mode de combat</label>
        <select
          v-model="combatMode"
          class="w-full border border-border rounded-card px-2 py-1.5 bg-surface-2 text-text-primary"
        >
          <option value="blitz">Blitz (rapide)</option>
          <option value="slow">Lent</option>
          <option value="timed">Chronométré</option>
        </select>
      </div>

      <label class="flex items-center gap-2 text-sm text-text-primary">
        <input type="checkbox" v-model="specialCardsEnabled" />
        Jouer avec les cartes spéciales
      </label>
      <div v-if="specialCardsEnabled">
        <label class="text-sm text-text-secondary block mb-1">Cartes spéciales par joueur</label>
        <input
          v-model.number="specialCardsPerPlayer"
          type="number"
          min="1"
          max="20"
          class="w-full border border-border rounded-card px-2 py-1.5 bg-surface-2 text-text-primary"
        />
      </div>

      <label class="flex items-center gap-2 text-sm text-text-primary">
        <input type="checkbox" v-model="modifiersEnabled" />
        Cartes piégées / bonus dans le Bocal
      </label>
      <div v-if="modifiersEnabled">
        <label class="text-sm text-text-secondary block mb-1">Intensité</label>
        <select
          v-model="modifiersIntensity"
          class="w-full border border-border rounded-card px-2 py-1.5 bg-surface-2 text-text-primary"
        >
          <option value="rare">Rare — surprise occasionnelle</option>
          <option value="frequent">Fréquent — presque à chaque carte</option>
        </select>
      </div>

      <div>
        <label class="text-sm text-text-secondary block mb-1">Connexion</label>
        <select
          v-model="connection"
          class="w-full border border-border rounded-card px-2 py-1.5 bg-surface-2 text-text-primary"
        >
          <option value="solo">Même écran</option>
          <option value="lan">En réseau local</option>
        </select>

        <div v-if="connection === 'lan'" class="mt-2">
          <BaseButton v-if="!lanInfo" variant="secondary" :disabled="lanStarting" @click="startLan">
            {{ lanStarting ? 'Démarrage…' : 'Démarrer le serveur local' }}
          </BaseButton>
          <div v-else class="flex flex-col items-center gap-2 p-3 bg-surface-2 rounded-card">
            <img :src="lanInfo.qrCode" alt="Code QR de connexion" class="w-36 h-36" />
            <p class="text-xs text-text-secondary break-all">{{ lanInfo.url }}</p>
            <p class="text-xs text-text-muted text-center">
              À scanner depuis le téléphone du deuxième joueur, sur le même réseau WiFi.
            </p>
          </div>
        </div>
      </div>

      <BaseButton @click="launch">Démarrer la partie</BaseButton>
    </BasePanel>
  </div>
</template>
