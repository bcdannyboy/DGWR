package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bcdannyboy/dgws/risk"
	"github.com/bcdannyboy/dgws/risk/analysis"
	"github.com/bcdannyboy/dgws/risk/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Assuming IDs are generated correctly
	vulnerabilityIntroductionID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for vulnerability introduction: %w", err))
	}
	vulnerabilityDiscoveryID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for vulnerability discovery: %w", err))
	}
	vulnerabilityExploitID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for vulnerability exploit: %w", err))
	}

	vulnerabilityIntroduction := &risk.Event{
		ID:          vulnerabilityIntroductionID,
		Name:        "Vulnerability Introduction",
		Description: "New vulnerabilities are introduced into the codebase monthly.",
		Probability: &risk.Probability{
			ExpectedFrequency: "monthly",
			Minimum:           0.27,
			MinimumConfidence: 0.9,
			Maximum:           0.27,
			MaximumConfidence: 0.9,
		},
	}

	vulnerabilityDiscovery := &risk.Event{
		ID:          vulnerabilityDiscoveryID,
		Name:        "Vulnerability Discovery",
		Description: "Some vulnerabilities are discovered before being exploited.",
		Probability: &risk.Probability{
			ExpectedFrequency: "yearly",
			Minimum:           0.5, // Assuming a better scenario for discovery
			MinimumConfidence: 0.8,
			Maximum:           0.8,
			MaximumConfidence: 0.8,
		},
	}

	vulnerabilityExploit := &risk.Event{
		ID:          vulnerabilityExploitID,
		Name:        "Vulnerability Exploit",
		Description: "Undiscovered vulnerabilities may be exploited.",
		Probability: &risk.Probability{
			ExpectedFrequency: "yearly",
			Minimum:           0.05, // Assuming lower exploit rate due to effective mitigations
			MinimumConfidence: 0.7,
			Maximum:           0.2,
			MaximumConfidence: 0.7,
		},
	}

	events := []*risk.Event{vulnerabilityIntroduction, vulnerabilityDiscovery, vulnerabilityExploit}

	probabilityMap, _, err := analysis.MonteCarlo(events, 100000)
	if err != nil {
		panic(fmt.Errorf("error running Monte Carlo analysis: %w", err))
	}

	outputResults(events, probabilityMap)
}

func outputResults(events []*risk.Event, probabilityMap map[int]float64) {
	ProbMap := make(map[string]float64)
	for k, v := range probabilityMap {
		for _, e := range events {
			if e.ID == k {
				ProbMap[e.Name] = v
			}
		}
	}

	probJSON, _ := json.MarshalIndent(ProbMap, "", "    ")

	fmt.Println("Yearly Probabilities:", string(probJSON))

	err := os.WriteFile("yearly_probabilities.json", probJSON, 0644)
	if err != nil {
		fmt.Println("Error writing yearly probabilities to file:", err)
	}
}
