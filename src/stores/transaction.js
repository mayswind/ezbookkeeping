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

            if (i === 0 && state.transactions.length > 0) {
                const lastMonthList = state.transactions[state.transactions.length - 1];

                if (lastMonthList.totalAmount.incompleteExpense || lastMonthList.totalAmount.incompleteIncome) {
                    // calculate the total amount of last month which has incomplete total amount before starting to process a new request
                    calculateMonthTotalAmount(exchangeRatesStore, lastMonthList, defaultCurrency, state.transactionsFilter.accountIds, false);
                }
            }

            if (currentMonthList && currentMonthList.year === transactionYear && currentMonthList.month === transactionMonth) {
                currentMonthList.items.push(Object.freeze(item));

                if (i === transactions.items.length - 1) {
                    // calculate the total amount of current month when processing the last transaction item of this request
                    calculateMonthTotalAmount(exchangeRatesStore, currentMonthList, defaultCurrency, state.transactionsFilter.accountIds, true);
                }
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
                // calculate the total amount of current month when processing the first transaction item of the next month
                calculateMonthTotalAmount(exchangeRatesStore, currentMonthList, defaultCurrency, state.transactionsFilter.accountIds, false);

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
            // init the total amount struct of current month when processing the first transaction item of current month
            calculateMonthTotalAmount(exchangeRatesStore, currentMonthList, defaultCurrency, state.transactionsFilter.accountIds, true);
        }
    }

    if (transactions.nextTimeSequenceId) {
        state.transactionsNextTimeId = transactions.nextTimeSequenceId;
    } else {
        calculateMonthTotalAmount(exchangeRatesStore, state.transactions[state.transactions.length - 1], defaultCurrency, state.transactionsFilter.accountIds, false);
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

                if ((state.transactionsFilter.categoryIds && !state.allFilterCategoryIds[transaction.categoryId]) ||
                    (state.transactionsFilter.accountIds && !state.allFilterAccountIds[transaction.sourceAccountId] && !state.allFilterAccountIds[transaction.destinationAccountId] &&
                        (!transaction.sourceAccount || !state.allFilterAccountIds[transaction.sourceAccount.parentId]) &&
                        (!transaction.destinationAccount || !state.allFilterAccountIds[transaction.destinationAccount.parentId])
                    )
                ) {
                    transactionMonthList.items.splice(j, 1);
                } else {
                    transactionMonthList.items.splice(j, 1, transaction);
                }

                if (transactionMonthList.items.length < 1) {
                    state.transactions.splice(i, 1);
                } else {
                    calculateMonthTotalAmount(exchangeRatesStore, transactionMonthList, defaultCurrency, state.transactionsFilter.accountIds, i >= state.transactions.length - 1 && state.transactionsNextTimeId > 0);
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
            calculateMonthTotalAmount(exchangeRatesStore, transactionMonthList, defaultCurrency, state.transactionsFilter.accountIds, i >= state.transactions.length - 1 && state.transactionsNextTimeId > 0);
        }
    }
}

function calculateMonthTotalAmount(exchangeRatesStore, transactionMonthList, defaultCurrency, accountIds, incomplete) {
    if (!transactionMonthList) {
        return;
    }

    let totalExpense = 0;
    let totalIncome = 0;
    let hasUnCalculatedTotalExpense = false;
    let hasUnCalculatedTotalIncome = false;

    const allAccountIdsMap = {};
    let totalAccountIdsCount = 0;

    if (accountIds && accountIds !== '0') {
        const allAccountIdsArray = accountIds.split(',');

        for (let i = 0; i < allAccountIdsArray.length; i++) {
            if (allAccountIdsArray[i]) {
                allAccountIdsMap[allAccountIdsArray[i]] = true;
                totalAccountIdsCount++;
            }
        }
    }

    for (let i = 0; i < transactionMonthList.items.length; i++) {
        const transaction = transactionMonthList.items[i];

        let amount = transaction.sourceAmount;
        let account = transaction.sourceAccount;

        if (totalAccountIdsCount > 0 && transaction.destinationAccount
            && (!allAccountIdsMap[transaction.sourceAccount.id] && !allAccountIdsMap[transaction.sourceAccount.parentId])
            && (allAccountIdsMap[transaction.destinationAccount.id] || allAccountIdsMap[transaction.destinationAccount.parentId])) {
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
        } else if (transaction.type === transactionConstants.allTransactionTypes.Transfer && totalAccountIdsCount > 0) {
            if (allAccountIdsMap[transaction.sourceAccountId] && allAccountIdsMap[transaction.destinationAccountId]) {
                // Do Nothing
            } else if (transaction.sourceAccount && transaction.destinationAccount && allAccountIdsMap[transaction.sourceAccount.parentId] && allAccountIdsMap[transaction.destinationAccount.parentId]) {
                // Do Nothing
            } else if (transaction.sourceAccount && allAccountIdsMap[transaction.sourceAccount.parentId] && allAccountIdsMap[transaction.destinationAccountId]) {
                // Do Nothing
            } else if (transaction.destinationAccount && allAccountIdsMap[transaction.sourceAccountId] && allAccountIdsMap[transaction.destinationAccount.parentId]) {
                // Do Nothing
            } else if (allAccountIdsMap[transaction.sourceAccountId] || (transaction.sourceAccount && allAccountIdsMap[transaction.sourceAccount.parentId])) {
                totalExpense += amount;
            } else if (allAccountIdsMap[transaction.destinationAccountId] || (transaction.destinationAccount && allAccountIdsMap[transaction.destinationAccount.parentId])) {
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

function buildBasicSubmitTransaction(transaction, dummyTime) {
    const submitTransaction = {
        type: transaction.type,
        time: dummyTime ? getActualUnixTimeForStore(transaction.time, transaction.utcOffset, getBrowserTimezoneOffsetMinutes()) : transaction.time,
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

    if (transaction.type === transactionConstants.allTransactionTypes.Transfer) {
        submitTransaction.destinationAccountId = transaction.destinationAccountId;
        submitTransaction.destinationAmount = transaction.destinationAmount;
    }

    return submitTransaction;
}

export const useTransactionsStore = defineStore('transactions', {
    state: () => ({
        transactionsFilter: {
            dateType: datetimeConstants.allDateRanges.All.type,
            maxTime: 0,
            minTime: 0,
            type: 0,
            categoryIds: '',
            accountIds: '',
            tagIds: '',
            amountFilter: '',
            keyword: ''
        },
        transactions: [],
        transactionsNextTimeId: 0,
        transactionListStateInvalid: true,
    }),
    getters: {
        allFilterCategoryIds(state) {
            if (!state.transactionsFilter.categoryIds) {
                return {};
            }

            const allCategoryIds = state.transactionsFilter.categoryIds.split(',');
            const ret = {};

            for (let i = 0; i < allCategoryIds.length; i++) {
                if (allCategoryIds[i]) {
                    ret[allCategoryIds[i]] = true;
                }
            }

            return ret;
        },
        allFilterAccountIds(state) {
            if (!state.transactionsFilter.accountIds) {
                return {};
            }

            const allAccountIds = state.transactionsFilter.accountIds.split(',');
            const ret = {};

            for (let i = 0; i < allAccountIds.length; i++) {
                if (allAccountIds[i]) {
                    ret[allAccountIds[i]] = true;
                }
            }

            return ret;
        },
        allFilterTagIds(state) {
            if (!state.transactionsFilter.tagIds) {
                return {};
            }

            const allTagIds = state.transactionsFilter.tagIds.split(',');
            const ret = {};

            for (let i = 0; i < allTagIds.length; i++) {
                if (allTagIds[i]) {
                    ret[allTagIds[i]] = true;
                }
            }

            return ret;
        },
        allFilterCategoryIdsCount(state) {
            if (!state.transactionsFilter.categoryIds) {
                return 0;
            }

            const allCategoryIds = state.transactionsFilter.categoryIds.split(',');
            let count = 0;

            for (let i = 0; i < allCategoryIds.length; i++) {
                if (allCategoryIds[i]) {
                    count++;
                }
            }

            return count;
        },
        allFilterAccountIdsCount(state) {
            if (!state.transactionsFilter.accountIds) {
                return 0;
            }

            const allAccountIds = state.transactionsFilter.accountIds.split(',');
            let count = 0;

            for (let i = 0; i < allAccountIds.length; i++) {
                if (allAccountIds[i]) {
                    count++;
                }
            }

            return count;
        },
        allFilterTagIdsCount(state) {
            if (!state.transactionsFilter.tagIds) {
                return 0;
            }

            const allTagIds = state.transactionsFilter.tagIds.split(',');
            let count = 0;

            for (let i = 0; i < allTagIds.length; i++) {
                if (allTagIds[i]) {
                    count++;
                }
            }

            return count;
        },
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
                pictures: [],
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
                    (transactionConstants.minAmountNumber <= newValue && newValue <= transactionConstants.maxAmountNumber)) {
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
            this.transactionsFilter.categoryIds = '';
            this.transactionsFilter.accountIds = '';
            this.transactionsFilter.tagIds = '';
            this.transactionsFilter.amountFilter = '';
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

            if (filter && isString(filter.categoryIds)) {
                this.transactionsFilter.categoryIds = filter.categoryIds;
            } else {
                this.transactionsFilter.categoryIds = '';
            }

            if (filter && isString(filter.accountIds)) {
                this.transactionsFilter.accountIds = filter.accountIds;
            } else {
                this.transactionsFilter.accountIds = '';
            }

            if (filter && isString(filter.tagIds)) {
                this.transactionsFilter.tagIds = filter.tagIds;
            } else {
                this.transactionsFilter.tagIds = '';
            }

            if (filter && isString(filter.amountFilter)) {
                this.transactionsFilter.amountFilter = filter.amountFilter;
            } else {
                this.transactionsFilter.amountFilter = '';
            }

            if (filter && isString(filter.keyword)) {
                this.transactionsFilter.keyword = filter.keyword;
            } else {
                this.transactionsFilter.keyword = '';
            }
        },
        updateTransactionListFilter(filter) {
            let changed = false;

            if (filter && isNumber(filter.dateType) && this.transactionsFilter.dateType !== filter.dateType) {
                this.transactionsFilter.dateType = filter.dateType;
                changed = true;
            }

            if (filter && isNumber(filter.maxTime) && this.transactionsFilter.maxTime !== filter.maxTime) {
                this.transactionsFilter.maxTime = filter.maxTime;
                changed = true;
            }

            if (filter && isNumber(filter.minTime) && this.transactionsFilter.minTime !== filter.minTime) {
                this.transactionsFilter.minTime = filter.minTime;
                changed = true;
            }

            if (filter && isNumber(filter.type) && this.transactionsFilter.type !== filter.type) {
                this.transactionsFilter.type = filter.type;
                changed = true;
            }

            if (filter && isString(filter.categoryIds) && this.transactionsFilter.categoryIds !== filter.categoryIds) {
                this.transactionsFilter.categoryIds = filter.categoryIds;
                changed = true;
            }

            if (filter && isString(filter.accountIds) && this.transactionsFilter.accountIds !== filter.accountIds) {
                this.transactionsFilter.accountIds = filter.accountIds;
                changed = true;
            }

            if (filter && isString(filter.tagIds) && this.transactionsFilter.tagIds !== filter.tagIds) {
                this.transactionsFilter.tagIds = filter.tagIds;
                changed = true;
            }

            if (filter && isString(filter.amountFilter) && this.transactionsFilter.amountFilter !== filter.amountFilter) {
                this.transactionsFilter.amountFilter = filter.amountFilter;
                changed = true;
            }

            if (filter && isString(filter.keyword) && this.transactionsFilter.keyword !== filter.keyword) {
                this.transactionsFilter.keyword = filter.keyword;
                changed = true;
            }

            return changed;
        },
        getTransactionListPageParams() {
            const querys = [];

            if (this.transactionsFilter.type) {
                querys.push('type=' + this.transactionsFilter.type);
            }

            if (this.transactionsFilter.accountIds) {
                querys.push('accountIds=' + this.transactionsFilter.accountIds);
            }

            if (this.transactionsFilter.categoryIds) {
                querys.push('categoryIds=' + this.transactionsFilter.categoryIds);
            }

            if (this.transactionsFilter.tagIds) {
                querys.push('tagIds=' + this.transactionsFilter.tagIds);
            }

            querys.push('dateType=' + this.transactionsFilter.dateType);

            if (this.transactionsFilter.dateType === datetimeConstants.allDateRanges.Custom.type) {
                querys.push('maxTime=' + this.transactionsFilter.maxTime);
                querys.push('minTime=' + this.transactionsFilter.minTime);
            }

            if (this.transactionsFilter.amountFilter) {
                querys.push('amountFilter=' + encodeURIComponent(this.transactionsFilter.amountFilter));
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
                    categoryIds: self.transactionsFilter.categoryIds,
                    accountIds: self.transactionsFilter.accountIds,
                    tagIds: self.transactionsFilter.tagIds,
                    amountFilter: self.transactionsFilter.amountFilter,
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
                    categoryIds: self.transactionsFilter.categoryIds,
                    accountIds: self.transactionsFilter.accountIds,
                    tagIds: self.transactionsFilter.tagIds,
                    amountFilter: self.transactionsFilter.amountFilter,
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
        saveTransaction({ transaction, defaultCurrency, isEdit, clientSessionId }) {
            const self = this;
            const settingsStore = useSettingsStore();
            const exchangeRatesStore = useExchangeRatesStore();

            const submitTransaction = buildBasicSubmitTransaction(transaction, true);

            if (transaction.type === transactionConstants.allTransactionTypes.Expense) {
                submitTransaction.categoryId = transaction.expenseCategory;
            } else if (transaction.type === transactionConstants.allTransactionTypes.Income) {
                submitTransaction.categoryId = transaction.incomeCategory;
            } else if (transaction.type === transactionConstants.allTransactionTypes.Transfer) {
                submitTransaction.categoryId = transaction.transferCategory;
            } else {
                return Promise.reject('An error occurred');
            }

            if (clientSessionId) {
                submitTransaction.clientSessionId = clientSessionId;
            }

            if (transaction.pictures && transaction.pictures.length > 0) {
                const pictureIds = [];

                for (let i = 0; i < transaction.pictures.length; i++) {
                    if (transaction.pictures[i].pictureId) {
                        pictureIds.push(transaction.pictures[i].pictureId);
                    }
                }

                submitTransaction.pictureIds = pictureIds;
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
        parseImportTransaction({ fileType, importFile }) {
            return new Promise((resolve, reject) => {
                services.parseImportTransaction({ fileType, importFile }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to parse import file' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('Unable to parse import file', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to parse import file' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        importTransactions({ transactions, clientSessionId }) {
            const submitTransactions = [];

            if (transactions) {
                for (let i = 0; i < transactions.length; i++) {
                    const transaction = transactions[i];
                    const submitTransaction = buildBasicSubmitTransaction(transaction, false);

                    submitTransaction.categoryId = transaction.categoryId;
                    submitTransactions.push(submitTransaction);
                }
            }

            return new Promise((resolve, reject) => {
                services.importTransactions({
                    transactions: submitTransactions,
                    clientSessionId: clientSessionId
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to import transactions' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('Unable to import transactions', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to import transactions' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        uploadTransactionPicture({ pictureFile, clientSessionId }) {
            return new Promise((resolve, reject) => {
                services.uploadTransactionPicture({ pictureFile, clientSessionId }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to upload transaction picture' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('Unable to upload transaction picture', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to upload transaction picture' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        removeUnusedTransactionPicture({ pictureInfo }) {
            return new Promise((resolve, reject) => {
                services.removeUnusedTransactionPicture({
                    id: pictureInfo.pictureId
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to remove transaction picture' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to remove transaction picture', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to remove transaction picture' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        getTransactionPictureUrl(pictureInfo, disableBrowserCache) {
            if (!pictureInfo || !pictureInfo.originalUrl) {
                return null;
            }

            return services.getTransactionPictureUrlWithToken(pictureInfo.originalUrl, disableBrowserCache);
        },
        collapseMonthInTransactionList({ month, collapse }) {
            if (month) {
                month.opened = !collapse;
            }
        }
    }
});
