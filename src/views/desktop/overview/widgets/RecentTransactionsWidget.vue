<template>
    <v-card class="h-100" :class="{ disabled: loading }">
        <template #title>
            <span class="text-title-medium">{{ title || tt('Recent Transactions') }}</span>
        </template>

        <v-card-text class="pt-0 px-2">
            <v-list class="py-0" density="compact" v-if="transactions.length">
                <v-list-item class="ps-2 pe-3 mb-h1" :key="transaction.id"
                             :title="getTransactionCategoryName(transaction)"
                             v-for="transaction in transactions"
                             @click="showTransaction(transaction)">
                    <template #prepend>
                        <ItemIcon class="me-2" size="24px" :icon-type="getCategoryIconType(getTransactionCategory(transaction)?.iconType)"
                                  :icon-id="getTransactionCategory(transaction)?.icon ?? ''"
                                  :color="getTransactionCategory(transaction)?.color"
                                  v-if="getTransactionCategory(transaction) && getTransactionCategory(transaction)?.color" />
                        <v-icon size="24" :icon="mdiPencilBoxOutline" v-else-if="!getTransactionCategory(transaction) || !getTransactionCategory(transaction)?.color" />
                    </template>

                    <template #append>
                        <div :class="getAmountClass(transaction)">{{ getDisplayTransactionAmount(transaction) }}</div>
                    </template>

                    <div class="text-truncate text-medium-emphasis mt-h1">{{ getDisplayDescription(transaction) }}</div>
                </v-list-item>
            </v-list>
            <div v-if="loading && !transactions.length">
                <v-skeleton-loader class="py-2" type="text" :key="idx" :loading="true" v-for="idx in props.itemCount"></v-skeleton-loader>
            </div>
            <div class="text-medium-emphasis text-center pt-4" v-else-if="!loading && !transactions.length">{{ tt('No data') }}</div>
        </v-card-text>

        <edit-dialog ref="editDialog" :type="TransactionEditPageType.Transaction" />

        <snack-bar ref="snackbar" />
    </v-card>
</template>

<script setup lang="ts">
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

import type { TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionInfoResponse, Transaction } from '@/models/transaction.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';
import { getCategoryIconType } from '@/lib/icon.ts';

import {
    mdiPencilBoxOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;
type EditDialogType = InstanceType<typeof EditDialog>;

const props = defineProps<{
    loading: boolean;
    title?: string;
    itemCount: number;
    editing?: boolean
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
const transactions = computed<TransactionInfoResponse[]>(() => overviewStore.recentTransactions.slice(0, props.itemCount));

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
