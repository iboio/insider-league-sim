package repositories

import (
	"errors"
	"testing"

	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewMatchResultRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)
	assert.NotNil(t, repo)
	assert.Implements(t, (*interfaces.MatchResultRepository)(nil), repo)
}

func TestMatchResultRepository_GetMatchResults_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"
	expectedResults := []models.MatchResult{
		{
			Home:      "Team A",
			HomeScore: 2,
			Away:      "Team B",
			AwayScore: 1,
			Winner:    "Team A",
			MatchWeek: 1,
		},
		{
			Home:      "Team C",
			HomeScore: 0,
			Away:      "Team D",
			AwayScore: 3,
			Winner:    "Team D",
			MatchWeek: 1,
		},
	}

	// Mock expectations
	rows := sqlmock.NewRows([]string{"homeTeam", "homeGoals", "awayTeam", "awayGoals", "winnerName", "matchWeek"}).
		AddRow("Team A", 2, "Team B", 1, "Team A", 1).
		AddRow("Team C", 0, "Team D", 3, "Team D", 1)

	mock.ExpectQuery("SELECT homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek FROM match_results WHERE leagueId = \\? ORDER BY matchWeek").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetMatchResults(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResults, result)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_GetMatchResults_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "empty-league-id"

	// Mock expectations
	rows := sqlmock.NewRows([]string{"homeTeam", "homeGoals", "awayTeam", "awayGoals", "winnerName", "matchWeek"})

	mock.ExpectQuery("SELECT homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek FROM match_results WHERE leagueId = \\? ORDER BY matchWeek").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetMatchResults(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_GetMatchResults_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations
	expectedError := errors.New("database query failed")
	mock.ExpectQuery("SELECT homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek FROM match_results WHERE leagueId = \\? ORDER BY matchWeek").
		WithArgs(leagueId).
		WillReturnError(expectedError)

	// Execute
	result, err := repo.GetMatchResults(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_GetMatchResults_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations - invalid data types that will cause scan error
	rows := sqlmock.NewRows([]string{"homeTeam", "homeGoals", "awayTeam", "awayGoals", "winnerName", "matchWeek"}).
		AddRow("Team A", "invalid_score", "Team B", 1, "Team A", 1)

	mock.ExpectQuery("SELECT homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek FROM match_results WHERE leagueId = \\? ORDER BY matchWeek").
		WithArgs(leagueId).
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetMatchResults(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_SetMatchResults_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"
	matchResults := []models.MatchResult{
		{
			Home:      "Team A",
			HomeScore: 2,
			Away:      "Team B",
			AwayScore: 1,
			Winner:    "Team A",
			MatchWeek: 1,
		},
		{
			Home:      "Team C",
			HomeScore: 0,
			Away:      "Team D",
			AwayScore: 3,
			Winner:    "Team D",
			MatchWeek: 1,
		},
	}

	// Mock expectations
	mock.ExpectExec("INSERT INTO match_results \\(leagueId, homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\),\\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(
			leagueId, "Team A", 2, "Team B", 1, "Team A", 1,
			leagueId, "Team C", 0, "Team D", 3, "Team D", 1,
		).
		WillReturnResult(sqlmock.NewResult(1, 2))

	// Execute
	err = repo.SetMatchResults(leagueId, matchResults)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_SetMatchResults_EmptyResults(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"
	matchResults := []models.MatchResult{}

	// Execute
	err = repo.SetMatchResults(leagueId, matchResults)

	// Assert
	assert.NoError(t, err) // Should return nil for empty results
}

func TestMatchResultRepository_SetMatchResults_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"
	matchResults := []models.MatchResult{
		{
			Home:      "Team A",
			HomeScore: 2,
			Away:      "Team B",
			AwayScore: 1,
			Winner:    "Team A",
			MatchWeek: 1,
		},
	}

	// Mock expectations
	expectedError := errors.New("database insert failed")
	mock.ExpectExec("INSERT INTO match_results \\(leagueId, homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(leagueId, "Team A", 2, "Team B", 1, "Team A", 1).
		WillReturnError(expectedError)

	// Execute
	err = repo.SetMatchResults(leagueId, matchResults)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_EditMatchScore_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		Winner:    "Team A",
		MatchWeek: 1,
	}

	// Mock expectations
	mock.ExpectExec("UPDATE match_results SET homeGoals = \\?, awayGoals = \\?, winnerName = \\? WHERE leagueId = \\? AND matchWeek = \\? AND homeTeam = \\? AND awayTeam = \\?").
		WithArgs(3, 1, "Team A", "test-league-id", 1, "Team A", "Team B").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute
	err = repo.EditMatchScore(editData)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_EditMatchScore_AwayTeam(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 1,
		AwayScore: 2,
		Winner:    "Team B",
		MatchWeek: 1,
	}

	// Mock expectations
	mock.ExpectExec("UPDATE match_results SET homeGoals = \\?, awayGoals = \\?, winnerName = \\? WHERE leagueId = \\? AND matchWeek = \\? AND homeTeam = \\? AND awayTeam = \\?").
		WithArgs(1, 2, "Team B", "test-league-id", 1, "Team A", "Team B").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute
	err = repo.EditMatchScore(editData)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_EditMatchScore_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	editData := models.EditMatchResult{
		LeagueId:  "test-league-id",
		Home:      "Team A",
		Away:      "Team B",
		HomeScore: 3,
		AwayScore: 1,
		Winner:    "Team A",
		MatchWeek: 1,
	}

	// Mock expectations
	expectedError := errors.New("update failed")
	mock.ExpectExec("UPDATE match_results SET homeGoals = \\?, awayGoals = \\?, winnerName = \\? WHERE leagueId = \\? AND matchWeek = \\? AND homeTeam = \\? AND awayTeam = \\?").
		WithArgs(3, 1, "Team A", "test-league-id", 1, "Team A", "Team B").
		WillReturnError(expectedError)

	// Execute
	err = repo.EditMatchScore(editData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_DeleteMatchResults_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations
	mock.ExpectExec("DELETE FROM match_results WHERE leagueId = \\?").
		WithArgs(leagueId).
		WillReturnResult(sqlmock.NewResult(0, 5)) // 5 rows deleted

	// Execute
	err = repo.DeleteMatchResults(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_DeleteMatchResults_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "non-existent-league"

	// Mock expectations
	mock.ExpectExec("DELETE FROM match_results WHERE leagueId = \\?").
		WithArgs(leagueId).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows deleted

	// Execute
	err = repo.DeleteMatchResults(leagueId)

	// Assert
	assert.NoError(t, err) // Should not error even if no rows affected
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_DeleteMatchResults_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations
	expectedError := errors.New("delete failed")
	mock.ExpectExec("DELETE FROM match_results WHERE leagueId = \\?").
		WithArgs(leagueId).
		WillReturnError(expectedError)

	// Execute
	err = repo.DeleteMatchResults(leagueId)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMatchResultRepository_SetMatchResults_SingleResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMatchResultRepository(db)

	// Test data
	leagueId := "test-league-id"
	matchResults := []models.MatchResult{
		{
			Home:      "Team A",
			HomeScore: 1,
			Away:      "Team B",
			AwayScore: 1,
			Winner:    "Draw",
			MatchWeek: 2,
		},
	}

	// Mock expectations
	mock.ExpectExec("INSERT INTO match_results \\(leagueId, homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
		WithArgs(leagueId, "Team A", 1, "Team B", 1, "Draw", 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.SetMatchResults(leagueId, matchResults)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Benchmark tests
func BenchmarkMatchResultRepository_GetMatchResults(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewMatchResultRepository(db)

	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{"homeTeam", "homeGoals", "awayTeam", "awayGoals", "winnerName", "matchWeek"}).
			AddRow("Team A", 2, "Team B", 1, "Team A", 1)

		mock.ExpectQuery("SELECT homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek FROM match_results WHERE leagueId = \\? ORDER BY matchWeek").
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(rows)

		repo.GetMatchResults("benchmark-league")
	}
}

func BenchmarkMatchResultRepository_SetMatchResults(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewMatchResultRepository(db)
	matchResults := []models.MatchResult{
		{
			Home:      "Team A",
			HomeScore: 2,
			Away:      "Team B",
			AwayScore: 1,
			Winner:    "Team A",
			MatchWeek: 1,
		},
	}

	for i := 0; i < b.N; i++ {
		mock.ExpectExec("INSERT INTO match_results \\(leagueId, homeTeam, homeGoals, awayTeam, awayGoals, winnerName, matchWeek\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\)").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo.SetMatchResults("benchmark-league", matchResults)
	}
}

func BenchmarkMatchResultRepository_DeleteMatchResults(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewMatchResultRepository(db)

	for i := 0; i < b.N; i++ {
		mock.ExpectExec("DELETE FROM match_results WHERE leagueId = \\?").
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo.DeleteMatchResults("benchmark-league")
	}
}
