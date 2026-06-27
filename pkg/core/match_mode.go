package core

import "fmt"

// MatchMode represents the type of match mode used for string comparisons
type MatchMode byte

// Match Mode
const (
	MATCH_MODE_DEFAULT     MatchMode = 0
	MATCH_MODE_IGNORE_CASE MatchMode = 1
)

// String returns a textual representation of the match mode enum
func (f MatchMode) String() string {
	switch f {
	case MATCH_MODE_DEFAULT:
		return "Default"
	case MATCH_MODE_IGNORE_CASE:
		return "Ignore Case"
	default:
		return fmt.Sprintf("Invalid(%d)", int(f))
	}
}
