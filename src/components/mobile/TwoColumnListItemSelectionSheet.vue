<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              style="height: auto" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="$t('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <div class="grid grid-cols-2 grid-gap">
                <div>
                    <div class="primary-list-container">
                        <f7-list dividers class="primary-list no-margin-vertical">
                            <f7-list-item link="#" no-chevron
                                          :class="{ 'primary-list-item-selected': item === selectedPrimaryItem }"
                                          :value="primaryValueField ? item[primaryValueField] : item"
                                          :title="$tIf(item[primaryTitleField], primaryTitleI18n)"
                                          :header="$tIf(item[primaryHeaderField], primaryHeaderI18n)"
                                          :footer="$tIf(item[primaryFooterField], primaryFooterI18n)"
                                          :key="primaryKeyField ? item[primaryKeyField] : item"
                                          v-for="item in items"
                                          @click="onPrimaryItemClicked(item)">
                                <template #media>
                                    <ItemIcon :icon-type="primaryIconType" :icon-id="item[primaryIconField]" :color="item[primaryColorField]"></ItemIcon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-showing" f7="chevron_right" v-if="item === selectedPrimaryItem"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </div>
                </div>
                <div>
                    <div class="secondary-list-container">
                        <f7-list dividers class="secondary-list no-margin-vertical" v-if="selectedPrimaryItem && primarySubItemsField && selectedPrimaryItem[primarySubItemsField]">
                            <f7-list-item link="#" no-chevron
                                          :class="{ 'secondary-list-item-selected': isSecondarySelected(subItem) }"
                                          :value="secondaryValueField ? subItem[secondaryValueField] : subItem"
                                          :title="$tIf(subItem[secondaryTitleField], secondaryTitleI18n)"
                                          :header="$tIf(subItem[secondaryHeaderField], secondaryHeaderI18n)"
                                          :footer="$tIf(subItem[secondaryFooterField], secondaryFooterI18n)"
                                          :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                          v-for="subItem in selectedPrimaryItem[primarySubItemsField]"
                                          @click="onSecondaryItemClicked(subItem)">
                                <template #media>
                                    <ItemIcon :icon-type="secondaryIconType" :icon-id="subItem[secondaryIconField]" :color="subItem[secondaryColorField]"></ItemIcon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="isSecondarySelected(subItem)"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </div>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
import { isArray } from '@/lib/common.js';

export default {
    props: [
        'modelValue',
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
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        const self = this;

        return {
            currentPrimaryValue: self.getPrimaryValueBySecondaryValue(self.modelValue),
            currentSecondaryValue: self.modelValue
        }
    },
    computed: {
        selectedPrimaryItem() {
            if (this.primaryValueField) {
                if (isArray(this.items)) {
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
        onSheetOpen(event) {
            this.currentPrimaryValue = this.getPrimaryValueBySecondaryValue(this.modelValue);
            this.currentSecondaryValue = this.modelValue;
            this.scrollToSelectedItem(event.$el, '.primary-list-container', 'li.primary-list-item-selected');
            this.scrollToSelectedItem(event.$el, '.secondary-list-container', 'li.secondary-list-item-selected');
        },
        onSheetClosed() {
            this.close();
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

            this.$emit('update:modelValue', this.currentSecondaryValue);
            this.close();
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
                if (isArray(this.items)) {
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
        },
        scrollToSelectedItem(parent, containerSelector, selectedItemSelector) {
            if (!parent || !parent.length) {
                return;
            }

            const container = parent.find(containerSelector);
            const selectedItem = parent.find(selectedItemSelector);

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
