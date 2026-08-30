<template>
    <v-dialog width="600" v-model="showState">
        <one-column-dialog-layout :title="tt('Widget Settings')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined" @click="confirm">{{ tt('Apply') }}</v-btn>
            </template>

            <template #content>
                <div class="mt-1">
                    <template :key="setting.settingName" v-for="(setting, index) in supportsSettings">
                        <v-select :class="{ 'mt-4': index > 0 }" item-title="name" item-value="value"
                                  :label="tt(setting.displayName)" :items="getItemCountOptions(setting.itemCountValues)"
                                  :model-value="getSettingValue(setting.settingName)"
                                  @update:model-value="updateSettingValue(setting, $event)"
                                  v-if="setting.settingType === 'itemCountSelect'" />

                        <v-select :class="{ 'mt-4': index > 0 }" item-title="name" item-value="value"
                                  :label="tt(setting.displayName)" :items="getMonthOptions(setting.monthValues)"
                                  :model-value="getSettingValue(setting.settingName)"
                                  @update:model-value="updateSettingValue(setting, $event)"
                                  v-else-if="setting.settingType === 'monthSelect'" />

                        <v-text-field class="always-cursor-pointer text-field-truncate"
                                      :class="{ 'mt-4': index > 0 }"
                                      item-title="displayName"
                                      item-value="type"
                                      persistent-placeholder
                                      :readonly="true"
                                      :disabled="loading"
                                      :label="tt(setting.displayName)"
                                      :placeholder="tt('Default')"
                                      :model-value="getFilteredAccountsDisplayContent(setting)"
                                      @click="currentSettingItem = setting; showFilterAccountsDialog = true"
                                      v-else-if="setting.settingType === 'accountSelect'" />

                        <v-text-field class="always-cursor-pointer text-field-truncate"
                                      :class="{ 'mt-4': index > 0 }"
                                      item-title="displayName"
                                      item-value="type"
                                      persistent-placeholder
                                      :readonly="true"
                                      :disabled="loading"
                                      :label="tt(setting.displayName)"
                                      :placeholder="tt('Default')"
                                      :model-value="getFilteredTransactionCategoriesDisplayContent(setting)"
                                      @click="currentSettingItem = setting; showFilterTransactionCategoriesDialog = true"
                                      v-else-if="setting.settingType === 'categorySelect'" />

                        <v-text-field class="always-cursor-pointer text-field-truncate"
                                      :class="{ 'mt-4': index > 0 }"
                                      item-title="displayName"
                                      item-value="type"
                                      persistent-placeholder
                                      :readonly="true"
                                      :disabled="loading"
                                      :label="tt(setting.displayName)"
                                      :placeholder="tt('Default')"
                                      :model-value="getSettingValue(setting.settingName) ? tt('Custom') : tt('Default')"
                                      @click="currentSettingItem = setting; showFilterTransactionTagsDialog = true"
                                      v-else-if="setting.settingType === 'tagSelect'" />

                        <v-select :class="{ 'mt-4': index > 0 }" item-title="name" item-value="value"
                                  :label="tt(setting.displayName)" :items="getCustomSelectOptions(setting)"
                                  :multiple="setting.multiple" :chips="setting.multiple" :closable-chips="setting.multiple"
                                  :model-value="getSettingValue(setting.settingName)"
                                  @update:model-value="updateSettingValue(setting, $event)"
                                  v-else-if="setting.settingType === 'customSelect'" />

                        <v-switch :class="{ 'mt-2': index > 0 }" :label="tt(setting.displayName)"
                                  :model-value="getSettingValue(setting.settingName)"
                                  @update:model-value="updateSettingValue(setting, $event)"
                                  v-else-if="setting.settingType === 'switch'" />

                        <amount-filter-input :class="{ 'mt-4': index > 0 }" :label="tt(setting.displayName)"
                                             :model-value="getSettingValue(setting.settingName) as string || ''"
                                             @update:model-value="updateSettingValue(setting, $event)"
                                             v-else-if="setting.settingType === 'amount'" />

                        <v-text-field :class="{ 'mt-4': index > 0 }" :label="tt(setting.displayName)"
                                      :placeholder="setting.placeholder ? tt(setting.placeholder) : undefined"
                                      :persistent-placeholder="!!setting.placeholder"
                                      :model-value="getSettingValue(setting.settingName)"
                                      @update:model-value="updateSettingValue(setting, $event)"
                                      v-else-if="setting.settingType === 'textbox'" />
                    </template>
                </div>

                <div class="text-body-large text-medium-emphasis text-center" v-if="widget && !supportsSettings.length">{{ tt('This widget has no configurable parameters') }}</div>
            </template>
        </one-column-dialog-layout>

        <snack-bar ref="snackbar" />

        <account-filter-settings-dialog type="custom"
                                        :selected-account-ids="(currentSettingItem?.settingType === 'accountSelect' && currentSettingItem?.settingName && isArray(getSettingValue(currentSettingItem.settingName))) ? getSettingValue(currentSettingItem.settingName) as string[] : []"
                                        :disable-hidden-account="(currentSettingItem?.settingType === 'accountSelect') ? (currentSettingItem?.disableHiddenAccounts ?? false) : false"
                                        v-model:show="showFilterAccountsDialog"
                                        @settings:change="updateAccountValue"/>

        <category-filter-settings-dialog type="custom"
                                         :selected-category-ids="(currentSettingItem?.settingType === 'categorySelect' && currentSettingItem?.settingName && isArray(getSettingValue(currentSettingItem.settingName))) ? getSettingValue(currentSettingItem.settingName) as string[] : []"
                                         v-model:show="showFilterTransactionCategoriesDialog"
                                         @settings:change="updateCategoryValue"/>

        <transaction-tag-filter-settings-dialog type="custom"
                                                :tag-filter="(currentSettingItem?.settingType === 'tagSelect' && currentSettingItem?.settingName && isString(getSettingValue(currentSettingItem.settingName))) ? getSettingValue(currentSettingItem.settingName) as string : ''"
                                                v-model:show="showFilterTransactionTagsDialog"
                                                @settings:change="updateTagValue"/>
    </v-dialog>
</template>

<script setup lang="ts">
import AccountFilterSettingsDialog from '@/views/desktop/common/dialogs/AccountFilterSettingsDialog.vue';
import CategoryFilterSettingsDialog from '@/views/desktop/common/dialogs/CategoryFilterSettingsDialog.vue';
import TransactionTagFilterSettingsDialog from '@/views/desktop/common/dialogs/TransactionTagFilterSettingsDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { type GenericNameValue, type NameNumeralValue, values } from '@/core/base.ts';
import { AccountType } from '@/core/account.ts';
import {
    type OverviewWidgetSettingValue,
    type OverviewWidgetCustomSelectSettingItem,
    type OverviewWidgetSettingItem,
    type DesktopOverviewWidgetLayout
} from '@/core/overview_layout.ts';

import { DESKTOP_OVERVIEW_WIDGET_DEFINITIONS } from '@/consts/overview_layout.ts';

import { isDefined, isArray, isString, isObjectEmpty, arrayItemToObjectField } from '@/lib/common.ts';
import { cloneWidget } from '@/lib/overview_layout.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    joinMultiText,
    formatNumberToLocalizedNumerals,
    getTablePageOptions
} = useI18n();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

let resolveFunc: ((widget: DesktopOverviewWidgetLayout) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const widget = ref<DesktopOverviewWidgetLayout | null>(null);
const currentSettingItem = ref<OverviewWidgetSettingItem | undefined>(undefined);
const showFilterAccountsDialog = ref<boolean>(false);
const showFilterTransactionCategoriesDialog = ref<boolean>(false);
const showFilterTransactionTagsDialog = ref<boolean>(false);

const hasAnyAccount = computed<boolean>(() => accountsStore.allPlainAccounts.length > 0);
const hasAnyTransactionCategory = computed<boolean>(() => !isObjectEmpty(transactionCategoriesStore.allTransactionCategoriesMap));
const supportsSettings = computed<OverviewWidgetSettingItem[]>(() => widget.value ? DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]?.supportsSettings ?? [] : []);

function getItemCountOptions(values: number[]): NameNumeralValue[] {
    return getTablePageOptions(values, undefined, false, true);
}

function getMonthOptions(values: number[]): NameNumeralValue[] {
    return values.map(value => ({
        name: tt('format.misc.nMonths', { n: formatNumberToLocalizedNumerals(value) }),
        value: value
    }));
}

function getCustomSelectOptions(setting: OverviewWidgetCustomSelectSettingItem): GenericNameValue<string | number>[] {
    return setting.selectValues.map(item => ({ name: tt(item.name), value: item.value }));
}

function getSettingValue(settingName: string): OverviewWidgetSettingValue | undefined {
    return widget.value?.settings[settingName];
}

function updateSettingValue(setting: OverviewWidgetSettingItem, value: OverviewWidgetSettingValue | null): void {
    if (!widget.value || value === null) {
        return;
    }

    if (setting.settingType === 'accountSelect' || setting.settingType === 'categorySelect') {
        if (isDefined(value) && isArray(value)) {
            widget.value.settings[setting.settingName] = value;
        } else {
            widget.value.settings[setting.settingName] = [];
        }
    } else if (setting.settingType === 'customSelect' && setting.multiple && isDefined(setting.allValue) && isArray(value)) {
        const previousValue = widget.value.settings[setting.settingName];
        const previousValues = isArray(previousValue) ? previousValue : [];
        const allValue = setting.allValue as number;

        if (value.includes(allValue) && !previousValues.includes(allValue)) {
            widget.value.settings[setting.settingName] = [allValue];
        } else {
            const selectedValues = value.filter(item => item !== allValue);
            widget.value.settings[setting.settingName] = selectedValues.length ? selectedValues : [allValue];
        }
    } else if (setting.settingType === 'customSelect' && setting.multiple && isArray(value)) {
        if (value.length >= (setting.minSelections ?? 0)) {
            widget.value.settings[setting.settingName] = value;
        }
    } else {
        widget.value.settings[setting.settingName] = value;
    }
}

function getFilteredAccountsDisplayContent(setting: OverviewWidgetSettingItem): string {
    if ((loading.value && !hasAnyAccount.value) || !accountsStore.allVisiblePlainAccounts || !accountsStore.allVisiblePlainAccounts.length) {
        return '';
    }

    const settingValue = getSettingValue(setting.settingName) as string[] | undefined;

    if (!settingValue || !isArray(settingValue) || !settingValue.length) {
        return tt('All');
    }

    const filterAccountIds: Record<string, boolean> = arrayItemToObjectField(settingValue, true);

    let allAccountSelected = true;
    const selectedAccountNames: string[] = [];

    for (const account of accountsStore.allVisiblePlainAccounts) {
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

function getFilteredTransactionCategoriesDisplayContent(setting: OverviewWidgetSettingItem): string {
    if ((loading.value && !hasAnyTransactionCategory.value) || !transactionCategoriesStore.allTransactionCategoriesMap) {
        return '';
    }

    const settingValue = getSettingValue(setting.settingName) as string[] | undefined;

    if (!settingValue || !isArray(settingValue) || !settingValue.length) {
        return tt('All');
    }

    const filterTransactionCategoryIds: Record<string, boolean> = arrayItemToObjectField(settingValue, true);

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

function updateAccountValue(changed: boolean, selectedAccountIds?: string[]): void {
    if (!changed || !currentSettingItem.value || currentSettingItem.value.settingType !== 'accountSelect') {
        showFilterAccountsDialog.value = false;
        return;
    }

    updateSettingValue(currentSettingItem.value, selectedAccountIds ?? []);
    currentSettingItem.value = undefined;
    showFilterAccountsDialog.value = false;
}

function updateCategoryValue(changed: boolean, selectedCategoryIds?: string[]): void {
    if (!changed || !currentSettingItem.value || currentSettingItem.value.settingType !== 'categorySelect') {
        showFilterTransactionCategoriesDialog.value = false;
        return;
    }

    updateSettingValue(currentSettingItem.value, selectedCategoryIds ?? []);
    currentSettingItem.value = undefined;
    showFilterTransactionCategoriesDialog.value = false;
}

function updateTagValue(changed: boolean, textualTagFilter?: string): void {
    if (!changed || !currentSettingItem.value || currentSettingItem.value.settingType !== 'tagSelect') {
        showFilterTransactionTagsDialog.value = false;
        return;
    }

    updateSettingValue(currentSettingItem.value, textualTagFilter ?? '');
    currentSettingItem.value = undefined;
    showFilterTransactionTagsDialog.value = false;
}

function open(value: DesktopOverviewWidgetLayout): Promise<DesktopOverviewWidgetLayout> {
    widget.value = cloneWidget(value);

    if (DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]) {
        const defaultSettings = DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]?.defaultSettings ?? {};
        widget.value.settings = { ...defaultSettings, ...widget.value.settings };
    }

    showState.value = true;

    const promises: Promise<unknown>[] = [];

    if (supportsSettings.value.some(setting => setting.settingType === 'accountSelect')) {
        promises.push(accountsStore.loadAllAccounts({ force: false }));
    }

    if (supportsSettings.value.some(setting => setting.settingType === 'categorySelect')) {
        promises.push(transactionCategoriesStore.loadAllCategories({ force: false }));
    }

    if (supportsSettings.value.some(setting => setting.settingType === 'tagSelect')) {
        promises.push(transactionTagsStore.loadAllTags({ force: false }));
    }

    if (promises.length > 0) {
        loading.value = true;

        Promise.all(promises).then(() => {
            loading.value = false;
        }).catch(error => {
            loading.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (!widget.value) {
        return;
    }

    resolveFunc?.(widget.value);
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
