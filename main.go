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
						ComponentID: 1,
						Name:        "Phishing Success",
						Probability: &types.DecomposedItem{
							SingleNumber: &types.SingleNumber{
								Value:      0.2,
								Confidence: utils.Float64toPointer(0.7),
							},
						},
					},
					{
						ComponentID: 2,
						Name:        "External Exploit",
						Probability: &types.DecomposedItem{
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

	// Forensic Analysis (unchanged from your initial code)
	ForensicAnalysis := &types.Event{
		ID:   5,
		Name: "Forensic Analysis",
		AssociatedProbability: &types.Probability{
			SingleNumber: &types.SingleNumber{
				Value:      1.0,
				Confidence: utils.Float64toPointer(1.0),
			},
		},
		AssociatedCost: &types.Cost{
			SingleNumber: &types.SingleNumber{
				Value:      15000, // Cost of conducting forensic analysis
				Confidence: utils.Float64toPointer(0.9),
			},
		},
		DependsOnEvent: []*types.EventDependency{
			{
				Name:             "Data Leak",
				Type:             types.Happens,
				DependentEventID: DataLeakRisk.ID,
			},
		},
		DependsOnImpact: []*types.ImpactDependency{
			{
				Name:             "High Impact Data Leak",
				Type:             types.GTE, // Greater than or equal to
				DependentEventID: utils.Uint64toPointer(DataLeakRisk.ID),
				SingleValue: &types.SingleNumber{
					Value: 10000, // Trigger if impact is $10,000 or more
				},
			},
		},
	}

	// New Event: APT Detection
	APTDetection := &types.Event{
		ID:   6,
		Name: "APT Detection",
		AssociatedProbability: &types.Probability{
			Decomposed: &types.Decomposed{
				Components: []*types.DecomposedComponent{
					{
						ComponentID: 1, // Matching the Phishing Success component
						Name:        "Effective Phishing Detection",
						Probability: &types.DecomposedItem{
							SingleNumber: &types.SingleNumber{
								Value:      0.5, // Adjusted probability for detection
								Confidence: utils.Float64toPointer(0.8),
							},
						},
					},
					{
						ComponentID: 2, // Matching the External Exploit component
						Name:        "Effective Exploit Detection",
						Probability: &types.DecomposedItem{
							SingleNumber: &types.SingleNumber{
								Value:      0.4, // Adjusted probability for detection
								Confidence: utils.Float64toPointer(0.85),
							},
						},
					},
				},
			},
		},
	}

	Events := []*types.Event{
		PhishingAttack,
		SystemCompromise,
		ForensicAnalysis,
		APTDetection, // Include the new event
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
