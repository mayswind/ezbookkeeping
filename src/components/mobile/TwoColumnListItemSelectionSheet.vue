<template>
    <f7-sheet ref="sheet" swipe-to-close swipe-handler=".swipe-handler"
              style="height: auto" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="tt('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-searchbar ref="searchbar" custom-searchs
                      :value="filterContent"
                      :placeholder="filterPlaceholder"
                      :disable-button="false"
                      v-if="enableFilter"
                      @input="filterContent = $event.target.value"
                      @focus="onSearchBarFocus">
        </f7-searchbar>
        <f7-page-content class="no-padding-top">
            <div class="grid grid-gap" :class="{ 'grid-cols-2': filteredItems && filteredItems.length }">
                <div>
                    <div class="primary-list-container">
                        <f7-list dividers class="primary-list no-margin-vertical">
                            <f7-list-item link="#" no-chevron
                                          :class="{ 'primary-list-item-selected': item === selectedPrimaryItem }"
                                          :value="primaryValueField ? item[primaryValueField] : item"
                                          :title="primaryTitleField ? ti(item[primaryTitleField] as string, !!primaryTitleI18n) : ''"
                                          :header="primaryHeaderField ? ti(item[primaryHeaderField] as string, !!primaryHeaderI18n) : ''"
                                          :footer="primaryFooterField ? ti(item[primaryFooterField] as string, !!primaryFooterI18n) : ''"
                                          :key="primaryKeyField ? item[primaryKeyField] : item"
                                          v-for="item in filteredItems"
                                          @click="onPrimaryItemClicked(item)">
                                <template #media>
                                    <ItemIcon :icon-type="primaryIconType" :icon-id="primaryIconField ? item[primaryIconField] : undefined" :color="primaryColorField ? item[primaryColorField] : undefined"></ItemIcon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-showing" f7="chevron_right" v-if="item === selectedPrimaryItem"></f7-icon>
                                </template>
                            </f7-list-item>
                            <f7-list-item v-if="!filteredItems || !filteredItems.length"
                                            :title="filterNoItemsText"></f7-list-item>
                        </f7-list>
                    </div>
                </div>
                <div v-show="filteredItems && filteredItems.length">
                    <div class="secondary-list-container">
                        <f7-list dividers class="secondary-list no-margin-vertical" v-if="selectedPrimaryItem && primarySubItemsField && (selectedPrimaryItem as Record<string, unknown>)[primarySubItemsField]">
                            <f7-list-item link="#" no-chevron
                                          :class="{ 'secondary-list-item-selected': isSecondarySelected(subItem) }"
                                          :value="secondaryValueField ? subItem[secondaryValueField] : subItem"
                                          :title="secondaryTitleField ? ti(subItem[secondaryTitleField] as string, !!secondaryTitleI18n) : ''"
                                          :header="secondaryHeaderField ? ti(subItem[secondaryHeaderField] as string, !!secondaryHeaderI18n) : ''"
                                          :footer="secondaryFooterField ? ti(subItem[secondaryFooterField] as string, !!secondaryFooterI18n) : ''"
                                          :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                          v-for="subItem in filteredSubItems"
                                          @click="onSecondaryItemClicked(subItem)">
                                <template #media>
                                    <ItemIcon :icon-type="secondaryIconType" :icon-id="secondaryIconField ? subItem[secondaryIconField] : undefined" :color="secondaryColorField ? subItem[secondaryColorField] : undefined"></ItemIcon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="isSecondarySelected(subItem)"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </div>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef } from 'vue';
import type { Sheet, Searchbar } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonTwoColumnListItemSelectionProps, useTwoColumnListItemSelectionBase } from '@/components/base/TwoColumnListItemSelectionBase.ts';

import { type Framework7Dom, scrollToSelectedItem, scrollSheetToTop } from '@/lib/ui/mobile.ts';

interface MobileTwoColumnListItemSelectionProps extends CommonTwoColumnListItemSelectionProps {
    show: boolean;
}

const props = defineProps<MobileTwoColumnListItemSelectionProps>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, ti } = useI18n();

const {
    filterContent,
    filteredItems,
    getFilteredSubItems,
    getCurrentPrimaryValueBySecondaryValue,
    isSecondaryValueSelected,
    getSelectedPrimaryItem,
    updateCurrentPrimaryValue,
    updateCurrentSecondaryValue
} = useTwoColumnListItemSelectionBase(props);

const sheet = useTemplateRef<Sheet.Sheet>('sheet');
const searchbar = useTemplateRef<Searchbar.Searchbar>('searchbar');

const currentPrimaryValue = ref<unknown>(getCurrentPrimaryValueBySecondaryValue(props.modelValue));
const currentSecondaryValue = ref<unknown>(props.modelValue);

const filteredSubItems = computed<Record<string, unknown>[]>(() => getFilteredSubItems(selectedPrimaryItem.value));
const selectedPrimaryItem = computed<unknown>(() => getSelectedPrimaryItem(currentPrimaryValue.value));

function isSecondarySelected(subItem: unknown): boolean {
    return isSecondaryValueSelected(currentSecondaryValue.value, subItem);
}

function close(): void {
    emit('update:show', false);
}

function onPrimaryItemClicked(item: unknown): void {
    updateCurrentPrimaryValue(currentPrimaryValue, item);
}

function onSecondaryItemClicked(subItem: unknown): void {
    updateCurrentSecondaryValue(currentSecondaryValue, subItem);
    emit('update:modelValue', currentSecondaryValue.value);
    close();
}

function onSearchBarFocus(): void {
    scrollSheetToTop(sheet.value?.$el as HTMLElement, window.innerHeight); // $el is not Framework7 Dom
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    currentPrimaryValue.value = getCurrentPrimaryValueBySecondaryValue(props.modelValue);
    currentSecondaryValue.value = props.modelValue;
    scrollToSelectedItem(event.$el, '.primary-list-container', 'li.primary-list-item-selected');
    scrollToSelectedItem(event.$el, '.secondary-list-container', 'li.secondary-list-item-selected');
}

function onSheetClosed(): void {
    close();
    filterContent.value = '';
    searchbar.value?.clear();
}
</script>

<style>
.primary-list-container, .secondary-list-container {
    height: 260px;
    overflow-y: auto;
}

@media (max-height: 629px) {
    .primary-list-container, .secondary-list-container {
        height: 240px;
    }
}

.primary-list.list .item-inner {
    padding-right: 6px;
}

.secondary-list.list .item-content {
    padding-left: 0;
}
</style>
