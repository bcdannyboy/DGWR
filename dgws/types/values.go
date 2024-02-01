package types

type SingleNumber struct {
	Value             float64  `json:"value,omitempty"`
	StandardDeviation *float64 `json:"standard_deviation,omitempty"`
	Confidence        *float64 `json:"confidence,omitempty"`
}

type Range struct {
	Minimum struct {
		Value             float64  `json:"value,omitempty"`
		StandardDeviation *float64 `json:"standard_deviation,omitempty"`
		Confidence        *float64 `json:"confidence,omitempty"`
	}
	Maximum struct {
		Value             float64  `json:"value,omitempty"`
		StandardDeviation *float64 `json:"standard_deviation,omitempty"`
		Confidence        *float64 `json:"confidence,omitempty"`
	}
}

type DecomposedComponent struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Probability *struct {
		SingleNumber *SingleNumber `json:"single_number,omitempty"`
		Range        *Range        `json:"range,omitempty"`
		Decomposed   *Decomposed   `json:"decomposed,omitempty"`
	} `json:"probability,omitempty"`
	Impact *struct {
		SingleNumber *SingleNumber `json:"single_number,omitempty"`
		Range        *Range        `json:"range,omitempty"`
		Decomposed   *Decomposed   `json:"decomposed,omitempty"`
	} `json:"impact,omitempty"`
	Cost *struct {
		SingleNumber *SingleNumber `json:"single_number,omitempty"`
		Range        *Range        `json:"range,omitempty"`
		Decomposed   *Decomposed   `json:"decomposed,omitempty"`
	} `json:"cost,omitempty"`
	TimeFrame uint64
}

type Decomposed struct {
	Components []*DecomposedComponent `json:"components,omitempty"`
}
