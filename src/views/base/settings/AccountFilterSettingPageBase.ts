import { ref, computed } from 'vue';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import type { Account, AccountCategoriesWithVisibleCount } from '@/models/account.ts';

import {
    getCategorizedAccountsWithVisibleCount,
    selectAccountOrSubAccounts,
    isAccountOrSubAccountsAllChecked
} from '@/lib/account.ts';

export function useAccountFilterSettingPageBase(type?: string) {
    const settingsStore = useSettingsStore();
    const accountsStore = useAccountsStore();
    const transactionsStore = useTransactionsStore();
    const statisticsStore = useStatisticsStore();

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
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
        return type === 'statisticsDefault' || type === 'statisticsCurrent' || type === 'transactionListCurrent';
    });

    const allCategorizedAccounts = computed<AccountCategoriesWithVisibleCount[]>(() => getCategorizedAccountsWithVisibleCount(accountsStore.allCategorizedAccountsMap));
    const hasAnyAvailableAccount = computed<boolean>(() => accountsStore.allAvailableAccountsCount > 0);

    const hasAnyVisibleAccount = computed<boolean>(() => {
        if (showHidden.value) {
            return accountsStore.allAvailableAccountsCount > 0;
        } else {
            return accountsStore.allVisibleAccountsCount > 0;
        }
    });

    function isAccountChecked(account: Account, filterAccountIds: Record<string, boolean>): boolean {
        return !filterAccountIds[account.id];
    }

    function loadFilterAccountIds(): boolean {
        const allAccountIds: Record<string, boolean> = {};

        for (const accountId in accountsStore.allAccountsMap) {
            if (!Object.prototype.hasOwnProperty.call(accountsStore.allAccountsMap, accountId)) {
                continue;
            }

            const account = accountsStore.allAccountsMap[accountId];

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
        } else if (type === 'transactionListCurrent') {
            for (const accountId in transactionsStore.allFilterAccountIds) {
                if (!Object.prototype.hasOwnProperty.call(transactionsStore.allFilterAccountIds, accountId)) {
                    continue;
                }

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

        for (const accountId in filterAccountIds.value) {
            if (!Object.prototype.hasOwnProperty.call(filterAccountIds.value, accountId)) {
                continue;
            }

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
        filterAccountIds,
        // computed states
        title,
        applyText,
        allowHiddenAccount,
        allCategorizedAccounts,
        hasAnyAvailableAccount,
        hasAnyVisibleAccount,
        // functions
        isAccountChecked,
        loadFilterAccountIds,
        saveFilterAccountIds
    };
}
