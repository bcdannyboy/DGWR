package simulator

import (
	"errors"

	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

func SimulateIndependentSingleNumer(sn *types.SingleNumber) (*float64, error) {
	// Initialize the random seed
	randomFloat, err := utils.CryptoRandFloat64()
	if err != nil {
		return nil, err
	}

	// probability is impacted by the user's confidence
	randomFloat = randomFloat - (randomFloat * *sn.Confidence)

	confidenceAffectOnValue, err := utils.CoinFlip()
	if err != nil {
		return nil, err
	}

	if confidenceAffectOnValue {
		sn.Value = sn.Value + randomFloat
	} else {
		sn.Value = sn.Value - randomFloat
		if sn.Value < 0 {
			sn.Value = 0
		}
	}

	if sn.StandardDeviation != nil {
		confidenceAffectOnStandardDeviation, err := utils.CoinFlip()
		if err != nil {
			return nil, err
		}

		if confidenceAffectOnStandardDeviation {
			sn.StandardDeviation = utils.Float64toPointer(*sn.StandardDeviation + randomFloat)
		} else {
			sn.StandardDeviation = utils.Float64toPointer(*sn.StandardDeviation - randomFloat)
			if *sn.StandardDeviation < 0 {
				sn.StandardDeviation = utils.Float64toPointer(0)
			}
		}

		return utils.Float64toPointer(utils.LogNormalSample(sn.Value, *sn.StandardDeviation)), nil
	}

	return utils.Float64toPointer((utils.LogNormalSample(sn.Value, 0))), nil
}

func SimulateIndependentRange(rng *types.Range) (*float64, error) {
	randomMinimumFloat, err := utils.CryptoRandFloat64()
	if err != nil {
		return nil, err
	}
	randomMaximumFloat, err := utils.CryptoRandFloat64()
	if err != nil {
		return nil, err
	}

	// probability is impacted by the user's confidence
	randomMinimumFloat = randomMinimumFloat - (randomMinimumFloat * *rng.Minimum.Confidence)
	randomMaximumFloat = randomMaximumFloat - (randomMaximumFloat * *rng.Maximum.Confidence)

	minimumStandardDeviation := 0.0
	maximumStandardDeviation := 0.0

	confidenceAffectOnMinimumValue, err := utils.CoinFlip()
	if err != nil {
		return nil, err
	}

	if confidenceAffectOnMinimumValue {
		rng.Minimum.Value = rng.Minimum.Value + randomMinimumFloat
	} else {
		rng.Minimum.Value = rng.Minimum.Value - randomMinimumFloat
		if rng.Minimum.Value < 0 {
			rng.Minimum.Value = 0
		}
	}

	confidenceAffectOnMaximumValue, err := utils.CoinFlip()
	if err != nil {
		return nil, err
	}

	if confidenceAffectOnMaximumValue {
		rng.Maximum.Value = rng.Maximum.Value + randomMaximumFloat
	} else {
		rng.Maximum.Value = rng.Maximum.Value - randomMaximumFloat
		if rng.Maximum.Value < 0 {
			rng.Maximum.Value = 0
		}
	}

	if rng.Minimum.StandardDeviation != nil {
		confidenceAffectOnMinimumStandardDeviation, err := utils.CoinFlip()
		if err != nil {
			return nil, err
		}

		if confidenceAffectOnMinimumStandardDeviation {
			rng.Minimum.StandardDeviation = utils.Float64toPointer(*rng.Minimum.StandardDeviation + randomMinimumFloat)
		} else {
			rng.Minimum.StandardDeviation = utils.Float64toPointer(*rng.Minimum.StandardDeviation - randomMinimumFloat)
			if *rng.Minimum.StandardDeviation < 0 {
				rng.Minimum.StandardDeviation = utils.Float64toPointer(0)
			}
		}

		minimumStandardDeviation = *rng.Minimum.StandardDeviation
	}

	if rng.Maximum.StandardDeviation != nil {
		confidenceAffectOnMaximumStandardDeviation, err := utils.CoinFlip()
		if err != nil {
			return nil, err
		}

		if confidenceAffectOnMaximumStandardDeviation {
			rng.Maximum.StandardDeviation = utils.Float64toPointer(*rng.Maximum.StandardDeviation + randomMaximumFloat)
		} else {
			rng.Maximum.StandardDeviation = utils.Float64toPointer(*rng.Maximum.StandardDeviation - randomMaximumFloat)
			if *rng.Maximum.StandardDeviation < 0 {
				rng.Maximum.StandardDeviation = utils.Float64toPointer(0)
			}
		}

		minimumStandardDeviation = *rng.Minimum.StandardDeviation
	}

	mean := (rng.Minimum.Value + rng.Maximum.Value) / 2

	return utils.Float64toPointer(utils.LogNormalSampleInRange(mean, minimumStandardDeviation, maximumStandardDeviation, rng.Minimum.Value, rng.Maximum.Value)), nil
}

func SimulateIndependentDecomposed(decomposed *types.Decomposed) (float64, float64, float64, float64, float64, float64, error) {
	if decomposed == nil {
		return 0, 0, 0, 0, 0, 0, errors.New("decomposed is nil")
	}

	var probValues, impactValues, costValues []float64
	var probStdDevs, impactStdDevs, costStdDevs []float64

	for _, component := range decomposed.Components {
		// Handle Probability
		if component.Probability != nil {
			pVal, pStdDev, err := handleComponent(&Component{
				SingleNumber: component.Probability.SingleNumber,
				Range:        component.Probability.Range,
				Decomposed:   component.Probability.Decomposed,
			})
			if err != nil {
				return 0, 0, 0, 0, 0, 0, err
			}
			probValues = append(probValues, pVal)
			probStdDevs = append(probStdDevs, pStdDev)
		}

		// Handle Impact
		if component.Impact != nil {
			iVal, iStdDev, err := handleComponent(&Component{
				SingleNumber: component.Impact.SingleNumber,
				Range:        component.Impact.Range,
				Decomposed:   component.Impact.Decomposed,
			})
			if err != nil {
				return 0, 0, 0, 0, 0, 0, err
			}
			impactValues = append(impactValues, iVal)
			impactStdDevs = append(impactStdDevs, iStdDev)
		}

		// Handle Cost
		if component.Cost != nil {
			cVal, cStdDev, err := handleComponent(&Component{
				SingleNumber: component.Cost.SingleNumber,
				Range:        component.Cost.Range,
				Decomposed:   component.Cost.Decomposed,
			})
			if err != nil {
				return 0, 0, 0, 0, 0, 0, err
			}
			costValues = append(costValues, cVal)
			costStdDevs = append(costStdDevs, cStdDev)
		}
	}

	probComposite, probStdDev := utils.ComputeCompositeLogNormal(probValues, probStdDevs)
	impactComposite, impactStdDev := utils.ComputeCompositeLogNormal(impactValues, impactStdDevs)
	costComposite, costStdDev := utils.ComputeCompositeLogNormal(costValues, costStdDevs)

	return probComposite, probStdDev, impactComposite, impactStdDev, costComposite, costStdDev, nil
}

func handleComponent(comp *Component) (float64, float64, error) {
	if comp.SingleNumber != nil {
		value, stdDev, err := simulateSingleNumber(comp.SingleNumber)
		if err != nil {
			return 0, 0, err
		}
		return value, stdDev, nil
	} else if comp.Range != nil {
		value, stdDev, err := simulateRange(comp.Range)
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

func simulateSingleNumber(sn *types.SingleNumber) (float64, float64, error) {
	value, err := SimulateIndependentSingleNumer(sn)
	if err != nil {
		return 0, 0, err
	}
	return *value, *sn.StandardDeviation, nil
}

func simulateRange(rng *types.Range) (float64, float64, error) {
	value, err := SimulateIndependentRange(rng)
	if err != nil {
		return 0, 0, err
	}
	stdDev := (*rng.Minimum.StandardDeviation + *rng.Maximum.StandardDeviation) / 2
	return *value, stdDev, nil
}
