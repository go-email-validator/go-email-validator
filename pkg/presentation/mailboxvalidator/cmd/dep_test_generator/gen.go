// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/mailboxvalidator"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/mailboxvalidator/addition"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/presentation_test"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	var bodyBytes []byte
	var err error
	emails := converter.EmailsForTests()
	deps := make([]interface{}, len(emails))
	depsForView := make([]interface{}, len(emails))

	err = godotenv.Load()
	die(err)

	apiKey := os.Getenv("MAIL_BOX_VALIDATOR_API")
	if apiKey == "" {
		panic("MAIL_BOX_VALIDATOR_API should be set")
	}

	for i, email := range emails {
		req, err := http.NewRequest(
			"GET",
			"https://api.mailboxvalidator.com/v1/validation/single?email="+url.QueryEscape(email)+"&key="+url.QueryEscape(apiKey),
			nil,
		)
		die(err)

		func() {
			resp, err := http.DefaultClient.Do(req)
			die(err)
			defer resp.Body.Close()
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			die(err)
		}()

		var depForView mailboxvalidator.DepPresentationForView
		var depFor mailboxvalidator.DepPresentation
		err = json.Unmarshal(bodyBytes, &depForView)
		die(err)

		depFor.IsFree = mailboxvalidator.ToBool(depForView.IsFree)
		depFor.IsSyntax = mailboxvalidator.ToBool(depForView.IsSyntax)
		depFor.IsDomain = mailboxvalidator.ToBool(depForView.IsDomain)
		depFor.IsSmtp = mailboxvalidator.ToBool(depForView.IsSmtp)
		depFor.IsVerified = mailboxvalidator.ToBool(depForView.IsVerified)
		depFor.IsServerDown = mailboxvalidator.ToBool(depForView.IsServerDown)
		depFor.IsGreylisted = mailboxvalidator.ToBool(depForView.IsGreylisted)
		depFor.IsDisposable = mailboxvalidator.ToBool(depForView.IsDisposable)
		depFor.IsSuppressed = mailboxvalidator.ToBool(depForView.IsSuppressed)
		depFor.IsRole = mailboxvalidator.ToBool(depForView.IsRole)
		depFor.IsHighRisk = mailboxvalidator.ToBool(depForView.IsHighRisk)
		depFor.Status = mailboxvalidator.ToBool(depForView.Status)

		if depForView.ErrorCode != "" {
			panic(fmt.Sprint(email, depForView.ErrorMessage))
		}

		depsForView[i] = depForView
	}

	f, err := os.Create(presentation_test.DefaultDepFixtureFile)
	die(err)
	defer f.Close()

	bytes, err := json.MarshalIndent(deps, "", "  ")
	die(err)
	_, err = f.Write(bytes)
	die(err)

	fForView, err := os.Create(addition.DepFixtureForViewFile)
	die(err)
	defer fForView.Close()

	bytes, err := json.MarshalIndent(depsForView, "", "  ")
	die(err)
	_, err = f.Write(bytes)
	die(err)
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
