import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import { type WritableStartEndTime, DateRange } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import type { TransactionType } from '@/core/transaction.ts';

import type {
    TransactionAmountsRequestType,
    TransactionAmountsRequestParams,
    TransactionAmountsResponse,
    TransactionOverviewData
} from '@/models/transaction.ts';
import {
    ALL_TRANSACTION_AMOUNTS_REQUEST_TYPE,
    LATEST_12MONTHS_TRANSACTION_AMOUNTS_REQUEST_TYPES
} from '@/models/transaction.ts';

import {
    isDefined,
    isNumber,
    isEquals,
    isObjectEmpty,
    normalizeInteger,
    objectFieldWithValueToArrayItem
} from '@/lib/common.ts';
import {
    BIG_DECIMAL_ZERO,
    parseBigDecimal
} from '@/lib/numeral.ts';
import {
    getUnixTimeBeforeUnixTime,
    getTodayFirstUnixTime,
    getTodayLastUnixTime,
    getThisWeekFirstUnixTime,
    getThisWeekLastUnixTime,
    getThisMonthFirstUnixTime,
    getThisMonthLastUnixTime,
    getThisYearFirstUnixTime,
    getThisYearLastUnixTime
} from '@/lib/datetime.ts';
import { getFinalAccountIdsByFilteredAccountIds } from '@/lib/account.ts';
import { getFinalCategoryIdsByFilteredCategoryIds } from '@/lib/category.ts';
import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

interface TransactionDataRange extends Record<TransactionAmountsRequestType, WritableStartEndTime> {
    today: {
        startTime: number;
        endTime: number;
    };
    thisWeek: {
        startTime: number;
        endTime: number;
    };
    thisMonth: {
        startTime: number;
        endTime: number;
    };
    thisYear: {
        startTime: number;
        endTime: number;
    };
    lastMonth: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLastMonth: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast2Months: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast3Months: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast4Months: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast5Months: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast6Months: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast7Months: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast8Months: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast9Months: {
        startTime: number;
        endTime: number;
    };
    monthBeforeLast10Months: {
        startTime: number;
        endTime: number;
    };
}

interface TransactionOverviewOptions {
    loadedMonths: number;
}

export const useOverviewStore = defineStore('overview', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const transactionDataRange = ref<TransactionDataRange>(getTransactionDateRange());

    const transactionOverviewOptions = ref<TransactionOverviewOptions>({
        loadedMonths: 1
    });

    const transactionOverviewData = ref<TransactionAmountsResponse>({});
    const transactionOverviewStateInvalid = ref<boolean>(true);

    const transactionOverview = computed<TransactionOverviewData>(() => {
        const overviewData = transactionOverviewData.value;

        if (!overviewData || !overviewData.thisMonth) {
            return {
                thisMonth: {
                    valid: false,
                    incomeAmount: BIG_DECIMAL_ZERO,
                    expenseAmount: BIG_DECIMAL_ZERO,
                    incompleteIncomeAmount: false,
                    incompleteExpenseAmount: false
                }
            } as TransactionOverviewData;
        }

        const finalOverviewData: TransactionOverviewData = {};
        const defaultCurrency = userStore.currentUserDefaultCurrency;

        ALL_TRANSACTION_AMOUNTS_REQUEST_TYPE.forEach(field => {
            const item = overviewData[field];

            if (!item) {
                return;
            }

            let totalIncomeAmount = BIG_DECIMAL_ZERO;
            let totalExpenseAmount = BIG_DECIMAL_ZERO;
            let hasUnCalculatedTotalIncome = false;
            let hasUnCalculatedTotalExpense = false;

            if (item.amounts) {
                for (const amount of item.amounts) {
                    if (amount.currency !== defaultCurrency) {
                        const incomeAmount = exchangeRatesStore.getExchangedAmount(parseBigDecimal(amount.incomeAmount), amount.currency, defaultCurrency);
                        const expenseAmount = exchangeRatesStore.getExchangedAmount(parseBigDecimal(amount.expenseAmount), amount.currency, defaultCurrency);

                        if (incomeAmount) {
                            totalIncomeAmount = totalIncomeAmount.add(incomeAmount.truncate());
                        } else {
                            hasUnCalculatedTotalIncome = true;
                        }

                        if (expenseAmount) {
                            totalExpenseAmount = totalExpenseAmount.add(expenseAmount.truncate());
                        } else {
                            hasUnCalculatedTotalExpense = true;
                        }
                    } else {
                        totalIncomeAmount = totalIncomeAmount.add(parseBigDecimal(amount.incomeAmount));
                        totalExpenseAmount = totalExpenseAmount.add(parseBigDecimal(amount.expenseAmount));
                    }
                }
            }

            finalOverviewData[field] = {
                valid: true,
                incomeAmount: totalIncomeAmount,
                expenseAmount: totalExpenseAmount,
                incompleteIncomeAmount: hasUnCalculatedTotalIncome,
                incompleteExpenseAmount: hasUnCalculatedTotalExpense,
                amounts: item.amounts || []
            };
        });

        return finalOverviewData;
    });

    function getTransactionDateRange(): TransactionDataRange {
        const dateRange: TransactionDataRange = {
            today: { startTime: 0, endTime: 0 },
            thisWeek: { startTime: 0, endTime: 0 },
            thisMonth: { startTime: 0, endTime: 0 },
            thisYear: { startTime: 0, endTime: 0 },
            lastMonth: { startTime: 0, endTime: 0 },
            monthBeforeLastMonth: { startTime: 0, endTime: 0 },
            monthBeforeLast2Months: { startTime: 0, endTime: 0 },
            monthBeforeLast3Months: { startTime: 0, endTime: 0 },
            monthBeforeLast4Months: { startTime: 0, endTime: 0 },
            monthBeforeLast5Months: { startTime: 0, endTime: 0 },
            monthBeforeLast6Months: { startTime: 0, endTime: 0 },
            monthBeforeLast7Months: { startTime: 0, endTime: 0 },
            monthBeforeLast8Months: { startTime: 0, endTime: 0 },
            monthBeforeLast9Months: { startTime: 0, endTime: 0 },
            monthBeforeLast10Months: { startTime: 0, endTime: 0 }
        };

        initTransactionDateRange(dateRange);
        return dateRange;
    }

    function initTransactionDateRange(dateRange: TransactionDataRange): void {
        dateRange.today.startTime = getTodayFirstUnixTime();
        dateRange.today.endTime = getTodayLastUnixTime();

        dateRange.thisWeek.startTime = getThisWeekFirstUnixTime(userStore.currentUserFirstDayOfWeek);
        dateRange.thisWeek.endTime = getThisWeekLastUnixTime(userStore.currentUserFirstDayOfWeek);

        dateRange.thisMonth.startTime = getThisMonthFirstUnixTime();
        dateRange.thisMonth.endTime = getThisMonthLastUnixTime();

        dateRange.thisYear.startTime = getThisYearFirstUnixTime();
        dateRange.thisYear.endTime = getThisYearLastUnixTime();

        dateRange.lastMonth.startTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'months');
        dateRange.lastMonth.endTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'seconds');

        dateRange.monthBeforeLastMonth.startTime = getUnixTimeBeforeUnixTime(dateRange.lastMonth.startTime, 1, 'months');
        dateRange.monthBeforeLastMonth.endTime = getUnixTimeBeforeUnixTime(dateRange.lastMonth.startTime, 1, 'seconds');

        dateRange.monthBeforeLast2Months.startTime = getUnixTimeBeforeUnixTime(dateRange.monthBeforeLastMonth.startTime, 1, 'months');
        dateRange.monthBeforeLast2Months.endTime = getUnixTimeBeforeUnixTime(dateRange.monthBeforeLastMonth.startTime, 1, 'seconds');

        for (let i = 3; i <= 10; i++) {
            dateRange[`monthBeforeLast${i}Months` as TransactionAmountsRequestType].startTime = getUnixTimeBeforeUnixTime(dateRange[`monthBeforeLast${i - 1}Months` as TransactionAmountsRequestType].startTime, 1, 'months');
            dateRange[`monthBeforeLast${i}Months` as TransactionAmountsRequestType].endTime = getUnixTimeBeforeUnixTime(dateRange[`monthBeforeLast${i - 1}Months` as TransactionAmountsRequestType].startTime, 1, 'seconds');
        }
    }

    function updateTransactionDateRange(): void {
        initTransactionDateRange(transactionDataRange.value);
    }

    function updateTransactionOverviewInvalidState(invalidState: boolean): void {
        transactionOverviewStateInvalid.value = invalidState;
    }

    function resetTransactionOverview(): void {
        updateTransactionDateRange();
        transactionOverviewOptions.value.loadedMonths = 1;
        transactionOverviewData.value = {};
        transactionOverviewStateInvalid.value = true;
    }

    function loadTransactionOverview({ force, months }: { force: boolean, months?: number }): Promise<TransactionAmountsResponse> {
        const requestedMonths: number = normalizeInteger(months, 1, 1, 12);
        let dateChanged: boolean = false;
        let rangeChanged: boolean = false;

        if (transactionDataRange.value.today.startTime !== getTodayFirstUnixTime()) {
            dateChanged = true;
            updateTransactionDateRange();
        }

        if (requestedMonths > transactionOverviewOptions.value.loadedMonths) {
            rangeChanged = true;
        }

        if (!dateChanged && !rangeChanged && !force && !transactionOverviewStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(transactionOverviewData.value);
            });
        }

        const requestParams: TransactionAmountsRequestParams = {
            useTransactionTimezone: settingsStore.appSettings.timezoneUsedForStatisticsInHomePage === TimezoneTypeForStatistics.TransactionTimezone.type,
            today: transactionDataRange.value.today,
            thisWeek: transactionDataRange.value.thisWeek,
            thisMonth: transactionDataRange.value.thisMonth,
            thisYear: transactionDataRange.value.thisYear
        };

        const requestedMonthTypes: TransactionAmountsRequestType[] = LATEST_12MONTHS_TRANSACTION_AMOUNTS_REQUEST_TYPES.slice(-requestedMonths, -1);

        for (const requestType of requestedMonthTypes) {
            requestParams[requestType] = transactionDataRange.value[requestType];
        }

        const excludeAccountIds: string[] = objectFieldWithValueToArrayItem(settingsStore.appSettings.overviewAccountFilterInHomePage, true);
        const excludeCategoryIds: string[] = objectFieldWithValueToArrayItem(settingsStore.appSettings.overviewTransactionCategoryFilterInHomePage, true);

        return new Promise((resolve, reject) => {
            services.getTransactionAmounts(requestParams, excludeAccountIds, excludeCategoryIds).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve transaction overview' });
                    return;
                }

                if (transactionOverviewStateInvalid.value) {
                    updateTransactionOverviewInvalidState(false);
                }

                if (force && data.result && isEquals(transactionOverviewData.value, data.result)) {
                    reject({ message: 'Data is up to date', isUpToDate: true });
                    return;
                }

                transactionOverviewData.value = data.result;
                transactionOverviewOptions.value.loadedMonths = requestedMonths;

                resolve(data.result);
            }).catch(error => {
                if (force) {
                    logger.error('failed to force load transaction overview', error);
                } else {
                    logger.error('failed to load transaction overview', error);
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve transaction overview' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getTransactionListPageParams({ type, dateType, minTime, maxTime }: { type?: TransactionType, dateType?: number, minTime?: number, maxTime?: number }): string {
        const querys: string[] = [];

        if (isDefined(type)) {
            querys.push('type=' + type);
        }

        if (isDefined(dateType)) {
            querys.push('dateType=' + dateType);

            if (dateType === DateRange.Custom.type) {
                if (isNumber(minTime) && minTime > 0) {
                    querys.push('minTime=' + minTime);
                }

                if (isNumber(maxTime) && maxTime > 0) {
                    querys.push('maxTime=' + maxTime);
                }
            }
        }

        if (!isObjectEmpty(settingsStore.appSettings.overviewTransactionCategoryFilterInHomePage)) {
            querys.push('categoryIds=' + getFinalCategoryIdsByFilteredCategoryIds(transactionCategoriesStore.allTransactionCategoriesMap, settingsStore.appSettings.overviewTransactionCategoryFilterInHomePage));
        }

        if (!isObjectEmpty(settingsStore.appSettings.overviewAccountFilterInHomePage)) {
            querys.push('accountIds=' + getFinalAccountIdsByFilteredAccountIds(accountsStore.allAccountsMap, settingsStore.appSettings.overviewAccountFilterInHomePage));
        }

        return querys.join('&');
    }

    return {
        // states
        transactionDataRange,
        transactionOverviewOptions,
        transactionOverviewData,
        transactionOverviewStateInvalid,
        // computed states,
        transactionOverview,
        // functions
        updateTransactionOverviewInvalidState,
        resetTransactionOverview,
        loadTransactionOverview,
        getTransactionListPageParams
    };
});
