import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type {
    TextualYearMonth,
    Year1BasedMonth,
    YearMonthDay,
    TimeRangeAndDateType,
    YearUnixTime,
    YearQuarterUnixTime,
    YearMonthUnixTime,
    YearMonthDayUnixTime
} from '@/core/datetime.ts';
import type { FiscalYearUnixTime } from '@/core/fiscalyear.ts';
import { ChartDataAggregationType, ChartDateAggregationType } from '@/core/statistics.ts';
import type { YearMonthItems, YearMonthDayItems } from '@/models/transaction.ts';

import {
    getYearMonthDayDateTime,
    getGregorianCalendarYearAndMonthFromUnixTime,
    getAllDaysStartAndEndUnixTimes
} from '@/lib/datetime.ts';
import {
    getAllDateRangesFromItems,
    getAllDateRangesByYearMonthRange
} from '@/lib/statistics.ts';

export type TrendsChartDateType = 'daily' | 'monthly';

interface TrendsChartTypes {
    daily: {
        ItemsType: YearMonthDayItems<YearMonthDay>;
        DateTimeRangeType: number;
        MonthRangeType: undefined;
    };
    monthly: {
        ItemsType: YearMonthItems<Year1BasedMonth>;
        DateTimeRangeType: undefined;
        MonthRangeType: TextualYearMonth | '';
    };
}

export interface CommonTrendsChartProps<T extends TrendsChartDateType> {
    chartMode: T;
    items: TrendsChartTypes[T]['ItemsType'][];
    stacked?: boolean;
    startTime: TrendsChartTypes[T]['DateTimeRangeType'];
    endTime: TrendsChartTypes[T]['DateTimeRangeType'];
    startYearMonth: TrendsChartTypes[T]['MonthRangeType'];
    endYearMonth: TrendsChartTypes[T]['MonthRangeType'];
    fiscalYearStart: number;
    sortingType: number;
    dataAggregationType: ChartDataAggregationType;
    dateAggregationType: number;
    idField?: string;
    nameField: string;
    valueField: string;
    colorField?: string;
    hiddenField?: string;
    displayOrdersField?: string;
    translateName?: boolean;
    defaultCurrency?: string;
    enableClickItem?: boolean;
}

export interface TrendsBarChartClickEvent {
    itemId: string;
    dateRange: TimeRangeAndDateType;
}

function buildDailyAllDateRanges(props: CommonTrendsChartProps<'daily'>): YearUnixTime[] | FiscalYearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] | YearMonthDayUnixTime[] {
    let startTime: number = props.startTime;
    let endTime: number = props.endTime;

    if ((!startTime || !endTime) && props.items && props.items.length) {
        let minUnixTime = Number.MAX_SAFE_INTEGER, maxUnixTime = 0;

        for (const accountItem of props.items) {
            for (const dataItem of accountItem.items) {
                const dateTime = getYearMonthDayDateTime(dataItem.year, dataItem.month, dataItem.day);
                const unixTime = dateTime.getUnixTime();

                if (unixTime < minUnixTime) {
                    minUnixTime = unixTime;
                }

                if (unixTime > maxUnixTime) {
                    maxUnixTime = unixTime;
                }
            }
        }

        if (minUnixTime < Number.MAX_SAFE_INTEGER && maxUnixTime > 0) {
            startTime = minUnixTime;
            endTime = maxUnixTime;
        }
    }

    if (props.dateAggregationType === ChartDateAggregationType.Day.type) {
        return getAllDaysStartAndEndUnixTimes(startTime, endTime);
    } else {
        const startYearMonth = getGregorianCalendarYearAndMonthFromUnixTime(startTime);
        const endYearMonth = getGregorianCalendarYearAndMonthFromUnixTime(endTime);
        return getAllDateRangesByYearMonthRange(startYearMonth, endYearMonth, props.fiscalYearStart, props.dateAggregationType);
    }
}

function buildMonthlyAllDateRanges(props: CommonTrendsChartProps<'monthly'>): YearUnixTime[] | FiscalYearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] {
    return getAllDateRangesFromItems(props.items, props.startYearMonth, props.endYearMonth, props.fiscalYearStart, props.dateAggregationType);
}

export function useTrendsChartBase<T extends TrendsChartDateType>(props: CommonTrendsChartProps<T>) {
    const { tt } = useI18n();

    const allDateRanges = computed<YearUnixTime[] | FiscalYearUnixTime[] | YearQuarterUnixTime[] | YearMonthUnixTime[] | YearMonthDayUnixTime[]>(() => {
        if (props.chartMode === 'daily') {
            return buildDailyAllDateRanges(props as CommonTrendsChartProps<'daily'>);
        } else if (props.chartMode === 'monthly') {
            return buildMonthlyAllDateRanges(props as CommonTrendsChartProps<'monthly'>);
        } else {
            return [];
        }
    });

    function getItemName(name: string): string {
        return props.translateName ? tt(name) : name;
    }

    return {
        // computed states
        allDateRanges,
        // functions
        getItemName
    };
}
