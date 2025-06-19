package beancount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeancountAccount_IsOpeningBalanceEquityAccount_True(t *testing.T) {
	account := beancountAccount{
		AccountType: beancountEquityAccountType,
		Name:        "Equity:Opening-Balances",
	}
	assert.True(t, account.isOpeningBalanceEquityAccount())

	account = beancountAccount{
		AccountType: beancountEquityAccountType,
		Name:        "E:Opening-Balances",
	}
	assert.True(t, account.isOpeningBalanceEquityAccount())
}

func TestBeancountAccount_IsOpeningBalanceEquityAccount_False(t *testing.T) {
	account := beancountAccount{
		AccountType: beancountAssetsAccountType,
		Name:        "Equity:Opening-Balances",
	}
	assert.False(t, account.isOpeningBalanceEquityAccount())

	account = beancountAccount{
		AccountType: beancountEquityAccountType,
		Name:        "Opening-Balances",
	}
	assert.False(t, account.isOpeningBalanceEquityAccount())

	account = beancountAccount{
		AccountType: beancountEquityAccountType,
		Name:        "Equity:Other",
	}
	assert.False(t, account.isOpeningBalanceEquityAccount())
}
