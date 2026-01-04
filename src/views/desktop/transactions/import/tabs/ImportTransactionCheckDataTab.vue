<template>
    <v-data-table
        fixed-header
        fixed-footer
        show-select
        multi-sort
        density="compact"
        item-value="index"
        :class="{ 'import-transaction-table': true, 'disabled': !!disabled }"
        :height="importTransactionsTableHeight"
        :headers="importTransactionHeaders"
        :items="importTransactions"
        :hover="true"
        :search="JSON.stringify(filters)"
        :custom-filter="importTransactionsFilter"
        :no-data-text="tt('No data to import')"
        v-model:items-per-page="countPerPage"
        v-model:page="currentPage"
    >
        <template #header.data-table-select>
            <v-checkbox readonly class="always-cursor-pointer"
                        density="compact" width="28"
                        :disabled="!!disabled"
                        :indeterminate="anyButNotAllTransactionSelected"
                        v-model="allTransactionSelected"
            >
                <v-menu activator="parent" location="bottom">
                    <v-list>
                        <v-list-item :prepend-icon="mdiSelectAll"
                                     :title="tt('Select All Valid Items')"
                                     :disabled="!!disabled"
                                     @click="selectAllValid"></v-list-item>
                        <v-list-item :prepend-icon="mdiSelectAll"
                                     :title="tt('Select All Invalid Items')"
                                     :disabled="!!disabled"
                                     @click="selectAllInvalid"></v-list-item>
                        <v-divider class="my-2"/>
                        <v-list-item :prepend-icon="mdiSelectAll"
                                     :title="tt('Select All')"
                                     :disabled="!!disabled"
                                     @click="selectAll"></v-list-item>
                        <v-list-item :prepend-icon="mdiSelect"
                                     :title="tt('Select None')"
                                     :disabled="!!disabled"
                                     @click="selectNone"></v-list-item>
                        <v-list-item :prepend-icon="mdiSelectInverse"
                                     :title="tt('Invert Selection')"
                                     :disabled="!!disabled"
                                     @click="selectInvert"></v-list-item>
                        <v-divider class="my-2"/>
                        <v-list-item :prepend-icon="mdiSelectAll"
                                     :title="tt('Select All in This Page')"
                                     :disabled="!!disabled"
                                     @click="selectAllInThisPage"></v-list-item>
                        <v-list-item :prepend-icon="mdiSelect"
                                     :title="tt('Select None in This Page')"
                                     :disabled="!!disabled"
                                     @click="selectNoneInThisPage"></v-list-item>
                        <v-list-item :prepend-icon="mdiSelectInverse"
                                     :title="tt('Invert Selection in This Page')"
                                     :disabled="!!disabled"
                                     @click="selectInvertInThisPage"></v-list-item>
                    </v-list>
                </v-menu>
            </v-checkbox>
        </template>
        <template #item.data-table-select="{ item }">
            <v-checkbox density="compact"
                        :color="!item.valid ? 'error' : 'primary'"
                        :disabled="!!disabled"
                        v-model="item.selected"></v-checkbox>
        </template>
        <template #item.valid="{ item }">
            <v-icon size="small" :class="{ 'text-error': !item.valid }"
                    :disabled="!!disabled"
                    :icon="editingTransaction === item ? mdiCheck : mdiPencilOutline"
                    @click="editTransaction(item)">
            </v-icon>
            <v-tooltip activator="parent" v-if="!disabled">{{ tt('Edit') }}</v-tooltip>
        </template>
        <template #item.time="{ item }">
            <span>{{ getDisplayDateTime(item) }}</span>
            <v-chip class="ms-1" variant="flat" color="grey" size="x-small"
                    v-if="!isSameAsDefaultTimezoneOffsetMinutes(item)">{{ getDisplayTimezone(item) }}</v-chip>
        </template>
        <template #item.type="{ value }">
            <v-chip label color="secondary" variant="outlined" size="x-small" v-if="value === TransactionType.ModifyBalance">{{ tt('Modify Balance') }}</v-chip>
            <v-chip label class="text-income" variant="outlined" size="x-small" v-else-if="value === TransactionType.Income">{{ tt('Income') }}</v-chip>
            <v-chip label class="text-expense" variant="outlined" size="x-small" v-else-if="value === TransactionType.Expense">{{ tt('Expense') }}</v-chip>
            <v-chip label color="primary" variant="outlined" size="x-small" v-else-if="value === TransactionType.Transfer">{{ tt('Transfer') }}</v-chip>
            <v-chip label color="default" variant="outlined" size="x-small" v-else>{{ tt('Unknown') }}</v-chip>
        </template>
        <template #item.actualCategoryName="{ item }">
            <div class="d-flex align-center" v-if="editingTransaction !== item || item.type === TransactionType.ModifyBalance">
                <span v-if="item.type === TransactionType.ModifyBalance">-</span>
                <ItemIcon size="24px" icon-type="category"
                          :icon-id="allCategoriesMap[item.categoryId]?.icon ?? ''"
                          :color="allCategoriesMap[item.categoryId]?.color ?? ''"
                          v-if="item.type !== TransactionType.ModifyBalance && item.categoryId && item.categoryId !== '0' && allCategoriesMap[item.categoryId]"></ItemIcon>
                <span class="ms-2" v-if="item.type !== TransactionType.ModifyBalance && item.categoryId && item.categoryId !== '0' && allCategoriesMap[item.categoryId]">
                                    {{ allCategoriesMap[item.categoryId]?.name }}
                                </span>
                <div class="text-error font-italic" v-else-if="item.type !== TransactionType.ModifyBalance && (!item.categoryId || item.categoryId === '0' || !allCategoriesMap[item.categoryId])">
                    <v-icon class="me-1" :icon="mdiAlertOutline"/>
                    <span>{{ item.originalCategoryName }}</span>
                </div>
            </div>
            <div style="width: 260px" v-if="editingTransaction === item && item.type === TransactionType.Expense">
                <two-column-select density="compact" variant="plain"
                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                   secondary-hidden-field="hidden"
                                   :disabled="!!disabled || !hasVisibleExpenseCategories"
                                   :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                   :show-selection-primary-text="true"
                                   :custom-selection-primary-text="getTransactionPrimaryCategoryName(item.categoryId, allCategories[CategoryType.Expense])"
                                   :custom-selection-secondary-text="getTransactionSecondaryCategoryName(item.categoryId, allCategories[CategoryType.Expense])"
                                   :placeholder="tt('Category')"
                                   :items="allCategories[CategoryType.Expense]"
                                   v-model="item.categoryId">
                </two-column-select>
            </div>
            <div style="width: 260px" v-if="editingTransaction === item && item.type === TransactionType.Income">
                <two-column-select density="compact" variant="plain"
                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                   secondary-hidden-field="hidden"
                                   :disabled="!!disabled || !hasVisibleIncomeCategories"
                                   :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                   :show-selection-primary-text="true"
                                   :custom-selection-primary-text="getTransactionPrimaryCategoryName(item.categoryId, allCategories[CategoryType.Income])"
                                   :custom-selection-secondary-text="getTransactionSecondaryCategoryName(item.categoryId, allCategories[CategoryType.Income])"
                                   :placeholder="tt('Category')"
                                   :items="allCategories[CategoryType.Income]"
                                   v-model="item.categoryId">
                </two-column-select>
            </div>
            <div style="width: 260px" v-if="editingTransaction === item && item.type === TransactionType.Transfer">
                <two-column-select density="compact" variant="plain"
                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                   secondary-hidden-field="hidden"
                                   :disabled="!!disabled || !hasVisibleTransferCategories"
                                   :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                   :show-selection-primary-text="true"
                                   :custom-selection-primary-text="getTransactionPrimaryCategoryName(item.categoryId, allCategories[CategoryType.Transfer])"
                                   :custom-selection-secondary-text="getTransactionSecondaryCategoryName(item.categoryId, allCategories[CategoryType.Transfer])"
                                   :placeholder="tt('Category')"
                                   :items="allCategories[CategoryType.Transfer]"
                                   v-model="item.categoryId">
                </two-column-select>
            </div>
        </template>
        <template #item.sourceAmount="{ item }">
            <div class="d-flex align-center" v-if="editingTransaction !== item">
                <span>{{ getTransactionDisplayAmount(item) }}</span>
                <v-icon class="icon-with-direction mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId"></v-icon>
                <span v-if="item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId">{{ getTransactionDisplayDestinationAmount(item) }}</span>
            </div>
            <div class="d-flex align-center" :style="`width: ${item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId ? 250 : 100}px`" v-if="editingTransaction === item">
                <amount-input density="compact" variant="plain"
                              persistent-placeholder
                              :currency="item.originalSourceAccountCurrency || defaultCurrency"
                              :show-currency="true"
                              :disabled="!!disabled"
                              :placeholder="tt('Amount')"
                              v-model="item.sourceAmount"/>
                <v-icon class="icon-with-direction mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId"></v-icon>
                <amount-input density="compact" variant="plain"
                              persistent-placeholder
                              :currency="item.originalDestinationAccountCurrency || defaultCurrency"
                              :show-currency="true"
                              :disabled="!!disabled"
                              :placeholder="tt('Transfer In Amount')"
                              v-model="item.destinationAmount"
                              v-if="item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId"/>
            </div>
        </template>
        <template #item.actualSourceAccountName="{ item }">
            <div class="d-flex align-center" v-if="editingTransaction !== item">
                <span v-if="item.sourceAccountId && item.sourceAccountId !== '0' && allAccountsMap[item.sourceAccountId]">{{ allAccountsMap[item.sourceAccountId]?.name }}</span>
                <div class="text-error font-italic" v-else>
                    <v-icon class="me-1" :icon="mdiAlertOutline"/>
                    <span>{{ item.originalSourceAccountName }}</span>
                </div>
                <v-icon class="icon-with-direction mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer"></v-icon>
                <span v-if="item.type === TransactionType.Transfer && item.destinationAccountId && item.destinationAccountId !== '0' && allAccountsMap[item.destinationAccountId]">{{allAccountsMap[item.destinationAccountId]?.name }}</span>
                <div class="text-error font-italic" v-else-if="item.type === TransactionType.Transfer && (!item.destinationAccountId || item.destinationAccountId === '0' || !allAccountsMap[item.destinationAccountId])">
                    <v-icon class="me-1" :icon="mdiAlertOutline"/>
                    <span>{{ item.originalDestinationAccountName }}</span>
                </div>
            </div>
            <div class="d-flex align-center" :style="`width: ${item.type === TransactionType.Transfer ? 450 : 200}px`"  v-if="editingTransaction === item">
                <two-column-select density="compact" variant="plain"
                                   primary-key-field="id" primary-value-field="category"
                                   primary-title-field="name" primary-footer-field="displayBalance"
                                   primary-icon-field="icon" primary-icon-type="account"
                                   primary-sub-items-field="accounts"
                                   :primary-title-i18n="true"
                                   secondary-key-field="id" secondary-value-field="id"
                                   secondary-title-field="name" secondary-footer-field="displayBalance"
                                   secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                   :disabled="!!disabled || !allVisibleAccounts.length"
                                   :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                   :custom-selection-primary-text="getSourceAccountDisplayName(item)"
                                   :placeholder="getSourceAccountTitle(item)"
                                   :items="allVisibleCategorizedAccounts"
                                   v-model="item.sourceAccountId">
                </two-column-select>
                <v-icon class="icon-with-direction mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer"></v-icon>
                <two-column-select density="compact" variant="plain"
                                   primary-key-field="id" primary-value-field="category"
                                   primary-title-field="name" primary-footer-field="displayBalance"
                                   primary-icon-field="icon" primary-icon-type="account"
                                   primary-sub-items-field="accounts"
                                   :primary-title-i18n="true"
                                   secondary-key-field="id" secondary-value-field="id"
                                   secondary-title-field="name" secondary-footer-field="displayBalance"
                                   secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                   :disabled="!!disabled || !allVisibleAccounts.length"
                                   :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                   :custom-selection-primary-text="getDestinationAccountDisplayName(item)"
                                   :placeholder="tt('Destination Account')"
                                   :items="allVisibleCategorizedAccounts"
                                   v-model="item.destinationAccountId"
                                   v-if="item.type === TransactionType.Transfer">
                </two-column-select>
            </div>
        </template>
        <template #item.geoLocation="{ item }">
            <span v-if="item.geoLocation">{{ `(${formatCoordinate(item.geoLocation, coordinateDisplayType)})` }}</span>
            <span v-else-if="!item.geoLocation">{{ tt('None') }}</span>
        </template>
        <template #item.tagIds="{ item }">
            <div v-if="editingTransaction !== item">
                <v-chip class="transaction-tag" size="small"
                        :class="{ 'font-italic': !tagId || tagId === '0' || !allTagsMap[tagId] }"
                        :prepend-icon="tagId && tagId !== '0' && allTagsMap[tagId] ? mdiPound : mdiAlertOutline"
                        :color="tagId && tagId !== '0' && allTagsMap[tagId] ? 'default' : 'error'"
                        :text="tagId && tagId !== '0' && allTagsMap[tagId] ? allTagsMap[tagId].name : item.originalTagNames[index]"
                        :key="tagId"
                        v-for="(tagId, index) in item.tagIds"/>
                <v-chip class="transaction-tag" size="small"
                        :text="tt('None')"
                        v-if="!item.tagIds || !item.tagIds.length"/>
            </div>
            <div style="width: 200px" v-if="editingTransaction === item">
                <v-autocomplete
                    item-title="name"
                    item-value="id"
                    auto-select-first
                    persistent-placeholder
                    multiple
                    chips
                    closable-chips
                    density="compact" variant="plain"
                    :disabled="!!disabled"
                    :placeholder="tt('None')"
                    :items="allTags"
                    :no-data-text="tt('No available tag')"
                    v-model="editingTags"
                >
                    <template #chip="{ props, index }">
                        <v-chip :class="{ 'font-italic': !isTagValid(editingTags, index) }"
                                :prepend-icon="isTagValid(editingTags, index) ? mdiPound : mdiAlertOutline"
                                :color="isTagValid(editingTags, index) ? 'default' : 'error'"
                                :text="isTagValid(editingTags, index) ? allTagsMap[editingTags[index] as string]?.name : item.originalTagNames[index]"
                                v-bind="props"/>
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
            </div>
        </template>
        <template #item.comment="{ item }">
            <span v-if="editingTransaction !== item">{{ item.comment || '' }}</span>
            <div v-if="editingTransaction === item">
                <v-text-field style="width: 200px" type="text"
                              density="compact" variant="plain"
                              persistent-placeholder
                              :placeholder="tt('Description')"
                              :disabled="!!disabled"
                              v-model="item.comment" />
            </div>
        </template>
        <template #bottom>
            <div class="title-and-toolbar d-flex align-center text-no-wrap mt-2" v-if="importTransactions">
                <span :class="{ 'text-error': selectedInvalidTransactionCount > 0 }">
                    {{ tt('format.misc.selectedCount', { count: getDisplayCount(selectedImportTransactionCount), totalCount: getDisplayCount(importTransactions.length) }) }}
                </span>
                <v-spacer v-if="importTransactions.length > 10"/>
                <span v-if="importTransactions.length > 10">{{ tt('Transactions Per Page') }}</span>
                <v-select class="ms-2" density="compact" max-width="100"
                          item-title="name"
                          item-value="value"
                          :disabled="!!disabled"
                          :items="importTransactionsTablePageOptions"
                          v-model="countPerPage"
                          v-if="importTransactions.length > 10"
                />
                <pagination-buttons density="compact"
                                    :disabled="!!disabled"
                                    :totalPageCount="totalPageCount"
                                    v-model="currentPage"
                                    v-if="importTransactions.length > 10"></pagination-buttons>
            </div>
        </template>
    </v-data-table>

    <v-dialog width="640" v-model="showCustomAmountFilterDialog">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex align-center">
                    <h4 class="text-h4">{{ tt('Filter Amount') }}</h4>
                </div>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
                <div class="me-2 d-flex flex-column justify-center" v-if="currentAmountFilterType">
                    {{ tt(currentAmountFilterType.name) }}
                </div>
                <amount-input :currency="defaultCurrency"
                              v-model="currentAmountFilterValue1"/>
                <div class="ms-2 me-2 d-flex flex-column justify-center" v-if="currentAmountFilterType && currentAmountFilterType.paramCount === 2">
                    ~
                </div>
                <amount-input :currency="defaultCurrency"
                              v-model="currentAmountFilterValue2"
                              v-if="currentAmountFilterType && currentAmountFilterType.paramCount === 2"/>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn @click="showCustomAmountFilterDialog = false; filters.amount = currentAmountFilterType?.toTextualFilter(currentAmountFilterValue1, currentAmountFilterValue2) ?? null">{{ tt('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="showCustomAmountFilterDialog = false">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <v-dialog width="640" v-model="showCustomDescriptionDialog">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex align-center">
                    <h4 class="text-h4">{{ tt('Filter Description') }}</h4>
                </div>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
                <v-text-field
                    type="text"
                    persistent-placeholder
                    :label="tt('Description')"
                    :placeholder="tt('Description')"
                    v-model="currentDescriptionFilterValue"
                />
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :disabled="!currentDescriptionFilterValue" @click="showCustomDescriptionDialog = false; filters.description = currentDescriptionFilterValue">{{ tt('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="showCustomDescriptionDialog = false; currentDescriptionFilterValue = ''">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <date-range-selection-dialog :title="tt('Custom Date Range')"
                                 :min-time="filters.minDatetime"
                                 :max-time="filters.maxDatetime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="changeCustomDateFilter"
                                 @error="onShowDateRangeError" />
    <batch-replace-dialog ref="batchReplaceDialog" />
    <batch-replace-all-types-dialog ref="batchReplaceAllTypesDialog" />
    <batch-create-dialog ref="batchCreateDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import PaginationButtons from '@/components/desktop/PaginationButtons.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import BatchReplaceDialog, { type BatchReplaceDialogDataType } from '../dialogs/BatchReplaceDialog.vue';
import BatchReplaceAllTypesDialog from '../dialogs/BatchReplaceAllTypesDialog.vue';
import BatchCreateDialog, { type BatchCreateDialogDataType } from '../dialogs/BatchCreateDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { type NameValue, type NameNumeralValue, itemAndIndex, reversed, keys } from '@/core/base.ts';
import { type NumeralSystem, AmountFilterType } from '@/core/numeral.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';

import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';
import { ImportTransaction } from '@/models/imported_transaction.ts';

import {
    isString,
    isNumber,
    objectFieldToArrayItem
} from '@/lib/common.ts';
import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    parseDateTimeFromUnixTime,
    parseDateTimeFromUnixTimeWithTimezoneOffset
} from '@/lib/datetime.ts';
import { formatCoordinate } from '@/lib/coordinate.ts';
import {
    getAccountMapByName
} from '@/lib/account.ts';
import {
    transactionTypeToCategoryType,
    getSecondaryTransactionMapByName,
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName
} from '@/lib/category.ts';

import {
    mdiCheck,
    mdiArrowRight,
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiPencilOutline,
    mdiAlertOutline,
    mdiPound,
    mdiTextBoxEditOutline,
    mdiFilterOffOutline,
    mdiShapePlusOutline,
    mdiPencilBoxMultipleOutline,
    mdiNumericPositive1,
    mdiNumericNegative1
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;
type BatchReplaceDialogType = InstanceType<typeof BatchReplaceDialog>;
type BatchReplaceAllTypesDialogType = InstanceType<typeof BatchReplaceAllTypesDialog>;
type BatchCreateDialogType = InstanceType<typeof BatchCreateDialog>;

interface ImportTransactionCheckDataFilter {
    minDatetime: number | null; // minDatetime or maxDatetime is null for 'All Date Range', all are not null for 'Custom Date Range'
    maxDatetime: number | null;
    transactionType: TransactionType | null; // null for 'All Transaction Type'
    category: string | null | undefined; // null for 'All Category', undefined for 'Invalid Category'
    amount: string | null; // null for 'All Amount'
    account: string | null | undefined; // null for 'All Account', undefined for 'Invalid Account'
    tag: string | null | undefined; // null for 'All Tag', undefined for 'Invalid Tag'
    description: string | null; // null for 'All Description'
}

interface ImportTransactionCheckDataMenuGroup {
    title?: string;
    items: ImportTransactionCheckDataMenu[];
}

interface ImportTransactionCheckDataMenu {
    prependIcon?: string;
    title: string;
    subTitle?: string;
    appendIcon?: string;
    disabled?: boolean;
    divider?: boolean;
    onClick: () => void;
}

const props = defineProps<{
    importTransactions?: ImportTransaction[]
    disabled?: boolean;
}>();

const {
    tt,
    getCurrentNumeralSystemType,
    formatDateTimeToLongDateTime,
    formatAmountToLocalizedNumeralsWithCurrency,
    getCategorizedAccountsWithDisplayBalance
} = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const batchReplaceDialog = useTemplateRef<BatchReplaceDialogType>('batchReplaceDialog');
const batchReplaceAllTypesDialog = useTemplateRef<BatchReplaceAllTypesDialogType>('batchReplaceAllTypesDialog');
const batchCreateDialog = useTemplateRef<BatchCreateDialogType>('batchCreateDialog');

const editingTransaction = ref<ImportTransaction | null>(null);
const editingTags = ref<string[]>([]);
const filters = ref<ImportTransactionCheckDataFilter>({
    minDatetime: null,
    maxDatetime: null,
    transactionType: null,
    category: null,
    amount: null,
    account: null,
    tag: null,
    description: null
});

const currentPage = ref<number>(1);
const countPerPage = ref<number>(10);
const showCustomDateRangeDialog = ref<boolean>(false);
const showCustomAmountFilterDialog = ref<boolean>(false);
const showCustomDescriptionDialog = ref<boolean>(false);
const currentAmountFilterType = ref<AmountFilterType | null>(null);
const currentAmountFilterValue1 = ref<number>(0);
const currentAmountFilterValue2 = ref<number>(0);
const currentDescriptionFilterValue = ref<string | null>(null);

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);

const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
const coordinateDisplayType = computed<number>(() => userStore.currentUserCoordinateDisplayType);

const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value, customAccountCategoryOrder.value));
const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
const allAccountsMapByName = computed<Record<string, Account>>(() => getAccountMapByName(accountsStore.allAccounts));
const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
const allCategoriesMap = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);
const allTags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);
const allTagsMap = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);

const hasVisibleExpenseCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleExpenseCategories);
const hasVisibleIncomeCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleIncomeCategories);
const hasVisibleTransferCategories = computed<boolean>(() => transactionCategoriesStore.hasVisibleTransferCategories);

const isEditing = computed<boolean>(() => !!editingTransaction.value);
const canImport = computed<boolean>(() => selectedImportTransactionCount.value > 0 && selectedInvalidTransactionCount.value < 1);

const filterMenus = computed<ImportTransactionCheckDataMenuGroup[]>(() => [
    {
        items: [
            {
                title: tt('Clear All Filters'),
                prependIcon: mdiFilterOffOutline,
                disabled: filters.value.minDatetime === null
                    && filters.value.maxDatetime === null
                    && filters.value.transactionType === null
                    && filters.value.category === null
                    && filters.value.amount === null
                    && filters.value.account === null
                    && filters.value.tag === null
                    && filters.value.description === null,
                onClick: () => {
                    filters.value.minDatetime = null;
                    filters.value.maxDatetime = null;
                    filters.value.transactionType = null;
                    filters.value.category = null;
                    filters.value.amount = null;
                    filters.value.account = null;
                    filters.value.tag = null;
                    filters.value.description = null;
                }
            }
        ]
    },
    {
        title: tt('Date Range'),
        items: [
            {
                title: tt('All'),
                appendIcon: filters.value.minDatetime === null || filters.value.maxDatetime === null ? mdiCheck : undefined,
                onClick: () => {
                    filters.value.minDatetime = null;
                    filters.value.maxDatetime = null;
                }
            },
            {
                title: tt('Custom'),
                subTitle: displayFilterCustomDateRange.value,
                appendIcon: filters.value.minDatetime !== null && filters.value.maxDatetime !== null ? mdiCheck : undefined,
                onClick: () => showCustomDateRangeDialog.value = true
            }
        ]
    },
    {
        title: tt('Type'),
        items: [
            {
                title: tt('All'),
                appendIcon: filters.value.transactionType === null ? mdiCheck : undefined,
                onClick: () => filters.value.transactionType = null
            },
            {
                title: tt('Income'),
                appendIcon: filters.value.transactionType === TransactionType.Income ? mdiCheck : undefined,
                onClick: () => filters.value.transactionType = TransactionType.Income
            },
            {
                title: tt('Expense'),
                appendIcon: filters.value.transactionType === TransactionType.Expense ? mdiCheck : undefined,
                onClick: () => filters.value.transactionType = TransactionType.Expense
            },
            {
                title: tt('Transfer'),
                appendIcon: filters.value.transactionType === TransactionType.Transfer ? mdiCheck : undefined,
                onClick: () => filters.value.transactionType = TransactionType.Transfer
            }
        ]
    },
    {
        title: tt('Category'),
        items: [
            {
                title: tt('All'),
                appendIcon: filters.value.category === null ? mdiCheck : undefined,
                onClick: () => filters.value.category = null
            },
            {
                title: tt('Invalid Category'),
                appendIcon: filters.value.category === undefined ? mdiCheck : undefined,
                onClick: () => filters.value.category = undefined
            },
            {
                title: tt('None'),
                appendIcon: filters.value.category === '' ? mdiCheck : undefined,
                onClick: () => filters.value.category = ''
            },
            ...allUsedCategoryNames.value.map(name => ({
                title: name,
                appendIcon: filters.value.category === name ? mdiCheck : undefined,
                onClick: () => filters.value.category = name
            }))
        ]
    },
    {
        title: tt('Amount'),
        items: [
            {
                title: tt('All'),
                appendIcon: !filters.value.amount ? mdiCheck : undefined,
                onClick: () => filters.value.amount = null
            },
            ...AmountFilterType.values().map(filterType => ({
                title: tt(filterType.name),
                appendIcon: filters.value.amount && filters.value.amount.startsWith(`${filterType.type}:`) ? mdiCheck : undefined,
                onClick: () => {
                    let filterValue1: number = 0;
                    let filterValue2: number = 0;

                    if (filters.value.amount) {
                        const parts = filters.value.amount.split(':');

                        if (parts.length >= 2) {
                            filterValue1 = parseInt(parts[1] as string);
                        }

                        if (parts.length >= 3) {
                            filterValue2 = parseInt(parts[2] as string);
                        }
                    }

                    if (Number.isNaN(filterValue1) || !Number.isFinite(filterValue1)) {
                        filterValue1 = 0;
                    }

                    if (Number.isNaN(filterValue2) || !Number.isFinite(filterValue2)) {
                        filterValue2 = 0;
                    }

                    currentAmountFilterType.value = filterType;
                    currentAmountFilterValue1.value = filterValue1;
                    currentAmountFilterValue2.value = filterValue2;
                    showCustomAmountFilterDialog.value = true;
                }
            }))
        ]
    },
    {
        title: tt('Account'),
        items: [
            {
                title: tt('All'),
                appendIcon: filters.value.account === null ? mdiCheck : undefined,
                onClick: () => filters.value.account = null
            },
            {
                title: tt('Invalid Account'),
                appendIcon: filters.value.account === undefined ? mdiCheck : undefined,
                onClick: () => filters.value.account = undefined
            },
            {
                title: tt('None'),
                appendIcon: filters.value.account === '' ? mdiCheck : undefined,
                onClick: () => filters.value.account = ''
            },
            ...allUsedAccountNames.value.map(name => ({
                title: name,
                appendIcon: filters.value.account === name ? mdiCheck : undefined,
                onClick: () => filters.value.account = name
            }))
        ]
    },
    {
        title: tt('Tags'),
        items: [
            {
                title: tt('All'),
                appendIcon: filters.value.tag === null ? mdiCheck : undefined,
                onClick: () => filters.value.tag = null
            },
            {
                title: tt('Invalid Tag'),
                appendIcon: filters.value.tag === undefined ? mdiCheck : undefined,
                onClick: () => filters.value.tag = undefined
            },
            {
                title: tt('None'),
                appendIcon: filters.value.tag === '' ? mdiCheck : undefined,
                onClick: () => filters.value.tag = ''
            },
            ...allUsedTagNames.value.map(name => ({
                title: name,
                appendIcon: filters.value.tag === name ? mdiCheck : undefined,
                onClick: () => filters.value.tag = name
            }))
        ]
    },
    {
        title: tt('Description'),
        items: [
            {
                title: tt('All'),
                appendIcon: filters.value.description === null ? mdiCheck : undefined,
                onClick: () => filters.value.description = null
            },
            {
                title: tt('None'),
                appendIcon: filters.value.description === '' ? mdiCheck : undefined,
                onClick: () => filters.value.description = ''
            },
            {
                title: tt('Custom'),
                subTitle: filters.value.description !== null ? filters.value.description : undefined,
                appendIcon: filters.value.description !== null && filters.value.description !== '' ? mdiCheck : undefined,
                onClick: () => {
                    currentDescriptionFilterValue.value = filters.value.description || '';
                    showCustomDescriptionDialog.value = true;
                }
            }
        ]
    }
]);

const toolMenus = computed<ImportTransactionCheckDataMenu[]>(() => [
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Batch Replace Categories / Accounts / Tags'),
        disabled: isEditing.value,
        onClick: showReplaceAllTypesDialog
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Batch Replace Selected Expense Categories'),
        disabled: isEditing.value || selectedExpenseTransactionCount.value < 1,
        divider: true,
        onClick: () => showBatchReplaceDialog('expenseCategory')
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Batch Replace Selected Income Categories'),
        disabled: isEditing.value || selectedIncomeTransactionCount.value < 1,
        onClick: () => showBatchReplaceDialog('incomeCategory')
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Batch Replace Selected Transfer Categories'),
        disabled: isEditing.value || selectedTransferTransactionCount.value < 1,
        onClick: () => showBatchReplaceDialog('transferCategory')
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Batch Replace Selected Accounts'),
        disabled: isEditing.value || selectedImportTransactionCount.value < 1,
        onClick: () => showBatchReplaceDialog('account')
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Batch Replace Selected Destination Accounts'),
        disabled: isEditing.value || selectedTransferTransactionCount.value < 1,
        onClick: () => showBatchReplaceDialog('destinationAccount')
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Batch Replace Selected Transaction Tags'),
        disabled: isEditing.value || selectedImportTransactionCount.value < 1,
        onClick: () => showBatchReplaceDialog('tag', allOriginalTransactionTagNames.value)
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Batch Add Transaction Tags'),
        disabled: isEditing.value || selectedImportTransactionCount.value < 1,
        onClick: () => showBatchAddDialog('tag')
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Replace Invalid Expense Categories'),
        disabled: isEditing.value || !allInvalidExpenseCategoryNames.value || allInvalidExpenseCategoryNames.value.length < 1,
        divider: true,
        onClick: () => showReplaceInvalidItemDialog('expenseCategory', allInvalidExpenseCategoryNames.value)
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Replace Invalid Income Categories'),
        disabled: isEditing.value || !allInvalidIncomeCategoryNames.value || allInvalidIncomeCategoryNames.value.length < 1,
        onClick: () => showReplaceInvalidItemDialog('incomeCategory', allInvalidIncomeCategoryNames.value)
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Replace Invalid Transfer Categories'),
        disabled: isEditing.value || !allInvalidTransferCategoryNames.value || allInvalidTransferCategoryNames.value.length < 1,
        onClick: () => showReplaceInvalidItemDialog('transferCategory', allInvalidTransferCategoryNames.value)
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Replace Invalid Accounts'),
        disabled: isEditing.value || !allInvalidAccountNames.value || allInvalidAccountNames.value.length < 1,
        onClick: () => showReplaceInvalidItemDialog('account', allInvalidAccountNames.value)
    },
    {
        prependIcon: mdiTextBoxEditOutline,
        title: tt('Replace Invalid Transaction Tags'),
        disabled: isEditing.value || !allInvalidTransactionTagNames.value || allInvalidTransactionTagNames.value.length < 1,
        onClick: () => showReplaceInvalidItemDialog('tag', allInvalidTransactionTagNames.value)
    },
    {
        prependIcon: mdiShapePlusOutline,
        title: tt('Create Nonexistent Expense Categories'),
        disabled: isEditing.value || !allInvalidExpenseCategoryNames.value || allInvalidExpenseCategoryNames.value.length < 1,
        divider: true,
        onClick: () => showBatchCreateInvalidItemDialog('expenseCategory', allInvalidExpenseCategoryNames.value)
    },
    {
        prependIcon: mdiShapePlusOutline,
        title: tt('Create Nonexistent Income Categories'),
        disabled: isEditing.value || !allInvalidIncomeCategoryNames.value || allInvalidIncomeCategoryNames.value.length < 1,
        onClick: () => showBatchCreateInvalidItemDialog('incomeCategory', allInvalidIncomeCategoryNames.value)
    },
    {
        prependIcon: mdiShapePlusOutline,
        title: tt('Create Nonexistent Transfer Categories'),
        disabled: isEditing.value || !allInvalidTransferCategoryNames.value || allInvalidTransferCategoryNames.value.length < 1,
        onClick: () => showBatchCreateInvalidItemDialog('transferCategory', allInvalidTransferCategoryNames.value)
    },
    {
        prependIcon: mdiShapePlusOutline,
        title: tt('Create Nonexistent Transaction Tags'),
        disabled: isEditing.value || !allInvalidTransactionTagNames.value || allInvalidTransactionTagNames.value.length < 1,
        onClick: () => showBatchCreateInvalidItemDialog('tag', allInvalidTransactionTagNames.value)
    },
    {
        prependIcon: mdiPencilBoxMultipleOutline,
        title: tt('Batch Convert Expense Transaction to Income Transaction'),
        disabled: isEditing.value || selectedExpenseTransactionCount.value < 1,
        divider: true,
        onClick: () => convertTransactionType(TransactionType.Expense, TransactionType.Income)
    },
    {
        prependIcon: mdiPencilBoxMultipleOutline,
        title: tt('Batch Convert Expense Transaction to Transfer Transaction'),
        disabled: isEditing.value || selectedExpenseTransactionCount.value < 1,
        onClick: () => convertTransactionType(TransactionType.Expense, TransactionType.Transfer)
    },
    {
        prependIcon: mdiPencilBoxMultipleOutline,
        title: tt('Batch Convert Income Transaction to Expense Transaction'),
        disabled: isEditing.value || selectedIncomeTransactionCount.value < 1,
        onClick: () => convertTransactionType(TransactionType.Income, TransactionType.Expense)
    },
    {
        prependIcon: mdiPencilBoxMultipleOutline,
        title: tt('Batch Convert Income Transaction to Transfer Transaction'),
        disabled: isEditing.value || selectedIncomeTransactionCount.value < 1,
        onClick: () => convertTransactionType(TransactionType.Income, TransactionType.Transfer)
    },
    {
        prependIcon: mdiPencilBoxMultipleOutline,
        title: tt('Batch Convert Transfer Transaction to Expense Transaction'),
        disabled: isEditing.value || selectedTransferTransactionCount.value < 1,
        onClick: () => convertTransactionType(TransactionType.Transfer, TransactionType.Expense)
    },
    {
        prependIcon: mdiPencilBoxMultipleOutline,
        title: tt('Batch Convert Transfer Transaction to Income Transaction'),
        disabled: isEditing.value || selectedTransferTransactionCount.value < 1,
        onClick: () => convertTransactionType(TransactionType.Transfer, TransactionType.Income)
    },
    {
        prependIcon: mdiNumericPositive1,
        title: tt('Batch Convert Selected Amounts to Positive Values'),
        disabled: isEditing.value || selectedImportTransactionCount.value < 1,
        divider: true,
        onClick: () => convertTransactionAmountSign(1)
    },
    {
        prependIcon: mdiNumericNegative1,
        title: tt('Batch Convert Selected Amounts to Negative Values'),
        disabled: isEditing.value || selectedImportTransactionCount.value < 1,
        onClick: () => convertTransactionAmountSign(-1)
    }
]);

const importTransactionsTableHeight = computed<number | undefined>(() => {
    if (countPerPage.value <= 10 || !props.importTransactions || props.importTransactions.length <= 10) {
        return undefined;
    } else {
        return 380;
    }
});

const importTransactionHeaders = computed<object[]>(() => {
    return [
        { key: 'data-table-select', fixed: true },
        { value: 'valid', sortable: true, nowrap: true, width: 35, fixed: true },
        { value: 'time', title: tt('Transaction Time'), sortable: true, nowrap: true },
        { value: 'type', title: tt('Type'), sortable: true, nowrap: true },
        { value: 'actualCategoryName', title: tt('Category'), sortable: true, nowrap: true },
        { value: 'sourceAmount', title: tt('Amount'), sortable: true, nowrap: true },
        { value: 'actualSourceAccountName', title: tt('Account'), sortable: true, nowrap: true },
        { value: 'geoLocation', title: tt('Geographic Location'), sortable: true, nowrap: true },
        { value: 'tagIds', title: tt('Tags'), sortable: true, nowrap: true },
        { value: 'comment', title: tt('Description'), sortable: true, nowrap: true },
    ];
});

const importTransactionsTablePageOptions = computed<NameNumeralValue[]>(() => getTablePageOptions(props.importTransactions?.length));

const totalPageCount = computed<number>(() => {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return 1;
    }

    let count = 0;

    for (const importTransaction of props.importTransactions) {
        if (isTransactionDisplayed(importTransaction)) {
            count++;
        }
    }

    return Math.ceil(count / countPerPage.value);
});

const currentPageTransactions = computed<ImportTransaction[]>(() => {
    const ret: ImportTransaction[] = [];

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return ret;
    }

    const previousCount = Math.max(0, (currentPage.value - 1) * countPerPage.value);
    let count = 0;

    for (const importTransaction of props.importTransactions) {
        if (ret.length >= countPerPage.value) {
            break;
        }

        if (isTransactionDisplayed(importTransaction)) {
            if (count >= previousCount) {
                ret.push(importTransaction);
            }

            count++;
        }
    }

    return ret;
});

const selectedImportTransactionCount = computed<number>(() => {
    let count = 0;

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return count;
    }

    for (const importTransaction of props.importTransactions) {
        if (importTransaction.selected) {
            count++;
        }
    }

    return count;
});

const selectedExpenseTransactionCount = computed<number>(() => {
    let count = 0;

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return count;
    }

    for (const importTransaction of props.importTransactions) {
        if (importTransaction.selected && importTransaction.type === TransactionType.Expense) {
            count++;
        }
    }

    return count;
});

const selectedIncomeTransactionCount = computed<number>(() => {
    let count = 0;

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return count;
    }

    for (const importTransaction of props.importTransactions) {
        if (importTransaction.selected && importTransaction.type === TransactionType.Income) {
            count++;
        }
    }

    return count;
});

const selectedTransferTransactionCount = computed<number>(() => {
    let count = 0;

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return count;
    }

    for (const importTransaction of props.importTransactions) {
        if (importTransaction.selected && importTransaction.type === TransactionType.Transfer) {
            count++;
        }
    }

    return count;
});

const selectedInvalidTransactionCount = computed<number>(() => {
    let count = 0;

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return count;
    }

    for (const importTransaction of props.importTransactions) {
        if (!importTransaction.valid && importTransaction.selected) {
            count++;
        }
    }

    return count;
});

const anyButNotAllTransactionSelected = computed<boolean>(() => !!props.importTransactions && selectedImportTransactionCount.value > 0 && selectedImportTransactionCount.value !== props.importTransactions.length);
const allTransactionSelected = computed<boolean>(() => !!props.importTransactions && selectedImportTransactionCount.value === props.importTransactions.length);

const allUsedCategoryNames = computed<string[]>(() => {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return [];
    }

    const categoryNames: Record<string, boolean> = {};

    for (const transaction of props.importTransactions) {
        if (transaction.actualCategoryName && transaction.actualCategoryName !== '') {
            categoryNames[transaction.actualCategoryName] = true;
        }
    }

    return objectFieldToArrayItem(categoryNames);
});

const allUsedAccountNames = computed<string[]>(() => {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return [];
    }

    const accountNames: Record<string, boolean> = {};

    for (const transaction of props.importTransactions) {
        if (transaction.actualSourceAccountName && transaction.actualSourceAccountName !== '') {
            accountNames[transaction.actualSourceAccountName] = true;
        }

        if (transaction.actualDestinationAccountName && transaction.actualDestinationAccountName !== '') {
            accountNames[transaction.actualDestinationAccountName] = true;
        }
    }

    return objectFieldToArrayItem(accountNames);
});

const allUsedTagNames = computed<string[]>(() => {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return [];
    }

    const tagNames: Record<string, boolean> = {};

    for (const transaction of props.importTransactions) {
        if (!transaction.tagIds || !transaction.originalTagNames) {
            continue;
        }

        for (const [tagId, tagIndex] of itemAndIndex(transaction.tagIds)) {
            const originalTagName = transaction.originalTagNames[tagIndex] as string | undefined;

            if (tagId && tagId !== '0' && allTagsMap.value[tagId] && allTagsMap.value[tagId].name) {
                tagNames[allTagsMap.value[tagId].name] = true;
            } else if (originalTagName) {
                tagNames[originalTagName] = true;
            }
        }
    }

    return objectFieldToArrayItem(tagNames);
});

const allInvalidExpenseCategoryNames = computed<NameValue[]>(() => getCurrentInvalidCategoryNames(TransactionType.Expense));
const allInvalidIncomeCategoryNames = computed<NameValue[]>(() => getCurrentInvalidCategoryNames(TransactionType.Income));
const allInvalidTransferCategoryNames = computed<NameValue[]>(() => getCurrentInvalidCategoryNames(TransactionType.Transfer));
const allInvalidAccountNames = computed<NameValue[]>(() => getCurrentInvalidAccountNames());
const allInvalidTransactionTagNames = computed<NameValue[]>(() => getCurrentInvalidTagNames());
const allOriginalTransactionTagNames = computed<NameValue[]>(() => getAllOriginalTagNames());

const displayFilterCustomDateRange = computed<string>(() => {
    if (filters.value.minDatetime === null || filters.value.maxDatetime === null) {
        return '';
    }

    const minDisplayTime = formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(filters.value.minDatetime));
    const maxDisplayTime = formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(filters.value.maxDatetime));

    return `${minDisplayTime} - ${maxDisplayTime}`
});

function getDisplayCount(count: number): string {
    return numeralSystem.value.formatNumber(count);
}

function getTablePageOptions(linesCount?: number): NameNumeralValue[] {
    const pageOptions: NameNumeralValue[] = [];

    if (!linesCount || linesCount < 1) {
        pageOptions.push({ value: -1, name: tt('All') });
        return pageOptions;
    }

    for (const count of [ 5, 10, 15, 20, 25, 30, 50 ]) {
        if (linesCount < count) {
            break;
        }

        pageOptions.push({ value: count, name: getDisplayCount(count) });
    }

    pageOptions.push({ value: -1, name: tt('All') });

    return pageOptions;
}

function isTransactionDisplayed(transaction: ImportTransaction): boolean {
    if (isNumber(filters.value.minDatetime) && isNumber(filters.value.maxDatetime) && (transaction.time < filters.value.minDatetime || transaction.time > filters.value.maxDatetime)) {
        return false;
    }

    if (isNumber(filters.value.transactionType) && transaction.type !== filters.value.transactionType) {
        return false;
    }

    if (isString(filters.value.category)) {
        if (filters.value.category === '' && transaction.actualCategoryName !== '') {
            return false;
        } else if (filters.value.category !== '' && transaction.actualCategoryName !== filters.value.category) {
            return false;
        }
    } else if (filters.value.category === undefined) {
        if (transaction.type !== TransactionType.ModifyBalance && transaction.categoryId && transaction.categoryId !== '0') {
            return false;
        }
    }

    if (isString(filters.value.amount)) {
        const match: boolean = AmountFilterType.match(filters.value.amount, transaction.sourceAmount);

        if (!match) {
            return false;
        }
    }

    if (isString(filters.value.account)) {
        if (filters.value.account === '' && (transaction.actualSourceAccountName !== '' || transaction.actualDestinationAccountName !== '')) {
            return false;
        } else if (filters.value.account !== '' && transaction.actualSourceAccountName !== filters.value.account && transaction.actualDestinationAccountName !== filters.value.account) {
            return false;
        }
    } else if (filters.value.account === undefined) {
        if (transaction.type !== TransactionType.Transfer && transaction.sourceAccountId && transaction.sourceAccountId !== '0') {
            return false;
        } else if (transaction.type === TransactionType.Transfer && transaction.sourceAccountId && transaction.sourceAccountId !== '0' && transaction.destinationAccountId && transaction.destinationAccountId !== '0') {
            return false;
        }
    }

    if (isString(filters.value.tag)) {
        if (filters.value.tag === '' && transaction.tagIds && transaction.tagIds.length) {
            return false;
        } else if (filters.value.tag !== '') {
            let hasTagName = false;

            if (transaction.tagIds && transaction.tagIds.length) {
                for (const [tagId, tagIndex] of itemAndIndex(transaction.tagIds)) {
                    let tagName: string = transaction.originalTagNames ? (transaction.originalTagNames[tagIndex] ?? '') : '';

                    if (tagId && tagId !== '0' && allTagsMap.value[tagId] && allTagsMap.value[tagId].name) {
                        tagName = allTagsMap.value[tagId].name;
                    }

                    if (tagName === filters.value.tag) {
                        hasTagName = true;
                        break;
                    }
                }
            }

            if (!hasTagName) {
                return false;
            }
        }
    } else if (filters.value.tag === undefined) {
        if (transaction.tagIds && transaction.tagIds.length) {
            let hasInvalidTag = false;

            for (const tagId of transaction.tagIds) {
                if (!tagId || tagId === '0') {
                    hasInvalidTag = true;
                    break;
                }
            }

            if (!hasInvalidTag) {
                return false;
            }
        } else {
            return false;
        }
    }

    if (isString(filters.value.description)) {
        if (filters.value.description === '' && transaction.comment !== '') {
            return false;
        } else if (filters.value.description !== '' && transaction.comment.indexOf(filters.value.description) < 0) {
            return false;
        }
    }

    return true;
}

function isTagValid(tagIds: string[], tagIndex: number): boolean {
    if (!tagIds || !tagIds[tagIndex]) {
        return false;
    }

    if (tagIds[tagIndex] === '0') {
        return false;
    }

    const tagId = tagIds[tagIndex];
    return !!allTagsMap.value[tagId];
}

function getDisplayDateTime(transaction: ImportTransaction): string {
    const dateTime = parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset)
    return formatDateTimeToLongDateTime(dateTime);
}

function isSameAsDefaultTimezoneOffsetMinutes(transaction: ImportTransaction): boolean {
    return transaction.utcOffset === getTimezoneOffsetMinutes(transaction.time);
}

function getDisplayTimezone(transaction: ImportTransaction): string {
    return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
}

function getDisplayCurrency(value: number, currencyCode: string): string {
    return formatAmountToLocalizedNumeralsWithCurrency(value, currencyCode);
}

function getTransactionDisplayAmount(transaction: ImportTransaction): string {
    let currency = transaction.originalSourceAccountCurrency || defaultCurrency.value;

    if (transaction.sourceAccountId && transaction.sourceAccountId !== '0' && allAccountsMap.value[transaction.sourceAccountId]) {
        currency = allAccountsMap.value[transaction.sourceAccountId]!.currency;
    }

    return getDisplayCurrency(transaction.sourceAmount, currency);
}

function getTransactionDisplayDestinationAmount(transaction: ImportTransaction): string {
    if (transaction.type !== TransactionType.Transfer) {
        return '-';
    }

    let currency = transaction.originalDestinationAccountCurrency || defaultCurrency.value;

    if (transaction.destinationAccountId && transaction.destinationAccountId !== '0' && allAccountsMap.value[transaction.destinationAccountId]) {
        currency = allAccountsMap.value[transaction.destinationAccountId]!.currency;
    }

    return getDisplayCurrency(transaction.destinationAmount, currency);
}

function getSourceAccountTitle(transaction: ImportTransaction): string {
    if (transaction.type === TransactionType.Expense || transaction.type === TransactionType.Income) {
        return tt('Account');
    } else if (transaction.type === TransactionType.Transfer) {
        return tt('Source Account');
    } else {
        return tt('Account');
    }
}

function getSourceAccountDisplayName(transaction: ImportTransaction): string {
    if (transaction.sourceAccountId) {
        return Account.findAccountNameById(allAccounts.value, transaction.sourceAccountId) || '';
    } else {
        return tt('None');
    }
}

function getDestinationAccountDisplayName(transaction: ImportTransaction): string {
    if (transaction.destinationAccountId) {
        return Account.findAccountNameById(allAccounts.value, transaction.destinationAccountId) || '';
    } else {
        return tt('None');
    }
}

function getCurrentInvalidCategoryNames(transactionType: TransactionType): NameValue[] {
    const invalidCategoryNames: Record<string, boolean> = {};
    const invalidCategories: NameValue[] = [];

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return invalidCategories;
    }

    for (const importTransaction of props.importTransactions) {
        const categoryId = importTransaction.categoryId;

        if (importTransaction.type === transactionType && (!categoryId || categoryId === '0' || !allCategoriesMap.value[categoryId])) {
            invalidCategoryNames[importTransaction.originalCategoryName] = true;
        }
    }

    for (const name of keys(invalidCategoryNames)) {
        invalidCategories.push({
            name: name || tt('(Empty)'),
            value: name
        });
    }

    return invalidCategories;
}

function getCurrentInvalidAccountNames(): NameValue[] {
    const invalidAccountNames: Record<string, boolean> = {};
    const invalidAccounts: NameValue[] = [];

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return invalidAccounts;
    }

    for (const importTransaction of props.importTransactions) {
        const sourceAccountId = importTransaction.sourceAccountId;
        const destinationAccountId = importTransaction.destinationAccountId;

        if (!sourceAccountId || sourceAccountId === '0' || !allAccountsMap.value[sourceAccountId]) {
            invalidAccountNames[importTransaction.originalSourceAccountName] = true;
        }

        if (importTransaction.type === TransactionType.Transfer && isString(importTransaction.originalDestinationAccountName) && (!destinationAccountId || destinationAccountId === '0' || !allAccountsMap.value[destinationAccountId])) {
            invalidAccountNames[importTransaction.originalDestinationAccountName] = true;
        }
    }

    for (const name of keys(invalidAccountNames)) {
        invalidAccounts.push({
            name: name || tt('(Empty)'),
            value: name
        });
    }

    return invalidAccounts;
}

function getCurrentInvalidTagNames(): NameValue[] {
    const invalidTagNames: Record<string, boolean> = {};
    const invalidTags: NameValue[] = [];

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return invalidTags;
    }

    for (const importTransaction of props.importTransactions) {
        if (!importTransaction.tagIds || !importTransaction.originalTagNames) {
            continue;
        }

        for (const [tagId, tagIndex] of itemAndIndex(importTransaction.tagIds)) {
            const originalTagName = importTransaction.originalTagNames[tagIndex] as string | undefined;

            if (!originalTagName) {
                continue;
            }

            if (!tagId || tagId === '0' || !allTagsMap.value[tagId]) {
                invalidTagNames[originalTagName] = true;
            }
        }
    }

    for (const name of keys(invalidTagNames)) {
        invalidTags.push({
            name: name || tt('(Empty)'),
            value: name
        });
    }

    return invalidTags;
}

function getAllOriginalTagNames(): NameValue[] {
    const allOriginalTagNames: Record<string, boolean> = {};
    const allOriginalTags: NameValue[] = [];

    if (!props.importTransactions || props.importTransactions.length < 1) {
        return allOriginalTags;
    }

    for (const importTransaction of props.importTransactions) {
        if (!importTransaction.originalTagNames) {
            continue;
        }

        for (const tagName of importTransaction.originalTagNames) {
            allOriginalTagNames[tagName] = true;
        }
    }

    for (const name of keys(allOriginalTagNames)) {
        allOriginalTags.push({
            name: name || tt('(Empty)'),
            value: name
        });
    }

    return allOriginalTags;
}

function importTransactionsFilter(value: string, query: string, item?: { value: unknown, raw: ImportTransaction }): boolean {
    if (!item || !item.raw) {
        return false;
    }

    return isTransactionDisplayed(item.raw);
}

function selectAllValid(): void {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return;
    }

    for (const importTransaction of props.importTransactions) {
        if (importTransaction.valid && isTransactionDisplayed(importTransaction)) {
            importTransaction.selected = true;
        }
    }
}

function selectAllInvalid(): void {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return;
    }

    for (const importTransaction of props.importTransactions) {
        if (!importTransaction.valid && isTransactionDisplayed(importTransaction)) {
            importTransaction.selected = true;
        }
    }
}

function selectAll(): void {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return;
    }

    for (const importTransaction of props.importTransactions) {
        if (isTransactionDisplayed(importTransaction)) {
            importTransaction.selected = true;
        }
    }
}

function selectNone(): void {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return;
    }

    for (const importTransaction of props.importTransactions) {
        if (isTransactionDisplayed(importTransaction)) {
            importTransaction.selected = false;
        }
    }
}

function selectInvert(): void {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return;
    }

    for (const importTransaction of props.importTransactions) {
        if (isTransactionDisplayed(importTransaction)) {
            importTransaction.selected = !importTransaction.selected;
        }
    }
}

function selectAllInThisPage(): void {
    for (const importTransaction of currentPageTransactions.value) {
        importTransaction.selected = true;
    }
}

function selectNoneInThisPage(): void {
    for (const importTransaction of currentPageTransactions.value) {
        importTransaction.selected = false;
    }
}

function selectInvertInThisPage(): void {
    for (const importTransaction of currentPageTransactions.value) {
        importTransaction.selected = !importTransaction.selected;
    }
}

function editTransaction(transaction: ImportTransaction): void {
    if (editingTransaction.value) {
        editingTransaction.value.tagIds = editingTags.value;
        updateTransactionData(editingTransaction.value);
    }

    if (editingTransaction.value === transaction) {
        editingTags.value = [];
        editingTransaction.value = null;
    } else {
        editingTransaction.value = transaction;
        editingTags.value = editingTransaction.value.tagIds;
    }
}

function updateTransactionData(transaction: ImportTransaction): void {
    transaction.valid = transaction.isTransactionValid();

    if (transaction.categoryId && allCategoriesMap.value[transaction.categoryId]) {
        transaction.actualCategoryName = allCategoriesMap.value[transaction.categoryId]!.name;
    }

    if (transaction.sourceAccountId && allAccountsMap.value[transaction.sourceAccountId]) {
        transaction.actualSourceAccountName = allAccountsMap.value[transaction.sourceAccountId]!.name;
    }

    if (transaction.destinationAccountId && allAccountsMap.value[transaction.destinationAccountId]) {
        transaction.actualDestinationAccountName = allAccountsMap.value[transaction.destinationAccountId]!.name;
    }
}

function showBatchReplaceDialog(type: BatchReplaceDialogDataType, allSourceTagItems?: NameValue[]): void {
    if (isEditing.value) {
        return;
    }

    batchReplaceDialog.value?.open({
        mode: 'batchReplace',
        type: type,
        allSourceTagItems: allSourceTagItems
    }).then(result => {
        if (!result) {
            return;
        }

        if (type !== 'tag') {
            if (!result.targetItem) {
                return;
            }
        }

        let updatedCount = 0;

        if (props.importTransactions) {
            for (const importTransaction of props.importTransactions) {
                if (!importTransaction.selected) {
                    continue;
                }

                let updated = false;

                if (type === 'expenseCategory') {
                    if (importTransaction.type === TransactionType.Expense) {
                        importTransaction.categoryId = result.targetItem as string;
                        updated = true;
                    }
                } else if (type === 'incomeCategory') {
                    if (importTransaction.type === TransactionType.Income) {
                        importTransaction.categoryId = result.targetItem as string;
                        updated = true;
                    }
                } else if (type === 'transferCategory') {
                    if (importTransaction.type === TransactionType.Transfer) {
                        importTransaction.categoryId = result.targetItem as string;
                        updated = true;
                    }
                } else if (type === 'account') {
                    importTransaction.sourceAccountId = result.targetItem as string;
                    updated = true;
                } else if (type === 'destinationAccount') {
                    if (importTransaction.type === TransactionType.Transfer) {
                        importTransaction.destinationAccountId = result.targetItem as string;
                        updated = true;
                    }
                } else if (type === 'tag') {
                    const removeIndex: number[] = [];

                    for (let tagIndex = 0; tagIndex < importTransaction.originalTagNames.length; tagIndex++) {
                        const originalTagName = importTransaction.originalTagNames ? (importTransaction.originalTagNames[tagIndex] ?? '') : '';

                        if (originalTagName === result.sourceItem) {
                            if (result.targetItem) {
                                importTransaction.tagIds[tagIndex] = result.targetItem;
                                importTransaction.originalTagNames[tagIndex] = allTagsMap.value[result.targetItem]?.name || '';
                            } else {
                                removeIndex.push(tagIndex);
                            }
                            updated = true;
                        }
                    }

                    for (const tagIndex of reversed(removeIndex)) {
                        importTransaction.tagIds.splice(tagIndex, 1);
                        importTransaction.originalTagNames.splice(tagIndex, 1);
                    }
                }

                if (updated) {
                    updatedCount++;
                    updateTransactionData(importTransaction);
                }
            }
        }

        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: getDisplayCount(updatedCount)
            });
        }
    });
}

function showBatchAddDialog(type: BatchReplaceDialogDataType): void {
    if (isEditing.value) {
        return;
    }

    batchReplaceDialog.value?.open({
        mode: 'batchAdd',
        type: type
    }).then(result => {
        if (!result || !result.targetItem) {
            return;
        }

        let updatedCount = 0;

        if (props.importTransactions) {
            for (const importTransaction of props.importTransactions) {
                if (!importTransaction.selected) {
                    continue;
                }

                let updated = false;

                if (type === 'tag') {
                    let containsTag = false;

                    for (const tagName of importTransaction.originalTagNames) {
                        if (tagName === result.targetItem) {
                            containsTag = true;
                            break;
                        }
                    }

                    if (!containsTag) {
                        if (!importTransaction.tagIds) {
                            importTransaction.tagIds = [];
                        }

                        if (!importTransaction.originalTagNames) {
                            importTransaction.originalTagNames = [];
                        }

                        importTransaction.tagIds.push(result.targetItem);
                        importTransaction.originalTagNames.push(allTagsMap.value[result.targetItem]?.name ?? '');
                        updated = true;
                    }
                }

                if (updated) {
                    updatedCount++;
                    updateTransactionData(importTransaction);
                }
            }
        }

        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: getDisplayCount(updatedCount)
            });
        }
    });
}

function showReplaceInvalidItemDialog(type: BatchReplaceDialogDataType, invalidItems: NameValue[]): void {
    if (isEditing.value) {
        return;
    }

    batchReplaceDialog.value?.open({
        mode: 'replaceInvalidItems',
        type: type,
        invalidItems: invalidItems
    }).then(result => {
        if (!result || (!result.sourceItem && result.sourceItem !== '')) {
            return;
        }

        if (type !== 'tag') {
            if (!result.targetItem) {
                return;
            }
        }

        let updatedCount = 0;

        if (props.importTransactions) {
            for (const importTransaction of props.importTransactions) {
                if (importTransaction.valid) {
                    continue;
                }

                let updated = false;

                if (type === 'expenseCategory' || type === 'incomeCategory' || type === 'transferCategory') {
                    const categoryId = importTransaction.categoryId;
                    const originalCategoryName = importTransaction.originalCategoryName;

                    if (importTransaction.type !== TransactionType.ModifyBalance && originalCategoryName === result.sourceItem && (!categoryId || categoryId === '0' || !allCategoriesMap.value[categoryId])) {
                        if (type === 'expenseCategory' && importTransaction.type === TransactionType.Expense) {
                            importTransaction.categoryId = result.targetItem as string;
                            updated = true;
                        } else if (type === 'incomeCategory' && importTransaction.type === TransactionType.Income) {
                            importTransaction.categoryId = result.targetItem as string;
                            updated = true;
                        } else if (type === 'transferCategory' && importTransaction.type === TransactionType.Transfer) {
                            importTransaction.categoryId = result.targetItem as string;
                            updated = true;
                        }
                    }
                } else if (type === 'account') {
                    const sourceAccountId = importTransaction.sourceAccountId;
                    const originalSourceAccountName = importTransaction.originalSourceAccountName;
                    const destinationAccountId = importTransaction.destinationAccountId;
                    const originalDestinationAccountName = importTransaction.originalDestinationAccountName;

                    if (originalSourceAccountName === result.sourceItem && (!sourceAccountId || sourceAccountId === '0' || !allAccountsMap.value[sourceAccountId])) {
                        importTransaction.sourceAccountId = result.targetItem as string;
                        updated = true;
                    }

                    if (importTransaction.type === TransactionType.Transfer && originalDestinationAccountName === result.sourceItem && (!destinationAccountId || destinationAccountId === '0' || !allAccountsMap.value[destinationAccountId])) {
                        importTransaction.destinationAccountId = result.targetItem as string;
                        updated = true;
                    }
                } else if (type === 'tag' && importTransaction.tagIds) {
                    const removeIndex: number[] = [];

                    for (let tagIndex = 0; tagIndex < importTransaction.tagIds.length; tagIndex++) {
                        const tagId = importTransaction.tagIds[tagIndex] as string;
                        const originalTagName = importTransaction.originalTagNames ? (importTransaction.originalTagNames[tagIndex] ?? '') : '';

                        if (originalTagName === result.sourceItem && (!tagId || tagId === '0' || !allTagsMap.value[tagId])) {
                            if (result.targetItem) {
                                importTransaction.tagIds[tagIndex] = result.targetItem;
                                importTransaction.originalTagNames[tagIndex] = allTagsMap.value[result.targetItem]?.name || '';
                            } else {
                                removeIndex.push(tagIndex);
                            }
                            updated = true;
                        }
                    }

                    for (const tagIndex of reversed(removeIndex)) {
                        importTransaction.tagIds.splice(tagIndex, 1);
                        importTransaction.originalTagNames.splice(tagIndex, 1);
                    }
                }

                if (updated) {
                    updatedCount++;
                    updateTransactionData(importTransaction);
                }
            }
        }

        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: getDisplayCount(updatedCount)
            });
        }
    });
}

function showReplaceAllTypesDialog(): void {
    if (isEditing.value) {
        return;
    }

    batchReplaceAllTypesDialog.value?.open({
        expenseCategoryNames: allInvalidExpenseCategoryNames.value,
        incomeCategoryNames: allInvalidIncomeCategoryNames.value,
        transferCategoryNames: allInvalidTransferCategoryNames.value,
        accountNames: allInvalidAccountNames.value,
        tagNames: allInvalidTransactionTagNames.value
    }).then(result => {
        if (!result || !result.rules) {
            return;
        }

        let updatedCount = 0;

        if (props.importTransactions) {
            for (const importTransaction of props.importTransactions) {
                let updated = false;

                for (const rule of result.rules) {
                    if (!rule || !rule.dataType || !rule.targetId) {
                        continue;
                    }

                    if (rule.dataType === 'expenseCategory' || rule.dataType === 'incomeCategory' || rule.dataType === 'transferCategory') {
                        if (importTransaction.type !== TransactionType.ModifyBalance && importTransaction.originalCategoryName === rule.sourceValue) {
                            if (rule.dataType === 'expenseCategory' && importTransaction.type === TransactionType.Expense) {
                                importTransaction.categoryId = rule.targetId;
                                updated = true;
                            } else if (rule.dataType === 'incomeCategory' && importTransaction.type === TransactionType.Income) {
                                importTransaction.categoryId = rule.targetId;
                                updated = true;
                            } else if (rule.dataType === 'transferCategory' && importTransaction.type === TransactionType.Transfer) {
                                importTransaction.categoryId = rule.targetId;
                                updated = true;
                            }
                        }
                    } else if (rule.dataType === 'account') {
                        if (importTransaction.originalSourceAccountName === rule.sourceValue) {
                            importTransaction.sourceAccountId = rule.targetId;
                            updated = true;
                        }

                        if (importTransaction.type === TransactionType.Transfer && importTransaction.originalDestinationAccountName === rule.sourceValue) {
                            importTransaction.destinationAccountId = rule.targetId;
                            updated = true;
                        }
                    } else if (rule.dataType === 'tag' && importTransaction.tagIds) {
                        for (let tagIndex = 0; tagIndex < importTransaction.tagIds.length; tagIndex++) {
                            const originalTagName = importTransaction.originalTagNames ? (importTransaction.originalTagNames[tagIndex] ?? '') : '';

                            if (originalTagName === rule.sourceValue) {
                                importTransaction.tagIds[tagIndex] = rule.targetId;
                                updated = true;
                            }
                        }
                    }
                }

                if (updated) {
                    updatedCount++;
                    updateTransactionData(importTransaction);
                }
            }
        }

        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: getDisplayCount(updatedCount)
            });
        }
    });
}

function showBatchCreateInvalidItemDialog(type: BatchCreateDialogDataType, invalidItems: NameValue[]): void {
    if (isEditing.value) {
        return;
    }

    batchCreateDialog.value?.open({
        type: type,
        invalidItems: invalidItems
    }).then(result => {
        if (!result || !result.sourceTargetMap) {
            return;
        }

        let updatedCount = 0;

        if (props.importTransactions) {
            const sourceTargetMap: Record<string, string> = result.sourceTargetMap;

            for (const importTransaction of props.importTransactions) {
                if (importTransaction.valid) {
                    continue;
                }

                let updated = false;

                if (type === 'expenseCategory' || type === 'incomeCategory' || type === 'transferCategory') {
                    const categoryId = importTransaction.categoryId;
                    const originalCategoryName = importTransaction.originalCategoryName;
                    const targetItem = sourceTargetMap[originalCategoryName];

                    if (importTransaction.type !== TransactionType.ModifyBalance && targetItem && (!categoryId || categoryId === '0' || !allCategoriesMap.value[categoryId])) {
                        if (type === 'expenseCategory' && importTransaction.type === TransactionType.Expense) {
                            importTransaction.categoryId = targetItem;
                            updated = true;
                        } else if (type === 'incomeCategory' && importTransaction.type === TransactionType.Income) {
                            importTransaction.categoryId = targetItem;
                            updated = true;
                        } else if (type === 'transferCategory' && importTransaction.type === TransactionType.Transfer) {
                            importTransaction.categoryId = targetItem;
                            updated = true;
                        }
                    }
                } else if (type === 'tag' && importTransaction.tagIds) {
                    for (let tagIndex = 0; tagIndex < importTransaction.tagIds.length; tagIndex++) {
                        const tagId = importTransaction.tagIds[tagIndex] as string;
                        const originalTagName = importTransaction.originalTagNames ? (importTransaction.originalTagNames[tagIndex] ?? '') : '';
                        const targetItem = sourceTargetMap[originalTagName];

                        if (targetItem && (!tagId || tagId === '0' || !allTagsMap.value[tagId])) {
                            importTransaction.tagIds[tagIndex] = targetItem;
                            updated = true;
                        }
                    }
                }

                if (updated) {
                    updatedCount++;
                    updateTransactionData(importTransaction);
                }
            }
        }

        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: getDisplayCount(updatedCount)
            });
        }
    });
}

function convertTransactionType(fromType: TransactionType, toType: TransactionType): void {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return;
    }

    const categoryType = transactionTypeToCategoryType(toType);

    if (!categoryType) {
        return;
    }

    const categoryMapByName: Record<string, TransactionCategory> = getSecondaryTransactionMapByName(allCategories.value[categoryType]);

    for (const importTransaction of props.importTransactions) {
        if (!importTransaction.selected || importTransaction.type !== fromType) {
            continue;
        }

        importTransaction.type = toType;
        importTransaction.categoryId = categoryMapByName[importTransaction.originalCategoryName]?.id || '0';

        if (importTransaction.type === TransactionType.Transfer) {
            importTransaction.destinationAccountId = allAccountsMapByName.value[importTransaction.originalDestinationAccountName || '']?.id || '0';
            importTransaction.destinationAmount = importTransaction.sourceAmount;
        } else {
            if (fromType === TransactionType.Transfer && toType === TransactionType.Income) {
                importTransaction.sourceAccountId = importTransaction.destinationAccountId;
                importTransaction.sourceAmount = importTransaction.destinationAmount;
            }

            importTransaction.destinationAccountId = '0';
            importTransaction.destinationAmount = 0;
        }

        updateTransactionData(importTransaction);
    }
}

function convertTransactionAmountSign(toSign: number): void {
    if (!props.importTransactions || props.importTransactions.length < 1) {
        return;
    }

    for (const importTransaction of props.importTransactions) {
        if (!importTransaction.selected) {
            continue;
        }

        if (toSign > 0) {
            importTransaction.sourceAmount = Math.abs(importTransaction.sourceAmount);
            importTransaction.destinationAmount = Math.abs(importTransaction.destinationAmount);
        } else if (toSign < 0) {
            importTransaction.sourceAmount = -Math.abs(importTransaction.sourceAmount);
            importTransaction.destinationAmount = -Math.abs(importTransaction.destinationAmount);
        }

        updateTransactionData(importTransaction);
    }
}

function changeCustomDateFilter(minTime: number, maxTime: number): void {
    filters.value.minDatetime = minTime;
    filters.value.maxDatetime = maxTime;
    showCustomDateRangeDialog.value = false;
}

function onShowDateRangeError(message: string): void {
    snackbar.value?.showError(message);
}

function reset(): void {
    editingTransaction.value = null;
    editingTags.value = [];
    filters.value.minDatetime = null;
    filters.value.maxDatetime = null;
    filters.value.transactionType = null;
    filters.value.category = null;
    filters.value.account = null;
    filters.value.tag = null;
    filters.value.description = null;
    currentPage.value = 1;
    countPerPage.value = 10;
}

function setCountPerPage(count: number): void {
    countPerPage.value = count;
}

defineExpose({
    filterMenus,
    toolMenus,
    isEditing,
    canImport,
    reset,
    setCountPerPage
});
</script>

<style>
.import-transaction-table > .v-table__wrapper > table {
    th:not(:last-child),
    td:not(:last-child) {
        width: auto !important;
        white-space: nowrap;
    }

    th:last-child,
    td:last-child {
        width: 100% !important;
    }
}

.import-transaction-table .v-autocomplete.v-input.v-input--density-compact:not(.v-textarea) .v-field__input,
.import-transaction-table .v-select.v-input.v-input--density-compact:not(.v-textarea) .v-field__input {
    min-height: inherit;
    padding-top: 4px;
}

.import-transaction-table .v-chip.transaction-tag {
    margin-inline-end: 4px;
    margin-top: 2px;
    margin-bottom: 2px;
}

.import-transaction-table .v-chip.transaction-tag > .v-chip__content {
    display: block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}

.import-transaction-table .v-text-field.v-input.v-input--density-compact .v-field__input {
    padding-top: 0;
}
</style>
