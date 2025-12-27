<template>
    <v-card-subtitle class="px-5">
        <div class="title-and-toolbar d-flex">
            <v-btn color="default" variant="outlined"
                   :disabled="loading || !!editingQuery"
                   @click="addQuery">{{ tt('Add Query') }}</v-btn>
            <v-spacer />
            <v-btn color="secondary" variant="tonal"
                   :disabled="loading || !!editingQuery || queries.length < 1"
                   @click="clearAllQueries">{{ tt('Clear All') }}</v-btn>
        </div>
    </v-card-subtitle>
    <v-card-text class="pt-0">
        <div :key="queryIndex" v-for="(query, queryIndex) in queries">
            <v-card border class="card-title-with-bg mt-4">
                <v-card-title class="d-flex align-center py-2 px-5">
                    <v-icon :icon="mdiTextBoxSearchOutline" size="20" />
                    <span class="query-name text-subtitle-1 ms-2" v-if="editingQuery !== query">{{ query.name || `${tt('Query')} #${queryIndex + 1}` }}</span>
                    <div class="query-name-edit ms-2" v-if="editingQuery === query">
                        <v-text-field autofocus type="text" density="compact" variant="underlined"
                                      :disabled="loading"
                                      :placeholder="`${tt('Query')} #${queryIndex + 1}`"
                                      v-text-field-auto-width="{ minWidth: 20, maxWidth: 300, auxSpanId: `query-name-aux-span-${queryIndex + 1}` }"
                                      v-model="editingQueryName"
                                      @keyup.esc="cancelUpdateQueryName"
                                      @keyup.enter="updateQueryName(query)" />
                        <span :id="`query-name-aux-span-${queryIndex + 1}`" />
                    </div>
                    <v-btn class="ms-2" density="compact" color="primary" variant="text" size="small"
                           :icon="true" :disabled="loading"
                           @click="updateQueryName(query)"
                           v-if="editingQuery === query">
                        <v-icon :icon="mdiCheck" size="18" />
                        <v-tooltip activator="parent">{{ tt('Update') }}</v-tooltip>
                    </v-btn>
                    <v-btn class="ms-2" density="compact" color="default" variant="text" size="small"
                           :icon="true" :disabled="loading"
                           @click="cancelUpdateQueryName"
                           v-if="editingQuery === query">
                        <v-icon :icon="mdiClose" size="18" />
                        <v-tooltip activator="parent">{{ tt('Cancel') }}</v-tooltip>
                    </v-btn>
                    <v-btn class="ms-2" density="compact" color="default" variant="text" size="small"
                           :icon="true" :disabled="loading || !!editingQuery"
                           @click="editingQueryName = query.name; editingQuery = query"
                           v-if="!editingQuery || editingQuery !== query">
                        <v-icon :icon="mdiPencilOutline" size="18" />
                        <v-tooltip activator="parent">{{ tt('Modify Query Name') }}</v-tooltip>
                    </v-btn>
                    <v-btn class="ms-2" density="compact" color="default" variant="text" size="small"
                           :icon="true" :disabled="loading || !!editingQuery"
                           @click="duplicateQuery(query)"
                           v-if="!editingQuery || editingQuery !== query">
                        <v-icon :icon="mdiContentCopy" size="18" />
                        <v-tooltip activator="parent">{{ tt('Duplicate') }}</v-tooltip>
                    </v-btn>
                    <v-spacer />
                    <v-switch class="bidirectional-switch ms-2" color="secondary"
                              :disabled="loading || !!editingQuery || !query.conditions || query.conditions.length < 1"
                              :label="tt('Expression')"
                              v-model="showExpression[queryIndex]"
                              @click="showExpression[queryIndex] = !showExpression[queryIndex]">
                        <template #prepend>
                            <span>{{ tt('Editor') }}</span>
                        </template>
                    </v-switch>
                    <v-btn class="ms-2" density="compact" color="default" variant="text" size="small"
                           :icon="true" :disabled="loading || !!editingQuery || queries.length < 1"
                           @click="removeQuery(queryIndex)">
                        <v-icon :icon="mdiClose" size="18" />
                        <v-tooltip activator="parent">{{ tt('Remove Query') }}</v-tooltip>
                    </v-btn>
                </v-card-title>

                <v-divider />

                <v-card-text>
                    <v-row>
                        <v-col cols="12">
                            <div class="text-center py-4" v-if="!query.conditions || query.conditions.length < 1">
                                {{ tt('No conditions defined. All transactions will match.') }}
                            </div>

                            <div v-else-if="query.conditions && query.conditions.length > 0 && !showExpression[queryIndex]">
                                <div :key="conditionIndex" v-for="(conditionWithRelation, conditionIndex) in query.conditions">
                                    <div class="d-flex overflow-x-auto align-center gap-2 mb-4">
                                        <v-select
                                            disabled
                                            class="flex-0-0"
                                            width="120px"
                                            density="compact"
                                            item-title="displayName"
                                            item-value="value"
                                            :items="[{ value: TransactionExploreConditionRelation.First, displayName: tt('WHERE') }]"
                                            :model-value="TransactionExploreConditionRelation.First"
                                            v-if="conditionIndex < 1"
                                        />

                                        <v-select
                                            class="flex-0-0"
                                            width="120px"
                                            density="compact"
                                            item-title="displayName"
                                            item-value="value"
                                            :disabled="loading || !!editingQuery"
                                            :items="[
                                                { value: TransactionExploreConditionRelation.And, displayName: tt('AND') },
                                                { value: TransactionExploreConditionRelation.Or, displayName: tt('OR') }
                                            ]"
                                            v-model="conditionWithRelation.relation"
                                            v-else-if="conditionIndex >= 1"
                                        />

                                        <v-select
                                            class="flex-0-0"
                                            density="compact"
                                            item-title="name"
                                            item-value="value"
                                            :disabled="loading || !!editingQuery"
                                            :items="allTransactionExploreConditionFields"
                                            @update:model-value="updateConditionField(queryIndex, conditionIndex, TransactionExploreConditionField.valueOf($event))"
                                            v-model="conditionWithRelation.condition.field"
                                        />

                                        <v-select
                                            class="flex-0-0"
                                            density="compact"
                                            item-title="name"
                                            item-value="value"
                                            :disabled="loading || !!editingQuery"
                                            :items="getAllTransactionExploreConditionOperators(conditionWithRelation.getSupportedOperators())"
                                            v-model="conditionWithRelation.condition.operator"
                                        />

                                        <div class="d-flex w-100 flex-1-1" style="min-width: 280px;">
                                            <v-select
                                                multiple chips closable-chips
                                                density="compact"
                                                item-title="displayName"
                                                item-value="type"
                                                :disabled="loading || !!editingQuery"
                                                :placeholder="tt('None')"
                                                :items="[
                                                    { type: TransactionType.Expense, displayName: tt('Expense') },
                                                    { type: TransactionType.Income, displayName: tt('Income') },
                                                    { type: TransactionType.Transfer, displayName: tt('Transfer') }
                                                ]"
                                                v-model="conditionWithRelation.condition.value"
                                                v-if="conditionWithRelation.condition.field === TransactionExploreConditionField.TransactionType.value"
                                            >
                                                <template #item="{ props, item }">
                                                    <v-list-item :value="item.value" v-bind="props">
                                                        <template #title>
                                                            <v-list-item-title>
                                                                <div class="d-flex align-center">{{ item.title }}</div>
                                                            </v-list-item-title>
                                                        </template>
                                                    </v-list-item>
                                                </template>
                                            </v-select>

                                            <v-text-field
                                                class="always-cursor-pointer text-field-truncate"
                                                density="compact"
                                                item-title="displayName"
                                                item-value="type"
                                                persistent-placeholder
                                                :readonly="true"
                                                :disabled="loading || !!editingQuery || !hasAnyTransactionCategory"
                                                :placeholder="tt('None')"
                                                :model-value="getFilteredTransactionCategoriesDisplayContent(arrayItemToObjectField(conditionWithRelation.condition.value as string[], true))"
                                                @click="currentCondition = conditionWithRelation.condition; showFilterTransactionCategoriesDialog = true"
                                                v-else-if="conditionWithRelation.condition.field === TransactionExploreConditionField.TransactionCategory.value"
                                            />

                                            <v-text-field
                                                class="always-cursor-pointer text-field-truncate"
                                                density="compact"
                                                item-title="displayName"
                                                item-value="type"
                                                persistent-placeholder
                                                :readonly="true"
                                                :disabled="loading || !!editingQuery || !hasAnyAccount"
                                                :placeholder="tt('None')"
                                                :model-value="getFilteredAccountsDisplayContent(arrayItemToObjectField(conditionWithRelation.condition.value as string[], true))"
                                                @click="currentCondition = conditionWithRelation.condition; showFilterSourceAccountsDialog = true"
                                                v-else-if="conditionWithRelation.condition.field === TransactionExploreConditionField.SourceAccount.value"
                                            />

                                            <v-text-field
                                                class="always-cursor-pointer text-field-truncate"
                                                density="compact"
                                                item-title="displayName"
                                                item-value="type"
                                                persistent-placeholder
                                                :readonly="true"
                                                :disabled="loading || !!editingQuery || !hasAnyAccount"
                                                :placeholder="tt('None')"
                                                :model-value="getFilteredAccountsDisplayContent(arrayItemToObjectField(conditionWithRelation.condition.value as string[], true))"
                                                @click="currentCondition = conditionWithRelation.condition; showFilterDestinationAccountsDialog = true"
                                                v-else-if="conditionWithRelation.condition.field === TransactionExploreConditionField.DestinationAccount.value"
                                            />

                                            <div class="d-flex w-100 align-center gap-2"
                                                 v-else-if="conditionWithRelation.condition.field === TransactionExploreConditionField.SourceAmount.value ||
                                                            conditionWithRelation.condition.field === TransactionExploreConditionField.DestinationAmount.value">
                                                <amount-input density="compact"
                                                              :currency="defaultCurrency"
                                                              :disabled="loading || !!editingQuery"
                                                              v-model="conditionWithRelation.condition.value[0]"
                                                />
                                                <span class="ms-2 me-2"
                                                      v-if="conditionWithRelation.condition.operator === TransactionExploreConditionOperator.Between.value ||
                                                            conditionWithRelation.condition.operator === TransactionExploreConditionOperator.NotBetween.value">~</span>
                                                <amount-input density="compact"
                                                              :currency="defaultCurrency"
                                                              :disabled="loading || !!editingQuery"
                                                              v-model="conditionWithRelation.condition.value[1]"
                                                              v-if="conditionWithRelation.condition.operator === TransactionExploreConditionOperator.Between.value ||
                                                                    conditionWithRelation.condition.operator === TransactionExploreConditionOperator.NotBetween.value"
                                                />
                                            </div>

                                            <div class="d-flex w-100" v-else-if="conditionWithRelation.condition.field === TransactionExploreConditionField.TransactionTag.value">
                                                <v-text-field
                                                    disabled
                                                    persistent-placeholder
                                                    density="compact"
                                                    :placeholder="tt('None')"
                                                    v-if="conditionWithRelation.condition.field === TransactionExploreConditionField.TransactionTag.value &&
                                                         (conditionWithRelation.condition.operator === TransactionExploreConditionOperator.IsEmpty.value || conditionWithRelation.condition.operator === TransactionExploreConditionOperator.IsNotEmpty.value)"
                                                />

                                                <v-autocomplete
                                                    density="compact"
                                                    item-title="name"
                                                    item-value="id"
                                                    auto-select-first
                                                    persistent-placeholder
                                                    multiple
                                                    chips
                                                    closable-chips
                                                    :disabled="loading || !!editingQuery"
                                                    :placeholder="tt('None')"
                                                    :items="allTags"
                                                    v-model="conditionWithRelation.condition.value"
                                                    v-model:search="tagSearchContent"
                                                    v-else-if="conditionWithRelation.condition.operator !== TransactionExploreConditionOperator.IsEmpty.value && conditionWithRelation.condition.operator !== TransactionExploreConditionOperator.IsNotEmpty.value"
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
                                                        <v-list-item :disabled="true" v-bind="props"
                                                                     v-if="item.raw.hidden && item.raw.name.toLowerCase().indexOf(tagSearchContent.toLowerCase()) >= 0 && isAllFilteredTagHidden">
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

                                                    <template #no-data>
                                                        <v-list class="py-0">
                                                            <v-list-item>{{ tt('No available tag') }}</v-list-item>
                                                        </v-list>
                                                    </template>
                                                </v-autocomplete>
                                            </div>

                                            <v-text-field disabled density="compact"
                                                          :placeholder="tt('None')"
                                                          v-else-if="conditionWithRelation.condition.field === TransactionExploreConditionField.Description.value &&
                                                                     conditionWithRelation.condition.operator === TransactionExploreConditionOperator.IsEmpty.value || conditionWithRelation.condition.operator === TransactionExploreConditionOperator.IsNotEmpty.value"
                                            />

                                            <v-text-field density="compact"
                                                          :disabled="loading || !!editingQuery"
                                                          :placeholder="tt('None')"
                                                          v-model="conditionWithRelation.condition.value"
                                                          v-else-if="conditionWithRelation.condition.field === TransactionExploreConditionField.Description.value &&
                                                                     conditionWithRelation.condition.operator !== TransactionExploreConditionOperator.IsEmpty.value && conditionWithRelation.condition.operator !== TransactionExploreConditionOperator.IsNotEmpty.value"
                                            />
                                        </div>

                                        <v-btn color="default" density="compact"
                                               variant="text" size="small"
                                               :icon="true"
                                               :disabled="loading || !!editingQuery"
                                               @click="removeCondition(queryIndex, conditionIndex)">
                                            <v-icon :icon="mdiClose" size="18" />
                                            <v-tooltip activator="parent">{{ tt('Remove Condition') }}</v-tooltip>
                                        </v-btn>
                                    </div>
                                </div>
                            </div>
                            <div v-else-if="query.conditions && query.conditions.length > 0 && showExpression[queryIndex]">
                                <div class="w-100 code-container">
                                    <v-textarea class="w-100 always-cursor-text mb-4" :readonly="true"
                                                :value="getExpression(queryIndex)"></v-textarea>
                                </div>
                            </div>

                            <v-btn class="px-2" density="comfortable" color="primary" variant="text" size="small"
                                   :prepend-icon="mdiPlus"
                                   :disabled="loading || !!editingQuery || showExpression[queryIndex]"
                                   @click="addCondition(queryIndex)">
                                {{ tt('Add Condition') }}
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>

            <div class="query-group-separator d-flex align-center justify-center my-4"
                 v-if="queries.length > 1 && queryIndex < queries.length - 1">
                <v-chip color="primary" variant="outlined" size="small">
                    {{ tt('or') }}
                </v-chip>
            </div>
        </div>
    </v-card-text>

    <v-dialog width="800" v-model="showFilterSourceAccountsDialog">
        <account-filter-settings-card type="custom" :dialog-mode="true"
                                      :selected-account-ids="isArray(currentCondition?.value) ? currentCondition?.value as string[] : []"
                                      @settings:change="updateSourceAccount" />
    </v-dialog>

    <v-dialog width="800" v-model="showFilterDestinationAccountsDialog">
        <account-filter-settings-card type="custom" :dialog-mode="true"
                                      :selected-account-ids="isArray(currentCondition?.value) ? currentCondition?.value as string[] : []"
                                      @settings:change="updateDestinationAccount" />
    </v-dialog>

    <v-dialog width="800" v-model="showFilterTransactionCategoriesDialog">
        <category-filter-settings-card type="custom" :dialog-mode="true"
                                       :selected-category-ids="isArray(currentCondition?.value) ? currentCondition?.value as string[] : []"
                                       @settings:change="updateTransactionCategories" />
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import AccountFilterSettingsCard from '@/views/desktop/common/cards/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from '@/views/desktop/common/cards/CategoryFilterSettingsCard.vue';

import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useExploresStore } from '@/stores/explore.ts';

import { type NameValue, entries, values } from '@/core/base.ts';
import { AccountType } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    TransactionExploreConditionRelation,
    TransactionExploreConditionField,
    TransactionExploreConditionOperator
} from '@/core/explore.ts';

import {
    type TransactionTag
} from '@/models/transaction_tag.ts';

import {
    type TransactionExploreCondition,
    TransactionExploreQuery
} from '@/models/explore.ts';

import {
    isArray,
    isObjectEmpty,
    arrayItemToObjectField
} from '@/lib/common.ts';

import logger from '@/lib/logger.ts';

import {
    mdiTextBoxSearchOutline,
    mdiPlus,
    mdiPencilOutline,
    mdiContentCopy,
    mdiCheck,
    mdiClose,
    mdiPound
} from '@mdi/js';

interface ExploreQueryTabProps {
    loading?: boolean;
}

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<ExploreQueryTabProps>();

const {
    tt,
    joinMultiText,
    getAllTransactionExploreConditionFields,
    getAllTransactionExploreConditionOperators
} = useI18n();

const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const exploresStore = useExploresStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const currentCondition = ref<TransactionExploreCondition | undefined>(undefined);
const showExpression = ref<Record<number, boolean>>({});
const showFilterSourceAccountsDialog = ref<boolean>(false);
const showFilterDestinationAccountsDialog = ref<boolean>(false);
const showFilterTransactionCategoriesDialog = ref<boolean>(false);
const tagSearchContent = ref<string>('');
const editingQuery = ref<TransactionExploreQuery | undefined>(undefined);
const editingQueryName = ref<string>('');

const queries = computed<TransactionExploreQuery[]>(() => exploresStore.transactionExploreFilter.query);

const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
const hasAnyAccount = computed<boolean>(() => accountsStore.allPlainAccounts.length > 0);
const hasAnyTransactionCategory = computed<boolean>(() => !isObjectEmpty(transactionCategoriesStore.allTransactionCategoriesMap));
const allTags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);

const allTransactionExploreConditionFields = computed<NameValue[]>(() => getAllTransactionExploreConditionFields());

const isAllFilteredTagHidden = computed<boolean>(() => {
    const lowerCaseTagSearchContent = tagSearchContent.value.toLowerCase();
    let hiddenCount = 0;

    for (const tag of allTags.value) {
        if (!lowerCaseTagSearchContent || tag.name.toLowerCase().indexOf(lowerCaseTagSearchContent) >= 0) {
            if (!tag.hidden) {
                return false;
            }

            hiddenCount++;
        }
    }

    return hiddenCount > 0;
});

function getFilteredAccountsDisplayContent(filterAccountIds?: Record<string, boolean>): string {
    if ((props.loading && !hasAnyAccount.value) || !accountsStore.allVisiblePlainAccounts || !accountsStore.allVisiblePlainAccounts.length) {
        return '';
    }

    if (!filterAccountIds) {
        return tt('All');
    }

    let allAccountSelected = true;
    const selectedAccountNames: string[] = [];

    for (const account of accountsStore.allPlainAccounts) {
        if (account.type === AccountType.MultiSubAccounts.type) {
            continue;
        }

        if (!filterAccountIds[account.id]) {
            allAccountSelected = false;
        } else {
            selectedAccountNames.push(account.name);
        }
    }

    if (allAccountSelected) {
        return tt('All');
    } else if (selectedAccountNames.length < 1) {
        return '';
    }

    return joinMultiText(selectedAccountNames);
}

function getFilteredTransactionCategoriesDisplayContent(filterTransactionCategoryIds?: Record<string, boolean>): string {
    if ((props.loading && !hasAnyTransactionCategory.value) || !transactionCategoriesStore.allTransactionCategoriesMap) {
        return '';
    }

    if (!filterTransactionCategoryIds) {
        return tt('All');
    }

    let allCategorySelected = true;
    const selectedCategoryNames: string[] = [];

    for (const transactionCategory of values(transactionCategoriesStore.allTransactionCategoriesMap)) {
        if (!transactionCategory.parentId || transactionCategory.parentId === '0') {
            continue;
        }

        if (!filterTransactionCategoryIds[transactionCategory.id]) {
            allCategorySelected = false;
        } else {
            selectedCategoryNames.push(transactionCategory.name);
        }
    }

    if (allCategorySelected) {
        return tt('All');
    } else if (selectedCategoryNames.length < 1) {
        return '';
    }

    return joinMultiText(selectedCategoryNames);
}

function addQuery(): void {
    queries.value.push(TransactionExploreQuery.create());
}

function updateQueryName(query: TransactionExploreQuery): void {
    query.name = editingQueryName.value;
    editingQuery.value = undefined;
    editingQueryName.value = '';
}

function cancelUpdateQueryName(): void {
    editingQuery.value = undefined;
    editingQueryName.value = '';
}

function duplicateQuery(query: TransactionExploreQuery): void {
    queries.value.push(query.clone());
}

function removeQuery(queryIndex: number): void {
    if (queries.value.length > 0) {
        queries.value.splice(queryIndex, 1);
    }

    const newShowExpression: Record<number, boolean> = {};

    for (const [key, state] of entries(showExpression.value)) {
        const index = parseInt(key);

        if (queryIndex > index) {
            newShowExpression[index] = state;
        } else if (queryIndex < index) {
            newShowExpression[index - 1] = state;
        }
    }

    showExpression.value = newShowExpression;

    if (queries.value.length < 1) {
        queries.value.push(TransactionExploreQuery.create());
    }
}

function clearAllQueries(): void {
    queries.value.length = 0;
    queries.value.push(TransactionExploreQuery.create());
}

function addCondition(queryIndex: number): void {
    const query = queries.value[queryIndex];

    if (!query) {
        return;
    }

    const newCondition = query.addNewCondition(TransactionExploreConditionField.TransactionType, query.conditions.length < 1);
    query.conditions.push(newCondition);
}

function removeCondition(queryIndex: number, conditionIndex: number): void {
    const query = queries.value[queryIndex];

    if (!query) {
        return;
    }

    query.conditions.splice(conditionIndex, 1);

    if (conditionIndex === 0 && query.conditions.length > 0) {
        const newFirstCondition = query.conditions[0];

        if (newFirstCondition) {
            newFirstCondition.relation = TransactionExploreConditionRelation.First;
        }
    }
}

function updateConditionField(queryIndex: number, conditionIndex: number, newField: TransactionExploreConditionField | undefined): void {
    if (!newField) {
        return;
    }

    const query = queries.value[queryIndex];

    if (!query) {
        return;
    }

    const oldConditionWithRelation = query.conditions[conditionIndex];

    if (!oldConditionWithRelation) {
        return;
    }

    const newConditionWithRelation = query.addNewCondition(newField, conditionIndex < 1);
    oldConditionWithRelation.condition = newConditionWithRelation.condition;
}

function updateSourceAccount(changed: boolean, selectedAccountIds?: string[]): void {
    if (!changed || !currentCondition.value || currentCondition.value.field !== TransactionExploreConditionField.SourceAccount.value) {
        showFilterSourceAccountsDialog.value = false;
        return;
    }

    currentCondition.value.value = selectedAccountIds || [];
    currentCondition.value = undefined;
    showFilterSourceAccountsDialog.value = false;
}

function updateDestinationAccount(changed: boolean, selectedAccountIds?: string[]): void {
    if (!changed || !currentCondition.value || currentCondition.value.field !== TransactionExploreConditionField.DestinationAccount.value) {
        showFilterDestinationAccountsDialog.value = false;
        return;
    }

    currentCondition.value.value = selectedAccountIds || [];
    currentCondition.value = undefined;
    showFilterDestinationAccountsDialog.value = false;
}

function updateTransactionCategories(changed: boolean, selectedCategoryIds?: string[]): void {
    if (!changed || !currentCondition.value || currentCondition.value.field !== TransactionExploreConditionField.TransactionCategory.value) {
        showFilterTransactionCategoriesDialog.value = false;
        return;
    }

    currentCondition.value.value = selectedCategoryIds || [];
    currentCondition.value = undefined;
    showFilterTransactionCategoriesDialog.value = false;
}

function getExpression(queryIndex: number): string {
    const query = queries.value[queryIndex];

    if (!query || !query.conditions || query.conditions.length < 1) {
        return '';
    }

    try {
        return query.toExpression(transactionCategoriesStore.allTransactionCategoriesMap, accountsStore.allAccountsMap, transactionTagsStore.allTransactionTagsMap);
    } catch (ex) {
        logger.error('failed to generate expression for explore query#' + queryIndex, ex);
        snackbar.value?.showError(tt('Failed to generate expression'));
        return tt('Failed to generate expression');
    }
}

if (queries.value.length === 0) {
    queries.value.push(TransactionExploreQuery.create());
}
</script>

<style>
.query-name {
    white-space: pre;
}

.query-name-edit {
    height: 36px;

    > .v-text-field {
        .v-field__input {
            margin-top: -2px;
            margin-bottom: -3px;
            padding-top: 0;
        }
    }
}
</style>
