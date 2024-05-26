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
