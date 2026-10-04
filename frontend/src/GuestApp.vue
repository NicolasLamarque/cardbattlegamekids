<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import GuestArenaView from './views/GuestArenaView.vue'

const state = ref(null)
const connected = ref(true)
let ws = null

onMounted(() => {
  ws = new WebSocket(`ws://${location.host}/ws`)
  ws.onmessage = (event) => {
    state.value = JSON.parse(event.data)
  }
  ws.onclose = () => {
    connected.value = false
  }
  ws.onerror = () => {
    connected.value = false
  }
})

onUnmounted(() => {
  ws?.close()
})

function sendAction(action) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(action))
  }
}
</script>

<template>
  <div v-if="!connected" class="min-h-screen flex items-center justify-center p-6">
    <p class="text-force text-center">
      Connexion perdue — l'hôte a fermé la partie ou son appareil a redémarré.
      <br />Redémarre le serveur côté hôte et rescanne le code QR.
    </p>
  </div>
  <div v-else-if="!state" class="min-h-screen flex items-center justify-center p-6">
    <p class="text-text-secondary text-center">Connecté — en attente que l'hôte lance le combat…</p>
  </div>
  <GuestArenaView v-else-if="state.type === 'arenaState'" :state="state" @action="sendAction" />
</template>
