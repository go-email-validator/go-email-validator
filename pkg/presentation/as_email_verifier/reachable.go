package as_email_verifier

type Reachable string

func (a Reachable) String() string {
	return string(a)
}

const (
	ReachableYes     Reachable = "yes"
	ReachableNo      Reachable = "no"
	ReachableUnknown Reachable = "unknown"
)

type FuncReachable func(depPresentation DepPresentation) Reachable

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
