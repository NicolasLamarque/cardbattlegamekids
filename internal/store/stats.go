package store

import "math"

// Miroir exact de frontend/src/game/stats.js — les deux doivent rester synchronisés.

func budgetFromFavourites(fav int) int {
	base := math.Log10(math.Max(float64(fav), 1))
	return int(math.Round(20 + (base/5)*80))
}

func ratioFromID(id int) float64 {
	return 0.3 + float64(id%41)/100
}

func statsFromFavourites(fav int, ratio float64) (force int, pv int) {
	budget := budgetFromFavourites(fav)
	force = int(math.Round(float64(budget) * ratio))
	pv = int(math.Round(float64(budget) * (1 - ratio) * 3))
	return force, pv
}
