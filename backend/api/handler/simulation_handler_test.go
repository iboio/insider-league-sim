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

	appContext "league-sim/internal/contexts/appContexts"
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
type MockAppContextSim struct {
	mock.Mock
}

func (m *MockAppContextSim) LeagueRepository() interfaces.LeagueRepository {
	args := m.Called()
	return args.Get(0).(interfaces.LeagueRepository)
}

func (m *MockAppContextSim) ActiveLeagueRepository() interfaces.ActiveLeagueRepository {
	args := m.Called()
	return args.Get(0).(interfaces.ActiveLeagueRepository)
}

func (m *MockAppContextSim) MatchResultRepository() interfaces.MatchResultRepository {
	args := m.Called()
	return args.Get(0).(interfaces.MatchResultRepository)
}

func (m *MockAppContextSim) DB() *appContext.DB {
	args := m.Called()
	return args.Get(0).(*appContext.DB)
}

// MockService for testing
type MockServiceSim struct {
	mock.Mock
}

func (m *MockServiceSim) LeagueService() leagueInterfaces.LeagueServiceInterface {
	args := m.Called()
	return args.Get(0).(leagueInterfaces.LeagueServiceInterface)
}

func (m *MockServiceSim) SimulationService() simulationInterfaces.SimulationServiceInterface {
	args := m.Called()
	return args.Get(0).(simulationInterfaces.SimulationServiceInterface)
}

func (m *MockServiceSim) PredictService() predictInterfaces.PredictServiceInterface {
	args := m.Called()
	return args.Get(0).(predictInterfaces.PredictServiceInterface)
}

func TestStartSimulation_Success(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.SimulateLeagueRequest{
		PlayAllFixture: false,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league/test-league/simulation", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
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
		PlayedFixtures: []models.Week{
			{Number: 1, Matches: []models.Match{}},
		},
	}
	mockSimulationService := &simulationInterfaces.MockSimulationServiceInterface{}
	mockService := &MockServiceSim{}
	mockAppCtx := &MockAppContextSim{}

	// Configure mocks
	mockService.On("SimulationService").Return(mockSimulationService)
	mockSimulationService.On("Simulation", "test-league", false).Return(expectedResponse, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := StartSimulation(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.SimulationResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedResponse, response)

	// Verify mocks
	mockService.AssertExpectations(t)
	mockSimulationService.AssertExpectations(t)
}

func TestStartSimulation_PlayAllFixtures(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.SimulateLeagueRequest{
		PlayAllFixture: true,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league/test-league/simulation", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	expectedResponse := models.SimulationResponse{
		Matches: []models.MatchResult{
			{WeekNumber: 1, Home: "Team A", HomeScore: 2, Away: "Team B", AwayScore: 1, Winner: "Team A"},
			{WeekNumber: 2, Home: "Team C", HomeScore: 1, Away: "Team D", AwayScore: 3, Winner: "Team D"},
		},
		UpcomingFixtures: []models.Week{},
		PlayedFixtures: []models.Week{
			{Number: 1, Matches: []models.Match{}},
			{Number: 2, Matches: []models.Match{}},
		},
	}
	mockSimulationService := &simulationInterfaces.MockSimulationServiceInterface{}
	mockService := &MockServiceSim{}
	mockAppCtx := &MockAppContextSim{}

	// Configure mocks
	mockService.On("SimulationService").Return(mockSimulationService)
	mockSimulationService.On("Simulation", "test-league", true).Return(expectedResponse, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := StartSimulation(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.SimulationResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedResponse, response)

	// Verify mocks
	mockService.AssertExpectations(t)
	mockSimulationService.AssertExpectations(t)
}

func TestStartSimulation_InvalidRequestBody(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league/test-league/simulation", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Execute
	err := StartSimulation(c)

	// Assert
	assert.Error(t, err)
	httpError := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
	assert.Equal(t, "Invalid request body", httpError.Message)
}

func TestStartSimulation_MissingServiceContext(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.SimulateLeagueRequest{PlayAllFixture: false}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league/test-league/simulation", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	mockAppCtx := &MockAppContextSim{}

	// Set context - only app context, no services (nil services will cause panic)
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	ctx = context.WithValue(ctx, "services", nil)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute and expect panic due to type assertion
	assert.Panics(t, func() {
		StartSimulation(c)
	}, "Should panic when service context is nil")
}

func TestStartSimulation_SimulationServiceError(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.SimulateLeagueRequest{PlayAllFixture: false}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league/test-league/simulation", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	mockSimulationService := &simulationInterfaces.MockSimulationServiceInterface{}
	mockService := &MockServiceSim{}
	mockAppCtx := &MockAppContextSim{}
	expectedError := errors.New("simulation failed")

	// Configure mocks
	mockService.On("SimulationService").Return(mockSimulationService)
	mockSimulationService.On("Simulation", "test-league", false).Return(models.SimulationResponse{}, expectedError)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := StartSimulation(c)

	// Assert
	assert.Error(t, err)
	httpError := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, httpError.Code)
	assert.Contains(t, httpError.Message.(string), "Failed to start simulation")

	// Verify mocks
	mockService.AssertExpectations(t)
	mockSimulationService.AssertExpectations(t)
}

func TestStartSimulation_EmptyResponse(t *testing.T) {
	// Setup
	e := echo.New()
	requestBody := models.SimulateLeagueRequest{PlayAllFixture: false}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/league/test-league/simulation", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	emptyResponse := models.SimulationResponse{
		Matches:          []models.MatchResult{},
		UpcomingFixtures: []models.Week{},
		PlayedFixtures:   []models.Week{},
	}
	mockSimulationService := &simulationInterfaces.MockSimulationServiceInterface{}
	mockService := &MockServiceSim{}
	mockAppCtx := &MockAppContextSim{}

	// Configure mocks
	mockService.On("SimulationService").Return(mockSimulationService)
	mockSimulationService.On("Simulation", "test-league", false).Return(emptyResponse, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	ctx = context.WithValue(ctx, "services", mockService)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := StartSimulation(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.SimulationResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, emptyResponse, response)
	assert.Empty(t, response.Matches)

	// Verify mocks
	mockService.AssertExpectations(t)
	mockSimulationService.AssertExpectations(t)
}

func TestGetMatchResults_Success(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league/match-results", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	expectedResults := []models.MatchResult{
		{
			WeekNumber: 1,
			Home:       "Team A",
			HomeScore:  2,
			Away:       "Team B",
			AwayScore:  1,
			Winner:     "Team A",
		},
		{
			WeekNumber: 1,
			Home:       "Team C",
			HomeScore:  0,
			Away:       "Team D",
			AwayScore:  3,
			Winner:     "Team D",
		},
	}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContextSim{}

	// Configure mocks
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo)
	mockMatchResultRepo.On("GetMatchResults", "test-league").Return(expectedResults, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetMatchResults(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []models.MatchResult
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, expectedResults, response)
	assert.Len(t, response, 2)

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestGetMatchResults_MissingAppContext(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league/match-results", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Set context with nil value to trigger the error
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", nil)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute and expect panic due to type assertion
	assert.Panics(t, func() {
		GetMatchResults(c)
	}, "Should panic when app context is nil")
}

func TestGetMatchResults_RepositoryError(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league/match-results", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContextSim{}
	expectedError := errors.New("database connection failed")

	// Configure mocks
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo)
	mockMatchResultRepo.On("GetMatchResults", "test-league").Return([]models.MatchResult{}, expectedError)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetMatchResults(c)

	// Assert
	assert.Error(t, err)
	httpError := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusInternalServerError, httpError.Code)

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}

func TestGetMatchResults_EmptyResults(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/league/test-league/match-results", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("leagueId")
	c.SetParamValues("test-league")

	// Mock data
	emptyResults := []models.MatchResult{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockAppCtx := &MockAppContextSim{}

	// Configure mocks
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo)
	mockMatchResultRepo.On("GetMatchResults", "test-league").Return(emptyResults, nil)

	// Set context
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, "appContext", mockAppCtx)
	c.SetRequest(c.Request().WithContext(ctx))

	// Execute
	err := GetMatchResults(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []models.MatchResult
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, emptyResults, response)
	assert.Empty(t, response)

	// Verify mocks
	mockAppCtx.AssertExpectations(t)
	mockMatchResultRepo.AssertExpectations(t)
}
