<template>
    <f7-sheet :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-block class="margin-vertical">
                <f7-row class="padding-vertical padding-horizontal-half" v-for="(row, idx) in allColorRows" :key="idx">
                    <f7-col class="text-align-center" v-for="colorInfo in row" :key="colorInfo.color">
                        <f7-icon f7="app_fill"
                                 :style="colorInfo.color | iconStyle('default', 'var(--default-icon-color)')"
                                 @click.native="onColorClicked(colorInfo)">
                            <f7-badge color="default" class="right-bottom-icon" v-if="currentValue && currentValue === colorInfo.color">
                                <f7-icon f7="checkmark_alt"></f7-icon>
                            </f7-badge>
                        </f7-icon>
                    </f7-col>
                    <f7-col v-for="idx in (itemPerRow - row.length)" :key="idx"></f7-col>
                </f7-row>
            </f7-block>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'value',
        'columnCount',
        'show',
        'allColorInfos'
    ],
    data() {
        const self = this;

        return {
            currentValue: self.value,
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
            this.$emit('input', this.currentValue);
            this.$emit('update:show', false);
        },
        onSheetOpen() {
            this.currentValue = this.value;
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        }
    }
}
</script>
