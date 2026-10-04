export function resolveRound(statsA, statsB) {
  if (statsA.force === statsB.force) {
    return { winner: 'draw' }
  }
  const winner = statsA.force > statsB.force ? 'a' : 'b'
  const winnerStats = winner === 'a' ? statsA : statsB
  const loserStats = winner === 'a' ? statsB : statsA
  return {
    winner,
    winnerStatsAfter: {
      force: winnerStats.force - loserStats.force,
      pv: winnerStats.pv - loserStats.pv,
    },
  }
}
