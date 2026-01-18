<template>
    <v-dialog width="1000" :persistent="loading || !!rules.length || !!newRule.targetId" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex flex-wrap align-center justify-center">
                    <div class="d-flex flex-wrap align-center">
                        <h4 class="text-h4 text-wrap">{{ tt('Batch Replace Categories / Accounts / Tags') }}</h4>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" :disabled="loading"
                               :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                    <v-spacer/>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :icon="true" :disabled="loading">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <v-list-item :prepend-icon="mdiFolderOpenOutline"
                                             :title="tt('Load Replace Rule File')"
                                             @click="loadReplaceRuleFile()"></v-list-item>
                                <v-list-item :prepend-icon="mdiContentSaveOutline"
                                             :title="tt('Save Replace Rule File')"
                                             @click="saveReplaceRuleFile()"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
                <v-row>
                    <v-col cols="12">
                        <v-table density="comfortable" fixed-header fixed-footer height="400" striped="even">
                            <thead>
                                <tr>
                                    <th class="text-left">{{ tt('Type') }}</th>
                                    <th class="text-left">{{ tt('Source Value') }}</th>
                                    <th class="text-left">{{ tt('Target Value') }}</th>
                                    <th class="text-right">{{ tt('Operation') }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="(rule, index) in rules" :key="index">
                                    <td class="text-left">{{ getRuleTypeDisplayName(rule) }}</td>
                                    <td class="text-left">{{ rule.sourceValue || tt('(Empty)') }}</td>
                                    <td class="text-left">{{ getRuleTargetValueDisplayName(rule) }}</td>
                                    <td class="text-right">
                                        <v-btn density="comfortable" variant="tonal" color="error"
                                               :disabled="loading" @click="removeRule(index)">{{ tt('Delete') }}</v-btn>
                                    </td>
                                </tr>
                            </tbody>
                            <tfoot>
                                <tr style="background-color: rgb(var(--v-theme-surface))">
                                    <td>
                                        <v-select class="w-100" density="compact" variant="underlined"
                                                  item-title="name"
                                                  item-value="value"
                                                  :disabled="loading"
                                                  :items="[
                                                      {
                                                          value: 'expenseCategory',
                                                          name: tt('Expense Category')
                                                      },
                                                      {
                                                          value: 'incomeCategory',
                                                          name: tt('Income Category')
                                                      },
                                                      {
                                                          value: 'transferCategory',
                                                          name: tt('Transfer Category')
                                                      },
                                                      {
                                                          value: 'account',
                                                          name: tt('Account')
                                                      },
                                                      {
                                                          value: 'tag',
                                                          name: tt('Transaction Tag')
                                                      }
                                                  ]"
                                                  v-model="newRule.dataType"
                                                  @update:model-value="newRule.sourceValue = ''; newRule.targetId = ''"
                                        />
                                    </td>
                                    <td>
                                        <v-autocomplete class="w-100" density="compact" variant="underlined"
                                                        item-title="name" item-value="value" persistent-placeholder
                                                        :disabled="loading" :items="sourceItems"
                                                        :no-data-text="noSourceItemText"
                                                        v-model="newRule.sourceValue">
                                        </v-autocomplete>
                                    </td>
                                    <td>
                                        <two-column-select density="compact" variant="underlined"
                                                           primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                           secondary-hidden-field="hidden"
                                                           :disabled="loading || !hasVisibleExpenseCategories"
                                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                           :show-selection-primary-text="true"
                                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(newRule.targetId, allCategories[CategoryType.Expense])"
                                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(newRule.targetId, allCategories[CategoryType.Expense])"
                                                           :items="allCategories[CategoryType.Expense]"
                                                           v-model="newRule.targetId"
                                                           v-if="newRule.dataType === 'expenseCategory'">
                                        </two-column-select>
                                        <two-column-select density="compact" variant="underlined"
                                                           primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                           secondary-hidden-field="hidden"
                                                           :disabled="loading || !hasVisibleIncomeCategories"
                                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                           :show-selection-primary-text="true"
                                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(newRule.targetId, allCategories[CategoryType.Income])"
                                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(newRule.targetId, allCategories[CategoryType.Income])"
                                                           :items="allCategories[CategoryType.Income]"
                                                           v-model="newRule.targetId"
                                                           v-if="newRule.dataType === 'incomeCategory'">
                                        </two-column-select>
                                        <two-column-select density="compact" variant="underlined"
                                                           primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                           secondary-hidden-field="hidden"
                                                           :disabled="loading || !hasVisibleTransferCategories"
                                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                           :show-selection-primary-text="true"
                                                           :custom-selection-primary-text="getTransactionPrimaryCategoryName(newRule.targetId, allCategories[CategoryType.Transfer])"
                                                           :custom-selection-secondary-text="getTransactionSecondaryCategoryName(newRule.targetId, allCategories[CategoryType.Transfer])"
                                                           :items="allCategories[CategoryType.Transfer]"
                                                           v-model="newRule.targetId"
                                                           v-if="newRule.dataType === 'transferCategory'">
                                        </two-column-select>
                                        <two-column-select density="compact" variant="underlined"
                                                           primary-key-field="id" primary-value-field="category"
                                                           primary-title-field="name" primary-footer-field="displayBalance"
                                                           primary-icon-field="icon" primary-icon-type="account"
                                                           primary-sub-items-field="accounts"
                                                           :primary-title-i18n="true"
                                                           secondary-key-field="id" secondary-value-field="id"
                                                           secondary-title-field="name" secondary-footer-field="displayBalance"
                                                           secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                           :disabled="loading || !allVisibleAccounts.length"
                                                           :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                           :custom-selection-primary-text="getAccountDisplayName(newRule.targetId)"
                                                           :items="allVisibleCategorizedAccounts"
                                                           v-model="newRule.targetId"
                                                           v-if="newRule.dataType === 'account'">
                                        </two-column-select>
                                        <v-autocomplete density="compact" variant="underlined"
                                                        item-title="name" item-value="id"
                                                        persistent-placeholder chips
                                                        :disabled="loading" :items="allTagsWithGroupHeader"
                                                        :no-data-text="tt('No available tag')"
                                                        v-model="newRule.targetId"
                                                        v-if="newRule.dataType == 'tag'">
                                            <template #chip="{ props, item }">
                                                <v-chip :prepend-icon="mdiPound" :text="item.title" v-bind="props" v-if="newRule.targetId"/>
                                            </template>

                                            <template #subheader="{ props }">
                                                <v-list-subheader>{{ props['title'] }}</v-list-subheader>
                                            </template>

                                            <template #item="{ props, item }">
                                                <v-list-item :value="item.value" v-bind="props" v-if="item.raw instanceof TransactionTag && !item.raw.hidden">
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
                                    </td>
                                    <td class="text-right">
                                        <v-btn density="comfortable" variant="tonal" color="primary"
                                               :disabled="loading || !newRule.dataType || !newRule.targetId"
                                               @click="addNewRule()">{{ tt('Add Rule') }}</v-btn>
                                    </td>
                                </tr>
                            </tfoot>
                        </v-table>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :disabled="loading" @click="confirm">{{ tt('OK') }}</v-btn>
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
import { useTransactionTagSelectionBase } from '@/components/base/TransactionTagSelectionBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import type { NameValue } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';
import { ImportTransactionReplaceRule, ImportTransactionReplaceRules } from '@/core/import_transaction.ts';
import { KnownFileType } from '@/core/file.ts';

import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';

import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName
} from '@/lib/category.ts';

import logger from '@/lib/logger.ts';

import {
    openTextFileContent,
    startDownloadFile
} from '@/lib/ui/common.ts';

import {
    mdiRefresh,
    mdiDotsVertical,
    mdiFolderOpenOutline,
    mdiContentSaveOutline, mdiPound
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

interface BatchReplaceAllTypesDialogResponse {
    rules: ImportTransactionReplaceRule[]
}

const { tt, getCategorizedAccountsWithDisplayBalance } = useI18n();

const { allTagsWithGroupHeader } = useTransactionTagSelectionBase({ modelValue: [] }, false);

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const rules = ref<ImportTransactionReplaceRule[]>([]);
const newRule = ref<ImportTransactionReplaceRule>(ImportTransactionReplaceRule.of('expenseCategory', '', ''));

const sourceExpenseCategoryNames = ref<NameValue[]>([]);
const sourceIncomeCategoryNames = ref<NameValue[]>([]);
const sourceTransferCategoryNames = ref<NameValue[]>([]);
const sourceAccountNames = ref<NameValue[]>([]);
const sourceTagNames = ref<NameValue[]>([]);

let resolveFunc: ((response: BatchReplaceAllTypesDialogResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);
const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value, customAccountCategoryOrder.value));
const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
const allTagsMap = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);

const hasVisibleExpenseCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleExpenseCategories);
const hasVisibleIncomeCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleIncomeCategories);
const hasVisibleTransferCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleTransferCategories);

const sourceItems = computed<NameValue[]>(() => {
    switch (newRule.value.dataType) {
        case 'expenseCategory':
            return sourceExpenseCategoryNames.value;
        case 'incomeCategory':
            return sourceIncomeCategoryNames.value;
        case 'transferCategory':
            return sourceTransferCategoryNames.value;
        case 'account':
            return sourceAccountNames.value;
        case 'tag':
            return sourceTagNames.value;
        default:
            return [];
    }
});

const noSourceItemText = computed<string>(() => {
    switch (newRule.value.dataType) {
        case 'expenseCategory':
            return tt('No available category');
        case 'incomeCategory':
            return tt('No available category');
        case 'transferCategory':
            return tt('No available category');
        case 'account':
            return tt('No available account');
        case 'tag':
            return tt('No available tag');
        default:
            return '';
    }
});

function getRuleTypeDisplayName(rule: ImportTransactionReplaceRule): string {
    switch (rule.dataType) {
        case 'expenseCategory':
            return tt('Expense Category');
        case 'incomeCategory':
            return tt('Income Category');
        case 'transferCategory':
            return tt('Transfer Category');
        case 'account':
            return tt('Account');
        case 'tag':
            return tt('Transaction Tag');
        default:
            return '';
    }
}

function getRuleTargetValueDisplayName(rule: ImportTransactionReplaceRule): string {
    switch (rule.dataType) {
        case 'expenseCategory':
            return getTransactionSecondaryCategoryName(rule.targetId, allCategories.value[CategoryType.Expense]) || '';
        case 'incomeCategory':
            return getTransactionSecondaryCategoryName(rule.targetId, allCategories.value[CategoryType.Income]) || '';
        case 'transferCategory':
            return getTransactionSecondaryCategoryName(rule.targetId, allCategories.value[CategoryType.Transfer]) || '';
        case 'account':
            return getAccountDisplayName(rule.targetId);
        case 'tag':
            return allTagsMap.value[rule.targetId]?.name ?? '';
        default:
            return '';
    }
}

function getAccountDisplayName(accountId?: string): string {
    if (accountId) {
        return Account.findAccountNameById(allAccounts.value, accountId) || '';
    } else {
        return tt('None');
    }
}

function open(options: { expenseCategoryNames: NameValue[], incomeCategoryNames: NameValue[], transferCategoryNames: NameValue[], accountNames: NameValue[], tagNames: NameValue[] }): Promise<BatchReplaceAllTypesDialogResponse> {
    rules.value = [];
    newRule.value = ImportTransactionReplaceRule.of('expenseCategory', '', '');
    sourceExpenseCategoryNames.value = options.expenseCategoryNames;
    sourceIncomeCategoryNames.value = options.incomeCategoryNames;
    sourceTransferCategoryNames.value = options.transferCategoryNames;
    sourceAccountNames.value = options.accountNames;
    sourceTagNames.value = options.tagNames;
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function reload(): void {
    loading.value = true;

    Promise.all([
        accountsStore.loadAllAccounts({ force: true }),
        transactionCategoriesStore.loadAllCategories({ force: true }),
        transactionTagsStore.loadAllTags({ force: true })
    ]).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function loadReplaceRuleFile(): void {
    openTextFileContent({
        allowedExtensions: KnownFileType.JSON.contentType
    }).then(content => {
        const result = ImportTransactionReplaceRules.parseFromJson(content);

        if (result) {
            rules.value = result.getRules();
        } else {
            logger.error('Failed to parse replace rule file');
            snackbar.value?.showError('Replace rule file is invalid');
        }
    }).catch(error => {
        logger.error('Failed to open replace rule file', error);
        snackbar.value?.showError('Replace rule file is invalid');
    });
}

function saveReplaceRuleFile(): void {
    const fileName = KnownFileType.JSON.formatFileName(tt('dataExport.defaultImportReplaceRuleFileName'));
    startDownloadFile(fileName, KnownFileType.JSON.createBlob(ImportTransactionReplaceRules.of(rules.value).toJson()));
}

function removeRule(index: number): void {
    rules.value.splice(index, 1);
}

function addNewRule(): void {
    if (!newRule.value.dataType || !newRule.value.targetId) {
        return;
    }

    rules.value.push(newRule.value);
    newRule.value = ImportTransactionReplaceRule.of('expenseCategory', '', '');
}

function confirm(): void {
    resolveFunc?.({
        rules: rules.value
    });
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
