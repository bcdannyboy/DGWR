package types

const (
	Happens = iota
	DoesNotHappen
	Exists
	DoesNotExist
	In
	Out
	Has
	HasNot
	EQ
	NEQ
	LT
	GT
	LTE
	GTE
)

type EventDependency struct {
	ID               uint64 `json:"id"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Type             uint64 `json:"type"` // happens, does not happen
	DependentEventID uint64 `json:"dependent_event_id"`
}

type ProbabilityDependency struct {
	Name             string        `json:"name,omitempty"`
	Description      string        `json:"description,omitempty"`
	Type             uint64        `json:"type"` // ==, !=, <, >, <=, >=, in, out, has, has not
	DependentEventID *uint64       `json:"dependent_event_id,omitempty"`
	SingleValue      *SingleNumber `json:"single_value,omitempty"`
	Range            *Range        `json:"range,omitempty"`
	Decomposed       *Decomposed   `json:"decomposed,omitempty"`
}

type ImpactDependency struct {
	Name             string        `json:"name,omitempty"`
	Description      string        `json:"description,omitempty"`
	Type             uint64        `json:"type"` // ==, !=, <, >, <=, >=, in, out, has, has not
	DependentEventID *uint64       `json:"dependent_event_id,omitempty"`
	SingleValue      *SingleNumber `json:"single_value,omitempty"`
	Range            *Range        `json:"range,omitempty"`
	Decomposed       *Decomposed   `json:"decomposed,omitempty"`
}

type CostDependency struct {
	Name                        string        `json:"name,omitempty"`
	Description                 string        `json:"description,omitempty"`
	Type                        uint64        `json:"type"` // ==, !=, <, >, <=, >=, in, out, has, has not, exists, does not exist
	DependentEventID            *uint64       `json:"dependent_event_id,omitempty"`
	DependentMitigationOrRiskID *uint64       `json:"dependent_mitigation_or_risk_id,omitempty"`
	SingleValue                 *SingleNumber `json:"single_value,omitempty"`
	Range                       *Range        `json:"range,omitempty"`
	Decomposed                  *Decomposed   `json:"decomposed,omitempty"`
}

type RiskDependency struct {
	Name             string  `json:"name,omitempty"`
	Description      string  `json:"description,omitempty"`
	Type             uint64  `json:"type"`  // exists, does not exist
	DependentRiskID  *uint64 `json:"value"` // hash of the risk that this risk depends on
	DependentEventID *uint64 `json:"dependent_event_id,omitempty"`
}

type MitigationDependency struct {
	Name                        string  `json:"name,omitempty"`
	Description                 string  `json:"description,omitempty"`
	Type                        uint64  `json:"type"`  // ==, !=, <, >, <=, >=, in, out, has, has not, exists, does not exist
	DependentMitigationOrRiskID *uint64 `json:"value"` // hash of the mitigation or risk that this mitigation depends on
	DependentEventID            *uint64 `json:"dependent_event_id,omitempty"`
}
