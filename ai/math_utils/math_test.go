package math_utils_test

import (
	"testing"

	"github.com/eliot-louet/go-utility-ai/ai/math_utils"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	testsCases := []struct {
		name     string
		x        float64
		min      float64
		max      float64
		expected float64
	}{
		{"within range", 5, 0, 10, 0.5},
		{"at min", 0, 0, 10, 0},
		{"at max", 10, 0, 10, 1},
		{"below range", -5, 0, 10, 0},
		{"above range", 15, 0, 10, 1},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			result := math_utils.Normalize(tc.x, tc.min, tc.max)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGeometricMean(t *testing.T) {
	testCases := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{"empty slice", []float64{}, 0},
		{"single value", []float64{5}, 5},
		{"multiple values", []float64{1, 2, 3, 4}, 2.213363839400643},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := math_utils.GeometricMean(tc.values)
			assert.InDelta(t, tc.expected, result, 1e-9)
		})
	}
}
