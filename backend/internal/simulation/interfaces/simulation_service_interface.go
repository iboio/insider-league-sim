package interfaces

import "league-sim/internal/models"

type SimulationServiceInterface interface {
	Simulation(leagueId string, playAllFixture bool) (models.SimulationResponse, error)
	EditMatch(data models.EditMatchResult) error
}
