package types

type Probability struct {
	ID              uint64             `json:"id"`
	Name            string             `json:"name,omitempty"`
	Description     string             `json:"description,omitempty"`
	SingleNumber    *SingleNumber      `json:"single_number,omitempty"`
	Range           *Range             `json:"range,omitempty"`
	Decomposed      *Decomposed        `json:"decomposed,omitempty"`
	DependsOnEvents []*EventDependency `json:"depends_on_events,omitempty"`
}

type Impact struct {
	ID              uint64              `json:"id"`
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
	ID              uint64              `json:"id"`
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
	AssociatedRisk        *Risk                    `json:"associated_risk,omitempty"`
	AssociatedMitigation  *Mitigation              `json:"associated_mitigation,omitempty"`
	DependsOnEvent        []*EventDependency       `json:"depends_on_event,omitempty"`
	DependsOnProbability  []*ProbabilityDependency `json:"depends_on_probability,omitempty"`
	DependsOnImpact       []*ImpactDependency      `json:"depends_on_impact,omitempty"`
	DependsOnRisk         []*RiskDependency        `json:"depends_on_risk,omitempty"`
	DependsOnCost         []*CostDependency        `json:"depends_on_cost,omitempty"`
	DependsOnMitigation   []*MitigationDependency  `json:"depends_on_mitigation,omitempty"`
}
