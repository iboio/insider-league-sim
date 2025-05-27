package simulation

import (
	"league-sim/internal/models"
)

func LoserTeamAttributeChanging(standings *models.Standings, team *models.Team, matchOutCome models.MatchOutcome) {
	standings.Goals += matchOutCome.LoserGoals
	standings.Against += matchOutCome.WinnerGoals
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

func WinnerTeamAttributeChanging(standings *models.Standings, team *models.Team, matchOutCome models.MatchOutcome) {
	standings.Goals += matchOutCome.WinnerGoals
	standings.Against += matchOutCome.LoserGoals
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

func DrawTeamAttributeChanging(standings *models.Standings, team *models.Team, matchOutCome models.MatchOutcome) {
	standings.Goals += matchOutCome.LoserGoals
	standings.Against += matchOutCome.WinnerGoals
	standings.Played += 1
	standings.Points += 1

	team.Stamina -= 5
	if team.Stamina < 0 {
		team.Stamina = 0
	}
}
