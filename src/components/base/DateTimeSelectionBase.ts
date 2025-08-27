import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { type NameValue } from '@/core/base.ts';

export interface TimePickerValue {
    value: string;
    itemsIndex: number;
}

export function useDateTimeSelectionBase() {
    const {
        getAllMeridiemIndicators,
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

    const meridiemItems = computed<NameValue[]>(() => getAllMeridiemIndicators());

    function getDisplayTimeValue(value: number, forceTwoDigits: boolean): string {
        if (forceTwoDigits && value < 10) {
            return `0${value}`;
        } else {
            return value.toString();
        }
    }

    function generateAllHours(count: number, forceTwoDigits: boolean): TimePickerValue[] {
        const ret: TimePickerValue[] = [];
        const startHour = is24Hour.value ? 0 : 1;
        const endHour = is24Hour.value ? 23 : 11;

        for (let i = 0; i < count; i++) {
            if (!is24Hour.value) {
                ret.push({
                    value: '12',
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
        getDisplayTimeValue,
        generateAllHours,
        generateAllMinutesOrSeconds
    };
}
