export function budgetFromFavourites(fav) {
  const base = Math.log10(Math.max(fav, 1)) // ~0 à ~5
  return Math.round(20 + (base / 5) * 80)   // 20 à 100
}

export function ratioFromId(id) {
  return 0.3 + (id % 41) / 100 // stable, entre 0.30 et 0.70
}

export function statsFromCharacter(character, ratio) {
  const budget = budgetFromFavourites(character.favourites)
  const force = Math.round(budget * ratio)
  const pv = Math.round(budget * (1 - ratio) * 3)
  return { force, pv }
}

// Miroir de internal/store/stats.go — un personnage à haute valeur
// débloque jusqu'à 5 emplacements d'attaque, un personnage faible n'en a qu'un.
export function maxAttacksFromBudget(budget) {
  const n = Math.round(1 + ((budget - 20) / 80) * 4)
  return Math.min(5, Math.max(1, n))
}

// Le mana max réutilise directement le budget : un personnage fort a aussi
// plus de mana.
export function maxManaFromBudget(budget) {
  return budget
}
