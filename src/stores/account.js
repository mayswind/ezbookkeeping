import { defineStore } from 'pinia';

import { useUserStore } from './user.js';
import { useExchangeRatesStore } from './exchangeRates.js';

import accountConstants from '@/consts/account.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';
import { isNumber, isEquals } from '@/lib/common.js';
import { getCategorizedAccounts, getAllFilteredAccountsBalance } from '@/lib/account.js';

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

    state.allCategorizedAccounts = getCategorizedAccounts(accounts);
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

    if (state.allCategorizedAccounts[account.category]) {
        const accountList = state.allCategorizedAccounts[account.category].accounts;
        accountList.push(account);
    } else {
        state.allCategorizedAccounts = getCategorizedAccounts(state.allAccounts);
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

    if (state.allCategorizedAccounts[account.category]) {
        const accountList = state.allCategorizedAccounts[account.category].accounts;

        for (let i = 0; i < accountList.length; i++) {
            if (accountList[i].id === account.id) {
                accountList.splice(i, 1, account);
                break;
            }
        }
    }
}

function updateAccountDisplayOrderInAccountList(state, { account, from, to }) {
    let fromAccount = null;
    let toAccount = null;

    if (state.allCategorizedAccounts[account.category]) {
        const accountList = state.allCategorizedAccounts[account.category].accounts;
        fromAccount = accountList[from];
        toAccount = accountList[to];

        accountList.splice(to, 0, accountList.splice(from, 1)[0]);
    }

    if (fromAccount && toAccount) {
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

    if (state.allCategorizedAccounts[account.category]) {
        const accountList = state.allCategorizedAccounts[account.category].accounts;

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
        allCategorizedAccounts: {},
        accountListStateInvalid: true,
    }),
    getters: {
        allPlainAccounts(state) {
            const allAccounts = [];

            for (let i = 0; i < state.allAccounts.length; i++) {
                const account = state.allAccounts[i];

                if (account.type === accountConstants.allAccountTypes.SingleAccount) {
                    allAccounts.push(account);
                } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts) {
                    for (let j = 0; j < account.subAccounts.length; j++) {
                        const subAccount = account.subAccounts[j];
                        allAccounts.push(subAccount);
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

                if (account.type === accountConstants.allAccountTypes.SingleAccount) {
                    allVisibleAccounts.push(account);
                } else if (account.type === accountConstants.allAccountTypes.MultiSubAccounts) {
                    for (let j = 0; j < account.subAccounts.length; j++) {
                        const subAccount = account.subAccounts[j];
                        allVisibleAccounts.push(subAccount);
                    }
                }
            }

            return allVisibleAccounts;
        },
        allAvailableAccountsCount(state) {
            let allAccountCount = 0;

            for (let category in state.allCategorizedAccounts) {
                if (!Object.prototype.hasOwnProperty.call(state.allCategorizedAccounts, category)) {
                    continue;
                }

                allAccountCount += state.allCategorizedAccounts[category].accounts.length;
            }

            return allAccountCount;
        },
        allVisibleAccountsCount(state) {
            let shownAccountCount = 0;

            for (let category in state.allCategorizedAccounts) {
                if (!Object.prototype.hasOwnProperty.call(state.allCategorizedAccounts, category)) {
                    continue;
                }

                const accountList = state.allCategorizedAccounts[category].accounts;

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
        updateAccountListInvalidState(invalidState) {
            this.accountListStateInvalid = invalidState;
        },
        resetAccounts() {
            this.allAccounts = [];
            this.allAccountsMap = {};
            this.allCategorizedAccounts = {};
            this.accountListStateInvalid = true;
        },
        getFirstShowingIds(showHidden) {
            const ret = {
                accounts: {},
                subAccounts: {}
            };

            for (let category in this.allCategorizedAccounts) {
                if (!Object.prototype.hasOwnProperty.call(this.allCategorizedAccounts, category)) {
                    continue;
                }

                if (!this.allCategorizedAccounts[category] || !this.allCategorizedAccounts[category].accounts) {
                    continue;
                }

                const accounts = this.allCategorizedAccounts[category].accounts;

                for (let i = 0; i < accounts.length; i++) {
                    const account = accounts[i];

                    if (account.type === accountConstants.allAccountTypes.MultiSubAccounts && account.subAccounts) {
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

            for (let category in this.allCategorizedAccounts) {
                if (!Object.prototype.hasOwnProperty.call(this.allCategorizedAccounts, category)) {
                    continue;
                }

                if (!this.allCategorizedAccounts[category] || !this.allCategorizedAccounts[category].accounts) {
                    continue;
                }

                const accounts = this.allCategorizedAccounts[category].accounts;

                for (let i = accounts.length - 1; i >= 0; i--) {
                    const account = accounts[i];

                    if (account.type === accountConstants.allAccountTypes.MultiSubAccounts && account.subAccounts) {
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
        getNetAssets(showAccountBalance) {
            if (!showAccountBalance) {
                return '***';
            }

            const userStore = useUserStore();
            const exchangeRatesStore = useExchangeRatesStore();
            const accountsBalance = getAllFilteredAccountsBalance(this.allCategorizedAccounts, () => true);
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
            const accountsBalance = getAllFilteredAccountsBalance(this.allCategorizedAccounts, account => account.isAsset);
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
            const accountsBalance = getAllFilteredAccountsBalance(this.allCategorizedAccounts, account => account.isLiability);
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
            const accountsBalance = getAllFilteredAccountsBalance(this.allCategorizedAccounts, account => account.category === accountCategory.id);
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
            if (account.type !== accountConstants.allAccountTypes.SingleAccount) {
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
            if (account.type !== accountConstants.allAccountTypes.MultiSubAccounts) {
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

            return {
                balance: showAccountBalance ? totalBalance : '***',
                currency: resultCurrency
            };
        },
        hasAccount(accountCategory, visibleOnly) {
            if (!this.allCategorizedAccounts[accountCategory.id] ||
                !this.allCategorizedAccounts[accountCategory.id].accounts ||
                !this.allCategorizedAccounts[accountCategory.id].accounts.length) {
                return false;
            }

            let shownCount = 0;

            for (let i = 0; i < this.allCategorizedAccounts[accountCategory.id].accounts.length; i++) {
                const account = this.allCategorizedAccounts[accountCategory.id].accounts[i];

                if (!visibleOnly || !account.hidden) {
                    shownCount++;
                }
            }

            return shownCount > 0;
        },
        hasVisibleSubAccount(showHidden, account) {
            if (!account || account.type !== accountConstants.allAccountTypes.MultiSubAccounts || !account.subAccounts) {
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
                        reject({ message: 'Unable to get account list' });
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
                        reject({ message: 'Unable to get account list' });
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
                        reject({ message: 'Unable to get account' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to load account info', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to get account' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        saveAccount({ account }) {
            const self = this;
            const oldAccount = account.id ? self.allAccountsMap[account.id] : null;

            return new Promise((resolve, reject) => {
                let promise = null;

                if (!account.id) {
                    promise = services.addAccount(account);
                } else {
                    promise = services.modifyAccount(account);
                }

                promise.then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (!account.id) {
                            reject({ message: 'Unable to add account' });
                        } else {
                            reject({ message: 'Unable to save account' });
                        }
                        return;
                    }

                    if (!account.id) {
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
                        if (!account.id) {
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
        changeAccountDisplayOrder({ accountId, from, to }) {
            const self = this;
            const account = self.allAccountsMap[accountId];

            return new Promise((resolve, reject) => {
                if (!account ||
                    !self.allCategorizedAccounts[account.category] ||
                    !self.allCategorizedAccounts[account.category].accounts ||
                    !self.allCategorizedAccounts[account.category].accounts[to]) {
                    reject({ message: 'Unable to move account' });
                    return;
                }

                if (!self.accountListStateInvalid) {
                    self.updateAccountListInvalidState(true);
                }

                updateAccountDisplayOrderInAccountList(self, {
                    account: account,
                    from: from,
                    to: to
                });

                resolve();
            });
        },
        updateAccountDisplayOrders() {
            const self = this;
            const newDisplayOrders = [];

            for (let category in self.allCategorizedAccounts) {
                if (!Object.prototype.hasOwnProperty.call(self.allCategorizedAccounts, category)) {
                    continue;
                }

                const accountList = self.allCategorizedAccounts[category].accounts;

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
