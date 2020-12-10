<template>
    <f7-sheet :class="{ 'icon-selection-huge-sheet': hugeIconRows }" :opened="show" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-block class="margin-vertical">
                <f7-row class="padding-vertical-half padding-horizontal-half" v-for="(row, idx) in allIconRows" :key="idx">
                    <f7-col class="text-align-center" v-for="iconInfo in row" :key="iconInfo.id">
                        <f7-icon :icon="iconInfo.icon" :style="{ color: '#' + (color || '000000') }" @click.native="onIconClicked(iconInfo)">
                            <f7-badge color="default" class="right-bottom-icon" v-if="icon && icon === iconInfo.id">
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
        'icon',
        'color',
        'columnCount',
        'show',
        'allIconInfos'
    ],
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
                } else if (ret[rowCount] && ret[rowCount].length >= this.columnCount) {
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
        hugeIconRows() {
            return this.allIconRows.length > 10;
        }
    },
    methods: {
        onIconClicked(iconInfo) {
            this.$emit('icon:change', iconInfo.id);
        },
        onSheetClosed() {
            this.$emit('icon:closed');
        }
    }
}
</script>

<style>
@media (min-height: 630px) {
    .icon-selection-huge-sheet {
        height: 400px;
    }
}
</style>
