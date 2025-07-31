import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { type NameValue } from '@/core/base.ts';
import { type WeekDayValue } from '@/core/datetime.ts';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import { getCurrentYear } from '@/lib/datetime.ts';

export interface TimePickerValue {
    value: string;
    itemsIndex: number;
}

export function useDateTimeSelectionBase() {
    const {
        getAllMinWeekdayNames,
        getAllMeridiemIndicators,
        isLongDateMonthAfterYear,
        isLongTime24HourFormat,
        isLongTimeMeridiemIndicatorFirst,
        isLongTimeHourTwoDigits,
        isLongTimeMinuteTwoDigits,
        isLongTimeSecondTwoDigits
    } = useI18n();

    const userStore = useUserStore();

    const is24Hour = ref<boolean>(isLongTime24HourFormat());
    const isHourTwoDigits = ref<boolean>(isLongTimeHourTwoDigits());
    const isMinuteTwoDigits = ref<boolean>(isLongTimeMinuteTwoDigits());
    const isSecondTwoDigits = ref<boolean>(isLongTimeSecondTwoDigits());
    const isMeridiemIndicatorFirst = ref<boolean>(isLongTimeMeridiemIndicatorFirst() || false);

    const yearRange = ref<number[]>([
        2000,
        getCurrentYear() + 1
    ]);

    const meridiemItems = computed<NameValue[]>(() => getAllMeridiemIndicators());

    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));
    const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());

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
        yearRange,
        // computed
        meridiemItems,
        firstDayOfWeek,
        dayNames,
        isYearFirst,
        // functions
        getDisplayTimeValue,
        generateAllHours,
        generateAllMinutesOrSeconds
    };
}
