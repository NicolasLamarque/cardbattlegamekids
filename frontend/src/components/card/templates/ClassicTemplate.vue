<script setup>
import CardFrame from '../parts/CardFrame.vue'
import CardArt from '../parts/CardArt.vue'
import CardName from '../parts/CardName.vue'
import CardStats from '../parts/CardStats.vue'

defineProps({
  character: { type: Object, required: true },
  stats: { type: Object, required: true },
  accentColor: { type: String, default: '' },
})
</script>

<template>
  <CardFrame :accent-color="accentColor">
    <CardArt :src="character.image?.large" :alt="character.name?.full" />
    <CardName :name="character.name?.full" />
    <CardStats :force="stats.force" :pv="stats.pv" />
  </CardFrame>
</template>
