package analysis

import (
	"math/rand"
	"time"

	"github.com/bcdannyboy/dgws/risk"
	"github.com/bcdannyboy/dgws/risk/statistics"
)

// UpdateEventProbabilityWithDependency updates the event probability based on the outcome of its dependencies using Bayesian principles.
func UpdateEventProbabilityWithDependency(event *risk.Event, eventsOccurred map[int]bool, eventProbabilities map[int]float64) float64 {
	if len(event.Dependencies) == 0 {
		return eventProbabilities[event.ID]
	}

	// Adjust event probability based on dependencies
	for _, dependency := range event.Dependencies {
		dependencyEventProbability := eventProbabilities[dependency.DependsOnEventID]
		dependencyOccurred := eventsOccurred[dependency.DependsOnEventID]

		if dependency.Happens && !dependencyOccurred {
			// If the event depended on this happening and it didn't, event cannot happen
			return 0
		} else if !dependency.Happens && dependencyOccurred {
			// If the event depended on this not happening but it did, greatly reduce the probability
			return eventProbabilities[event.ID] * (1 - dependencyEventProbability)
		}
	}

	return eventProbabilities[event.ID]
}

// SimulateEvent checks if an event happens based on its probability and dependencies.
func SimulateEvent(event *risk.Event, eventsOccurred map[int]bool, eventProbabilities map[int]float64) (bool, map[string]float64) {
	adjustedProbability := UpdateEventProbabilityWithDependency(event, eventsOccurred, eventProbabilities)

	rand.Seed(time.Now().UnixNano())
	if rand.Float64() <= adjustedProbability {
		eventsOccurred[event.ID] = true
		impacts := calculateImpacts(event)
		return true, impacts
	}
	eventsOccurred[event.ID] = false
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
