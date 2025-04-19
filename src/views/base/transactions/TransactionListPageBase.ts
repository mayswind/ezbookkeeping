import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { type TransactionListFilter, type TransactionMonthList, useTransactionsStore } from '@/stores/transaction.ts';

import { type LocalizedDateRange, DateRange, DateRangeScene } from '@/core/datetime.ts';
import { AccountType } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';

import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';
import type { Transaction } from '@/models/transaction.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    parseDateFromUnixTime,
    getUnixTime,
    getYearMonthFirstUnixTime
} from '@/lib/datetime.ts';

import {
    getUnifiedSelectedAccountsCurrencyOrDefaultCurrency
} from '@/lib/account.ts';

import {
    categoryTypeToTransactionType
} from '@/lib/category.ts';

export function useTransactionListPageBase() {
    const {
        tt,
        getAllDateRanges,
        formatUnixTimeToLongDateTime,
        formatUnixTimeToLongDate,
        formatUnixTimeToLongYearMonth,
        formatUnixTimeToShortTime,
        formatDateRange,
        formatAmountWithCurrency
    } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionTagsStore = useTransactionTagsStore();
    const transactionsStore = useTransactionsStore();

    const loading = ref<boolean>(true);
    const customMinDatetime = ref<number>(0);
    const customMaxDatetime = ref<number>(0);

    const currentTimezoneOffsetMinutes = computed<number>(() => getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone));
    const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);
    const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);
    const defaultCurrency = computed<string>(() => getUnifiedSelectedAccountsCurrencyOrDefaultCurrency(allAccountsMap.value, queryAllFilterAccountIds.value, userStore.currentUserDefaultCurrency));
    const showTotalAmountInTransactionListPage = computed<boolean>(() => settingsStore.appSettings.showTotalAmountInTransactionListPage);
    const showTagInTransactionListPage = computed<boolean>(() => settingsStore.appSettings.showTagInTransactionListPage);

    const allDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.Normal, true, !!accountsStore.getAccountStatementDate(query.value.accountIds)));

    const allAccounts = computed<Account[]>(() => accountsStore.allMixedPlainAccounts);
    const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
    const allAvailableAccountsCount = computed<number>(() => accountsStore.allAvailableAccountsCount);
    const allPrimaryCategories = computed<Record<number, TransactionCategory[]>>(() => {
        const primaryCategories: Record<number, TransactionCategory[]> = {};

        for (const categoryType in transactionCategoriesStore.allTransactionCategories) {
            if (!Object.prototype.hasOwnProperty.call(transactionCategoriesStore.allTransactionCategories, categoryType)) {
                continue;
            }

            if (query.value.type && categoryTypeToTransactionType(parseInt(categoryType)) !== query.value.type) {
                continue;
            }

            primaryCategories[categoryType] = transactionCategoriesStore.allTransactionCategories[categoryType];
        }

        return primaryCategories;
    });
    const allCategories = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);
    const allAvailableCategoriesCount = computed<number>(() => {
        let totalCount = 0;

        for (const categoryType in transactionCategoriesStore.allTransactionCategories) {
            if (!Object.prototype.hasOwnProperty.call(transactionCategoriesStore.allTransactionCategories, categoryType)) {
                continue;
            }

            if (query.value.type && categoryTypeToTransactionType(parseInt(categoryType)) !== query.value.type) {
                continue;
            }

            if (transactionCategoriesStore.allTransactionCategories[categoryType]) {
                totalCount += transactionCategoriesStore.allTransactionCategories[categoryType].length;
            }
        }

        return totalCount;

    });
    const allTransactionTags = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);
    const allAvailableTagsCount = computed<number>(() => transactionTagsStore.allAvailableTagsCount);

    const query = computed<TransactionListFilter>(() => transactionsStore.transactionsFilter);
    const queryDateRangeName = computed<string>(() => {
        if (query.value.dateType === DateRange.All.type) {
            return tt('Date');
        }

        return formatDateRange(query.value.dateType, query.value.minTime, query.value.maxTime);
    });
    const queryMinTime = computed<string>(() => formatUnixTimeToLongDateTime(query.value.minTime));
    const queryMaxTime = computed<string>(() => formatUnixTimeToLongDateTime(query.value.maxTime));
    const queryAllFilterCategoryIds = computed<Record<string, boolean>>(() => transactionsStore.allFilterCategoryIds);
    const queryAllFilterAccountIds = computed<Record<string, boolean>>(() => transactionsStore.allFilterAccountIds);
    const queryAllFilterTagIds = computed<Record<string, boolean>>(() => transactionsStore.allFilterTagIds);
    const queryAllFilterCategoryIdsCount = computed<number>(() => transactionsStore.allFilterCategoryIdsCount);
    const queryAllFilterAccountIdsCount = computed<number>(() => transactionsStore.allFilterAccountIdsCount);
    const queryAllFilterTagIdsCount = computed<number>(() => transactionsStore.allFilterTagIdsCount);

    const queryAccountName = computed<string>(() => {
        if (queryAllFilterAccountIdsCount.value > 1) {
            return tt('Multiple Accounts');
        }

        return allAccountsMap.value[query.value.accountIds]?.name || tt('Account');
    });

    const queryCategoryName = computed<string>(() => {
        if (queryAllFilterCategoryIdsCount.value > 1) {
            return tt('Multiple Categories');
        }

        return allCategories.value[query.value.categoryIds]?.name || tt('Category');
    });

    const queryTagName = computed<string>(() => {
        if (query.value.tagIds === 'none') {
            return tt('Without Tags');
        }

        if (queryAllFilterTagIdsCount.value > 1) {
            return tt('Multiple Tags');
        }

        return allTransactionTags.value[query.value.tagIds]?.name || tt('Tags');
    });

    const queryAmount = computed<string>(() => {
        if (!query.value.amountFilter) {
            return '';
        }

        const amountFilterItems = query.value.amountFilter.split(':');

        if (amountFilterItems.length < 2) {
            return '';
        }

        const displayAmount: string[] = [];

        for (let i = 1; i < amountFilterItems.length; i++) {
            displayAmount.push(formatAmountWithCurrency(amountFilterItems[i], false));
        }

        return displayAmount.join(' ~ ');
    });

    const canAddTransaction = computed<boolean>(() => {
        if (query.value.accountIds && queryAllFilterAccountIdsCount.value === 1) {
            const account = allAccountsMap.value[query.value.accountIds];

            if (account && account.type === AccountType.MultiSubAccounts.type) {
                return false;
            }
        }

        return true;
    });

    function formatAmount(amount: number, hideAmount: boolean, currencyCode: string): string {
        if (hideAmount) {
            return formatAmountWithCurrency('***', currencyCode);
        }

        return formatAmountWithCurrency(amount, currencyCode);
    }

    function getDisplayTime(transaction: Transaction): string {
        return formatUnixTimeToShortTime(transaction.time, transaction.utcOffset, currentTimezoneOffsetMinutes.value);
    }

    function getDisplayLongDate(transaction: Transaction): string {
        const transactionTime = getUnixTime(parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentTimezoneOffsetMinutes.value));
        return formatUnixTimeToLongDate(transactionTime);
    }

    function getDisplayLongYearMonth(transactionMonthList: TransactionMonthList): string {
        return formatUnixTimeToLongYearMonth(getYearMonthFirstUnixTime(transactionMonthList.yearMonth));
    }

    function getDisplayTimezone(transaction: Transaction): string {
        return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
    }

    function getDisplayTimeInDefaultTimezone(transaction: Transaction): string {
        return `${formatUnixTimeToLongDateTime(transaction.time)} (UTC${getTimezoneOffset(settingsStore.appSettings.timeZone)})`;
    }

    function getDisplayAmount(transaction: Transaction): string {
        if (queryAllFilterAccountIdsCount.value < 1) {
            if (transaction.sourceAccount) {
                return formatAmount(transaction.sourceAmount, transaction.hideAmount, transaction.sourceAccount.currency);
            }
        } else if (queryAllFilterAccountIdsCount.value === 1) {
            if (transaction.sourceAccount && (queryAllFilterAccountIds.value[transaction.sourceAccount.id] || queryAllFilterAccountIds.value[transaction.sourceAccount.parentId])) {
                return formatAmount(transaction.sourceAmount, transaction.hideAmount, transaction.sourceAccount.currency);
            } else if (transaction.destinationAccount && (queryAllFilterAccountIds.value[transaction.destinationAccount.id] || queryAllFilterAccountIds.value[transaction.destinationAccount.parentId])) {
                return formatAmount(transaction.destinationAmount, transaction.hideAmount, transaction.destinationAccount.currency);
            }
        } else { // queryAllFilterAccountIdsCount.value > 1
            if (transaction.sourceAccount && transaction.destinationAccount) {
                if ((queryAllFilterAccountIds.value[transaction.sourceAccount.id] || queryAllFilterAccountIds.value[transaction.sourceAccount.parentId])
                    && !queryAllFilterAccountIds.value[transaction.destinationAccount.id] && !queryAllFilterAccountIds.value[transaction.destinationAccount.parentId]) {
                    return formatAmount(transaction.sourceAmount, transaction.hideAmount, transaction.sourceAccount.currency);
                } else if ((queryAllFilterAccountIds.value[transaction.destinationAccount.id] || queryAllFilterAccountIds.value[transaction.destinationAccount.parentId])
                    && !queryAllFilterAccountIds.value[transaction.sourceAccount.id] && !queryAllFilterAccountIds.value[transaction.sourceAccount.parentId]) {
                    return formatAmount(transaction.destinationAmount, transaction.hideAmount, transaction.destinationAccount.currency);
                }
            }
        }

        if (transaction.sourceAccount) {
            return formatAmount(transaction.sourceAmount, transaction.hideAmount, transaction.sourceAccount.currency);
        }

        return '';
    }

    function getDisplayMonthTotalAmount(amount: number, currency: string, symbol: string, incomplete: boolean): string {
        const displayAmount = formatAmountWithCurrency(amount, currency);
        return symbol + displayAmount + (incomplete ? '+' : '');
    }

    function getTransactionTypeName(type: number | null, defaultName: string): string {
        switch (type){
            case TransactionType.ModifyBalance:
                return tt('Modify Balance');
            case TransactionType.Income:
                return tt('Income');
            case TransactionType.Expense:
                return tt('Expense');
            case TransactionType.Transfer:
                return tt('Transfer');
            default:
                return tt(defaultName);
        }
    }

    return {
        // states
        loading,
        customMinDatetime,
        customMaxDatetime,
        // computed states
        currentTimezoneOffsetMinutes,
        firstDayOfWeek,
        fiscalYearStart,
        defaultCurrency,
        showTotalAmountInTransactionListPage,
        showTagInTransactionListPage,
        allDateRanges,
        allAccounts,
        allAccountsMap,
        allAvailableAccountsCount,
        allCategories,
        allPrimaryCategories,
        allAvailableCategoriesCount,
        allTransactionTags,
        allAvailableTagsCount,
        query,
        queryDateRangeName,
        queryMinTime,
        queryMaxTime,
        queryAllFilterCategoryIds,
        queryAllFilterAccountIds,
        queryAllFilterTagIds,
        queryAllFilterCategoryIdsCount,
        queryAllFilterAccountIdsCount,
        queryAllFilterTagIdsCount,
        queryAccountName,
        queryCategoryName,
        queryTagName,
        queryAmount,
        canAddTransaction,
        // functions
        getDisplayTime,
        getDisplayLongDate,
        getDisplayLongYearMonth,
        getDisplayTimezone,
        getDisplayTimeInDefaultTimezone,
        getDisplayAmount,
        getDisplayMonthTotalAmount,
        getTransactionTypeName,
    };
}
