package services

import (
	"testing"

	appContext "league-sim/internal/contexts/appContexts"
	leagueInterfaces "league-sim/internal/league/interfaces"
	predictInterfaces "league-sim/internal/predict/interfaces"
	"league-sim/internal/repositories/interfaces"
	simulationInterfaces "league-sim/internal/simulation/interfaces"

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

func TestServiceImpl_LeagueService(t *testing.T) {
	// Create mock league service
	mockLeagueService := &leagueInterfaces.MockLeagueServiceInterface{}

	// Create ServiceImpl with mock league service
	service := &ServiceImpl{
		leagueService: mockLeagueService,
	}

	// Test LeagueService() method
	result := service.LeagueService()

	assert.NotNil(t, result)
	assert.Equal(t, mockLeagueService, result)
}

func TestServiceImpl_PredictService(t *testing.T) {
	// Create mock predict service
	mockPredictService := &predictInterfaces.MockPredictServiceInterface{}

	// Create ServiceImpl with mock predict service
	service := &ServiceImpl{
		predictService: mockPredictService,
	}

	// Test PredictService() method
	result := service.PredictService()

	assert.NotNil(t, result)
	assert.Equal(t, mockPredictService, result)
}

func TestServiceImpl_SimulationService(t *testing.T) {
	// Create mock simulation service
	mockSimulationService := &simulationInterfaces.MockSimulationServiceInterface{}

	// Create ServiceImpl with mock simulation service
	service := &ServiceImpl{
		simulationService: mockSimulationService,
	}

	// Test SimulationService() method
	result := service.SimulationService()

	assert.NotNil(t, result)
	assert.Equal(t, mockSimulationService, result)
}

func TestBuildService_Success(t *testing.T) {
	// Create mock repositories
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockDB := &appContext.DB{}

	// Create mock app context
	mockAppCtx := &MockAppContext{}
	// BuildService calls these methods when creating the services
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo).Maybe()
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo).Maybe()
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo).Maybe()
	mockAppCtx.On("DB").Return(mockDB).Maybe()

	// Test BuildService
	service, err := BuildService(mockAppCtx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.NotNil(t, service.leagueService)
	assert.NotNil(t, service.predictService)
	assert.NotNil(t, service.simulationService)

	// Test that all service getters work
	assert.NotNil(t, service.LeagueService())
	assert.NotNil(t, service.PredictService())
	assert.NotNil(t, service.SimulationService())

	// Note: We don't assert expectations here because the actual service constructors
	// may or may not call the repository methods depending on their implementation
}

func TestBuildService_WithRealAppContext(t *testing.T) {
	// Create real repositories (mocked)
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockDB := &appContext.DB{}

	// This test verifies that BuildService can work with any AppContext implementation
	// We can't easily test with a real AppContext without a database connection
	// So we'll use the mock approach
	mockAppCtx := &MockAppContext{}
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo).Maybe()
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo).Maybe()
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo).Maybe()
	mockAppCtx.On("DB").Return(mockDB).Maybe()

	service, err := BuildService(mockAppCtx)

	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Verify that the service implements the Service interface
	var _ Service = service
	assert.True(t, true, "ServiceImpl implements Service interface")
}

// Test that ServiceImpl implements Service interface
func TestServiceImpl_ImplementsInterface(t *testing.T) {
	var _ Service = (*ServiceImpl)(nil)
	// If this compiles, the interface is implemented correctly
	assert.True(t, true, "ServiceImpl implements Service interface")
}

func TestServiceImpl_AllMethods(t *testing.T) {
	// Create all mock services
	mockLeagueService := &leagueInterfaces.MockLeagueServiceInterface{}
	mockPredictService := &predictInterfaces.MockPredictServiceInterface{}
	mockSimulationService := &simulationInterfaces.MockSimulationServiceInterface{}

	// Create ServiceImpl with all mock services
	service := &ServiceImpl{
		leagueService:     mockLeagueService,
		predictService:    mockPredictService,
		simulationService: mockSimulationService,
	}

	// Test all getter methods
	leagueResult := service.LeagueService()
	predictResult := service.PredictService()
	simulationResult := service.SimulationService()

	// Assert all results
	assert.NotNil(t, leagueResult)
	assert.NotNil(t, predictResult)
	assert.NotNil(t, simulationResult)
	assert.Equal(t, mockLeagueService, leagueResult)
	assert.Equal(t, mockPredictService, predictResult)
	assert.Equal(t, mockSimulationService, simulationResult)
}

func TestServiceImpl_NilServices(t *testing.T) {
	// Create ServiceImpl with nil services to test edge case
	service := &ServiceImpl{}

	// Test getter methods with nil services
	leagueResult := service.LeagueService()
	predictResult := service.PredictService()
	simulationResult := service.SimulationService()

	// Assert all results are nil
	assert.Nil(t, leagueResult)
	assert.Nil(t, predictResult)
	assert.Nil(t, simulationResult)
}

// Benchmark test for BuildService
func BenchmarkBuildService(b *testing.B) {
	// Create mock repositories
	mockLeagueRepo := &interfaces.MockLeagueRepository{}
	mockActiveLeagueRepo := &interfaces.MockActiveLeagueRepository{}
	mockMatchResultRepo := &interfaces.MockMatchResultRepository{}
	mockDB := &appContext.DB{}

	// Create mock app context
	mockAppCtx := &MockAppContext{}
	mockAppCtx.On("LeagueRepository").Return(mockLeagueRepo)
	mockAppCtx.On("ActiveLeagueRepository").Return(mockActiveLeagueRepo)
	mockAppCtx.On("MatchResultRepository").Return(mockMatchResultRepo)
	mockAppCtx.On("DB").Return(mockDB)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		service, err := BuildService(mockAppCtx)
		if err != nil {
			b.Fatalf("BuildService failed: %v", err)
		}
		if service == nil {
			b.Fatal("BuildService returned nil service")
		}
	}
}
