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
    const openingBalance = ref<number>(0);
    const closingBalance = ref<number>(0);

    const currentTimezoneOffsetMinutes = computed<number>(() => getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone));
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

    const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
    const allCategoriesMap = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);

    const accountCurrency = computed<string>(() => {
        let currency = defaultCurrency.value;

        if (allAccountsMap.value[accountId.value]) {
            currency = allAccountsMap.value[accountId.value].currency;
        }

        return currency;
    });

    const totalOutflows = computed<number>(() => {
        let totalOutflows = 0;

        for (let i = 0; i < reconciliationStatements.value.length; i++) {
            const transaction = reconciliationStatements.value[i];

            if (transaction.type === TransactionType.Expense) {
                totalOutflows += transaction.sourceAmount;
            } else if (transaction.type === TransactionType.Transfer && transaction.sourceAccountId === accountId.value) {
                totalOutflows += transaction.sourceAmount;
            }
        }

        return totalOutflows;
    });

    const totalInflows = computed<number>(() => {
        let totalInflows = 0;

        for (let i = 0; i < reconciliationStatements.value.length; i++) {
            const transaction = reconciliationStatements.value[i];

            if (transaction.type === TransactionType.Income) {
                totalInflows += transaction.sourceAmount;
            } else if (transaction.type === TransactionType.Transfer && transaction.destinationAccountId === accountId.value) {
                totalInflows += transaction.destinationAmount;
            }
        }

        return totalInflows;
    });

    const displayStartDateTime = computed<string>(() => {
        return formatUnixTimeToLongDateTime(startTime.value);
    });

    const displayEndDateTime = computed<string>(() => {
        return formatUnixTimeToLongDateTime(endTime.value);
    });

    const displayTotalOutflows = computed<string>(() => {
        return formatAmountWithCurrency(totalOutflows.value, accountCurrency.value);
    });

    const displayTotalInflows = computed<string>(() => {
        return formatAmountWithCurrency(totalInflows.value, accountCurrency.value);
    });

    const displayTotalBalance = computed<string>(() => {
        return formatAmountWithCurrency(totalInflows.value - totalOutflows.value, accountCurrency.value);
    });

    const displayOpeningBalance = computed<string>(() => {
        let isLiabilityAccount = false;

        if (allAccountsMap.value[accountId.value]) {
            isLiabilityAccount = allAccountsMap.value[accountId.value].isLiability;
        }

        if (isLiabilityAccount) {
            return formatAmountWithCurrency(-openingBalance.value, accountCurrency.value);
        } else {
            return formatAmountWithCurrency(openingBalance.value, accountCurrency.value);
        }
    });

    const displayClosingBalance = computed<string>(() => {
        let isLiabilityAccount = false;

        if (allAccountsMap.value[accountId.value]) {
            isLiabilityAccount = allAccountsMap.value[accountId.value].isLiability;
        }

        if (isLiabilityAccount) {
            return formatAmountWithCurrency(-closingBalance.value, accountCurrency.value);
        } else {
            return formatAmountWithCurrency(closingBalance.value, accountCurrency.value);
        }
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
        openingBalance,
        closingBalance,
        // computed states
        currentTimezoneOffsetMinutes,
        defaultCurrency,
        allAccountsMap,
        allCategoriesMap,
        displayStartDateTime,
        displayEndDateTime,
        displayTotalOutflows,
        displayTotalInflows,
        displayTotalBalance,
        displayOpeningBalance,
        displayClosingBalance,
        // functions
        getDisplayDateTime,
        getDisplayTimezone,
        getDisplaySourceAmount,
        getDisplayDestinationAmount,
        getDisplayAccountBalance
    };
}
