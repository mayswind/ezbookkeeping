<template>
    <f7-popup push :opened="show" @popup:open="onPopupOpen" @popup:closed="onPopupClosed">
        <f7-page>
            <f7-navbar :outline="false">
                <f7-nav-left></f7-nav-left>
                <f7-nav-title :title="title" v-if="title"></f7-nav-title>
                <f7-nav-right>
                    <f7-link popup-close :text="tt('Done')"></f7-link>
                </f7-nav-right>
            </f7-navbar>
            <f7-searchbar ref="searchbar" custom-searchs
                          :value="filterContent"
                          :placeholder="filterPlaceholder"
                          :disable-button="false"
                          v-if="enableFilter"
                          @input="filterContent = $event.target.value">
            </f7-searchbar>
            <f7-block class="no-padding">
                <f7-list strong outline dividers>
                    <f7-list-item link="#" no-chevron
                                  :title="ti((titleField ? (item as Record<string, unknown>)[titleField] : item) as string, !!titleI18n)"
                                  :value="getItemValue(item, index, valueField, valueType)"
                                  :class="{ 'list-item-selected': isSelected(item, index) }"
                                  :key="getItemValue(item, index, keyField, valueType)"
                                  v-for="(item, index) in filteredItems"
                                  v-show="item && (!hiddenField || !(item as Record<string, unknown>)[hiddenField])"
                                  @click="onItemClicked(item, index)">
                        <template #content-start>
                            <f7-icon class="list-item-checked-icon" f7="checkmark_alt" :style="{ 'color': isSelected(item, index) ? '' : 'transparent' }"></f7-icon>
                        </template>
                        <template #media v-if="iconField">
                            <ItemIcon :icon-type="iconType" :icon-id="(item as Record<string, unknown>)[iconField]" :color="colorField ? (item as Record<string, unknown>)[colorField] : undefined"></ItemIcon>
                        </template>
                        <template #after>
                            <small v-if="afterField">{{ getItemAfterText(item) }}</small>
                        </template>
                    </f7-list-item>
                    <f7-list-item v-if="!filteredItems || !filteredItems.length"
                                  :title="filterNoItemsText"></f7-list-item>
                </f7-list>
            </f7-block>
        </f7-page>
    </f7-popup>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef } from 'vue';
import type { Searchbar } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';

import { type Framework7Dom, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: unknown;
    title?: string;
    valueType: string; // item or index
    keyField?: string; // for value type == item
    valueField?: string; // for value type == item
    titleField: string;
    titleI18n?: boolean;
    afterField?: string;
    afterI18n?: boolean;
    iconField?: string;
    iconType?: string;
    colorField?: string;
    hiddenField?: string;
    enableFilter?: boolean;
    filterPlaceholder?: string;
    filterNoItemsText?: string;
    items: unknown[];
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, ti } = useI18n();

const searchbar = useTemplateRef<Searchbar.Searchbar>('searchbar');

const currentValue = ref<unknown>(props.modelValue);
const filterContent = ref<string>('');

const filteredItems = computed<unknown[]>(() => {
    const finalItems: unknown[] = [];
    const items = props.items;
    const lowerCaseFilterContent = filterContent.value?.toLowerCase() ?? '';

    for (const item of items) {
        if (props.valueType === 'index') {
            if (!props.enableFilter || !lowerCaseFilterContent || String(item).toLowerCase().indexOf(lowerCaseFilterContent) >= 0) {
                finalItems.push(item);
                continue;
            }
        } else {
            const itemRecord = item as Record<string, unknown>;

            if (props.hiddenField && itemRecord[props.hiddenField]) {
                continue;
            }

            if (!props.enableFilter || !lowerCaseFilterContent) {
                finalItems.push(item);
                continue;
            }

            const title = ti(itemRecord[props.titleField] as string, !!props.titleI18n);

            if (title.toLowerCase().indexOf(lowerCaseFilterContent) >= 0) {
                finalItems.push(item);
                continue;
            }

            const afterText = getItemAfterText(item);

            if (afterText.toLowerCase().indexOf(lowerCaseFilterContent) >= 0) {
                finalItems.push(item);
                continue;
            }
        }
    }
    return finalItems;
});

function isSelected(item: unknown, index: number): boolean {
    if (props.valueType === 'index') {
        return currentValue.value === index;
    } else {
        if (props.valueField) {
            return currentValue.value === (item as Record<string, unknown>)[props.valueField];
        } else {
            return currentValue.value === item;
        }
    }
}

function getItemValue(item: unknown, index: number, fieldName: string | undefined, valueType: string): unknown {
    if (valueType === 'index') {
        return index;
    } else if (fieldName) {
        return (item as Record<string, unknown>)[fieldName];
    } else {
        return item;
    }
}

function getItemAfterText(item: unknown): string {
    if (props.valueType === 'index') {
        return '';
    } else if (props.afterField) {
        return ti((item as Record<string, unknown>)[props.afterField] as string, !!props.afterI18n);
    } else {
        return '';
    }
}

function close(): void {
    emit('update:show', false);
}

function onItemClicked(item: unknown, index: number): void {
    if (props.valueType === 'index') {
        currentValue.value = index;
    } else {
        if (props.valueField) {
            currentValue.value = (item as Record<string, unknown>)[props.valueField];
        } else {
            currentValue.value = item;
        }
    }

    emit('update:modelValue', currentValue.value);
    close();
}

function onPopupOpen(event: { $el: Framework7Dom }): void {
    currentValue.value = props.modelValue;
    scrollToSelectedItem(event.$el, '.page-content', 'li.list-item-selected');
}

function onPopupClosed(): void {
    close();
    filterContent.value = '';
    searchbar.value?.clear();
}
</script>
