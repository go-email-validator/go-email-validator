package mailboxvalidator

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/presentation_test"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func TestDepConverter_Functional_Convert(t *testing.T) {
	evtests.FunctionalSkip(t)

	validator := NewDepValidator(nil)
	d := NewDepConverterDefault()

	tests := detPresenters(t)

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
			opts := converter.NewOptions(tt.TimeTaken)

			resultValidator := validator.Validate(ev.NewInput(email))
			if gotResult := d.Convert(email, resultValidator, opts); !reflect.DeepEqual(gotResult, tt) {
				t.Errorf("Convert()\n%#v, \n want\n%#v", gotResult, tt)
			}
		})
	}
}

func detPresenters(t *testing.T) []DepPresentation {
	tests := make([]DepPresentation, 0)
	presentation_test.TestDepPresentations(t, &tests, "")

	return tests
}
