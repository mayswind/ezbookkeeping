import { keys, keysIfValueEquals, values } from '@/core/base.ts';
import { AccountType, AccountCategory } from '@/core/account.ts';
import { PARENT_ACCOUNT_CURRENCY_PLACEHOLDER } from '@/consts/currency.ts';
import { type AccountBalance, type CategorizedAccount, Account } from '@/models/account.ts';

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

export function filterCategorizedAccounts(categorizedAccountsMap: Record<number, CategorizedAccount>, allowAccountName?: string, showHidden?: boolean): Record<number, CategorizedAccount> {
    const ret: Record<number, CategorizedAccount> = {};
    const allCategories = AccountCategory.values();
    const lowercaseFilterContent = allowAccountName ? allowAccountName.toLowerCase() : '';

    for (const accountCategory of allCategories) {
        const categorizedAccount = categorizedAccountsMap[accountCategory.type];

        if (!categorizedAccount || !categorizedAccount.accounts || categorizedAccount.accounts.length < 1) {
            continue;
        }

        const allFilteredAccounts: Account[] = [];

        for (const account of categorizedAccount.accounts) {
            if (!showHidden && account.hidden) {
                continue;
            }

            const accountMatchesName = !lowercaseFilterContent || account.name.toLowerCase().includes(lowercaseFilterContent);
            const filteredSubAccounts: Account[] = [];

            if (account.subAccounts) {
                for (const subAccount of account.subAccounts) {
                    if (!showHidden && subAccount.hidden) {
                        continue;
                    }

                    if (!accountMatchesName && lowercaseFilterContent && !subAccount.name.toLowerCase().includes(lowercaseFilterContent)) {
                        continue;
                    }

                    const filteredSubAccount = subAccount.clone();
                    filteredSubAccounts.push(filteredSubAccount);
                }
            }

            if (!accountMatchesName && filteredSubAccounts.length < 1) {
                continue;
            }

            const filteredAccount = account.cloneSelf();

            if (filteredAccount.type === AccountType.MultiSubAccounts.type) {
                filteredAccount.subAccounts = filteredSubAccounts;
            }

            allFilteredAccounts.push(filteredAccount);
        }

        if (allFilteredAccounts.length > 0) {
            ret[accountCategory.type] = {
                category: categorizedAccount.category,
                name: categorizedAccount.name,
                icon: categorizedAccount.icon,
                accounts: allFilteredAccounts
            };
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

export function selectAll(filterAccountIds: Record<string, boolean>, allAccountsMap: Record<string, Account>): void {
    for (const accountId of keys(filterAccountIds)) {
        const account = allAccountsMap[accountId];

        if (account && account.type === AccountType.SingleAccount.type) {
            filterAccountIds[account.id] = false;
        }
    }
}

export function selectNone(filterAccountIds: Record<string, boolean>, allAccountsMap: Record<string, Account>): void {
    for (const accountId of keys(filterAccountIds)) {
        const account = allAccountsMap[accountId];

        if (account && account.type === AccountType.SingleAccount.type) {
            filterAccountIds[account.id] = true;
        }
    }
}

export function selectInvert(filterAccountIds: Record<string, boolean>, allAccountsMap: Record<string, Account>): void {
    for (const accountId of keys(filterAccountIds)) {
        const account = allAccountsMap[accountId];

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
