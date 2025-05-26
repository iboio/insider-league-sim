package interfaces

import (
	"testing"

	"league-sim/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestMockLeagueRepository_ImplementsInterface(t *testing.T) {
	// Test that MockLeagueRepository implements LeagueRepository interface
	var _ LeagueRepository = (*MockLeagueRepository)(nil)
	assert.True(t, true, "MockLeagueRepository implements LeagueRepository interface")
}

func TestMockActiveLeagueRepository_ImplementsInterface(t *testing.T) {
	// Test that MockActiveLeagueRepository implements ActiveLeagueRepository interface
	var _ ActiveLeagueRepository = (*MockActiveLeagueRepository)(nil)
	assert.True(t, true, "MockActiveLeagueRepository implements ActiveLeagueRepository interface")
}

func TestMockMatchResultRepository_ImplementsInterface(t *testing.T) {
	// Test that MockMatchResultRepository implements MatchResultRepository interface
	var _ MatchResultRepository = (*MockMatchResultRepository)(nil)
	assert.True(t, true, "MockMatchResultRepository implements MatchResultRepository interface")
}

func TestMockLeagueRepository_SetLeague(t *testing.T) {
	// Create mock
	mockRepo := &MockLeagueRepository{}

	// Setup expectations
	testID := "test-id"
	testData := models.CreateLeagueRequest{} // Assuming this struct exists

	mockRepo.On("SetLeague", testID, testData).Return(nil)

	// Call method
	err := mockRepo.SetLeague(testID, testData)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockLeagueRepository_GetLeague(t *testing.T) {
	// Create mock
	mockRepo := &MockLeagueRepository{}

	// Setup expectations
	expectedResult := []models.GetLeaguesIdsWithNameResponse{}
	mockRepo.On("GetLeague").Return(expectedResult, nil)

	// Call method
	result, err := mockRepo.GetLeague()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
}

func TestMockLeagueRepository_DeleteLeague(t *testing.T) {
	// Create mock
	mockRepo := &MockLeagueRepository{}

	// Setup expectations
	testID := "test-id"
	mockRepo.On("DeleteLeague", testID).Return(nil)

	// Call method
	err := mockRepo.DeleteLeague(testID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockActiveLeagueRepository_GetActiveLeague(t *testing.T) {
	// Create mock
	mockRepo := &MockActiveLeagueRepository{}

	// Setup expectations
	testID := "test-id"
	expectedLeague := models.League{}
	mockRepo.On("GetActiveLeague", testID).Return(expectedLeague, nil)

	// Call method
	result, err := mockRepo.GetActiveLeague(testID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedLeague, result)
	mockRepo.AssertExpectations(t)
}

func TestMockActiveLeagueRepository_SetActiveLeague(t *testing.T) {
	// Create mock
	mockRepo := &MockActiveLeagueRepository{}

	// Setup expectations
	testData := models.League{}
	mockRepo.On("SetActiveLeague", testData).Return(nil)

	// Call method
	err := mockRepo.SetActiveLeague(testData)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockMatchResultRepository_EditMatchScore(t *testing.T) {
	// Create mock
	mockRepo := &MockMatchResultRepository{}

	// Setup expectations
	testData := models.EditMatchResult{}
	mockRepo.On("EditMatchScore", testData).Return(nil)

	// Call method
	err := mockRepo.EditMatchScore(testData)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockMatchResultRepository_SetMatchResults(t *testing.T) {
	// Create mock
	mockRepo := &MockMatchResultRepository{}

	// Setup expectations
	testLeagueID := "test-league-id"
	testResults := []models.MatchResult{}
	mockRepo.On("SetMatchResults", testLeagueID, testResults).Return(nil)

	// Call method
	err := mockRepo.SetMatchResults(testLeagueID, testResults)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockMatchResultRepository_GetMatchResults(t *testing.T) {
	// Create mock
	mockRepo := &MockMatchResultRepository{}

	// Setup expectations
	testLeagueID := "test-league-id"
	expectedResults := []models.MatchResult{}
	mockRepo.On("GetMatchResults", testLeagueID).Return(expectedResults, nil)

	// Call method
	result, err := mockRepo.GetMatchResults(testLeagueID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, result)
	mockRepo.AssertExpectations(t)
}

func TestMockMatchResultRepository_DeleteMatchResults(t *testing.T) {
	// Create mock
	mockRepo := &MockMatchResultRepository{}

	// Setup expectations
	testLeagueID := "test-league-id"
	mockRepo.On("DeleteMatchResults", testLeagueID).Return(nil)

	// Call method
	err := mockRepo.DeleteMatchResults(testLeagueID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
