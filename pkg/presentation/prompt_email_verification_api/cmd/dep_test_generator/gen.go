// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/prompt_email_verification_api/cmd/dep_test_generator/struct"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var bodyBytes []byte
	var err error
	emails := converter.EmailsForTests()
	deps := make([]interface{}, len(emails))

	err = godotenv.Load()
	die(err)

	apiKey := os.Getenv("PROMPT_EMAIL_VERIFICATION_API")
	if apiKey == "" {
		panic("PROMPT_EMAIL_VERIFICATION_API should be set")
	}

	for i, email := range emails {
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("https://api.promptapi.com/email_verification/%s", email),
			nil,
		)
		die(err)
		req.Header.Set("apikey", apiKey)

		func() {
			resp, err := http.DefaultClient.Do(req)
			die(err)
			defer resp.Body.Close()
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			die(err)
		}()

		var depTest _struct.DepPresentationTest

		if strings.Contains(dep.Message, "API rate limit") {
			panic(fmt.Sprint(email, dep.Message))
		}
		depTest = _struct.DepPresentationTest{Email: email, Dep: dep}

		deps[i] = depTest
	}

	f, err := os.Create(presentation_test.DefaultDepFixtureFile)
	die(err)
	defer f.Close()

	bytes, err := json.MarshalIndent(deps, "", "  ")
	die(err)
	_, err = f.Write(bytes)
	die(err)
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
