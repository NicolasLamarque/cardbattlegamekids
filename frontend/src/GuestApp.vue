<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import GuestArenaView from './views/GuestArenaView.vue'
import GuestDraftView from './views/GuestDraftView.vue'

const state = ref(null)
const connected = ref(true)
// Un téléphone qui se verrouille suspend sa page et ferme le WebSocket —
// inévitable côté navigateur mobile. On ne peut pas empêcher la coupure,
// seulement la réparer vite : on retente automatiquement, avec un délai qui
// grandit à chaque échec, et on force une tentative immédiate dès que la
// page redevient visible (déverrouillage du téléphone).
let ws = null
let reconnectTimer = null
let reconnectDelay = 1000
const RECONNECT_DELAY_MAX = 8000

function connect() {
  clearTimeout(reconnectTimer)
  ws = new WebSocket(`ws://${location.host}/ws`)
  ws.onopen = () => {
    connected.value = true
    reconnectDelay = 1000
  }
  ws.onmessage = (event) => {
    state.value = JSON.parse(event.data)
  }
  ws.onclose = scheduleReconnect
  ws.onerror = scheduleReconnect
}

function scheduleReconnect() {
  connected.value = false
  clearTimeout(reconnectTimer)
  reconnectTimer = setTimeout(connect, reconnectDelay)
  reconnectDelay = Math.min(reconnectDelay * 2, RECONNECT_DELAY_MAX)
}

function onVisibilityChange() {
  if (document.visibilityState === 'visible' && (!ws || ws.readyState !== WebSocket.OPEN)) {
    reconnectDelay = 1000
    connect()
  }
}

onMounted(() => {
  connect()
  document.addEventListener('visibilitychange', onVisibilityChange)
})

onUnmounted(() => {
  document.removeEventListener('visibilitychange', onVisibilityChange)
  clearTimeout(reconnectTimer)
  ws?.close()
})

const actionFailed = ref(false)
let actionFailedTimer = null

// Si le socket n'est pas ouvert (reconnexion en cours), l'action ne part
// pas — mieux vaut le dire tout de suite que laisser le joueur croire que
// son tap a compté.
function sendAction(action) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(action))
    return
  }
  actionFailed.value = true
  clearTimeout(actionFailedTimer)
  actionFailedTimer = setTimeout(() => {
    actionFailed.value = false
  }, 2500)
}
</script>

<template>
  <div v-if="!connected" class="min-h-screen flex items-center justify-center p-6">
    <p class="text-force text-center">
      Connexion perdue — reconnexion automatique en cours…
      <br />
      <span class="text-xs text-text-muted">
        Si ça persiste, l'hôte a peut-être fermé la partie — redémarre le serveur et rescanne le code QR.
      </span>
    </p>
  </div>
  <div v-else-if="!state" class="min-h-screen flex items-center justify-center p-6">
    <p class="text-text-secondary text-center">Connecté — en attente de l'hôte…</p>
  </div>
  <GuestDraftView v-else-if="state.type === 'draftState'" :state="state" @action="sendAction" />
  <GuestArenaView v-else-if="state.type === 'arenaState'" :state="state" @action="sendAction" />

  <div v-if="actionFailed" class="fixed bottom-4 left-1/2 -translate-x-1/2 bg-force text-white text-sm px-3 py-1.5 rounded-card shadow-lg">
    Action non envoyée — reconnexion en cours…
  </div>
</template>
