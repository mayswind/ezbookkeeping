import {
    type TextualYearMonth,
    type Year1BasedMonth,
    type YearUnixTime,
    type YearQuarterUnixTime,
    type YearMonthUnixTime,
    YearMonthDayUnixTime
} from '@/core/datetime.ts';
import {
    type FiscalYearUnixTime
} from '@/core/fiscalyear.ts';
import {
    ChartSortingType,
    ChartDateAggregationType
} from '@/core/statistics.ts';

import type {
    YearMonthItems,
    SortableTransactionStatisticDataItem,
    TransactionStatisticResponse,
    TransactionStatisticResponseItem,
    TransactionStatisticTrendsResponseItem,
    TransactionStatisticAssetTrendsResponseItem
} from '@/models/transaction.ts';

import {
    isArrayEquals
} from '@/lib/common.ts';

import {
    getAllMonthsStartAndEndUnixTimes,
    getAllQuartersStartAndEndUnixTimes,
    getAllYearsStartAndEndUnixTimes,
    getAllFiscalYearsStartAndEndUnixTimes
} from '@/lib/datetime.ts';

function isTransactionStatisticItemsEquals(items1: TransactionStatisticResponseItem[], items2: TransactionStatisticResponseItem[]): boolean {
    return isArrayEquals(items1, items2,
        item => `${item.categoryId}_${item.accountId}_${item.relatedAccountId}_${item.relatedAccountType}`,
        (item1, item2) => item1.categoryId === item2.categoryId &&
            item1.accountId === item2.accountId &&
            item1.relatedAccountId === item2.relatedAccountId &&
            item1.relatedAccountType === item2.relatedAccountType &&
            item1.amount === item2.amount);
}

export function sortStatisticsItems<T extends SortableTransactionStatisticDataItem>(items: T[], sortingType: number): void {
    if (sortingType === ChartSortingType.DisplayOrder.type) {
        items.sort(function (data1, data2) {
            for (let i = 0; i < Math.min(data1.displayOrders.length, data2.displayOrders.length); i++) {
                if (data1.displayOrders[i] !== data2.displayOrders[i]) {
                    return (data1.displayOrders[i] as number) - (data2.displayOrders[i] as number); // asc
                }
            }

            return data1.name.localeCompare(data2.name, undefined, { // asc
                numeric: true,
                sensitivity: 'base'
            });
        });
    } else if (sortingType === ChartSortingType.Name.type) {
        items.sort(function (data1, data2) {
            return data1.name.localeCompare(data2.name, undefined, { // asc
                numeric: true,
                sensitivity: 'base'
            });
        });
    } else {
        items.sort(function (data1, data2) {
            const compareValue = data2.value.compareTo(data1.value); // desc

            if (compareValue !== 0) {
                return compareValue;
            }

            for (let i = 0; i < Math.min(data1.displayOrders.length, data2.displayOrders.length); i++) {
                if (data1.displayOrders[i] !== data2.displayOrders[i]) {
                    return (data1.displayOrders[i] as number) - (data2.displayOrders[i] as number); // asc
                }
            }

            return data1.name.localeCompare(data2.name, undefined, { // asc
                numeric: true,
                sensitivity: 'base'
            });
        });
    }
}

export function isTransactionStatisticsEquals(statistics1: TransactionStatisticResponse | null, statistics2: TransactionStatisticResponse | null): boolean {
    if (statistics1 === statistics2) {
        return true;
    }

    if (!statistics1 || !statistics2 || statistics1.startTime !== statistics2.startTime || statistics1.endTime !== statistics2.endTime) {
        return false;
    }

    return isTransactionStatisticItemsEquals(statistics1.items, statistics2.items);
}

export function isTransactionStatisticsTrendsEquals(statistics1: TransactionStatisticTrendsResponseItem[], statistics2: TransactionStatisticTrendsResponseItem[]): boolean {
    if (statistics1.length !== statistics2.length) {
        return false;
    }

    for (let i = 0; i < statistics1.length; i++) {
        const item1 = statistics1[i]!;
        const item2 = statistics2[i]!;

        if (item1.year !== item2.year || item1.month !== item2.month || !isTransactionStatisticItemsEquals(item1.items, item2.items)) {
            return false;
        }
    }

    return true;
}

export function isTransactionStatisticsAssetTrendsEquals(statistics1: TransactionStatisticAssetTrendsResponseItem[], statistics2: TransactionStatisticAssetTrendsResponseItem[]): boolean {
    if (statistics1.length !== statistics2.length) {
        return false;
    }

    for (let i = 0; i < statistics1.length; i++) {
        const item1 = statistics1[i]!;
        const item2 = statistics2[i]!;

        if (item1.year !== item2.year || item1.month !== item2.month || item1.day !== item2.day) {
            return false;
        }

        if (!isArrayEquals(item1.items, item2.items,
            item => item.accountId,
            (dataItem1, dataItem2) => dataItem1.accountId === dataItem2.accountId &&
                dataItem1.accountOpeningBalance === dataItem2.accountOpeningBalance &&
                dataItem1.accountClosingBalance === dataItem2.accountClosingBalance)) {
            return false;
        }
    }

    return true;
}

export function getAllDateRangesFromItems<T extends Year1BasedMonth>(items: YearMonthItems<T>[], startYearMonth: Year1BasedMonth | TextualYearMonth | '', endYearMonth: Year1BasedMonth | TextualYearMonth | '', fiscalYearStart: number, dateAggregationType: number): YearUnixTime[] | FiscalYearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] {
    if ((!startYearMonth || !endYearMonth) && items && items.length) {
        let minYear = Number.MAX_SAFE_INTEGER, minMonth = Number.MAX_SAFE_INTEGER, maxYear = 0, maxMonth = 0;

        for (const item of items) {
            for (const dataItem of item.items) {
                if (dataItem.year < minYear || (dataItem.year === minYear && dataItem.month1base < minMonth)) {
                    minYear = dataItem.year;
                    minMonth = dataItem.month1base;
                }

                if (dataItem.year > maxYear || (dataItem.year === maxYear && dataItem.month1base > maxMonth)) {
                    maxYear = dataItem.year;
                    maxMonth = dataItem.month1base;
                }
            }
        }

        startYearMonth = `${minYear}-${minMonth}` as TextualYearMonth;
        endYearMonth = `${maxYear}-${maxMonth}` as TextualYearMonth;
    }

    return getAllDateRangesByYearMonthRange(startYearMonth, endYearMonth, fiscalYearStart, dateAggregationType);
}

export function getAllDateRangesByYearMonthRange(startYearMonth: Year1BasedMonth | TextualYearMonth | '', endYearMonth: Year1BasedMonth | TextualYearMonth | '', fiscalYearStart: number, dateAggregationType: number): YearUnixTime[] | FiscalYearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] {
    if (!startYearMonth || !endYearMonth) {
        return [];
    }

    if (dateAggregationType === ChartDateAggregationType.Year.type) {
        return getAllYearsStartAndEndUnixTimes(startYearMonth, endYearMonth);
    } else if (dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
        return getAllFiscalYearsStartAndEndUnixTimes(startYearMonth, endYearMonth, fiscalYearStart);
    } else if (dateAggregationType === ChartDateAggregationType.Quarter.type) {
        return getAllQuartersStartAndEndUnixTimes(startYearMonth, endYearMonth);
    } else { // if (dateAggregationType === ChartDateAggregationType.Month.type) {
        return getAllMonthsStartAndEndUnixTimes(startYearMonth, endYearMonth);
    }
}

export function getDateRangeKeyWithYearOffset(dateRange: YearUnixTime | FiscalYearUnixTime | YearQuarterUnixTime | YearMonthUnixTime | YearMonthDayUnixTime, dateAggregationType: number, yearOffset?: number): string | undefined {
    if (dateAggregationType === ChartDateAggregationType.Year.type) {
        return (dateRange.year + (yearOffset ?? 0)).toString();
    } else if (dateAggregationType === ChartDateAggregationType.FiscalYear.type && 'year' in dateRange) {
        return (dateRange.year + (yearOffset ?? 0)).toString();
    } else if (dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
        return `${dateRange.year + (yearOffset ?? 0)}-${dateRange.quarter}`;
    } else if ((dateAggregationType === ChartDateAggregationType.Month.type || dateAggregationType === ChartDateAggregationType.BillingCycle.type) && 'month0base' in dateRange) {
        return `${dateRange.year + (yearOffset ?? 0)}-${dateRange.month0base + 1}`;
    } else if (dateAggregationType === ChartDateAggregationType.Day.type && 'day' in dateRange) {
        return `${dateRange.year + (yearOffset ?? 0)}-${dateRange.month}-${dateRange.day}`;
    } else {
        return undefined;
    }
}
