<template>
    <v-data-table
        fixed-header
        fixed-footer
        multi-sort
        item-value="index"
        :class="{ 'insights-explore-table': true, 'text-sm': true, 'disabled': loading, 'loading-skeleton': loading }"
        :headers="dataTableHeaders"
        :items="filteredTransactions"
        :hover="true"
        v-model:items-per-page="itemsPerPage"
        v-model:page="currentPage"
    >
        <template #item.time="{ item }">
            <span>{{ getDisplayDateTime(item) }}</span>
            <v-chip class="ms-1" variant="flat" color="secondary" size="x-small"
                    v-if="item.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimezone(item) }}</v-chip>
        </template>
        <template #item.type="{ item }">
            <v-chip label variant="outlined" size="x-small"
                    :class="{ 'text-income' : item.type === TransactionType.Income, 'text-expense': item.type === TransactionType.Expense }"
                    :color="getTransactionTypeColor(item)">{{ getDisplayTransactionType(item) }}</v-chip>
        </template>
        <template #item.secondaryCategoryName="{ item }">
            <div class="d-flex align-center">
                <ItemIcon size="24px" icon-type="category"
                          :icon-id="item.secondaryCategory?.icon ?? ''"
                          :color="item.secondaryCategory?.color ?? ''"
                          v-if="item.secondaryCategory?.color"></ItemIcon>
                <v-icon size="24" :icon="mdiPencilBoxOutline" v-else-if="!item.secondaryCategory || !item.secondaryCategory?.color" />
                <span class="ms-2" v-if="item.type === TransactionType.ModifyBalance">
                    {{ tt('Modify Balance') }}
                </span>
                <span class="ms-2" v-else-if="item.type !== TransactionType.ModifyBalance && item.secondaryCategory">
                    {{ item.secondaryCategory?.name }}
                </span>
            </div>
        </template>
        <template #item.sourceAmount="{ item }">
            <span :class="{ 'text-expense': item.type === TransactionType.Expense, 'text-income': item.type === TransactionType.Income }">{{ getDisplaySourceAmount(item) }}</span>
            <v-icon class="icon-with-direction mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer && item.sourceAccount?.id !== item.destinationAccount?.id && getDisplaySourceAmount(item) !== getDisplayDestinationAmount(item)"></v-icon>
            <span v-if="item.type === TransactionType.Transfer && item.sourceAccount?.id !== item.destinationAccount?.id && getDisplaySourceAmount(item) !== getDisplayDestinationAmount(item)">{{ getDisplayDestinationAmount(item) }}</span>
        </template>
        <template #item.sourceAccountName="{ item }">
            <div class="d-flex align-center">
                <span v-if="item.sourceAccount">{{ item.sourceAccount?.name }}</span>
                <v-icon class="icon-with-direction mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer"></v-icon>
                <span v-if="item.type === TransactionType.Transfer && item.destinationAccount">{{ item.destinationAccount?.name }}</span>
            </div>
        </template>
        <template #no-data>
            <div v-if="loading && (!filteredTransactions || filteredTransactions.length < 1)">
                <div class="ms-1" style="padding-top: 3px; padding-bottom: 3px" :key="itemIdx" v-for="itemIdx in skeletonData">
                    <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                </div>
            </div>
            <div v-else>
                {{ tt('No transaction data') }}
            </div>
        </template>
        <template #bottom>
            <div class="title-and-toolbar d-flex align-center justify-center text-no-wrap mt-2 mb-4">
                <pagination-buttons :disabled="loading"
                                    :totalPageCount="totalPageCount"
                                    v-model="currentPage">
                </pagination-buttons>
            </div>
        </template>
    </v-data-table>
</template>

<script setup lang="ts">
import PaginationButtons from '@/components/desktop/PaginationButtons.vue';

import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useExploresStore } from '@/stores/explore.ts';

import { TransactionType } from '@/core/transaction.ts';

import {
    type TransactionInsightDataItem
} from '@/models/transaction.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    parseDateTimeFromUnixTime
} from '@/lib/datetime.ts';

import {
    mdiArrowRight,
    mdiPencilBoxOutline
} from '@mdi/js';

interface InsightsExploreDataTableTabProps {
    loading?: boolean;
    countPerPage: number;
}

const props = defineProps<InsightsExploreDataTableTabProps>();
const emit = defineEmits<{
    (e: 'update:countPerPage', value: number): void;
}>();

const {
    tt,
    formatUnixTimeToLongDateTime,
    formatUnixTimeToGregorianDefaultDateTime,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatAmountToLocalizedNumeralsWithCurrency
} = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const exploresStore = useExploresStore();

const currentPage = ref<number>(1);

const currentTimezoneOffsetMinutes = computed<number>(() => getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone));
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const filteredTransactions = computed<TransactionInsightDataItem[]>(() => exploresStore.filteredTransactions);

const itemsPerPage = computed<number>({
    get: () => props.countPerPage,
    set: (value: number) => emit('update:countPerPage', value)
})

const skeletonData = computed<number[]>(() => {
    const data: number[] = [];

    for (let i = 0; i < itemsPerPage.value; i++) {
        data.push(i);
    }

    return data;
});

const totalPageCount = computed<number>(() => {
    if (!filteredTransactions.value || filteredTransactions.value.length < 1) {
        return 1;
    }

    const count = filteredTransactions.value.length;
    return Math.ceil(count / itemsPerPage.value);
});

const dataTableHeaders = computed<object[]>(() => {
    const headers: object[] = [];

    headers.push({ key: 'time', value: 'time', title: tt('Transaction Time'), sortable: true, nowrap: true });
    headers.push({ key: 'type', value: 'type', title: tt('Type'), sortable: true, nowrap: true });
    headers.push({ key: 'secondaryCategoryName', value: 'secondaryCategoryName', title: tt('Category'), sortable: true, nowrap: true });
    headers.push({ key: 'sourceAmount', value: 'sourceAmount', title: tt('Amount'), sortable: true, nowrap: true });
    headers.push({ key: 'sourceAccountName', value: 'sourceAccountName', title: tt('Account'), sortable: true, nowrap: true });
    headers.push({ key: 'comment', value: 'comment', title: tt('Description'), sortable: true, nowrap: true });
    return headers;
});

function getDisplayDateTime(transaction: TransactionInsightDataItem): string {
    return formatUnixTimeToLongDateTime(transaction.time, transaction.utcOffset, currentTimezoneOffsetMinutes.value);
}

function getDisplayTimezone(transaction: TransactionInsightDataItem): string {
    return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
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

function buildExportResults(): { headers: string[], data: string[][] } | undefined {
    if (!filteredTransactions.value) {
        return undefined;
    }

    return {
        headers: [
            tt('Transaction Time'),
            tt('Type'),
            tt('Category'),
            tt('Amount'),
            tt('Account'),
            tt('Description')
        ],
        data: filteredTransactions.value
            .map(transaction => {
                const transactionTime = parseDateTimeFromUnixTime(transaction.time, transaction.utcOffset, currentTimezoneOffsetMinutes.value).getUnixTime();
                const type = getDisplayTransactionType(transaction);

                let categoryName = transaction.secondaryCategoryName;
                let displayAmount = formatAmountToWesternArabicNumeralsWithoutDigitGrouping(transaction.sourceAmount);
                let displayAccountName = transaction.sourceAccountName;

                if (transaction.type === TransactionType.ModifyBalance) {
                    categoryName = tt('Modify Balance');
                } else if (transaction.type === TransactionType.Transfer && transaction.sourceAccount?.id !== transaction.destinationAccount?.id && getDisplaySourceAmount(transaction) !== getDisplayDestinationAmount(transaction)) {
                    displayAmount = displayAmount + ' → ' + formatAmountToWesternArabicNumeralsWithoutDigitGrouping(transaction.destinationAmount);
                }

                if (transaction.type === TransactionType.Transfer && transaction.destinationAccount) {
                    displayAccountName = displayAccountName + ' → ' + (transaction.destinationAccount?.name || '');
                }

                const description = transaction.comment || '';

                return [
                    formatUnixTimeToGregorianDefaultDateTime(transactionTime),
                    type,
                    categoryName,
                    displayAmount,
                    displayAccountName,
                    description
                ];
            }
        )
    };
}

defineExpose({
    buildExportResults
});
</script>

<style>
.v-table.insights-explore-table > .v-table__wrapper > table {
    th:not(:last-child),
    td:not(:last-child) {
        width: auto !important;
        white-space: nowrap;
    }

    th:last-child,
    td:last-child {
        width: 100% !important;
    }
}

.v-table.insights-explore-table.loading-skeleton tr.v-data-table-rows-no-data > td {
    padding: 0;
}
</style>
