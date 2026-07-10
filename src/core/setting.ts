import { type WeekDayValue, WeekDay } from './datetime.ts';
import { TimezoneTypeForStatistics } from './timezone.ts';
import { CurrencySortingType } from './currency.ts';
import { KeywordMatchMode } from './text.ts';
import { ImageUploadQualityType } from './image.ts';
import {
    TransactionQuickSaveButtonStyle,
    TransactionQuickAddButtonActionType
} from './transaction.ts';
import {
    CategoricalChartType,
    TrendChartType,
    ChartDataType,
    ChartSortingType,
    DEFAULT_CATEGORICAL_CHART_DATA_RANGE,
    DEFAULT_TREND_CHART_DATA_RANGE,
    DEFAULT_ASSET_TRENDS_CHART_DATA_RANGE,
    DEFAULT_RECONCILIATION_STATEMENT_DATE_RANGE_IN_DESKTOP,
    DEFAULT_RECONCILIATION_STATEMENT_DATE_RANGE_IN_MOBILE,
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
    // General Settings
    chartColors: string;
    // Navigation Bar
    showAddTransactionButtonInDesktopNavbar: boolean;
    // Overview Page
    showAmountInHomePage: boolean;
    timezoneUsedForStatisticsInHomePage: number;
    overviewAccountFilterInHomePage: Record<string, boolean>;
    overviewTransactionCategoryFilterInHomePage: Record<string, boolean>;
    // Transaction List Page
    quickSaveButtonStyleInMobileTransactionListPage: number;
    quickAddButtonActionInMobileTransactionEditPage: number;
    itemsCountInTransactionListPage: number;
    showTotalAmountInTransactionListPage: boolean;
    showTagInTransactionListPage: boolean;
    defaultKeywordMatchModeInTransactionListPage: number;
    // Transaction Edit Page
    autoSaveTransactionDraft: string;
    autoGetCurrentGeoLocation: boolean;
    alwaysShowTransactionPicturesInMobileTransactionEditPage: boolean;
    transactionPictureQuality: number;
    // AI Clipboard Text Recognition
    alwaysRequireConfirmationOfClipboardContentBeforeSubmission: boolean;
    // AI Image Recognition
    autoUploadTransactionPictureForAIRecognition: boolean;
    // Import Transaction Dialog
    rememberLastSelectedFileTypeInImportTransactionDialog: boolean;
    lastSelectedFileTypeInImportTransactionDialog: string;
    // Insights Explorer Page
    insightsExplorerDefaultDateRangeType: number;
    showTagInInsightsExplorerPage: boolean;
    // Account List Page
    totalAmountExcludeAccountIds: Record<string, boolean>;
    accountCategoryOrders: string;
    hideCategoriesWithoutAccounts: boolean;
    reconciliationStatementButtonDefaultDateRangeTypeInDesktop: number;
    reconciliationStatementPageDefaultDateRangeTypeInMobile: number;
    // Exchange Rates Data Page
    currencySortByInExchangeRatesPage: number;
    // Browser Cache Management
    mapCacheExpiration: number,
    exchangeRatesDataCacheExpiration: number,
    // Statistics Settings
    statistics: {
        defaultChartDataType: number;
        defaultTimezoneType: number;
        defaultAccountFilter: Record<string, boolean>;
        defaultTransactionCategoryFilter: Record<string, boolean>;
        defaultKeywordMatchMode: number;
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
    'autoUpdateExchangeRatesData': UserApplicationCloudSettingType.Boolean,
    // General Settings
    'chartColors': UserApplicationCloudSettingType.String,
    // Navigation Bar
    'showAddTransactionButtonInDesktopNavbar': UserApplicationCloudSettingType.Boolean,
    // Overview Page
    'showAmountInHomePage': UserApplicationCloudSettingType.Boolean,
    'timezoneUsedForStatisticsInHomePage': UserApplicationCloudSettingType.Number,
    'overviewAccountFilterInHomePage': UserApplicationCloudSettingType.StringBooleanMap,
    'overviewTransactionCategoryFilterInHomePage': UserApplicationCloudSettingType.StringBooleanMap,
    // Transaction List Page
    'itemsCountInTransactionListPage': UserApplicationCloudSettingType.Number,
    'showTotalAmountInTransactionListPage': UserApplicationCloudSettingType.Boolean,
    'showTagInTransactionListPage': UserApplicationCloudSettingType.Boolean,
    'defaultKeywordMatchModeInTransactionListPage': UserApplicationCloudSettingType.Number,
    // Transaction Edit Page
    'quickSaveButtonStyleInMobileTransactionListPage': UserApplicationCloudSettingType.Number,
    'quickAddButtonActionInMobileTransactionEditPage': UserApplicationCloudSettingType.Number,
    'autoSaveTransactionDraft': UserApplicationCloudSettingType.String,
    'autoGetCurrentGeoLocation': UserApplicationCloudSettingType.Boolean,
    'alwaysShowTransactionPicturesInMobileTransactionEditPage': UserApplicationCloudSettingType.Boolean,
    'transactionPictureQuality': UserApplicationCloudSettingType.Number,
    // AI Clipboard Text Recognition
    'alwaysRequireConfirmationOfClipboardContentBeforeSubmission': UserApplicationCloudSettingType.Boolean,
    // AI Image Recognition
    'autoUploadTransactionPictureForAIRecognition': UserApplicationCloudSettingType.Boolean,
    // Import Transaction Dialog
    'rememberLastSelectedFileTypeInImportTransactionDialog': UserApplicationCloudSettingType.Boolean,
    'lastSelectedFileTypeInImportTransactionDialog': UserApplicationCloudSettingType.String,
    // Insights Explorer Page
    'insightsExplorerDefaultDateRangeType': UserApplicationCloudSettingType.Number,
    'showTagInInsightsExplorerPage': UserApplicationCloudSettingType.Boolean,
    // Account List Page
    'totalAmountExcludeAccountIds': UserApplicationCloudSettingType.StringBooleanMap,
    'accountCategoryOrders': UserApplicationCloudSettingType.String,
    'hideCategoriesWithoutAccounts': UserApplicationCloudSettingType.Boolean,
    'reconciliationStatementButtonDefaultDateRangeTypeInDesktop': UserApplicationCloudSettingType.Number,
    'reconciliationStatementPageDefaultDateRangeTypeInMobile': UserApplicationCloudSettingType.Number,
    // Exchange Rates Data Page
    'currencySortByInExchangeRatesPage': UserApplicationCloudSettingType.Number,
    // Browser Cache Management
    'mapCacheExpiration': UserApplicationCloudSettingType.Number,
    'exchangeRatesDataCacheExpiration': UserApplicationCloudSettingType.Number,
    // Statistics Settings
    'statistics.defaultChartDataType': UserApplicationCloudSettingType.Number,
    'statistics.defaultTimezoneType': UserApplicationCloudSettingType.Number,
    'statistics.defaultAccountFilter': UserApplicationCloudSettingType.StringBooleanMap,
    'statistics.defaultTransactionCategoryFilter': UserApplicationCloudSettingType.StringBooleanMap,
    'statistics.defaultKeywordMatchMode': UserApplicationCloudSettingType.Number,
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
    // General Settings
    chartColors: '',
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
    defaultKeywordMatchModeInTransactionListPage: KeywordMatchMode.Default.type,
    // Transaction Edit Page
    quickSaveButtonStyleInMobileTransactionListPage: TransactionQuickSaveButtonStyle.Default.type,
    quickAddButtonActionInMobileTransactionEditPage: TransactionQuickAddButtonActionType.Default.type,
    autoSaveTransactionDraft: 'disabled',
    autoGetCurrentGeoLocation: false,
    alwaysShowTransactionPicturesInMobileTransactionEditPage: false,
    transactionPictureQuality: ImageUploadQualityType.Default.type,
    // AI Clipboard Text Recognition
    alwaysRequireConfirmationOfClipboardContentBeforeSubmission: true,
    // AI Image Recognition
    autoUploadTransactionPictureForAIRecognition: false,
    // Import Transaction Dialog
    rememberLastSelectedFileTypeInImportTransactionDialog: true,
    lastSelectedFileTypeInImportTransactionDialog: '',
    // Insights Explorer Page
    insightsExplorerDefaultDateRangeType: DEFAULT_TRANSACTION_EXPLORER_DATE_RANGE.type,
    showTagInInsightsExplorerPage: true,
    // Account List Page
    totalAmountExcludeAccountIds: {},
    accountCategoryOrders: '',
    hideCategoriesWithoutAccounts: false,
    reconciliationStatementButtonDefaultDateRangeTypeInDesktop: DEFAULT_RECONCILIATION_STATEMENT_DATE_RANGE_IN_DESKTOP.type,
    reconciliationStatementPageDefaultDateRangeTypeInMobile: DEFAULT_RECONCILIATION_STATEMENT_DATE_RANGE_IN_MOBILE.type,
    // Exchange Rates Data Page
    currencySortByInExchangeRatesPage: CurrencySortingType.Default.type,
    // Browser Cache Management
    mapCacheExpiration: -1,
    exchangeRatesDataCacheExpiration: 0,
    // Statistics Settings
    statistics: {
        defaultChartDataType: ChartDataType.Default.type,
        defaultTimezoneType: TimezoneTypeForStatistics.Default.type,
        defaultAccountFilter: {},
        defaultTransactionCategoryFilter: {},
        defaultKeywordMatchMode: KeywordMatchMode.Default.type,
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
