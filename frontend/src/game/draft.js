import { budgetFromFavourites } from './stats.js'

export function suggestedPrice(character) {
  const budget = budgetFromFavourites(character.favourites)
  return budget * 3000
}

export const rewardTiers = {
  faible: 10000,
  moyen: 25000,
  fort: 50000,
}

export function rewardForVictory(tier, customAmount) {
  if (tier === 'custom') return customAmount ?? 0
  return rewardTiers[tier] ?? 0
}
