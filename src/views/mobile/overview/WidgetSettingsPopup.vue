<template>
    <f7-popup push swipe-to-close :opened="show" @popup:open="onPopupOpen" @popup:closed="onPopupClosed">
        <f7-page>
            <f7-navbar>
                <f7-nav-left>
                    <f7-link popup-close icon-f7="xmark"></f7-link>
                </f7-nav-left>
                <f7-nav-title :title="tt('Widget Settings')"></f7-nav-title>
                <f7-nav-right>
                    <f7-link icon-f7="checkmark_alt" @click="confirm"></f7-link>
                </f7-nav-right>
            </f7-navbar>

            <f7-list strong inset dividers class="settings-list margin-top-half">
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
                                  link="#" popover-open=".widget-setting-selection-popover"
                                  @click="selectedSetting = setting"
                                  v-else>
                        <template #after-title>
                            <div class="item-actual-title">
                                <span>{{ tt(setting.displayName) }}</span>
                            </div>
                        </template>
                        <template #after>
                            <div>{{ getSingleSettingDisplayName(setting) }}</div>
                        </template>
                    </f7-list-item>
                </template>
            </f7-list>

            <f7-popover class="widget-setting-selection-popover">
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
    </f7-popup>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import type { GenericNameValue } from '@/core/base.ts';

import {
    type OverviewWidgetSettingValue,
    type OverviewWidgetCustomSelectSettingItem,
    type OverviewWidgetSettingItem,
    type MobileOverviewWidgetLayout
} from '@/core/overview_layout.ts';
import { MOBILE_OVERVIEW_WIDGET_DEFINITIONS } from '@/consts/overview_layout.ts';

import { isDefined, isArray } from '@/lib/common.ts';
import { getDisplayColor } from '@/lib/color.ts';
import { cloneWidget } from '@/lib/overview_layout.ts';

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
    getAllAccountCategories,
    formatNumberToLocalizedNumerals,
    getTablePageOptions
} = useI18n();

const settingsStore = useSettingsStore();

const widget = ref<MobileOverviewWidgetLayout | null>(null);
const selectedSetting = ref<OverviewWidgetSettingItem | null>(null);

const supportsSettings = computed<OverviewWidgetSettingItem[]>(() => widget.value ? MOBILE_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]?.supportsSettings ?? [] : []);
const selectedSettingOptions = computed<GenericNameValue<string | number>[]>(() => selectedSetting.value ? getSettingOptions(selectedSetting.value) : []);
const selectedSettingValue = computed<OverviewWidgetSettingValue | undefined>(() => selectedSetting.value ? getSettingValue(selectedSetting.value.settingName) : undefined);

function getSettingOptions(setting: OverviewWidgetSettingItem): GenericNameValue<string | number>[] {
    if (setting.settingType === 'itemCountSelect') {
        return getTablePageOptions(setting.itemCountValues, undefined, false, true);
    } else if (setting.settingType === 'monthSelect') {
        return setting.monthValues.map(value => ({
            name: tt('format.misc.nMonths', { n: formatNumberToLocalizedNumerals(value) }),
            value: value
        }));
    } else if (setting.settingType === 'customSelect') {
        if (setting.selectValueSource === 'accountCategories') {
            return [
                { name: tt('All'), value: 0 },
                ...getAllAccountCategories(settingsStore.appSettings.accountCategoryOrders).map(category => ({ name: category.displayName, value: category.type }))
            ];
        }

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

function updateSelectedSettingValue(value: string | number): void {
    if (selectedSetting.value) {
        updateSettingValue(selectedSetting.value.settingName, value);
    }
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

function onPopupOpen(): void {
    if (!props.modelValue) {
        close();
        return;
    }

    widget.value = cloneWidget(props.modelValue);
    selectedSetting.value = null;

    if (MOBILE_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]) {
        const defaultSettings = MOBILE_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]?.defaultSettings ?? {};
        widget.value.settings = { ...defaultSettings, ...widget.value.settings };
    }
}

function onPopupClosed(): void {
    widget.value = null;
    selectedSetting.value = null;
    close();
}
</script>

<style>
.widget-setting-selection-popover .popover-inner {
    max-height: 350px;
    overflow-y: auto;
}
</style>
