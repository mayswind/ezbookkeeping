<template>
    <v-select
        density="comfortable"
        item-title="icon"
        item-value="id"
        persistent-placeholder
        :disabled="disabled"
        :label="label"
        v-model="color"
        @update:menu="onMenuStateChanged"
    >
        <template #selection="{ item }">
            <v-label class="cursor-pointer" style="padding-top: 3px">
                <v-icon size="28" :icon="mdiSquareRounded" :color="getDisplayColor(item.raw)"/>
            </v-label>
        </template>

        <template #no-data>
            <div class="color-select-dropdown px-2" ref="dropdownMenu">
                <div class="color-item" :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                     :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                     :key="idx" v-for="(row, idx) in allColorRows">
                    <div class="text-center" :key="colorInfo.color" v-for="colorInfo in row">
                        <div class="cursor-pointer" @click="color = colorInfo.color">
                            <v-icon class="ma-2" size="28"
                                    :icon="mdiSquareRounded" :color="getDisplayColor(colorInfo.color)"
                                    v-if="!modelValue || modelValue !== colorInfo.color" />
                            <v-badge class="right-bottom-icon" color="primary"
                                     offset-x="8" offset-y="8"
                                     :location="`bottom ${textDirection === TextDirection.LTR ? 'right' : 'left'}`"
                                     :icon="mdiCheck"
                                     v-if="modelValue && modelValue === colorInfo.color">
                                <v-icon class="ma-2" size="28" :icon="mdiSquareRounded" :color="getDisplayColor(colorInfo.color)" />
                            </v-badge>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { TextDirection } from '@/core/text.ts';
import type { ColorValue, ColorInfo } from '@/core/color.ts';

import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getColorsInRows, getDisplayColor } from '@/lib/color.ts';
import { scrollToSelectedItem } from '@/lib/ui/desktop.ts';

import {
    mdiSquareRounded,
    mdiCheck
} from '@mdi/js';

const props = defineProps<{
    modelValue: ColorValue;
    disabled?: boolean;
    label?: string;
    columnCount?: number;
    allColorInfos: ColorValue[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: ColorValue): void;
}>();

const { getCurrentLanguageTextDirection } = useI18n();

const dropdownMenu = useTemplateRef<HTMLElement>('dropdownMenu');
const itemPerRow = ref<number>(props.columnCount || 7);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const allColorRows = computed<ColorInfo[][]>(() => getColorsInRows(props.allColorInfos, itemPerRow.value));

const color = computed<ColorValue>({
    get: () => props.modelValue,
    set: (value: ColorValue) => emit('update:modelValue', value)
});

function hasSelectedIcon(row: ColorInfo[]): boolean {
    return arrayContainsFieldValue(row, 'id', props.modelValue);
}

function onMenuStateChanged(state: boolean): void {
    if (state) {
        nextTick(() => {
            if (dropdownMenu.value && dropdownMenu.value.parentElement) {
                scrollToSelectedItem(dropdownMenu.value.parentElement, null, '.row-has-selected-item');
            }
        });
    }
}
</script>

<style>
.color-select-dropdown .color-item {
    display: grid;
}
</style>
