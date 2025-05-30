package league

import (
	"league-sim/internal/league/interfaces"
	repoInterfaces "league-sim/internal/repositories/interfaces"
	"strconv"

	"league-sim/internal/models"

	"github.com/google/uuid"
	adaptInterface "league-sim/internal/layers/adapt/interfaces"
)

type LeagueService struct {
	leagueRepo       repoInterfaces.LeagueRepository
	activeLeagueRepo repoInterfaces.ActiveLeagueRepository
	matchResultRepo  repoInterfaces.MatchResultRepository
}

func NewLeagueService(adapt adaptInterface.AdaptInterface) interfaces.LeagueServiceInterface {
	return &LeagueService{
		leagueRepo:       adapt.LeagueRepository(),
		activeLeagueRepo: adapt.ActiveLeagueRepository(),
		matchResultRepo:  adapt.MatchResultRepository(),
	}
}

func (ls *LeagueService) CreateLeague(n string, leagueName string) (models.GetLeaguesIdsWithNameResponse, error) {
	i, err := strconv.Atoi(n)

	if err != nil {
		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	leagueId := uuid.New()
	teams := TeamGenerate(i)
	fixtures := GenerateFixtures(teams)
	standings := CreateStandingsTable(teams)
	league := models.League{
		LeagueID:         leagueId.String(),
		LeagueName:       leagueName,
		Teams:            teams,
		Standings:        standings,
		TotalWeeks:       len(fixtures),
		CurrentWeek:      0,
		UpcomingFixtures: fixtures,
		PlayedFixtures:   []models.Week{},
	}

	err = ls.leagueRepo.SetLeague(leagueId.String(), models.CreateLeagueRequest{LeagueName: leagueName})

	if err != nil {

		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	err = ls.activeLeagueRepo.SetActiveLeague(league)

	if err != nil {

		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	return models.GetLeaguesIdsWithNameResponse{
		LeagueName: leagueName,
		LeagueId:   leagueId.String(),
	}, nil
}

func (ls *LeagueService) ResetLeague(leagueId string) error {
	league, err := ls.activeLeagueRepo.GetActiveLeague(leagueId)
	if err != nil {

		return err
	}

	teams := TeamGenerate(len(league.Teams))
	fixtures := GenerateFixtures(teams)
	standings := CreateStandingsTable(teams)
	league.LeagueID = leagueId
	league.Teams = teams
	league.Standings = standings
	league.TotalWeeks = len(fixtures)
	league.CurrentWeek = 0
	league.UpcomingFixtures = fixtures
	league.PlayedFixtures = []models.Week{}

	err = ls.activeLeagueRepo.SetActiveLeague(league)
	if err != nil {

		return err
	}
	err = ls.matchResultRepo.DeleteMatchResults(leagueId)
	if err != nil {

		return err
	}

	return nil
}
