package checkifemailexist

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
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
	d := NewDepConverterDefault()

	tests := make([]DepPresentation, 0)
	test.DepPresentations(t, &tests, "")

	// Some data or functional cannot be matched, see more nearby DepPresentation of emails
	skipEmail := hashset.New(
		// TODO problem with SMTP, CIEE think that email is not is_catch_all. Need to run and research source code on RUST
		"sewag33689@itymail.com",
		/* TODO add proxy to test
		5.7.1 Service unavailable, Client host [94.181.152.110] blocked using Spamhaus. To request removal from this list see https://www.spamhaus.org/query/ip/94.181.152.110 (AS3130). [BN8NAM12FT053.eop-nam12.prod.protection.outlook.com]
		*/
		"salestrade86@hotmail.com",
		"monicaramirezrestrepo@hotmail.com",
		// TODO CIEE banned
		"credit@mail.ru",
		// TODO need check source code
		"asdasd@tradepro.net",
		// TODO need check source code
		"y-numata@senko.ed.jp",
		"theofanisgiotis@12pm.gr",
		"theofanis.giot2is@12pm.gr",
	)

	opts := converter.NewOptions(0)
	for _, tt := range tests {
		tt := tt
		if skipEmail.Contains(tt.Input) {
			t.Logf("skipped %v", tt.Input)
			continue
		}
		t.Run(tt.Input, func(t *testing.T) {
			t.Parallel()
			email := evmail.FromString(tt.Input)

			resultValidator := validator.Validate(ev.NewInput(email))
			got := d.Convert(email, resultValidator, opts)
			gotPresenter := got.(DepPresentation)

			sort.Strings(gotPresenter.MX.Records)
			sort.Strings(tt.MX.Records)
			if !reflect.DeepEqual(got, tt) {
				t.Errorf("Convert()\n%#v, \n want\n%#v", got, tt)
			}
		})
	}
}
