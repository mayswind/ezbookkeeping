package models

import "fmt"

// CurrencyDisplayType represents the display type of amount with currency
type CurrencyDisplayType byte

// Currency Display Type
const (
	CURRENCY_DISPLAY_TYPE_DEFAULT                            CurrencyDisplayType = 0
	CURRENCY_DISPLAY_TYPE_NONE                               CurrencyDisplayType = 1
	CURRENCY_DISPLAY_TYPE_SYMBOL_BEFORE_AMOUNT               CurrencyDisplayType = 2
	CURRENCY_DISPLAY_TYPE_SYMBOL_AFTER_AMOUNT                CurrencyDisplayType = 3
	CURRENCY_DISPLAY_TYPE_SYMBOL_BEFORE_AMOUNT_WITHOUT_SPACE CurrencyDisplayType = 4
	CURRENCY_DISPLAY_TYPE_SYMBOL_AFTER_AMOUNT_WITHOUT_SPACE  CurrencyDisplayType = 5
	CURRENCY_DISPLAY_TYPE_CODE_BEFORE_AMOUNT                 CurrencyDisplayType = 6
	CURRENCY_DISPLAY_TYPE_CODE_AFTER_AMOUNT                  CurrencyDisplayType = 7
	CURRENCY_DISPLAY_TYPE_NAME_BEFORE_AMOUNT                 CurrencyDisplayType = 8
	CURRENCY_DISPLAY_TYPE_NAME_AFTER_AMOUNT                  CurrencyDisplayType = 9
	CURRENCY_DISPLAY_TYPE_INVALID                            CurrencyDisplayType = 255
)

// String returns a textual representation of the currency display type enum
func (d CurrencyDisplayType) String() string {
	switch d {
	case CURRENCY_DISPLAY_TYPE_DEFAULT:
		return "Default"
	case CURRENCY_DISPLAY_TYPE_NONE:
		return "None"
	case CURRENCY_DISPLAY_TYPE_SYMBOL_BEFORE_AMOUNT:
		return "Symbol Before Amount"
	case CURRENCY_DISPLAY_TYPE_SYMBOL_AFTER_AMOUNT:
		return "Symbol After Amount"
	case CURRENCY_DISPLAY_TYPE_SYMBOL_BEFORE_AMOUNT_WITHOUT_SPACE:
		return "Symbol Before Amount Without Space"
	case CURRENCY_DISPLAY_TYPE_SYMBOL_AFTER_AMOUNT_WITHOUT_SPACE:
		return "Symbol After Amount Without Space"
	case CURRENCY_DISPLAY_TYPE_CODE_BEFORE_AMOUNT:
		return "Code Before Amount"
	case CURRENCY_DISPLAY_TYPE_CODE_AFTER_AMOUNT:
		return "Code After Amount"
	case CURRENCY_DISPLAY_TYPE_NAME_BEFORE_AMOUNT:
		return "Name Before Amount"
	case CURRENCY_DISPLAY_TYPE_NAME_AFTER_AMOUNT:
		return "Name After Amount"
	case CURRENCY_DISPLAY_TYPE_INVALID:
		return "Invalid"
	default:
		return fmt.Sprintf("Invalid(%d)", int(d))
	}
}
