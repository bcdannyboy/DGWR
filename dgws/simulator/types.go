package simulator

import "github.com/bcdannyboy/montecargo/dgws/types"

type Component struct {
	SingleNumber *types.SingleNumber
	Range        *types.Range
	Decomposed   *types.Decomposed
}

const (
	ProbabilityAttribute = iota
	ImpactAttribute
	CostAttribute
)
