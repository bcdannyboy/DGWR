package simulator

import (
	"errors"

	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

func SimulateDependentEvent(event *types.Event, Events []*utils.FilteredEvent, Risks []*types.Risk, Mitigations []*types.Mitigation, IndependentResults []*types.SimulationResults) error {
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
			return errors.New("dependent event does not exist for cost dependency")
		}

		// 2. simulate the cost dependencies and store the results
		switch costDependency.Type {
		case types.Exists:
			// check if the event dependency exists
			if depEvent == nil {
				// event does not exist, so the dependent event does not happen
				DependenciesMissed++
				break
			}

			// event exists, so the dependent event happens
			DependenciesMet++

			break
		case types.DoesNotExist:
			// check if the event dependency does not exist
			if depEvent != nil {
				// event exists, so the dependent event does not happen
				DependenciesMissed++
				break
			}

			// event does not exist, so the dependent event happens
			DependenciesMet++
			break
		case types.In:
			// check if the event dependency is in a range
			if deiRange == nil {
				return errors.New("cost dependency range is nil but type is In")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is In for cost dependency")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
					return err
				}
				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
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
					return err
				}
				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()

					if err != nil {
						return err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
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
				return errors.New("dependent event does not have a cost value or range but type is In for cost dependency")
			}
		case types.Out:
			// check if the event dependency is out of a range
			if deiRange == nil {
				return errors.New("cost dependency range is nil but type is Out")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is Out for cost dependency")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
					return err
				}

				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedCost.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
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
					return err
				}

				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedCost.Range.Minimum.Confidence)
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedCost.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
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
				return errors.New("dependent event does not have a cost value or range but type is Out for cost dependency")
			}
		case types.Has:
			// check if the event decomposition has a component
			if deiDecomposed.Components == nil {
				return errors.New("cost dependency decomposed components is nil but type is Has")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is Has for cost dependency")
			}
			if depEvent.Event.AssociatedCost.Decomposed == nil {
				return errors.New("dependent event does not have a cost value but type is Has for cost dependency")
			}

			// check if the dependent event has a cost value
			found := false
			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedCost.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedCost.Decomposed.Components[j].Name {
						found = true
						break
					}
				}
			}
			if found {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		case types.HasNot:
			// check if the event decomposition does not have a component
			if deiDecomposed.Components == nil {
				return errors.New("cost dependency decomposed components is nil but type is HasNot")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is HasNot for cost dependency")
			}
			if depEvent.Event.AssociatedCost.Decomposed == nil {
				return errors.New("dependent event does not have a cost value but type is HasNot for cost dependency")
			}

			// check if the dependent event has a cost value
			found := false
			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedCost.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedCost.Decomposed.Components[j].Name {
						found = true
						break
					}
				}
			}
			if !found {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		case types.EQ:
			// check if the event value is equal to the dependency value
			if deiSingle == nil {
				return errors.New("cost dependency single value is nil but type is EQ")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is EQ for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return errors.New("dependent event does not have a cost value but type is EQ for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("cost dependency single value is nil but type is NEQ")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is NEQ for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return errors.New("dependent event does not have a cost value but type is NEQ for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()

				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()

				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("cost dependency single value is nil but type is LT")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is LT for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return errors.New("dependent event does not have a cost value but type is LT for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("cost dependency single value is nil but type is GT")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is GT for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return errors.New("dependent event does not have a cost value but type is GT for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("cost dependency single value is nil but type is LTE")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is LTE for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return errors.New("dependent event does not have a cost value but type is LTE for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("cost dependency single value is nil but type is GTE")
			}
			if depEvent.Event.AssociatedCost == nil {
				return errors.New("dependent event does not have an associated cost but type is GTE for cost dependency")
			}
			if depEvent.Event.AssociatedCost.SingleNumber == nil {
				return errors.New("dependent event does not have a cost value but type is GTE for cost dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedCost.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
			return errors.New("invalid cost dependency type")
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
			return errors.New("dependent event is nil")
		}

		// 2. simulate the impact dependencies and store the results
		switch impactDependency.Type {
		case types.In:
			// check if the event dependency is in a range
			if deiRange == nil {
				return errors.New("impact dependency range is nil but type is In")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
					return err
				}

				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
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
					return err
				}

				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				}

				if maxLower < 0 {
					maxLower = 0
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxUpper = maxUpper - (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				}

				if minLower >= min && minUpper <= max && maxLower >= min && maxUpper <= max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}
			} else {
				return errors.New("dependent event does not have an impact value or range but type is In for impact dependency")
			}
			break
		case types.Out:
			// check if the event dependency is out of a range
			if deiRange == nil {
				return errors.New("impact dependency range is nil but type is Out")
			}

			dist, sd, err := simulateRange(deiRange, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
					return err
				}

				if minCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minNum = minNum + (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()

					if err != nil {
						return err
					}

					minNum = minNum - (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				}

				maxCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxNum = maxNum + (r * *depEvent.Event.AssociatedImpact.SingleNumber.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
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
					return err
				}

				if minLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minLower = minLower + (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minLower = minLower - (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				}

				minUpperCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if minUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minUpper = minUpper + (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					minUpper = minUpper - (r * *depEvent.Event.AssociatedImpact.Range.Minimum.Confidence)
				}

				if maxLower < 0 {
					maxLower = 0
				}

				maxLowerCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxLowerCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxLower = maxLower + (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxLower = maxLower - (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				}

				maxUpperCF, err := utils.CoinFlip()
				if err != nil {
					return err
				}

				if maxUpperCF {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxUpper = maxUpper + (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				} else {
					r, err := utils.CryptoRandFloat64()
					if err != nil {
						return err
					}

					maxUpper = maxUpper - (r * *depEvent.Event.AssociatedImpact.Range.Maximum.Confidence)
				}

				if minLower < min && minUpper < max && maxLower < min && maxUpper < max {
					DependenciesMet++
				} else {
					DependenciesMissed++
				}

			} else {
				return errors.New("dependent event does not have an impact value or range but type is Out for impact dependency")
			}

			break
		case types.Has:
			// check if the event decomposition has a component
			if deiDecomposed.Components == nil {
				return errors.New("impact dependency decomposed components is nil but type is Has")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return errors.New("dependent event does not have an associated impact but type is Has for impact dependency")
			}
			if depEvent.Event.AssociatedImpact.Decomposed == nil {
				return errors.New("dependent event does not have an impact value but type is Has for impact dependency")
			}

			// check if the dependent event has a impact value
			found := false

			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedImpact.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedImpact.Decomposed.Components[j].Name {
						found = true
						break
					}
				}
			}

			if found {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		case types.HasNot:
			// check if the event decomposition does not have a component
			if deiDecomposed.Components == nil {
				return errors.New("impact dependency decomposed components is nil but type is HasNot")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return errors.New("dependent event does not have an associated impact but type is HasNot for impact dependency")
			}

			// check if the dependent event has a impact value
			found := false

			for i := 0; i < len(deiDecomposed.Components); i++ {
				for j := 0; j < len(depEvent.Event.AssociatedImpact.Decomposed.Components); j++ {
					if deiDecomposed.Components[i].Name == depEvent.Event.AssociatedImpact.Decomposed.Components[j].Name {
						found = true
						break
					}
				}
			}

			if !found {
				DependenciesMet++
			} else {
				DependenciesMissed++
			}

			break
		case types.EQ:
			// check if the event value is equal to the dependency value
			if deiSingle == nil {
				return errors.New("impact dependency single value is nil but type is EQ")
			}

			if depEvent.Event.AssociatedImpact.SingleNumber == nil {
				return errors.New("dependent event does not have an impact value but type is EQ for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)

			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("impact dependency single value is nil but type is NEQ")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return errors.New("dependent event does not have an associated impact but type is NEQ for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("impact dependency single value is nil but type is LT")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return errors.New("dependent event does not have an associated impact but type is LT for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("impact dependency single value is nil but type is GT")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return errors.New("dependent event does not have an associated impact but type is GT for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("impact dependency single value is nil but type is LTE")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return errors.New("dependent event does not have an associated impact but type is LTE for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {

				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
				return errors.New("impact dependency single value is nil but type is GTE")
			}
			if depEvent.Event.AssociatedImpact == nil {
				return errors.New("dependent event does not have an associated impact but type is GTE for impact dependency")
			}

			base, err := SimulateIndependentSingleNumer(depEvent.Event.AssociatedImpact.SingleNumber, depEvent.Event.Timeframe)
			if err != nil {
				return err
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
				return err
			}

			if minCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMin = deiMin - (r * *deiSingle.Confidence)
			}

			maxCF, err := utils.CoinFlip()
			if err != nil {
				return err
			}

			if maxCF {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
				}

				deiMax = deiMax + (r * *deiSingle.Confidence)
			} else {
				r, err := utils.CryptoRandFloat64()
				if err != nil {
					return err
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
			return errors.New("invalid impact dependency type")
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
			depEvent.Event.AssociatedProbability.SingleNumber.Value = inResults.Probability
			depEvent.Event.AssociatedProbability.SingleNumber.StandardDeviation = &inResults.ProbabilityStandardDeviation
			depEvent.Event.AssociatedProbability.SingleNumber.Confidence = utils.Float64toPointer(0.95)
		}

		if depEvent == nil {
			return errors.New("dependent event is nil")
		}

		// 2. simulate the probability dependencies and store the results
		switch probabilityDependency.Type {
		case types.In:
			// check if the event dependency is in a range
			break
		case types.Out:
			// check if the event dependency is out of a range
			break
		case types.Has:
			// check if the event decomposition has a component
			break
		case types.HasNot:
			// check if the event decomposition does not have a component
			break
		case types.EQ:
			// check if the event value is equal to the dependency value
			break
		case types.NEQ:
			// check if the event value is not equal to the dependency value
			break
		case types.LT:
			// check if the event value is less than the dependency value
			break
		case types.GT:
			// check if the event value is greater than the dependency value
			break
		case types.LTE:
			// check if the event value is less than or equal to the dependency value
			break
		case types.GTE:
			// check if the event value is greater than or equal to the dependency value
			break
		default:
			return errors.New("invalid probability dependency type")
		}
	}

	for _, riskDependency := range event.DependsOnRisk {
		dei := riskDependency.DependentEventID

		// 2. simulate the risk dependencies and store the results
		switch riskDependency.Type {
		case types.Exists:
			// check if the event dependency exists
			break
		case types.DoesNotExist:
			// check if the event dependency does not exist
			break
		default:
			return errors.New("invalid probability dependency type")
		}
	}

	for _, eventDependency := range event.DependsOnEvent {
		dei := eventDependency.DependentEventID

		// 2. simulate the event dependencies and store the results
		switch eventDependency.Type {
		case types.Happens:
			// check if the event dependency happens
			break
		case types.DoesNotHappen:
			// check if the event dependency does not happen
			break
		default:
			return errors.New("invalid event dependency type")
		}
	}

	for _, mitigationDependency := range event.DependsOnMitigation {
		dei := mitigationDependency.DependentEventID

		// 2. simulate the mitigation dependencies and store the results
		switch mitigationDependency.Type {
		case types.Exists:
			// check if the event dependency exists
			break
		case types.DoesNotExist:
			// check if the event dependency does not exist
			break
		case types.In:
			// check if the event dependency is in a range
			break
		case types.Out:
			// check if the event dependency is out of a range
			break
		case types.Has:
			// check if the event decomposition has a component
			break
		case types.HasNot:
			// check if the event decomposition does not have a component
			break
		case types.EQ:
			// check if the event value is equal to the dependency value
			break
		case types.NEQ:
			// check if the event value is not equal to the dependency value
			break
		case types.LT:
			// check if the event value is less than the dependency value
			break
		case types.GT:
			// check if the event value is greater than the dependency value
			break
		case types.LTE:
			// check if the event value is less than or equal to the dependency value
			break
		case types.GTE:
			// check if the event value is greater than or equal to the dependency value
			break
		default:
			return errors.New("invalid impact dependency type")
		}
	}

	return nil
}
