<template>
    <div class="pie-chart-container">
        <svg class="pie-chart" :viewBox="`${-diameter} ${-diameter} ${diameter * 2} ${diameter * 2}`">
            <circle class="pie-chart-background" cx="0" cy="0" :r="diameter"></circle>

            <circle class="pie-chart-item"
                    fill="transparent"
                    cx="0" cy="0"
                    :r="diameter / 2"
                    :stroke="getColor(item.color)"
                    :stroke-width="diameter"
                    :stroke-dasharray="getItemStrokeDash(item)"
                    :stroke-dashoffset="getItemDashOffset(item, validItems, itemCommonDashOffset)"
                    :key="idx"
                    v-for="(item, idx) in validItems"
                    @click="switchSelectedIndex(idx)">
            </circle>

            <circle class="pie-chart-text-background"
                    cx="0" cy="0"
                    stroke="#ddd"
                    :style="{ '--pie-chart-text-background': centerTextBackground ? centerTextBackground : '#7f2020' }"
                    :r="diameter / 2.5"
                    v-if="showCenterText"/>

            <circle cx="0" cy="0"
                    stroke="#ddd"
                    fill="transparent"
                    :r="diameter / 2.5"
                    v-if="showCenterText"/>

            <clipPath id="pie-chart-text-clip">
                <rect :x="-diameter / 2.5 + 2" :y="-diameter / 2.5 + 2" :width="diameter / 1.25 - 4" :height="diameter / 1.25 -4 "/>
            </clipPath>

            <g class="pie-chart-text-group" clip-path="url(#pie-chart-text-clip)" v-if="showCenterText">
                <slot></slot>
            </g>
        </svg>
        <div class="pie-chart-toolbox-container padding-horizontal" v-if="showSelectedItemInfo">
            <div class="pie-chart-toolbox">
                <f7-link class="pie-chart-toolbox-button" :class="{ 'disabled': !!skeleton || !validItems || validItems.length <= 1 }" @click="switchSelectedItem(1)">
                    <f7-icon f7="arrow_left"></f7-icon>
                </f7-link>

                <div class="pie-chart-toolbox-info">
                    <p v-if="selectedItem">
                        <f7-chip class="chip-placeholder" outline v-if="skeleton">
                            <span class="skeleton-text">Percent</span>
                        </f7-chip>
                        <f7-chip outline
                                 :text="selectedItem.displayPercent"
                                 :style="getColorStyle(selectedItem ? selectedItem.color : '', '--f7-chip-outline-border-color')"
                                 v-else-if="!skeleton"></f7-chip>
                    </p>
                    <p v-else-if="!validItems || !validItems.length">
                        <f7-chip outline text="---"></f7-chip>
                    </p>
                    <f7-link class="pie-chart-selected-item-info" :no-link-class="!enableClickItem" v-if="selectedItem" @click="clickItem(selectedItem)">
                        <span class="skeleton-text" v-if="skeleton">Name</span>
                        <span v-else-if="!skeleton && selectedItem.name">{{ selectedItem.name }}</span>
                        <span class="skeleton-text" v-if="skeleton">Value</span>
                        <span v-else-if="!skeleton && showValue" :style="getColorStyle(selectedItem ? selectedItem.color : '')">{{ selectedItem.displayValue }}</span>
                        <f7-icon class="item-navigate-icon" f7="chevron_right" v-if="enableClickItem"></f7-icon>
                    </f7-link>
                    <f7-link :no-link-class="true" v-else-if="!validItems || !validItems.length">
                        {{ $t('No transaction data') }}
                    </f7-link>
                </div>

                <f7-link class="pie-chart-toolbox-button" :class="{ 'disabled': !!skeleton || !validItems || validItems.length <= 1 }" @click="switchSelectedItem(-1)">
                    <f7-icon f7="arrow_right"></f7-icon>
                </f7-link>
            </div>
        </div>
    </div>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import { DEFAULT_ICON_COLOR, DEFAULT_CHART_COLORS } from '@/consts/color.ts';
import { formatPercent } from '@/lib/numeral.js';

export default {
    props: [
        'skeleton',
        'items',
        'nameField',
        'valueField',
        'percentField',
        'colorField',
        'hiddenField',
        'minValidPercent',
        'defaultCurrency',
        'showValue',
        'showCenterText',
        'showSelectedItemInfo',
        'enableClickItem',
        'centerTextBackground',
    ],
    emits: [
        'click'
    ],
    data: function () {
        const diameter = 100;

        return {
            diameter: diameter,
            circumference: diameter * Math.PI,
            selectedIndex: 0
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore),
        validItems: function () {
            let totalValidValue = 0;

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (item[this.valueField] && item[this.valueField] > 0 && (!this.hiddenField || !item[this.hiddenField])) {
                    totalValidValue += item[this.valueField];
                }
            }

            const validItems = [];

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (item[this.valueField] && item[this.valueField] > 0 &&
                    (!this.hiddenField || !item[this.hiddenField]) &&
                    (!this.minValidPercent || item[this.valueField] / totalValidValue > this.minValidPercent)) {
                    const finalItem = {
                        name: item[this.nameField],
                        value: item[this.valueField],
                        percent: (item[this.percentField] > 0 || item[this.percentField] === 0 || item[this.percentField] === '0') ? item[this.percentField] : (item[this.valueField] / totalValidValue * 100),
                        actualPercent: item[this.valueField] / totalValidValue,
                        color: item[this.colorField] ? item[this.colorField] : DEFAULT_CHART_COLORS[validItems.length % DEFAULT_CHART_COLORS.length],
                        sourceItem: item
                    };

                    finalItem.displayPercent = formatPercent(finalItem.percent, 2, '&lt;0.01');
                    finalItem.displayValue = this.getDisplayCurrency(finalItem.value, this.defaultCurrency);

                    validItems.push(finalItem);
                }
            }

            return validItems;
        },
        totalValidValue: function () {
            let totalValidValue = 0;

            for (let i = 0; i < this.validItems.length; i++) {
                totalValidValue += this.validItems[i].value;
            }

            return totalValidValue;
        },
        itemCommonDashOffset: function () {
            if (this.totalValidValue <= 0) {
                return 0;
            }

            let offset = 0;

            for (let i = 0; i < Math.min(this.selectedIndex + 1, this.validItems.length); i++) {
                const item = this.validItems[i];

                if (item.actualPercent > 0) {
                    if (i === this.selectedIndex) {
                        offset += -this.circumference * (1 - item.actualPercent) / 2;
                    } else {
                        offset += -this.circumference * (1 - item.actualPercent);
                    }
                }
            }

            return offset;
        },
        selectedItem: function () {
            if (!this.validItems || !this.validItems.length) {
                return null;
            }

            let selectedIndex = this.selectedIndex;

            if (selectedIndex < 0 || selectedIndex >= this.validItems.length) {
                selectedIndex = 0;
            }

            return this.validItems[selectedIndex];
        }
    },
    watch: {
        'items': {
            handler() {
                this.selectedIndex = 0;
            },
            deep: true
        }
    },
    methods: {
        switchSelectedIndex: function (index) {
            this.selectedIndex = index;
        },
        switchSelectedItem: function (offset) {
            let newSelectedIndex = this.selectedIndex + offset;

            while (newSelectedIndex < 0) {
                newSelectedIndex += this.validItems.length;
            }

            this.selectedIndex = newSelectedIndex % this.validItems.length;
        },
        clickItem: function (item) {
            if (this.enableClickItem) {
                this.$emit('click', item.sourceItem);
            }
        },
        getColor: function (color) {
            if (color && color !== DEFAULT_ICON_COLOR) {
                color = '#' + color;
            } else {
                color = 'var(--default-icon-color)';
            }

            return color;
        },
        getColorStyle: function (color, additionalFieldName) {
            const ret = {
                color: this.getColor(color)
            };

            if (additionalFieldName) {
                ret[additionalFieldName] = ret.color;
            }

            return ret;
        },
        getItemStrokeDash(item) {
            const length = item.actualPercent * this.circumference;
            return `${length} ${this.circumference - length}`;
        },
        getItemDashOffset(item, items, offset) {
            let allPreviousPercent = 0;

            for (let i = 0; i < items.length; i++) {
                const curItem = items[i];

                if (curItem === item) {
                    break;
                }

                allPreviousPercent += curItem.actualPercent;
            }

            if (offset) {
                offset += this.circumference / 4;
            } else {
                offset = this.circumference / 4;
            }

            if (allPreviousPercent <= 0) {
                return offset;
            }

            const allPreviousLength = allPreviousPercent * this.circumference;
            return this.circumference - allPreviousLength + offset;
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        }
    }
}
</script>

<style scoped>
.pie-chart-container {
    width: 100%;
    height: 100%;
}

.pie-chart {
    margin: 0 24px 0 24px;
}

.pie-chart-toolbox-container {
    margin-top: 16px;
    padding-bottom: 16px;
}

.pie-chart-toolbox {
    display: inline-flex;
    width: 100%;
    justify-content: space-between;
}

.pie-chart-toolbox-info {
    --f7-chip-height: var(--ebk-pie-chart-toolbox-percentage-height);
    --f7-chip-font-size: var(--ebk-pie-chart-toolbox-percentage-font-size);
    font-size: var(--ebk-pie-chart-toolbox-text-font-size);
    align-self: center;
}

.pie-chart-toolbox-info p {
    text-align: center;
    margin: 0 0 4px 0;
}

.pie-chart-toolbox-info a {
    color: var(--f7-text-color);
}

.pie-chart-toolbox-info a > span + span {
    padding-left: 8px;
}

.pie-chart-toolbox-info .item-navigate-icon {
    color: rgba(0, 0, 0, 0.2);
    font-size: 18px;
    font-weight: bold;
    padding-left: 4px;
}

.pie-chart-toolbox-button {
    color: var(--f7-text-color);
}

.pie-chart-selected-item-info {
    display: inline-block;
    text-align: center;
}

.pie-chart-background {
    fill: #f0f0f0;
}

.dark .pie-chart-background {
    fill: #181818;
}

.pie-chart-text-background {
    --pie-chart-text-background: #7f2020;
    fill: var(--pie-chart-text-background);
}

.pie-chart-text-group {
    font-size: 0.7em;
    -moz-transform: translateY(-1em);
    -ms-transform: translateY(-1em);
    -webkit-transform: translateY(-1em);
    transform: translateY(-1em);
}
</style>
