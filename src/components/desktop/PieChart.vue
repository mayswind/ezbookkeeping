<template>
    <v-chart autoresize class="pie-chart-container" :class="{ 'transition-in': skeleton }" :option="chartOptions"
             @click="clickItem" @legendselectchanged="onLegendSelectChanged" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useTheme } from 'vuetify';

import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonPieChartDataItem, type CommonPieChartProps, usePieChartBase } from '@/components/base/PieChartBase.ts'

import type { ColorValue } from '@/core/color.ts';
import { ThemeType } from '@/core/theme.ts';
import { DEFAULT_ICON_COLOR } from '@/consts/color.ts';

interface DesktopPieChartDataItem extends CommonPieChartDataItem {
    itemStyle: {
        color: ColorValue;
    };
    selected: boolean;
}

const props = defineProps<CommonPieChartProps>();

const emit = defineEmits<{
    (e: 'click', value: Record<string, unknown>): void;
}>();

const theme = useTheme();

const { formatAmountToLocalizedNumeralsWithCurrency } = useI18n();
const { selectedIndex, validItems } = usePieChartBase(props);

const selectedLegends = ref<Record<string, boolean> | null>(null);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const itemsMap = computed<Record<string, Record<string, unknown>>>(() => {
    const map: Record<string, Record<string, unknown>> = {};

    for (let i = 0; i < props.items.length; i++) {
        const item = props.items[i];
        let id = '';

        if (props.idField && item[props.idField]) {
            id = item[props.idField] as string;
        } else {
            id = item[props.nameField] as string;;
        }

        map[id] = item;
    }

    return map;
});

const seriesData = computed<DesktopPieChartDataItem[]>(() => {
    const ret: DesktopPieChartDataItem[] = [];

    for (let i = 0; i < validItems.value.length; i++) {
        const item = validItems.value[i];
        ret.push({
            ...item,
            itemStyle: {
                color: getColor(item.color),
            },
            selected: true
        });
    }

    return ret;
});

const hasUnselectedItem = computed<boolean>(() => {
    for (let i = 0; i < validItems.value.length; i++) {
        const item = validItems.value[i];

        if (selectedLegends.value && !selectedLegends.value[item.id]) {
            return true;
        }
    }

    return false;
});

const firstItemAndHalfCurrentItemTotalPercent = computed<number>(() => {
    let totalValue = 0;
    let firstValue = null;
    let firstToCurrentTotalValue = 0;

    for (let i = 0; i < validItems.value.length; i++) {
        const item = validItems.value[i];

        if (selectedLegends.value && !selectedLegends.value[item.id]) {
            continue;
        }

        if (firstValue === null) {
            firstValue = item.value;
        }

        if (firstValue !== null) {
            if (i < selectedIndex.value) {
                firstToCurrentTotalValue += item.value;
            } else if (i === selectedIndex.value) {
                firstToCurrentTotalValue += item.value / 2;
            }
        }

        totalValue += item.value;
    }

    if (firstToCurrentTotalValue && totalValue > 0) {
        return firstToCurrentTotalValue / totalValue;
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
                const name = dataItem ? dataItem.displayName : '';
                const value = dataItem ? dataItem.displayValue : formatAmountToLocalizedNumeralsWithCurrency(params.value as number);
                let percent = dataItem ? dataItem.displayPercent : (params.percent + '%');

                if (hasUnselectedItem.value) {
                    percent = params.percent + '%';
                }

                let tooltip = `<div><span class="chart-pointer" style="background-color: ${params.color}"></span>`;

                if (name) {
                    tooltip += `<div class="d-inline-flex">${name}</div><br/>`;
                }

                if (props.showValue) {
                    tooltip += `<div class="d-inline-flex"><span>${value}</span><span class="ms-1">(${percent})</span></div>`;
                } else {
                    tooltip += `<div class="d-inline-flex">${percent}</div>`;
                }

                tooltip += '</div>';

                return tooltip;
            }
        },
        legend: {
            orient: 'horizontal',
            data: validItems.value.map(item => item.name),
            selected: selectedLegends.value,
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (id: string) => {
                const item = itemsMap.value[id];
                return item && props.nameField && item[props.nameField] ? item[props.nameField] as string : id;
            }
        },
        series: [
            {
                type: 'pie',
                data: seriesData.value,
                top: 50,
                startAngle: -90 + firstItemAndHalfCurrentItemTotalPercent.value * 360,
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

function getColor(color: string): ColorValue {
    if (color && color !== DEFAULT_ICON_COLOR) {
        color = '#' + color;
    }

    return color;
}

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

        for (let i = 0; i < validItems.value.length; i++) {
            const item = validItems.value[i];

            if (selectedLegends.value[item.id]) {
                newSelectedIndex = i;
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
    height: 400px;
}

@media (min-width: 600px) {
    .pie-chart-container {
        height: 500px;
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
