<template>
    <v-select
        persistent-placeholder
        :density="density"
        :variant="variant"
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ 'content-class': 'two-column-select-menu' }"
        v-model="currentSecondaryValue"
        v-model:menu="menuState"
        @update:menu="onMenuStateChanged"
    >
        <template #selection>
            <div class="d-flex align-center text-truncate cursor-pointer">
                <span class="text-truncate" v-if="customSelectionPrimaryText">{{ customSelectionPrimaryText }}</span>
                <v-icon class="disabled" :icon="icons.chevronRight" size="23" v-if="customSelectionPrimaryText && customSelectionSecondaryText" />
                <span class="text-truncate" v-if="customSelectionPrimaryText && customSelectionSecondaryText">{{ customSelectionSecondaryText }}</span>
                <span class="text-truncate" v-if="!customSelectionPrimaryText && !selectedPrimaryItem && !selectedSecondaryItem">{{ noSelectionText }}</span>
                <span class="text-truncate" v-if="!customSelectionPrimaryText && showSelectionPrimaryText && selectedPrimaryItem">{{ selectionPrimaryItemText }}</span>
                <v-icon class="disabled" :icon="icons.chevronRight" size="23" v-if="!customSelectionPrimaryText && showSelectionPrimaryText && selectedPrimaryItem && selectedSecondaryItem" />
                <ItemIcon class="mr-2" icon-type="account" size="21.5px"
                          :icon-id="selectedSecondaryItem ? selectedSecondaryItem[secondaryIconField] : null"
                          :color="selectedSecondaryItem ? selectedSecondaryItem[secondaryColorField] : null"
                          v-if="!customSelectionPrimaryText && selectedSecondaryItem && showSelectionSecondaryIcon" />
                <span class="text-truncate" v-if="!customSelectionPrimaryText && selectedSecondaryItem">{{ selectionSecondaryItemText }}</span>
            </div>
        </template>

        <template #no-data>
            <div ref="dropdownMenu" class="two-column-list-container">
                <div class="primary-list-container">
                    <v-list :class="{ 'list-item-with-header': !!primaryHeaderField, 'list-item-with-footer': !!primaryFooterField }">
                        <v-list-item :class="{ 'primary-list-item-selected v-list-item--active text-primary': item === selectedPrimaryItem }"
                                     :key="primaryKeyField ? item[primaryKeyField] : item"
                                     v-for="item in items"
                                     v-show="item && (!primaryHiddenField || !item[primaryHiddenField])"
                                     @click="onPrimaryItemClicked(item)">
                            <template #prepend>
                                <ItemIcon class="mr-2" :icon-type="primaryIconType"
                                          :icon-id="item[primaryIconField]" :color="item[primaryColorField]"></ItemIcon>
                            </template>
                            <template #title>
                                <div class="list-item-header text-truncate" v-if="primaryHeaderField">{{ $tIf(item[primaryHeaderField], primaryHeaderI18n) }}</div>
                                <div class="text-truncate">{{ $tIf(item[primaryTitleField], primaryTitleI18n) }}</div>
                                <div class="list-item-footer text-truncate" v-if="primaryFooterField">{{ $tIf(item[primaryFooterField], primaryFooterI18n) }}</div>
                            </template>
                        </v-list-item>
                    </v-list>
                </div>
                <div class="secondary-list-container">
                    <v-list :class="{ 'list-item-with-header': !!secondaryHeaderField, 'list-item-with-footer': !!secondaryFooterField }"
                            v-if="selectedPrimaryItem && primarySubItemsField && selectedPrimaryItem[primarySubItemsField]">
                        <v-list-item :class="{ 'secondary-list-item-selected v-list-item--active text-primary': isSecondarySelected(subItem) }"
                                     :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                     v-for="subItem in selectedPrimaryItem[primarySubItemsField]"
                                     v-show="subItem && (!secondaryHiddenField || !subItem[secondaryHiddenField])"
                                     @click="onSecondaryItemClicked(subItem)">
                            <template #prepend>
                                <ItemIcon class="mr-2" :icon-type="secondaryIconType"
                                          :icon-id="subItem[secondaryIconField]" :color="subItem[secondaryColorField]"></ItemIcon>
                            </template>
                            <template #title>
                                <div class="list-item-header text-truncate" v-if="secondaryHeaderField">{{ $tIf(subItem[secondaryHeaderField], secondaryHeaderI18n) }}</div>
                                <div class="text-truncate">{{ $tIf(subItem[secondaryTitleField], secondaryTitleI18n) }}</div>
                                <div class="list-item-footer text-truncate" v-if="secondaryFooterField">{{ $tIf(subItem[secondaryFooterField], secondaryFooterI18n) }}</div>
                            </template>
                        </v-list-item>
                    </v-list>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script>
import {
    getFirstVisibleItem,
    getItemByKeyValue,
    getNameByKeyValue,
    getPrimaryValueBySecondaryValue
} from '@/lib/common.js';
import { scrollToSelectedItem } from '@/lib/ui.desktop.js';

import {
    mdiChevronRight
} from '@mdi/js';

export default {
    props: [
        'modelValue',
        'density',
        'variant',
        'disabled',
        'readonly',
        'label',
        'showSelectionPrimaryText',
        'showSelectionSecondaryIcon',
        'customSelectionPrimaryText',
        'customSelectionSecondaryText',
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
        'primaryHiddenField',
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
        'secondaryHiddenField',
        'noItemText',
        'items'
    ],
    emits: [
        'update:modelValue'
    ],
    data() {
        return {
            menuState: false,
            icons: {
                chevronRight: mdiChevronRight
            }
        }
    },
    computed: {
        currentPrimaryValue: {
            get: function () {
                return this.getPrimaryValueBySecondaryValue(this.modelValue);
            },
            set: function (value) {
                const primaryItem = getItemByKeyValue(this.items, value, this.primaryValueField);
                const secondaryItem = getFirstVisibleItem(primaryItem[this.primarySubItemsField], this.primaryHiddenField);

                if (secondaryItem) {
                    if (this.secondaryValueField) {
                        this.$emit('update:modelValue', secondaryItem[this.secondaryValueField]);
                    }
                }
            }
        },
        currentSecondaryValue: {
            get: function () {
                return this.modelValue;
            },
            set: function (value) {
                this.menuState = false;
                this.$emit('update:modelValue', value);
            }
        },
        selectedPrimaryItem() {
            if (this.primaryValueField) {
                return getItemByKeyValue(this.items, this.currentPrimaryValue, this.primaryValueField);
            } else {
                return this.currentPrimaryValue;
            }
        },
        selectedSecondaryItem() {
            if (this.currentSecondaryValue && this.selectedPrimaryItem && this.selectedPrimaryItem[this.primarySubItemsField]) {
                return getItemByKeyValue(this.selectedPrimaryItem[this.primarySubItemsField], this.currentSecondaryValue, this.secondaryValueField);
            } else {
                return null;
            }
        },
        noSelectionText() {
            return this.noItemText ? this.noItemText : this.$t('None');
        },
        selectionPrimaryItemText() {
            if (this.primaryValueField && this.primaryTitleField) {
                if (this.currentPrimaryValue) {
                    return getNameByKeyValue(this.items, this.currentPrimaryValue, this.primaryValueField, this.primaryTitleField, this.noSelectionText);
                } else {
                    return this.noSelectionText;
                }
            } else {
                return this.currentPrimaryValue;
            }
        },
        selectionSecondaryItemText() {
            if (this.secondaryValueField && this.secondaryTitleField) {
                if (this.currentSecondaryValue && this.selectedPrimaryItem && this.selectedPrimaryItem[this.primarySubItemsField]) {
                    return getNameByKeyValue(this.selectedPrimaryItem[this.primarySubItemsField], this.currentSecondaryValue, this.secondaryValueField, this.secondaryTitleField, this.noSelectionText);
                } else {
                    return this.noSelectionText;
                }
            } else {
                return this.currentSecondaryValue;
            }
        }
    },
    methods: {
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
        },
        isSecondarySelected(subItem) {
            if (this.secondaryValueField) {
                return this.currentSecondaryValue === subItem[this.secondaryValueField];
            } else {
                return this.currentSecondaryValue === subItem;
            }
        },
        getPrimaryValueBySecondaryValue(secondaryValue) {
            return getPrimaryValueBySecondaryValue(this.items, this.primarySubItemsField, this.primaryValueField, this.primaryHiddenField, this.secondaryValueField, this.secondaryHiddenField, secondaryValue);
        },
        onMenuStateChanged(state) {
            const self = this;

            if (state) {
                self.$nextTick(() => {
                    if (self.$refs.dropdownMenu && self.$refs.dropdownMenu.parentElement) {
                        scrollToSelectedItem(self.$refs.dropdownMenu.parentElement, '.primary-list-container', '.primary-list-item-selected');
                        scrollToSelectedItem(self.$refs.dropdownMenu.parentElement, '.secondary-list-container', '.secondary-list-item-selected');
                    }
                });
            }
        }
    }
}
</script>

<style>
.two-column-select-menu {
    max-height: inherit !important;
}

.two-column-select-menu > .v-list {
    padding: 0;
}

.two-column-select-menu .two-column-list-container {
    width: 100%;
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
}

.two-column-select-menu .primary-list-container,
.two-column-select-menu .secondary-list-container {
    width: 100%;
    max-height: 310px;
    overflow-y: scroll;
}

.two-column-select-menu .list-item-with-header > .v-list-item,
.two-column-select-menu .list-item-with-footer > .v-list-item {
    min-height: 58px;
    padding-top: 6px;
    padding-bottom: 6px;
}

.two-column-select-menu .list-item-with-header.list-item-with-footer > .v-list-item {
    min-height: 78px;
    padding-top: 8px;
    padding-bottom: 8px;
}

.two-column-select-menu .list-item-header,
.two-column-select-menu .list-item-footer {
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity));
    font-size: 0.75rem;
    line-height: 1.2rem;
}
</style>
