<template>
    <v-select
        density="comfortable"
        item-title="icon"
        item-value="id"
        persistent-placeholder
        :disabled="disabled"
        :label="label"
        v-model="icon"
        @update:menu="onMenuStateChanged"
    >
        <template #selection>
            <v-label class="cursor-pointer">
                <ItemIcon :icon-type="iconType" :icon-id="icon" :color="color" />
            </v-label>
        </template>

        <template #no-data>
            <div class="icon-select-dropdown" ref="dropdownMenu">
                <div class="icon-item" :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                     :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                     :key="idx" v-for="(row, idx) in allIconRows">
                    <div class="text-center" :key="iconInfo.id" v-for="iconInfo in row">
                        <div class="cursor-pointer" @click="icon = iconInfo.id">
                            <ItemIcon class="ma-2" icon-type="fixed" :icon-id="iconInfo.icon" :color="color" v-if="!modelValue || modelValue !== iconInfo.id" />
                            <v-badge class="right-bottom-icon" color="primary"
                                     location="bottom right" offset-x="8" offset-y="10" :icon="icons.checked"
                                     v-if="modelValue && modelValue === iconInfo.id">
                                <ItemIcon class="ma-2" icon-type="fixed" :icon-id="iconInfo.icon" :color="color" />
                            </v-badge>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { type Ref, ref, computed, useTemplateRef, nextTick } from 'vue';

import type { ColorValue } from '@/core/color.ts';
import type {IconInfo, IconInfoWithId} from '@/core/icon.ts';
import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getIconsInRows } from '@/lib/icon.ts';
import { scrollToSelectedItem } from '@/lib/ui/desktop.ts';

import {
    mdiCheck
} from '@mdi/js';

const props = defineProps<{
    modelValue: string;
    disabled?: boolean;
    label?: string;
    iconType: string;
    color: ColorValue;
    columnCount?: number;
    allIconInfos: Record<string, IconInfo>;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void
}>();

const icons = {
    checked: mdiCheck
};

const dropdownMenu: Ref<HTMLElement | null> = useTemplateRef('dropdownMenu');
const itemPerRow: Ref<number> = ref(props.columnCount || 7);

const allIconRows = computed<IconInfoWithId[][]>(() => {
    return getIconsInRows(props.allIconInfos, itemPerRow.value);
});

const icon = computed({
    get: () => props.modelValue,
    set: (value: string) => emit('update:modelValue', value)
});

function hasSelectedIcon(row: IconInfoWithId[]): boolean {
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
.icon-select-dropdown {
    padding-left: 8px;
    padding-right: 8px;
}

.icon-select-dropdown .icon-item {
    display: grid;
}
</style>
