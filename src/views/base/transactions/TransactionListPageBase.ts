import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { type TransactionListFilter, type TransactionMonthList, useTransactionsStore } from '@/stores/transaction.ts';

import { type TypeAndName, keys, entries } from '@/core/base.ts';
import type { NumeralSystem } from '@/core/numeral.ts';
import { type TextualYearMonthDay, type Year0BasedMonth, type LocalizedDateRange, type WeekDayValue, DateRange, DateRangeScene } from '@/core/datetime.ts';
import { AccountType } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import { DISPLAY_HIDDEN_AMOUNT, INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';
import { DEFAULT_TAG_GROUP_ID } from '@/consts/tag.ts';

import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTagGroup } from '@/models/transaction_tag_group.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';
import { type Transaction, TransactionTagFilter } from '@/models/transaction.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    getLocalDatetimeFromUnixTime,
    getSameDateTimeWithBrowserTimezone,
    getCurrentDateTime,
    parseDateTimeFromUnixTime,
    parseDateTimeFromUnixTimeWithTimezoneOffset,
    getYearMonthFirstUnixTime,
    isDateRangeMatchOneMonth
} from '@/lib/datetime.ts';

import {
    getUnifiedSelectedAccountsCurrencyOrDefaultCurrency
} from '@/lib/account.ts';

import {
    categoryTypeToTransactionType
} from '@/lib/category.ts';

export class TransactionListPageType implements TypeAndName {
    private static readonly allInstances: TransactionListPageType[] = [];
    private static readonly allInstancesByType: Record<number, TransactionListPageType> = {};

    public static readonly List = new TransactionListPageType(0, 'Transaction List');
    public static readonly Calendar = new TransactionListPageType(1, 'Transaction Calendar');

    public static readonly Default = TransactionListPageType.List;

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        TransactionListPageType.allInstances.push(this);
        TransactionListPageType.allInstancesByType[type] = this;
    }

    public static values(): TransactionListPageType[] {
        return TransactionListPageType.allInstances;
    }

    public static valueOf(type: number): TransactionListPageType | undefined {
        return TransactionListPageType.allInstancesByType[type];
    }
}

export function useTransactionListPageBase() {
    const {
        tt,
        getAllDateRanges,
        getCurrentNumeralSystemType,
        formatDateTimeToLongDateTime,
        formatDateTimeToLongDate,
        formatDateTimeToShortTime,
        formatDateTimeToGregorianLikeLongYearMonth,
        formatDateRange,
        formatAmountToLocalizedNumeralsWithCurrency
    } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionTagsStore = useTransactionTagsStore();
    const transactionsStore = useTransactionsStore();

    const pageType = ref<number>(TransactionListPageType.List.type);
    const loading = ref<boolean>(true);
    const customMinDatetime = ref<number>(0);
    const customMaxDatetime = ref<number>(0);
    const currentCalendarDate = ref<TextualYearMonthDay | ''>('');

    const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);
    const defaultCurrency = computed<string>(() => getUnifiedSelectedAccountsCurrencyOrDefaultCurrency(allAccountsMap.value, queryAllFilterAccountIds.value, userStore.currentUserDefaultCurrency));
    const showTotalAmountInTransactionListPage = computed<boolean>(() => settingsStore.appSettings.showTotalAmountInTransactionListPage);
    const showTagInTransactionListPage = computed<boolean>(() => settingsStore.appSettings.showTagInTransactionListPage);

    const allDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.Normal, true, !!accountsStore.getAccountStatementDate(query.value.accountIds)));

    const allAccounts = computed<Account[]>(() => accountsStore.allMixedPlainAccounts);
    const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
    const allAvailableAccountsCount = computed<number>(() => accountsStore.allAvailableAccountsCount);
    const allPrimaryCategories = computed<Record<string, TransactionCategory[]>>(() => {
        const primaryCategories: Record<string, TransactionCategory[]> = {};

        for (const [categoryType, categories] of entries(transactionCategoriesStore.allTransactionCategories)) {
            if (query.value.type && categoryTypeToTransactionType(parseInt(categoryType)) !== query.value.type) {
                continue;
            }

            primaryCategories[categoryType] = categories;
        }

        return primaryCategories;
    });
    const allCategories = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);
    const allAvailableCategoriesCount = computed<number>(() => {
        let totalCount = 0;

        for (const [categoryType, categories] of entries(transactionCategoriesStore.allTransactionCategories)) {
            if (query.value.type && categoryTypeToTransactionType(parseInt(categoryType)) !== query.value.type) {
                continue;
            }

            if (categories) {
                totalCount += categories.length;
            }
        }

        return totalCount;
    });
    const allTransactionTagGroupsWithDefault = computed<TransactionTagGroup[]>(() => {
        const allGroups: TransactionTagGroup[] = [];
        const defaultGroup = TransactionTagGroup.createNewTagGroup(tt('Default Group'));
        defaultGroup.id = DEFAULT_TAG_GROUP_ID;
        allGroups.push(defaultGroup);
        allGroups.push(...transactionTagsStore.allTransactionTagGroups);
        return allGroups;
    });
    const allTransactionTagsByGroup = computed<Record<string, TransactionTag[]>>(() => transactionTagsStore.allTransactionTagsByGroupMap);
    const allTransactionTags = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);
    const allAvailableTagsCount = computed<number>(() => transactionTagsStore.allAvailableTagsCount);

    const displayPageTypeName = computed<string>(() => {
        const type = TransactionListPageType.valueOf(pageType.value);

        if (type) {
            return tt(type.name);
        }

        return tt(TransactionListPageType.List.name);
    });

    const query = computed<TransactionListFilter>(() => transactionsStore.transactionsFilter);
    const queryDateRangeName = computed<string>(() => {
        if (query.value.dateType === DateRange.All.type) {
            return tt('Date');
        }

        return formatDateRange(query.value.dateType, query.value.minTime, query.value.maxTime);
    });
    const queryMinTime = computed<string>(() => formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(query.value.minTime)));
    const queryMaxTime = computed<string>(() => formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(query.value.maxTime)));
    const queryMonthlyData = computed<boolean>(() => isDateRangeMatchOneMonth(query.value.minTime, query.value.maxTime));
    const queryMonth = computed<Year0BasedMonth>(() => {
        if (!query.value.minTime || !query.value.maxTime) {
            return getCurrentDateTime().toGregorianCalendarYear0BasedMonth();
        }

        return parseDateTimeFromUnixTime(query.value.minTime).toGregorianCalendarYear0BasedMonth();
    });

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
        if (query.value.tagFilter === TransactionTagFilter.TransactionNoTagFilterValue) {
            return tt('Without Tags');
        }

        if (queryAllFilterTagIdsCount.value > 1) {
            return tt('Multiple Tags');
        }

        for (const tagId of keys(queryAllFilterTagIds.value)) {
            const tagName = allTransactionTags.value[tagId]?.name;

            if (tagName) {
                return tagName;
            }
        }

        return tt('Tags');
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
            displayAmount.push(formatAmountToLocalizedNumeralsWithCurrency(parseInt(amountFilterItems[i] as string), false));
        }

        return displayAmount.join(' ~ ');
    });

    const transactionCalendarMinDate = computed<Date>(() => getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(parseDateTimeFromUnixTime(query.value.minTime)).getUnixTime()));
    const transactionCalendarMaxDate = computed<Date>(() => getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(parseDateTimeFromUnixTime(query.value.maxTime)).getUnixTime()));

    const currentMonthTransactionData = computed<TransactionMonthList | null>(() => {
        const allTransactions = transactionsStore.transactions;

        if (!allTransactions || !allTransactions.length) {
            return null;
        }

        const currentMonthMinDate = parseDateTimeFromUnixTime(query.value.minTime);
        const currentYear = currentMonthMinDate.getGregorianCalendarYear();
        const currentMonth = currentMonthMinDate.getGregorianCalendarMonth();

        for (const transactionMonth of allTransactions) {
            if (transactionMonth.year === currentYear && transactionMonth.month === currentMonth) {
                return transactionMonth;
            }
        }

        return null;
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

    function hasSubCategoryInQuery(category: TransactionCategory): boolean {
        if (!category.subCategories || !category.subCategories.length) {
            return false;
        }

        for (const subCategory of category.subCategories) {
            if (queryAllFilterCategoryIds.value[subCategory.id]) {
                return true;
            }
        }

        return false;
    }

    function hasVisibleTagsInTagGroup(tagGroup: TransactionTagGroup): boolean {
        const tagsInGroup = allTransactionTagsByGroup.value[tagGroup.id];

        if (!tagsInGroup || !tagsInGroup.length) {
            return false;
        }

        for (const tag of tagsInGroup) {
            if (!tag.hidden || queryAllFilterTagIds.value[tag.id]) {
                return true;
            }
        }

        return false;
    }

    function isSameAsDefaultTimezoneOffsetMinutes(transaction: Transaction): boolean {
        return transaction.utcOffset === getTimezoneOffsetMinutes(transaction.time);
    }

    function formatAmount(amount: number, hideAmount: boolean, currencyCode: string): string {
        if (hideAmount) {
            return formatAmountToLocalizedNumeralsWithCurrency(DISPLAY_HIDDEN_AMOUNT, currencyCode);
        }

        return formatAmountToLocalizedNumeralsWithCurrency(amount, currencyCode);
    }

    function getDisplayTime(transaction: Transaction): string {
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
        return formatDateTimeToShortTime(dateTime);
    }

    function getDisplayLongDate(transaction: Transaction): string {
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
        return formatDateTimeToLongDate(dateTime);
    }

    function getDisplayLongYearMonth(transactionMonthList: TransactionMonthList): string {
        const yearMonthDateTime = parseDateTimeFromUnixTime(getYearMonthFirstUnixTime(transactionMonthList.yearDashMonth));
        return formatDateTimeToGregorianLikeLongYearMonth(yearMonthDateTime);
    }

    function getDisplayTimezone(transaction: Transaction): string {
        const utcOffset = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset));
        return `UTC${utcOffset}`;
    }

    function getDisplayTimeInDefaultTimezone(transaction: Transaction): string {
        const timezoneOffsetMinutes = getTimezoneOffsetMinutes(transaction.time);
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, timezoneOffsetMinutes);
        const utcOffset = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(getUtcOffsetByUtcOffsetMinutes(timezoneOffsetMinutes));
        return `${formatDateTimeToLongDateTime(dateTime)} (UTC${utcOffset})`;
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

    function getDisplayMonthTotalAmount(amount: number, currency: string | false, symbol: string, incomplete: boolean): string {
        const displayAmount = formatAmountToLocalizedNumeralsWithCurrency(amount, currency);
        return symbol + displayAmount + (incomplete ? INCOMPLETE_AMOUNT_SUFFIX : '');
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
        pageType,
        loading,
        customMinDatetime,
        customMaxDatetime,
        currentCalendarDate,
        // computed states
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
        allTransactionTagGroupsWithDefault,
        allTransactionTagsByGroup,
        allTransactionTags,
        allAvailableTagsCount,
        displayPageTypeName,
        query,
        queryDateRangeName,
        queryMinTime,
        queryMaxTime,
        queryMonthlyData,
        queryMonth,
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
        transactionCalendarMinDate,
        transactionCalendarMaxDate,
        currentMonthTransactionData,
        canAddTransaction,
        // functions
        hasSubCategoryInQuery,
        hasVisibleTagsInTagGroup,
        isSameAsDefaultTimezoneOffsetMinutes,
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
