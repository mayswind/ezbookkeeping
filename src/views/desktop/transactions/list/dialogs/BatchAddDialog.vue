<template>
    <v-dialog :persistent="isModified" v-model="showState">
        <one-column-dialog-layout content-class="pa-0"
                                  :title="tt('Add Multiple Transactions')" :cancel-button-title="tt('Close')"
                                  :disabled="loading || submitting" :loading="loading"
                                  @cancel="cancel">
            <template #after-title>
                <v-btn density="compact" color="primary" variant="outlined" class="ms-2"
                       :disabled="loading || submitting || rows.length >= 10"
                       @click="addTransactionRow()">{{ tt('Add') }}</v-btn>
            </template>

            <template #content>
                <v-table fixed-header density="compact"
                         class="w-100 table-striped batch-add-transaction-table" height="min(65dvh, 416px)">
                    <thead>
                    <tr>
                        <th class="batch-add-time-column">
                            <div class="d-inline-block h-100 align-content-center">
                                {{ tt('Transaction Time') }}
                            </div>
                        </th>
                        <th class="batch-add-type-column">
                            <div class="d-inline-block h-100 align-content-center batch-add-component-header">
                                <v-select density="compact" variant="plain"
                                          item-title="displayName" item-value="type"
                                          :disabled="loading || submitting"
                                          :items="allAvailableTransactionTypes"
                                          v-model="globalTransactionType">
                                    <template #selection>
                                        <span>{{ tt('Type') }}</span>
                                    </template>
                                </v-select>
                            </div>
                        </th>
                        <th class="batch-add-category-column">
                            <div class="d-inline-block h-100 align-content-center batch-add-component-header" v-if="allTransactionSameType">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type-field="iconType" primary-icon-type="category" primary-color-field="color"
                                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type-field="iconType" secondary-icon-type="category" secondary-color-field="color"
                                                   secondary-hidden-field="hidden"
                                                   :disabled="loading || submitting || !hasVisibleExpenseCategories"
                                                   :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                   :show-selection-primary-text="true"
                                                   :custom-selection-primary-text="tt('Category')"
                                                   :placeholder="tt('Category')"
                                                   :items="allCategories[CategoryType.Expense] || []"
                                                   v-model="globalCategoryId"
                                                   v-if="globalTransactionType === TransactionType.Expense">
                                </two-column-select>

                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type-field="iconType" primary-icon-type="category" primary-color-field="color"
                                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type-field="iconType" secondary-icon-type="category" secondary-color-field="color"
                                                   secondary-hidden-field="hidden"
                                                   :disabled="loading || submitting || !hasVisibleIncomeCategories"
                                                   :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                   :show-selection-primary-text="true"
                                                   :custom-selection-primary-text="tt('Category')"
                                                   :placeholder="tt('Category')"
                                                   :items="allCategories[CategoryType.Income] || []"
                                                   v-model="globalCategoryId"
                                                   v-else-if="globalTransactionType === TransactionType.Income">
                                </two-column-select>
                            </div>

                            <div class="d-inline-block h-100 align-content-center" v-else-if="!allTransactionSameType">
                                {{ tt('Category') }}
                            </div>
                        </th>
                        <th class="batch-add-account-column">
                            <div class="d-inline-block h-100 align-content-center batch-add-component-header">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="category"
                                                   primary-title-field="name" primary-footer-field="displayBalance"
                                                   primary-icon-field="icon" primary-icon-type-field="iconType" primary-icon-type="account"
                                                   primary-sub-items-field="accounts"
                                                   :primary-title-i18n="true"
                                                   secondary-key-field="id" secondary-value-field="id"
                                                   secondary-title-field="name" secondary-footer-field="displayBalance"
                                                   secondary-icon-field="icon" secondary-icon-type-field="iconType" secondary-icon-type="account" secondary-color-field="color"
                                                   :disabled="loading || submitting || !allVisibleAccounts.length"
                                                   :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                   :custom-selection-primary-text="tt('Account')"
                                                   :placeholder="tt('Account')"
                                                   :items="allVisibleCategorizedAccounts"
                                                   v-model="globalAccountId">
                                </two-column-select>
                            </div>
                        </th>
                        <th class="batch-add-amount-column">
                            <div class="d-inline-block h-100 align-content-center">
                                {{ tt('Amount') }}
                            </div>
                        </th>
                        <th class="batch-add-tags-column">
                            <div class="d-inline-block h-100 align-content-center batch-add-component-header">
                                <transaction-tag-auto-complete select-only density="compact" variant="plain"
                                                               :disabled="loading || submitting"
                                                               :custom-selection-text="tt('Tags')"
                                                               v-model="globalTagIds" />
                            </div>
                        </th>
                        <th class="batch-add-description-column">
                            <div class="d-inline-block h-100 align-content-center">
                                {{ tt('Description') }}
                            </div>
                        </th>
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr :key="index" v-for="(row, index) in rows">
                        <td>
                            <date-time-select density="compact" variant="plain"
                                              :disabled="loading || submitting"
                                              :timezone-utc-offset="row.transaction.utcOffset"
                                              :model-value="row.transaction.time"
                                              @update:model-value="updateTime(row.transaction, $event)"
                                              @error="snackbar?.showError($event)" />
                        </td>
                        <td>
                            <v-select density="compact" variant="plain"
                                      item-title="displayName" item-value="type"
                                      :disabled="loading || submitting"
                                      :items="allAvailableTransactionTypes"
                                      :model-value="row.transaction.type"
                                      @update:model-value="updateType(row.transaction, $event)" />
                        </td>
                        <td>
                            <two-column-select density="compact" variant="plain"
                                               primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                               primary-icon-field="icon" primary-icon-type-field="iconType" primary-icon-type="category" primary-color-field="color"
                                               primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                               secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                               secondary-icon-field="icon" secondary-icon-type-field="iconType" secondary-icon-type="category" secondary-color-field="color"
                                               secondary-hidden-field="hidden"
                                               :disabled="loading || submitting || !hasVisibleExpenseCategories"
                                               :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                               :show-selection-primary-text="true"
                                               :custom-selection-primary-text="getTransactionPrimaryCategoryName(row.transaction.expenseCategoryId, allCategories[CategoryType.Expense])"
                                               :custom-selection-secondary-text="getTransactionSecondaryCategoryName(row.transaction.expenseCategoryId, allCategories[CategoryType.Expense])"
                                               :placeholder="tt('Category')"
                                               :items="allCategories[CategoryType.Expense] || []"
                                               v-model="row.transaction.expenseCategoryId"
                                               v-if="row.transaction.type === TransactionType.Expense">
                            </two-column-select>

                            <two-column-select density="compact" variant="plain"
                                               primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                               primary-icon-field="icon" primary-icon-type-field="iconType" primary-icon-type="category" primary-color-field="color"
                                               primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                               secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                               secondary-icon-field="icon" secondary-icon-type-field="iconType" secondary-icon-type="category" secondary-color-field="color"
                                               secondary-hidden-field="hidden"
                                               :disabled="loading || submitting || !hasVisibleIncomeCategories"
                                               :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                               :show-selection-primary-text="true"
                                               :custom-selection-primary-text="getTransactionPrimaryCategoryName(row.transaction.incomeCategoryId, allCategories[CategoryType.Income])"
                                               :custom-selection-secondary-text="getTransactionSecondaryCategoryName(row.transaction.incomeCategoryId, allCategories[CategoryType.Income])"
                                               :placeholder="tt('Category')"
                                               :items="allCategories[CategoryType.Income] || []"
                                               v-model="row.transaction.incomeCategoryId"
                                               v-else-if="row.transaction.type === TransactionType.Income">
                            </two-column-select>
                        </td>
                        <td>
                            <two-column-select density="compact" variant="plain"
                                               primary-key-field="id" primary-value-field="category"
                                               primary-title-field="name" primary-footer-field="displayBalance"
                                               primary-icon-field="icon" primary-icon-type-field="iconType" primary-icon-type="account"
                                               primary-sub-items-field="accounts"
                                               :primary-title-i18n="true"
                                               secondary-key-field="id" secondary-value-field="id"
                                               secondary-title-field="name" secondary-footer-field="displayBalance"
                                               secondary-icon-field="icon" secondary-icon-type-field="iconType" secondary-icon-type="account" secondary-color-field="color"
                                               :disabled="loading || submitting || !allVisibleAccounts.length"
                                               :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                               :custom-selection-primary-text="Account.findAccountNameById(allAccounts, row.transaction.sourceAccountId) || ''"
                                               :placeholder="tt('Account')"
                                               :items="allVisibleCategorizedAccounts"
                                               v-model="row.transaction.sourceAccountId">
                            </two-column-select>
                        </td>
                        <td>
                            <amount-input density="compact" variant="plain"
                                          persistent-placeholder
                                          :color="getAccountColor(row.transaction)"
                                          :currency="getAccountCurrency(row.transaction)"
                                          :show-currency="true"
                                          :disabled="loading || submitting"
                                          :hide="row.transaction.hideAmount"
                                          :placeholder="tt('Amount')"
                                          :enable-formula="true"
                                          v-model="row.transaction.sourceAmount"/>
                        </td>
                        <td>
                            <transaction-tag-auto-complete density="compact" variant="plain"
                                                           :disabled="loading || submitting"
                                                           v-model="row.transaction.tagIds" />
                        </td>
                        <td>
                            <v-text-field density="compact" variant="plain"
                                          type="text" autocomplete="off" persistent-placeholder
                                          :disabled="loading || submitting"
                                          :placeholder="tt('Your transaction description (optional)')"
                                          v-model="row.transaction.comment" />
                        </td>
                        <td>
                            <v-btn class="px-2" color="default"
                                   density="comfortable" variant="text"
                                   :prepend-icon="mdiContentCopy"
                                   :disabled="loading || submitting || rows.length >= 10"
                                   @click="duplicateTransactionRow(row.transaction, index)">
                                {{ tt('Duplicate') }}
                            </v-btn>
                            <v-btn class="px-2" color="default"
                                   density="comfortable" variant="text"
                                   :prepend-icon="mdiDeleteOutline"
                                   :disabled="loading || submitting || rows.length <= 1"
                                   @click="removeTransactionRow(index)">
                                {{ tt('Delete') }}
                            </v-btn>
                        </td>
                    </tr>
                    </tbody>
                </v-table>
            </template>

            <template #footer>
                <span class="align-self-center text-body-medium">
                    {{ tt('format.misc.batchTransactionCount', { count: formatNumberToLocalizedNumerals(rows.length) }) }}
                </span>
                <v-spacer />
                <v-btn :disabled="loading || submitting || !rows.length" @click="submit">
                    {{ tt('Save') }}
                    <v-progress-circular indeterminate size="20" class="ms-2" v-if="submitting" />
                </v-btn>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <confirm-dialog ref="confirmDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { TransactionType } from '@/core/transaction.ts';
import { CategoryType } from '@/core/category.ts';

import { type CategorizedAccountWithDisplayBalance,  Account } from '@/models/account.ts';
import { Transaction } from '@/models/transaction.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';

import { isArray1SubsetOfArray2 } from '@/lib/common.ts';
import {
    getTimezoneOffsetMinutes,
    getSameDateTimeWithCurrentTimezone,
    getCurrentUnixTime,
    parseDateTimeFromUnixTimeWithBrowserTimezone
} from '@/lib/datetime.ts';
import { type SetTransactionOptions, setTransactionModelByTransaction } from '@/lib/transaction.ts';
import { getTransactionPrimaryCategoryName, getTransactionSecondaryCategoryName } from '@/lib/category.ts';
import logger from '@/lib/logger.ts';

import {
    mdiContentCopy,
    mdiDeleteOutline
} from '@mdi/js';

interface BatchAddTransactionRow {
    transaction: Transaction;
}

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    formatNumberToLocalizedNumerals,
    getCategorizedAccountsWithDisplayBalance
} = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionsStore = useTransactionsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

let resolveFunc: ((changed: boolean) => void) | null = null;
let rejectFunc: ((reason?: string) => void) | null = null;

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const submitting = ref<boolean>(false);
const initOptions = ref<SetTransactionOptions | undefined>(undefined);
const rows = ref<BatchAddTransactionRow[]>([]);

const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);

const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value, customAccountCategoryOrder.value));
const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
const allCategoriesMap = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);
const allTagsMap = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);
const firstVisibleAccountId = computed<string | undefined>(() => allVisibleAccounts.value && allVisibleAccounts.value[0] ? allVisibleAccounts.value[0].id : undefined);

const hasVisibleExpenseCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleExpenseCategories);
const hasVisibleIncomeCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleIncomeCategories);

const allAvailableTransactionTypes = computed<TypeAndDisplayName[]>(() => [
    { type: TransactionType.Expense, displayName: tt('Expense') },
    { type: TransactionType.Income, displayName: tt('Income') }
]);

const isModified = computed<boolean>(() => {
    if (rows.value.length > 1) {
        return true;
    }

    return transactionsStore.isTransactionDraftModified(rows.value[0]?.transaction, initOptions.value?.amount, initOptions.value?.categoryId, initOptions.value?.accountId, initOptions.value?.tagIds, firstVisibleAccountId.value);
});

const allTransactionSameType = computed<boolean>(() => {
    const firstTransactionType = rows.value[0]?.transaction.type;

    for (const row of rows.value) {
        if (row.transaction.type !== firstTransactionType) {
            return false;
        }
    }

    return true;
});

const globalTransactionType = computed<number>({
    get: () => {
        const firstTransactionType = rows.value[0]?.transaction.type;

        for (const row of rows.value) {
            if (row.transaction.type !== firstTransactionType) {
                return 0;
            }
        }

        return firstTransactionType ?? TransactionType.Expense;
    },
    set: (type: number) => {
        for (const row of rows.value) {
            if (row.transaction.type === type) {
                continue;
            }

            updateType(row.transaction, type);
        }
    }
});

const globalCategoryId = computed<string>({
    get: () => {
        const firstTransactionType = rows.value[0]?.transaction.type;
        let firstCategoryId: string | undefined = undefined;

        if (firstTransactionType === TransactionType.Expense) {
            firstCategoryId = rows.value[0]?.transaction.expenseCategoryId;
        } else if (firstTransactionType === TransactionType.Income) {
            firstCategoryId = rows.value[0]?.transaction.incomeCategoryId;
        }

        for (const row of rows.value) {
            if (row.transaction.type !== firstTransactionType) {
                return '';
            }

            let categoryId: string | undefined;

            if (row.transaction.type === TransactionType.Expense) {
                categoryId = row.transaction.expenseCategoryId;
            } else if (row.transaction.type === TransactionType.Income) {
                categoryId = row.transaction.incomeCategoryId;
            }

            if (categoryId !== firstCategoryId) {
                return '';
            }
        }

        return firstCategoryId ?? '';
    },
    set: (categoryId: string) => {
        for (const row of rows.value) {
            if (row.transaction.type === TransactionType.Expense && allCategoriesMap.value[categoryId]?.type === CategoryType.Expense) {
                row.transaction.expenseCategoryId = categoryId;
            } else if (row.transaction.type === TransactionType.Income && allCategoriesMap.value[categoryId]?.type === CategoryType.Income) {
                row.transaction.incomeCategoryId = categoryId;
            }
        }
    }
});

const globalAccountId = computed<string>({
    get: () => {
        const firstAccountId = rows.value[0]?.transaction.sourceAccountId;

        for (const row of rows.value) {
            if (row.transaction.sourceAccountId !== firstAccountId) {
                return '';
            }
        }

        return firstAccountId ?? '';
    },
    set: (accountId: string) => {
        for (const row of rows.value) {
            row.transaction.sourceAccountId = accountId;
        }
    }
});

const globalTagIds = computed<string[]>({
    get: () => {
        const firstTagIds: string[] = rows.value[0]?.transaction.tagIds ?? [];

        for (const row of rows.value) {
            if ((row.transaction.tagIds ?? []).length !== firstTagIds.length || !isArray1SubsetOfArray2(row.transaction.tagIds ?? [], firstTagIds)) {
                return [];
            }
        }

        return firstTagIds;
    },
    set: (tagIds: string[]) => {
        for (const row of rows.value) {
            row.transaction.tagIds = [...tagIds];
        }
    }
});

function open(options: SetTransactionOptions): Promise<boolean> {
    showState.value = true;
    loading.value = true;
    initOptions.value = options;
    rows.value = [];

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false })
    ]).then(() => {
        addTransactionRow();
        loading.value = false;
    }).catch(error => {
        logger.error('failed to load essential data for batch add transaction dialog', error);

        loading.value = false;
        showState.value = false;

        if (!error.processed) {
            if (rejectFunc) {
                rejectFunc(error);
            }
        }
    });

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function getAccountColor(transaction: Transaction): string | undefined {
    if (transaction.type === TransactionType.Expense) {
        return 'expense';
    } else if (transaction.type === TransactionType.Income) {
        return 'income';
    }

    return undefined;
}

function getAccountCurrency(transaction: Transaction): string {
    const sourceAccount = allAccountsMap.value[transaction.sourceAccountId];

    if (sourceAccount) {
        return sourceAccount.currency;
    }

    return defaultCurrency.value;
}

function createNewTransactionModel(): Transaction {
    const time = getSameDateTimeWithCurrentTimezone(parseDateTimeFromUnixTimeWithBrowserTimezone(getCurrentUnixTime())).getUnixTime();
    const timezone = settingsStore.appSettings.timeZone;
    const transaction = Transaction.createNewTransaction(initOptions.value?.type ?? TransactionType.Expense, time, timezone, getTimezoneOffsetMinutes(time, timezone));

    setTransactionModelByTransaction(
        transaction,
        null,
        allCategories.value,
        allCategoriesMap.value,
        allVisibleAccounts.value,
        allAccountsMap.value,
        allTagsMap.value,
        userStore.currentUserDefaultAccountId,
        initOptions.value ?? {},
        false
    );

    if (transaction.type !== TransactionType.Expense && transaction.type !== TransactionType.Income) {
        transaction.type = TransactionType.Expense;
    }

    return transaction;
}

function addTransactionRow(): void {
    const newTransaction = createNewTransactionModel();
    const newTransactionRow: BatchAddTransactionRow = {
        transaction: newTransaction
    };

    rows.value.push(newTransactionRow);
}

function duplicateTransactionRow(transaction: Transaction, index: number): void {
    if (rows.value.length >= 10) {
        return;
    }

    const newTransaction = Transaction.of(transaction);
    const newTransactionRow: BatchAddTransactionRow = {
        transaction: newTransaction
    };

    rows.value.splice(index + 1, 0, newTransactionRow);
}

function removeTransactionRow(index: number): void {
    if (rows.value.length <= 1) {
        return;
    }

    rows.value.splice(index, 1);
}

function submit(): void {
    if (loading.value || submitting.value || !rows.value.length) {
        return;
    }

    confirmDialog.value?.open('format.misc.confirmBatchAddTransactions', {
        count: formatNumberToLocalizedNumerals(rows.value.length)
    }).then(() => {
        submitting.value = true;

        const transactions: Transaction[] = [];

        for (const row of rows.value) {
            transactions.push(row.transaction);
        }

        transactionsStore.batchCreateTransactions({
            transactions: transactions
        }).then(() => {
            submitting.value = false;
            showState.value = false;

            resolveFunc?.(true);
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function updateTime(transaction: Transaction, time: number): void {
    transaction.time = time;
    transaction.utcOffset = getTimezoneOffsetMinutes(time, transaction.timeZone);
}

function updateType(transaction: Transaction, type: number): void {
    setTransactionModelByTransaction(
        transaction,
        null,
        allCategories.value,
        allCategoriesMap.value,
        allVisibleAccounts.value,
        allAccountsMap.value,
        allTagsMap.value,
        userStore.currentUserDefaultAccountId,
        {
            type: type
        },
        false
    );
}

function cancel(): void {
    if (loading.value || submitting.value) {
        return;
    }

    showState.value = false;
    resolveFunc?.(false);
}

defineExpose({
    open
});
</script>

<style>
.batch-add-transaction-table {
    table {
        .v-input--density-compact {
            --v-input-padding-top: 0;
        }

        th {
            .v-field.v-field--appended {
                --v-field-padding-end: 0;
            }
        }

        th .batch-add-component-header {
            margin-top: 2px;

            .v-field.v-field--appended {
                --v-field-padding-end: 0;
            }

            .v-input__prepend,
            .v-select__selection {
                font-size: 0.8125rem;
            }
        }

        th.batch-add-time-column {
            width: 330px;
            min-width: 200px;
        }

        th.batch-add-type-column {
            width: 140px;
            min-width: 100px;
        }

        th.batch-add-category-column {
            width: 300px;
            min-width: 140px;
        }

        th.batch-add-account-column {
            width: 300px;
            min-width: 120px;
        }

        th.batch-add-amount-column {
            width: 260px;
            min-width: 150px;
        }

        th.batch-add-tags-column {
            width: 260px;
            min-width: 110px;
        }

        th.batch-add-description-column {
            min-width: min(240px, 25dvh);
        }

        td:last-child {
            white-space: nowrap;
        }
    }
}
</style>
