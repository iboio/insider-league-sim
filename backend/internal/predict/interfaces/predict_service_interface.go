package interfaces

import "league-sim/internal/models"

type PredictServiceInterface interface {
	PredictChampionShipSession(id string) ([]models.PredictedStanding, error)
}
