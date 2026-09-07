import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';

import type { BigDecimal, HiddenAmount, BigDecimalWithSuffix } from '@/core/numeral.ts';
import type { WeekDayValue } from '@/core/datetime.ts';
import { AccountCategory, AccountType } from '@/core/account.ts';
import { ACCOUNT_CURRENCY_NOT_SET_VALUE } from '@/consts/currency.ts';

import type { Account, CategorizedAccount } from '@/models/account.ts';

import { isDefined, isObject, isString } from '@/lib/common.ts';
import { isBigDecimal, parseBigDecimal } from '@/lib/numeral.ts';

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
        const netAssets: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getNetAssets(showAccountBalance.value, settingsStore.appSettings.totalAmountExcludeAccountIds);
        return formatAmountToLocalizedNumeralsWithCurrency(netAssets, defaultCurrency.value);
    });

    const totalAssets = computed<string>(() => {
        const totalAssets: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getTotalAssets(showAccountBalance.value, settingsStore.appSettings.totalAmountExcludeAccountIds);
        return formatAmountToLocalizedNumeralsWithCurrency(totalAssets, defaultCurrency.value);
    });

    const totalLiabilities = computed<string>(() => {
        const totalLiabilities: BigDecimal | HiddenAmount | BigDecimalWithSuffix = accountsStore.getTotalLiabilities(showAccountBalance.value, settingsStore.appSettings.totalAmountExcludeAccountIds);
        return formatAmountToLocalizedNumeralsWithCurrency(totalLiabilities, defaultCurrency.value);
    });

    function canShowAvailableCredit(account: Account): boolean {
        return account.category === AccountCategory.CreditCard.type && account.numericCreditCardLimit > 0 && !!account.currency && account.currency !== ACCOUNT_CURRENCY_NOT_SET_VALUE;
    }

    function accountCategoryTotalBalance(accountCategory: AccountCategory | undefined, showAvailableCreditForCreditCard: boolean): string {
        if (!accountCategory) {
            return '';
        }

        const totalBalance: BigDecimal | HiddenAmount | BigDecimalWithSuffix | undefined = accountsStore.getAccountCategoryTotalBalance(showAccountBalance.value, accountCategory, showAvailableCreditForCreditCard);

        if (isDefined(totalBalance)) {
            return formatAmountToLocalizedNumeralsWithCurrency(totalBalance, defaultCurrency.value);
        } else {
            return '';
        }
    }

    function accountBalance(account: Account, currentSubAccountId: string | undefined, showBalance: boolean, onlyShowSelectedAccountIds?: string[]): string | null {
        if (account.type === AccountType.SingleAccount.type) {
            const balance: BigDecimal | HiddenAmount | null = accountsStore.getAccountBalance(showBalance, account);

            if (isBigDecimal(balance) || isString(balance)) {
                return formatAmountToLocalizedNumeralsWithCurrency(balance, account.currency);
            } else {
                return '';
            }
        } else if (account.type === AccountType.MultiSubAccounts.type) {
            const balanceResult = accountsStore.getAccountSubAccountBalance(showBalance, showHidden.value, account, currentSubAccountId, onlyShowSelectedAccountIds);

            if (!isObject(balanceResult)) {
                return '';
            }

            return formatAmountToLocalizedNumeralsWithCurrency(balanceResult.balance, balanceResult.currency);
        } else {
            return null;
        }
    }

    function accountAvailableCredit(account: Account, showBalance: boolean): string | null {
        if (!canShowAvailableCredit(account)) {
            return null;
        }

        const creditCardLimit = parseBigDecimal(account.creditCardLimit);

        if (account.type === AccountType.SingleAccount.type) {
            const balance: BigDecimal | HiddenAmount | null = accountsStore.getAccountBalance(showBalance, account);

            if (isBigDecimal(balance)) {
                return formatAmountToLocalizedNumeralsWithCurrency(creditCardLimit.subtract(balance), account.currency);
            } else if (isString(balance)) {
                return formatAmountToLocalizedNumeralsWithCurrency(balance, account.currency);
            } else {
                return null;
            }
        } else if (account.type === AccountType.MultiSubAccounts.type) {
            const balanceResult = accountsStore.getAccountSubAccountBalance(showBalance, showHidden.value, account, undefined, undefined);

            if (!isObject(balanceResult)) {
                return null;
            }

            const balance: BigDecimal | HiddenAmount | BigDecimalWithSuffix = balanceResult.balance;

            if (isBigDecimal(balance)) {
                return formatAmountToLocalizedNumeralsWithCurrency(creditCardLimit.subtract(balance), balanceResult.currency);
            } else if (isObject(balance) && !balance.suffix) {
                return formatAmountToLocalizedNumeralsWithCurrency(creditCardLimit.subtract(balance.value), balanceResult.currency);
            } else if (isString(balance)) {
                return formatAmountToLocalizedNumeralsWithCurrency(balance, balanceResult.currency);
            } else {
                return null;
            }
        } else {
            return null;
        }
    }

    function accountBalanceOrAvailableCredit(account: Account, currentSubAccountId: string | undefined, showAvailableCreditForCreditCard: boolean, showBalance: boolean, onlyShowSelectedAccountIds?: string[]): string | null {
        if (showAvailableCreditForCreditCard && account.category === AccountCategory.CreditCard.type) {
            return accountAvailableCredit(account, showBalance);
        } else {
            return accountBalance(account, currentSubAccountId, showBalance, onlyShowSelectedAccountIds);
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
        canShowAvailableCredit,
        accountCategoryTotalBalance,
        accountBalance,
        accountAvailableCredit,
        accountBalanceOrAvailableCredit
    };
}
