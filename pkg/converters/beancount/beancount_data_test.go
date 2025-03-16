package beancount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeancountAccount_IsOpeningBalanceEquityAccount_True(t *testing.T) {
	account := beancountAccount{
		accountType: beancountEquityAccountType,
		name:        "Equity:Opening-Balances",
	}
	assert.True(t, account.isOpeningBalanceEquityAccount())

	account = beancountAccount{
		accountType: beancountEquityAccountType,
		name:        "E:Opening-Balances",
	}
	assert.True(t, account.isOpeningBalanceEquityAccount())
}

func TestBeancountAccount_IsOpeningBalanceEquityAccount_False(t *testing.T) {
	account := beancountAccount{
		accountType: beancountAssetsAccountType,
		name:        "Equity:Opening-Balances",
	}
	assert.False(t, account.isOpeningBalanceEquityAccount())

	account = beancountAccount{
		accountType: beancountEquityAccountType,
		name:        "Opening-Balances",
	}
	assert.False(t, account.isOpeningBalanceEquityAccount())

	account = beancountAccount{
		accountType: beancountEquityAccountType,
		name:        "Equity:Other",
	}
	assert.False(t, account.isOpeningBalanceEquityAccount())
}
