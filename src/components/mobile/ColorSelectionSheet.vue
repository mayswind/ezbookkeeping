<template>
    <f7-sheet :opened="show" @sheet:closed="onSheetClosed">
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
                        <f7-icon f7="app_fill" :style="{ color: '#' + colorInfo.color }" @click.native="onColorClicked(colorInfo)">
                            <f7-badge color="default" class="right-bottom-icon" v-if="color && color === colorInfo.color">
                                <f7-icon f7="checkmark_alt"></f7-icon>
                            </f7-badge>
                        </f7-icon>
                    </f7-col>
                    <f7-col v-for="idx in (columnCount - row.length)" :key="idx"></f7-col>
                </f7-row>
            </f7-block>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'color',
        'columnCount',
        'show',
        'allColorInfos'
    ],
    computed: {
        allColorRows() {
            const ret = [];
            let rowCount = -1;

            for (let i = 0; i < this.allColorInfos.length; i++) {
                if (i % this.columnCount === 0) {
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
            this.$emit('color:change', colorInfo.color);
        },
        onSheetClosed() {
            this.$emit('color:closed');
        }
    }
}
</script>
