package core

import "fmt"

// CalendarDisplayType represents calendar display type
type CalendarDisplayType byte

// Calendar Display Type
const (
	CALENDAR_DISPLAY_TYPE_DEFAULT   CalendarDisplayType = 0
	CALENDAR_DISPLAY_TYPE_GREGORAIN CalendarDisplayType = 1
	CALENDAR_DISPLAY_TYPE_INVALID   CalendarDisplayType = 255
)

// String returns a textual representation of the calendar display type enum
func (f CalendarDisplayType) String() string {
	switch f {
	case CALENDAR_DISPLAY_TYPE_DEFAULT:
		return "Default"
	case CALENDAR_DISPLAY_TYPE_GREGORAIN:
		return "Gregorian"
	case CALENDAR_DISPLAY_TYPE_INVALID:
		return "Invalid"
	default:
		return fmt.Sprintf("Invalid(%d)", int(f))
	}
}

// DateDisplayType represents date display type
type DateDisplayType byte

// Date Display Type
const (
	DATE_DISPLAY_TYPE_DEFAULT   DateDisplayType = 0
	DATE_DISPLAY_TYPE_GREGORAIN DateDisplayType = 1
	DATE_DISPLAY_TYPE_BUDDHIST  DateDisplayType = 2
	DATE_DISPLAY_TYPE_INVALID   DateDisplayType = 255
)

// String returns a textual representation of the date display type enum
func (f DateDisplayType) String() string {
	switch f {
	case DATE_DISPLAY_TYPE_DEFAULT:
		return "Default"
	case DATE_DISPLAY_TYPE_GREGORAIN:
		return "Gregorian"
	case DATE_DISPLAY_TYPE_BUDDHIST:
		return "Buddhist"
	case DATE_DISPLAY_TYPE_INVALID:
		return "Invalid"
	default:
		return fmt.Sprintf("Invalid(%d)", int(f))
	}
}
