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
import type { TransactionReconciliationStatementResponseItem } from '@/models/transaction.ts';

import { isDefined, isArray } from '@/lib/common.ts';
import {
    getYearAndMonthFromUnixTime,
    getYearFirstUnixTimeBySpecifiedUnixTime,
    getQuarterFirstUnixTimeBySpecifiedUnixTime,
    getMonthFirstUnixTimeBySpecifiedUnixTime,
    getDayFirstUnixTimeBySpecifiedUnixTime,
    getAllDaysStartAndEndUnixTimes,
    getFiscalYearStartUnixTime
} from '@/lib/datetime.ts';
import { getAllDateRangesByYearMonthRange } from '@/lib/statistics.ts';

export interface AccountBalanceUnixTimeAndBalanceRange extends UnixTimeRange {
    minUnixTimeBalance: number;
    maxUnixTimeBalance: number;
}

export interface AccountBalanceTrendsChartItem {
    displayDate: string;
    amount: number;
}

export interface CommonAccountBalanceTrendsChartProps {
    items: TransactionReconciliationStatementResponseItem[] | undefined;
    dateAggregationType?: number;
    fiscalYearStart: number;
    accountCurrency: string;
}

export function useAccountBalanceTrendsChartBase(props: CommonAccountBalanceTrendsChartProps) {
    const { formatUnixTimeToShortDate, formatUnixTimeToShortYear, formatUnixTimeToShortYearMonth, formatUnixTimeToYearQuarter, formatUnixTimeToFiscalYear } = useI18n();

    const dataDateRange = computed<AccountBalanceUnixTimeAndBalanceRange | null>(() => {
        if (!props.items || props.items.length < 1) {
            return null;
        }

        let minUnixTime = Number.MAX_SAFE_INTEGER, maxUnixTime = 0;
        let minUnixTimeBalance = 0, maxUnixTimeBalance = 0;

        for (let i = 0; i < props.items.length; i++) {
            const item = props.items[i];

            if (item.time < minUnixTime) {
                minUnixTime = item.time;
                minUnixTimeBalance = item.accountBalance;
            }

            if (item.time > maxUnixTime) {
                maxUnixTime = item.time;
                maxUnixTimeBalance = item.accountBalance;
            }
        }

        if (minUnixTime >= Number.MAX_SAFE_INTEGER || maxUnixTime <= 0) {
            return null;
        }

        return {
            minUnixTime: minUnixTime,
            maxUnixTime: maxUnixTime,
            minUnixTimeBalance: minUnixTimeBalance,
            maxUnixTimeBalance: maxUnixTimeBalance
        };
    });

    const allDateRanges = computed<YearUnixTime[] | FiscalYearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] | YearMonthDayUnixTime[]>(() => {
        if (!dataDateRange.value) {
            return [];
        }

        if (!isDefined(props.dateAggregationType)) {
            return getAllDaysStartAndEndUnixTimes(dataDateRange.value.minUnixTime, dataDateRange.value.maxUnixTime);
        } else {
            const startYearMonth = getYearAndMonthFromUnixTime(dataDateRange.value.minUnixTime);
            const endYearMonth = getYearAndMonthFromUnixTime(dataDateRange.value.maxUnixTime);
            return getAllDateRangesByYearMonthRange(startYearMonth, endYearMonth, props.fiscalYearStart, props.dateAggregationType);
        }
    });

    const allDataItems = computed<AccountBalanceTrendsChartItem[]>(() => {
        const ret: AccountBalanceTrendsChartItem[] = [];

        if (!dataDateRange.value || !allDateRanges.value || allDateRanges.value.length < 1 || !props.items || props.items.length < 1) {
            return ret;
        }

        const dayDataItemsMap: Record<number, TransactionReconciliationStatementResponseItem[]> = {};

        for (let i = 0; i < props.items.length; i++) {
            const dateItem = props.items[i];
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

        let lastAmount = dataDateRange.value.minUnixTimeBalance;

        for (let i = 0; i < allDateRanges.value.length; i++) {
            const dateRange = allDateRanges.value[i];
            const dataItems = dayDataItemsMap[dateRange.minUnixTime];

            let displayDate = '';

            if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                displayDate = formatUnixTimeToShortYear(dateRange.minUnixTime);
            } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
                displayDate = formatUnixTimeToFiscalYear(dateRange.minUnixTime);
            } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                displayDate = formatUnixTimeToYearQuarter(dateRange.minUnixTime);
            } else if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                displayDate = formatUnixTimeToShortYearMonth(dateRange.minUnixTime);
            } else {
                displayDate = formatUnixTimeToShortDate(dateRange.minUnixTime);
            }

            if (isArray(dataItems)) {
                let lastUnixTime = 0;

                for (let i = 0; i < dataItems.length; i++) {
                    const dataItem = dataItems[i];

                    if (dataItem.time >= lastUnixTime) {
                        lastUnixTime = dataItem.time;
                        lastAmount = dataItem.accountBalance;
                    }
                }
            }

            ret.push({
                displayDate: displayDate,
                amount: lastAmount
            });
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
