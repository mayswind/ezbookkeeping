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
                        :label="tt('Data Source')"
                        :items="allDataTableQuerySources"
                        v-model="currentExplorer.datatableQuerySource"
                    />
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
                    <v-spacer/>
                    <div class="d-flex align-center">
                        <span class="text-subtitle-1">{{ tt('Total Transactions') }}</span>
                        <span v-if="loading">
                            <v-skeleton-loader class="skeleton-no-margin ms-2" type="text" style="width: 50px" :loading="true"></v-skeleton-loader>
                        </span>
                        <span class="text-subtitle-1 ms-2" v-else-if="!loading">
                            {{ formatNumberToLocalizedNumerals(filteredTransactions.length) }}
                        </span>
                        <span class="text-subtitle-1 ms-3" v-if="loading || filteredTransactionsStatistic">{{ tt('Total Amount') }}</span>
                        <span v-if="loading">
                            <v-skeleton-loader class="skeleton-no-margin ms-2" type="text" style="width: 80px" :loading="true"></v-skeleton-loader>
                        </span>
                        <span class="text-subtitle-1 ms-2" v-else-if="!loading && filteredTransactionsStatistic">
                            {{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.totalAmount) }}
                        </span>
                        <v-tooltip interactive class="table-tooltip" activator="parent" v-if="!loading && filteredTransactions.length > 0 && filteredTransactionsStatistic">
                            <v-table density="compact">
                                <tbody>
                                <tr>
                                    <td>{{ tt('Total Amount') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.totalAmount) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Total Income') }}</td>
                                    <td class="text-end text-income">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.totalIncome) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Total Expense') }}</td>
                                    <td class="text-end text-expense">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.totalExpense) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Net Income') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.netIncome) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Average Amount') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.averageAmount) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Median Amount') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.medianAmount) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Minimum Amount') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.minimumAmount) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Maximum Amount') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.maximumAmount) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('90th Percentile Amount') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.p90Amount) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Range (Max - Min)') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.range) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Interquartile Range (Q3 - Q1)') }}</td>
                                    <td class="text-end">{{ formatAmountToLocalizedNumeralsWithCurrency(filteredTransactionsStatistic.interquartileRange) }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Top 5 Amount Share') }}</td>
                                    <td class="text-end">{{ isDefined(filteredTransactionsStatistic.top5AmountShare) ? formatPercentToLocalizedNumerals(filteredTransactionsStatistic.top5AmountShare, 2, '<0.01') : '-' }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Transactions for 80% of Amount') }}</td>
                                    <td class="text-end">{{ isDefined(filteredTransactionsStatistic.transactionsFor80PercentAmount) ? formatPercentToLocalizedNumerals(filteredTransactionsStatistic.transactionsFor80PercentAmount, 2, '<0.01') : '-' }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Variance') }}</td>
                                    <td class="text-end">{{ isDefined(filteredTransactionsStatistic.variance) ? formatNumberToLocalizedNumerals(filteredTransactionsStatistic.variance, 2) : '-' }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Standard Deviation') }}</td>
                                    <td class="text-end">{{ isDefined(filteredTransactionsStatistic.standardDeviation) ? formatNumberToLocalizedNumerals(filteredTransactionsStatistic.standardDeviation, 2) : '-' }}</td>
                                </tr>
                                <tr>
                                    <td>{{ tt('Coefficient of Variation') }}</td>
                                    <td class="text-end">{{ isDefined(filteredTransactionsStatistic.coefficientOfVariation) ? formatNumberToLocalizedNumerals(filteredTransactionsStatistic.coefficientOfVariation, 2) : '-' }}</td>
                                </tr>
                                </tbody>
                            </v-table>
                        </v-tooltip>
                    </div>
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
            <v-chip class="ms-1" variant="flat" color="grey" size="x-small"
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

import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useExplorerDataTablePageBase } from '@/views/base/explorer/ExplorerDataTablePageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { type InsightsExplorerTransactionStatisticData, useExplorersStore } from '@/stores/explorer.ts';

import { TransactionType } from '@/core/transaction.ts';
import type { TransactionInsightDataItem } from '@/models/transaction.ts';

import { isDefined, replaceAll } from '@/lib/common.ts';

import {
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
    formatDateTimeToGregorianDefaultDateTime,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatNumberToLocalizedNumerals,
    formatPercentToLocalizedNumerals
} = useI18n();

const {
    currentPage,
    currentExplorer,
    filteredTransactions,
    allDataTableQuerySources,
    allPageCounts,
    skeletonData,
    totalPageCount,
    dataTableHeaders,
    getDisplayDateTime,
    isSameAsDefaultTimezoneOffsetMinutes,
    getDisplayTimezone,
    getDisplayTimeInDefaultTimezone,
    getDisplayTransactionType,
    getTransactionTypeColor,
    getDisplaySourceAmount,
    getDisplayDestinationAmount
} = useExplorerDataTablePageBase();

const settingsStore = useSettingsStore();
const explorersStore = useExplorersStore();

const filteredTransactionsStatistic = computed<InsightsExplorerTransactionStatisticData | undefined>(() => explorersStore.filteredTransactionsInDataTableStatistic);

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
                let displayAmount = formatAmountToWesternArabicNumeralsWithoutDigitGrouping(transaction.sourceAmount, transaction.sourceAccount?.currency);
                let displayAccountName = transaction.sourceAccountName;

                if (transaction.type === TransactionType.ModifyBalance) {
                    categoryName = tt('Modify Balance');
                } else if (transaction.type === TransactionType.Transfer && transaction.sourceAccount?.id !== transaction.destinationAccount?.id && getDisplaySourceAmount(transaction) !== getDisplayDestinationAmount(transaction)) {
                    displayAmount = displayAmount + ' → ' + formatAmountToWesternArabicNumeralsWithoutDigitGrouping(transaction.destinationAmount, transaction.destinationAccount?.currency);
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
