<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="heightClass" :opened="show"
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
                     :key="idx" v-for="(row, idx) in allIconRows">
                    <div class="text-align-center" :key="iconInfo.id" v-for="iconInfo in row">
                        <ItemIcon icon-type="fixed" :icon-id="iconInfo.icon" :color="color" @click="onIconClicked(iconInfo)">
                            <f7-badge color="default" class="right-bottom-icon" v-if="currentValue && currentValue === iconInfo.id">
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

import type { IconInfo, IconInfoWithId } from '@/core/icon.ts';
import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getIconsInRows } from '@/lib/icon.ts';
import { type Framework7Dom, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: string;
    show: boolean;
    columnCount?: number;
    color: string;
    allIconInfos: Record<string, IconInfo>;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
    (e: 'update:show', value: boolean): void;
}>();

const currentValue = ref<string>(props.modelValue);
const itemPerRow = ref<number>(props.columnCount || 7);

const allIconRows = computed<IconInfoWithId[][]>(() => getIconsInRows(props.allIconInfos, itemPerRow.value));

const heightClass = computed<string>(() => {
    if (allIconRows.value.length > 10) {
        return 'icon-selection-huge-sheet';
    } else if (allIconRows.value.length > 6) {
        return 'icon-selection-large-sheet';
    } else {
        return '';
    }
});

function onIconClicked(iconInfo: IconInfoWithId): void {
    currentValue.value = iconInfo.id;
    emit('update:modelValue', currentValue.value);
}

function hasSelectedIcon(row: IconInfoWithId[]): boolean {
    return arrayContainsFieldValue(row, 'id', currentValue.value);
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    scrollToSelectedItem(event.$el, '.page-content', '.row-has-selected-item');
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>

<style>
@media (min-height: 630px) {
    .icon-selection-large-sheet {
        height: 310px;
    }

    .icon-selection-huge-sheet {
        height: 400px;
    }
}
</style>
