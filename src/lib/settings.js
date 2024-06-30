import timezoneConstants from '@/consts/timezone.js';
import statisticsConstants from '@/consts/statistics.js';

const settingsLocalStorageKey = 'ebk_app_settings';

const defaultSettings = {
    theme: 'auto',
    fontSize: 1,
    timeZone: '',
    debug: false,
    applicationLock: false,
    applicationLockWebAuthn: false,
    autoUpdateExchangeRatesData: true,
    autoGetCurrentGeoLocation: false,
    showAmountInHomePage: true,
    timezoneUsedForStatisticsInHomePage: timezoneConstants.defaultTimezoneTypesUsedForStatistics,
    itemsCountInTransactionListPage: 15,
    showTotalAmountInTransactionListPage: true,
    showAccountBalance: true,
    statistics: {
        defaultChartDataType: statisticsConstants.defaultChartDataType,
        defaultTimezoneType: timezoneConstants.defaultTimezoneTypesUsedForStatistics,
        defaultAccountFilter: {},
        defaultTransactionCategoryFilter: {},
        defaultSortingType: statisticsConstants.defaultSortingType,
        defaultCategoricalChartType: statisticsConstants.defaultCategoricalChartType,
        defaultCategoricalChartDataRangeType: statisticsConstants.defaultCategoricalChartDataRangeType,
        defaultTrendChartType: statisticsConstants.defaultTrendChartType,
        defaultTrendChartDataRangeType: statisticsConstants.defaultTrendChartDataRangeType,
    },
    animate: true
};

function getOriginalSettings() {
    try {
        const storageData = localStorage.getItem(settingsLocalStorageKey) || '{}';
        return JSON.parse(storageData);
    } catch (ex) {
        console.warn('settings in local storage is invalid', ex);
        return {};
    }
}

function getFinalSettings() {
    const originalSettings = getOriginalSettings();

    for (let key in originalSettings) {
        if (!Object.prototype.hasOwnProperty.call(originalSettings, key)) {
            continue;
        }

        if (typeof(defaultSettings[key]) === 'object') {
            originalSettings[key] = Object.assign({}, defaultSettings[key], originalSettings[key]);
        }
    }

    return Object.assign({}, defaultSettings, originalSettings);
}

function setSettings(settings) {
    const storageData = JSON.stringify(settings);
    return localStorage.setItem(settingsLocalStorageKey, storageData);
}

function getOption(key) {
    return getFinalSettings()[key];
}

function getSubOption(key, subKey) {
    const options = getFinalSettings()[key] || {};
    return options[subKey];
}

function setOption(key, value) {
    if (!Object.prototype.hasOwnProperty.call(defaultSettings, key)) {
        return;
    }

    const settings = getFinalSettings();
    settings[key] = value;

    return setSettings(settings);
}

function setSubOption(key, subKey, value) {
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

    options[subKey] = value;
    settings[key] = options;

    return setSettings(settings);
}

export function isEnableDebug() {
    return getOption('debug');
}

export function getTheme() {
    return getOption('theme');
}

export function setTheme(value) {
    return setOption('theme', value);
}

export function getFontSize() {
    return getOption('fontSize');
}

export function setFontSize(value) {
    return setOption('fontSize', value);
}

export function getTimeZone() {
    return getOption('timeZone');
}

export function setTimeZone(value) {
    return setOption('timeZone', value);
}

export function isEnableApplicationLock() {
    return getOption('applicationLock');
}

export function setEnableApplicationLock(value) {
    return setOption('applicationLock', value);
}

export function isEnableApplicationLockWebAuthn() {
    return getOption('applicationLockWebAuthn');
}

export function setEnableApplicationLockWebAuthn(value) {
    return setOption('applicationLockWebAuthn', value);
}

export function isAutoUpdateExchangeRatesData() {
    return getOption('autoUpdateExchangeRatesData');
}

export function setAutoUpdateExchangeRatesData(value) {
    setOption('autoUpdateExchangeRatesData', value);
}

export function isAutoGetCurrentGeoLocation() {
    return getOption('autoGetCurrentGeoLocation');
}

export function setAutoGetCurrentGeoLocation(value) {
    setOption('autoGetCurrentGeoLocation', value);
}

export function isShowAmountInHomePage() {
    return getOption('showAmountInHomePage');
}

export function setShowAmountInHomePage(value) {
    setOption('showAmountInHomePage', value);
}

export function getTimezoneUsedForStatisticsInHomePage() {
    return getOption('timezoneUsedForStatisticsInHomePage');
}

export function setTimezoneUsedForStatisticsInHomePage(value) {
    setOption('timezoneUsedForStatisticsInHomePage', value);
}

export function getItemsCountInTransactionListPage() {
    return getOption('itemsCountInTransactionListPage');
}

export function setItemsCountInTransactionListPage(value) {
    setOption('itemsCountInTransactionListPage', value);
}

export function isShowTotalAmountInTransactionListPage() {
    return getOption('showTotalAmountInTransactionListPage');
}

export function setShowTotalAmountInTransactionListPage(value) {
    setOption('showTotalAmountInTransactionListPage', value);
}

export function isShowAccountBalance() {
    return getOption('showAccountBalance');
}

export function setShowAccountBalance(value) {
    setOption('showAccountBalance', value);
}

export function getStatisticsDefaultChartDataType() {
    return getSubOption('statistics', 'defaultChartDataType');
}

export function setStatisticsDefaultChartDataType(value) {
    setSubOption('statistics', 'defaultChartDataType', value);
}

export function getStatisticsDefaultTimezoneType() {
    return getSubOption('statistics', 'defaultTimezoneType');
}

export function setStatisticsDefaultTimezoneType(value) {
    setSubOption('statistics', 'defaultTimezoneType', value);
}

export function getStatisticsDefaultAccountFilter() {
    return getSubOption('statistics', 'defaultAccountFilter');
}

export function setStatisticsDefaultAccountFilter(value) {
    setSubOption('statistics', 'defaultAccountFilter', value);
}

export function getStatisticsDefaultTransactionCategoryFilter() {
    return getSubOption('statistics', 'defaultTransactionCategoryFilter');
}

export function setStatisticsDefaultTransactionCategoryFilter(value) {
    setSubOption('statistics', 'defaultTransactionCategoryFilter', value);
}

export function getStatisticsSortingType() {
    return getSubOption('statistics', 'defaultSortingType');
}

export function setStatisticsSortingType(value) {
    setSubOption('statistics', 'defaultSortingType', value);
}

export function getStatisticsDefaultCategoricalChartType() {
    return getSubOption('statistics', 'defaultCategoricalChartType');
}

export function setStatisticsDefaultCategoricalChartType(value) {
    setSubOption('statistics', 'defaultCategoricalChartType', value);
}

export function getStatisticsDefaultCategoricalChartDataRange() {
    return getSubOption('statistics', 'defaultCategoricalChartDataRangeType');
}

export function setStatisticsDefaultCategoricalChartDataRange(value) {
    setSubOption('statistics', 'defaultCategoricalChartDataRangeType', value);
}

export function getStatisticsDefaultTrendChartType() {
    return getSubOption('statistics', 'defaultTrendChartType');
}

export function setStatisticsDefaultTrendChartType(value) {
    setSubOption('statistics', 'defaultTrendChartType', value);
}

export function getStatisticsDefaultTrendChartDataRange() {
    return getSubOption('statistics', 'defaultTrendChartDataRangeType');
}

export function setStatisticsDefaultTrendChartDataRange(value) {
    setSubOption('statistics', 'defaultTrendChartDataRangeType', value);
}

export function isEnableAnimate() {
    return getOption('animate');
}

export function setEnableAnimate(value) {
    return setOption('animate', value);
}

export function clearSettings() {
    localStorage.removeItem(settingsLocalStorageKey);
}
