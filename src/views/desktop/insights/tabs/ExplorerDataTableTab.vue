<template>
    <v-card-text class="px-5 py-0 mb-4">
        <v-row>
            <v-col cols="12">
                <div class="d-flex overflow-x-auto align-center gap-2 pt-2">
                    <v-select
                        class="flex-0-0"
                        min-width="150"
                        item-title="name"
                        item-value="value"
                        density="compact"
                        :disabled="loading || disabled"
                        :label="tt('Transactions Per Page')"
                        :items="allPageCounts"
                        v-model="currentExplorer.countPerPage"
                    />
                </div>
            </v-col>
        </v-row>
    </v-card-text>
    <v-data-table
        fixed-header
        fixed-footer
        multi-sort
        item-value="index"
        :class="{ 'insights-explorer-table': true, 'text-sm': true, 'disabled': loading || disabled, 'loading-skeleton': loading }"
        :headers="dataTableHeaders"
        :items="filteredTransactions"
        :hover="true"
        v-model:items-per-page="currentExplorer.countPerPage"
        v-model:page="currentPage"
    >
        <template #item.time="{ item }">
            <span>{{ getDisplayDateTime(item) }}</span>
            <v-chip class="ms-1" variant="flat" color="secondary" size="x-small"
                    v-if="!isSameAsDefaultTimezoneOffsetMinutes(item)">{{ getDisplayTimezone(item) }}</v-chip>
            <v-tooltip activator="parent" v-if="!isSameAsDefaultTimezoneOffsetMinutes(item)">{{ getDisplayTimeInDefaultTimezone(item) }}</v-tooltip>
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
        <template #item.tags="{ item }">
            <div class="d-flex">
                <v-chip class="transaction-tag" size="small"
                        :key="tag.id" :prepend-icon="mdiPound"
                        :text="tag.name"
                        v-for="tag in item.tags"/>
                <v-chip class="transaction-tag" size="small"
                        :text="tt('None')"
                        v-if="!item.tagIds || !item.tagIds.length"/>
            </div>
        </template>
        <template #item.operation="{ item }">
            <v-btn density="compact" variant="text" color="default" :disabled="loading || disabled"
                   @click="showTransaction(item)">
                {{ tt('View') }}
            </v-btn>
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
                <pagination-buttons :disabled="loading || disabled"
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
import { useExplorersStore } from '@/stores/explorer.ts';

import type { NameNumeralValue } from '@/core/base.ts';
import type { NumeralSystem } from '@/core/numeral.ts';
import { TransactionType } from '@/core/transaction.ts';

import type { TransactionInsightDataItem } from '@/models/transaction.ts';
import type { InsightsExplorer} from '@/models/explorer.ts';

import { replaceAll } from '@/lib/common.ts';

import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    parseDateTimeFromUnixTimeWithTimezoneOffset
} from '@/lib/datetime.ts';

import {
    mdiArrowRight,
    mdiPencilBoxOutline,
    mdiPound
} from '@mdi/js';

interface InsightsExplorerDataTableTabProps {
    loading?: boolean;
    disabled?: boolean;
}

defineProps<InsightsExplorerDataTableTabProps>();

const emit = defineEmits<{
    (e: 'click:transaction', value: TransactionInsightDataItem): void;
}>();

const {
    tt,
    getCurrentNumeralSystemType,
    formatDateTimeToLongDateTime,
    formatDateTimeToGregorianDefaultDateTime,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatAmountToLocalizedNumeralsWithCurrency
} = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const explorersStore = useExplorersStore();

const currentPage = ref<number>(1);

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const currentExplorer = computed<InsightsExplorer>(() => explorersStore.currentInsightsExplorer);

const filteredTransactions = computed<TransactionInsightDataItem[]>(() => explorersStore.filteredTransactions);

const allPageCounts = computed<NameNumeralValue[]>(() => {
    const pageCounts: NameNumeralValue[] = [];
    const availableCountPerPage: number[] = [ 5, 10, 15, 20, 25, 30, 50 ];

    for (const count of availableCountPerPage) {
        pageCounts.push({ value: count, name: numeralSystem.value.formatNumber(count) });
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

function showTransaction(transaction: TransactionInsightDataItem): void {
    emit('click:transaction', transaction);
}

function buildExportResults(): { headers: string[], data: string[][] } | undefined {
    if (!filteredTransactions.value) {
        return undefined;
    }

    const includeTags = settingsStore.appSettings.showTagInInsightsExplorerPage;

    const headers = [
        tt('Transaction Time'),
        tt('Type'),
        tt('Category'),
        tt('Amount'),
        tt('Account')
    ];

    if (includeTags) {
        headers.push(tt('Tags'));
    }

    headers.push(tt('Description'));

    return {
        headers: headers,
        data: filteredTransactions.value
            .map(transaction => {
                const transactionTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
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

                const data = [
                    formatDateTimeToGregorianDefaultDateTime(transactionTime),
                    type,
                    categoryName,
                    displayAmount,
                    displayAccountName
                ];

                if (includeTags) {
                    const tags = transaction.tags && transaction.tags.length ? transaction.tags.map(tag => replaceAll(tag.name, ';', ' ')).join(';') : tt('None');
                    data.push(tags);
                }

                data.push(description);

                return data;
            }
        )
    };
}

defineExpose({
    buildExportResults
});
</script>

<style>
.v-table.insights-explorer-table > .v-table__wrapper > table {
    th:not(:nth-last-child(2)),
    td:not(:nth-last-child(2)) {
        width: auto !important;
        white-space: nowrap;
    }

    th:nth-last-child(2),
    td:nth-last-child(2) {
        width: 100% !important;
    }
}

.v-table.insights-explorer-table.loading-skeleton tr.v-data-table-rows-no-data > td {
    padding: 0;
}

.v-table.insights-explorer-table .v-chip.transaction-tag {
    margin-inline-end: 4px;
    margin-top: 2px;
    margin-bottom: 2px;
}

.v-table.insights-explorer-table .v-chip.transaction-tag > .v-chip__content {
    display: block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
