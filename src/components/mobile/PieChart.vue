<template>
    <div class="pie-chart-container">
        <svg class="pie-chart" :viewBox="`${-diameter} ${-diameter} ${diameter * 2} ${diameter * 2}`">
            <circle class="pie-chart-background" cx="0" cy="0" :r="diameter"></circle>
            <circle class="pie-chart-item"
                    v-for="(item, idx) in validItems" :key="idx"
                    fill="transparent"
                    cx="0" cy="0"
                    :r="diameter / 2"
                    :stroke="item.color | defaultIconColor('var(--default-icon-color)')"
                    :stroke-width="diameter"
                    :stroke-dasharray="item | itemStrokeDash(circumference)"
                    :stroke-dashoffset="item | itemDashOffset(validItems, circumference, itemCommonDashOffset)">
            </circle>
            <circle class="pie-chart-text-background"
                    cx="0" cy="0"
                    :style="{ '--pie-chart-text-background': centerTextBackground ? centerTextBackground : 'var(--f7-theme-color)' }"
                    :r="diameter / 2.5"
                    v-if="showCenterText"/>
            <g class="pie-chart-text-group" v-if="showCenterText">
                <slot></slot>
            </g>
        </svg>
        <div class="pie-chart-toolbox-container padding-horizontal" v-if="showSelectedItemInfo">
            <div class="pie-chart-toolbox">
                <f7-link class="pie-chart-toolbox-button" :class="{ 'disabled': !!skeleton || !validItems || !validItems.length }" @click="switchSelectedItem(1)">
                    <f7-icon f7="arrow_left"></f7-icon>
                </f7-link>

                <div class="pie-chart-toolbox-info">
                    <p v-if="selectedItem">
                        <f7-chip class="chip-placeholder" outline v-if="skeleton">
                            <span class="skeleton-text">Percent</span>
                        </f7-chip>
                        <f7-chip outline
                                 :text="(selectedItem.percent * 100) | percent(2, '&lt;0.01')"
                                 :style="(selectedItem ? selectedItem.color : '') | iconStyle('default', 'var(--default-icon-color)', '--f7-chip-outline-border-color')"
                                 v-else-if="!skeleton"></f7-chip>
                    </p>
                    <p v-else-if="!validItems || !validItems.length">
                        <f7-chip outline text="0%"></f7-chip>
                    </p>
                    <p v-if="selectedItem">
                        <span class="skeleton-text" v-if="skeleton">Name</span>
                        <span v-else-if="!skeleton && selectedItem.name">{{ selectedItem.name }}</span>
                        <span class="skeleton-text" v-if="skeleton">Value</span>
                        <span v-else-if="!skeleton" :style="(selectedItem ? selectedItem.color : '') | iconStyle('default', 'var(--default-icon-color)')">{{ selectedItem.value | currency(defaultCurrency) }}</span>
                    </p>
                    <p v-else-if="!validItems || !validItems.length">
                        {{ $t('No transaction data') }}
                    </p>
                </div>

                <f7-link class="pie-chart-toolbox-button" :class="{ 'disabled': !!skeleton || !validItems || !validItems.length }" @click="switchSelectedItem(-1)">
                    <f7-icon f7="arrow_right"></f7-icon>
                </f7-link>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: [
        'skeleton',
        'items',
        'nameField',
        'valueField',
        'colorField',
        'minValidPercent',
        'defaultCurrency',
        'showCenterText',
        'showSelectedItemInfo',
        'centerTextBackground',
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
        validItems: function () {
            let totalValidValue = 0;

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (item[this.valueField] && item[this.valueField] > 0) {
                    totalValidValue += item[this.valueField];
                }
            }

            const validItems = [];

            for (let i = 0; i < this.items.length; i++) {
                const item = this.items[i];

                if (item[this.valueField] && item[this.valueField] > 0 &&
                    (!this.minValidPercent || item[this.valueField] / totalValidValue > this.minValidPercent)) {
                    validItems.push({
                        name: item[this.nameField],
                        value: item[this.valueField],
                        percent: item[this.valueField] / totalValidValue,
                        color: item[this.colorField] ? item[this.colorField] : 'c8c8c8'
                    });
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

                if (item.percent > 0) {
                    if (i === this.selectedIndex) {
                        offset += -this.circumference * (1 - item.percent) / 2;
                    } else {
                        offset += -this.circumference * (1 - item.percent);
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
        'items': function () {
            this.selectedIndex = 0;
        }
    },
    methods: {
        switchSelectedItem: function (offset) {
            let newSelectedIndex = this.selectedIndex + offset;

            while (newSelectedIndex < 0) {
                newSelectedIndex += this.validItems.length;
            }

            this.selectedIndex = newSelectedIndex % this.validItems.length;
        }
    },
    filters: {
        itemStrokeDash(item, circumference) {
            const length = item.percent * circumference;
            return `${length} ${circumference - length}`;
        },
        itemDashOffset(item, items, circumference, offset) {
            let allPreviousPercent = 0;

            for (let i = 0; i < items.length; i++) {
                const curItem = items[i];

                if (curItem === item) {
                    break;
                }

                allPreviousPercent += curItem.percent;
            }

            if (offset) {
                offset += circumference / 4;
            } else {
                offset = circumference / 4;
            }

            if (allPreviousPercent <= 0) {
                return offset;
            }

            const allPreviousLength = allPreviousPercent * circumference;
            return circumference - allPreviousLength + offset;
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
    margin: 24px 24px 0 24px;
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
    --f7-chip-height: 30px;
    --f7-chip-font-size: 18px;
    font-size: 16px;
    align-self: center;
}

.pie-chart-toolbox-info p {
    text-align: center;
    margin: 0 0 4px 0;
}

.pie-chart-toolbox-info p > span {
    padding-right: 8px;
}

.pie-chart-toolbox-info p > span:last-child {
    padding-right: 0;
}

.pie-chart-toolbox-button {
    color: var(--f7-text-color);
}

.pie-chart-background {
    fill: #f0f0f0;
}

.theme-dark .pie-chart-background {
    fill: #181818;
}

.pie-chart-text-background {
    --pie-chart-text-background: var(--f7-theme-color);
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