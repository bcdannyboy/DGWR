package simulator

import (
	"errors"

	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

func SimulateIndependentSingleNumer(sn *types.SingleNumber, timeFrame uint64) (*float64, error) {
	if utils.CheckForValidTimeframe(timeFrame) {
		adjustedProbability, err := utils.AdjustProbabilityForTimeFrame(sn.Value, timeFrame)
		if err != nil {
			return nil, err
		}
		sn.Value = adjustedProbability
	} else {
		return nil, errors.New("invalid timeframe")
	}

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

func SimulateIndependentRange(rng *types.Range, timeFrame uint64) (*float64, error) {
	if utils.CheckForValidTimeframe(timeFrame) {
		adjustedMinimum, err := utils.AdjustProbabilityForTimeFrame(rng.Minimum.Value, timeFrame)
		if err != nil {
			return nil, err
		}
		rng.Minimum.Value = adjustedMinimum

		adjustedMaximum, err := utils.AdjustProbabilityForTimeFrame(rng.Maximum.Value, timeFrame)
		if err != nil {
			return nil, err
		}
		rng.Maximum.Value = adjustedMaximum
	} else {
		return nil, errors.New("invalid timeframe")
	}

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
			if utils.CheckForValidTimeframe(component.TimeFrame) {
				if component.Probability.SingleNumber != nil {
					adjustedProbability, err := utils.AdjustProbabilityForTimeFrame(component.Probability.SingleNumber.Value, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}

					component.Probability.SingleNumber.Value = adjustedProbability
				} else if component.Probability.Range != nil {
					adjustedMinimum, err := utils.AdjustProbabilityForTimeFrame(component.Probability.Range.Minimum.Value, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}

					component.Probability.Range.Minimum.Value = adjustedMinimum
				} else {
					err := adjustTimeFrameForDecomposed(component.Probability.Decomposed, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}

				}
			} else {
				return 0, 0, 0, 0, 0, 0, errors.New("invalid timeframe")
			}

			pVal, pStdDev, err := handleComponent(&Component{
				SingleNumber: component.Probability.SingleNumber,
				Range:        component.Probability.Range,
				Decomposed:   component.Probability.Decomposed,
			}, component.TimeFrame)
			if err != nil {
				return 0, 0, 0, 0, 0, 0, err
			}
			probValues = append(probValues, pVal)
			probStdDevs = append(probStdDevs, pStdDev)
		}

		// Handle Impact
		if component.Impact != nil {
			if utils.CheckForValidTimeframe(component.TimeFrame) {
				if component.Impact.SingleNumber != nil {
					adjustedImpact, err := utils.AdjustProbabilityForTimeFrame(component.Impact.SingleNumber.Value, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}
					component.Impact.SingleNumber.Value = adjustedImpact
				} else if component.Impact.Range != nil {
					adjustedMinimum, err := utils.AdjustProbabilityForTimeFrame(component.Impact.Range.Minimum.Value, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}
					component.Impact.Range.Minimum.Value = adjustedMinimum

					adjustedMaximum, err := utils.AdjustProbabilityForTimeFrame(component.Impact.Range.Maximum.Value, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}
					component.Impact.Range.Maximum.Value = adjustedMaximum
				} else {
					err := adjustTimeFrameForDecomposed(component.Impact.Decomposed, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}
				}
			} else {
				return 0, 0, 0, 0, 0, 0, errors.New("invalid timeframe")
			}

			iVal, iStdDev, err := handleComponent(&Component{
				SingleNumber: component.Impact.SingleNumber,
				Range:        component.Impact.Range,
				Decomposed:   component.Impact.Decomposed,
			}, component.TimeFrame)
			if err != nil {
				return 0, 0, 0, 0, 0, 0, err
			}
			impactValues = append(impactValues, iVal)
			impactStdDevs = append(impactStdDevs, iStdDev)
		}

		// Handle Cost
		if component.Cost != nil {
			if utils.CheckForValidTimeframe(component.TimeFrame) {
				if component.Cost.SingleNumber != nil {
					adjustedCost, err := utils.AdjustProbabilityForTimeFrame(component.Cost.SingleNumber.Value, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}
					component.Cost.SingleNumber.Value = adjustedCost
				} else if component.Cost.Range != nil {
					adjustedMinimum, err := utils.AdjustProbabilityForTimeFrame(component.Cost.Range.Minimum.Value, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}
					component.Cost.Range.Minimum.Value = adjustedMinimum

					adjustedMaximum, err := utils.AdjustProbabilityForTimeFrame(component.Cost.Range.Maximum.Value, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}
					component.Cost.Range.Maximum.Value = adjustedMaximum
				} else {
					err := adjustTimeFrameForDecomposed(component.Cost.Decomposed, component.TimeFrame)
					if err != nil {
						return 0, 0, 0, 0, 0, 0, err
					}
				}
			} else {
				return 0, 0, 0, 0, 0, 0, errors.New("invalid timeframe")
			}

			cVal, cStdDev, err := handleComponent(&Component{
				SingleNumber: component.Cost.SingleNumber,
				Range:        component.Cost.Range,
				Decomposed:   component.Cost.Decomposed,
			}, component.TimeFrame)
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
