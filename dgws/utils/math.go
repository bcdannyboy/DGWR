package utils

import (
	"math"
	"math/rand"
)

// LogNormalSample generates a sample from a log-normal distribution
func LogNormalSample(mean, stdDev float64) float64 {
	if stdDev <= 0 {
		return math.Exp(mean)
	}
	normalSample := rand.NormFloat64()*stdDev + mean
	return math.Exp(normalSample)
}

// LogNormalSampleInRange generates a sample from a log-normal distribution within a specified range
func LogNormalSampleInRange(mean, stdDevMin, stdDevMax, min, max float64) float64 {
	// Generate a random standard deviation within the specified range
	stdDev := stdDevMin + rand.Float64()*(stdDevMax-stdDevMin)

	// Generate a log-normal sample using the random standard deviation
	sample := LogNormalSample(mean, stdDev)

	// Scale the sample to fit within the specified range
	return min + (sample-math.Exp(mean))*(max-min)/(math.Exp(mean+stdDev*stdDev)-math.Exp(mean))
}

// ComputeCompositeLogNormal computes the composite log-normal distribution of a set of values
func ComputeCompositeLogNormal(values, stdDevs []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}

	var logSum, logStdDevSum float64
	for i, value := range values {
		logValue := math.Log(value)
		logSum += logValue
		logStdDevSum += math.Pow(stdDevs[i], 2)
	}

	mean := logSum / float64(len(values))
	stdDev := math.Sqrt(logStdDevSum / float64(len(values)))

	return math.Exp(mean), stdDev
}
