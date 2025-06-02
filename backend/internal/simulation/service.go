package simulation

import (
	"fmt"
	adaptInterface "league-sim/internal/layers/adapt/interfaces"
	"league-sim/internal/models"
	repoInterfaces "league-sim/internal/repositories/interfaces"
	"league-sim/internal/simulation/interfaces"
	"sort"
)

type SimulationService struct {
	matchesRepo   repoInterfaces.MatchesRepository
	standingsRepo repoInterfaces.StandingRepository
	teamsRepo     repoInterfaces.TeamsRepository
}

func NewSimulationService(adapt adaptInterface.AdaptInterface) interfaces.SimulationServiceInterface {
	return &SimulationService{
		matchesRepo:   adapt.MatchesRepository(),
		standingsRepo: adapt.StandingsRepository(),
		teamsRepo:     adapt.TeamRepository(),
	}
}

func (ss *SimulationService) Simulation(leagueId string, playAllFixture bool) (models.SimulationResponse, error) {

	teams, err := ss.teamsRepo.GetTeams(leagueId)
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("teams fetch error: %w", err)
	}

	standings, err := ss.standingsRepo.GetStandings(leagueId)
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("standings fetch error: %w", err)
	}

	fixtures, err := ss.matchesRepo.GetFixtures(leagueId)
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("fixtures fetch error: %w", err)
	}

	if len(fixtures.UpcomingFixtures) == 0 {
		return models.SimulationResponse{}, nil
	}

	teamMap := make(map[string]*models.Team)
	for i := range teams {
		teamMap[teams[i].TeamName] = &teams[i]
	}

	standingsMap := make(map[string]*models.Standings)
	for i := range standings {
		standingsMap[standings[i].TeamName] = &standings[i]
	}

	playingWeekCount := 1
	if playAllFixture {
		playingWeekCount = len(fixtures.UpcomingFixtures)
	}

	weeklyMatches := groupMatchesByWeek(fixtures.UpcomingFixtures)

	var playedMatches []models.Matches
	var updatedTeams []models.Team
	var updatedStandings []models.Standings

	var weekNumbers []int
	for weekNum := range weeklyMatches {
		weekNumbers = append(weekNumbers, weekNum)
	}
	sort.Ints(weekNumbers)

	for i := 0; i < playingWeekCount && i < len(weekNumbers); i++ {
		currentWeekNumber := weekNumbers[i]
		currentWeekMatches := weeklyMatches[currentWeekNumber]

		for _, match := range currentWeekMatches {

			matchResult := GenerateMatchResult(*teamMap[match.Home], *teamMap[match.Away])

			ss.updateTeamStats(standingsMap, teamMap, matchResult, match)

			matchData := ss.createMatchData(match, matchResult, currentWeekNumber, leagueId)
			playedMatches = append(playedMatches, matchData)
		}
	}

	for _, team := range teamMap {
		updatedTeams = append(updatedTeams, *team)
	}

	for _, standing := range standingsMap {

		standing.Diff = standing.Goals - standing.Against
		updatedStandings = append(updatedStandings, *standing)
	}

	err = ss.updateDatabase(leagueId, playedMatches, updatedTeams, updatedStandings)
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("database update error: %w", err)
	}

	fixtures, err = ss.matchesRepo.GetFixtures(leagueId)
	if err != nil {
		return models.SimulationResponse{}, fmt.Errorf("fixtures fetch error: %w", err)
	}

	return models.SimulationResponse{
		Matches:          playedMatches,
		UpcomingFixtures: fixtures.UpcomingFixtures,
		PlayedFixtures:   fixtures.PlayedFixtures,
	}, nil
}

func (ss *SimulationService) updateTeamStats(standingsMap map[string]*models.Standings,
	teamMap map[string]*models.Team,
	matchResult models.MatchOutcome,
	match models.Match) {
	if matchResult.IsDraw {

		DrawTeamAttributeChanging(standingsMap[match.Home], teamMap[match.Home], matchResult)
		DrawTeamAttributeChanging(standingsMap[match.Away], teamMap[match.Away], matchResult)
	} else {
		WinnerTeamAttributeChanging(
			standingsMap[matchResult.Winner.TeamName],
			teamMap[matchResult.Winner.TeamName],
			matchResult)
		LoserTeamAttributeChanging(
			standingsMap[matchResult.Loser.TeamName],
			teamMap[matchResult.Loser.TeamName],
			matchResult)
	}
}

func (ss *SimulationService) createMatchData(match models.Match,
	matchResult models.MatchOutcome,
	weekNumber int,
	leagueId string) models.Matches {
	var homeScore, awayScore int
	var winner string

	if matchResult.IsDraw {

		homeScore = matchResult.WinnerGoals
		awayScore = matchResult.LoserGoals
		winner = "draw"
	} else {

		if matchResult.Winner.TeamName == match.Home {
			homeScore = matchResult.WinnerGoals
			awayScore = matchResult.LoserGoals
			winner = match.Home
		} else {
			homeScore = matchResult.LoserGoals
			awayScore = matchResult.WinnerGoals
			winner = match.Away
		}
	}

	return models.Matches{
		LeagueId:  leagueId,
		MatchWeek: weekNumber,
		Home:      match.Home,
		HomeScore: homeScore,
		Away:      match.Away,
		AwayScore: awayScore,
		Winner:    winner,
		IsPlayed:  true,
	}
}

func (ss *SimulationService) updateDatabase(leagueId string,
	matches []models.Matches,
	teams []models.Team,
	standings []models.Standings) error {
	for _, match := range matches {
		err := ss.matchesRepo.EditMatch(match)
		if err != nil {
			return fmt.Errorf("match update error: %w", err)
		}
	}

	for _, team := range teams {
		err := ss.teamsRepo.UpdateTeam(team)
		if err != nil {
			return fmt.Errorf("team update error: %w", err)
		}
	}

	for _, standing := range standings {
		standing.LeagueId = leagueId
		err := ss.standingsRepo.UpdateStanding(standing)
		if err != nil {
			return fmt.Errorf("standing update error: %w", err)
		}
	}

	return nil
}

func (ss *SimulationService) EditMatch(data models.EditMatchResult) error {

	matchData, err := ss.matchesRepo.GetMatchByTeams(data)
	if err != nil {
		return fmt.Errorf("error fetching match data: %w", err)
	}

	standings, err := ss.standingsRepo.GetStandings(data.LeagueId)
	if err != nil {
		return fmt.Errorf("error fetching standings: %w", err)
	}

	teams, err := ss.teamsRepo.GetTeams(data.LeagueId)
	if err != nil {
		return fmt.Errorf("error fetching teams: %w", err)
	}

	standingsMap := make(map[string]*models.Standings)
	for i := range standings {
		standingsMap[standings[i].TeamName] = &standings[i]
	}

	teamMap := make(map[string]*models.Team)
	for i := range teams {
		teamMap[teams[i].TeamName] = &teams[i]
	}

	err = ss.revertMatchResult(matchData, standingsMap, teamMap)
	if err != nil {
		return fmt.Errorf("error reverting old match result: %w", err)
	}

	var newMatchOutcome models.MatchOutcome

	if data.HomeScore > data.AwayScore {
		newMatchOutcome = models.MatchOutcome{
			Winner:      *teamMap[data.Home],
			Loser:       *teamMap[data.Away],
			IsDraw:      false,
			WinnerGoals: data.HomeScore,
			LoserGoals:  data.AwayScore,
		}
		data.Winner = data.Home
	} else if data.HomeScore < data.AwayScore {
		newMatchOutcome = models.MatchOutcome{
			Winner:      *teamMap[data.Away],
			Loser:       *teamMap[data.Home],
			IsDraw:      false,
			WinnerGoals: data.AwayScore,
			LoserGoals:  data.HomeScore,
		}
		data.Winner = data.Away
	} else {
		newMatchOutcome = models.MatchOutcome{
			Winner:      *teamMap[data.Home],
			Loser:       *teamMap[data.Away],
			IsDraw:      true,
			WinnerGoals: data.HomeScore,
			LoserGoals:  data.AwayScore,
		}
		data.Winner = "draw"
	}

	err = ss.applyNewMatchResult(newMatchOutcome, standingsMap, teamMap)
	if err != nil {
		return fmt.Errorf("error applying new match result: %w", err)
	}

	for _, standing := range standingsMap {
		standing.Diff = standing.Goals - standing.Against
	}

	err = ss.updateDatabaseAfterEdit(data, standingsMap)
	if err != nil {
		return fmt.Errorf("error updating database: %w", err)
	}

	return nil
}

func (ss *SimulationService) revertMatchResult(matchData models.Matches,
	standingsMap map[string]*models.Standings,
	teamMap map[string]*models.Team) error {
	homeStanding := standingsMap[matchData.Home]
	awayStanding := standingsMap[matchData.Away]

	homeStanding.Played -= 1
	awayStanding.Played -= 1

	homeStanding.Goals -= matchData.HomeScore
	awayStanding.Goals -= matchData.AwayScore
	homeStanding.Against -= matchData.AwayScore
	awayStanding.Against -= matchData.HomeScore

	if matchData.Winner == "draw" {
		homeStanding.Points -= 1
		awayStanding.Points -= 1
		homeStanding.Draws -= 1
		awayStanding.Draws -= 1
	} else if matchData.Winner == matchData.Home {
		homeStanding.Points -= 3
		homeStanding.Wins -= 1
		awayStanding.Losses -= 1
	} else {
		awayStanding.Points -= 3
		awayStanding.Wins -= 1
		homeStanding.Losses -= 1
	}

	return nil
}

func (ss *SimulationService) applyNewMatchResult(matchOutcome models.MatchOutcome,
	standingsMap map[string]*models.Standings,
	teamMap map[string]*models.Team) error {
	if matchOutcome.IsDraw {
		DrawTeamAttributeChanging(
			standingsMap[matchOutcome.Winner.TeamName],
			teamMap[matchOutcome.Winner.TeamName],
			matchOutcome)
		DrawTeamAttributeChanging(
			standingsMap[matchOutcome.Loser.TeamName],
			teamMap[matchOutcome.Loser.TeamName],
			matchOutcome)
	} else {

		WinnerTeamAttributeChanging(
			standingsMap[matchOutcome.Winner.TeamName],
			teamMap[matchOutcome.Winner.TeamName],
			matchOutcome)
		LoserTeamAttributeChanging(
			standingsMap[matchOutcome.Loser.TeamName],
			teamMap[matchOutcome.Loser.TeamName],
			matchOutcome)
	}

	return nil
}

func (ss *SimulationService) updateDatabaseAfterEdit(data models.EditMatchResult,
	standingsMap map[string]*models.Standings) error {

	matchToUpdate := models.Matches{
		LeagueId:  data.LeagueId,
		MatchWeek: data.MatchWeek,
		Home:      data.Home,
		HomeScore: data.HomeScore,
		Away:      data.Away,
		AwayScore: data.AwayScore,
		Winner:    data.Winner,
		IsPlayed:  true,
	}

	err := ss.matchesRepo.EditMatch(matchToUpdate)
	if err != nil {
		return fmt.Errorf("match update error: %w", err)
	}

	for _, standing := range standingsMap {
		err := ss.standingsRepo.UpdateStanding(*standing)
		if err != nil {
			return fmt.Errorf("standing update error for %s: %w", standing.TeamName, err)
		}
	}

	return nil
}
