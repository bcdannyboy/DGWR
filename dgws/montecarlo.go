package dgws

import (
	"errors"
	"fmt"

	"github.com/bcdannyboy/montecargo/dgws/types"
	"github.com/bcdannyboy/montecargo/dgws/utils"
)

type MonteCarlo struct {
	Iterations int            `json:"iterations"`
	Events     []*types.Event `json:"events"`
}

func (m *MonteCarlo) Simulate() ([]*types.SimulationResults, error) {
	fmt.Printf("Simulating %d iterations with %d events\n", m.Iterations, len(m.Events))

	// 1. filter out dependent and independent events
	filteredEvents, badEvents := utils.FilterDependencies(m.Events)
	if len(badEvents) > 0 {
		fmt.Printf("Found %d bad events\n", len(badEvents))
	}
	if len(filteredEvents) == 0 {
		return nil, errors.New("No events to simulate")
	}

	// 2. simulate all the events and store the results
	Results, err := SimulateDependentEvents(filteredEvents, m.Iterations)
	if err != nil {
		return nil, fmt.Errorf("Error simulating dependent events: %s", err)
	}

	// 3. return the results
	return Results, nil
}
