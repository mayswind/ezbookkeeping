import datetime from './datetime.js';

const allChartTypes = {
    Pie: 0,
    Bar: 1
};

const defaultChartType = allChartTypes.Pie;

const allChartDataTypes = {
    ExpenseByAccount: 0,
    ExpenseByPrimaryCategory: 1,
    ExpenseBySecondaryCategory: 2,
    IncomeByAccount: 3,
    IncomeByPrimaryCategory: 4,
    IncomeBySecondaryCategory: 5
};

const defaultChartDataType = allChartDataTypes.ExpenseByPrimaryCategory;

export default {
    allChartTypes: allChartTypes,
    defaultChartType: defaultChartType,
    allChartDataTypes: allChartDataTypes,
    defaultChartDataType: defaultChartDataType,
    defaultDataRangeType: datetime.allDateRanges.ThisMonth.type,
};
