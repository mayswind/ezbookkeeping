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
    </v-dialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useSettingsStore } from '@/stores/setting.ts';

import type { GenericNameValue, NameNumeralValue } from '@/core/base.ts';
import {
    type OverviewWidgetSettingValue,
    type OverviewWidgetCustomSelectSettingItem,
    type OverviewWidgetSettingItem,
    type DesktopOverviewWidgetLayout
} from '@/core/overview_layout.ts';

import { DESKTOP_OVERVIEW_WIDGET_DEFINITIONS } from '@/consts/overview_layout.ts';

import { isDefined, isArray } from '@/lib/common.ts';
import { cloneWidget } from '@/lib/overview_layout.ts';

const {
    tt,
    getAllAccountCategories,
    formatNumberToLocalizedNumerals,
    getTablePageOptions
} = useI18n();

const settingsStore = useSettingsStore();

let resolveFunc: ((widget: DesktopOverviewWidgetLayout) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const widget = ref<DesktopOverviewWidgetLayout | null>(null);

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
    if (setting.selectValueSource === 'accountCategories') {
        return [
            { name: tt('All'), value: 0 },
            ...getAllAccountCategories(settingsStore.appSettings.accountCategoryOrders).map(category => ({ name: category.displayName, value: category.type }))
        ];
    }

    return setting.selectValues.map(item => ({ name: tt(item.name), value: item.value }));
}

function getSettingValue(settingName: string): OverviewWidgetSettingValue | undefined {
    return widget.value?.settings[settingName];
}

function updateSettingValue(setting: OverviewWidgetSettingItem, value: OverviewWidgetSettingValue | null): void {
    if (!widget.value || value === null) {
        return;
    }

    if (setting.settingType === 'customSelect' && setting.multiple && isDefined(setting.allValue) && isArray(value)) {
        const previousValue = widget.value.settings[setting.settingName];
        const previousValues = isArray(previousValue) ? previousValue : [];
        const allValue = setting.allValue as number;

        if (value.includes(allValue) && !previousValues.includes(allValue)) {
            widget.value.settings[setting.settingName] = [allValue];
        } else {
            const selectedValues = value.filter(item => item !== allValue);
            widget.value.settings[setting.settingName] = selectedValues.length ? selectedValues : [allValue];
        }
    } else {
        widget.value.settings[setting.settingName] = value;
    }
}

function open(value: DesktopOverviewWidgetLayout): Promise<DesktopOverviewWidgetLayout> {
    widget.value = cloneWidget(value);

    if (DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]) {
        const defaultSettings = DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[widget.value.type]?.defaultSettings ?? {};
        widget.value.settings = { ...defaultSettings, ...widget.value.settings };
    }

    showState.value = true;

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
