package adapt

import (
	"league-sim/internal/layers/infra"
	"league-sim/internal/repositories"
	repoInterfaces "league-sim/internal/repositories/interfaces"
)

type Adapt struct {
	leagueRepository       repoInterfaces.LeagueRepository
	activeLeagueRepository repoInterfaces.ActiveLeagueRepository
	matchResultRepository  repoInterfaces.MatchResultRepository
}

func BuildAdaptLayer(infra *infra.Infra) *Adapt {
	newActiveLeagueRepo := repositories.NewActiveLeagueRepository(infra)
	newLeagueRepo := repositories.NewLeagueRepository(infra)
	newMatchResultRepo := repositories.NewMatchResultRepository(infra)

	return &Adapt{
		leagueRepository:       newLeagueRepo,
		activeLeagueRepository: newActiveLeagueRepo,
		matchResultRepository:  newMatchResultRepo,
	}
}

func (a *Adapt) LeagueRepository() repoInterfaces.LeagueRepository {
	return a.leagueRepository
}

func (a *Adapt) ActiveLeagueRepository() repoInterfaces.ActiveLeagueRepository {
	return a.activeLeagueRepository
}

func (a *Adapt) MatchResultRepository() repoInterfaces.MatchResultRepository {
	return a.matchResultRepository
}
