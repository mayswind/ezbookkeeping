import type { GenericNameValue } from './base.ts';

export enum OverviewWidgetType {
    AssetSummary = 'asset-summary',
    AccountBalanceList = 'account-balance-list',
    CurrentMonthOverview = 'current-month-overview',
    CurrentMonthExpenseProgress = 'current-month-expense-progress',
    PeriodIncomeExpense = 'period-income-expense',
    PeriodNetIncomeAndSavingsRate = 'period-net-income-and-savings-rate',
    IncomeExpenseTrend = 'income-expense-trend',
    NetAssetsTrend = 'net-assets-trend',
    ExpenseCategoryRanking = 'expense-category-ranking',
    RecentTransactions = 'recent-transactions',
    TransactionCalendarHeatmap = 'transaction-calendar-heatmap'
}

export enum OverviewWidgetDataRequirement {
    Accounts = 'accounts',
    TransactionCategories = 'transactionCategories',
    TransactionOverview = 'transactionOverview',
    TransactionOverviewLast2Months = 'transactionOverviewLast2Months',
    TransactionOverviewLast12Months = 'transactionOverviewLast12Months',
    TransactionCategoryStatistics = 'transactionCategoryStatistics',
    AssetTrends = 'assetTrends',
    RecentTransactions = 'recentTransactions',
    DailyTransactionAmounts = 'dailyTransactionAmounts'
}

export type OverviewWidgetSettingValue = string | number | boolean | (string | number)[];

interface DesktopOverviewWidgetSettingBase {
    settingType: 'itemCountSelect' | 'monthSelect' | 'customSelect' | 'switch' | 'textbox';
    settingName: string;
    displayName: string;
}

export interface DesktopOverviewWidgetItemCountSelectSetting extends DesktopOverviewWidgetSettingBase {
    settingType: 'itemCountSelect';
    itemCountValues: number[];
}

export interface DesktopOverviewWidgetMonthSelectSetting extends DesktopOverviewWidgetSettingBase {
    settingType: 'monthSelect';
    monthValues: number[];
}

export interface DesktopOverviewWidgetCustomSelectSetting extends DesktopOverviewWidgetSettingBase {
    settingType: 'customSelect';
    selectValues: GenericNameValue<string | number>[];
    multiple?: boolean;
    allValue?: string | number;
    selectValueSource?: 'accountCategories';
}

export interface DesktopOverviewWidgetSwitchSetting extends DesktopOverviewWidgetSettingBase {
    settingType: 'switch';
}

export interface DesktopOverviewWidgetTextboxSetting extends DesktopOverviewWidgetSettingBase {
    settingType: 'textbox';
    placeholder?: string;
}

export type DesktopOverviewWidgetSetting = DesktopOverviewWidgetItemCountSelectSetting |
    DesktopOverviewWidgetMonthSelectSetting |
    DesktopOverviewWidgetCustomSelectSetting |
    DesktopOverviewWidgetSwitchSetting |
    DesktopOverviewWidgetTextboxSetting;

export interface DesktopOverviewWidgetDefinition {
    type: OverviewWidgetType;
    name: string;
    supportsSettings: DesktopOverviewWidgetSetting[];
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
