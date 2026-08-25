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

interface OverviewWidgetSettingItemBase {
    settingType: 'itemCountSelect' | 'monthSelect' | 'customSelect' | 'switch' | 'color' | 'textbox';
    settingName: string;
    displayName: string;
}

export interface OverviewWidgetItemCountSelectSettingItem extends OverviewWidgetSettingItemBase {
    settingType: 'itemCountSelect';
    itemCountValues: number[];
}

export interface OverviewWidgetMonthSelectSettingItem extends OverviewWidgetSettingItemBase {
    settingType: 'monthSelect';
    monthValues: number[];
}

export interface OverviewWidgetCustomSelectSettingItem extends OverviewWidgetSettingItemBase {
    settingType: 'customSelect';
    selectValues: GenericNameValue<string | number>[];
    multiple?: boolean;
    allValue?: string | number;
    selectValueSource?: 'accountCategories';
}

export interface OverviewWidgetSwitchSettingItem extends OverviewWidgetSettingItemBase {
    settingType: 'switch';
}

export interface OverviewWidgetColorSettingItem extends OverviewWidgetSettingItemBase {
    settingType: 'color';
}

export interface OverviewWidgetTextboxSettingItem extends OverviewWidgetSettingItemBase {
    settingType: 'textbox';
    placeholder?: string;
}

export type OverviewWidgetSettingItem = OverviewWidgetItemCountSelectSettingItem |
    OverviewWidgetMonthSelectSettingItem |
    OverviewWidgetCustomSelectSettingItem |
    OverviewWidgetSwitchSettingItem |
    OverviewWidgetColorSettingItem |
    OverviewWidgetTextboxSettingItem;

export interface OverviewWidgetDefinitionBase {
    type: OverviewWidgetType;
    name: string;
    supportsSettings: OverviewWidgetSettingItem[];
    defaultSettings: Record<string, OverviewWidgetSettingValue>;
    dataRequirements: OverviewWidgetDataRequirement[];
}

export interface OverviewLayoutBase {
    widgets: OverviewWidgetLayoutBase[];
}

export interface OverviewWidgetLayoutBase {
    id: string;
    type: OverviewWidgetType;
    settings: Record<string, OverviewWidgetSettingValue>;
}

export interface DesktopOverviewWidgetDefinition extends OverviewWidgetDefinitionBase{
    type: OverviewWidgetType;
    name: string;
    supportsSettings: OverviewWidgetSettingItem[];
    defaultWidth: number;
    defaultHeight: number;
    minWidth: number;
    minHeight: number;
    maxWidth?: number;
    maxHeight?: number;
    defaultSettings: Record<string, OverviewWidgetSettingValue>;
    dataRequirements: OverviewWidgetDataRequirement[];
}

export interface DesktopOverviewLayout extends OverviewLayoutBase {
    widgets: DesktopOverviewWidgetLayout[];
}

export interface DesktopOverviewWidgetLayout extends OverviewWidgetLayoutBase {
    id: string;
    type: OverviewWidgetType;
    x: number;
    y: number;
    w: number;
    h: number;
    settings: Record<string, OverviewWidgetSettingValue>;
}

export interface MobileOverviewWidgetDefinition extends OverviewWidgetDefinitionBase {
    type: OverviewWidgetType;
    name: string;
    supportsSettings: OverviewWidgetSettingItem[];
    defaultSettings: Record<string, OverviewWidgetSettingValue>;
    dataRequirements: OverviewWidgetDataRequirement[];
}

export interface MobileOverviewLayout extends OverviewLayoutBase {
    widgets: MobileOverviewWidgetLayout[];
}

export interface MobileOverviewWidgetLayout extends OverviewWidgetLayoutBase {
    id: string;
    type: OverviewWidgetType;
    settings: Record<string, OverviewWidgetSettingValue>;
}
