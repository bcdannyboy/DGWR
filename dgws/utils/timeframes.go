package utils

import (
	"errors"

	"github.com/bcdannyboy/montecargo/dgws/types"
)

func CheckForValidTimeframe(timeframe uint64) bool {
	switch timeframe {
	case types.Year:
		return true
	case types.Quarter:
		return true
	case types.Month:
		return true
	case types.Week:
		return true
	case types.Day:
		return true
	case types.TwoYears:
		return true
	case types.FiveYears:
		return true
	case types.TenYears:
		return true
	default:
		return false
	}
}

func AdjustProbabilityForTimeFrame(probability float64, timeframe uint64) (float64, error) {
	// adjust the probability based on the timeframe of the event to standardize at a yearly rate
	switch timeframe {
	case types.Year:
		return probability, nil
	case types.Quarter:
		return probability * 4, nil
	case types.Month:
		return probability * 12, nil
	case types.Week:
		return probability * 52, nil
	case types.Day:
		return probability * 365, nil
	case types.TwoYears:
		return probability / 2, nil
	case types.FiveYears:
		return probability / 5, nil
	case types.TenYears:
		return probability / 10, nil
	default:
		return 0, errors.New("invalid timeframe")
	}
}
