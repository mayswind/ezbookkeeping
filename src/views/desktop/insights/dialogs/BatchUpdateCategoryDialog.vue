<template>
    <v-dialog width="600" :persistent="true" v-model="showState">
        <one-column-dialog-layout content-class="pt-0" :disabled="loading || submitting"
                                  :title="title" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #after-title>
                <v-btn density="compact" color="default" variant="text"
                       class="ms-2" :icon="true" :disabled="loading || submitting"
                       :loading="loading" @click="reload">
                    <template #loader>
                        <v-progress-circular indeterminate size="20"/>
                    </template>
                    <v-icon :icon="mdiRefresh" size="22" />
                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                </v-btn>
            </template>

            <template #content>
                <div class="mt-5">
                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                       secondary-hidden-field="hidden"
                                       :disabled="loading || submitting || !hasVisibleExpenseCategories"
                                       :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                       :show-selection-primary-text="true"
                                       :custom-selection-primary-text="getTransactionPrimaryCategoryName(categoryId, allCategories[CategoryType.Expense])"
                                       :custom-selection-secondary-text="getTransactionSecondaryCategoryName(categoryId, allCategories[CategoryType.Expense])"
                                       :label="tt('Expense Category')"
                                       :placeholder="tt('Expense Category')"
                                       :items="allCategories[CategoryType.Expense]"
                                       v-model="categoryId"
                                       v-if="type === CategoryType.Expense">
                    </two-column-select>
                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                       secondary-hidden-field="hidden"
                                       :disabled="loading || submitting || !hasVisibleIncomeCategories"
                                       :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                       :show-selection-primary-text="true"
                                       :custom-selection-primary-text="getTransactionPrimaryCategoryName(categoryId, allCategories[CategoryType.Income])"
                                       :custom-selection-secondary-text="getTransactionSecondaryCategoryName(categoryId, allCategories[CategoryType.Income])"
                                       :label="tt('Income Category')"
                                       :placeholder="tt('Income Category')"
                                       :items="allCategories[CategoryType.Income]"
                                       v-model="categoryId"
                                       v-if="type === CategoryType.Income">
                    </two-column-select>
                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                       secondary-hidden-field="hidden"
                                       :disabled="loading || submitting || !hasVisibleTransferCategories"
                                       :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                       :show-selection-primary-text="true"
                                       :custom-selection-primary-text="getTransactionPrimaryCategoryName(categoryId, allCategories[CategoryType.Transfer])"
                                       :custom-selection-secondary-text="getTransactionSecondaryCategoryName(categoryId, allCategories[CategoryType.Transfer])"
                                       :label="tt('Transfer Category')"
                                       :placeholder="tt('Transfer Category')"
                                       :items="allCategories[CategoryType.Transfer]"
                                       v-model="categoryId"
                                       v-if="type === CategoryType.Transfer">
                    </two-column-select>
                </div>
            </template>

            <template #footer>
                <v-btn color="secondary" variant="tonal" :disabled="loading || submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                <v-spacer/>
                <v-btn :disabled="loading || submitting || updateIds.length < 1 || !categoryId" @click="confirm">
                    {{ tt('OK') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                </v-btn>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import { CategoryType } from '@/core/category.ts';

import type { TransactionCategory } from '@/models/transaction_category.ts';

import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName
} from '@/lib/category.ts';

import {
    mdiRefresh
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt
} = useI18n();

const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionsStore = useTransactionsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const submitting = ref<boolean>(false);
const type = ref<CategoryType>(CategoryType.Expense);
const updateIds = ref<string[]>([]);
const categoryId = ref<string>('');

let resolveFunc: ((response: number) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);

const hasVisibleExpenseCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleExpenseCategories);
const hasVisibleIncomeCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleIncomeCategories);
const hasVisibleTransferCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleTransferCategories);

const title = computed<string>(() => {
    if (type.value === CategoryType.Expense) {
        return tt('Update Categories for Expense Transactions');
    } else if (type.value === CategoryType.Income) {
        return tt('Update Categories for Income Transactions');
    } else if (type.value === CategoryType.Transfer) {
        return tt('Update Categories for Transfer Transactions');
    } else {
        return '';
    }
});

function open(options: { type: CategoryType; updateIds: string[] }): Promise<number> {
    type.value = options.type;
    updateIds.value = options.updateIds;
    categoryId.value = '';
    submitting.value = false;
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function reload(): void {
    loading.value = true;

    transactionCategoriesStore.loadAllCategories({ force: true }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function confirm(): void {
    submitting.value = true;

    transactionsStore.batchUpdateTransactionCategories({
        transactionIds: updateIds.value,
        categoryId: categoryId.value
    }).then(() => {
        submitting.value = false;
        showState.value = false;
        resolveFunc?.(updateIds.value.length);
    }).catch(error => {
        submitting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
