<template>
    <v-select
        class="color-select"
        density="comfortable"
        item-title="icon"
        item-value="id"
        persistent-placeholder
        :disabled="disabled"
        :label="label"
        :menu-props="{ contentClass: 'color-select-menu' }"
        v-model="color"
        @update:menu="onMenuStateChanged"
    >
        <template #selection="{ internalItem }">
            <v-label class="cursor-pointer" style="padding-top: 3px">
                <v-icon size="28" :icon="mdiSquareRounded" :color="getDisplayColor(internalItem.raw)"/>
            </v-label>
        </template>

        <template #no-data>
            <div ref="dropdownMenu" class="color-select-dropdown px-2" :class="{ 'color-select-dropdown-custom': currentTab === 'custom' }">
                <div class="color-select-tabs-container">
                    <v-tabs grow density="compact" v-model="currentTab">
                        <v-tab value="system">{{ tt('System Colors') }}</v-tab>
                        <v-tab value="custom">{{ tt('Custom Color') }}</v-tab>
                    </v-tabs>
                </div>
                <div v-if="currentTab === 'system'">
                    <div class="color-item" :class="{ 'row-has-selected-item': hasSelectedColor(row) }"
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
                <div class="d-flex justify-center" v-if="currentTab === 'custom'">
                    <v-color-picker hide-input-labels elevation="0" :width="Math.max(320, dropdownMenuWidth - 16)"
                                    hide-alpha mode="hex" :modes="['hex']" v-model="customColor" />
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
import { scrollToSelectedItem } from '@/lib/ui/common.ts';

import {
    mdiSquareRounded,
    mdiCheck
} from '@mdi/js';

const props = defineProps<{
    modelValue: ColorValue;
    disabled?: boolean;
    label?: string;
    columnCount?: number;
    allSystemColorInfos: ColorValue[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: ColorValue): void;
}>();

const { tt, getCurrentLanguageTextDirection } = useI18n();

const dropdownMenu = useTemplateRef<HTMLElement>('dropdownMenu');

const currentTab = ref<'system' | 'custom'>(isSystemColor(props.modelValue) ? 'system' : 'custom');
const itemPerRow = ref<number>(props.columnCount || 7);
const dropdownMenuWidth = ref<number>(320);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const allColorRows = computed<ColorInfo[][]>(() => getColorsInRows(props.allSystemColorInfos, itemPerRow.value));

const color = computed<ColorValue>({
    get: () => props.modelValue,
    set: (value: ColorValue) => emit('update:modelValue', value)
});

const customColor = computed<string>({
    get: () => `#${props.modelValue}`,
    set: (value: string) => emit('update:modelValue', value.replace(/^#/, '').substring(0, 6).toLowerCase())
});

function isSystemColor(value: ColorValue): boolean {
    return props.allSystemColorInfos.includes(value);
}

function hasSelectedColor(row: ColorInfo[]): boolean {
    return arrayContainsFieldValue(row, 'color', props.modelValue);
}

function onMenuStateChanged(state: boolean): void {
    if (state) {
        currentTab.value = isSystemColor(props.modelValue) ? 'system' : 'custom';

        nextTick(() => {
            if (dropdownMenu.value && dropdownMenu.value.parentElement) {
                scrollToSelectedItem(dropdownMenu.value.parentElement, null, null, '.row-has-selected-item');
            }

            window.requestAnimationFrame(() => {
                if (dropdownMenu.value) {
                    dropdownMenuWidth.value = Math.floor(dropdownMenu.value.clientWidth);
                }
            });
        });
    }
}
</script>

<style>
.color-select:not(.v-input--disabled) .v-field__input,
.color-select:not(.v-input--disabled) .v-label {
    opacity: 1;
}

.color-select-menu {
    .color-select-dropdown {
        .color-select-tabs-container {
            position: sticky;
            top: -8px;
            z-index: 1;
            margin: -8px -8px 8px;
            padding: 8px 8px 0;
            background-color: rgb(var(--v-theme-surface));
        }

        .color-item {
            display: grid;
        }

        .v-color-picker__controls {
            padding: 4px;

            .v-color-picker-edit {
                margin-top: 4px;

                > .v-color-picker-edit__input {
                    > input {
                        margin-bottom: 0;
                    }
                }
            }
        }
    }
}
</style>
