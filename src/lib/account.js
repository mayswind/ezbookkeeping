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
    account.comment = account2.comment;
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

export function getCategorizedAccounts(allAccounts) {
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

export function getVisibleCategorizedAccounts(categorizedAccounts) {
    const ret = {};

    for (let i = 0; i < accountConstants.allCategories.length; i++) {
        const accountCategory = accountConstants.allCategories[i];

        if (!categorizedAccounts[accountCategory.id] || !categorizedAccounts[accountCategory.id].accounts) {
            continue;
        }

        const allAccounts = categorizedAccounts[accountCategory.id].accounts;
        const visibleAccounts = [];
        const allVisibleSubAccounts = {};

        for (let j = 0; j < allAccounts.length; j++) {
            const account = allAccounts[j];

            if (account.hidden) {
                continue;
            }

            visibleAccounts.push(account);

            if (account.type === accountConstants.allAccountTypes.MultiSubAccounts && account.subAccounts) {
                const visibleSubAccounts = [];

                for (let k = 0; k < account.subAccounts.length; k++) {
                    const subAccount = account.subAccounts[k];

                    if (!subAccount.hidden) {
                        visibleSubAccounts.push(subAccount);
                    }
                }

                if (visibleSubAccounts.length > 0) {
                    allVisibleSubAccounts[account.id] = visibleSubAccounts;
                }
            }
        }

        if (visibleAccounts.length > 0) {
            ret[accountCategory.id] = {
                category: accountCategory.id,
                name: accountCategory.name,
                icon: accountCategory.defaultAccountIconId,
                visibleAccounts: visibleAccounts,
                visibleSubAccounts: allVisibleSubAccounts
            };
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
