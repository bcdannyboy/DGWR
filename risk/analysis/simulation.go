package analysis

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/bcdannyboy/dgws/risk"
	"github.com/bcdannyboy/dgws/risk/statistics"
	"github.com/bcdannyboy/dgws/risk/utils"
)

// returns:
// bool - whether the event happens
// float64 - a map of impacts by unit
// error - if the event with the given ID is not found
func SimulateEvent(ID int, events []*risk.Event) (bool, map[string]float64, error) {
	event := utils.FindEvent(ID, events)
	if event == nil {
		return false, nil, errors.New(fmt.Sprintf("event with ID %d not found", ID))
	}

	if len(event.Dependencies) > 0 {
		for _, dependency := range event.Dependencies {
			depHappens, _, err := SimulateEvent(dependency.DependsOnEventID, events)
			if err != nil {
				return false, nil, err
			}

			if (dependency.Happens && !depHappens) || (!dependency.Happens && depHappens) {
				// If the dependency happens and the event doesn't want it to, or the dependency doesn't happen and the event does, then the event doesn't happen
				return false, nil, nil
			}
		}
	}

	rand.Seed(time.Now().UnixNano())
	MinProbabilityCoinFlip := rand.Float64() > 0.5
	MaxProbabilityCoinFlip := rand.Float64() > 0.5

	MinConfidenceImpact := rand.Float64() * event.Probability.MinimumConfidence
	MaxConfidenceImpact := rand.Float64() * event.Probability.MaximumConfidence

	MinProbability := event.Probability.Minimum
	MaxProbability := event.Probability.Maximum

	if MinProbabilityCoinFlip {
		MinProbability = MinProbability - MinConfidenceImpact
	} else {
		MinProbability = MinProbability + MinConfidenceImpact
	}

	if MaxProbabilityCoinFlip {
		MaxProbability = MaxProbability - MaxConfidenceImpact
	} else {
		MaxProbability = MaxProbability + MaxConfidenceImpact
	}

	if MinProbability < 0 {
		MinProbability = 0
	}

	AverageProbability := (MinProbability + MaxProbability) / 2

	TimeFrameAdjustedMinProbability := utils.AdjustForTime(MinProbability, event.Probability.ExpectedFrequency)
	TimeFrameAdjustedMaxProbability := utils.AdjustForTime(MaxProbability, event.Probability.ExpectedFrequency)
	TimeFrameAdjustedAverageProbability := utils.AdjustForTime(AverageProbability, event.Probability.ExpectedFrequency)

	ProbabilityPERTSample := statistics.CalculatePERTSample(TimeFrameAdjustedMinProbability, TimeFrameAdjustedMaxProbability, TimeFrameAdjustedAverageProbability)

	if rand.Float64() < ProbabilityPERTSample {
		if len(event.Impact) == 0 {
			return true, nil, nil
		}

		totalImpacts := make(map[string]float64)

		for _, impact := range event.Impact {
			MinIndividualUnitImpactCoinFlip := rand.Float64() > 0.5
			MaxIndividualUnitImpactCoinFlip := rand.Float64() > 0.5
			MinImpactEventsCoinFlip := rand.Float64() > 0.5
			MaxImpactEventsCoinFlip := rand.Float64() > 0.5

			MinConfidenceIndividualUnitImpact := rand.Float64() * impact.MinimumIndividualUnitImpactConfidence
			MaxConfidenceIndividualUnitImpact := rand.Float64() * impact.MaximumIndividualUnitImpactConfidence
			MinConfidenceImpactEvents := rand.Float64() * impact.MinimumImpactEventsConfidence
			MaxConfidenceImpactEvents := rand.Float64() * impact.MaximumImpactEventsConfidence

			MinIndividualUnitImpact := impact.MinimumIndividualUnitImpact
			MaxIndividualUnitImpact := impact.MaximumIndividualUnitImpact
			MinImpactEvents := impact.MinimumImpactEvents
			MaxImpactEvents := impact.MaximumImpactEvents

			if MinIndividualUnitImpactCoinFlip {
				MinIndividualUnitImpact = MinIndividualUnitImpact - MinConfidenceIndividualUnitImpact
			} else {
				MinIndividualUnitImpact = MinIndividualUnitImpact + MinConfidenceIndividualUnitImpact
			}

			if MaxIndividualUnitImpactCoinFlip {
				MaxIndividualUnitImpact = MaxIndividualUnitImpact - MaxConfidenceIndividualUnitImpact
			} else {
				MaxIndividualUnitImpact = MaxIndividualUnitImpact + MaxConfidenceIndividualUnitImpact
			}

			if MinImpactEventsCoinFlip {
				MinImpactEvents = MinImpactEvents - MinConfidenceImpactEvents
			} else {
				MinImpactEvents = MinImpactEvents + MinConfidenceImpactEvents
			}

			if MaxImpactEventsCoinFlip {
				MaxImpactEvents = MaxImpactEvents - MaxConfidenceImpactEvents
			} else {
				MaxImpactEvents = MaxImpactEvents + MaxConfidenceImpactEvents
			}

			if !impact.PositiveImpact {
				// if its not a postivie impact then we need to ensure the value is no less than 0
				if MinIndividualUnitImpact < 0 {
					MinIndividualUnitImpact = 0
				}
				if MaxIndividualUnitImpact < 0 {
					MaxIndividualUnitImpact = 0
				}
				if MinImpactEvents < 0 {
					MinImpactEvents = 0
				}
				if MaxImpactEvents < 0 {
					MaxImpactEvents = 0
				}
			}

			AverageIndividualUnitImpact := (MinIndividualUnitImpact + MaxIndividualUnitImpact) / 2
			AverageImpactEvents := (MinImpactEvents + MaxImpactEvents) / 2

			TimeFrameAdjustedMinIndividualUnitImpact := utils.AdjustForTime(MinIndividualUnitImpact, impact.ExpectedFrequency)
			TimeFrameAdjustedMaxIndividualUnitImpact := utils.AdjustForTime(MaxIndividualUnitImpact, impact.ExpectedFrequency)
			TimeFrameAdjustedAverageIndividualUnitImpact := utils.AdjustForTime(AverageIndividualUnitImpact, impact.ExpectedFrequency)

			TimeFrameAdjustedMinImpactEvents := utils.AdjustForTime(MinImpactEvents, impact.ExpectedFrequency)
			TimeFrameAdjustedMaxImpactEvents := utils.AdjustForTime(MaxImpactEvents, impact.ExpectedFrequency)
			TimeFrameAdjustedAverageImpactEvents := utils.AdjustForTime(AverageImpactEvents, impact.ExpectedFrequency)

			IndividualUnitImpactPERTSample := statistics.CalculatePERTSample(TimeFrameAdjustedMinIndividualUnitImpact, TimeFrameAdjustedMaxIndividualUnitImpact, TimeFrameAdjustedAverageIndividualUnitImpact)
			ImpactEventsPERTSample := statistics.CalculatePERTSample(TimeFrameAdjustedMinImpactEvents, TimeFrameAdjustedMaxImpactEvents, TimeFrameAdjustedAverageImpactEvents)

			if ImpactEventsPERTSample < 0 {
				ImpactEventsPERTSample = 1
			}

			if IndividualUnitImpactPERTSample < 0 {
				IndividualUnitImpactPERTSample = 0
			}

			TotalImpact := float64(IndividualUnitImpactPERTSample * ImpactEventsPERTSample)

			if !impact.PositiveImpact {
				if v, ok := totalImpacts[impact.Unit]; ok {
					totalImpacts[impact.Unit] = v + TotalImpact
				} else {
					totalImpacts[impact.Unit] = TotalImpact
				}
			} else {
				if v, ok := totalImpacts[impact.Unit]; ok {
					totalImpacts[impact.Unit] = v - TotalImpact
				} else {
					totalImpacts[impact.Unit] = -TotalImpact
				}
			}
		}

		return true, totalImpacts, nil
	}

	return false, nil, nil
}

// returns:
// map[int]float64 - a map of event IDs to their probabilities
// map[string]float64 - a map of impact units to their total impacts
// error - if the event with the given ID is not found
func MonteCarlo(events []*risk.Event, iterations int) (map[int]float64, map[string]float64, error) {
	eventProbabilities := make(map[int]float64)
	totalImpacts := make(map[string]float64)

	for _, event := range events {
		eventProbabilities[event.ID] = 0
	}

	for i := 0; i < iterations; i++ {
		for _, event := range events {
			happens, impacts, err := SimulateEvent(event.ID, events)
			if err != nil {
				return nil, nil, err
			}

			if happens {
				eventProbabilities[event.ID]++
				for k, v := range impacts {
					if val, ok := totalImpacts[k]; ok {
						totalImpacts[k] = val + v
					} else {
						totalImpacts[k] = v
					}
				}
			}
		}
	}

	// Normalize the probabilities and impacts
	for k, v := range eventProbabilities {
		eventProbabilities[k] = v / float64(iterations)
	}

	for k, v := range totalImpacts {
		totalImpacts[k] = v / float64(iterations)
	}

	return eventProbabilities, totalImpacts, nil
}
