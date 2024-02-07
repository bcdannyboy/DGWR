package simulator

import (
	"errors"
	"fmt"

	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

func simulateDecomposedByAttribute(decomposed *types.Decomposed, attribute int) (float64, float64, error) {
	if decomposed == nil {
		return 0, 0, errors.New("decomposed is nil")
	}

	var values, stdDevs []float64

	for _, component := range decomposed.Components {
		var compResult float64
		var compStdDev float64
		var err error

		// Updated to use an integer attribute following the constants defined earlier
		switch attribute {
		case ProbabilityAttribute:
			if component.Probability != nil {
				comp := &Component{SingleNumber: component.Probability.SingleNumber, Range: component.Probability.Range, Decomposed: component.Probability.Decomposed}
				compResult, compStdDev, err = handleAttributeSimulation(comp, ProbabilityAttribute, component.TimeFrame)
			}
		case ImpactAttribute:
			if component.Impact != nil {
				comp := &Component{SingleNumber: component.Impact.SingleNumber, Range: component.Impact.Range, Decomposed: component.Impact.Decomposed}
				compResult, compStdDev, err = handleAttributeSimulation(comp, ImpactAttribute, component.TimeFrame)
			}
		case CostAttribute:
			if component.Cost != nil {
				comp := &Component{SingleNumber: component.Cost.SingleNumber, Range: component.Cost.Range, Decomposed: component.Cost.Decomposed}
				compResult, compStdDev, err = handleAttributeSimulation(comp, CostAttribute, component.TimeFrame)
			}
		default:
			return 0, 0, errors.New("invalid attribute specified")
		}

		if err != nil {
			return 0, 0, err
		}

		// Directly appending values as they are no longer pointers
		values = append(values, compResult)
		stdDevs = append(stdDevs, compStdDev)
	}

	if len(values) == 0 || len(stdDevs) == 0 {
		return 0, 0, errors.New("no valid components for simulation")
	}

	// Compute composite result using the utility function
	compositeValue, compositeStdDev := utils.ComputeCompositeLogNormal(values, stdDevs)

	return compositeValue, compositeStdDev, nil
}

func handleAttributeSimulation(component *Component, attribute int, timeFrame uint64) (float64, float64, error) {
	if component.SingleNumber != nil {
		result, err := SimulateIndependentSingleNumer(component.SingleNumber, timeFrame)
		if err != nil {
			return 0, 0, err
		}
		stdDev := 0.0 // Default to 0 if not provided
		if component.SingleNumber.StandardDeviation != nil {
			stdDev = *component.SingleNumber.StandardDeviation
		}
		return *result, stdDev, nil
	} else if component.Range != nil {
		result, err := SimulateIndependentRange(component.Range, timeFrame)
		if err != nil {
			return 0, 0, err
		}
		// Compute the average standard deviation of the range's min and max if they exist
		stdDev := 0.0
		if component.Range.Minimum.StandardDeviation != nil && component.Range.Maximum.StandardDeviation != nil {
			stdDev = (*component.Range.Minimum.StandardDeviation + *component.Range.Maximum.StandardDeviation) / 2
		}
		return *result, stdDev, nil
	} else if component.Decomposed != nil {
		probComposite, probStdDev, impactComposite, impactStdDev, costComposite, costStdDev, err := SimulateIndependentDecomposed(component.Decomposed)
		if err != nil {
			return 0, 0, err
		}
		// Select the appropriate attribute's composite value and standard deviation
		switch attribute {
		case ProbabilityAttribute:
			return probComposite, probStdDev, nil
		case ImpactAttribute:
			return impactComposite, impactStdDev, nil
		case CostAttribute:
			return costComposite, costStdDev, nil
		default:
			return 0, 0, errors.New("invalid attribute specified")
		}
	}

	return 0, 0, errors.New("invalid component")
}

func handleComponent(comp *Component, timeFrame uint64) (float64, float64, error) {
	if comp.SingleNumber != nil {
		value, stdDev, err := simulateSingleNumber(comp.SingleNumber, timeFrame)
		if err != nil {
			return 0, 0, err
		}
		return value, stdDev, nil
	} else if comp.Range != nil {
		value, stdDev, err := simulateRange(comp.Range, timeFrame)
		if err != nil {
			return 0, 0, err
		}
		return value, stdDev, nil
	} else if comp.Decomposed != nil {
		pVal, pStdDev, _, _, _, _, err := SimulateIndependentDecomposed(comp.Decomposed)
		if err != nil {
			return 0, 0, err
		}
		return pVal, pStdDev, nil
	}
	return 0, 0, errors.New("invalid component")
}

func simulateSingleNumber(sn *types.SingleNumber, timeFrame uint64) (float64, float64, error) {
	value, err := SimulateIndependentSingleNumer(sn, timeFrame)
	if err != nil {
		return 0, 0, err
	}
	return *value, *sn.StandardDeviation, nil
}

func simulateRange(rng *types.Range, timeFrame uint64) (float64, float64, error) {
	value, err := SimulateIndependentRange(rng, timeFrame)
	if err != nil {
		return 0, 0, err
	}
	stdDev := (*rng.Minimum.StandardDeviation + *rng.Maximum.StandardDeviation) / 2
	return *value, stdDev, nil
}

func adjustTimeFrameForDecomposed(decomposed *types.Decomposed, timeFrame uint64) error {
	for _, component := range decomposed.Components {
		// Adjust Probability
		if component.Probability != nil && component.Probability.Decomposed != nil {
			err := adjustTimeFrameForDecomposed(component.Probability.Decomposed, timeFrame)
			if err != nil {
				return err
			}
		}

		// Adjust Impact
		if component.Impact != nil && component.Impact.Decomposed != nil {
			err := adjustTimeFrameForDecomposed(component.Impact.Decomposed, timeFrame)
			if err != nil {
				return err
			}
		}

		// Adjust Cost
		if component.Cost != nil && component.Cost.Decomposed != nil {
			err := adjustTimeFrameForDecomposed(component.Cost.Decomposed, timeFrame)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SimulateIndependentEvent(event *utils.FilteredEvent) (*types.SimulationResults, error) {
	results := &types.SimulationResults{
		EventID:        event.ID,
		EventTimeFrame: event.Event.Timeframe,
	}

	if event.Event.AssociatedProbability != nil {
		if event.Event.AssociatedProbability.SingleNumber != nil {
			probability, err := SimulateIndependentSingleNumer(event.Event.AssociatedProbability.SingleNumber, event.Event.Timeframe)
			if err != nil {
				return nil, err
			}
			results.Probability = *probability
		}

		if event.Event.AssociatedProbability.Range != nil {
			probability, err := SimulateIndependentRange(event.Event.AssociatedProbability.Range, event.Event.Timeframe)
			if err != nil {
				return nil, err
			}
			results.Probability = *probability
		}

		if event.Event.AssociatedProbability.Decomposed != nil {
			probability, probabilityStdDev, _, _, _, _, err := SimulateIndependentDecomposed(event.Event.AssociatedProbability.Decomposed)
			if err != nil {
				return nil, err
			}
			results.Probability = probability
			results.ProbabilityStandardDeviation = probabilityStdDev
		}

	}

	if event.Event.AssociatedImpact != nil {
		if event.Event.AssociatedImpact.SingleNumber != nil {
			impact, err := SimulateIndependentSingleNumer(event.Event.AssociatedImpact.SingleNumber, event.Event.Timeframe)
			if err != nil {
				return nil, err
			}
			results.Impact = *impact
		}

		if event.Event.AssociatedImpact.Range != nil {
			impact, err := SimulateIndependentRange(event.Event.AssociatedImpact.Range, event.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			results.Impact = *impact

		}

		if event.Event.AssociatedImpact.Decomposed != nil {
			_, _, impact, impactStdDev, _, _, err := SimulateIndependentDecomposed(event.Event.AssociatedImpact.Decomposed)
			if err != nil {
				return nil, err
			}
			results.Impact = impact
			results.ImpactStandardDeviation = impactStdDev
		}

	}

	if event.Event.AssociatedCost != nil {

		if event.Event.AssociatedCost.SingleNumber != nil {
			cost, err := SimulateIndependentSingleNumer(event.Event.AssociatedCost.SingleNumber, event.Event.Timeframe)
			if err != nil {
				return nil, err
			}
			results.Cost = *cost
		}

		if event.Event.AssociatedCost.Range != nil {
			cost, err := SimulateIndependentRange(event.Event.AssociatedCost.Range, event.Event.Timeframe)
			if err != nil {
				return nil, err
			}
			results.Cost = *cost
		}

		if event.Event.AssociatedCost.Decomposed != nil {
			_, _, _, _, cost, costStdDev, err := SimulateIndependentDecomposed(event.Event.AssociatedCost.Decomposed)
			if err != nil {
				return nil, err
			}
			results.Cost = cost
			results.CostStandardDeviation = costStdDev
		}

	}

	if event.Event.AssociatedRisk != nil {
		if event.Event.AssociatedRisk.Probability != nil {
			if event.Event.AssociatedRisk.Probability.SingleNumber != nil {
				probability, err := SimulateIndependentSingleNumer(event.Event.AssociatedRisk.Probability.SingleNumber, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				results.AssociatedRisk.Probability = *probability
			}
			if event.Event.AssociatedRisk.Probability.Range != nil {
				probability, err := SimulateIndependentRange(event.Event.AssociatedRisk.Probability.Range, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}
				results.AssociatedRisk.Probability = *probability

			}
			if event.Event.AssociatedRisk.Probability.Decomposed != nil {
				probability, probabilityStdDev, _, _, _, _, err := SimulateIndependentDecomposed(event.Event.AssociatedRisk.Probability.Decomposed)
				if err != nil {
					return nil, err
				}

				results.AssociatedRisk.Probability = probability
				results.AssociatedRisk.ProbabilityStandardDeviation = probabilityStdDev
			}
		}

		if event.Event.AssociatedRisk.Impact != nil {
			if event.Event.AssociatedRisk.Impact.SingleNumber != nil {
				impact, err := SimulateIndependentSingleNumer(event.Event.AssociatedRisk.Impact.SingleNumber, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}
				results.AssociatedRisk.Impact = *impact
			}
			if event.Event.AssociatedRisk.Impact.Range != nil {
				impact, err := SimulateIndependentRange(event.Event.AssociatedRisk.Impact.Range, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				results.AssociatedRisk.Impact = *impact

			}
			if event.Event.AssociatedRisk.Impact.Decomposed != nil {
				_, _, impact, impactStdDev, _, _, err := SimulateIndependentDecomposed(event.Event.AssociatedRisk.Impact.Decomposed)
				if err != nil {
					return nil, err
				}
				results.AssociatedRisk.Impact = impact
				results.AssociatedRisk.ImpactStandardDeviation = impactStdDev
			}
		}
	}

	if event.Event.AssociatedMitigation != nil {
		if event.Event.AssociatedMitigation.Probability != nil {
			if event.Event.AssociatedMitigation.Probability.SingleNumber != nil {
				probability, err := SimulateIndependentSingleNumer(event.Event.AssociatedMitigation.Probability.SingleNumber, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}
				results.AssociatedMitigation.Probability = *probability
			}
			if event.Event.AssociatedMitigation.Probability.Range != nil {
				probability, err := SimulateIndependentRange(event.Event.AssociatedMitigation.Probability.Range, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				results.AssociatedMitigation.Probability = *probability
			}
			if event.Event.AssociatedMitigation.Probability.Decomposed != nil {
				probability, probabilityStdDev, _, _, _, _, err := SimulateIndependentDecomposed(event.Event.AssociatedMitigation.Probability.Decomposed)
				if err != nil {
					return nil, err
				}

				results.AssociatedMitigation.Probability = probability
				results.AssociatedMitigation.ProbabilityStandardDeviation = probabilityStdDev

			}
		}
		if event.Event.AssociatedMitigation.Impact != nil {
			if event.Event.AssociatedMitigation.Impact.SingleNumber != nil {
				imapct, err := SimulateIndependentSingleNumer(event.Event.AssociatedMitigation.Impact.SingleNumber, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}
				results.AssociatedMitigation.Impact = *imapct
			}
			if event.Event.AssociatedMitigation.Impact.Range != nil {
				impact, err := SimulateIndependentRange(event.Event.AssociatedMitigation.Impact.Range, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				results.AssociatedMitigation.Impact = *impact
			}
			if event.Event.AssociatedMitigation.Impact.Decomposed != nil {
				_, _, impact, impactStdDev, _, _, err := SimulateIndependentDecomposed(event.Event.AssociatedMitigation.Impact.Decomposed)
				if err != nil {
					return nil, err
				}
				results.AssociatedMitigation.Impact = impact
				results.AssociatedMitigation.ImpactStandardDeviation = impactStdDev
			}

		}
		if event.Event.AssociatedMitigation.AssociatedCost != nil {
			if event.Event.AssociatedMitigation.AssociatedCost.SingleNumber != nil {
				cost, err := SimulateIndependentSingleNumer(event.Event.AssociatedMitigation.AssociatedCost.SingleNumber, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				results.AssociatedMitigation.Cost = *cost
			}
			if event.Event.AssociatedMitigation.AssociatedCost.Range != nil {
				cost, err := SimulateIndependentRange(event.Event.AssociatedMitigation.AssociatedCost.Range, event.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				results.AssociatedMitigation.Cost = *cost
			}
			if event.Event.AssociatedMitigation.AssociatedCost.Decomposed != nil {
				_, _, _, _, cost, costStdDev, err := SimulateIndependentDecomposed(event.Event.AssociatedMitigation.AssociatedCost.Decomposed)
				if err != nil {
					return nil, err
				}
				results.AssociatedMitigation.Cost = cost
				results.AssociatedMitigation.CostStandardDeviation = costStdDev
			}
		}
	}

	return results, nil
}

func CheckRiskDependencies(risk *types.Risk, events []*utils.FilteredEvent, risks []*types.Risk, mitigations []*types.Mitigation, dependenciesType uint64) (bool, error) {

	for _, riskDep := range risk.DependsOnRisk {
		if riskDep.DependentRiskID == nil {
			return false, errors.New("dependent risk ID is nil")
		}

		if *riskDep.DependentRiskID == risk.ID {
			return false, errors.New("risk cannot depend on itself")
		}

		depRisk, err := utils.FindRiskByID(*riskDep.DependentRiskID, risks)
		if err != nil {
			return false, fmt.Errorf("error finding dependent risk %d: %s", *riskDep.DependentRiskID, err)
		}
		if depRisk == nil {
			return false, fmt.Errorf("dependent risk %d not found", *riskDep.DependentRiskID)
		}

		hit, err := CheckRiskDependencies(depRisk, events, risks, mitigations, riskDep.Type)
		if err != nil || !hit {
			return false, err
		}
	}

	for _, eventDep := range risk.DependsOnEvent {
		depEvent, err := utils.FindEventByID(eventDep.DependentEventID, events)
		if err != nil {
			return false, fmt.Errorf("error finding dependent event %d: %s", eventDep.DependentEventID, err)
		}
		if depEvent == nil {
			return false, fmt.Errorf("dependent event %d not found", eventDep.DependentEventID)
		}

		hit, err := DependencyCheck(depEvent, dependenciesType, events, risks, mitigations)
		if err != nil || !hit {
			return false, err
		}
	}

	switch dependenciesType {
	case types.Exists:
		// check risk has non-zero probability and impact
		if risk.Probability.SingleNumber != nil {
			base, err := SimulateIndependentSingleNumer(risk.Probability.SingleNumber, risk.TimeFrame)
			if err != nil {
				return false, err
			}

			if *base <= 0 {
				return false, fmt.Errorf("risk %d has zero probability", risk.ID)
			}

		} else if risk.Probability.Range != nil {
			base, err := SimulateIndependentRange(risk.Probability.Range, risk.TimeFrame)
			if err != nil {
				return false, err
			}

			if *base <= 0 {
				return false, fmt.Errorf("risk %d has zero probability", risk.ID)
			}

		} else if risk.Probability.Decomposed != nil {
			base, _, _, _, _, _, err := SimulateIndependentDecomposed(risk.Probability.Decomposed)
			if err != nil {
				return false, err
			}

			if base <= 0 {
				return false, fmt.Errorf("risk %d has zero probability", risk.ID)
			}

		} else {
			return false, fmt.Errorf("risk %d has no probability", risk.ID)
		}
	case types.DoesNotExist:
		// check risk has zero probability and impact
		if risk.Probability.SingleNumber != nil {
			base, err := SimulateIndependentSingleNumer(risk.Probability.SingleNumber, risk.TimeFrame)
			if err != nil {
				return false, err
			}

			if *base > 0 {
				return false, fmt.Errorf("risk %d has non-zero probability", risk.ID)
			}

		} else if risk.Probability.Range != nil {
			base, err := SimulateIndependentRange(risk.Probability.Range, risk.TimeFrame)
			if err != nil {
				return false, err
			}

			if *base > 0 {
				return false, fmt.Errorf("risk %d has non-zero probability", risk.ID)
			}

		} else if risk.Probability.Decomposed != nil {
			base, _, _, _, _, _, err := SimulateIndependentDecomposed(risk.Probability.Decomposed)
			if err != nil {
				return false, err
			}

			if base > 0 {
				return false, fmt.Errorf("risk %d has non-zero probability", risk.ID)
			}

		} else {
			return false, fmt.Errorf("risk %d has no probability", risk.ID)
		}
	default:
		return false, fmt.Errorf("invalid dependency type %d", dependenciesType)
	}

	return true, nil
}

// CheckMitigationDependencies checks if a mitigation's dependencies on risks, events, or other mitigations are met.
func CheckMitigationDependencies(mitigation *types.Mitigation, events []*utils.FilteredEvent, risks []*types.Risk, mitigations []*types.Mitigation, dependenciesType uint64) (bool, error) {
	// Check dependencies on events
	for _, eventDep := range mitigation.DependsOnEvent {
		depEvent, err := utils.FindEventByID(eventDep.DependentEventID, events)
		if err != nil {
			return false, fmt.Errorf("error finding dependent event %d: %s", eventDep.DependentEventID, err)
		}
		if depEvent == nil {
			return false, fmt.Errorf("dependent event %d not found", eventDep.DependentEventID)
		}

		hit, err := DependencyCheck(depEvent, dependenciesType, events, risks, mitigations)
		if err != nil || !hit {
			return false, err
		}
	}

	// Check dependencies on risks
	for _, riskDep := range mitigation.DependsOnRisk {
		if riskDep.DependentRiskID == nil {
			return false, errors.New("dependent risk ID is nil")
		}

		depRisk, err := utils.FindRiskByID(*riskDep.DependentRiskID, risks)
		if err != nil {
			return false, fmt.Errorf("error finding dependent risk %d: %s", *riskDep.DependentRiskID, err)
		}
		if depRisk == nil {
			return false, fmt.Errorf("dependent risk %d not found", *riskDep.DependentRiskID)
		}

		// Assuming risks can be evaluated similarly to events
		hit, err := CheckRiskDependencies(depRisk, events, risks, mitigations, riskDep.Type)
		if err != nil || !hit {
			return false, err
		}
	}

	// Check dependencies on other mitigations (if applicable)
	for _, mitDep := range mitigation.DependsOnMitigation {
		if mitDep.DependentMitigationOrRiskID == nil {
			return false, errors.New("dependent mitigation ID is nil")
		}

		depMitigation, err := utils.FindMitigationByID(*mitDep.DependentMitigationOrRiskID, mitigations)
		if err != nil {
			return false, fmt.Errorf("error finding dependent mitigation %d: %s", *mitDep.DependentMitigationOrRiskID, err)
		}
		if depMitigation == nil {
			return false, fmt.Errorf("dependent mitigation %d not found", *mitDep.DependentMitigationOrRiskID)
		}

		// Recursively check the dependencies of the dependent mitigation
		hit, err := CheckMitigationDependencies(depMitigation, events, risks, mitigations, mitDep.Type)
		if err != nil || !hit {
			return false, err
		}
	}

	switch dependenciesType {
	case types.Exists:
		// check mitigation has non-zero probability and impact
		if mitigation.Probability.SingleNumber != nil {
			base, err := SimulateIndependentSingleNumer(mitigation.Probability.SingleNumber, mitigation.TimeFrame)
			if err != nil {
				return false, err
			}

			if *base <= 0 {
				return false, fmt.Errorf("mitigation %d has zero probability", mitigation.ID)
			}

		} else if mitigation.Probability.Range != nil {
			base, err := SimulateIndependentRange(mitigation.Probability.Range, mitigation.TimeFrame)
			if err != nil {
				return false, err
			}

			if *base <= 0 {
				return false, fmt.Errorf("mitigation %d has zero probability", mitigation.ID)
			}

		} else if mitigation.Probability.Decomposed != nil {
			base, _, _, _, _, _, err := SimulateIndependentDecomposed(mitigation.Probability.Decomposed)
			if err != nil {
				return false, err
			}

			if base <= 0 {
				return false, fmt.Errorf("mitigation %d has zero probability", mitigation.ID)
			}

		} else {
			return false, fmt.Errorf("mitigation %d has no probability", mitigation.ID)
		}

		break
	case types.DoesNotExist:
		// check mitigation has zero probability
		if mitigation.Probability.SingleNumber != nil {
			base, err := SimulateIndependentSingleNumer(mitigation.Probability.SingleNumber, mitigation.TimeFrame)
			if err != nil {
				return false, err
			}

			if *base > 0 {
				return false, fmt.Errorf("mitigation %d has non-zero probability", mitigation.ID)
			}

		} else if mitigation.Probability.Range != nil {
			base, err := SimulateIndependentRange(mitigation.Probability.Range, mitigation.TimeFrame)
			if err != nil {
				return false, err
			}

			if *base > 0 {
				return false, fmt.Errorf("mitigation %d has non-zero probability", mitigation.ID)
			}

		} else if mitigation.Probability.Decomposed != nil {
			base, _, _, _, _, _, err := SimulateIndependentDecomposed(mitigation.Probability.Decomposed)
			if err != nil {
				return false, err
			}

			if base > 0 {
				return false, fmt.Errorf("mitigation %d has non-zero probability", mitigation.ID)
			}

		} else {
			return false, fmt.Errorf("mitigation %d has no probability", mitigation.ID)
		}
		break
	default:
		return false, fmt.Errorf("invalid dependency type %d", dependenciesType)
	}

	return true, nil
}
