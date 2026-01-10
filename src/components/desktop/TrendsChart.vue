<template>
    <axis-chart class="trends-chart-container" ref="axisChart" values-field="values"
                :skeleton="skeleton" :type="chartDisplayType" :stacked="stacked" :sorting-type="sortingType"
                :show-value="showValue"
                :show-total-amount-in-tooltip="showTotalAmountInTooltip" :total-name-in-tooltip="tt('Total Amount')"
                :category-type-name="tt('Date')" :all-category-names="allDisplayDateRanges" :items="allSeriesData"
                :id-field="idField" :name-field="nameField" :color-field="colorField" :hidden-field="hiddenField"
                :display-orders-field="displayOrdersField"
                :translate-name="translateName"
                :amount-value="true" :default-currency="defaultCurrency"
                :enable-click-item="enableClickItem"
                @click="clickItem"
                v-if="chartDisplayType"
    />
</template>

<script setup lang="ts">
import AxisChart, { type AxisChartDisplayType } from './AxisChart.vue';

import { computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type TrendsChartDateType,
    type CommonTrendsChartProps,
    type TrendsBarChartClickEvent,
    useTrendsChartBase
} from '@/components/base/TrendsChartBase.ts'

import { useUserStore } from '@/stores/user.ts';

import {
    type Year1BasedMonth,
    type YearMonthDay,
    DateRangeScene
} from '@/core/datetime.ts';
import {
    ChartDataAggregationType,
    TrendChartType,
    ChartDateAggregationType
} from '@/core/statistics.ts';

import { isArray, isNumber } from '@/lib/common.ts';
import {
    parseDateTimeFromUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getDateTypeByDateRange,
    getFiscalYearFromUnixTime
} from '@/lib/datetime.ts';

type AxisChartType = InstanceType<typeof AxisChart>;

interface DesktopTrendsChartProps<T extends TrendsChartDateType> extends CommonTrendsChartProps<T> {
    skeleton?: boolean;
    type?: number;
    showValue?: boolean;
    showTotalAmountInTooltip?: boolean;
}

const props = defineProps<DesktopTrendsChartProps<TrendsChartDateType>>();

const emit = defineEmits<{
    (e: 'click', value: TrendsBarChartClickEvent): void;
}>();

const {
    tt,
    formatDateTimeToShortDate,
    formatDateTimeToGregorianLikeShortYear,
    formatDateTimeToGregorianLikeShortYearMonth,
    formatYearQuarterToGregorianLikeYearQuarter,
    formatDateTimeToGregorianLikeFiscalYear
} = useI18n();

const { allDateRanges } = useTrendsChartBase(props);

const userStore = useUserStore();

const axisChart = useTemplateRef<AxisChartType>('axisChart');

const chartDisplayType = computed<AxisChartDisplayType | undefined>(() => {
    if (props.type === TrendChartType.Area.type) {
        return 'area';
    } else if (props.type === TrendChartType.Column.type) {
        return 'column';
    } else if (props.type === TrendChartType.Bubble.type) {
        return 'bubble';
    } else {
        return undefined;
    }
});

const allDisplayDateRanges = computed<string[]>(() => {
    const allDisplayDateRanges: string[] = [];

    for (const dateRange of allDateRanges.value) {
        const minDateTime = parseDateTimeFromUnixTime(dateRange.minUnixTime);

        if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
            allDisplayDateRanges.push(formatDateTimeToGregorianLikeShortYear(minDateTime));
        } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type && 'year' in dateRange) {
            allDisplayDateRanges.push(formatDateTimeToGregorianLikeFiscalYear(minDateTime));
        } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
            allDisplayDateRanges.push(formatYearQuarterToGregorianLikeYearQuarter(dateRange.year, dateRange.quarter));
        } else if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
            allDisplayDateRanges.push(formatDateTimeToGregorianLikeShortYearMonth(minDateTime));
        } else if (props.dateAggregationType === ChartDateAggregationType.Day.type && props.chartMode === 'daily') {
            allDisplayDateRanges.push(formatDateTimeToShortDate(minDateTime));
        }
    }

    return allDisplayDateRanges;
});

const allSeriesData = computed<Record<string, unknown>[]>(() => {
    const result: Record<string, unknown>[] = [];

    for (const item of props.items) {
        if (props.hiddenField && item[props.hiddenField]) {
            continue;
        }

        const finalItem: Record<string, unknown> = {};

        if (props.idField) {
            finalItem[props.idField] = item[props.idField];
        }

        if (props.nameField) {
            finalItem[props.nameField] = item[props.nameField];
        }

        if (props.colorField) {
            finalItem[props.colorField] = item[props.colorField];
        }

        if (props.hiddenField) {
            finalItem[props.hiddenField] = item[props.hiddenField];
        }

        if (props.displayOrdersField) {
            finalItem[props.displayOrdersField] = item[props.displayOrdersField];
        }

        const allAmounts: number[] = [];
        const dateRangeAmountMap: Record<string, (Year1BasedMonth | YearMonthDay)[]> = {};

        for (const dataItem of item.items) {
            let dateRangeKey = '';

            if (props.chartMode === 'daily' && 'month' in dataItem) {
                if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                    dateRangeKey = dataItem.year.toString();
                } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
                    const fiscalYear = getFiscalYearFromUnixTime(
                        getYearMonthFirstUnixTime({ year: dataItem.year, month1base: dataItem.month }),
                        props.fiscalYearStart
                    );
                    dateRangeKey = fiscalYear.toString();
                } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                    dateRangeKey = `${dataItem.year}-${Math.floor((dataItem.month - 1) / 3) + 1}`;
                } else if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                    dateRangeKey = `${dataItem.year}-${dataItem.month}`;
                } else { // if (props.dateAggregationType === ChartDateAggregationType.Day.type) {
                    dateRangeKey = `${dataItem.year}-${dataItem.month}-${dataItem.day}`;
                }
            } else if (props.chartMode === 'monthly' && 'month1base' in dataItem) {
                if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                    dateRangeKey = dataItem.year.toString();
                } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type) {
                    const fiscalYear = getFiscalYearFromUnixTime(
                        getYearMonthFirstUnixTime({ year: dataItem.year, month1base: dataItem.month1base }),
                        props.fiscalYearStart
                    );
                    dateRangeKey = fiscalYear.toString();
                } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type) {
                    dateRangeKey = `${dataItem.year}-${Math.floor((dataItem.month1base - 1) / 3) + 1}`;
                } else { // if (props.dateAggregationType === ChartDateAggregationType.Month.type) {
                    dateRangeKey = `${dataItem.year}-${dataItem.month1base}`;
                }
            }

            const dataItems = dateRangeAmountMap[dateRangeKey] || [];
            dataItems.push(dataItem);

            dateRangeAmountMap[dateRangeKey] = dataItems;
        }

        for (const dateRange of allDateRanges.value) {
            let dateRangeKey = '';

            if (props.dateAggregationType === ChartDateAggregationType.Year.type) {
                dateRangeKey = dateRange.year.toString();
            } else if (props.dateAggregationType === ChartDateAggregationType.FiscalYear.type && 'year' in dateRange) {
                dateRangeKey = dateRange.year.toString();
            } else if (props.dateAggregationType === ChartDateAggregationType.Quarter.type && 'quarter' in dateRange) {
                dateRangeKey = `${dateRange.year}-${dateRange.quarter}`;
            } else if (props.dateAggregationType === ChartDateAggregationType.Month.type && 'month0base' in dateRange) {
                dateRangeKey = `${dateRange.year}-${dateRange.month0base + 1}`;
            } else if (props.dateAggregationType === ChartDateAggregationType.Day.type && 'day' in dateRange && props.chartMode === 'daily') {
                dateRangeKey = `${dateRange.year}-${dateRange.month}-${dateRange.day}`;
            }

            let amount = 0;
            const dataItems = dateRangeAmountMap[dateRangeKey];

            if (isArray(dataItems)) {
                for (const dataItem of dataItems) {
                    const value = (dataItem as unknown as Record<string, unknown>)[props.valueField];

                    if (isNumber(value)) {
                        if (props.dataAggregationType === ChartDataAggregationType.Sum) {
                            amount += value;
                        } else if (props.dataAggregationType === ChartDataAggregationType.Last) {
                            amount = value;
                        }
                    }
                }
            }

            allAmounts.push(amount);
        }

        finalItem['values'] = allAmounts;
        result.push(finalItem);
    }

    return result;
});

function clickItem(itemId: string, categoryIndex: number): void {
    const dateRange = allDateRanges.value[categoryIndex];

    if (!dateRange) {
        return;
    }

    let minUnixTime = dateRange.minUnixTime;
    let maxUnixTime = dateRange.maxUnixTime;

    if (props.chartMode === 'daily') {
        if (props.startTime) {
            if (props.startTime > minUnixTime) {
                minUnixTime = props.startTime;
            }
        }

        if (props.endTime) {
            if (props.endTime < maxUnixTime) {
                maxUnixTime = props.endTime;
            }
        }
    } else if (props.chartMode === 'monthly') {
        if (props.startYearMonth) {
            const startMinUnixTime = getYearMonthFirstUnixTime(props.startYearMonth);

            if (startMinUnixTime > minUnixTime) {
                minUnixTime = startMinUnixTime;
            }
        }

        if (props.endYearMonth) {
            const endMaxUnixTime = getYearMonthLastUnixTime(props.endYearMonth);

            if (endMaxUnixTime < maxUnixTime) {
                maxUnixTime = endMaxUnixTime;
            }
        }
    }

    const dateRangeType = getDateTypeByDateRange(minUnixTime, maxUnixTime, userStore.currentUserFirstDayOfWeek, userStore.currentUserFiscalYearStart, DateRangeScene.Normal);

    emit('click', {
        itemId: itemId,
        dateRange: {
            minTime: minUnixTime,
            maxTime: maxUnixTime,
            dateType: dateRangeType
        }
    });
}

function exportData(): { headers: string[], data: string[][] } {
    return axisChart.value?.exportData() ?? { headers: [], data: [] };
}

defineExpose({
    exportData
})
</script>

<style scoped>
.trends-chart-container {
    width: 100%;
    height: 720px;
    margin-top: 10px;
}

@media (min-width: 600px) {
    .trends-chart-container {
        height: 790px;
    }
}
</style>
