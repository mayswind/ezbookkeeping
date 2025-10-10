import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';

import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';

export function useMoveAllTransactionsPageBase() {
    const { tt, getCategorizedAccountsWithDisplayBalance } = useI18n();

    const settingsStore = useSettingsStore();
    const accountsStore = useAccountsStore();

    const moving = ref<boolean>(false);
    const fromAccount = ref<Account | undefined>(undefined);
    const toAccountId = ref<string>('');
    const toAccountName = ref<string>('');

    const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
    const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
    const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
    const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value));

    const displayToAccountName = computed<string>(() => {
        if (!toAccountId.value) {
            return tt('Target Account');
        }

        return Account.findAccountNameById(allAccounts.value, toAccountId.value, tt('Target Account')) || tt('Target Account');
    });

    const isToAccountNameValid = computed<boolean>(() => {
        if (!toAccountId.value || !toAccountName.value) {
            return false;
        }

        const expectedAccountName = Account.findAccountNameById(allAccounts.value, toAccountId.value);

        if (!expectedAccountName) {
            return false;
        }

        return expectedAccountName === toAccountName.value;
    });

    return {
        // states
        moving,
        fromAccount,
        toAccountId,
        toAccountName,
        // computed states
        allAccounts,
        allVisibleAccounts,
        allVisibleCategorizedAccounts,
        displayToAccountName,
        isToAccountNameValid
    };
}
