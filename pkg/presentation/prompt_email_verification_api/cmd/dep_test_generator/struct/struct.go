package _struct

import "github.com/go-email-validator/go-email-validator/pkg/presentation/prompt_email_verification_api"

// DepPresentationTest is used for test, see prompt_email_verification_api/dep_functional_test.go
type DepPresentationTest struct {
	Email string
	Dep   prompt_email_verification_api.DepPresentation
}
