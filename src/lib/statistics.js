import statisticsConstants from '@/consts/statistics.js';

import {
    getAllMonthsStartAndEndUnixTimes,
    getAllQuartersStartAndEndUnixTimes,
    getAllYearsStartAndEndUnixTimes
} from '@/lib/datetime.js';

export function isChartDataTypeAvailableForAnalysisType(chartDataType, analysisType) {
    for (const dataTypeField in statisticsConstants.allChartDataTypes) {
        if (!Object.prototype.hasOwnProperty.call(statisticsConstants.allChartDataTypes, dataTypeField)) {
            continue;
        }

        const dataTypeItem = statisticsConstants.allChartDataTypes[dataTypeField];

        if (dataTypeItem.type !== chartDataType) {
            continue;
        }

        return !!dataTypeItem.availableAnalysisTypes[analysisType];
    }

    return false;
}

export function sortStatisticsItems(items, sortingType) {
    if (sortingType === statisticsConstants.allSortingTypes.DisplayOrder.type) {
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
    } else if (sortingType === statisticsConstants.allSortingTypes.Name.type) {
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

export function getAllDateRanges(items, startYearMonth, endYearMonth, dateAggregationType) {
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
    if (dateAggregationType === statisticsConstants.allDateAggregationTypes.Year.type) {
        return getAllYearsStartAndEndUnixTimes(startYearMonth, endYearMonth);
    } else if (dateAggregationType === statisticsConstants.allDateAggregationTypes.Quarter.type) {
        return getAllQuartersStartAndEndUnixTimes(startYearMonth, endYearMonth);
    } else { // if (dateAggregationType === statisticsConstants.allDateAggregationTypes.Month.type) {
        return getAllMonthsStartAndEndUnixTimes(startYearMonth, endYearMonth);
    }
}
