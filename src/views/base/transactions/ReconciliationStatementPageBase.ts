import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { TransactionType } from '@/core/transaction.ts';
import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionReconciliationStatementResponseItem } from '@/models/transaction.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    parseDateFromUnixTime,
    getUnixTime
} from '@/lib/datetime.ts';

export function useReconciliationStatementPageBase() {
    const {
        formatUnixTimeToLongDateTime,
        formatAmountWithCurrency
    } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();

    const accountId = ref<string>('');
    const startTime = ref<number>(0);
    const endTime = ref<number>(0);
    const reconciliationStatements = ref<TransactionReconciliationStatementResponseItem[]>([]);

    const currentTimezoneOffsetMinutes = computed<number>(() => getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone));
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

    const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
    const allCategoriesMap = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);

    const displayStartDateTime = computed<string>(() => {
        return formatUnixTimeToLongDateTime(startTime.value);
    });

    const displayEndDateTime = computed<string>(() => {
        return formatUnixTimeToLongDateTime(endTime.value);
    });

    const displayTotalOutflows = computed<string>(() => {
        let totalOutflows = 0;

        for (let i = 0; i < reconciliationStatements.value.length; i++) {
            const transaction = reconciliationStatements.value[i];

            if (transaction.type === TransactionType.Expense) {
                totalOutflows += transaction.sourceAmount;
            } else if (transaction.type === TransactionType.Transfer && transaction.sourceAccountId === accountId.value) {
                totalOutflows += transaction.sourceAmount;
            }
        }

        let currency = defaultCurrency.value;

        if (allAccountsMap.value[accountId.value]) {
            currency = allAccountsMap.value[accountId.value].currency;
        }

        return formatAmountWithCurrency(totalOutflows, currency);
    });

    const displayTotalInflows = computed<string>(() => {
        let totalInflows = 0;

        for (let i = 0; i < reconciliationStatements.value.length; i++) {
            const transaction = reconciliationStatements.value[i];

            if (transaction.type === TransactionType.Income) {
                totalInflows += transaction.sourceAmount;
            } else if (transaction.type === TransactionType.Transfer && transaction.destinationAccountId === accountId.value) {
                totalInflows += transaction.destinationAmount;
            }
        }

        let currency = defaultCurrency.value;

        if (allAccountsMap.value[accountId.value]) {
            currency = allAccountsMap.value[accountId.value].currency;
        }

        return formatAmountWithCurrency(totalInflows, currency);
    });

    function getDisplayDateTime(transaction: TransactionReconciliationStatementResponseItem): string {
        const transactionTime = getUnixTime(parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentTimezoneOffsetMinutes.value));
        return formatUnixTimeToLongDateTime(transactionTime);
    }

    function getDisplayTimezone(transaction: TransactionReconciliationStatementResponseItem): string {
        return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
    }

    function getDisplaySourceAmount(transaction: TransactionReconciliationStatementResponseItem): string {
        let currency = defaultCurrency.value;

        if (allAccountsMap.value[transaction.sourceAccountId]) {
            currency = allAccountsMap.value[transaction.sourceAccountId].currency;
        }

        return formatAmountWithCurrency(transaction.sourceAmount, currency);
    }

    function getDisplayDestinationAmount(transaction: TransactionReconciliationStatementResponseItem): string {
        let currency = defaultCurrency.value;

        if (allAccountsMap.value[transaction.destinationAccountId]) {
            currency = allAccountsMap.value[transaction.destinationAccountId].currency;
        }

        return formatAmountWithCurrency(transaction.destinationAmount, currency);
    }

    function getDisplayAccountBalance(transaction: TransactionReconciliationStatementResponseItem): string {
        let currency = defaultCurrency.value;
        let isLiabilityAccount = false;

        if (transaction.type === TransactionType.Transfer && transaction.destinationAccountId === accountId.value) {
            if (allAccountsMap.value[transaction.destinationAccountId]) {
                currency = allAccountsMap.value[transaction.destinationAccountId].currency;
                isLiabilityAccount = allAccountsMap.value[transaction.destinationAccountId].isLiability;
            }
        } else if (allAccountsMap.value[transaction.sourceAccountId]) {
            currency = allAccountsMap.value[transaction.sourceAccountId].currency;
            isLiabilityAccount = allAccountsMap.value[transaction.sourceAccountId].isLiability;
        }

        if (isLiabilityAccount) {
            return formatAmountWithCurrency(-transaction.accountBalance, currency);
        } else {
            return formatAmountWithCurrency(transaction.accountBalance, currency);
        }
    }

    return {
        // states
        accountId,
        startTime,
        endTime,
        reconciliationStatements,
        // computed states
        currentTimezoneOffsetMinutes,
        defaultCurrency,
        allAccountsMap,
        allCategoriesMap,
        displayStartDateTime,
        displayEndDateTime,
        displayTotalOutflows,
        displayTotalInflows,
        // functions
        getDisplayDateTime,
        getDisplayTimezone,
        getDisplaySourceAmount,
        getDisplayDestinationAmount,
        getDisplayAccountBalance
    };
}
