package interfaces

import (
	"league-sim/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockLeagueRepository is a mock implementation of LeagueRepository
type MockLeagueRepository struct {
	mock.Mock
}

func (m *MockLeagueRepository) SetLeague(id string, data models.CreateLeagueRequest) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockLeagueRepository) GetLeague() ([]models.GetLeaguesIdsWithNameResponse, error) {
	args := m.Called()
	return args.Get(0).([]models.GetLeaguesIdsWithNameResponse), args.Error(1)
}

func (m *MockLeagueRepository) DeleteLeague(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockActiveLeagueRepository is a mock implementation of ActiveLeagueRepository
type MockActiveLeagueRepository struct {
	mock.Mock
}

func (m *MockActiveLeagueRepository) GetActiveLeague(id string) (models.League, error) {
	args := m.Called(id)
	return args.Get(0).(models.League), args.Error(1)
}

func (m *MockActiveLeagueRepository) SetActiveLeague(data models.League) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockActiveLeagueRepository) GetActiveLeagueTeams(id string) ([]models.Team, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Team), args.Error(1)
}

func (m *MockActiveLeagueRepository) GetActiveLeaguesFixtures(id string) (
	models.GetFixturesResponse, error) {
	args := m.Called(id)
	return args.Get(0).(models.GetFixturesResponse), args.Error(1)
}

func (m *MockActiveLeagueRepository) GetActiveLeaguesStandings(id string) ([]models.Standings, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Standings), args.Error(1)
}

// MockMatchResultRepository is a mock implementation of MatchesRepository
type MockMatchResultRepository struct {
	mock.Mock
}

func (m *MockMatchResultRepository) EditMatchScore(data models.EditMatchResult) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockMatchResultRepository) SetMatchResults(leagueId string, matchResults []models.Matches) error {
	args := m.Called(leagueId, matchResults)
	return args.Error(0)
}

func (m *MockMatchResultRepository) GetMatchResults(leagueId string) ([]models.Matches, error) {
	args := m.Called(leagueId)
	return args.Get(0).([]models.Matches), args.Error(1)
}

func (m *MockMatchResultRepository) DeleteMatchResults(leagueId string) error {
	args := m.Called(leagueId)
	return args.Error(0)
}

func (m *MockMatchResultRepository) GetMatchResultByWeekAndTeam(data models.EditMatchResult) (models.Matches, error) {
	args := m.Called(data)
	return args.Get(0).(models.Matches), args.Error(1)
}
