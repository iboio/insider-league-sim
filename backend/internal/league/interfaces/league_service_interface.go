package interfaces

import "league-sim/internal/models"

type LeagueServiceInterface interface {
	CreateLeague(request models.CreateLeagueRequest) (models.GetLeaguesIdsWithNameResponse, error)
	ResetLeague(leagueId string) error
}
