import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useUserStore } from './user.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import type { BeforeResolveFunction } from '@/core/base.ts';

import { AccountType, AccountCategory } from '@/core/account.ts';
import {
    type AccountNewDisplayOrderRequest,
    type AccountDisplayBalance,
    type CategorizedAccount,
    type AccountShowingIds,
    Account
} from '@/models/account.ts';

import { isNumber, isEquals } from '@/lib/common.ts';
import { getCategorizedAccountsMap, getAllFilteredAccountsBalance } from '@/lib/account.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useAccountsStore = defineStore('accounts', () => {
    const userStore = useUserStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const allAccounts = ref<Account[]>([]);
    const allAccountsMap = ref<Record<string, Account>>({});
    const allCategorizedAccountsMap = ref<Record<number, CategorizedAccount>>({});
    const accountListStateInvalid = ref<boolean>(true);

    const allPlainAccounts = computed<Account[]>(() => {
        const allAccountsList: Account[] = [];

        for (let i = 0; i < allAccounts.value.length; i++) {
            const account = allAccounts.value[i];

            if (account.type === AccountType.SingleAccount.type) {
                allAccountsList.push(account);
            } else if (account.type === AccountType.MultiSubAccounts.type) {
                if (account.subAccounts) {
                    for (let j = 0; j < account.subAccounts.length; j++) {
                        const subAccount = account.subAccounts[j];
                        allAccountsList.push(subAccount);
                    }
                }
            }
        }

        return Account.sortAccounts(allAccountsList, allAccountsMap.value);
    });

    const allMixedPlainAccounts = computed<Account[]>(() => {
        const allAccountsList: Account[] = [];

        for (let i = 0; i < allAccounts.value.length; i++) {
            const account = allAccounts.value[i];

            if (account.type === AccountType.SingleAccount.type) {
                allAccountsList.push(account);
            } else if (account.type === AccountType.MultiSubAccounts.type) {
                allAccountsList.push(account);

                if (account.subAccounts) {
                    for (let j = 0; j < account.subAccounts.length; j++) {
                        const subAccount = account.subAccounts[j];
                        allAccountsList.push(subAccount);
                    }
                }
            }
        }

        return Account.sortAccounts(allAccountsList, allAccountsMap.value);
    });

    const allVisiblePlainAccounts = computed<Account[]>(() => {
        const allVisibleAccounts: Account[] = [];

        for (let i = 0; i < allAccounts.value.length; i++) {
            const account = allAccounts.value[i];

            if (account.hidden) {
                continue;
            }

            if (account.type === AccountType.SingleAccount.type) {
                allVisibleAccounts.push(account);
            } else if (account.type === AccountType.MultiSubAccounts.type) {
                if (account.subAccounts) {
                    for (let j = 0; j < account.subAccounts.length; j++) {
                        const subAccount = account.subAccounts[j];
                        allVisibleAccounts.push(subAccount);
                    }
                }
            }
        }

        return Account.sortAccounts(allVisibleAccounts, allAccountsMap.value);
    });

    const allAvailableAccountsCount = computed<number>(() => {
        let allAccountCount = 0;

        for (const category in allCategorizedAccountsMap.value) {
            if (!Object.prototype.hasOwnProperty.call(allCategorizedAccountsMap.value, category)) {
                continue;
            }

            allAccountCount += allCategorizedAccountsMap.value[category].accounts.length;
        }

        return allAccountCount;
    });

    const allVisibleAccountsCount = computed<number>(() => {
        let shownAccountCount = 0;

        for (const category in allCategorizedAccountsMap.value) {
            if (!Object.prototype.hasOwnProperty.call(allCategorizedAccountsMap.value, category)) {
                continue;
            }

            const accountList = allCategorizedAccountsMap.value[category].accounts;

            for (let i = 0; i < accountList.length; i++) {
                if (!accountList[i].hidden) {
                    shownAccountCount++;
                }
            }
        }

        return shownAccountCount;
    });

    function loadAccountList(accounts: Account[]): void {
        allAccounts.value = accounts;
        allAccountsMap.value = {};

        for (let i = 0; i < accounts.length; i++) {
            const account = accounts[i];
            allAccountsMap.value[account.id] = account;

            if (account.subAccounts) {
                for (let j = 0; j < account.subAccounts.length; j++) {
                    const subAccount = account.subAccounts[j];
                    allAccountsMap.value[subAccount.id] = subAccount;
                }
            }
        }

        allCategorizedAccountsMap.value = getCategorizedAccountsMap(accounts);
    }

    function addAccountToAccountList(account: Account): void {
        const newAccountCategory = AccountCategory.valueOf(account.category);
        let insertIndexToAllList = allAccounts.value.length;

        if (newAccountCategory) {
            for (let i = 0; i < allAccounts.value.length; i++) {
                const accountCategory = AccountCategory.valueOf(allAccounts.value[i].category);

                if (accountCategory && accountCategory.displayOrder > newAccountCategory.displayOrder) {
                    insertIndexToAllList = i;
                    break;
                }
            }
        }

        allAccounts.value.splice(insertIndexToAllList, 0, account);

        allAccountsMap.value[account.id] = account;

        if (account.subAccounts) {
            for (let i = 0; i < account.subAccounts.length; i++) {
                const subAccount = account.subAccounts[i];
                allAccountsMap.value[subAccount.id] = subAccount;
            }
        }

        if (allCategorizedAccountsMap.value[account.category]) {
            const accountList = allCategorizedAccountsMap.value[account.category].accounts;
            accountList.push(account);
        } else {
            allCategorizedAccountsMap.value = getCategorizedAccountsMap(allAccounts.value);
        }
    }

    function updateAccountToAccountList(oldAccount: Account, newAccount: Account): void {
        for (let i = 0; i < allAccounts.value.length; i++) {
            if (allAccounts.value[i].id === newAccount.id) {
                allAccounts.value.splice(i, 1, newAccount);
                break;
            }
        }

        if (oldAccount.subAccounts) {
            for (let i = 0; i < oldAccount.subAccounts.length; i++) {
                const subAccount = oldAccount.subAccounts[i];
                if (allAccountsMap.value[subAccount.id]) {
                    delete allAccountsMap.value[subAccount.id];
                }
            }
        }

        allAccountsMap.value[newAccount.id] = newAccount;

        if (newAccount.subAccounts) {
            for (let i = 0; i < newAccount.subAccounts.length; i++) {
                const subAccount = newAccount.subAccounts[i];
                allAccountsMap.value[subAccount.id] = subAccount;
            }
        }

        if (allCategorizedAccountsMap.value[newAccount.category]) {
            const accountList = allCategorizedAccountsMap.value[newAccount.category].accounts;

            for (let i = 0; i < accountList.length; i++) {
                if (accountList[i].id === newAccount.id) {
                    accountList.splice(i, 1, newAccount);
                    break;
                }
            }
        }
    }

    function updateAccountDisplayOrderInAccountList({ account, from, to, updateListOrder, updateGlobalListOrder }: { account: Account, from: number, to: number, updateListOrder: boolean, updateGlobalListOrder: boolean }): void {
        let fromAccount = null;
        let toAccount = null;

        if (allCategorizedAccountsMap.value[account.category]) {
            const accountList = allCategorizedAccountsMap.value[account.category].accounts;

            if (updateListOrder) {
                fromAccount = accountList[from];
                toAccount = accountList[to];
                accountList.splice(to, 0, accountList.splice(from, 1)[0]);
            } else {
                fromAccount = accountList[to];

                if (from < to) {
                    toAccount = accountList[to - 1];
                } else if (from > to) {
                    toAccount = accountList[to + 1];
                }
            }
        }

        if (updateGlobalListOrder && fromAccount && toAccount) {
            let globalFromIndex = -1;
            let globalToIndex = -1;

            for (let i = 0; i < allAccounts.value.length; i++) {
                if (allAccounts.value[i].id === fromAccount.id) {
                    globalFromIndex = i;
                } else if (allAccounts.value[i].id === toAccount.id) {
                    globalToIndex = i;
                }
            }

            if (globalFromIndex >= 0 && globalToIndex >= 0) {
                allAccounts.value.splice(globalToIndex, 0, allAccounts.value.splice(globalFromIndex, 1)[0]);
            }
        }
    }

    function updateAccountVisibilityInAccountList({ account, hidden }: { account: Account, hidden: boolean }): void {
        if (allAccountsMap.value[account.id]) {
            allAccountsMap.value[account.id].visible = !hidden;
        }
    }

    function removeAccountFromAccountList(account: Account): void {
        for (let i = 0; i < allAccounts.value.length; i++) {
            if (allAccounts.value[i].id === account.id) {
                allAccounts.value.splice(i, 1);
                break;
            }
        }

        if (allAccountsMap.value[account.id] && allAccountsMap.value[account.id].subAccounts) {
            const subAccounts = allAccountsMap.value[account.id].subAccounts as Account[];

            for (let i = 0; i < subAccounts.length; i++) {
                const subAccount = subAccounts[i];
                if (allAccountsMap.value[subAccount.id]) {
                    delete allAccountsMap.value[subAccount.id];
                }
            }
        }

        if (allAccountsMap.value[account.id]) {
            delete allAccountsMap.value[account.id];
        }

        if (allCategorizedAccountsMap.value[account.category]) {
            const accountList = allCategorizedAccountsMap.value[account.category].accounts;

            for (let i = 0; i < accountList.length; i++) {
                if (accountList[i].id === account.id) {
                    accountList.splice(i, 1);
                    break;
                }
            }
        }
    }

    function removeSubAccountFromAccountList(subAccount: Account): void {
        for (let i = 0; i < allAccounts.value.length; i++) {
            if (allAccounts.value[i].type !== AccountType.MultiSubAccounts.type || !allAccounts.value[i].subAccounts) {
                continue;
            }

            const subAccounts = allAccounts.value[i].subAccounts as Account[];

            for (let j = 0; j < subAccounts.length; j++) {
                if (subAccounts[j].id === subAccount.id) {
                    subAccounts.splice(j, 1);
                    break;
                }
            }
        }

        if (allAccountsMap.value[subAccount.id]) {
            delete allAccountsMap.value[subAccount.id];
        }

        if (allCategorizedAccountsMap.value[subAccount.category]) {
            const accountList = allCategorizedAccountsMap.value[subAccount.category].accounts;

            for (let i = 0; i < accountList.length; i++) {
                if (accountList[i].type !== AccountType.MultiSubAccounts.type || !accountList[i].subAccounts) {
                    continue;
                }

                const subAccounts = accountList[i].subAccounts as Account[];

                for (let j = 0; j < subAccounts.length; j++) {
                    if (subAccounts[j].id === subAccount.id) {
                        subAccounts.splice(j, 1);
                        break;
                    }
                }
            }
        }
    }

    function updateAccountListInvalidState(invalidState: boolean): void {
        accountListStateInvalid.value = invalidState;
    }

    function resetAccounts(): void {
        allAccounts.value = [];
        allAccountsMap.value = {};
        allCategorizedAccountsMap.value = {};
        accountListStateInvalid.value = true;
    }

    function getFirstShowingIds(showHidden: boolean): AccountShowingIds {
        const ret: AccountShowingIds = {
            accounts: {},
            subAccounts: {}
        };

        for (const category in allCategorizedAccountsMap.value) {
            if (!Object.prototype.hasOwnProperty.call(allCategorizedAccountsMap.value, category)) {
                continue;
            }

            if (!allCategorizedAccountsMap.value[category] || !allCategorizedAccountsMap.value[category].accounts) {
                continue;
            }

            const accounts = allCategorizedAccountsMap.value[category].accounts;

            for (let i = 0; i < accounts.length; i++) {
                const account = accounts[i];

                if (account.type === AccountType.MultiSubAccounts.type && account.subAccounts) {
                    for (let j = 0; j < account.subAccounts.length; j++) {
                        const subAccount = account.subAccounts[j];

                        if (showHidden || !subAccount.hidden) {
                            ret.subAccounts[account.id] = subAccount.id;
                            break;
                        }
                    }
                }

                if (showHidden || !account.hidden) {
                    ret.accounts[category] = account.id;
                    break;
                }
            }
        }

        return ret;
    }

    function getLastShowingIds(showHidden: boolean): AccountShowingIds {
        const ret: AccountShowingIds = {
            accounts: {},
            subAccounts: {}
        };

        for (const category in allCategorizedAccountsMap.value) {
            if (!Object.prototype.hasOwnProperty.call(allCategorizedAccountsMap.value, category)) {
                continue;
            }

            if (!allCategorizedAccountsMap.value[category] || !allCategorizedAccountsMap.value[category].accounts) {
                continue;
            }

            const accounts = allCategorizedAccountsMap.value[category].accounts;

            for (let i = accounts.length - 1; i >= 0; i--) {
                const account = accounts[i];

                if (account.type === AccountType.MultiSubAccounts.type && account.subAccounts) {
                    for (let j = account.subAccounts.length - 1; j >= 0; j--) {
                        const subAccount = account.subAccounts[j];

                        if (showHidden || !subAccount.hidden) {
                            ret.subAccounts[account.id] = subAccount.id;
                            break;
                        }
                    }
                }

                if (showHidden || !account.hidden) {
                    ret.accounts[category] = account.id;
                    break;
                }
            }
        }

        return ret;
    }

    function getAccountStatementDate(accountId?: string): number | undefined | null {
        if (!accountId) {
            return null;
        }

        const accountIds = accountId.split(',');
        let mainAccount = null;

        for (let i = 0; i < accountIds.length; i++) {
            const id = accountIds[i];
            let account = allAccountsMap.value[id];

            if (!account) {
                return null;
            }

            if (account.parentId !== '0') {
                account = allAccountsMap.value[account.parentId];
            }

            if (mainAccount !== null) {
                if (mainAccount.id !== account.id) {
                    return null;
                } else {
                    continue;
                }
            }

            mainAccount = account;
        }

        if (!mainAccount) {
            return null;
        }

        if (mainAccount.category === AccountCategory.CreditCard.type) {
            return mainAccount.creditCardStatementDate;
        }

        return null;
    }

    function getNetAssets(showAccountBalance: boolean): string {
        if (!showAccountBalance) {
            return '***';
        }

        const accountsBalance = getAllFilteredAccountsBalance(allCategorizedAccountsMap.value, () => true);
        let netAssets = 0;
        let hasUnCalculatedAmount = false;

        for (let i = 0; i < accountsBalance.length; i++) {
            if (accountsBalance[i].currency === userStore.currentUserDefaultCurrency) {
                netAssets += accountsBalance[i].balance;
            } else {
                const balance = exchangeRatesStore.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, userStore.currentUserDefaultCurrency);

                if (!isNumber(balance)) {
                    hasUnCalculatedAmount = true;
                    continue;
                }

                netAssets += Math.floor(balance);
            }
        }

        if (hasUnCalculatedAmount) {
            return netAssets + '+';
        } else {
            return netAssets.toString();
        }
    }

    function getTotalAssets(showAccountBalance: boolean): string {
        if (!showAccountBalance) {
            return '***';
        }

        const accountsBalance = getAllFilteredAccountsBalance(allCategorizedAccountsMap.value, account => account.isAsset || false);
        let totalAssets = 0;
        let hasUnCalculatedAmount = false;

        for (let i = 0; i < accountsBalance.length; i++) {
            if (accountsBalance[i].currency === userStore.currentUserDefaultCurrency) {
                totalAssets += accountsBalance[i].balance;
            } else {
                const balance = exchangeRatesStore.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, userStore.currentUserDefaultCurrency);

                if (!isNumber(balance)) {
                    hasUnCalculatedAmount = true;
                    continue;
                }

                totalAssets += Math.floor(balance);
            }
        }

        if (hasUnCalculatedAmount) {
            return totalAssets + '+';
        } else {
            return totalAssets.toString();
        }
    }

    function getTotalLiabilities(showAccountBalance: boolean): string {
        if (!showAccountBalance) {
            return '***';
        }

        const accountsBalance = getAllFilteredAccountsBalance(allCategorizedAccountsMap.value, account => account.isLiability || false);
        let totalLiabilities = 0;
        let hasUnCalculatedAmount = false;

        for (let i = 0; i < accountsBalance.length; i++) {
            if (accountsBalance[i].currency === userStore.currentUserDefaultCurrency) {
                totalLiabilities -= accountsBalance[i].balance;
            } else {
                const balance = exchangeRatesStore.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, userStore.currentUserDefaultCurrency);

                if (!isNumber(balance)) {
                    hasUnCalculatedAmount = true;
                    continue;
                }

                totalLiabilities -= Math.floor(balance);
            }
        }

        if (hasUnCalculatedAmount) {
            return totalLiabilities + '+';
        } else {
            return totalLiabilities.toString();
        }
    }

    function getAccountCategoryTotalBalance(showAccountBalance: boolean, accountCategory: AccountCategory): string {
        if (!showAccountBalance) {
            return '***';
        }

        const accountsBalance = getAllFilteredAccountsBalance(allCategorizedAccountsMap.value, account => account.category === accountCategory.type);
        let totalBalance = 0;
        let hasUnCalculatedAmount = false;

        for (let i = 0; i < accountsBalance.length; i++) {
            if (accountsBalance[i].currency === userStore.currentUserDefaultCurrency) {
                if (accountsBalance[i].isAsset) {
                    totalBalance += accountsBalance[i].balance;
                } else if (accountsBalance[i].isLiability) {
                    totalBalance -= accountsBalance[i].balance;
                } else {
                    totalBalance += accountsBalance[i].balance;
                }
            } else {
                const balance = exchangeRatesStore.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, userStore.currentUserDefaultCurrency);

                if (!isNumber(balance)) {
                    hasUnCalculatedAmount = true;
                    continue;
                }

                if (accountsBalance[i].isAsset) {
                    totalBalance += Math.floor(balance);
                } else if (accountsBalance[i].isLiability) {
                    totalBalance -= Math.floor(balance);
                } else {
                    totalBalance += Math.floor(balance);
                }
            }
        }

        if (hasUnCalculatedAmount) {
            return totalBalance + '+';
        } else {
            return totalBalance.toString();
        }
    }

    function getAccountBalance(showAccountBalance: boolean, account: Account): string | null {
        if (account.type !== AccountType.SingleAccount.type) {
            return null;
        }

        if (showAccountBalance) {
            if (account.isAsset) {
                return account.balance.toString();
            } else if (account.isLiability) {
                return (-account.balance).toString();
            } else {
                return account.balance.toString();
            }
        } else {
            return '***';
        }
    }

    function getAccountSubAccountBalance(showAccountBalance: boolean, showHidden: boolean, account: Account, subAccountId: string): AccountDisplayBalance | null {
        if (account.type !== AccountType.MultiSubAccounts.type) {
            return null;
        }

        let resultCurrency = userStore.currentUserDefaultCurrency;

        if (!account.subAccounts || !account.subAccounts.length) {
            return {
                balance: showAccountBalance ? '0' : '***',
                currency: resultCurrency
            };
        }

        const allSubAccountCurrenciesMap: Record<string, boolean> = {};
        const allSubAccountCurrencies: string[] = [];
        let totalBalance = 0;

        for (let i = 0; i < account.subAccounts.length; i++) {
            const subAccount = account.subAccounts[i];

            if (!showHidden && subAccount.hidden) {
                continue;
            }

            if (!allSubAccountCurrenciesMap[subAccount.currency]) {
                allSubAccountCurrenciesMap[subAccount.currency] = true;
                allSubAccountCurrencies.push(subAccount.currency);
            }
        }

        if (allSubAccountCurrencies.length === 0) {
            return {
                balance: showAccountBalance ? '0' : '***',
                currency: resultCurrency
            };
        }

        if (allSubAccountCurrencies.length === 1) {
            resultCurrency = allSubAccountCurrencies[0];
        }

        let hasUnCalculatedAmount = false;

        for (let i = 0; i < account.subAccounts.length; i++) {
            const subAccount = account.subAccounts[i];

            if (!showHidden && subAccount.hidden) {
                continue;
            }

            if (subAccountId) {
                if (subAccountId === subAccount.id) {
                    return {
                        balance: showAccountBalance ? getAccountBalance(showAccountBalance, subAccount) as string : '***',
                        currency: subAccount.currency
                    };
                }
            }

            if (subAccount.currency === resultCurrency) {
                if (subAccount.isAsset) {
                    totalBalance += subAccount.balance;
                } else if (subAccount.isLiability) {
                    totalBalance -= subAccount.balance;
                } else {
                    totalBalance += subAccount.balance;
                }
            } else {
                const balance = exchangeRatesStore.getExchangedAmount(subAccount.balance, subAccount.currency, resultCurrency);

                if (!isNumber(balance)) {
                    hasUnCalculatedAmount = true;
                    continue;
                }

                if (subAccount.isAsset) {
                    totalBalance += Math.floor(balance);
                } else if (subAccount.isLiability) {
                    totalBalance -= Math.floor(balance);
                } else {
                    totalBalance += Math.floor(balance);
                }
            }
        }

        if (subAccountId) { // not found specified id in sub accounts
            return null;
        }

        let displayTotalBalance = totalBalance.toString();

        if (hasUnCalculatedAmount) {
            displayTotalBalance += '+';
        }

        return {
            balance: showAccountBalance ? displayTotalBalance : '***',
            currency: resultCurrency
        };
    }

    function hasAccount(accountCategory: AccountCategory, visibleOnly: boolean): boolean {
        if (!allCategorizedAccountsMap.value[accountCategory.type] ||
            !allCategorizedAccountsMap.value[accountCategory.type].accounts ||
            !allCategorizedAccountsMap.value[accountCategory.type].accounts.length) {
            return false;
        }

        let shownCount = 0;

        for (let i = 0; i < allCategorizedAccountsMap.value[accountCategory.type].accounts.length; i++) {
            const account = allCategorizedAccountsMap.value[accountCategory.type].accounts[i];

            if (!visibleOnly || !account.hidden) {
                shownCount++;
            }
        }

        return shownCount > 0;
    }

    function hasVisibleSubAccount(showHidden: boolean, account: Account): boolean {
        if (!account || account.type !== AccountType.MultiSubAccounts.type || !account.subAccounts) {
            return false;
        }

        for (let i = 0; i < account.subAccounts.length; i++) {
            if (showHidden || !account.subAccounts[i].hidden) {
                return true;
            }
        }

        return false;
    }

    function loadAllAccounts({ force }: { force: boolean }): Promise<Account[]> {
        if (!force && !accountListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allAccounts.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllAccounts({
                visibleOnly: false
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve account list' });
                    return;
                }

                if (accountListStateInvalid.value) {
                    updateAccountListInvalidState(false);
                }

                const accounts = Account.sortAccounts(Account.ofMulti(data.result));

                if (force && data.result && isEquals(allAccounts.value, accounts)) {
                    reject({ message: 'Account list is up to date', isUpToDate: true });
                    return;
                }

                loadAccountList(accounts);

                resolve(accounts);
            }).catch(error => {
                if (force) {
                    logger.error('failed to force load account list', error);
                } else {
                    logger.error('failed to load account list', error);
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve account list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getAccount({ accountId }: { accountId: string }): Promise<Account> {
        return new Promise((resolve, reject) => {
            services.getAccount({
                id: accountId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve account' });
                    return;
                }

                const account = Account.of(data.result);

                resolve(account);
            }).catch(error => {
                logger.error('failed to load account info', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve account' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function saveAccount({ account, subAccounts, isEdit, clientSessionId }: { account: Account, subAccounts: Account[], isEdit: boolean, clientSessionId: string }): Promise<Account> {
        return new Promise((resolve, reject) => {
            const oldAccount = isEdit ? allAccountsMap.value[account.id] : null;
            let promise = null;

            if (!isEdit) {
                promise = services.addAccount(account.toCreateRequest(clientSessionId, subAccounts));
            } else {
                promise = services.modifyAccount(account.toModifyRequest(clientSessionId, subAccounts));
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add account' });
                    } else {
                        reject({ message: 'Unable to save account' });
                    }
                    return;
                }

                const newAccount = Account.of(data.result);

                if (!isEdit) {
                    addAccountToAccountList(newAccount);
                } else {
                    if (oldAccount && oldAccount.category === newAccount.category) {
                        updateAccountToAccountList(oldAccount, newAccount);
                    } else {
                        updateAccountListInvalidState(true);
                    }
                }

                resolve(newAccount);
            }).catch(error => {
                logger.error('failed to save account', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add account' });
                    } else {
                        reject({ message: 'Unable to save account' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function changeAccountDisplayOrder({ accountId, from, to, updateListOrder, updateGlobalListOrder }: { accountId: string, from: number, to: number, updateListOrder: boolean, updateGlobalListOrder: boolean }): Promise<void> {
        const account = allAccountsMap.value[accountId];

        return new Promise((resolve, reject) => {
            if (!account ||
                !allCategorizedAccountsMap.value[account.category] ||
                !allCategorizedAccountsMap.value[account.category].accounts ||
                !allCategorizedAccountsMap.value[account.category].accounts[to]) {
                reject({ message: 'Unable to move account' });
                return;
            }

            if (!accountListStateInvalid.value) {
                updateAccountListInvalidState(true);
            }

            updateAccountDisplayOrderInAccountList({ account, from, to, updateListOrder, updateGlobalListOrder });

            resolve();
        });
    }

    function updateAccountDisplayOrders(): Promise<boolean> {
        const newDisplayOrders: AccountNewDisplayOrderRequest[] = [];

        for (const category in allCategorizedAccountsMap.value) {
            if (!Object.prototype.hasOwnProperty.call(allCategorizedAccountsMap.value, category)) {
                continue;
            }

            const accountList = allCategorizedAccountsMap.value[category].accounts;

            for (let i = 0; i < accountList.length; i++) {
                newDisplayOrders.push({
                    id: accountList[i].id,
                    displayOrder: i + 1
                });
            }
        }

        return new Promise((resolve, reject) => {
            services.moveAccount({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to move account' });
                    return;
                }

                if (accountListStateInvalid.value) {
                    updateAccountListInvalidState(false);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to save accounts display order', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to move account' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function hideAccount({ account, hidden }: { account: Account, hidden: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.hideAccount({
                id: account.id,
                hidden: hidden
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (hidden) {
                        reject({ message: 'Unable to hide this account' });
                    } else {
                        reject({ message: 'Unable to unhide this account' });
                    }

                    return;
                }

                updateAccountVisibilityInAccountList({ account, hidden });

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to change account visibility', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (hidden) {
                        reject({ message: 'Unable to hide this account' });
                    } else {
                        reject({ message: 'Unable to unhide this account' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteAccount({ account, beforeResolve }: { account: Account, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteAccount({
                id: account.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this account' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeAccountFromAccountList(account);
                    });
                } else {
                    removeAccountFromAccountList(account);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete account', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this account' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteSubAccount({ subAccount, beforeResolve }: { subAccount: Account, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteSubAccount({
                id: subAccount.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this sub-account' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeSubAccountFromAccountList(subAccount);
                    });
                } else {
                    removeSubAccountFromAccountList(subAccount);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete sub-account', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this sub-account' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // states
        allAccounts,
        allAccountsMap,
        allCategorizedAccountsMap,
        accountListStateInvalid,
        // computed states
        allPlainAccounts,
        allMixedPlainAccounts,
        allVisiblePlainAccounts,
        allAvailableAccountsCount,
        allVisibleAccountsCount,
        // functions
        updateAccountListInvalidState,
        resetAccounts,
        getFirstShowingIds,
        getLastShowingIds,
        getAccountStatementDate,
        getNetAssets,
        getTotalAssets,
        getTotalLiabilities,
        getAccountCategoryTotalBalance,
        getAccountBalance,
        getAccountSubAccountBalance,
        hasAccount,
        hasVisibleSubAccount,
        loadAllAccounts,
        getAccount,
        saveAccount,
        changeAccountDisplayOrder,
        updateAccountDisplayOrders,
        hideAccount,
        deleteAccount,
        deleteSubAccount
    }
});
