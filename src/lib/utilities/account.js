import accountConstants from '../../consts/account.js';

export function getAccountCategoryInfo(categoryId) {
    for (let i = 0; i < accountConstants.allCategories.length; i++) {
        if (accountConstants.allCategories[i].id === categoryId) {
            return accountConstants.allCategories[i];
        }
    }

    return null;
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
