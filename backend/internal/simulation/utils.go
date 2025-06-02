package simulation

import (
	"league-sim/internal/league"
	"league-sim/internal/models"
	"math/rand"
)

func LoserTeamAttributeChanging(standings *models.Standings, team *models.Team, matchOutcome models.MatchOutcome) {
	standings.Goals += matchOutcome.LoserGoals
	standings.Against += matchOutcome.WinnerGoals
	standings.Played += 1
	standings.Losses += 1

	team.Stamina -= 5
	if team.Stamina < 0 {
		team.Stamina = 0
	}
	team.Morale -= 5
	if team.Morale < 0 {
		team.Morale = 0
	}
}

func WinnerTeamAttributeChanging(standings *models.Standings, team *models.Team, matchOutcome models.MatchOutcome) {
	standings.Goals += matchOutcome.WinnerGoals
	standings.Against += matchOutcome.LoserGoals
	standings.Played += 1
	standings.Wins += 1
	standings.Points += 3

	team.Stamina -= 5
	if team.Stamina < 0 {
		team.Stamina = 0
	}
	team.Morale += 5
	if team.Morale > 100 {
		team.Morale = 100
	}
}

func DrawTeamAttributeChanging(standings *models.Standings, team *models.Team, matchOutcome models.MatchOutcome) {

	standings.Goals += matchOutcome.WinnerGoals
	standings.Against += matchOutcome.LoserGoals
	standings.Played += 1
	standings.Points += 1
	standings.Draws += 1

	team.Stamina -= 5
	if team.Stamina < 0 {
		team.Stamina = 0
	}
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

func groupMatchesByWeek(upcomingFixtures []models.Week) map[int][]models.Match {
	weeklyMatches := make(map[int][]models.Match)

	for _, fixture := range upcomingFixtures {
		weekNumber := fixture.Number

		if _, exists := weeklyMatches[weekNumber]; !exists {
			weeklyMatches[weekNumber] = []models.Match{}
		}
		
		for _, match := range fixture.Matches {
			weeklyMatches[weekNumber] = append(weeklyMatches[weekNumber], match)
		}
	}

	return weeklyMatches
}
