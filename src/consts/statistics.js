import datetime from './datetime.js';

const allAnalysisTypes = {
    CategoricalAnalysis: 0,
    TrendAnalysis: 1
};

const allCategoricalChartTypes = {
    Pie: 0,
    Bar: 1
};

const allCategoricalChartTypesArray = [
    {
        name: 'Pie Chart',
        type: allCategoricalChartTypes.Pie
    },
    {
        name: 'Bar Chart',
        type: allCategoricalChartTypes.Bar
    }
];

const defaultCategoricalChartType = allCategoricalChartTypes.Pie;

const allTrendChartTypes = {
    Area: 0,
    Column: 1
};

const allTrendChartTypesArray = [
    {
        name: 'Area Chart',
        type: allTrendChartTypes.Area
    },
    {
        name: 'Column Chart',
        type: allTrendChartTypes.Column
    }
];

const defaultTrendChartType = allTrendChartTypes.Area;

const allChartDataTypes = {
    ExpenseByAccount: {
        type: 0,
        name: 'Expense By Account',
        availableAnalysisTypes: {
            [allAnalysisTypes.CategoricalAnalysis]: true,
            [allAnalysisTypes.TrendAnalysis]: true,
        }
    },
    ExpenseByPrimaryCategory: {
        type: 1,
        name: 'Expense By Primary Category',
        availableAnalysisTypes: {
            [allAnalysisTypes.CategoricalAnalysis]: true,
            [allAnalysisTypes.TrendAnalysis]: true,
        }
    },
    ExpenseBySecondaryCategory: {
        type: 2,
        name: 'Expense By Secondary Category',
        availableAnalysisTypes: {
            [allAnalysisTypes.CategoricalAnalysis]: true,
            [allAnalysisTypes.TrendAnalysis]: true,
        }
    },
    IncomeByAccount: {
        type: 3,
        name: 'Income By Account',
        availableAnalysisTypes: {
            [allAnalysisTypes.CategoricalAnalysis]: true,
            [allAnalysisTypes.TrendAnalysis]: true,
        }
    },
    IncomeByPrimaryCategory: {
        type: 4,
        name: 'Income By Primary Category',
        availableAnalysisTypes: {
            [allAnalysisTypes.CategoricalAnalysis]: true,
            [allAnalysisTypes.TrendAnalysis]: true,
        }
    },
    IncomeBySecondaryCategory: {
        type: 5,
        name: 'Income By Secondary Category',
        availableAnalysisTypes: {
            [allAnalysisTypes.CategoricalAnalysis]: true,
            [allAnalysisTypes.TrendAnalysis]: true,
        }
    },
    AccountTotalAssets: {
        type: 6,
        name: 'Account Total Assets',
        availableAnalysisTypes: {
            [allAnalysisTypes.CategoricalAnalysis]: true
        }
    },
    AccountTotalLiabilities: {
        type: 7,
        name: 'Account Total Liabilities',
        availableAnalysisTypes: {
            [allAnalysisTypes.CategoricalAnalysis]: true
        }
    }
};

const defaultChartDataType = allChartDataTypes.ExpenseByPrimaryCategory.type;

const allSortingTypes = {
    Amount: {
        type: 0,
        name: 'Amount',
        fullName: 'Sort by Amount'
    },
    DisplayOrder: {
        type: 1,
        name: 'Display Order',
        fullName: 'Sort by Display Order'
    },
    Name: {
        type: 2,
        name: 'Name',
        fullName: 'Sort by Name'
    }
};

const allSortingTypesArray = [
    allSortingTypes.Amount,
    allSortingTypes.DisplayOrder,
    allSortingTypes.Name
]

const defaultSortingType = allSortingTypes.Amount.type;

export default {
    allAnalysisTypes: allAnalysisTypes,
    allCategoricalChartTypes: allCategoricalChartTypes,
    allCategoricalChartTypesArray: allCategoricalChartTypesArray,
    defaultCategoricalChartType: defaultCategoricalChartType,
    allTrendChartTypes: allTrendChartTypes,
    allTrendChartTypesArray: allTrendChartTypesArray,
    defaultTrendChartType: defaultTrendChartType,
    allChartDataTypes: allChartDataTypes,
    defaultChartDataType: defaultChartDataType,
    defaultCategoricalChartDataRangeType: datetime.allDateRanges.ThisMonth.type,
    defaultTrendChartDataRangeType: datetime.allDateRanges.RecentTwelveMonths.type,
    allSortingTypes: allSortingTypes,
    allSortingTypesArray: allSortingTypesArray,
    defaultSortingType: defaultSortingType,
};
