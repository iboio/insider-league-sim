package simulation

import (
	"errors"
	"testing"

	appContext "league-sim/internal/contexts/appContexts"
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

// Helper function to create a test AppContext with mock repositories
func createTestAppContext(activeLeagueRepo interfaces.ActiveLeagueRepository, matchResultRepo interfaces.MatchResultRepository) *MockAppContext {
	mockAppCtx := &MockAppContext{}
	mockAppCtx.On("ActiveLeagueRepository").Return(activeLeagueRepo)
	mockAppCtx.On("MatchResultRepository").Return(matchResultRepo)
	return mockAppCtx
}

func TestNewSimulationService(t *testing.T) {
	// Create mock app context
	mockAppCtx := &MockAppContext{}

	// Test NewSimulationService
	service := NewSimulationService(mockAppCtx)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, mockAppCtx, service.appCtx)
}

func TestSimulationService_Simulation_SingleWeek_Success(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data
	leagueId := "test-league-id"
	activeLeague := models.League{
		LeagueID:   leagueId,
		LeagueName: "Test League",
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
		},
		Standings: []models.Standings{
			{
				Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
				Points:  0,
				Played:  0,
				Wins:    0,
				Losses:  0,
				Goals:   0,
				Against: 0,
			},
			{
				Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
				Points:  0,
				Played:  0,
				Wins:    0,
				Losses:  0,
				Goals:   0,
				Against: 0,
			},
		},
		UpcomingFixtures: []models.Week{
			{
				Number: 1,
				Matches: []models.Match{
					{
						Home: &models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
						Away: &models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
					},
				},
			},
		},
		PlayedFixtures: []models.Week{},
		CurrentWeek:    0,
	}

	// Configure mock expectations
	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("SetMatchResults", leagueId, mock.AnythingOfType("[]models.MatchResult")).Return(nil)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute - simulate single week
	result, err := service.Simulation(leagueId, false)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.Matches, 1, "Should simulate one match")
	assert.Len(t, result.PlayedFixtures, 1, "Should have one played fixture")
	assert.Empty(t, result.UpcomingFixtures, "Should have no upcoming fixtures left")

	// Verify match result structure
	match := result.Matches[0]
	assert.Equal(t, 1, match.MatchWeek)
	assert.Equal(t, "Team A", match.Home)
	assert.Equal(t, "Team B", match.Away)
	assert.GreaterOrEqual(t, match.HomeScore, 0)
	assert.GreaterOrEqual(t, match.AwayScore, 0)
	assert.NotEmpty(t, match.Winner)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_Simulation_AllWeeks_Success(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data with multiple weeks
	leagueId := "test-league-id"
	activeLeague := models.League{
		LeagueID:   leagueId,
		LeagueName: "Test League",
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
			{Name: "Team C", AttackPower: 75, DefensePower: 85, Stamina: 75, Morale: 85},
		},
		Standings: []models.Standings{
			{Team: models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80}, Points: 0, Played: 0},
			{Team: models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75}, Points: 0, Played: 0},
			{Team: models.Team{Name: "Team C", AttackPower: 75, DefensePower: 85, Stamina: 75, Morale: 85}, Points: 0, Played: 0},
		},
		UpcomingFixtures: []models.Week{
			{
				Number: 1,
				Matches: []models.Match{
					{
						Home: &models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
						Away: &models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
					},
				},
			},
			{
				Number: 2,
				Matches: []models.Match{
					{
						Home: &models.Team{Name: "Team C", AttackPower: 75, DefensePower: 85, Stamina: 75, Morale: 85},
						Away: &models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
					},
				},
			},
		},
		PlayedFixtures: []models.Week{},
		CurrentWeek:    0,
	}

	// Configure mock expectations
	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("SetMatchResults", leagueId, mock.AnythingOfType("[]models.MatchResult")).Return(nil)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute - simulate all weeks
	result, err := service.Simulation(leagueId, true)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.Matches, 2, "Should simulate all matches")
	assert.Len(t, result.PlayedFixtures, 2, "Should have all fixtures played")
	assert.Empty(t, result.UpcomingFixtures, "Should have no upcoming fixtures left")

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_Simulation_GetActiveLeagueError(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	expectedError := errors.New("league not found")

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(models.League{}, expectedError)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute - should return error instead of panic
	result, err := service.Simulation(leagueId, false)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result.Matches)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestSimulationService_Simulation_EmptyUpcomingFixtures(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data - league with no upcoming fixtures
	leagueId := "test-league-id"
	activeLeague := models.League{
		LeagueID: leagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
		},
		Standings: []models.Standings{
			{Team: models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80}, Points: 0, Played: 0},
			{Team: models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75}, Points: 0, Played: 0},
		},
		UpcomingFixtures: []models.Week{}, // Empty upcoming fixtures
		PlayedFixtures:   []models.Week{},
	}

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(activeLeague, nil)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	result, err := service.Simulation(leagueId, false)

	// Assert - should return empty response without error
	assert.NoError(t, err)
	assert.Empty(t, result.Matches, "Should have no matches")
	assert.Empty(t, result.UpcomingFixtures, "Should have no upcoming fixtures")
	assert.Empty(t, result.PlayedFixtures, "Should have no played fixtures")

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestSimulationService_Simulation_SetActiveLeagueError(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	activeLeague := models.League{
		LeagueID: leagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
		},
		Standings: []models.Standings{
			{Team: models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80}, Points: 0, Played: 0},
			{Team: models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75}, Points: 0, Played: 0},
		},
		UpcomingFixtures: []models.Week{
			{
				Number: 1,
				Matches: []models.Match{
					{
						Home: &models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
						Away: &models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
					},
				},
			},
		},
		PlayedFixtures: []models.Week{},
	}
	expectedError := errors.New("failed to set active league")

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(expectedError)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute and expect panic
	assert.Panics(t, func() {
		service.Simulation(leagueId, false)
	}, "Should panic when SetActiveLeague fails")

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestSimulationService_EditMatch_Success_Win(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		MatchWeek: 1,
		Winner:    "Team A",
	}

	// Existing match result that will be returned by GetMatchResultByWeekAndTeam
	existingMatch := models.MatchResult{
		MatchWeek: 1,
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 1,
		AwayScore: 2,
		Winner:    "Team B",
	}

	activeLeague := models.League{
		LeagueID: editData.LeagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 75},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 70},
		},
		Standings: []models.Standings{
			{
				Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 75},
				Points:  0,
				Wins:    0,
				Losses:  1,
				Goals:   1,
				Against: 2,
			},
			{
				Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 70},
				Points:  3,
				Wins:    1,
				Losses:  0,
				Goals:   2,
				Against: 1,
			},
		},
	}

	// Configure mock expectations
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(existingMatch, nil)
	mockActiveLeagueRepo.On("GetActiveLeague", editData.LeagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("EditMatchScore", editData).Return(nil)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_EditMatch_Success_Draw(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 2,
		AwayScore: 2,
		MatchWeek: 1,
		Winner:    "draw",
	}

	// Existing match result that will be returned by GetMatchResultByWeekAndTeam
	existingMatch := models.MatchResult{
		MatchWeek: 1,
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 1,
		AwayScore: 2,
		Winner:    "Team B",
	}

	activeLeague := models.League{
		LeagueID: editData.LeagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 75},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 70},
		},
		Standings: []models.Standings{
			{
				Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 75},
				Points:  0,
				Wins:    0,
				Losses:  1,
				Goals:   1,
				Against: 2,
			},
			{
				Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 70},
				Points:  3,
				Wins:    1,
				Losses:  0,
				Goals:   2,
				Against: 1,
			},
		},
	}

	// Configure mock expectations
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(existingMatch, nil)
	mockActiveLeagueRepo.On("GetActiveLeague", editData.LeagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("EditMatchScore", editData).Return(nil)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_EditMatch_MoraleCapLimits(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		MatchWeek: 1,
		Winner:    "Team A",
	}

	// Existing match result that will be returned by GetMatchResultByWeekAndTeam
	existingMatch := models.MatchResult{
		MatchWeek: 1,
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 1,
		AwayScore: 2,
		Winner:    "Team B",
	}

	activeLeague := models.League{
		LeagueID: editData.LeagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 98}, // High morale to test cap
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 3},  // Low morale to test floor
		},
		Standings: []models.Standings{
			{
				Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 98},
				Points:  0,
				Wins:    0,
				Losses:  1,
				Goals:   1,
				Against: 2,
			},
			{
				Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 3},
				Points:  3,
				Wins:    1,
				Losses:  0,
				Goals:   2,
				Against: 1,
			},
		},
	}

	// Configure mock expectations
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(existingMatch, nil)
	mockActiveLeagueRepo.On("GetActiveLeague", editData.LeagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("EditMatchScore", editData).Return(nil)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_EditMatch_DetailedWinScenario(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data - Team A wins with 2 goals difference
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 4,
		AwayScore: 2,
		MatchWeek: 1,
		Winner:    "Team A",
	}

	// Existing match result that will be returned by GetMatchResultByWeekAndTeam
	existingMatch := models.MatchResult{
		MatchWeek: 1,
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 2,
		AwayScore: 3,
		Winner:    "Team B",
	}

	initialTeamAMorale := float64(60)
	initialTeamBMorale := float64(70)

	activeLeague := models.League{
		LeagueID: editData.LeagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: initialTeamAMorale},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: initialTeamBMorale},
		},
		Standings: []models.Standings{
			{
				Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: initialTeamAMorale},
				Points:  1,
				Wins:    0,
				Losses:  0,
				Goals:   2,
				Against: 3,
			},
			{
				Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: initialTeamBMorale},
				Points:  1,
				Wins:    0,
				Losses:  0,
				Goals:   3,
				Against: 2,
			},
		},
	}

	// Configure mock expectations
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(existingMatch, nil)
	mockActiveLeagueRepo.On("GetActiveLeague", editData.LeagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("EditMatchScore", editData).Return(nil)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_EditMatch_DetailedDrawScenario(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data - Draw scenario
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 2,
		AwayScore: 2,
		MatchWeek: 1,
		Winner:    "draw",
	}

	// Existing match result that will be returned by GetMatchResultByWeekAndTeam
	existingMatch := models.MatchResult{
		MatchWeek: 1,
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 1,
		AwayScore: 2,
		Winner:    "Team B",
	}

	initialTeamAMorale := float64(50)
	initialTeamBMorale := float64(80)

	activeLeague := models.League{
		LeagueID: editData.LeagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: initialTeamAMorale},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: initialTeamBMorale},
		},
		Standings: []models.Standings{
			{
				Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: initialTeamAMorale},
				Points:  0,
				Wins:    0,
				Losses:  1,
				Goals:   1,
				Against: 2,
			},
			{
				Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: initialTeamBMorale},
				Points:  3,
				Wins:    1,
				Losses:  0,
				Goals:   2,
				Against: 1,
			},
		},
	}

	// Configure mock expectations
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(existingMatch, nil)
	mockActiveLeagueRepo.On("GetActiveLeague", editData.LeagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("EditMatchScore", editData).Return(nil)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestGenerateMatchResult_HomeWins(t *testing.T) {
	// Test data - strong home team vs weak away team
	homeTeam := models.Team{
		Name:         "Strong Home",
		AttackPower:  95,
		DefensePower: 95,
		Stamina:      95,
		Morale:       95,
	}
	awayTeam := models.Team{
		Name:         "Weak Away",
		AttackPower:  20,
		DefensePower: 20,
		Stamina:      20,
		Morale:       20,
	}

	// Run multiple times to test distribution
	homeWins := 0
	totalTests := 100

	for i := 0; i < totalTests; i++ {
		result := GenerateMatchResult(homeTeam, awayTeam)

		// Verify basic structure
		assert.GreaterOrEqual(t, result.WinnerGoals, 0, "Winner goals should be non-negative")
		assert.GreaterOrEqual(t, result.LoserGoals, 0, "Loser goals should be non-negative")

		if !result.IsDraw {
			assert.NotEqual(t, result.Winner.Name, result.Loser.Name, "Winner and loser should be different")
			assert.GreaterOrEqual(t, result.WinnerGoals, result.LoserGoals, "Winner should have more or equal goals")
		} else {
			assert.Equal(t, result.WinnerGoals, result.LoserGoals, "Draw should have equal goals")
		}

		// Count home wins
		if !result.IsDraw && result.Winner.Name == homeTeam.Name {
			homeWins++
		}
	}

	// Strong home team should win more often (at least 50% considering 20% draw chance and home advantage)
	homeWinPercentage := float64(homeWins) / float64(totalTests) * 100
	assert.Greater(t, homeWinPercentage, 50.0, "Strong home team should win more often")
}

func TestGenerateMatchResult_ZeroStrengthTeams(t *testing.T) {
	// Test data - both teams with very low strength (not zero to avoid division by zero)
	homeTeam := models.Team{
		Name:         "Low Home",
		AttackPower:  1,
		DefensePower: 1,
		Stamina:      1,
		Morale:       1,
	}
	awayTeam := models.Team{
		Name:         "Low Away",
		AttackPower:  1,
		DefensePower: 1,
		Stamina:      1,
		Morale:       1,
	}

	// Run multiple times to test distribution
	draws := 0
	homeWins := 0
	totalTests := 100

	for i := 0; i < totalTests; i++ {
		result := GenerateMatchResult(homeTeam, awayTeam)

		// Verify basic structure
		assert.GreaterOrEqual(t, result.WinnerGoals, 0, "Winner goals should be non-negative")
		assert.GreaterOrEqual(t, result.LoserGoals, 0, "Loser goals should be non-negative")

		if result.IsDraw {
			draws++
			assert.Equal(t, result.WinnerGoals, result.LoserGoals, "Draw should have equal goals")
		} else {
			if result.Winner.Name == homeTeam.Name {
				homeWins++
			}
		}
	}

	// With very low strength teams, we should see some draws (around 20% due to draw chance)
	// and home team should still have slight advantage due to 1.05 multiplier
	drawPercentage := float64(draws) / float64(totalTests) * 100
	homeWinPercentage := float64(homeWins) / float64(totalTests) * 100

	assert.Greater(t, drawPercentage, 10.0, "Should have some draws")
	assert.Greater(t, homeWinPercentage, 25.0, "Home team should have some advantage even with low strength")
}

func TestGenerateMatchResult_EqualStrengthTeams(t *testing.T) {
	// Test data - equal strength teams
	homeTeam := models.Team{
		Name:         "Equal Home",
		AttackPower:  80,
		DefensePower: 80,
		Stamina:      80,
		Morale:       80,
	}
	awayTeam := models.Team{
		Name:         "Equal Away",
		AttackPower:  80,
		DefensePower: 80,
		Stamina:      80,
		Morale:       80,
	}

	// Run multiple times to test distribution
	homeWins := 0
	awayWins := 0
	draws := 0
	totalTests := 100

	for i := 0; i < totalTests; i++ {
		result := GenerateMatchResult(homeTeam, awayTeam)

		if result.IsDraw {
			draws++
		} else if result.Winner.Name == homeTeam.Name {
			homeWins++
		} else {
			awayWins++
		}
	}

	// With equal teams, home should still have slight advantage due to 1.05 multiplier
	// But results should be somewhat distributed
	assert.Greater(t, homeWins, 0, "Home team should win some matches")
	assert.Greater(t, awayWins, 0, "Away team should win some matches")

	// Home team should have advantage but not overwhelming
	// With 1.05 multiplier and 20% draw chance, expect around 35-45% home wins
	homeWinPercentage := float64(homeWins) / float64(totalTests) * 100
	assert.Greater(t, homeWinPercentage, 30.0, "Home team should have reasonable win rate")
	assert.Less(t, homeWinPercentage, 80.0, "Home advantage shouldn't be overwhelming")
}

// Benchmark tests
func BenchmarkSimulationService_Simulation(b *testing.B) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContext{}

	activeLeague := models.League{
		LeagueID: "test-league",
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
		},
		Standings: []models.Standings{
			{Team: models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80}, Points: 0, Played: 0},
			{Team: models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75}, Points: 0, Played: 0},
		},
		UpcomingFixtures: []models.Week{
			{
				Number: 1,
				Matches: []models.Match{
					{
						Home: &models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 80},
						Away: &models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 75},
					},
				},
			},
		},
		PlayedFixtures: []models.Week{},
	}

	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo)
	mockActiveLeagueRepo.On("GetActiveLeague", "test-league").Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("SetMatchResults", "test-league", mock.AnythingOfType("[]models.MatchResult")).Return(nil)

	service := NewSimulationService(mockAppCtx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.Simulation("test-league", false)
	}
}

func BenchmarkGenerateMatchResult(b *testing.B) {
	homeTeam := models.Team{
		Name:         "Home Team",
		AttackPower:  80,
		DefensePower: 80,
		Stamina:      80,
		Morale:       80,
	}
	awayTeam := models.Team{
		Name:         "Away Team",
		AttackPower:  85,
		DefensePower: 75,
		Stamina:      85,
		Morale:       75,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateMatchResult(homeTeam, awayTeam)
	}
}

func TestSimulationService_EditMatch_GetActiveLeagueError(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		MatchWeek: 1,
		Winner:    "Team A",
	}

	// Existing match result that will be returned by GetMatchResultByWeekAndTeam
	existingMatch := models.MatchResult{
		MatchWeek: 1,
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 1,
		AwayScore: 2,
		Winner:    "Team B",
	}
	expectedError := errors.New("league not found")

	// Configure mock expectations
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo)
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(existingMatch, nil)
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeague", editData.LeagueId).Return(models.League{}, expectedError)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_EditMatch_SetActiveLeagueError(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		MatchWeek: 1,
		Winner:    "Team A",
	}

	// Existing match result that will be returned by GetMatchResultByWeekAndTeam
	existingMatch := models.MatchResult{
		MatchWeek: 1,
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 1,
		AwayScore: 2,
		Winner:    "Team B",
	}

	activeLeague := models.League{
		LeagueID: editData.LeagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 75},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 70},
		},
		Standings: []models.Standings{
			{
				Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 75},
				Points:  0,
				Wins:    0,
				Losses:  1,
				Goals:   1,
				Against: 2,
			},
			{
				Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 70},
				Points:  3,
				Wins:    1,
				Losses:  0,
				Goals:   2,
				Against: 1,
			},
		},
	}
	expectedError := errors.New("failed to set active league")

	// Configure mock expectations
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(existingMatch, nil)
	mockActiveLeagueRepo.On("GetActiveLeague", editData.LeagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(expectedError)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert - The current implementation returns the error from SetActiveLeague
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_EditMatch_EditMatchScoreError(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		MatchWeek: 1,
		Winner:    "Team A",
	}

	// Existing match result that will be returned by GetMatchResultByWeekAndTeam
	existingMatch := models.MatchResult{
		MatchWeek: 1,
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 1,
		AwayScore: 2,
		Winner:    "Team B",
	}

	activeLeague := models.League{
		LeagueID: editData.LeagueId,
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 75},
			{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 70},
		},
		Standings: []models.Standings{
			{
				Team:    models.Team{Name: "Team A", AttackPower: 80, DefensePower: 80, Stamina: 80, Morale: 75},
				Points:  0,
				Wins:    0,
				Losses:  1,
				Goals:   1,
				Against: 2,
			},
			{
				Team:    models.Team{Name: "Team B", AttackPower: 85, DefensePower: 75, Stamina: 85, Morale: 70},
				Points:  3,
				Wins:    1,
				Losses:  0,
				Goals:   2,
				Against: 1,
			},
		},
	}
	expectedError := errors.New("failed to edit match score")

	// Configure mock expectations
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(existingMatch, nil)
	mockActiveLeagueRepo.On("GetActiveLeague", editData.LeagueId).Return(activeLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("EditMatchScore", editData).Return(expectedError)

	// Create test AppContext
	mockAppCtx := createTestAppContext(mockActiveLeagueRepo, mockMatchResultRepo)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestSimulationService_EditMatch_GetMatchResultError(t *testing.T) {
	// Setup mocks
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		MatchWeek: 1,
		Winner:    "Team A",
	}
	expectedError := errors.New("match not found")

	// Configure mock expectations
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo)
	mockMatchResultRepo.On("GetMatchResultByWeekAndTeam", editData).Return(models.MatchResult{}, expectedError)

	// Create service
	service := NewSimulationService(mockAppCtx)

	// Execute
	err := service.EditMatch(editData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}
