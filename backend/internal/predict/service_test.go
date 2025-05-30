package predict

import (
	"errors"
	"testing"

	"league-sim/config"
	appContext "league-sim/internal/appContext/appContexts"
	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAppContext is a mock implementation of AppContext for testing
type MockAppContext struct {
	mock.Mock
}

func (m *MockAppContext) LeagueRepository() interfaces.LeagueRepository {
	args := m.Called()
	return args.Get(0).(interfaces.LeagueRepository)
}

func (m *MockAppContext) ActiveLeagueRepository() interfaces.ActiveLeagueRepository {
	args := m.Called()
	return args.Get(0).(interfaces.ActiveLeagueRepository)
}

func (m *MockAppContext) MatchResultRepository() interfaces.MatchResultRepository {
	args := m.Called()
	return args.Get(0).(interfaces.MatchResultRepository)
}

func (m *MockAppContext) DB() *appContext.DB {
	args := m.Called()
	return args.Get(0).(*appContext.DB)
}

func TestNewPredictService(t *testing.T) {
	// Create mock app context
	mockAppCtx := &MockAppContext{}

	// Test NewPredictService
	service := NewPredictService(mockAppCtx)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, mockAppCtx, service.appCtx)
}

func TestPredict_PredictChampionShipSession_Success(t *testing.T) {
	// Setup config values for testing
	config.WeightPoints = 0.4
	config.WeightsStrength = 0.6

	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	standings := []models.Standings{
		{
			Team: models.Team{
				Name:         "Team A",
				AttackPower:  80.0,
				DefensePower: 80.0,
				Stamina:      80.0,
				Morale:       80.0,
			},
			Points:  9,
			Played:  3,
			Wins:    3,
			Losses:  0,
			Goals:   8,
			Against: 2,
		},
		{
			Team: models.Team{
				Name:         "Team B",
				AttackPower:  85.0,
				DefensePower: 75.0,
				Stamina:      80.0,
				Morale:       80.0,
			},
			Points:  6,
			Played:  3,
			Wins:    2,
			Losses:  1,
			Goals:   6,
			Against: 4,
		},
		{
			Team: models.Team{
				Name:         "Team C",
				AttackPower:  75.0,
				DefensePower: 85.0,
				Stamina:      80.0,
				Morale:       80.0,
			},
			Points:  3,
			Played:  3,
			Wins:    1,
			Losses:  2,
			Goals:   4,
			Against: 6,
		},
		{
			Team: models.Team{
				Name:         "Team D",
				AttackPower:  70.0,
				DefensePower: 70.0,
				Stamina:      70.0,
				Morale:       70.0,
			},
			Points:  0,
			Played:  3,
			Wins:    0,
			Losses:  3,
			Goals:   2,
			Against: 8,
		},
	}

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeaguesStandings", leagueId).Return(standings, nil)

	// Create service
	service := NewPredictService(mockAppCtx)

	// Execute
	result, err := service.PredictChampionShipSession(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 4, "Should return predictions for all teams")

	// Verify all teams are included
	teamNames := make(map[string]bool)
	totalOdds := 0.0
	for _, prediction := range result {
		teamNames[prediction.TeamName] = true
		totalOdds += prediction.Odds
		assert.GreaterOrEqual(t, prediction.Odds, 0.0, "Odds should be non-negative")
		assert.Greater(t, prediction.Strength, 0.0, "Strength should be positive")
	}

	// Verify all teams are present
	assert.True(t, teamNames["Team A"])
	assert.True(t, teamNames["Team B"])
	assert.True(t, teamNames["Team C"])
	assert.True(t, teamNames["Team D"])

	// Verify odds sum to approximately 100% (allowing for eliminated teams)
	assert.LessOrEqual(t, totalOdds, 100.0, "Total odds should not exceed 100%")

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestPredict_PredictChampionShipSession_AllMatchesPlayed(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data - all teams have played all matches (3 matches each in 4-team league)
	leagueId := "test-league-id"
	standings := []models.Standings{
		{
			Team: models.Team{
				Name:         "Team A",
				AttackPower:  80.0,
				DefensePower: 80.0,
				Stamina:      80.0,
				Morale:       80.0,
			},
			Points:  9,
			Played:  3, // All matches played (4-1)
			Wins:    3,
			Losses:  0,
			Goals:   8,
			Against: 2,
		},
		{
			Team: models.Team{
				Name:         "Team B",
				AttackPower:  85.0,
				DefensePower: 75.0,
				Stamina:      80.0,
				Morale:       80.0,
			},
			Points:  6,
			Played:  3,
			Wins:    2,
			Losses:  1,
			Goals:   6,
			Against: 4,
		},
	}

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeaguesStandings", leagueId).Return(standings, nil)

	// Create service
	service := NewPredictService(mockAppCtx)

	// Execute
	result, err := service.PredictChampionShipSession(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2, "Should return predictions for all teams")

	// Find the leader (Team A with 9 points)
	var leaderPrediction *models.PredictedStanding
	var otherPrediction *models.PredictedStanding

	for i := range result {
		if result[i].TeamName == "Team A" {
			leaderPrediction = &result[i]
		} else {
			otherPrediction = &result[i]
		}
	}

	// Verify leader has 100% odds and is not eliminated
	assert.NotNil(t, leaderPrediction)
	assert.Equal(t, 100.0, leaderPrediction.Odds, "Leader should have 100% odds when all matches are played")
	assert.False(t, leaderPrediction.Eliminated, "Leader should not be eliminated")

	// Verify other team has 0% odds and is eliminated
	assert.NotNil(t, otherPrediction)
	assert.Equal(t, 0.0, otherPrediction.Odds, "Non-leader should have 0% odds when all matches are played")
	assert.True(t, otherPrediction.Eliminated, "Non-leader should be eliminated")

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestPredict_PredictChampionShipSession_GetStandingsError(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	expectedError := errors.New("failed to get standings")

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeaguesStandings", leagueId).Return([]models.Standings{}, expectedError)

	// Create service
	service := NewPredictService(mockAppCtx)

	// Execute
	result, err := service.PredictChampionShipSession(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestPredict_PredictChampionShipSession_EmptyStandings(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	standings := []models.Standings{}

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeaguesStandings", leagueId).Return(standings, nil)

	// Create service
	service := NewPredictService(mockAppCtx)

	// Execute
	result, err := service.PredictChampionShipSession(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result, "Should return empty result for empty standings")

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestFindLeaderPoints(t *testing.T) {
	tests := []struct {
		name              string
		standings         []models.Standings
		expectedMaxPoints int
	}{
		{
			name: "Normal standings",
			standings: []models.Standings{
				{Points: 9}, {Points: 6}, {Points: 3}, {Points: 0},
			},
			expectedMaxPoints: 9,
		},
		{
			name: "All teams have same points",
			standings: []models.Standings{
				{Points: 6}, {Points: 6}, {Points: 6}, {Points: 6},
			},
			expectedMaxPoints: 6,
		},
		{
			name: "Single team",
			standings: []models.Standings{
				{Points: 12},
			},
			expectedMaxPoints: 12,
		},
		{
			name:              "Empty standings",
			standings:         []models.Standings{},
			expectedMaxPoints: 0,
		},
		{
			name: "Zero points",
			standings: []models.Standings{
				{Points: 0}, {Points: 0},
			},
			expectedMaxPoints: 0,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := findLeaderPoints(tt.standings)
				assert.Equal(t, tt.expectedMaxPoints, result, "Should return correct maximum points")
			})
	}
}

func TestFindLeader(t *testing.T) {
	tests := []struct {
		name           string
		standings      []models.Standings
		expectedLeader string
	}{
		{
			name: "Clear leader by points",
			standings: []models.Standings{
				{
					Team:    models.Team{Name: "Team A"},
					Points:  9,
					Goals:   8,
					Against: 2,
				},
				{
					Team:    models.Team{Name: "Team B"},
					Points:  6,
					Goals:   6,
					Against: 4,
				},
			},
			expectedLeader: "Team A",
		},
		{
			name: "Tie in points, leader by goal difference",
			standings: []models.Standings{
				{
					Team:    models.Team{Name: "Team A"},
					Points:  6,
					Goals:   5,
					Against: 3, // Goal diff: +2
				},
				{
					Team:    models.Team{Name: "Team B"},
					Points:  6,
					Goals:   8,
					Against: 2, // Goal diff: +6
				},
				{
					Team:    models.Team{Name: "Team C"},
					Points:  6,
					Goals:   4,
					Against: 4, // Goal diff: 0
				},
			},
			expectedLeader: "Team B", // Best goal difference
		},
		{
			name: "Single team",
			standings: []models.Standings{
				{
					Team:    models.Team{Name: "Solo Team"},
					Points:  3,
					Goals:   2,
					Against: 1,
				},
			},
			expectedLeader: "Solo Team",
		},
		{
			name: "All teams tied on points and goal difference",
			standings: []models.Standings{
				{
					Team:    models.Team{Name: "Team A"},
					Points:  3,
					Goals:   2,
					Against: 1, // Goal diff: +1
				},
				{
					Team:    models.Team{Name: "Team B"},
					Points:  3,
					Goals:   2,
					Against: 1, // Goal diff: +1
				},
			},
			expectedLeader: "Team A", // First team in case of complete tie
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				result := FindLeader(tt.standings)
				assert.Equal(t, tt.expectedLeader, result.Name, "Should return correct leader")
			})
	}
}

func TestFindLeader_EmptyStandings(t *testing.T) {
	// Test with empty standings
	standings := []models.Standings{}

	result := FindLeader(standings)

	// Should return empty team
	assert.Equal(t, "", result.Name, "Should return empty team for empty standings")
}

func BenchmarkPredict_PredictChampionShipSession(b *testing.B) {
	// Setup config values
	config.WeightPoints = 0.4
	config.WeightsStrength = 0.6

	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	standings := []models.Standings{
		{
			Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
			Points:  9,
			Played:  3,
			Goals:   8,
			Against: 2,
		},
		{
			Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 80, Morale: 80},
			Points:  6,
			Played:  3,
			Goals:   6,
			Against: 4,
		},
	}

	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeaguesStandings", "test-league").Return(standings, nil)

	service := NewPredictService(mockAppCtx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.PredictChampionShipSession("test-league")
	}
}

func BenchmarkFindLeader(b *testing.B) {
	standings := []models.Standings{
		{Team: models.Team{Name: "Team A"}, Points: 9, Goals: 8, Against: 2},
		{Team: models.Team{Name: "Team B"}, Points: 6, Goals: 6, Against: 4},
		{Team: models.Team{Name: "Team C"}, Points: 3, Goals: 4, Against: 6},
		{Team: models.Team{Name: "Team D"}, Points: 0, Goals: 2, Against: 8},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindLeader(standings)
	}
}
