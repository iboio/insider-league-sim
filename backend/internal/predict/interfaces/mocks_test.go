package interfaces

import (
	"testing"

	"league-sim/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestMockPredictServiceInterface_ImplementsInterface(t *testing.T) {
	// Test that MockPredictServiceInterface implements PredictServiceInterface
	var _ PredictServiceInterface = (*MockPredictServiceInterface)(nil)
	assert.True(t, true, "MockPredictServiceInterface implements PredictServiceInterface interface")
}

func TestMockPredictServiceInterface_PredictChampionShipSession(t *testing.T) {
	// Create mock
	mockService := &MockPredictServiceInterface{}

	// Setup expectations
	testID := "test-league-id"
	expectedResponse := []models.PredictedStanding{
		{
			TeamName: "Team A",
			Points:   90,
		},
		{
			TeamName: "Team B",
			Points:   85,
		},
	}

	mockService.On("PredictChampionShipSession", testID).Return(expectedResponse, nil)

	// Call method
	result, err := mockService.PredictChampionShipSession(testID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Team A", result[0].TeamName)
	assert.Equal(t, 90, result[0].Points)
	mockService.AssertExpectations(t)
}
