package view

type SMTPView struct {
	CanConnectSmtp bool `json:"can_connect_smtp"`
	HasFullInbox   bool `json:"has_full_inbox"`
	IsCatchAll     bool `json:"is_catch_all"`
	IsDeliverable  bool `json:"is_deliverable"`
	IsDisabled     bool `json:"is_disabled"`
}

/*
TODO
var NewSMTPView = func(_ email.EmailAddressInterface, result ev.ValidationResultInterface) ViewInterface {
	view := SMTPView{
		CanConnectSmtp
		HasFullInbox
		IsCatchAll
		IsDeliverable
		IsDisabled,
	}

	return ViewInterface()
}*/
