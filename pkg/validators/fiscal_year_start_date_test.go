package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/stretchr/testify/assert"
)

type fiscalYearStartContainer struct {
	FiscalYearStart core.FiscalYearStart `validate:"validFiscalYearStart"`
}

func TestValidateFiscalYearStart_ValidValues(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("validFiscalYearStart", ValidateFiscalYearStart)

	testCases := []struct {
		name  string
		value core.FiscalYearStart
	}{
		{"January 1", 0x0101},   // January 1
		{"December 31", 0x0C1F}, // December 31
		{"July 1", 0x0701},      // July 1
		{"April 15", 0x040F},    // April 15
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container := fiscalYearStartContainer{FiscalYearStart: tc.value}
			err := validate.Struct(container)
			assert.Nil(t, err)
		})
	}
}

func TestValidateFiscalYearStart_InvalidValues(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("validFiscalYearStart", ValidateFiscalYearStart)

	testCases := []struct {
		name  string
		value core.FiscalYearStart
	}{
		{"Zero value", 0},             // Zero value
		{"Month 0", 0x0001},           // Month 0 (invalid)
		{"Month 13", 0x0D01},          // Month 13 (invalid)
		{"Day 0", 0x0100},             // Day 0 (invalid)
		{"January 32", 0x0120},        // January 32 (invalid)
		{"February 29", 0x021D},       // February 29 (not permitted)
		{"February 30", 0x021E},       // February 30 (invalid)
		{"April 31", 0x041F},          // April 31 (invalid)
		{"June 31", 0x061F},           // June 31 (invalid)
		{"September 31", 0x091F},      // September 31 (invalid)
		{"November 32", 0x0B20},       // November 32 (invalid)
		{"Invalid month 255", 0xFF01}, // Invalid month
		{"Invalid day 255", 0x01FF},   // Invalid day
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			container := fiscalYearStartContainer{FiscalYearStart: tc.value}
			err := validate.Struct(container)
			assert.NotNil(t, err)
		})
	}
}
