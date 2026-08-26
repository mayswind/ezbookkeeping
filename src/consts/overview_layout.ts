import type { PartialRecord } from '@/core/base.ts';
import { DateRange } from '@/core/datetime.ts';
import { AccountCategory } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    type OverviewWidgetColorSettingItem,
    type OverviewWidgetTextboxSettingItem,
    type DesktopOverviewLayout,
    type DesktopOverviewWidgetDefinition,
    type MobileOverviewLayout,
    type MobileOverviewWidgetDefinition,
    OverviewWidgetType,
    OverviewWidgetDataRequirement
} from '@/core/overview_layout.ts';

import {
    DEFAULT_MOBILE_OVERVIEW_WIDGET_DARK_BACKGROUND_COLOR,
    DEFAULT_MOBILE_OVERVIEW_WIDGET_LIGHT_BACKGROUND_COLOR
} from '@/consts/color.ts';

export const DESKTOP_OVERVIEW_LAYOUT_COLUMNS: number = 12;
export const DESKTOP_OVERVIEW_LAYOUT_MAX_WIDGETS: number = 100;
export const DESKTOP_OVERVIEW_LAYOUT_MAX_ROWS: number = 1000;
export const MOBILE_OVERVIEW_LAYOUT_MAX_WIDGETS: number = 100;

const WIDGET_TITLE_SETTING: OverviewWidgetTextboxSettingItem = {
    settingType: 'textbox',
    settingName: 'title',
    displayName: 'Widget Title',
    placeholder: 'Widget Title'
};

const WIDGET_BACKGROUND_COLOR_SETTINGS: OverviewWidgetColorSettingItem[] = [
    {
        settingType: 'color',
        settingName: 'lightBackgroundColor',
        displayName: 'Light Mode Background Color'
    },
    {
        settingType: 'color',
        settingName: 'darkBackgroundColor',
        displayName: 'Dark Mode Background Color'
    }
];

export const DESKTOP_OVERVIEW_WIDGET_DEFINITIONS: PartialRecord<OverviewWidgetType, DesktopOverviewWidgetDefinition> = {
    [OverviewWidgetType.AssetSummary]: {
        type: OverviewWidgetType.AssetSummary,
        name: 'Asset Summary',
        supportsSettings: [
            WIDGET_TITLE_SETTING
        ],
        defaultSettings: {},
        defaultWidth: 8,
        defaultHeight: 3,
        minWidth: 3,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.Accounts
        ]
    },
    [OverviewWidgetType.AccountBalanceList]: {
        type: OverviewWidgetType.AccountBalanceList,
        name: 'Account Balance List',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'customSelect',
                settingName: 'accountCategories',
                displayName: 'Account Category',
                selectValues: [
                    { name: 'All', value: 0 },
                    ...AccountCategory.values().map(category => ({
                        name: category.name,
                        value: category.type
                    }))
                ],
                multiple: true,
                allValue: 0,
                selectValueSource: 'accountCategories'
            },
            {
                settingType: 'itemCountSelect',
                settingName: 'itemCount',
                displayName: 'Item Count',
                itemCountValues: [3, 4, 5, 6, 7, 8, 9, 10]
            },
            {
                settingType: 'customSelect',
                settingName: 'sortBy',
                displayName: 'Sort By',
                selectValues: [
                    {
                        name: 'Display Order',
                        value: 'displayOrder'
                    },
                    {
                        name: 'Balance',
                        value: 'balance'
                    }
                ]
            }
        ],
        defaultSettings: {
            accountCategories: [0],
            itemCount: 4,
            sortBy: 'displayOrder',
        },
        defaultWidth: 3,
        defaultHeight: 3,
        minWidth: 2,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.Accounts
        ]
    },
    [OverviewWidgetType.CurrentMonthOverview]: {
        type: OverviewWidgetType.CurrentMonthOverview,
        name: 'This Month\'s Income and Expense Overview',
        supportsSettings: [],
        defaultSettings: {},
        defaultWidth: 4,
        defaultHeight: 3,
        minWidth: 3,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverview
        ]
    },
    [OverviewWidgetType.CurrentMonthExpenseProgress]: {
        type: OverviewWidgetType.CurrentMonthExpenseProgress,
        name: 'This Month\'s Expense Progress',
        supportsSettings: [
            WIDGET_TITLE_SETTING
        ],
        defaultSettings: {},
        defaultWidth: 3,
        defaultHeight: 3,
        minWidth: 2,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverviewLast2Months
        ]
    },
    [OverviewWidgetType.PeriodIncomeExpense]: {
        type: OverviewWidgetType.PeriodIncomeExpense,
        name: 'Period Income and Expense',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'customSelect',
                settingName: 'dateRange',
                displayName: 'Date Range',
                selectValues: [
                    DateRange.Today,
                    DateRange.ThisWeek,
                    DateRange.ThisMonth,
                    DateRange.ThisYear
                ].map(dateRange => ({
                    name: dateRange.name,
                    value: dateRange.type
                }))
            }
        ],
        defaultSettings: {
            dateRange: DateRange.Today.type
        },
        defaultWidth: 3,
        defaultHeight: 3,
        minWidth: 2,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverview
        ]
    },
    [OverviewWidgetType.PeriodNetIncomeAndSavingsRate]: {
        type: OverviewWidgetType.PeriodNetIncomeAndSavingsRate,
        name: 'Period Net Income and Savings Rate',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'customSelect',
                settingName: 'dateRange',
                displayName: 'Date Range',
                selectValues: [
                    DateRange.Today,
                    DateRange.ThisWeek,
                    DateRange.ThisMonth,
                    DateRange.ThisYear
                ].map(dateRange => ({
                    name: dateRange.name,
                    value: dateRange.type
                }))
            }
        ],
        defaultSettings: {
            dateRange: DateRange.ThisMonth.type
        },
        defaultWidth: 3,
        defaultHeight: 3,
        minWidth: 2,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverview
        ]
    },
    [OverviewWidgetType.IncomeExpenseTrend]: {
        type: OverviewWidgetType.IncomeExpenseTrend,
        name: 'Income and Expense Trends',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'monthSelect',
                settingName: 'months',
                displayName: 'Date Range',
                monthValues: [6, 12]
            },
            {
                settingType: 'switch',
                settingName: 'showXAxisLabels',
                displayName: 'Show Horizontal Axis Labels'
            },
            {
                settingType: 'switch',
                settingName: 'showLegend',
                displayName: 'Show Legend'
            }
        ],
        defaultSettings: {
            months: 12,
            showXAxisLabels: true,
            showLegend: true
        },
        defaultWidth: 6,
        defaultHeight: 6,
        minWidth: 3,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverviewLast12Months
        ]
    },
    [OverviewWidgetType.NetAssetsTrend]: {
        type: OverviewWidgetType.NetAssetsTrend,
        name: 'Net Assets Trends',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'monthSelect',
                settingName: 'months',
                displayName: 'Date Range',
                monthValues: [6, 12]
            },
            {
                settingType: 'switch',
                settingName: 'showXAxisLabels',
                displayName: 'Show Horizontal Axis Labels'
            },
            {
                settingType: 'switch',
                settingName: 'showLegend',
                displayName: 'Show Legend'
            }
        ],
        defaultSettings: {
            months: 12,
            showXAxisLabels: true,
            showLegend: true
        },
        defaultWidth: 6,
        defaultHeight: 6,
        minWidth: 3,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.Accounts,
            OverviewWidgetDataRequirement.AssetTrends
        ]
    },
    [OverviewWidgetType.ExpenseCategoryRanking]: {
        type: OverviewWidgetType.ExpenseCategoryRanking,
        name: 'Expense Category Ranking',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'customSelect',
                settingName: 'dateRange',
                displayName: 'Date Range',
                selectValues: [
                    DateRange.ThisMonth,
                    DateRange.ThisYear
                ].map(dateRange => ({
                    name: dateRange.name,
                    value: dateRange.type
                }))
            },
            {
                settingType: 'customSelect',
                settingName: 'categoryLevel',
                displayName: 'Category Level',
                selectValues: [
                    {
                        name: 'Primary Category',
                        value: 'primary'
                    },
                    {
                        name: 'Secondary Category',
                        value: 'secondary'
                    }
                ]
            },
            {
                settingType: 'itemCountSelect',
                settingName: 'itemCount',
                displayName: 'Item Count',
                itemCountValues: [3, 4, 5, 6, 7, 8, 9, 10]
            }
        ],
        defaultSettings: {
            dateRange: DateRange.ThisMonth.type,
            categoryLevel: 'primary',
            itemCount: 3
        },
        defaultWidth: 3,
        defaultHeight: 3,
        minWidth: 2,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.Accounts,
            OverviewWidgetDataRequirement.TransactionCategories,
            OverviewWidgetDataRequirement.TransactionCategoryStatistics
        ]
    },
    [OverviewWidgetType.RecentTransactions]: {
        type: OverviewWidgetType.RecentTransactions,
        name: 'Recent Transactions',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'itemCountSelect',
                settingName: 'itemCount',
                displayName: 'Item Count',
                itemCountValues: [3, 4, 5, 6, 7, 8, 9, 10]
            }
        ],
        defaultSettings: {
            itemCount: 3
        },
        defaultWidth: 3,
        defaultHeight: 3,
        minWidth: 2,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.RecentTransactions
        ]
    },
    [OverviewWidgetType.TransactionCalendarHeatmap]: {
        type: OverviewWidgetType.TransactionCalendarHeatmap,
        name: 'Transaction Calendar Heatmap',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'customSelect',
                settingName: 'transactionType',
                displayName: 'Transaction Type',
                selectValues: [
                    {
                        name: 'Income',
                        value: TransactionType.Income
                    },
                    {
                        name: 'Expense',
                        value: TransactionType.Expense
                    }
                ]
            },
            {
                settingType: 'monthSelect',
                settingName: 'months',
                displayName: 'Date Range',
                monthValues: [6, 12]
            }
        ],
        defaultSettings: {
            transactionType: TransactionType.Expense,
            months: 12
        },
        defaultWidth: 6,
        defaultHeight: 3,
        minWidth: 3,
        minHeight: 3,
        dataRequirements: [
            OverviewWidgetDataRequirement.DailyTransactionAmounts
        ]
    }
};

export const DEFAULT_DESKTOP_OVERVIEW_LAYOUT: DesktopOverviewLayout = {
    widgets: [
        {
            id: 'default-current-month-overview',
            type: OverviewWidgetType.CurrentMonthOverview,
            x: 0,
            y: 0,
            w: 4,
            h: 3,
            settings: {}
        },
        {
            id: 'default-asset-summary',
            type: OverviewWidgetType.AssetSummary,
            x: 4,
            y: 0,
            w: 8,
            h: 3,
            settings: {}
        },
        {
            id: 'default-today-income-expense',
            type: OverviewWidgetType.PeriodIncomeExpense,
            x: 0,
            y: 3,
            w: 3,
            h: 3,
            settings: {
                dateRange: DateRange.Today.type
            }
        },
        {
            id: 'default-week-income-expense',
            type: OverviewWidgetType.PeriodIncomeExpense,
            x: 3,
            y: 3,
            w: 3,
            h: 3,
            settings: {
                dateRange: DateRange.ThisWeek.type
            }
        },
        {
            id: 'default-month-income-expense',
            type: OverviewWidgetType.PeriodIncomeExpense,
            x: 0,
            y: 6,
            w: 3,
            h: 3,
            settings: {
                dateRange: DateRange.ThisMonth.type
            }
        },
        {
            id: 'default-year-income-expense',
            type: OverviewWidgetType.PeriodIncomeExpense,
            x: 3,
            y: 6,
            w: 3,
            h: 3,
            settings: {
                dateRange: DateRange.ThisYear.type
            }
        },
        {
            id: 'default-income-expense-trend',
            type: OverviewWidgetType.IncomeExpenseTrend,
            x: 6,
            y: 3,
            w: 6,
            h: 6,
            settings: {
                months: 12,
                showXAxisLabels: true,
                showLegend: true
            }
        }
    ]
};

export const MOBILE_OVERVIEW_WIDGET_DEFINITIONS: PartialRecord<OverviewWidgetType, MobileOverviewWidgetDefinition> = {
    [OverviewWidgetType.AssetSummary]: {
        type: OverviewWidgetType.AssetSummary,
        name: 'Asset Summary',
        supportsSettings: [
            {
                settingType: 'customSelect',
                settingName: 'height',
                displayName: 'Widget Height',
                selectValues: [
                    { name: 'Small', value: 1 },
                    { name: 'Medium', value: 2 },
                    { name: 'Large', value: 3 }
                ]
            },
            ...WIDGET_BACKGROUND_COLOR_SETTINGS
        ],
        defaultSettings: {
            height: 3,
            lightBackgroundColor: DEFAULT_MOBILE_OVERVIEW_WIDGET_LIGHT_BACKGROUND_COLOR,
            darkBackgroundColor: DEFAULT_MOBILE_OVERVIEW_WIDGET_DARK_BACKGROUND_COLOR
        },
        dataRequirements: [
            OverviewWidgetDataRequirement.Accounts
        ]
    },
    [OverviewWidgetType.AccountBalanceList]: {
        type: OverviewWidgetType.AccountBalanceList,
        name: 'Account Balance List',
        supportsSettings: [
            WIDGET_TITLE_SETTING,
            {
                settingType: 'itemCountSelect',
                settingName: 'itemCount',
                displayName: 'Item Count',
                itemCountValues: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
            },
            {
                settingType: 'customSelect',
                settingName: 'sortBy',
                displayName: 'Sort By',
                selectValues: [
                    {
                        name: 'Display Order',
                        value: 'displayOrder'
                    },
                    {
                        name: 'Balance',
                        value: 'balance'
                    }
                ]
            },
            {
                settingType: 'customSelect',
                settingName: 'accountCategories',
                displayName: 'Account Category',
                selectValues: [
                    { name: 'All', value: 0 },
                    ...AccountCategory.values().map(category => ({
                        name: category.name,
                        value: category.type
                    }))
                ],
                multiple: true,
                allValue: 0,
                selectValueSource: 'accountCategories'
            }
        ],
        defaultSettings: {
            accountCategories: [0],
            itemCount: 4,
            sortBy: 'displayOrder'
        },
        dataRequirements: [
            OverviewWidgetDataRequirement.Accounts
        ]
    },
    [OverviewWidgetType.CurrentMonthOverview]: {
        type: OverviewWidgetType.CurrentMonthOverview,
        name: 'This Month\'s Income and Expense Overview',
        supportsSettings: [
            {
                settingType: 'customSelect',
                settingName: 'height',
                displayName: 'Widget Height',
                selectValues: [
                    { name: 'Small', value: 1 },
                    { name: 'Medium', value: 2 },
                    { name: 'Large', value: 3 }
                ]
            },
            ...WIDGET_BACKGROUND_COLOR_SETTINGS
        ],
        defaultSettings: {
            height: 3,
            lightBackgroundColor: DEFAULT_MOBILE_OVERVIEW_WIDGET_LIGHT_BACKGROUND_COLOR,
            darkBackgroundColor: DEFAULT_MOBILE_OVERVIEW_WIDGET_DARK_BACKGROUND_COLOR
        },
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverview
        ]
    },
    [OverviewWidgetType.PeriodIncomeExpense]: {
        type: OverviewWidgetType.PeriodIncomeExpense,
        name: 'Period Income and Expense',
        supportsSettings: [
            {
                settingType: 'customSelect',
                settingName: 'dateRanges',
                displayName: 'Date Range',
                selectValues: [
                    DateRange.Today,
                    DateRange.ThisWeek,
                    DateRange.ThisMonth,
                    DateRange.ThisYear
                ].map(dateRange => ({
                    name: dateRange.name,
                    value: dateRange.type
                })),
                multiple: true
            }
        ],
        defaultSettings: {
            dateRanges: [
                DateRange.Today.type,
                DateRange.ThisWeek.type,
                DateRange.ThisMonth.type,
                DateRange.ThisYear.type
            ]
        },
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverview
        ]
    }
};

export const DEFAULT_MOBILE_OVERVIEW_LAYOUT: MobileOverviewLayout = {
    widgets: [
        {
            id: 'default-current-month-overview',
            type: OverviewWidgetType.CurrentMonthOverview,
            settings: {
                height: 3
            }
        },
        {
            id: 'default-period-income-expense',
            type: OverviewWidgetType.PeriodIncomeExpense,
            settings: {
                dateRanges: [
                    DateRange.Today.type,
                    DateRange.ThisWeek.type,
                    DateRange.ThisMonth.type,
                    DateRange.ThisYear.type
                ]
            }
        }
    ]
};
