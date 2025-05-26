package interfaces

import (
	"testing"

	"league-sim/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestMockLeagueServiceInterface_ImplementsInterface(t *testing.T) {
	// Test that MockLeagueServiceInterface implements LeagueServiceInterface
	var _ LeagueServiceInterface = (*MockLeagueServiceInterface)(nil)
	assert.True(t, true, "MockLeagueServiceInterface implements LeagueServiceInterface interface")
}

func TestMockLeagueServiceInterface_CreateLeague(t *testing.T) {
	// Create mock
	mockService := &MockLeagueServiceInterface{}

	// Setup expectations
	testN := "4"
	testLeagueName := "Test League"
	expectedResponse := models.GetLeaguesIdsWithNameResponse{
		LeagueId:   "test-id",
		LeagueName: testLeagueName,
	}

	mockService.On("CreateLeague", testN, testLeagueName).Return(expectedResponse, nil)

	// Call method
	result, err := mockService.CreateLeague(testN, testLeagueName)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	assert.Equal(t, "test-id", result.LeagueId)
	assert.Equal(t, testLeagueName, result.LeagueName)
	mockService.AssertExpectations(t)
}

func TestMockLeagueServiceInterface_ResetLeague(t *testing.T) {
	// Create mock
	mockService := &MockLeagueServiceInterface{}

	// Setup expectations
	testLeagueID := "test-league-id"
	mockService.On("ResetLeague", testLeagueID).Return(nil)

	// Call method
	err := mockService.ResetLeague(testLeagueID)

	// Assert
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}
