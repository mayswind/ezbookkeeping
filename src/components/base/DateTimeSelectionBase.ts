import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { type NameValue } from '@/core/base.ts';
import { NumeralSystem } from '@/core/numeral.ts';

import {
    getLocalDatetimeFromUnixTime,
    getUnixTimeFromLocalDatetime,
    getSameDateTimeWithBrowserTimezone,
    getSameDateTimeWithTimezoneOffset,
    parseDateTimeFromUnixTimeWithBrowserTimezone,
    parseDateTimeFromUnixTimeWithTimezoneOffset
} from '@/lib/datetime.ts';

export interface TimePickerValue {
    value: string;
    itemsIndex: number;
}

export function useDateTimeSelectionBase() {
    const {
        getAllMeridiemIndicators,
        getCurrentNumeralSystemType,
        isLongTime24HourFormat,
        isLongTimeMeridiemIndicatorFirst,
        isLongTimeHourTwoDigits,
        isLongTimeMinuteTwoDigits,
        isLongTimeSecondTwoDigits
    } = useI18n();

    const is24Hour = ref<boolean>(isLongTime24HourFormat());
    const isHourTwoDigits = ref<boolean>(isLongTimeHourTwoDigits());
    const isMinuteTwoDigits = ref<boolean>(isLongTimeMinuteTwoDigits());
    const isSecondTwoDigits = ref<boolean>(isLongTimeSecondTwoDigits());
    const isMeridiemIndicatorFirst = ref<boolean>(isLongTimeMeridiemIndicatorFirst() || false);

    const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
    const meridiemItems = computed<NameValue[]>(() => getAllMeridiemIndicators());

    function getLocalDatetimeFromSameDateTimeOfUnixTime(unixTime: number, utcOffset: number): Date {
        return getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(parseDateTimeFromUnixTimeWithTimezoneOffset(unixTime, utcOffset)).getUnixTime());
    }

    function getUnixTimeFromSameDateTimeOfLocalDatetime(localDatetime: Date, utcOffset: number): number {
        return getSameDateTimeWithTimezoneOffset(parseDateTimeFromUnixTimeWithBrowserTimezone(getUnixTimeFromLocalDatetime(localDatetime)), utcOffset).getUnixTime();
    }

    function getDisplayTimeValue(value: number, forceTwoDigits: boolean): string {
        let textualValue = value.toString();

        if (forceTwoDigits) {
            textualValue = textualValue.padStart(2, NumeralSystem.WesternArabicNumerals.digitZero);
        }

        return numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(textualValue);
    }

    function generateAllHours(count: number, forceTwoDigits: boolean): TimePickerValue[] {
        const ret: TimePickerValue[] = [];
        const startHour = is24Hour.value ? 0 : 1;
        const endHour = is24Hour.value ? 23 : 11;

        for (let i = 0; i < count; i++) {
            if (!is24Hour.value) {
                ret.push({
                    value: getDisplayTimeValue(12, forceTwoDigits),
                    itemsIndex: i
                });
            }

            for (let j = startHour; j <= endHour; j++) {
                ret.push({
                    value: getDisplayTimeValue(j, forceTwoDigits),
                    itemsIndex: i
                });
            }
        }

        return ret;
    }

    function generateAllMinutesOrSeconds(count: number, forceTwoDigits: boolean): TimePickerValue[] {
        const ret: TimePickerValue[] = [];

        for (let i = 0; i < count; i++) {
            for (let j = 0; j < 60; j++) {
                ret.push({
                    value: getDisplayTimeValue(j, forceTwoDigits),
                    itemsIndex: i
                });
            }
        }

        return ret;
    }

    return {
        // states
        is24Hour,
        isHourTwoDigits,
        isMinuteTwoDigits,
        isSecondTwoDigits,
        isMeridiemIndicatorFirst,
        // computed
        meridiemItems,
        // functions
        getLocalDatetimeFromSameDateTimeOfUnixTime,
        getUnixTimeFromSameDateTimeOfLocalDatetime,
        getDisplayTimeValue,
        generateAllHours,
        generateAllMinutesOrSeconds
    };
}
