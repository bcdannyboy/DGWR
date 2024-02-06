package simulator

import (
	"errors"
	"fmt"

	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

func DependencyCheck(
	DepEvent *utils.FilteredEvent,
	DType uint64,
	Events []*utils.FilteredEvent,
	ExpectedValue *types.SingleNumber,
	ExpectedRange *types.Range,
	ExpectedDecomp *types.Decomposed) (bool, error) {

	// 1. check if the dependent event is valid
	if DepEvent == nil || DepEvent.Event == nil {
		return false, fmt.Errorf("dependent event is nil")
	}

	// 2. process the dependency type
	for _, DoE := range DepEvent.Event.DependsOnEvent {
		DoEvent, err := utils.FindEventByID(DoE.DependentEventID, Events)
		if err != nil {
			return false, fmt.Errorf("dependent event %d not found", DoE.DependentEventID)
		}

		// Process Depends on Event
		switch DType {
		case types.Happens:
			// happens means a non-zero probability for the event or its associated risk / mitigation if no event probability is provided
			if DepEvent.Event.AssociatedProbability == nil {
				if DepEvent.Event.AssociatedRisk != nil {
					if DepEvent.Event.AssociatedRisk.Probability == nil {
						return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
					}

					if DepEvent.Event.AssociatedRisk.Probability.SingleNumber != nil {
						base, std, err := simulateSingleNumber(DepEvent.Event.AssociatedRisk.Probability.SingleNumber, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						min := base - std
						if min <= 0 {
							return false, nil // missed dependency
						}
					} else if DepEvent.Event.AssociatedRisk.Probability.Range != nil {
						base, std, err := simulateRange(DepEvent.Event.AssociatedRisk.Probability.Range, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating range for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						min := base - std
						if min <= 0 {
							return false, nil // missed dependency
						}
					} else if DepEvent.Event.AssociatedRisk.Probability.Decomposed != nil {
						base, std, err := simulateDecomposedByAttribute(DepEvent.Event.AssociatedRisk.Probability.Decomposed, ProbabilityAttribute, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						min := base - std
						if min <= 0 {
							return false, nil // missed dependency
						}
					} else {
						return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
					}

				} else if DepEvent.Event.AssociatedMitigation != nil {
					if DepEvent.Event.AssociatedMitigation.Probability == nil {
						return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
					}

					if DepEvent.Event.AssociatedMitigation.Probability.SingleNumber != nil {
						base, std, err := simulateSingleNumber(DepEvent.Event.AssociatedMitigation.Probability.SingleNumber, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						min := base - std
						if min <= 0 {
							return false, nil // missed dependency
						}

					} else if DepEvent.Event.AssociatedMitigation.Probability.Range != nil {

						base, std, err := simulateRange(DepEvent.Event.AssociatedMitigation.Probability.Range, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating range for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						min := base - std
						if min <= 0 {
							return false, nil // missed dependency

						}

					} else if DepEvent.Event.AssociatedMitigation.Probability.Decomposed != nil {

						base, std, err := simulateDecomposedByAttribute(DepEvent.Event.AssociatedMitigation.Probability.Decomposed, ProbabilityAttribute, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						min := base - std
						if min <= 0 {
							return false, nil // missed dependency
						}

					} else {
						return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
					}
				} else {
					return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
				}
			} else {
				if DoEvent.Event.AssociatedProbability == nil {
					return false, fmt.Errorf("dependent event %d has no probability to compare with %d", DoEvent.Event.ID, DepEvent.Event.ID)
				}

				if DoEvent.Event.AssociatedProbability.SingleNumber != nil {
					base, std, err := simulateSingleNumber(DoEvent.Event.AssociatedProbability.SingleNumber, DoEvent.Event.Timeframe)
					if err != nil {
						return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoEvent.Event.ID, err.Error())
					}

					min := base - std
					if min <= 0 {
						return false, nil // missed dependency
					}
				} else if DoEvent.Event.AssociatedProbability.Range != nil {
					base, std, err := simulateRange(DoEvent.Event.AssociatedProbability.Range, DoEvent.Event.Timeframe)
					if err != nil {
						return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoEvent.Event.ID, err.Error())
					}

					min := base - std
					if min <= 0 {
						return false, nil // missed dependency
					}
				} else if DoEvent.Event.AssociatedProbability.Decomposed != nil {
					base, std, err := simulateDecomposedByAttribute(DoEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoEvent.Event.Timeframe)
					if err != nil {
						return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoEvent.Event.ID, err.Error())
					}

					min := base - std
					if min <= 0 {
						return false, nil // missed dependency
					}
				} else {
					return false, fmt.Errorf("dependent event %d has no probability to compare with %d", DoEvent.Event.ID, DepEvent.Event.ID)
				}
			}
			break
		case types.DoesNotHappen:
			// does not happen means a zero probability for the event or its associated risk / mitigation if no event probability is provided
			if DepEvent.Event.AssociatedProbability == nil {
				if DepEvent.Event.AssociatedRisk != nil {
					if DepEvent.Event.AssociatedRisk.Probability == nil {
						return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
					}

					if DepEvent.Event.AssociatedRisk.Probability.SingleNumber != nil {
						base, std, err := simulateSingleNumber(DepEvent.Event.AssociatedRisk.Probability.SingleNumber, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						if base-std > 0 {
							return false, nil // missed dependency
						}
					} else if DepEvent.Event.AssociatedRisk.Probability.Range != nil {
						base, std, err := simulateRange(DepEvent.Event.AssociatedRisk.Probability.Range, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating range for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						if base-std > 0 {
							return false, nil // missed dependency
						}
					} else if DepEvent.Event.AssociatedRisk.Probability.Decomposed != nil {
						base, std, err := simulateDecomposedByAttribute(DepEvent.Event.AssociatedRisk.Probability.Decomposed, ProbabilityAttribute, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						if base-std > 0 {
							return false, nil // missed dependency
						}
					} else {
						return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
					}

				} else if DepEvent.Event.AssociatedMitigation != nil {
					if DepEvent.Event.AssociatedMitigation.Probability == nil {
						return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
					}

					if DepEvent.Event.AssociatedMitigation.Probability.SingleNumber != nil {
						base, std, err := simulateSingleNumber(DepEvent.Event.AssociatedMitigation.Probability.SingleNumber, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						if base-std > 0 {
							return false, nil // missed dependency
						}

					} else if DepEvent.Event.AssociatedMitigation.Probability.Range != nil {

						base, std, err := simulateRange(DepEvent.Event.AssociatedMitigation.Probability.Range, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating range for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						if base-std > 0 {
							return false, nil // missed dependency
						}

					} else if DepEvent.Event.AssociatedMitigation.Probability.Decomposed != nil {

						base, std, err := simulateDecomposedByAttribute(DepEvent.Event.AssociatedMitigation.Probability.Decomposed, ProbabilityAttribute, DepEvent.Event.Timeframe)
						if err != nil {
							return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DepEvent.Event.ID, err.Error())
						}

						if base-std > 0 {
							return false, nil // missed dependency
						}

					} else {
						return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
					}
				} else {
					return false, fmt.Errorf("dependent event has no probability, risk, or mitigation")
				}
			} else {
				if DoEvent.Event.AssociatedProbability == nil {
					return false, fmt.Errorf("dependent event %d has no probability to compare with %d", DoEvent.Event.ID, DepEvent.Event.ID)
				}
			}
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}

	}

	for _, DoP := range DepEvent.Event.DependsOnProbability {
		// Process Depends on Probability
		switch DType {
		case types.Has:
			break
		case types.HasNot:
			break
		case types.In:
			break
		case types.Out:
			break
		case types.EQ:
			break
		case types.NEQ:
			break
		case types.LT:
			break
		case types.GT:
			break
		case types.LTE:
			break
		case types.GTE:
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	for _, DoI := range DepEvent.Event.DependsOnImpact {
		// Process Depends on Impact

		switch DType {
		case types.Has:
			break
		case types.HasNot:
			break
		case types.In:
			break
		case types.Out:
			break
		case types.EQ:
			break
		case types.NEQ:
			break
		case types.LT:
			break
		case types.GT:
			break
		case types.LTE:
			break
		case types.GTE:
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	for _, DoC := range DepEvent.Event.DependsOnCost {
		// Process Depends on Cost

		switch DType {
		case types.Has:
			break
		case types.HasNot:
			break
		case types.In:
			break
		case types.Out:
			break
		case types.EQ:
			break
		case types.NEQ:
			break
		case types.LT:
			break
		case types.GT:
			break
		case types.LTE:
			break
		case types.GTE:
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	for _, DoR := range DepEvent.Event.DependsOnRisk {
		// Process Depends on Risk
		switch DType {
		case types.Exists:
			break
		case types.DoesNotExist:
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	for _, DoM := range DepEvent.Event.DependsOnMitigation {
		// Process Depends on Mitigation
		switch DType {
		case types.Exists:
			break
		case types.DoesNotExist:
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	return true, nil
}

func simulateDecomposedByAttribute(decomposed *types.Decomposed, attribute int, timeFrame uint64) (float64, float64, error) {
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
