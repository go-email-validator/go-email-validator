package evtests

import (
	"flag"
	"os"
	"testing"
)

const usageFunctionalMessage = "run functional tests"

var functionalFlag = flag.Bool("functional", false, usageFunctionalMessage)
var functionalShortFlag = flag.Bool("func", false, usageFunctionalMessage)

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func FunctionalSkip(t *testing.T) {
	if !*functionalFlag && !*functionalShortFlag {
		t.Skip()
	}
}

func ToError(ret interface{}) error {
	if ret != nil {
		return ret.(error)
	}
	return nil
}
