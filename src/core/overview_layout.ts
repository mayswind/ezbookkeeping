export enum OverviewWidgetType {
    CurrentMonthOverview = 'current-month-overview',
    AssetSummary = 'asset-summary',
    PeriodIncomeExpense = 'period-income-expense',
    IncomeExpenseTrend = 'income-expense-trend'
}

export enum OverviewWidgetDataRequirement {
    Accounts = 'accounts',
    TransactionCategories = 'transactionCategories',
    TransactionOverview = 'transactionOverview',
    TransactionOverviewLast12Months = 'transactionOverviewLast12Months'
}

export type OverviewPeriod = 'today' | 'thisWeek' | 'thisMonth' | 'thisYear';
export type OverviewWidgetSettingValue = string | number | boolean;

export interface DesktopOverviewWidgetDefinition {
    type: OverviewWidgetType;
    name: string;
    supportsSettings: boolean;
    defaultWidth: number;
    defaultHeight: number;
    minWidth: number;
    minHeight: number;
    maxWidth?: number;
    maxHeight?: number;
    defaultSettings: Record<string, OverviewWidgetSettingValue>;
    dataRequirements: OverviewWidgetDataRequirement[];
}

export interface DesktopOverviewLayout {
    widgets: DesktopOverviewWidgetLayout[];
}

export interface DesktopOverviewWidgetLayout {
    id: string;
    type: OverviewWidgetType;
    x: number;
    y: number;
    w: number;
    h: number;
    settings: Record<string, OverviewWidgetSettingValue>;
}
