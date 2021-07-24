package mailboxvalidator

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/test"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func TestDepConverter_Functional_Convert(t *testing.T) {
	evtests.FunctionalSkip(t)

	validator := NewDepValidator(nil)
	d := NewDepConverterForViewDefault()

	tests := detPresentersFoView(t)

	// Some data or functional cannot be matched, see more nearby DepPresentation of emails
	skipEmail := hashset.New(
		"zxczxczxc@joycasinoru", // TODO syntax is valid
		"sewag33689@itymail.com",
		"derduzikne@nedoz.com",
		"tvzamhkdc@emlhub.com",
		"admin@gmail.com",
		"salestrade86@hotmail.com",
		"monicaramirezrestrepo@hotmail.com",
		"y-numata@senko.ed.jp",
		"pr@yandex-team.ru",
		"asdasd@tradepro.net",
		"theofanisgiotis@12pm.gr",
		"theofanis.giot2is@12pm.gr",
	)

	for _, tt := range tests {
		tt := tt
		if skipEmail.Contains(tt.EmailAddress) {
			t.Logf("skipped %v", tt.EmailAddress)
			continue
		}

		t.Run(tt.EmailAddress, func(t *testing.T) {
			t.Parallel()

			email := EmailFromString(tt.EmailAddress)

			opts := converter.NewOptions(0)
			tt.TimeTaken = "0"

			resultValidator := validator.Validate(ev.NewInput(email))
			if gotResult := d.Convert(email, resultValidator, opts); !reflect.DeepEqual(gotResult, tt) {
				t.Errorf("Convert()\n%#v, \n want\n%#v", gotResult, tt)
			}
		})
	}
}

func detPresenters(t *testing.T) []DepPresentation {
	tests := make([]DepPresentation, 0)
	test.DepPresentations(t, &tests, "")

	return tests
}

func detPresentersFoView(t *testing.T) []DepPresentationForView {
	tests := make([]DepPresentationForView, 0)
	test.DepPresentations(t, &tests, "dep_view_fixture_test.json")

	return tests
}
