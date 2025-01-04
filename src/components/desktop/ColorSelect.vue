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
                <v-icon size="28" :icon="icons.square" :color="getFinalColor(item.raw)"/>
            </v-label>
        </template>

        <template #no-data>
            <div class="color-select-dropdown" ref="dropdownMenu">
                <div class="color-item" :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                     :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                     :key="idx" v-for="(row, idx) in allColorRows">
                    <div class="text-center" :key="colorInfo.color" v-for="colorInfo in row">
                        <div class="cursor-pointer" @click="color = colorInfo.color">
                            <v-icon class="ma-2" size="28"
                                    :icon="icons.square" :color="getFinalColor(colorInfo.color)"
                                    v-if="!modelValue || modelValue !== colorInfo.color" />
                            <v-badge class="right-bottom-icon" color="primary"
                                     location="bottom right" offset-x="8" offset-y="8" :icon="icons.checked"
                                     v-if="modelValue && modelValue === colorInfo.color">
                                <v-icon class="ma-2" size="28" :icon="icons.square" :color="getFinalColor(colorInfo.color)" />
                            </v-badge>
                        </div>
                    </div>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script>
import { DEFAULT_ICON_COLOR } from '@/consts/color.ts';
import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getColorsInRows } from '@/lib/color.ts';
import { scrollToSelectedItem } from '@/lib/ui/desktop.ts';

import {
    mdiSquareRounded,
    mdiCheck
} from '@mdi/js';

export default {
    props: [
        'modelValue',
        'disabled',
        'label',
        'columnCount',
        'allColorInfos'
    ],
    emits: [
        'update:modelValue',
    ],
    data() {
        const self = this;

        return {
            itemPerRow: self.columnCount || 7,
            icons: {
                square: mdiSquareRounded,
                checked: mdiCheck
            }
        }
    },
    computed: {
        allColorRows() {
            return getColorsInRows(this.allColorInfos, this.itemPerRow);
        },
        color: {
            get: function () {
                return this.modelValue;
            },
            set: function (value) {
                this.$emit('update:modelValue', value);
            }
        }
    },
    methods: {
        hasSelectedIcon(row) {
            return arrayContainsFieldValue(row, 'id', this.modelValue);
        },
        getFinalColor(color) {
            if (color && color !== DEFAULT_ICON_COLOR) {
                return '#' + color;
            } else {
                return 'var(--default-icon-color)';
            }
        },
        onMenuStateChanged(state) {
            const self = this;

            if (state) {
                self.$nextTick(() => {
                    if (self.$refs.dropdownMenu && self.$refs.dropdownMenu.parentElement) {
                        scrollToSelectedItem(self.$refs.dropdownMenu.parentElement, null, '.row-has-selected-item');
                    }
                });
            }
        }
    }
}
</script>

<style>
.color-select-dropdown {
    padding-left: 8px;
    padding-right: 8px;
}

.color-select-dropdown .color-item {
    display: grid;
}
</style>
