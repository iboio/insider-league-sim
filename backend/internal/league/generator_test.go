package league

import (
	"testing"

	"league-sim/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestTeamGenerate(t *testing.T) {
	tests := []struct {
		name          string
		numberOfTeams int
	}{
		{
			name:          "Generate 4 teams",
			numberOfTeams: 4,
		},
		{
			name:          "Generate 8 teams",
			numberOfTeams: 8,
		},
		{
			name:          "Generate 16 teams",
			numberOfTeams: 16,
		},
		{
			name:          "Generate 1 team",
			numberOfTeams: 1,
		},
		{
			name:          "Generate 0 teams",
			numberOfTeams: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teams := TeamGenerate(tt.numberOfTeams)

			// Check correct number of teams generated
			assert.Len(t, teams, tt.numberOfTeams, "Should generate correct number of teams")

			// Check each team has proper values
			for i, team := range teams {
				// Check team name follows pattern
				expectedName := string(rune('A' + i))
				assert.Equal(t, "Team "+expectedName, team.Name, "Team name should follow pattern")

				// Check all stats are within expected range (70-100)
				assert.GreaterOrEqual(t, team.AttackPower, 70.0, "Attack power should be >= 70")
				assert.LessOrEqual(t, team.AttackPower, 100.0, "Attack power should be <= 100")

				assert.GreaterOrEqual(t, team.DefensePower, 70.0, "Defense power should be >= 70")
				assert.LessOrEqual(t, team.DefensePower, 100.0, "Defense power should be <= 100")

				assert.GreaterOrEqual(t, team.Stamina, 70.0, "Stamina should be >= 70")
				assert.LessOrEqual(t, team.Stamina, 100.0, "Stamina should be <= 100")

				assert.GreaterOrEqual(t, team.Morale, 70.0, "Morale should be >= 70")
				assert.LessOrEqual(t, team.Morale, 100.0, "Morale should be <= 100")
			}
		})
	}
}

func TestTeamGenerate_TeamNames(t *testing.T) {
	// Test that team names are generated correctly
	teams := TeamGenerate(5)

	expectedNames := []string{"Team A", "Team B", "Team C", "Team D", "Team E"}

	for i, team := range teams {
		assert.Equal(t, expectedNames[i], team.Name, "Team name should match expected pattern")
	}
}

func TestTeamGenerate_StatsVariation(t *testing.T) {
	// Test that teams have different stats (not all identical)
	teams := TeamGenerate(10)

	if len(teams) < 2 {
		t.Skip("Need at least 2 teams to test variation")
	}

	// Check that not all teams have identical stats
	firstTeam := teams[0]
	hasVariation := false

	for i := 1; i < len(teams); i++ {
		if teams[i].AttackPower != firstTeam.AttackPower ||
			teams[i].DefensePower != firstTeam.DefensePower ||
			teams[i].Stamina != firstTeam.Stamina ||
			teams[i].Morale != firstTeam.Morale {
			hasVariation = true
			break
		}
	}

	assert.True(t, hasVariation, "Teams should have different stats due to randomization")
}

func TestCreateStandingsTable(t *testing.T) {
	tests := []struct {
		name  string
		teams []models.Team
	}{
		{
			name: "Create standings for 4 teams",
			teams: []models.Team{
				{Name: "Team A", AttackPower: 80, DefensePower: 75, Stamina: 85, Morale: 90},
				{Name: "Team B", AttackPower: 85, DefensePower: 80, Stamina: 75, Morale: 85},
				{Name: "Team C", AttackPower: 75, DefensePower: 85, Stamina: 90, Morale: 80},
				{Name: "Team D", AttackPower: 90, DefensePower: 90, Stamina: 80, Morale: 75},
			},
		},
		{
			name: "Create standings for 1 team",
			teams: []models.Team{
				{Name: "Solo Team", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
			},
		},
		{
			name:  "Create standings for empty teams",
			teams: []models.Team{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			standings := CreateStandingsTable(tt.teams)

			// Check correct number of standings entries
			assert.Len(t, standings, len(tt.teams), "Should create standings entry for each team")

			// Check each standings entry
			for i, standing := range standings {
				if i < len(tt.teams) {
					// Check team reference
					assert.Equal(t, tt.teams[i], standing.Team, "Standing should reference correct team")
				}

				// Check initial values are zero
				assert.Equal(t, 0, standing.Goals, "Initial goals should be 0")
				assert.Equal(t, 0, standing.Against, "Initial goals against should be 0")
				assert.Equal(t, 0, standing.Played, "Initial games played should be 0")
				assert.Equal(t, 0, standing.Wins, "Initial wins should be 0")
				assert.Equal(t, 0, standing.Losses, "Initial losses should be 0")
				assert.Equal(t, 0, standing.Points, "Initial points should be 0")
			}
		})
	}
}

func TestGenerateFixtures(t *testing.T) {
	tests := []struct {
		name          string
		teams         []models.Team
		expectedWeeks int
	}{
		{
			name: "Generate fixtures for 4 teams",
			teams: []models.Team{
				{Name: "Team A"}, {Name: "Team B"}, {Name: "Team C"}, {Name: "Team D"},
			},
			expectedWeeks: 3, // n-1 weeks for n teams
		},
		{
			name: "Generate fixtures for 6 teams",
			teams: []models.Team{
				{Name: "Team A"}, {Name: "Team B"}, {Name: "Team C"},
				{Name: "Team D"}, {Name: "Team E"}, {Name: "Team F"},
			},
			expectedWeeks: 5, // n-1 weeks for n teams
		},
		{
			name: "Generate fixtures for 2 teams",
			teams: []models.Team{
				{Name: "Team A"}, {Name: "Team B"},
			},
			expectedWeeks: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weeks := GenerateFixtures(tt.teams)

			// Check correct number of weeks
			assert.Len(t, weeks, tt.expectedWeeks, "Should generate correct number of weeks")

			// Check week numbering
			for i, week := range weeks {
				assert.Equal(t, i+1, week.Number, "Week numbers should be sequential starting from 1")
			}

			// Check that each team plays against every other team exactly once
			teamMatchCount := make(map[string]map[string]int)
			for _, team := range tt.teams {
				teamMatchCount[team.Name] = make(map[string]int)
			}

			for _, week := range weeks {
				for _, match := range week.Matches {
					// Skip BYE matches
					if match.Home.Name == "BYE" || match.Away.Name == "BYE" {
						continue
					}

					teamMatchCount[match.Home.Name][match.Away.Name]++
					teamMatchCount[match.Away.Name][match.Home.Name]++
				}
			}

			// Verify each team plays every other team exactly once
			for _, team1 := range tt.teams {
				for _, team2 := range tt.teams {
					if team1.Name != team2.Name {
						count := teamMatchCount[team1.Name][team2.Name]
						assert.Equal(t, 1, count, "Each team should play every other team exactly once")
					}
				}
			}
		})
	}
}

func TestGenerateFixtures_OddNumberOfTeams(t *testing.T) {
	// Test with odd number of teams (should add BYE team)
	teams := []models.Team{
		{Name: "Team A"}, {Name: "Team B"}, {Name: "Team C"},
	}

	weeks := GenerateFixtures(teams)

	// Should generate 3 weeks (n weeks for n teams when odd)
	assert.Len(t, weeks, 3, "Should generate 3 weeks for 3 teams")

	// Check that BYE team doesn't appear in actual matches
	for _, week := range weeks {
		for _, match := range week.Matches {
			assert.NotEqual(t, "BYE", match.Home.Name, "BYE team should not appear in matches")
			assert.NotEqual(t, "BYE", match.Away.Name, "BYE team should not appear in matches")
		}
	}
}

func TestGenerateFixtures_EmptyTeams(t *testing.T) {
	// Test with empty teams slice
	teams := []models.Team{}

	weeks := GenerateFixtures(teams)

	// With empty teams, after adding BYE team, n=1, so n-1=0 rounds, resulting in empty weeks
	assert.Len(t, weeks, 0, "Should generate 0 weeks for empty teams")
}

func TestGenerateFixtures_SingleTeam(t *testing.T) {
	// Test with single team
	teams := []models.Team{
		{Name: "Solo Team"},
	}

	weeks := GenerateFixtures(teams)

	// With 1 team, after adding BYE team, n=2, so n-1=1 round, resulting in 1 week
	assert.Len(t, weeks, 1, "Should generate 1 week for 1 team")

	for _, week := range weeks {
		assert.Len(t, week.Matches, 0, "Should have no matches when only one team (Solo vs BYE gets filtered out)")
	}
}

func TestGenerateFixtures_MatchIntegrity(t *testing.T) {
	// Test that no team plays against itself
	teams := TeamGenerate(6)
	weeks := GenerateFixtures(teams)

	for _, week := range weeks {
		for _, match := range week.Matches {
			assert.NotEqual(t, match.Home.Name, match.Away.Name, "Team should not play against itself")
		}
	}
}

// Benchmark tests
func BenchmarkTeamGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TeamGenerate(16)
	}
}

func BenchmarkCreateStandingsTable(b *testing.B) {
	teams := TeamGenerate(16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateStandingsTable(teams)
	}
}

func BenchmarkGenerateFixtures(b *testing.B) {
	teams := TeamGenerate(16)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateFixtures(teams)
	}
}
