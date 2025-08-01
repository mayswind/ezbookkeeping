<template>
    <v-dialog :min-height="400" :persistent="loading" v-model="showState">
        <v-card class="pa-6 pa-sm-10 pa-md-12">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ tt('Reconciliation Statement') }}</h4>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ml-2" :icon="true" :loading="loading" @click="reload()">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                           :icon="true" :disabled="loading">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiInvoiceTextPlusOutline"
                                             :title="tt('Add Transaction')"
                                             @click="addTransaction()"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-item :prepend-icon="mdiInvoiceTextEditOutline"
                                             :title="tt('Update Closing Balance')"
                                             @click="updateClosingBalance()"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-item :prepend-icon="mdiComma"
                                             :disabled="!reconciliationStatements || !reconciliationStatements.transactions || reconciliationStatements.transactions.length < 1"
                                             @click="exportReconciliationStatements(KnownFileType.CSV)">
                                    <v-list-item-title>{{ tt('Export to CSV (Comma-separated values) File') }}</v-list-item-title>
                                </v-list-item>
                                <v-list-item :prepend-icon="mdiKeyboardTab"
                                             :disabled="!reconciliationStatements || !reconciliationStatements.transactions || reconciliationStatements.transactions.length < 1"
                                             @click="exportReconciliationStatements(KnownFileType.TSV)">
                                    <v-list-item-title>{{ tt('Export to TSV (Tab-separated values) File') }}</v-list-item-title>
                                </v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>

            <template #subtitle>
                <div class="text-body-1 text-center text-wrap mt-2" v-if="!startTime && !endTime">
                    <span>{{ tt('All') }}</span>
                </div>
                <div class="text-body-1 text-center text-wrap mt-2" v-if="startTime || endTime">
                    <span>{{ displayStartDateTime }}</span>
                    <span> - </span>
                    <span>{{ displayEndDateTime }}</span>
                </div>
            </template>

            <div class="d-flex align-center mb-4">
                <div class="d-flex align-center text-body-1">
                    <span class="ml-2">{{ tt('Opening Balance') }}</span>
                    <span class="text-primary" v-if="loading">
                        <v-skeleton-loader class="skeleton-no-margin ml-3" type="text" style="width: 80px" :loading="true"></v-skeleton-loader>
                    </span>
                    <span class="text-primary ml-2" v-else-if="!loading">
                        {{ displayOpeningBalance }}
                    </span>
                    <span class="ml-3">{{ tt('Closing Balance') }}</span>
                    <span class="text-primary" v-if="loading">
                        <v-skeleton-loader class="skeleton-no-margin ml-3" type="text" style="width: 80px" :loading="true"></v-skeleton-loader>
                    </span>
                    <span class="text-primary ml-2" v-else-if="!loading">
                        {{ displayClosingBalance }}
                    </span>
                </div>
                <v-spacer/>
                <div class="d-flex align-center text-body-1">
                    <span class="ml-2">{{ tt('Total Inflows') }}</span>
                    <span class="text-income" v-if="loading">
                        <v-skeleton-loader class="skeleton-no-margin ml-3" type="text" style="width: 80px" :loading="true"></v-skeleton-loader>
                    </span>
                    <span class="text-income ml-2" v-else-if="!loading">
                        {{ displayTotalInflows }}
                    </span>
                    <span class="ml-3">{{ tt('Total Outflows') }}</span>
                    <span class="text-expense" v-if="loading">
                        <v-skeleton-loader class="skeleton-no-margin ml-3" type="text" style="width: 80px" :loading="true"></v-skeleton-loader>
                    </span>
                    <span class="text-expense ml-2" v-else-if="!loading">
                        {{ displayTotalOutflows }}
                    </span>
                    <span class="ml-3">{{ tt('Net Cash Flow') }}</span>
                    <span class="text-primary" v-if="loading">
                        <v-skeleton-loader class="skeleton-no-margin ml-3" type="text" style="width: 80px" :loading="true"></v-skeleton-loader>
                    </span>
                    <span class="text-primary ml-2" v-else-if="!loading">
                        {{ displayTotalBalance }}
                    </span>
                </div>
            </div>

            <v-data-table
                fixed-header
                fixed-footer
                multi-sort
                density="compact"
                item-value="index"
                :class="{ 'disabled': loading }"
                :headers="dataTableHeaders"
                :items="reconciliationStatements?.transactions ?? []"
                :no-data-text="loading ? '' : tt('No transaction data')"
                v-model:items-per-page="countPerPage"
                v-model:page="currentPage"
            >
                <template #item.time="{ item }">
                    <span>{{ getDisplayDateTime(item) }}</span>
                    <v-chip class="ml-1" variant="flat" color="secondary" size="x-small"
                            v-if="item.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimezone(item) }}</v-chip>
                </template>
                <template #item.type="{ item }">
                    <v-chip label color="secondary" variant="outlined" size="x-small" v-if="item.type === TransactionType.ModifyBalance">{{ tt('Modify Balance') }}</v-chip>
                    <v-chip label class="text-income" variant="outlined" size="x-small" v-else-if="item.type === TransactionType.Income">{{ tt('Income') }}</v-chip>
                    <v-chip label class="text-expense" variant="outlined" size="x-small" v-else-if="item.type === TransactionType.Expense">{{ tt('Expense') }}</v-chip>
                    <v-chip label color="primary" variant="outlined" size="x-small" v-else-if="item.type === TransactionType.Transfer && item.destinationAccountId === accountId">{{ tt('Transfer In') }}</v-chip>
                    <v-chip label color="primary" variant="outlined" size="x-small" v-else-if="item.type === TransactionType.Transfer && item.sourceAccountId === accountId">{{ tt('Transfer Out') }}</v-chip>
                    <v-chip label color="primary" variant="outlined" size="x-small" v-else-if="item.type === TransactionType.Transfer">{{ tt('Transfer') }}</v-chip>
                    <v-chip label color="default" variant="outlined" size="x-small" v-else>{{ tt('Unknown') }}</v-chip>
                </template>
                <template #item.categoryId="{ item }">
                    <div class="d-flex align-center">
                        <ItemIcon size="24px" icon-type="category"
                                  :icon-id="allCategoriesMap[item.categoryId].icon"
                                  :color="allCategoriesMap[item.categoryId].color"
                                  v-if="allCategoriesMap[item.categoryId] && allCategoriesMap[item.categoryId]?.color"></ItemIcon>
                        <v-icon size="24" :icon="mdiPencilBoxOutline" v-else-if="!allCategoriesMap[item.categoryId] || !allCategoriesMap[item.categoryId]?.color" />
                        <span class="ml-2" v-if="item.type === TransactionType.ModifyBalance">
                            {{ tt('Modify Balance') }}
                        </span>
                        <span class="ml-2" v-else-if="item.type !== TransactionType.ModifyBalance && allCategoriesMap[item.categoryId]">
                            {{ allCategoriesMap[item.categoryId].name }}
                        </span>
                    </div>
                </template>
                <template #item.sourceAmount="{ item }">
                    <span :class="{ 'text-expense': item.type === TransactionType.Expense, 'text-income': item.type === TransactionType.Income }">{{ getDisplaySourceAmount(item) }}</span>
                    <v-icon class="mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId && getDisplaySourceAmount(item) !== getDisplayDestinationAmount(item)"></v-icon>
                    <span v-if="item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId && getDisplaySourceAmount(item) !== getDisplayDestinationAmount(item)">{{ getDisplayDestinationAmount(item) }}</span>
                </template>
                <template #item.sourceAccountId="{ item }">
                    <div class="d-flex align-center">
                        <span v-if="item.sourceAccountId && allAccountsMap[item.sourceAccountId]">{{ allAccountsMap[item.sourceAccountId].name }}</span>
                        <v-icon class="mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer"></v-icon>
                        <span v-if="item.type === TransactionType.Transfer && item.destinationAccountId && allAccountsMap[item.destinationAccountId]">{{ allAccountsMap[item.destinationAccountId].name }}</span>
                    </div>
                </template>
                <template #item.accountBalance="{ item }">
                    <span>{{ getDisplayAccountBalance(item) }}</span>
                </template>
                <template #item.operation="{ item }">
                    <v-btn density="compact" variant="text" color="default" :disabled="loading || item.type === TransactionType.ModifyBalance"
                           @click="showTransaction(item)">
                        {{ tt('View') }}
                    </v-btn>
                </template>
                <template #bottom>
                    <div class="title-and-toolbar d-flex align-center text-no-wrap mt-2" v-if="loading || (reconciliationStatements && reconciliationStatements.transactions && reconciliationStatements.transactions.length)">
                        <span class="ml-2">{{ tt('Total Transactions') }}</span>
                        <span v-if="loading">
                            <v-skeleton-loader class="skeleton-no-margin ml-3" type="text" style="width: 80px" :loading="true"></v-skeleton-loader>
                        </span>
                        <span class="ml-2" v-else-if="!loading">
                            {{ reconciliationStatements?.transactions.length ?? 0 }}
                        </span>
                        <v-spacer/>
                        <span v-if="reconciliationStatements && reconciliationStatements.transactions && reconciliationStatements.transactions.length > 10">
                            {{ tt('Transactions Per Page') }}
                        </span>
                        <v-select class="ml-2" density="compact" max-width="100"
                                  item-title="title"
                                  item-value="value"
                                  :disabled="loading"
                                  :items="reconciliationStatementsTablePageOptions"
                                  v-model="countPerPage"
                                  v-if="reconciliationStatements && reconciliationStatements.transactions && reconciliationStatements.transactions.length > 10"
                        />
                        <pagination-buttons density="compact"
                                            :disabled="loading"
                                            :totalPageCount="totalPageCount"
                                            v-model="currentPage"
                                            v-if="reconciliationStatements && reconciliationStatements.transactions && reconciliationStatements.transactions.length > 10">
                        </pagination-buttons>
                    </div>
                </template>
            </v-data-table>

            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn color="secondary" variant="tonal"
                           :disabled="loading" @click="close">{{ tt('Close') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <amount-input-dialog ref="amountInputDialog" />
    <edit-dialog ref="editDialog" :type="TransactionEditPageType.Transaction" />

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import PaginationButtons from '@/components/desktop/PaginationButtons.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import AmountInputDialog from '@/components/desktop/AmountInputDialog.vue';
import EditDialog from '@/views/desktop/transactions/list/dialogs/EditDialog.vue';
import { TransactionEditPageType } from '@/views/base/transactions/TransactionEditPageBase.ts';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useReconciliationStatementPageBase } from '@/views/base/accounts/ReconciliationStatementPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import { TransactionType } from '@/core/transaction.ts';
import { KnownFileType } from '@/core/file.ts';
import { Transaction, type TransactionReconciliationStatementResponseItem } from '@/models/transaction.ts';

import { getCurrentUnixTime } from '@/lib/datetime.ts';
import { startDownloadFile } from '@/lib/ui/common.ts';

import {
    mdiRefresh,
    mdiArrowRight,
    mdiDotsVertical,
    mdiInvoiceTextPlusOutline,
    mdiInvoiceTextEditOutline,
    mdiComma,
    mdiKeyboardTab,
    mdiPencilBoxOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;
type AmountInputDialogType = InstanceType<typeof AmountInputDialog>;
type EditDialogType = InstanceType<typeof EditDialog>;

interface ReconciliationStatementDialogTablePageOption {
    value: number;
    title: string;
}

const emit = defineEmits<{
    (e: 'error', message: string): void;
}>();

const { tt } = useI18n();

const {
    accountId,
    startTime,
    endTime,
    reconciliationStatements,
    currentTimezoneOffsetMinutes,
    allAccountsMap,
    allCategoriesMap,
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
    getDisplayDateTime,
    getDisplayTimezone,
    getDisplaySourceAmount,
    getDisplayDestinationAmount,
    getDisplayAccountBalance,
    getExportedData
} = useReconciliationStatementPageBase();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionsStore = useTransactionsStore();

const amountInputDialog = useTemplateRef<AmountInputDialogType>('amountInputDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const currentPage = ref<number>(1);
const countPerPage = ref<number>(10);

let rejectFunc: ((reason?: unknown) => void) | null = null;

const reconciliationStatementsTablePageOptions = computed<ReconciliationStatementDialogTablePageOption[]>(() => getTablePageOptions(reconciliationStatements.value?.transactions.length));

const totalPageCount = computed<number>(() => {
    if (!reconciliationStatements.value || !reconciliationStatements.value.transactions || reconciliationStatements.value.transactions.length < 1) {
        return 1;
    }

    let count = 0;

    for (let i = 0; i < reconciliationStatements.value.transactions.length; i++) {
        count++;
    }

    return Math.ceil(count / countPerPage.value);
});

const dataTableHeaders = computed<object[]>(() => {
    const headers: object[] = [];
    const accountBalanceName = isCurrentLiabilityAccount.value ? 'Account Outstanding Balance' : 'Account Balance';

    headers.push({ key: 'time', value: 'time', title: tt('Transaction Time'), sortable: true, nowrap: true });
    headers.push({ key: 'type', value: 'type', title: tt('Type'), sortable: true, nowrap: true });
    headers.push({ key: 'categoryId', value: 'categoryId', title: tt('Category'), sortable: true, nowrap: true });
    headers.push({ key: 'sourceAmount', value: 'sourceAmount', title: tt('Amount'), sortable: true, nowrap: true });
    headers.push({ key: 'sourceAccountId', value: 'sourceAccountId', title: tt('Account'), sortable: true, nowrap: true });
    headers.push({ key: 'accountBalance', value: 'accountBalance', title: tt(accountBalanceName), sortable: true, nowrap: true });
    headers.push({ key: 'comment', value: 'comment', title: tt('Description'), sortable: true, nowrap: true });
    headers.push({ key: 'operation', title: tt('Operation'), sortable: false, nowrap: true, align: 'end' });
    return headers;
});

function getTablePageOptions(linesCount?: number): ReconciliationStatementDialogTablePageOption[] {
    const pageOptions: ReconciliationStatementDialogTablePageOption[] = [];

    if (!linesCount || linesCount < 1) {
        pageOptions.push({ value: -1, title: tt('All') });
        return pageOptions;
    }

    const availableCountPerPage = [ 5, 10, 15, 20, 25, 30, 50 ];

    for (let i = 0; i < availableCountPerPage.length; i++) {
        const count = availableCountPerPage[i];

        if (linesCount < count) {
            break;
        }

        pageOptions.push({ value: count, title: count.toString() });
    }

    pageOptions.push({ value: -1, title: tt('All') });

    return pageOptions;
}

function open(options: { accountId: string, startTime: number, endTime: number }): Promise<void> {
    accountId.value = options.accountId;
    startTime.value = options.startTime;
    endTime.value = options.endTime;
    reconciliationStatements.value = undefined;
    currentPage.value = 1;
    countPerPage.value = 10;
    showState.value = true;
    loading.value = true;

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false })
    ]).then(() => {
        return transactionsStore.getReconciliationStatements({
            accountId: options.accountId,
            startTime: options.startTime,
            endTime: options.endTime
        });
    }).then(result => {
        reconciliationStatements.value = result;
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            emit('error', error);
            showState.value = false;
        }
    });

    return new Promise<void>((resolve, reject) => {
        rejectFunc = reject;
    });
}

function reload(): void {
    loading.value = true;

    transactionsStore.getReconciliationStatements({
        accountId: accountId.value,
        startTime: startTime.value,
        endTime: endTime.value
    }).then(result => {
        reconciliationStatements.value = result;
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    })
}

function addTransaction(): void {
    editDialog.value?.open({
        accountId: accountId.value
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        reload();
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function updateClosingBalance(): void {
    let currentClosingBalance = reconciliationStatements.value?.closingBalance ?? 0;

    if (isCurrentLiabilityAccount.value) {
        currentClosingBalance = -currentClosingBalance;
    }

    amountInputDialog.value?.open({
        text: tt('Please enter the new closing balance for the account'),
        inputLabel: tt('Closing Balance'),
        inputPlaceholder: tt('Closing Balance'),
        currency: currentAccountCurrency.value,
        initAmount: currentClosingBalance
    }).then(newClosingBalance => {
        if (!newClosingBalance) {
            return;
        }

        const currentUnixTime = getCurrentUnixTime();
        let newTransactionTime: number | undefined = undefined;

        if (endTime.value < currentUnixTime) {
            newTransactionTime = endTime.value;
        } else if (currentUnixTime < startTime.value) {
            newTransactionTime = startTime.value;
        }

        let newTransactionType: TransactionType = isCurrentLiabilityAccount.value ? TransactionType.Expense : TransactionType.Income;
        let newTransactionAmount: number = newClosingBalance - currentClosingBalance;

        if (newTransactionAmount < 0) {
            newTransactionType = isCurrentLiabilityAccount.value ? TransactionType.Income : TransactionType.Expense;
            newTransactionAmount = -newTransactionAmount;
        }

        editDialog.value?.open({
            time: newTransactionTime,
            type: newTransactionType,
            amount: newTransactionAmount,
            accountId: accountId.value,
            noTransactionDraft: true
        }).then(result => {
            if (result && result.message) {
                snackbar.value?.showMessage(result.message);
            }

            reload();
        }).catch(error => {
            if (error) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function exportReconciliationStatements(fileType: KnownFileType): void {
    if (!reconciliationStatements.value || !reconciliationStatements.value.transactions || reconciliationStatements.value.transactions.length < 1) {
        return;
    }

    const exportedData = getExportedData(fileType);
    startDownloadFile(fileType.formatFileName(exportFileName.value), fileType.createBlob(exportedData));
}

function showTransaction(transaction: TransactionReconciliationStatementResponseItem): void {
    if (transaction.type === TransactionType.ModifyBalance) {
        return;
    }

    editDialog.value?.open({
        id: transaction.id,
        currentTransaction: Transaction.of(transaction)
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        reload();
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function close(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
