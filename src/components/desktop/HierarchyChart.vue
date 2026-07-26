<template>
    <v-chart autoresize :class="finalClass" :option="chartOptions"
             @click="clickItem" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import { itemAndIndex } from '@/core/base.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import type { ColorValue, ColorStyleValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { type AxisChartSourceDataItem, ChartValueType } from '@/core/chart.ts';

import { isArray, isString, isNumber } from '@/lib/common.ts';
import { BIG_DECIMAL_ZERO, parseBigDecimal, isBigDecimal } from '@/lib/numeral.ts';
import { getDisplayColor } from '@/lib/color.ts';

export type HierarchyChartDisplayType = 'treemap' | 'sunburst';

interface HierarchyData {
    data: HierarchyDataItem[];
    totalAmount: BigDecimal;
}

interface HierarchyDataItem {
    name: string;
    value: number; // only used for echarts calculation, the actual value is originalValue
    originalValue: string;
    parentName?: string;
    parentOrginalValue?: string;
    categoryIndex?: number;
    seriesIndex?: number;
    children?: HierarchyDataItem[];
    itemStyle: {
        color: ColorStyleValue;
    };
}

const props = defineProps<{
    class?: string;
    skeleton?: boolean;
    type: HierarchyChartDisplayType;
    showValue?: boolean;
    enableClickItem?: boolean;
    categoryTypeName: string;
    allCategoryNames: string[];
    items: AxisChartSourceDataItem[];
    valueType: ChartValueType;
    translateName?: boolean;
    useCustomColor?: boolean;
    defaultCurrency?: string;
}>();

const emit = defineEmits<{
    (e: 'click', categoryIndex: number, seriesIndex: number): void;
}>();

const theme = useTheme();

const {
    tt,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping,
    formatPercentToLocalizedNumerals,
    formatChartValueToLocalizedNumerals
} = useI18n();

const settingsStore = useSettingsStore();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const chartColors = computed<ColorValue[]>(() => settingsStore.chartColorList);

const finalClass = computed<string>(() => {
    let finalClass = '';

    if (props.skeleton) {
        finalClass += 'transition-in';
    }

    if (props.class) {
        finalClass += ` ${props.class}`;
    } else {
        finalClass += ' hierarchy-chart-container';
    }

    return finalClass;
});

const hierarchyData = computed<HierarchyData>(() => {
    const ret: HierarchyDataItem[] = [];
    let totalAmount: BigDecimal = BIG_DECIMAL_ZERO;

    for (const [item, seriesIndex] of itemAndIndex(props.items)) {
        if (item.hidden) {
            continue;
        }

        if (!isArray(item.values)) {
            continue;
        }

        const color: ColorStyleValue = getDisplayColor(props.useCustomColor && item.color ? item.color : chartColors.value[seriesIndex % chartColors.value.length]);

        const hierarchyItem: HierarchyDataItem = {
            name: props.translateName ? tt(item.name) : item.name,
            value: 0,
            originalValue: '0',
            children: [],
            itemStyle: {
                color: color
            }
        };

        const allAmounts: BigDecimal[] = item.values;
        let childrenTotalAmount: BigDecimal = BIG_DECIMAL_ZERO;

        for (const [amount, categoryIndex] of itemAndIndex(allAmounts)) {
            childrenTotalAmount = childrenTotalAmount.add(amount);
            hierarchyItem.children?.push({
                name: props.allCategoryNames[categoryIndex] ?? '',
                value: amount.toDoubleNumber(),
                originalValue: amount.toString(),
                categoryIndex: categoryIndex,
                seriesIndex: seriesIndex,
                itemStyle: {
                    color: color
                }
            });
        }

        hierarchyItem.value = childrenTotalAmount.toDoubleNumber();
        hierarchyItem.originalValue = childrenTotalAmount.toString();

        for (const child of hierarchyItem.children ?? []) {
            child.parentName = hierarchyItem.name;
            child.parentOrginalValue = hierarchyItem.originalValue;
        }

        totalAmount = totalAmount.add(childrenTotalAmount);
        ret.push(hierarchyItem);
    }

    for (const item of ret) {
        item.parentOrginalValue = totalAmount.toString();
    }

    const hierarchyData: HierarchyData = {
        data: ret,
        totalAmount: totalAmount
    };

    return hierarchyData;
});

const chartOptions = computed<object>(() => {
    const seriesOptions: Record<string, unknown> = {
        type: props.type,
        width: '100%',
        height: '100%',
        right: 20,
        top: 0,
        bottom: 20,
        data: hierarchyData.value.data,
        levels: [
            {
                itemStyle: {
                    gapWidth: 2
                }
            },
            {
                itemStyle: {
                    gapWidth: 1
                }
            }
        ],
        animation: !props.skeleton,
        nodeClick: false
    };

    if (props.type === 'treemap') {
        seriesOptions['breadcrumb'] = {
            show: false
        };
    } if (props.type === 'sunburst') {
        seriesOptions['radius'] = [60, '95%'];
        seriesOptions['itemStyle'] = {
            borderRadius: 7,
            borderWidth: 2
        };
    }

    return {
        tooltip: {
            backgroundColor: isDarkMode.value ? '#333' : '#fff',
            borderColor: isDarkMode.value ? '#333' : '#fff',
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams) => {
                if (!props.showValue || !params.name) {
                    return '';
                }

                const dataItem = params.data as HierarchyDataItem;
                const rootValue: BigDecimal = hierarchyData.value.totalAmount;
                const parentName: string | undefined = dataItem.parentName;
                const parentValue: BigDecimal | undefined = isString(dataItem.parentOrginalValue) ? parseBigDecimal(dataItem.parentOrginalValue) : undefined;
                const parentDisplayValue: string | undefined = isBigDecimal(parentValue) ? formatChartValueToLocalizedNumerals(parentValue, props.valueType, props.defaultCurrency) : undefined;
                const parentDisplayPercent: string | undefined = isBigDecimal(parentValue) && isBigDecimal(rootValue) && rootValue.isPositive() ? formatPercentToLocalizedNumerals(parentValue.divide(rootValue).multiply(100).toDoubleNumber(), 2, '<0.01') : undefined;

                const name = params.name;
                const displayValue = isString(dataItem.originalValue) ? formatChartValueToLocalizedNumerals(parseBigDecimal(dataItem.originalValue), props.valueType, props.defaultCurrency) : '';
                const displayPercent = isString(dataItem.originalValue) && isBigDecimal(parentValue) && parentValue.isPositive() ? formatPercentToLocalizedNumerals(parseBigDecimal(dataItem.originalValue).divide(parentValue).multiply(100).toDoubleNumber(), 2, '<0.01') : undefined;

                let tooltip = `<tr><td><span class="chart-pointer" style="background-color: ${params.color}"></span><span>${name}</span></td>`
                    + `<td><span class="ms-5">${displayValue}</span>`
                    + (isString(displayPercent) ? `<span class="ms-1">(${displayPercent})</span>` : '')
                    + `</td></tr>`;

                if (isString(parentName) && isString(parentDisplayValue) && parentValue?.notEquals(rootValue)) {
                    tooltip = `<tr><td><span class="chart-pointer" style="background-color: ${params.color}"></span><span>${parentName}</span></td>`
                        + `<td><span class="ms-5">${parentDisplayValue}</span>`
                        + (isString(parentDisplayPercent) ? `<span class="ms-1">(${parentDisplayPercent})</span>` : '')
                        + `</td></tr>`
                        + tooltip;
                }

                tooltip = `<table class="chart-tooltip-table"><tbody>` + tooltip + `</tbody></table>`;
                return tooltip;
            }
        },
        series: [ seriesOptions ]
    };
});

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series' || !e.data) {
        return;
    }

    const dataItem = e.data as HierarchyDataItem;

    if (isNumber(dataItem.categoryIndex) && isNumber(dataItem.seriesIndex)) {
        emit('click', dataItem.categoryIndex, dataItem.seriesIndex);
    }
}

function exportData(): { headers: string[], data: string[][] } {
    const headers: string[] = [];
    const data: string[][] = [];

    headers.push(props.categoryTypeName);

    for (const categoryName of props.allCategoryNames) {
        headers.push(categoryName);
    }

    for (const item of hierarchyData.value.data) {
        const row: string[] = [];
        row.push(item.name);

        for (const child of item.children ?? []) {
            if (props.valueType === ChartValueType.Amount) {
                row.push(formatAmountToWesternArabicNumeralsWithoutDigitGrouping(parseBigDecimal(child.value), props.defaultCurrency));
            } else {
                row.push(formatBigDecimalToWesternArabicNumeralsWithoutDigitGrouping(parseBigDecimal(child.value)));
            }
        }

        data.push(row);
    }

    return {
        headers: headers,
        data: data
    };
}

defineExpose({
    exportData
});
</script>

<style scoped>
.hierarchy-chart-container {
    width: 100%;
    height: 560px;
    margin-top: 10px;
}

@media (min-width: 600px) {
    .hierarchy-chart-container {
        height: 630px;
    }
}
</style>
