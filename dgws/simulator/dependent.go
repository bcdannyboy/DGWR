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
		// Process Depends on Event
		switch DType {
		case types.Happens:
			break
		case types.DoesNotHappen:
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
