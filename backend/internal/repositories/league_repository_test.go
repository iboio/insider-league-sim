package repositories

import (
	"errors"
	"testing"

	"league-sim/internal/models"
	"league-sim/internal/repositories/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewLeagueRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)
	assert.NotNil(t, repo)
	assert.Implements(t, (*interfaces.LeagueRepository)(nil), repo)
}

func TestLeagueRepository_SetLeague_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	request := models.CreateLeagueRequest{
		LeagueName: "Test League",
		TeamCount:  "8",
	}

	// Mock expectations
	mock.ExpectExec("INSERT INTO league \\(leagueId, name\\) VALUES \\(\\?,\\?\\)").
		WithArgs(leagueId, request.LeagueName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	err = repo.SetLeague(leagueId, request)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeagueRepository_SetLeague_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"
	request := models.CreateLeagueRequest{
		LeagueName: "Test League",
		TeamCount:  "8",
	}

	// Mock expectations
	expectedError := errors.New("database connection failed")
	mock.ExpectExec("INSERT INTO league \\(leagueId, name\\) VALUES \\(\\?,\\?\\)").
		WithArgs(leagueId, request.LeagueName).
		WillReturnError(expectedError)

	// Execute
	err = repo.SetLeague(leagueId, request)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeagueRepository_GetLeague_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Test data
	expectedLeagues := []models.GetLeaguesIdsWithNameResponse{
		{LeagueId: "league1", LeagueName: "League 1"},
		{LeagueId: "league2", LeagueName: "League 2"},
	}

	// Mock expectations
	rows := sqlmock.NewRows([]string{"leagueId", "name"}).
		AddRow("league1", "League 1").
		AddRow("league2", "League 2")

	mock.ExpectQuery("SELECT leagueId, name FROM league").
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetLeague()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedLeagues, result)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeagueRepository_GetLeague_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Mock expectations
	rows := sqlmock.NewRows([]string{"leagueId", "name"})
	mock.ExpectQuery("SELECT leagueId, name FROM league").
		WillReturnRows(rows)

	// Execute
	result, err := repo.GetLeague()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeagueRepository_GetLeague_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Mock expectations
	expectedError := errors.New("query failed")
	mock.ExpectQuery("SELECT leagueId, name FROM league").
		WillReturnError(expectedError)

	// Execute and expect panic
	assert.Panics(t, func() {
		repo.GetLeague()
	}, "Should panic when query fails")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeagueRepository_GetLeague_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Mock expectations - invalid data type for scan
	rows := sqlmock.NewRows([]string{"leagueId", "name"}).
		AddRow(123, nil) // Invalid types that will cause scan error

	mock.ExpectQuery("SELECT leagueId, name FROM league").
		WillReturnRows(rows)

	// Execute and expect panic
	assert.Panics(t, func() {
		repo.GetLeague()
	}, "Should panic when scan fails")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeagueRepository_DeleteLeague_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations
	mock.ExpectExec("DELETE FROM league WHERE leagueId = \\?").
		WithArgs(leagueId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute
	err = repo.DeleteLeague(leagueId)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeagueRepository_DeleteLeague_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Test data
	leagueId := "test-league-id"

	// Mock expectations
	expectedError := errors.New("delete failed")
	mock.ExpectExec("DELETE FROM league WHERE leagueId = \\?").
		WithArgs(leagueId).
		WillReturnError(expectedError)

	// Execute and expect panic
	assert.Panics(t, func() {
		repo.DeleteLeague(leagueId)
	}, "Should panic when delete fails")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLeagueRepository_DeleteLeague_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLeagueRepository(db)

	// Test data
	leagueId := "non-existent-league"

	// Mock expectations
	mock.ExpectExec("DELETE FROM league WHERE leagueId = \\?").
		WithArgs(leagueId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Execute
	err = repo.DeleteLeague(leagueId)

	// Assert
	assert.NoError(t, err) // Should not error even if no rows affected
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Benchmark tests
func BenchmarkLeagueRepository_SetLeague(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewLeagueRepository(db)
	request := models.CreateLeagueRequest{
		LeagueName: "Benchmark League",
		TeamCount:  "8",
	}

	for i := 0; i < b.N; i++ {
		mock.ExpectExec("INSERT INTO league \\(leagueId, name\\) VALUES \\(\\?,\\?\\)").
			WithArgs(sqlmock.AnyArg(), request.LeagueName).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo.SetLeague("benchmark-league", request)
	}
}

func BenchmarkLeagueRepository_GetLeague(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewLeagueRepository(db)

	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{"leagueId", "name"}).
			AddRow("league1", "League 1").
			AddRow("league2", "League 2")

		mock.ExpectQuery("SELECT leagueId, name FROM league").
			WillReturnRows(rows)

		repo.GetLeague()
	}
}

func BenchmarkLeagueRepository_DeleteLeague(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := NewLeagueRepository(db)

	for i := 0; i < b.N; i++ {
		mock.ExpectExec("DELETE FROM league WHERE leagueId = \\?").
			WithArgs(sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo.DeleteLeague("benchmark-league")
	}
}
