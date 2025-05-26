package interfaces

import "league-sim/internal/models"

type LeagueServiceInterface interface {
	CreateLeague(n string, leagueName string) (models.GetLeaguesIdsWithNameResponse, error)
	ResetLeague(leagueId string) error
}
