import { ref, computed } from 'vue';

import {
    type DateTime,
    type UnixTimeRange,
    type TimeRangeAndDateType,
    type PresetDateRange,
    type WeekDayValue,
    DateRange,
} from '@/core/datetime.ts';

import {
    getCurrentUnixTime,
    getLocalDatetimeFromUnixTime,
    getUnixTimeFromLocalDatetime,
    getTodayFirstUnixTime,
    getSameDateTimeWithCurrentTimezone,
    getSameDateTimeWithBrowserTimezone,
    parseDateTimeFromUnixTime,
    parseDateTimeFromUnixTimeWithBrowserTimezone,
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
    const { tt, formatDateTimeToLongDateTime } = useI18n();
    const userStore = useUserStore();
    const { minDate, maxDate } = getDateRangeFromProps(props);

    const dateRange = ref<Date[]>([
        getLocalDatetimeFromSameDateTimeOfUnixTime(minDate),
        getLocalDatetimeFromSameDateTimeOfUnixTime(maxDate)
    ]);

    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const beginDateTime = computed<string>(() => {
        return formatDateTimeToLongDateTime(getDateTimeFromSameDateTimeOfLocalDatetime(dateRange.value[0] as Date));
    });
    const endDateTime = computed<string>(() => {
        return formatDateTimeToLongDateTime(getDateTimeFromSameDateTimeOfLocalDatetime(dateRange.value[1] as Date));
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
                    getLocalDatetimeFromSameDateTimeOfUnixTime(dateRange.minTime),
                    getLocalDatetimeFromSameDateTimeOfUnixTime(dateRange.maxTime)
                ]
            });
        });

        return presetRanges;
    });

    function getLocalDatetimeFromSameDateTimeOfUnixTime(unixTime: number): Date {
        return getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(parseDateTimeFromUnixTime(unixTime)).getUnixTime());
    }

    function getDateTimeFromSameDateTimeOfLocalDatetime(localDatetime: Date): DateTime {
        return getSameDateTimeWithCurrentTimezone(parseDateTimeFromUnixTimeWithBrowserTimezone(getUnixTimeFromLocalDatetime(localDatetime)));
    }

    function getFinalDateRange(): UnixTimeRange | null {
        if (!dateRange.value[0] || !dateRange.value[1]) {
            return null;
        }

        const currentMinDate = dateRange.value[0];
        const currentMaxDate = dateRange.value[1];

        const minUnixTime = getDateTimeFromSameDateTimeOfLocalDatetime(currentMinDate).getUnixTime();
        const maxUnixTime = getDateTimeFromSameDateTimeOfLocalDatetime(currentMaxDate).getUnixTime();

        if (minUnixTime < 0 || maxUnixTime < 0) {
            throw new Error('Date is too early');
        }

        return {
            minUnixTime,
            maxUnixTime
        };
    }

    return {
        // states
        dateRange,
        // computed states
        beginDateTime,
        endDateTime,
        presetRanges,
        // functions
        getLocalDatetimeFromSameDateTimeOfUnixTime,
        getDateTimeFromSameDateTimeOfLocalDatetime,
        getFinalDateRange
    };
}
