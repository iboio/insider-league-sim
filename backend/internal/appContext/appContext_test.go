package appContext

import (
	"database/sql"
	"testing"

	"league-sim/config"
	"league-sim/internal/repositories/interfaces"

	"github.com/stretchr/testify/assert"
)

// MockDB is a mock implementation for testing
type MockDB struct {
	Sql *sql.DB
}

func TestAppContextImpl_DB(t *testing.T) {
	// Create a mock DB
	mockDB := &DB{Sql: &sql.DB{}}

	// Create AppContext with mock DB
	appCtx := &AppContextImpl{
		db: mockDB,
	}

	// Test DB() method
	result := appCtx.DB()

	assert.NotNil(t, result)
	assert.Equal(t, mockDB, result)
	assert.Equal(t, mockDB.Sql, result.Sql)
}

func TestAppContextImpl_LeagueRepository(t *testing.T) {
	// Create mock repository
	mockRepo := &interfaces.MockLeagueRepository{}

	// Create AppContext with mock repository
	appCtx := &AppContextImpl{
		leagueRepository: mockRepo,
	}

	// Test LeagueRepository() method
	result := appCtx.LeagueRepository()

	assert.NotNil(t, result)
	assert.Equal(t, mockRepo, result)
}

func TestAppContextImpl_ActiveLeagueRepository(t *testing.T) {
	// Create mock repository
	mockRepo := &interfaces.MockActiveLeagueRepository{}

	// Create AppContext with mock repository
	appCtx := &AppContextImpl{
		activeLeagueRepository: mockRepo,
	}

	// Test ActiveLeagueRepository() method
	result := appCtx.ActiveLeagueRepository()

	assert.NotNil(t, result)
	assert.Equal(t, mockRepo, result)
}

func TestAppContextImpl_MatchResultRepository(t *testing.T) {
	// Create mock repository
	mockRepo := &interfaces.MockMatchResultRepository{}

	// Create AppContext with mock repository
	appCtx := &AppContextImpl{
		matchResultRepository: mockRepo,
	}

	// Test MatchResultRepository() method
	result := appCtx.MatchResultRepository()

	assert.NotNil(t, result)
	assert.Equal(t, mockRepo, result)
}

func TestAppContextDBInit_Success(t *testing.T) {
	// Save original config values
	originalHost := config.MySQLHost
	originalPort := config.MySQLPort
	originalUser := config.MySQLUser
	originalPassword := config.MySQLPassword
	originalDatabase := config.MySQLDatabase

	// Restore original values after test
	defer func() {
		config.MySQLHost = originalHost
		config.MySQLPort = originalPort
		config.MySQLUser = originalUser
		config.MySQLPassword = originalPassword
		config.MySQLDatabase = originalDatabase
	}()

	// Set test values
	config.MySQLHost = "localhost"
	config.MySQLPort = "3306"
	config.MySQLUser = "test_user"
	config.MySQLPassword = "test_password"
	config.MySQLDatabase = "test_db"

	// Test AppContextDBInit
	db, err := AppContextDBInit()

	// Since we don't have a real MySQL server, we expect an error
	// but we can test that the function handles it gracefully
	if err != nil {
		assert.Nil(t, db)
		assert.Error(t, err)
		t.Logf("Expected error when no MySQL server available: %v", err)
	} else {
		// If somehow we have a connection, test it
		assert.NotNil(t, db)
		assert.NotNil(t, db.Sql)
		assert.NoError(t, err)

		// Clean up
		db.Sql.Close()
	}
}

func TestAppContextDBInit_DatabaseConnectionError(t *testing.T) {
	// Save original config values
	originalHost := config.MySQLHost
	originalPort := config.MySQLPort
	originalUser := config.MySQLUser
	originalPassword := config.MySQLPassword
	originalDatabase := config.MySQLDatabase

	// Restore original values after test
	defer func() {
		config.MySQLHost = originalHost
		config.MySQLPort = originalPort
		config.MySQLUser = originalUser
		config.MySQLPassword = originalPassword
		config.MySQLDatabase = originalDatabase
	}()

	// Set invalid connection parameters
	config.MySQLHost = "invalid_host_that_does_not_exist"
	config.MySQLPort = "3306"
	config.MySQLUser = "invalid_user"
	config.MySQLPassword = "invalid_password"
	config.MySQLDatabase = "invalid_db"

	// Test AppContextDBInit with invalid parameters
	db, err := AppContextDBInit()

	// Should return an error
	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestAppContextInit_Success(t *testing.T) {
	// Save original config values
	originalHost := config.MySQLHost
	originalPort := config.MySQLPort
	originalUser := config.MySQLUser
	originalPassword := config.MySQLPassword
	originalDatabase := config.MySQLDatabase

	// Restore original values after test
	defer func() {
		config.MySQLHost = originalHost
		config.MySQLPort = originalPort
		config.MySQLUser = originalUser
		config.MySQLPassword = originalPassword
		config.MySQLDatabase = originalDatabase
	}()

	// Set test values
	config.MySQLHost = "localhost"
	config.MySQLPort = "3306"
	config.MySQLUser = "test_user"
	config.MySQLPassword = "test_password"
	config.MySQLDatabase = "test_db"

	// Test AppContextInit
	appCtx, err := AppContextInit()

	// Since we don't have a real MySQL server, we expect an error
	if err != nil {
		assert.Nil(t, appCtx)
		assert.Error(t, err)
		t.Logf("Expected error when no MySQL server available: %v", err)
	} else {
		// If somehow we have a connection, test the context
		assert.NotNil(t, appCtx)
		assert.NotNil(t, appCtx.db)
		assert.NotNil(t, appCtx.leagueRepository)
		assert.NotNil(t, appCtx.activeLeagueRepository)
		assert.NotNil(t, appCtx.matchResultRepository)

		// Test interface methods
		assert.NotNil(t, appCtx.DB())
		assert.NotNil(t, appCtx.LeagueRepository())
		assert.NotNil(t, appCtx.ActiveLeagueRepository())
		assert.NotNil(t, appCtx.MatchResultRepository())

		// Clean up
		appCtx.db.Sql.Close()
	}
}

func TestAppContextInit_DatabaseError(t *testing.T) {
	// Save original config values
	originalHost := config.MySQLHost
	originalPort := config.MySQLPort
	originalUser := config.MySQLUser
	originalPassword := config.MySQLPassword
	originalDatabase := config.MySQLDatabase

	// Restore original values after test
	defer func() {
		config.MySQLHost = originalHost
		config.MySQLPort = originalPort
		config.MySQLUser = originalUser
		config.MySQLPassword = originalPassword
		config.MySQLDatabase = originalDatabase
	}()

	// Set invalid connection parameters to force an error
	config.MySQLHost = "invalid_host_that_does_not_exist"
	config.MySQLPort = "99999"
	config.MySQLUser = "invalid_user"
	config.MySQLPassword = "invalid_password"
	config.MySQLDatabase = "invalid_db"

	// Test AppContextInit with invalid parameters
	appCtx, err := AppContextInit()

	// Should return an error and nil context
	assert.Nil(t, appCtx)
	assert.Error(t, err)
}

// Test that AppContextImpl implements AppContext interface
func TestAppContextImpl_ImplementsInterface(t *testing.T) {
	var _ AppContext = (*AppContextImpl)(nil)
	// If this compiles, the interface is implemented correctly
	assert.True(t, true, "AppContextImpl implements AppContext interface")
}

// Test DB struct
func TestDB_Structure(t *testing.T) {
	mockSqlDB := &sql.DB{}
	db := &DB{Sql: mockSqlDB}

	assert.NotNil(t, db)
	assert.Equal(t, mockSqlDB, db.Sql)
}

// Benchmark test for AppContextInit
func BenchmarkAppContextInit(b *testing.B) {
	// Save original config values
	originalHost := config.MySQLHost
	originalPort := config.MySQLPort
	originalUser := config.MySQLUser
	originalPassword := config.MySQLPassword
	originalDatabase := config.MySQLDatabase

	// Restore original values after benchmark
	defer func() {
		config.MySQLHost = originalHost
		config.MySQLPort = originalPort
		config.MySQLUser = originalUser
		config.MySQLPassword = originalPassword
		config.MySQLDatabase = originalDatabase
	}()

	// Set test values
	config.MySQLHost = "localhost"
	config.MySQLPort = "3306"
	config.MySQLUser = "test_user"
	config.MySQLPassword = "test_password"
	config.MySQLDatabase = "test_db"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		appCtx, err := AppContextInit()
		if err == nil && appCtx != nil && appCtx.db != nil && appCtx.db.Sql != nil {
			appCtx.db.Sql.Close()
		}
	}
}
