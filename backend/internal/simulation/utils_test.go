package simulation

import (
	"testing"

	"league-sim/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestLoserTeamAttributeChanging(t *testing.T) {
	tests := []struct {
		name         string
		initialTeam  models.Team
		standings    models.Standings
		matchOutcome models.MatchOutcome
		expectedTeam models.Team
	}{
		{
			name: "Normal losing team",
			initialTeam: models.Team{
				Name:         "Loser Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      80,
				Morale:       80,
			},
			standings: models.Standings{
				Team:    models.Team{Name: "Loser Team"},
				Points:  0,
				Played:  0,
				Wins:    0,
				Losses:  0,
				Goals:   0,
				Against: 0,
			},
			matchOutcome: models.MatchOutcome{
				Winner:      models.Team{Name: "Winner Team"},
				Loser:       models.Team{Name: "Loser Team"},
				IsDraw:      false,
				WinnerGoals: 3,
				LoserGoals:  1,
			},
			expectedTeam: models.Team{
				Name:         "Loser Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      75, // 80 - 5
				Morale:       75, // 80 - 5
			},
		},
		{
			name: "Low stamina team losing",
			initialTeam: models.Team{
				Name:         "Low Stamina Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      3, // Will go below 0
				Morale:       80,
			},
			standings: models.Standings{
				Team:   models.Team{Name: "Low Stamina Team"},
				Points: 0,
			},
			matchOutcome: models.MatchOutcome{
				WinnerGoals: 2,
				LoserGoals:  0,
			},
			expectedTeam: models.Team{
				Name:         "Low Stamina Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      0, // Should be capped at 0
				Morale:       75,
			},
		},
		{
			name: "Low morale team losing",
			initialTeam: models.Team{
				Name:         "Low Morale Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      80,
				Morale:       2, // Will go below 0
			},
			standings: models.Standings{
				Team:   models.Team{Name: "Low Morale Team"},
				Points: 0,
			},
			matchOutcome: models.MatchOutcome{
				WinnerGoals: 1,
				LoserGoals:  0,
			},
			expectedTeam: models.Team{
				Name:         "Low Morale Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      75,
				Morale:       0, // Should be capped at 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create copies to avoid modifying original data
			team := tt.initialTeam
			standings := tt.standings

			// Execute
			LoserTeamAttributeChanging(&standings, &team, tt.matchOutcome)

			// Assert team attributes
			assert.Equal(t, tt.expectedTeam.Stamina, team.Stamina, "Stamina should be correctly updated")
			assert.Equal(t, tt.expectedTeam.Morale, team.Morale, "Morale should be correctly updated")
			assert.Equal(t, tt.expectedTeam.AttackPower, team.AttackPower, "Attack power should remain unchanged")
			assert.Equal(t, tt.expectedTeam.DefensePower, team.DefensePower, "Defense power should remain unchanged")

			// Assert standings updates
			assert.Equal(t, tt.matchOutcome.LoserGoals, standings.Goals, "Goals should be updated")
			assert.Equal(t, tt.matchOutcome.WinnerGoals, standings.Against, "Goals against should be updated")
			assert.Equal(t, 1, standings.Played, "Played should be incremented")
			assert.Equal(t, 1, standings.Losses, "Losses should be incremented")
			assert.Equal(t, 0, standings.Wins, "Wins should remain unchanged")
			assert.Equal(t, 0, standings.Points, "Points should remain unchanged for loser")
		})
	}
}

func TestWinnerTeamAttributeChanging(t *testing.T) {
	tests := []struct {
		name         string
		initialTeam  models.Team
		standings    models.Standings
		matchOutcome models.MatchOutcome
		expectedTeam models.Team
	}{
		{
			name: "Normal winning team",
			initialTeam: models.Team{
				Name:         "Winner Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      80,
				Morale:       80,
			},
			standings: models.Standings{
				Team:    models.Team{Name: "Winner Team"},
				Points:  0,
				Played:  0,
				Wins:    0,
				Losses:  0,
				Goals:   0,
				Against: 0,
			},
			matchOutcome: models.MatchOutcome{
				Winner:      models.Team{Name: "Winner Team"},
				Loser:       models.Team{Name: "Loser Team"},
				IsDraw:      false,
				WinnerGoals: 3,
				LoserGoals:  1,
			},
			expectedTeam: models.Team{
				Name:         "Winner Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      75, // 80 - 5
				Morale:       85, // 80 + 5
			},
		},
		{
			name: "High morale team winning",
			initialTeam: models.Team{
				Name:         "High Morale Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      80,
				Morale:       98, // Will go above 100
			},
			standings: models.Standings{
				Team:   models.Team{Name: "High Morale Team"},
				Points: 0,
			},
			matchOutcome: models.MatchOutcome{
				WinnerGoals: 2,
				LoserGoals:  0,
			},
			expectedTeam: models.Team{
				Name:         "High Morale Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      75,
				Morale:       100, // Should be capped at 100
			},
		},
		{
			name: "Low stamina team winning",
			initialTeam: models.Team{
				Name:         "Low Stamina Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      3, // Will go below 0
				Morale:       80,
			},
			standings: models.Standings{
				Team:   models.Team{Name: "Low Stamina Team"},
				Points: 0,
			},
			matchOutcome: models.MatchOutcome{
				WinnerGoals: 1,
				LoserGoals:  0,
			},
			expectedTeam: models.Team{
				Name:         "Low Stamina Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      0, // Should be capped at 0
				Morale:       85,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create copies to avoid modifying original data
			team := tt.initialTeam
			standings := tt.standings

			// Execute
			WinnerTeamAttributeChanging(&standings, &team, tt.matchOutcome)

			// Assert team attributes
			assert.Equal(t, tt.expectedTeam.Stamina, team.Stamina, "Stamina should be correctly updated")
			assert.Equal(t, tt.expectedTeam.Morale, team.Morale, "Morale should be correctly updated")
			assert.Equal(t, tt.expectedTeam.AttackPower, team.AttackPower, "Attack power should remain unchanged")
			assert.Equal(t, tt.expectedTeam.DefensePower, team.DefensePower, "Defense power should remain unchanged")

			// Assert standings updates
			assert.Equal(t, tt.matchOutcome.WinnerGoals, standings.Goals, "Goals should be updated")
			assert.Equal(t, tt.matchOutcome.LoserGoals, standings.Against, "Goals against should be updated")
			assert.Equal(t, 1, standings.Played, "Played should be incremented")
			assert.Equal(t, 1, standings.Wins, "Wins should be incremented")
			assert.Equal(t, 0, standings.Losses, "Losses should remain unchanged")
			assert.Equal(t, 3, standings.Points, "Points should be incremented by 3")
		})
	}
}

func TestDrawTeamAttributeChanging(t *testing.T) {
	tests := []struct {
		name         string
		initialTeam  models.Team
		standings    models.Standings
		matchOutcome models.MatchOutcome
		expectedTeam models.Team
	}{
		{
			name: "Normal draw",
			initialTeam: models.Team{
				Name:         "Draw Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      80,
				Morale:       80,
			},
			standings: models.Standings{
				Team:    models.Team{Name: "Draw Team"},
				Points:  0,
				Played:  0,
				Wins:    0,
				Losses:  0,
				Goals:   0,
				Against: 0,
			},
			matchOutcome: models.MatchOutcome{
				IsDraw:      true,
				WinnerGoals: 2, // In draw, both teams have same goals
				LoserGoals:  2,
			},
			expectedTeam: models.Team{
				Name:         "Draw Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      75, // 80 - 5
				Morale:       80, // Unchanged in draw
			},
		},
		{
			name: "Low stamina team in draw",
			initialTeam: models.Team{
				Name:         "Low Stamina Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      3, // Will go below 0
				Morale:       80,
			},
			standings: models.Standings{
				Team:   models.Team{Name: "Low Stamina Team"},
				Points: 0,
			},
			matchOutcome: models.MatchOutcome{
				IsDraw:      true,
				WinnerGoals: 1,
				LoserGoals:  1,
			},
			expectedTeam: models.Team{
				Name:         "Low Stamina Team",
				AttackPower:  80,
				DefensePower: 80,
				Stamina:      0,  // Should be capped at 0
				Morale:       80, // Unchanged in draw
			},
		},
		{
			name: "Zero-zero draw",
			initialTeam: models.Team{
				Name:         "Boring Team",
				AttackPower:  70,
				DefensePower: 90,
				Stamina:      85,
				Morale:       75,
			},
			standings: models.Standings{
				Team:   models.Team{Name: "Boring Team"},
				Points: 3, // Already has some points
				Played: 0, // Start with 0 played
			},
			matchOutcome: models.MatchOutcome{
				IsDraw:      true,
				WinnerGoals: 0,
				LoserGoals:  0,
			},
			expectedTeam: models.Team{
				Name:         "Boring Team",
				AttackPower:  70,
				DefensePower: 90,
				Stamina:      80, // 85 - 5
				Morale:       75, // Unchanged
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create copies to avoid modifying original data
			team := tt.initialTeam
			standings := tt.standings
			initialPoints := standings.Points

			// Execute
			DrawTeamAttributeChanging(&standings, &team, tt.matchOutcome)

			// Assert team attributes
			assert.Equal(t, tt.expectedTeam.Stamina, team.Stamina, "Stamina should be correctly updated")
			assert.Equal(t, tt.expectedTeam.Morale, team.Morale, "Morale should remain unchanged in draw")
			assert.Equal(t, tt.expectedTeam.AttackPower, team.AttackPower, "Attack power should remain unchanged")
			assert.Equal(t, tt.expectedTeam.DefensePower, team.DefensePower, "Defense power should remain unchanged")

			// Assert standings updates
			assert.Equal(t, tt.matchOutcome.LoserGoals, standings.Goals, "Goals should be updated")
			assert.Equal(t, tt.matchOutcome.WinnerGoals, standings.Against, "Goals against should be updated")
			assert.Equal(t, 1, standings.Played, "Played should be incremented")
			assert.Equal(t, initialPoints+1, standings.Points, "Points should be incremented by 1 for draw")
			// Note: Wins and Losses are not updated in draw function
		})
	}
}

func TestAttributeChanging_EdgeCases(t *testing.T) {
	t.Run("Zero stamina and morale team losing", func(t *testing.T) {
		team := models.Team{
			Name:         "Zero Team",
			AttackPower:  50,
			DefensePower: 50,
			Stamina:      0,
			Morale:       0,
		}
		standings := models.Standings{
			Team:   models.Team{Name: "Zero Team"},
			Points: 0,
		}
		matchOutcome := models.MatchOutcome{
			WinnerGoals: 1,
			LoserGoals:  0,
		}

		LoserTeamAttributeChanging(&standings, &team, matchOutcome)

		assert.Equal(t, 0.0, team.Stamina, "Stamina should remain 0")
		assert.Equal(t, 0.0, team.Morale, "Morale should remain 0")
	})

	t.Run("Maximum morale team winning", func(t *testing.T) {
		team := models.Team{
			Name:         "Max Team",
			AttackPower:  100,
			DefensePower: 100,
			Stamina:      100,
			Morale:       100,
		}
		standings := models.Standings{
			Team:   models.Team{Name: "Max Team"},
			Points: 0,
		}
		matchOutcome := models.MatchOutcome{
			WinnerGoals: 3,
			LoserGoals:  0,
		}

		WinnerTeamAttributeChanging(&standings, &team, matchOutcome)

		assert.Equal(t, 95.0, team.Stamina, "Stamina should decrease by 5")
		assert.Equal(t, 100.0, team.Morale, "Morale should be capped at 100")
	})
}

// Benchmark tests
func BenchmarkLoserTeamAttributeChanging(b *testing.B) {
	team := models.Team{
		Name:         "Benchmark Team",
		AttackPower:  80,
		DefensePower: 80,
		Stamina:      80,
		Morale:       80,
	}
	standings := models.Standings{
		Team:   models.Team{Name: "Benchmark Team"},
		Points: 0,
	}
	matchOutcome := models.MatchOutcome{
		WinnerGoals: 2,
		LoserGoals:  1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset values for each iteration
		teamCopy := team
		standingsCopy := standings
		LoserTeamAttributeChanging(&standingsCopy, &teamCopy, matchOutcome)
	}
}

func BenchmarkWinnerTeamAttributeChanging(b *testing.B) {
	team := models.Team{
		Name:         "Benchmark Team",
		AttackPower:  80,
		DefensePower: 80,
		Stamina:      80,
		Morale:       80,
	}
	standings := models.Standings{
		Team:   models.Team{Name: "Benchmark Team"},
		Points: 0,
	}
	matchOutcome := models.MatchOutcome{
		WinnerGoals: 2,
		LoserGoals:  1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset values for each iteration
		teamCopy := team
		standingsCopy := standings
		WinnerTeamAttributeChanging(&standingsCopy, &teamCopy, matchOutcome)
	}
}

func BenchmarkDrawTeamAttributeChanging(b *testing.B) {
	team := models.Team{
		Name:         "Benchmark Team",
		AttackPower:  80,
		DefensePower: 80,
		Stamina:      80,
		Morale:       80,
	}
	standings := models.Standings{
		Team:   models.Team{Name: "Benchmark Team"},
		Points: 0,
	}
	matchOutcome := models.MatchOutcome{
		IsDraw:      true,
		WinnerGoals: 1,
		LoserGoals:  1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset values for each iteration
		teamCopy := team
		standingsCopy := standings
		DrawTeamAttributeChanging(&standingsCopy, &teamCopy, matchOutcome)
	}
}
