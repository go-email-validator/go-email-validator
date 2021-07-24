package checkifemailexist

// FuncAvailability is an interface to calculate availability
type FuncAvailability func(depPresentation DepPresentation) Availability

// Availability is a type
type Availability string

func (a Availability) String() string {
	return string(a)
}

// Availability constants
const (
	Risky   Availability = "risky"
	Invalid Availability = "invalid"
	Safe    Availability = "safe"
	Unknown Availability = "unknown"
)

// CalculateAvailability calculates availability status
func CalculateAvailability(depPresentation DepPresentation) Availability {
	if depPresentation.Misc.IsDisposable ||
		depPresentation.Misc.IsRoleAccount ||
		depPresentation.SMTP.IsCatchAll ||
		depPresentation.SMTP.HasFullInbox {
		return Risky
	}

	if !depPresentation.SMTP.IsDeliverable ||
		!depPresentation.SMTP.CanConnectSMTP ||
		depPresentation.SMTP.IsDisabled {
		return Invalid
	}
	return Safe
	/*
		TODO run rust code to understand when Unknown should be used
		return Unknown
	*/
}
