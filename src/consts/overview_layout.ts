import {
    type DesktopOverviewLayout,
    type DesktopOverviewWidgetDefinition,
    OverviewWidgetType,
    OverviewWidgetDataRequirement
} from '@/core/overview_layout.ts';

export const DESKTOP_OVERVIEW_LAYOUT_COLUMNS: number = 12;
export const DESKTOP_OVERVIEW_LAYOUT_MAX_WIDGETS: number = 100;
export const DESKTOP_OVERVIEW_LAYOUT_MAX_ROWS: number = 1000;

export const DESKTOP_OVERVIEW_WIDGET_DEFINITIONS: Record<OverviewWidgetType, DesktopOverviewWidgetDefinition> = {
    [OverviewWidgetType.CurrentMonthOverview]: {
        type: OverviewWidgetType.CurrentMonthOverview,
        name: 'Monthly Expense Overview',
        supportsSettings: false,
        defaultWidth: 4,
        defaultHeight: 3,
        minWidth: 3,
        minHeight: 3,
        defaultSettings: {},
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverview
        ]
    },
    [OverviewWidgetType.AssetSummary]: {
        type: OverviewWidgetType.AssetSummary,
        name: 'Asset Summary',
        supportsSettings: false,
        defaultWidth: 8,
        defaultHeight: 3,
        minWidth: 3,
        minHeight: 3,
        defaultSettings: {},
        dataRequirements: [
            OverviewWidgetDataRequirement.Accounts
        ]
    },
    [OverviewWidgetType.PeriodIncomeExpense]: {
        type: OverviewWidgetType.PeriodIncomeExpense,
        name: 'Period Income and Expense',
        supportsSettings: true,
        defaultWidth: 3,
        defaultHeight: 3,
        minWidth: 2,
        minHeight: 3,
        defaultSettings: {
            dateRange: 'today'
        },
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverview
        ]
    },
    [OverviewWidgetType.IncomeExpenseTrend]: {
        type: OverviewWidgetType.IncomeExpenseTrend,
        name: 'Income and Expense Trends',
        supportsSettings: true,
        defaultWidth: 6,
        defaultHeight: 6,
        minWidth: 3,
        minHeight: 6,
        defaultSettings: {
            months: 12
        },
        dataRequirements: [
            OverviewWidgetDataRequirement.TransactionOverviewLast12Months
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
            id: 'default-period-today',
            type: OverviewWidgetType.PeriodIncomeExpense,
            x: 0,
            y: 3,
            w: 3,
            h: 3,
            settings: {
                dateRange: 'today'
            }
        },
        {
            id: 'default-period-week',
            type: OverviewWidgetType.PeriodIncomeExpense,
            x: 3,
            y: 3,
            w: 3,
            h: 3,
            settings: {
                dateRange: 'thisWeek'
            }
        },
        {
            id: 'default-period-month',
            type: OverviewWidgetType.PeriodIncomeExpense,
            x: 0,
            y: 6,
            w: 3,
            h: 3,
            settings: {
                dateRange: 'thisMonth'
            }
        },
        {
            id: 'default-period-year',
            type: OverviewWidgetType.PeriodIncomeExpense,
            x: 3,
            y: 6,
            w: 3,
            h: 3,
            settings: {
                dateRange: 'thisYear'
            }
        },
        {
            id: 'default-income-expense-trend',
            type: OverviewWidgetType.IncomeExpenseTrend,
            x: 6,
            y: 3,
            w: 6,
            h: 6,
            settings: {months: 12}
        }
    ]
};
