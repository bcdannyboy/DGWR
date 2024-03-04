package analysis

import (
	"math"
	"math/rand"
	"time"

	"github.com/bcdannyboy/dgws/risk"
	"github.com/bcdannyboy/dgws/risk/statistics"
	"github.com/bcdannyboy/dgws/risk/utils"
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
	rand.Seed(time.Now().UnixNano())
	adjustedProbability := UpdateEventProbabilityWithDependency(event, eventsOccurred, eventProbabilities)

	if rand.Float64() <= adjustedProbability {
		eventsOccurred[event.ID] = true
		// Adjust impacts based on the event's role in the simulation, such as filtering phishing emails.
		impacts := calculateImpacts(event, eventProbabilities)
		return true, impacts
	}
	eventsOccurred[event.ID] = false
	return false, nil
}

// calculateImpacts has been updated to include confidence levels in its calculations.
func calculateImpacts(event *risk.Event, eventProbabilities map[int]float64) map[string]float64 {
	impacts := make(map[string]float64)
	for _, impact := range event.Impact {
		// Adjust the impact calculation to factor in confidence levels for both unit impacts and event numbers.
		actualUnitImpact, actualEvents := adjustImpactBasedOnEventProbability(impact, eventProbabilities[event.ID])

		totalImpact := actualUnitImpact * actualEvents
		if impact.PositiveImpact {
			totalImpact = -totalImpact // Adjust for positive impacts if necessary.
		}

		impacts[impact.Unit] += totalImpact
	}

	return impacts
}

// adjustImpactBasedOnEventProbability is refined to consider confidence levels.
func adjustImpactBasedOnEventProbability(impact *risk.Impact, eventProbability float64) (float64, float64) {
	// Calculate the average impact and events considering the confidence intervals.

	scaledMinIndividualUnitImpact := utils.AdjustForTime(impact.MinimumIndividualUnitImpact, impact.ExpectedFrequency)
	scaledMaxIndividualUnitImpact := utils.AdjustForTime(impact.MaximumIndividualUnitImpact, impact.ExpectedFrequency)
	scaledMinImpactEvents := utils.AdjustForTime(impact.MinimumImpactEvents, impact.ExpectedFrequency)
	scaledMaxImpactEvents := utils.AdjustForTime(impact.MaximumImpactEvents, impact.ExpectedFrequency)

	avgUnitImpact := weightedAverageWithConfidence(
		scaledMinIndividualUnitImpact, impact.MaximumIndividualUnitImpact,
		scaledMaxIndividualUnitImpact, impact.MaximumIndividualUnitImpactConfidence,
	)

	avgEvents := weightedAverageWithConfidence(
		scaledMinImpactEvents, impact.MaximumImpactEvents,
		scaledMaxImpactEvents, impact.MaximumImpactEventsConfidence,
	)

	// Optionally adjust avgEvents based on the event type and probability.
	if impact.Name == "Phishing Emails Detected" {
		// Example: Scale the number of detected emails based on the filter's effectiveness and confidence.
		confidenceScale := (impact.MinimumImpactEventsConfidence + impact.MaximumImpactEventsConfidence) / 2
		avgEvents *= eventProbability * confidenceScale
	}

	// Ensure the calculated values are within expected bounds.
	avgUnitImpact = clamp(avgUnitImpact, impact.MinimumIndividualUnitImpact, impact.MaximumIndividualUnitImpact)
	avgEvents = math.Round(clamp(avgEvents, impact.MinimumImpactEvents, impact.MaximumImpactEvents)) // Round events to nearest integer.

	return avgUnitImpact, avgEvents
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
		scaledMin := utils.AdjustForTime(event.Probability.Minimum, event.Probability.ExpectedFrequency)
		scaledMax := utils.AdjustForTime(event.Probability.Maximum, event.Probability.ExpectedFrequency)
		initMin := statistics.GenerateBetaSample(scaledMin, event.Probability.MinimumConfidence)
		initMax := statistics.GenerateBetaSample(scaledMax, event.Probability.MaximumConfidence)
		sampleProb := statistics.GenerateLHSSamples(initMin, initMax, 100)
		avgProb := calcAvg(sampleProb)
		eventProbabilities[event.ID] = clampProbability(avgProb)
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
