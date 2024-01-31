package dgws

import (
	"fmt"

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

	// 2. simulate all the indepdendent events seperately and store the results
	for _, event := range filteredEvents {
		if event.Independent {
			fmt.Printf("Simulating event %d\n", event.ID)
		}
	}

	// 3. simulate all the depedenent events seperately and store the results

	// 4. combine the results of the independent and dependent events

	// 5. return the results
}
