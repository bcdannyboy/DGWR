package types

// Data Exposure Risks are the risks associated with the exposure of sensitive data.
type Risk struct {
	ID             uint64             `json:"id"`
	Name           string             `json:"name`
	RiskType       string             `json:"risk_type"`
	CIACategory    string             `json:"cia_category"`
	Description    string             `json:"description,omitempty"`
	Probability    *Probability       `json:"probability,omitempty"`
	Impact         *Impact            `json:"impact,omitempty"`
	DependsOnEvent []*EventDependency `json:"depends_on_event,omitempty"`
	DependsOnRisk  []*RiskDependency  `json:"depends_on_risk,omitempty"`
}

type Mitigation struct {
	ID             uint64             `json:"id"`
	Name           string             `json:"name"`
	Description    string             `json:"description,omitempty"`
	Probability    *Probability       `json:"probability,omitempty"`
	Impact         *Impact            `json:"impact,omitempty"`
	Mitigates      []uint64           `json:"mitigates,omitempty"`
	AssociatedCost *Cost              `json:"associated_cost,omitempty"`
	DependsOnCost  []*CostDependency  `json:"depends_on_cost,omitempty"`
	DependsOnEvent []*EventDependency `json:"depends_on_event,omitempty"`
}
