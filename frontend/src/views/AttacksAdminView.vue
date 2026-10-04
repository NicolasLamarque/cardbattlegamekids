<script setup>
import { ref, computed, onMounted } from 'vue'
import BaseImage from '../components/base/BaseImage.vue'
import BaseButton from '../components/base/BaseButton.vue'
import { localCharacterImage } from '../data/images.js'
import {
  GetCharacters,
  GetCharacterAttacks,
  SaveCharacterAttack,
  DeleteCharacterAttack,
  ResetCharacterAttacks,
  UpdateCharacterManaSettings,
} from '../../wailsjs/go/main/App.js'

const characters = ref([])
const attacksByCharacter = ref({})
const search = ref('')
const expanded = ref(null)

onMounted(load)

async function load() {
  const [chars, attacks] = await Promise.all([GetCharacters(), GetCharacterAttacks()])
  characters.value = chars.map((c) => ({
    ...c,
    draftManaRegen: c.manaRegen,
    draftMaxAttacksOverride: c.maxAttacksOverride ?? 'auto',
  }))
  const grouped = {}
  for (const a of attacks) {
    ;(grouped[a.characterId] ??= []).push({ ...a })
  }
  attacksByCharacter.value = grouped
}

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  if (!q) return characters.value
  return characters.value.filter((c) => c.nameFull.toLowerCase().includes(q))
})

function attacksFor(characterId) {
  return attacksByCharacter.value[characterId] ?? []
}

function toggleExpanded(id) {
  expanded.value = expanded.value === id ? null : id
}

async function saveManaSettings(row) {
  const override = row.draftMaxAttacksOverride === 'auto' ? -1 : Number(row.draftMaxAttacksOverride)
  await UpdateCharacterManaSettings(row.id, Number(row.draftManaRegen), override)
  await load()
}

async function saveAttack(row, attack) {
  const saved = await SaveCharacterAttack({
    id: attack.id ?? 0,
    characterId: row.id,
    name: attack.name,
    color: attack.color,
    forceBonus: Number(attack.forceBonus),
    manaCost: Number(attack.manaCost),
    sortOrder: Number(attack.sortOrder ?? 0),
  })
  attack.id = saved.id
}

async function removeAttack(row, attack) {
  if (attack.id) await DeleteCharacterAttack(attack.id)
  attacksByCharacter.value[row.id] = attacksFor(row.id).filter((a) => a !== attack)
}

function addAttack(row) {
  const list = attacksFor(row.id)
  if (list.length >= row.maxAttacks) return
  list.push({ id: 0, characterId: row.id, name: 'Nouvelle attaque', color: '#f87171', forceBonus: 5, manaCost: 2, sortOrder: list.length })
  attacksByCharacter.value[row.id] = list
}

async function resetAttacks(row) {
  await ResetCharacterAttacks(row.id)
  await load()
}
</script>

<template>
  <div class="p-6">
    <h1 class="text-lg font-medium text-text-primary mb-1">Table des attaques</h1>
    <p class="text-sm text-text-secondary mb-4">
      Chaque personnage peut lier une attaque à sa carte au combat, en échange de mana. Le nombre d'attaques max
      et la jauge de mana sont calculés depuis sa valeur — réglables ici si besoin.
    </p>

    <input
      v-model="search"
      type="text"
      placeholder="Rechercher un personnage…"
      class="w-full max-w-sm border border-border rounded-card px-3 py-1.5 mb-4 bg-surface-2 text-text-primary"
    />

    <div class="flex flex-col gap-2">
      <div v-for="row in filtered" :key="row.id" class="border border-border rounded-card">
        <button class="w-full flex items-center gap-3 p-3 text-left cursor-pointer" @click="toggleExpanded(row.id)">
          <div class="w-10 h-10 rounded-card overflow-hidden shrink-0">
            <BaseImage :src="localCharacterImage(row.imageLocal) || row.imageLarge" :alt="row.nameFull" />
          </div>
          <div class="flex-1">
            <p class="text-text-primary">{{ row.nameFull }}</p>
            <p class="text-xs text-text-secondary">
              Mana max {{ row.maxMana }} · {{ attacksFor(row.id).length }}/{{ row.maxAttacks }} attaques
            </p>
          </div>
        </button>

        <div v-if="expanded === row.id" class="p-3 border-t border-border bg-surface-1 flex flex-col gap-4">
          <div class="flex items-end gap-4">
            <div>
              <label class="text-xs text-text-secondary block mb-1">Recharge de mana / manche</label>
              <input
                v-model.number="row.draftManaRegen"
                type="number"
                min="0"
                class="w-24 border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              />
            </div>
            <div>
              <label class="text-xs text-text-secondary block mb-1">Nombre max d'attaques</label>
              <select
                v-model="row.draftMaxAttacksOverride"
                class="border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              >
                <option value="auto">Auto (calcul via la valeur)</option>
                <option v-for="n in 5" :key="n" :value="n">{{ n }}</option>
              </select>
            </div>
            <BaseButton variant="secondary" @click="saveManaSettings(row)">Enregistrer</BaseButton>
          </div>

          <div class="flex flex-col gap-2">
            <div
              v-for="attack in attacksFor(row.id)"
              :key="attack.id || attack.name"
              class="flex items-center gap-2 text-sm"
            >
              <input type="color" v-model="attack.color" class="w-8 h-8 rounded-card border border-border bg-surface-2" />
              <input
                v-model="attack.name"
                type="text"
                placeholder="Nom de l'attaque"
                class="flex-1 border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              />
              <input
                v-model.number="attack.forceBonus"
                type="number"
                title="Bonus de Force"
                class="w-20 border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              />
              <span class="text-text-muted text-xs">FOR</span>
              <input
                v-model.number="attack.manaCost"
                type="number"
                min="1"
                title="Coût en mana"
                class="w-20 border border-border rounded-card px-2 py-1 bg-surface-2 text-text-primary"
              />
              <span class="text-text-muted text-xs">mana</span>
              <BaseButton variant="secondary" @click="saveAttack(row, attack)">Enregistrer</BaseButton>
              <BaseButton variant="ghost" @click="removeAttack(row, attack)">Supprimer</BaseButton>
            </div>
            <div v-if="!attacksFor(row.id).length" class="text-sm text-text-muted">Aucune attaque pour l'instant.</div>
          </div>

          <div class="flex gap-2">
            <BaseButton
              variant="secondary"
              :disabled="attacksFor(row.id).length >= row.maxAttacks"
              @click="addAttack(row)"
            >
              + Ajouter une attaque
            </BaseButton>
            <BaseButton variant="ghost" @click="resetAttacks(row)">Réinitialiser aux valeurs par défaut</BaseButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
