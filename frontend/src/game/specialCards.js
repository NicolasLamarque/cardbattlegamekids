export const specialCards = [
  { id: 'heal-10', name: 'Petit soin', type: 'heal', amount: 10 },
  { id: 'heal-20', name: 'Soin', type: 'heal', amount: 20 },
  { id: 'heal-30', name: 'Grand soin', type: 'heal', amount: 30 },
  { id: 'double-attack', name: 'Rage', type: 'buff', stat: 'force', multiplier: 2 },
  { id: 'mana-recharge', name: 'Recharge mana', type: 'manaRecharge', amount: 5 },
]

export function dealHands(catalog, playerIds, countPerPlayer) {
  const hands = {}
  for (const id of playerIds) {
    hands[id] = Array.from(
      { length: countPerPlayer },
      () => catalog[Math.floor(Math.random() * catalog.length)]
    )
  }
  return hands
}

export function applyEffect(card, targetStats) {
  if (card.type === 'heal') {
    return { ...targetStats, pv: targetStats.pv + card.amount }
  }
  if (card.type === 'buff' && card.stat === 'force') {
    return { ...targetStats, force: targetStats.force * card.multiplier }
  }
  if (card.type === 'manaRecharge') {
    return { ...targetStats, mana: targetStats.mana + card.amount }
  }
  return targetStats
}
