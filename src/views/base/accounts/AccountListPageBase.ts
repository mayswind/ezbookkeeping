import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import { type AccountCategory, AccountType } from '@/core/account.ts';
import type { Account, CategorizedAccount } from '@/models/account.ts';

import { isObject, isString } from '@/lib/common.ts';

export function useAccountListPageBaseBase() {
    const { formatAmountWithCurrency } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const displayOrderModified = ref<boolean>(false);

    const showAccountBalance = computed<boolean>({
        get: () => settingsStore.appSettings.showAccountBalance,
        set: (value) => settingsStore.setShowAccountBalance(value)
    });

    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

    const allAccounts = computed<Account[]>(() => accountsStore.allAccounts);
    const allCategorizedAccountsMap = computed<Record<number, CategorizedAccount>>(() => accountsStore.allCategorizedAccountsMap);
    const allAccountCount = computed<number>(() => accountsStore.allAvailableAccountsCount);

    const netAssets = computed<string>(() => {
        const netAssets = accountsStore.getNetAssets(showAccountBalance.value);
        return formatAmountWithCurrency(netAssets, defaultCurrency.value);
    });

    const totalAssets = computed<string>(() => {
        const totalAssets = accountsStore.getTotalAssets(showAccountBalance.value);
        return formatAmountWithCurrency(totalAssets, defaultCurrency.value);
    });

    const totalLiabilities = computed<string>(() => {
        const totalLiabilities = accountsStore.getTotalLiabilities(showAccountBalance.value);
        return formatAmountWithCurrency(totalLiabilities, defaultCurrency.value);
    });

    function accountCategoryTotalBalance(accountCategory?: AccountCategory): string {
        if (!accountCategory) {
            return '';
        }

        const totalBalance = accountsStore.getAccountCategoryTotalBalance(showAccountBalance.value, accountCategory);
        return formatAmountWithCurrency(totalBalance, defaultCurrency.value);
    }

    function accountBalance(account: Account, currentSubAccountId?: string): string | null {
        if (account.type === AccountType.SingleAccount.type) {
            const balance = accountsStore.getAccountBalance(showAccountBalance.value, account);

            if (!isString(balance)) {
                return '';
            }

            return formatAmountWithCurrency(balance, account.currency);
        } else if (account.type === AccountType.MultiSubAccounts.type) {
            const balanceResult = accountsStore.getAccountSubAccountBalance(showAccountBalance.value, showHidden.value, account, currentSubAccountId);

            if (!isObject(balanceResult)) {
                return '';
            }

            return formatAmountWithCurrency(balanceResult.balance, balanceResult.currency);
        } else {
            return null;
        }
    }

    return {
        // states
        loading,
        showHidden,
        displayOrderModified,
        // computed states
        showAccountBalance,
        defaultCurrency,
        allAccounts,
        allCategorizedAccountsMap,
        allAccountCount,
        netAssets,
        totalAssets,
        totalLiabilities,
        // functions
        accountCategoryTotalBalance,
        accountBalance
    };
}
