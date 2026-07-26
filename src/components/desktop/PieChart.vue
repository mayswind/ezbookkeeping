<template>
    <v-chart autoresize class="pie-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"
             @click="clickItem" @legendselectchanged="onLegendSelectChanged" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useTheme } from 'vuetify';

import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { type CommonPieChartDataItem, type CommonPieChartProps, usePieChartBase } from '@/components/base/PieChartBase.ts'

import { itemAndIndex } from '@/core/base.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import type { ColorStyleValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';

import { getObjectOwnFieldCount } from '@/lib/common.ts';
import { BIG_DECIMAL_ZERO, parseBigDecimal } from '@/lib/numeral.ts';

interface DesktopPieChartDataItem extends CommonPieChartDataItem {
    itemStyle: {
        color: ColorStyleValue;
    };
    selected: boolean;
}

const props = defineProps<CommonPieChartProps>();

const emit = defineEmits<{
    (e: 'click', value: Record<string, unknown>): void;
}>();

const theme = useTheme();

const { selectedIndex, validItems, allItemsMap } = usePieChartBase(props);

const selectedLegends = ref<Record<string, boolean>>({});

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const seriesData = computed<DesktopPieChartDataItem[]>(() => {
    const ret: DesktopPieChartDataItem[] = [];

    for (const item of validItems.value) {
        ret.push({
            ...item,
            itemStyle: {
                color: item.color,
            },
            selected: true
        });
    }

    return ret;
});

const hasUnselectedItem = computed<boolean>(() => {
    for (const item of validItems.value) {
        if (getObjectOwnFieldCount(selectedLegends.value) && !selectedLegends.value[item.id]) {
            return true;
        }
    }

    return false;
});

const firstItemAndHalfCurrentItemTotalPercent = computed<number>(() => {
    let totalValue: BigDecimal = BIG_DECIMAL_ZERO;
    let firstValue: string | null = null;
    let firstToCurrentTotalValue: BigDecimal = BIG_DECIMAL_ZERO;

    for (const [item, index] of itemAndIndex(validItems.value)) {
        if (getObjectOwnFieldCount(selectedLegends.value) && !selectedLegends.value[item.id]) {
            continue;
        }

        if (firstValue === null) {
            firstValue = item.value;
        }

        if (firstValue !== null) {
            if (index < selectedIndex.value) {
                firstToCurrentTotalValue = firstToCurrentTotalValue.add(parseBigDecimal(item.value));
            } else if (index === selectedIndex.value) {
                firstToCurrentTotalValue = firstToCurrentTotalValue.add(parseBigDecimal(item.value).divide(2));
            }
        }

        totalValue = totalValue.add(parseBigDecimal(item.value));
    }

    if (firstToCurrentTotalValue && totalValue.isPositive()) {
        return firstToCurrentTotalValue.divide(totalValue).toDoubleNumber();
    } else {
        return 0;
    }
});

const chartOptions = computed<object>(() => {
    return {
        tooltip: {
            trigger: 'item',
            backgroundColor: isDarkMode.value ? '#333' : '#fff',
            borderColor: isDarkMode.value ? '#333' : '#fff',
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams) => {
                const dataItem = params.data as DesktopPieChartDataItem;
                let percent = dataItem.displayPercent;

                if (hasUnselectedItem.value) {
                    percent = params.percent + '%';
                }

                let tooltip = `<div><span class="chart-pointer" style="background-color: ${params.color}"></span>`;

                if (dataItem.displayName) {
                    tooltip += `<div class="d-inline-flex">${dataItem.displayName}</div><br/>`;
                }

                const showValue = props.showValue;
                const showPercent = props.showPercent && dataItem.originalValuePositive;

                if (showValue && showPercent) {
                    tooltip += `<div class="d-inline-flex"><span>${dataItem.displayValue}</span><span class="ms-1">(${percent})</span></div>`;
                } else if (showValue && !showPercent) {
                    tooltip += `<div class="d-inline-flex">${dataItem.displayValue}</div>`;
                } else if (!showValue && showPercent) {
                    tooltip += `<div class="d-inline-flex">${percent}</div>`;
                }

                tooltip += '</div>';

                return tooltip;
            }
        },
        legend: {
            orient: 'horizontal',
            type: 'scroll',
            top: 0,
            data: validItems.value.map(item => item.name),
            selected: selectedLegends.value,
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (id: string) => allItemsMap.value[id]?.name ?? id
        },
        series: [
            {
                type: 'pie',
                data: seriesData.value,
                top: 50,
                startAngle: -90 + firstItemAndHalfCurrentItemTotalPercent.value * 360,
                radius: [0, '75%'],
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)',
                    }
                },
                label: {
                    color: isDarkMode.value ? '#eee' : '#333',
                    formatter: (params: CallbackDataParams) => {
                        const dataItem = params.data as DesktopPieChartDataItem;
                        return dataItem ? dataItem.displayName : '';
                    }
                },
                animation: !props.skeleton
            }
        ],
        media: [
            {
                query: {
                    minWidth: 600,
                },
                option: {
                    legend: {
                        orient: 'vertical',
                        left: 'left'
                    },
                    series: [
                        {
                            type: 'pie',
                            top: 0
                        }
                    ]
                },
            }
        ]
    };
});

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series' || e.seriesType !=='pie') {
        return;
    }

    if (e.event && e.event.target && e.event.target.currentStates && e.event.target.currentStates[0] && e.event.target.currentStates[0] === 'emphasis') {
        selectedIndex.value = e.dataIndex;
        return;
    }

    if (!e.data) {
        return;
    }

    const data = e.data as object;

    if ('sourceItem' in data) {
        emit('click', data.sourceItem as Record<string, unknown>);
    }
}

function onLegendSelectChanged(e: { selected: Record<string, boolean> }): void {
    selectedLegends.value = e.selected;
    const selectedItem = validItems.value[selectedIndex.value];

    if (!selectedItem || !selectedLegends.value[selectedItem.id]) {
        let newSelectedIndex = 0;

        for (const [item, index] of itemAndIndex(validItems.value)) {
            if (selectedLegends.value[item.id]) {
                newSelectedIndex = index;
                break;
            }
        }

        selectedIndex.value = newSelectedIndex;
    }
}
</script>

<style scoped>
.pie-chart-container {
    width: 100%;
    height: 460px;
}

@media (min-width: 600px) {
    .pie-chart-container {
        height: 650px;
    }
}

.pie-chart-container.transition-in {
    animation: pie-chart-skeleton-fade-in 2s 1;
}

@keyframes pie-chart-skeleton-fade-in {
    0% {
        opacity: 0;
    }
    20% {
        opacity: 0;
    }
    100% {
        opacity: 1;
    }
}
</style>
