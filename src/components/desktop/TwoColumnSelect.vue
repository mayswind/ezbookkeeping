<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ 'content-class': 'two-column-select-menu' }"
        v-model="currentSecondaryValue"
        @update:menu="onMenuStateChanged"
    >
        <template #selection>
            <div class="d-flex align-center text-truncate cursor-pointer">
                <span class="text-truncate" v-if="showPrimaryName">{{ primaryItemDisplayName }}</span>
                <v-icon class="disabled" :icon="icons.chevronRight" size="23" v-if="showPrimaryName" />
                <ItemIcon class="mr-2" icon-type="account" size="21.5px"
                          :icon-id="selectedSecondaryItem ? selectedSecondaryItem[secondaryIconField] : null"
                          :color="selectedSecondaryItem ? selectedSecondaryItem[secondaryColorField] : null"
                          v-if="selectedSecondaryItem && showSecondaryIcon" />
                <span class="text-truncate">{{ secondaryItemDisplayName }}</span>
            </div>
        </template>

        <template #no-data>
            <div ref="dropdownMenu" class="two-column-list-container">
                <div class="primary-list-container">
                    <v-list>
                        <v-list-item :class="{ 'primary-list-item-selected v-list-item--active text-primary': item === selectedPrimaryItem }"
                                     :title="$tIf(item[primaryTitleField], primaryTitleI18n)"
                                     :key="primaryKeyField ? item[primaryKeyField] : item"
                                     v-for="item in items"
                                     @click="onPrimaryItemClicked(item)">
                            <template #prepend>
                                <ItemIcon class="mr-2" :icon-type="primaryIconType"
                                          :icon-id="item[primaryIconField]" :color="item[primaryColorField]"></ItemIcon>
                            </template>
                        </v-list-item>
                    </v-list>
                </div>
                <div class="secondary-list-container">
                    <v-list v-if="selectedPrimaryItem && primarySubItemsField && selectedPrimaryItem[primarySubItemsField]">
                        <v-list-item :class="{ 'secondary-list-item-selected v-list-item--active text-primary': isSecondarySelected(subItem) }"
                                     :title="$tIf(subItem[secondaryTitleField], secondaryTitleI18n)"
                                     :key="secondaryKeyField ? subItem[secondaryKeyField] : subItem"
                                     v-for="subItem in selectedPrimaryItem[primarySubItemsField]"
                                     @click="onSecondaryItemClicked(subItem)">
                            <template #prepend>
                                <ItemIcon class="mr-2" :icon-type="secondaryIconType"
                                          :icon-id="subItem[secondaryIconField]" :color="subItem[secondaryColorField]"></ItemIcon>
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
    getFirstItem,
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
        'disabled',
        'readonly',
        'label',
        'showPrimaryName',
        'showSecondaryIcon',
        'primaryKeyField',
        'primaryValueField',
        'primaryTitleField',
        'primaryTitleI18n',
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
        'secondaryFooterField',
        'secondaryFooterI18n',
        'secondaryIconField',
        'secondaryIconType',
        'secondaryColorField',
        'noItemText',
        'items'
    ],
    emits: [
        'update:modelValue'
    ],
    data() {
        return {
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
                const secondaryItem = getFirstItem(primaryItem[this.primarySubItemsField]);

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
        noItemDisplayName() {
            return this.noItemText ? this.noItemText : this.$t('None');
        },
        primaryItemDisplayName() {
            if (this.primaryValueField && this.primaryTitleField) {
                if (this.currentPrimaryValue) {
                    return getNameByKeyValue(this.items, this.currentPrimaryValue, this.primaryValueField, this.primaryTitleField, this.noItemDisplayName);
                } else {
                    return this.noItemDisplayName;
                }
            } else {
                return this.currentPrimaryValue;
            }
        },
        secondaryItemDisplayName() {
            if (this.secondaryValueField && this.secondaryTitleField) {
                if (this.currentSecondaryValue && this.selectedPrimaryItem && this.selectedPrimaryItem[this.primarySubItemsField]) {
                    return getNameByKeyValue(this.selectedPrimaryItem[this.primarySubItemsField], this.currentSecondaryValue, this.secondaryValueField, this.secondaryTitleField, this.noItemDisplayName);
                } else {
                    return this.noItemDisplayName;
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
            return getPrimaryValueBySecondaryValue(this.items, this.primarySubItemsField, this.primaryValueField, this.secondaryValueField, secondaryValue);
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
</style>
