<template>
    <v-select
        class="two-column-main-select"
        persistent-placeholder
        :density="density"
        :variant="variant"
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ contentClass: 'two-column-select-menu' }"
        v-model="currentSecondaryValue"
        v-model:menu="menuState"
        @update:menu="onMenuStateChanged"
    >
        <template #selection>
            <div class="d-flex align-center text-truncate cursor-pointer">
                <span class="text-truncate" v-if="customSelectionPrimaryText">{{ customSelectionPrimaryText }}</span>
                <v-icon class="icon-with-direction disabled" :icon="mdiChevronRight" size="23" v-if="customSelectionPrimaryText && customSelectionSecondaryText" />
                <span class="text-truncate" v-if="customSelectionPrimaryText && customSelectionSecondaryText">{{ customSelectionSecondaryText }}</span>
                <span class="text-truncate" v-if="!customSelectionPrimaryText && !selectedPrimaryItem && !selectedSecondaryItem">{{ noSelectionText }}</span>
                <span class="text-truncate" v-if="!customSelectionPrimaryText && showSelectionPrimaryText && selectedPrimaryItem">{{ selectionPrimaryItemText }}</span>
                <v-icon class="icon-with-direction disabled" :icon="mdiChevronRight" size="23" v-if="!customSelectionPrimaryText && showSelectionPrimaryText && selectedPrimaryItem && selectedSecondaryItem" />
                <ItemIcon class="me-2" size="21.5px"
                          :icon-type="secondaryIconType"
                          :icon-id="selectedSecondaryItem && secondaryIconField ? (selectedSecondaryItem as Record<string, unknown>)[secondaryIconField] : null"
                          :color="selectedSecondaryItem && secondaryColorField ? (selectedSecondaryItem as Record<string, unknown>)[secondaryColorField] : null"
                          v-if="!customSelectionPrimaryText && selectedSecondaryItem && showSelectionSecondaryIcon" />
                <span class="text-truncate" v-if="!customSelectionPrimaryText && selectedSecondaryItem">{{ selectionSecondaryItemText }}</span>
            </div>
        </template>

        <template #no-data>
            <div class="mx-2 mt-2" v-if="enableFilter">
                <v-text-field eager ref="filterInput" density="compact"
                              :prepend-inner-icon="mdiMagnify"
                              :placeholder="filterPlaceholder"
                              v-model="filterContent"
                              @click:control="onInputFocused(filterInput, true)"
                              @update:focused="onInputFocused(filterInput, $event)"></v-text-field>
            </div>
            <div class="mx-4 my-3" v-show="!filteredItems || !filteredItems.length">
                {{ filterNoItemsText }}
            </div>
            <div ref="dropdownMenu" class="two-column-list-container" v-show="filteredItems && filteredItems.length">
                <div class="primary-list-container">
                    <v-list :class="{ 'list-item-with-header': !!primaryHeaderField, 'list-item-with-footer': !!primaryFooterField }">
                        <v-list-item :class="{ 'primary-list-item-selected v-list-item--active text-primary': item === selectedPrimaryItem }"
                                     :key="primaryKeyField ? (item as Record<string, unknown>)[primaryKeyField] as string : JSON.stringify(item)"
                                     v-for="item in filteredItems"
                                     @click="onPrimaryItemClicked(item)">
                            <template #prepend>
                                <ItemIcon class="me-2" :icon-type="primaryIconType"
                                          :icon-id="primaryIconField ? (item as Record<string, unknown>)[primaryIconField] : undefined"
                                          :color="primaryColorField ? (item as Record<string, unknown>)[primaryColorField] : undefined"
                                          v-if="primaryIconField"></ItemIcon>
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
                                     :key="secondaryKeyField ? subItem[secondaryKeyField] as string : JSON.stringify(subItem)"
                                     v-for="subItem in filteredSubItems"
                                     @click="onSecondaryItemClicked(subItem)">
                            <template #prepend>
                                <ItemIcon class="me-2" :icon-type="secondaryIconType"
                                          :icon-id="secondaryIconField ? subItem[secondaryIconField] : undefined"
                                          :color="secondaryColorField ? subItem[secondaryColorField] : undefined"
                                          v-if="secondaryIconField"></ItemIcon>
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
import { VTextField } from 'vuetify/components/VTextField';

import { ref, computed, useTemplateRef, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonTwoColumnListItemSelectionProps, useTwoColumnListItemSelectionBase } from '@/components/base/TwoColumnListItemSelectionBase.ts';

import {
    getFirstVisibleItem,
    getItemByKeyValue,
    getNameByKeyValue
} from '@/lib/common.ts';
import { type ComponentDensity, type InputVariant, setChildInputFocus, scrollToSelectedItem } from '@/lib/ui/desktop.ts';

import {
    mdiChevronRight,
    mdiMagnify
} from '@mdi/js';

interface DesktopTwoColumnListItemSelectionProps extends CommonTwoColumnListItemSelectionProps {
    density?: ComponentDensity;
    variant?: InputVariant;
    disabled?: boolean;
    readonly?: boolean;
    label?: string;
    showSelectionPrimaryText?: boolean;
    showSelectionSecondaryIcon?: boolean;
    customSelectionPrimaryText?: string;
    customSelectionSecondaryText?: string;
    noItemText?: string;
    autoUpdateMenuPosition?: boolean;
}

const props = defineProps<DesktopTwoColumnListItemSelectionProps>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void;
}>();

const { tt, ti } = useI18n();

const {
    filterContent,
    filteredItems,
    getFilteredSubItems,
    getCurrentPrimaryValueBySecondaryValue,
    isSecondaryValueSelected,
    getSelectedPrimaryItem,
    getSelectedSecondaryItem,
    updateCurrentPrimaryValue,
    updateCurrentSecondaryValue
} = useTwoColumnListItemSelectionBase(props);

const filterInput = useTemplateRef<VTextField>('filterInput');
const dropdownMenu = useTemplateRef<HTMLElement>('dropdownMenu');

const menuState = ref<boolean>(false);

const filteredSubItems = computed<Record<string, unknown>[]>(() => getFilteredSubItems(selectedPrimaryItem.value));

const currentPrimaryValue = computed<unknown>({
    get: () => {
        return getCurrentPrimaryValueBySecondaryValue(props.modelValue);
    },
    set: (value) => {
        const primaryItem = getItemByKeyValue(filteredItems.value, value, props.primaryValueField as string);

        if (!primaryItem) {
            return;
        }

        const secondaryItem = getFirstVisibleItem(getFilteredSubItems(primaryItem), props.primaryHiddenField as string);

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

const selectedPrimaryItem = computed<unknown>(() => getSelectedPrimaryItem(currentPrimaryValue.value));
const selectedSecondaryItem = computed<unknown>(() => getSelectedSecondaryItem(currentSecondaryValue.value, selectedPrimaryItem.value));

const noSelectionText = computed<string>(() => props.noItemText ? props.noItemText : tt('None'));

const selectionPrimaryItemText = computed<string>(() => {
    if (props.primaryValueField && props.primaryTitleField) {
        if (currentPrimaryValue.value) {
            return getNameByKeyValue(props.items as Record<string, string>[], currentPrimaryValue.value, props.primaryValueField, props.primaryTitleField, noSelectionText.value) as string;
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
            return getNameByKeyValue((selectedPrimaryItem.value as Record<string, unknown>)[props.primarySubItemsField] as Record<string, string>[], currentSecondaryValue.value, props.secondaryValueField, props.secondaryTitleField, noSelectionText.value) as string;
        } else {
            return noSelectionText.value;
        }
    } else {
        return currentSecondaryValue.value as string;
    }
});

function isSecondarySelected(subItem: unknown): boolean {
    return isSecondaryValueSelected(currentSecondaryValue.value, subItem);
}

function onPrimaryItemClicked(item: unknown): void {
    updateCurrentPrimaryValue(currentPrimaryValue, item);

    if (props.autoUpdateMenuPosition) {
        nextTick(() => {
            const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
            const mainSelectRect = document.querySelector('.two-column-main-select')?.getBoundingClientRect();
            const selectMenu = document.querySelector('.two-column-select-menu') as (HTMLElement | null);
            const selectMenuRect = selectMenu?.getBoundingClientRect();

            if (mainSelectRect && selectMenu && selectMenuRect) {
                const newTop = scrollTop + mainSelectRect.top + mainSelectRect.height + 0.5;

                if (newTop + selectMenuRect.height < document.documentElement.scrollHeight) {
                    selectMenu.style.top = newTop + 'px';
                }
            }
        });
    }
}

function onSecondaryItemClicked(subItem: unknown): void {
    updateCurrentSecondaryValue(currentSecondaryValue, subItem);
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

function onInputFocused(input: VTextField | null | undefined, focused: boolean): void {
    if (input && focused) {
        nextTick(() => {
            setChildInputFocus(input?.$el, 'input');
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
