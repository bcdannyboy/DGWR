package statistics

import (
	"math/rand"
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

// GenerateBetaSample generates a sample from a beta distribution
// for a given probability p and confidence level c.
func GenerateBetaSample(p float64, c float64) float64 {
	rand.Seed(time.Now().UnixNano())

	// Ensure p is within valid range for a beta distribution
	if p <= 0 {
		p = 0.01 // Assign a small probability if p is less or equal to 0
	} else if p >= 1 {
		p = 0.99 // Assign a high probability close to 1 if p is greater or equal to 1
	}

	// Ensure confidence level c is positive
	if c <= 0 {
		c = -c // Use the absolute value of c if it's negative
	}

	// Calculate alpha and beta
	alpha := p * c
	beta := (1 - p) * c

	// Create and sample from the beta distribution
	betaDist := distuv.Beta{Alpha: alpha, Beta: beta}
	sample := betaDist.Rand()

	return sample
}

// GenerateLHSSamples generates Latin Hypercube Samples for a given range and sample size.
func GenerateLHSSamples(min, max float64, n int) []float64 {
	rand.Seed(time.Now().UnixNano())

	step := (max - min) / float64(n)
	samples := make([]float64, n)
	for i := range samples {
		samples[i] = rand.Float64()*step + float64(i)*step + min
	}

	rand.Shuffle(len(samples), func(i, j int) { samples[i], samples[j] = samples[j], samples[i] })
	return samples
}
