package check_if_email_exist

type FuncAvailability func(depPresentation DepPresentation) Availability

type Availability string

func (a Availability) String() string {
	return string(a)
}

const (
	Risky   Availability = "risky"
	Invalid Availability = "invalid"
	Safe    Availability = "safe"
	Unknown Availability = "unknown"
)

func CalculateAvailability(depPresentation DepPresentation) Availability {
	if depPresentation.Misc.IsDisposable ||
		depPresentation.Misc.IsRoleAccount ||
		depPresentation.SMTP.IsCatchAll ||
		depPresentation.SMTP.HasFullInbox {
		return Risky
	}

	if !depPresentation.SMTP.IsDeliverable ||
		!depPresentation.SMTP.CanConnectSmtp ||
		depPresentation.SMTP.IsDisabled {
		return Invalid
	}
	return Safe
	/*
		TODO run rust code to understand when Unknown should be used
		return Unknown
	*/
}
