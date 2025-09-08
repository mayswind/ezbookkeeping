import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type { MonthDay } from '@/core/datetime.ts';
import { FiscalYearStart } from '@/core/fiscalyear.ts';

import {
    getLocalDatetimeFromUnixTime,
    getThisYearFirstUnixTime,
    getThisYearLastUnixTime
} from '@/lib/datetime.ts';

export interface CommonFiscalYearStartSelectionProps {
    modelValue?: number;
    disabled?: boolean;
    readonly?: boolean;
    label?: string;
}

export interface CommonFiscalYearStartSelectionEmits {
    (e: 'update:modelValue', value: number): void;
}

function getFiscalYearStartFromProps(props: CommonFiscalYearStartSelectionProps): number {
    if (!props.modelValue) {
        return FiscalYearStart.Default.value;
    }

    const fiscalYearStart = FiscalYearStart.valueOf(props.modelValue);

    if (!fiscalYearStart) {
        return FiscalYearStart.Default.value;
    }

    return fiscalYearStart.value;
}

export function useFiscalYearStartSelectionBase(props: CommonFiscalYearStartSelectionProps) {
    const { formatGregorianTextualMonthDayToGregorianLikeLongMonthDay } = useI18n();

    const disabledDates = (date: Date) => {
        // Disable February 29 (leap day)
        return date.getMonth() === 1 && date.getDate() === 29;
    };

    const selectedFiscalYearStart = ref<number>(getFiscalYearStartFromProps(props));

    const selectedFiscalYearStartValue = computed<Date>({
        get: () => {
            const fiscalYearStart = FiscalYearStart.valueOf(selectedFiscalYearStart.value);
            const monthDay: MonthDay = fiscalYearStart?.toMonthDay() ?? FiscalYearStart.Default.toMonthDay();

            return new Date(new Date().getFullYear(), monthDay.month - 1, monthDay.day);
        },
        set: (value: Date) => {
            const fiscalYearStart = FiscalYearStart.of(value.getMonth() + 1, value.getDate());

            if (fiscalYearStart) {
                selectedFiscalYearStart.value = fiscalYearStart.value;
            } else {
                selectedFiscalYearStart.value = FiscalYearStart.Default.value;
            }
        }
    });

    const displayFiscalYearStartDate = computed<string>(() => {
        let fiscalYearStart = FiscalYearStart.valueOf(selectedFiscalYearStart.value);

        if (!fiscalYearStart) {
            fiscalYearStart = FiscalYearStart.Default;
        }

        return formatGregorianTextualMonthDayToGregorianLikeLongMonthDay(fiscalYearStart.toMonthDashDayString());
    });

    const allowedMinDate = computed<Date>(() => getLocalDatetimeFromUnixTime(getThisYearFirstUnixTime()));
    const allowedMaxDate = computed<Date>(() => getLocalDatetimeFromUnixTime(getThisYearLastUnixTime()));

    return {
        // constants
        disabledDates,
        // states,
        selectedFiscalYearStart,
        // computed states
        selectedFiscalYearStartValue,
        displayFiscalYearStartDate,
        allowedMinDate,
        allowedMaxDate
    };
}
