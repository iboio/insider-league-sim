package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	WeightPoints    float64
	WeightsStrength float64
	HTTPPort        string
	RedisHost       string
	RedisPort       string
	MySQLHost       string
	MySQLPort       string
	MySQLUser       string
	MySQLPassword   string
	MySQLDatabase   string
	CTXTimeout      context.Context
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		panic("Error loading .env file")
	}
	WeightPoints = 0.4
	WeightsStrength = 0.6
	HTTPPort = getEnv("HTTP_PORT", "8080")
	RedisHost = getEnv("REDIS_HOST", "localhost")
	RedisPort = getEnv("REDIS_PORT", "4000")
	MySQLHost = getEnv("MYSQL_HOST", "localhost")
	MySQLPort = getEnv("MYSQL_PORT", "4050")
	MySQLUser = getEnv("MYSQL_USER", "iboio")
	MySQLPassword = getEnv("MYSQL_PASSWORD", "1234")
	MySQLDatabase = getEnv("MYSQL_DATABASE", "league_sim")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	CTXTimeout = ctx
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
