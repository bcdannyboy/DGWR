package risk

type Probability struct {
	ExpectedFrequency string  `json:"ExpectedFrequency"`
	Minimum           float64 `json:"Minimum"`
	MinimumConfidence float64 `json:"MinimumConfidence"`
	Maximum           float64 `json:"Maximum"`
	MaximumConfidence float64 `json:"MaximumConfidence"`
}

type Impact struct {
	ImpactID       int    `json:"ImpactID"`
	Name           string `json:"Name"`
	Unit           string `json:"Unit"`
	PositiveImpact bool   `json:"PositiveImpact"`

	Description       string `json:"Description"`
	ExpectedFrequency string `json:"ExpectedFrequency"`

	MinimumIndividualUnitImpact           float64 `json:"MinimumIndividualUnitImpact"`
	MinimumIndividualUnitImpactConfidence float64 `json:"MinimumIndividualUnitImpactConfidence"`
	MaximumIndividualUnitImpact           float64 `json:"MaximumIndividualUnitImpact"`
	MaximumIndividualUnitImpactConfidence float64 `json:"MaximumIndividualUnitImpactConfidence"`
	MinimumImpactEvents                   float64 `json:"MinimumImpactEvents"`
	MinimumImpactEventsConfidence         float64 `json:"MinimumImpactEventsConfidence"`
	MaximumImpactEvents                   float64 `json:"MaximumImpactEvents"`
	MaximumImpactEventsConfidence         float64 `json:"MaximumImpactEventsConfidence"`
}

type Dependency struct {
	DependsOnEventID int  `json:"DependsOnEventID"`
	Happens          bool `json:"Happens"`
}

type Event struct {
	ID           int           `json:"ID"`
	Name         string        `json:"Name"`
	Description  string        `json:"Description"`
	Probability  *Probability  `json:"Probability"`
	Impact       []*Impact     `json:"Impact,omitempty"`
	Dependencies []*Dependency `json:"Dependencies,omitempty"`
}
