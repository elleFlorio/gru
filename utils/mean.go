package utils

func Mean(values []float64) float64 {
	if len(values) < 1 {
		return 0.0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}
