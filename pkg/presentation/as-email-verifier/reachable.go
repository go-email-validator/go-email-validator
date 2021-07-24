package asemailverifier

// Reachable is a type
type Reachable string

func (a Reachable) String() string {
	return string(a)
}

// Reachable constants
const (
	ReachableYes     Reachable = "yes"
	ReachableNo      Reachable = "no"
	ReachableUnknown Reachable = "unknown"
)

// FuncReachable is an interface to calculate Reachable
type FuncReachable func(depPresentation DepPresentation) Reachable

// CalculateReachable returns Reachable status
func CalculateReachable(depPresentation DepPresentation) Reachable {
	smtp := depPresentation.SMTP

	if smtp == nil {
		return ReachableUnknown
	}

	if smtp.Deliverable {
		return ReachableYes
	}
	if smtp.CatchAll {
		return ReachableUnknown
	}
	return ReachableNo
}
