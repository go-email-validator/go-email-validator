package mailboxvalidator

import (
	"regexp"
	"strings"
)

var hasNumberInUserNameRE = regexp.MustCompile(`.*\d.*?`)
var hasDotInUsername = regexp.MustCompile(`.*\..*?`)
var minScore = 1

func CalculateScore(presentation DepPresentation) float64 {
	score := minScore

	if presentation.IsDomain && presentation.IsSyntax {
		score += 9
		if presentation.IsSmtp {
			score += 10
		}
		if presentation.IsVerified {
			score += 40

			if presentation.IsDisposable {
				score = 30
				if presentation.IsCatchall {
					score -= 5
				}
			} else {
				if !presentation.IsFree {
					score += 39
					if presentation.IsCatchall {
						score -= 44
					} else if presentation.IsRole {
						score -= 39
					}
				} else if presentation.IsCatchall {
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
		score += 1
	}

	return float64(score) / 100.0
}
