import { DateRange } from '@/core/datetime.ts';

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

const defaultTrendChartType = allTrendChartTypes.Column;

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
    },
    TotalExpense: {
        type: 8,
        name: 'Total Expense',
        availableAnalysisTypes: {
            [allAnalysisTypes.TrendAnalysis]: true
        }
    },
    TotalIncome: {
        type: 9,
        name: 'Total Income',
        availableAnalysisTypes: {
            [allAnalysisTypes.TrendAnalysis]: true
        }
    },
    TotalBalance: {
        type: 10,
        name: 'Total Balance',
        availableAnalysisTypes: {
            [allAnalysisTypes.TrendAnalysis]: true
        }
    }
};

const allChartDataTypesMap = {
    [allChartDataTypes.ExpenseByAccount.type]: allChartDataTypes.ExpenseByAccount,
    [allChartDataTypes.ExpenseByPrimaryCategory.type]: allChartDataTypes.ExpenseByPrimaryCategory,
    [allChartDataTypes.ExpenseBySecondaryCategory.type]: allChartDataTypes.ExpenseBySecondaryCategory,
    [allChartDataTypes.IncomeByAccount.type]: allChartDataTypes.IncomeByAccount,
    [allChartDataTypes.IncomeByPrimaryCategory.type]: allChartDataTypes.IncomeByPrimaryCategory,
    [allChartDataTypes.IncomeBySecondaryCategory.type]: allChartDataTypes.IncomeBySecondaryCategory,
    [allChartDataTypes.AccountTotalAssets.type]: allChartDataTypes.AccountTotalAssets,
    [allChartDataTypes.AccountTotalLiabilities.type]: allChartDataTypes.AccountTotalLiabilities,
    [allChartDataTypes.TotalExpense.type]: allChartDataTypes.TotalExpense,
    [allChartDataTypes.TotalIncome.type]: allChartDataTypes.TotalIncome,
    [allChartDataTypes.TotalBalance.type]: allChartDataTypes.TotalBalance
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

const allDateAggregationTypes = {
    Month: {
        type: 0,
        name: 'Aggregate by Month'
    },
    Quarter: {
        type: 1,
        name: 'Aggregate by Quarter'
    },
    Year: {
        type: 2,
        name: 'Aggregate by Year'
    }
};

const allDateAggregationTypesArray = [
    allDateAggregationTypes.Month,
    allDateAggregationTypes.Quarter,
    allDateAggregationTypes.Year
]

const defaultDateAggregationType = allDateAggregationTypes.Month.type;

export default {
    allAnalysisTypes: allAnalysisTypes,
    allCategoricalChartTypes: allCategoricalChartTypes,
    allCategoricalChartTypesArray: allCategoricalChartTypesArray,
    defaultCategoricalChartType: defaultCategoricalChartType,
    allTrendChartTypes: allTrendChartTypes,
    allTrendChartTypesArray: allTrendChartTypesArray,
    defaultTrendChartType: defaultTrendChartType,
    allChartDataTypes: allChartDataTypes,
    allChartDataTypesMap: allChartDataTypesMap,
    defaultChartDataType: defaultChartDataType,
    defaultCategoricalChartDataRangeType: DateRange.ThisMonth.type,
    defaultTrendChartDataRangeType: DateRange.ThisYear.type,
    allSortingTypes: allSortingTypes,
    allSortingTypesArray: allSortingTypesArray,
    defaultSortingType: defaultSortingType,
    allDateAggregationTypes: allDateAggregationTypes,
    allDateAggregationTypesArray: allDateAggregationTypesArray,
    defaultDateAggregationType: defaultDateAggregationType,
};
