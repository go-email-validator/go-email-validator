package checkifemailexist

// SMTPPresentation is a smtp presentation for check-if-email-exists
type SMTPPresentation struct {
	CanConnectSMTP bool `json:"can_connect_smtp"`
	HasFullInbox   bool `json:"has_full_inbox"`
	IsCatchAll     bool `json:"is_catch_all"`
	IsDeliverable  bool `json:"is_deliverable"`
	IsDisabled     bool `json:"is_disabled"`
}
