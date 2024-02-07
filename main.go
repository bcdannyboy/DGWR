package main

import (
	"fmt"

	"github.com/bcdannyboy/montecargo/dgws"
	"github.com/bcdannyboy/montecargo/dgws/types"
)

func main() {

	phishingAttackProbability := &types.Probability{
		ID:           1,
		Name:         "Phishing Attack Likelihood",
		SingleNumber: &types.SingleNumber{Value: 0.3, StandardDeviation: &[]float64{0.05}[0], Confidence: &[]float64{95}[0]},
	}

	phishingAttackImpact := &types.Impact{
		ID:           1,
		Name:         "Phishing Attack Impact",
		SingleNumber: &types.SingleNumber{Value: 5000, StandardDeviation: &[]float64{1000}[0]}, // Assuming some unit of impact
	}

	phishingAttack := &types.Event{
		ID:                    1,
		Name:                  "Phishing Attack",
		AssociatedProbability: phishingAttackProbability,
		AssociatedImpact:      phishingAttackImpact,
		Timeframe:             types.Month,
	}

	dataBreachProbability := &types.Probability{
		ID:   2,
		Name: "Data Breach Likelihood",
		Decomposed: &types.Decomposed{
			Components: []*types.DecomposedComponent{
				{
					ComponentID: 1,
					Name:        "Due to Phishing",
					Probability: &types.DecomposedItem{
						SingleNumber: &types.SingleNumber{Value: 0.2, StandardDeviation: &[]float64{0.1}[0]},
					},
				},
			},
		},
	}

	dataBreachImpact := &types.Impact{
		ID:           2,
		Name:         "Data Breach Impact",
		SingleNumber: &types.SingleNumber{Value: 20000, StandardDeviation: &[]float64{5000}[0]},
	}

	dataBreach := &types.Event{
		ID:                    2,
		Name:                  "Data Breach",
		AssociatedProbability: dataBreachProbability,
		AssociatedImpact:      dataBreachImpact,
		DependsOnEvent: []*types.EventDependency{
			{ID: 1, DependentEventID: phishingAttack.ID, Type: types.Happens},
		},
		Timeframe: types.Month,
	}

	forensicCost := &types.Cost{
		ID:           1,
		Name:         "Forensic Analysis Cost",
		SingleNumber: &types.SingleNumber{Value: 10000, StandardDeviation: &[]float64{2000}[0]},
	}

	forensicAnalysis := &types.Event{
		ID:             3,
		Name:           "Forensic Analysis",
		AssociatedCost: forensicCost,
		DependsOnEvent: []*types.EventDependency{
			{ID: 2, DependentEventID: dataBreach.ID, Type: types.Happens},
		},
		Timeframe: types.Month,
	}

	// Setup MonteCarlo simulation
	Events := []*types.Event{phishingAttack, dataBreach, forensicAnalysis}

	MC := &dgws.MonteCarlo{
		Iterations: 1000,
		Events:     Events,
	}

	Results, err := MC.Simulate()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Simulation completed with %d results\n", len(Results))
}
