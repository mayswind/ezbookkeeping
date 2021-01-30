<template>
    <f7-sheet style="height: auto" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-row>
                <f7-col width="50">
                    <div class="primary-list-container">
                        <f7-list no-hairlines class="primary-list no-margin-top no-margin-bottom">
                            <f7-list-item link="#" no-chevron
                                          v-for="item in items"
                                          :key="item | itemFieldContent(primaryKeyField, item, false)"
                                          :value="item | itemFieldContent(primaryValueField, item, false)"
                                          :title="item | itemFieldContent(primaryTitleField, null, primaryTitleI18n)"
                                          :header="item | itemFieldContent(primaryHeaderField, null, primaryHeaderI18n)"
                                          :footer="item | itemFieldContent(primaryFooterField, null, primaryFooterI18n)"
                                          @click="onPrimaryItemClicked(item)">
                                <f7-icon slot="media"
                                         :icon="item[primaryIconField] | icon(primaryIconType)"
                                         :style="item[primaryColorField] | iconStyle(primaryIconType, 'var(--default-icon-color)')"
                                         v-if="primaryIconField"></f7-icon>
                                <f7-icon slot="after" class="list-item-showing" f7="chevron_right" v-if="item === selectedPrimaryItem"></f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </div>
                </f7-col>
                <f7-col width="50">
                    <div class="secondary-list-container">
                        <f7-list no-hairlines class="secondary-list no-margin-top no-margin-bottom" v-if="selectedPrimaryItem && primarySubItemsField && selectedPrimaryItem[primarySubItemsField]">
                            <f7-list-item link="#" no-chevron
                                          v-for="subItem in selectedPrimaryItem[primarySubItemsField]"
                                          :key="subItem | itemFieldContent(secondaryKeyField, subItem, false)"
                                          :value="subItem | itemFieldContent(secondaryValueField, subItem, false)"
                                          :title="subItem | itemFieldContent(secondaryTitleField, null, secondaryTitleI18n)"
                                          :header="subItem | itemFieldContent(secondaryHeaderField, null, secondaryHeaderI18n)"
                                          :footer="subItem | itemFieldContent(secondaryFooterField, null, secondaryFooterI18n)"
                                          @click="onSecondaryItemClicked(subItem)">
                                <f7-icon slot="media"
                                         :icon="subItem[secondaryIconField] | icon(secondaryIconType)"
                                         :style="subItem[secondaryColorField] | iconStyle(secondaryIconType, 'var(--default-icon-color)')"
                                         v-if="secondaryIconField"></f7-icon>
                                <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="isSecondarySelected(subItem)"></f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </div>
                </f7-col>
            </f7-row>
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
        'primaryHeaderField',
        'primaryHeaderI18n',
        'primaryFooterField',
        'primaryFooterI18n',
        'primaryIconField',
        'primaryIconType',
        'primaryColorField',
        'primarySubItemsField',
        'secondaryKeyField',
        'secondaryValueField',
        'secondaryTitleField',
        'secondaryTitleI18n',
        'secondaryHeaderField',
        'secondaryHeaderI18n',
        'secondaryFooterField',
        'secondaryFooterI18n',
        'secondaryIconField',
        'secondaryIconType',
        'secondaryColorField',
        'items',
        'show'
    ],
    data() {
        const self = this;

        return {
            currentPrimaryValue: null,
            currentSecondaryValue: self.value
        }
    },
    computed: {
        selectedPrimaryItem() {
            if (this.primaryValueField) {
                if (this.$utilities.isArray(this.items)) {
                    for (let i = 0; i < this.items.length; i++) {
                        const item = this.items[i];

                        if (this.currentPrimaryValue === item[this.primaryValueField]) {
                            return item;
                        }
                    }
                } else {
                    for (let field in this.items) {
                        if (!Object.prototype.hasOwnProperty.call(this.items, field)) {
                            continue;
                        }

                        const item = this.items[field];

                        if (this.currentPrimaryValue === item[this.primaryValueField]) {
                            return item;
                        }
                    }
                }
            } else {
                return this.currentPrimaryValue;
            }

            return null;
        }
    },
    methods: {
        onSheetOpen() {
            this.currentPrimaryValue = this.getPrimaryValueBySecondaryValue(this.value);
            this.currentSecondaryValue = this.value;
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        onPrimaryItemClicked(item) {
            if (this.primaryValueField) {
                this.currentPrimaryValue = item[this.primaryValueField];
            } else {
                this.currentPrimaryValue = item;
            }
        },
        onSecondaryItemClicked(subItem) {
            if (this.secondaryValueField) {
                this.currentSecondaryValue = subItem[this.secondaryValueField];
            } else {
                this.currentSecondaryValue = subItem;
            }

            this.$emit('input', this.currentSecondaryValue);
            this.$emit('update:show', false);
        },
        isSecondarySelected(subItem) {
            if (this.secondaryValueField) {
                return this.currentSecondaryValue === subItem[this.secondaryValueField];
            } else {
                return this.currentSecondaryValue === subItem;
            }
        },
        isPrimaryItemHasSecondaryValue(primaryItem, secondaryValue) {
            for (let i = 0; i < primaryItem[this.primarySubItemsField].length; i++) {
                const secondaryItem = primaryItem[this.primarySubItemsField][i];

                if (this.secondaryValueField && secondaryItem[this.secondaryValueField] === secondaryValue) {
                    return true;
                } else if (!this.secondaryValueField && secondaryItem === secondaryValue) {
                    return true;
                }
            }

            return false;
        },
        getPrimaryValueBySecondaryValue(secondaryValue) {
            if (this.primarySubItemsField) {
                if (this.$utilities.isArray(this.items)) {
                    for (let i = 0; i < this.items.length; i++) {
                        const primaryItem = this.items[i];

                        if (this.isPrimaryItemHasSecondaryValue(primaryItem, secondaryValue)) {
                            if (this.primaryValueField) {
                                return primaryItem[this.primaryValueField];
                            } else {
                                return primaryItem;
                            }
                        }
                    }
                } else {
                    for (let field in this.items) {
                        if (!Object.prototype.hasOwnProperty.call(this.items, field)) {
                            continue;
                        }

                        const primaryItem = this.items[field];

                        if (this.isPrimaryItemHasSecondaryValue(primaryItem, secondaryValue)) {
                            if (this.primaryValueField) {
                                return primaryItem[this.primaryValueField];
                            } else {
                                return primaryItem;
                            }
                        }
                    }
                }
            }

            return null;
        }
    }
}
</script>

<style>
.primary-list-container, .secondary-list-container {
    height: 260px;
    overflow-y: auto;
}

.primary-list.list .item-inner {
    padding-right: 6px;
}

.secondary-list.list .item-content {
    padding-left: 0;
}
</style>
