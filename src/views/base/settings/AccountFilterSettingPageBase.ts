import { ref, computed } from 'vue';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { keys, keysIfValueEquals, values } from '@/core/base.ts';
import type {Account, CategorizedAccount} from '@/models/account.ts';

import {
    filterCategorizedAccounts,
    selectAccountOrSubAccounts,
    isAccountOrSubAccountsAllChecked
} from '@/lib/account.ts';

export type AccountFilterType = 'statisticsDefault' | 'statisticsCurrent' | 'homePageOverview' | 'transactionListCurrent' | 'accountListTotalAmount';

export function useAccountFilterSettingPageBase(type?: AccountFilterType) {
    const settingsStore = useSettingsStore();
    const accountsStore = useAccountsStore();
    const transactionsStore = useTransactionsStore();
    const statisticsStore = useStatisticsStore();
    const overviewStore = useOverviewStore();

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const filterContent = ref<string>('');
    const filterAccountIds = ref<Record<string, boolean>>({});

    const title = computed<string>(() => {
        if (type === 'statisticsDefault') {
            return 'Default Account Filter';
        } else {
            return 'Filter Accounts';
        }
    });

    const applyText = computed<string>(() => {
        if (type === 'statisticsDefault') {
            return 'Save';
        } else {
            return 'Apply';
        }
    });

    const allowHiddenAccount = computed<boolean>(() => {
        return type === 'statisticsDefault' || type === 'statisticsCurrent' || type === 'homePageOverview' || type === 'transactionListCurrent';
    });

    const allCategorizedAccounts = computed<Record<number, CategorizedAccount>>(() => filterCategorizedAccounts(accountsStore.allCategorizedAccountsMap, filterContent.value, showHidden.value));
    const allVisibleAccountMap = computed<Record<string, Account>>(() => {
        const accountMap: Record<string, Account> = {};

        for (const accountCategory of values(allCategorizedAccounts.value)) {
            for (const account of accountCategory.accounts) {
                accountMap[account.id] = account;

                if (account.subAccounts) {
                    for (const subAccount of account.subAccounts) {
                        accountMap[subAccount.id] = subAccount;
                    }
                }
            }
        }

        return accountMap;
    });
    const hasAnyAvailableAccount = computed<boolean>(() => accountsStore.allAvailableAccountsCount > 0);
    const hasAnyVisibleAccount = computed<boolean>(() => {
        for (const accountCategory of values(allCategorizedAccounts.value)) {
            if (accountCategory.accounts.length > 0) {
                return true;
            }
        }

        return false;
    });

    function isAccountChecked(account: Account, filterAccountIds: Record<string, boolean>): boolean {
        return !filterAccountIds[account.id];
    }

    function loadFilterAccountIds(): boolean {
        const allAccountIds: Record<string, boolean> = {};

        for (const account of values(accountsStore.allAccountsMap)) {
            if (!allowHiddenAccount.value && account.hidden) {
                continue;
            }

            if (type === 'transactionListCurrent' && transactionsStore.allFilterAccountIdsCount > 0) {
                allAccountIds[account.id] = true;
            } else {
                allAccountIds[account.id] = false;
            }
        }

        if (type === 'statisticsDefault') {
            filterAccountIds.value = Object.assign(allAccountIds, settingsStore.appSettings.statistics.defaultAccountFilter);
            return true;
        } else if (type === 'statisticsCurrent') {
            filterAccountIds.value = Object.assign(allAccountIds, statisticsStore.transactionStatisticsFilter.filterAccountIds);
            return true;
        } else if (type === 'homePageOverview') {
            filterAccountIds.value = Object.assign(allAccountIds, settingsStore.appSettings.overviewAccountFilterInHomePage);
            return true;
        } else if (type === 'transactionListCurrent') {
            for (const accountId of keysIfValueEquals(transactionsStore.allFilterAccountIds, true)) {
                const account = accountsStore.allAccountsMap[accountId];

                if (account) {
                    selectAccountOrSubAccounts(allAccountIds, account, false);
                }
            }
            filterAccountIds.value = allAccountIds;
            return true;
        } else if (type === 'accountListTotalAmount') {
            filterAccountIds.value = Object.assign(allAccountIds, settingsStore.appSettings.totalAmountExcludeAccountIds);
            return true;
        } else {
            return false;
        }
    }

    function saveFilterAccountIds(): boolean {
        const filteredAccountIds: Record<string, boolean> = {};
        let isAllSelected = true;
        let finalAccountIds = '';
        let changed = true;

        for (const accountId of keys(filterAccountIds.value)) {
            const account = accountsStore.allAccountsMap[accountId];

            if (!account) {
                continue;
            }

            if (!allowHiddenAccount.value && account.hidden) {
                continue;
            }

            if (!isAccountOrSubAccountsAllChecked(account, filterAccountIds.value)) {
                filteredAccountIds[accountId] = true;
                isAllSelected = false;
            } else {
                if (finalAccountIds.length > 0) {
                    finalAccountIds += ',';
                }

                finalAccountIds += accountId;
            }
        }

        if (type === 'statisticsDefault') {
            settingsStore.setStatisticsDefaultAccountFilter(filteredAccountIds);
        } else if (type === 'statisticsCurrent') {
            changed = statisticsStore.updateTransactionStatisticsFilter({
                filterAccountIds: filteredAccountIds
            });
        } else if (type === 'homePageOverview') {
            settingsStore.setOverviewAccountFilterInHomePage(filteredAccountIds);
            overviewStore.updateTransactionOverviewInvalidState(true);
        } else if (type === 'transactionListCurrent') {
            changed = transactionsStore.updateTransactionListFilter({
                accountIds: isAllSelected ? '' : finalAccountIds
            });

            if (changed) {
                transactionsStore.updateTransactionListInvalidState(true);
            }
        } else if (type === 'accountListTotalAmount') {
            settingsStore.setTotalAmountExcludeAccountIds(filteredAccountIds);
        }

        return changed;
    }

    return {
        // states
        loading,
        showHidden,
        filterContent,
        filterAccountIds,
        // computed states
        title,
        applyText,
        allowHiddenAccount,
        allCategorizedAccounts,
        allVisibleAccountMap,
        hasAnyAvailableAccount,
        hasAnyVisibleAccount,
        // functions
        isAccountChecked,
        loadFilterAccountIds,
        saveFilterAccountIds
    };
}
