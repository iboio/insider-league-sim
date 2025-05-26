package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success(t *testing.T) {
	// Create a temporary .env file
	envContent := `HTTP_PORT=9090
REDIS_HOST=test-redis
REDIS_PORT=6379
MYSQL_HOST=test-mysql
MYSQL_PORT=3306
MYSQL_USER=testuser
MYSQL_PASSWORD=testpass
MYSQL_DATABASE=testdb`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Execute
	LoadConfig()

	// Assert
	assert.Equal(t, 0.4, WeightPoints)
	assert.Equal(t, 0.6, WeightsStrength)
	assert.Equal(t, "9090", HTTPPort)
	assert.Equal(t, "test-redis", RedisHost)
	assert.Equal(t, "6379", RedisPort)
	assert.Equal(t, "test-mysql", MySQLHost)
	assert.Equal(t, "3306", MySQLPort)
	assert.Equal(t, "testuser", MySQLUser)
	assert.Equal(t, "testpass", MySQLPassword)
	assert.Equal(t, "testdb", MySQLDatabase)
	assert.NotNil(t, CTXTimeout)
}

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Remove .env file if exists
	os.Remove(".env")

	// Clear environment variables
	os.Unsetenv("HTTP_PORT")
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("MYSQL_HOST")
	os.Unsetenv("MYSQL_PORT")
	os.Unsetenv("MYSQL_USER")
	os.Unsetenv("MYSQL_PASSWORD")
	os.Unsetenv("MYSQL_DATABASE")

	// Create empty .env file
	err := os.WriteFile(".env", []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(".env")

	// Execute
	LoadConfig()

	// Assert defaults
	assert.Equal(t, 0.4, WeightPoints)
	assert.Equal(t, 0.6, WeightsStrength)
	assert.Equal(t, "8080", HTTPPort)
	assert.Equal(t, "localhost", RedisHost)
	assert.Equal(t, "4000", RedisPort)
	assert.Equal(t, "localhost", MySQLHost)
	assert.Equal(t, "4050", MySQLPort)
	assert.Equal(t, "iboio", MySQLUser)
	assert.Equal(t, "1234", MySQLPassword)
	assert.Equal(t, "league_sim", MySQLDatabase)
}

func TestLoadConfig_NoEnvFile(t *testing.T) {
	// Remove .env file
	os.Remove(".env")

	// Execute and expect panic
	assert.Panics(t, func() {
		LoadConfig()
	}, "Should panic when .env file is missing")
}

func TestGetEnv_ExistingValue(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	// Execute
	result := getEnv("TEST_VAR", "default_value")

	// Assert
	assert.Equal(t, "test_value", result)
}

func TestGetEnv_DefaultValue(t *testing.T) {
	// Ensure environment variable doesn't exist
	os.Unsetenv("NON_EXISTENT_VAR")

	// Execute
	result := getEnv("NON_EXISTENT_VAR", "default_value")

	// Assert
	assert.Equal(t, "default_value", result)
}

func TestGetEnvAsInt_ExistingValidValue(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_INT_VAR", "123")
	defer os.Unsetenv("TEST_INT_VAR")

	// Execute
	result := getEnvAsInt("TEST_INT_VAR", 456)

	// Assert
	assert.Equal(t, 123, result)
}

func TestGetEnvAsInt_ExistingInvalidValue(t *testing.T) {
	// Set environment variable with invalid int
	os.Setenv("TEST_INT_VAR", "not_a_number")
	defer os.Unsetenv("TEST_INT_VAR")

	// Execute
	result := getEnvAsInt("TEST_INT_VAR", 456)

	// Assert
	assert.Equal(t, 456, result) // Should return default
}

func TestGetEnvAsInt_DefaultValue(t *testing.T) {
	// Ensure environment variable doesn't exist
	os.Unsetenv("NON_EXISTENT_INT_VAR")

	// Execute
	result := getEnvAsInt("NON_EXISTENT_INT_VAR", 789)

	// Assert
	assert.Equal(t, 789, result)
}

func TestGetIntEnv_ExistingValidValue(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_INT_VAR2", "321")
	defer os.Unsetenv("TEST_INT_VAR2")

	// Execute
	result := getIntEnv("TEST_INT_VAR2", 654)

	// Assert
	assert.Equal(t, 321, result)
}

func TestGetIntEnv_ExistingInvalidValue(t *testing.T) {
	// Set environment variable with invalid int
	os.Setenv("TEST_INT_VAR2", "invalid")
	defer os.Unsetenv("TEST_INT_VAR2")

	// Execute
	result := getIntEnv("TEST_INT_VAR2", 654)

	// Assert
	assert.Equal(t, 654, result) // Should return default
}

func TestGetIntEnv_DefaultValue(t *testing.T) {
	// Ensure environment variable doesn't exist
	os.Unsetenv("NON_EXISTENT_INT_VAR2")

	// Execute
	result := getIntEnv("NON_EXISTENT_INT_VAR2", 987)

	// Assert
	assert.Equal(t, 987, result)
}

// Benchmark tests
func BenchmarkLoadConfig(b *testing.B) {
	// Create a temporary .env file
	envContent := `HTTP_PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(".env")

	for i := 0; i < b.N; i++ {
		LoadConfig()
	}
}

func BenchmarkGetEnv(b *testing.B) {
	os.Setenv("BENCH_VAR", "bench_value")
	defer os.Unsetenv("BENCH_VAR")

	for i := 0; i < b.N; i++ {
		getEnv("BENCH_VAR", "default")
	}
}

func BenchmarkGetEnvAsInt(b *testing.B) {
	os.Setenv("BENCH_INT_VAR", "123")
	defer os.Unsetenv("BENCH_INT_VAR")

	for i := 0; i < b.N; i++ {
		getEnvAsInt("BENCH_INT_VAR", 456)
	}
}
