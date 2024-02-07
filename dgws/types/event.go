package types

const (
	Year = iota
	Quarter
	Month
	Week
	Day
	TwoYears
	FiveYears
	TenYears
)

type Probability struct {
	Name            string             `json:"name,omitempty"`
	Description     string             `json:"description,omitempty"`
	SingleNumber    *SingleNumber      `json:"single_number,omitempty"`
	Range           *Range             `json:"range,omitempty"`
	Decomposed      *Decomposed        `json:"decomposed,omitempty"`
	DependsOnEvents []*EventDependency `json:"depends_on_events,omitempty"`
}

type Impact struct {
	Name            string              `json:"name,omitempty"`
	Description     string              `json:"description,omitempty"`
	IsCostSaving    bool                `json:"is_cost_saving,omitempty"`
	SingleNumber    *SingleNumber       `json:"single_number,omitempty"`
	Range           *Range              `json:"range,omitempty"`
	Decomposed      *Decomposed         `json:"decomposed,omitempty"`
	DependsOnEvent  []*EventDependency  `json:"depends_on_events,omitempty"`
	DependsOnImpact []*ImpactDependency `json:"depends_on_impact,omitempty"`
}

type Cost struct {
	Name            string              `json:"name,omitempty"`
	Description     string              `json:"description,omitempty"`
	SingleNumber    *SingleNumber       `json:"single_number,omitempty"`
	Range           *Range              `json:"range,omitempty"`
	Decomposed      *Decomposed         `json:"decomposed,omitempty"`
	DependsOnEvent  []*EventDependency  `json:"depends_on_events,omitempty"`
	DependsOnImpact []*ImpactDependency `json:"depends_on_impact,omitempty"`
}

type Event struct {
	ID                    uint64                   `json:"id"`
	Name                  string                   `json:"name"`
	Description           string                   `json:"description,omitempty"`
	AssociatedProbability *Probability             `json:"associated_probability,omitempty"`
	AssociatedImpact      *Impact                  `json:"associated_impact,omitempty"`
	AssociatedCost        *Cost                    `json:"associated_cost,omitempty"`
	DependsOnEvent        []*EventDependency       `json:"depends_on_event,omitempty"`
	DependsOnProbability  []*ProbabilityDependency `json:"depends_on_probability,omitempty"`
	DependsOnImpact       []*ImpactDependency      `json:"depends_on_impact,omitempty"`
	DependsOnCost         []*CostDependency        `json:"depends_on_cost,omitempty"`
	Timeframe             uint64                   `json:"timeframe"`
}
