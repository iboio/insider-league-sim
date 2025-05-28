package interfaces

import (
	"league-sim/internal/models"
)

type LeagueRepository interface {
	SetLeague(id string, data models.CreateLeagueRequest) error
	GetLeague() ([]models.GetLeaguesIdsWithNameResponse, error)
	DeleteLeague(id string) error
}

type ActiveLeagueRepository interface {
	GetActiveLeague(id string) (models.League, error)
	SetActiveLeague(data models.League) error
	GetActiveLeagueTeams(id string) ([]models.Team, error)
	GetActiveLeaguesFixtures(id string) (models.GetActiveLeagueFixturesResponse, error)
	GetActiveLeaguesStandings(id string) ([]models.Standings, error)
}

type MatchResultRepository interface {
	EditMatchScore(data models.EditMatchResult) error
	SetMatchResults(leagueId string, matchResults []models.MatchResult) error
	GetMatchResults(leagueId string) ([]models.MatchResult, error)
	DeleteMatchResults(leagueId string) error
	GetMatchResultByWeekAndTeam(data models.EditMatchResult) (models.MatchResult, error)
}
