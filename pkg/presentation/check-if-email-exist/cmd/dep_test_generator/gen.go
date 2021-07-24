// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/check_if_email_exist"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// TODO DRY for gen fixtures
func main() {
	var bodyBytes []byte
	var err error
	emails := converter.EmailsForTests()
	deps := make([]interface{}, len(emails))

	err = godotenv.Load()
	die(err)

	apiKey := os.Getenv("CHECK_IF_EMAIL_EXIST")

	for i, email := range emails {
		message := map[string]interface{}{
			"to_email": email,
		}

		bytesRepresentation, _ := json.Marshal(message)
		req, err := http.NewRequest(
			"POST",
			"https://ssfy.sh/amaurymartiny/reacher@2d2ce35c/check_email",
			bytes.NewBuffer(bytesRepresentation),
		)
		die(err)
		req.Header.Set("Content-Type", "application/json")
		if apiKey != "" {
			req.Header.Set("authorization", apiKey)
		}

		func() {
			resp, err := http.DefaultClient.Do(req)
			die(err)
			defer resp.Body.Close()
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			die(err)
		}()

		var dep checkifemailexist.DepPresentation
		err = json.Unmarshal(bodyBytes, &dep)
		die(err)

		if dep.Error != "" {
			panic(fmt.Sprint(email, dep.Error))
		}

		deps[i] = dep
	}

	f, err := os.Create(test.DefaultDepFixtureFile)
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
