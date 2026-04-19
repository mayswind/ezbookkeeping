<template>
    <v-chart autoresize :class="finalClass" :style="finalStyle" :option="chartOptions" />
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import type { CallbackDataParams } from 'echarts/types/dist/shared';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { type WeekDayValue, KnownDateTimeFormat } from '@/core/datetime.ts';
import { ThemeType } from '@/core/theme.ts';

import {
    isNumber,
    getObjectOwnFieldCount,
    mapObjectToArray
} from '@/lib/common.ts';
import { parseDateTimeFromKnownDateTimeFormat } from '@/lib/datetime.ts';

interface HeatMapData {
    data: Record<number, YearlyHeatmapData>;
    minValue: number;
    maxValue: number;
}

interface YearlyHeatmapData {
    gregorianYear: number;
    displayYear: string;
    data: [string, number][];
}

const props = defineProps<{
    class?: string;
    skeleton?: boolean;
    showValue?: boolean;
    items: Record<string, unknown>[];
    idField: string;
    valueField: string;
    hiddenField?: string;
    translateName?: boolean;
    valueTypeName: string;
    amountValue?: boolean;
    percentValue?: boolean;
    defaultCurrency?: string;
}>();

const theme = useTheme();

const {
    tt,
    getAllShortMonthNames,
    getAllMinWeekdayNames,
    formatDateTimeToLongDate,
    getCalendarDisplayLongYearFromDateTime,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatNumberToLocalizedNumerals,
    formatPercentToLocalizedNumerals
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
    const allData: Record<number, YearlyHeatmapData> = {};
    let minValue: number = Number.POSITIVE_INFINITY;
    let maxValue: number = 0;

    for (const item of props.items) {
        const id = getItemName(item[props.idField] as string);
        const dateTime = parseDateTimeFromKnownDateTimeFormat(id, KnownDateTimeFormat.DefaultDate);
        const value = item[props.valueField];

        if (dateTime && isNumber(value) && (!props.hiddenField || !item[props.hiddenField])) {
            if (value > maxValue) {
                maxValue = value;
            }

            if (value < minValue) {
                minValue = value;
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

            data.data.push([dateTime.getGregorianCalendarYearDashMonthDashDay(), value]);
        }
    }

    const ret: HeatMapData = {
        data: allData,
        minValue: minValue === Number.POSITIVE_INFINITY ? 0 : minValue,
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
                const value = dataItem && isNumber(dataItem[1]) ? getDisplayValue(dataItem[1]) : '';

                return (dateTime ? `<div class="d-inline-flex">${formatDateTimeToLongDate(dateTime)}</div><br/>` : '')
                    + `<div><span class="chart-pointer" style="background-color: ${params.color}"></span>`
                    + `<span>${name}</span>`
                    + `<span class="ms-5">${value}</span>`
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
                min: heatMapData.value.minValue,
                max: heatMapData.value.maxValue,
                calculable: true,
                inRange: {
                    color: isDarkMode.value ? [ '#1a1a1a', '#c67e48' ] : [ '#faf8f4', '#c67e48' ]
                },
                textStyle: {
                    color: isDarkMode.value ? '#888' : '#666'
                },
                formatter: (value: string) => {
                    if (!props.showValue) {
                        return '';
                    }

                    return getDisplayValue(parseInt(value));
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

function getDisplayValue(value: number): string {
    if (props.percentValue) {
        return formatPercentToLocalizedNumerals(value, 2, '<0.01');
    }

    if (props.amountValue) {
        return formatAmountToLocalizedNumeralsWithCurrency(value, props.defaultCurrency);
    }

    return formatNumberToLocalizedNumerals(value, 2);
}
</script>

<style scoped>
.calendar-heatmap-chart-container {
    width: 100%;
    margin-top: 10px;
}
</style>
