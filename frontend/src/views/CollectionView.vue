<script setup>
import { ref, onMounted } from 'vue'
import CardRenderer from '../components/card/CardRenderer.vue'
import { localCharacterImage } from '../data/images.js'
import { GetCharacters } from '../../wailsjs/go/main/App.js'

const characters = ref([])

onMounted(async () => {
  const rows = await GetCharacters()
  characters.value = rows.map((c) => ({
    character: {
      id: c.id,
      name: { full: c.nameFull, native: c.nameNative },
      image: { large: localCharacterImage(c.imageLocal) || c.imageLarge },
      favourites: c.favourites,
    },
    stats: { force: c.force, pv: c.pv },
    accentColor: c.seriesColor,
  }))
})
</script>

<template>
  <div class="p-6">
    <h1 class="text-lg font-medium text-text-primary mb-4">Collection — Jujutsu Kaisen</h1>
    <div class="flex flex-wrap gap-4">
      <CardRenderer
        v-for="entry in characters"
        :key="entry.character.id"
        :character="entry.character"
        :stats="entry.stats"
        :accent-color="entry.accentColor"
      />
    </div>
  </div>
</template>
