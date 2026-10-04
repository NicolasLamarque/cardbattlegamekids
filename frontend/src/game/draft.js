export const rewardTiers = {
  faible: 10000,
  moyen: 25000,
  fort: 50000,
}

export function rewardForVictory(tier, customAmount) {
  if (tier === 'custom') return customAmount ?? 0
  return rewardTiers[tier] ?? 0
}
