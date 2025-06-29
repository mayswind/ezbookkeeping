import { ref, computed } from 'vue';

import { useSettingsStore } from '@/stores/setting.ts';

import type { ApplicationCloudSetting } from '@/core/setting.ts';

export interface CategorizedApplicationCloudSettingItems {
    readonly categoryName: string;
    readonly categorySubName?: string;
    readonly items: ApplicationCloudSettingItem[];
}

export interface ApplicationCloudSettingItem {
    readonly settingKey: string;
    readonly settingName: string;
    readonly mobile: boolean;
    readonly desktop: boolean;
}

export const ALL_APPLICATION_CLOUD_SETTINGS: CategorizedApplicationCloudSettingItems[] = [
    {
        categoryName: 'Basic Settings',
        items: [
            { settingKey: 'showAccountBalance', settingName: 'Show Account Balance', mobile: true, desktop: true }
        ]
    },
    {
        categoryName: 'Overview Page',
        items: [
            { settingKey: 'showAmountInHomePage', settingName: 'Show Amount', mobile: true, desktop: true },
            { settingKey: 'timezoneUsedForStatisticsInHomePage', settingName: 'Timezone Used for Statistics', mobile: true, desktop: true }
        ]
    },
    {
        categoryName: 'Transaction List Page',
        items: [
            { settingKey: 'itemsCountInTransactionListPage', settingName: 'Transactions Per Page', mobile: false, desktop: true },
            { settingKey: 'showTotalAmountInTransactionListPage', settingName: 'Show Monthly Total Amount', mobile: true, desktop: true },
            { settingKey: 'showTagInTransactionListPage', settingName: 'Show Transaction Tag', mobile: true, desktop: true }
        ]
    },
    {
        categoryName: 'Transaction Edit Page',
        items: [
            { settingKey: 'autoSaveTransactionDraft', settingName: 'Automatically Save Draft', mobile: true, desktop: true },
            { settingKey: 'autoGetCurrentGeoLocation', settingName: 'Automatically Add Geolocation', mobile: true, desktop: true },
            { settingKey: 'alwaysShowTransactionPicturesInMobileTransactionEditPage', settingName: 'Always Show Transaction Pictures', mobile: true, desktop: false }
        ]
    },
    {
        categoryName: 'Account List Page',
        items: [
            { settingKey: 'totalAmountExcludeAccountIds', settingName: 'Accounts Included in Total', mobile: true, desktop: true },
        ]
    },
    {
        categoryName: 'Exchange Rates Data Page',
        items: [
            { settingKey: 'currencySortByInExchangeRatesPage', settingName: 'Sort by', mobile: true, desktop: true }
        ]
    },
    {
        categoryName: 'Statistics Settings',
        categorySubName: 'Common Settings',
        items: [
            { settingKey: 'statistics.defaultChartDataType', settingName: 'Default Chart Data Type', mobile: true, desktop: true },
            { settingKey: 'statistics.defaultTimezoneType', settingName: 'Timezone Used for Date Range', mobile: true, desktop: true },
            { settingKey: 'statistics.defaultAccountFilter', settingName: 'Default Account Filter', mobile: true, desktop: true },
            { settingKey: 'statistics.defaultTransactionCategoryFilter', settingName: 'Default Transaction Category Filter', mobile: true, desktop: true },
            { settingKey: 'statistics.defaultSortingType', settingName: 'Default Sort Order', mobile: true, desktop: true }
        ]
    },
    {
        categoryName: 'Statistics Settings',
        categorySubName: 'Categorical Analysis Settings',
        items: [
            { settingKey: 'statistics.defaultCategoricalChartType', settingName: 'Default Chart Type', mobile: true, desktop: true },
            { settingKey: 'statistics.defaultCategoricalChartDataRangeType', settingName: 'Default Date Range', mobile: true, desktop: true }
        ]
    },
    {
        categoryName: 'Statistics Settings',
        categorySubName: 'Trend Analysis Settings',
        items: [
            { settingKey: 'statistics.defaultTrendChartType', settingName: 'Default Chart Type', mobile: false, desktop: true },
            { settingKey: 'statistics.defaultTrendChartDataRangeType', settingName: 'Default Date Range', mobile: true, desktop: true }
        ]
    }
];

export function useAppCloudSyncBase() {
    const settingsStore = useSettingsStore();

    const loading = ref<boolean>(false);
    const enabling = ref<boolean>(false);
    const disabling = ref<boolean>(false);
    const enabledApplicationCloudSettings = ref<Record<string, boolean>>(Object.assign({}, settingsStore.syncedAppSettings));

    const isEnableCloudSync = computed<boolean>(() => settingsStore.enableApplicationCloudSync);

    const hasEnabledApplicationCloudSettings = computed<boolean>(() => {
        for (const key in enabledApplicationCloudSettings.value) {
            if (!Object.prototype.hasOwnProperty.call(enabledApplicationCloudSettings.value, key)) {
                continue;
            }

            if (enabledApplicationCloudSettings.value[key]) {
                return true;
            }
        }

        return false;
    });

    const enabledApplicationCloudSettingKeys = computed<string[]>(() => {
        const keys: string[] = [];

        for (const key in enabledApplicationCloudSettings.value) {
            if (!Object.prototype.hasOwnProperty.call(enabledApplicationCloudSettings.value, key)) {
                continue;
            }

            if (enabledApplicationCloudSettings.value[key]) {
                keys.push(key);
            }
        }

        return keys;
    });

    function isAllSettingsSelected(categorizedItems: CategorizedApplicationCloudSettingItems): boolean {
        for (let i = 0; i < categorizedItems.items.length; i++) {
            const item = categorizedItems.items[i];
            if (!enabledApplicationCloudSettings.value[item.settingKey]) {
                return false;
            }
        }

        return true;
    }

    function hasSettingSelectedButNotAllChecked(categorizedItems: CategorizedApplicationCloudSettingItems): boolean {
        let checkedCount = 0;

        for (let i = 0; i < categorizedItems.items.length; i++) {
            const item = categorizedItems.items[i];
            if (!enabledApplicationCloudSettings.value[item.settingKey]) {
                checkedCount++;
            }
        }

        return checkedCount > 0 && checkedCount < categorizedItems.items.length;
    }

    function updateSettingsSelected(categorizedItems: CategorizedApplicationCloudSettingItems, value: boolean): void {
        for (let i = 0; i < categorizedItems.items.length; i++) {
            const item = categorizedItems.items[i];
            enabledApplicationCloudSettings.value[item.settingKey] = value;
        }
    }

    function selectAllSettings(): void {
        for (let i = 0; i < ALL_APPLICATION_CLOUD_SETTINGS.length; i++) {
            const categorizedItems = ALL_APPLICATION_CLOUD_SETTINGS[i];

            for (let j = 0; j < categorizedItems.items.length; j++) {
                const item = categorizedItems.items[j];
                enabledApplicationCloudSettings.value[item.settingKey] = true;
            }
        }
    }

    function selectNoneSettings(): void {
        for (let i = 0; i < ALL_APPLICATION_CLOUD_SETTINGS.length; i++) {
            const categorizedItems = ALL_APPLICATION_CLOUD_SETTINGS[i];

            for (let j = 0; j < categorizedItems.items.length; j++) {
                const item = categorizedItems.items[j];
                enabledApplicationCloudSettings.value[item.settingKey] = false;
            }
        }
    }

    function selectInvertSettings(): void {
        for (let i = 0; i < ALL_APPLICATION_CLOUD_SETTINGS.length; i++) {
            const categorizedItems = ALL_APPLICATION_CLOUD_SETTINGS[i];

            for (let j = 0; j < categorizedItems.items.length; j++) {
                const item = categorizedItems.items[j];
                enabledApplicationCloudSettings.value[item.settingKey] = !enabledApplicationCloudSettings.value[item.settingKey];
            }
        }
    }

    function setUserApplicationCloudSettings(settings: ApplicationCloudSetting[] | false) {
        if (settings && settings.length > 0) {
            settingsStore.setApplicationSettingsFromCloudSettings(settings);

            for (let i = 0; i < settings.length; i++) {
                const setting = settings[i];
                if (setting && setting.settingKey) {
                    enabledApplicationCloudSettings.value[setting.settingKey] = true;
                }
            }
        } else {
            settingsStore.setApplicationSettingsFromCloudSettings(undefined);
            enabledApplicationCloudSettings.value = {};
        }
    }

    return {
        // constants
        ALL_APPLICATION_CLOUD_SETTINGS,
        // states
        loading,
        enabling,
        disabling,
        enabledApplicationCloudSettings,
        // computed states
        isEnableCloudSync,
        hasEnabledApplicationCloudSettings,
        enabledApplicationCloudSettingKeys,
        // functions
        isAllSettingsSelected,
        hasSettingSelectedButNotAllChecked,
        updateSettingsSelected,
        selectAllSettings,
        selectNoneSettings,
        selectInvertSettings,
        setUserApplicationCloudSettings
    };
}
