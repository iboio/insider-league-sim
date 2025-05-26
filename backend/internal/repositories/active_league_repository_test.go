package repositories

import (
	"database/sql"
	"errors"
	"testing"

	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewActiveLeagueRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)
	assert.NotNil(t, repo)
	assert.Implements(t, (*interfaces.ActiveLeagueRepository)(nil), repo)
}

func TestActiveLeagueRepository_GetActiveLeague_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	upcomingFixturesJson := `[{"Number":1,"Matches":[]}]`
	playedFixturesJson := `[]`
	standingsJson := `[{"Team":{"Name":"Team A"},"Points":3,"Wins":1}]`
	teamsJson := `[{"Name":"Team A","AttackPower":80,"DefensePower":75,"Stamina":90,"Morale":85}]`
	currentWeek := 1

	// Mock expectations
	rows := sqlmock.NewRows([]string{"upcomingFixtures", "playedFixtures", "currentWeek", "teams", "standings"}).
		AddRow(upcomingFixturesJson, playedFixturesJson, currentWeek, teamsJson, standingsJson)

	mock.ExpectQuery("SELECT upcomingFixtures, playedFixtures, currentWeek,teams, standings FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetActiveLeague(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, leagueId, result.LeagueID)
	assert.Equal(t, currentWeek, result.CurrentWeek)
	assert.Len(t, result.UpcomingFixtures, 1)
	assert.Len(t, result.PlayedFixtures, 0)
	assert.Len(t, result.Standings, 1)
	assert.Len(t, result.Teams, 1)
	assert.Equal(t, "Team A", result.Teams[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeague_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "non-existent-league"

	// Mock expectations
	mock.ExpectQuery("SELECT upcomingFixtures, playedFixtures, currentWeek,teams, standings FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnError(sql.ErrNoRows)

	// Execute
	result, err := repo.GetActiveLeague(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Equal(t, models.League{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeague_InvalidJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	invalidJson := `invalid json`

	// Mock expectations
	rows := sqlmock.NewRows([]string{"upcomingFixtures", "playedFixtures", "currentWeek", "teams", "standings"}).
		AddRow(invalidJson, `[]`, 1, `[]`, `[]`)

	mock.ExpectQuery("SELECT upcomingFixtures, playedFixtures, currentWeek,teams, standings FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetActiveLeague(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, models.League{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeagueTeams_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	teamsJson := `[{"Name":"Team A","AttackPower":80,"DefensePower":75,"Stamina":90,"Morale":85},{"Name":"Team B","AttackPower":75,"DefensePower":80,"Stamina":85,"Morale":90}]`

	// Mock expectations
	rows := sqlmock.NewRows([]string{"teams"}).
		AddRow(teamsJson)

	mock.ExpectQuery("SELECT teams FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetActiveLeagueTeams(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Team A", result[0].Name)
	assert.Equal(t, "Team B", result[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeagueTeams_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations
	mock.ExpectQuery("SELECT teams FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnError(sql.ErrNoRows)

	// Execute and expect panic
	assert.Panics(t, func() {
		repo.GetActiveLeagueTeams(leagueId)
	}, "Should panic when scan fails")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeagueTeams_InvalidJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	invalidJson := `invalid json`

	// Mock expectations
	rows := sqlmock.NewRows([]string{"teams"}).
		AddRow(invalidJson)

	mock.ExpectQuery("SELECT teams FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute and expect panic
	assert.Panics(t, func() {
		repo.GetActiveLeagueTeams(leagueId)
	}, "Should panic when JSON parsing fails")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_SetActiveLeague_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	league := models.League{
		LeagueID:    "test-league-id",
		CurrentWeek: 1,
		UpcomingFixtures: []models.Week{
			{Number: 1, Matches: []models.Match{}},
		},
		PlayedFixtures: []models.Week{},
		Standings: []models.Standings{
			{Team: models.Team{Name: "Team A"}, Points: 3, Wins: 1},
		},
		Teams: []models.Team{
			{Name: "Team A", AttackPower: 80, DefensePower: 75, Stamina: 90, Morale: 85},
		},
	}

	// Mock expectations
	mock.ExpectExec("INSERT INTO active_league \\(leagueId, upcomingFixtures ,playedFixtures,teams, currentWeek, standings\\) VALUES \\(\\?, \\?, \\?, \\?, \\?,\\?\\)").
		WithArgs(
			league.LeagueID,
			sqlmock.AnyArg(), // upcomingFixtures JSON
			sqlmock.AnyArg(), // playedFixtures JSON
			sqlmock.AnyArg(), // teams JSON
			league.CurrentWeek,
			sqlmock.AnyArg(), // standings JSON
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.SetActiveLeague(league)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_SetActiveLeague_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	league := models.League{
		LeagueID:    "test-league-id",
		CurrentWeek: 1,
	}

	// Mock expectations
	expectedError := errors.New("database insert failed")
	mock.ExpectExec("INSERT INTO active_league \\(leagueId, upcomingFixtures ,playedFixtures,teams, currentWeek, standings\\) VALUES \\(\\?, \\?, \\?, \\?, \\?,\\?\\)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(expectedError)

	// Execute and expect panic
	assert.Panics(t, func() {
		repo.SetActiveLeague(league)
	}, "Should panic when database insert fails")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeaguesFixtures_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	upcomingFixturesJson := `[{"Number":2,"Matches":[]}]`
	playedFixturesJson := `[{"Number":1,"Matches":[]}]`

	// Mock expectations
	rows := sqlmock.NewRows([]string{"upcomingFixtures", "playedFixtures"}).
		AddRow(upcomingFixturesJson, playedFixturesJson)

	mock.ExpectQuery("SELECT upcomingFixtures, playedFixtures FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetActiveLeaguesFixtures(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.UpcomingFixtures, 1)
	assert.Len(t, result.PlayedFixtures, 1)
	assert.Equal(t, 2, result.UpcomingFixtures[0].Number)
	assert.Equal(t, 1, result.PlayedFixtures[0].Number)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeaguesFixtures_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations
	expectedError := errors.New("database query failed")
	mock.ExpectQuery("SELECT upcomingFixtures, playedFixtures FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnError(expectedError)

	// Execute
	result, err := repo.GetActiveLeaguesFixtures(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, models.GetActiveLeagueFixturesResponse{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeaguesStandings_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	standingsJson := `[{"Team":{"Name":"Team A"},"Points":6,"Wins":2,"Draws":0,"Losses":0},{"Team":{"Name":"Team B"},"Points":3,"Wins":1,"Draws":0,"Losses":1}]`

	// Mock expectations
	rows := sqlmock.NewRows([]string{"standings"}).
		AddRow(standingsJson)

	mock.ExpectQuery("SELECT standings FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetActiveLeaguesStandings(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Team A", result[0].Team.Name)
	assert.Equal(t, 6, result[0].Points)
	assert.Equal(t, "Team B", result[1].Team.Name)
	assert.Equal(t, 3, result[1].Points)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeaguesStandings_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations
	expectedError := errors.New("database query failed")
	mock.ExpectQuery("SELECT standings FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnError(expectedError)

	// Execute
	result, err := repo.GetActiveLeaguesStandings(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestActiveLeagueRepository_GetActiveLeaguesStandings_InvalidJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	invalidJson := `invalid json`

	// Mock expectations
	rows := sqlmock.NewRows([]string{"standings"}).
		AddRow(invalidJson)

	mock.ExpectQuery("SELECT standings FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetActiveLeaguesStandings(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Benchmark tests
func BenchmarkActiveLeagueRepository_GetActiveLeague(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{"upcomingFixtures", "playedFixtures", "currentWeek", "teams", "standings"}).
			AddRow(`[]`, `[]`, 1, `[]`, `[]`)

		mock.ExpectQuery("SELECT upcomingFixtures, playedFixtures, currentWeek,teams, standings FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(rows)

		repo.GetActiveLeague("benchmark-league")
	}
}

func BenchmarkActiveLeagueRepository_SetActiveLeague(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewActiveLeagueRepository(db)
	league := models.League{
		LeagueID:    "benchmark-league",
		CurrentWeek: 1,
	}

	for i := 0; i < b.N; i++ {
		mock.ExpectExec("INSERT INTO active_league \\(leagueId, upcomingFixtures ,playedFixtures,teams, currentWeek, standings\\) VALUES \\(\\?, \\?, \\?, \\?, \\?,\\?\\)").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo.SetActiveLeague(league)
	}
}

func BenchmarkActiveLeagueRepository_GetActiveLeaguesStandings(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewActiveLeagueRepository(db)

	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{"standings"}).
			AddRow(`[]`)

		mock.ExpectQuery("SELECT standings FROM active_league WHERE leagueId = \\? order by createdAt desc limit 1").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(rows)

		repo.GetActiveLeaguesStandings("benchmark-league")
	}
}
