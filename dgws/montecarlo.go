package dgws

import (
	"fmt"

	"github.com/bcdannyboy/montecargo/dgws/simulator"
	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

type MonteCarlo struct {
	Iterations  int                 `json:"iterations"`
	Events      []*types.Event      `json:"events"`
	Risks       []*types.Risk       `json:"risks"`
	Mitigations []*types.Mitigation `json:"mitigations"`
}

func (m *MonteCarlo) Simulate() {
	fmt.Printf("Simulating %d iterations with %d events\n", m.Iterations, len(m.Events))

	// 1. filter out dependent and independent events
	filteredEvents, badEvents := utils.FilterDependencies(m.Events, m.Risks, m.Mitigations)
	if len(badEvents) > 0 {
		fmt.Printf("Found %d bad events\n", len(badEvents))
	}
	if len(filteredEvents) == 0 {
		fmt.Println("No events to simulate")
		return
	}

	IndependentEvents := []*utils.FilteredEvent{}
	DependentEvents := []*utils.FilteredEvent{}
	for _, event := range filteredEvents {
		if event.Independent {
			IndependentEvents = append(IndependentEvents, event)
		} else {
			DependentEvents = append(DependentEvents, event)
		}
	}

	// 2. simulate all the indepdendent events seperately and store the results
	IndependentResults, err := SimulateIndependentevents(IndependentEvents, m.Iterations)
	if err != nil {
		fmt.Println("Error simulating independent events")
		return
	}

	// 3. simulate all the depedenent events seperately and store the results
	DependentResults := []*types.SimulationResults{}
	for i := 0; i < m.Iterations; i++ {
		for _, event := range DependentEvents {
			DependentResult, err := simulator.SimulateDependentEvent(event.Event, filteredEvents, m.Risks, m.Mitigations, IndependentResults)
			if err != nil {
				fmt.Printf("Error simulating dependent event %d\n", event.ID)
				return
			}

			DependentResults = append(DependentResults, DependentResult)
		}
	}

	// 4. combine the results of the independent and dependent events
	AllResults := append(IndependentResults, DependentResults...)

	// 5. return the results
}
