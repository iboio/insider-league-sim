package league

import (
	"math/rand"

	"league-sim/internal/models"
)

func RandomNumberGenerator(min float64, max float64) float64 {
	if min >= max {
		return min
	}
	return min + rand.Float64()*(max-min)
}

func CalculateStrength(team models.Team) float64 {
	return team.AttackPower*0.3 + team.DefensePower*0.3 + team.Morale*0.2 + team.Stamina*0.2
}
