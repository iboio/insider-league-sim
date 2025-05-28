package simulation

import (
	appContext "league-sim/internal/contexts/appContexts"
	"league-sim/internal/league"
	"league-sim/internal/models"
	"math/rand"
)

type SimulationService struct {
	appCtx appContext.AppContext
}

func NewSimulationService(ctx appContext.AppContext) *SimulationService {
	return &SimulationService{
		appCtx: ctx,
	}
}

func (ss *SimulationService) Simulation(leagueId string, playAllFixture bool) (models.SimulationResponse, error) {
	var matches []models.MatchResult
	activeLeague, err := ss.appCtx.ActiveLeagueRepository().GetActiveLeague(leagueId)
	if err != nil {
		return models.SimulationResponse{}, err
	}

	if len(activeLeague.UpcomingFixtures) == 0 {
		return models.SimulationResponse{}, nil
	}

	activeLeague.TotalWeeks = len(activeLeague.UpcomingFixtures) + len(activeLeague.PlayedFixtures)
	teamMap := make(map[string]*models.Team)
	for _, t := range activeLeague.Teams {
		teamMap[t.Name] = &t
	}

	standingsMap := make(map[string]*models.Standings)
	for i := range activeLeague.Standings {
		t := activeLeague.Standings[i].Team.Name
		standingsMap[t] = &activeLeague.Standings[i]
	}

	playingWeekCount := 1

	if playAllFixture {
		playingWeekCount = len(activeLeague.UpcomingFixtures)
	}

	for i := 0; i < playingWeekCount; i++ {

		var currentFixtureWeek models.Week

		if playAllFixture {
			currentFixtureWeek = activeLeague.UpcomingFixtures[0]
		} else {
			currentFixtureWeek = activeLeague.UpcomingFixtures[i]
		}

		for _, match := range currentFixtureWeek.Matches {
			matchResult := GenerateMatchResult(*match.Home, *match.Away)
			homeStanding := standingsMap[match.Home.Name]
			awayStanding := standingsMap[match.Away.Name]

			winnerTeam := teamMap[matchResult.Winner.Name]
			loserTeam := teamMap[matchResult.Loser.Name]

			if matchResult.IsDraw {
				DrawTeamAttributeChanging(
					standingsMap[matchResult.Winner.Name],
					winnerTeam,
					matchResult)
				DrawTeamAttributeChanging(
					standingsMap[matchResult.Loser.Name],
					loserTeam,
					matchResult)
			} else {
				WinnerTeamAttributeChanging(
					standingsMap[matchResult.Winner.Name],
					winnerTeam,
					matchResult)
				LoserTeamAttributeChanging(
					standingsMap[matchResult.Loser.Name],
					loserTeam,
					matchResult)
			}
			homeStanding.Team = *teamMap[match.Home.Name]
			awayStanding.Team = *teamMap[match.Away.Name]

			var homeScore, awayScore int
			matchWinner := matchResult.Winner.Name
			if matchResult.IsDraw {
				homeScore = matchResult.WinnerGoals
				awayScore = matchResult.LoserGoals
			} else if match.Home.Name == matchResult.Winner.Name {
				homeScore = matchResult.WinnerGoals
				awayScore = matchResult.LoserGoals

			} else if matchResult.WinnerGoals == matchResult.LoserGoals {
				matchWinner = "draw"
			} else {
				homeScore = matchResult.LoserGoals
				awayScore = matchResult.WinnerGoals
			}

			matches = append(
				matches, models.MatchResult{
					MatchWeek: currentFixtureWeek.Number,
					Home:      match.Home.Name,
					HomeScore: homeScore,
					Away:      match.Away.Name,
					AwayScore: awayScore,
					Winner:    matchWinner,
				})
		}

		activeLeague.PlayedFixtures = append(activeLeague.PlayedFixtures, currentFixtureWeek)
		activeLeague.UpcomingFixtures = activeLeague.UpcomingFixtures[1:]
		activeLeague.CurrentWeek = currentFixtureWeek.Number
	}

	err = ss.appCtx.ActiveLeagueRepository().SetActiveLeague(activeLeague)

	if err != nil {
		panic(err)
		return models.SimulationResponse{}, err
	}

	err = ss.appCtx.MatchResultRepository().SetMatchResults(activeLeague.LeagueID, matches)

	if err != nil {
		panic(err)
		return models.SimulationResponse{}, err
	}

	return models.SimulationResponse{
		Matches:          matches,
		UpcomingFixtures: activeLeague.UpcomingFixtures,
		PlayedFixtures:   activeLeague.PlayedFixtures,
	}, nil
}

func GenerateMatchResult(home models.Team, away models.Team) models.MatchOutcome {
	homeScore := league.CalculateStrength(home) * 1.05
	awayScore := league.CalculateStrength(away)

	total := homeScore + awayScore

	drawChance := 0.2
	if rand.Float64() < drawChance {
		goals := rand.Intn(3)
		return models.MatchOutcome{
			Winner:      home,
			Loser:       away,
			IsDraw:      true,
			LoserGoals:  goals,
			WinnerGoals: goals,
		}
	}

	homeChance := homeScore / total
	homeWins := rand.Float64() < homeChance

	winnerGoals := rand.Intn(5) + 1
	loserGoals := rand.Intn(winnerGoals)

	if homeWins {
		return models.MatchOutcome{
			Winner:      home,
			Loser:       away,
			IsDraw:      false,
			WinnerGoals: winnerGoals,
			LoserGoals:  loserGoals,
		}
	} else {
		return models.MatchOutcome{
			Winner:      away,
			Loser:       home,
			IsDraw:      false,
			WinnerGoals: winnerGoals,
			LoserGoals:  loserGoals,
		}
	}
}

func (ss *SimulationService) EditMatch(data models.EditMatchResult) error {
	matching, err := ss.appCtx.MatchResultRepository().GetMatchResultByWeekAndTeam(data)
	if err != nil {
		return err
	}
	activeLeague, err := ss.appCtx.ActiveLeagueRepository().GetActiveLeague(data.LeagueId)
	if err != nil {
		return err
	}

	standingsMap := make(map[string]*models.Standings)
	for i := range activeLeague.Standings {
		t := activeLeague.Standings[i].Team.Name
		standingsMap[t] = &activeLeague.Standings[i]
	}

	teamMap := make(map[string]*models.Team)
	for _, t := range activeLeague.Teams {
		teamMap[t.Name] = &t
	}

	awayStanding := standingsMap[matching.Away]
	homeStanding := standingsMap[matching.Home]

	homeTeam := teamMap[matching.Home]
	awayTeam := teamMap[matching.Away]

	if matching.HomeScore == matching.AwayScore {
		awayStanding.Points -= 1
		homeStanding.Points -= 1

	}
	if matching.HomeScore == matching.AwayScore {
		awayStanding.Wins -= 1
		homeStanding.Wins -= 1
	} else if matching.HomeScore > matching.AwayScore {
		homeStanding.Points -= 3
		homeStanding.Wins -= 1

	} else if matching.HomeScore < matching.AwayScore {
		awayStanding.Points -= 3
		awayStanding.Wins -= 1
	}
	awayStanding.Goals -= matching.AwayScore
	homeStanding.Goals -= matching.HomeScore

	awayStanding.Played -= 1
	homeStanding.Played -= 1

	awayStanding.Against -= matching.HomeScore
	homeStanding.Against -= matching.AwayScore

	if data.HomeScore == data.AwayScore {
		DrawTeamAttributeChanging(
			awayStanding, awayTeam, models.MatchOutcome{
				Winner:      *awayTeam,
				Loser:       *homeTeam,
				IsDraw:      true,
				WinnerGoals: data.HomeScore,
				LoserGoals:  data.AwayScore,
			})
		DrawTeamAttributeChanging(
			awayStanding, awayTeam, models.MatchOutcome{
				Winner:      *awayTeam,
				Loser:       *homeTeam,
				IsDraw:      true,
				WinnerGoals: data.HomeScore,
				LoserGoals:  data.AwayScore,
			})
		data.Winner = "draw"
	}
	if data.HomeScore > data.AwayScore {
		WinnerTeamAttributeChanging(
			homeStanding, homeTeam, models.MatchOutcome{
				Winner:      *homeTeam,
				Loser:       *awayTeam,
				IsDraw:      false,
				WinnerGoals: data.HomeScore,
				LoserGoals:  data.AwayScore,
			})
		LoserTeamAttributeChanging(
			awayStanding, awayTeam, models.MatchOutcome{
				Winner:      *homeTeam,
				Loser:       *awayTeam,
				IsDraw:      false,
				WinnerGoals: data.HomeScore,
				LoserGoals:  data.AwayScore,
			})
		data.Winner = data.Home
	}
	if data.HomeScore < data.AwayScore {
		WinnerTeamAttributeChanging(
			awayStanding, awayTeam, models.MatchOutcome{
				Winner:      *awayTeam,
				Loser:       *homeTeam,
				IsDraw:      false,
				WinnerGoals: data.AwayScore,
				LoserGoals:  data.HomeScore,
			})
		LoserTeamAttributeChanging(
			homeStanding, homeTeam, models.MatchOutcome{
				Winner:      *awayTeam,
				Loser:       *homeTeam,
				IsDraw:      false,
				WinnerGoals: data.AwayScore,
				LoserGoals:  data.HomeScore,
			})
		data.Winner = data.Away
	}

	err = ss.appCtx.ActiveLeagueRepository().SetActiveLeague(activeLeague)
	if err != nil {
		return err
	}
	err = ss.appCtx.MatchResultRepository().EditMatchScore(data)
	if err != nil {
		return err
	}
	return nil
}
