<template>
    <f7-popup push swipe-to-close :opened="show" @popup:open="onPopupOpen" @popup:closed="onPopupClosed">
        <f7-page v-show="!showFilterAccountsPopup">
            <f7-navbar>
                <f7-nav-left>
                    <f7-link popup-close icon-f7="xmark"></f7-link>
                </f7-nav-left>
                <f7-nav-title :title="tt('Widget Settings')"></f7-nav-title>
                <f7-nav-right>
                    <f7-link icon-f7="checkmark_alt" @click="confirm"></f7-link>
                </f7-nav-right>
            </f7-navbar>

            <f7-list strong inset dividers class="settings-list margin-top-half" :class="{ 'disabled': loading }" >
                <template :key="setting.settingName" v-for="setting in supportsSettings">
                    <template v-if="setting.settingType === 'customSelect' && setting.multiple">
                        <f7-list-item group-title>
                            <small>{{ tt(setting.displayName) }}</small>
                        </f7-list-item>
                        <f7-list-item checkbox :key="`${setting.settingName}-${option.value}`"
                                      :title="tt(option.name)"
                                      :value="option.value"
                                      :checked="isMultipleValueSelected(setting, option.value)"
                                      :disabled="isLastSelectedMultipleValue(setting, option.value)"
                                      v-for="option in setting.selectValues"
                                      @change="updateMultipleValue(setting, option.value, $event.target.checked)"></f7-list-item>
                    </template>

                    <f7-list-item v-else-if="setting.settingType === 'switch'">
                        <template #after-title>
                            {{ tt(setting.displayName) }}
                        </template>
                        <template #after>
                            <f7-toggle :checked="getSettingValue(setting.settingName) as boolean"
                                       @toggle:change="updateSettingValue(setting.settingName, $event)"></f7-toggle>
                        </template>
                    </f7-list-item>

                    <f7-list-input readonly type="colorpicker" class="list-color-picker-input"
                                   :color-picker-params="{
                                       modules: ['sb-spectrum', 'hue-slider', 'hex'],
                                       hexLabel: false,
                                       hexValueEditable: true
                                   }"
                                   :value="{ hex: getDisplayColor(getSettingValue(setting.settingName) as string) }"
                                   @colorpicker:change="updateColorSettingValue(setting.settingName, $event)"
                                   v-else-if="setting.settingType === 'color'">
                        <template #inner-start>
                            <div class="item-actual-title">
                                <span>{{ tt(setting.displayName) }}</span>
                            </div>
                        </template>
                        <template #inner-end>
                            <div class="color-picker-display">
                                <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="getSettingValue(setting.settingName)"></ItemIcon>
                                <f7-icon class="color-picker-chevron" f7="chevron_up_chevron_down"></f7-icon>
                            </div>
                        </template>
                    </f7-list-input>

                    <f7-list-input type="text"
                                   :label="tt(setting.displayName)"
                                   :placeholder="setting.placeholder ? tt(setting.placeholder) : undefined"
                                   :value="getSettingValue(setting.settingName) as string"
                                   @input="updateSettingValue(setting.settingName, $event.target.value)"
                                   v-else-if="setting.settingType === 'textbox'"></f7-list-input>

                    <f7-list-item class="item-truncate-after-text"
                                  link="#" :disabled="isSettingDisabled(setting)"
                                  @click="openSettingSelection(setting, $event)"
                                  v-else>
                        <template #after-title>
                            <div class="item-actual-title">
                                <span>{{ tt(setting.displayName) }}</span>
                            </div>
                        </template>
                        <template #after>
                            <f7-preloader v-if="loading && (setting.settingType === 'accountSelect' || setting.settingType === 'categorySelect' || setting.settingType === 'tagSelect')" />
                            <div v-else>{{ getSingleSettingDisplayName(setting) }}</div>
                        </template>
                    </f7-list-item>
                </template>
            </f7-list>

            <f7-popover class="widget-setting-selection-popover" @popover:open="onPopoverOpen">
                <f7-list dividers>
                    <f7-list-item link="#" no-chevron popover-close
                                  :title="option.name"
                                  :class="{ 'list-item-selected': selectedSettingValue === option.value }"
                                  :key="option.value" v-for="option in selectedSettingOptions"
                                  @click="updateSelectedSettingValue(option.value)">
                        <template #after>
                            <f7-icon class="list-item-checked-icon" f7="checkmark_alt"
                                     v-if="selectedSettingValue === option.value"></f7-icon>
                        </template>
                    </f7-list-item>
                </f7-list>
            </f7-popover>
        </f7-page>

        <account-filter-settings-page disable-hidden-account
                                      v-model:custom-selected-account-ids="customSelectedAccountIds"
                                      @save="updateAccountValue"
                                      v-if="showFilterAccountsPopup" />
    </f7-popup>
</template>

<script setup lang="ts">
import AccountFilterSettingsPage from '@/views/mobile/settings/AccountFilterSettingsPage.vue';

import { ref, computed, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type Framework7Dom, openPopover, useI18nUIComponents } from '@/lib/ui/mobile.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import type { GenericNameValue } from '@/core/base.ts';

import {
    type OverviewWidgetSettingValue,
    type OverviewWidgetCustomSelectSettingItem,
    type OverviewWidgetSettingItem,
    type MobileOverviewWidgetLayout
} from '@/core/overview_layout.ts';
import { MOBILE_OVERVIEW_WIDGET_DEFINITIONS } from '@/consts/overview_layout.ts';

import { isDefined, isArray, isObjectEmpty, arrayItemToObjectField } from '@/lib/common.ts';
import { getDisplayColor } from '@/lib/color.ts';
import { isAllAccountsChecked } from '@/lib/account.ts';
import { cloneWidget } from '@/lib/overview_layout.ts';
import { scrollToSelectedItem } from '@/lib/ui/common.ts';

const props = defineProps<{
    modelValue: MobileOverviewWidgetLayout | null;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: MobileOverviewWidgetLayout): void;
    (e: 'update:show', value: boolean): void;
}>();

const {
    tt,
    formatNumberToLocalizedNumerals,
    getTablePageOptions
} = useI18n();
const { showToast } = useI18nUIComponents();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

const loading = ref<boolean>(false);
const widget = ref<MobileOverviewWidgetLayout | null>(null);
const currentSettingItem = ref<OverviewWidgetSettingItem | undefined>(undefined);
const showFilterAccountsPopup = ref<boolean>(false);
const customSelectedAccountIds = ref<string[]>([]);

const hasAnyAccount = computed<boolean>(() => accountsStore.allPlainAccounts.length > 0);
const hasAnyVisibleAccount = computed<boolean>(() => accountsStore.allVisibleAccountsCount > 0);
const hasAnyTransactionCategory = computed<boolean>(() => !isObjectEmpty(transactionCategoriesStore.allTransactionCategoriesMap));
const hasAnyAvailableTag = computed<boolean>(() => transactionTagsStore.allAvailableTagsCount > 0);

const supportsSettings = computed<OverviewWidgetSettingItem[]>(() => widget.value ? MOBILE_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]?.supportsSettings ?? [] : []);
const selectedSettingOptions = computed<GenericNameValue<string | number>[]>(() => currentSettingItem.value ? getSettingOptions(currentSettingItem.value) : []);
const selectedSettingValue = computed<OverviewWidgetSettingValue | undefined>(() => currentSettingItem.value ? getSettingValue(currentSettingItem.value.settingName) : undefined);

function isSettingDisabled(setting: OverviewWidgetSettingItem): boolean {
    if (setting.settingType === 'accountSelect' && setting.disableHiddenAccounts) {
        return !hasAnyVisibleAccount.value;
    } else if (setting.settingType === 'accountSelect' && !setting.disableHiddenAccounts) {
        return !hasAnyAccount.value;
    } else if (setting.settingType === 'categorySelect') {
        return !hasAnyTransactionCategory.value;
    } else if (setting.settingType === 'tagSelect') {
        return !hasAnyAvailableTag.value;
    }

    return false;
}

function getSettingOptions(setting: OverviewWidgetSettingItem): GenericNameValue<string | number>[] {
    if (setting.settingType === 'itemCountSelect') {
        return getTablePageOptions(setting.itemCountValues, undefined, false, true);
    } else if (setting.settingType === 'monthSelect') {
        return setting.monthValues.map(value => ({
            name: tt('format.misc.nMonths', { n: formatNumberToLocalizedNumerals(value) }),
            value: value
        }));
    } else if (setting.settingType === 'customSelect') {
        return setting.selectValues.map(option => ({ name: tt(option.name), value: option.value }));
    }

    return [];
}

function getSettingValue(settingName: string): OverviewWidgetSettingValue | undefined {
    return widget.value?.settings[settingName];
}

function updateSettingValue(settingName: string, value: OverviewWidgetSettingValue): void {
    if (widget.value) {
        widget.value.settings[settingName] = value;
    }
}

function updateColorSettingValue(settingName: string, value: { hex: string }): void {
    if (widget.value && value && value.hex) {
        widget.value.settings[settingName] = value.hex.replace(/^#/, '').substring(0, 6).toLowerCase();
    }
}

function getSingleSettingDisplayName(setting: OverviewWidgetSettingItem): string {
    const value = getSettingValue(setting.settingName);

    if (setting.settingType === 'accountSelect') {
        if (!isArray(value) || value.length < 1) {
            return tt('All');
        }

        return isAllAccountsChecked(accountsStore.allVisiblePlainAccounts, arrayItemToObjectField(value as string[], true)) ? tt('All') : tt('Partial');
    }

    return getSettingOptions(setting).find(option => option.value === value)?.name ?? '';
}

function isMultipleValueSelected(setting: OverviewWidgetCustomSelectSettingItem, value: string | number): boolean {
    const selectedValues = getSettingValue(setting.settingName);
    return isArray(selectedValues) && selectedValues.includes(value);
}

function isLastSelectedMultipleValue(setting: OverviewWidgetCustomSelectSettingItem, value: string | number): boolean {
    const selectedValues = getSettingValue(setting.settingName);
    return isArray(selectedValues) && selectedValues.length === 1 && selectedValues[0] === value;
}

function updateMultipleValue(setting: OverviewWidgetCustomSelectSettingItem, value: string | number, checked: boolean): void {
    const currentValue = getSettingValue(setting.settingName);
    const selectedValues = isArray(currentValue) ? [...currentValue] : [];
    const index = selectedValues.indexOf(value);

    if (checked && index < 0) {
        if (isDefined(setting.allValue) && value === setting.allValue) {
            updateSettingValue(setting.settingName, [value]);
            return;
        }

        selectedValues.push(value);
    } else if (!checked && index >= 0) {
        if (selectedValues.length <= 1) {
            return;
        }

        selectedValues.splice(index, 1);
    }

    if (isDefined(setting.allValue) && value !== setting.allValue) {
        const allValueIndex = selectedValues.indexOf(setting.allValue);

        if (allValueIndex >= 0) {
            selectedValues.splice(allValueIndex, 1);
        }
    }

    updateSettingValue(setting.settingName, selectedValues);
}

function openSettingSelection(setting: OverviewWidgetSettingItem, event: MouseEvent): void {
    if (loading.value) {
        return;
    }

    currentSettingItem.value = setting;

    if (setting.settingType === 'accountSelect') {
        const selectedAccountIds = getSettingValue(setting.settingName);
        customSelectedAccountIds.value = isArray(selectedAccountIds) ? [...selectedAccountIds] as string[] : [];
        showFilterAccountsPopup.value = true;
    } else {
        nextTick(() => {
            openPopover('.widget-setting-selection-popover', event.currentTarget as HTMLElement);
        });
    }
}

function updateSelectedSettingValue(value: string | number): void {
    if (currentSettingItem.value) {
        updateSettingValue(currentSettingItem.value.settingName, value);
    }
}

function updateAccountValue(): void {
    if (currentSettingItem.value && currentSettingItem.value.settingType === 'accountSelect') {
        updateSettingValue(currentSettingItem.value.settingName, customSelectedAccountIds.value);
    }

    currentSettingItem.value = undefined;
    customSelectedAccountIds.value = [];
    showFilterAccountsPopup.value = false;
}

function confirm(): void {
    if (widget.value) {
        emit('update:modelValue', widget.value);
    }

    close();
}

function close(): void {
    emit('update:show', false);
}

function onPopoverOpen(event: { $el: Framework7Dom }): void {
    scrollToSelectedItem(event.$el[0], '.popover-inner', '.popover-inner', 'li.list-item-selected');
}

function onPopupOpen(): void {
    if (!props.modelValue) {
        close();
        return;
    }

    widget.value = cloneWidget(props.modelValue);
    currentSettingItem.value = undefined;
    customSelectedAccountIds.value = [];

    if (MOBILE_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]) {
        const defaultSettings = MOBILE_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]?.defaultSettings ?? {};
        widget.value.settings = { ...defaultSettings, ...widget.value.settings };
    }

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
                showToast(error.message || error);
            }
        });
    }
}

function onPopupClosed(): void {
    widget.value = null;
    currentSettingItem.value = undefined;
    showFilterAccountsPopup.value = false;
    customSelectedAccountIds.value = [];
    close();
}
</script>

<style>
.widget-setting-selection-popover .popover-inner {
    max-height: 350px;
    overflow-y: auto;
}
</style>
