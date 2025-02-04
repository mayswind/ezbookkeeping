import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import type { WritableStartEndTime } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import type {
    TransactionAmountsRequestType,
    TransactionAmountsRequestParams,
    TransactionAmountsResponse,
    TransactionOverviewResponse
} from '@/models/transaction.ts';
import { ALL_TRANSACTION_AMOUNTS_REQUEST_TYPE } from '@/models/transaction.ts';

import { isNumber, isEquals } from '@/lib/common.ts';
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
    loadLast11Months: boolean;
}

export const useOverviewStore = defineStore('overview', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const transactionDataRange = ref<TransactionDataRange>({
        today: {
            startTime: getTodayFirstUnixTime(),
            endTime: getTodayLastUnixTime()
        },
        thisWeek: {
            startTime: getThisWeekFirstUnixTime(userStore.currentUserFirstDayOfWeek),
            endTime: getThisWeekLastUnixTime(userStore.currentUserFirstDayOfWeek)
        },
        thisMonth: {
            startTime: getThisMonthFirstUnixTime(),
            endTime: getThisMonthLastUnixTime()
        },
        thisYear: {
            startTime: getThisYearFirstUnixTime(),
            endTime: getThisYearLastUnixTime()
        },
        lastMonth: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 1, 'months')
        },
        monthBeforeLastMonth: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 2, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 2, 'months')
        },
        monthBeforeLast2Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 3, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 3, 'months')
        },
        monthBeforeLast3Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 4, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 4, 'months')
        },
        monthBeforeLast4Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 5, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 5, 'months')
        },
        monthBeforeLast5Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 6, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 6, 'months')
        },
        monthBeforeLast6Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 7, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 7, 'months')
        },
        monthBeforeLast7Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 8, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 8, 'months')
        },
        monthBeforeLast8Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 9, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 9, 'months')
        },
        monthBeforeLast9Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 10, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 10, 'months')
        },
        monthBeforeLast10Months: {
            startTime: getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 11, 'months'),
            endTime: getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 11, 'months')
        }
    });

    const transactionOverviewOptions = ref<TransactionOverviewOptions>({
        loadLast11Months: false
    });

    const transactionOverviewData = ref<TransactionAmountsResponse>({});
    const transactionOverviewStateInvalid = ref<boolean>(true);

    const transactionOverview = computed<TransactionOverviewResponse>(() => {
        const overviewData = transactionOverviewData.value;

        if (!overviewData || !overviewData.thisMonth) {
            return {
                thisMonth: {
                    valid: false,
                    incomeAmount: 0,
                    expenseAmount: 0,
                    incompleteIncomeAmount: false,
                    incompleteExpenseAmount: false
                }
            } as TransactionOverviewResponse;
        }

        const finalOverviewData: TransactionOverviewResponse = {};
        const defaultCurrency = userStore.currentUserDefaultCurrency;

        ALL_TRANSACTION_AMOUNTS_REQUEST_TYPE.forEach(field => {
            const item = overviewData[field];

            if (!item) {
                return;
            }

            let totalIncomeAmount = 0;
            let totalExpenseAmount = 0;
            let hasUnCalculatedTotalIncome = false;
            let hasUnCalculatedTotalExpense = false;

            if (item.amounts) {
                for (let i = 0; i < item.amounts.length; i++) {
                    const amount = item.amounts[i];

                    if (amount.currency !== defaultCurrency) {
                        const incomeAmount = exchangeRatesStore.getExchangedAmount(amount.incomeAmount, amount.currency, defaultCurrency);
                        const expenseAmount = exchangeRatesStore.getExchangedAmount(amount.expenseAmount, amount.currency, defaultCurrency);

                        if (isNumber(incomeAmount)) {
                            totalIncomeAmount += Math.floor(incomeAmount);
                        } else {
                            hasUnCalculatedTotalIncome = true;
                        }

                        if (isNumber(expenseAmount)) {
                            totalExpenseAmount += Math.floor(expenseAmount);
                        } else {
                            hasUnCalculatedTotalExpense = true;
                        }
                    } else {
                        totalIncomeAmount += amount.incomeAmount;
                        totalExpenseAmount += amount.expenseAmount;
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

    function updateTransactionDateRange(): void {
        transactionDataRange.value.today.startTime = getTodayFirstUnixTime();
        transactionDataRange.value.today.endTime = getTodayLastUnixTime();

        transactionDataRange.value.thisWeek.startTime = getThisWeekFirstUnixTime(userStore.currentUserFirstDayOfWeek);
        transactionDataRange.value.thisWeek.endTime = getThisWeekLastUnixTime(userStore.currentUserFirstDayOfWeek);

        transactionDataRange.value.thisMonth.startTime = getThisMonthFirstUnixTime();
        transactionDataRange.value.thisMonth.endTime = getThisMonthLastUnixTime();

        transactionDataRange.value.thisYear.startTime = getThisYearFirstUnixTime();
        transactionDataRange.value.thisYear.endTime = getThisYearLastUnixTime();

        transactionDataRange.value.lastMonth.startTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 1, 'months');
        transactionDataRange.value.lastMonth.endTime = getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 1, 'months');

        transactionDataRange.value.monthBeforeLastMonth.startTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), 2, 'months');
        transactionDataRange.value.monthBeforeLastMonth.endTime = getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), 2, 'months');

        for (let i = 2; i <= 10; i++) {
            transactionDataRange.value[`monthBeforeLast${i}Months` as TransactionAmountsRequestType].startTime = getUnixTimeBeforeUnixTime(getThisMonthFirstUnixTime(), i + 1, 'months');
            transactionDataRange.value[`monthBeforeLast${i}Months` as TransactionAmountsRequestType].endTime = getUnixTimeBeforeUnixTime(getThisMonthLastUnixTime(), i + 1, 'months');
        }
    }

    function updateTransactionOverviewInvalidState(invalidState: boolean): void {
        transactionOverviewStateInvalid.value = invalidState;
    }

    function resetTransactionOverview(): void {
        updateTransactionDateRange();
        transactionOverviewOptions.value.loadLast11Months = false;
        transactionOverviewData.value = {};
        transactionOverviewStateInvalid.value = true;
    }

    function loadTransactionOverview({ force, loadLast11Months }: { force: boolean, loadLast11Months?: boolean }): Promise<TransactionAmountsResponse> {
        let dateChanged = false;
        let rangeChanged = false;

        if (transactionDataRange.value.today.startTime !== getTodayFirstUnixTime()) {
            dateChanged = true;
            updateTransactionDateRange();
        }

        if (loadLast11Months && !transactionOverviewOptions.value.loadLast11Months) {
            rangeChanged = true;
        }

        if (!dateChanged && !rangeChanged && !force && !transactionOverviewStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(transactionOverviewData.value);
            });
        }

        const requestParams: TransactionAmountsRequestParams = {
            useTransactionTimezone: settingsStore.appSettings.timezoneUsedForStatisticsInHomePage == TimezoneTypeForStatistics.TransactionTimezone.type,
            today: transactionDataRange.value.today,
            thisWeek: transactionDataRange.value.thisWeek,
            thisMonth: transactionDataRange.value.thisMonth,
            thisYear: transactionDataRange.value.thisYear
        };

        if (loadLast11Months) {
            requestParams.lastMonth = transactionDataRange.value.lastMonth;
            requestParams.monthBeforeLastMonth = transactionDataRange.value.monthBeforeLastMonth;
            requestParams.monthBeforeLast2Months = transactionDataRange.value.monthBeforeLast2Months;
            requestParams.monthBeforeLast3Months = transactionDataRange.value.monthBeforeLast3Months;
            requestParams.monthBeforeLast4Months = transactionDataRange.value.monthBeforeLast4Months;
            requestParams.monthBeforeLast5Months = transactionDataRange.value.monthBeforeLast5Months;
            requestParams.monthBeforeLast6Months = transactionDataRange.value.monthBeforeLast6Months;
            requestParams.monthBeforeLast7Months = transactionDataRange.value.monthBeforeLast7Months;
            requestParams.monthBeforeLast8Months = transactionDataRange.value.monthBeforeLast8Months;
            requestParams.monthBeforeLast9Months = transactionDataRange.value.monthBeforeLast9Months;
            requestParams.monthBeforeLast10Months = transactionDataRange.value.monthBeforeLast10Months;
        }

        return new Promise((resolve, reject) => {
            services.getTransactionAmounts(requestParams).then(response => {
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
                transactionOverviewOptions.value.loadLast11Months = !!loadLast11Months;

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
        loadTransactionOverview
    };
});
