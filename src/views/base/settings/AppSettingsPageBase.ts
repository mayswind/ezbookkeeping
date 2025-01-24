import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
// @ts-expect-error the above file is migrating to ts
import { useTransactionsStore } from '@/stores/transaction.js';
import { useOverviewStore } from '@/stores/overview.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import type { LocalizedTimezoneInfo } from '@/core/timezone.ts';

export function useAppSettingPageBase() {
    const { getAllTimezones, getAllTimezoneTypesUsedForStatistics, getAllCurrencySortingTypes, setTimeZone } = useI18n();

    const settingsStore = useSettingsStore();
    const transactionsStore = useTransactionsStore();
    const overviewStore = useOverviewStore();
    const statisticsStore = useStatisticsStore();

    const allTimezones = computed<LocalizedTimezoneInfo[]>(() => getAllTimezones(true));
    const allTimezoneTypesUsedForStatistics = computed<TypeAndDisplayName[]>(() => getAllTimezoneTypesUsedForStatistics());
    const allCurrencySortingTypes = computed<TypeAndDisplayName[]>(() => getAllCurrencySortingTypes());

    const timeZone = computed<string>({
        get: () => settingsStore.appSettings.timeZone,
        set: (value) => {
            settingsStore.setTimeZone(value);
            setTimeZone(value);
            transactionsStore.updateTransactionListInvalidState(true);
            overviewStore.updateTransactionOverviewInvalidState(true);
            statisticsStore.updateTransactionStatisticsInvalidState(true);
        }
    });

    const isAutoUpdateExchangeRatesData = computed<boolean>({
        get: () => settingsStore.appSettings.autoUpdateExchangeRatesData,
        set: (value) => settingsStore.setAutoUpdateExchangeRatesData(value)
    });

    const showAccountBalance = computed<boolean>({
        get: () => settingsStore.appSettings.showAccountBalance,
        set: (value) => settingsStore.setShowAccountBalance(value)
    });

    const showAmountInHomePage = computed<boolean>({
        get: () => settingsStore.appSettings.showAmountInHomePage,
        set: (value) => settingsStore.setShowAmountInHomePage(value)
    });

    const timezoneUsedForStatisticsInHomePage = computed<number>({
        get: () => settingsStore.appSettings.timezoneUsedForStatisticsInHomePage,
        set: (value: number) => {
            settingsStore.setTimezoneUsedForStatisticsInHomePage(value);
            overviewStore.updateTransactionOverviewInvalidState(true);
        }
    });

    const showTotalAmountInTransactionListPage = computed<boolean>({
        get: () => settingsStore.appSettings.showTotalAmountInTransactionListPage,
        set: (value) => settingsStore.setShowTotalAmountInTransactionListPage(value)
    });

    const showTagInTransactionListPage = computed<boolean>({
        get: () => settingsStore.appSettings.showTagInTransactionListPage,
        set: (value) => settingsStore.setShowTagInTransactionListPage(value)
    });

    const itemsCountInTransactionListPage = computed<number>({
        get: () => settingsStore.appSettings.itemsCountInTransactionListPage,
        set: (value) => settingsStore.setItemsCountInTransactionListPage(value)
    });

    const autoSaveTransactionDraft = computed<string>({
        get: () => settingsStore.appSettings.autoSaveTransactionDraft,
        set: (value: string) => {
            settingsStore.setAutoSaveTransactionDraft(value);

            if (value === 'disabled') {
                transactionsStore.clearTransactionDraft();
            }
        }
    });

    const isAutoGetCurrentGeoLocation = computed<boolean>({
        get: () => settingsStore.appSettings.autoGetCurrentGeoLocation,
        set: (value) => settingsStore.setAutoGetCurrentGeoLocation(value)
    });

    const currencySortByInExchangeRatesPage = computed<number>({
        get: () => settingsStore.appSettings.currencySortByInExchangeRatesPage,
        set: (value: number) => settingsStore.setCurrencySortByInExchangeRatesPage(value)
    });

    return {
        // computed states
        allTimezones,
        allTimezoneTypesUsedForStatistics,
        allCurrencySortingTypes,
        timeZone,
        isAutoUpdateExchangeRatesData,
        showAccountBalance,
        showAmountInHomePage,
        itemsCountInTransactionListPage,
        timezoneUsedForStatisticsInHomePage,
        showTotalAmountInTransactionListPage,
        showTagInTransactionListPage,
        autoSaveTransactionDraft,
        isAutoGetCurrentGeoLocation,
        currencySortByInExchangeRatesPage
    };
}
