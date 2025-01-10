<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="heightClass"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-treeview>
                <f7-treeview-item item-toggle
                                  :opened="isPrimaryItemHasSecondaryValue(item)"
                                  :label="$tIf((primaryTitleField ? item[primaryTitleField] : item), primaryTitleI18n)"
                                  :key="primaryKeyField ? item[primaryKeyField] : item"
                                  v-for="item in items"
                                  v-show="item && (!primaryHiddenField || !item[primaryHiddenField])">
                    <template #media>
                        <ItemIcon :icon-type="primaryIconType" :icon-id="item[primaryIconField]"
                                  :color="item[primaryColorField]" v-if="primaryIconField"></ItemIcon>
                    </template>

                    <f7-treeview-item selectable
                                      :selected="isSecondarySelected(subItem)"
                                      :label="$tIf((secondaryTitleField ? subItem[secondaryTitleField] : subItem), secondaryTitleI18n)"
                                      :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                      v-for="subItem in item[primarySubItemsField]"
                                      v-show="subItem && (!secondaryHiddenField || !subItem[secondaryHiddenField])"
                                      @click="onSecondaryItemClicked(subItem)">
                        <template #media>
                            <ItemIcon :icon-type="secondaryIconType" :icon-id="subItem[secondaryIconField]"
                                      :color="subItem[secondaryColorField]" v-if="secondaryIconField"></ItemIcon>
                        </template>
                    </f7-treeview-item>
                </f7-treeview-item>
            </f7-treeview>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
import { isArray } from '@/lib/common.ts';
import { scrollToSelectedItem } from '@/lib/ui/mobile.ts';

export default {
    props: [
        'modelValue',
        'primaryKeyField',
        'primaryValueField',
        'primaryTitleField',
        'primaryTitleI18n',
        'primaryIconField',
        'primaryIconType',
        'primaryColorField',
        'primaryHiddenField',
        'primarySubItemsField',
        'secondaryKeyField',
        'secondaryValueField',
        'secondaryTitleField',
        'secondaryTitleI18n',
        'secondaryIconField',
        'secondaryIconType',
        'secondaryColorField',
        'secondaryHiddenField',
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
            let count = 0;

            if (isArray(this.items)) {
                count = this.items.length;
            } else {
                for (const field in this.items) {
                    if (!Object.prototype.hasOwnProperty.call(this.items, field)) {
                        continue;
                    }

                    count++;
                }
            }

            if (count > 6) {
                return 'tree-view-selection-huge-sheet';
            } else if (count > 2) {
                return 'tree-view-selection-large-sheet';
            } else {
                return '';
            }
        }
    },
    methods: {
        onSheetOpen(event) {
            this.currentValue = this.modelValue;
            scrollToSelectedItem(event.$el, '.page-content', '.treeview-item .treeview-item-selected');
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        onSecondaryItemClicked(subItem) {
            if (this.secondaryValueField) {
                this.currentValue = subItem[this.secondaryValueField];
            } else {
                this.currentValue = subItem;
            }

            this.$emit('update:modelValue', this.currentValue);
            this.$emit('update:show', false);
        },
        isPrimaryItemHasSecondaryValue(primaryItem) {
            for (let i = 0; i < primaryItem[this.primarySubItemsField].length; i++) {
                const secondaryItem = primaryItem[this.primarySubItemsField][i];

                if (this.secondaryHiddenField && secondaryItem[this.secondaryHiddenField]) {
                    continue;
                }

                if (this.secondaryValueField && secondaryItem[this.secondaryValueField] === this.currentValue) {
                    return true;
                } else if (!this.secondaryValueField && secondaryItem === this.currentValue) {
                    return true;
                }
            }

            return false;
        },
        isSecondarySelected(subItem) {
            if (this.secondaryValueField) {
                return this.currentValue === subItem[this.secondaryValueField];
            } else {
                return this.currentValue === subItem;
            }
        }
    }
}
</script>

<style>
@media (min-height: 630px) {
    .tree-view-selection-large-sheet {
        height: 310px;
    }

    .tree-view-selection-huge-sheet {
        height: 400px;
    }
}
</style>
