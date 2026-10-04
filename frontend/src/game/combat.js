// resolveRound compare deux cartes posées, chacune avec sa Force de base et
// (optionnel) un bonus d'attaque liée. Le gagnant infge à l'adversaire des
// dégâts égaux à sa Force effective — un personnage tombé à 0 PV ou moins
// est éliminé.
export function resolveRound(a, b) {
  const forceA = a.force + (a.attackBonus ?? 0)
  const forceB = b.force + (b.attackBonus ?? 0)

  if (forceA === forceB) {
    return { winner: 'draw', forceA, forceB, pvAAfter: a.pv, pvBAfter: b.pv, eliminatedA: false, eliminatedB: false }
  }

  const winner = forceA > forceB ? 'a' : 'b'
  const damage = winner === 'a' ? forceA : forceB
  const pvAAfter = winner === 'a' ? a.pv : a.pv - damage
  const pvBAfter = winner === 'b' ? b.pv : b.pv - damage

  return {
    winner,
    forceA,
    forceB,
    pvAAfter,
    pvBAfter,
    eliminatedA: pvAAfter <= 0,
    eliminatedB: pvBAfter <= 0,
  }
}

// Un personnage ne peut lier une attaque que s'il a assez de mana pour la payer.
export function canUseAttack(mana, attack) {
  return !attack || mana >= attack.manaCost
}

// Recharge de mana entre deux manches, plafonnée au max du personnage.
export function regenMana(current, regen, max) {
  return Math.min(max, current + regen)
}

// Une équipe a perdu quand plus aucun personnage ne lui reste : ni en main,
// ni en file d'attente, ni sur l'emplacement de combat.
export function isTeamDefeated(hand, queue, slot) {
  return hand.length === 0 && queue.length === 0 && !slot
}
