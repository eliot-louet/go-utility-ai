package math_utils

import "math"

func Normalize(x, min, max float64) float64 {

	v := (x - min) / (max - min)

	if v < 0 {
		return 0
	}

	if v > 1 {
		return 1
	}

	return v
}

func GeometricMean(values []float64) float64 {

	if len(values) == 0 {
		return 0
	}

	product := 1.0

	for _, v := range values {
		product *= v
	}

	return math.Pow(product, 1.0/float64(len(values)))
}
