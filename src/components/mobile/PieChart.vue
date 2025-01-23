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
                        {{ tt('No transaction data') }}
                    </f7-link>
                </div>

                <f7-link class="pie-chart-toolbox-button" :class="{ 'disabled': !!skeleton || !validItems || validItems.length <= 1 }" @click="switchSelectedItem(-1)">
                    <f7-icon f7="arrow_right"></f7-icon>
                </f7-link>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type { ColorValue } from '@/core/color.ts';
import { DEFAULT_ICON_COLOR, DEFAULT_CHART_COLORS } from '@/consts/color.ts';

import { isNumber } from '@/lib/common.ts';
import { formatPercent } from '@/lib/numeral.ts';

interface MobilePieChartDataItem {
    name: string;
    value: number;
    percent: number;
    actualPercent: number;
    color: ColorValue;
    sourceItem: Record<string, unknown>;
    displayPercent?: string;
    displayValue?: string;
}

const props = defineProps<{
    skeleton?: boolean;
    items: Record<string, unknown>[];
    nameField: string;
    valueField: string;
    percentField?: string;
    colorField?: string;
    hiddenField?: string;
    minValidPercent?: number;
    defaultCurrency?: string;
    showValue?: boolean;
    showCenterText?: boolean;
    showSelectedItemInfo?: boolean;
    enableClickItem?: boolean;
    centerTextBackground?: ColorValue;
}>();

const emit = defineEmits<{
    (e: 'click', value: Record<string, unknown>): void;
}>();

const { tt, formatAmountWithCurrency } = useI18n();

const diameter: number = 100;
const circumference: number = diameter * Math.PI;

const selectedIndex = ref<number>(0);

const validItems = computed<MobilePieChartDataItem[]>(() => {
    let totalValidValue = 0;

    for (let i = 0; i < props.items.length; i++) {
        const item = props.items[i];
        const value = item[props.valueField];

        if (isNumber(value) && value > 0 && (!props.hiddenField || !item[props.hiddenField])) {
            totalValidValue += value;
        }
    }

    const validItems: MobilePieChartDataItem[] = [];

    for (let i = 0; i < props.items.length; i++) {
        const item = props.items[i];
        const value = item[props.valueField];
        const percent = props.percentField ? item[props.percentField] : -1;

        if (isNumber(value) && value > 0 &&
            (!props.hiddenField || !item[props.hiddenField]) &&
            (!props.minValidPercent || value / totalValidValue > props.minValidPercent)) {
            const finalItem: MobilePieChartDataItem = {
                name: item[props.nameField] as string,
                value: value,
                percent: (isNumber(percent) && percent >= 0) ? percent : (value / totalValidValue * 100),
                actualPercent: value / totalValidValue,
                color: (props.colorField && item[props.colorField]) ? item[props.colorField] as ColorValue : DEFAULT_CHART_COLORS[validItems.length % DEFAULT_CHART_COLORS.length],
                sourceItem: item
            };

            finalItem.displayPercent = formatPercent(finalItem.percent, 2, '&lt;0.01');
            finalItem.displayValue = formatAmountWithCurrency(finalItem.value, props.defaultCurrency) as string;

            validItems.push(finalItem);
        }
    }

    return validItems;
});

const totalValidValue = computed<number>(() => {
    let totalValidValue = 0;

    for (let i = 0; i < validItems.value.length; i++) {
        totalValidValue += validItems.value[i].value;
    }

    return totalValidValue;
});

const itemCommonDashOffset = computed<number>(() => {
    if (totalValidValue.value <= 0) {
        return 0;
    }

    let offset = 0;

    for (let i = 0; i < Math.min(selectedIndex.value + 1, validItems.value.length); i++) {
        const item = validItems.value[i];

        if (item.actualPercent > 0) {
            if (i === selectedIndex.value) {
                offset += -circumference * (1 - item.actualPercent) / 2;
            } else {
                offset += -circumference * (1 - item.actualPercent);
            }
        }
    }

    return offset;
});

const selectedItem = computed<MobilePieChartDataItem | null>(() => {
    if (!validItems.value || !validItems.value.length) {
        return null;
    }

    let index = selectedIndex.value;

    if (index < 0 || index >= validItems.value.length) {
        index = 0;
    }

    return validItems.value[index];
});

watch(() => props.items, () => {
    selectedIndex.value = 0;
});

function switchSelectedIndex(index: number): void {
    selectedIndex.value = index;
}

function switchSelectedItem(offset: number): void {
    let newSelectedIndex = selectedIndex.value + offset;

    while (newSelectedIndex < 0) {
        newSelectedIndex += validItems.value.length;
    }

    selectedIndex.value = newSelectedIndex % validItems.value.length;
}

function clickItem(item: MobilePieChartDataItem): void {
    if (props.enableClickItem) {
        emit('click', item.sourceItem);
    }
}

function getColor(color: ColorValue): ColorValue {
    if (color && color !== DEFAULT_ICON_COLOR) {
        color = '#' + color;
    } else {
        color = 'var(--default-icon-color)';
    }

    return color;
}

function getColorStyle(color: ColorValue, additionalFieldName?: string): Record<string, string> {
    const ret: Record<string, string> = {
        color: getColor(color)
    };

    if (additionalFieldName) {
        ret[additionalFieldName] = ret.color;
    }

    return ret;
}

function getItemStrokeDash(item: MobilePieChartDataItem): string {
    const length = item.actualPercent * circumference;
    return `${length} ${circumference - length}`;
}

function getItemDashOffset(item: MobilePieChartDataItem, items: MobilePieChartDataItem[], offset?: number): number {
    let allPreviousPercent = 0;

    for (let i = 0; i < items.length; i++) {
        const curItem = items[i];

        if (curItem === item) {
            break;
        }

        allPreviousPercent += curItem.actualPercent;
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
