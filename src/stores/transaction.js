import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.js';
import { useAccountsStore } from './account.js';
import { useTransactionCategoriesStore } from './transactionCategory.js';
import { useOverviewStore } from './overview.js';
import { useStatisticsStore } from './statistics.js';
import { useExchangeRatesStore } from './exchangeRates.js';

import datetimeConstants from '@/consts/datetime.js';
import transactionConstants from '@/consts/transaction.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';
import { isNumber, isString } from '@/lib/common.js';
import {
    getCurrentUnixTime,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore,
    parseDateFromUnixTime,
    getShortDate,
    getYear,
    getMonth,
    getYearAndMonth,
    getDay,
    getDayOfWeekName
} from '@/lib/datetime.js';
import {
    stringCurrencyToNumeric
} from '@/lib/currency.js';

const emptyTransactionResult = {
    items: [],
    transactionsNextTimeId: 0
};

function loadTransactionList(state, settingsStore, exchangeRatesStore, { transactions, reload, autoExpand, defaultCurrency }) {
    if (reload) {
        state.transactions = [];
    }

    if (transactions.items && transactions.items.length) {
        const currentUtcOffset = getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone);
        let currentMonthListIndex = -1;
        let currentMonthList = null;

        for (let i = 0; i < transactions.items.length; i++) {
            const item = transactions.items[i];
            fillTransactionObject(state, item, currentUtcOffset);

            const transactionTime = parseDateFromUnixTime(item.time, item.utcOffset, currentUtcOffset);
            const transactionYear = getYear(transactionTime);
            const transactionMonth = getMonth(transactionTime);
            const transactionYearMonth = getYearAndMonth(transactionTime);

            if (currentMonthList && currentMonthList.year === transactionYear && currentMonthList.month === transactionMonth) {
                currentMonthList.items.push(Object.freeze(item));
                continue;
            }

            for (let j = currentMonthListIndex + 1; j < state.transactions.length; j++) {
                if (state.transactions[j].year === transactionYear && state.transactions[j].month === transactionMonth) {
                    currentMonthListIndex = j;
                    currentMonthList = state.transactions[j];
                    break;
                }
            }

            if (!currentMonthList || currentMonthList.year !== transactionYear || currentMonthList.month !== transactionMonth) {
                calculateMonthTotalAmount(state, exchangeRatesStore, currentMonthList, defaultCurrency, state.transactionsFilter.accountId, false);

                state.transactions.push({
                    year: transactionYear,
                    month: transactionMonth,
                    yearMonth: transactionYearMonth,
                    opened: autoExpand,
                    items: []
                });

                currentMonthListIndex = state.transactions.length - 1;
                currentMonthList = state.transactions[state.transactions.length - 1];
            }

            currentMonthList.items.push(Object.freeze(item));
            calculateMonthTotalAmount(state, exchangeRatesStore, currentMonthList, defaultCurrency, state.transactionsFilter.accountId, true);
        }
    }

    if (transactions.nextTimeSequenceId) {
        state.transactionsNextTimeId = transactions.nextTimeSequenceId;
    } else {
        calculateMonthTotalAmount(state, exchangeRatesStore, state.transactions[state.transactions.length - 1], defaultCurrency, state.transactionsFilter.accountId, false);
        state.transactionsNextTimeId = -1;
    }
}

function updateTransactionInTransactionList(state, settingsStore, exchangeRatesStore, { transaction, defaultCurrency }) {
    const currentUtcOffset = getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone);
    const transactionTime = parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentUtcOffset);
    const transactionYear = getYear(transactionTime);
    const transactionMonth = getMonth(transactionTime);

    for (let i = 0; i < state.transactions.length; i++) {
        const transactionMonthList = state.transactions[i];

        if (!transactionMonthList.items) {
            continue;
        }

        for (let j = 0; j < transactionMonthList.items.length; j++) {
            if (transactionMonthList.items[j].id === transaction.id) {
                fillTransactionObject(state, transaction, currentUtcOffset);

                if (transactionYear !== transactionMonthList.year ||
                    transactionMonth !== transactionMonthList.month ||
                    transaction.day !== transactionMonthList.items[j].day) {
                    state.transactionListStateInvalid = true;
                    return;
                }

                if ((state.transactionsFilter.categoryId && state.transactionsFilter.categoryId !== '0' && state.transactionsFilter.categoryId !== transaction.categoryId) ||
                    (state.transactionsFilter.accountId && state.transactionsFilter.accountId !== '0' &&
                        state.transactionsFilter.accountId !== transaction.sourceAccountId &&
                        state.transactionsFilter.accountId !== transaction.destinationAccountId &&
                        (!transaction.sourceAccount || state.transactionsFilter.accountId !== transaction.sourceAccount.parentId) &&
                        (!transaction.destinationAccount || state.transactionsFilter.accountId !== transaction.destinationAccount.parentId)
                    )
                ) {
                    transactionMonthList.items.splice(j, 1);
                } else {
                    transactionMonthList.items.splice(j, 1, transaction);
                }

                if (transactionMonthList.items.length < 1) {
                    state.transactions.splice(i, 1);
                } else {
                    calculateMonthTotalAmount(state, exchangeRatesStore, transactionMonthList, defaultCurrency, state.transactionsFilter.accountId, i >= state.transactions.length - 1 && state.transactionsNextTimeId > 0);
                }

                return;
            }
        }
    }
}

function removeTransactionFromTransactionList(state, exchangeRatesStore, { transaction, defaultCurrency }) {
    for (let i = 0; i < state.transactions.length; i++) {
        const transactionMonthList = state.transactions[i];

        if (!transactionMonthList.items ||
            transactionMonthList.items[0].time < transaction.time ||
            transactionMonthList.items[transactionMonthList.items.length - 1].time > transaction.time) {
            continue;
        }

        for (let j = 0; j < transactionMonthList.items.length; j++) {
            if (transactionMonthList.items[j].id === transaction.id) {
                transactionMonthList.items.splice(j, 1);
            }
        }

        if (transactionMonthList.items.length < 1) {
            state.transactions.splice(i, 1);
        } else {
            calculateMonthTotalAmount(state, exchangeRatesStore, transactionMonthList, defaultCurrency, state.transactionsFilter.accountId, i >= state.transactions.length - 1 && state.transactionsNextTimeId > 0);
        }
    }
}

function calculateMonthTotalAmount(state, exchangeRatesStore, transactionMonthList, defaultCurrency, accountId, incomplete) {
    if (!transactionMonthList) {
        return;
    }

    let totalExpense = 0;
    let totalIncome = 0;
    let hasUnCalculatedTotalExpense = false;
    let hasUnCalculatedTotalIncome = false;

    for (let i = 0; i < transactionMonthList.items.length; i++) {
        const transaction = transactionMonthList.items[i];

        let amount = transaction.sourceAmount;
        let account = transaction.sourceAccount;

        if (accountId && transaction.destinationAccount && (transaction.destinationAccount.id === accountId || transaction.destinationAccount.parentId === accountId)) {
            amount = transaction.destinationAmount;
            account = transaction.destinationAccount;
        }

        if (!account) {
            continue;
        }

        if (account.currency !== defaultCurrency) {
            const balance = exchangeRatesStore.getExchangedAmount(amount, account.currency, defaultCurrency);

            if (!isNumber(balance)) {
                if (transaction.type === transactionConstants.allTransactionTypes.Expense) {
                    hasUnCalculatedTotalExpense = true;
                } else if (transaction.type === transactionConstants.allTransactionTypes.Income) {
                    hasUnCalculatedTotalIncome = true;
                }

                continue;
            }

            amount = balance;
        }

        if (transaction.type === transactionConstants.allTransactionTypes.Expense) {
            totalExpense += amount;
        } else if (transaction.type === transactionConstants.allTransactionTypes.Income) {
            totalIncome += amount;
        } else if (transaction.type === transactionConstants.allTransactionTypes.Transfer && accountId && accountId !== '0') {
            if (accountId === transaction.sourceAccountId) {
                totalExpense += amount;
            } else if (accountId === transaction.destinationAccountId) {
                totalIncome += amount;
            } else if (transaction.sourceAccount && accountId === transaction.sourceAccount.parentId &&
                transaction.destinationAccount && accountId === transaction.destinationAccount.parentId) {
                // Do Nothing
            } else if (transaction.sourceAccount && accountId === transaction.sourceAccount.parentId) {
                totalExpense += amount;
            } else if (transaction.destinationAccount && accountId === transaction.destinationAccount.parentId) {
                totalIncome += amount;
            }
        }
    }

    transactionMonthList.totalAmount = {
        expense: Math.floor(totalExpense),
        incompleteExpense: incomplete || hasUnCalculatedTotalExpense,
        income: Math.floor(totalIncome),
        incompleteIncome: incomplete || hasUnCalculatedTotalIncome
    };
}

function fillTransactionObject(state, transaction, currentUtcOffset) {
    if (!transaction) {
        return;
    }

    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionTime = parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentUtcOffset);

    transaction.date = getShortDate(transactionTime);
    transaction.day = getDay(transactionTime);
    transaction.dayOfWeek = getDayOfWeekName(transactionTime);

    if (transaction.sourceAccountId) {
        transaction.sourceAccount = accountsStore.allAccountsMap[transaction.sourceAccountId];
    }

    if (transaction.destinationAccountId) {
        transaction.destinationAccount = accountsStore.allAccountsMap[transaction.destinationAccountId];
    }

    if (transaction.categoryId) {
        transaction.category = transactionCategoriesStore.allTransactionCategoriesMap[transaction.categoryId];
    }

    return transaction;
}

export const useTransactionsStore = defineStore('transactions', {
    state: () => ({
        transactionsFilter: {
            dateType: datetimeConstants.allDateRanges.All.type,
            maxTime: 0,
            minTime: 0,
            type: 0,
            categoryId: '0',
            accountId: '0',
            keyword: ''
        },
        transactions: [],
        transactionsNextTimeId: 0,
        transactionListStateInvalid: true,
    }),
    getters: {
        noTransaction(state) {
            for (let i = 0; i < state.transactions.length; i++) {
                const transactionMonthList = state.transactions[i];

                for (let j = 0; j < transactionMonthList.items.length; j++) {
                    if (transactionMonthList.items[j]) {
                        return false;
                    }
                }
            }

            return true;
        },
        hasMoreTransaction(state) {
            return state.transactionsNextTimeId > 0;
        }
    },
    actions: {
        generateNewTransactionModel(type) {
            const settingsStore = useSettingsStore();
            const now = getCurrentUnixTime();
            const currentTimezone = settingsStore.appSettings.timeZone;

            let defaultType = transactionConstants.allTransactionTypes.Expense;

            if (type === transactionConstants.allTransactionTypes.Income.toString()) {
                defaultType = transactionConstants.allTransactionTypes.Income;
            } else if (type === transactionConstants.allTransactionTypes.Transfer.toString()) {
                defaultType = transactionConstants.allTransactionTypes.Transfer;
            }

            return {
                type: defaultType,
                time: now,
                timeZone: currentTimezone,
                utcOffset: getTimezoneOffsetMinutes(currentTimezone),
                expenseCategory: '',
                incomeCategory: '',
                transferCategory: '',
                sourceAccountId: '',
                destinationAccountId: '',
                sourceAmount: 0,
                destinationAmount: 0,
                hideAmount: false,
                tagIds: [],
                comment: '',
                geoLocation: null
            };
        },
        setTransactionSuitableDestinationAmount(transaction, oldValue, newValue) {
            const accountsStore = useAccountsStore();
            const exchangeRatesStore = useExchangeRatesStore();

            if (transaction.type === transactionConstants.allTransactionTypes.Expense || transaction.type === transactionConstants.allTransactionTypes.Income) {
                transaction.destinationAmount = newValue;
            } else if (transaction.type === transactionConstants.allTransactionTypes.Transfer) {
                const sourceAccount = accountsStore.allAccountsMap[transaction.sourceAccountId];
                const destinationAccount = accountsStore.allAccountsMap[transaction.destinationAccountId];

                if (sourceAccount && destinationAccount && sourceAccount.currency !== destinationAccount.currency) {
                    const exchangedOldValue = exchangeRatesStore.getExchangedAmount(oldValue, sourceAccount.currency, destinationAccount.currency);
                    const exchangedNewValue = exchangeRatesStore.getExchangedAmount(newValue, sourceAccount.currency, destinationAccount.currency);

                    if (isNumber(exchangedOldValue)) {
                        oldValue = Math.floor(exchangedOldValue);
                    }

                    if (isNumber(exchangedNewValue)) {
                        newValue = Math.floor(exchangedNewValue);
                    }
                }

                if ((!sourceAccount || !destinationAccount || transaction.destinationAmount === oldValue || transaction.destinationAmount === 0) &&
                    (stringCurrencyToNumeric(transactionConstants.minAmount) <= newValue &&
                        newValue <= stringCurrencyToNumeric(transactionConstants.maxAmount))) {
                    transaction.destinationAmount = newValue;
                }
            }
        },
        updateTransactionListInvalidState(invalidState) {
            this.transactionListStateInvalid = invalidState;
        },
        resetTransactions() {
            this.transactionsFilter.dateType = datetimeConstants.allDateRanges.All.type;
            this.transactionsFilter.maxTime = 0;
            this.transactionsFilter.minTime = 0;
            this.transactionsFilter.type = 0;
            this.transactionsFilter.categoryId = '0';
            this.transactionsFilter.accountId = '0';
            this.transactionsFilter.keyword = '';
            this.transactions = [];
            this.transactionsNextTimeId = 0;
            this.transactionListStateInvalid = true;
        },
        clearTransactions() {
            this.transactions = [];
            this.transactionsNextTimeId = 0;
            this.transactionListStateInvalid = true;
        },
        initTransactionListFilter(filter) {
            if (filter && isNumber(filter.dateType)) {
                this.transactionsFilter.dateType = filter.dateType;
            } else {
                this.transactionsFilter.dateType = datetimeConstants.allDateRanges.All.type;
            }

            if (filter && isNumber(filter.maxTime)) {
                this.transactionsFilter.maxTime = filter.maxTime;
            } else {
                this.transactionsFilter.maxTime = 0;
            }

            if (filter && isNumber(filter.minTime)) {
                this.transactionsFilter.minTime = filter.minTime;
            } else {
                this.transactionsFilter.minTime = 0;
            }

            if (filter && isNumber(filter.type)) {
                this.transactionsFilter.type = filter.type;
            } else {
                this.transactionsFilter.type = 0;
            }

            if (filter && isString(filter.categoryId)) {
                this.transactionsFilter.categoryId = filter.categoryId;
            } else {
                this.transactionsFilter.categoryId = '0';
            }

            if (filter && isString(filter.accountId)) {
                this.transactionsFilter.accountId = filter.accountId;
            } else {
                this.transactionsFilter.accountId = '0';
            }

            if (filter && isString(filter.keyword)) {
                this.transactionsFilter.keyword = filter.keyword;
            } else {
                this.transactionsFilter.keyword = '';
            }
        },
        updateTransactionListFilter(filter) {
            if (filter && isNumber(filter.dateType)) {
                this.transactionsFilter.dateType = filter.dateType;
            }

            if (filter && isNumber(filter.maxTime)) {
                this.transactionsFilter.maxTime = filter.maxTime;
            }

            if (filter && isNumber(filter.minTime)) {
                this.transactionsFilter.minTime = filter.minTime;
            }

            if (filter && isNumber(filter.type)) {
                this.transactionsFilter.type = filter.type;
            }

            if (filter && isString(filter.categoryId)) {
                this.transactionsFilter.categoryId = filter.categoryId;
            }

            if (filter && isString(filter.accountId)) {
                this.transactionsFilter.accountId = filter.accountId;
            }

            if (filter && isString(filter.keyword)) {
                this.transactionsFilter.keyword = filter.keyword;
            }
        },
        getTransactionListPageParams() {
            const querys = [];

            if (this.transactionsFilter.type) {
                querys.push('type=' + this.transactionsFilter.type);
            }

            if (this.transactionsFilter.accountId && this.transactionsFilter.accountId !== '0') {
                querys.push('accountId=' + this.transactionsFilter.accountId);
            }

            if (this.transactionsFilter.categoryId && this.transactionsFilter.categoryId !== '0') {
                querys.push('categoryId=' + this.transactionsFilter.categoryId);
            }

            querys.push('dateType=' + this.transactionsFilter.dateType);

            if (this.transactionsFilter.dateType === datetimeConstants.allDateRanges.Custom.type) {
                querys.push('maxTime=' + this.transactionsFilter.maxTime);
                querys.push('minTime=' + this.transactionsFilter.minTime);
            }

            if (this.transactionsFilter.keyword) {
                querys.push('keyword=' + encodeURIComponent(this.transactionsFilter.keyword));
            }

            return querys.join('&');
        },
        loadTransactions({ reload, count, page, withCount, autoExpand, defaultCurrency }) {
            const self = this;
            const settingsStore = useSettingsStore();
            const exchangeRatesStore = useExchangeRatesStore();
            let actualMaxTime = self.transactionsNextTimeId;

            if (reload && self.transactionsFilter.maxTime > 0) {
                actualMaxTime = self.transactionsFilter.maxTime * 1000 + 999;
            } else if (reload && self.transactionsFilter.maxTime <= 0) {
                actualMaxTime = 0;
            }

            return new Promise((resolve, reject) => {
                services.getTransactions({
                    maxTime: actualMaxTime,
                    minTime: self.transactionsFilter.minTime * 1000,
                    count: count || 50,
                    page: page || 1,
                    withCount: (!!withCount) || false,
                    type: self.transactionsFilter.type,
                    categoryId: self.transactionsFilter.categoryId,
                    accountId: self.transactionsFilter.accountId,
                    keyword: self.transactionsFilter.keyword
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (reload) {
                            loadTransactionList(self, settingsStore, exchangeRatesStore, {
                                transactions: emptyTransactionResult,
                                reload: reload,
                                autoExpand: autoExpand,
                                defaultCurrency: defaultCurrency
                            });

                            if (!self.transactionListStateInvalid) {
                                self.updateTransactionListInvalidState(true);
                            }
                        }

                        reject({ message: 'Unable to retrieve transaction list' });
                        return;
                    }

                    loadTransactionList(self, settingsStore, exchangeRatesStore, {
                        transactions: data.result,
                        reload: reload,
                        autoExpand: autoExpand,
                        defaultCurrency: defaultCurrency
                    });

                    if (reload) {
                        if (self.transactionListStateInvalid) {
                            self.updateTransactionListInvalidState(false);
                        }
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to load transaction list', error);

                    if (reload) {
                        loadTransactionList(self, settingsStore, exchangeRatesStore, {
                            transactions: emptyTransactionResult,
                            reload: reload,
                            autoExpand: autoExpand,
                            defaultCurrency: defaultCurrency
                        });

                        if (!self.transactionListStateInvalid) {
                            self.updateTransactionListInvalidState(true);
                        }
                    }

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to retrieve transaction list' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        loadMonthlyAllTransactions({ year, month, autoExpand, defaultCurrency }) {
            const self = this;
            const settingsStore = useSettingsStore();
            const exchangeRatesStore = useExchangeRatesStore();

            return new Promise((resolve, reject) => {
                services.getAllTransactionsByMonth({
                    year: year,
                    month: month,
                    type: self.transactionsFilter.type,
                    categoryId: self.transactionsFilter.categoryId,
                    accountId: self.transactionsFilter.accountId,
                    keyword: self.transactionsFilter.keyword
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        loadTransactionList(self, settingsStore, exchangeRatesStore, {
                            transactions: emptyTransactionResult,
                            reload: true,
                            autoExpand: autoExpand,
                            defaultCurrency: defaultCurrency
                        });

                        if (!self.transactionListStateInvalid) {
                            self.updateTransactionListInvalidState(true);
                        }

                        reject({ message: 'Unable to retrieve transaction list' });
                        return;
                    }

                    loadTransactionList(self, settingsStore, exchangeRatesStore, {
                        transactions: data.result,
                        reload: true,
                        autoExpand: autoExpand,
                        defaultCurrency: defaultCurrency
                    });

                    if (self.transactionListStateInvalid) {
                        self.updateTransactionListInvalidState(false);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to load monthly all transaction list', error);

                    loadTransactionList(self, settingsStore, exchangeRatesStore, {
                        transactions: emptyTransactionResult,
                        reload: true,
                        autoExpand: autoExpand,
                        defaultCurrency: defaultCurrency
                    });

                    if (!self.transactionListStateInvalid) {
                        self.updateTransactionListInvalidState(true);
                    }

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to retrieve transaction list' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        getTransaction({ transactionId }) {
            return new Promise((resolve, reject) => {
                services.getTransaction({
                    id: transactionId
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to retrieve transaction' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to load transaction info', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to retrieve transaction' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        saveTransaction({ transaction, defaultCurrency, isEdit }) {
            const self = this;
            const settingsStore = useSettingsStore();
            const exchangeRatesStore = useExchangeRatesStore();

            const submitTransaction = {
                type: transaction.type,
                time: getActualUnixTimeForStore(transaction.time, transaction.utcOffset, getBrowserTimezoneOffsetMinutes()),
                sourceAccountId: transaction.sourceAccountId,
                sourceAmount: transaction.sourceAmount,
                destinationAccountId: '0',
                destinationAmount: 0,
                hideAmount: transaction.hideAmount,
                tagIds: transaction.tagIds,
                comment: transaction.comment,
                geoLocation: transaction.geoLocation,
                utcOffset: transaction.utcOffset
            };

            if (transaction.type === transactionConstants.allTransactionTypes.Expense) {
                submitTransaction.categoryId = transaction.expenseCategory;
            } else if (transaction.type === transactionConstants.allTransactionTypes.Income) {
                submitTransaction.categoryId = transaction.incomeCategory;
            } else if (transaction.type === transactionConstants.allTransactionTypes.Transfer) {
                submitTransaction.categoryId = transaction.transferCategory;
                submitTransaction.destinationAccountId = transaction.destinationAccountId;
                submitTransaction.destinationAmount = transaction.destinationAmount;
            } else {
                return Promise.reject('An error occurred');
            }

            if (isEdit) {
                submitTransaction.id = transaction.id;
            }

            return new Promise((resolve, reject) => {
                let promise = null;

                if (!submitTransaction.id) {
                    promise = services.addTransaction(submitTransaction);
                } else {
                    promise = services.modifyTransaction(submitTransaction);
                }

                promise.then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (!submitTransaction.id) {
                            reject({ message: 'Unable to add transaction' });
                        } else {
                            reject({ message: 'Unable to save transaction' });
                        }
                        return;
                    }

                    if (!submitTransaction.id) {
                        if (!self.transactionListStateInvalid) {
                            self.updateTransactionListInvalidState(true);
                        }
                    } else {
                        updateTransactionInTransactionList(self, settingsStore, exchangeRatesStore, {
                            transaction: data.result,
                            defaultCurrency: defaultCurrency
                        });
                    }

                    const accountsStore = useAccountsStore();
                    if (!accountsStore.accountListStateInvalid) {
                        accountsStore.updateAccountListInvalidState(true);
                    }

                    const overviewStore = useOverviewStore();
                    if (!overviewStore.transactionOverviewStateInvalid) {
                        overviewStore.updateTransactionOverviewInvalidState(true);
                    }

                    const statisticsStore = useStatisticsStore();
                    if (!statisticsStore.transactionStatisticsStateInvalid) {
                        statisticsStore.updateTransactionStatisticsInvalidState(true);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to save transaction', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        if (!submitTransaction.id) {
                            reject({ message: 'Unable to add transaction' });
                        } else {
                            reject({ message: 'Unable to save transaction' });
                        }
                    } else {
                        reject(error);
                    }
                });
            });
        },
        deleteTransaction({ transaction, defaultCurrency, beforeResolve }) {
            const self = this;
            const exchangeRatesStore = useExchangeRatesStore();

            return new Promise((resolve, reject) => {
                services.deleteTransaction({
                    id: transaction.id
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to delete this transaction' });
                        return;
                    }

                    if (beforeResolve) {
                        beforeResolve(() => {
                            removeTransactionFromTransactionList(self, exchangeRatesStore, {
                                transaction: transaction,
                                defaultCurrency: defaultCurrency
                            });
                        });
                    } else {
                        removeTransactionFromTransactionList(self, exchangeRatesStore, {
                            transaction: transaction,
                            defaultCurrency: defaultCurrency
                        });
                    }

                    const accountsStore = useAccountsStore();
                    if (!accountsStore.accountListStateInvalid) {
                        accountsStore.updateAccountListInvalidState(true);
                    }

                    const overviewStore = useOverviewStore();
                    if (!overviewStore.transactionOverviewStateInvalid) {
                        overviewStore.updateTransactionOverviewInvalidState(true);
                    }

                    const statisticsStore = useStatisticsStore();
                    if (!statisticsStore.transactionStatisticsStateInvalid) {
                        statisticsStore.updateTransactionStatisticsInvalidState(true);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to delete transaction', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to delete this transaction' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        collapseMonthInTransactionList({ month, collapse }) {
            if (month) {
                month.opened = !collapse;
            }
        }
    }
});
