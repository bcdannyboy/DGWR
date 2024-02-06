package types

type SingleNumber struct {
	Value             float64  `json:"value,omitempty"`
	StandardDeviation *float64 `json:"standard_deviation,omitempty"`
	Confidence        *float64 `json:"confidence,omitempty"`
}

type Minimum struct {
	Value             float64  `json:"value,omitempty"`
	StandardDeviation *float64 `json:"standard_deviation,omitempty"`
	Confidence        *float64 `json:"confidence,omitempty"`
}

type Maximum struct {
	Value             float64  `json:"value,omitempty"`
	StandardDeviation *float64 `json:"standard_deviation,omitempty"`
	Confidence        *float64 `json:"confidence,omitempty"`
}

type Range struct {
	Minimum Minimum
	Maximum Maximum
}

type DecomposedItem struct {
	SingleNumber *SingleNumber `json:"single_number,omitempty"`
	Range        *Range        `json:"range,omitempty"`
	Decomposed   *Decomposed   `json:"decomposed,omitempty"`
}

type DecomposedComponent struct {
	ComponentID uint64          `json:"component_id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Probability *DecomposedItem `json:"probability,omitempty"`
	Impact      *DecomposedItem `json:"impact,omitempty"`
	Cost        *DecomposedItem `json:"cost,omitempty"`
	TimeFrame   uint64
}

type Decomposed struct {
	Components []*DecomposedComponent `json:"components,omitempty"`
}
