import { defineStore } from 'pinia';

import currencyConstants from '@/consts/currency.js';
import datetimeConstants from '@/consts/datetime.js';
import * as settings from '@/lib/settings.js';

export const useSettingsStore = defineStore('settings', {
    state: () => ({
        appSettings: {
            theme: settings.getTheme(),
            fontSize: settings.getFontSize(),
            timeZone: settings.getTimeZone(),
            applicationLock: settings.isEnableApplicationLock(),
            applicationLockWebAuthn: settings.isEnableApplicationLockWebAuthn(),
            autoUpdateExchangeRatesData: settings.isAutoUpdateExchangeRatesData(),
            autoGetCurrentGeoLocation: settings.isAutoGetCurrentGeoLocation(),
            thousandsSeparator: settings.isEnableThousandsSeparator(),
            currencyDisplayMode: settings.getCurrencyDisplayMode(),
            showAmountInHomePage: settings.isShowAmountInHomePage(),
            timezoneUsedForStatisticsInHomePage: settings.getTimezoneUsedForStatisticsInHomePage(),
            itemsCountInTransactionListPage: settings.getItemsCountInTransactionListPage(),
            showTotalAmountInTransactionListPage: settings.isShowTotalAmountInTransactionListPage(),
            showAccountBalance: settings.isShowAccountBalance(),
            statistics: {
                defaultChartDataType: settings.getStatisticsDefaultChartDataType(),
                defaultDataRangeType: settings.getStatisticsDefaultDateRange(),
                defaultTimezoneType: settings.getStatisticsDefaultTimezoneType(),
                defaultAccountFilter: settings.getStatisticsDefaultAccountFilter(),
                defaultTransactionCategoryFilter: settings.getStatisticsDefaultTransactionCategoryFilter(),
                defaultSortingType: settings.getStatisticsSortingType(),
                defaultCategoricalChartType: settings.getStatisticsDefaultCategoricalChartType(),
                defaultTrendChartType: settings.getStatisticsDefaultTrendChartType(),
            },
            animate: settings.isEnableAnimate()
        },
        localeDefaultSettings: {
            currency: currencyConstants.defaultCurrency,
            firstDayOfWeek: datetimeConstants.defaultFirstDayOfWeek
        }
    }),
    actions: {
        setTheme(value) {
            settings.setTheme(value);
            this.appSettings.theme = value;
        },
        setFontSize(value) {
            settings.setFontSize(value);
            this.appSettings.fontSize = value;
        },
        setTimeZone(value) {
            settings.setTimeZone(value);
            this.appSettings.timeZone = value;
        },
        setEnableApplicationLock(value) {
            settings.setEnableApplicationLock(value);
            this.appSettings.applicationLock = value;
        },
        setEnableApplicationLockWebAuthn(value) {
            settings.setEnableApplicationLockWebAuthn(value);
            this.appSettings.applicationLockWebAuthn = value;
        },
        setAutoUpdateExchangeRatesData(value) {
            settings.setAutoUpdateExchangeRatesData(value);
            this.appSettings.autoUpdateExchangeRatesData = value;
        },
        setAutoGetCurrentGeoLocation(value) {
            settings.setAutoGetCurrentGeoLocation(value);
            this.appSettings.autoGetCurrentGeoLocation = value;
        },
        setEnableThousandsSeparator(value) {
            settings.setEnableThousandsSeparator(value);
            this.appSettings.thousandsSeparator = value;
        },
        setCurrencyDisplayMode(value) {
            settings.setCurrencyDisplayMode(value);
            this.appSettings.currencyDisplayMode = value;
        },
        setShowAmountInHomePage(value) {
            settings.setShowAmountInHomePage(value);
            this.appSettings.showAmountInHomePage = value;
        },
        setTimezoneUsedForStatisticsInHomePage(value) {
            settings.setTimezoneUsedForStatisticsInHomePage(value);
            this.appSettings.timezoneUsedForStatisticsInHomePage = value;
        },
        setItemsCountInTransactionListPage(value) {
            settings.setItemsCountInTransactionListPage(value);
            this.appSettings.itemsCountInTransactionListPage = value;
        },
        setShowTotalAmountInTransactionListPage(value) {
            settings.setShowTotalAmountInTransactionListPage(value);
            this.appSettings.showTotalAmountInTransactionListPage = value;
        },
        setShowAccountBalance(value) {
            settings.setShowAccountBalance(value);
            this.appSettings.showAccountBalance = value;
        },
        setStatisticsDefaultChartDataType(value) {
            settings.setStatisticsDefaultChartDataType(value);
            this.appSettings.statistics.defaultChartDataType = value;
        },
        setStatisticsDefaultDateRange(value) {
            settings.setStatisticsDefaultDateRange(value);
            this.appSettings.statistics.defaultDataRangeType = value;
        },
        setStatisticsDefaultTimezoneType(value) {
            settings.setStatisticsDefaultTimezoneType(value);
            this.appSettings.statistics.defaultTimezoneType = value;
        },
        setStatisticsDefaultAccountFilter(value) {
            settings.setStatisticsDefaultAccountFilter(value);
            this.appSettings.statistics.defaultAccountFilter = value;
        },
        setStatisticsDefaultTransactionCategoryFilter(value) {
            settings.setStatisticsDefaultTransactionCategoryFilter(value);
            this.appSettings.statistics.defaultTransactionCategoryFilter = value;
        },
        setStatisticsSortingType(value) {
            settings.setStatisticsSortingType(value);
            this.appSettings.statistics.defaultSortingType = value;
        },
        setStatisticsDefaultCategoricalChartType(value) {
            settings.setStatisticsDefaultCategoricalChartType(value);
            this.appSettings.statistics.defaultCategoricalChartType = value;
        },
        setStatisticsDefaultTrendChartType(value) {
            settings.setStatisticsDefaultTrendChartType(value);
            this.appSettings.statistics.defaultTrendChartType = value;
        },
        setEnableAnimate(value) {
            settings.setEnableAnimate(value);
            this.appSettings.animate = value;
        },
        clearAppSettings() {
            settings.clearSettings();
        },
        updateLocalizedDefaultSettings(localeDefaultSettings) {
            if (!localeDefaultSettings) {
                return;
            }

            this.localeDefaultSettings.currency = localeDefaultSettings.defaultCurrency;
            this.localeDefaultSettings.firstDayOfWeek = localeDefaultSettings.defaultFirstDayOfWeek;
        }
    }
});
