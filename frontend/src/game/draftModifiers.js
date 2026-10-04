const intensityRates = {
  rare: 0.12,
  frequent: 0.3,
}

export function rollModifier(catalog, intensity) {
  const rate = intensityRates[intensity] ?? 0
  if (!catalog?.length || Math.random() > rate) return null
  const pool = catalog.filter((m) => !m.rare || Math.random() < 0.25)
  if (!pool.length) return null
  return pool[Math.floor(Math.random() * pool.length)]
}
