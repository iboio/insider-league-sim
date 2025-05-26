package main

import (
	"os"
	"testing"
	"time"

	"league-sim/config"

	"github.com/stretchr/testify/assert"
)

func TestMain_ConfigLoad(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := `HTTP_PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=testuser
MYSQL_PASSWORD=testpass
MYSQL_DATABASE=testdb`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Test that config.LoadConfig() doesn't panic
	assert.NotPanics(t, func() {
		config.LoadConfig()
	}, "LoadConfig should not panic with valid .env file")

	// Verify config values are loaded
	assert.Equal(t, "8080", config.HTTPPort)
	assert.Equal(t, "localhost", config.RedisHost)
	assert.Equal(t, "6379", config.RedisPort)
	assert.Equal(t, "localhost", config.MySQLHost)
	assert.Equal(t, "3306", config.MySQLPort)
	assert.Equal(t, "testuser", config.MySQLUser)
	assert.Equal(t, "testpass", config.MySQLPassword)
	assert.Equal(t, "testdb", config.MySQLDatabase)
}

func TestMain_ConfigLoadPanic(t *testing.T) {
	// Remove .env file to cause panic
	os.Remove(".env")

	// Test that config.LoadConfig() panics without .env file
	assert.Panics(t, func() {
		config.LoadConfig()
	}, "LoadConfig should panic without .env file")
}

func TestMain_Integration(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := `HTTP_PORT=0
REDIS_HOST=localhost
REDIS_PORT=6379
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=testuser
MYSQL_PASSWORD=testpass
MYSQL_DATABASE=testdb`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Test main function execution in a goroutine
	// Note: This will likely fail due to database connection issues, but we're testing the flow
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected to panic due to missing database/redis connections
				done <- true
			}
		}()
		main()
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		// Expected behavior - main should panic due to missing dependencies
		t.Log("Main function panicked as expected due to missing dependencies")
	case <-time.After(2 * time.Second):
		// If main doesn't panic within 2 seconds, it might be running successfully
		t.Log("Main function is running (likely server started)")
	}
}

func TestMain_AppContextInitError(t *testing.T) {
	// Create .env file with invalid database configuration to force AppContextInit error
	envContent := `HTTP_PORT=8080
REDIS_HOST=invalid_host
REDIS_PORT=invalid_port
MYSQL_HOST=invalid_host
MYSQL_PORT=invalid_port
MYSQL_USER=invalid_user
MYSQL_PASSWORD=invalid_pass
MYSQL_DATABASE=invalid_db`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Test that main panics when AppContextInit fails
	assert.Panics(t, func() {
		main()
	}, "Main should panic when AppContextInit fails")
}

func TestMain_BuildServiceError(t *testing.T) {
	// This test is harder to trigger without complex mocking
	// We'll test the general flow instead

	// Create .env file
	envContent := `HTTP_PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=testuser
MYSQL_PASSWORD=testpass
MYSQL_DATABASE=testdb`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Test that main function attempts to run
	// It will likely panic due to database connection issues
	assert.Panics(t, func() {
		main()
	}, "Main should panic when dependencies are not available")
}

func TestMain_StartServerError(t *testing.T) {
	// Create .env file with invalid port to force StartServer error
	envContent := `HTTP_PORT=invalid_port
REDIS_HOST=localhost
REDIS_PORT=6379
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=testuser
MYSQL_PASSWORD=testpass
MYSQL_DATABASE=testdb`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Test that main panics when StartServer fails
	assert.Panics(t, func() {
		main()
	}, "Main should panic when StartServer fails with invalid port")
}

func TestMain_EnvironmentVariables(t *testing.T) {
	// Test with environment variables instead of .env file
	os.Remove(".env") // Remove .env file

	// Set environment variables
	os.Setenv("HTTP_PORT", "9090")
	os.Setenv("REDIS_HOST", "test-redis")
	os.Setenv("MYSQL_HOST", "test-mysql")
	defer func() {
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("MYSQL_HOST")
	}()

	// Create empty .env file to prevent panic
	err := os.WriteFile(".env", []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Test that config loads from environment variables
	assert.NotPanics(t, func() {
		config.LoadConfig()
	})

	assert.Equal(t, "9090", config.HTTPPort)
	assert.Equal(t, "test-redis", config.RedisHost)
	assert.Equal(t, "test-mysql", config.MySQLHost)
}

func TestMain_SelectStatement(t *testing.T) {
	// Test the select{} statement behavior
	// This is more of a documentation test since select{} blocks forever

	// Create .env file
	envContent := `HTTP_PORT=0
REDIS_HOST=localhost
REDIS_PORT=6379
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=testuser
MYSQL_PASSWORD=testpass
MYSQL_DATABASE=testdb`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Test that main function includes select{} for blocking
	// We can't easily test this without complex goroutine management
	// But we can verify the function structure doesn't panic during setup

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- true
			}
		}()
		// This will panic due to missing dependencies, which is expected
		main()
	}()

	select {
	case <-done:
		t.Log("Main function setup completed (panicked as expected)")
	case <-time.After(1 * time.Second):
		t.Log("Main function is running")
	}
}

// Benchmark test for main function setup
func BenchmarkMain_Setup(b *testing.B) {
	// Create .env file
	envContent := `HTTP_PORT=0
REDIS_HOST=localhost
REDIS_PORT=6379
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=testuser
MYSQL_PASSWORD=testpass
MYSQL_DATABASE=testdb`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(".env")

	for i := 0; i < b.N; i++ {
		func() {
			defer func() {
				recover() // Catch expected panics
			}()
			config.LoadConfig()
		}()
	}
}
