import { ref, computed } from 'vue';

import type { TextualYearMonth, Year0BasedMonth } from '@/core/datetime.ts';

import {
    getYear0BasedMonthObjectFromUnixTime,
    getYear0BasedMonthObjectFromString,
    getYearMonthStringFromYear0BasedMonthObject,
    getCurrentUnixTime,
    getAllowedYearRange,
    getThisYearFirstUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime
} from '@/lib/datetime.ts';

import { useI18n } from '@/locales/helpers.ts';

export interface MonthSelectionValue {
    year: number;
    month: number; // 0-based month (0 = January, 11 = December)
}

export interface CommonMonthRangeSelectionProps {
    minTime?: TextualYearMonth;
    maxTime?: TextualYearMonth;
    title?: string;
    hint?: string;
    show: boolean;
}

function getMonthRangeFromProps(props: CommonMonthRangeSelectionProps): { minDate: MonthSelectionValue; maxDate: MonthSelectionValue } {
    let minDate: Year0BasedMonth = getYear0BasedMonthObjectFromUnixTime(getThisYearFirstUnixTime());
    let maxDate: Year0BasedMonth = getYear0BasedMonthObjectFromUnixTime(getCurrentUnixTime());

    if (props.minTime) {
        const yearMonth = getYear0BasedMonthObjectFromString(props.minTime);

        if (yearMonth) {
            minDate = yearMonth;
        }
    }

    if (props.maxTime) {
        const yearMonth = getYear0BasedMonthObjectFromString(props.maxTime);

        if (yearMonth) {
            maxDate = yearMonth;
        }
    }

    return {
        minDate: {
            year: minDate.year,
            month: minDate.month0base
        },
        maxDate: {
            year: maxDate.year,
            month: maxDate.month0base
        }
    };
}

export function useMonthRangeSelectionBase(props: CommonMonthRangeSelectionProps) {
    const { formatUnixTimeToLongYearMonth, isLongDateMonthAfterYear } = useI18n();
    const { minDate, maxDate } = getMonthRangeFromProps(props);

    const yearRange = ref<number[]>(getAllowedYearRange());
    const dateRange = ref<MonthSelectionValue[]>([
        minDate,
        maxDate
    ]);

    const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());
    const beginDateTime = computed<string>(() => formatUnixTimeToLongYearMonth(getYearMonthFirstUnixTime({
        year: dateRange.value[0].year,
        month0base: dateRange.value[0].month
    })));
    const endDateTime = computed<string>(() => formatUnixTimeToLongYearMonth(getYearMonthLastUnixTime({
        year: dateRange.value[1].year,
        month0base: dateRange.value[1].month
    })));

    function getMonthSelectionValue(yearMonth: TextualYearMonth): MonthSelectionValue | null {
        const yearMonthObj = getYear0BasedMonthObjectFromString(yearMonth);

        if (!yearMonthObj) {
            return null;
        }

        return {
            year: yearMonthObj.year,
            month: yearMonthObj.month0base
        };
    }

    function getFinalMonthRange(): { minYearMonth: TextualYearMonth | '', maxYearMonth: TextualYearMonth | '' } | null {
        if (!dateRange.value[0] || !dateRange.value[1]) {
            return null;
        }

        if (dateRange.value[0].year <= 0 || dateRange.value[0].month < 0 || dateRange.value[1].year <= 0 || dateRange.value[1].month < 0) {
            throw new Error('Date is too early');
        }

        const minYearMonth = getYearMonthStringFromYear0BasedMonthObject({
            year: dateRange.value[0].year,
            month0base: dateRange.value[0].month
        });
        const maxYearMonth = getYearMonthStringFromYear0BasedMonthObject({
            year: dateRange.value[1].year,
            month0base: dateRange.value[1].month
        });

        return {
            minYearMonth,
            maxYearMonth
        };
    }

    return {
        // states
        yearRange,
        dateRange,
        // computed states
        isYearFirst,
        beginDateTime,
        endDateTime,
        // functions
        getMonthSelectionValue,
        getFinalMonthRange
    };
}
