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
                <div class="grid grid-cols-7 padding-vertical-half padding-horizontal-half"
                     :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
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
            const ret = [];
            let rowCount = -1;

            for (let i = 0; i < this.allColorInfos.length; i++) {
                if (i % this.itemPerRow === 0) {
                    ret[++rowCount] = [];
                }

                ret[rowCount].push({
                    color: this.allColorInfos[i]
                });
            }

            return ret;
        }
    },
    methods: {
        onColorClicked(colorInfo) {
            this.currentValue = colorInfo.color;
            this.$emit('update:modelValue', this.currentValue);
            this.$emit('update:show', false);
        },
        onSheetOpen(event) {
            this.currentValue = this.modelValue;
            this.scrollToSelectedItem(event.$el);
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        hasSelectedIcon(row) {
            if (!this.currentValue || !row || !row.length) {
                return false;
            }

            for (let i = 0; i < row.length; i++) {
                if (row[i].id === this.currentValue) {
                    return true;
                }
            }

            return false;
        },
        scrollToSelectedItem(parent) {
            if (!parent || !parent.length) {
                return;
            }

            const container = parent.find('.page-content');
            const selectedItem = parent.find('.row.row-has-selected-item');

            if (!container.length || !selectedItem.length) {
                return;
            }

            let targetPos = selectedItem.offset().top - container.offset().top - parseInt(container.css('padding-top'), 10)
                - (container.outerHeight() - selectedItem.outerHeight()) / 2;

            if (targetPos <= 0) {
                return;
            }

            container.scrollTop(targetPos);
        }
    }
}
</script>
