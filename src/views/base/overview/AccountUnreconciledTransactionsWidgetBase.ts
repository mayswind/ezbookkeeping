import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { type TimeRangeAndDateType, DateRange } from '@/core/datetime.ts';

import type { Account } from '@/models/account.ts';

import { isArray, isNumber, arrayItemToObjectField } from '@/lib/common.ts';
import { getDateRangeByLastReconciledTimeRangeDateType } from '@/lib/datetime.ts';

export interface AccountAndUnreconciledTransactions {
    account: Account;
    count: number;
}

export interface CommonAccountUnreconciledTransactionsWidgetProps {
    loading: boolean;
    title?: string;
    accountIds: string[];
    itemCount: number;
    sortBy: string;
}

export function useAccountUnreconciledTransactionsWidgetBase(props: CommonAccountUnreconciledTransactionsWidgetProps) {
    const { tt } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const overviewStore = useOverviewStore();

    const noDataText = computed<string>(() => {
        if (props.loading) {
            return '';
        }

        if (!userStore.currentUserUseLastReconciledTime) {
            return tt('Last reconciled time is not enabled');
        }

        return tt('No data');
    });

    const unreconciledCounts = computed<Record<string, number>>(() => {
        const counts: Record<string, number> = {};

        for (const item of overviewStore.transactionUnreconciledCounts) {
            counts[item.accountId] = item.count;
        }

        return counts;
    });

    const displayAccounts = computed<AccountAndUnreconciledTransactions[]>(() => {
        const selectedAccountIds: Record<string, boolean> = isArray(props.accountIds) ? arrayItemToObjectField(props.accountIds, true) : {};
        const accounts: AccountAndUnreconciledTransactions[] = [];

        for (const account of accountsStore.allVisiblePlainAccounts) {
            if (account.hidden) {
                continue;
            }

            if (props.accountIds && props.accountIds.length && !selectedAccountIds[account.id]) {
                continue;
            }

            const count = unreconciledCounts.value[account.id];

            if (!isNumber(count) || count <= 0) {
                continue;
            }

            accounts.push({
                account: account,
                count: count
            });
        }

        if (props.sortBy === 'transactionCount') {
            accounts.sort((a, b) => {
                const aCount = unreconciledCounts.value[a.account.id] ?? 0;
                const bCount = unreconciledCounts.value[b.account.id] ?? 0;

                if (aCount === bCount) {
                    return a.account.compareTo(b.account, settingsStore.accountCategoryDisplayOrders, accountsStore.allAccountsMap);
                }

                return bCount - aCount;
            });
        }

        return accounts.slice(0, props.itemCount);
    });

    function getTransactionListDateRange(accountId: string): TimeRangeAndDateType | null {
        const account = accountsStore.allAccountsMap[accountId];

        if (!account) {
            return null;
        }

        return getDateRangeByLastReconciledTimeRangeDateType(DateRange.SinceLastReconciledTime.type, account.lastReconciledTime);
    }

    return {
        // computed states
        noDataText,
        displayAccounts,
        // functions
        getTransactionListDateRange
    };
}
