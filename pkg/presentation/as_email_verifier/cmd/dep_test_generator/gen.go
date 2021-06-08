package main

import (
	"encoding/json"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-ev-presenters/pkg/presentation/converter"
	"github.com/go-email-validator/go-ev-presenters/pkg/presentation/presentation_test"
	"log"
	"os"
)

// TODO DRY for gen fixtures
func main() {
	var err error
	emails := converter.EmailsForTests()
	deps := make([]interface{}, len(emails))

	verifier := emailverifier.NewVerifier().
		EnableGravatarCheck().
		EnableSMTPCheck()
	// .EnableDomainSuggest()

	incorrectGravatar := hashset.New(
		"amazedfuckporno@gmail.com",
		"theofanis.giot2is@12pm.gr",
		"asdasd@tradepro.net",
		"credit@mail.ru",
		"admin@gmail.com",
	)

	for i, email := range emails {
		verifyResult, _ := verifier.Verify(email)

		if incorrectGravatar.Contains(email) && verifyResult.Gravatar != nil {
			verifyResult.Gravatar.HasGravatar = false
			verifyResult.Gravatar.GravatarUrl = ""
		}

		deps[i] = verifyResult
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
