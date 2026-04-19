<template>
    <v-dialog width="600" :persistent="true" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex flex-wrap align-center">
                    <h4 class="text-h4 text-wrap" v-if="type === CategoryType.Expense">{{ tt('Update Categories for Expense Transactions') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === CategoryType.Income">{{ tt('Update Categories for Income Transactions') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === CategoryType.Transfer">{{ tt('Update Categories for Transfer Transactions') }}</h4>
                    <v-btn class="ms-2" density="compact" color="default" variant="text" size="24"
                           :icon="true" :disabled="loading || submitting" :loading="loading"
                           @click="reload">
                        <template #loader>
                            <v-progress-circular indeterminate size="20"/>
                        </template>
                        <v-icon :icon="mdiRefresh" size="24" />
                        <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
                <v-row>
                    <v-col cols="12">
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
                                           :label="tt('Target Category')"
                                           :placeholder="tt('Target Category')"
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
                                           :label="tt('Target Category')"
                                           :placeholder="tt('Target Category')"
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
                                           :label="tt('Target Category')"
                                           :placeholder="tt('Target Category')"
                                           :items="allCategories[CategoryType.Transfer]"
                                           v-model="categoryId"
                                           v-if="type === CategoryType.Transfer">
                        </two-column-select>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :disabled="loading || submitting || updateIds.length < 1 || !categoryId" @click="confirm">
                        {{ tt('OK') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="loading || submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
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

function open(options: { type: CategoryType; updateIds: string[] }): Promise<number> {
    type.value = options.type;
    updateIds.value = options.updateIds;
    categoryId.value = '';
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function reload(): void {
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
