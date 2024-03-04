package analysis

func weightedAverage(min, max, minConf, maxConf float64) float64 {
	return (min*minConf + max*maxConf) / (minConf + maxConf)
}

func averageConfidence(minConf, maxConf float64) float64 {
	return (minConf + maxConf) / 2.0
}

// adjustStddevBasedOnConfidence adjusts the standard deviation based on the confidence levels.
func adjustStddevBasedOnConfidence(min, max, minConf, maxConf float64) float64 {
	rangeVal := max - min
	avgConf := (minConf + maxConf) / 2
	// This is a heuristic: lower confidence means higher uncertainty, thus larger standard deviation.
	return rangeVal * (1 - avgConf) / 2
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func clampProbability(p float64) float64 {
	if p < 0 {
		return 0
	}
	if p > 1 {
		return 1
	}
	return p
}

func calcAvg(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// weightedAverageWithConfidence calculates a weighted average of min and max values based on their confidences.
func weightedAverageWithConfidence(min, max, minConf, maxConf float64) float64 {
	weightedMin := min * minConf
	weightedMax := max * maxConf
	totalConf := minConf + maxConf
	if totalConf == 0 { // Avoid division by zero.
		return (min + max) / 2
	}
	return (weightedMin + weightedMax) / totalConf
}
