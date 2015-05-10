package action

func computeCpuAverage(cpuPerc []float64) float32 {
	sum := 0.0
	n := float64(len(cpuPerc))
	for _, cpu := range cpuPerc {
		sum += cpu
	}
	return float32(sum / n)
}
