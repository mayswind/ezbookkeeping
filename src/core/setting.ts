import { type WeekDayValue, WeekDay } from './datetime.ts';
import { TimezoneTypeForStatistics } from './timezone.ts';
import { CurrencySortingType } from './currency.ts';
import {
    CategoricalChartType,
    TrendChartType,
    ChartDataType,
    ChartSortingType,
    DEFAULT_CATEGORICAL_CHART_DATA_RANGE,
    DEFAULT_TREND_CHART_DATA_RANGE
} from './statistics.ts';
import { DEFAULT_CURRENCY_CODE } from '@/consts/currency.ts';

export type ApplicationSettingKey = string;
export type ApplicationSettingValue = string | number | boolean | Record<string, ApplicationSettingSubValue>;
export type ApplicationSettingSubValue = string | number | boolean | Record<string, boolean> | Record<string, number>;

export interface BaseApplicationSetting {
    [key: ApplicationSettingKey]: ApplicationSettingValue;
}

export interface ApplicationSettings extends BaseApplicationSetting {
    theme: string;
    fontSize: number;
    timeZone: string;
    debug: boolean;
    applicationLock: boolean;
    applicationLockWebAuthn: boolean;
    autoUpdateExchangeRatesData: boolean;
    autoSaveTransactionDraft: string;
    autoGetCurrentGeoLocation: boolean;
    alwaysShowTransactionPicturesInMobileTransactionEditPage: boolean;
    showAddTransactionButtonInDesktopNavbar: boolean;
    showAmountInHomePage: boolean;
    timezoneUsedForStatisticsInHomePage: number;
    itemsCountInTransactionListPage: number;
    showTotalAmountInTransactionListPage: boolean;
    showTagInTransactionListPage: boolean;
    showAccountBalance: boolean;
    currencySortByInExchangeRatesPage: number;
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
    };
    animate: boolean;
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

export const DEFAULT_APPLICATION_SETTINGS: ApplicationSettings = {
    theme: 'auto',
    fontSize: 1,
    timeZone: '',
    debug: false,
    applicationLock: false,
    applicationLockWebAuthn: false,
    autoUpdateExchangeRatesData: true,
    autoSaveTransactionDraft: 'disabled',
    autoGetCurrentGeoLocation: false,
    alwaysShowTransactionPicturesInMobileTransactionEditPage: false,
    showAddTransactionButtonInDesktopNavbar: true,
    showAmountInHomePage: true,
    timezoneUsedForStatisticsInHomePage: TimezoneTypeForStatistics.Default.type,
    itemsCountInTransactionListPage: 15,
    showTotalAmountInTransactionListPage: true,
    showTagInTransactionListPage: true,
    showAccountBalance: true,
    currencySortByInExchangeRatesPage: CurrencySortingType.Default.type,
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
    },
    animate: true
};

export const DEFAULT_LOCALE_SETTINGS: LocaleDefaultSettings = {
    currency: DEFAULT_CURRENCY_CODE,
    firstDayOfWeek: WeekDay.DefaultFirstDay.type
};
