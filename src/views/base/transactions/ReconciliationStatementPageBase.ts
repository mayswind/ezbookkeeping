import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { KnownDateTimeFormat } from '@/core/datetime.ts';
import { TransactionType } from '@/core/transaction.ts';
import { KnownFileType } from '@/core/file.ts';
import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionReconciliationStatementResponseItem } from '@/models/transaction.ts';

import {
    replaceAll,
    removeAll
} from '@/lib/common.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    parseDateFromUnixTime,
    formatUnixTime,
    getUnixTime
} from '@/lib/datetime.ts';

export function useReconciliationStatementPageBase() {
    const {
        tt,
        getCurrentDigitGroupingSymbol,
        formatUnixTimeToLongDateTime,
        formatAmount,
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

    const currentAccount = computed(() => allAccountsMap.value[accountId.value]);
    const currentAccountCurrency = computed<string>(() => currentAccount.value?.currency ?? defaultCurrency.value);
    const isCurrentLiabilityAccount = computed<boolean>(() => currentAccount.value?.isLiability ?? false);

    const exportFileName = computed<string>(() => {
        const nickname = userStore.currentUserNickname;

        if (nickname) {
            return tt('dataExport.exportReconciliationStatementsFileName', {
                nickname: nickname
            });
        }

        return tt('dataExport.defaultExportReconciliationStatementsFileName');
    });

    const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
    const allCategoriesMap = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);

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
        return formatAmountWithCurrency(totalOutflows.value, currentAccountCurrency.value);
    });

    const displayTotalInflows = computed<string>(() => {
        return formatAmountWithCurrency(totalInflows.value, currentAccountCurrency.value);
    });

    const displayTotalBalance = computed<string>(() => {
        return formatAmountWithCurrency(totalInflows.value - totalOutflows.value, currentAccountCurrency.value);
    });

    const displayOpeningBalance = computed<string>(() => {
        if (isCurrentLiabilityAccount.value) {
            return formatAmountWithCurrency(-openingBalance.value, currentAccountCurrency.value);
        } else {
            return formatAmountWithCurrency(openingBalance.value, currentAccountCurrency.value);
        }
    });

    const displayClosingBalance = computed<string>(() => {
        if (isCurrentLiabilityAccount.value) {
            return formatAmountWithCurrency(-closingBalance.value, currentAccountCurrency.value);
        } else {
            return formatAmountWithCurrency(closingBalance.value, currentAccountCurrency.value);
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

    function getExportedData(fileType: KnownFileType): string {
        let separator = ',';

        if (fileType === KnownFileType.TSV) {
            separator = '\t';
        }

        const digitGroupingSymbol = getCurrentDigitGroupingSymbol();
        const accountBalanceName = isCurrentLiabilityAccount.value ? 'Account Outstanding Balance' : 'Account Balance';

        const header = [
            tt('Transaction Time'),
            tt('Type'),
            tt('Category'),
            tt('Amount'),
            tt('Account'),
            tt(accountBalanceName),
            tt('Description')
        ].join(separator) + '\n';

        const rows = reconciliationStatements.value.map(transaction => {
            const transactionTime = getUnixTime(parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentTimezoneOffsetMinutes.value));
            let type = '';
            let categoryName = allCategoriesMap.value[transaction.categoryId]?.name || '';
            let displayAmount = removeAll(formatAmount(transaction.sourceAmount), digitGroupingSymbol);
            let displayAccountName = allAccountsMap.value[transaction.sourceAccountId]?.name || '';

            if (transaction.type === TransactionType.ModifyBalance) {
                type = tt('Modify Balance');
                categoryName = '-';
            } else if (transaction.type === TransactionType.Income) {
                type = tt('Income');
            } else if (transaction.type === TransactionType.Expense) {
                type = tt('Expense');
            } else if (transaction.type === TransactionType.Transfer && transaction.destinationAccountId === accountId.value) {
                type = tt('Transfer In');
                displayAmount = removeAll(formatAmount(transaction.destinationAmount), digitGroupingSymbol);
            } else if (transaction.type === TransactionType.Transfer && transaction.sourceAccountId === accountId.value) {
                type = tt('Transfer Out');
            } else if (transaction.type === TransactionType.Transfer) {
                type = tt('Transfer');
            } else {
                type = tt('Unknown');
            }

            if (transaction.type === TransactionType.Transfer && allAccountsMap.value[transaction.destinationAccountId]) {
                displayAccountName = displayAccountName + ' → ' + (allAccountsMap.value[transaction.destinationAccountId]?.name || '');
            }

            let displayAccountBalance = '';

            if (isCurrentLiabilityAccount.value) {
                displayAccountBalance = removeAll(formatAmount(-transaction.accountBalance), digitGroupingSymbol);
            } else {
                displayAccountBalance = removeAll(formatAmount(transaction.accountBalance), digitGroupingSymbol);
            }

            let description = transaction.comment || '';

            if (fileType === KnownFileType.CSV) {
                description = replaceAll(description, ',', ' ');
            } else if (fileType === KnownFileType.TSV) {
                description = replaceAll(description, '\t', ' ');
            }

            return [
                formatUnixTime(transactionTime, KnownDateTimeFormat.DefaultDateTime.format),
                type,
                categoryName,
                displayAmount,
                displayAccountName,
                displayAccountBalance,
                description
            ].join(separator);
        });

        return header + rows.join('\n');
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
        currentAccount,
        currentAccountCurrency,
        isCurrentLiabilityAccount,
        exportFileName,
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
        getDisplayAccountBalance,
        getExportedData
    };
}
