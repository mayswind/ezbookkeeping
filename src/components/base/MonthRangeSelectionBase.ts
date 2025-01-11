import { ref, computed } from 'vue';

import type { YearMonth } from '@/core/datetime.ts';

import {
    getYearMonthObjectFromUnixTime,
    getYearMonthObjectFromString,
    getYearMonthStringFromObject,
    getCurrentUnixTime,
    getCurrentYear,
    getThisYearFirstUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime
} from '@/lib/datetime.ts';

import { useI18n } from '@/locales/helpers.ts';

export interface CommonMonthRangeSelectionProps {
    minTime?: string;
    maxTime?: string;
    title?: string;
    hint?: string;
    show: boolean;
}

function getMonthRangeFromProps(props: CommonMonthRangeSelectionProps): { minDate: YearMonth; maxDate: YearMonth } {
    let minDate: YearMonth = getYearMonthObjectFromUnixTime(getThisYearFirstUnixTime());
    let maxDate: YearMonth = getYearMonthObjectFromUnixTime(getCurrentUnixTime());

    if (props.minTime) {
        const yearMonth = getYearMonthObjectFromString(props.minTime);

        if (yearMonth) {
            minDate = yearMonth;
        }
    }

    if (props.maxTime) {
        const yearMonth = getYearMonthObjectFromString(props.maxTime);

        if (yearMonth) {
            maxDate = yearMonth;
        }
    }

    return {
        minDate,
        maxDate
    };
}

export function useMonthRangeSelectionBase(props: CommonMonthRangeSelectionProps) {
    const { formatUnixTimeToLongYearMonth, isLongDateMonthAfterYear } = useI18n();
    const { minDate, maxDate } = getMonthRangeFromProps(props);

    const yearRange = ref<number[]>([
        2000,
        getCurrentYear() + 1
    ]);

    const dateRange = ref<YearMonth[]>([
        minDate,
        maxDate
    ]);

    const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());
    const beginDateTime = computed<string>(() => formatUnixTimeToLongYearMonth(getYearMonthFirstUnixTime(dateRange.value[0])));
    const endDateTime = computed<string>(() => formatUnixTimeToLongYearMonth(getYearMonthLastUnixTime(dateRange.value[1])));

    function getFinalMonthRange(): { minYearMonth: string, maxYearMonth: string } | null {
        if (!dateRange.value[0] || !dateRange.value[1]) {
            return null;
        }

        if (dateRange.value[0].year <= 0 || dateRange.value[0].month < 0 || dateRange.value[1].year <= 0 || dateRange.value[1].month < 0) {
            throw new Error('Date is too early');
        }

        const minYearMonth = getYearMonthStringFromObject(dateRange.value[0]);
        const maxYearMonth = getYearMonthStringFromObject(dateRange.value[1]);

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
        getFinalMonthRange
    };
}
