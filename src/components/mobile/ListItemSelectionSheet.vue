<template>
    <f7-sheet swipe-to-close backdrop swipe-handler=".swipe-handler"
              :class="{ 'list-item-selection-huge-sheet': hugeListItemRows }" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler" @click="close"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content class="margin-top no-padding-top">
            <f7-list dividers no-hairlines class="no-margin-top no-margin-bottom">
                <f7-list-item link="#" no-chevron
                              v-for="(item, index) in items"
                              :key="getItemValue(item, index, keyField, valueType)"
                              :class="{ 'list-item-selected': isSelected(item, index) }"
                              :value="getItemValue(item, index, valueField, valueType)"
                              :title="$tIf((titleField ? item[titleField] : item), titleI18n)"
                              @click="onItemClicked(item, index)">
                    <template #media>
                        <ItemIcon :icon-type="iconType" :icon-id="item[iconField]" :color="item[colorField]" v-if="iconField"></ItemIcon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="isSelected(item, index)"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'modelValue',
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
            currentValue: self.modelValue
        }
    },
    computed: {
        hugeListItemRows() {
            return this.items.length > 10;
        }
    },
    methods: {
        getItemValue(item, index, fieldName, valueType) {
            if (valueType === 'index') {
                return index;
            } else if (fieldName) {
                return item[fieldName];
            } else {
                return item;
            }
        },
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

            this.$emit('update:modelValue', this.currentValue);
            this.close();
        },
        onSheetOpen(event) {
            this.currentValue = this.modelValue;
            this.scrollToSelectedItem(event.$el);
        },
        onSheetClosed() {
            this.close();
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
        },
        scrollToSelectedItem(parent) {
            if (!parent || !parent.length) {
                return;
            }

            const container = parent.find('.page-content');
            const selectedItem = parent.find('li.list-item-selected');

            if (!container.length || !selectedItem.length) {
                return;
            }

            let targetPos = selectedItem.offset().top - container.offset().top - parseInt(container.css('padding-top'), 10)
                - (container.outerHeight() - selectedItem.outerHeight()) / 2;

            if (targetPos <= 0) {
                return;
            }

            container.scrollTop(targetPos);
        },
        close() {
            this.$emit('update:show', false);
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
