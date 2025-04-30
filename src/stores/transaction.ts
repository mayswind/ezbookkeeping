import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useOverviewStore } from './overview.ts';
import { useStatisticsStore } from './statistics.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import type { BeforeResolveFunction } from '@/core/base.ts';
import { DateRange } from '@/core/datetime.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType, TransactionTagFilterType } from '@/core/transaction.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import {
    type TransactionDraft,
    type TransactionCreateRequest,
    type TransactionInfoResponse,
    type TransactionPageWrapper,
    Transaction,
    EMPTY_TRANSACTION_RESULT
} from '@/models/transaction.ts';
import type {
    TransactionPictureInfoBasicResponse
} from '@/models/transaction_picture_info.ts';
import {
    type ImportTransactionResponsePageWrapper,
    ImportTransaction
} from '@/models/imported_transaction.ts';

import {
    getUserTransactionDraft,
    updateUserTransactionDraft,
    clearUserTransactionDraft
} from '@/lib/userstate.ts';
import {
    isDefined,
    isNumber,
    isString,
    isArray1SubsetOfArray2,
    splitItemsToMap,
    countSplitItems
} from '@/lib/common.ts';
import {
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
} from '@/lib/datetime.ts';
import { getAmountWithDecimalNumberCount } from '@/lib/numeral.ts';
import { getCurrencyFraction } from '@/lib/currency.ts';
import { getFirstAvailableCategoryId } from '@/lib/category.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export interface TransactionListPartialFilter {
    dateType?: number;
    maxTime?: number;
    minTime?: number;
    type?: number;
    categoryIds?: string;
    accountIds?: string;
    tagIds?: string;
    tagFilterType?: number;
    amountFilter?: string;
    keyword?: string;
}

export interface TransactionListFilter extends TransactionListPartialFilter {
    dateType: number;
    maxTime: number;
    minTime: number;
    type: number;
    categoryIds: string;
    accountIds: string;
    tagIds: string;
    tagFilterType: number;
    amountFilter: string;
    keyword: string;
}

export interface TransactionMonthList {
    readonly year: number;
    readonly month: number;
    readonly yearMonth: string;
    opened: boolean;
    readonly items: Transaction[];
    readonly totalAmount: {
        expense: number;
        incompleteExpense: boolean;
        income: number;
        incompleteIncome: boolean;
    };
}

export const useTransactionsStore = defineStore('transactions', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const overviewStore = useOverviewStore();
    const statisticsStore = useStatisticsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const transactionDraft = ref<TransactionDraft | null>(getUserTransactionDraft());

    const transactionsFilter = ref<TransactionListFilter>({
        dateType: DateRange.All.type,
        maxTime: 0,
        minTime: 0,
        type: 0,
        categoryIds: '',
        accountIds: '',
        tagIds: '',
        tagFilterType: TransactionTagFilterType.Default.type,
        amountFilter: '',
        keyword: ''
    });

    const transactions = ref<TransactionMonthList[]>([]);
    const transactionsNextTimeId = ref<number>(0);
    const transactionListStateInvalid = ref<boolean>(true);

    const allFilterCategoryIds = computed<Record<string, boolean>>(() => splitItemsToMap(transactionsFilter.value.categoryIds, ','));
    const allFilterAccountIds = computed<Record<string, boolean>>(() => splitItemsToMap(transactionsFilter.value.accountIds, ','));
    const allFilterTagIds = computed<Record<string, boolean>>(() => splitItemsToMap(transactionsFilter.value.tagIds, ','));

    const allFilterCategoryIdsCount = computed<number>(() => countSplitItems(transactionsFilter.value.categoryIds, ','));
    const allFilterAccountIdsCount = computed<number>(() => countSplitItems(transactionsFilter.value.accountIds, ','));
    const allFilterTagIdsCount = computed<number>(() => countSplitItems(transactionsFilter.value.tagIds, ','));

    const noTransaction = computed<boolean>(() => {
        for (let i = 0; i < transactions.value.length; i++) {
            const transactionMonthList = transactions.value[i];

            for (let j = 0; j < transactionMonthList.items.length; j++) {
                if (transactionMonthList.items[j]) {
                    return false;
                }
            }
        }

        return true;
    });

    const hasMoreTransaction = computed<boolean>(() => {
        return transactionsNextTimeId.value > 0;
    });

    function loadTransactionList({ transactionPageWrapper, reload, autoExpand, defaultCurrency, nextTimeSequenceId }: { transactionPageWrapper: TransactionPageWrapper, reload: boolean, autoExpand: boolean, defaultCurrency: string, nextTimeSequenceId?: number }): void {
        if (reload) {
            transactions.value = [];
        }

        if (transactionPageWrapper.items && transactionPageWrapper.items.length) {
            const currentUtcOffset = getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone);
            let currentMonthListIndex = -1;
            let currentMonthList: TransactionMonthList | null = null;

            for (let i = 0; i < transactionPageWrapper.items.length; i++) {
                const item = transactionPageWrapper.items[i];
                fillTransactionObject(item, currentUtcOffset);

                const transactionTime = parseDateFromUnixTime(item.time, item.utcOffset, currentUtcOffset);
                const transactionYear = getYear(transactionTime);
                const transactionMonth = getMonth(transactionTime);
                const transactionYearMonth = getYearAndMonth(transactionTime);

                if (i === 0 && transactions.value.length > 0) {
                    const lastMonthList = transactions.value[transactions.value.length - 1];

                    if (lastMonthList.totalAmount.incompleteExpense || lastMonthList.totalAmount.incompleteIncome) {
                        // calculate the total amount of last month which has incomplete total amount before starting to process a new request
                        calculateMonthTotalAmount(lastMonthList, defaultCurrency, transactionsFilter.value.accountIds, false);
                    }
                }

                if (currentMonthList && currentMonthList.year === transactionYear && currentMonthList.month === transactionMonth) {
                    currentMonthList.items.push(Object.freeze(item));

                    if (i === transactionPageWrapper.items.length - 1) {
                        // calculate the total amount of current month when processing the last transaction item of this request
                        calculateMonthTotalAmount(currentMonthList, defaultCurrency, transactionsFilter.value.accountIds, true);
                    }
                    continue;
                }

                for (let j = currentMonthListIndex + 1; j < transactions.value.length; j++) {
                    if (transactions.value[j].year === transactionYear && transactions.value[j].month === transactionMonth) {
                        currentMonthListIndex = j;
                        currentMonthList = transactions.value[j];
                        break;
                    }
                }

                if (!currentMonthList || currentMonthList.year !== transactionYear || currentMonthList.month !== transactionMonth) {
                    // calculate the total amount of current month when processing the first transaction item of the next month
                    calculateMonthTotalAmount(currentMonthList, defaultCurrency, transactionsFilter.value.accountIds, false);

                    const monthList: TransactionMonthList = {
                        year: transactionYear,
                        month: transactionMonth,
                        yearMonth: transactionYearMonth,
                        opened: autoExpand,
                        items: [],
                        totalAmount: {
                            expense: 0,
                            incompleteExpense: true,
                            income: 0,
                            incompleteIncome: true
                        }
                    };

                    transactions.value.push(monthList);

                    currentMonthListIndex = transactions.value.length - 1;
                    currentMonthList = transactions.value[transactions.value.length - 1];
                }

                currentMonthList.items.push(Object.freeze(item));
                // init the total amount struct of current month when processing the first transaction item of current month
                calculateMonthTotalAmount(currentMonthList, defaultCurrency, transactionsFilter.value.accountIds, true);
            }
        }

        if (nextTimeSequenceId) {
            transactionsNextTimeId.value = nextTimeSequenceId;
        } else {
            calculateMonthTotalAmount(transactions.value[transactions.value.length - 1], defaultCurrency, transactionsFilter.value.accountIds, false);
            transactionsNextTimeId.value = -1;
        }
    }

    function updateTransactionInTransactionList({ transaction, defaultCurrency }: { transaction: Transaction, defaultCurrency: string }): void {
        const currentUtcOffset = getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone);
        const transactionTime = parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentUtcOffset);
        const transactionYear = getYear(transactionTime);
        const transactionMonth = getMonth(transactionTime);

        for (let i = 0; i < transactions.value.length; i++) {
            const transactionMonthList = transactions.value[i];

            if (!transactionMonthList.items) {
                continue;
            }

            for (let j = 0; j < transactionMonthList.items.length; j++) {
                if (transactionMonthList.items[j].id === transaction.id) {
                    fillTransactionObject(transaction, currentUtcOffset);

                    if (transactionYear !== transactionMonthList.year ||
                        transactionMonth !== transactionMonthList.month ||
                        transaction.day !== transactionMonthList.items[j].day) {
                        transactionListStateInvalid.value = true;
                        return;
                    }

                    if ((transactionsFilter.value.categoryIds && !allFilterCategoryIds.value[transaction.categoryId]) ||
                        (transactionsFilter.value.accountIds && !allFilterAccountIds.value[transaction.sourceAccountId] && !allFilterAccountIds.value[transaction.destinationAccountId] &&
                            (!transaction.sourceAccount || !allFilterAccountIds.value[transaction.sourceAccount.parentId]) &&
                            (!transaction.destinationAccount || !allFilterAccountIds.value[transaction.destinationAccount.parentId])
                        )
                    ) {
                        transactionMonthList.items.splice(j, 1);
                    } else {
                        transactionMonthList.items.splice(j, 1, transaction);
                    }

                    if (transactionMonthList.items.length < 1) {
                        transactions.value.splice(i, 1);
                    } else {
                        calculateMonthTotalAmount(transactionMonthList, defaultCurrency, transactionsFilter.value.accountIds, i >= transactions.value.length - 1 && transactionsNextTimeId.value > 0);
                    }

                    return;
                }
            }
        }
    }

    function removeTransactionFromTransactionList({ transaction, defaultCurrency }: { transaction: TransactionInfoResponse, defaultCurrency: string }): void {
        for (let i = 0; i <transactions.value.length; i++) {
            const transactionMonthList = transactions.value[i];

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
                transactions.value.splice(i, 1);
            } else {
                calculateMonthTotalAmount(transactionMonthList, defaultCurrency, transactionsFilter.value.accountIds, i >= transactions.value.length - 1 && transactionsNextTimeId.value > 0);
            }
        }
    }

    function calculateMonthTotalAmount(transactionMonthList: TransactionMonthList | null, defaultCurrency: string, accountIds: string, incomplete: boolean): void {
        if (!transactionMonthList) {
            return;
        }

        let totalExpense = 0;
        let totalIncome = 0;
        let hasUnCalculatedTotalExpense = false;
        let hasUnCalculatedTotalIncome = false;

        const allAccountIdsMap: Record<string, boolean> = {};
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
                && (!allAccountIdsMap[transaction.sourceAccount?.id || ''] && !allAccountIdsMap[transaction.sourceAccount?.parentId || ''])
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
                    if (transaction.type === TransactionType.Expense) {
                        hasUnCalculatedTotalExpense = true;
                    } else if (transaction.type === TransactionType.Income) {
                        hasUnCalculatedTotalIncome = true;
                    }

                    continue;
                }

                amount = balance;
            }

            if (transaction.type === TransactionType.Expense) {
                totalExpense += amount;
            } else if (transaction.type === TransactionType.Income) {
                totalIncome += amount;
            } else if (transaction.type === TransactionType.Transfer && totalAccountIdsCount > 0) {
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

        transactionMonthList.totalAmount.expense = Math.floor(totalExpense);
        transactionMonthList.totalAmount.incompleteExpense = incomplete || hasUnCalculatedTotalExpense;
        transactionMonthList.totalAmount.income = Math.floor(totalIncome);
        transactionMonthList.totalAmount.incompleteIncome = incomplete || hasUnCalculatedTotalIncome;
    }

    function fillTransactionObject(transaction: Transaction, currentUtcOffset: number): void {
        if (!transaction) {
            return;
        }

        const transactionTime = parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentUtcOffset);
        transaction.setDisplayDate(getShortDate(transactionTime), getDay(transactionTime), getDayOfWeekName(transactionTime));

        if (transaction.sourceAccountId) {
            transaction.setSourceAccount(accountsStore.allAccountsMap[transaction.sourceAccountId]);
        }

        if (transaction.destinationAccountId) {
            transaction.setDestinationAccount(accountsStore.allAccountsMap[transaction.destinationAccountId]);
        }

        if (transaction.categoryId) {
            transaction.setCategory(transactionCategoriesStore.allTransactionCategoriesMap[transaction.categoryId]);
        }
    }

    function initTransactionDraft(): void {
        if (settingsStore.appSettings.autoSaveTransactionDraft === 'enabled' || settingsStore.appSettings.autoSaveTransactionDraft === 'confirmation') {
            transactionDraft.value = getUserTransactionDraft();
        } else {
            transactionDraft.value = null;
        }
    }

    function isTransactionDraftModified(transaction?: Transaction, initCategoryId?: string, initAccountId?: string, initTagIds?: string): boolean {
        if (!transaction) {
            return false;
        }

        if (transaction.sourceAmount !== 0) {
            return true;
        }

        if (transaction.type === TransactionType.Transfer && transaction.destinationAmount !== 0) {
            return true;
        }

        if (transaction.sourceAccountId && transaction.sourceAccountId !== '0' && transaction.sourceAccountId !== userStore.currentUserDefaultAccountId && transaction.sourceAccountId !== initAccountId) {
            return true;
        }

        if (transaction.type === TransactionType.Transfer && transaction.destinationAccountId && transaction.destinationAccountId !== '0' && transaction.destinationAccountId !== userStore.currentUserDefaultAccountId && transaction.destinationAccountId !== initAccountId) {
            return true;
        }

        const allCategories = transactionCategoriesStore.allTransactionCategories;

        if (allCategories) {
            if (transaction.type === TransactionType.Expense) {
                const defaultCategoryId = getFirstAvailableCategoryId(allCategories[CategoryType.Expense]);

                if (transaction.expenseCategoryId && transaction.expenseCategoryId !== '0' && transaction.expenseCategoryId !== defaultCategoryId && transaction.expenseCategoryId !== initCategoryId) {
                    return true;
                }
            } else if (transaction.type === TransactionType.Income) {
                const defaultCategoryId = getFirstAvailableCategoryId(allCategories[CategoryType.Income]);

                if (transaction.incomeCategoryId && transaction.incomeCategoryId !== '0' && transaction.incomeCategoryId !== defaultCategoryId && transaction.incomeCategoryId !== initCategoryId) {
                    return true;
                }
            } else if (transaction.type === TransactionType.Transfer) {
                const defaultCategoryId = getFirstAvailableCategoryId(allCategories[CategoryType.Transfer]);

                if (transaction.transferCategoryId && transaction.transferCategoryId !== '0' && transaction.transferCategoryId !== defaultCategoryId && transaction.transferCategoryId !== initCategoryId) {
                    return true;
                }
            }
        }

        if (transaction.hideAmount) {
            return true;
        }

        if (transaction.tagIds && transaction.tagIds.length > 0) {
            return !initTagIds || !isArray1SubsetOfArray2(transaction.tagIds, initTagIds.split(','));
        }

        if (transaction.pictures && transaction.pictures.length > 0) {
            return true;
        }

        if (transaction.comment && transaction.comment.trim()) {
            return true;
        }

        return false;
    }

    function saveTransactionDraft(transaction?: Transaction, initCategoryId?: string, initAccountId?: string, initTagIds?: string): void {
        if (settingsStore.appSettings.autoSaveTransactionDraft !== 'enabled' && settingsStore.appSettings.autoSaveTransactionDraft !== 'confirmation') {
            clearTransactionDraft();
            return;
        }

        if (transaction) {
            if (!isTransactionDraftModified(transaction, initCategoryId, initAccountId, initTagIds)) {
                clearTransactionDraft();
                return;
            }

            transactionDraft.value = transaction.toTransactionDraft();
        }

        updateUserTransactionDraft(transactionDraft.value);
    }

    function clearTransactionDraft(): void {
        transactionDraft.value = null;
        clearUserTransactionDraft();
    }

    function setTransactionSuitableDestinationAmount(transaction: Transaction, oldValue: number, newValue: number): void {
        if (transaction.type === TransactionType.Expense || transaction.type === TransactionType.Income) {
            transaction.destinationAmount = newValue;
        } else if (transaction.type === TransactionType.Transfer) {
            const sourceAccount = accountsStore.allAccountsMap[transaction.sourceAccountId];
            const destinationAccount = accountsStore.allAccountsMap[transaction.destinationAccountId];

            if (sourceAccount && destinationAccount && sourceAccount.currency !== destinationAccount.currency) {
                const decimalNumberCount = getCurrencyFraction(destinationAccount.currency);
                const exchangedOldValue = exchangeRatesStore.getExchangedAmount(oldValue, sourceAccount.currency, destinationAccount.currency);
                const exchangedNewValue = exchangeRatesStore.getExchangedAmount(newValue, sourceAccount.currency, destinationAccount.currency);

                if (isNumber(decimalNumberCount) && isNumber(exchangedOldValue)) {
                    oldValue = Math.floor(exchangedOldValue);
                    oldValue = getAmountWithDecimalNumberCount(oldValue, decimalNumberCount);
                }

                if (isNumber(decimalNumberCount) && isNumber(exchangedNewValue)) {
                    newValue = Math.floor(exchangedNewValue);
                    newValue = getAmountWithDecimalNumberCount(newValue, decimalNumberCount);
                } else {
                    return;
                }
            }

            if ((!sourceAccount || !destinationAccount || transaction.destinationAmount === oldValue || transaction.destinationAmount === 0) &&
                (TRANSACTION_MIN_AMOUNT <= newValue && newValue <= TRANSACTION_MAX_AMOUNT)) {
                transaction.destinationAmount = newValue;
            }
        }
    }

    function updateTransactionListInvalidState(invalidState: boolean): void {
        transactionListStateInvalid.value = invalidState;
    }

    function resetTransactions(): void {
        transactionsFilter.value.dateType = DateRange.All.type;
        transactionsFilter.value.maxTime = 0;
        transactionsFilter.value.minTime = 0;
        transactionsFilter.value.type = 0;
        transactionsFilter.value.categoryIds = '';
        transactionsFilter.value.accountIds = '';
        transactionsFilter.value.tagIds = '';
        transactionsFilter.value.tagFilterType = TransactionTagFilterType.Default.type;
        transactionsFilter.value.amountFilter = '';
        transactionsFilter.value.keyword = '';
        transactions.value = [];
        transactionsNextTimeId.value = 0;
        transactionListStateInvalid.value = true;
    }

    function clearTransactions(): void {
        transactions.value = [];
        transactionsNextTimeId.value = 0;
        transactionListStateInvalid.value = true;
    }

    function initTransactionListFilter(filter: TransactionListPartialFilter): void {
        if (filter && isNumber(filter.dateType)) {
            transactionsFilter.value.dateType = filter.dateType;
        } else {
            transactionsFilter.value.dateType = DateRange.All.type;
        }

        if (filter && isNumber(filter.maxTime)) {
            transactionsFilter.value.maxTime = filter.maxTime;
        } else {
            transactionsFilter.value.maxTime = 0;
        }

        if (filter && isNumber(filter.minTime)) {
            transactionsFilter.value.minTime = filter.minTime;
        } else {
            transactionsFilter.value.minTime = 0;
        }

        if (filter && isNumber(filter.type)) {
            transactionsFilter.value.type = filter.type;
        } else {
            transactionsFilter.value.type = 0;
        }

        if (filter && isString(filter.categoryIds)) {
            transactionsFilter.value.categoryIds = filter.categoryIds;
        } else {
            transactionsFilter.value.categoryIds = '';
        }

        if (filter && isString(filter.accountIds)) {
            transactionsFilter.value.accountIds = filter.accountIds;
        } else {
            transactionsFilter.value.accountIds = '';
        }

        if (filter && isString(filter.tagIds)) {
            transactionsFilter.value.tagIds = filter.tagIds;
        } else {
            transactionsFilter.value.tagIds = '';
        }

        if (filter && isNumber(filter.tagFilterType)) {
            transactionsFilter.value.tagFilterType = filter.tagFilterType;
        } else {
            transactionsFilter.value.tagFilterType = TransactionTagFilterType.Default.type;
        }

        if (filter && isString(filter.amountFilter)) {
            transactionsFilter.value.amountFilter = filter.amountFilter;
        } else {
            transactionsFilter.value.amountFilter = '';
        }

        if (filter && isString(filter.keyword)) {
            transactionsFilter.value.keyword = filter.keyword;
        } else {
            transactionsFilter.value.keyword = '';
        }
    }

    function updateTransactionListFilter(filter: TransactionListPartialFilter): boolean {
        let changed = false;

        if (filter && isNumber(filter.dateType) && transactionsFilter.value.dateType !== filter.dateType) {
            transactionsFilter.value.dateType = filter.dateType;
            changed = true;
        }

        if (filter && isNumber(filter.maxTime) && transactionsFilter.value.maxTime !== filter.maxTime) {
            transactionsFilter.value.maxTime = filter.maxTime;
            changed = true;
        }

        if (filter && isNumber(filter.minTime) && transactionsFilter.value.minTime !== filter.minTime) {
            transactionsFilter.value.minTime = filter.minTime;
            changed = true;
        }

        if (filter && isNumber(filter.type) && transactionsFilter.value.type !== filter.type) {
            transactionsFilter.value.type = filter.type;
            changed = true;
        }

        if (filter && isString(filter.categoryIds) && transactionsFilter.value.categoryIds !== filter.categoryIds) {
            transactionsFilter.value.categoryIds = filter.categoryIds;
            changed = true;
        }

        if (filter && isString(filter.accountIds) && transactionsFilter.value.accountIds !== filter.accountIds) {
            if (DateRange.isBillingCycle(transactionsFilter.value.dateType) &&
                (!accountsStore.getAccountStatementDate(filter.accountIds) || accountsStore.getAccountStatementDate(filter.accountIds) !== accountsStore.getAccountStatementDate(transactionsFilter.value.accountIds))) {
                transactionsFilter.value.dateType = DateRange.Custom.type;
            }

            transactionsFilter.value.accountIds = filter.accountIds;
            changed = true;
        }

        if (filter && isString(filter.tagIds) && transactionsFilter.value.tagIds !== filter.tagIds) {
            transactionsFilter.value.tagIds = filter.tagIds;
            changed = true;
        }

        if (filter && isNumber(filter.tagFilterType) && transactionsFilter.value.tagFilterType !== filter.tagFilterType) {
            transactionsFilter.value.tagFilterType = filter.tagFilterType;
            changed = true;
        }

        if (filter && isString(filter.amountFilter) && transactionsFilter.value.amountFilter !== filter.amountFilter) {
            transactionsFilter.value.amountFilter = filter.amountFilter;
            changed = true;
        }

        if (filter && isString(filter.keyword) && transactionsFilter.value.keyword !== filter.keyword) {
            transactionsFilter.value.keyword = filter.keyword;
            changed = true;
        }

        return changed;
    }

    function getTransactionListPageParams(): string {
        const querys: string[] = [];

        if (transactionsFilter.value.type) {
            querys.push('type=' + transactionsFilter.value.type);
        }

        if (transactionsFilter.value.accountIds) {
            querys.push('accountIds=' + transactionsFilter.value.accountIds);
        }

        if (transactionsFilter.value.categoryIds) {
            querys.push('categoryIds=' + transactionsFilter.value.categoryIds);
        }

        if (transactionsFilter.value.tagIds) {
            querys.push('tagIds=' + transactionsFilter.value.tagIds);
        }

        if (transactionsFilter.value.tagFilterType) {
            querys.push('tagFilterType=' + transactionsFilter.value.tagFilterType);
        }

        querys.push('dateType=' + transactionsFilter.value.dateType);

        if (DateRange.isBillingCycle(transactionsFilter.value.dateType) || transactionsFilter.value.dateType === DateRange.Custom.type) {
            querys.push('maxTime=' + transactionsFilter.value.maxTime);
            querys.push('minTime=' + transactionsFilter.value.minTime);
        }

        if (transactionsFilter.value.amountFilter) {
            querys.push('amountFilter=' + encodeURIComponent(transactionsFilter.value.amountFilter));
        }

        if (transactionsFilter.value.keyword) {
            querys.push('keyword=' + encodeURIComponent(transactionsFilter.value.keyword));
        }

        return querys.join('&');
    }

    function loadTransactions({ reload, count, page, withCount, autoExpand, defaultCurrency }: { reload?: boolean, count?: number, page?: number, withCount?: boolean, autoExpand: boolean, defaultCurrency: string }): Promise<TransactionPageWrapper> {
        let actualMaxTime = transactionsNextTimeId.value;

        if (reload && transactionsFilter.value.maxTime > 0) {
            actualMaxTime = transactionsFilter.value.maxTime * 1000 + 999;
        } else if (reload && transactionsFilter.value.maxTime <= 0) {
            actualMaxTime = 0;
        }

        return new Promise((resolve, reject) => {
            services.getTransactions({
                maxTime: actualMaxTime,
                minTime: transactionsFilter.value.minTime * 1000,
                count: count || 50,
                page: page || 1,
                withCount: !!withCount,
                type: transactionsFilter.value.type,
                categoryIds: transactionsFilter.value.categoryIds,
                accountIds: transactionsFilter.value.accountIds,
                tagIds: transactionsFilter.value.tagIds,
                tagFilterType: transactionsFilter.value.tagFilterType,
                amountFilter: transactionsFilter.value.amountFilter,
                keyword: transactionsFilter.value.keyword
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (reload) {
                        loadTransactionList({
                            transactionPageWrapper: EMPTY_TRANSACTION_RESULT,
                            reload: reload,
                            autoExpand: autoExpand,
                            defaultCurrency: defaultCurrency
                        });

                        if (!transactionListStateInvalid.value) {
                            updateTransactionListInvalidState(true);
                        }
                    }

                    reject({ message: 'Unable to retrieve transaction list' });
                    return;
                }

                const transactionPageWrapper: TransactionPageWrapper = {
                    items: Transaction.ofMulti(data.result.items),
                    totalCount: data.result.totalCount
                };

                loadTransactionList({
                    transactionPageWrapper: transactionPageWrapper,
                    reload: !!reload,
                    autoExpand: autoExpand,
                    defaultCurrency: defaultCurrency,
                    nextTimeSequenceId: data.result.nextTimeSequenceId
                });

                if (reload) {
                    if (transactionListStateInvalid.value) {
                        updateTransactionListInvalidState(false);
                    }
                }

                resolve(transactionPageWrapper);
            }).catch(error => {
                logger.error('failed to load transaction list', error);

                if (reload) {
                    loadTransactionList({
                        transactionPageWrapper: EMPTY_TRANSACTION_RESULT,
                        reload: reload,
                        autoExpand: autoExpand,
                        defaultCurrency: defaultCurrency
                    });

                    if (!transactionListStateInvalid.value) {
                        updateTransactionListInvalidState(true);
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
    }

    function loadMonthlyAllTransactions({ year, month, autoExpand, defaultCurrency }: { year: number, month: number, autoExpand: boolean, defaultCurrency: string }): Promise<TransactionPageWrapper> {
        return new Promise((resolve, reject) => {
            services.getAllTransactionsByMonth({
                year: year,
                month: month,
                type: transactionsFilter.value.type,
                categoryIds: transactionsFilter.value.categoryIds,
                accountIds: transactionsFilter.value.accountIds,
                tagIds: transactionsFilter.value.tagIds,
                tagFilterType: transactionsFilter.value.tagFilterType,
                amountFilter: transactionsFilter.value.amountFilter,
                keyword: transactionsFilter.value.keyword
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    loadTransactionList({
                        transactionPageWrapper: EMPTY_TRANSACTION_RESULT,
                        reload: true,
                        autoExpand: autoExpand,
                        defaultCurrency: defaultCurrency
                    });

                    if (!transactionListStateInvalid.value) {
                        updateTransactionListInvalidState(true);
                    }

                    reject({ message: 'Unable to retrieve transaction list' });
                    return;
                }

                const transactionPageWrapper: TransactionPageWrapper = {
                    items: Transaction.ofMulti(data.result.items),
                    totalCount: data.result.totalCount
                };

                loadTransactionList({
                    transactionPageWrapper: transactionPageWrapper,
                    reload: true,
                    autoExpand: autoExpand,
                    defaultCurrency: defaultCurrency
                });

                if (transactionListStateInvalid.value) {
                    updateTransactionListInvalidState(false);
                }

                resolve(transactionPageWrapper);
            }).catch(error => {
                logger.error('failed to load monthly all transaction list', error);

                loadTransactionList({
                    transactionPageWrapper: EMPTY_TRANSACTION_RESULT,
                    reload: true,
                    autoExpand: autoExpand,
                    defaultCurrency: defaultCurrency
                });

                if (!transactionListStateInvalid.value) {
                    updateTransactionListInvalidState(true);
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
    }

    function getTransaction({ transactionId, withPictures }: { transactionId: string, withPictures?: boolean }): Promise<Transaction> {
        return new Promise((resolve, reject) => {
            if (!isDefined(withPictures)) {
                withPictures = true;
            }

            services.getTransaction({
                id: transactionId,
                withPictures: withPictures
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve transaction' });
                    return;
                }

                const transaction = Transaction.of(data.result);

                resolve(transaction);
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
    }

    function saveTransaction({ transaction, defaultCurrency, isEdit, clientSessionId }: { transaction: Transaction, defaultCurrency: string, isEdit: boolean, clientSessionId: string }): Promise<Transaction> {
        return new Promise((resolve, reject) => {
            const actualTime = getActualUnixTimeForStore(transaction.time, transaction.utcOffset, getBrowserTimezoneOffsetMinutes());
            let promise: ApiResponsePromise<TransactionInfoResponse>;

            if (transaction.type !== TransactionType.Expense &&
                transaction.type !== TransactionType.Income &&
                transaction.type !== TransactionType.Transfer) {
                reject({ message: 'An error occurred' });
                return;
            }

            if (!isEdit) {
                promise = services.addTransaction(transaction.toCreateRequest(clientSessionId, actualTime));
            } else {
                promise = services.modifyTransaction(transaction.toModifyRequest(actualTime));
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add transaction' });
                    } else {
                        reject({ message: 'Unable to save transaction' });
                    }
                }

                const transaction = Transaction.of(data.result);

                if (!isEdit) {
                    if (!transactionListStateInvalid.value) {
                        updateTransactionListInvalidState(true);
                    }
                } else {
                    updateTransactionInTransactionList({
                        transaction: transaction,
                        defaultCurrency: defaultCurrency
                    });
                }

                if (!accountsStore.accountListStateInvalid) {
                    accountsStore.updateAccountListInvalidState(true);
                }

                if (!overviewStore.transactionOverviewStateInvalid) {
                    overviewStore.updateTransactionOverviewInvalidState(true);
                }

                if (!statisticsStore.transactionStatisticsStateInvalid) {
                    statisticsStore.updateTransactionStatisticsInvalidState(true);
                }

                resolve(transaction);
            }).catch(error => {
                logger.error('failed to save transaction', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add transaction' });
                    } else {
                        reject({ message: 'Unable to save transaction' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteTransaction({ transaction, defaultCurrency, beforeResolve }: { transaction: TransactionInfoResponse, defaultCurrency: string, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
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
                        removeTransactionFromTransactionList({
                            transaction: transaction,
                            defaultCurrency: defaultCurrency
                        });
                    });
                } else {
                    removeTransactionFromTransactionList({
                        transaction: transaction,
                        defaultCurrency: defaultCurrency
                    });
                }

                if (!accountsStore.accountListStateInvalid) {
                    accountsStore.updateAccountListInvalidState(true);
                }

                if (!overviewStore.transactionOverviewStateInvalid) {
                    overviewStore.updateTransactionOverviewInvalidState(true);
                }

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
    }

    function parseImportDsvFile({ fileType, fileEncoding, importFile }: { fileType: string, fileEncoding?: string, importFile: File }): Promise<string[][]> {
        return new Promise((resolve, reject) => {
            services.parseImportDsvFile({ fileType, fileEncoding, importFile }).then(response => {
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
    }

    function parseImportTransaction({ fileType, fileEncoding, importFile, columnMapping, transactionTypeMapping, hasHeaderLine, timeFormat, timezoneFormat, amountDecimalSeparator, amountDigitGroupingSymbol, geoSeparator, tagSeparator }: { fileType: string, fileEncoding?: string, importFile: File, columnMapping?: Record<number, number>, transactionTypeMapping?: Record<string, TransactionType>, hasHeaderLine?: boolean, timeFormat?: string, timezoneFormat?: string, amountDecimalSeparator?: string, amountDigitGroupingSymbol?: string, geoSeparator?: string, tagSeparator?: string }): Promise<ImportTransactionResponsePageWrapper> {
        return new Promise((resolve, reject) => {
            services.parseImportTransaction({ fileType, fileEncoding, importFile, columnMapping, transactionTypeMapping, hasHeaderLine, timeFormat, timezoneFormat, amountDecimalSeparator, amountDigitGroupingSymbol, geoSeparator, tagSeparator }).then(response => {
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
    }

    function importTransactions({ transactions, clientSessionId }: { transactions: ImportTransaction[], clientSessionId: string }): Promise<number> {
        const submitTransactions: TransactionCreateRequest[] = [];

        if (transactions) {
            for (const transaction of transactions) {
                const submitTransaction = transaction.toCreateRequest();
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
    }

    function getImportTransactionsProcess({ clientSessionId }: { clientSessionId: string }): Promise<number | null> {
        return new Promise((resolve, reject) => {
            services.getImportTransactionsProcess(clientSessionId).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to get transactions import process' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('Unable to get transactions import process', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to get transactions import process' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function uploadTransactionPicture({ pictureFile, clientSessionId }: { pictureFile: File, clientSessionId?: string }): Promise<TransactionPictureInfoBasicResponse> {
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
    }

    function removeUnusedTransactionPicture({ pictureInfo }: { pictureInfo: TransactionPictureInfoBasicResponse }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.removeUnusedTransactionPicture({ id: pictureInfo.pictureId }).then(response => {
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
    }

    function getTransactionPictureUrl(pictureInfo?: TransactionPictureInfoBasicResponse | null, disableBrowserCache?: boolean | string): string | undefined {
        if (!pictureInfo || !pictureInfo.originalUrl) {
            return undefined;
        }

        return services.getTransactionPictureUrlWithToken(pictureInfo.originalUrl, disableBrowserCache);
    }

    function collapseMonthInTransactionList({ month, collapse }: { month: TransactionMonthList, collapse: boolean }): void {
        if (month) {
            month.opened = !collapse;
        }
    }

    return {
        // states
        transactionDraft,
        transactionsFilter,
        transactions,
        transactionsNextTimeId,
        transactionListStateInvalid,
        // computed states
        allFilterCategoryIds,
        allFilterAccountIds,
        allFilterTagIds,
        allFilterCategoryIdsCount,
        allFilterAccountIdsCount,
        allFilterTagIdsCount,
        noTransaction,
        hasMoreTransaction,
        // functions
        initTransactionDraft,
        isTransactionDraftModified,
        saveTransactionDraft,
        clearTransactionDraft,
        setTransactionSuitableDestinationAmount,
        updateTransactionListInvalidState,
        resetTransactions,
        clearTransactions,
        initTransactionListFilter,
        updateTransactionListFilter,
        getTransactionListPageParams,
        loadTransactions,
        loadMonthlyAllTransactions,
        getTransaction,
        saveTransaction,
        deleteTransaction,
        parseImportDsvFile,
        parseImportTransaction,
        importTransactions,
        getImportTransactionsProcess,
        uploadTransactionPicture,
        removeUnusedTransactionPicture,
        getTransactionPictureUrl,
        collapseMonthInTransactionList
    };
});
