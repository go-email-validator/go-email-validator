package mailboxvalidator

import (
	"regexp"
	"strings"
)

var hasNumberInUserNameRE = regexp.MustCompile(`.*\d.*?`)
var hasDotInUsername = regexp.MustCompile(`.*\..*?`)
var minScore = 1

// CalculateScore calculates scores for the MailboxvalidatorScore field
func CalculateScore(presentation DepPresentation) float64 {
	score := minScore

	if presentation.IsDomain.ToBool() && presentation.IsSyntax.ToBool() {
		score += 9
		if presentation.IsSMTP.ToBool() {
			score += 10
		}
		if presentation.IsVerified.ToBool() {
			score += 40

			if presentation.IsDisposable.ToBool() {
				score = 30
				if presentation.IsCatchall.ToBool() {
					score -= 5
				}
			} else {
				if !presentation.IsFree.ToBool() {
					score += 39
					if presentation.IsCatchall.ToBool() {
						score -= 44
					} else if presentation.IsRole.ToBool() {
						score -= 39
					}
				} else if presentation.IsCatchall.ToBool() {
					score -= 5
				}
			}
		}
	}
	if score < minScore {
		score = minScore
	}

	pos := strings.IndexByte(presentation.EmailAddress, '@')
	if pos == -1 {
		pos = len(presentation.EmailAddress) - 1
	}
	username := presentation.EmailAddress[:pos]

	if hasNumberInUserNameRE.MatchString(username) {
		score -= 2
	}
	if hasDotInUsername.MatchString(username) {
		score++
	}

	return float64(score) / 100.0
}
