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
                    :stroke-dashoffset="item | itemDashOffset(validItems, circumference, firstValidItemDashOffset)">
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
    </div>
</template>

<script>
export default {
    props: [
        'items',
        'nameField',
        'valueField',
        'colorField',
        'minValidPercent',
        'showCenterText',
        'centerTextBackground',
    ],
    data: function () {
        const diameter = 100;

        return {
            diameter: diameter,
            circumference: diameter * Math.PI
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
        firstValidItemDashOffset: function () {
            if (this.totalValidValue <= 0) {
                return 0;
            }

            for (let i = 0; i < this.validItems.length; i++) {
                const item = this.validItems[i];

                if (item.percent > 0) {
                    return -this.circumference * (1 - item.percent) / 2;
                }
            }

            return 0;
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
    margin: 24px;
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
