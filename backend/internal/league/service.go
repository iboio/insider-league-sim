package league

import (
	"league-sim/internal/league/interfaces"
	repoInterfaces "league-sim/internal/repositories/interfaces"
	"strconv"

	"league-sim/internal/models"

	adaptInterface "league-sim/internal/layers/adapt/interfaces"
)

type leagueService struct {
	leagueRepo    repoInterfaces.LeagueRepository
	matchesRepo   repoInterfaces.MatchesRepository
	standingsRepo repoInterfaces.StandingRepository
	teamsRepo     repoInterfaces.TeamsRepository
}

func NewLeagueService(adapt adaptInterface.AdaptInterface) interfaces.LeagueServiceInterface {
	return &leagueService{
		leagueRepo:    adapt.LeagueRepository(),
		matchesRepo:   adapt.MatchesRepository(),
		standingsRepo: adapt.StandingsRepository(),
		teamsRepo:     adapt.TeamRepository(),
	}
}

func (ls *leagueService) CreateLeague(data models.CreateLeagueRequest) (models.GetLeaguesIdsWithNameResponse, error) {
	i, err := strconv.Atoi(data.TeamCount)

	if err != nil {
		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	teams := TeamGenerate(i)
	fixtures := GenerateFixtures(teams)

	err = ls.leagueRepo.SetLeague(
		models.CreateLeagueRequest{
			LeagueId:   data.LeagueId,
			LeagueName: data.LeagueName,
			TeamCount:  data.TeamCount,
		})

	if err != nil {
		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	err = ls.SetMatches(data.LeagueId, &fixtures)
	if err != nil {
		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	err = ls.SetStandings(data.LeagueId, &teams)
	if err != nil {
		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	err = ls.teamsRepo.SetTeams(data.LeagueId, teams)
	if err != nil {
		return models.GetLeaguesIdsWithNameResponse{}, err
	}

	return models.GetLeaguesIdsWithNameResponse{
		LeagueName: data.LeagueName,
		LeagueId:   data.LeagueId,
	}, nil
}

func (ls *leagueService) SetMatches(leagueId string, fixtures *[]models.Week) error {
	var matches []models.Matches
	for _, week := range *fixtures {
		for _, match := range week.Matches {
			matches = append(
				matches, models.Matches{
					LeagueId:  leagueId,
					Home:      match.Home,
					HomeScore: 0,
					Away:      match.Away,
					AwayScore: 0,
					MatchWeek: week.Number,
				})
		}
	}

	return ls.matchesRepo.SetMatches(matches)
}

func (ls *leagueService) SetStandings(leagueId string, teams *[]models.Team) error {
	var standings []models.Standings
	for _, team := range *teams {
		standings = append(
			standings, models.Standings{
				LeagueId: leagueId,
				TeamName: team.TeamName,
				Goals:    0,
				Against:  0,
				Played:   0,
				Wins:     0,
				Losses:   0,
				Draws:    0,
				Points:   0,
			})
	}
	return ls.standingsRepo.SetStandings(standings)
}

func (ls *leagueService) ResetLeague(leagueId string) error {
	result, err := ls.leagueRepo.GetLeagueById(leagueId)
	if err != nil {
		return err
	}
	err = ls.leagueRepo.DeleteLeague(leagueId)
	if err != nil {
		return err
	}
	_, err = ls.CreateLeague(result)
	return nil
}
