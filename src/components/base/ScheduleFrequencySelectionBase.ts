import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { type TypeAndDisplayName, itemAndIndex } from '@/core/base.ts';
import { type DateTime } from '@/core/datetime.ts';

import { sortNumbersArray } from '@/lib/common.ts';
import { getCurrentDateTime } from '@/lib/datetime.ts';

export interface CommonScheduleFrequencySelectionProps {
    type: number;
    modelValue: string;
    disabled?: boolean;
    readonly?: boolean;
    label?: string;
}

export function useScheduleFrequencySelectionBase() {
    const {
        getAllWeekDays,
        getAvailableMonthDays,
        getAllTransactionScheduledFrequencyTypes,
        formatDateTimeToLongMonthDay
    } = useI18n();
    const userStore = useUserStore();

    const allTransactionScheduledFrequencyTypes = computed<TypeAndDisplayName[]>(() => getAllTransactionScheduledFrequencyTypes());
    const allWeekDays = computed<TypeAndDisplayName[]>(() => getAllWeekDays(userStore.currentUserFirstDayOfWeek));

    const allAvailableMonthDays = computed<TypeAndDisplayName[]>(() => getAvailableMonthDays(28, 3));
    const allAvailableMonthAndDays = computed<TypeAndDisplayName[]>(() => {
        const ret: TypeAndDisplayName[] = [];
        const now: DateTime = getCurrentDateTime();
        const maxDaysOfMonth: number[] = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];

        for (const [days, index] of itemAndIndex(maxDaysOfMonth)) {
            const month = index + 1;

            for (let day = 1; day <= days; day++) {
                const dateTime = now.set({
                    month: month,
                    dayOfMonth: day,
                    hour: 0,
                    minute: 0,
                    second: 0,
                    millisecond: 0
                });

                ret.push({
                    type: month * 100 + day,
                    displayName: formatDateTimeToLongMonthDay(dateTime)
                });
            }
        }

        return ret;
    });

    function getFrequencyValues(value: string): number[] {
        const values = value.split(',');
        const ret: number[] = [];

        for (const value of values) {
            if (value) {
                ret.push(parseInt(value));
            }
        }

        return sortNumbersArray(ret);
    }

    return {
        // computed states
        allTransactionScheduledFrequencyTypes,
        allWeekDays,
        allAvailableMonthDays,
        allAvailableMonthAndDays,
        // functions
        getFrequencyValues
    };
}
