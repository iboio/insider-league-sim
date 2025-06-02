package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appContext "league-sim/internal/appContext/appContexts"
	leagueInterfaces "league-sim/internal/league/interfaces"
	"league-sim/internal/models"
	predictInterfaces "league-sim/internal/predict/interfaces"
	"league-sim/internal/repositories/interfaces"
	simulationInterfaces "league-sim/internal/simulation/interfaces"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAppContext for testing
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

// MockService for testing
type MockService struct {
	mock.Mock
}

func (m *MockService) LeagueService() leagueInterfaces.LeagueServiceInterface {
	args := m.Called()
	return args.Get(0).(leagueInterfaces.LeagueServiceInterface)
}

func (m *MockService) SimulationService() simulationInterfaces.SimulationServiceInterface {
	args := m.Called()
	return args.Get(0).(simulationInterfaces.SimulationServiceInterface)
}

func (m *MockService) PredictService() predictInterfaces.PredictServiceInterface {
	args := m.Called()
	return args.Get(0).(predictInterfaces.PredictServiceInterface)
}

func TestGetLeagueIds_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock data
	expectedLeagueIds := []models.GetLeaguesIdsWithNameResponse{
		{LeagueId: "league1", LeagueName: "League 1"},
		{LeagueId: "league2", LeagueName: "League 2"},
	}
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Configure mocks
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)
	mockLeagueRepo.On("GetLeague").Return(expectedLeagueIds, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetLeagueIds(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []models.GetLeaguesIdsWithNameResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedLeagueIds, response)

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockLeagueRepo.AssertExpectations(t)
}

func TestGetLeagueIds_MissingAppContext(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set context with nil value to trigger the error
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", nil)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute and expect panic due to type assertion
	assert.Panics(
		t, func() {
			GetLeagueIds(c)
		}, "Should panic when app context is nil")
}

func TestGetLeagueIds_RepositoryError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock data
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockAppCtx := &MockAppContext{}
	expectedError := errors.New("database connection failed")

	// Configure mocks
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)
	mockLeagueRepo.On("GetLeague").Return([]models.GetLeaguesIdsWithNameResponse{}, expectedError)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetLeagueIds(c)

	// Assert
	assert.Error(t, err)
	httpError := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, httpError.Code)
	assert.Equal(t, "Failed to get league IDs", httpError.Message)

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockLeagueRepo.AssertExpectations(t)
}

func TestGetLeague_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league-id")

	// Execute
	err := GetLeague(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `"test-league-id"`, strings.TrimSpace(rec.Body.String()))
}

func TestCreateLeague_Success(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.CreateLeagueRequest{
		TeamCount:  "8",
		LeagueName: "Test League",
	}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock data
	expectedResponse := models.GetLeaguesIdsWithNameResponse{
		LeagueId:   "new-league-id",
		LeagueName: "Test League",
	}
	mockLeagueService := &leagueInterfaces.MockLeagueServiceInterface{}
	mockService := &MockService{}
	mockAppCtx := &MockAppContext{}

	// Configure mocks
	mockService.On("LeagueService").Return(mockLeagueService)
	mockLeagueService.On("CreateLeague", "8", "Test League").Return(expectedResponse, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := CreateLeague(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.GetLeaguesIdsWithNameResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedResponse, response)

	// Verify mocks
	mockService.AssertExpectations(t)
	mockLeagueService.AssertExpectations(t)
}

func TestCreateLeague_InvalidRequestBody(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := CreateLeague(c)

	// Assert
	assert.Error(t, err)
	httpError := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
}

func TestCreateLeague_MissingAppContext(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.CreateLeagueRequest{TeamCount: "8", LeagueName: "Test League"}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set context with nil value to trigger the error
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", nil)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute and expect panic due to type assertion
	assert.Panics(
		t, func() {
			CreateLeague(c)
		}, "Should panic when app context is nil")
}

func TestCreateLeague_ServiceError(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.CreateLeagueRequest{TeamCount: "8", LeagueName: "Test League"}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Mock data
	mockLeagueService := &leagueInterfaces.MockLeagueServiceInterface{}
	mockService := &MockService{}
	mockAppCtx := &MockAppContext{}
	expectedError := errors.New("failed to create league")

	// Configure mocks
	mockService.On("LeagueService").Return(mockLeagueService)
	mockLeagueService.On("CreateLeague", "8", "Test League").Return(
		models.GetLeaguesIdsWithNameResponse{},
		expectedError)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := CreateLeague(c)

	// Assert
	assert.Error(t, err)
	httpError := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, httpError.Code)

	// Verify mocks
	mockService.AssertExpectations(t)
	mockLeagueService.AssertExpectations(t)
}

func TestGetStanding_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league/standing", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	expectedStandings := []models.Standings{
		{Team: models.Team{Name: "Team A"}, Points: 9, Wins: 3},
		{Team: models.Team{Name: "Team B"}, Points: 6, Wins: 2},
	}
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Configure mocks
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeaguesStandings", "test-league").Return(expectedStandings, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetStanding(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []models.Standings
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedStandings, response)

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestGetStanding_RepositoryError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league/standing", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}
	expectedError := errors.New("database error")

	// Configure mocks
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeaguesStandings", "test-league").Return([]models.Standings{}, expectedError)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetStanding(c)

	// Assert
	assert.Error(t, err)
	httpError := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, httpError.Code)
	assert.Equal(t, "Failed to get standings", httpError.Message)

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestGetFixtures_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league/fixtures", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	expectedFixtures := models.GetActiveLeagueFixturesResponse{
		UpcomingFixtures: []models.Week{
			{
				Number:  1,
				Matches: []models.Match{{Home: &models.Team{Name: "Team A"}, Away: &models.Team{Name: "Team B"}}},
			},
		},
		PlayedFixtures: []models.Week{},
	}
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Configure mocks
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockActiveLeagueRepo.On("GetActiveLeaguesFixtures", "test-league").Return(expectedFixtures, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetFixtures(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.GetActiveLeagueFixturesResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedFixtures, response)

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockActiveLeagueRepo.AssertExpectations(t)
}

func TestDeleteLeague_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/league/test-league", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockAppCtx := &MockAppContext{}

	// Configure mocks
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)
	mockLeagueRepo.On("DeleteLeague", "test-league").Return(nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := DeleteLeague(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `"test-league"`, strings.TrimSpace(rec.Body.String()))

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockLeagueRepo.AssertExpectations(t)
}

func TestEditMatch_Success(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.EditMatchResult{
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		MatchWeek: 1,
		Winner:    "Team A",
	}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/league/test-league", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	mockSimulationService := &simulationInterfaces.MockSimulationServiceInterface{}
	mockService := &MockService{}

	// Configure mocks
	mockService.On("SimulationService").Return(mockSimulationService)
	expectedEditData := requestBody
	expectedEditData.LeagueId = "test-league"
	mockSimulationService.On("EditMatch", expectedEditData).Return(nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := EditMatch(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `"Standings updated successfully"`, strings.TrimSpace(rec.Body.String()))

	// Verify mocks
	mockService.AssertExpectations(t)
	mockSimulationService.AssertExpectations(t)
}

func TestGetPredictTable_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league/predict", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	expectedPredictTable := []models.PredictedStanding{
		{TeamName: "Team A", Odds: 75.5},
		{TeamName: "Team B", Odds: 24.5},
	}
	mockPredictService := &predictInterfaces.MockPredictServiceInterface{}
	mockService := &MockService{}

	// Configure mocks
	mockService.On("PredictService").Return(mockPredictService)
	mockPredictService.On("PredictChampionShipSession", "test-league").Return(expectedPredictTable, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetPredictTable(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []models.PredictedStanding
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedPredictTable, response)

	// Verify mocks
	mockService.AssertExpectations(t)
	mockPredictService.AssertExpectations(t)
}

func TestResetLeague_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league/test-league/reset", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	mockLeagueService := &leagueInterfaces.MockLeagueServiceInterface{}
	mockService := &MockService{}

	// Configure mocks
	mockService.On("LeagueService").Return(mockLeagueService)
	mockLeagueService.On("ResetLeague", "test-league").Return(nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := ResetLeague(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `"League reset successfully"`, strings.TrimSpace(rec.Body.String()))

	// Verify mocks
	mockService.AssertExpectations(t)
	mockLeagueService.AssertExpectations(t)
}
