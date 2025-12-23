import { ref, computed } from 'vue';

import type { TextualYearMonth, Year0BasedMonth } from '@/core/datetime.ts';

import {
    getYear0BasedMonthObjectFromUnixTime,
    getYear0BasedMonthObjectFromString,
    getYearMonthStringFromYear0BasedMonthObject,
    getCurrentUnixTime,
    parseDateTimeFromUnixTime,
    getThisYearFirstUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime
} from '@/lib/datetime.ts';

import { useI18n } from '@/locales/helpers.ts';

export interface CommonMonthRangeSelectionProps {
    minTime?: TextualYearMonth;
    maxTime?: TextualYearMonth;
    title?: string;
    hint?: string;
    show: boolean;
}

function getMonthRangeFromProps(props: CommonMonthRangeSelectionProps): { minDate: Year0BasedMonth; maxDate: Year0BasedMonth } {
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
        minDate: minDate,
        maxDate: maxDate
    };
}

export function useMonthRangeSelectionBase(props: CommonMonthRangeSelectionProps) {
    const { formatDateTimeToGregorianLikeLongYearMonth } = useI18n();
    const { minDate, maxDate } = getMonthRangeFromProps(props);

    const dateRange = ref<Year0BasedMonth[]>([
        minDate,
        maxDate
    ]);

    const beginDateTime = computed<string>(() => formatDateTimeToGregorianLikeLongYearMonth(parseDateTimeFromUnixTime(getYearMonthFirstUnixTime(dateRange.value[0] as Year0BasedMonth))));
    const endDateTime = computed<string>(() => formatDateTimeToGregorianLikeLongYearMonth(parseDateTimeFromUnixTime(getYearMonthLastUnixTime(dateRange.value[1] as Year0BasedMonth))));

    function getFinalMonthRange(): { minYearMonth: TextualYearMonth | '', maxYearMonth: TextualYearMonth | '' } | null {
        if (!dateRange.value[0] || !dateRange.value[1]) {
            return null;
        }

        if (dateRange.value[0].year <= 0 || dateRange.value[0].month0base < 0 || dateRange.value[1].year <= 0 || dateRange.value[1].month0base < 0) {
            throw new Error('Date is too early');
        }

        const minYearMonth = getYearMonthStringFromYear0BasedMonthObject(dateRange.value[0] as Year0BasedMonth);
        const maxYearMonth = getYearMonthStringFromYear0BasedMonthObject(dateRange.value[1] as Year0BasedMonth);

        return {
            minYearMonth,
            maxYearMonth
        };
    }

    return {
        // states
        dateRange,
        // computed states
        beginDateTime,
        endDateTime,
        // functions
        getFinalMonthRange
    };
}
