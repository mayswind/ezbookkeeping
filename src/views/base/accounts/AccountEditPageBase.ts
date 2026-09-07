import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { AccountCategory, AccountType } from '@/core/account.ts';
import type { LocalizedAccountCategory } from '@/core/account.ts';
import { ACCOUNT_CURRENCY_NOT_SET_VALUE } from '@/consts/currency.ts';

import { Account } from '@/models/account.ts';

import { isDefined, isNumber } from '@/lib/common.ts';
import { parseBigDecimal } from '@/lib/numeral.ts';
import {
    getTimezoneOffsetMinutes,
    getSameDateTimeWithCurrentTimezone,
    parseDateTimeFromUnixTimeWithBrowserTimezone,
    getCurrentUnixTime
} from '@/lib/datetime.ts';

export function useAccountEditPageBase() {
    const {
        tt,
        getAvailableMonthDays,
        getAllAccountCategories,
        getAllAccountTypes,
        formatAmountToLocalizedNumeralsWithCurrency
    } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();

    const defaultAccountCategory = AccountCategory.values(settingsStore.appSettings.accountCategoryOrders)[0] ?? AccountCategory.Default;

    const editAccountId = ref<string | null>(null);
    const clientSessionId = ref<string>('');
    const loading = ref<boolean>(false);
    const submitting = ref<boolean>(false);
    const account = ref<Account>(Account.createNewAccount(defaultAccountCategory, userStore.currentUserDefaultCurrency, getCurrentUnixTimeForNewAccount()));
    const subAccounts = ref<Account[]>([]);

    const useLastReconciledTime = computed(() => userStore.currentUserUseLastReconciledTime);

    const title = computed<string>(() => {
        if (!editAccountId.value) {
            return 'Add Account';
        } else {
            return 'Edit Account';
        }
    });

    const saveButtonTitle = computed<string>(() => {
        if (!editAccountId.value) {
            return 'Add';
        } else {
            return 'Save';
        }
    });

    const inputEmptyProblemMessage = computed<string | null>(() => {
        let problemMessage = getInputEmptyProblemMessage(account.value, false);

        if (problemMessage) {
            return problemMessage;
        }

        if (account.value.type === AccountType.MultiSubAccounts.type) {
            for (const subAccount of subAccounts.value) {
                problemMessage = getInputEmptyProblemMessage(subAccount, true);

                if (problemMessage) {
                    return problemMessage;
                }
            }
        }

        return null;
    });

    const inputIsEmpty = computed<boolean>(() => !!inputEmptyProblemMessage.value);

    const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);
    const allAccountCategories = computed<LocalizedAccountCategory[]>(() => getAllAccountCategories(customAccountCategoryOrder.value));
    const allAccountTypes = computed<TypeAndDisplayName[]>(() => getAllAccountTypes());

    const allAvailableMonthDays = computed<TypeAndDisplayName[]>(() => {
        const allAvailableDays: TypeAndDisplayName[] = getAvailableMonthDays(28);

        allAvailableDays.splice(0, 0, {
            type: 0,
            displayName: tt('Not set'),
        });

        return allAvailableDays;
    });

    function getCurrentUnixTimeForNewAccount(): number {
        return getSameDateTimeWithCurrentTimezone(parseDateTimeFromUnixTimeWithBrowserTimezone(getCurrentUnixTime())).getUnixTime();
    }

    function getDefaultTimezoneOffsetMinutes(unixTime?: number): number {
        if (!unixTime) {
            return getTimezoneOffsetMinutes(getCurrentUnixTime());
        }

        return getTimezoneOffsetMinutes(unixTime);
    }

    function getAccountCreditCardStatementDate(statementDate?: number): string | null {
        if (!isDefined(statementDate)) {
            return tt('Not set');
        }

        for (const item of allAvailableMonthDays.value) {
            if (item.type === statementDate) {
                return item.displayName;
            }
        }

        return null;
    }

    function getAccountCreditCardCreditLimitDisplayValue(creditLimit: number | undefined, currency: string): string {
        if (!isNumber(creditLimit) || creditLimit <= 0 || !currency || currency === ACCOUNT_CURRENCY_NOT_SET_VALUE) {
            return tt('Not set');
        }

        return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(creditLimit), currency);
    }

    function updateAccountBalanceTime(account: Account, balanceTime: number): void {
        if (!isDefined(account.balanceTime)) {
            account.balanceTime = balanceTime;
            return;
        }

        const oldUtcOffset = getTimezoneOffsetMinutes(account.balanceTime);
        const newUtcOffset = getTimezoneOffsetMinutes(balanceTime);

        if (oldUtcOffset === newUtcOffset) {
            account.balanceTime = balanceTime;
            return;
        }

        account.balanceTime = balanceTime - (newUtcOffset - oldUtcOffset) * 60;
    }

    function updateAccountLastReconciledTime(account: Account, lastReconciledTime: number): void {
        if (!isDefined(account.lastReconciledTime)) {
            account.lastReconciledTime = lastReconciledTime;
            return;
        }

        const oldUtcOffset = getTimezoneOffsetMinutes(account.lastReconciledTime);
        const newUtcOffset = getTimezoneOffsetMinutes(lastReconciledTime);

        if (oldUtcOffset === newUtcOffset) {
            account.lastReconciledTime = lastReconciledTime;
            return;
        }

        account.lastReconciledTime = lastReconciledTime - (newUtcOffset - oldUtcOffset) * 60;
    }

    function getInputEmptyProblemMessage(account: Account, isSubAccount: boolean): string | null {
        if (!isSubAccount && !account.category) {
            return 'Account category cannot be blank';
        } else if (!isSubAccount && !account.type) {
            return 'Account type cannot be blank';
        } else if (!account.name) {
            return 'Account name cannot be blank';
        } else if (account.type === AccountType.SingleAccount.type && !account.currency) {
            return 'Account currency cannot be blank';
        } else {
            return null;
        }
    }

    function isNewAccount(account: Account): boolean {
        return account.id === '' || account.id === '0';
    }

    function addSubAccount(): boolean {
        if (account.value.type !== AccountType.MultiSubAccounts.type) {
            return false;
        }

        const subAccount = account.value.createNewSubAccount(userStore.currentUserDefaultCurrency, getCurrentUnixTimeForNewAccount());
        subAccounts.value.push(subAccount);
        return true;
    }

    function setAccount(newAccount: Account): void {
        account.value.fillFrom(newAccount);
        subAccounts.value = [];

        if (newAccount.subAccounts && newAccount.subAccounts.length > 0) {
            for (const oldSubAccount of newAccount.subAccounts) {
                const subAccount: Account = account.value.createNewSubAccount(userStore.currentUserDefaultCurrency, getCurrentUnixTimeForNewAccount());
                subAccount.fillFrom(oldSubAccount);

                subAccounts.value.push(subAccount);
            }
        }
    }

    watch(() => account.value.category, (newValue, oldValue) => {
        account.value.setSuitableIcon(oldValue, newValue);
    });

    watch(() => account.value.currency, (newValue) => {
        if (account.value.category === AccountCategory.CreditCard.type && (!newValue || newValue === ACCOUNT_CURRENCY_NOT_SET_VALUE)) {
            account.value.numericCreditCardLimit = 0;
        }
    });

    return {
        // constants
        defaultAccountCategory,
        // states
        editAccountId,
        clientSessionId,
        loading,
        submitting,
        account,
        subAccounts,
        // computed states
        useLastReconciledTime,
        title,
        saveButtonTitle,
        inputEmptyProblemMessage,
        inputIsEmpty,
        allAccountCategories,
        allAccountTypes,
        allAvailableMonthDays,
        // functions
        getCurrentUnixTimeForNewAccount,
        getDefaultTimezoneOffsetMinutes,
        getAccountCreditCardStatementDate,
        getAccountCreditCardCreditLimitDisplayValue,
        updateAccountBalanceTime,
        updateAccountLastReconciledTime,
        isNewAccount,
        addSubAccount,
        setAccount
    };
}
