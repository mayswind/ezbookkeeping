<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              style="height: auto" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="tt('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <div class="grid grid-cols-2 grid-gap">
                <div>
                    <div class="primary-list-container">
                        <f7-list dividers class="primary-list no-margin-vertical">
                            <f7-list-item link="#" no-chevron
                                          :class="{ 'primary-list-item-selected': item === selectedPrimaryItem }"
                                          :value="primaryValueField ? (item as Record<string, unknown>)[primaryValueField] : item"
                                          :title="primaryTitleField ? ti((item as Record<string, unknown>)[primaryTitleField] as string, !!primaryTitleI18n) : ''"
                                          :header="primaryHeaderField ? ti((item as Record<string, unknown>)[primaryHeaderField] as string, !!primaryHeaderI18n) : ''"
                                          :footer="primaryFooterField ? ti((item as Record<string, unknown>)[primaryFooterField] as string, !!primaryFooterI18n) : ''"
                                          :key="primaryKeyField ? (item as Record<string, unknown>)[primaryKeyField] : item"
                                          v-for="item in items"
                                          v-show="item && (!primaryHiddenField || !(item as Record<string, unknown>)[primaryHiddenField])"
                                          @click="onPrimaryItemClicked(item)">
                                <template #media>
                                    <ItemIcon :icon-type="primaryIconType" :icon-id="primaryIconField ? (item as Record<string, unknown>)[primaryIconField] : undefined" :color="primaryColorField ? (item as Record<string, unknown>)[primaryColorField] : undefined"></ItemIcon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-showing" f7="chevron_right" v-if="item === selectedPrimaryItem"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </div>
                </div>
                <div>
                    <div class="secondary-list-container">
                        <f7-list dividers class="secondary-list no-margin-vertical" v-if="selectedPrimaryItem && primarySubItemsField && (selectedPrimaryItem as Record<string, unknown>)[primarySubItemsField]">
                            <f7-list-item link="#" no-chevron
                                          :class="{ 'secondary-list-item-selected': isSecondarySelected(subItem) }"
                                          :value="secondaryValueField ? subItem[secondaryValueField] : subItem"
                                          :title="secondaryTitleField ? ti(subItem[secondaryTitleField] as string, !!secondaryTitleI18n) : ''"
                                          :header="secondaryHeaderField ? ti(subItem[secondaryHeaderField] as string, !!secondaryHeaderI18n) : ''"
                                          :footer="secondaryFooterField ? ti(subItem[secondaryFooterField] as string, !!secondaryFooterI18n) : ''"
                                          :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                          v-for="subItem in (selectedPrimaryItem as Record<string, unknown>)[primarySubItemsField]"
                                          v-show="subItem && (!secondaryHiddenField || !subItem[secondaryHiddenField])"
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
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import {
    getItemByKeyValue,
    getPrimaryValueBySecondaryValue
} from '@/lib/common.ts';
import { type Framework7Dom, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: unknown;
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
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, ti } = useI18n();

const currentPrimaryValue = ref<unknown>(getCurrentPrimaryValueBySecondaryValue(props.modelValue));
const currentSecondaryValue = ref<unknown>(props.modelValue);

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

const selectedPrimaryItem = computed<unknown>(() => {
    if (props.primaryValueField) {
        return getItemByKeyValue(props.items as Record<string, unknown>[] | Record<string, Record<string, unknown>>, currentPrimaryValue.value, props.primaryValueField);
    } else {
        return currentPrimaryValue.value;
    }
});

function close() {
    emit('update:show', false);
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

    emit('update:modelValue', currentSecondaryValue.value);
    close();
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    currentPrimaryValue.value = getCurrentPrimaryValueBySecondaryValue(props.modelValue);
    currentSecondaryValue.value = props.modelValue;
    scrollToSelectedItem(event.$el, '.primary-list-container', 'li.primary-list-item-selected');
    scrollToSelectedItem(event.$el, '.secondary-list-container', 'li.secondary-list-item-selected');
}

function onSheetClosed() {
    close();
}
</script>

<style>
.primary-list-container, .secondary-list-container {
    height: 260px;
    overflow-y: auto;
}

.primary-list.list .item-inner {
    padding-right: 6px;
}

.secondary-list.list .item-content {
    padding-left: 0;
}
</style>
