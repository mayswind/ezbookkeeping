import type { YearMonth, YearUnixTime, YearQuarterUnixTime, YearMonthUnixTime } from '@/core/datetime.ts';
import type { FiscalYearUnixTime } from '@/core/fiscalyear.ts';
import { ChartSortingType, ChartDateAggregationType } from '@/core/statistics.ts';
import type {
    YearMonthItems,
    SortableTransactionStatisticDataItem
} from '@/models/transaction.ts';

import {
    getAllMonthsStartAndEndUnixTimes,
    getAllQuartersStartAndEndUnixTimes,
    getAllYearsStartAndEndUnixTimes,
    getAllFiscalYearsStartAndEndUnixTimes
} from '@/lib/datetime.ts';

export function sortStatisticsItems<T extends SortableTransactionStatisticDataItem>(items: T[], sortingType: number): void {
    if (sortingType === ChartSortingType.DisplayOrder.type) {
        items.sort(function (data1, data2) {
            for (let i = 0; i < Math.min(data1.displayOrders.length, data2.displayOrders.length); i++) {
                if (data1.displayOrders[i] !== data2.displayOrders[i]) {
                    return data1.displayOrders[i] - data2.displayOrders[i]; // asc
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
            if (data1.totalAmount !== data2.totalAmount) {
                return data2.totalAmount - data1.totalAmount; // desc
            }

            return data1.name.localeCompare(data2.name, undefined, { // asc
                numeric: true,
                sensitivity: 'base'
            });
        });
    }
}

export function getAllDateRanges<T extends YearMonth>(items: YearMonthItems<T>[], startYearMonth: YearMonth | string, endYearMonth: YearMonth | string, fiscalYearStart: number, dateAggregationType: number): YearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] | FiscalYearUnixTime[] {
    if ((!startYearMonth || !endYearMonth) && items && items.length) {
        let minYear = Number.MAX_SAFE_INTEGER, minMonth = Number.MAX_SAFE_INTEGER, maxYear = 0, maxMonth = 0;

        for (let i = 0; i < items.length; i++) {
            const item = items[i];

            for (let j = 0; j < item.items.length; j++) {
                const dataItem = item.items[j];

                if (dataItem.year < minYear || (dataItem.year === minYear && dataItem.month < minMonth)) {
                    minYear = dataItem.year;
                    minMonth = dataItem.month;
                }

                if (dataItem.year > maxYear || (dataItem.year === maxYear && dataItem.month > maxMonth)) {
                    maxYear = dataItem.year;
                    maxMonth = dataItem.month;
                }
            }
        }

        startYearMonth = `${minYear}-${minMonth}`;
        endYearMonth = `${maxYear}-${maxMonth}`;
    }

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
