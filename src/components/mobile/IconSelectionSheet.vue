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
                <div class="grid grid-cols-7 padding-vertical-half padding-horizontal-half"
                     :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                     :key="idx"
                     v-for="(row, idx) in allIconRows">
                    <div class="text-align-center" v-for="iconInfo in row" :key="iconInfo.id">
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
            const ret = [];
            let rowCount = 0;

            for (let iconInfoId in this.allIconInfos) {
                if (!Object.prototype.hasOwnProperty.call(this.allIconInfos, iconInfoId)) {
                    continue;
                }

                const iconInfo = this.allIconInfos[iconInfoId];

                if (!ret[rowCount]) {
                    ret[rowCount] = [];
                } else if (ret[rowCount] && ret[rowCount].length >= this.itemPerRow) {
                    rowCount++;
                    ret[rowCount] = [];
                }

                ret[rowCount].push({
                    id: iconInfoId,
                    icon: iconInfo.icon
                });
            }

            return ret;
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
