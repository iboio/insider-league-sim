package league

import (
	"strconv"

	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/models"

	"github.com/google/uuid"
)

type LeagueService struct {
	appCtx appContext.AppContext
}

func NewLeagueService(ctx appContext.AppContext) *LeagueService {
	return &LeagueService{
		appCtx: ctx,
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

	err = ls.appCtx.LeagueRepository().SetLeague(leagueId.String(), models.CreateLeagueRequest{LeagueName: leagueName})

	if err != nil {

		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	err = ls.appCtx.ActiveLeagueRepository().SetActiveLeague(league)

	if err != nil {

		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	return models.GetLeaguesIdsWithNameResponse{
		LeagueName: leagueName,
		LeagueId:   leagueId.String(),
	}, nil
}

func (ls *LeagueService) ResetLeague(leagueId string) error {
	league, err := ls.appCtx.ActiveLeagueRepository().GetActiveLeague(leagueId)
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

	err = ls.appCtx.ActiveLeagueRepository().SetActiveLeague(league)
	if err != nil {

		return err
	}
	err = ls.appCtx.MatchResultRepository().DeleteMatchResults(leagueId)
	if err != nil {

		return err
	}

	return nil
}
