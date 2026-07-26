import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import type { BigDecimal, HiddenAmount, BigDecimalWithSuffix } from '@/core/numeral.ts';
import type { WeekDayValue } from '@/core/datetime.ts';
import { AccountCategory, AccountType } from '@/core/account.ts';
import type { Account, CategorizedAccount } from '@/models/account.ts';

import { isObject, isString } from '@/lib/common.ts';
import { isBigDecimal } from '@/lib/numeral.ts';

export function useAccountListPageBase() {
    const { formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

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

    const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);
    const defaultAccountCategory = computed<AccountCategory>(() => AccountCategory.values(customAccountCategoryOrder.value)[0] ?? AccountCategory.Default);

    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
    const useLastReconciledTime = computed(() => userStore.currentUserUseLastReconciledTime);

    const allAccounts = computed<Account[]>(() => accountsStore.allAccounts);
    const allCategorizedAccountsMap = computed<Record<number, CategorizedAccount>>(() => accountsStore.allCategorizedAccountsMap);
    const allAccountCount = computed<number>(() => accountsStore.allAvailableAccountsCount);
    const maxCategoryAccountCount = computed<number>(() => accountsStore.maxCategoryAccountCount);

    const netAssets = computed<string>(() => {
        const netAssets: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getNetAssets(showAccountBalance.value);
        return formatAmountToLocalizedNumeralsWithCurrency(netAssets, defaultCurrency.value);
    });

    const totalAssets = computed<string>(() => {
        const totalAssets: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getTotalAssets(showAccountBalance.value);
        return formatAmountToLocalizedNumeralsWithCurrency(totalAssets, defaultCurrency.value);
    });

    const totalLiabilities = computed<string>(() => {
        const totalLiabilities: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getTotalLiabilities(showAccountBalance.value);
        return formatAmountToLocalizedNumeralsWithCurrency(totalLiabilities, defaultCurrency.value);
    });

    function accountCategoryTotalBalance(accountCategory?: AccountCategory): string {
        if (!accountCategory) {
            return '';
        }

        const totalBalance: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getAccountCategoryTotalBalance(showAccountBalance.value, accountCategory);
        return formatAmountToLocalizedNumeralsWithCurrency(totalBalance, defaultCurrency.value);
    }

    function accountBalance(account: Account, currentSubAccountId?: string): string | null {
        if (account.type === AccountType.SingleAccount.type) {
            const balance: BigDecimal | HiddenAmount | null = accountsStore.getAccountBalance(showAccountBalance.value, account);

            if (isBigDecimal(balance) || isString(balance)) {
                return formatAmountToLocalizedNumeralsWithCurrency(balance, account.currency);
            } else {
                return '';
            }
        } else if (account.type === AccountType.MultiSubAccounts.type) {
            const balanceResult = accountsStore.getAccountSubAccountBalance(showAccountBalance.value, showHidden.value, account, currentSubAccountId);

            if (!isObject(balanceResult)) {
                return '';
            }

            return formatAmountToLocalizedNumeralsWithCurrency(balanceResult.balance, balanceResult.currency);
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
        customAccountCategoryOrder,
        defaultAccountCategory,
        firstDayOfWeek,
        fiscalYearStart,
        defaultCurrency,
        useLastReconciledTime,
        allAccounts,
        allCategorizedAccountsMap,
        allAccountCount,
        maxCategoryAccountCount,
        netAssets,
        totalAssets,
        totalLiabilities,
        // functions
        accountCategoryTotalBalance,
        accountBalance
    };
}
