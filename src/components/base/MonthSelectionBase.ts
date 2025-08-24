import { ref, computed } from 'vue';

import type { Year0BasedMonth } from '@/core/datetime.ts';

import {
    getYear0BasedMonthObjectFromUnixTime,
    getAllowedYearRange,
    getThisMonthFirstUnixTime
} from '@/lib/datetime.ts';

import { useI18n } from '@/locales/helpers.ts';

export interface MonthSelectionValue {
    year: number;
    month: number; // 0-based month (0 = January, 11 = December)
}

export interface CommonMonthSelectionProps {
    modelValue?: Year0BasedMonth;
    title?: string;
    hint?: string;
    show: boolean;
}

function getYearMonthValueFromProps(props: CommonMonthSelectionProps): MonthSelectionValue {
    let value: Year0BasedMonth = getYear0BasedMonthObjectFromUnixTime(getThisMonthFirstUnixTime());

    if (props.modelValue) {
        value = props.modelValue;
    }

    return {
        year: value.year,
        month: value.month0base
    };
}

export function useMonthSelectionBase(props: CommonMonthSelectionProps) {
    const { isLongDateMonthAfterYear } = useI18n();

    const yearRange = ref<number[]>(getAllowedYearRange());
    const monthValue = ref<MonthSelectionValue>(getYearMonthValueFromProps(props));

    const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());

    function getMonthSelectionValue(yearMonth: Year0BasedMonth): MonthSelectionValue | null {
        if (!yearMonth) {
            return null;
        }

        return {
            year: yearMonth.year,
            month: yearMonth.month0base
        };
    }

    function getYear0BasedMonth(): Year0BasedMonth | null {
        if (!monthValue.value) {
            return null;
        }

        if (monthValue.value.year <= 0 || monthValue.value.month < 0) {
            throw new Error('Date is too early');
        }

        return {
            year: monthValue.value.year,
            month0base: monthValue.value.month
        };
    }

    return {
        // states
        yearRange,
        monthValue,
        // computed states
        isYearFirst,
        // functions
        getMonthSelectionValue,
        getYear0BasedMonth
    };
}
