package analysis

import (
	"math/rand"
	"time"

	"github.com/bcdannyboy/dgws/risk"
	"github.com/bcdannyboy/dgws/risk/statistics"
)

// SimulateEvent checks if an event happens based on its probability and dependencies, using Bayesian statistics.
func SimulateEvent(event *risk.Event, eventsOccurred map[int]bool, eventProbabilities map[int]float64) (bool, map[string]float64) {
	// Initial probability of the event
	baseProbability := eventProbabilities[event.ID]

	// Adjust the probability based on dependencies using Bayesian statistics
	for _, dependency := range event.Dependencies {
		if depEventProb, exists := eventProbabilities[dependency.DependsOnEventID]; exists {
			if dependency.Happens && eventsOccurred[dependency.DependsOnEventID] {
				// Increase the probability if the dependency happened as required
				baseProbability = baseProbability * depEventProb / clampProbability(baseProbability)
			} else if !dependency.Happens && !eventsOccurred[dependency.DependsOnEventID] {
				// Adjust the probability if the dependency did not happen as required
				// This could be more complex depending on the specific logic needed
				baseProbability = baseProbability * (1 - depEventProb) / clampProbability(baseProbability)
			}
		}
	}

	// Ensure the probability is within [0, 1]
	adjustedProbability := clampProbability(baseProbability)

	// Decide if the event happens based on the adjusted probability
	rand.Seed(time.Now().UnixNano())
	if rand.Float64() <= adjustedProbability {
		// Event happens, calculate impacts
		impacts := calculateImpacts(event)
		return true, impacts
	}
	return false, nil
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
