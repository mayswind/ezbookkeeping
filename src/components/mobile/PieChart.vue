<template>
    <div class="pie-chart-container">
        <svg class="pie-chart" :viewBox="`${-diameter} ${-diameter} ${diameter * 2} ${diameter * 2}`">
            <circle class="pie-chart-background" cx="0" cy="0" :r="diameter"></circle>

            <circle class="pie-chart-item"
                    fill="transparent"
                    cx="0" cy="0"
                    :r="diameter / 2"
                    :stroke="item.color"
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
                    <f7-icon class="icon-with-direction" f7="arrow_left"></f7-icon>
                </f7-link>

                <div class="pie-chart-toolbox-info">
                    <p v-if="selectedItem">
                        <f7-chip class="chip-placeholder" outline v-if="skeleton">
                            <span class="skeleton-text">Percent</span>
                        </f7-chip>
                        <f7-chip outline
                                 :text="selectedItem.displayPercent"
                                 :style="getColorStyle(selectedItem?.color, '--f7-chip-outline-border-color')"
                                 v-else-if="!skeleton"></f7-chip>
                    </p>
                    <p v-else-if="!validItems || !validItems.length">
                        <f7-chip outline text="---"></f7-chip>
                    </p>
                    <f7-link class="pie-chart-selected-item-info" :no-link-class="!enableClickItem" v-if="selectedItem" @click="clickItem(selectedItem)">
                        <span class="skeleton-text" v-if="skeleton">Name</span>
                        <span v-else-if="!skeleton && selectedItem.displayName">{{ selectedItem.displayName }}</span>
                        <span class="skeleton-text" v-if="skeleton">Value</span>
                        <span v-else-if="!skeleton && showValue" :style="getColorStyle(selectedItem?.color)">{{ selectedItem.displayValue }}</span>
                        <f7-icon class="item-navigate-icon icon-with-direction" f7="chevron_right" v-if="enableClickItem"></f7-icon>
                    </f7-link>
                    <f7-link :no-link-class="true" v-else-if="!validItems || !validItems.length">
                        {{ tt('No transaction data') }}
                    </f7-link>
                </div>

                <f7-link class="pie-chart-toolbox-button" :class="{ 'disabled': !!skeleton || !validItems || validItems.length <= 1 }" @click="switchSelectedItem(-1)">
                    <f7-icon class="icon-with-direction" f7="arrow_right"></f7-icon>
                </f7-link>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonPieChartDataItem, type CommonPieChartProps, usePieChartBase } from '@/components/base/PieChartBase.ts'

import type { ColorStyleValue } from '@/core/color.ts';

interface MobilePieChartProps extends CommonPieChartProps {
    showCenterText?: boolean;
    showSelectedItemInfo?: boolean;
    centerTextBackground?: ColorStyleValue;
}

const props = defineProps<MobilePieChartProps>();

const emit = defineEmits<{
    (e: 'click', value: Record<string, unknown>): void;
}>();

const { tt } = useI18n();
const { selectedIndex, validItems } = usePieChartBase(props);

const diameter: number = 100;
const circumference: number = diameter * Math.PI;

const totalValidValue = computed<number>(() => {
    let totalValidValue = 0;

    for (const item of validItems.value) {
        totalValidValue += item.value;
    }

    return totalValidValue;
});

const itemCommonDashOffset = computed<number>(() => {
    if (totalValidValue.value <= 0) {
        return 0;
    }

    let offset = 0;

    for (let i = 0; i < Math.min(selectedIndex.value + 1, validItems.value.length); i++) {
        const item = validItems.value[i] as CommonPieChartDataItem;

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

function getColorStyle(color: ColorStyleValue, additionalFieldName?: string): Record<string, string> {
    const ret: Record<string, string> = {
        color: color
    };

    if (additionalFieldName) {
        ret[additionalFieldName] = color;
    }

    return ret;
}

function getItemStrokeDash(item: CommonPieChartDataItem): string {
    const length = item.actualPercent * circumference;
    return `${length} ${circumference - length}`;
}

function getItemDashOffset(item: CommonPieChartDataItem, items: CommonPieChartDataItem[], offset?: number): number {
    let allPreviousPercent = 0;

    for (const curItem of items) {
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

const selectedItem = computed<CommonPieChartDataItem | null>(() => {
    if (!validItems.value || !validItems.value.length) {
        return null;
    }

    let index = selectedIndex.value;

    if (index < 0 || index >= validItems.value.length) {
        index = 0;
    }

    return validItems.value[index] ?? null;
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

function clickItem(item: CommonPieChartDataItem): void {
    if (props.enableClickItem) {
        emit('click', item.sourceItem);
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
    padding-inline-start: 8px;
}

.pie-chart-toolbox-info .item-navigate-icon {
    color: rgba(0, 0, 0, 0.2);
    font-size: 18px;
    font-weight: bold;
    padding-inline-start: 4px;
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
