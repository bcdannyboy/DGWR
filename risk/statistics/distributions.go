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

// CalculatePERTSample safely calculates a PERT sample.
func CalculatePERTSample(min, max, mode float64) float64 {
	// Seed RNG - consider seeding globally or during initialization instead
	rand.Seed(time.Now().UnixNano())

	// Ensure min < max to avoid division by zero
	if min == max {
		max = min + 0.0001 // Slight adjustment to avoid division by zero
	}

	// Ensure mode is within [min, max]
	if mode < min {
		mode = min
	} else if mode > max {
		mode = max
	}

	// Calculate alpha and beta using the PERT formula
	alpha := 1 + ((mode-min)/(max-min))*4
	beta := 1 + ((max-mode)/(max-min))*4

	// Generate and return sample
	betaDist := distuv.Beta{Alpha: alpha, Beta: beta}
	pertSample := min + betaDist.Rand()*(max-min) // Scale and shift the beta sample
	return pertSample
}
