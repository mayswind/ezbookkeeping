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
    },
    AccountTotalAssets: {
        type: 6,
        name: 'Account Total Assets'
    },
    AccountTotalLiabilities: {
        type: 7,
        name: 'Account Total Liabilities'
    }
};

const defaultChartDataType = allChartDataTypes.ExpenseByPrimaryCategory.type;

const allSortingTypes = {
    Amount: {
        type: 0,
        name: 'Amount'
    },
    DisplayOrder: {
        type: 1,
        name: 'Display Order'
    },
    Name: {
        type: 2,
        name: 'Name'
    }
};

const defaultSortingType = allSortingTypes.Amount.type;

export default {
    allChartTypes: allChartTypes,
    defaultChartType: defaultChartType,
    allChartDataTypes: allChartDataTypes,
    defaultChartDataType: defaultChartDataType,
    defaultDataRangeType: datetime.allDateRanges.ThisMonth.type,
    allSortingTypes: allSortingTypes,
    defaultSortingType: defaultSortingType,
};
