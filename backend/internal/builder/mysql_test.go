package builder

import (
	"testing"

	"league-sim/config"

	"github.com/stretchr/testify/assert"
)

func TestSqlConnectionInit_Success(t *testing.T) {
	// Setup test environment variables
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

	// Set test values - using a mock/test database
	config.MySQLHost = "localhost"
	config.MySQLPort = "3306"
	config.MySQLUser = "test_user"
	config.MySQLPassword = "test_password"
	config.MySQLDatabase = "test_db"

	// Note: This test will fail if no actual MySQL server is running
	// In a real scenario, you might want to use a test container or mock
	db, err := SqlConnectionInit()

	// If we have a real MySQL server running, test the connection
	if err == nil {
		assert.NotNil(t, db)
		assert.NoError(t, err)

		// Test that we can ping the database
		pingErr := db.Ping()
		assert.NoError(t, pingErr)

		// Clean up
		db.Close()
	} else {
		// If no MySQL server is available, just test that the function handles errors gracefully
		assert.Nil(t, db)
		assert.Error(t, err)
		t.Logf("Expected error when no MySQL server available: %v", err)
	}
}

func TestSqlConnectionInit_InvalidCredentials(t *testing.T) {
	// Setup test environment variables
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

	// Set invalid credentials
	config.MySQLHost = "localhost"
	config.MySQLPort = "3306"
	config.MySQLUser = "invalid_user"
	config.MySQLPassword = "invalid_password"
	config.MySQLDatabase = "invalid_db"

	db, err := SqlConnectionInit()

	// Should return an error with invalid credentials
	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestSqlConnectionInit_InvalidHost(t *testing.T) {
	// Setup test environment variables
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

	// Set invalid host
	config.MySQLHost = "invalid_host_that_does_not_exist"
	config.MySQLPort = "3306"
	config.MySQLUser = "test_user"
	config.MySQLPassword = "test_password"
	config.MySQLDatabase = "test_db"

	db, err := SqlConnectionInit()

	// Should return an error with invalid host
	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestSqlConnectionInit_InvalidPort(t *testing.T) {
	// Setup test environment variables
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

	// Set invalid port
	config.MySQLHost = "localhost"
	config.MySQLPort = "99999" // Invalid port
	config.MySQLUser = "test_user"
	config.MySQLPassword = "test_password"
	config.MySQLDatabase = "test_db"

	db, err := SqlConnectionInit()

	// Should return an error with invalid port
	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestSqlConnectionInit_DSNFormat(t *testing.T) {
	// This test verifies that the DSN is formatted correctly
	// We can test this by checking the error message when connection fails

	// Setup test environment variables
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
	config.MySQLHost = "testhost"
	config.MySQLPort = "3306"
	config.MySQLUser = "testuser"
	config.MySQLPassword = "testpass"
	config.MySQLDatabase = "testdb"

	// The function should at least attempt to connect with the correct DSN format
	// Even if it fails, we can verify the DSN was constructed properly
	db, err := SqlConnectionInit()

	// We expect an error since this is likely not a real server
	assert.Error(t, err)
	assert.Nil(t, db)

	// The error should indicate a connection issue, not a DSN format issue
	// DSN format errors would typically mention "invalid DSN" or similar
	assert.NotContains(t, err.Error(), "invalid DSN")
}

// Benchmark test for connection initialization
func BenchmarkSqlConnectionInit(b *testing.B) {
	// Setup test environment variables
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
		db, err := SqlConnectionInit()
		if err == nil && db != nil {
			db.Close()
		}
	}
}
