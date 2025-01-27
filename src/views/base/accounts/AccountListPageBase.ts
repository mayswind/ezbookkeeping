import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import type { AccountCategory } from '@/core/account.ts';
import type { Account, CategorizedAccount } from '@/models/account.ts';

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
        accountCategoryTotalBalance
    };
}
