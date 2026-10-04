<script setup>
import BaseStatBar from '../../base/BaseStatBar.vue'

defineProps({
  force: { type: Number, required: true },
  pv: { type: Number, required: true },
  maxForce: { type: Number, default: 100 },
  maxPv: { type: Number, default: 300 },
})
</script>

<template>
  <div class="px-2 pb-2 pt-1 flex flex-col gap-1">
    <BaseStatBar label="FOR" :value="force" :max="maxForce" tone="force" />
    <BaseStatBar label="PV" :value="pv" :max="maxPv" tone="pv" />
  </div>
</template>
