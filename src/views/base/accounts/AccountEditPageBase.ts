import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { AccountCategory, AccountType } from '@/core/account.ts';
import type { LocalizedAccountCategory } from '@/core/account.ts';
import { Account } from '@/models/account.ts';

import { isDefined } from '@/lib/common.ts';
import {
    getTimezoneOffsetMinutes,
    getSameDateTimeWithCurrentTimezone,
    parseDateTimeFromUnixTimeWithBrowserTimezone,
    getCurrentUnixTime
} from '@/lib/datetime.ts';

export interface DayAndDisplayName {
    readonly day: number;
    readonly displayName: string;
}

export function useAccountEditPageBase() {
    const { tt, getAllAccountCategories, getAllAccountTypes, getMonthdayShortName } = useI18n();

    const userStore = useUserStore();

    const editAccountId = ref<string | null>(null);
    const clientSessionId = ref<string>('');
    const loading = ref<boolean>(false);
    const submitting = ref<boolean>(false);
    const account = ref<Account>(Account.createNewAccount(userStore.currentUserDefaultCurrency, getCurrentUnixTimeForNewAccount()));
    const subAccounts = ref<Account[]>([]);

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

    const allAccountCategories = computed<LocalizedAccountCategory[]>(() => getAllAccountCategories());
    const allAccountTypes = computed<TypeAndDisplayName[]>(() => getAllAccountTypes());

    const allAvailableMonthDays = computed<DayAndDisplayName[]>(() => {
        const allAvailableDays: DayAndDisplayName[] = [];

        allAvailableDays.push({
            day: 0,
            displayName: tt('Not set'),
        });

        for (let i = 1; i <= 28; i++) {
            allAvailableDays.push({
                day: i,
                displayName: getMonthdayShortName(i),
            });
        }

        return allAvailableDays;
    });

    const isAccountSupportCreditCardStatementDate = computed<boolean>(() => account.value && account.value.category === AccountCategory.CreditCard.type);

    function getCurrentUnixTimeForNewAccount(): number {
        return getSameDateTimeWithCurrentTimezone(parseDateTimeFromUnixTimeWithBrowserTimezone(getCurrentUnixTime())).getUnixTime();
    }

    function getDefaultTimezoneOffsetMinutes(account: Account): number {
        if (!account.balanceTime) {
            return 0;
        }

        return getTimezoneOffsetMinutes(account.balanceTime);
    }

    function getAccountCreditCardStatementDate(statementDate?: number): string | null {
        for (const item of allAvailableMonthDays.value) {
            if (item.day === statementDate) {
                return item.displayName;
            }
        }

        return null;
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

    return {
        // states
        editAccountId,
        clientSessionId,
        loading,
        submitting,
        account,
        subAccounts,
        // computed states
        title,
        saveButtonTitle,
        inputEmptyProblemMessage,
        inputIsEmpty,
        allAccountCategories,
        allAccountTypes,
        allAvailableMonthDays,
        isAccountSupportCreditCardStatementDate,
        // functions
        getCurrentUnixTimeForNewAccount,
        getDefaultTimezoneOffsetMinutes,
        getAccountCreditCardStatementDate,
        updateAccountBalanceTime,
        isNewAccount,
        addSubAccount,
        setAccount
    };
}
