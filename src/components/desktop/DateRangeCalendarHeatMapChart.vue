<template>
    <v-chart autoresize :class="finalClass"
             :option="chartOptions" :update-options="{ notMerge: true }"
             @click="clickItem" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import type { ECElementEvent } from 'echarts/core';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { type WeekDayValue, KnownDateTimeFormat } from '@/core/datetime.ts';
import { ThemeType } from '@/core/theme.ts';
import { ChartValueType, type CalendarChartSourceDataItem } from '@/core/chart.ts';

import {
    BIG_DECIMAL_ZERO,
    BIG_DECIMAL_POSITIVE_INFINITY,
    isBigDecimal
} from '@/lib/numeral.ts';

import {
    parseDateTimeFromUnixTime,
    parseDateTimeFromKnownDateTimeFormat
} from '@/lib/datetime.ts';

interface HeatMapData {
    allOriginalDataMap: Record<string, BigDecimal>;
    data: [string, number][]; // second value only used for echarts rendering, the actual value is in allOriginalDataMap
    minValue: BigDecimal;
    maxValue: BigDecimal;
}

const props = defineProps<{
    class?: string;
    skeleton?: boolean;
    showValue?: boolean;
    enableClickItem?: boolean;
    startTime: number;
    endTime: number;
    items: CalendarChartSourceDataItem[];
    valueType: ChartValueType;
    valueTypeName: string;
    translateName?: boolean;
    defaultCurrency?: string;
}>();

const emit = defineEmits<{
    (e: 'click', date: string, displayDate: string, value: BigDecimal): void;
}>();

const theme = useTheme();

const {
    tt,
    getAllMinWeekdayNames,
    formatDateTimeToLongDate,
    formatChartValueToLocalizedNumerals
} = useI18n();

const userStore = useUserStore();

const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
const dayNames = computed<string[]>(() => getAllMinWeekdayNames());

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const finalClass = computed<string>(() => {
    let finalClass = '';

    if (props.skeleton) {
        finalClass += 'transition-in';
    }

    if (props.class) {
        finalClass += ` ${props.class}`;
    } else {
        finalClass += ' date-range-calendar-heatmap-chart-container';
    }

    return finalClass;
});

const startDate = computed<string>(() => parseDateTimeFromUnixTime(props.startTime).getGregorianCalendarYearDashMonthDashDay());
const endDate = computed<string>(() => parseDateTimeFromUnixTime(props.endTime).getGregorianCalendarYearDashMonthDashDay());

const heatMapData = computed<HeatMapData>(() => {
    const allOriginalDataMap: Record<string, BigDecimal> = {};
    const data: [string, number][] = [];
    let minValue: BigDecimal = BIG_DECIMAL_POSITIVE_INFINITY;
    let maxValue: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of props.items) {
        const id = getItemName(item.id);
        const dateTime = parseDateTimeFromKnownDateTimeFormat(id, KnownDateTimeFormat.DefaultDate);

        if (!dateTime || dateTime.getUnixTime() < props.startTime || dateTime.getUnixTime() > props.endTime || !isBigDecimal(item.value) || item.hidden) {
            continue;
        }

        if (item.value.greaterThan(maxValue)) {
            maxValue = item.value;
        }

        if (item.value.lessThan(minValue)) {
            minValue = item.value;
        }

        const date = dateTime.getGregorianCalendarYearDashMonthDashDay();
        allOriginalDataMap[date] = item.value;
        data.push([date, item.value.toDoubleNumber()]);
    }

    return {
        allOriginalDataMap: allOriginalDataMap,
        data: data,
        minValue: minValue.isPositiveInfinity() ? BIG_DECIMAL_ZERO : minValue,
        maxValue: maxValue
    };
});

const chartOptions = computed<object>(() => {
    return {
        tooltip: {
            backgroundColor: isDarkMode.value ? '#333' : '#fff',
            borderColor: isDarkMode.value ? '#333' : '#fff',
            textStyle: {
                color: isDarkMode.value ? '#eee' : '#333'
            },
            formatter: (params: CallbackDataParams) => {
                if (!props.showValue) {
                    return '';
                }

                const dataItem = params.data as [string, number];
                const dateTime = dataItem && dataItem[0] ? parseDateTimeFromKnownDateTimeFormat(dataItem[0], KnownDateTimeFormat.DefaultDate) : '';
                const value: BigDecimal | undefined = dataItem && dataItem[0] ? heatMapData.value.allOriginalDataMap[dataItem[0]] : undefined;
                const displayValue: string = value ? formatChartValueToLocalizedNumerals(value, props.valueType, props.defaultCurrency) : '';

                return (dateTime ? `<div class="d-inline-flex">${formatDateTimeToLongDate(dateTime)}</div><br/>` : '')
                    + `<div><span class="chart-pointer" style="background-color: ${params.color}"></span>`
                    + `<span>${props.valueTypeName}</span>`
                    + `<span class="ms-5">${displayValue}</span>`
                    + '</div>';
            }
        },
        visualMap: {
            show: false,
            min: heatMapData.value.minValue.toDoubleNumber(),
            max: heatMapData.value.maxValue.toDoubleNumber(),
            inRange: {
                color: isDarkMode.value ? [ '#1a1a1a', '#c67e48' ] : [ '#faf8f4', '#c67e48' ]
            }
        },
        calendar: {
            range: [ startDate.value, endDate.value ],
            orient: 'horizontal',
            left: 40,
            top: 0,
            right: 20,
            bottom: 10,
            cellSize: ['auto', 20],
            itemStyle: {
                color: isDarkMode.value ? '#060504' : '#ffffff',
                borderColor: isDarkMode.value ? '#4f4f4f' : '#e1e6f2'
            },
            splitLine: {
                show: false
            },
            dayLabel: {
                firstDay: firstDayOfWeek.value,
                nameMap: dayNames.value,
                color: isDarkMode.value ? '#888' : '#666'
            },
            monthLabel: {
                show: false
            },
            yearLabel: {
                show: false
            }
        },
        series: {
            type: 'heatmap',
            animation: !props.skeleton,
            coordinateSystem: 'calendar',
            data: heatMapData.value.data,
            label: {
                show: false
            },
            emphasis: {
                itemStyle: {
                    shadowBlur: 6,
                    shadowColor: isDarkMode.value ? 'rgba(255, 255, 255, 0.5)' : 'rgba(0, 0, 0, 0.5)'
                }
            }
        }
    };
});

function getItemName(name: string): string {
    return props.translateName ? tt(name) : name;
}

function clickItem(e: ECElementEvent): void {
    if (!props.enableClickItem || e.componentType !== 'series') {
        return;
    }

    const dataItem = e.data as [string, number];

    if (!dataItem || !dataItem[0]) {
        return;
    }

    const date = dataItem[0];
    const dateTime = parseDateTimeFromKnownDateTimeFormat(date, KnownDateTimeFormat.DefaultDate);
    const displayDate = dateTime ? formatDateTimeToLongDate(dateTime) : '';
    const value: BigDecimal = heatMapData.value.allOriginalDataMap[date] ?? BIG_DECIMAL_ZERO;
    emit('click', date, displayDate, value);
}
</script>

<style scoped>
.date-range-calendar-heatmap-chart-container {
    width: 100%;
    height: 100%;
}
</style>
