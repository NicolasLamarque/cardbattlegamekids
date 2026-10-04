<script setup>
defineProps({
  color: { type: String, default: null },
  active: { type: Boolean, default: false },
})
</script>

<template>
  <Transition name="attack-burst">
    <div v-if="active && color" class="absolute inset-0 pointer-events-none flex items-center justify-center overflow-hidden" :style="{ color }">
      <svg viewBox="0 0 160 256" class="w-full h-full swirl-group">
        <path
          class="swirl-line swirl-line-1"
          d="M10,128 C 50,20 110,236 150,128"
          stroke="currentColor"
          stroke-width="5"
          fill="none"
          stroke-linecap="round"
        />
        <path
          class="swirl-line swirl-line-2"
          d="M10,128 C 50,236 110,20 150,128"
          stroke="currentColor"
          stroke-width="5"
          fill="none"
          stroke-linecap="round"
        />
        <polyline
          class="bolt bolt-1"
          points="34,8 58,96 26,104 74,214"
          stroke="currentColor"
          stroke-width="3"
          fill="none"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
        <polyline
          class="bolt bolt-2"
          points="126,16 96,112 134,122 86,224"
          stroke="currentColor"
          stroke-width="3"
          fill="none"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
      <div class="glow"></div>
    </div>
  </Transition>
</template>

<style scoped>
.attack-burst-enter-active {
  animation: burst-fade-in 0.15s ease-out;
}
.attack-burst-leave-active {
  animation: burst-fade-in 0.2s ease-in reverse;
}

.swirl-group {
  filter: drop-shadow(0 0 6px currentColor);
}

.swirl-line {
  transform-origin: 80px 128px;
  animation: swirl-spin 0.9s cubic-bezier(0.2, 0.8, 0.3, 1) 1, swirl-pulse 1.6s ease-in-out infinite;
  opacity: 0.85;
}
.swirl-line-2 {
  animation-direction: reverse, alternate;
}

.bolt {
  opacity: 0;
  animation: bolt-flicker 0.9s steps(1, end) 1;
}
.bolt-2 {
  animation-delay: 0.08s;
}

.glow {
  position: absolute;
  inset: 0;
  background: radial-gradient(circle, currentColor 0%, transparent 65%);
  opacity: 0.18;
  animation: swirl-pulse 1.6s ease-in-out infinite;
}

@keyframes burst-fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes swirl-spin {
  from {
    transform: rotate(0deg) scale(0.3);
  }
  60% {
    transform: rotate(300deg) scale(1.15);
  }
  to {
    transform: rotate(360deg) scale(1);
  }
}

@keyframes swirl-pulse {
  0%, 100% {
    opacity: 0.6;
  }
  50% {
    opacity: 1;
  }
}

@keyframes bolt-flicker {
  0%, 100% {
    opacity: 0;
  }
  10%, 30%, 50% {
    opacity: 1;
  }
  20%, 40%, 60% {
    opacity: 0;
  }
  70% {
    opacity: 0;
  }
}
</style>
