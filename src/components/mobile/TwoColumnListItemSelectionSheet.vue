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
                                          :key="primaryKeyField ? item[primaryKeyField] : item"
                                          :value="primaryValueField ? item[primaryValueField] : item"
                                          :title="primaryTitleField ? item[primaryTitleField] : item"
                                          @click="onPrimaryItemClicked(item)">
                                <f7-icon slot="media"
                                         :icon="item[primaryIconField] | icon(primaryIconType)"
                                         :style="{ color: (primaryColorField && item[primaryColorField] && item[primaryColorField] !== '000000' ? '#' + item[primaryColorField] : 'var(--default-icon-color)') }"
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
                                          :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                          :value="secondaryValueField ? subItem[secondaryValueField] : subItem"
                                          :title="secondaryTitleField ? subItem[secondaryTitleField] : subItem"
                                          @click="onSecondaryItemClicked(subItem)">
                                <f7-icon slot="media"
                                         :icon="subItem[secondaryIconField] | icon(secondaryIconType)"
                                         :style="{ color: (secondaryColorField && subItem[secondaryColorField] && subItem[secondaryColorField] !== '000000' ? '#' + subItem[secondaryColorField] : 'var(--default-icon-color)') }"
                                         v-if="secondaryIconField"></f7-icon>
                                <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="isSecondarySelected(subItem)"></f7-icon>
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
        'primaryIconField',
        'primaryIconType',
        'primaryColorField',
        'primarySubItemsField',
        'secondaryKeyField',
        'secondaryValueField',
        'secondaryTitleField',
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
                for (let i = 0; i < this.items.length; i++) {
                    const item = this.items[i];
                    if (this.currentPrimaryValue === item[this.primaryValueField]) {
                        return item;
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
        getPrimaryValueBySecondaryValue(secondaryValue) {
            if (this.primarySubItemsField) {
                for (let i = 0; i < this.items.length; i++) {
                    const primaryItem = this.items[i];

                    for (let j = 0; j < primaryItem[this.primarySubItemsField].length; j++) {
                        const secondaryItem = primaryItem[this.primarySubItemsField][j];

                        if (this.secondaryValueField && secondaryItem[this.secondaryValueField] === secondaryValue) {
                            if (this.primaryValueField) {
                                return primaryItem[this.primaryValueField];
                            } else {
                                return primaryItem;
                            }
                        } else if (!this.secondaryValueField && secondaryItem === secondaryValue) {
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
