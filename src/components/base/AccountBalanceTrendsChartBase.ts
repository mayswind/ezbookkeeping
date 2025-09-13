import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import {
    type UnixTimeRange,
    type YearUnixTime,
    type YearQuarterUnixTime,
    type YearMonthUnixTime,
    YearMonthDayUnixTime,
} from '@/core/datetime.ts';
import type { FiscalYearUnixTime } from '@/core/fiscalyear.ts';
import { ChartDateAggregationType } from '@/core/statistics.ts';
import type { AccountInfoResponse } from '@/models/account.ts';
import type { TransactionReconciliationStatementResponseItem } from '@/models/transaction.ts';

import { isDefined, isArray } from '@/lib/common.ts';
import { sumAmounts } from '@/lib/numeral.ts';
import {
    getGregorianCalendarYearAndMonthFromUnixTime,
    getYearFirstUnixTimeBySpecifiedUnixTime,
    getQuarterFirstUnixTimeBySpecifiedUnixTime,
    getMonthFirstUnixTimeBySpecifiedUnixTime,
    getDayFirstUnixTimeBySpecifiedUnixTime,
    getAllDaysStartAndEndUnixTimes,
    getFiscalYearStartUnixTime
} from '@/lib/datetime.ts';
import { getAllDateRangesByYearMonthRange } from '@/lib/statistics.ts';

export interface AccountBalanceUnixTimeAndBalanceRange extends UnixTimeRange {
    minUnixTimeOpeningBalance: number;
    minUnixTimeClosingBalance: number;
    maxUnixTimeClosingBalance: number;
}

export interface AccountBalanceTrendsChartItem {
    displayDate: string;
    openingBalance: number;
    closingBalance: number;
    minimumBalance: number;
    maximumBalance: number;
    medianBalance: number;
    averageBalance: number;
}

export interface CommonAccountBalanceTrendsChartProps {
    items: TransactionReconciliationStatementResponseItem[] | undefined;
    dateAggregationType?: number;
    fiscalYearStart: number;
    account: AccountInfoResponse;
}

export function useAccountBalanceTrendsChartBase(props: CommonAccountBalanceTrendsChartProps) {
    const {
        formatUnixTimeToShortDate,
        formatUnixTimeToGregorianLikeShortYear,
        formatUnixTimeToGregorianLikeShortYearMonth,
        formatUnixTimeToGregorianLikeYearQuarter,
        formatUnixTimeToGregorianLikeFiscalYear
    } = useI18n();

    const dataDateRange = computed<AccountBalanceUnixTimeAndBalanceRange | null>(() => {
        if (!props.items || props.items.length < 1) {
            return null;
        }

        let minUnixTime = Number.MAX_SAFE_INTEGER, maxUnixTime = 0;
        let minUnixTimeOpeningBalance = 0;
        let minUnixTimeClosingBalance = 0;
        let maxUnixTimeClosingBalance = 0;

        for (const item of props.items) {
            if (item.time < minUnixTime) {
                minUnixTime = item.time;
                minUnixTimeOpeningBalance = item.accountOpeningBalance;
                minUnixTimeClosingBalance = item.accountClosingBalance;
            }

            if (item.time > maxUnixTime) {
                maxUnixTime = item.time;
                maxUnixTimeClosingBalance = item.accountClosingBalance;
            }
        }

        if (minUnixTime >= Number.MAX_SAFE_INTEGER || maxUnixTime <= 0) {
            return null;
        }

        return {
            minUnixTime: minUnixTime,
            maxUnixTime: maxUnixTime,
            minUnixTimeOpeningBalance: minUnixTimeOpeningBalance,
            minUnixTimeClosingBalance: minUnixTimeClosingBalance,
            maxUnixTimeClosingBalance: maxUnixTimeClosingBalance
        };
    });

    const allDateRanges = computed<YearUnixTime[] | FiscalYearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] | YearMonthDayUnixTime[]>(() => {
        if (!dataDateRange.value) {
            return [];
        }

        if (!isDefined(props.dateAggregationType)) {
            return getAllDaysStartAndEndUnixTimes(dataDateRange.value.minUnixTime, dataDateRange.value.maxUnixTime);
        } else {
            const startYearMonth = getGregorianCalendarYearAndMonthFromUnixTime(dataDateRange.value.minUnixTime);
            const endYearMonth = getGregorianCalendarYearAndMonthFromUnixTime(dataDateRange.value.maxUnixTime);
            return getAllDateRangesByYearMonthRange(startYearMonth, endYearMonth, props.fiscalYearStart, props.dateAggregationType);
        }
    });

    const allDataItems = computed<AccountBalanceTrendsChartItem[]>(() => {
        const ret: AccountBalanceTrendsChartItem[] = [];

        if (!dataDateRange.value || !allDateRanges.value || allDateRanges.value.length < 1 || !props.items || props.items.length < 1) {
            return ret;
        }

        const dayDataItemsMap: Record<number, TransactionReconciliationStatementResponseItem[]> = {};

        for (const dateItem of props.items) {
            let dateRangeMinUnixTime = 0;

            if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                dateRangeMinUnixTime = getYearFirstUnixTimeBySpecifiedUnixTime(dateItem.time);
            } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
                dateRangeMinUnixTime = getFiscalYearStartUnixTime(dateItem.time, props.fiscalYearStart);
            } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                dateRangeMinUnixTime = getQuarterFirstUnixTimeBySpecifiedUnixTime(dateItem.time);
            } else if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                dateRangeMinUnixTime = getMonthFirstUnixTimeBySpecifiedUnixTime(dateItem.time);
            } else {
                dateRangeMinUnixTime = getDayFirstUnixTimeBySpecifiedUnixTime(dateItem.time);
            }

            const dataItems: TransactionReconciliationStatementResponseItem[] = dayDataItemsMap[dateRangeMinUnixTime] || [];
            dataItems.push(dateItem);

            dayDataItemsMap[dateRangeMinUnixTime] = dataItems;
        }

        let lastOpeningBalance = dataDateRange.value.minUnixTimeOpeningBalance;
        let lastClosingBalance = dataDateRange.value.minUnixTimeClosingBalance;
        let lastMinimumBalance = lastClosingBalance;
        let lastMaximumBalance = lastClosingBalance;
        let lastMedianBalance = lastClosingBalance;
        let lastAverageBalance = lastClosingBalance;

        for (const dateRange of allDateRanges.value) {
            const dataItems = dayDataItemsMap[dateRange.minUnixTime];

            let displayDate = '';

            if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                displayDate = formatUnixTimeToGregorianLikeShortYear(dateRange.minUnixTime);
            } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
                displayDate = formatUnixTimeToGregorianLikeFiscalYear(dateRange.minUnixTime);
            } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                displayDate = formatUnixTimeToGregorianLikeYearQuarter(dateRange.minUnixTime);
            } else if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                displayDate = formatUnixTimeToGregorianLikeShortYearMonth(dateRange.minUnixTime);
            } else {
                displayDate = formatUnixTimeToShortDate(dateRange.minUnixTime);
            }

            if (isArray(dataItems)) {
                if (dataItems.length < 1) {
                    continue;
                }

                dataItems.sort(function (data1: TransactionReconciliationStatementResponseItem, data2: TransactionReconciliationStatementResponseItem) {
                    return data1.time - data2.time;
                });

                const openingBalance = dataItems[0]!.accountOpeningBalance;
                const closingBalance = dataItems[dataItems.length - 1]!.accountClosingBalance;
                const minimumBalance = Math.min(...dataItems.map(item => item.accountClosingBalance));
                const maximumBalance = Math.max(...dataItems.map(item => item.accountClosingBalance));
                const medianBalance = dataItems[Math.floor(dataItems.length / 2)]!.accountClosingBalance;
                const averageBalance = Math.trunc(sumAmounts(dataItems.map(item => item.accountClosingBalance)) / dataItems.length);

                if (props.account.isAsset) {
                    lastOpeningBalance = openingBalance;
                    lastClosingBalance = closingBalance;
                    lastMinimumBalance = minimumBalance;
                    lastMaximumBalance = maximumBalance;
                    lastMedianBalance = medianBalance;
                    lastAverageBalance = averageBalance;
                } else if (props.account.isLiability) {
                    lastOpeningBalance = -openingBalance;
                    lastClosingBalance = -closingBalance;
                    lastMinimumBalance = -minimumBalance;
                    lastMaximumBalance = -maximumBalance;
                    lastMedianBalance = -medianBalance;
                    lastAverageBalance = -averageBalance;
                } else {
                    lastOpeningBalance = openingBalance;
                    lastClosingBalance = closingBalance;
                    lastMinimumBalance = minimumBalance;
                    lastMaximumBalance = maximumBalance;
                    lastMedianBalance = medianBalance;
                    lastAverageBalance = averageBalance;
                }
            }

            ret.push({
                displayDate: displayDate,
                openingBalance: lastOpeningBalance,
                closingBalance: lastClosingBalance,
                minimumBalance: lastMinimumBalance,
                maximumBalance: lastMaximumBalance,
                medianBalance: lastMedianBalance,
                averageBalance: lastAverageBalance
            });

            lastOpeningBalance = lastClosingBalance;
        }

        return ret;
    });

    const allDisplayDateRanges = computed<string[]>(() => {
        if (!allDataItems.value || allDataItems.value.length < 1) {
            return [];
        }

        return allDataItems.value.map(item => item.displayDate);
    });

    return {
        // computed states
        allDateRanges,
        allDataItems,
        allDisplayDateRanges
    };
}
