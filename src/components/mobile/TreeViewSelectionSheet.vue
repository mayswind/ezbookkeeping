<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="{ 'tree-view-selection-huge-sheet': hugeTreeViewItems }"
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
                                  v-for="item in items">
                    <template #media>
                        <ItemIcon :icon-type="primaryIconType" :icon-id="item[primaryIconField]"
                                  :color="item[primaryColorField]" v-if="primaryIconField"></ItemIcon>
                    </template>

                    <f7-treeview-item selectable
                                      :selected="isSecondarySelected(subItem)"
                                      :label="$tIf((secondaryTitleField ? subItem[secondaryTitleField] : subItem), secondaryTitleI18n)"
                                      :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                      v-for="subItem in item[primarySubItemsField]"
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
        'primarySubItemsField',
        'secondaryKeyField',
        'secondaryValueField',
        'secondaryTitleField',
        'secondaryTitleI18n',
        'secondaryIconField',
        'secondaryIconType',
        'secondaryColorField',
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
        hugeTreeViewItems() {
            if (this.$utilities.isArray(this.items)) {
                return this.items.length > 10;
            } else {
                let count = 0;

                for (let field in this.items) {
                    if (!Object.prototype.hasOwnProperty.call(this.items, field)) {
                        continue;
                    }

                    count++;
                }

                return count > 10;
            }
        }
    },
    methods: {
        onSheetOpen(event) {
            this.currentValue = this.modelValue;
            this.scrollToSelectedItem(event.$el);
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
        },
        scrollToSelectedItem(parent) {
            if (!parent || !parent.length) {
                return;
            }

            const container = parent.find('.page-content');
            const selectedItem = parent.find('.treeview-item .treeview-item-selected');

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
    .tree-view-selection-huge-sheet {
        height: 400px;
    }
}
</style>
