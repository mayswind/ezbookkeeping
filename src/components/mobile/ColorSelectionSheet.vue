<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :style="sheetStyle" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar class="toolbar-with-swipe-handler">
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close icon-f7="xmark"></f7-link>
            </div>
            <f7-segmented strong round class="width-100">
                <f7-button :active="currentTab === 'system'" @click="showSystemColors">{{ tt('System Colors') }}</f7-button>
                <f7-button :active="currentTab === 'custom'" @click="showCustomColor">{{ tt('Custom Color') }}</f7-button>
            </f7-segmented>
        </f7-toolbar>
        <f7-page-content>
            <f7-block class="margin-vertical no-padding" v-if="currentTab === 'system'">
                <div class="grid padding-vertical-half padding-horizontal-half"
                     :class="{ 'row-has-selected-item': hasSelectedColor(row) }"
                     :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                     :key="idx" v-for="(row, idx) in allColorRows">
                    <div class="text-align-center" :key="colorInfo.color" v-for="colorInfo in row">
                        <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="colorInfo.color" @click="onColorClicked(colorInfo)">
                            <f7-badge color="default" class="right-bottom-icon" v-if="currentValue && currentValue === colorInfo.color">
                                <f7-icon f7="checkmark_alt"></f7-icon>
                            </f7-badge>
                        </ItemIcon>
                    </div>
                </div>
            </f7-block>
            <f7-block class="no-margin no-padding" v-if="currentTab === 'custom'">
                <div class="custom-color-picker" ref="customColorPickerContainer"></div>
            </f7-block>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef, nextTick, onBeforeUnmount } from 'vue';
import type { ColorPicker } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import type { ColorValue, ColorInfo } from '@/core/color.ts';
import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getColorsInRows } from '@/lib/color.ts';
import { scrollToSelectedItem } from '@/lib/ui/common.ts';
import { type Framework7Dom, createInlineColorPicker } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: ColorValue;
    show: boolean;
    columnCount?: number;
    allSystemColorInfos: ColorValue[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: ColorValue): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const customColorPickerContainer = useTemplateRef<HTMLElement>('customColorPickerContainer');

let customColorPicker: ColorPicker.ColorPicker | null = null;

const currentTab = ref<'system' | 'custom'>(isSystemColor(props.modelValue) ? 'system' : 'custom');
const currentValue = ref<ColorValue>(props.modelValue);
const itemPerRow = ref<number>(props.columnCount || 7);
const customSheetHeight = ref<number>(420);

const allColorRows = computed<ColorInfo[][]>(() => getColorsInRows(props.allSystemColorInfos, itemPerRow.value));

const sheetStyle = computed<Record<string, string> | undefined>(() => {
    const style: Record<string, string> = {};

    if (currentTab.value === 'custom' && customSheetHeight.value) {
        style['height'] = `${customSheetHeight.value}px`;
    }

    return style;
});

function onColorClicked(colorInfo: ColorInfo): void {
    currentValue.value = colorInfo.color;
    emit('update:modelValue', currentValue.value);
}

function isSystemColor(value: ColorValue): boolean {
    return props.allSystemColorInfos.includes(value);
}

function hasSelectedColor(row: ColorInfo[]): boolean {
    return arrayContainsFieldValue(row, 'color', currentValue.value);
}

function createCustomColorPicker(): void {
    nextTick(() => {
        if (!customColorPickerContainer.value || customColorPicker) {
            return;
        }

        customColorPicker = createInlineColorPicker(
            customColorPickerContainer.value,
            currentValue.value,
            (value: ColorValue) => {
                currentValue.value = value;
                emit('update:modelValue', currentValue.value);
            }
        );

        nextTick(() => {
            const pageContent = customColorPickerContainer.value?.closest('.page-content');

            if (pageContent) {
                customSheetHeight.value = Math.ceil(pageContent.scrollHeight);
            }
        });
    });
}

function destroyCustomColorPicker(): void {
    if (customColorPicker) {
        customColorPicker.destroy();
        customColorPicker = null;
    }

    if (customColorPickerContainer.value) {
        customColorPickerContainer.value.innerHTML = '';
    }
}

function showSystemColors(): void {
    currentTab.value = 'system';
    destroyCustomColorPicker();
}

function showCustomColor(): void {
    currentTab.value = 'custom';
    createCustomColorPicker();
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    currentValue.value = props.modelValue;
    currentTab.value = isSystemColor(currentValue.value) ? 'system' : 'custom';

    if (currentTab.value === 'custom') {
        createCustomColorPicker();
    } else {
        scrollToSelectedItem(event.$el[0], '.sheet-modal-inner', '.page-content', '.row-has-selected-item');
    }
}

function onSheetClosed(): void {
    destroyCustomColorPicker();
    emit('update:show', false);
}

onBeforeUnmount(() => {
    destroyCustomColorPicker();
});
</script>

<style>
.custom-color-picker {
    .color-picker {
        width: 100%;

        .color-picker-hex-value {
            width: 100%;
        }
    }
}
</style>
