package core

import (
	"testing"

	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestNewFiscalYearStart_ValidMonthDay(t *testing.T) {
	testCases := []struct {
		month    uint8
		day      uint8
		expected FiscalYearStart
	}{
		{1, 1, 0x0101},   // January 1
		{4, 15, 0x040F},  // April 15
		{7, 1, 0x0701},   // July 1
		{12, 31, 0x0C1F}, // December 31
	}

	for _, tc := range testCases {
		fiscal, err := NewFiscalYearStart(tc.month, tc.day)
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, fiscal)
	}
}

func TestNewFiscalYearStart_InvalidMonthDay(t *testing.T) {
	testCases := []struct {
		month uint8
		day   uint8
	}{
		{0, 1},    // Month 0 (invalid)
		{13, 1},   // Month 13 (invalid)
		{1, 0},    // Day 0 (invalid)
		{1, 32},   // Day 32 (invalid for January)
		{2, 30},   // Day 30 (invalid for February)
		{2, 29},   // Day 29 (leap day not permitted)
		{4, 31},   // Day 31 (invalid for April)
		{6, 31},   // Day 31 (invalid for June)
		{9, 31},   // Day 31 (invalid for September)
		{11, 32},  // Day 32 (invalid for November)
		{255, 15}, // Invalid month
		{5, 255},  // Invalid day
	}

	for _, tc := range testCases {
		fiscal, err := NewFiscalYearStart(tc.month, tc.day)
		assert.Equal(t, FiscalYearStart(0), fiscal)
		assert.Equal(t, errs.ErrFormatInvalid, err)
	}
}

func TestGetMonthDay_ValidFiscalYearStart(t *testing.T) {
	testCases := []struct {
		fiscalYear FiscalYearStart
		month      uint8
		day        uint8
	}{
		{0x0101, 1, 1},   // January 1st
		{0x0C1F, 12, 31}, // December 31st
		{0x0701, 7, 1},   // July 1st
		{0x040F, 4, 15},  // April 15th
	}

	for _, tc := range testCases {
		month, day, err := tc.fiscalYear.GetMonthDay()
		assert.Nil(t, err)
		assert.Equal(t, tc.month, month)
		assert.Equal(t, tc.day, day)
	}
}

func TestGetMonthDay_InvalidFiscalYearStart(t *testing.T) {
	testCases := []struct {
		fiscalYear FiscalYearStart
	}{
		{0x0000}, // 0/0 (invalid)
		{0x0D01}, // Month 13 (invalid)
		{0x0100}, // Day 0 (invalid)
		{0x0120}, // January 32 (invalid)
		{0x021D}, // February 29 (not permitted)
		{0x021E}, // February 30 (invalid)
		{0x041F}, // April 31 (invalid)
		{0x061F}, // June 31 (invalid)
		{0x091F}, // September 31 (invalid)
		{0x0B20}, // November 32 (invalid)
		{0xFF01}, // Invalid month
		{0x01FF}, // Invalid day
		{0},      // Zero value
	}

	for _, tc := range testCases {
		month, day, err := tc.fiscalYear.GetMonthDay()
		assert.Equal(t, uint8(0), month)
		assert.Equal(t, uint8(0), day)
		assert.Equal(t, errs.ErrFormatInvalid, err)
	}
}

func TestFiscalYearStart_String(t *testing.T) {
	testCases := []struct {
		fiscalYear FiscalYearStart
		expected   string
	}{
		{0x0101, "01-01"},   // January 1st
		{0x0C1F, "12-31"},   // December 31st
		{0x0701, "07-01"},   // July 1st
		{0x040F, "04-15"},   // April 15th
		{0x021D, "Invalid"}, // February 29th (leap day not permitted)
		{0x0000, "Invalid"}, // Invalid date
		{0x0D01, "Invalid"}, // Invalid month
		{0x0120, "Invalid"}, // Invalid day
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, tc.fiscalYear.String())
	}
}

func TestFiscalYearStartConstants(t *testing.T) {
	assert.Equal(t, FiscalYearStart(0x0000), FISCAL_YEAR_START_INVALID)
	assert.Equal(t, FiscalYearStart(0x0101), FISCAL_YEAR_START_DEFAULT)
	assert.Equal(t, FiscalYearStart(0x0101), FISCAL_YEAR_START_MIN)
	assert.Equal(t, FiscalYearStart(0x0C1F), FISCAL_YEAR_START_MAX)
}
