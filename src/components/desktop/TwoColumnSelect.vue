<template>
    <v-select
        persistent-placeholder
        :density="density"
        :variant="variant"
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ 'content-class': 'two-column-select-menu' }"
        v-model="currentSecondaryValue"
        v-model:menu="menuState"
        @update:menu="onMenuStateChanged"
    >
        <template #selection>
            <div class="d-flex align-center text-truncate cursor-pointer">
                <span class="text-truncate" v-if="customSelectionPrimaryText">{{ customSelectionPrimaryText }}</span>
                <v-icon class="disabled" :icon="icons.chevronRight" size="23" v-if="customSelectionPrimaryText && customSelectionSecondaryText" />
                <span class="text-truncate" v-if="customSelectionPrimaryText && customSelectionSecondaryText">{{ customSelectionSecondaryText }}</span>
                <span class="text-truncate" v-if="!customSelectionPrimaryText && !selectedPrimaryItem && !selectedSecondaryItem">{{ noSelectionText }}</span>
                <span class="text-truncate" v-if="!customSelectionPrimaryText && showSelectionPrimaryText && selectedPrimaryItem">{{ selectionPrimaryItemText }}</span>
                <v-icon class="disabled" :icon="icons.chevronRight" size="23" v-if="!customSelectionPrimaryText && showSelectionPrimaryText && selectedPrimaryItem && selectedSecondaryItem" />
                <ItemIcon class="mr-2" icon-type="account" size="21.5px"
                          :icon-id="selectedSecondaryItem && secondaryIconField ? (selectedSecondaryItem as Record<string, unknown>)[secondaryIconField] : null"
                          :color="selectedSecondaryItem && secondaryColorField ? (selectedSecondaryItem as Record<string, unknown>)[secondaryColorField] : null"
                          v-if="!customSelectionPrimaryText && selectedSecondaryItem && showSelectionSecondaryIcon" />
                <span class="text-truncate" v-if="!customSelectionPrimaryText && selectedSecondaryItem">{{ selectionSecondaryItemText }}</span>
            </div>
        </template>

        <template #no-data>
            <div ref="dropdownMenu" class="two-column-list-container">
                <div class="primary-list-container">
                    <v-list :class="{ 'list-item-with-header': !!primaryHeaderField, 'list-item-with-footer': !!primaryFooterField }">
                        <v-list-item :class="{ 'primary-list-item-selected v-list-item--active text-primary': item === selectedPrimaryItem }"
                                     :key="primaryKeyField ? (item as Record<string, unknown>)[primaryKeyField] : item"
                                     v-for="item in items"
                                     v-show="item && (!primaryHiddenField || !(item as Record<string, unknown>)[primaryHiddenField])"
                                     @click="onPrimaryItemClicked(item)">
                            <template #prepend>
                                <ItemIcon class="mr-2" :icon-type="primaryIconType"
                                          :icon-id="primaryIconField ? (item as Record<string, unknown>)[primaryIconField] : undefined" :color="primaryColorField ? (item as Record<string, unknown>)[primaryColorField] : undefined"></ItemIcon>
                            </template>
                            <template #title>
                                <div class="list-item-header text-truncate" v-if="primaryHeaderField">{{ primaryHeaderField ? ti(item[primaryHeaderField] as string, !!primaryHeaderI18n) : '' }}</div>
                                <div class="text-truncate">{{ primaryTitleField ? ti(item[primaryTitleField] as string, !!primaryTitleI18n) : '' }}</div>
                                <div class="list-item-footer text-truncate" v-if="primaryFooterField">{{ primaryFooterField ? ti(item[primaryFooterField] as string, !!primaryFooterI18n) : '' }}</div>
                            </template>
                        </v-list-item>
                    </v-list>
                </div>
                <div class="secondary-list-container">
                    <v-list :class="{ 'list-item-with-header': !!secondaryHeaderField, 'list-item-with-footer': !!secondaryFooterField }"
                            v-if="selectedPrimaryItem && primarySubItemsField && (selectedPrimaryItem as Record<string, unknown>)[primarySubItemsField]">
                        <v-list-item :class="{ 'secondary-list-item-selected v-list-item--active text-primary': isSecondarySelected(subItem) }"
                                     :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                     v-for="subItem in (selectedPrimaryItem as Record<string, unknown>)[primarySubItemsField]"
                                     v-show="subItem && (!secondaryHiddenField || !subItem[secondaryHiddenField])"
                                     @click="onSecondaryItemClicked(subItem)">
                            <template #prepend>
                                <ItemIcon class="mr-2" :icon-type="secondaryIconType"
                                          :icon-id="secondaryIconField ? subItem[secondaryIconField] : undefined" :color="secondaryColorField ? subItem[secondaryColorField] : undefined"></ItemIcon>
                            </template>
                            <template #title>
                                <div class="list-item-header text-truncate" v-if="secondaryHeaderField">{{ secondaryHeaderField ? ti(subItem[secondaryHeaderField] as string, !!secondaryHeaderI18n) : '' }}</div>
                                <div class="text-truncate">{{ ti(secondaryTitleField ? subItem[secondaryTitleField] as string : '', !!secondaryTitleI18n) }}</div>
                                <div class="list-item-footer text-truncate" v-if="secondaryFooterField">{{ secondaryFooterField ? ti(subItem[secondaryFooterField] as string, !!secondaryFooterI18n) : '' }}</div>
                            </template>
                        </v-list-item>
                    </v-list>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import {
    getFirstVisibleItem,
    getItemByKeyValue,
    getNameByKeyValue,
    getPrimaryValueBySecondaryValue
} from '@/lib/common.ts';
import { scrollToSelectedItem } from '@/lib/ui/desktop.ts';

import {
    mdiChevronRight
} from '@mdi/js';

const props = defineProps<{
    modelValue: unknown;
    density?: string;
    variant?: string;
    disabled?: boolean;
    readonly?: boolean;
    label?: string;
    showSelectionPrimaryText?: boolean;
    showSelectionSecondaryIcon?: boolean;
    customSelectionPrimaryText?: string;
    customSelectionSecondaryText?: string;
    primaryKeyField?: string;
    primaryValueField?: string;
    primaryTitleField?: string;
    primaryTitleI18n?: boolean;
    primaryHeaderField?: string;
    primaryHeaderI18n?: boolean;
    primaryFooterField?: string;
    primaryFooterI18n?: boolean;
    primaryIconField?: string;
    primaryIconType?: string;
    primaryColorField?: string;
    primaryHiddenField?: string;
    primarySubItemsField: string;
    secondaryKeyField?: string;
    secondaryValueField?: string;
    secondaryTitleField?: string;
    secondaryTitleI18n?: boolean;
    secondaryHeaderField?: string;
    secondaryHeaderI18n?: boolean;
    secondaryFooterField?: string;
    secondaryFooterI18n?: boolean;
    secondaryIconField?: string;
    secondaryIconType?: string;
    secondaryColorField?: string;
    secondaryHiddenField?: string;
    items: unknown[];
    noItemText?: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void;
}>();

const { tt, ti } = useI18n();

const icons = {
    chevronRight: mdiChevronRight
};

const dropdownMenu = useTemplateRef<HTMLElement>('dropdownMenu');

const menuState = ref<boolean>(false);

const currentPrimaryValue = computed<unknown>({
    get: () => {
        return getCurrentPrimaryValueBySecondaryValue(props.modelValue);
    },
    set: (value) => {
        const primaryItem = getItemByKeyValue(props.items as Record<string, unknown>[] | Record<string, Record<string, unknown>>, value, props.primaryValueField as string);

        if (!primaryItem) {
            return;
        }

        const secondaryItem = getFirstVisibleItem(primaryItem[props.primarySubItemsField] as Record<string, unknown>[] | Record<string, Record<string, unknown>>, props.primaryHiddenField as string);

        if (secondaryItem) {
            if (props.secondaryValueField) {
                emit('update:modelValue', secondaryItem[props.secondaryValueField]);
            }
        }
    }
});

const currentSecondaryValue = computed<unknown>({
    get: () => {
        return props.modelValue;
    },
    set: (value) => {
        menuState.value = false;
        emit('update:modelValue', value);
    }
});

const selectedPrimaryItem = computed<unknown>(() => {
    if (props.primaryValueField) {
        return getItemByKeyValue(props.items as Record<string, unknown>[] | Record<string, Record<string, unknown>>, currentPrimaryValue.value, props.primaryValueField);
    } else {
        return currentPrimaryValue.value;
    }
});

const selectedSecondaryItem = computed<unknown>(() => {
    if (currentSecondaryValue.value && selectedPrimaryItem.value && (selectedPrimaryItem.value as Record<string, unknown>)[props.primarySubItemsField]) {
        return getItemByKeyValue((selectedPrimaryItem.value as Record<string, unknown>)[props.primarySubItemsField] as Record<string, unknown>[] | Record<string, Record<string, unknown>>, currentSecondaryValue.value, props.secondaryValueField as string);
    } else {
        return null;
    }
});

const noSelectionText = computed<string>(() => props.noItemText ? props.noItemText : tt('None'));

const selectionPrimaryItemText = computed<string>(() => {
    if (props.primaryValueField && props.primaryTitleField) {
        if (currentPrimaryValue.value) {
            return getNameByKeyValue(props.items as Record<string, unknown>[] | Record<string, Record<string, unknown>>, currentPrimaryValue.value, props.primaryValueField, props.primaryTitleField, noSelectionText.value) as string;
        } else {
            return noSelectionText.value;
        }
    } else {
        return currentPrimaryValue.value as string;
    }
});

const selectionSecondaryItemText = computed<string>(() => {
    if (props.secondaryValueField && props.secondaryTitleField) {
        if (currentSecondaryValue.value && selectedPrimaryItem.value && (selectedPrimaryItem.value as Record<string, unknown>)[props.primarySubItemsField]) {
            return getNameByKeyValue((selectedPrimaryItem.value as Record<string, unknown>)[props.primarySubItemsField] as Record<string, unknown>[] | Record<string, Record<string, unknown>>, currentSecondaryValue.value, props.secondaryValueField, props.secondaryTitleField, noSelectionText.value) as string;
        } else {
            return noSelectionText.value;
        }
    } else {
        return currentSecondaryValue.value as string;
    }
});

function getCurrentPrimaryValueBySecondaryValue(secondaryValue: unknown): unknown {
    return getPrimaryValueBySecondaryValue(props.items as Record<string, Record<string, unknown>[]>[] | Record<string, Record<string, Record<string, unknown>[]>>, props.primarySubItemsField, props.primaryValueField, props.primaryHiddenField, props.secondaryValueField, props.secondaryHiddenField, secondaryValue);
}

function isSecondarySelected(subItem: unknown): boolean {
    if (props.secondaryValueField) {
        return currentSecondaryValue.value === (subItem as Record<string, unknown>)[props.secondaryValueField];
    } else {
        return currentSecondaryValue.value === subItem;
    }
}

function onPrimaryItemClicked(item: unknown): void {
    if (props.primaryValueField) {
        currentPrimaryValue.value = (item as Record<string, unknown>)[props.primaryValueField];
    } else {
        currentPrimaryValue.value = item;
    }
}

function onSecondaryItemClicked(subItem: unknown): void {
    if (props.secondaryValueField) {
        currentSecondaryValue.value = (subItem as Record<string, unknown>)[props.secondaryValueField];
    } else {
        currentSecondaryValue.value = subItem;
    }
}

function onMenuStateChanged(state: boolean): void {
    if (state) {
        nextTick(() => {
            if (dropdownMenu.value && dropdownMenu.value.parentElement) {
                scrollToSelectedItem(dropdownMenu.value.parentElement, '.primary-list-container', '.primary-list-item-selected');
                scrollToSelectedItem(dropdownMenu.value.parentElement, '.secondary-list-container', '.secondary-list-item-selected');
            }
        });
    }
}
</script>

<style>
.two-column-select-menu {
    max-height: inherit !important;
}

.two-column-select-menu > .v-list {
    padding: 0;
}

.two-column-select-menu .two-column-list-container {
    width: 100%;
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
}

.two-column-select-menu .primary-list-container,
.two-column-select-menu .secondary-list-container {
    width: 100%;
    max-height: 310px;
    overflow-y: scroll;
}

.two-column-select-menu .list-item-with-header > .v-list-item,
.two-column-select-menu .list-item-with-footer > .v-list-item {
    min-height: 58px;
    padding-top: 6px;
    padding-bottom: 6px;
}

.two-column-select-menu .list-item-with-header.list-item-with-footer > .v-list-item {
    min-height: 78px;
    padding-top: 8px;
    padding-bottom: 8px;
}

.two-column-select-menu .list-item-header,
.two-column-select-menu .list-item-footer {
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity));
    font-size: 0.75rem;
    line-height: 1.2rem;
}
</style>
