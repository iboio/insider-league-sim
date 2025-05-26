package interfaces

import (
	"league-sim/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockLeagueServiceInterface is a mock implementation of LeagueServiceInterface
type MockLeagueServiceInterface struct {
	mock.Mock
}

func (m *MockLeagueServiceInterface) CreateLeague(n string, leagueName string) (models.GetLeaguesIdsWithNameResponse, error) {
	args := m.Called(n, leagueName)
	return args.Get(0).(models.GetLeaguesIdsWithNameResponse), args.Error(1)
}

func (m *MockLeagueServiceInterface) ResetLeague(leagueId string) error {
	args := m.Called(leagueId)
	return args.Error(0)
}
