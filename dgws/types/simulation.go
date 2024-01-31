package types

type SimulationResults struct {
	EventID                      uint64  `json:"event_id"`
	Probability                  float64 `json:"probability"`
	ProbabilityStandardDeviation float64 `json:"probability_standard_deviation"`
	Impact                       float64 `json:"impact"`
	ImpactStandardDeviation      float64 `json:"impact_standard_deviation"`
	IsCostSaving                 bool    `json:"is_cost_saving"`
	Cost                         float64 `json:"cost"`
	CostStandardDeviation        float64 `json:"cost_standard_deviation"`
}
