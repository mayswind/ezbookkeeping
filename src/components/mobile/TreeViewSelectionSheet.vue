<template>
    <f7-sheet :class="{ 'tree-view-selection-huge-sheet': hugeTreeViewItems }" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-treeview>
                <f7-treeview-item v-for="item in items"
                                  item-toggle
                                  :opened="isPrimaryItemHasSecondaryValue(item)"
                                  :key="primaryKeyField ? item[primaryKeyField] : item"
                                  :label="primaryTitleField ? (primaryTitleI18n ? $t(item[primaryTitleField]) : item[primaryTitleField]) : (primaryTitleI18n ? $t(item) : item)">
                    <f7-icon slot="media"
                             :icon="item[primaryIconField] | icon(primaryIconType)"
                             :style="{ color: (primaryColorField && item[primaryColorField] && item[primaryColorField] !== '000000' ? '#' + item[primaryColorField] : 'var(--default-icon-color)') }"
                             v-if="primaryIconField"></f7-icon>

                    <f7-treeview-item v-for="subItem in item[primarySubItemsField]"
                                      selectable
                                      :selected="isSecondarySelected(subItem)"
                                      :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                      :label="secondaryTitleField ? (secondaryTitleI18n ? $t(subItem[secondaryTitleField]) : subItem[secondaryTitleField]) : (secondaryTitleI18n ? $t(subItem) : subItem)"
                                      @click="onSecondaryItemClicked(subItem)">
                        <f7-icon slot="media"
                                 :icon="subItem[secondaryIconField] | icon(secondaryIconType)"
                                 :style="{ color: (secondaryColorField && subItem[secondaryColorField] && subItem[secondaryColorField] !== '000000' ? '#' + subItem[secondaryColorField] : 'var(--default-icon-color)') }"
                                 v-if="secondaryIconField"></f7-icon>
                    </f7-treeview-item>
                </f7-treeview-item>
            </f7-treeview>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'value',
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
    data() {
        const self = this;

        return {
            currentValue: self.value
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
        onSheetOpen() {
            this.currentValue = this.value;
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

            this.$emit('input', this.currentValue);
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
