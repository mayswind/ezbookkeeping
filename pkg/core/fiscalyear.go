package core

import (
	"fmt"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// FiscalYearStart represents the fiscal year start date as a uint16 (month: high byte, day: low byte)
type FiscalYearStart uint16

// Fiscal Year Start Date Type
const (
	FISCAL_YEAR_START_DEFAULT FiscalYearStart = 0x0101 // January 1
	FISCAL_YEAR_START_MIN     FiscalYearStart = 0x0101 // January 1
	FISCAL_YEAR_START_MAX     FiscalYearStart = 0x0C1F // December 31
	FISCAL_YEAR_START_INVALID FiscalYearStart = 0xFFFF // Invalid
)

var MONTH_MAX_DAYS = []uint8{
	uint8(31), // January
	uint8(28), // February (Disallow fiscal year start on leap day)
	uint8(31), // March
	uint8(30), // April
	uint8(31), // May
	uint8(30), // June
	uint8(31), // July
	uint8(31), // August
	uint8(30), // September
	uint8(31), // October
	uint8(30), // November
	uint8(31), // December
}

// NewFiscalYearStart creates a new FiscalYearStart from month and day values
func NewFiscalYearStart(month uint8, day uint8) (FiscalYearStart, error) {
	if !isValidFiscalYearMonthDay(month, day) {
		return 0, errs.ErrFormatInvalid
	}

	return FiscalYearStart(uint16(month)<<8 | uint16(day)), nil
}

// GetMonthDay extracts the month and day from FiscalYearType
func (f FiscalYearStart) GetMonthDay() (uint8, uint8, error) {
	if f < FISCAL_YEAR_START_MIN || f > FISCAL_YEAR_START_MAX {
		return 0, 0, errs.ErrFormatInvalid
	}

	// Extract month and day (month in high byte, day in low byte)
	month := uint8(f >> 8)
	day := uint8(f & 0xFF)

	if !isValidFiscalYearMonthDay(month, day) {
		return 0, 0, errs.ErrFormatInvalid
	}

	return month, day, nil
}

// String returns a string representation of FiscalYearStart in MM/DD format
func (f FiscalYearStart) String() string {
	month, day, err := f.GetMonthDay()

	if err != nil {
		return "Invalid"
	}

	return fmt.Sprintf("%02d-%02d", month, day)
}

// isValidFiscalYearMonthDay returns whether the specified month and day is valid
func isValidFiscalYearMonthDay(month uint8, day uint8) bool {
	return uint8(1) <= month && month <= uint8(12) && uint8(1) <= day && day <= MONTH_MAX_DAYS[int(month)-1]
}

// FiscalYearFormat represents the fiscal year format as a uint8
type FiscalYearFormat uint8

// Fiscal Year Format Type Name
const (
	FISCAL_YEAR_FORMAT_DEFAULT           FiscalYearFormat = 0
	FISCAL_YEAR_FORMAT_STARTYYYY_ENDYYYY FiscalYearFormat = 1
	FISCAL_YEAR_FORMAT_STARTYYYY_ENDYY   FiscalYearFormat = 2
	FISCAL_YEAR_FORMAT_STARTYY_ENDYY     FiscalYearFormat = 3
	FISCAL_YEAR_FORMAT_ENDYYYY           FiscalYearFormat = 4
	FISCAL_YEAR_FORMAT_ENDYY             FiscalYearFormat = 5
	FISCAL_YEAR_FORMAT_INVALID           FiscalYearFormat = 255 // Invalid
)

// String returns a textual representation of the long date format enum
func (f FiscalYearFormat) String() string {
	switch f {
	case FISCAL_YEAR_FORMAT_DEFAULT:
		return "Default"
	case FISCAL_YEAR_FORMAT_STARTYYYY_ENDYYYY:
		return "StartYYYY-EndYYYY"
	case FISCAL_YEAR_FORMAT_STARTYYYY_ENDYY:
		return "StartYYYY-EndYY"
	case FISCAL_YEAR_FORMAT_STARTYY_ENDYY:
		return "StartYY-EndYY"
	case FISCAL_YEAR_FORMAT_ENDYYYY:
		return "EndYYYY"
	case FISCAL_YEAR_FORMAT_ENDYY:
		return "EndYY"
	case FISCAL_YEAR_FORMAT_INVALID:
		return "Invalid"
	default:
		return fmt.Sprintf("Invalid(%d)", int(f))
	}
}
