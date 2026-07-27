<template>
    <v-chart autoresize :class="finalClass" :style="finalStyle"
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
    getObjectOwnFieldCount,
    mapObjectToArray
} from '@/lib/common.ts';
import {
    BIG_DECIMAL_ZERO,
    BIG_DECIMAL_POSITIVE_INFINITY,
    parseBigDecimal,
    isBigDecimal
} from '@/lib/numeral.ts';
import { parseDateTimeFromKnownDateTimeFormat } from '@/lib/datetime.ts';

interface HeatMapData {
    allOriginalDataMap: Record<string, BigDecimal>;
    data: Record<number, YearlyHeatmapData>;
    minValue: BigDecimal;
    maxValue: BigDecimal;
}

interface YearlyHeatmapData {
    gregorianYear: number;
    displayYear: string;
    data: [string, number][]; // second value only used for echarts rendering, the actual value is in allOriginalDataMap of HeatMapData
}

const props = defineProps<{
    class?: string;
    skeleton?: boolean;
    showValue?: boolean;
    enableClickItem?: boolean;
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
    getAllShortMonthNames,
    getAllMinWeekdayNames,
    formatDateTimeToLongDate,
    getCalendarDisplayLongYearFromDateTime,
    formatChartValueToLocalizedNumerals
} = useI18n();

const userStore = useUserStore();

const visualMapHeight: number = 100;
const calendarHeight: number = 180;
const calendarBottomMargin: number = 10;

const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
const dayNames = computed<string[]>(() => getAllMinWeekdayNames());
const monthNames = computed<string[]>(() => getAllShortMonthNames());

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const finalClass = computed<string>(() => {
    let finalClass = '';

    if (props.skeleton) {
        finalClass += 'transition-in';
    }

    if (props.class) {
        finalClass += ` ${props.class}`;
    } else {
        finalClass += ' calendar-heatmap-chart-container';
    }

    return finalClass;
});
const finalStyle = computed<Record<string, string>>(() => {
    const style: Record<string, string> = {};

    if (heatMapData.value.data) {
        const calendarCount = getObjectOwnFieldCount(heatMapData.value.data);
        style['height'] = `${visualMapHeight + calendarCount * calendarHeight + (calendarCount - 1) * calendarBottomMargin}px`;
    }

    return style;
});

const heatMapData = computed<HeatMapData>(() => {
    const allOriginalDataMap: Record<string, BigDecimal> = {};
    const allData: Record<number, YearlyHeatmapData> = {};
    let minValue: BigDecimal = BIG_DECIMAL_POSITIVE_INFINITY;
    let maxValue: BigDecimal = BIG_DECIMAL_ZERO;

    for (const item of props.items) {
        const id = getItemName(item.id);
        const dateTime = parseDateTimeFromKnownDateTimeFormat(id, KnownDateTimeFormat.DefaultDate);

        if (dateTime && isBigDecimal(item.value) && !item.hidden) {
            if (item.value.greaterThan(maxValue)) {
                maxValue = item.value;
            }

            if (item.value.lessThan(minValue)) {
                minValue = item.value;
            }

            const year: number = dateTime.getGregorianCalendarYear();
            let data: YearlyHeatmapData | undefined = allData[year];

            if (!data) {
                data = {
                    gregorianYear: year,
                    displayYear: getCalendarDisplayLongYearFromDateTime(dateTime),
                    data: []
                };
                allData[year] = data;
            }

            allOriginalDataMap[dateTime.getGregorianCalendarYearDashMonthDashDay()] = item.value;
            data.data.push([dateTime.getGregorianCalendarYearDashMonthDashDay(), item.value.toDoubleNumber()]);
        }
    }

    const ret: HeatMapData = {
        allOriginalDataMap: allOriginalDataMap,
        data: allData,
        minValue: minValue.isPositiveInfinity() ? BIG_DECIMAL_ZERO : minValue,
        maxValue: maxValue
    };

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
            formatter: (params: CallbackDataParams) => {
                if (!props.showValue) {
                    return '';
                }

                const dataItem = params.data as [string, number];
                const dateTime = dataItem && dataItem[0] ? parseDateTimeFromKnownDateTimeFormat(dataItem[0], KnownDateTimeFormat.DefaultDate) : '';
                const name = props.valueTypeName;
                const value: BigDecimal | undefined = dataItem && dataItem[0] ? heatMapData.value.allOriginalDataMap[dataItem[0]] : undefined;
                const displayValue: string = value ? formatChartValueToLocalizedNumerals(value, props.valueType, props.defaultCurrency) : '';

                return (dateTime ? `<div class="d-inline-flex">${formatDateTimeToLongDate(dateTime)}</div><br/>` : '')
                    + `<div><span class="chart-pointer" style="background-color: ${params.color}"></span>`
                    + `<span>${name}</span>`
                    + `<span class="ms-5">${displayValue}</span>`
                    + '</div>';
            }
        },
        visualMap: [
            {
                type: 'continuous',
                orient: 'horizontal',
                top: 0,
                left: 'center',
                itemHeight: 320,
                min: heatMapData.value.minValue.toDoubleNumber(),
                max: heatMapData.value.maxValue.toDoubleNumber(),
                calculable: true,
                inRange: {
                    color: isDarkMode.value ? [ '#1a1a1a', '#c67e48' ] : [ '#faf8f4', '#c67e48' ]
                },
                textStyle: {
                    color: isDarkMode.value ? '#888' : '#666'
                },
                formatter: (value: number) => {
                    if (!props.showValue) {
                        return '';
                    }

                    let actualValue: BigDecimal = parseBigDecimal(value);

                    if (value === heatMapData.value.minValue.toDoubleNumber()) {
                        actualValue = heatMapData.value.minValue;
                    } else if (value === heatMapData.value.maxValue.toDoubleNumber()) {
                        actualValue = heatMapData.value.maxValue;
                    }

                    return formatChartValueToLocalizedNumerals(actualValue, props.valueType, props.defaultCurrency);
                }
            }
        ],
        calendar: mapObjectToArray(heatMapData.value.data, (item, _, index) => {
            return {
                range: item.gregorianYear,
                orient: 'horizontal',
                left: 70,
                top: visualMapHeight + index * (calendarHeight + calendarBottomMargin),
                right: 20,
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
                    nameMap: monthNames.value,
                    color: isDarkMode.value ? '#888' : '#666'
                },
                yearLabel: {
                    formatter: item.displayYear,
                    color: isDarkMode.value ? '#888' : '#666'
                }
            };
        }),
        series: mapObjectToArray(heatMapData.value.data, (item, _, index) => {
            return {
                type: 'heatmap',
                animation: !props.skeleton,
                coordinateSystem: 'calendar',
                calendarIndex: index,
                data: item.data,
                label: {
                    show: false
                },
                emphasis: {
                    itemStyle: {
                        shadowBlur: 6,
                        shadowColor: isDarkMode.value ? 'rgba(255, 255, 255, 0.5)' : 'rgba(0, 0, 0, 0.5)'
                    }
                }
            };
        })
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
    const value: BigDecimal = dataItem && dataItem[0] ? (heatMapData.value.allOriginalDataMap[dataItem[0]] ?? BIG_DECIMAL_ZERO) : BIG_DECIMAL_ZERO;
    emit('click', date, displayDate, value);
}
</script>

<style scoped>
.calendar-heatmap-chart-container {
    width: 100%;
    margin-top: 10px;
}
</style>
