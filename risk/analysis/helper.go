package analysis

func clampProbability(p float64) float64 {
	if p < 0 {
		return 0
	}
	if p > 1 {
		return 1
	}
	return p
}
