import statisticsConstants from '@/consts/statistics.js';

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
