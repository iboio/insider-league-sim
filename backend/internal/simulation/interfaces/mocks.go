package interfaces

import (
	"league-sim/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockSimulationServiceInterface is a mock implementation of SimulationServiceInterface
type MockSimulationServiceInterface struct {
	mock.Mock
}

func (m *MockSimulationServiceInterface) Simulation(leagueId string, playAllFixture bool) (models.SimulationResponse, error) {
	args := m.Called(leagueId, playAllFixture)
	return args.Get(0).(models.SimulationResponse), args.Error(1)
}

func (m *MockSimulationServiceInterface) EditMatch(data models.EditMatchResult) error {
	args := m.Called(data)
	return args.Error(0)
}
