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

		if DoEvent == nil || DoEvent.Event == nil {
			return false, fmt.Errorf("dependent event %d is nil", DoE.DependentEventID)
		}

		// recursive dependency check
		if !DoEvent.Independent {
			dID := DoEvent.DependentEventID
			dValue := DoEvent.DependencyValue
			dRange := DoEvent.DependencyRange
			dDecomp := DoEvent.DependencyDecomp

			if dID == nil {
				return false, fmt.Errorf("dependent event %d has no dependency event", DoEvent.Event.ID)
			}

			if dValue == nil && dRange == nil && dDecomp == nil {
				return false, fmt.Errorf("dependent event %d has no dependency value, range, or decomposed", DoEvent.Event.ID)
			}

			dEvent, err := utils.FindEventByID(*dID, Events)
			if err != nil {
				return false, fmt.Errorf("dependent event %d not found", *dID)
			}

			if dEvent == nil || dEvent.Event == nil {
				return false, fmt.Errorf("dependent event %d is nil", *dID)
			}

			hitormiss, err := DependencyCheck(dEvent, DoEvent.DependencyType, Events, dValue, dRange, dDecomp)
			if err != nil {
				return false, fmt.Errorf("error checking dependency for dependent event %d: %s", DoEvent.Event.ID, err.Error())
			}

			if !hitormiss {
				return false, nil // missed dependency
			}
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
		DoPEvent, err := utils.FindEventByID(*DoP.DependentEventID, Events)
		if err != nil {
			return false, fmt.Errorf("dependent event %d not found", *DoP.DependentEventID)
		}

		if DoPEvent == nil || DoPEvent.Event == nil {
			return false, fmt.Errorf("dependent event %d is nil", *DoP.DependentEventID)
		}

		// recursive dependency check
		if !DoPEvent.Independent {
			dID := DoPEvent.DependentEventID
			dValue := DoPEvent.DependencyValue
			dRange := DoPEvent.DependencyRange
			dDecomp := DoPEvent.DependencyDecomp

			if dID == nil {
				return false, fmt.Errorf("dependent event %d has no dependency event", DoPEvent.Event.ID)
			}

			if dValue == nil && dRange == nil && dDecomp == nil {
				return false, fmt.Errorf("dependent event %d has no dependency value, range, or decomposed", DoPEvent.Event.ID)
			}

			dEvent, err := utils.FindEventByID(*dID, Events)
			if err != nil {
				return false, fmt.Errorf("dependent event %d not found", *dID)

			}

			if dEvent == nil || dEvent.Event == nil {
				return false, fmt.Errorf("dependent event %d is nil", *dID)
			}

			hitormiss, err := DependencyCheck(dEvent, DoPEvent.DependencyType, Events, dValue, dRange, dDecomp)
			if err != nil {
				return false, fmt.Errorf("error checking dependency for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
			}

			if !hitormiss {
				return false, nil // missed dependency
			}
		}

		// Process Depends on Probability
		switch DType {
		case types.Has:
			// has means the probability is dependent on a non-zero specific component of the decomposed attribute of the dependent event

			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.Decomposed == nil {
				return false, fmt.Errorf("dependent event has no decomposed probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyDecomp == nil {
				return false, fmt.Errorf("dependent event has no decomposed dependency to compare with %d", DepEvent.Event.ID)
			}

			found := false
			for _, component := range DoPEvent.Event.AssociatedProbability.Decomposed.Components {
				for _, expectedComponent := range DoPEvent.DependencyDecomp.Components {
					if component.ComponentID == expectedComponent.ComponentID {
						found = true
						if component.Probability != nil {
							if component.Probability.SingleNumber != nil {
								base, std, err := simulateSingleNumber(component.Probability.SingleNumber, DoPEvent.Event.Timeframe)
								if err != nil {
									return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
								}

								min := base - std
								if min <= 0 {
									return false, nil // missed dependency
								}
							} else if component.Probability.Range != nil {
								base, std, err := simulateRange(component.Probability.Range, DoPEvent.Event.Timeframe)
								if err != nil {
									return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
								}

								min := base - std
								if min <= 0 {
									return false, nil // missed dependency
								}
							} else if component.Probability.Decomposed != nil {
								base, std, err := simulateDecomposedByAttribute(component.Probability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)
								if err != nil {
									return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
								}

								min := base - std
								if min <= 0 {
									return false, nil // missed dependency
								}
							} else {
								return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
							}
						} else {
							return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
						}
					}
				}
			}

			if !found {
				return false, fmt.Errorf("dependent event has no decomposed dependency to compare with %d", DepEvent.Event.ID)
			}

			break
		case types.HasNot:
			// has not means the probability is dependent on a zero specific component of the decomposed attribute of the dependent event

			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.Decomposed == nil {
				return false, fmt.Errorf("dependent event has no decomposed probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyDecomp == nil {
				return false, fmt.Errorf("dependent event has no decomposed dependency to compare with %d", DepEvent.Event.ID)
			}

			found := false
			for _, component := range DoPEvent.Event.AssociatedProbability.Decomposed.Components {
				for _, expectedComponent := range DoPEvent.DependencyDecomp.Components {
					if component.ComponentID == expectedComponent.ComponentID {
						found = true
						if component.Probability != nil {
							if component.Probability.SingleNumber != nil {
								base, std, err := simulateSingleNumber(component.Probability.SingleNumber, DoPEvent.Event.Timeframe)
								if err != nil {
									return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
								}

								if base-std > 0 {
									return false, nil // missed dependency
								}

							} else if component.Probability.Range != nil {
								base, std, err := simulateRange(component.Probability.Range, DoPEvent.Event.Timeframe)
								if err != nil {
									return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
								}

								if base-std > 0 {
									return false, nil // missed dependency
								}

							} else if component.Probability.Decomposed != nil {
								base, std, err := simulateDecomposedByAttribute(component.Probability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)
								if err != nil {
									return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
								}

								if base-std > 0 {
									return false, nil // missed dependency

								}

							} else {
								return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
							}

						} else {
							return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
						}

					}

				}

			}

			if !found {
				return false, fmt.Errorf("dependent event has no decomposed dependency to compare with %d", DepEvent.Event.ID)
			}

			break
		case types.In:
			// in means the probability is in a specific range
			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyRange == nil {
				return false, fmt.Errorf("dependent event has no range dependency to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.SingleNumber != nil {
				base, std, err := simulateSingleNumber(DoPEvent.Event.AssociatedProbability.SingleNumber, DoPEvent.Event.Timeframe)
				if err != nil {
					return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depRangeAbsoluteMin := DoPEvent.DependencyRange.Minimum.Value
				if DoPEvent.DependencyRange.Minimum.StandardDeviation != nil {
					depRangeAbsoluteMin = depRangeAbsoluteMin - *DoPEvent.DependencyRange.Minimum.StandardDeviation
				}
				if DoPEvent.DependencyRange.Minimum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin + (depRangeAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin - (depRangeAbsoluteMin * ci)
					}
				}

				depRangeAbsoluteMax := DoPEvent.DependencyRange.Maximum.Value
				if DoPEvent.DependencyRange.Maximum.StandardDeviation != nil {
					depRangeAbsoluteMax = depRangeAbsoluteMax + *DoPEvent.DependencyRange.Maximum.StandardDeviation
				}
				if DoPEvent.DependencyRange.Maximum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax + (depRangeAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax - (depRangeAbsoluteMax * ci)
					}
				}

				if !(min < depRangeAbsoluteMin && max > depRangeAbsoluteMax) { // partial in is ok
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Range != nil {
				base, std, err := simulateRange(DoPEvent.Event.AssociatedProbability.Range, DoPEvent.Event.Timeframe)
				if err != nil {
					return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depRangeAbsoluteMin := DoPEvent.DependencyRange.Minimum.Value
				if DoPEvent.DependencyRange.Minimum.StandardDeviation != nil {
					depRangeAbsoluteMin = depRangeAbsoluteMin - *DoPEvent.DependencyRange.Minimum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Minimum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin + (depRangeAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin - (depRangeAbsoluteMin * ci)
					}
				}

				depRangeAbsoluteMax := DoPEvent.DependencyRange.Maximum.Value
				if DoPEvent.DependencyRange.Maximum.StandardDeviation != nil {
					depRangeAbsoluteMax = depRangeAbsoluteMax + *DoPEvent.DependencyRange.Maximum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Maximum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax + (depRangeAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax - (depRangeAbsoluteMax * ci)
					}
				}

				if !(min < depRangeAbsoluteMin && max > depRangeAbsoluteMax) { // partial in is ok
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Decomposed != nil {
				base, std, err := simulateDecomposedByAttribute(DoPEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)
				if err != nil {
					return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depRangeAbsoluteMin := DoPEvent.DependencyRange.Minimum.Value
				if DoPEvent.DependencyRange.Minimum.StandardDeviation != nil {
					depRangeAbsoluteMin = depRangeAbsoluteMin - *DoPEvent.DependencyRange.Minimum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Minimum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin + (depRangeAbsoluteMin * ci)

					} else {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin - (depRangeAbsoluteMin * ci)
					}
				}

				depRangeAbsoluteMax := DoPEvent.DependencyRange.Maximum.Value
				if DoPEvent.DependencyRange.Maximum.StandardDeviation != nil {
					depRangeAbsoluteMax = depRangeAbsoluteMax + *DoPEvent.DependencyRange.Maximum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Maximum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax + (depRangeAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax - (depRangeAbsoluteMax * ci)
					}

				}

				if !(min < depRangeAbsoluteMin && max > depRangeAbsoluteMax) { // partial in is ok
					return false, nil // missed dependency
				}

			} else {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}
			break
		case types.Out:
			// out means the probability is outside a specific range
			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyRange == nil {
				return false, fmt.Errorf("dependent event has no range dependency to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.SingleNumber != nil {
				base, std, err := simulateSingleNumber(DoPEvent.Event.AssociatedProbability.SingleNumber, DoPEvent.Event.Timeframe)
				if err != nil {
					return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depRangeAbsoluteMin := DoPEvent.DependencyRange.Minimum.Value
				if DoPEvent.DependencyRange.Minimum.StandardDeviation != nil {
					depRangeAbsoluteMin = depRangeAbsoluteMin - *DoPEvent.DependencyRange.Minimum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Minimum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin + (depRangeAbsoluteMin * ci)

					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin - (depRangeAbsoluteMin * ci)
					}

				}

				depRangeAbsoluteMax := DoPEvent.DependencyRange.Maximum.Value
				if DoPEvent.DependencyRange.Maximum.StandardDeviation != nil {
					depRangeAbsoluteMax = depRangeAbsoluteMax + *DoPEvent.DependencyRange.Maximum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Maximum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax + (depRangeAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax - (depRangeAbsoluteMax * ci)
					}

				}

				if !(min > depRangeAbsoluteMin && max < depRangeAbsoluteMax) { // partial out is ok
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Range != nil {
				base, std, err := simulateRange(DoPEvent.Event.AssociatedProbability.Range, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depRangeAbsoluteMin := DoPEvent.DependencyRange.Minimum.Value
				if DoPEvent.DependencyRange.Minimum.StandardDeviation != nil {
					depRangeAbsoluteMin = depRangeAbsoluteMin - *DoPEvent.DependencyRange.Minimum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Minimum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin + (depRangeAbsoluteMin * ci)

					} else {

						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin - (depRangeAbsoluteMin * ci)

					}

				}

				depRangeAbsoluteMax := DoPEvent.DependencyRange.Maximum.Value
				if DoPEvent.DependencyRange.Maximum.StandardDeviation != nil {

					depRangeAbsoluteMax = depRangeAbsoluteMax + *DoPEvent.DependencyRange.Maximum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Maximum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax + (depRangeAbsoluteMax * ci)
					} else {

						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax - (depRangeAbsoluteMax * ci)

					}

				}

				if !(min > depRangeAbsoluteMin && max < depRangeAbsoluteMax) { // partial out is ok
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Decomposed != nil {
				base, std, err := simulateDecomposedByAttribute(DoPEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depRangeAbsoluteMin := DoPEvent.DependencyRange.Minimum.Value
				if DoPEvent.DependencyRange.Minimum.StandardDeviation != nil {
					depRangeAbsoluteMin = depRangeAbsoluteMin - *DoPEvent.DependencyRange.Minimum.StandardDeviation
				}

				if DoPEvent.DependencyRange.Minimum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin + (depRangeAbsoluteMin * ci)

					} else {

						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Minimum.Confidence)

						depRangeAbsoluteMin = depRangeAbsoluteMin - (depRangeAbsoluteMin * ci)
					}

				}

				depRangeAbsoluteMax := DoPEvent.DependencyRange.Maximum.Value
				if DoPEvent.DependencyRange.Maximum.StandardDeviation != nil {
					depRangeAbsoluteMax = depRangeAbsoluteMax + *DoPEvent.DependencyRange.Maximum.StandardDeviation

				}

				if DoPEvent.DependencyRange.Maximum.Confidence != nil {
					confCF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if confCF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax + (depRangeAbsoluteMax * ci)

					} else {

						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyRange.Maximum.Confidence)

						depRangeAbsoluteMax = depRangeAbsoluteMax - (depRangeAbsoluteMax * ci)
					}

				}

				if !(min > depRangeAbsoluteMin && max < depRangeAbsoluteMax) { // partial out is ok
					return false, nil // missed dependency
				}

			} else {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			break
		case types.EQ:
			// eq means the probability is equal to a specific value
			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyValue == nil {
				return false, fmt.Errorf("dependent event has no value dependency to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.SingleNumber != nil {
				base, std, err := simulateSingleNumber(DoPEvent.Event.AssociatedProbability.SingleNumber, DoPEvent.Event.Timeframe)
				if err != nil {
					return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depAbsoluteMin := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)

					} else {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				depAbsoluteMax := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()
					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)

					} else {
						confimpact, err := utils.CryptoRandFloat64()
						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(min > depAbsoluteMin && max < depAbsoluteMax) { // value should be in std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Range != nil {
				base, std, err := simulateRange(DoPEvent.Event.AssociatedProbability.Range, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depAbsoluteMin := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				depAbsoluteMax := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {

					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {

						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(min > depAbsoluteMin && max < depAbsoluteMax) { // value should be in std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Decomposed != nil {
				base, std, err := simulateDecomposedByAttribute(DoPEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depAbsoluteMin := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				depAbsoluteMax := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(min > depAbsoluteMin && max < depAbsoluteMax) { // value should be in std range of dependency
					return false, nil // missed dependency
				}

			} else {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			break
		case types.NEQ:
			// neq means the probability is not equal to a specific value
			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyValue == nil {
				return false, fmt.Errorf("dependent event has no value dependency to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.SingleNumber != nil {
				base, std, err := simulateSingleNumber(DoPEvent.Event.AssociatedProbability.SingleNumber, DoPEvent.Event.Timeframe)
				if err != nil {
					return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depAbsoluteMin := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation

				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				depAbsoluteMax := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(min < depAbsoluteMin && max > depAbsoluteMax) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Range != nil {
				base, std, err := simulateRange(DoPEvent.Event.AssociatedProbability.Range, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depAbsoluteMin := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				depAbsoluteMax := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation

				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())

					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(min < depAbsoluteMin && max > depAbsoluteMax) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Decomposed != nil {
				base, std, err := simulateDecomposedByAttribute(DoPEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std
				max := base + std

				depAbsoluteMin := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation

				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				depAbsoluteMax := DoPEvent.DependencyValue.Value
				if DoPEvent.DependencyValue.StandardDeviation != nil {

					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}
				}

				if !(min < depAbsoluteMin && max > depAbsoluteMax) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}
			break
		case types.LT:
			// lt means the probability is less than a specific value
			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)

			}

			if DoPEvent.DependencyValue == nil {
				return false, fmt.Errorf("dependent event has no value dependency to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.SingleNumber != nil {
				base, std, err := simulateSingleNumber(DoPEvent.Event.AssociatedProbability.SingleNumber, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				max := base + std

				depAbsoluteMax := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(max < depAbsoluteMax) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Range != nil {
				base, std, err := simulateRange(DoPEvent.Event.AssociatedProbability.Range, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				max := base + std

				depAbsoluteMax := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(max < depAbsoluteMax) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Decomposed != nil {
				base, std, err := simulateDecomposedByAttribute(DoPEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				max := base + std

				depAbsoluteMax := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {

						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(max < depAbsoluteMax) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			break
		case types.GT:
			// gt means the probability is greater than a specific value
			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyValue == nil {
				return false, fmt.Errorf("dependent event has no value dependency to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.SingleNumber != nil {
				base, std, err := simulateSingleNumber(DoPEvent.Event.AssociatedProbability.SingleNumber, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std

				depAbsoluteMin := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)

					} else {

						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}
				}

				if !(min > depAbsoluteMin) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Range != nil {
				base, std, err := simulateRange(DoPEvent.Event.AssociatedProbability.Range, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std

				depAbsoluteMin := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation

				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {

						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)

					}

				}

				if !(min > depAbsoluteMin) { // value should be fully out of std range of dependency

					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Decomposed != nil {
				base, std, err := simulateDecomposedByAttribute(DoPEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std

				depAbsoluteMin := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				if !(min > depAbsoluteMin) { // value should be fully out of std range of dependency

					return false, nil // missed dependency
				}

			} else {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			break
		case types.LTE:
			// lte means the probability is less than or equal to a specific value
			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyValue == nil {
				return false, fmt.Errorf("dependent event has no value dependency to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.SingleNumber != nil {
				base, std, err := simulateSingleNumber(DoPEvent.Event.AssociatedProbability.SingleNumber, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				max := base + std

				depAbsoluteMax := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {

					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation

				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(max <= depAbsoluteMax) { // value should be fully out of std range of dependency

					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Range != nil {
				base, std, err := simulateRange(DoPEvent.Event.AssociatedProbability.Range, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				max := base + std

				depAbsoluteMax := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {

							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(max <= depAbsoluteMax) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Decomposed != nil {
				base, std, err := simulateDecomposedByAttribute(DoPEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				max := base + std

				depAbsoluteMax := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMax = depAbsoluteMax + *DoPEvent.DependencyValue.StandardDeviation

				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {

						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax + (depAbsoluteMax * ci)

					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMax = depAbsoluteMax - (depAbsoluteMax * ci)
					}

				}

				if !(max <= depAbsoluteMax) { // value should be fully out of std range of dependency

					return false, nil // missed dependency
				}

			} else {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			break
		case types.GTE:
			// gte means the probability is greater than or equal to a specific value
			if DoPEvent.Event.AssociatedProbability == nil {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.DependencyValue == nil {
				return false, fmt.Errorf("dependent event has no value dependency to compare with %d", DepEvent.Event.ID)
			}

			if DoPEvent.Event.AssociatedProbability.SingleNumber != nil {
				base, std, err := simulateSingleNumber(DoPEvent.Event.AssociatedProbability.SingleNumber, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating single number for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std

				depAbsoluteMin := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())

					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				if !(min >= depAbsoluteMin) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Range != nil {
				base, std, err := simulateRange(DoPEvent.Event.AssociatedProbability.Range, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating range for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std

				depAbsoluteMin := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation

				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				if !(min >= depAbsoluteMin) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else if DoPEvent.Event.AssociatedProbability.Decomposed != nil {
				base, std, err := simulateDecomposedByAttribute(DoPEvent.Event.AssociatedProbability.Decomposed, ProbabilityAttribute, DoPEvent.Event.Timeframe)

				if err != nil {
					return false, fmt.Errorf("error simulating decomposed for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
				}

				min := base - std

				depAbsoluteMin := DoPEvent.DependencyValue.Value

				if DoPEvent.DependencyValue.StandardDeviation != nil {
					depAbsoluteMin = depAbsoluteMin - *DoPEvent.DependencyValue.StandardDeviation
				}

				if DoPEvent.DependencyValue.Confidence != nil {
					CF, err := utils.CoinFlip()

					if err != nil {
						return false, fmt.Errorf("error simulating coin flip for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
					}

					if CF {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin + (depAbsoluteMin * ci)
					} else {
						confimpact, err := utils.CryptoRandFloat64()

						if err != nil {
							return false, fmt.Errorf("error simulating crypto rand float64 for dependent event %d: %s", DoPEvent.Event.ID, err.Error())
						}

						ci := confimpact - (confimpact * *DoPEvent.DependencyValue.Confidence)

						depAbsoluteMin = depAbsoluteMin - (depAbsoluteMin * ci)
					}

				}

				if !(min >= depAbsoluteMin) { // value should be fully out of std range of dependency
					return false, nil // missed dependency
				}

			} else {
				return false, fmt.Errorf("dependent event has no probability to compare with %d", DepEvent.Event.ID)
			}

			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	for _, DoI := range DepEvent.Event.DependsOnImpact {
		// Process Depends on Impact
		DoIEvent, err := utils.FindEventByID(*DoI.DependentEventID, Events)
		if err != nil {
			return false, fmt.Errorf("dependent event %d not found", *DoI.DependentEventID)
		}

		if DoIEvent == nil || DoIEvent.Event == nil {
			return false, fmt.Errorf("dependent event %d is nil", *DoI.DependentEventID)
		}

		// recursive dependency check
		if !DoIEvent.Independent {
			dID := DoIEvent.DependentEventID
			dValue := DoIEvent.DependencyValue
			dRange := DoIEvent.DependencyRange
			dDecomp := DoIEvent.DependencyDecomp

			if dID == nil {
				return false, fmt.Errorf("dependent event %d has no dependency event", DoIEvent.Event.ID)
			}

			if dValue == nil && dRange == nil && dDecomp == nil {
				return false, fmt.Errorf("dependent event %d has no dependency value, range, or decomposed", DoIEvent.Event.ID)
			}

			dEvent, err := utils.FindEventByID(*dID, Events)
			if err != nil {
				return false, fmt.Errorf("dependent event %d not found", *dID)
			}

			if dEvent == nil || dEvent.Event == nil {
				return false, fmt.Errorf("dependent event %d is nil", *dID)
			}

			hitormiss, err := DependencyCheck(dEvent, DoIEvent.DependencyType, Events, dValue, dRange, dDecomp)

			if err != nil {
				return false, fmt.Errorf("error checking dependency for dependent event %d: %s", DoIEvent.Event.ID, err.Error())
			}

			if !hitormiss {
				return false, nil // missed dependency
			}

		}

		switch DType {
		case types.Has:
			// has means the impact is dependent on a non-zero specific component of the decomposed attribute of the dependent event
			break
		case types.HasNot:
			// has not means the impact is dependent on a zero specific component of the decomposed attribute of the dependent event
			break
		case types.In:
			// in means the impact is in a specific range
			break
		case types.Out:
			// out means the impact is outside a specific range
			break
		case types.EQ:
			// eq means the impact is equal to a specific value
			break
		case types.NEQ:
			// neq means the impact is not equal to a specific value
			break
		case types.LT:
			// lt means the impact is less than a specific value
			break
		case types.GT:
			// gt means the impact is greater than a specific value
			break
		case types.LTE:
			// lte means the impact is less than or equal to a specific value
			break
		case types.GTE:
			// gte means the impact is greater than or equal to a specific value
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	for _, DoC := range DepEvent.Event.DependsOnCost {
		// Process Depends on Cost
		DoCEvent, err := utils.FindEventByID(*DoC.DependentEventID, Events)
		if err != nil {
			return false, fmt.Errorf("dependent event %d not found", *DoC.DependentEventID)
		}

		if DoCEvent == nil || DoCEvent.Event == nil {
			return false, fmt.Errorf("dependent event %d is nil", *DoC.DependentEventID)
		}

		// recursive dependency check
		if !DoCEvent.Independent {
			dID := DoCEvent.DependentEventID
			dValue := DoCEvent.DependencyValue
			dRange := DoCEvent.DependencyRange
			dDecomp := DoCEvent.DependencyDecomp

			if dID == nil {
				return false, fmt.Errorf("dependent event %d has no dependency event", DoCEvent.Event.ID)
			}

			if dValue == nil && dRange == nil && dDecomp == nil {
				return false, fmt.Errorf("dependent event %d has no dependency value, range, or decomposed", DoCEvent.Event.ID)
			}

			dEvent, err := utils.FindEventByID(*dID, Events)
			if err != nil {
				return false, fmt.Errorf("dependent event %d not found", *dID)
			}

			if dEvent == nil || dEvent.Event == nil {
				return false, fmt.Errorf("dependent event %d is nil", *dID)
			}

			hitormiss, err := DependencyCheck(dEvent, DoCEvent.DependencyType, Events, dValue, dRange, dDecomp)
			if err != nil {
				return false, fmt.Errorf("error checking dependency for dependent event %d: %s", DoCEvent.Event.ID, err.Error())
			}

			if !hitormiss {
				return false, nil // missed dependency
			}
		}

		switch DType {
		case types.Has:
			// has means the cost is dependent on a non-zero specific component of the decomposed attribute of the dependent event
			break
		case types.HasNot:
			// has not means the cost is dependent on a zero specific component of the decomposed attribute of the dependent event
			break
		case types.In:
			// in means the cost is in a specific range
			break
		case types.Out:
			// out means the cost is outside a specific range
			break
		case types.EQ:
			// eq means the cost is equal to a specific value
			break
		case types.NEQ:
			// neq means the cost is not equal to a specific value
			break
		case types.LT:
			// lt means the cost is less than a specific value
			break
		case types.GT:
			// gt means the cost is greater than a specific value
			break
		case types.LTE:
			// lte means the cost is less than or equal to a specific value
			break
		case types.GTE:
			// gte means the cost is greater than or equal to a specific value
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	for _, DoR := range DepEvent.Event.DependsOnRisk {
		// Process Depends on Risk

		switch DType {
		case types.Exists:
			// exists means the risk has a non-zero probability and impact
			break
		case types.DoesNotExist:
			// does not exist means the risk has a zero probability and impact
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	for _, DoM := range DepEvent.Event.DependsOnMitigation {
		// Process Depends on Mitigation
		switch DType {
		case types.Exists:
			// exists means the mitigation has a non-zero probability, impact, and cost
			break
		case types.DoesNotExist:
			// does not exist means the mitigation has a zero probability, impact, and cost
			break
		default:
			return false, fmt.Errorf("invalid dependency type")
		}
	}

	return true, nil
}
