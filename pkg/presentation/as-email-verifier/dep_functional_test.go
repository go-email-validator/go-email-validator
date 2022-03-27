package asemailverifier

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/test"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func TestDepConverter_Convert(t *testing.T) {
	evtests.FunctionalSkip(t)

	validator := NewDepValidator(nil)
	d := NewDepConverterDefault()

	tests := make([]DepPresentation, 0)
	test.DepPresentations(t, &tests, "")
	// Some data or functional cannot be matched, see more nearby DepPresentation of emails
	skipEmail := hashset.New(
		// todo disposable verification
		"derduzikne@nedoz.com",
		// TODO change presenter, if there is error before mail stage then smtp is nil
		"salestrade86@hotmail.com",
		"monicaramirezrestrepo@hotmail.com",
		"theofanisgiotis@12pm.gr",
		"theofanis.giot2is@12pm.gr",
		"admin@huntgear.ru",
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
			email := evmail.FromString(tt.Email)

			resultValidator := validator.Validate(ev.NewInput(email))
			got := d.Convert(email, resultValidator, opts)

			if !reflect.DeepEqual(got, tt) {
				t.Errorf("Convert()\n%#v, \n want\n%#v", got, tt)
			}
		})
	}
}
