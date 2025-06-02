package interfaces

import repoInterfaces "league-sim/internal/repositories/interfaces"

type AdaptInterface interface {
	LeagueRepository() repoInterfaces.LeagueRepository
	StandingsRepository() repoInterfaces.StandingRepository
	MatchesRepository() repoInterfaces.MatchesRepository
	TeamRepository() repoInterfaces.TeamsRepository
}
