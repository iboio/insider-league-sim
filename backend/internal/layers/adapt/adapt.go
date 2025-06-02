package adapt

import (
	"league-sim/internal/layers/infra"
	"league-sim/internal/repositories"
	repoInterfaces "league-sim/internal/repositories/interfaces"
)

type Adapt struct {
	leagueRepository   repoInterfaces.LeagueRepository
	matchesRepository  repoInterfaces.MatchesRepository
	standingRepository repoInterfaces.StandingRepository
	teamsRepository    repoInterfaces.TeamsRepository
}

func BuildAdaptLayer(infra *infra.Infra) *Adapt {
	newLeagueRepo := repositories.NewLeagueRepository(infra)
	newMatchesRepo := repositories.NewMatchResultRepository(infra)
	newStandingRepo := repositories.NewStandingRepository(infra)
	teamsRepo := repositories.NewTeamsRepository(infra)

	return &Adapt{
		leagueRepository:   newLeagueRepo,
		matchesRepository:  newMatchesRepo,
		standingRepository: newStandingRepo,
		teamsRepository:    teamsRepo,
	}
}

func (a *Adapt) LeagueRepository() repoInterfaces.LeagueRepository {
	return a.leagueRepository
}

func (a *Adapt) StandingsRepository() repoInterfaces.StandingRepository {
	return a.standingRepository
}

func (a *Adapt) MatchesRepository() repoInterfaces.MatchesRepository {
	return a.matchesRepository
}

func (a *Adapt) TeamRepository() repoInterfaces.TeamsRepository {
	return a.teamsRepository
}
