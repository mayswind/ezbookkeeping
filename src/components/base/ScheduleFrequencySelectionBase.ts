import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { sortNumbersArray } from '@/lib/common.ts';

export interface CommonScheduleFrequencySelectionProps {
    type: number;
    modelValue: string;
    disabled?: boolean;
    readonly?: boolean;
    label?: string;
}

export interface AvailableMonthDay {
    day: number;
    displayName: string;
}

export function useScheduleFrequencySelectionBase() {
    const { getAllWeekDays, getAllTransactionScheduledFrequencyTypes, getMonthdayShortName } = useI18n();
    const userStore = useUserStore();

    const allTransactionScheduledFrequencyTypes = computed<TypeAndDisplayName[]>(() => getAllTransactionScheduledFrequencyTypes());
    const allWeekDays = computed<TypeAndDisplayName[]>(() => getAllWeekDays(userStore.currentUserFirstDayOfWeek));

    const allAvailableMonthDays = computed<AvailableMonthDay[]>(() => {
        const allAvailableDays = [];

        for (let i = 1; i <= 28; i++) {
            allAvailableDays.push({
                day: i,
                displayName: getMonthdayShortName(i),
            });
        }

        return allAvailableDays;
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
        // functions
        getFrequencyValues
    };
}
