import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { type NameValue, type TypeAndDisplayName, keysIfValueEquals, values } from '@/core/base.ts';
import type { LocalizedTimezoneInfo } from '@/core/timezone.ts';
import { CategoryType } from '@/core/category.ts';
import type { Account } from '@/models/account.ts';

import { isObjectEmpty } from '@/lib/common.ts';

export function useAppSettingPageBase() {
    const { tt, getAllTimezones, getAllTimezoneTypesUsedForStatistics, getAllCurrencySortingTypes, setTimeZone } = useI18n();

    const settingsStore = useSettingsStore();
    const accountsStore = useAccountsStore();
    const transactionsStore = useTransactionsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const overviewStore = useOverviewStore();
    const statisticsStore = useStatisticsStore();

    const loadingAccounts = ref<boolean>(false);
    const loadingTransactionCategories = ref<boolean>(false);

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

    const hasAnyAccount = computed<boolean>(() => accountsStore.allPlainAccounts.length > 0);
    const hasAnyVisibleAccount = computed<boolean>(() => accountsStore.allVisibleAccountsCount > 0);
    const hasAnyTransactionCategory = computed<boolean>(() => !isObjectEmpty(transactionCategoriesStore.allTransactionCategoriesMap));

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

    const accountsIncludedInHomePageOverviewDisplayContent = computed<string>(() => {
        const excludeAccountIds = settingsStore.appSettings.overviewAccountFilterInHomePage;
        return getIncludedAccountsDisplayContent(excludeAccountIds, accountsStore.allPlainAccounts);
    });

    const accountsIncludedInTotalDisplayContent = computed<string>(() => {
        const excludeAccountIds = settingsStore.appSettings.totalAmountExcludeAccountIds;
        return getIncludedAccountsDisplayContent(excludeAccountIds, accountsStore.allVisiblePlainAccounts);
    });

    const transactionCategoriesIncludedInHomePageOverviewDisplayContent = computed<string>(() => {
        const excludeAccountIds = settingsStore.appSettings.overviewTransactionCategoryFilterInHomePage;
        return getIncludedTransactionCategoriesDisplayContent(excludeAccountIds);
    });

    function getIncludedAccountsDisplayContent(excludeAccountIds: Record<string, boolean>, allAccounts: Account[]): string {
        if (loadingAccounts.value || !allAccounts || !allAccounts.length) {
            return '';
        }

        let hasExcludeAccount = false;

        for (const accountId of keysIfValueEquals(excludeAccountIds, true)) {
            if (accountsStore.allAccountsMap[accountId]) {
                hasExcludeAccount = true;
                break;
            }
        }

        if (!hasExcludeAccount) {
            return tt('All');
        }

        let allAccountExcluded = true;

        for (const account of allAccounts) {
            if (!excludeAccountIds[account.id]) {
                allAccountExcluded = false;
                break;
            }
        }

        if (allAccountExcluded) {
            return tt('None');
        }

        return tt('Partial');
    }

    function getIncludedTransactionCategoriesDisplayContent(excludeTransactionCategoryIds: Record<string, boolean>): string {
        if (loadingTransactionCategories.value || !transactionCategoriesStore.allTransactionCategoriesMap) {
            return '';
        }

        let hasExcludeTransactionCategory = false;

        for (const transactionCategoryId of keysIfValueEquals(excludeTransactionCategoryIds, true)) {
            if (transactionCategoriesStore.allTransactionCategoriesMap[transactionCategoryId]) {
                hasExcludeTransactionCategory = true;
                break;
            }
        }

        if (!hasExcludeTransactionCategory) {
            return tt('All');
        }

        let allTransactionCategoryExcluded = true;

        for (const transactionCategory of values(transactionCategoriesStore.allTransactionCategoriesMap)) {
            if (transactionCategory.type !== CategoryType.Income && transactionCategory.type !== CategoryType.Expense) {
                continue;
            }

            if (!excludeTransactionCategoryIds[transactionCategory.id]) {
                allTransactionCategoryExcluded = false;
                break;
            }
        }

        if (allTransactionCategoryExcluded) {
            return tt('None');
        }

        return tt('Partial');
    }

    return {
        // states
        loadingAccounts,
        loadingTransactionCategories,
        // computed states
        allThemes,
        allTimezones,
        allTimezoneTypesUsedForStatistics,
        allCurrencySortingTypes,
        allAutoSaveTransactionDraftTypes,
        timeZone,
        hasAnyAccount,
        hasAnyVisibleAccount,
        hasAnyTransactionCategory,
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
        accountsIncludedInHomePageOverviewDisplayContent,
        accountsIncludedInTotalDisplayContent,
        transactionCategoriesIncludedInHomePageOverviewDisplayContent
    };
}
