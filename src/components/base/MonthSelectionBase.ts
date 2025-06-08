import { ref, computed } from 'vue';

import type { Year0BasedMonth } from '@/core/datetime.ts';

import {
    getYear0BasedMonthObjectFromUnixTime,
    getYear0BasedMonthObjectFromString,
    getYearMonthStringFromYear0BasedMonthObject,
    getCurrentYear,
    getThisMonthFirstUnixTime
} from '@/lib/datetime.ts';

import { useI18n } from '@/locales/helpers.ts';

export interface MonthSelectionValue {
    year: number;
    month: number; // 0-based month (0 = January, 11 = December)
}

export interface CommonMonthSelectionProps {
    modelValue?: string;
    title?: string;
    hint?: string;
    show: boolean;
}

function getYearMonthValueFromProps(props: CommonMonthSelectionProps): MonthSelectionValue {
    let value: Year0BasedMonth = getYear0BasedMonthObjectFromUnixTime(getThisMonthFirstUnixTime());

    if (props.modelValue) {
        const yearMonth = getYear0BasedMonthObjectFromString(props.modelValue);

        if (yearMonth) {
            value = yearMonth;
        }
    }

    return {
        year: value.year,
        month: value.month0base
    };
}

export function useMonthSelectionBase(props: CommonMonthSelectionProps) {
    const { isLongDateMonthAfterYear } = useI18n();

    const yearRange = ref<number[]>([
        2000,
        getCurrentYear() + 1
    ]);

    const monthValue = ref<MonthSelectionValue>(getYearMonthValueFromProps(props));

    const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());

    function getMonthSelectionValue(yearMonth: string): MonthSelectionValue | null {
        const yearMonthObj = getYear0BasedMonthObjectFromString(yearMonth);

        if (!yearMonthObj) {
            return null;
        }

        return {
            year: yearMonthObj.year,
            month: yearMonthObj.month0base
        };
    }

    function getTextualYearMonth(): string | null {
        if (!monthValue.value) {
            return null;
        }

        if (monthValue.value.year <= 0 || monthValue.value.month < 0) {
            throw new Error('Date is too early');
        }

        return getYearMonthStringFromYear0BasedMonthObject({
            year: monthValue.value.year,
            month0base: monthValue.value.month
        });
    }

    return {
        // states
        yearRange,
        monthValue,
        // computed states
        isYearFirst,
        // functions
        getMonthSelectionValue,
        getTextualYearMonth
    };
}
