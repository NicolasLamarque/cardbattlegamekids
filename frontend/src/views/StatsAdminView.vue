<script setup>
import { ref, computed, onMounted } from 'vue'
import BaseImage from '../components/base/BaseImage.vue'
import BaseButton from '../components/base/BaseButton.vue'
import { localCharacterImage } from '../data/images.js'
import { GetCharacters, UpdateCharacterStats, ResetCharacterStats } from '../../wailsjs/go/main/App.js'

const rows = ref([])
const search = ref('')

onMounted(load)

async function load() {
  const data = await GetCharacters()
  rows.value = data.map((c) => ({ ...c, draftForce: c.force, draftPV: c.pv }))
}

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return rows.value
  return rows.value.filter((c) => c.nameFull.toLowerCase().includes(q))
})

async function save(row) {
  await UpdateCharacterStats(row.id, Number(row.draftForce), Number(row.draftPV))
  row.force = Number(row.draftForce)
  row.pv = Number(row.draftPV)
  row.isCustomStats = true
}

async function reset(row) {
  await ResetCharacterStats(row.id)
  await load()
}
</script>

<template>
  <div class="p-6">
    <h1 class="text-lg font-medium text-text-primary mb-1">Gestion des stats</h1>
    <p class="text-sm text-text-secondary mb-4">
      Force et PV sont calculés depuis les favoris AniList. Modifier une valeur ici la fige en "custom" — elle
      ne sera plus jamais recalculée automatiquement, sauf réinitialisation.
    </p>

    <input
      v-model="search"
      type="text"
      placeholder="Rechercher un personnage…"
      class="w-full max-w-sm border border-border rounded-card px-3 py-1.5 mb-4 bg-surface-2 text-text-primary"
    />

    <div class="overflow-x-auto">
      <table class="w-full text-sm border-collapse">
        <thead>
          <tr class="text-left text-text-secondary border-b border-border">
            <th class="py-2 pr-3"></th>
            <th class="py-2 pr-3">Nom</th>
            <th class="py-2 pr-3">Série</th>
            <th class="py-2 pr-3">Favoris</th>
            <th class="py-2 pr-3">Force</th>
            <th class="py-2 pr-3">PV</th>
            <th class="py-2 pr-3">Source</th>
            <th class="py-2 pr-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in filtered" :key="row.id" class="border-b border-border">
            <td class="py-2 pr-3">
              <div class="w-10 h-10 rounded-card overflow-hidden">
                <BaseImage :src="localCharacterImage(row.imageLocal) || row.imageLarge" :alt="row.nameFull" />
              </div>
            </td>
            <td class="py-2 pr-3 text-text-primary">{{ row.nameFull }}</td>
            <td class="py-2 pr-3 text-text-secondary">{{ row.seriesTitle }}</td>
            <td class="py-2 pr-3 text-text-secondary">{{ row.favourites.toLocaleString() }}</td>
            <td class="py-2 pr-3">
              <input
                v-model.number="row.draftForce"
                type="number"
                class="w-20 border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              />
            </td>
            <td class="py-2 pr-3">
              <input
                v-model.number="row.draftPV"
                type="number"
                class="w-20 border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              />
            </td>
            <td class="py-2 pr-3">
              <span
                class="text-xs px-2 py-0.5 rounded-full"
                :class="row.isCustomStats ? 'bg-gold-soft text-gold' : 'bg-surface-2 text-text-muted'"
              >
                {{ row.isCustomStats ? 'Custom' : 'Auto' }}
              </span>
            </td>
            <td class="py-2 pr-3">
              <div class="flex gap-2">
                <BaseButton
                  variant="secondary"
                  :disabled="row.draftForce === row.force && row.draftPV === row.pv"
                  @click="save(row)"
                >
                  Enregistrer
                </BaseButton>
                <BaseButton v-if="row.isCustomStats" variant="ghost" @click="reset(row)">Réinitialiser</BaseButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
