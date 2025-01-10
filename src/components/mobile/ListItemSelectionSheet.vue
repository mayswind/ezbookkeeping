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
            <f7-list dividers class="no-margin-vertical">
                <f7-list-item link="#" no-chevron
                              :title="$tIf((titleField ? item[titleField] : item), titleI18n)"
                              :value="getItemValue(item, index, valueField, valueType)"
                              :class="{ 'list-item-selected': isSelected(item, index) }"
                              :key="getItemValue(item, index, keyField, valueType)"
                              v-for="(item, index) in items"
                              v-show="item && (!hiddenField || !item[hiddenField])"
                              @click="onItemClicked(item, index)">
                    <template #content-start>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" :style="{ 'color': isSelected(item, index) ? '' : 'transparent' }"></f7-icon>
                    </template>
                    <template #media v-if="iconField">
                        <ItemIcon :icon-type="iconType" :icon-id="item[iconField]" :color="item[colorField]"></ItemIcon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
import { scrollToSelectedItem } from '@/lib/ui/mobile.ts';

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
        'hiddenField',
        'items',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        const self = this;

        return {
            currentValue: self.modelValue
        }
    },
    computed: {
        heightClass() {
            if (this.items.length > 10) {
                return 'list-item-selection-huge-sheet';
            } else if (this.items.length > 6) {
                return 'list-item-selection-large-sheet';
            } else {
                return '';
            }
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
            scrollToSelectedItem(event.$el, '.page-content', 'li.list-item-selected');
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
        close() {
            this.$emit('update:show', false);
        }
    }
}
</script>

<style>
@media (min-height: 630px) {
    .list-item-selection-large-sheet {
        height: 310px;
    }

    .list-item-selection-huge-sheet {
        height: 400px;
    }
}
</style>
