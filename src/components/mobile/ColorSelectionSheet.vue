<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
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
import { type Ref, ref, computed } from 'vue';

import type { ColorValue, ColorInfo } from '@/core/color.ts';
import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getColorsInRows } from '@/lib/color.ts';
import { scrollToSelectedItem } from '@/lib/ui/mobile.js';

const props = defineProps<{
    modelValue: ColorValue;
    show: boolean;
    columnCount?: number;
    allColorInfos: ColorValue[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: ColorValue): void
    (e: 'update:show', value: boolean): void
}>();

const currentValue: Ref<ColorValue> = ref(props.modelValue);
const itemPerRow: Ref<number> = ref(props.columnCount || 7);

const allColorRows = computed<ColorInfo[][]>(() => {
    return getColorsInRows(props.allColorInfos, itemPerRow.value);
});

function onColorClicked(colorInfo: ColorInfo) {
    currentValue.value = colorInfo.color;
    emit('update:modelValue', currentValue.value);
}

function hasSelectedIcon(row: ColorInfo[]) {
    return arrayContainsFieldValue(row, 'id', currentValue.value);
}

function onSheetOpen(event: { $el: HTMLElement }) {
    currentValue.value = props.modelValue;
    scrollToSelectedItem(event.$el, '.page-content', '.row-has-selected-item');
}

function onSheetClosed() {
    emit('update:show', false);
}
</script>
