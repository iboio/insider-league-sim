package interfaces

import repoInterfaces "league-sim/internal/repositories/interfaces"

type AdaptInterface interface {
	LeagueRepository() repoInterfaces.LeagueRepository
	ActiveLeagueRepository() repoInterfaces.ActiveLeagueRepository
	MatchResultRepository() repoInterfaces.MatchResultRepository
}
