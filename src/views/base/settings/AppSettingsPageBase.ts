import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import type { NameValue, TypeAndDisplayName } from '@/core/base.ts';
import type { LocalizedTimezoneInfo } from '@/core/timezone.ts';

export function useAppSettingPageBase() {
    const { tt, getAllTimezones, getAllTimezoneTypesUsedForStatistics, getAllCurrencySortingTypes, setTimeZone } = useI18n();

    const settingsStore = useSettingsStore();
    const accountsStore = useAccountsStore();
    const transactionsStore = useTransactionsStore();
    const overviewStore = useOverviewStore();
    const statisticsStore = useStatisticsStore();

    const loadingAccounts = ref<boolean>(false);

    const allThemes = computed<NameValue[]>(() => {
        return [
            { name: tt('System Default'), value: 'auto' },
            { name: tt('Light'), value: 'light' },
            { name: tt('Dark'), value: 'dark' }
        ];
    });

    const allTimezones = computed<LocalizedTimezoneInfo[]>(() => getAllTimezones(true));
    const allTimezoneTypesUsedForStatistics = computed<TypeAndDisplayName[]>(() => getAllTimezoneTypesUsedForStatistics());
    const allCurrencySortingTypes = computed<TypeAndDisplayName[]>(() => getAllCurrencySortingTypes());

    const allAutoSaveTransactionDraftTypes = computed<NameValue[]>(() => {
        return [
            { name: tt('Disabled'), value: 'disabled' },
            { name: tt('Enabled'), value: 'enabled' },
            { name: tt('Show Confirmation Every Time'), value: 'confirmation' }
        ];
    });

    const hasAnyVisibleAccount = computed<boolean>(() => accountsStore.allVisibleAccountsCount > 0);

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

    const accountsIncludedInTotalDisplayContent = computed<string>(() => {
        if (loadingAccounts.value || !accountsStore.allVisiblePlainAccounts || !accountsStore.allVisiblePlainAccounts.length) {
            return '';
        }

        const excludeAccountIds = settingsStore.appSettings.totalAmountExcludeAccountIds;
        let hasExcludeAccount = false;

        for (const accountId in excludeAccountIds) {
            if (!Object.prototype.hasOwnProperty.call(excludeAccountIds, accountId)) {
                continue;
            }

            if (excludeAccountIds[accountId]) {
                hasExcludeAccount = true;
                break;
            }
        }

        if (!hasExcludeAccount) {
            return tt('All');
        }

        let allVisibleAccountExcluded = true;

        for (let i = 0; i < accountsStore.allVisiblePlainAccounts.length; i++) {
            const account = accountsStore.allVisiblePlainAccounts[i];

            if (!excludeAccountIds[account.id]) {
                allVisibleAccountExcluded = false;
                break;
            }
        }

        if (allVisibleAccountExcluded) {
            return tt('None');
        }

        return tt('Partial');
    });

    return {
        // states
        loadingAccounts,
        // computed states
        allThemes,
        allTimezones,
        allTimezoneTypesUsedForStatistics,
        allCurrencySortingTypes,
        allAutoSaveTransactionDraftTypes,
        timeZone,
        hasAnyVisibleAccount,
        isAutoUpdateExchangeRatesData,
        showAccountBalance,
        showAmountInHomePage,
        itemsCountInTransactionListPage,
        timezoneUsedForStatisticsInHomePage,
        showTotalAmountInTransactionListPage,
        showTagInTransactionListPage,
        autoSaveTransactionDraft,
        isAutoGetCurrentGeoLocation,
        currencySortByInExchangeRatesPage,
        accountsIncludedInTotalDisplayContent
    };
}
