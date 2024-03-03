package statistics

import (
	"math/rand"
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

// CalculatePERTSample calculates a random value from the PERT distribution given a minimum, maximum, and mode (most likely value).
func CalculatePERTSample(min, max, mode float64) float64 {
	// Seed the random number generator (consider doing this once during initialization)
	rand.Seed(time.Now().UnixNano())

	// Calculate the parameters for the Beta distribution
	alpha := 1 + ((4 * (mode - min)) / (max - min))
	beta := 1 + ((4 * (max - mode)) / (max - min))

	// Create a Beta distribution using the standard math/rand package
	betaDist := distuv.Beta{
		Alpha: alpha,
		Beta:  beta,
	}

	// Generate a sample from the Beta distribution
	sample := betaDist.Rand()

	// Scale and shift the Beta distribution sample to fit the PERT distribution
	pertSample := min + (sample * (max - min))

	return pertSample
}
