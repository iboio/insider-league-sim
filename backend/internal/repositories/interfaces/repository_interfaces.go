package interfaces

import (
	"league-sim/internal/models"
)

type TeamsRepository interface {
	GetTeamByName(leagueId, teamName string) (models.Team, error)
	GetTeams(leagueId string) ([]models.Team, error)
	SetTeams(leagueId string, teams []models.Team) error
	UpdateTeam(team models.Team) error
}

type StandingRepository interface {
	GetStandings(leagueId string) ([]models.Standings, error)
	GetStandingsByTeam(leagueId, teamName string) (models.Standings, error)
	SetStandings(standings []models.Standings) error
	UpdateStanding(standing models.Standings) error
}

type LeagueRepository interface {
	GetLeagueById(id string) (models.CreateLeagueRequest, error)
	SetLeague(data models.CreateLeagueRequest) error
	GetLeagues() ([]models.GetLeaguesIdsWithNameResponse, error)
	DeleteLeague(id string) error
}

type MatchesRepository interface {
	SetMatches(matches []models.Matches) error
	GetMatches(leagueId string) ([]models.Matches, error)
	EditMatch(data models.Matches) error
	GetPlayedMatches(leagueId string) ([]models.Matches, error)
	GetFixtures(leagueId string) (models.GetFixturesResponse, error)
	GetMatchByTeams(data models.EditMatchResult) (models.Matches, error)
}
