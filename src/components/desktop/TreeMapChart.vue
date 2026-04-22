<template>
    <v-chart autoresize :class="finalClass" :option="chartOptions" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { itemAndIndex } from '@/core/base.ts';
import type { ColorValue, ColorStyleValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { DEFAULT_CHART_COLORS } from '@/consts/color.ts';

import { isArray, isString, isNumber } from '@/lib/common.ts';
import { getDisplayColor } from '@/lib/color.ts';

interface TreeMapDataItem {
    name: string;
    value: number;
    children?: TreeMapDataItem[];
    itemStyle: {
        color: ColorStyleValue;
    };
}

const props = defineProps<{
    class?: string;
    skeleton?: boolean;
    showValue?: boolean;
    categoryTypeName: string;
    allCategoryNames: string[];
    items: Record<string, unknown>[];
    nameField: string;
    valuesField: string;
    colorField?: string;
    hiddenField?: string;
    translateName?: boolean;
    amountValue?: boolean;
    percentValue?: boolean;
    defaultCurrency?: string;
}>();

const theme = useTheme();

const {
    tt,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatNumberToLocalizedNumerals,
    formatPercentToLocalizedNumerals
} = useI18n();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const finalClass = computed<string>(() => {
    let finalClass = '';

    if (props.skeleton) {
        finalClass += 'transition-in';
    }

    if (props.class) {
        finalClass += ` ${props.class}`;
    } else {
        finalClass += ' treemap-chart-container';
    }

    return finalClass;
});

const treeMapData = computed<TreeMapDataItem[]>(() => {
    const ret: TreeMapDataItem[] = [];

    for (const [item, seriesIndex] of itemAndIndex(props.items)) {
        if (props.hiddenField && item[props.hiddenField]) {
            continue;
        }

        if (!isArray(item[props.valuesField])) {
            continue;
        }

        const color: ColorStyleValue = getDisplayColor((props.colorField && item[props.colorField]) ? item[props.colorField] as ColorValue : DEFAULT_CHART_COLORS[seriesIndex % DEFAULT_CHART_COLORS.length]);

        const treeMapItem: TreeMapDataItem = {
            name: getItemName(item[props.nameField] as string),
            value: 0,
            children: [],
            itemStyle: {
                color: color
            }
        };

        const allAmounts: number[] = item[props.valuesField] as number[];

        for (const [amount, categoryIndex] of itemAndIndex(allAmounts)) {
            treeMapItem.value += amount;
            treeMapItem.children?.push({
                name: props.allCategoryNames[categoryIndex] ?? '',
                value: amount,
                itemStyle: {
                    color: color
                }
            });
        }

        ret.push(treeMapItem);
    }

    return ret;
});

const chartOptions = computed<object>(() => {
    return {
        tooltip: {
            backgroundColor: isDarkMode.value ? '#333' : '#fff',
            borderColor: isDarkMode.value ? '#333' : '#fff',
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams & { treePathInfo: { name: string, value: number }[] }) => {
                if (!props.showValue || !params.name) {
                    return '';
                }

                const rootValue = params.treePathInfo.length > 0 ? params.treePathInfo[0]?.value : 0;
                const parentName = params.treePathInfo.length > 1 ? params.treePathInfo[params.treePathInfo.length - 2]?.name : undefined;
                const parentValue = params.treePathInfo.length > 1 ? params.treePathInfo[params.treePathInfo.length - 2]?.value : undefined;
                const parentDisplayValue = isNumber(parentValue) ? getDisplayValue(parentValue) : undefined;
                const parentPercent = isNumber(parentValue) && isNumber(rootValue) && rootValue > 0 ? formatPercentToLocalizedNumerals(100.0 * parentValue / rootValue, 2, '<0.01') : undefined;

                const name = params.name;
                const displayValue = isNumber(params.value) ? getDisplayValue(params.value) : '';
                const percent = isNumber(params.value) && isNumber(parentValue) && parentValue > 0 ? formatPercentToLocalizedNumerals(100.0 * params.value / parentValue, 2, '<0.01') : undefined;


                let tooltip = `<tr><td><span class="chart-pointer" style="background-color: ${params.color}"></span><span>${name}</span></td>`
                    + `<td><span class="ms-5">${displayValue}</span>`
                    + (isString(percent) ? `<span class="ms-1">(${percent})</span>` : '')
                    + `</td></tr>`;

                if (isString(parentName) && isString(parentDisplayValue) && parentValue !== rootValue) {
                    tooltip = `<tr><td><span class="chart-pointer" style="background-color: ${params.color}"></span><span>${parentName}</span></td>`
                        + `<td><span class="ms-5">${parentDisplayValue}</span>`
                        + (isString(parentPercent) ? `<span class="ms-1">(${parentPercent})</span>` : '')
                        + `</td></tr>`
                        + tooltip;
                }

                tooltip = `<table class="chart-tooltip-table"><tbody>` + tooltip + `</tbody></table>`;
                return tooltip;
            }
        },
        series: [
            {
                type: 'treemap',
                animation: !props.skeleton,
                nodeClick: false,
                breadcrumb: {
                    show: false
                },
                data: treeMapData.value,
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
                ]
            }
        ]
    };
});

function getItemName(name: string): string {
    return props.translateName ? tt(name) : name;
}

function getDisplayValue(value: number): string {
    if (props.percentValue) {
        return formatPercentToLocalizedNumerals(value, 2, '<0.01');
    }

    if (props.amountValue) {
        return formatAmountToLocalizedNumeralsWithCurrency(value, props.defaultCurrency);
    }

    return formatNumberToLocalizedNumerals(value, 2);
}

function exportData(): { headers: string[], data: string[][] } {
    const headers: string[] = [];
    const data: string[][] = [];

    headers.push(props.categoryTypeName);

    for (const categoryName of props.allCategoryNames) {
        headers.push(categoryName);
    }

    for (const item of treeMapData.value) {
        const row: string[] = [];
        row.push(item.name);

        for (const child of item.children ?? []) {
            row.push(formatAmountToWesternArabicNumeralsWithoutDigitGrouping(child.value));
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
.treemap-chart-container {
    width: 100%;
    height: 560px;
    margin-top: 10px;
}

@media (min-width: 600px) {
    .treemap-chart-container {
        height: 630px;
    }
}
</style>
