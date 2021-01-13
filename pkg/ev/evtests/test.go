package evtests

import (
	"flag"
	"os"
	"testing"
)

const usageFunctionalMessage = "run functional tests"

var functionalFlag = flag.Bool("functional", false, usageFunctionalMessage)
var functionalShortFlag = flag.Bool("func", false, usageFunctionalMessage)

// TestMain initialize flags
func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

// FunctionalSkip skips test if "-func" or "-functional" was not put in `go test`
func FunctionalSkip(t *testing.T) {
	if !*functionalFlag && !*functionalShortFlag {
		t.Skip()
	}
}

// ToError casts interface{} to error
func ToError(ret interface{}) error {
	if ret != nil {
		return ret.(error)
	}
	return nil
}
