<template>
    <v-card class="overview-widget overview-widget--list h-100" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || tt('Recent Transactions')" :icon="mdiHistory" />
        </template>

        <v-card-text class="overview-widget__body overview-widget__list-body">
            <v-list class="py-0" density="compact" v-if="transactions.length">
                <v-list-item class="overview-widget__list-item" :key="transaction.id"
                             :title="getTransactionCategoryName(transaction)"
                             v-for="transaction in transactions"
                             @click="showTransaction(transaction)">
                    <template #prepend>
                        <span class="overview-widget__item-icon">
                            <ItemIcon size="24px" :icon-type="getCategoryIconType(getTransactionCategory(transaction)?.iconType)"
                                      :icon-id="getTransactionCategory(transaction)?.icon ?? ''"
                                      :color="getTransactionCategory(transaction)?.color"
                                      v-if="getTransactionCategory(transaction) && getTransactionCategory(transaction)?.color" />
                            <v-icon size="24" :icon="mdiPencilBoxOutline" v-else-if="!getTransactionCategory(transaction) || !getTransactionCategory(transaction)?.color" />
                        </span>
                    </template>

                    <template #append>
                        <div class="overview-widget__list-amount overview-widget__amount" :class="getAmountClass(transaction)">{{ getDisplayTransactionAmount(transaction) }}</div>
                    </template>

                    <div class="text-truncate overview-widget__caption text-body-small mt-h1" :title="getDisplayDescription(transaction)">{{ getDisplayDescription(transaction) }}</div>
                </v-list-item>
            </v-list>
            <div v-if="loading && !transactions.length">
                <v-skeleton-loader class="skeleton-no-margin mx-2 py-4 my-1" type="text" :key="idx" :loading="true" v-for="idx in props.itemCount"></v-skeleton-loader>
            </div>
            <div class="overview-widget__empty" v-else-if="!loading && !transactions.length">
                <v-icon :icon="mdiHistory" size="32" />
                <span>{{ tt('No data') }}</span>
            </div>
        </v-card-text>

        <edit-dialog ref="editDialog" :type="TransactionEditPageType.Transaction" />

        <snack-bar ref="snackbar" />
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import SnackBar from '@/components/desktop/SnackBar.vue';
import EditDialog from '@/views/desktop/transactions/list/dialogs/EditDialog.vue';

import { computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { TransactionEditPageType } from '@/views/base/transactions/TransactionEditPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { TransactionType } from '@/core/transaction.ts';
import { DISPLAY_HIDDEN_AMOUNT } from '@/consts/numeral.ts';
import type { OverviewRecentTransactionsQuery } from '@/core/overview_layout.ts';

import type { TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionInfoResponse, Transaction } from '@/models/transaction.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';
import { getCategoryIconType } from '@/lib/icon.ts';
import { getOverviewRecentTransactionsQuery } from '@/lib/overview_layout.ts';

import {
    mdiHistory,
    mdiPencilBoxOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;
type EditDialogType = InstanceType<typeof EditDialog>;

const props = defineProps<{
    loading: boolean;
    title?: string;
    editing?: boolean;
    itemCount: number;
    accountIds?: string[];
    categoryIds?: string[];
    tagFilter?: string;
    amountFilter?: string;
    keyword?: string;
}>();

const emit = defineEmits<{
    (e: 'refresh'): void
}>();

const { tt, formatDateTimeToShortDateTime, formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const overviewStore = useOverviewStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
const recentTransactionsQuery = computed<OverviewRecentTransactionsQuery>(() => getOverviewRecentTransactionsQuery({
    accountIds: props.accountIds ?? [],
    categoryIds: props.categoryIds ?? [],
    tagFilter: props.tagFilter ?? '',
    amountFilter: props.amountFilter ?? '',
    keyword: props.keyword ?? ''
}));
const transactions = computed<TransactionInfoResponse[]>(() => overviewStore.getRecentTransactions(recentTransactionsQuery.value).slice(0, props.itemCount));

function getDisplayDescription(transaction: TransactionInfoResponse): string {
    const accountName = accountsStore.allAccountsMap[transaction.sourceAccountId]?.name ?? '';
    return `${formatDateTimeToShortDateTime(parseDateTimeFromUnixTime(transaction.time))}${accountName ? ` · ${accountName}` : ''}${transaction.comment ? ` · ${transaction.comment}` : ''}`;
}

function getDisplayTransactionAmount(transaction: TransactionInfoResponse): string {
    const amount = !showAmountInHomePage.value || transaction.hideAmount ? DISPLAY_HIDDEN_AMOUNT : parseBigDecimal(transaction.sourceAmount);
    return formatAmountToLocalizedNumeralsWithCurrency(amount, accountsStore.allAccountsMap[transaction.sourceAccountId]?.currency ?? '');
}

function getTransactionCategory(transaction: TransactionInfoResponse): TransactionCategory | undefined {
    if (transaction.type === TransactionType.ModifyBalance) {
        return undefined;
    }

    return transactionCategoriesStore.allTransactionCategoriesMap[transaction.categoryId];
}

function getTransactionCategoryName(transaction: TransactionInfoResponse): string {
    if (transaction.type === TransactionType.ModifyBalance) {
        return tt('Modify Balance');
    }

    const category = getTransactionCategory(transaction);

    if (category) {
        return category.name;
    } else if (transaction.type === TransactionType.Income) {
        return tt('Income');
    } else if (transaction.type === TransactionType.Expense) {
        return tt('Expense');
    } else if (transaction.type === TransactionType.Transfer) {
        return tt('Transfer');
    } else {
        return tt('Transaction');
    }
}

function getAmountClass(transaction: TransactionInfoResponse): string {
    if (transaction.type === TransactionType.Income) {
        return 'text-income';
    } else if (transaction.type === TransactionType.Expense) {
        return 'text-expense';
    } else {
        return '';
    }
}

function showTransaction(transaction: TransactionInfoResponse): void {
    if (props.editing) {
        return;
    }

    editDialog.value?.open({
        id: transaction.id,
        currentTransaction: Transaction.of(transaction)
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        emit('refresh');
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}
</script>
