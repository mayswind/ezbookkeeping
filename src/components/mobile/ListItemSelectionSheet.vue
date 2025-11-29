<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="heightClass" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="right">
                <f7-link sheet-close :text="tt('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-list dividers class="no-margin-vertical">
                <f7-list-item link="#" no-chevron
                              :title="ti((titleField ? (item as Record<string, unknown>)[titleField] : item) as string, !!titleI18n)"
                              :value="getItemValue(item, index, valueField, valueType)"
                              :after="ti((afterField ? (item as Record<string, unknown>)[afterField] : '') as string, !!afterI18n)"
                              :class="{ 'list-item-selected': isSelected(item, index) }"
                              :key="getItemValue(item, index, keyField, valueType)"
                              v-for="(item, index) in items"
                              v-show="item && (!hiddenField || !(item as Record<string, unknown>)[hiddenField])"
                              @click="onItemClicked(item, index)">
                    <template #content-start>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" :style="{ 'color': isSelected(item, index) ? '' : 'transparent' }"></f7-icon>
                    </template>
                    <template #media v-if="iconField">
                        <ItemIcon :icon-type="iconType" :icon-id="(item as Record<string, unknown>)[iconField]" :color="colorField ? (item as Record<string, unknown>)[colorField] : undefined"></ItemIcon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { type Framework7Dom, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: unknown;
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
    items: unknown[];
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, ti } = useI18n();

const currentValue = ref<unknown>(props.modelValue);

const heightClass = computed<string>(() => {
    if (props.items.length > 10) {
        return 'list-item-selection-huge-sheet';
    } else if (props.items.length > 6) {
        return 'list-item-selection-large-sheet';
    } else {
        return '';
    }
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

function onSheetOpen(event: { $el: Framework7Dom }): void {
    currentValue.value = props.modelValue;
    scrollToSelectedItem(event.$el, '.page-content', 'li.list-item-selected');
}

function onSheetClosed(): void {
    close();
}
</script>

<style>
@media (min-height: 630px) {
    .list-item-selection-large-sheet {
        height: 310px;
    }

    .list-item-selection-huge-sheet {
        height: 400px;
    }
}
</style>
