package core

import "fmt"

// IconType represents icon source type
type IconType byte

// Icon source types
const (
	ICON_TYPE_SYSTEM      IconType = 0
	ICON_TYPE_USER_CUSTOM IconType = 1
)

// String returns a textual representation of the icon type enum
func (t IconType) String() string {
	switch t {
	case ICON_TYPE_SYSTEM:
		return "System"
	case ICON_TYPE_USER_CUSTOM:
		return "User Custom"
	default:
		return fmt.Sprintf("Invalid(%d)", int(t))
	}
}

// IsValid returns whether icon type is valid
func (t IconType) IsValid() bool {
	return t == ICON_TYPE_SYSTEM || t == ICON_TYPE_USER_CUSTOM
}
