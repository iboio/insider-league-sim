package interfaces

import (
	"testing"

	"league-sim/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestMockSimulationServiceInterface_ImplementsInterface(t *testing.T) {
	// Test that MockSimulationServiceInterface implements SimulationServiceInterface
	var _ SimulationServiceInterface = (*MockSimulationServiceInterface)(nil)
	assert.True(t, true, "MockSimulationServiceInterface implements SimulationServiceInterface interface")
}

func TestMockSimulationServiceInterface_Simulation(t *testing.T) {
	// Create mock
	mockService := &MockSimulationServiceInterface{}

	// Setup expectations
	testLeagueID := "test-league-id"
	testPlayAllFixture := true
	expectedResponse := models.SimulationResponse{
		Matches: []models.MatchResult{
			{
				WeekNumber: 1,
				Home:       "Team A",
				HomeScore:  2,
				Away:       "Team B",
				AwayScore:  1,
				Winner:     "Team A",
			},
		},
		UpcomingFixtures: []models.Week{},
		PlayedFixtures:   []models.Week{},
	}

	mockService.On("Simulation", testLeagueID, testPlayAllFixture).Return(expectedResponse, nil)

	// Call method
	result, err := mockService.Simulation(testLeagueID, testPlayAllFixture)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	assert.Len(t, result.Matches, 1)
	assert.Equal(t, "Team A", result.Matches[0].Home)
	assert.Equal(t, "Team B", result.Matches[0].Away)
	mockService.AssertExpectations(t)
}

func TestMockSimulationServiceInterface_EditMatch(t *testing.T) {
	// Create mock
	mockService := &MockSimulationServiceInterface{}

	// Setup expectations
	testData := models.EditMatchResult{
		WeekNumber:   1,
		LeagueId:     "test-league-id",
		TeamName:     "Team A",
		TeamType:     "home",
		AgainstTeam:  "Team B",
		TeamOldGoals: 1,
		Goals:        2,
		IsDraw:       false,
	}

	mockService.On("EditMatch", testData).Return(nil)

	// Call method
	err := mockService.EditMatch(testData)

	// Assert
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}
