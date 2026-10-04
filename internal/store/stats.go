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

// maxAttacksFromBudget donne de 1 à 5 emplacements d'attaque selon la valeur
// du personnage — un budget élevé (ex: Gojo) débloque les 5, un budget
// faible n'en débloque qu'un seul. Miroir de game/stats.js.
func maxAttacksFromBudget(budget int) int {
	n := int(math.Round(1 + (float64(budget-20)/80)*4))
	if n < 1 {
		n = 1
	}
	if n > 5 {
		n = 5
	}
	return n
}

// maxManaFromBudget réutilise directement le budget comme jauge de mana
// max : un personnage fort a aussi plus de mana.
func maxManaFromBudget(budget int) int {
	return budget
}

// defaultManaRegen donne une recharge par manche raisonnable par défaut —
// réglable ensuite à la main dans l'admin.
func defaultManaRegen(maxMana int) int {
	r := int(math.Round(float64(maxMana) / 10))
	if r < 1 {
		r = 1
	}
	return r
}
