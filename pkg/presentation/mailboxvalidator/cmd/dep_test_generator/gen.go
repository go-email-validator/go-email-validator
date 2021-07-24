// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/mailboxvalidator"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/mailboxvalidator/addition"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/test"
	"github.com/joho/godotenv"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/sets"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	var bodyBytes []byte
	var err error
	emails := converter.EmailsForTests()
	emailsHasResponseError := sets.NewString("")
	var deps, depsForView []interface{}

	err = godotenv.Load()
	die(err)

	apiKey := os.Getenv("MAIL_BOX_VALIDATOR_API")
	if apiKey == "" {
		panic("MAIL_BOX_VALIDATOR_API should be set")
	}

	for _, email := range emails {
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
		err = json.Unmarshal(bodyBytes, &depForView)
		die(err)

		if depForView.ErrorCode != "" && !emailsHasResponseError.Has(email) {
			panic(fmt.Sprint(email, depForView.ErrorMessage))
		}
		depForView = prepareDepView(depForView)
		depsForView = append(depsForView, depForView)

		dep := createDep(depForView)
		deps = append(deps, dep)
	}

	write(test.DefaultDepFixtureFile, deps)
	write(addition.DepFixtureForViewFile, depsForView)
}

// TODO add null string
func prepareDepView(depForView mailboxvalidator.DepPresentationForView) mailboxvalidator.DepPresentationForView {
	depForView.IsFree = emptyBoolToFalse(depForView.IsFree)
	depForView.IsSyntax = emptyBoolToFalse(depForView.IsSyntax)
	depForView.IsDomain = emptyBoolToFalse(depForView.IsDomain)
	depForView.IsSMTP = emptyBoolToFalse(depForView.IsSMTP)
	depForView.IsVerified = emptyBoolToFalse(depForView.IsVerified)
	depForView.IsServerDown = emptyBoolToFalse(depForView.IsServerDown)
	depForView.IsGreylisted = emptyBoolToFalse(depForView.IsGreylisted)
	depForView.IsDisposable = emptyBoolToFalse(depForView.IsDisposable)
	depForView.IsSuppressed = emptyBoolToFalse(depForView.IsSuppressed)
	depForView.IsRole = emptyBoolToFalse(depForView.IsRole)
	depForView.IsHighRisk = emptyBoolToFalse(depForView.IsHighRisk)
	depForView.Status = emptyBoolToFalse(depForView.Status)

	depForView.CreditsAvailable = 4294967295

	return depForView
}

func createDep(depForView mailboxvalidator.DepPresentationForView) mailboxvalidator.DepPresentation {
	var dep mailboxvalidator.DepPresentation

	dep.EmailAddress = depForView.EmailAddress
	dep.Domain = depForView.Domain
	dep.IsFree = mailboxvalidator.ToBool(depForView.IsFree)
	dep.IsSyntax = mailboxvalidator.ToBool(depForView.IsSyntax)
	dep.IsDomain = mailboxvalidator.ToBool(depForView.IsDomain)
	dep.IsSMTP = mailboxvalidator.ToBool(depForView.IsSMTP)
	dep.IsVerified = mailboxvalidator.ToBool(depForView.IsVerified)
	dep.IsServerDown = mailboxvalidator.ToBool(depForView.IsServerDown)
	dep.IsGreylisted = mailboxvalidator.ToBool(depForView.IsGreylisted)
	dep.IsDisposable = mailboxvalidator.ToBool(depForView.IsDisposable)
	dep.IsSuppressed = mailboxvalidator.ToBool(depForView.IsSuppressed)
	dep.IsRole = mailboxvalidator.ToBool(depForView.IsRole)
	dep.IsHighRisk = mailboxvalidator.ToBool(depForView.IsHighRisk)
	dep.IsCatchall = mailboxvalidator.ToBool(depForView.IsCatchall)

	dep.MailboxvalidatorScore, _ = strconv.ParseFloat(depForView.MailboxvalidatorScore, 64)

	dep.TimeTaken = 568300000
	if depForView.TimeTaken == 0 {
		dep.TimeTaken = 0
	}

	dep.Status = mailboxvalidator.ToBool(depForView.Status)
	dep.CreditsAvailable = depForView.CreditsAvailable
	dep.ErrorCode = depForView.ErrorCode
	dep.ErrorMessage = depForView.ErrorMessage

	return dep
}

func emptyBoolToFalse(boolean string) string {
	if boolean == "-" {
		return "False"
	}

	return boolean
}

func write(filepath string, deps interface{}) {
	file, err := os.Create(filepath)
	die(err)
	defer file.Close()

	bytes, err := json.MarshalIndent(deps, "", "  ")
	die(err)
	_, err = file.Write(bytes)
	die(err)
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
