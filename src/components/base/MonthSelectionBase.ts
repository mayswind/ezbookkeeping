import { ref, computed } from 'vue';

import type { YearMonth } from '@/core/datetime.ts';

import {
    getYearMonthObjectFromUnixTime,
    getYearMonthObjectFromString,
    getYearMonthStringFromObject,
    getCurrentYear,
    getThisMonthFirstUnixTime
} from '@/lib/datetime.ts';

import { useI18n } from '@/locales/helpers.ts';

export interface CommonMonthSelectionProps {
    modelValue?: string;
    title?: string;
    hint?: string;
    show: boolean;
}

function getYearMonthValueFromProps(props: CommonMonthSelectionProps): YearMonth {
    let value: YearMonth = getYearMonthObjectFromUnixTime(getThisMonthFirstUnixTime());

    if (props.modelValue) {
        const yearMonth = getYearMonthObjectFromString(props.modelValue);

        if (yearMonth) {
            value = yearMonth;
        }
    }

    return value;
}

export function useMonthSelectionBase(props: CommonMonthSelectionProps) {
    const { isLongDateMonthAfterYear } = useI18n();

    const yearRange = ref<number[]>([
        2000,
        getCurrentYear() + 1
    ]);

    const monthValue = ref<YearMonth>(getYearMonthValueFromProps(props));

    const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());

    function getTextualYearMonth(): string | null {
        if (!monthValue.value) {
            return null;
        }

        if (monthValue.value.year <= 0 || monthValue.value.month < 0) {
            throw new Error('Date is too early');
        }

        return getYearMonthStringFromObject(monthValue.value);
    }

    return {
        // states
        yearRange,
        monthValue,
        // computed states
        isYearFirst,
        // functions
        getTextualYearMonth
    };
}
