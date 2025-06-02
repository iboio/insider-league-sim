package league

import (
	"errors"
	"testing"

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

func (m *MockAppContext) MatchResultRepository() interfaces.MatchesRepository {
	args := m.Called()
	return args.Get(0).(interfaces.MatchesRepository)
}

func (m *MockAppContext) DB() *appContext.DB {
	args := m.Called()
	return args.Get(0).(*appContext.DB)
}

func TestNewLeagueService(t *testing.T) {
	// Create mock app context
	mockAppCtx := &MockAppContext{}

	// Test NewLeagueService
	service := NewLeagueService(mockAppCtx)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, mockAppCtx, service.appCtx)
}

func TestLeagueService_CreateLeague_Success(t *testing.T) {
	// Setup mocks
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Configure mock expectations
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)

	mockLeagueRepo.On(
		"SetLeague",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("models.CreateLeagueRequest")).Return(nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)

	// Create service
	service := NewLeagueService(mockAppCtx)

	// Test data
	numberOfTeams := "4"
	leagueName := "Test League"

	// Execute
	result, err := service.CreateLeague(numberOfTeams, leagueName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, leagueName, result.LeagueName)
	assert.NotEmpty(t, result.LeagueId)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockLeagueRepo.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestLeagueService_CreateLeague_InvalidNumberFormat(t *testing.T) {
	// Setup mocks
	mockAppCtx := &MockAppContext{}

	// Create service
	service := NewLeagueService(mockAppCtx)

	// Test data with invalid number format
	numberOfTeams := "invalid"
	leagueName := "Test League"

	// Execute
	result, err := service.CreateLeague(numberOfTeams, leagueName)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result.LeagueName)
	assert.Empty(t, result.LeagueId)
}

func TestLeagueService_CreateLeague_LeagueRepositoryError(t *testing.T) {
	// Setup mocks
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Configure mock expectations
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)

	expectedError := errors.New("league repository error")
	mockLeagueRepo.On(
		"SetLeague",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("models.CreateLeagueRequest")).Return(expectedError)

	// Create service
	service := NewLeagueService(mockAppCtx)

	// Test data
	numberOfTeams := "4"
	leagueName := "Test League"

	// Execute
	result, err := service.CreateLeague(numberOfTeams, leagueName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result.LeagueName)
	assert.Empty(t, result.LeagueId)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockLeagueRepo.AssertExpectations(t)
}

func TestLeagueService_CreateLeague_ActiveLeagueRepositoryError(t *testing.T) {
	// Setup mocks
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Configure mock expectations
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)

	mockLeagueRepo.On(
		"SetLeague",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("models.CreateLeagueRequest")).Return(nil)

	expectedError := errors.New("active league repository error")
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(expectedError)

	// Create service
	service := NewLeagueService(mockAppCtx)

	// Test data
	numberOfTeams := "4"
	leagueName := "Test League"

	// Execute
	result, err := service.CreateLeague(numberOfTeams, leagueName)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result.LeagueName)
	assert.Empty(t, result.LeagueId)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockLeagueRepo.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestLeagueService_CreateLeague_DifferentTeamCounts(t *testing.T) {
	tests := []struct {
		name          string
		numberOfTeams string
		expectedTeams int
	}{
		{
			name:          "Create league with 2 teams",
			numberOfTeams: "2",
			expectedTeams: 2,
		},
		{
			name:          "Create league with 8 teams",
			numberOfTeams: "8",
			expectedTeams: 8,
		},
		{
			name:          "Create league with 16 teams",
			numberOfTeams: "16",
			expectedTeams: 16,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Setup mocks
				mockLeagueRepo := &interfaces.MockLeagueRepository{}
				mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
				mockAppCtx := &MockAppContext{}

				// Configure mock expectations
				mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)
				mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)

				mockLeagueRepo.On(
					"SetLeague",
					mock.AnythingOfType("string"),
					mock.AnythingOfType("models.CreateLeagueRequest")).Return(nil)

				// Capture the league data to verify team count
				var capturedLeague models.League
				mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Run(
					func(args mock.Arguments) {
						capturedLeague = args.Get(0).(models.League)
					}).Return(nil)

				// Create service
				service := NewLeagueService(mockAppCtx)

				// Execute
				result, err := service.CreateLeague(tt.numberOfTeams, "Test League")

				// Assert
				assert.NoError(t, err)
				assert.NotEmpty(t, result.LeagueId)
				assert.Equal(t, "Test League", result.LeagueName)

				// Verify the captured league has correct number of teams
				assert.Len(t, capturedLeague.Teams, tt.expectedTeams)
				assert.Len(t, capturedLeague.Standings, tt.expectedTeams)
				assert.Equal(t, 0, capturedLeague.CurrentWeek)
				assert.Empty(t, capturedLeague.PlayedFixtures)
				assert.NotEmpty(t, capturedLeague.UpcomingFixtures)

				// Verify mock expectations
				mockAppCtx.AssertExpectations(t)
				mockLeagueRepo.AssertExpectations(t)
				mockActiveLeagueRepo.AssertExpectations(t)
			})
	}
}

func TestLeagueService_ResetLeague_Success(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	existingLeague := models.League{
		LeagueID:   leagueId,
		LeagueName: "Test League",
		Teams: []models.Team{
			{Name: "Team A"}, {Name: "Team B"}, {Name: "Team C"}, {Name: "Team D"},
		},
		CurrentWeek: 5,
		PlayedFixtures: []models.Week{
			{Number: 1, Matches: []models.Match{}},
		},
	}

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockAppCtx.On("MatchesRepository").Return(mockMatchResultRepo)

	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(existingLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("DeleteMatchResults", leagueId).Return(nil)

	// Create service
	service := NewLeagueService(mockAppCtx)

	// Execute
	err := service.ResetLeague(leagueId)

	// Assert
	assert.NoError(t, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestLeagueService_ResetLeague_GetActiveLeagueError(t *testing.T) {
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
	service := NewLeagueService(mockAppCtx)

	// Execute
	err := service.ResetLeague(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestLeagueService_ResetLeague_SetActiveLeagueError(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	existingLeague := models.League{
		LeagueID:   leagueId,
		LeagueName: "Test League",
		Teams: []models.Team{
			{Name: "Team A"}, {Name: "Team B"},
		},
	}
	expectedError := errors.New("failed to set active league")

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)

	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(existingLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(expectedError)

	// Create service
	service := NewLeagueService(mockAppCtx)

	// Execute
	err := service.ResetLeague(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestLeagueService_ResetLeague_DeleteMatchResultsError(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	existingLeague := models.League{
		LeagueID:   leagueId,
		LeagueName: "Test League",
		Teams: []models.Team{
			{Name: "Team A"}, {Name: "Team B"},
		},
	}
	expectedError := errors.New("failed to delete match results")

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockAppCtx.On("MatchesRepository").Return(mockMatchResultRepo)

	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(existingLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("DeleteMatchResults", leagueId).Return(expectedError)

	// Create service
	service := NewLeagueService(mockAppCtx)

	// Execute
	err := service.ResetLeague(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestLeagueService_ResetLeague_VerifyResetData(t *testing.T) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContext{}

	// Test data
	leagueId := "test-league-id"
	existingLeague := models.League{
		LeagueID:   leagueId,
		LeagueName: "Test League",
		Teams: []models.Team{
			{Name: "Team A"}, {Name: "Team B"}, {Name: "Team C"}, {Name: "Team D"},
		},
		CurrentWeek: 10,
		PlayedFixtures: []models.Week{
			{Number: 1, Matches: []models.Match{}},
			{Number: 2, Matches: []models.Match{}},
		},
	}

	// Configure mock expectations
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockAppCtx.On("MatchesRepository").Return(mockMatchResultRepo)

	mockActiveLeagueRepo.On("GetActiveLeague", leagueId).Return(existingLeague, nil)

	// Capture the reset league data
	var capturedLeague models.League
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Run(
		func(args mock.Arguments) {
			capturedLeague = args.Get(0).(models.League)
		}).Return(nil)

	mockMatchResultRepo.On("DeleteMatchResults", leagueId).Return(nil)

	// Create service
	service := NewLeagueService(mockAppCtx)

	// Execute
	err := service.ResetLeague(leagueId)

	// Assert
	assert.NoError(t, err)

	// Verify the reset league data
	assert.Equal(t, leagueId, capturedLeague.LeagueID)
	assert.Equal(t, "Test League", capturedLeague.LeagueName)
	assert.Len(t, capturedLeague.Teams, 4)              // Same number of teams as original
	assert.Len(t, capturedLeague.Standings, 4)          // New standings for all teams
	assert.Equal(t, 0, capturedLeague.CurrentWeek)      // Reset to 0
	assert.Empty(t, capturedLeague.PlayedFixtures)      // Should be empty
	assert.NotEmpty(t, capturedLeague.UpcomingFixtures) // Should have new fixtures
	assert.Greater(t, capturedLeague.TotalWeeks, 0)     // Should have total weeks

	// Verify all teams have new random stats
	for _, team := range capturedLeague.Teams {
		assert.GreaterOrEqual(t, team.AttackPower, 70.0)
		assert.LessOrEqual(t, team.AttackPower, 100.0)
		assert.GreaterOrEqual(t, team.DefensePower, 70.0)
		assert.LessOrEqual(t, team.DefensePower, 100.0)
		assert.GreaterOrEqual(t, team.Stamina, 70.0)
		assert.LessOrEqual(t, team.Stamina, 100.0)
		assert.GreaterOrEqual(t, team.Morale, 70.0)
		assert.LessOrEqual(t, team.Morale, 100.0)
	}

	// Verify all standings are reset
	for _, standing := range capturedLeague.Standings {
		assert.Equal(t, 0, standing.Goals)
		assert.Equal(t, 0, standing.Against)
		assert.Equal(t, 0, standing.Played)
		assert.Equal(t, 0, standing.Wins)
		assert.Equal(t, 0, standing.Losses)
		assert.Equal(t, 0, standing.Points)
	}

	// Verify mock expectations
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

// Benchmark tests
func BenchmarkLeagueService_CreateLeague(b *testing.B) {
	// Setup mocks
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockLeagueRepo.On(
		"SetLeague",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("models.CreateLeagueRequest")).Return(nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)

	service := NewLeagueService(mockAppCtx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreateLeague("8", "Benchmark League")
	}
}

func BenchmarkLeagueService_ResetLeague(b *testing.B) {
	// Setup mocks
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContext{}

	existingLeague := models.League{
		LeagueID:   "test-league",
		LeagueName: "Test League",
		Teams: []models.Team{
			{Name: "Team A"}, {Name: "Team B"}, {Name: "Team C"}, {Name: "Team D"},
			{Name: "Team E"}, {Name: "Team F"}, {Name: "Team G"}, {Name: "Team H"},
		},
	}

	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockAppCtx.On("MatchesRepository").Return(mockMatchResultRepo)
	mockActiveLeagueRepo.On("GetActiveLeague", "test-league").Return(existingLeague, nil)
	mockActiveLeagueRepo.On("SetActiveLeague", mock.AnythingOfType("models.League")).Return(nil)
	mockMatchResultRepo.On("DeleteMatchResults", "test-league").Return(nil)

	service := NewLeagueService(mockAppCtx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ResetLeague("test-league")
	}
}
