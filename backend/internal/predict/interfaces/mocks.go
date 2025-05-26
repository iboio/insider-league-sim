package interfaces

import (
	"league-sim/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockPredictServiceInterface is a mock implementation of PredictServiceInterface
type MockPredictServiceInterface struct {
	mock.Mock
}

func (m *MockPredictServiceInterface) PredictChampionShipSession(id string) ([]models.PredictedStanding, error) {
	args := m.Called(id)
	return args.Get(0).([]models.PredictedStanding), args.Error(1)
}
