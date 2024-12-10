import currencyConstants from '@/consts/currency.js';
import accountConstants from '@/consts/account.js';

export function setAccountModelByAnotherAccount(account, account2) {
    account.id = account2.id;
    account.category = account2.category;
    account.type = account2.type;
    account.name = account2.name;
    account.icon = account2.icon;
    account.color = account2.color;
    account.currency = account2.currency;
    account.balance = account2.balance;
    account.balanceTime = account2.balanceTime;
    account.comment = account2.comment;
    account.creditCardStatementDate = account2.creditCardStatementDate;
    account.visible = !account2.hidden;
}

export function getAccountCategoryInfo(categoryId) {
    for (let i = 0; i < accountConstants.allCategories.length; i++) {
        if (accountConstants.allCategories[i].id === categoryId) {
            return accountConstants.allCategories[i];
        }
    }

    return null;
}

export function getAccountOrSubAccountId(account, subAccountId) {
    if (account.type === accountConstants.allAccountTypes.SingleAccount) {
        return account.id;
    } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts && !subAccountId) {
        return account.id;
    } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts && subAccountId) {
        if (!account.subAccounts || !account.subAccounts.length) {
            return null;
        }

        for (let i = 0; i < account.subAccounts.length; i++) {
            const subAccount = account.subAccounts[i];

            if (subAccountId && subAccountId === subAccount.id) {
                return subAccount.id;
            }
        }

        return null;
    } else {
        return null;
    }
}

export function getAccountOrSubAccountComment(account, subAccountId) {
    if (account.type === accountConstants.allAccountTypes.SingleAccount) {
        return account.comment;
    } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts && !subAccountId) {
        return account.comment;
    } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts && subAccountId) {
        if (!account.subAccounts || !account.subAccounts.length) {
            return null;
        }

        for (let i = 0; i < account.subAccounts.length; i++) {
            const subAccount = account.subAccounts[i];

            if (subAccountId && subAccountId === subAccount.id) {
                return subAccount.comment;
            }
        }

        return null;
    } else {
        return null;
    }
}

export function getSubAccountCurrencies(account, showHidden, subAccountId) {
    if (!account.subAccounts || !account.subAccounts.length) {
        return [];
    }

    const subAccountCurrenciesMap = {};
    const subAccountCurrencies = [];

    for (let i = 0; i < account.subAccounts.length; i++) {
        const subAccount = account.subAccounts[i];

        if (!showHidden && subAccount.hidden) {
            continue;
        }

        if (subAccountId && subAccountId === subAccount.id) {
            return [subAccount.currency];
        } else {
            if (!subAccountCurrenciesMap[subAccount.currency]) {
                subAccountCurrenciesMap[subAccount.currency] = true;
                subAccountCurrencies.push(subAccount.currency);
            }
        }
    }

    return subAccountCurrencies;
}

export function getCategorizedAccountsMap(allAccounts) {
    const ret = {};

    for (let i = 0; i < allAccounts.length; i++) {
        const account = allAccounts[i];

        if (!ret[account.category]) {
            const categoryInfo = getAccountCategoryInfo(account.category);

            if (categoryInfo) {
                ret[account.category] = {
                    category: account.category,
                    name: categoryInfo.name,
                    icon: categoryInfo.defaultAccountIconId,
                    accounts: []
                };
            }
        }

        if (ret[account.category]) {
            const accountList = ret[account.category].accounts;
            accountList.push(account);
        }
    }

    return ret;
}

export function getCategorizedAccounts(allAccounts) {
    const ret = [];
    const categorizedAccounts = getCategorizedAccountsMap(allAccounts);

    for (let i = 0; i < accountConstants.allCategories.length; i++) {
        const category = accountConstants.allCategories[i];

        if (!categorizedAccounts[category.id]) {
            continue;
        }

        const accountCategory = categorizedAccounts[category.id];
        ret.push(accountCategory);
    }

    return ret;
}

export function getCategorizedAccountsWithVisibleCount(categorizedAccountsMap) {
    const ret = [];

    for (let i = 0; i < accountConstants.allCategories.length; i++) {
        const accountCategory = accountConstants.allCategories[i];

        if (!categorizedAccountsMap[accountCategory.id] || !categorizedAccountsMap[accountCategory.id].accounts) {
            continue;
        }

        const allAccounts = categorizedAccountsMap[accountCategory.id].accounts;
        const allSubAccounts = {};
        const allVisibleSubAccountCounts = {};
        const allFirstVisibleSubAccountIndexes = {};
        let allVisibleAccountCount = 0;
        let firstVisibleAccountIndex = -1;

        for (let j = 0; j < allAccounts.length; j++) {
            const account = allAccounts[j];

            if (!account.hidden) {
                allVisibleAccountCount++;

                if (firstVisibleAccountIndex === -1) {
                    firstVisibleAccountIndex = j;
                }
            }

            if (account.type === accountConstants.allAccountTypes.MultiSubAccounts && account.subAccounts) {
                let visibleSubAccountCount = 0;
                let firstVisibleSubAccountIndex = -1;

                for (let k = 0; k < account.subAccounts.length; k++) {
                    const subAccount = account.subAccounts[k];

                    if (!subAccount.hidden) {
                        visibleSubAccountCount++;

                        if (firstVisibleSubAccountIndex === -1) {
                            firstVisibleSubAccountIndex = k;
                        }
                    }
                }

                if (account.subAccounts.length > 0) {
                    allSubAccounts[account.id] = account.subAccounts;
                    allVisibleSubAccountCounts[account.id] = visibleSubAccountCount;
                    allFirstVisibleSubAccountIndexes[account.id] = firstVisibleSubAccountIndex;
                }
            }
        }

        if (allAccounts.length > 0) {
            ret.push({
                category: accountCategory.id,
                name: accountCategory.name,
                icon: accountCategory.defaultAccountIconId,
                allAccounts: allAccounts,
                allVisibleAccountCount: allVisibleAccountCount,
                firstVisibleAccountIndex: firstVisibleAccountIndex,
                allSubAccounts: allSubAccounts,
                allVisibleSubAccountCounts: allVisibleSubAccountCounts,
                allFirstVisibleSubAccountIndexes: allFirstVisibleSubAccountIndexes
            });
        }
    }

    return ret;
}

export function getAllFilteredAccountsBalance(categorizedAccounts, accountFilter) {
    const allAccountCategories = accountConstants.allCategories;
    const ret = [];

    for (let categoryIdx = 0; categoryIdx < allAccountCategories.length; categoryIdx++) {
        const accountCategory = allAccountCategories[categoryIdx];

        if (!categorizedAccounts[accountCategory.id] || !categorizedAccounts[accountCategory.id].accounts) {
            continue;
        }

        for (let accountIdx = 0; accountIdx < categorizedAccounts[accountCategory.id].accounts.length; accountIdx++) {
            const account = categorizedAccounts[accountCategory.id].accounts[accountIdx];

            if (account.hidden || !accountFilter(account)) {
                continue;
            }

            if (account.type === accountConstants.allAccountTypes.SingleAccount) {
                ret.push({
                    balance: account.balance,
                    isAsset: account.isAsset,
                    isLiability: account.isLiability,
                    currency: account.currency
                });
            } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts) {
                for (let subAccountIdx = 0; subAccountIdx < account.subAccounts.length; subAccountIdx++) {
                    const subAccount = account.subAccounts[subAccountIdx];

                    if (subAccount.hidden || !accountFilter(subAccount)) {
                        continue;
                    }

                    ret.push({
                        balance: subAccount.balance,
                        isAsset: subAccount.isAsset,
                        isLiability: subAccount.isLiability,
                        currency: subAccount.currency
                    });
                }
            }
        }
    }

    return ret;
}

export function getFinalAccountIdsByFilteredAccountIds(allAccountsMap, filteredAccountIds) {
    let finalAccountIds = '';

    if (!allAccountsMap) {
        return finalAccountIds;
    }

    for (let accountId in allAccountsMap) {
        if (!Object.prototype.hasOwnProperty.call(allAccountsMap, accountId)) {
            continue;
        }

        const account = allAccountsMap[accountId];

        if (filteredAccountIds && !isAccountOrSubAccountsAllChecked(account, filteredAccountIds)) {
            continue;
        }

        if (finalAccountIds.length > 0) {
            finalAccountIds += ',';
        }

        finalAccountIds += account.id;
    }

    return finalAccountIds;
}

export function getUnifiedSelectedAccountsCurrencyOrDefaultCurrency(allAccounts, selectedAccountIds, defaultCurrency) {
    if (!selectedAccountIds) {
        return defaultCurrency;
    }

    let accountCurrency = '';

    for (let accountId in selectedAccountIds) {
        if (!Object.prototype.hasOwnProperty.call(selectedAccountIds, accountId)) {
            continue;
        }

        const account = allAccounts[accountId];

        if (account.currency === currencyConstants.parentAccountCurrencyPlaceholder) {
            continue;
        }

        if (accountCurrency === '') {
            accountCurrency = account.currency;
        } else if (accountCurrency !== account.currency) {
            return defaultCurrency;
        }
    }

    if (accountCurrency) {
        return accountCurrency;
    }

    return defaultCurrency;
}

export function selectAccountOrSubAccounts(filterAccountIds, account, value) {
    if (account.type === accountConstants.allAccountTypes.SingleAccount) {
        filterAccountIds[account.id] = value;
    } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts) {
        if (!account.subAccounts || !account.subAccounts.length) {
            return;
        }

        for (let i = 0; i < account.subAccounts.length; i++) {
            const subAccount = account.subAccounts[i];
            filterAccountIds[subAccount.id] = value;
        }
    }
}

export function selectAll(filterAccountIds, allAccountsMap) {
    for (let accountId in filterAccountIds) {
        if (!Object.prototype.hasOwnProperty.call(filterAccountIds, accountId)) {
            continue;
        }

        const account = allAccountsMap[accountId];

        if (account && account.type === accountConstants.allAccountTypes.SingleAccount) {
            filterAccountIds[account.id] = false;
        }
    }
}

export function selectNone(filterAccountIds, allAccountsMap) {
    for (let accountId in filterAccountIds) {
        if (!Object.prototype.hasOwnProperty.call(filterAccountIds, accountId)) {
            continue;
        }

        const account = allAccountsMap[accountId];

        if (account && account.type === accountConstants.allAccountTypes.SingleAccount) {
            filterAccountIds[account.id] = true;
        }
    }
}

export function selectInvert(filterAccountIds, allAccountsMap) {
    for (let accountId in filterAccountIds) {
        if (!Object.prototype.hasOwnProperty.call(filterAccountIds, accountId)) {
            continue;
        }

        const account = allAccountsMap[accountId];

        if (account && account.type === accountConstants.allAccountTypes.SingleAccount) {
            filterAccountIds[account.id] = !filterAccountIds[account.id];
        }
    }
}

export function isAccountOrSubAccountsAllChecked(account, filterAccountIds) {
    if (!account.subAccounts) {
        return !filterAccountIds[account.id];
    }

    for (let i = 0; i < account.subAccounts.length; i++) {
        const subAccount = account.subAccounts[i];
        if (filterAccountIds[subAccount.id]) {
            return false;
        }
    }

    return true;
}

export function isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds) {
    if (!account.subAccounts) {
        return false;
    }

    let checkedCount = 0;

    for (let i = 0; i < account.subAccounts.length; i++) {
        const subAccount = account.subAccounts[i];
        if (!filterAccountIds[subAccount.id]) {
            checkedCount++;
        }
    }

    return checkedCount > 0 && checkedCount < account.subAccounts.length;
}

export function setAccountSuitableIcon(account, oldCategory, newCategory) {
    for (let i = 0; i < accountConstants.allCategories.length; i++) {
        if (accountConstants.allCategories[i].id === oldCategory) {
            if (account.icon !== accountConstants.allCategories[i].defaultAccountIconId) {
                return;
            } else {
                break;
            }
        }
    }

    for (let i = 0; i < accountConstants.allCategories.length; i++) {
        if (accountConstants.allCategories[i].id === newCategory) {
            account.icon = accountConstants.allCategories[i].defaultAccountIconId;
        }
    }
}
