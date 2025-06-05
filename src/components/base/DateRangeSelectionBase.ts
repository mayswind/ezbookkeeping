import { ref, computed } from 'vue';

import { type TimeRangeAndDateType, type PresetDateRange, type UnixTimeRange, DateRange } from '@/core/datetime.ts';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import {
    getCurrentUnixTime,
    getCurrentYear,
    getUnixTime,
    getLocalDatetimeFromUnixTime,
    getTodayFirstUnixTime,
    getDummyUnixTimeForLocalUsage,
    getActualUnixTimeForStore,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getDateRangeByDateType
} from '@/lib/datetime.ts';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

export interface CommonDateRangeSelectionProps {
    minTime?: number;
    maxTime?: number;
    title?: string;
    hint?: string;
    show: boolean;
}

function getDateRangeFromProps(props: CommonDateRangeSelectionProps): { minDate: number; maxDate: number } {
    let minDate = getTodayFirstUnixTime();
    let maxDate = getCurrentUnixTime();

    if (props.minTime) {
        minDate = props.minTime;
    }

    if (props.maxTime) {
        maxDate = props.maxTime;
    }

    return {
        minDate,
        maxDate
    };
}

export function useDateRangeSelectionBase(props: CommonDateRangeSelectionProps) {
    const { tt, getAllMinWeekdayNames, formatUnixTimeToLongDateTime, isLongDateMonthAfterYear, isLongTime24HourFormat } = useI18n();
    const userStore = useUserStore();
    const { minDate, maxDate } = getDateRangeFromProps(props);

    const yearRange = ref<number[]>([
        2000,
        getCurrentYear() + 1
    ]);

    const dateRange = ref<Date[]>([
        getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(minDate, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())),
        getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(maxDate, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
    ]);

    const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);
    const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));
    const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());
    const is24Hour = computed<boolean>(() => isLongTime24HourFormat());
    const beginDateTime = computed<string>(() => {
        const actualBeginUnixTime = getActualUnixTimeForStore(getUnixTime(dateRange.value[0]), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes());
        return formatUnixTimeToLongDateTime(actualBeginUnixTime);
    });
    const endDateTime = computed<string>(() => {
        const actualEndUnixTime = getActualUnixTimeForStore(getUnixTime(dateRange.value[1]), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes());
        return formatUnixTimeToLongDateTime(actualEndUnixTime);
    });
    const presetRanges = computed<PresetDateRange[]>(() => {
        const presetRanges:PresetDateRange[] = [];

        [
            DateRange.Today,
            DateRange.LastSevenDays,
            DateRange.LastThirtyDays,
            DateRange.ThisWeek,
            DateRange.ThisMonth,
            DateRange.ThisYear,
            DateRange.LastYear,
            DateRange.ThisFiscalYear,
            DateRange.LastFiscalYear
        ].forEach(dateRangeType => {
            const dateRange = getDateRangeByDateType(dateRangeType.type, firstDayOfWeek.value, userStore.currentUserFiscalYearStart) as TimeRangeAndDateType;

            presetRanges.push({
                label: tt(dateRangeType.name),
                value: [
                    getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(dateRange.minTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())),
                    getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(dateRange.maxTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
                ]
            });
        });

        return presetRanges;
    });

    function getFinalDateRange(): UnixTimeRange | null {
        if (!dateRange.value[0] || !dateRange.value[1]) {
            return null;
        }

        const currentMinDate = dateRange.value[0];
        const currentMaxDate = dateRange.value[1];

        let minUnixTime = getUnixTime(currentMinDate);
        let maxUnixTime = getUnixTime(currentMaxDate);

        if (minUnixTime < 0 || maxUnixTime < 0) {
            throw new Error('Date is too early');
        }

        minUnixTime = getActualUnixTimeForStore(minUnixTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes());
        maxUnixTime = getActualUnixTimeForStore(maxUnixTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes());

        return {
            minUnixTime,
            maxUnixTime
        };
    }

    return {
        // states
        yearRange,
        dateRange,
        // computed states
        dayNames,
        isYearFirst,
        is24Hour,
        beginDateTime,
        endDateTime,
        presetRanges,
        // functions
        getFinalDateRange
    };
}
