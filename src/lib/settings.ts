import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { CurrencySortingType } from '@/core/currency.ts';
import {
    CategoricalChartType,
    TrendChartType,
    ChartDataType,
    ChartSortingType,
    DEFAULT_CATEGORICAL_CHART_DATA_RANGE,
    DEFAULT_TREND_CHART_DATA_RANGE
} from '@/core/statistics.ts';
import { isObject } from './common.ts';

const settingsLocalStorageKey: string = 'ebk_app_settings';

export type ApplicationSettingKey = string;
export type ApplicationSettingValue = string | number | boolean | Record<string, ApplicationSettingSubValue>;
export type ApplicationSettingSubValue = string | number | boolean | Record<string, boolean> | Record<string, number>;

export interface ApplicationSettings {
    [key: ApplicationSettingKey]: ApplicationSettingValue;
}

const defaultSettings: ApplicationSettings = {
    theme: 'auto',
    fontSize: 1,
    timeZone: '',
    debug: false,
    applicationLock: false,
    applicationLockWebAuthn: false,
    autoUpdateExchangeRatesData: true,
    autoSaveTransactionDraft: 'disabled',
    autoGetCurrentGeoLocation: false,
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

function getOriginalSettings(): ApplicationSettings {
    try {
        const storageData = localStorage.getItem(settingsLocalStorageKey) || '{}';
        return JSON.parse(storageData);
    } catch (ex) {
        console.warn('settings in local storage is invalid', ex);
        return {};
    }
}

function getFinalSettings(): ApplicationSettings {
    const originalSettings = getOriginalSettings();

    for (const key in originalSettings) {
        if (!Object.prototype.hasOwnProperty.call(originalSettings, key)) {
            continue;
        }

        if (typeof(defaultSettings[key]) === 'object') {
            originalSettings[key] = Object.assign({}, defaultSettings[key], originalSettings[key]);
        }
    }

    return Object.assign({}, defaultSettings, originalSettings);
}

function setSettings(settings: ApplicationSettings): void {
    const storageData = JSON.stringify(settings);
    return localStorage.setItem(settingsLocalStorageKey, storageData);
}

function getOption(key: ApplicationSettingKey): ApplicationSettingValue | undefined {
    return getFinalSettings()[key];
}

function getSubOption(key: ApplicationSettingKey, subKey: ApplicationSettingKey): ApplicationSettingSubValue | undefined {
    const options = getFinalSettings()[key];

    if (!isObject(options)) {
        return undefined;
    }

    return (options as Record<string, ApplicationSettingSubValue>)[subKey];
}

function setOption(key: ApplicationSettingKey, value: ApplicationSettingValue): void {
    if (!Object.prototype.hasOwnProperty.call(defaultSettings, key)) {
        return;
    }

    const settings = getFinalSettings();
    settings[key] = value;

    return setSettings(settings);
}

function setSubOption(key: ApplicationSettingKey, subKey: ApplicationSettingKey, value: ApplicationSettingSubValue): void {
    if (!Object.prototype.hasOwnProperty.call(defaultSettings, key)) {
        return;
    }

    if (!Object.prototype.hasOwnProperty.call(defaultSettings[key], subKey)) {
        return;
    }

    const settings = getFinalSettings();
    let options = settings[key];

    if (!options) {
        options = {};
    }

    (options as Record<string, ApplicationSettingSubValue>)[subKey] = value;
    settings[key] = options;

    return setSettings(settings);
}

export function isEnableDebug(): boolean {
    return getOption('debug') as boolean;
}

export function getTheme(): string {
    return getOption('theme') as string;
}

export function setTheme(value: string): void {
    return setOption('theme', value);
}

export function getFontSize(): number {
    return getOption('fontSize') as number;
}

export function setFontSize(value: number): void {
    return setOption('fontSize', value);
}

export function getTimeZone(): string {
    return getOption('timeZone') as string;
}

export function setTimeZone(value: string): void {
    return setOption('timeZone', value);
}

export function isEnableApplicationLock(): boolean {
    return getOption('applicationLock') as boolean;
}

export function setEnableApplicationLock(value: boolean): void {
    return setOption('applicationLock', value);
}

export function isEnableApplicationLockWebAuthn(): boolean {
    return getOption('applicationLockWebAuthn') as boolean;
}

export function setEnableApplicationLockWebAuthn(value: boolean): void {
    return setOption('applicationLockWebAuthn', value);
}

export function isAutoUpdateExchangeRatesData(): boolean {
    return getOption('autoUpdateExchangeRatesData') as boolean;
}

export function setAutoUpdateExchangeRatesData(value: boolean): void {
    setOption('autoUpdateExchangeRatesData', value);
}

export function getAutoSaveTransactionDraft(): boolean {
    return getOption('autoSaveTransactionDraft') as boolean;
}

export function setAutoSaveTransactionDraft(value: boolean): void {
    setOption('autoSaveTransactionDraft', value);
}

export function isAutoGetCurrentGeoLocation(): boolean {
    return getOption('autoGetCurrentGeoLocation') as boolean;
}

export function setAutoGetCurrentGeoLocation(value: boolean): void {
    setOption('autoGetCurrentGeoLocation', value);
}

export function isShowAmountInHomePage(): boolean {
    return getOption('showAmountInHomePage') as boolean;
}

export function setShowAmountInHomePage(value: boolean): void {
    setOption('showAmountInHomePage', value);
}

export function getTimezoneUsedForStatisticsInHomePage(): number {
    return getOption('timezoneUsedForStatisticsInHomePage') as number;
}

export function setTimezoneUsedForStatisticsInHomePage(value: number): void {
    setOption('timezoneUsedForStatisticsInHomePage', value);
}

export function getItemsCountInTransactionListPage(): number {
    return getOption('itemsCountInTransactionListPage') as number;
}

export function setItemsCountInTransactionListPage(value: number): void {
    setOption('itemsCountInTransactionListPage', value);
}

export function isShowTotalAmountInTransactionListPage(): boolean {
    return getOption('showTotalAmountInTransactionListPage') as boolean;
}

export function setShowTotalAmountInTransactionListPage(value: boolean): void {
    setOption('showTotalAmountInTransactionListPage', value);
}

export function isShowTagInTransactionListPage(): boolean {
    return getOption('showTagInTransactionListPage') as boolean;
}

export function setShowTagInTransactionListPage(value: boolean): void {
    setOption('showTagInTransactionListPage', value);
}

export function isShowAccountBalance(): boolean {
    return getOption('showAccountBalance') as boolean;
}

export function setShowAccountBalance(value: boolean): void {
    setOption('showAccountBalance', value);
}

export function getCurrencySortByInExchangeRatesPage(): number {
    return getOption('currencySortByInExchangeRatesPage') as number;
}

export function setCurrencySortByInExchangeRatesPage(value: number): void {
    setOption('currencySortByInExchangeRatesPage', value);
}

export function getStatisticsDefaultChartDataType(): number {
    return getSubOption('statistics', 'defaultChartDataType') as number;
}

export function setStatisticsDefaultChartDataType(value: number): void {
    setSubOption('statistics', 'defaultChartDataType', value);
}

export function getStatisticsDefaultTimezoneType(): number {
    return getSubOption('statistics', 'defaultTimezoneType') as number;
}

export function setStatisticsDefaultTimezoneType(value: number): void {
    setSubOption('statistics', 'defaultTimezoneType', value);
}

export function getStatisticsDefaultAccountFilter(): Record<string, boolean> {
    return getSubOption('statistics', 'defaultAccountFilter') as Record<string, boolean>;
}

export function setStatisticsDefaultAccountFilter(value: Record<string, boolean>): void {
    setSubOption('statistics', 'defaultAccountFilter', value);
}

export function getStatisticsDefaultTransactionCategoryFilter(): Record<string, boolean> {
    return getSubOption('statistics', 'defaultTransactionCategoryFilter') as Record<string, boolean>;
}

export function setStatisticsDefaultTransactionCategoryFilter(value: Record<string, boolean>): void {
    setSubOption('statistics', 'defaultTransactionCategoryFilter', value);
}

export function getStatisticsSortingType(): number {
    return getSubOption('statistics', 'defaultSortingType') as number;
}

export function setStatisticsSortingType(value: number): void {
    setSubOption('statistics', 'defaultSortingType', value);
}

export function getStatisticsDefaultCategoricalChartType(): number {
    return getSubOption('statistics', 'defaultCategoricalChartType') as number;
}

export function setStatisticsDefaultCategoricalChartType(value: number): void {
    setSubOption('statistics', 'defaultCategoricalChartType', value);
}

export function getStatisticsDefaultCategoricalChartDataRange(): number {
    return getSubOption('statistics', 'defaultCategoricalChartDataRangeType') as number;
}

export function setStatisticsDefaultCategoricalChartDataRange(value: number): void {
    setSubOption('statistics', 'defaultCategoricalChartDataRangeType', value);
}

export function getStatisticsDefaultTrendChartType(): number {
    return getSubOption('statistics', 'defaultTrendChartType') as number;
}

export function setStatisticsDefaultTrendChartType(value: number): void {
    setSubOption('statistics', 'defaultTrendChartType', value);
}

export function getStatisticsDefaultTrendChartDataRange(): number {
    return getSubOption('statistics', 'defaultTrendChartDataRangeType') as number;
}

export function setStatisticsDefaultTrendChartDataRange(value: number): void {
    setSubOption('statistics', 'defaultTrendChartDataRangeType', value);
}

export function isEnableAnimate(): boolean {
    return getOption('animate') as boolean;
}

export function setEnableAnimate(value: boolean) {
    return setOption('animate', value);
}

export function clearSettings(): void {
    localStorage.removeItem(settingsLocalStorageKey);
}
