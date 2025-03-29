<template>
    <v-dialog width="600" :persistent="!!loading || (mode === 'replaceInvalidItems' && !!sourceItem) || !!targetItem" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'expenseCategory'">{{ tt('Batch Replace Selected Expense Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'incomeCategory'">{{ tt('Batch Replace Selected Income Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'transferCategory'">{{ tt('Batch Replace Selected Transfer Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'account'">{{ tt('Batch Replace Selected Accounts') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'batchReplace' && type === 'destinationAccount'">{{ tt('Batch Replace Selected Destination Accounts') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'expenseCategory'">{{ tt('Replace Invalid Expense Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'incomeCategory'">{{ tt('Replace Invalid Income Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'transferCategory'">{{ tt('Replace Invalid Transfer Categories') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'account'">{{ tt('Replace Invalid Accounts') }}</h4>
                    <h4 class="text-h4" v-if="mode === 'replaceInvalidItems' && type === 'tag'">{{ tt('Replace Invalid Transaction Tags') }}</h4>
                    <v-btn density="compact" color="default" variant="text" size="24"
                           class="ml-2" :icon="true" :disabled="loading"
                           :loading="loading" @click="reload">
                        <template #loader>
                            <v-progress-circular indeterminate size="20"/>
                        </template>
                        <v-icon :icon="mdiRefresh" size="24" />
                        <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="my-md-4 w-100 d-flex justify-center" v-if="type === 'expenseCategory' || type === 'incomeCategory' || type === 'transferCategory'">
                <v-row>
                    <v-col cols="12" v-if="mode === 'replaceInvalidItems'">
                        <v-autocomplete
                            item-title="name"
                            item-value="value"
                            persistent-placeholder
                            :disabled="loading"
                            :label="tt('Invalid Category')"
                            :placeholder="tt('Invalid Category')"
                            :items="invalidItems"
                            :no-data-text="tt('No available category')"
                            v-model="sourceItem">
                        </v-autocomplete>
                    </v-col>
                    <v-col cols="12">
                        <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :disabled="loading || !hasAvailableExpenseCategories"
                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                           :show-selection-primary-text="true"
                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(targetItem, allCategories[CategoryType.Expense])"
                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(targetItem, allCategories[CategoryType.Expense])"
                                           :label="tt('Target Category')"
                                           :placeholder="tt('Target Category')"
                                           :items="allCategories[CategoryType.Expense]"
                                           v-model="targetItem"
                                           v-if="type === 'expenseCategory'">
                        </two-column-select>
                        <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :disabled="loading || !hasAvailableIncomeCategories"
                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                           :show-selection-primary-text="true"
                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(targetItem, allCategories[CategoryType.Income])"
                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(targetItem, allCategories[CategoryType.Income])"
                                           :label="tt('Target Category')"
                                           :placeholder="tt('Target Category')"
                                           :items="allCategories[CategoryType.Income]"
                                           v-model="targetItem"
                                           v-if="type === 'incomeCategory'">
                        </two-column-select>
                        <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :disabled="loading || !hasAvailableTransferCategories"
                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                           :show-selection-primary-text="true"
                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(targetItem, allCategories[CategoryType.Transfer])"
                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(targetItem, allCategories[CategoryType.Transfer])"
                                           :label="tt('Target Category')"
                                           :placeholder="tt('Target Category')"
                                           :items="allCategories[CategoryType.Transfer]"
                                           v-model="targetItem"
                                           v-if="type === 'transferCategory'">
                        </two-column-select>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="my-md-4 w-100 d-flex justify-center" v-if="type === 'account' || type === 'destinationAccount'">
                <v-row>
                    <v-col cols="12" v-if="mode === 'replaceInvalidItems'">
                        <v-autocomplete
                            item-title="name"
                            item-value="value"
                            persistent-placeholder
                            :disabled="loading"
                            :label="tt('Invalid Account')"
                            :placeholder="tt('Invalid Account')"
                            :items="invalidItems"
                            :no-data-text="tt('No available account')"
                            v-model="sourceItem">
                        </v-autocomplete>
                    </v-col>
                    <v-col cols="12">
                        <two-column-select primary-key-field="id" primary-value-field="category"
                                           primary-title-field="name" primary-footer-field="displayBalance"
                                           primary-icon-field="icon" primary-icon-type="account"
                                           primary-sub-items-field="accounts"
                                           :primary-title-i18n="true"
                                           secondary-key-field="id" secondary-value-field="id"
                                           secondary-title-field="name" secondary-footer-field="displayBalance"
                                           secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                           :disabled="loading || !allVisibleAccounts.length"
                                           :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                           :custom-selection-primary-text="getAccountDisplayName(targetItem)"
                                           :label="tt('Target Account')"
                                           :placeholder="tt('Target Account')"
                                           :items="allVisibleCategorizedAccounts"
                                           v-model="targetItem">
                        </two-column-select>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="my-md-4 w-100 d-flex justify-center" v-if="type === 'tag'">
                <v-row>
                    <v-col cols="12" v-if="mode === 'replaceInvalidItems'">
                        <v-autocomplete
                            item-title="name"
                            item-value="value"
                            persistent-placeholder
                            :disabled="loading"
                            :label="tt('Invalid Tag')"
                            :placeholder="tt('Invalid Tag')"
                            :items="invalidItems"
                            :no-data-text="tt('No available tag')"
                            v-model="sourceItem">
                        </v-autocomplete>
                    </v-col>
                    <v-col cols="12">
                        <v-autocomplete
                            item-title="name"
                            item-value="id"
                            persistent-placeholder
                            chips
                            :disabled="loading"
                            :label="tt('Target Tag')"
                            :placeholder="tt('Target Tag')"
                            :items="allTags"
                            :no-data-text="tt('No available tag')"
                            v-model="targetItem"
                        >
                            <template #chip="{ props, item }">
                                <v-chip :prepend-icon="mdiPound" :text="item.title" v-bind="props"/>
                            </template>

                            <template #item="{ props, item }">
                                <v-list-item :value="item.value" v-bind="props" v-if="!item.raw.hidden">
                                    <template #title>
                                        <v-list-item-title>
                                            <div class="d-flex align-center">
                                                <v-icon size="20" start :icon="mdiPound"/>
                                                <span>{{ item.title }}</span>
                                            </div>
                                        </v-list-item-title>
                                    </template>
                                </v-list-item>
                            </template>
                        </v-autocomplete>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="loading || (mode === 'replaceInvalidItems' && !sourceItem && sourceItem !== '') || (!targetItem && targetItem !== '')" @click="confirm">{{ tt('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="loading" @click="cancel">{{ tt('Cancel') }}</v-btn>
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

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import type { NameValue } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';
import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';

import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName
} from '@/lib/category.ts';

import {
    mdiRefresh,
    mdiPound
} from '@mdi/js';

export type BatchReplaceDialogMode = 'batchReplace' | 'replaceInvalidItems';
export type BatchReplaceDialogDataType = 'expenseCategory' | 'incomeCategory' | 'transferCategory' | 'account' | 'destinationAccount' | 'tag';

type SnackBarType = InstanceType<typeof SnackBar>;

interface BatchReplaceDialogResponse {
    sourceItem?: string;
    targetItem?: string;
}

const { tt, getCategorizedAccountsWithDisplayBalance } = useI18n();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const mode = ref<BatchReplaceDialogMode | ''>('');
const type = ref<BatchReplaceDialogDataType | ''>('');
const invalidItems = ref<NameValue[] | undefined>([]);
const sourceItem = ref<string | undefined>(undefined);
const targetItem = ref<string | undefined>(undefined);

let resolveFunc: ((response: BatchReplaceDialogResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value));
const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
const allTags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);

const hasAvailableExpenseCategories = computed<boolean>(() => transactionCategoriesStore.hasAvailableExpenseCategories);
const hasAvailableIncomeCategories = computed<boolean>(() => transactionCategoriesStore.hasAvailableIncomeCategories);
const hasAvailableTransferCategories = computed<boolean>(() => transactionCategoriesStore.hasAvailableTransferCategories);

function getAccountDisplayName(accountId?: string): string {
    if (accountId) {
        return Account.findAccountNameById(allAccounts.value, accountId) || '';
    } else {
        return tt('None');
    }
}

function open(options: { mode: BatchReplaceDialogMode; type: BatchReplaceDialogDataType; invalidItems?: NameValue[] }): Promise<BatchReplaceDialogResponse> {
    mode.value = options.mode;
    type.value = options.type;
    sourceItem.value = undefined;

    if (mode.value === 'batchReplace') {
        invalidItems.value = undefined;
    } else if (mode.value === 'replaceInvalidItems') {
        invalidItems.value = options.invalidItems;
    }

    targetItem.value = undefined;
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function reload(): void {
    if (type.value === 'expenseCategory' || type.value === 'incomeCategory' || type.value === 'transferCategory') {
        loading.value = true;

        transactionCategoriesStore.loadAllCategories({ force: true }).then(() => {
            loading.value = false;
        }).catch(error => {
            loading.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    } else if (type.value === 'account' || type.value === 'destinationAccount') {
        loading.value = true;

        accountsStore.loadAllAccounts({ force: true }).then(() => {
            loading.value = false;
        }).catch(error => {
            loading.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    } else if (type.value === 'tag') {
        loading.value = true;

        transactionTagsStore.loadAllTags({ force: true }).then(() => {
            loading.value = false;
        }).catch(error => {
            loading.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }
}

function confirm(): void {
    if (mode.value === 'batchReplace') {
        resolveFunc?.({
            targetItem: targetItem.value
        });
    } else if (mode.value === 'replaceInvalidItems') {
        resolveFunc?.({
            sourceItem: sourceItem.value,
            targetItem: targetItem.value
        });
    }

    showState.value = false;
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
