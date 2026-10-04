<script setup>
import { ref, onMounted } from 'vue'
import BaseButton from '../components/base/BaseButton.vue'
import { GetDraftModifiers, SaveDraftModifier, DeleteDraftModifier } from '../../wailsjs/go/main/App.js'

const rows = ref([])

onMounted(load)

async function load() {
  rows.value = await GetDraftModifiers()
}

async function save(row) {
  await SaveDraftModifier({
    id: row.id,
    label: row.label,
    type: row.type,
    kind: row.kind,
    amount: Number(row.amount),
    rare: !!row.rare,
  })
  await load()
}

async function remove(row) {
  await DeleteDraftModifier(row.id)
  await load()
}

function addNew() {
  rows.value.push({
    id: `custom-${Date.now()}`,
    label: 'Nouveau modificateur',
    type: 'bonus',
    kind: 'money',
    amount: 25000,
    rare: false,
  })
}
</script>

<template>
  <div class="p-6">
    <h1 class="text-lg font-medium text-text-primary mb-1">Malus / Bonus du Bocal</h1>
    <p class="text-sm text-text-secondary mb-4">
      Ce qui peut apparaître sur une carte tirée du Bocal. "Argent" ajuste le budget de l'équipe qui remporte
      la carte, "Slot" ajuste son nombre de personnages à acheter pour cette partie.
    </p>

    <div class="overflow-x-auto">
      <table class="w-full text-sm border-collapse">
        <thead>
          <tr class="text-left text-text-secondary border-b border-border">
            <th class="py-2 pr-3">Nom</th>
            <th class="py-2 pr-3">Type</th>
            <th class="py-2 pr-3">Genre</th>
            <th class="py-2 pr-3">Montant</th>
            <th class="py-2 pr-3">Rare</th>
            <th class="py-2 pr-3"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-b border-border">
            <td class="py-2 pr-3">
              <input
                v-model="row.label"
                type="text"
                class="w-40 border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              />
            </td>
            <td class="py-2 pr-3">
              <select v-model="row.type" class="border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary">
                <option value="bonus">Bonus</option>
                <option value="malus">Malus</option>
              </select>
            </td>
            <td class="py-2 pr-3">
              <select v-model="row.kind" class="border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary">
                <option value="money">Argent</option>
                <option value="slot">Slot</option>
              </select>
            </td>
            <td class="py-2 pr-3">
              <input
                v-model.number="row.amount"
                type="number"
                class="w-24 border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              />
            </td>
            <td class="py-2 pr-3">
              <input type="checkbox" v-model="row.rare" />
            </td>
            <td class="py-2 pr-3">
              <div class="flex gap-2">
                <BaseButton variant="secondary" @click="save(row)">Enregistrer</BaseButton>
                <BaseButton variant="ghost" @click="remove(row)">Supprimer</BaseButton>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <BaseButton class="mt-4" variant="secondary" @click="addNew">+ Ajouter un modificateur</BaseButton>
  </div>
</template>
