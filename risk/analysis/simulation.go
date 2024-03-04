package analysis

import (
	"math/rand"
	"time"

	"github.com/bcdannyboy/dgws/risk"
	"github.com/bcdannyboy/dgws/risk/statistics"
)

func adjustProbabilityForDependencies(event *risk.Event, eventsOccurred map[int]bool, eventProbabilities map[int]float64) float64 {
	// The base probability of the event without considering dependencies.
	baseProbability := eventProbabilities[event.ID]

	for _, dependency := range event.Dependencies {
		dependencyOccurred := eventsOccurred[dependency.DependsOnEventID]

		if dependency.Happens {
			if !dependencyOccurred {
				return 0
			}
		} else {
			if dependencyOccurred {
				return 0
			}
		}
	}

	return baseProbability
}

// UpdateEventProbabilityWithDependency updates the event probability based on the outcome of its dependencies using Bayesian principles.
func UpdateEventProbabilityWithDependency(event *risk.Event, eventsOccurred map[int]bool, eventProbabilities map[int]float64) float64 {
	adjustedProbability := adjustProbabilityForDependencies(event, eventsOccurred, eventProbabilities)
	// Iterate over each dependency to apply Bayesian updating
	for _, dependency := range event.Dependencies {
		dependencyProbability := eventProbabilities[dependency.DependsOnEventID]
		if dependency.Happens {
			// For dependencies that must happen, use the positive impact on the probability.
			adjustedProbability *= dependencyProbability
		} else {
			// For dependencies that must not happen, use the negative impact on the probability.
			adjustedProbability *= (1 - dependencyProbability)
		}
	}
	// Ensure the probability is within [0,1]
	return clampProbability(adjustedProbability)
}

// SimulateEvent checks if an event happens based on its probability and dependencies.
func SimulateEvent(event *risk.Event, eventsOccurred map[int]bool, eventProbabilities map[int]float64) (bool, map[string]float64) {
	// Adjust the probability based on dependencies
	adjustedProbability := adjustProbabilityForDependencies(event, eventsOccurred, eventProbabilities)

	// Decide if the event happens based on the adjusted probability
	rand.Seed(time.Now().UnixNano())
	if rand.Float64() <= adjustedProbability {
		// Event happens, calculate impacts
		impacts := calculateImpacts(event)
		return true, impacts
	}
	return false, nil
}

// calculateImpacts calculates the impacts for an event.
func calculateImpacts(event *risk.Event) map[string]float64 {
	impacts := make(map[string]float64)
	for _, impact := range event.Impact {
		// Simplified example of impact calculation
		impactValue := (impact.MinimumIndividualUnitImpact + impact.MaximumIndividualUnitImpact) / 2
		impacts[impact.Unit] += impactValue
	}
	return impacts
}

// MonteCarlo simulates the risk event network a specified number of times,
// adjusting for dependencies using Bayesian statistics.
func MonteCarlo(events []*risk.Event, iterations int) (map[int]float64, map[string]float64, error) {
	eventProbabilities := make(map[int]float64)
	totalImpacts := make(map[string]float64)
	eventOccurrences := make(map[int]int)
	impactOccurrences := make(map[string]int)

	// Initialize probabilities with a reasonable estimate
	for _, event := range events {
		initMin := statistics.GenerateBetaSample(event.Probability.Minimum, event.Probability.MinimumConfidence)
		initMax := statistics.GenerateBetaSample(event.Probability.Maximum, event.Probability.MaximumConfidence)
		eventProbabilities[event.ID] = (initMax + initMin) / 2
	}

	for i := 0; i < iterations; i++ {
		eventsOccurred := make(map[int]bool)
		rand.Seed(time.Now().UnixNano())

		for _, event := range events {
			happened, impacts := SimulateEvent(event, eventsOccurred, eventProbabilities)
			eventsOccurred[event.ID] = happened
			if happened {
				eventOccurrences[event.ID]++
				for impactType, impactValue := range impacts {
					totalImpacts[impactType] += impactValue
					impactOccurrences[impactType]++ // Increment count for this impact type.
				}
			}
		}
	}

	// Normalize the total impacts based on the number of occurrences, not iterations.
	for impactType, totalValue := range totalImpacts {
		if occurrences, ok := impactOccurrences[impactType]; ok && occurrences > 0 {
			totalImpacts[impactType] = totalValue / float64(occurrences)
		}
	}

	// Adjust probabilities based on occurrences
	for eventID := range eventProbabilities {
		if occurrences, found := eventOccurrences[eventID]; found {
			eventProbabilities[eventID] = float64(occurrences) / float64(iterations)
		} else {
			eventProbabilities[eventID] = 0 // Event never occurred.
		}
	}

	return eventProbabilities, totalImpacts, nil
}
