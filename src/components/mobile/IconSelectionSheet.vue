<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="heightClass" :opened="show"
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

<script>
import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getIconsInRows } from '@/lib/icon.ts';
import { scrollToSelectedItem } from '@/lib/ui/mobile.js';

export default {
    props: [
        'modelValue',
        'color',
        'columnCount',
        'show',
        'allIconInfos'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        const self = this;

        return {
            currentValue: self.modelValue,
            itemPerRow: self.columnCount || 7
        }
    },
    computed: {
        allIconRows() {
            return getIconsInRows(this.allIconInfos, this.itemPerRow);
        },
        heightClass() {
            if (this.allIconRows.length > 10) {
                return 'icon-selection-huge-sheet';
            } else if (this.allIconRows.length > 6) {
                return 'icon-selection-large-sheet';
            } else {
                return '';
            }
        }
    },
    methods: {
        onIconClicked(iconInfo) {
            this.currentValue = iconInfo.id;
            this.$emit('update:modelValue', this.currentValue);
        },
        onSheetOpen(event) {
            this.currentValue = this.modelValue;
            scrollToSelectedItem(event.$el, '.page-content', '.row-has-selected-item');
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        hasSelectedIcon(row) {
            return arrayContainsFieldValue(row, 'id', this.currentValue);
        }
    }
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
