package league

import (
	"fmt"

	"league-sim/internal/models"
)

func TeamGenerate(n int) []models.Team {
	teams := make([]models.Team, n)
	for i := 0; i < n; i++ {
		teams[i] = models.Team{
			TeamName:     fmt.Sprintf("Team %c", 'A'+i),
			AttackPower:  RandomNumberGenerator(70, 100),
			DefensePower: RandomNumberGenerator(70, 100),
			Stamina:      RandomNumberGenerator(70, 100),
			Morale:       RandomNumberGenerator(70, 100),
		}
	}

	return teams
}

func GenerateFixtures(teams []models.Team) []models.Week {
	if len(teams)%2 != 0 {
		teams = append(teams, models.Team{TeamName: "BYE"})
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
			if home.TeamName != "BYE" && away.TeamName != "BYE" {
				matches = append(matches, models.Match{Home: home.TeamName, Away: away.TeamName})
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
