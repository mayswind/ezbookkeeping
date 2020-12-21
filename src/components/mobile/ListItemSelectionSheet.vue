<template>
    <f7-sheet :class="{ 'list-item-selection-huge-sheet': hugeListItemRows }" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-list no-hairlines class="no-margin-top no-margin-bottom">
                <f7-list-item link="#" no-chevron
                              v-for="(item, index) in items"
                              :key="valueType === 'index' ? index : (keyField ? item[keyField] : item)"
                              :value="valueType === 'index' ? index : (valueField ? item[valueField] : item)"
                              :title="titleField ? (titleI18n ? $t(item[titleField]) : item[titleField]) : (titleI18n ? $t(item) : item)"
                              @click="onItemClicked(item, index)">
                    <f7-icon slot="media"
                             :icon="item[iconField] | icon(iconType)"
                             :style="item[colorField] | iconStyle(iconType, 'var(--default-icon-color)')"
                             v-if="iconField"></f7-icon>
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="isSelected(item, index)"></f7-icon>
                </f7-list-item>
            </f7-list>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'value',
        'valueType', // item or index
        'keyField', // for value type == item
        'valueField', // for value type == item
        'titleField',
        'titleI18n',
        'iconField',
        'iconType',
        'colorField',
        'items',
        'show'
    ],
    data() {
        const self = this;

        return {
            currentValue: self.value
        }
    },
    computed: {
        hugeListItemRows() {
            return this.items.length > 10;
        }
    },
    methods: {
        onItemClicked(item, index) {
            if (this.valueType === 'index') {
                this.currentValue = index;
            } else {
                if (this.valueField) {
                    this.currentValue = item[this.valueField];
                } else {
                    this.currentValue = item;
                }
            }

            this.$emit('input', this.currentValue);
            this.$emit('update:show', false);
        },
        onSheetOpen() {
            this.currentValue = this.value;
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        isSelected(item, index) {
            if (this.valueType === 'index') {
                return this.currentValue === index;
            } else {
                if (this.valueField) {
                    return this.currentValue === item[this.valueField];
                } else {
                    return this.currentValue === item;
                }
            }
        }
    }
}
</script>

<style>
@media (min-height: 630px) {
    .list-item-selection-huge-sheet {
        height: 400px;
    }
}
</style>
