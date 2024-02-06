package simulator

import (
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
						for _, comp := range DepEvent.Event.AssociatedRisk.Probability.Decomposed.Components {

						}
					}

				} else if DepEvent.Event.AssociatedMitigation != nil {

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
					for _, comp := range DoEvent.Event.AssociatedProbability.Decomposed.Components {

					}
				}
			}
			break
		case types.DoesNotHappen:
			// does not happen means a zero probability for the event or its associated risk / mitigation if no event probability is provided
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
