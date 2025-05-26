package league

import (
	"fmt"

	"league-sim/internal/models"
)

func TeamGenerate(n int) []models.Team {
	teams := make([]models.Team, n)
	for i := 0; i < n; i++ {
		teams[i] = models.Team{
			Name:         fmt.Sprintf("Team %c", 'A'+i),
			AttackPower:  RandomNumberGenerator(70, 100),
			DefensePower: RandomNumberGenerator(70, 100),
			Stamina:      RandomNumberGenerator(70, 100),
			Morale:       RandomNumberGenerator(70, 100),
		}
	}

	return teams
}

func CreateStandingsTable(teams []models.Team) []models.Standings {
	var standings []models.Standings

	for _, teamData := range teams {
		standings = append(
			standings, models.Standings{
				Team:    teamData,
				Goals:   0,
				Against: 0,
				Played:  0,
				Wins:    0,
				Losses:  0,
				Points:  0,
			})
	}

	return standings
}

func GenerateFixtures(teams []models.Team) []models.Week {
	if len(teams)%2 != 0 {
		teams = append(teams, models.Team{Name: "BYE"})
	}

	n := len(teams)
	half := n / 2
	teamIndexes := make([]int, n)

	for i := 0; i < n; i++ {
		teamIndexes[i] = i
	}

	var weeks []models.Week

	for round := 0; round < n-1; round++ {
		var matches []models.Match

		for i := 0; i < half; i++ {
			home := &teams[teamIndexes[i]]
			away := &teams[teamIndexes[n-1-i]]
			if home.Name != "BYE" && away.Name != "BYE" {
				matches = append(matches, models.Match{Home: home, Away: away})
			}
		}

		weeks = append(
			weeks, models.Week{
				Number:  round + 1,
				Matches: matches,
			})

		last := teamIndexes[n-1]
		copy(teamIndexes[2:], teamIndexes[1:n-1])
		teamIndexes[1] = last
	}

	return weeks
}
