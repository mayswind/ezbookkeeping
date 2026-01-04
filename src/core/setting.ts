import { type WeekDayValue, WeekDay } from './datetime.ts';
import { TimezoneTypeForStatistics } from './timezone.ts';
import { CurrencySortingType } from './currency.ts';
import {
    CategoricalChartType,
    TrendChartType,
    ChartDataType,
    ChartSortingType,
    DEFAULT_CATEGORICAL_CHART_DATA_RANGE,
    DEFAULT_TREND_CHART_DATA_RANGE,
    DEFAULT_ASSET_TRENDS_CHART_DATA_RANGE
} from './statistics.ts';
import { DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE } from './explorer.ts';
import { DEFAULT_CURRENCY_CODE } from '@/consts/currency.ts';

export type ApplicationSettingKey = string;
export type ApplicationSettingValue = string | number | boolean | Record<string, ApplicationSettingSubValue>;
export type ApplicationSettingSubValue = string | number | boolean | Record<string, boolean> | Record<string, number>;

export interface BaseApplicationSetting {
    [key: ApplicationSettingKey]: ApplicationSettingValue;
}

export interface ApplicationSettings extends BaseApplicationSetting {
    // Debug Settings
    debug: boolean;
    // Basic Settings
    theme: string;
    fontSize: number;
    timeZone: string;
    autoUpdateExchangeRatesData: boolean;
    showAccountBalance: boolean;
    swipeBack: boolean;
    animate: boolean;
    // Application Lock
    applicationLock: boolean;
    applicationLockWebAuthn: boolean;
    // Navigation Bar
    showAddTransactionButtonInDesktopNavbar: boolean;
    // Overview Page
    showAmountInHomePage: boolean;
    timezoneUsedForStatisticsInHomePage: number;
    overviewAccountFilterInHomePage: Record<string, boolean>;
    overviewTransactionCategoryFilterInHomePage: Record<string, boolean>;
    // Transaction List Page
    itemsCountInTransactionListPage: number;
    showTotalAmountInTransactionListPage: boolean;
    showTagInTransactionListPage: boolean;
    // Transaction Edit Page
    autoSaveTransactionDraft: string;
    autoGetCurrentGeoLocation: boolean;
    alwaysShowTransactionPicturesInMobileTransactionEditPage: boolean;
    // Insights Explorer Page
    insightsExplorerDefaultDateRangeType: number;
    showTagInInsightsExplorerPage: boolean;
    // Account List Page
    totalAmountExcludeAccountIds: Record<string, boolean>;
    accountCategoryOrders: string;
    hideCategoriesWithoutAccounts: boolean;
    // Exchange Rates Data Page
    currencySortByInExchangeRatesPage: number;
    // Statistics Settings
    statistics: {
        defaultChartDataType: number;
        defaultTimezoneType: number;
        defaultAccountFilter: Record<string, boolean>;
        defaultTransactionCategoryFilter: Record<string, boolean>;
        defaultSortingType: number;
        defaultCategoricalChartType: number;
        defaultCategoricalChartDataRangeType: number;
        defaultTrendChartType: number;
        defaultTrendChartDataRangeType: number;
        defaultAssetTrendsChartType: number;
        defaultAssetTrendsChartDataRangeType: number;
    };
}

export enum UserApplicationCloudSettingType {
    String = 'string',
    Number = 'number',
    Boolean = 'boolean',
    StringBooleanMap = 'string_boolean_map',
}

export interface ApplicationCloudSetting {
    readonly settingKey: string;
    readonly settingValue: string;
}

export interface LocaleDefaultSettings {
    currency: string;
    firstDayOfWeek: WeekDayValue;
}

export interface ApplicationLockState {
    readonly username: string;
    readonly secret: string;
}

export interface WebAuthnConfig {
    readonly credentialId: string;
}

export const ALL_ALLOWED_CLOUD_SYNC_APP_SETTING_KEY_TYPES: Record<string, UserApplicationCloudSettingType> = {
    // Basic Settings
    'showAccountBalance': UserApplicationCloudSettingType.Boolean,
    // Overview Page
    'showAmountInHomePage': UserApplicationCloudSettingType.Boolean,
    'timezoneUsedForStatisticsInHomePage': UserApplicationCloudSettingType.Number,
    'overviewAccountFilterInHomePage': UserApplicationCloudSettingType.StringBooleanMap,
    'overviewTransactionCategoryFilterInHomePage': UserApplicationCloudSettingType.StringBooleanMap,
    // Transaction List Page
    'itemsCountInTransactionListPage': UserApplicationCloudSettingType.Number,
    'showTotalAmountInTransactionListPage': UserApplicationCloudSettingType.Boolean,
    'showTagInTransactionListPage': UserApplicationCloudSettingType.Boolean,
    // Transaction Edit Page
    'autoSaveTransactionDraft': UserApplicationCloudSettingType.String,
    'autoGetCurrentGeoLocation': UserApplicationCloudSettingType.Boolean,
    'alwaysShowTransactionPicturesInMobileTransactionEditPage': UserApplicationCloudSettingType.Boolean,
    // Insights Explorer Page
    'insightsExplorerDefaultDateRangeType': UserApplicationCloudSettingType.Number,
    'showTagInInsightsExplorerPage': UserApplicationCloudSettingType.Boolean,
    // Account List Page
    'totalAmountExcludeAccountIds': UserApplicationCloudSettingType.StringBooleanMap,
    'accountCategoryOrders': UserApplicationCloudSettingType.String,
    'hideCategoriesWithoutAccounts': UserApplicationCloudSettingType.Boolean,
    // Exchange Rates Data Page
    'currencySortByInExchangeRatesPage': UserApplicationCloudSettingType.Number,
    // Statistics Settings
    'statistics.defaultChartDataType': UserApplicationCloudSettingType.Number,
    'statistics.defaultTimezoneType': UserApplicationCloudSettingType.Number,
    'statistics.defaultAccountFilter': UserApplicationCloudSettingType.StringBooleanMap,
    'statistics.defaultTransactionCategoryFilter': UserApplicationCloudSettingType.StringBooleanMap,
    'statistics.defaultSortingType': UserApplicationCloudSettingType.Number,
    'statistics.defaultCategoricalChartType': UserApplicationCloudSettingType.Number,
    'statistics.defaultCategoricalChartDataRangeType': UserApplicationCloudSettingType.Number,
    'statistics.defaultTrendChartType': UserApplicationCloudSettingType.Number,
    'statistics.defaultTrendChartDataRangeType': UserApplicationCloudSettingType.Number,
    'statistics.defaultAssetTrendsChartType': UserApplicationCloudSettingType.Number,
    'statistics.defaultAssetTrendsChartDataRangeType': UserApplicationCloudSettingType.Number,
};

export const DEFAULT_APPLICATION_SETTINGS: ApplicationSettings = {
    // Debug Settings
    debug: false,
    // Basic Settings
    theme: 'auto',
    fontSize: 1,
    timeZone: '',
    autoUpdateExchangeRatesData: true,
    showAccountBalance: true,
    swipeBack: true,
    animate: true,
    // Application Lock
    applicationLock: false,
    applicationLockWebAuthn: false,
    // Navigation Bar
    showAddTransactionButtonInDesktopNavbar: true,
    // Overview Page
    showAmountInHomePage: true,
    timezoneUsedForStatisticsInHomePage: TimezoneTypeForStatistics.Default.type,
    overviewAccountFilterInHomePage: {},
    overviewTransactionCategoryFilterInHomePage: {},
    // Transaction List Page
    itemsCountInTransactionListPage: 15,
    showTotalAmountInTransactionListPage: true,
    showTagInTransactionListPage: true,
    // Transaction Edit Page
    autoSaveTransactionDraft: 'disabled',
    autoGetCurrentGeoLocation: false,
    alwaysShowTransactionPicturesInMobileTransactionEditPage: false,
    // Insights Explorer Page
    insightsExplorerDefaultDateRangeType: DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type,
    showTagInInsightsExplorerPage: true,
    // Account List Page
    totalAmountExcludeAccountIds: {},
    accountCategoryOrders: '',
    hideCategoriesWithoutAccounts: false,
    // Exchange Rates Data Page
    currencySortByInExchangeRatesPage: CurrencySortingType.Default.type,
    // Statistics Settings
    statistics: {
        defaultChartDataType: ChartDataType.Default.type,
        defaultTimezoneType: TimezoneTypeForStatistics.Default.type,
        defaultAccountFilter: {},
        defaultTransactionCategoryFilter: {},
        defaultSortingType: ChartSortingType.Default.type,
        defaultCategoricalChartType: CategoricalChartType.Default.type,
        defaultCategoricalChartDataRangeType: DEFAULT_CATEGORICAL_CHART_DATA_RANGE.type,
        defaultTrendChartType: TrendChartType.Default.type,
        defaultTrendChartDataRangeType: DEFAULT_TREND_CHART_DATA_RANGE.type,
        defaultAssetTrendsChartType: TrendChartType.Default.type,
        defaultAssetTrendsChartDataRangeType: DEFAULT_ASSET_TRENDS_CHART_DATA_RANGE.type,
    }
};

export const DEFAULT_LOCALE_SETTINGS: LocaleDefaultSettings = {
    currency: DEFAULT_CURRENCY_CODE,
    firstDayOfWeek: WeekDay.DefaultFirstDay.type
};
