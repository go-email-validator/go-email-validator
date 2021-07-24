package promptemailverificationapi

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/test"
	"reflect"
	"sort"
	"testing"
)

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func TestDepConverter_Convert(t *testing.T) {
	evtests.FunctionalSkip(t)

	validator := NewDepValidator(nil)
	d := NewDepConverter()
	tests := make([]DepPresenterTest, 0)
	test.DepPresentations(t, &tests, "")

	// Some data or functional cannot be matched, see more nearby DepPresentation of emails
	skipEmail := hashset.New(
		// Banned on disposable domain
		"sewag33689@itymail.com",
		// Banned on disposable domain
		"sewag33689@itymail.com",
		// Banned on disposable domain
		"asdasd@tradepro.net",
		// Banned on disposable domain
		"tvzamhkdc@emlhub.com",
		// Banned
		"credit@mail.ru",
		// Banned
		"salestrade86@hotmail.com",
		// Banned
		"monicaramirezrestrepo@hotmail.com",
		// TODO Cannot connect from my hosts pc
		"y-numata@senko.ed.jp",
		"theofanisgiotis@12pm.gr",
		"theofanis.giot2is@12pm.gr",
	)

	opts := converter.NewOptions(0)
	for _, tt := range tests {
		tt := tt
		if skipEmail.Contains(tt.Email) {
			t.Logf("skipped %v", tt.Email)
			continue
		}

		t.Run(tt.Email, func(t *testing.T) {
			t.Parallel()

			email := EmailFromString(tt.Email)

			resultValidator := validator.Validate(ev.NewInput(email))
			got := d.Convert(email, resultValidator, opts)
			gotPresenter := got.(DepPresentation)

			sort.Strings(gotPresenter.MxRecords.Records)
			sort.Strings(tt.Dep.MxRecords.Records)
			if !reflect.DeepEqual(got, tt.Dep) {
				t.Errorf("Convert()\n%#v, \n want\n%#v", got, tt.Dep)
			}
		})
	}
}

// see prompt_email_verification_api/cmd/dep_test_generator/gen.go
type DepPresenterTest struct {
	Email string
	Dep   DepPresentation
}
