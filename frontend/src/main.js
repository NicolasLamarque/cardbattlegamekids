import { createApp } from 'vue'
import App from './App.vue'
import GuestApp from './GuestApp.vue'
import { isGuest } from './network/session.js'
import './style.css'

createApp(isGuest ? GuestApp : App).mount('#app')
