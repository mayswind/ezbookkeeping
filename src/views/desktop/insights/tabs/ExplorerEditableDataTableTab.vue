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
                        :disabled="true"
                        :label="tt('Data Source')"
                        :items="allDataTableQuerySources"
                        :model-value="currentExplorer.datatableQuerySource"
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
                        <span class="text-subtitle-1">
                            {{ tt('format.misc.selectedCount', { count: formatNumberToLocalizedNumerals(selectedTransactionCount), totalCount: formatNumberToLocalizedNumerals(filteredTransactions.length) }) }}
                        </span>
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
        :class="{ 'insights-editable-explorer-table': true, 'text-sm': true, 'disabled': loading || disabled, 'loading-skeleton': loading }"
        :headers="editableDataTableHeaders"
        :items="filteredTransactions"
        :hover="true"
        v-model:items-per-page="currentExplorer.countPerPage"
        v-model:page="currentPage"
    >
        <template #header.data-table-select>
            <v-checkbox readonly class="always-cursor-pointer"
                        density="compact" width="28"
                        :disabled="!!disabled"
                        :indeterminate="anyButNotAllTransactionSelected"
                        v-model="allTransactionSelected"
            >
                <v-menu activator="parent" location="bottom">
                    <v-list>
                        <v-list-item :prepend-icon="mdiSelectAll"
                                     :title="tt('Select All')"
                                     :disabled="loading || disabled"
                                     @click="selectAll"></v-list-item>
                        <v-list-item :prepend-icon="mdiSelect"
                                     :title="tt('Select None')"
                                     :disabled="loading || disabled"
                                     @click="selectNone"></v-list-item>
                        <v-list-item :prepend-icon="mdiSelectInverse"
                                     :title="tt('Invert Selection')"
                                     :disabled="loading || disabled"
                                     @click="selectInvert"></v-list-item>
                    </v-list>
                </v-menu>
            </v-checkbox>
        </template>
        <template #header.operation>
            <div>
                <span>{{ tt('Operation') }}</span>
                <v-icon :icon="mdiMenuDown" size="20" />
                <v-menu activator="parent" location="bottom">
                    <v-list>
                        <v-list-item :prepend-icon="mdiTextBoxEditOutline"
                                     :title="tt('Update Categories for Expense Transactions')"
                                     :disabled="!isAllSelectedTransactionsExpense"
                                     @click="batchUpdateTransactionCategories(CategoryType.Expense)"></v-list-item>
                        <v-list-item :prepend-icon="mdiTextBoxEditOutline"
                                     :title="tt('Update Categories for Income Transactions')"
                                     :disabled="!isAllSelectedTransactionsIncome"
                                     @click="batchUpdateTransactionCategories(CategoryType.Income)"></v-list-item>
                        <v-list-item :prepend-icon="mdiTextBoxEditOutline"
                                     :title="tt('Update Categories for Transfer Transactions')"
                                     :disabled="!isAllSelectedTransactionsTransfer"
                                     @click="batchUpdateTransactionCategories(CategoryType.Transfer)"></v-list-item>
                    </v-list>
                </v-menu>
            </div>
        </template>
        <template #item.data-table-select="{ item }">
            <v-checkbox density="compact" :disabled="loading || disabled"
                        v-model="selectedTransactions[item.id]"></v-checkbox>
        </template>
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

    <batch-update-category-dialog ref="batchUpdateCategoryDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import PaginationButtons from '@/components/desktop/PaginationButtons.vue';
import BatchUpdateCategoryDialog from '@/views/desktop/insights/dialogs/BatchUpdateCategoryDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useExplorerDataTablePageBase } from '@/views/base/explorer/ExplorerDataTablePageBase.ts';

import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import type { TransactionInsightDataItem } from '@/models/transaction.ts';

import { getObjectOwnFieldWithValueCount } from '@/lib/common.ts';

import {
    mdiArrowRight,
    mdiPencilBoxOutline,
    mdiPound,
    mdiSelect,
    mdiSelectAll,
    mdiSelectInverse,
    mdiMenuDown,
    mdiTextBoxEditOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;
type BatchUpdateCategoryDialogType = InstanceType<typeof BatchUpdateCategoryDialog>;

interface InsightsExplorerDataTableTabProps {
    loading?: boolean;
    disabled?: boolean;
}

defineProps<InsightsExplorerDataTableTabProps>();

const emit = defineEmits<{
    (e: 'click:transaction', value: TransactionInsightDataItem): void;
    (e: 'update:transactions'): void;
}>();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const batchUpdateCategoryDialog = useTemplateRef<BatchUpdateCategoryDialogType>('batchUpdateCategoryDialog');

const {
    tt,
    formatNumberToLocalizedNumerals
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

const selectedTransactions = ref<Record<string, boolean>>({});

const selectedTransactionCount = computed<number>(() => getObjectOwnFieldWithValueCount(selectedTransactions.value, true));
const allTransactionSelected = computed<boolean>(() => selectedTransactionCount.value > 0 && selectedTransactionCount.value === filteredTransactions.value.length);
const anyButNotAllTransactionSelected = computed<boolean>(() => selectedTransactionCount.value > 0 && selectedTransactionCount.value < filteredTransactions.value.length);

const isAllSelectedTransactionsExpense = computed<boolean>(() => isAllSelectedTransactionsSpecificType(TransactionType.Expense));
const isAllSelectedTransactionsIncome = computed<boolean>(() => isAllSelectedTransactionsSpecificType(TransactionType.Income));
const isAllSelectedTransactionsTransfer = computed<boolean>(() => isAllSelectedTransactionsSpecificType(TransactionType.Transfer));

const editableDataTableHeaders = computed<object[]>(() => {
    const headers: object[] = [
        { key: 'data-table-select', fixed: true }
    ];

    headers.push(...dataTableHeaders.value);
    return headers;
});

function isAllSelectedTransactionsSpecificType(type: TransactionType): boolean {
    for (const transaction of filteredTransactions.value) {
        if (selectedTransactions.value[transaction.id] && transaction.type !== type) {
            return false;
        }
    }
    return selectedTransactionCount.value > 0;
}

function getAllSelectedTransactionIds(): string[] {
    const selectedIds: string[] = [];
    for (const transaction of filteredTransactions.value) {
        if (selectedTransactions.value[transaction.id]) {
            selectedIds.push(transaction.id);
        }
    }
    return selectedIds;
}

function selectAll(): void {
    for (const transaction of filteredTransactions.value) {
        selectedTransactions.value[transaction.id] = true;
    }
}

function selectNone(): void {
    for (const transaction of filteredTransactions.value) {
        selectedTransactions.value[transaction.id] = false;
    }
}

function selectInvert(): void {
    for (const transaction of filteredTransactions.value) {
        selectedTransactions.value[transaction.id] = !selectedTransactions.value[transaction.id];
    }
}

function batchUpdateTransactionCategories(type: CategoryType): void {
    batchUpdateCategoryDialog.value?.open({
        type: type,
        updateIds: getAllSelectedTransactionIds() }
    ).then(updatedCount => {
        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: formatNumberToLocalizedNumerals(updatedCount)
            });
        }
        selectedTransactions.value = {};
        emit('update:transactions');
    }).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function showTransaction(transaction: TransactionInsightDataItem): void {
    emit('click:transaction', transaction);
}
</script>

<style>
.v-table.insights-editable-explorer-table > .v-table__wrapper > table {
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

.v-table.insights-editable-explorer-table.loading-skeleton tr.v-data-table-rows-no-data > td {
    padding: 0;
}

.v-table.insights-editable-explorer-table .v-chip.transaction-tag {
    margin-inline-end: 4px;
    margin-top: 2px;
    margin-bottom: 2px;
}

.v-table.insights-editable-explorer-table .v-chip.transaction-tag > .v-chip__content {
    display: block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
