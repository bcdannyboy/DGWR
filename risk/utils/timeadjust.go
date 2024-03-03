package utils

func AdjustForTime(Value float64, TimeFrame string) float64 {
	// we scale everything up  or down to a yearly basis
	switch TimeFrame {
	case "hourly":
		return Value * 8760
	case "daily":
		return Value * 365
	case "weekly":
		return Value * 52
	case "monthly":
		return Value * 12
	case "quarterly":
		return Value * 4
	case "yearly":
		return Value
	case "2years":
		return Value / 2
	case "5years":
		return Value / 5
	case "10years":
		return Value / 10
	default:
		return Value
	}
}
