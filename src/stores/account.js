import { defineStore } from 'pinia';

import { useUserStore } from './user.ts';
import { useExchangeRatesStore } from './exchangeRates.js';

import { AccountType, AccountCategory } from '@/core/account.ts';
import { PARENT_ACCOUNT_CURRENCY_PLACEHOLDER } from '@/consts/currency.ts';
import { DEFAULT_ACCOUNT_ICON_ID } from '@/consts/icon.ts';
import { DEFAULT_ACCOUNT_COLOR } from '@/consts/color.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';
import { isNumber, isEquals } from '@/lib/common.ts';
import { getCurrentUnixTime } from '@/lib/datetime.ts';
import { getCategorizedAccountsMap, getAllFilteredAccountsBalance } from '@/lib/account.js';

function loadAccountList(state, accounts) {
    state.allAccounts = accounts;
    state.allAccountsMap = {};

    for (let i = 0; i < accounts.length; i++) {
        const account = accounts[i];
        state.allAccountsMap[account.id] = account;

        if (account.subAccounts) {
            for (let j = 0; j < account.subAccounts.length; j++) {
                const subAccount = account.subAccounts[j];
                state.allAccountsMap[subAccount.id] = subAccount;
            }
        }
    }

    state.allCategorizedAccountsMap = getCategorizedAccountsMap(accounts);
}

function addAccountToAccountList(state, account) {
    let insertIndexToAllList = 0;

    for (let i = 0; i < state.allAccounts.length; i++) {
        if (state.allAccounts[i].category > account.category) {
            insertIndexToAllList = i;
            break;
        }
    }

    state.allAccounts.splice(insertIndexToAllList, 0, account);

    state.allAccountsMap[account.id] = account;

    if (account.subAccounts) {
        for (let i = 0; i < account.subAccounts.length; i++) {
            const subAccount = account.subAccounts[i];
            state.allAccountsMap[subAccount.id] = subAccount;
        }
    }

    if (state.allCategorizedAccountsMap[account.category]) {
        const accountList = state.allCategorizedAccountsMap[account.category].accounts;
        accountList.push(account);
    } else {
        state.allCategorizedAccountsMap = getCategorizedAccountsMap(state.allAccounts);
    }
}

function updateAccountToAccountList(state, account) {
    for (let i = 0; i < state.allAccounts.length; i++) {
        if (state.allAccounts[i].id === account.id) {
            state.allAccounts.splice(i, 1, account);
            break;
        }
    }

    state.allAccountsMap[account.id] = account;

    if (account.subAccounts) {
        for (let i = 0; i < account.subAccounts.length; i++) {
            const subAccount = account.subAccounts[i];
            state.allAccountsMap[subAccount.id] = subAccount;
        }
    }

    if (state.allCategorizedAccountsMap[account.category]) {
        const accountList = state.allCategorizedAccountsMap[account.category].accounts;

        for (let i = 0; i < accountList.length; i++) {
            if (accountList[i].id === account.id) {
                accountList.splice(i, 1, account);
                break;
            }
        }
    }
}

function updateAccountDisplayOrderInAccountList(state, { account, from, to, updateListOrder, updateGlobalListOrder }) {
    let fromAccount = null;
    let toAccount = null;

    if (state.allCategorizedAccountsMap[account.category]) {
        const accountList = state.allCategorizedAccountsMap[account.category].accounts;

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

        for (let i = 0; i < state.allAccounts.length; i++) {
            if (state.allAccounts[i].id === fromAccount.id) {
                globalFromIndex = i;
            } else if (state.allAccounts[i].id === toAccount.id) {
                globalToIndex = i;
            }
        }

        if (globalFromIndex >= 0 && globalToIndex >= 0) {
            state.allAccounts.splice(globalToIndex, 0, state.allAccounts.splice(globalFromIndex, 1)[0]);
        }
    }
}

function updateAccountVisibilityInAccountList(state, { account, hidden }) {
    if (state.allAccountsMap[account.id]) {
        state.allAccountsMap[account.id].hidden = hidden;
    }
}

function removeAccountFromAccountList(state, account) {
    for (let i = 0; i < state.allAccounts.length; i++) {
        if (state.allAccounts[i].id === account.id) {
            state.allAccounts.splice(i, 1);
            break;
        }
    }

    if (state.allAccountsMap[account.id] && state.allAccountsMap[account.id].subAccounts) {
        const subAccounts = state.allAccountsMap[account.id].subAccounts;

        for (let i = 0; i < subAccounts.length; i++) {
            const subAccount = subAccounts[i];
            if (state.allAccountsMap[subAccount.id]) {
                delete state.allAccountsMap[subAccount.id];
            }
        }
    }

    if (state.allAccountsMap[account.id]) {
        delete state.allAccountsMap[account.id];
    }

    if (state.allCategorizedAccountsMap[account.category]) {
        const accountList = state.allCategorizedAccountsMap[account.category].accounts;

        for (let i = 0; i < accountList.length; i++) {
            if (accountList[i].id === account.id) {
                accountList.splice(i, 1);
                break;
            }
        }
    }
}

export const useAccountsStore = defineStore('accounts', {
    state: () => ({
        allAccounts: [],
        allAccountsMap: {},
        allCategorizedAccountsMap: {},
        accountListStateInvalid: true,
    }),
    getters: {
        allPlainAccounts(state) {
            const allAccounts = [];

            for (let i = 0; i < state.allAccounts.length; i++) {
                const account = state.allAccounts[i];

                if (account.type === AccountType.SingleAccount.type) {
                    allAccounts.push(account);
                } else if (account.type === AccountType.MultiSubAccounts.type) {
                    if (account.subAccounts) {
                        for (let j = 0; j < account.subAccounts.length; j++) {
                            const subAccount = account.subAccounts[j];
                            allAccounts.push(subAccount);
                        }
                    }
                }
            }

            return allAccounts;
        },
        allVisiblePlainAccounts(state) {
            const allVisibleAccounts = [];

            for (let i = 0; i < state.allAccounts.length; i++) {
                const account = state.allAccounts[i];

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

            return allVisibleAccounts;
        },
        allAvailableAccountsCount(state) {
            let allAccountCount = 0;

            for (let category in state.allCategorizedAccountsMap) {
                if (!Object.prototype.hasOwnProperty.call(state.allCategorizedAccountsMap, category)) {
                    continue;
                }

                allAccountCount += state.allCategorizedAccountsMap[category].accounts.length;
            }

            return allAccountCount;
        },
        allVisibleAccountsCount(state) {
            let shownAccountCount = 0;

            for (let category in state.allCategorizedAccountsMap) {
                if (!Object.prototype.hasOwnProperty.call(state.allCategorizedAccountsMap, category)) {
                    continue;
                }

                const accountList = state.allCategorizedAccountsMap[category].accounts;

                for (let i = 0; i < accountList.length; i++) {
                    if (!accountList[i].hidden) {
                        shownAccountCount++;
                    }
                }
            }

            return shownAccountCount;
        }
    },
    actions: {
        generateNewAccountModel() {
            const userStore = useUserStore();
            const now = getCurrentUnixTime();

            return {
                category: AccountCategory.Cash.type,
                type: AccountType.SingleAccount.type,
                name: '',
                icon: DEFAULT_ACCOUNT_ICON_ID,
                color: DEFAULT_ACCOUNT_COLOR,
                currency: userStore.currentUserDefaultCurrency,
                balance: 0,
                balanceTime: now,
                comment: '',
                creditCardStatementDate: 0,
                visible: true
            };
        },
        generateNewSubAccountModel(parentAccount) {
            const userStore = useUserStore();
            const now = getCurrentUnixTime();

            return {
                category: null,
                type: null,
                name: '',
                icon: parentAccount.icon,
                color: parentAccount.color,
                currency: userStore.currentUserDefaultCurrency,
                balance: 0,
                balanceTime: now,
                comment: '',
                visible: true
            };
        },
        updateAccountListInvalidState(invalidState) {
            this.accountListStateInvalid = invalidState;
        },
        resetAccounts() {
            this.allAccounts = [];
            this.allAccountsMap = {};
            this.allCategorizedAccountsMap = {};
            this.accountListStateInvalid = true;
        },
        getFirstShowingIds(showHidden) {
            const ret = {
                accounts: {},
                subAccounts: {}
            };

            for (let category in this.allCategorizedAccountsMap) {
                if (!Object.prototype.hasOwnProperty.call(this.allCategorizedAccountsMap, category)) {
                    continue;
                }

                if (!this.allCategorizedAccountsMap[category] || !this.allCategorizedAccountsMap[category].accounts) {
                    continue;
                }

                const accounts = this.allCategorizedAccountsMap[category].accounts;

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
        },
        getLastShowingIds(showHidden) {
            const ret = {
                accounts: {},
                subAccounts: {}
            };

            for (let category in this.allCategorizedAccountsMap) {
                if (!Object.prototype.hasOwnProperty.call(this.allCategorizedAccountsMap, category)) {
                    continue;
                }

                if (!this.allCategorizedAccountsMap[category] || !this.allCategorizedAccountsMap[category].accounts) {
                    continue;
                }

                const accounts = this.allCategorizedAccountsMap[category].accounts;

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
        },
        getAccountStatementDate(accountId) {
            if (!accountId) {
                return null;
            }

            const accountIds = accountId.split(',');
            let mainAccount = null;

            for (let i = 0; i < accountIds.length; i++) {
                const id = accountIds[i];
                let account = this.allAccountsMap[id];

                if (!account) {
                    return null;
                }

                if (account.parentId !== '0') {
                    account = this.allAccountsMap[account.parentId];
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
        },
        getNetAssets(showAccountBalance) {
            if (!showAccountBalance) {
                return '***';
            }

            const userStore = useUserStore();
            const exchangeRatesStore = useExchangeRatesStore();
            const accountsBalance = getAllFilteredAccountsBalance(this.allCategorizedAccountsMap, () => true);
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
                return netAssets;
            }
        },
        getTotalAssets(showAccountBalance) {
            if (!showAccountBalance) {
                return '***';
            }

            const userStore = useUserStore();
            const exchangeRatesStore = useExchangeRatesStore();
            const accountsBalance = getAllFilteredAccountsBalance(this.allCategorizedAccountsMap, account => account.isAsset);
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
                return totalAssets;
            }
        },
        getTotalLiabilities(showAccountBalance) {
            if (!showAccountBalance) {
                return '***';
            }

            const userStore = useUserStore();
            const exchangeRatesStore = useExchangeRatesStore();
            const accountsBalance = getAllFilteredAccountsBalance(this.allCategorizedAccountsMap, account => account.isLiability);
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
                return totalLiabilities;
            }
        },
        getAccountCategoryTotalBalance(showAccountBalance, accountCategory) {
            if (!showAccountBalance) {
                return '***';
            }

            const userStore = useUserStore();
            const exchangeRatesStore = useExchangeRatesStore();
            const accountsBalance = getAllFilteredAccountsBalance(this.allCategorizedAccountsMap, account => account.category === accountCategory.type);
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
                return totalBalance;
            }
        },
        getAccountBalance(showAccountBalance, account) {
            if (account.type !== AccountType.SingleAccount.type) {
                return null;
            }

            if (showAccountBalance) {
                if (account.isAsset) {
                    return account.balance;
                } else if (account.isLiability) {
                    return -account.balance;
                } else {
                    return account.balance;
                }
            } else {
                return '***';
            }
        },
        getAccountSubAccountBalance(showAccountBalance, showHidden, account, subAccountId) {
            if (account.type !== AccountType.MultiSubAccounts.type) {
                return null;
            }

            const userStore = useUserStore();
            const exchangeRatesStore = useExchangeRatesStore();
            let resultCurrency = userStore.currentUserDefaultCurrency;

            if (!account.subAccounts || !account.subAccounts.length) {
                return {
                    balance: showAccountBalance ? 0 : '***',
                    currency: resultCurrency
                };
            }

            const allSubAccountCurrenciesMap = {};
            const allSubAccountCurrencies = [];
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
                    balance: showAccountBalance ? 0 : '***',
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
                            balance: showAccountBalance ? this.getAccountBalance(showAccountBalance, subAccount) : '***',
                            currency: subAccount.currency
                        };
                    }
                }

                if (subAccount === resultCurrency) {
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

            if (hasUnCalculatedAmount) {
                totalBalance += '+';
            }

            return {
                balance: showAccountBalance ? totalBalance : '***',
                currency: resultCurrency
            };
        },
        hasAccount(accountCategory, visibleOnly) {
            if (!this.allCategorizedAccountsMap[accountCategory.type] ||
                !this.allCategorizedAccountsMap[accountCategory.type].accounts ||
                !this.allCategorizedAccountsMap[accountCategory.type].accounts.length) {
                return false;
            }

            let shownCount = 0;

            for (let i = 0; i < this.allCategorizedAccountsMap[accountCategory.type].accounts.length; i++) {
                const account = this.allCategorizedAccountsMap[accountCategory.type].accounts[i];

                if (!visibleOnly || !account.hidden) {
                    shownCount++;
                }
            }

            return shownCount > 0;
        },
        hasVisibleSubAccount(showHidden, account) {
            if (!account || account.type !== AccountType.MultiSubAccounts.type || !account.subAccounts) {
                return false;
            }

            for (let i = 0; i < account.subAccounts.length; i++) {
                if (showHidden || !account.subAccounts[i].hidden) {
                    return true;
                }
            }

            return false;
        },
        loadAllAccounts({ force }) {
            const self = this;

            if (!force && !self.accountListStateInvalid) {
                return new Promise((resolve) => {
                    resolve(self.allAccounts);
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

                    if (self.accountListStateInvalid) {
                        self.updateAccountListInvalidState(false);
                    }

                    if (force && data.result && isEquals(self.allAccounts, data.result)) {
                        reject({ message: 'Account list is up to date' });
                        return;
                    }

                    loadAccountList(self, data.result);

                    resolve(data.result);
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
        },
        getAccount({ accountId }) {
            return new Promise((resolve, reject) => {
                services.getAccount({
                    id: accountId
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to retrieve account' });
                        return;
                    }

                    resolve(data.result);
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
        },
        saveAccount({ account, subAccounts, isEdit, clientSessionId }) {
            const self = this;

            const submitSubAccounts = [];

            if (account.type === AccountType.MultiSubAccounts.type) {
                for (let i = 0; i < subAccounts.length; i++) {
                    const subAccount = subAccounts[i];
                    const submitAccount = {
                        category: account.category,
                        type: AccountType.SingleAccount.type,
                        name: subAccount.name,
                        icon: subAccount.icon,
                        color: subAccount.color,
                        currency: subAccount.currency,
                        balance: subAccount.balance,
                        comment: subAccount.comment
                    };

                    if (isEdit) {
                        submitAccount.id = subAccount.id;
                        submitAccount.hidden = !subAccount.visible;
                    } else {
                        submitAccount.balanceTime = subAccount.balanceTime;
                    }

                    submitSubAccounts.push(submitAccount);
                }
            }

            const submitAccount = {
                category: account.category,
                type: account.type,
                name: account.name,
                icon: account.icon,
                color: account.color,
                currency: account.type === AccountType.SingleAccount.type ? account.currency : PARENT_ACCOUNT_CURRENCY_PLACEHOLDER,
                balance: account.type === AccountType.SingleAccount.type ? account.balance : 0,
                comment: account.comment,
                subAccounts: account.type === AccountType.SingleAccount.type ? null : submitSubAccounts,
            };

            if (account.category === AccountCategory.CreditCard.type) {
                submitAccount.creditCardStatementDate = account.creditCardStatementDate;
            }

            if (clientSessionId) {
                submitAccount.clientSessionId = clientSessionId;
            }

            if (isEdit) {
                submitAccount.id = account.id;
                submitAccount.hidden = !account.visible;
            } else {
                if (account.type === AccountType.SingleAccount.type) {
                    submitAccount.balanceTime = account.balanceTime;
                }
            }

            const oldAccount = submitAccount.id ? self.allAccountsMap[submitAccount.id] : null;

            return new Promise((resolve, reject) => {
                let promise = null;

                if (!submitAccount.id) {
                    promise = services.addAccount(submitAccount);
                } else {
                    promise = services.modifyAccount(submitAccount);
                }

                promise.then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (!submitAccount.id) {
                            reject({ message: 'Unable to add account' });
                        } else {
                            reject({ message: 'Unable to save account' });
                        }
                        return;
                    }

                    if (!submitAccount.id) {
                        addAccountToAccountList(self, data.result);
                    } else {
                        if (oldAccount && oldAccount.category === data.result.category) {
                            updateAccountToAccountList(self, data.result);
                        } else {
                            self.updateAccountListInvalidState(true);
                        }
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to save account', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        if (!submitAccount.id) {
                            reject({ message: 'Unable to add account' });
                        } else {
                            reject({ message: 'Unable to save account' });
                        }
                    } else {
                        reject(error);
                    }
                });
            });
        },
        changeAccountDisplayOrder({ accountId, from, to, updateListOrder, updateGlobalListOrder }) {
            const self = this;
            const account = self.allAccountsMap[accountId];

            return new Promise((resolve, reject) => {
                if (!account ||
                    !self.allCategorizedAccountsMap[account.category] ||
                    !self.allCategorizedAccountsMap[account.category].accounts ||
                    !self.allCategorizedAccountsMap[account.category].accounts[to]) {
                    reject({ message: 'Unable to move account' });
                    return;
                }

                if (!self.accountListStateInvalid) {
                    self.updateAccountListInvalidState(true);
                }

                updateAccountDisplayOrderInAccountList(self, {
                    account: account,
                    from: from,
                    to: to,
                    updateListOrder: updateListOrder,
                    updateGlobalListOrder: updateGlobalListOrder
                });

                resolve();
            });
        },
        updateAccountDisplayOrders() {
            const self = this;
            const newDisplayOrders = [];

            for (let category in self.allCategorizedAccountsMap) {
                if (!Object.prototype.hasOwnProperty.call(self.allCategorizedAccountsMap, category)) {
                    continue;
                }

                const accountList = self.allCategorizedAccountsMap[category].accounts;

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

                    if (self.accountListStateInvalid) {
                        self.updateAccountListInvalidState(false);
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
        },
        hideAccount({ account, hidden }) {
            const self = this;

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

                    updateAccountVisibilityInAccountList(self, {
                        account: account,
                        hidden: hidden
                    });

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
        },
        deleteAccount({ account, beforeResolve }) {
            const self = this;

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
                            removeAccountFromAccountList(self, account);
                        });
                    } else {
                        removeAccountFromAccountList(self, account);
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
    }
});
