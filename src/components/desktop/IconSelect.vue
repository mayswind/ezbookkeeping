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

<script>
import { arrayContainsFieldvalue } from '@/lib/common.js';
import { getIconsInRows } from '@/lib/icon.js';
import { scrollToSelectedItem } from '@/lib/ui.desktop.js';

import {
    mdiCheck
} from '@mdi/js';

export default {
    props: [
        'modelValue',
        'disabled',
        'label',
        'iconType',
        'color',
        'columnCount',
        'allIconInfos'
    ],
    emits: [
        'update:modelValue',
    ],
    data() {
        const self = this;

        return {
            itemPerRow: self.columnCount || 7,
            icons: {
                checked: mdiCheck
            }
        }
    },
    computed: {
        allIconRows() {
            return getIconsInRows(this.allIconInfos, this.itemPerRow);
        },
        icon: {
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
            return arrayContainsFieldvalue(row, 'id', this.modelValue);
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
.icon-select-dropdown {
    padding-left: 8px;
    padding-right: 8px;
}

.icon-select-dropdown .icon-item {
    display: grid;
}
</style>
