package league

import (
	"testing"

	"league-sim/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestRandomNumberGenerator(t *testing.T) {
	tests := []struct {
		name string
		min  float64
		max  float64
	}{
		{
			name: "Normal range",
			min:  1.0,
			max:  10.0,
		},
		{
			name: "Zero to positive",
			min:  0.0,
			max:  5.0,
		},
		{
			name: "Negative to positive",
			min:  -5.0,
			max:  5.0,
		},
		{
			name: "Small range",
			min:  70.0,
			max:  100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run multiple times to test randomness
			for i := 0; i < 100; i++ {
				result := RandomNumberGenerator(tt.min, tt.max)

				// Check that result is within range
				assert.GreaterOrEqual(t, result, tt.min, "Result should be >= min")
				assert.LessOrEqual(t, result, tt.max, "Result should be <= max")
			}
		})
	}
}

func TestRandomNumberGenerator_EqualMinMax(t *testing.T) {
	// Test edge case where min equals max
	min := 5.0
	max := 5.0

	result := RandomNumberGenerator(min, max)

	assert.Equal(t, min, result, "When min equals max, should return min")
}

func TestRandomNumberGenerator_MinGreaterThanMax(t *testing.T) {
	// Test edge case where min is greater than max
	min := 10.0
	max := 5.0

	result := RandomNumberGenerator(min, max)

	assert.Equal(t, min, result, "When min > max, should return min")
}

func TestRandomNumberGenerator_Distribution(t *testing.T) {
	// Test that the function produces different values (not always the same)
	min := 1.0
	max := 100.0

	results := make(map[float64]bool)

	// Generate 50 random numbers
	for i := 0; i < 50; i++ {
		result := RandomNumberGenerator(min, max)
		results[result] = true
	}

	// We should have multiple different values (at least 10 different ones)
	assert.GreaterOrEqual(t, len(results), 10, "Should generate diverse random numbers")
}

func TestCalculateStrength(t *testing.T) {
	tests := []struct {
		name           string
		team           models.Team
		expectedResult float64
	}{
		{
			name: "Balanced team",
			team: models.Team{
				Name:         "Test Team",
				AttackPower:  80.0,
				DefensePower: 80.0,
				Stamina:      80.0,
				Morale:       80.0,
			},
			expectedResult: 80.0, // (80*0.3 + 80*0.3 + 80*0.2 + 80*0.2) = 80
		},
		{
			name: "High attack team",
			team: models.Team{
				Name:         "Attack Team",
				AttackPower:  100.0,
				DefensePower: 60.0,
				Stamina:      70.0,
				Morale:       70.0,
			},
			expectedResult: 76.0, // (100*0.3 + 60*0.3 + 70*0.2 + 70*0.2) = 30 + 18 + 14 + 14 = 76
		},
		{
			name: "High defense team",
			team: models.Team{
				Name:         "Defense Team",
				AttackPower:  60.0,
				DefensePower: 100.0,
				Stamina:      70.0,
				Morale:       70.0,
			},
			expectedResult: 76.0, // (60*0.3 + 100*0.3 + 70*0.2 + 70*0.2) = 18 + 30 + 14 + 14 = 76
		},
		{
			name: "Weak team",
			team: models.Team{
				Name:         "Weak Team",
				AttackPower:  50.0,
				DefensePower: 50.0,
				Stamina:      50.0,
				Morale:       50.0,
			},
			expectedResult: 50.0, // (50*0.3 + 50*0.3 + 50*0.2 + 50*0.2) = 50
		},
		{
			name: "Zero values team",
			team: models.Team{
				Name:         "Zero Team",
				AttackPower:  0.0,
				DefensePower: 0.0,
				Stamina:      0.0,
				Morale:       0.0,
			},
			expectedResult: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateStrength(tt.team)

			assert.InDelta(t, tt.expectedResult, result, 0.001, "Calculated strength should match expected value")
		})
	}
}

func TestCalculateStrength_WeightDistribution(t *testing.T) {
	// Test that the weights are correctly applied
	team := models.Team{
		Name:         "Weight Test Team",
		AttackPower:  100.0, // 30% weight
		DefensePower: 0.0,   // 30% weight
		Stamina:      0.0,   // 20% weight
		Morale:       0.0,   // 20% weight
	}

	result := CalculateStrength(team)
	expected := 100.0 * 0.3 // Only attack power contributes

	assert.InDelta(t, expected, result, 0.001, "Attack power should contribute 30% to total strength")
}

// Benchmark tests
func BenchmarkRandomNumberGenerator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomNumberGenerator(70.0, 100.0)
	}
}

func BenchmarkCalculateStrength(b *testing.B) {
	team := models.Team{
		Name:         "Benchmark Team",
		AttackPower:  85.5,
		DefensePower: 78.2,
		Stamina:      82.1,
		Morale:       79.8,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateStrength(team)
	}
}
