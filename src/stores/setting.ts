import { ref } from 'vue';
import { defineStore } from 'pinia';

import type { ApplicationSettings, LocaleDefaultSettings } from '@/core/setting.ts';
import {
    getApplicationSettings,
    getLocaleDefaultSettings,
    updateApplicationSettingsValue,
    updateApplicationSettingsSubValue,
    clearSettings
} from '@/lib/settings.ts';

export const useSettingsStore = defineStore('settings', () => {
    const appSettings = ref<ApplicationSettings>(getApplicationSettings());
    const localeDefaultSettings = ref<LocaleDefaultSettings>(getLocaleDefaultSettings());

    function setTheme(value: string): void {
        updateApplicationSettingsValue('theme', value);
        appSettings.value.theme = value;
    }

    function setFontSize(value: number): void {
        updateApplicationSettingsValue('fontSize', value);
        appSettings.value.fontSize = value;
    }

    function setTimeZone(value: string): void {
        updateApplicationSettingsValue('timeZone', value);
        appSettings.value.timeZone = value;
    }

    function setEnableApplicationLock(value: boolean): void {
        updateApplicationSettingsValue('applicationLock', value);
        appSettings.value.applicationLock = value;
    }

    function setEnableApplicationLockWebAuthn(value: boolean): void {
        updateApplicationSettingsValue('applicationLockWebAuthn', value);
        appSettings.value.applicationLockWebAuthn = value;
    }

    function setAutoUpdateExchangeRatesData(value: boolean): void {
        updateApplicationSettingsValue('autoUpdateExchangeRatesData', value);
        appSettings.value.autoUpdateExchangeRatesData = value;
    }

    function setAutoSaveTransactionDraft(value: string): void {
        updateApplicationSettingsValue('autoSaveTransactionDraft', value);
        appSettings.value.autoSaveTransactionDraft = value;
    }

    function setAutoGetCurrentGeoLocation(value: boolean): void {
        updateApplicationSettingsValue('autoGetCurrentGeoLocation', value);
        appSettings.value.autoGetCurrentGeoLocation = value;
    }

    function setAlwaysShowTransactionPicturesInMobileTransactionEditPage(value: boolean): void {
        updateApplicationSettingsValue('alwaysShowTransactionPicturesInMobileTransactionEditPage', value);
        appSettings.value.alwaysShowTransactionPicturesInMobileTransactionEditPage = value;
    }

    function setShowAddTransactionButtonInDesktopNavbar(value: boolean): void {
        updateApplicationSettingsValue('showAddTransactionButtonInDesktopNavbar', value);
        appSettings.value.showAddTransactionButtonInDesktopNavbar = value;
    }

    function setShowAmountInHomePage(value: boolean): void {
        updateApplicationSettingsValue('showAmountInHomePage', value);
        appSettings.value.showAmountInHomePage = value;
    }

    function setTimezoneUsedForStatisticsInHomePage(value: number): void {
        updateApplicationSettingsValue('timezoneUsedForStatisticsInHomePage', value);
        appSettings.value.timezoneUsedForStatisticsInHomePage = value;
    }

    function setItemsCountInTransactionListPage(value: number): void {
        updateApplicationSettingsValue('itemsCountInTransactionListPage', value);
        appSettings.value.itemsCountInTransactionListPage = value;
    }

    function setShowTotalAmountInTransactionListPage(value: boolean): void {
        updateApplicationSettingsValue('showTotalAmountInTransactionListPage', value);
        appSettings.value.showTotalAmountInTransactionListPage = value;
    }

    function setShowTagInTransactionListPage(value: boolean): void {
        updateApplicationSettingsValue('showTagInTransactionListPage', value);
        appSettings.value.showTagInTransactionListPage = value;
    }

    function setShowAccountBalance(value: boolean): void {
        updateApplicationSettingsValue('showAccountBalance', value);
        appSettings.value.showAccountBalance = value;
    }

    function setCurrencySortByInExchangeRatesPage(value: number): void {
        updateApplicationSettingsValue('currencySortByInExchangeRatesPage', value);
        appSettings.value.currencySortByInExchangeRatesPage = value;
    }

    function setStatisticsDefaultChartDataType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultChartDataType', value);
        appSettings.value.statistics.defaultChartDataType = value;
    }

    function setStatisticsDefaultTimezoneType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultTimezoneType', value);
        appSettings.value.statistics.defaultTimezoneType = value;
    }

    function setStatisticsDefaultAccountFilter(value: Record<string, boolean>): void {
        updateApplicationSettingsSubValue('statistics', 'defaultAccountFilter', value);
        appSettings.value.statistics.defaultAccountFilter = value;
    }

    function setStatisticsDefaultTransactionCategoryFilter(value: Record<string, boolean>): void {
        updateApplicationSettingsSubValue('statistics', 'defaultTransactionCategoryFilter', value);
        appSettings.value.statistics.defaultTransactionCategoryFilter = value;
    }

    function setStatisticsSortingType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultSortingType', value);
        appSettings.value.statistics.defaultSortingType = value;
    }

    function setStatisticsDefaultCategoricalChartType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultCategoricalChartType', value);
        appSettings.value.statistics.defaultCategoricalChartType = value;
    }

    function setStatisticsDefaultCategoricalChartDateRange(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultCategoricalChartDataRangeType', value);
        appSettings.value.statistics.defaultCategoricalChartDataRangeType = value;
    }

    function setStatisticsDefaultTrendChartType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultTrendChartType', value);
        appSettings.value.statistics.defaultTrendChartType = value;
    }

    function setStatisticsDefaultTrendChartDateRange(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultTrendChartDataRangeType', value);
        appSettings.value.statistics.defaultTrendChartDataRangeType = value;
    }

    function setEnableAnimate(value: boolean): void {
        updateApplicationSettingsValue('animate', value);
        appSettings.value.animate = value;
    }

    function clearAppSettings(): void {
        clearSettings();
        appSettings.value = getApplicationSettings();
    }

    function updateLocalizedDefaultSettings(newLocaleDefaultSettings: LocaleDefaultSettings | null) {
        if (!newLocaleDefaultSettings) {
            return;
        }

        localeDefaultSettings.value.currency = newLocaleDefaultSettings.currency;
        localeDefaultSettings.value.firstDayOfWeek = newLocaleDefaultSettings.firstDayOfWeek;
    }

    return {
        // states
        appSettings,
        localeDefaultSettings,
        // functions
        setTheme,
        setFontSize,
        setTimeZone,
        setEnableApplicationLock,
        setEnableApplicationLockWebAuthn,
        setAutoUpdateExchangeRatesData,
        setAutoSaveTransactionDraft,
        setAutoGetCurrentGeoLocation,
        setAlwaysShowTransactionPicturesInMobileTransactionEditPage,
        setShowAddTransactionButtonInDesktopNavbar,
        setShowAmountInHomePage,
        setTimezoneUsedForStatisticsInHomePage,
        setItemsCountInTransactionListPage,
        setShowTotalAmountInTransactionListPage,
        setShowTagInTransactionListPage,
        setShowAccountBalance,
        setCurrencySortByInExchangeRatesPage,
        setStatisticsDefaultChartDataType,
        setStatisticsDefaultTimezoneType,
        setStatisticsDefaultAccountFilter,
        setStatisticsDefaultTransactionCategoryFilter,
        setStatisticsSortingType,
        setStatisticsDefaultCategoricalChartType,
        setStatisticsDefaultCategoricalChartDateRange,
        setStatisticsDefaultTrendChartType,
        setStatisticsDefaultTrendChartDateRange,
        setEnableAnimate,
        clearAppSettings,
        updateLocalizedDefaultSettings
    };
});
