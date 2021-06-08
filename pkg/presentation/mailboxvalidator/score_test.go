package mailboxvalidator

import (
	"github.com/emirpasic/gods/sets/hashset"
	"testing"
)

func TestCalculateScore(t *testing.T) {
	skipEmail := hashset.New(
		"zxczxczxc@joycasinoru", //TODO syntax is valid
	)

	tests := detPresenters(t)

	for _, tt := range tests {
		if skipEmail.Contains(tt.EmailAddress) {
			t.Logf("skipped %v", tt.EmailAddress)
			continue
		}

		t.Run(tt.EmailAddress, func(t *testing.T) {
			if got := CalculateScore(tt); got != tt.MailboxvalidatorScore {
				t.Errorf("CalculateScore() = %v, want %v", got, tt.MailboxvalidatorScore)
			}
		})
	}
}
