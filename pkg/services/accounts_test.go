package services

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func TestGetAccountMapByList_EmptyList(t *testing.T) {
	accounts := make([]*models.Account, 0)
	actualAccountMap := Accounts.GetAccountMapByList(accounts)

	assert.NotNil(t, actualAccountMap)
	assert.Equal(t, 0, len(actualAccountMap))
}

func TestGetAccountMapByList_MultipleList(t *testing.T) {
	accounts := []*models.Account{
		{
			AccountId: 1001,
			Name:      "Cash Account",
			Category:  models.ACCOUNT_CATEGORY_CASH,
		},
		{
			AccountId: 1002,
			Name:      "Checking Account",
			Category:  models.ACCOUNT_CATEGORY_CHECKING_ACCOUNT,
		},
		{
			AccountId: 1003,
			Name:      "Credit Card",
			Category:  models.ACCOUNT_CATEGORY_CREDIT_CARD,
		},
	}
	actualAccountMap := Accounts.GetAccountMapByList(accounts)

	assert.Equal(t, 3, len(actualAccountMap))
	assert.Contains(t, actualAccountMap, int64(1001))
	assert.Contains(t, actualAccountMap, int64(1002))
	assert.Contains(t, actualAccountMap, int64(1003))
	assert.Equal(t, "Cash Account", actualAccountMap[1001].Name)
	assert.Equal(t, "Checking Account", actualAccountMap[1002].Name)
	assert.Equal(t, "Credit Card", actualAccountMap[1003].Name)
}

func TestGetVisibleAccountNameMapByList_EmptyList(t *testing.T) {
	accounts := make([]*models.Account, 0)
	actualAccountMap := Accounts.GetVisibleAccountNameMapByList(accounts)

	assert.NotNil(t, actualAccountMap)
	assert.Equal(t, 0, len(actualAccountMap))
}

func TestGetVisibleAccountNameMapByList_WithHiddenAccount(t *testing.T) {
	accounts := []*models.Account{
		{
			AccountId: 1001,
			Name:      "Visible Account",
			Type:      models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
			Hidden:    false,
		},
		{
			AccountId: 1002,
			Name:      "Hidden Account",
			Type:      models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
			Hidden:    true,
		},
	}
	actualAccountMap := Accounts.GetVisibleAccountNameMapByList(accounts)

	assert.Equal(t, 1, len(actualAccountMap))
	assert.Contains(t, actualAccountMap, "Visible Account")
	assert.NotContains(t, actualAccountMap, "Hidden Account")
}

func TestGetVisibleAccountNameMapByList_WithParentAccount(t *testing.T) {
	accounts := []*models.Account{
		{
			AccountId: 1001,
			Name:      "Single Account",
			Type:      models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
			Hidden:    false,
		},
		{
			AccountId: 1002,
			Name:      "Multi Sub Accounts",
			Type:      models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS,
			Hidden:    false,
		},
	}
	actualAccountMap := Accounts.GetVisibleAccountNameMapByList(accounts)
	assert.Equal(t, 1, len(actualAccountMap))
	assert.Contains(t, actualAccountMap, "Single Account")
	assert.NotContains(t, actualAccountMap, "Multi Sub Accounts")
}

func TestGetAccountNames_EmptyList(t *testing.T) {
	accounts := make([]*models.Account, 0)
	actualAccountMap := Accounts.GetAccountNames(accounts)

	assert.NotNil(t, actualAccountMap)
	assert.Equal(t, 0, len(actualAccountMap))
}

func TestGetAccountNames_MultipleList(t *testing.T) {
	accounts := []*models.Account{
		{
			AccountId: 1001,
			Name:      "Cash Account",
		},
		{
			AccountId: 1002,
			Name:      "Checking Account",
		},
		{
			AccountId: 1003,
			Name:      "Credit Card",
		},
	}
	actualAccountMap := Accounts.GetAccountNames(accounts)

	assert.Equal(t, 3, len(actualAccountMap))
	assert.Equal(t, "Cash Account", actualAccountMap[0])
	assert.Equal(t, "Checking Account", actualAccountMap[1])
	assert.Equal(t, "Credit Card", actualAccountMap[2])
}

func TestGetAccountOrSubAccountIdsByAccountName_EmptyList(t *testing.T) {
	accounts := make([]*models.Account, 0)
	actualAccountMap := Accounts.GetAccountOrSubAccountIdsByAccountName(accounts, "Test Account")

	assert.NotNil(t, actualAccountMap)
	assert.Equal(t, 0, len(actualAccountMap))
}

func TestGetAccountOrSubAccountIdsByAccountName_NotMatch(t *testing.T) {
	accounts := []*models.Account{
		{
			AccountId: 1001,
			Name:      "Cash Account",
			Type:      models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
		},
	}
	actualAccountMap := Accounts.GetAccountOrSubAccountIdsByAccountName(accounts, "Non-existent Account")

	assert.NotNil(t, actualAccountMap)
	assert.Equal(t, 0, len(actualAccountMap))
}

func TestGetAccountOrSubAccountIdsByAccountName_MatchSingle(t *testing.T) {
	accounts := []*models.Account{
		{
			AccountId: 1001,
			Name:      "Cash Account",
			Type:      models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
		},
	}
	actualAccountMap := Accounts.GetAccountOrSubAccountIdsByAccountName(accounts, "Cash Account")
	assert.Equal(t, 1, len(actualAccountMap))
	assert.Contains(t, actualAccountMap, int64(1001))
}

func TestGetAccountOrSubAccountIdsByAccountName_MatchMultiple(t *testing.T) {
	accounts := []*models.Account{
		{
			AccountId: 1001,
			Name:      "Test Account",
			Type:      models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
		},
		{
			AccountId:       2001,
			Name:            "Test Account",
			Type:            models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS,
			ParentAccountId: 0,
		},
		{
			AccountId:       2002,
			Name:            "Sub 1-1",
			Type:            models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
			ParentAccountId: 2001,
		},
		{
			AccountId:       2003,
			Name:            "Sub 1-2",
			Type:            models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
			ParentAccountId: 2001,
		},
		{
			AccountId:       3001,
			Name:            "Test Account",
			Type:            models.ACCOUNT_TYPE_MULTI_SUB_ACCOUNTS,
			ParentAccountId: 0,
		},
		{
			AccountId:       3002,
			Name:            "Sub 2-1",
			Type:            models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
			ParentAccountId: 3001,
		},
		{
			AccountId: 4001,
			Name:      "Other Account",
			Type:      models.ACCOUNT_TYPE_SINGLE_ACCOUNT,
		},
	}
	actualAccountMap := Accounts.GetAccountOrSubAccountIdsByAccountName(accounts, "Test Account")

	assert.Equal(t, 4, len(actualAccountMap))
	assert.Contains(t, actualAccountMap, int64(1001))
	assert.Contains(t, actualAccountMap, int64(2002))
	assert.Contains(t, actualAccountMap, int64(2003))
	assert.Contains(t, actualAccountMap, int64(3002))
	assert.NotContains(t, actualAccountMap, int64(2001))
	assert.NotContains(t, actualAccountMap, int64(3001))
	assert.NotContains(t, actualAccountMap, int64(4001))
}
