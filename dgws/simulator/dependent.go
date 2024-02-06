package simulator

import (
	"errors"

	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

func SimulateDependentEvent(event *types.Event, Events []*utils.FilteredEvent, Risks []*types.Risk, Mitigations []*types.Mitigation, IndependentResults []*types.SimulationResults) (*types.SimulationResults, error) {
	// 1. check what kind of dependencies the event has
	expectedDependencies := 0
	expectedDependencies += len(event.DependsOnCost)
	expectedDependencies += len(event.DependsOnImpact)
	expectedDependencies += len(event.DependsOnProbability)
	expectedDependencies += len(event.DependsOnRisk)
	expectedDependencies += len(event.DependsOnEvent)
	expectedDependencies += len(event.DependsOnMitigation)
	DependenciesMet := 0
	DependenciesMissed := 0

	for _, costDependency := range event.DependsOnCost {

		dei := costDependency.DependentEventID
		deiSingle := costDependency.SingleValue
		deiRange := costDependency.Range
		deiDecomposed := costDependency.Decomposed

		depEvent := &utils.FilteredEvent{}

		inResults := utils.EventInResults(*dei, IndependentResults)
		if inResults == nil {
			de, err := utils.FindEventByID(*dei, Events)
			if err != nil {
				de = nil
			}
			depEvent = de
		} else {
			depEvent.Event.AssociatedCost.SingleNumber.Value = inResults.Cost
			depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation = &inResults.CostStandardDeviation
			depEvent.Event.AssociatedCost.SingleNumber.Confidence = utils.Float64toPointer(0.95)
		}

		if depEvent == nil {
			return nil, errors.New("dependent event does not exist for cost dependency")
		}

		// 2. simulate the cost dependencies and store the results
		switch costDependency.Type {
		case types.Exists:
			// check if the event dependency exists
			// existence for a cost means the dependency event has a non-zero cost value after simulation
			if depEvent.Event.AssociatedCost == nil {
				DependenciesMissed++
				break
			}

			if depEvent.Event.AssociatedCost.SingleNumber != nil {
				val := depEvent.Event.AssociatedCost.SingleNumber.Value
				sd := *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
				conf := *depEvent.Event.AssociatedCost.SingleNumber.Confidence
				if conf == 0 {
					conf = 0.5
				}

				singleRange := &types.Range{
					Minimum: types.Minimum{
						Value:             val,
						StandardDeviation: &sd,
						Confidence:        &conf,
					},
					Maximum: types.Maximum{
						Value:             val,
						StandardDeviation: &sd,
						Confidence:        &conf,
					},
				}

				result, resultSD, err := simulateRange(singleRange, depEvent.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				if result+resultSD > 0 {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}

				break
			} else if depEvent.Event.AssociatedCost.Range != nil {
				result, resultSD, err := simulateRange(depEvent.Event.AssociatedCost.Range, depEvent.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				if result+resultSD > 0 {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}

				break
			} else {
				return nil, errors.New("dependent event does not have a cost value or range but type is Exists for cost dependency")
			}
		case types.DoesNotExist:
			// check if the event dependency does not exist
			// non-existence for a cost means the dependency event has a zero cost value after simulation
			if depEvent.Event.AssociatedCost == nil {
				DependenciesMet++
				break
			}

			if depEvent.Event.AssociatedCost.SingleNumber != nil {
				val := depEvent.Event.AssociatedCost.SingleNumber.Value
				sd := *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
				conf := *depEvent.Event.AssociatedCost.SingleNumber.Confidence
				if conf == 0 {
					conf = 0.5
				}

				singleRange := &types.Range{
					Minimum: types.Minimum{
						Value:             val,
						StandardDeviation: &sd,
						Confidence:        &conf,
					},

					Maximum: types.Maximum{
						Value:             val,
						StandardDeviation: &sd,
						Confidence:        &conf,
					},
				}

				result, resultSD, err := simulateRange(singleRange, depEvent.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				if result+resultSD <= 0 {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
			} else if depEvent.Event.AssociatedCost.Range != nil {
				result, resultSD, err := simulateRange(depEvent.Event.AssociatedCost.Range, depEvent.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				if result+resultSD <= 0 {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
				break
			} else {
				return nil, errors.New("dependent event does not have a cost value or range but type is DoesNotExist for cost dependency")
			}
		case types.In:
			// check if the event dependency is in a range
			if deiRange == nil {
				return nil, errors.New("cost dependency range is nil but type is In")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is In for cost dependency")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}
			min := dist - sd
			if min < 0 {
				min = 0
			}
			max := dist + sd

			if depEvent.Event.AssociatedCost.SingleNumber != nil {
				minNum := depEvent.Event.AssociatedCost.SingleNumber.Value - *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
				maxNum := depEvent.Event.AssociatedCost.SingleNumber.Value + *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
				if minNum < 0 {
					minNum = 0
				}

				minCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}
				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum - (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				}

				if minNum >= min && maxNum <= max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
				break
			} else if depEvent.Event.AssociatedCost.Range != nil {
				minLower := depEvent.Event.AssociatedCost.Range.Minimum.Value - *depEvent.Event.AssociatedCost.Range.Minimum.StandardDeviation
				minUpper := depEvent.Event.AssociatedCost.Range.Minimum.Value + *depEvent.Event.AssociatedCost.Range.Minimum.StandardDeviation
				maxLower := depEvent.Event.AssociatedCost.Range.Maximum.Value - *depEvent.Event.AssociatedCost.Range.Maximum.StandardDeviation
				maxUpper := depEvent.Event.AssociatedCost.Range.Maximum.Value + *depEvent.Event.AssociatedCost.Range.Maximum.StandardDeviation

				if minLower < 0 {
					minLower = 0
				}
				if maxLower < 0 {
					maxLower = 0
				}

				minLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}
				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()

					if err != nil {
						return nil, err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper - (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				}

				if minLower >= min && minUpper <= max && maxLower >= min && maxUpper <= max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
				break
			} else {
				return nil, errors.New("dependent event does not have a cost value or range but type is In for cost dependency")
			}
		case types.Out:
			// check if the event dependency is out of a range
			if deiRange == nil {
				return nil, errors.New("cost dependency range is nil but type is Out")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is Out for cost dependency")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}
			min := dist - sd
			if min < 0 {
				min = 0
			}
			max := dist + sd

			if depEvent.Event.AssociatedCost.SingleNumber != nil {
				minNum := depEvent.Event.AssociatedCost.SingleNumber.Value - *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
				maxNum := depEvent.Event.AssociatedCost.SingleNumber.Value + *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
				if minNum < 0 {
					minNum = 0
				}

				minCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum - (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				}

				if minNum < min && maxNum > max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
				break
			} else if depEvent.Event.AssociatedCost.Range != nil {
				minLower := depEvent.Event.AssociatedCost.Range.Minimum.Value - *depEvent.Event.AssociatedCost.Range.Minimum.StandardDeviation
				minUpper := depEvent.Event.AssociatedCost.Range.Minimum.Value + *depEvent.Event.AssociatedCost.Range.Minimum.StandardDeviation
				maxLower := depEvent.Event.AssociatedCost.Range.Maximum.Value - *depEvent.Event.AssociatedCost.Range.Maximum.StandardDeviation
				maxUpper := depEvent.Event.AssociatedCost.Range.Maximum.Value + *depEvent.Event.AssociatedCost.Range.Maximum.StandardDeviation

				if minLower < 0 {
					minLower = 0
				}
				if maxLower < 0 {
					maxLower = 0
				}

				minLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper - (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				}

				if minLower < min && minUpper > max && maxLower < min && maxUpper > max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
				break
			} else {
				return nil, errors.New("dependent event does not have a cost value or range but type is Out for cost dependency")
			}
		case types.Has:
			// check if the event decomposition has a component
			if deiDecomposed.Components == nil {
				return nil, errors.New("cost dependency decomposed components is nil but type is Has")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is Has for cost dependency")
			}
			if depEvent.Event.AssociatedCost.Decomposed == nil {
				return nil, errors.New("dependent event does not have a cost value but type is Has for cost dependency")
			}

			// check if the dependent event has a cost value
			found := -1
			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedCost.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedCost.Decomposed.Components[j].Name {
						found = j
						break
					}
				}
			}
			if found == -1 {
				return nil, errors.New("dependent event has a cost value but type is Has for cost dependency")
			}

			// check if the dependent event has a nonzero cost value
			if deiDecomposed.Components[found].Cost.SingleNumber != nil {
				result, err := SimulateIndependentSingleNumer(deiDecomposed.Components[found].Cost.SingleNumber, depEvent.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				if *result > 0 {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
				break
			} else if deiDecomposed.Components[found].Cost.Range != nil {
				result, resultSD, err := simulateRange(deiDecomposed.Components[found].Cost.Range, depEvent.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				if result+resultSD > 0 {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
				break
			} else {
				return nil, errors.New("dependent event does not have a cost value or range but type is Has for cost dependency")
			}
		case types.HasNot:
			// check if the event decomposition does not have a component
			if deiDecomposed.Components == nil {
				return nil, errors.New("cost dependency decomposed components is nil but type is HasNot")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is HasNot for cost dependency")
			}
			if depEvent.Event.AssociatedCost.Decomposed == nil {
				return nil, errors.New("dependent event does not have a cost value but type is HasNot for cost dependency")
			}

			// check if the dependent event has a cost value
			found := -1
			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedCost.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedCost.Decomposed.Components[j].Name {
						found = j
						break
					}
				}
			}
			if found != -1 {
				return nil, errors.New("dependent event has a cost value but type is HasNot for cost dependency")
			}

			// check if the dependent event has a zero cost value
			if deiDecomposed.Components[found].Cost.SingleNumber != nil {
				result, err := SimulateIndependentSingleNumer(deiDecomposed.Components[found].Cost.SingleNumber, depEvent.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				if *result <= 0 {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
				break
			} else if deiDecomposed.Components[found].Cost.Range != nil {
				result, resultSD, err := simulateRange(deiDecomposed.Components[found].Cost.Range, depEvent.Event.Timeframe)
				if err != nil {
					return nil, err
				}

				if result+resultSD <= 0 {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}

			} else {
				return nil, errors.New("dependent event does not have a cost value or range but type is HasNot for cost dependency")
			}

			break
		case types.EQ:
			// check if the event value is equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("cost dependency single value is nil but type is EQ")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is EQ for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return nil, errors.New("dependent event does not have a cost value but type is EQ for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin >= min && deiMax <= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		case types.NEQ:
			// check if the event value is not equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("cost dependency single value is nil but type is NEQ")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is NEQ for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return nil, errors.New("dependent event does not have a cost value but type is NEQ for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()

				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()

				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin < min && deiMax > max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}
			break
		case types.LT:
			// check if the event value is less than the dependency value
			if deiSingle == nil {
				return nil, errors.New("cost dependency single value is nil but type is LT")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is LT for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return nil, errors.New("dependent event does not have a cost value but type is LT for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin < min && deiMax < max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}
			break
		case types.GT:
			// check if the event value is greater than the dependency value
			if deiSingle == nil {
				return nil, errors.New("cost dependency single value is nil but type is GT")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is GT for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return nil, errors.New("dependent event does not have a cost value but type is GT for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin > min && deiMax > max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}
			break
		case types.LTE:
			// check if the event value is less than or equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("cost dependency single value is nil but type is LTE")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is LTE for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return nil, errors.New("dependent event does not have a cost value but type is LTE for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin <= min && deiMax <= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}
			break
		case types.GTE:
			// check if the event value is greater than or equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("cost dependency single value is nil but type is GTE")
			}
			if depEvent.Event.AssociatedCost == nil {
				return nil, errors.New("dependent event does not have an associated cost but type is GTE for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return nil, errors.New("dependent event does not have a cost value but type is GTE for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedCost.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin >= min && deiMax >= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}
			break
		default:
			return nil, errors.New("invalid cost dependency type")
		}
	}

	for _, impactDependency := range event.DependsOnImpact {
		dei := impactDependency.DependentEventID
		deiSingle := impactDependency.SingleValue
		deiRange := impactDependency.Range
		deiDecomposed := impactDependency.Decomposed

		depEvent := &utils.FilteredEvent{}

		inResults := utils.EventInResults(*dei, IndependentResults)
		if inResults == nil {
			de, err := utils.FindEventByID(*dei, Events)
			if err != nil {
				de = nil
			}
			depEvent = de
		} else {
			depEvent.Event.AssociatedImpact.SingleNumber.Value = inResults.Impact
			depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation = &inResults.ImpactStandardDeviation
			depEvent.Event.AssociatedImpact.SingleNumber.Confidence = utils.Float64toPointer(0.95)
		}

		if depEvent == nil {
			return nil, errors.New("dependent event is nil")
		}

		// 2. simulate the impact dependencies and store the results
		switch impactDependency.Type {
		case types.In:
			// check if the event dependency is in a range
			if deiRange == nil {
				return nil, errors.New("impact dependency range is nil but type is In")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := dist - sd
			if min < 0 {
				min = 0
			}
			max := dist + sd

			if depEvent.Event.AssociatedImpact.SingleNumber != nil {
				minNum := depEvent.Event.AssociatedImpact.SingleNumber.Value - *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation
				maxNum := depEvent.Event.AssociatedImpact.SingleNumber.Value + *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation

				if minNum < 0 {
					minNum = 0
				}

				minCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum - (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				}

				if minNum >= min && maxNum <= max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
			} else if depEvent.Event.AssociatedImpact.Range != nil {
				minLower := depEvent.Event.AssociatedImpact.Range.Minimum.Value - *depEvent.Event.AssociatedImpact.Range.Minimum.StandardDeviation
				minUpper := depEvent.Event.AssociatedImpact.Range.Minimum.Value + *depEvent.Event.AssociatedImpact.Range.Minimum.StandardDeviation
				maxLower := depEvent.Event.AssociatedImpact.Range.Maximum.Value - *depEvent.Event.AssociatedImpact.Range.Maximum.StandardDeviation
				maxUpper := depEvent.Event.AssociatedImpact.Range.Maximum.Value + *depEvent.Event.AssociatedImpact.Range.Maximum.StandardDeviation

				if minLower < 0 {
					minLower = 0
				}

				minLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				}

				if maxLower < 0 {
					maxLower = 0
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper - (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				}

				if minLower >= min && minUpper <= max && maxLower >= min && maxUpper <= max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
			} else {
				return nil, errors.New("dependent event does not have an impact value or range but type is In for impact dependency")
			}
			break
		case types.Out:
			// check if the event dependency is out of a range
			if deiRange == nil {
				return nil, errors.New("impact dependency range is nil but type is Out")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := dist - sd
			if min < 0 {
				min = 0
			}

			max := dist + sd

			if depEvent.Event.AssociatedImpact.SingleNumber != nil {
				minNum := depEvent.Event.AssociatedImpact.SingleNumber.Value - *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation
				maxNum := depEvent.Event.AssociatedImpact.SingleNumber.Value + *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation

				if minNum < 0 {
					minNum = 0
				}

				minCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()

					if err != nil {
						return nil, err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum - (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				}

				if minNum < min && maxNum > max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
			} else if depEvent.Event.AssociatedImpact.Range != nil {
				minLower := depEvent.Event.AssociatedImpact.Range.Minimum.Value - *depEvent.Event.AssociatedImpact.Range.Minimum.StandardDeviation
				minUpper := depEvent.Event.AssociatedImpact.Range.Minimum.Value + *depEvent.Event.AssociatedImpact.Range.Minimum.StandardDeviation
				maxLower := depEvent.Event.AssociatedImpact.Range.Maximum.Value - *depEvent.Event.AssociatedImpact.Range.Maximum.StandardDeviation
				maxUpper := depEvent.Event.AssociatedImpact.Range.Maximum.Value + *depEvent.Event.AssociatedImpact.Range.Maximum.StandardDeviation

				if minLower < 0 {
					minLower = 0
				}

				minLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				}

				if maxLower < 0 {
					maxLower = 0
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper - (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				}

				if minLower < min && minUpper < max && maxLower < min && maxUpper < max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}

			} else {
				return nil, errors.New("dependent event does not have an impact value or range but type is Out for impact dependency")
			}

			break
		case types.Has:
			// check if the event decomposition has a component
			if deiDecomposed.Components == nil {
				return nil, errors.New("impact dependency decomposed components is nil but type is Has")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return nil, errors.New("dependent event does not have an associated impact but type is Has for impact dependency")
			}
			if depEvent.Event.AssociatedImpact.Decomposed == nil {
				return nil, errors.New("dependent event does not have an impact value but type is Has for impact dependency")
			}

			// check if the dependent event has a impact value
			found := -1

			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedImpact.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedImpact.Decomposed.Components[j].Name {
						found = j
						break
					}
				}
			}

			if found == -1 {
				return nil, errors.New("dependent event does not have an impact value but type is Has for impact dependency")
			}

			if depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact != nil {
				if depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact.SingleNumber != nil {
					result, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact.SingleNumber, depEvent.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					if *result > 0 {
						DependenciesMet++
					} else {
						DependenciesMissed++
					}
				} else if depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact.Range != nil {
					result, resultSD, err := simulateRange(depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact.Range, depEvent.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					if result+resultSD > 0 {
						DependenciesMet++
					} else {
						DependenciesMissed++
					}
					break
				} else {
					return nil, errors.New("dependent event does not have an impact value but type is Has for impact dependency")
				}
			} else {
				return nil, errors.New("dependent event does not have an impact value but type is Has for impact dependency")
			}

			break
		case types.HasNot:
			// check if the event decomposition does not have a component
			if deiDecomposed.Components == nil {
				return nil, errors.New("impact dependency decomposed components is nil but type is HasNot")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return nil, errors.New("dependent event does not have an associated impact but type is HasNot for impact dependency")
			}

			// check if the dependent event has a impact value
			found := -1

			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedImpact.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedImpact.Decomposed.Components[j].Name {
						found = j
						break
					}
				}
			}

			if found == -1 {
				return nil, errors.New("dependent event does not have an impact value but type is HasNot for impact dependency")
			}

			if depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact != nil {
				if depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact.SingleNumber != nil {
					result, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact.SingleNumber, depEvent.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					if *result <= 0 {
						DependenciesMet++
					} else {
						DependenciesMissed++
					}
				} else if depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact.Range != nil {
					result, resultSD, err := simulateRange(depEvent.Event.AssociatedImpact.Decomposed.Components[found].Impact.Range, depEvent.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					if result+resultSD <= 0 {
						DependenciesMet++
					} else {
						DependenciesMissed++
					}
				} else {
					return nil, errors.New("dependent event does not have an impact value but type is HasNot for impact dependency")
				}
			} else {
				return nil, errors.New("dependent event does not have an impact value but type is HasNot for impact dependency")
			}

			break
		case types.EQ:
			// check if the event value is equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("impact dependency single value is nil but type is EQ")
			}

			if depEvent.Event.AssociatedImpact.SingleNumber == nil {
				return nil, errors.New("dependent event does not have an impact value but type is EQ for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)

			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin >= min && deiMax <= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}
			break
		case types.NEQ:
			// check if the event value is not equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("impact dependency single value is nil but type is NEQ")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return nil, errors.New("dependent event does not have an associated impact but type is NEQ for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin < min && deiMax > max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		case types.LT:
			// check if the event value is less than the dependency value
			if deiSingle == nil {
				return nil, errors.New("impact dependency single value is nil but type is LT")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return nil, errors.New("dependent event does not have an associated impact but type is LT for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin < min && deiMax < max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}
			break
		case types.GT:
			// check if the event value is greater than the dependency value
			if deiSingle == nil {
				return nil, errors.New("impact dependency single value is nil but type is GT")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return nil, errors.New("dependent event does not have an associated impact but type is GT for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin > min && deiMax > max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		case types.LTE:
			// check if the event value is less than or equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("impact dependency single value is nil but type is LTE")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return nil, errors.New("dependent event does not have an associated impact but type is LTE for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {

				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin <= min && deiMax <= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		case types.GTE:
			// check if the event value is greater than or equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("impact dependency single value is nil but type is GTE")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return nil, errors.New("dependent event does not have an associated impact but type is GTE for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedImpact.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin >= min && deiMax >= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		default:
			return nil, errors.New("invalid impact dependency type")
		}
	}

	for _, probabilityDependency := range event.DependsOnProbability {
		dei := probabilityDependency.DependentEventID
		deiSingle := probabilityDependency.SingleValue
		deiRange := probabilityDependency.Range
		deiDecomposed := probabilityDependency.Decomposed

		depEvent := &utils.FilteredEvent{}

		inResults := utils.EventInResults(*dei, IndependentResults)
		if inResults == nil {
			de, err := utils.FindEventByID(*dei, Events)
			if err != nil {
				de = nil
			}
			depEvent = de
		} else {

			resultID := inResults.EventID
			for _, ev := range Events {
				if ev.ID == resultID {
					depEvent.Event = ev.Event
					break
				}
			}

			if depEvent.Event.AssociatedProbability == nil {
				de, err := utils.FindEventByID(*dei, Events)
				if err != nil {
					de = nil
				}

				depEvent = de
			}

			if depEvent.Event.AssociatedProbability != nil {
				depEvent.Event.AssociatedProbability.SingleNumber.Value = inResults.Probability
				depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation = &inResults.ProbabilityStandardDeviation
				depEvent.Event.AssociatedProbability.SingleNumber.Confidence = utils.Float64toPointer(0.95)
			} else {
				de, err := utils.FindEventByID(*dei, Events)
				if err != nil {
					de = nil
				}
				depEvent = de
			}
		}

		if depEvent == nil {
			return nil, errors.New("dependent event is nil")
		}

		// 2. simulate the probability dependencies and store the results
		switch probabilityDependency.Type {
		case types.In:
			// check if the event dependency is in a range
			if deiRange == nil {
				return nil, errors.New("probability dependency range is nil but type is In")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := dist - sd
			if min < 0 {
				min = 0
			}

			max := dist + sd

			if depEvent.Event.AssociatedProbability.SingleNumber != nil {
				minNum := depEvent.Event.AssociatedProbability.SingleNumber.Value - *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation
				maxNum := depEvent.Event.AssociatedProbability.SingleNumber.Value + *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation

				if minNum < 0 {
					minNum = 0
				}

				minCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedProbability.SingleNumber.Confidence)
				} else {

					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedProbability.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedProbability.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum - (r * *depEvent.Event.AssociatedProbability.SingleNumber.Confidence)
				}

				if minNum >= min && maxNum <= max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}

			} else if depEvent.Event.AssociatedProbability.Range != nil {
				minLower := depEvent.Event.AssociatedProbability.Range.Minimum.Value - *depEvent.Event.AssociatedProbability.Range.Minimum.StandardDeviation
				minUpper := depEvent.Event.AssociatedProbability.Range.Minimum.Value + *depEvent.Event.AssociatedProbability.Range.Minimum.StandardDeviation
				maxLower := depEvent.Event.AssociatedProbability.Range.Maximum.Value - *depEvent.Event.AssociatedProbability.Range.Maximum.StandardDeviation
				maxUpper := depEvent.Event.AssociatedProbability.Range.Maximum.Value + *depEvent.Event.AssociatedProbability.Range.Maximum.StandardDeviation

				if minLower < 0 {
					minLower = 0
				}

				minLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedProbability.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedProbability.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedProbability.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedProbability.Range.Minimum.Confidence)
				}

				if maxLower < 0 {
					maxLower = 0
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedProbability.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedProbability.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedProbability.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper - (r * *depEvent.Event.AssociatedProbability.Range.Maximum.Confidence)
				}

				if minLower >= min && minUpper <= max && maxLower >= min && maxUpper <= max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
			} else {
				return nil, errors.New("dependent event does not have a probability value or range but type is In for probability dependency")
			}

			break
		case types.Out:
			// check if the event dependency is out of a range
			if deiRange == nil {
				return nil, errors.New("probability dependency range is nil but type is Out")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := dist - sd
			if min < 0 {
				min = 0
			}

			max := dist + sd

			if depEvent.Event.AssociatedProbability.SingleNumber != nil {
				minNum := depEvent.Event.AssociatedProbability.SingleNumber.Value - *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation
				maxNum := depEvent.Event.AssociatedProbability.SingleNumber.Value + *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation

				if minNum < 0 {
					minNum = 0
				}

				minCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedProbability.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()

					if err != nil {
						return nil, err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedProbability.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedProbability.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxNum = maxNum - (r * *depEvent.Event.AssociatedProbability.SingleNumber.Confidence)
				}

				if minNum < min && maxNum > max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}

			} else if depEvent.Event.AssociatedProbability.Range != nil {
				minLower := depEvent.Event.AssociatedProbability.Range.Minimum.Value - *depEvent.Event.AssociatedProbability.Range.Minimum.StandardDeviation
				minUpper := depEvent.Event.AssociatedProbability.Range.Minimum.Value + *depEvent.Event.AssociatedProbability.Range.Minimum.StandardDeviation
				maxLower := depEvent.Event.AssociatedProbability.Range.Maximum.Value - *depEvent.Event.AssociatedProbability.Range.Maximum.StandardDeviation
				maxUpper := depEvent.Event.AssociatedProbability.Range.Maximum.Value + *depEvent.Event.AssociatedProbability.Range.Maximum.StandardDeviation

				if minLower < 0 {
					minLower = 0
				}

				minLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedProbability.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedProbability.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedProbability.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedProbability.Range.Minimum.Confidence)
				}

				if maxLower < 0 {
					maxLower = 0
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedProbability.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedProbability.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return nil, err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedProbability.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return nil, err
					}

					maxUpper = maxUpper - (r * *depEvent.Event.AssociatedProbability.Range.Maximum.Confidence)
				}

				if minLower >= min && minUpper <= max && maxLower >= min && maxUpper <= max {
					DependenciesMissed++
				} else {
					DependenciesMet++
				}

			} else {
				return nil, errors.New("dependent event does not have a probability value or range but type is Out for probability dependency")
			}
		case types.Has:
			// check if the event decomposition has a component
			if deiDecomposed.Components == nil {
				return nil, errors.New("probability dependency decomposed components is nil but type is Has")
			}

			// check if the dependent event has a probability value
			found := -1

			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedProbability.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedProbability.Decomposed.Components[j].Name {
						found = j
						break
					}
				}
			}

			if found == -1 {
				return nil, errors.New("dependent event does not have a probability value but type is Has for probability dependency")
			}

			if depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability != nil {
				if depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability.SingleNumber != nil {
					result, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability.SingleNumber, depEvent.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					if *result > 0 {
						DependenciesMet++
					} else {
						DependenciesMissed++
					}

				} else if depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability.Range != nil {
					result, resultSD, err := simulateRange(depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability.Range, depEvent.Event.Timeframe)

					if err != nil {
						return nil, err
					}

					if result+resultSD > 0 {
						DependenciesMet++
					} else {
						DependenciesMissed++
					}

				} else {
					return nil, errors.New("dependent event does not have a probability value but type is Has for probability dependency")
				}
			} else {
				return nil, errors.New("dependent event does not have a probability value but type is Has for probability dependency")
			}

		case types.HasNot:
			if deiDecomposed.Components == nil {
				return nil, errors.New("probability dependency decomposed components is nil but type is HasNot")
			}

			// check if the dependent event has a probability value
			found := -1

			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedProbability.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedProbability.Decomposed.Components[j].Name {
						found = j
						break
					}
				}
			}

			if found == -1 {
				return nil, errors.New("dependent event does not have a probability value but type is HasNot for probability dependency")
			}

			if depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability != nil {
				if depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability.SingleNumber != nil {
					result, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability.SingleNumber, depEvent.Event.Timeframe)
					if err != nil {
						return nil, err
					}

					if *result <= 0 {
						DependenciesMet++
					} else {
						DependenciesMissed++
					}

				} else if depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability.Range != nil {
					result, resultSD, err := simulateRange(depEvent.Event.AssociatedProbability.Decomposed.Components[found].Probability.Range, depEvent.Event.Timeframe)

					if err != nil {
						return nil, err
					}

					if result+resultSD <= 0 {
						DependenciesMet++
					} else {
						DependenciesMissed++
					}

				} else {
					return nil, errors.New("dependent event does not have a probability value but type is HasNot for probability dependency")
				}
			} else {
				return nil, errors.New("dependent event does not have a probability value but type is HasNot for probability dependency")
			}

		case types.EQ:
			// check if the event value is equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("probability dependency single value is nil but type is EQ")
			}
			if depEvent.Event.AssociatedProbability.SingleNumber == nil {
				return nil, errors.New("dependent event does not have a probability value but type is EQ for probability dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedProbability.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {

				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin >= min && deiMax <= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

		case types.NEQ:
			// check if the event value is not equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("probability dependency single value is nil but type is NEQ")
			}
			if depEvent.Event.AssociatedProbability == nil {
				return nil, errors.New("dependent event does not have a probability value but type is NEQ for probability dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedProbability.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin < min && deiMax > max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

		case types.LT:
			// check if the event value is less than the dependency value
			if deiSingle == nil {
				return nil, errors.New("probability dependency single value is nil but type is LT")
			}
			if depEvent.Event.AssociatedProbability == nil {
				return nil, errors.New("dependent event does not have a probability value but type is LT for probability dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedProbability.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {

				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err

				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin < min && deiMax < max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

		case types.GT:
			// check if the event value is greater than the dependency value
			if deiSingle == nil {
				return nil, errors.New("probability dependency single value is nil but type is GT")
			}

			if depEvent.Event.AssociatedProbability == nil {
				return nil, errors.New("dependent event does not have a probability value but type is GT for probability dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedProbability.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin > min && deiMax > max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

		case types.LTE:
			// check if the event value is less than or equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("probability dependency single value is nil but type is LTE")
			}
			if depEvent.Event.AssociatedProbability == nil {
				return nil, errors.New("dependent event does not have a probability value but type is LTE for probability dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedProbability.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)

			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin <= min && deiMax <= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}
		case types.GTE:
			// check if the event value is greater than or equal to the dependency value
			if deiSingle == nil {
				return nil, errors.New("probability dependency single value is nil but type is GTE")
			}
			if depEvent.Event.AssociatedProbability == nil {
				return nil, errors.New("dependent event does not have a probability value but type is GTE for probability dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedProbability.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return nil, err
			}

			min := *base - *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation
			max := *base + *depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation

			deiMin := deiSingle.Value - *deiSingle.StandardDeviation
			deiMax := deiSingle.Value + *deiSingle.StandardDeviation

			if deiMin < 0 {
				deiMin = 0
			}

			minCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return nil, err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return nil, err
				}

				deiMax = deiMax - (r * *deiSingle.Confidence)
			}

			if deiMin >= min && deiMax >= max {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		default:
			return nil, errors.New("invalid probability dependency type")
		}
	}

	for _, riskDependency := range event.DependsOnRisk {
		// Simulate the risk dependencies and store the results
		switch riskDependency.Type {
		case types.Exists, types.DoesNotExist:
			// Find the risk by ID
			var riskExists bool
			for _, risk := range Risks {
				if risk.ID == *riskDependency.DependentRiskID {
					// Check if the risk's probability is non-zero (exists) or zero (does not exist)
					if risk.Probability != nil && risk.Probability.SingleNumber != nil {
						if risk.Probability.SingleNumber.Value > 0 {
							riskExists = true
						}
						break
					}
				}
			}

			if (riskDependency.Type == types.Exists && !riskExists) || (riskDependency.Type == types.DoesNotExist && riskExists) {
				DependenciesMissed++
			} else {
				DependenciesMet++
			}
		default:
			return nil, errors.New("invalid risk dependency type")
		}
	}

	for _, eventDependency := range event.DependsOnEvent {
		// Simulate the event dependencies and store the results
		var eventExists bool
		for _, independentResult := range IndependentResults {
			if independentResult.EventID == eventDependency.DependentEventID {
				if independentResult.Probability > 0 {
					eventExists = true
				}
				break
			}
		}

		switch eventDependency.Type {
		case types.Happens:
			if !eventExists {
				DependenciesMissed++
			} else {
				DependenciesMet++
			}
		case types.DoesNotHappen:
			if eventExists {
				DependenciesMissed++
			} else {
				DependenciesMet++
			}
		default:
			return nil, errors.New("invalid event dependency type")
		}
	}

	for _, mitigationDependency := range event.DependsOnMitigation {
		// Simulate the mitigation dependencies and store the results
		var mitigationExists bool
		for _, mitigation := range Mitigations {
			if mitigation.ID == *mitigationDependency.DependentMitigationOrRiskID {
				if mitigation.Probability != nil && mitigation.Probability.SingleNumber != nil {
					if mitigation.Probability.SingleNumber.Value > 0 {
						mitigationExists = true
					}
					break
				}
			}
		}

		switch mitigationDependency.Type {
		case types.Exists:
			if !mitigationExists {
				DependenciesMissed++
			} else {
				DependenciesMet++
			}
		case types.DoesNotExist:
			if mitigationExists {
				DependenciesMissed++
			} else {
				DependenciesMet++
			}
		default:
			return nil, errors.New("invalid mitigation dependency type")
		}
	}

	// Evaluate if all dependencies are met
	if DependenciesMet == expectedDependencies {

		eventAsIndependent := &utils.FilteredEvent{
			Event:       event,
			Independent: true,
		}

		return SimulateIndependentEvent(eventAsIndependent)
	} else {
		return nil, errors.New("not all dependencies are met")
	}
}
