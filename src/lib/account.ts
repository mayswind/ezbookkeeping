import { itemAndIndex, keys, keysIfValueEquals, values } from '@/core/base.ts';
import { AccountType, AccountCategory } from '@/core/account.ts';
import { PARENT_ACCOUNT_CURRENCY_PLACEHOLDER } from '@/consts/currency.ts';
import { type AccountBalance, type CategorizedAccount, type AccountCategoriesWithVisibleCount, Account } from '@/models/account.ts';

export function getCategorizedAccountsMap(allAccounts: Account[]): Record<number, CategorizedAccount> {
    const ret: Record<number, CategorizedAccount> = {};

    for (const account of allAccounts) {
        if (!ret[account.category]) {
            const categoryInfo = AccountCategory.valueOf(account.category);

            if (categoryInfo) {
                ret[account.category] = {
                    category: account.category,
                    name: categoryInfo.name,
                    icon: categoryInfo.defaultAccountIconId,
                    accounts: []
                };
            }
        }

        const categorizedAccount = ret[account.category];

        if (categorizedAccount) {
            const accountList = categorizedAccount.accounts;
            accountList.push(account);
        }
    }

    return ret;
}

export function getCategorizedAccounts(allAccounts: Account[]): CategorizedAccount[] {
    const ret: CategorizedAccount[] = [];
    const allCategories = AccountCategory.values();
    const categorizedAccounts = getCategorizedAccountsMap(allAccounts);

    for (const category of allCategories) {
        if (!categorizedAccounts[category.type]) {
            continue;
        }

        const accountCategory = categorizedAccounts[category.type];

        if (accountCategory) {
            ret.push(accountCategory);
        }
    }

    return ret;
}

export function getAccountMapByName(allAccounts: Account[]): Record<string, Account> {
    const ret: Record<string, Account> = {};

    if (!allAccounts) {
        return ret;
    }

    for (const account of allAccounts) {
        if (account.type === AccountType.SingleAccount.type) {
            ret[account.name] = account;
        } else if (account.type === AccountType.MultiSubAccounts.type && account.subAccounts) {
            for (const subAccount of account.subAccounts) {
                ret[subAccount.name] = subAccount;
            }
        }
    }

    return ret;
}

export function getCategorizedAccountsWithVisibleCount(categorizedAccountsMap: Record<number, CategorizedAccount>): AccountCategoriesWithVisibleCount[] {
    const ret: AccountCategoriesWithVisibleCount[] = [];
    const allCategories = AccountCategory.values();

    for (const accountCategory of allCategories) {
        const categorizedAccount = categorizedAccountsMap[accountCategory.type];

        if (!categorizedAccount || !categorizedAccount.accounts) {
            continue;
        }

        const allAccounts = categorizedAccount.accounts;
        const allSubAccounts: Record<string, Account[]> = {};
        const allVisibleSubAccountCounts: Record<string, number> = {};
        const allFirstVisibleSubAccountIndexes: Record<string, number> = {};
        let allVisibleAccountCount = 0;
        let firstVisibleAccountIndex = -1;

        for (const [account, accountIndex] of itemAndIndex(allAccounts)) {
            if (!account.hidden) {
                allVisibleAccountCount++;

                if (firstVisibleAccountIndex === -1) {
                    firstVisibleAccountIndex = accountIndex;
                }
            }

            if (account.type === AccountType.MultiSubAccounts.type && account.subAccounts) {
                let visibleSubAccountCount = 0;
                let firstVisibleSubAccountIndex = -1;

                for (const [subAccount, subAccountIndex] of itemAndIndex(account.subAccounts)) {
                    if (!subAccount.hidden) {
                        visibleSubAccountCount++;

                        if (firstVisibleSubAccountIndex === -1) {
                            firstVisibleSubAccountIndex = subAccountIndex;
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
                category: accountCategory.type,
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

export function getAllFilteredAccountsBalance(categorizedAccounts: Record<number, CategorizedAccount>, accountFilter: (account: Account) => boolean): AccountBalance[] {
    const allAccountCategories = AccountCategory.values();
    const ret: AccountBalance[] = [];

    for (const accountCategory of allAccountCategories) {
        const categorizedAccount = categorizedAccounts[accountCategory.type];

        if (!categorizedAccount || !categorizedAccount.accounts) {
            continue;
        }

        for (const account of categorizedAccount.accounts) {
            if (account.hidden || !accountFilter(account)) {
                continue;
            }

            if (account.type === AccountType.SingleAccount.type) {
                ret.push({
                    balance: account.balance,
                    isAsset: !!account.isAsset,
                    isLiability: !!account.isLiability,
                    currency: account.currency
                });
            } else if (account.type === AccountType.MultiSubAccounts.type && account.subAccounts) {
                for (const subAccount of account.subAccounts) {
                    if (subAccount.hidden || !accountFilter(subAccount)) {
                        continue;
                    }

                    ret.push({
                        balance: subAccount.balance,
                        isAsset: !!subAccount.isAsset,
                        isLiability: !!subAccount.isLiability,
                        currency: subAccount.currency
                    });
                }
            }
        }
    }

    return ret;
}

export function getFinalAccountIdsByFilteredAccountIds(allAccountsMap: Record<string, Account>, filteredAccountIds: Record<string, boolean>): string {
    let finalAccountIds = '';

    if (!allAccountsMap) {
        return finalAccountIds;
    }

    for (const account of values(allAccountsMap)) {
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

export function getUnifiedSelectedAccountsCurrencyOrDefaultCurrency(allAccountsMap: Record<string, Account>, selectedAccountIds: Record<string, boolean>, defaultCurrency: string): string {
    if (!selectedAccountIds) {
        return defaultCurrency;
    }

    let accountCurrency = '';

    for (const accountId of keysIfValueEquals(selectedAccountIds, true)) {
        const account = allAccountsMap[accountId];

        if (!account) {
            continue;
        }

        if (account.currency === PARENT_ACCOUNT_CURRENCY_PLACEHOLDER) {
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

export function selectAccountOrSubAccounts(filterAccountIds: Record<string, boolean>, account: Account, value: boolean): void {
    if (account.type === AccountType.SingleAccount.type) {
        filterAccountIds[account.id] = value;
    } else if (account.type === AccountType.MultiSubAccounts.type) {
        if (!account.subAccounts || !account.subAccounts.length) {
            return;
        }

        for (const subAccount of account.subAccounts) {
            filterAccountIds[subAccount.id] = value;
        }
    }
}

export function selectAll(filterAccountIds: Record<string, boolean>, allAccountsMap: Record<string, Account>, skipHiddenAccount: boolean): void {
    for (const accountId of keys(filterAccountIds)) {
        const account = allAccountsMap[accountId];

        if (skipHiddenAccount && account && account.hidden) {
            continue;
        }

        if (account && account.type === AccountType.SingleAccount.type) {
            filterAccountIds[account.id] = false;
        }
    }
}

export function selectNone(filterAccountIds: Record<string, boolean>, allAccountsMap: Record<string, Account>, skipHiddenAccount: boolean): void {
    for (const accountId of keys(filterAccountIds)) {
        const account = allAccountsMap[accountId];

        if (skipHiddenAccount && account && account.hidden) {
            continue;
        }

        if (account && account.type === AccountType.SingleAccount.type) {
            filterAccountIds[account.id] = true;
        }
    }
}

export function selectInvert(filterAccountIds: Record<string, boolean>, allAccountsMap: Record<string, Account>, skipHiddenAccount: boolean): void {
    for (const accountId of keys(filterAccountIds)) {
        const account = allAccountsMap[accountId];

        if (skipHiddenAccount && account && account.hidden) {
            continue;
        }

        if (account && account.type === AccountType.SingleAccount.type) {
            filterAccountIds[account.id] = !filterAccountIds[account.id];
        }
    }
}

export function isAccountOrSubAccountsAllChecked(account: Account, filterAccountIds: Record<string, boolean>): boolean {
    if (!account.subAccounts) {
        return !filterAccountIds[account.id];
    }

    for (const subAccount of account.subAccounts) {
        if (filterAccountIds[subAccount.id]) {
            return false;
        }
    }

    return true;
}

export function isAccountOrSubAccountsHasButNotAllChecked(account: Account, filterAccountIds: Record<string, boolean>): boolean {
    if (!account.subAccounts) {
        return false;
    }

    let checkedCount = 0;

    for (const subAccount of account.subAccounts) {
        if (!filterAccountIds[subAccount.id]) {
            checkedCount++;
        }
    }

    return checkedCount > 0 && checkedCount < account.subAccounts.length;
}
