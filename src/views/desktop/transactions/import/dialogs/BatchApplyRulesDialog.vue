<template>
    <v-dialog width="1000" :persistent="loading || isModified" v-model="showState">
        <one-column-dialog-layout content-class="pa-0" :title="tt('Batch Apply Rules')"
                                  :cancel-button-title="tt('Cancel')" @cancel="cancel">
            <template #after-title>
                <v-btn density="compact" color="default" variant="text" class="ms-2"
                       :aria-label="tt('Refresh')" :icon="true" :disabled="loading"
                       :loading="loading" @click="reload">
                    <template #loader>
                        <v-progress-circular indeterminate size="20" />
                    </template>
                    <v-icon :icon="mdiRefresh" size="22" />
                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                </v-btn>
                <v-btn density="compact" color="primary" variant="outlined" class="ms-2"
                       :disabled="loading || !!editingRule" @click="addRule">{{ tt('Add Rule') }}</v-btn>
            </template>

            <template #toolbar>
                <v-btn-group class="ms-2" density="compact" variant="outlined" color="primary"
                             :disabled="loading || !!editingRule || !canApply">
                    <v-btn :disabled="loading || !!editingRule || !canApply"
                           @click="apply('selected')">
                        {{ tt('Apply to Selected Transactions') }}
                    </v-btn>
                    <v-btn :aria-label="tt('More')" :icon="true"
                           :disabled="loading || !!editingRule || !canApply">
                        <v-icon :icon="mdiMenuDown" size="24" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :title="tt('Apply to All Transactions')"
                                             @click="apply('all')"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </v-btn-group>
                <v-btn density="compact" color="default" variant="text" class="ms-2"
                       :aria-label="tt('More')" :disabled="loading || !!editingRule" :icon="true">
                    <v-icon :icon="mdiDotsVertical" />
                    <v-menu activator="parent" max-height="500">
                        <v-list>
                            <v-list-item :prepend-icon="mdiFolderOpenOutline" :title="tt('Load Replace Rule File')"
                                         @click="loadReplaceRuleFile()" />
                            <v-list-item :prepend-icon="mdiContentSaveOutline" :title="tt('Save Replace Rule File')"
                                         :disabled="rules.length < 1" @click="saveReplaceRuleFile()" />
                        </v-list>
                    </v-menu>
                </v-btn>
            </template>

            <template #content>
                <div ref="rulesContainer" class="px-4 pb-2">
                    <div class="text-body-medium text-center py-5" v-if="rules.length < 1">{{ tt('No rules') }}</div>
                    <draggable-list item-key="id" handle=".drag-handle" ghost-class="dragging-item"
                                    :disabled="loading || !!editingRule || rules.length < 2" v-model="rules">
                        <template #item="{ element, index }">
                            <v-card border class="rule-card card-title-with-bg my-4">
                                <v-card-title class="rule-card-title d-flex align-center py-2 px-5">
                                    <v-icon :icon="mdiFunctionVariant" size="20" />
                                    <span class="rule-name text-body-large ms-2" v-if="editingRule !== element">{{ getRuleName(element, index) }}</span>
                                    <div class="rule-name-edit ms-2" v-if="editingRule === element">
                                        <v-text-field autofocus type="text" autocomplete="off"
                                                      density="compact" variant="underlined"
                                                      :disabled="loading" :placeholder="getDefaultRuleName(index)"
                                                      v-text-field-auto-width="{ minWidth: 20, maxWidth: 300, auxSpanId: `rule-name-aux-span-${index + 1}-${element.id}` }"
                                                      v-model="editingRuleName" @keyup.esc="cancelUpdateRuleName"
                                                      @keyup.enter="updateRuleName(element)" />
                                        <span :id="`rule-name-aux-span-${index + 1}-${element.id}`" />
                                    </div>
                                    <v-btn class="ms-2" density="compact" color="primary" variant="text"
                                           :aria-label="tt('Apply')" :disabled="loading" :icon="true"
                                           @click="updateRuleName(element)" v-if="editingRule === element">
                                        <v-icon :icon="mdiCheck" size="18" />
                                        <v-tooltip activator="parent">{{ tt('Apply') }}</v-tooltip>
                                    </v-btn>
                                    <v-btn class="ms-1" density="compact" color="default" variant="text"
                                           :aria-label="tt('Cancel')" :disabled="loading" :icon="true"
                                           @click="cancelUpdateRuleName" v-if="editingRule === element">
                                        <v-icon :icon="mdiClose" size="18" />
                                        <v-tooltip activator="parent">{{ tt('Cancel') }}</v-tooltip>
                                    </v-btn>
                                    <v-btn class="ms-2" density="compact" color="default" variant="text"
                                           :aria-label="tt('Modify Rule Name')" :disabled="loading || !!editingRule" :icon="true"
                                           @click="editingRuleName = element.name; editingRule = element" v-if="editingRule !== element">
                                        <v-icon :icon="mdiPencilOutline" size="18" />
                                        <v-tooltip activator="parent">{{ tt('Modify Rule Name') }}</v-tooltip>
                                    </v-btn>
                                    <v-spacer />
                                    <v-btn density="compact" color="default" variant="text" :aria-label="tt('Delete')"
                                           :disabled="loading || !!editingRule" :icon="true" @click="removeRule(index)">
                                        <v-icon :icon="mdiClose" size="18" />
                                        <v-tooltip activator="parent">{{ tt('Delete') }}</v-tooltip>
                                    </v-btn>
                                    <span class="ms-1">
                                        <v-icon :class="!loading && !editingRule && rules.length > 1 ? 'drag-handle' : 'disabled'"
                                                :aria-label="tt('Drag to Reorder')" :icon="mdiDrag" />
                                        <v-tooltip activator="parent" v-if="!loading && !editingRule && rules.length > 1">{{ tt('Drag to Reorder') }}</v-tooltip>
                                    </span>
                                </v-card-title>
                                <v-divider />
                                <v-card-text>
                                    <div class="text-body-small text-medium-emphasis mb-2">{{ tt('Match Condition') }}</div>
                                    <div class="d-flex overflow-x-auto align-center gap-2 mb-5">
                                        <v-select class="flex-0-0" min-width="220" density="compact"
                                                  item-title="name" item-value="value"
                                                  :disabled="loading || !!editingRule"
                                                  :items="allConditionFields"
                                                  :model-value="element.conditionField"
                                                  @update:model-value="updateConditionField(element, $event)" />

                                        <v-select class="flex-0-0" min-width="140" density="compact"
                                                  item-title="name" item-value="value"
                                                  :disabled="loading || !!editingRule"
                                                  :items="getConditionOperators(element)"
                                                  :model-value="element.conditionOperator"
                                                  @update:model-value="updateConditionOperator(element, $event)" />

                                        <v-autocomplete class="flex-1-1" min-width="240" density="compact"
                                                        item-title="name" item-value="value"
                                                        persistent-placeholder
                                                        :disabled="loading || !!editingRule"
                                                        :no-data-text="tt('No results')"
                                                        :items="getSourceItems(element.conditionField)"
                                                        :model-value="element.conditionValue"
                                                        @update:model-value="updateStringConditionValue(element, $event)"
                                                        v-if="isItemCondition(element.conditionField)" />

                                        <div class="d-flex flex-1-1 align-center gap-2" style="min-width: 280px"
                                             v-else-if="isAmountCondition(element.conditionField)">
                                            <amount-input density="compact" :currency="defaultCurrency" :disabled="loading || !!editingRule"
                                                          :model-value="getAmountConditionValue(element, 0)"
                                                          @update:model-value="updateAmountConditionValue(element, 0, $event)" />
                                            <span class="ms-2 me-2" v-if="isRangeOperator(element.conditionOperator)">{{ tt('format.misc.rangeSeparator') }}</span>
                                            <amount-input density="compact" :currency="defaultCurrency" :disabled="loading || !!editingRule"
                                                          :model-value="getAmountConditionValue(element, 1)"
                                                          @update:model-value="updateAmountConditionValue(element, 1, $event)"
                                                          v-if="isRangeOperator(element.conditionOperator)" />
                                        </div>

                                        <v-text-field class="flex-1-1" min-width="240" density="compact" disabled
                                                      :placeholder="tt('None')"
                                                      v-else-if="isEmptyTextConditionOperator(element.conditionOperator)" />

                                        <v-text-field class="flex-1-1" min-width="240" autocomplete="off" density="compact"
                                                      :disabled="loading || !!editingRule"
                                                      :placeholder="tt('None')" :model-value="element.conditionValue"
                                                      @update:model-value="updateStringConditionValue(element, $event)" v-else />
                                    </div>

                                    <div class="text-body-small text-medium-emphasis mb-2">{{ tt('Action') }}</div>
                                    <div class="d-flex overflow-x-auto align-center gap-2">

                                        <v-select class="flex-0-0" min-width="220" density="compact" item-title="name" item-value="value"
                                                  :disabled="loading || !!editingRule" :items="getActions(element)" :model-value="element.actionType"
                                                  @update:model-value="updateActionType(element, $event)">
                                            <template #selection>
                                                <span>{{ tt(ImportTransactionReplaceRuleAction.valueOf(element.actionType)?.name || '') }}</span>
                                            </template>
                                        </v-select>

                                        <v-select class="flex-1-1" min-width="280" density="compact" item-title="displayName" item-value="type"
                                                  :disabled="loading || !!editingRule" :items="transactionTypes"
                                                  v-model="element.targetValue" v-if="element.actionType === ImportTransactionReplaceRuleActionType.SetTransactionType" />

                                        <two-column-select class="flex-1-1" min-width="280" density="compact"
                                                           primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                           primary-icon-field="icon" primary-icon-type-field="iconType" primary-icon-type="category" primary-color-field="color"
                                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                           secondary-icon-field="icon" secondary-icon-type-field="iconType" secondary-icon-type="category" secondary-color-field="color"
                                                           secondary-hidden-field="hidden" :disabled="loading || !!editingRule"
                                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                           :show-selection-primary-text="true" :custom-selection-primary-text="getSelectedCategoryPrimaryName(element)"
                                                           :custom-selection-secondary-text="getSelectedCategorySecondaryName(element)"
                                                           :items="allCategories[getActionCategoryType(element.actionType)]"
                                                           v-model="element.targetValue"
                                                           v-if="element.actionType === ImportTransactionReplaceRuleActionType.SetExpenseCategory || element.actionType === ImportTransactionReplaceRuleActionType.SetIncomeCategory || element.actionType === ImportTransactionReplaceRuleActionType.SetTransferCategory" />

                                        <two-column-select class="flex-1-1" min-width="280" density="compact"
                                                           primary-key-field="id" primary-value-field="category" primary-title-field="name"
                                                           primary-footer-field="displayBalance" primary-icon-field="icon" primary-icon-type-field="iconType"
                                                           primary-icon-type="account" primary-sub-items-field="accounts" :primary-title-i18n="true"
                                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                           secondary-footer-field="displayBalance" secondary-icon-field="icon" secondary-icon-type-field="iconType"
                                                           secondary-icon-type="account" secondary-color-field="color"
                                                           :disabled="loading || !!editingRule || !allVisibleAccounts.length" :enable-filter="true"
                                                           :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                           :custom-selection-primary-text="getAccountDisplayName(element.targetValue)"
                                                           :items="allVisibleCategorizedAccounts" v-model="element.targetValue"
                                                           v-else-if="element.actionType === ImportTransactionReplaceRuleActionType.SetSourceAccount || element.actionType === ImportTransactionReplaceRuleActionType.SetDestinationAccount" />

                                        <v-autocomplete class="flex-1-1" min-width="280" density="compact" item-title="name" item-value="id"
                                                        persistent-placeholder chips :disabled="loading || !!editingRule"
                                                        :items="allTagsWithGroupHeader" :no-data-text="tt('No available tag')"
                                                        v-model="element.targetValue" v-else-if="element.actionType === ImportTransactionReplaceRuleActionType.AddTag || element.actionType === ImportTransactionReplaceRuleActionType.ReplaceTag">
                                            <template #chip="{ props, internalItem }">
                                                <v-chip :prepend-icon="mdiPound" :text="internalItem.title" v-bind="props" v-if="element.targetValue" />
                                            </template>
                                            <template #subheader="{ props }"><v-list-subheader class="text-body-small">{{ props['title'] }}</v-list-subheader></template>
                                            <template #item="{ props, internalItem }">
                                                <v-list-item :value="internalItem.value" v-bind="props"
                                                             v-if="internalItem.raw instanceof TransactionTag && !internalItem.raw.hidden">
                                                    <template #title>
                                                        <v-list-item-title>
                                                            <div class="d-flex align-center">
                                                                <v-icon size="20" start :icon="mdiPound" />
                                                                <span>{{ internalItem.title }}</span>
                                                            </div>
                                                        </v-list-item-title>
                                                    </template>
                                                </v-list-item>
                                            </template>
                                        </v-autocomplete>

                                        <v-select class="flex-1-1" min-width="280" density="compact" item-title="name" item-value="value"
                                                  :disabled="loading || !!editingRule" :items="amountSigns"
                                                  v-model="element.targetValue"
                                                  v-else-if="element.actionType === ImportTransactionReplaceRuleActionType.SetSourceAmount || element.actionType === ImportTransactionReplaceRuleActionType.SetDestinationAmount" />

                                        <v-autocomplete class="flex-1-1" min-width="280" density="compact"
                                                        item-title="displayNameWithUtcOffset" item-value="name"
                                                        persistent-placeholder auto-select-first
                                                        :disabled="loading || !!editingRule"
                                                        :items="allTimezones"
                                                        v-model="element.targetValue"
                                                        v-else-if="element.actionType === ImportTransactionReplaceRuleActionType.SetTimezone" />
                                    </div>
                                </v-card-text>
                            </v-card>
                        </template>
                    </draggable-list>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagSelectionBase } from '@/components/base/TransactionTagSelectionBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import type { NameValue, NameNumeralValue, TypeAndDisplayName } from '@/core/base.ts';
import type { LocalizedTimezoneInfo } from '@/core/timezone.ts';
import { KnownFileType } from '@/core/file.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    type ImportTransactionReplaceRuleApplyScope,
    ImportTransactionReplaceRuleConditionFieldType,
    ImportTransactionReplaceRuleConditionOperatorType,
    ImportTransactionReplaceRuleConditionOperator,
    ImportTransactionReplaceRuleActionType,
    ImportTransactionReplaceRuleAction
} from '@/core/rule.ts';

import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';
import { ImportTransactionReplaceRule, ImportTransactionReplaceRules } from '@/models/rule.ts';

import { getCurrentUnixTime } from '@/lib/datetime.ts';
import { getTransactionPrimaryCategoryName, getTransactionSecondaryCategoryName } from '@/lib/category.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import logger from '@/lib/logger.ts';
import { openTextFileContent, startDownloadFile } from '@/lib/ui/common.ts';

import {
    mdiRefresh,
    mdiDotsVertical,
    mdiFolderOpenOutline,
    mdiContentSaveOutline,
    mdiPound,
    mdiFunctionVariant,
    mdiPencilOutline,
    mdiCheck,
    mdiClose,
    mdiDrag,
    mdiMenuDown
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

interface BatchApplyRulesDialogResponse {
    rules: ImportTransactionReplaceRule[];
    scope: ImportTransactionReplaceRuleApplyScope;
}

const {
    tt,
    getAllTimezones,
    getAllImportTransactionReplaceRuleConditionFields,
    formatNumberToLocalizedNumeralsWithoutDigitGrouping,
    getCategorizedAccountsWithDisplayBalance
} = useI18n();

const {
    allTagsWithGroupHeader
} = useTransactionTagSelectionBase({ modelValue: [] }, false);

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const rulesContainer = useTemplateRef<HTMLDivElement>('rulesContainer');

let resolveFunc: ((response: BatchApplyRulesDialogResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const rules = ref<ImportTransactionReplaceRule[]>([]);
const editingRule = ref<ImportTransactionReplaceRule | undefined>(undefined);
const editingRuleName = ref<string>('');
const selectedTransactionCount = ref<number>(0);
const sourceExpenseCategoryNames = ref<NameValue[]>([]);
const sourceIncomeCategoryNames = ref<NameValue[]>([]);
const sourceTransferCategoryNames = ref<NameValue[]>([]);
const sourceAccountNames = ref<NameValue[]>([]);
const destinationAccountNames = ref<NameValue[]>([]);
const sourceTagNames = ref<NameValue[]>([]);

const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value, customAccountCategoryOrder.value));
const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
const allTimezones = computed<LocalizedTimezoneInfo[]>(() => getAllTimezones(getCurrentUnixTime(), false));

const allConditionFields = computed<NameValue[]>(() => getAllImportTransactionReplaceRuleConditionFields());
const transactionTypes = computed<TypeAndDisplayName[]>(() => [
    { displayName: tt('Income'), type: TransactionType.Income },
    { displayName: tt('Expense'), type: TransactionType.Expense },
    { displayName: tt('Transfer'), type: TransactionType.Transfer }
]);
const amountSigns = computed<NameNumeralValue[]>(() => [
    { name: tt('Positive'), value: 1 },
    { name: tt('Negative'), value: -1 }
]);

const canApply = computed<boolean>(() => rules.value.length > 0 && rules.value.every(rule => rule.isValid()));
const isModified = computed<boolean>(() => rules.value.length > 1 || (rules.value.length === 1 && !rules.value[0]?.equals(createNewRule(), false)));

function createNewRule(): ImportTransactionReplaceRule {
    return ImportTransactionReplaceRule.of(
        generateRandomUUID(), '',
        ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized,
        ImportTransactionReplaceRuleConditionOperatorType.Contains, '',
        ImportTransactionReplaceRuleActionType.SetExpenseCategory, '');
}

function getDefaultRuleName(index: number): string {
    return tt('format.misc.ruleIndex', { index: formatNumberToLocalizedNumeralsWithoutDigitGrouping(index + 1) });
}

function getRuleName(rule: ImportTransactionReplaceRule, index: number): string {
    return rule.name || getDefaultRuleName(index);
}

function getConditionOperators(rule: ImportTransactionReplaceRule): NameValue[] {
    let operators: ImportTransactionReplaceRuleConditionOperatorType[];

    if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.Tag) {
        operators = [
            ImportTransactionReplaceRuleConditionOperatorType.Contains
        ];
    } else if (isItemCondition(rule.conditionField)) {
        operators = [
            ImportTransactionReplaceRuleConditionOperatorType.Is
        ];
    } else if (isAmountCondition(rule.conditionField)) {
        operators = [
            ImportTransactionReplaceRuleConditionOperatorType.Equals,
            ImportTransactionReplaceRuleConditionOperatorType.NotEquals,
            ImportTransactionReplaceRuleConditionOperatorType.GreaterThan,
            ImportTransactionReplaceRuleConditionOperatorType.LessThan,
            ImportTransactionReplaceRuleConditionOperatorType.Between,
            ImportTransactionReplaceRuleConditionOperatorType.NotBetween
        ];
    } else {
        operators = [
            ImportTransactionReplaceRuleConditionOperatorType.IsEmpty,
            ImportTransactionReplaceRuleConditionOperatorType.IsNotEmpty,
            ImportTransactionReplaceRuleConditionOperatorType.Equals,
            ImportTransactionReplaceRuleConditionOperatorType.NotEquals,
            ImportTransactionReplaceRuleConditionOperatorType.Contains,
            ImportTransactionReplaceRuleConditionOperatorType.NotContains,
            ImportTransactionReplaceRuleConditionOperatorType.StartsWith,
            ImportTransactionReplaceRuleConditionOperatorType.NotStartsWith,
            ImportTransactionReplaceRuleConditionOperatorType.EndsWith,
            ImportTransactionReplaceRuleConditionOperatorType.NotEndsWith,
            ImportTransactionReplaceRuleConditionOperatorType.RegexMatch,
            ImportTransactionReplaceRuleConditionOperatorType.NotRegexMatch
        ];
    }

    return operators.map(operator => ({
        name: tt(ImportTransactionReplaceRuleConditionOperator.valueOf(operator)?.name || ''),
        value: operator
    }));
}

function isActionDisabledForRule(action: ImportTransactionReplaceRuleAction, rule: ImportTransactionReplaceRule): boolean {
    if (action.value === ImportTransactionReplaceRuleActionType.SetExpenseCategory) {
        return rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.IncomeCategory || rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.TransferCategory;
    } else if (action.value === ImportTransactionReplaceRuleActionType.SetIncomeCategory) {
        return rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory || rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.TransferCategory;
    } else if (action.value === ImportTransactionReplaceRuleActionType.SetTransferCategory) {
        return rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory || rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.IncomeCategory;
    } else if (action.value === ImportTransactionReplaceRuleActionType.ReplaceTag || action.value === ImportTransactionReplaceRuleActionType.DeleteTag) {
        return rule.conditionField !== ImportTransactionReplaceRuleConditionFieldType.Tag;
    }

    return false;
}

function getActions(rule: ImportTransactionReplaceRule): { name: string; value: string; props: { disabled: boolean } }[] {
    const result = [];

    for (const action of ImportTransactionReplaceRuleAction.values()) {
        result.push({
            name: tt(action.name),
            value: action.value,
            props: {
                disabled: isActionDisabledForRule(action, rule)
            }
        })
    }

    return result;
}

function getDefaultActionTypeForConditionField(field: ImportTransactionReplaceRuleConditionFieldType): ImportTransactionReplaceRuleActionType {
    if (field === ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory) {
        return ImportTransactionReplaceRuleActionType.SetExpenseCategory;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.IncomeCategory) {
        return ImportTransactionReplaceRuleActionType.SetIncomeCategory;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.TransferCategory) {
        return ImportTransactionReplaceRuleActionType.SetTransferCategory;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.SourceAccount) {
        return ImportTransactionReplaceRuleActionType.SetSourceAccount;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.DestinationAccount) {
        return ImportTransactionReplaceRuleActionType.SetDestinationAccount;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.Tag) {
        return ImportTransactionReplaceRuleActionType.ReplaceTag;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.Amount) {
        return ImportTransactionReplaceRuleActionType.SetSourceAmount;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.TransferInAmount) {
        return ImportTransactionReplaceRuleActionType.SetDestinationAmount;
    }

    return ImportTransactionReplaceRuleActionType.SetTransactionType;
}

function getSourceItems(field: ImportTransactionReplaceRuleConditionFieldType): NameValue[] {
    if (field === ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory) {
        return sourceExpenseCategoryNames.value;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.IncomeCategory) {
        return sourceIncomeCategoryNames.value;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.TransferCategory) {
        return sourceTransferCategoryNames.value;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.SourceAccount) {
        return sourceAccountNames.value;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.DestinationAccount) {
        return destinationAccountNames.value;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.Tag) {
        return sourceTagNames.value;
    } else {
        return [];
    }
}

function isItemCondition(field: ImportTransactionReplaceRuleConditionFieldType): boolean {
    return field === ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory
        || field === ImportTransactionReplaceRuleConditionFieldType.IncomeCategory
        || field === ImportTransactionReplaceRuleConditionFieldType.TransferCategory
        || field === ImportTransactionReplaceRuleConditionFieldType.SourceAccount
        || field === ImportTransactionReplaceRuleConditionFieldType.DestinationAccount
        || field === ImportTransactionReplaceRuleConditionFieldType.Tag;
}

function isAmountCondition(field: ImportTransactionReplaceRuleConditionFieldType): boolean {
    return field === ImportTransactionReplaceRuleConditionFieldType.Amount || field === ImportTransactionReplaceRuleConditionFieldType.TransferInAmount;
}

function isRangeOperator(operator: ImportTransactionReplaceRuleConditionOperatorType): boolean {
    return operator === ImportTransactionReplaceRuleConditionOperatorType.Between || operator === ImportTransactionReplaceRuleConditionOperatorType.NotBetween;
}

function isEmptyTextConditionOperator(operator: ImportTransactionReplaceRuleConditionOperatorType): boolean {
    return operator === ImportTransactionReplaceRuleConditionOperatorType.IsEmpty || operator === ImportTransactionReplaceRuleConditionOperatorType.IsNotEmpty;
}

function getActionCategoryType(action: ImportTransactionReplaceRuleActionType): CategoryType {
    if (action === ImportTransactionReplaceRuleActionType.SetIncomeCategory) {
        return CategoryType.Income;
    } else if (action === ImportTransactionReplaceRuleActionType.SetTransferCategory) {
        return CategoryType.Transfer;
    } else {
        return CategoryType.Expense;
    }
}

function getSelectedCategoryPrimaryName(rule: ImportTransactionReplaceRule): string {
    const targetValue = typeof rule.targetValue === 'string' ? rule.targetValue : '';
    return getTransactionPrimaryCategoryName(targetValue, allCategories.value[getActionCategoryType(rule.actionType)]);
}

function getSelectedCategorySecondaryName(rule: ImportTransactionReplaceRule): string {
    const targetValue = typeof rule.targetValue === 'string' ? rule.targetValue : '';
    return getTransactionSecondaryCategoryName(targetValue, allCategories.value[getActionCategoryType(rule.actionType)]);
}

function getAccountDisplayName(accountId: string | number): string {
    return typeof accountId === 'string' ? Account.findAccountNameById(allAccounts.value, accountId) || '' : '';
}

function open(options: { expenseCategoryNames: NameValue[], incomeCategoryNames: NameValue[], transferCategoryNames: NameValue[], sourceAccountNames: NameValue[], destinationAccountNames: NameValue[], tagNames: NameValue[], selectedTransactionCount: number }): Promise<BatchApplyRulesDialogResponse> {
    rules.value = [];
    editingRule.value = undefined;
    editingRuleName.value = '';
    sourceExpenseCategoryNames.value = options.expenseCategoryNames;
    sourceIncomeCategoryNames.value = options.incomeCategoryNames;
    sourceTransferCategoryNames.value = options.transferCategoryNames;
    sourceAccountNames.value = options.sourceAccountNames;
    destinationAccountNames.value = options.destinationAccountNames;
    sourceTagNames.value = options.tagNames;
    selectedTransactionCount.value = options.selectedTransactionCount;
    showState.value = true;

    addRule();

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function reload(): void {
    loading.value = true;

    Promise.allSettled([
        accountsStore.loadAllAccounts({ force: true }),
        transactionCategoriesStore.loadAllCategories({ force: true }),
        transactionTagsStore.loadAllTags({ force: true })
    ]).then(results => {
        loading.value = false;

        const isAllUpToDate = results.length === 3
            && results[0].status === 'rejected' && results[0].reason?.isUpToDate
            && results[1].status === 'rejected' && results[1].reason?.isUpToDate
            && results[2].status === 'rejected' && results[2].reason?.isUpToDate;

        // show info if all up to date
        if (isAllUpToDate) {
            snackbar.value?.showMessage('Data is up to date');
            return;
        }

        // show error if any
        for (const result of results) {
            if (result.status === 'rejected' && !result.reason?.isUpToDate) {
                snackbar.value?.showError(result.reason);
                return;
            }
        }

        // show info if one of them updated
        for (const result of results) {
            if (result.status === 'fulfilled') {
                snackbar.value?.showMessage('Data has been updated');
                return;
            }
        }
    });
}

function addRule(): void {
    const rule = createNewRule();
    rules.value.push(rule);
    editingRuleName.value = '';
    editingRule.value = undefined;

    nextTick(() => {
        const ruleCards = rulesContainer.value?.querySelectorAll<HTMLElement>('.rule-card');

        if (!ruleCards?.length) {
            return;
        }

        ruleCards[ruleCards.length - 1]?.scrollIntoView({
            behavior: 'smooth',
            block: 'nearest'
        });
    });
}

function removeRule(index: number): void {
    rules.value.splice(index, 1);
}

function updateRuleName(rule: ImportTransactionReplaceRule): void {
    rule.name = editingRuleName.value;
    editingRule.value = undefined;
    editingRuleName.value = '';
}

function cancelUpdateRuleName(): void {
    editingRule.value = undefined;
    editingRuleName.value = '';
}

function updateConditionField(rule: ImportTransactionReplaceRule, field: ImportTransactionReplaceRuleConditionFieldType): void {
    rule.conditionField = field;

    if (field === ImportTransactionReplaceRuleConditionFieldType.Tag) {
        rule.conditionOperator = ImportTransactionReplaceRuleConditionOperatorType.Contains;
        rule.conditionValue = '';
    } else if (isItemCondition(field)) {
        rule.conditionOperator = ImportTransactionReplaceRuleConditionOperatorType.Is;
        rule.conditionValue = '';
    } else if (isAmountCondition(field)) {
        rule.conditionOperator = ImportTransactionReplaceRuleConditionOperatorType.Between;
        rule.conditionValue = [0, 0];
    } else {
        rule.conditionOperator = ImportTransactionReplaceRuleConditionOperatorType.Contains;
        rule.conditionValue = '';
    }

    const currentAction = ImportTransactionReplaceRuleAction.valueOf(rule.actionType);

    if (!currentAction || isActionDisabledForRule(currentAction, rule)) {
        updateActionType(rule, getDefaultActionTypeForConditionField(field));
    }
}

function updateStringConditionValue(rule: ImportTransactionReplaceRule, value: unknown): void {
    rule.conditionValue = typeof value === 'string' ? value : '';
}

function updateConditionOperator(rule: ImportTransactionReplaceRule, operator: ImportTransactionReplaceRuleConditionOperatorType): void {
    rule.conditionOperator = operator;

    if (isEmptyTextConditionOperator(operator)) {
        rule.conditionValue = '';
    }
}

function getAmountConditionValue(rule: ImportTransactionReplaceRule, index: number): number {
    return Array.isArray(rule.conditionValue) ? rule.conditionValue[index] || 0 : 0;
}

function updateAmountConditionValue(rule: ImportTransactionReplaceRule, index: number, value: number): void {
    if (!Array.isArray(rule.conditionValue)) {
        rule.conditionValue = [0, 0];
    }

    rule.conditionValue[index] = value;

    if (!isRangeOperator(rule.conditionOperator)) {
        rule.conditionValue[1] = rule.conditionValue[0];
    }
}

function updateActionType(rule: ImportTransactionReplaceRule, actionType: ImportTransactionReplaceRuleActionType): void {
    rule.actionType = actionType;

    if (actionType === ImportTransactionReplaceRuleActionType.SetTransactionType) {
        rule.targetValue = TransactionType.Expense;
    } else if (actionType === ImportTransactionReplaceRuleActionType.SetSourceAmount || actionType === ImportTransactionReplaceRuleActionType.SetDestinationAmount) {
        rule.targetValue = 1;
    } else {
        rule.targetValue = '';
    }
}

function loadReplaceRuleFile(): void {
    openTextFileContent({
        allowedExtensions: KnownFileType.JSON.contentType
    }).then(content => {
        const result = ImportTransactionReplaceRules.parseFromJson(content, generateRandomUUID);

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

function apply(scope: ImportTransactionReplaceRuleApplyScope): void {
    resolveFunc?.({
        rules: rules.value,
        scope: scope
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

<style>
.rule-name {
    white-space: pre;
}

.rule-card-title {
    min-height: 52px;
}

.rule-name-edit {
    display: flex;
    align-items: center;
    height: 36px;
}

.rule-name-edit .v-field {
    font-size: 0.9375rem;
}

.rule-name-edit .v-field__input {
    padding-top: 0;
}
</style>
