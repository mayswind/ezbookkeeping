package core

import (
	"fmt"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

// FiscalYearStart represents the fiscal year start date as a uint16 (month: high byte, day: low byte)
type FiscalYearStart uint16

// Fiscal Year Start Date Type
const (
	FISCAL_YEAR_START_INVALID FiscalYearStart = 0x0000 // Invalid
	FISCAL_YEAR_START_DEFAULT FiscalYearStart = 0x0101 // January 1
	FISCAL_YEAR_START_MIN     FiscalYearStart = 0x0101 // January 1
	FISCAL_YEAR_START_MAX     FiscalYearStart = 0x0C1F // December 31
)

// NewFiscalYearStart creates a new FiscalYearStart from month and day values
func NewFiscalYearStart(month uint8, day uint8) (FiscalYearStart, error) {
	month, day, err := validateMonthDay(month, day)
	if err != nil {
		return 0, err
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

	return validateMonthDay(month, day)
}

// String returns a string representation of FiscalYearStart in MM/DD format
func (f FiscalYearStart) String() string {
	month, day, err := f.GetMonthDay()
	if err != nil {
		return "Invalid"
	}
	return fmt.Sprintf("%02d-%02d", month, day)
}

// validateMonthDay validates a month and day and returns them if valid
func validateMonthDay(month uint8, day uint8) (uint8, uint8, error) {
	if month < 1 || month > 12 || day < 1 {
		return 0, 0, errs.ErrFormatInvalid
	}

	maxDays := uint8(31)
	switch month {
	case 1, 3, 5, 7, 8, 10, 12: // January, March, May, July, August, October, December
		maxDays = 31
	case 4, 6, 9, 11: // April, June, September, November
		maxDays = 30
	case 2: // February
		maxDays = 28 // Disallow fiscal year start on leap day
	}

	if day > maxDays {
		return 0, 0, errs.ErrFormatInvalid
	}

	return month, day, nil
}
