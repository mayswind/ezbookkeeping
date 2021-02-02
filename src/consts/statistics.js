import datetime from './datetime.js';

const allChartTypes = {
    Pie: 0,
    Bar: 1
};

const defaultChartType = allChartTypes.Pie;

const allChartDataTypes = {
    ExpenseByAccount: {
        type: 0,
        name: 'Expense By Account'
    },
    ExpenseByPrimaryCategory: {
        type: 1,
        name: 'Expense By Primary Category'
    },
    ExpenseBySecondaryCategory: {
        type: 2,
        name: 'Expense By Secondary Category'
    },
    IncomeByAccount: {
        type: 3,
        name: 'Income By Account'
    },
    IncomeByPrimaryCategory: {
        type: 4,
        name: 'Income By Primary Category'
    },
    IncomeBySecondaryCategory: {
        type: 5,
        name: 'Income By Secondary Category'
    }
};

const defaultChartDataType = allChartDataTypes.ExpenseByPrimaryCategory;

export default {
    allChartTypes: allChartTypes,
    defaultChartType: defaultChartType,
    allChartDataTypes: allChartDataTypes,
    defaultChartDataType: defaultChartDataType,
    defaultDataRangeType: datetime.allDateRanges.ThisMonth.type,
};
