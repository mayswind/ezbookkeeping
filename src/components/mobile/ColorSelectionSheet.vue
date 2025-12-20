<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar class="toolbar-with-swipe-handler">
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close icon-f7="xmark"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-block class="margin-vertical no-padding">
                <div class="grid padding-vertical-half padding-horizontal-half"
                     :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
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
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import type { ColorValue, ColorInfo } from '@/core/color.ts';
import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getColorsInRows } from '@/lib/color.ts';
import { scrollToSelectedItem } from '@/lib/ui/common.ts';
import { type Framework7Dom } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: ColorValue;
    show: boolean;
    columnCount?: number;
    allColorInfos: ColorValue[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: ColorValue): void;
    (e: 'update:show', value: boolean): void;
}>();

const currentValue = ref<ColorValue>(props.modelValue);
const itemPerRow = ref<number>(props.columnCount || 7);

const allColorRows = computed<ColorInfo[][]>(() => getColorsInRows(props.allColorInfos, itemPerRow.value));

function onColorClicked(colorInfo: ColorInfo): void {
    currentValue.value = colorInfo.color;
    emit('update:modelValue', currentValue.value);
}

function hasSelectedIcon(row: ColorInfo[]): boolean {
    return arrayContainsFieldValue(row, 'id', currentValue.value);
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    currentValue.value = props.modelValue;
    scrollToSelectedItem(event.$el[0], '.sheet-modal-inner', '.page-content', '.row-has-selected-item');
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>
