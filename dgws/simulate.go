package dgws

import (
	"fmt"

	"github.com/bcdannyboy/montecargo/dgws/simulator"
	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

func SimulateIndependentEvents(Events []*utils.FilteredEvent, iterations int) ([]*types.SimulationResults, error) {
	SimulationResults := []*types.SimulationResults{}

	for _, event := range Events {
		for i := 0; i < iterations; i++ {
			if event.Event.AssociatedProbability == nil && event.Event.AssociatedImpact == nil && event.Event.AssociatedCost == nil {
				return nil, fmt.Errorf("event %v has no associated probability, impact or cost", event.ID)
			}

			// probability
			if event.Event.AssociatedProbability != nil {
				// default probability item
				if event.Event.AssociatedProbability.SingleNumber != nil {
					// associated probability single number probability
					associatedProbability, err := simulator.SimulateIndependentSingleNumer(event.Event.AssociatedProbability.SingleNumber, event.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID:     event.ID,
						Probability: *associatedProbability,
					})
				} else if event.Event.AssociatedProbability.Range != nil {
					// associated probability range probability
					associatedProbability, err := simulator.SimulateIndependentRange(event.Event.AssociatedProbability.Range, event.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID:     event.ID,
						Probability: *associatedProbability,
					})
				} else if event.Event.AssociatedProbability.Decomposed != nil {
					// associated probability decomposed probability
					associatedProbability, associatedProbabilityStandardDeviation, _, _, _, _, err := simulator.SimulateIndependentDecomposed(event.Event.AssociatedProbability.Decomposed)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID:                      event.ID,
						Probability:                  associatedProbability,
						ProbabilityStandardDeviation: associatedProbabilityStandardDeviation,
					})
				}
			}

			// impact
			if event.Event.AssociatedImpact != nil {
				// default impact item
				if event.Event.AssociatedImpact.SingleNumber != nil {
					// associated impact single number impact
					associatedImpact, err := simulator.SimulateIndependentSingleNumer(event.Event.AssociatedImpact.SingleNumber, event.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID:      event.ID,
						Impact:       *associatedImpact,
						IsCostSaving: event.Event.AssociatedImpact.IsCostSaving,
					})
				} else if event.Event.AssociatedImpact.Range != nil {
					// associated impact range impact
					associatedImpact, err := simulator.SimulateIndependentRange(event.Event.AssociatedImpact.Range, event.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID:      event.ID,
						Impact:       *associatedImpact,
						IsCostSaving: event.Event.AssociatedImpact.IsCostSaving,
					})
				} else if event.Event.AssociatedImpact.Decomposed != nil {
					// associated impact decomposed impact
					associatedImpact, associatedImpactStandardDeviation, _, _, _, _, err := simulator.SimulateIndependentDecomposed(event.Event.AssociatedImpact.Decomposed)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID:                 event.ID,
						Impact:                  associatedImpact,
						ImpactStandardDeviation: associatedImpactStandardDeviation,
						IsCostSaving:            event.Event.AssociatedImpact.IsCostSaving,
					})
				}
			}

			// cost
			if event.Event.AssociatedCost != nil {
				// default cost item
				if event.Event.AssociatedCost.SingleNumber != nil {
					// associated cost single number cost
					associatedCost, err := simulator.SimulateIndependentSingleNumer(event.Event.AssociatedCost.SingleNumber, event.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID: event.ID,
						Cost:    *associatedCost,
					})
				} else if event.Event.AssociatedCost.Range != nil {
					// associated cost range cost
					associatedCost, err := simulator.SimulateIndependentRange(event.Event.AssociatedCost.Range, event.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID: event.ID,
						Cost:    *associatedCost,
					})
				} else if event.Event.AssociatedCost.Decomposed != nil {
					// associated cost decomposed cost
					associatedCost, associatedCostStandardDeviation, _, _, _, _, err := simulator.SimulateIndependentDecomposed(event.Event.AssociatedCost.Decomposed)
					if err != nil {
						return nil, err
					}

					SimulationResults = append(SimulationResults, &types.SimulationResults{
						EventID:               event.ID,
						Cost:                  associatedCost,
						CostStandardDeviation: associatedCostStandardDeviation,
					})
				}
			}

		}
	}

	return SimulationResults, nil
}

func SimulateDependentEvents(Events []*utils.FilteredEvent, iterations int) ([]*types.SimulationResults, error) {
	SimulationResults := []*types.SimulationResults{}

	for _, event := range Events {
		if !event.Independent { // we're only interested in dependent events

			dType := event.DependencyType
			dID := event.DependentEventID

			if dID == nil {
				return nil, fmt.Errorf("dependent event ID is nil for event %v", event.ID)
			}

			dEvent, err := utils.FindEventByID(*dID, Events)
			if err != nil {
				return nil, fmt.Errorf("failed to find dependent event %v: %s", *dID, err.Error())
			}

			for i := 0; i < iterations; i++ {
				// check for dependencies until all are met or at least one is missed
				depMet, err := simulator.DependencyCheck(dEvent, dType, Events)
				if err != nil {
					return nil, fmt.Errorf("failed to check dependencies for event %v: %s", event.ID, err.Error())
				}

				// if a dependency is missed, skip the iteration
				if !depMet {
					continue
				}

				// if all dependencies are met, simulate the event and store the results
				simResults, err := simulator.SimulateIndependentEvent(event)
				if err != nil {
					return nil, fmt.Errorf("failed to simulate event %v: %s", event.ID, err.Error())
				}

				SimulationResults = append(SimulationResults, simResults)
			}
		}
	}

	return SimulationResults, nil
}
