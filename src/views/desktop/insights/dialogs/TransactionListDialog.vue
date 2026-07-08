<template>
    <v-dialog width="1200" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex flex-wrap align-center justify-center">
                    <h4 class="text-h4">{{ title }}</h4>
                    <v-spacer/>
                    <div class="title-and-toolbar d-flex align-center justify-center text-no-wrap">
                        <span class="text-body-1" v-if="transactions.length > 10">{{ tt('Transactions Per Page') }}</span>
                        <v-select class="ms-2" density="compact" max-width="100"
                                  item-title="name"
                                  item-value="value"
                                  :items="allPageCounts"
                                  v-model="countPerPage"
                                  v-if="transactions.length > 10"
                        />
                        <pagination-buttons density="compact"
                                            :totalPageCount="totalPageCount"
                                            v-model="currentPage"
                                            v-if="transactions.length > 10">
                        </pagination-buttons>
                    </div>
                </div>
            </template>

            <v-card-text>
                <v-data-table
                    fixed-header
                    fixed-footer
                    multi-sort
                    density="compact"
                    item-value="index"
                    :class="{ 'insights-explorer-transactions-dialog-table': true, 'text-sm': true }"
                    :headers="dataTableHeaders"
                    :items="transactions"
                    :hover="true"
                    v-model:items-per-page="countPerPage"
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
                        <v-btn density="compact" variant="text" color="default"
                               @click="showTransaction(item)">
                            {{ tt('View') }}
                        </v-btn>
                    </template>
                    <template #no-data>
                        {{ tt('No transaction data') }}
                    </template>
                    <template #bottom></template>
                </v-data-table>
            </v-card-text>

            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Close') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import PaginationButtons from '@/components/desktop/PaginationButtons.vue';

import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useExplorerDataTablePageBase } from '@/views/base/explorer/ExplorerDataTablePageBase.ts';

import { useExplorersStore } from '@/stores/explorer.ts';

import { values } from '@/core/base.ts';
import { TransactionType } from '@/core/transaction.ts';

import type { TransactionInsightDataItem } from '@/models/transaction.ts';

import {
    mdiArrowRight,
    mdiPencilBoxOutline,
    mdiPound
} from '@mdi/js';

const { tt } = useI18n();

const {
    allPageCounts,
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

const explorersStore = useExplorersStore();

const emit = defineEmits<{
    (e: 'click:transaction', value: TransactionInsightDataItem): void;
}>();

const showState = ref<boolean>(false);
const title = ref<string>('');
const categoryId = ref<string>('');
const seriesId = ref<string>('');
const currentPage = ref<number>(1);
const countPerPage = ref<number>(10);

const transactions = computed<TransactionInsightDataItem[]>(() => {
    const categoriedData = explorersStore.categoriedTransactions[categoryId.value];
    let transactionList: TransactionInsightDataItem[] = [];

    if (!categoriedData) {
        return transactionList;
    }

    const seriesData = seriesId.value ? categoriedData.trasactions[seriesId.value] : undefined;

    if (seriesData) {
        transactionList = seriesData.trasactions;
    } else {
        for (const seriesTransactions of values(categoriedData.trasactions)) {
            transactionList.push(...seriesTransactions.trasactions);
        }
    }

    return transactionList;
});

const totalPageCount = computed<number>(() => {
    if (!transactions.value || transactions.value.length < 1) {
        return 1;
    }

    return Math.ceil(transactions.value.length / countPerPage.value);
});

function open(options: { title: string, categoryId: string, seriesId?: string }): void {
    title.value = options.title;
    categoryId.value = options.categoryId;
    seriesId.value = options.seriesId || '';
    currentPage.value = 1;
    countPerPage.value = 10;
    showState.value = true;
}

function showTransaction(transaction: TransactionInsightDataItem): void {
    emit('click:transaction', transaction);
}

function cancel(): void {
    showState.value = false;
}

defineExpose({
    open
});
</script>

<style>
.v-table.insights-explorer-transactions-dialog-table > .v-table__wrapper > table {
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

.v-table.insights-explorer-transactions-dialog-table .v-chip.transaction-tag {
    margin-inline-end: 4px;
    margin-top: 2px;
    margin-bottom: 2px;
}

.v-table.insights-explorer-transactions-dialog-table .v-chip.transaction-tag > .v-chip__content {
    display: block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
