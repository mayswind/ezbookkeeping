import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import type { NumeralSystem } from '@/core/numeral.ts';
import type { WeekDayValue } from '@/core/datetime.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { TransactionType } from '@/core/transaction.ts';
import { StatisticsAnalysisType, ChartDateAggregationType } from '@/core/statistics.ts';
import { KnownFileType } from '@/core/file.ts';

import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type {
    TransactionReconciliationStatementResponse,
    TransactionReconciliationStatementResponseItemWithInfo,
    TransactionReconciliationStatementResponseWithInfo
} from '@/models/transaction.ts';

import { replaceAll } from '@/lib/common.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    parseDateTimeFromUnixTime,
    parseDateTimeFromUnixTimeWithTimezoneOffset
} from '@/lib/datetime.ts';

export function useReconciliationStatementPageBase() {
    const {
        tt,
        getAllAccountBalanceTrendChartTypes,
        getAllStatisticsDateAggregationTypesWithShortName,
        getAllTimezoneTypesUsedForStatistics,
        getCurrentNumeralSystemType,
        formatDateTimeToLongDateTime,
        formatDateTimeToLongDate,
        formatDateTimeToShortTime,
        formatDateTimeToGregorianDefaultDateTime,
        formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
        formatAmountToLocalizedNumeralsWithCurrency
    } = useI18n();

    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();

    const accountId = ref<string>('');
    const startTime = ref<number>(0);
    const endTime = ref<number>(0);
    const reconciliationStatements = ref<TransactionReconciliationStatementResponseWithInfo | undefined>(undefined);
    const chartDataDateAggregationType = ref<number>(ChartDateAggregationType.Day.type);
    const timezoneUsedForDateRange = ref<number>(TimezoneTypeForStatistics.ApplicationTimezone.type);

    const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

    const allChartTypes = computed<TypeAndDisplayName[]>(() => getAllAccountBalanceTrendChartTypes());
    const allDateAggregationTypes = computed<TypeAndDisplayName[]>(() => getAllStatisticsDateAggregationTypesWithShortName(StatisticsAnalysisType.AssetTrends));
    const allTimezoneTypesUsedForDateRange = computed<TypeAndDisplayName[]>(() => getAllTimezoneTypesUsedForStatistics());

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

    const displayStartDateTime = computed<string>(() => {
        const dateTime = parseDateTimeFromUnixTime(startTime.value);
        return formatDateTimeToLongDateTime(dateTime);
    });

    const displayEndDateTime = computed<string>(() => {
        const dateTime = parseDateTimeFromUnixTime(endTime.value);
        return formatDateTimeToLongDateTime(dateTime);
    });

    const displayTotalInflows = computed<string>(() => {
        return formatAmountToLocalizedNumeralsWithCurrency(reconciliationStatements.value?.totalInflows ?? 0, currentAccountCurrency.value);
    });

    const displayTotalOutflows = computed<string>(() => {
        return formatAmountToLocalizedNumeralsWithCurrency(reconciliationStatements.value?.totalOutflows ?? 0, currentAccountCurrency.value);
    });

    const displayTotalBalance = computed<string>(() => {
        return formatAmountToLocalizedNumeralsWithCurrency((reconciliationStatements?.value?.totalInflows ?? 0) - (reconciliationStatements.value?.totalOutflows ?? 0), currentAccountCurrency.value);
    });

    const displayOpeningBalance = computed<string>(() => {
        if (isCurrentLiabilityAccount.value) {
            return formatAmountToLocalizedNumeralsWithCurrency(-(reconciliationStatements?.value?.openingBalance ?? 0), currentAccountCurrency.value);
        } else {
            return formatAmountToLocalizedNumeralsWithCurrency(reconciliationStatements?.value?.openingBalance ?? 0, currentAccountCurrency.value);
        }
    });

    const displayClosingBalance = computed<string>(() => {
        if (isCurrentLiabilityAccount.value) {
            return formatAmountToLocalizedNumeralsWithCurrency(-(reconciliationStatements?.value?.closingBalance ?? 0), currentAccountCurrency.value);
        } else {
            return formatAmountToLocalizedNumeralsWithCurrency(reconciliationStatements?.value?.closingBalance ?? 0, currentAccountCurrency.value);
        }
    });

    function setReconciliationStatements(response: TransactionReconciliationStatementResponse | undefined) {
        if (!response) {
            reconciliationStatements.value = undefined;
            return;
        }

        const responseWithInfo: TransactionReconciliationStatementResponseWithInfo = {
            transactions: response.transactions.map(transaction => {
                const transactionWithInfo: TransactionReconciliationStatementResponseItemWithInfo = {
                    ...transaction,
                    sourceAccount: allAccountsMap.value[transaction.sourceAccountId],
                    sourceAccountName: allAccountsMap.value[transaction.sourceAccountId]?.name || '',
                    destinationAccount: transaction.destinationAccountId && transaction.destinationAccountId !== '0' ? allAccountsMap.value[transaction.destinationAccountId] : undefined,
                    category: allCategoriesMap.value[transaction.categoryId],
                    categoryName: allCategoriesMap.value[transaction.categoryId]?.name || ''
                };
                return transactionWithInfo;
            }),
            totalInflows: response.totalInflows,
            totalOutflows: response.totalOutflows,
            openingBalance: response.openingBalance,
            closingBalance: response.closingBalance
        };

        reconciliationStatements.value = responseWithInfo;
    }

    function getDisplayTransactionType(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        if (transaction.type === TransactionType.ModifyBalance) {
            return tt('Modify Balance');
        } else if (transaction.type === TransactionType.Income) {
            return tt('Income');
        } else if (transaction.type === TransactionType.Expense) {
            return tt('Expense');
        } else if (transaction.type === TransactionType.Transfer && transaction.destinationAccountId === accountId.value) {
            return tt('Transfer In');
        } else if (transaction.type === TransactionType.Transfer && transaction.sourceAccountId === accountId.value) {
            return tt('Transfer Out');
        } else if (transaction.type === TransactionType.Transfer) {
            return tt('Transfer');
        } else {
            return tt('Unknown');
        }
    }

    function getDisplayDateTime(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
        return formatDateTimeToLongDateTime(dateTime);
    }

    function getDisplayDate(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
        return formatDateTimeToLongDate(dateTime);
    }

    function getDisplayTime(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
        return formatDateTimeToShortTime(dateTime);
    }

    function isSameAsDefaultTimezoneOffsetMinutes(transaction: TransactionReconciliationStatementResponseItemWithInfo): boolean {
        return transaction.utcOffset === getTimezoneOffsetMinutes(transaction.time);
    }

    function getDisplayTimezone(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
    }

    function getDisplayTimeInDefaultTimezone(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        const timezoneOffsetMinutes = getTimezoneOffsetMinutes(transaction.time);
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, timezoneOffsetMinutes);
        const utcOffset = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(getUtcOffsetByUtcOffsetMinutes(timezoneOffsetMinutes));
        return `${formatDateTimeToLongDateTime(dateTime)} (UTC${utcOffset})`;
    }

    function getDisplaySourceAmount(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        const currency = transaction.sourceAccount?.currency ?? defaultCurrency.value;
        return formatAmountToLocalizedNumeralsWithCurrency(transaction.sourceAmount, currency);
    }

    function getDisplayDestinationAmount(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        const currency = transaction.destinationAccount?.currency ?? defaultCurrency.value;
        return formatAmountToLocalizedNumeralsWithCurrency(transaction.destinationAmount, currency);
    }

    function getDisplayAccountBalance(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
        let currency = defaultCurrency.value;
        let isLiabilityAccount = false;

        if (transaction.type === TransactionType.Transfer && transaction.destinationAccountId === accountId.value) {
            if (transaction.destinationAccount) {
                currency = transaction.destinationAccount.currency;
                isLiabilityAccount = transaction.destinationAccount.isLiability;
            }
        } else if (transaction.sourceAccount) {
            currency = transaction.sourceAccount.currency;
            isLiabilityAccount = transaction.sourceAccount.isLiability;
        }

        if (isLiabilityAccount) {
            return formatAmountToLocalizedNumeralsWithCurrency(-transaction.accountClosingBalance, currency);
        } else {
            return formatAmountToLocalizedNumeralsWithCurrency(transaction.accountClosingBalance, currency);
        }
    }

    function getExportedData(fileType: KnownFileType): string {
        let separator = ',';

        if (fileType === KnownFileType.TSV) {
            separator = '\t';
        }

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

        const transactions = reconciliationStatements.value?.transactions ?? [];
        const rows = transactions.map(transaction => {
            const transactionTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
            const type = getDisplayTransactionType(transaction);
            let categoryName = transaction.categoryName;
            let displayAmount = formatAmountToWesternArabicNumeralsWithoutDigitGrouping(transaction.sourceAmount);
            let displayAccountName = transaction.sourceAccountName;

            if (transaction.type === TransactionType.ModifyBalance) {
                categoryName = tt('Modify Balance');
            } else if (transaction.type === TransactionType.Transfer && transaction.destinationAccountId === accountId.value) {
                displayAmount = formatAmountToWesternArabicNumeralsWithoutDigitGrouping(transaction.destinationAmount);
            }

            if (transaction.type === TransactionType.Transfer && transaction.destinationAccount) {
                displayAccountName = displayAccountName + ' → ' + (transaction.destinationAccount?.name || '');
            }

            let displayAccountBalance = '';

            if (isCurrentLiabilityAccount.value) {
                displayAccountBalance = formatAmountToWesternArabicNumeralsWithoutDigitGrouping(-transaction.accountClosingBalance);
            } else {
                displayAccountBalance = formatAmountToWesternArabicNumeralsWithoutDigitGrouping(transaction.accountClosingBalance);
            }

            let description = transaction.comment || '';

            if (fileType === KnownFileType.CSV) {
                description = replaceAll(description, ',', ' ');
            } else if (fileType === KnownFileType.TSV) {
                description = replaceAll(description, '\t', ' ');
            }

            return [
                formatDateTimeToGregorianDefaultDateTime(transactionTime),
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
        chartDataDateAggregationType,
        timezoneUsedForDateRange,
        // computed states
        firstDayOfWeek,
        fiscalYearStart,
        defaultCurrency,
        allChartTypes,
        allDateAggregationTypes,
        allTimezoneTypesUsedForDateRange,
        currentAccount,
        currentAccountCurrency,
        isCurrentLiabilityAccount,
        exportFileName,
        displayStartDateTime,
        displayEndDateTime,
        displayTotalInflows,
        displayTotalOutflows,
        displayTotalBalance,
        displayOpeningBalance,
        displayClosingBalance,
        // functions
        setReconciliationStatements,
        getDisplayTransactionType,
        getDisplayDateTime,
        getDisplayDate,
        getDisplayTime,
        isSameAsDefaultTimezoneOffsetMinutes,
        getDisplayTimezone,
        getDisplayTimeInDefaultTimezone,
        getDisplaySourceAmount,
        getDisplayDestinationAmount,
        getDisplayAccountBalance,
        getExportedData
    };
}
