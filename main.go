package main

import (
	"fmt"

	"github.com/bcdannyboy/montecargo/dgws"
	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

func main() {
	// Independent Events
	PhishingAttack := &types.Event{
		ID:   3,
		Name: "Phishing Attack",
		AssociatedProbability: &types.Probability{
			SingleNumber: &types.SingleNumber{
				Value:      0.3,
				Confidence: utils.Float64toPointer(0.8),
			},
		},
		AssociatedImpact: &types.Impact{
			SingleNumber: &types.SingleNumber{
				Value:      20000, // Estimated financial impact
				Confidence: utils.Float64toPointer(0.6),
			},
		},
	}

	// Dependent Events
	SystemCompromise := &types.Event{
		ID:   4,
		Name: "System Compromise",
		AssociatedProbability: &types.Probability{
			Decomposed: &types.Decomposed{
				Components: []*types.DecomposedComponent{
					{
						Name: "Phishing Success",
						Probability: &struct {
							SingleNumber *types.SingleNumber `json:"single_number,omitempty"`
							Range        *types.Range        `json:"range,omitempty"`
							Decomposed   *types.Decomposed   `json:"decomposed,omitempty"`
						}{
							SingleNumber: &types.SingleNumber{
								Value:      0.2,
								Confidence: utils.Float64toPointer(0.7),
							},
						},
					},
					{
						Name: "External Exploit",
						Probability: &struct {
							SingleNumber *types.SingleNumber `json:"single_number,omitempty"`
							Range        *types.Range        `json:"range,omitempty"`
							Decomposed   *types.Decomposed   `json:"decomposed,omitempty"`
						}{
							SingleNumber: &types.SingleNumber{
								Value:      0.1,
								Confidence: utils.Float64toPointer(0.9),
							},
						},
					},
				},
			},
		},
		DependsOnEvent: []*types.EventDependency{
			{
				Name:             "Phishing Attack",
				Type:             types.Happens,
				DependentEventID: PhishingAttack.ID,
			},
		},
	}

	// Risks
	DataLeakRisk := &types.Risk{
		ID:   1,
		Name: "Data Leak",
		Probability: &types.Probability{
			Range: &types.Range{
				Minimum: types.Minimum{
					Value:      0.1,
					Confidence: utils.Float64toPointer(0.7),
				},
				Maximum: types.Maximum{
					Value:      0.4,
					Confidence: utils.Float64toPointer(0.7),
				},
			},
		},
		Impact: &types.Impact{
			Range: &types.Range{
				Minimum: types.Minimum{
					Value:      5000,
					Confidence: utils.Float64toPointer(0.6),
				},
				Maximum: types.Maximum{
					Value:      25000,
					Confidence: utils.Float64toPointer(0.6),
				},
			},
		},
		DependsOnEvent: []*types.EventDependency{
			{
				Name:             "System Compromise",
				Type:             types.Happens,
				DependentEventID: SystemCompromise.ID,
			},
		},
	}

	// Mitigations
	EmployeeTraining := &types.Mitigation{
		ID:   1,
		Name: "Employee Security Awareness Training",
		Probability: &types.Probability{
			SingleNumber: &types.SingleNumber{
				Value:      0.9,
				Confidence: utils.Float64toPointer(0.8),
			},
		},
		Impact: &types.Impact{
			SingleNumber: &types.SingleNumber{
				Value:      0.2, // Reduction in risk probability
				Confidence: utils.Float64toPointer(0.7),
			},
		},
		AssociatedCost: &types.Cost{
			SingleNumber: &types.SingleNumber{
				Value:      10000, // Cost of implementing the mitigation
				Confidence: utils.Float64toPointer(0.9),
			},
		},
		Mitigates: []uint64{DataLeakRisk.ID},
	}

	ForensicAnalysis := &types.Event{
		ID:   5,
		Name: "Forensic Analysis",
		AssociatedProbability: &types.Probability{
			// Assuming forensic analysis is certain if conditions are met
			SingleNumber: &types.SingleNumber{
				Value:      1.0,
				Confidence: utils.Float64toPointer(1.0),
			},
		},
		AssociatedCost: &types.Cost{
			// Assuming a fixed cost for forensic analysis
			SingleNumber: &types.SingleNumber{
				Value:      15000, // Cost of conducting forensic analysis
				Confidence: utils.Float64toPointer(0.9),
			},
		},
		DependsOnEvent: []*types.EventDependency{
			{
				// This dependency is on the occurrence of a data leak
				Name:             "Data Leak",
				Type:             types.Happens,
				DependentEventID: DataLeakRisk.ID,
			},
		},
		DependsOnImpact: []*types.ImpactDependency{
			{
				// Trigger forensic analysis if the impact of the data leak is within a certain range
				Name:             "High Impact Data Leak",
				Type:             types.GTE, // Greater than or equal to
				DependentEventID: utils.Uint64toPointer(DataLeakRisk.ID),
				SingleValue: &types.SingleNumber{
					Value: 10000, // Trigger if impact is $10,000 or more
				},
			},
		},
	}

	Events := []*types.Event{
		PhishingAttack,
		SystemCompromise,
		ForensicAnalysis,
	}

	Risks := []*types.Risk{
		DataLeakRisk,
	}

	Mitigations := []*types.Mitigation{
		EmployeeTraining,
	}

	MC := &dgws.MonteCarlo{
		Iterations:  1000,
		Events:      Events,
		Risks:       Risks,
		Mitigations: Mitigations,
	}

	Results, err := MC.Simulate()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Simulation completed with %d results\n", len(Results))
}
