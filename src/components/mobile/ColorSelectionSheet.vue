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

<script>
import { arrayContainsFieldvalue } from '@/lib/common.js';
import { getColorsInRows } from '@/lib/color.js';
import { scrollToSelectedItem } from '@/lib/ui.mobile.js';

export default {
    props: [
        'modelValue',
        'columnCount',
        'show',
        'allColorInfos'
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
        allColorRows() {
            return getColorsInRows(this.allColorInfos, this.itemPerRow);
        }
    },
    methods: {
        onColorClicked(colorInfo) {
            this.currentValue = colorInfo.color;
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
            return arrayContainsFieldvalue(row, 'id', this.currentValue);
        }
    }
}
</script>
