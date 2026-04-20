import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useExplorersStore } from '@/stores/explorer.ts';

import { type NameValue, type NameNumeralValue, itemAndIndex } from '@/core/base.ts';
import type { NumeralSystem } from '@/core/numeral.ts';

import { TransactionType } from '@/core/transaction.ts';

import type { TransactionInsightDataItem } from '@/models/transaction.ts';
import type { InsightsExplorer} from '@/models/explorer.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    parseDateTimeFromUnixTimeWithTimezoneOffset
} from '@/lib/datetime.ts';

export function useExplorerDataTablePageBase() {
    const {
        tt,
        getCurrentNumeralSystemType,
        formatDateTimeToLongDateTime,
        formatAmountToLocalizedNumeralsWithCurrency,
        formatNumberToLocalizedNumerals
    } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const explorersStore = useExplorersStore();

    const currentPage = ref<number>(1);

    const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

    const currentExplorer = computed<InsightsExplorer>(() => explorersStore.currentInsightsExplorer);

    const filteredTransactions = computed<TransactionInsightDataItem[]>(() => explorersStore.filteredTransactionsInDataTable);

    const allDataTableQuerySources = computed<NameValue[]>(() => {
        const sources: NameValue[] = [];

        sources.push({
            name: tt('All Queries'),
            value: ''
        });

        for (const [query, index] of itemAndIndex(currentExplorer.value.queries)) {
            if (query.name) {
                sources.push({
                    name: query.name,
                    value: query.id
                });
            } else {
                sources.push({
                    name: tt('format.misc.queryIndex', { index: index + 1 }),
                    value: query.id
                });
            }
        }

        return sources;
    });

    const allPageCounts = computed<NameNumeralValue[]>(() => {
        const pageCounts: NameNumeralValue[] = [];
        const availableCountPerPage: number[] = [ 5, 10, 15, 20, 25, 30, 50 ];

        for (const count of availableCountPerPage) {
            pageCounts.push({ value: count, name: formatNumberToLocalizedNumerals(count) });
        }

        pageCounts.push({ value: -1, name: tt('All') });

        return pageCounts;
    });

    const skeletonData = computed<number[]>(() => {
        const data: number[] = [];

        for (let i = 0; i < currentExplorer.value.countPerPage; i++) {
            data.push(i);
        }

        return data;
    });

    const totalPageCount = computed<number>(() => {
        if (!filteredTransactions.value || filteredTransactions.value.length < 1) {
            return 1;
        }

        const count = filteredTransactions.value.length;
        return Math.ceil(count / currentExplorer.value.countPerPage);
    });

    const dataTableHeaders = computed<object[]>(() => {
        const headers: object[] = [];

        headers.push({ key: 'time', value: 'time', title: tt('Transaction Time'), sortable: true, nowrap: true });
        headers.push({ key: 'type', value: 'type', title: tt('Type'), sortable: true, nowrap: true });
        headers.push({ key: 'secondaryCategoryName', value: 'secondaryCategoryName', title: tt('Category'), sortable: true, nowrap: true });
        headers.push({ key: 'sourceAmount', value: 'sourceAmount', title: tt('Amount'), sortable: true, nowrap: true });
        headers.push({ key: 'sourceAccountName', value: 'sourceAccountName', title: tt('Account'), sortable: true, nowrap: true });

        if (settingsStore.appSettings.showTagInInsightsExplorerPage) {
            headers.push({ key: 'tags', value: 'tags', title: tt('Tags'), sortable: true, nowrap: true });
        }

        headers.push({ key: 'comment', value: 'comment', title: tt('Description'), sortable: true, nowrap: true });
        headers.push({ key: 'operation', title: tt('Operation'), sortable: false, nowrap: true, align: 'center' });
        return headers;
    });

    function getDisplayDateTime(transaction: TransactionInsightDataItem): string {
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
        return formatDateTimeToLongDateTime(dateTime);
    }

    function isSameAsDefaultTimezoneOffsetMinutes(transaction: TransactionInsightDataItem): boolean {
        return transaction.utcOffset === getTimezoneOffsetMinutes(transaction.time);
    }

    function getDisplayTimezone(transaction: TransactionInsightDataItem): string {
        return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
    }

    function getDisplayTimeInDefaultTimezone(transaction: TransactionInsightDataItem): string {
        const timezoneOffsetMinutes = getTimezoneOffsetMinutes(transaction.time);
        const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, timezoneOffsetMinutes);
        const utcOffset = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(getUtcOffsetByUtcOffsetMinutes(timezoneOffsetMinutes));
        return `${formatDateTimeToLongDateTime(dateTime)} (UTC${utcOffset})`;
    }

    function getDisplayTransactionType(transaction: TransactionInsightDataItem): string {
        if (transaction.type === TransactionType.ModifyBalance) {
            return tt('Modify Balance');
        } else if (transaction.type === TransactionType.Income) {
            return tt('Income');
        } else if (transaction.type === TransactionType.Expense) {
            return tt('Expense');
        } else if (transaction.type === TransactionType.Transfer) {
            return tt('Transfer');
        } else {
            return tt('Unknown');
        }
    }

    function getTransactionTypeColor(transaction: TransactionInsightDataItem): string | undefined {
        if (transaction.type === TransactionType.ModifyBalance) {
            return 'secondary';
        } else if (transaction.type === TransactionType.Income) {
            return undefined;
        } else if (transaction.type === TransactionType.Expense) {
            return undefined;
        } else if (transaction.type === TransactionType.Transfer) {
            return 'primary';
        } else {
            return 'default';
        }
    }

    function getDisplaySourceAmount(transaction: TransactionInsightDataItem): string {
        let currency = defaultCurrency.value;

        if (transaction.sourceAccount) {
            currency = transaction.sourceAccount.currency;
        }

        return formatAmountToLocalizedNumeralsWithCurrency(transaction.sourceAmount, currency);
    }

    function getDisplayDestinationAmount(transaction: TransactionInsightDataItem): string {
        let currency = defaultCurrency.value;

        if (transaction.destinationAccount) {
            currency = transaction.destinationAccount.currency;
        }

        return formatAmountToLocalizedNumeralsWithCurrency(transaction.destinationAmount, currency);
    }

    return {
        // states
        currentPage,
        // computed states
        currentExplorer,
        filteredTransactions,
        allDataTableQuerySources,
        allPageCounts,
        skeletonData,
        totalPageCount,
        dataTableHeaders,
        // functions
        getDisplayDateTime,
        isSameAsDefaultTimezoneOffsetMinutes,
        getDisplayTimezone,
        getDisplayTimeInDefaultTimezone,
        getDisplayTransactionType,
        getTransactionTypeColor,
        getDisplaySourceAmount,
        getDisplayDestinationAmount
    };
}
