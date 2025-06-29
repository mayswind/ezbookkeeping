import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import {
    type ApplicationSettingValue,
    type ApplicationSettingSubValue,
    type ApplicationSettings,
    type ApplicationCloudSetting,
    type LocaleDefaultSettings,
    UserApplicationCloudSettingType,
    ALL_ALLOWED_CLOUD_SYNC_APP_SETTING_KEY_TYPES
} from '@/core/setting.ts';

import {
    isObject,
    isString,
    isBoolean,
    getObjectOwnFieldCount,
    arrayItemToObjectField
} from '@/lib/common.ts';

import {
    getApplicationSettings,
    getLocaleDefaultSettings,
    updateApplicationSettingsValue,
    updateApplicationSettingsSubValue,
    clearSettings
} from '@/lib/settings.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

export const useSettingsStore = defineStore('settings', () => {
    const appSettings = ref<ApplicationSettings>(getApplicationSettings());
    const syncedAppSettings = ref<Record<string, boolean>>({});
    const localeDefaultSettings = ref<LocaleDefaultSettings>(getLocaleDefaultSettings());

    const enableApplicationCloudSync = computed<boolean>(() => getObjectOwnFieldCount(syncedAppSettings.value) > 0);

    function updateApplicationSettingsValueAndAppSettingsFromCloudSetting(key: string, value: string | number | boolean | Record<string, boolean>): void {
        const keyItems = key.split('.');

        if (keyItems.length === 1) {
            updateApplicationSettingsValue(keyItems[0], value);
            appSettings.value[keyItems[0]] = value;
        } else if (keyItems.length === 2) {
            updateApplicationSettingsSubValue(keyItems[0], keyItems[1], value);
            (appSettings.value[keyItems[0]] as Record<string, ApplicationSettingSubValue>)[keyItems[1]] = value;
        } else {
            logger.warn(`cannot load application cloud setting "${key}", because it has invalid key format`);
        }
    }

    function createUserApplicationCloudSetting(key: string): ApplicationCloudSetting | null {
        const settingType = ALL_ALLOWED_CLOUD_SYNC_APP_SETTING_KEY_TYPES[key];

        if (!settingType) {
            logger.warn(`cannot get application cloud setting "${key}", because it is not supported to sync`);
            return null;
        }

        const keyItems = key.split('.');
        let value: ApplicationSettingValue | ApplicationSettingSubValue = appSettings.value[key];

        if (keyItems.length === 2) {
            value = (appSettings.value[keyItems[0]] as Record<string, ApplicationSettingSubValue>)[keyItems[1]];
        } else if (keyItems.length > 2) {
            logger.warn(`cannot get application cloud setting "${key}", because it has invalid key format`);
            return null;
        }

        let settingValue = '';

        if (settingType === UserApplicationCloudSettingType.String) {
            if (!value) {
                settingValue = '';
            } else {
                settingValue = value.toString();
            }
        } else {
            settingValue = JSON.stringify(value);
        }

        return {
            settingKey: key,
            settingValue: settingValue
        };
    }

    function updateUserApplicationCloudSettingValue(key: string, value: string | number | boolean | Record<string, boolean>): void {
        if (!syncedAppSettings.value || !syncedAppSettings.value[key]) {
            return;
        }

        const settingType = ALL_ALLOWED_CLOUD_SYNC_APP_SETTING_KEY_TYPES[key];

        if (!settingType) {
            return;
        }

        const settingValue = isString(value) ? value : JSON.stringify(value);

        services.updateUserApplicationCloudSettings({
            settings: [{
                settingKey: key,
                settingValue: settingValue
            }],
            fullUpdate: false
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                logger.debug(`failed to update user application cloud setting "${key}" with value "${settingValue}"`);
                return;
            }

            logger.debug(`update user application cloud setting "${key}" with value "${settingValue}" successfully`);
        }).catch(error => {
            logger.debug(`failed to update user application cloud setting "${key}" with value "${settingValue}"`, error);
        });
    }

    // Basic Settings
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

    function setAutoUpdateExchangeRatesData(value: boolean): void {
        updateApplicationSettingsValue('autoUpdateExchangeRatesData', value);
        appSettings.value.autoUpdateExchangeRatesData = value;
    }

    function setShowAccountBalance(value: boolean): void {
        updateApplicationSettingsValue('showAccountBalance', value);
        appSettings.value.showAccountBalance = value;
        updateUserApplicationCloudSettingValue('showAccountBalance', value);
    }

    function setEnableAnimate(value: boolean): void {
        updateApplicationSettingsValue('animate', value);
        appSettings.value.animate = value;
    }

    // Application Lock
    function setEnableApplicationLock(value: boolean): void {
        updateApplicationSettingsValue('applicationLock', value);
        appSettings.value.applicationLock = value;
    }

    function setEnableApplicationLockWebAuthn(value: boolean): void {
        updateApplicationSettingsValue('applicationLockWebAuthn', value);
        appSettings.value.applicationLockWebAuthn = value;
    }

    // Navigation Bar
    function setShowAddTransactionButtonInDesktopNavbar(value: boolean): void {
        updateApplicationSettingsValue('showAddTransactionButtonInDesktopNavbar', value);
        appSettings.value.showAddTransactionButtonInDesktopNavbar = value;
    }

    // Overview Page
    function setShowAmountInHomePage(value: boolean): void {
        updateApplicationSettingsValue('showAmountInHomePage', value);
        appSettings.value.showAmountInHomePage = value;
        updateUserApplicationCloudSettingValue('showAmountInHomePage', value);
    }

    function setTimezoneUsedForStatisticsInHomePage(value: number): void {
        updateApplicationSettingsValue('timezoneUsedForStatisticsInHomePage', value);
        appSettings.value.timezoneUsedForStatisticsInHomePage = value;
        updateUserApplicationCloudSettingValue('timezoneUsedForStatisticsInHomePage', value);
    }

    // Transaction List Page
    function setItemsCountInTransactionListPage(value: number): void {
        updateApplicationSettingsValue('itemsCountInTransactionListPage', value);
        appSettings.value.itemsCountInTransactionListPage = value;
        updateUserApplicationCloudSettingValue('itemsCountInTransactionListPage', value);
    }

    function setShowTotalAmountInTransactionListPage(value: boolean): void {
        updateApplicationSettingsValue('showTotalAmountInTransactionListPage', value);
        appSettings.value.showTotalAmountInTransactionListPage = value;
        updateUserApplicationCloudSettingValue('showTotalAmountInTransactionListPage', value);
    }

    function setShowTagInTransactionListPage(value: boolean): void {
        updateApplicationSettingsValue('showTagInTransactionListPage', value);
        appSettings.value.showTagInTransactionListPage = value;
        updateUserApplicationCloudSettingValue('showTagInTransactionListPage', value);
    }

    // Transaction Edit Page
    function setAutoSaveTransactionDraft(value: string): void {
        updateApplicationSettingsValue('autoSaveTransactionDraft', value);
        appSettings.value.autoSaveTransactionDraft = value;
        updateUserApplicationCloudSettingValue('autoSaveTransactionDraft', value);
    }

    function setAutoGetCurrentGeoLocation(value: boolean): void {
        updateApplicationSettingsValue('autoGetCurrentGeoLocation', value);
        appSettings.value.autoGetCurrentGeoLocation = value;
        updateUserApplicationCloudSettingValue('autoGetCurrentGeoLocation', value);
    }

    function setAlwaysShowTransactionPicturesInMobileTransactionEditPage(value: boolean): void {
        updateApplicationSettingsValue('alwaysShowTransactionPicturesInMobileTransactionEditPage', value);
        appSettings.value.alwaysShowTransactionPicturesInMobileTransactionEditPage = value;
        updateUserApplicationCloudSettingValue('alwaysShowTransactionPicturesInMobileTransactionEditPage', value);
    }

    // Exchange Rates Data Page
    function setCurrencySortByInExchangeRatesPage(value: number): void {
        updateApplicationSettingsValue('currencySortByInExchangeRatesPage', value);
        appSettings.value.currencySortByInExchangeRatesPage = value;
        updateUserApplicationCloudSettingValue('currencySortByInExchangeRatesPage', value);
    }

    // Statistics Settings
    function setStatisticsDefaultChartDataType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultChartDataType', value);
        appSettings.value.statistics.defaultChartDataType = value;
        updateUserApplicationCloudSettingValue('statistics.defaultChartDataType', value);
    }

    function setStatisticsDefaultTimezoneType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultTimezoneType', value);
        appSettings.value.statistics.defaultTimezoneType = value;
        updateUserApplicationCloudSettingValue('statistics.defaultTimezoneType', value);
    }

    function setStatisticsDefaultAccountFilter(value: Record<string, boolean>): void {
        updateApplicationSettingsSubValue('statistics', 'defaultAccountFilter', value);
        appSettings.value.statistics.defaultAccountFilter = value;
        updateUserApplicationCloudSettingValue('statistics.defaultAccountFilter', value);
    }

    function setStatisticsDefaultTransactionCategoryFilter(value: Record<string, boolean>): void {
        updateApplicationSettingsSubValue('statistics', 'defaultTransactionCategoryFilter', value);
        appSettings.value.statistics.defaultTransactionCategoryFilter = value;
        updateUserApplicationCloudSettingValue('statistics.defaultTransactionCategoryFilter', value);
    }

    function setStatisticsSortingType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultSortingType', value);
        appSettings.value.statistics.defaultSortingType = value;
        updateUserApplicationCloudSettingValue('statistics.defaultSortingType', value);
    }

    function setStatisticsDefaultCategoricalChartType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultCategoricalChartType', value);
        appSettings.value.statistics.defaultCategoricalChartType = value;
        updateUserApplicationCloudSettingValue('statistics.defaultCategoricalChartType', value);
    }

    function setStatisticsDefaultCategoricalChartDateRange(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultCategoricalChartDataRangeType', value);
        appSettings.value.statistics.defaultCategoricalChartDataRangeType = value;
        updateUserApplicationCloudSettingValue('statistics.defaultCategoricalChartDataRangeType', value);
    }

    function setStatisticsDefaultTrendChartType(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultTrendChartType', value);
        appSettings.value.statistics.defaultTrendChartType = value;
        updateUserApplicationCloudSettingValue('statistics.defaultTrendChartType', value);
    }

    function setStatisticsDefaultTrendChartDateRange(value: number): void {
        updateApplicationSettingsSubValue('statistics', 'defaultTrendChartDataRangeType', value);
        appSettings.value.statistics.defaultTrendChartDataRangeType = value;
        updateUserApplicationCloudSettingValue('statistics.defaultTrendChartDataRangeType', value);
    }

    function clearAppSettings(): void {
        clearSettings();
        appSettings.value = getApplicationSettings();
    }

    function createApplicationCloudSettings(applicationSettingKeys: string[]): ApplicationCloudSetting[] {
        if (!applicationSettingKeys || applicationSettingKeys.length < 1) {
            return [];
        }

        const settings: ApplicationCloudSetting[] = [];

        for (let i = 0; i < applicationSettingKeys.length; i++) {
            const settingKey = applicationSettingKeys[i];
            const cloudSetting = createUserApplicationCloudSetting(settingKey);

            if (cloudSetting) {
                settings.push(cloudSetting);
            }
        }

        return settings;
    }

    function setApplicationSettingsFromCloudSettings(cloudSettings?: ApplicationCloudSetting[]): void {
        if (!cloudSettings || cloudSettings.length < 1) {
            syncedAppSettings.value = {};
            return;
        }

        syncedAppSettings.value = arrayItemToObjectField(cloudSettings.map(item => item.settingKey), true);

        for (let i = 0; i < cloudSettings.length; i++) {
            const setting = cloudSettings[i];

            if (!setting || !setting.settingKey) {
                continue;
            }

            const settingType = ALL_ALLOWED_CLOUD_SYNC_APP_SETTING_KEY_TYPES[setting.settingKey];

            if (!settingType) {
                logger.warn(`cannot load application cloud setting "${setting.settingKey}", because it is not supported to sync`);
                continue;
            }

            if (settingType === UserApplicationCloudSettingType.String) {
                updateApplicationSettingsValueAndAppSettingsFromCloudSetting(setting.settingKey, setting.settingValue);
            } else if (settingType === UserApplicationCloudSettingType.Number) {
                const value = parseFloat(setting.settingValue);

                if (isNaN(value)) {
                    logger.warn(`cannot load application cloud setting "${setting.settingKey}", because it has invalid number value`);
                    continue;
                }

                updateApplicationSettingsValueAndAppSettingsFromCloudSetting(setting.settingKey, value);
            } else if (settingType === UserApplicationCloudSettingType.Boolean) {
                if (setting.settingValue !== 'true' && setting.settingValue !== 'false') {
                    logger.warn(`cannot load application cloud setting "${setting.settingKey}", because it has invalid boolean value`);
                    continue;
                }

                updateApplicationSettingsValueAndAppSettingsFromCloudSetting(setting.settingKey, setting.settingValue === 'true');
            } else if (settingType === UserApplicationCloudSettingType.StringBooleanMap) {
                try {
                    const map = JSON.parse(setting.settingValue);
                    let isValid = isObject(map);

                    if (isValid) {
                        for (const key in map) {
                            if (!Object.prototype.hasOwnProperty.call(map, key)) {
                                continue;
                            }

                            const value = map[key];

                            if (!isBoolean(value)) {
                                isValid = false;
                                break;
                            }
                        }
                    }

                    if (!isValid) {
                        logger.warn(`cannot load application cloud setting "${setting.settingKey}", because it has invalid map value`);
                        continue;
                    }

                    updateApplicationSettingsValueAndAppSettingsFromCloudSetting(setting.settingKey, map as Record<string, boolean>);
                } catch (error) {
                    logger.warn(`cannot load application cloud setting "${setting.settingKey}", because cannot parse JSON (${error})`);
                }
            } else {
                logger.warn(`cannot load application cloud setting "${setting.settingKey}", because it has unknown type "${settingType}"`);
            }
        }
    }

    function updateApplicationSyncSettingKeys(settingKeys?: string[]): void {
        if (!settingKeys || settingKeys.length < 1) {
            syncedAppSettings.value = {};
        } else {
            syncedAppSettings.value = arrayItemToObjectField(settingKeys, true);
        }
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
        syncedAppSettings,
        localeDefaultSettings,
        // computed states
        enableApplicationCloudSync,
        // functions
        // -- Basic Settings
        setTheme,
        setFontSize,
        setTimeZone,
        setAutoUpdateExchangeRatesData,
        setShowAccountBalance,
        setEnableAnimate,
        // -- Application Lock
        setEnableApplicationLock,
        setEnableApplicationLockWebAuthn,
        // -- Navigation Bar
        setShowAddTransactionButtonInDesktopNavbar,
        // -- Overview Page
        setShowAmountInHomePage,
        setTimezoneUsedForStatisticsInHomePage,
        // -- Transaction List Page
        setItemsCountInTransactionListPage,
        setShowTotalAmountInTransactionListPage,
        setShowTagInTransactionListPage,
        // -- Transaction Edit Page
        setAutoSaveTransactionDraft,
        setAutoGetCurrentGeoLocation,
        setAlwaysShowTransactionPicturesInMobileTransactionEditPage,
        // -- Exchange Rates Data Page
        setCurrencySortByInExchangeRatesPage,
        // -- Statistics Settings
        setStatisticsDefaultChartDataType,
        setStatisticsDefaultTimezoneType,
        setStatisticsDefaultAccountFilter,
        setStatisticsDefaultTransactionCategoryFilter,
        setStatisticsSortingType,
        setStatisticsDefaultCategoricalChartType,
        setStatisticsDefaultCategoricalChartDateRange,
        setStatisticsDefaultTrendChartType,
        setStatisticsDefaultTrendChartDateRange,
        clearAppSettings,
        createApplicationCloudSettings,
        setApplicationSettingsFromCloudSettings,
        updateApplicationSyncSettingKeys,
        updateLocalizedDefaultSettings
    };
});
