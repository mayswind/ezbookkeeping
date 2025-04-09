import { computed } from 'vue';

import { FiscalYearStart } from '@/core/fiscalyear.ts';

import { useI18n } from '@/locales/helpers.ts';

export interface FiscalYearStartSelectionBaseProps {
    modelValue?: number;
}

export interface FiscalYearStartSelectionBaseEmits {
    (e: 'update:modelValue', value: number): void;
}

export function useFiscalYearStartSelectionBase(props: FiscalYearStartSelectionBaseProps, emit?: FiscalYearStartSelectionBaseEmits) {
    const { getCurrentFiscalYearStart, formatMonthDayToLongDate } = useI18n();

    const effectiveModelValue = computed<number>(() => {
        return props.modelValue !== undefined ? props.modelValue : getCurrentFiscalYearStart().value;
    });

    function getterModelValue(input?: number): string {
        const valueToUse = input !== undefined ? input : effectiveModelValue.value;
        
        if (valueToUse !== 0 && valueToUse !== undefined) {
            const fy = FiscalYearStart.fromNumber(valueToUse);
            if (fy) {
                return fy.toMonthDashDayString();
            }
        }
        return getCurrentFiscalYearStart().toMonthDashDayString();
    }

    function setterModelValue(input: string): number {
        const fyString = FiscalYearStart.fromMonthDashDayString(input);
        if (fyString) {
            return fyString.value;
        }
        return getCurrentFiscalYearStart().value;
    }
    
    const displayName = computed<string>(() => {
        let fy = getCurrentFiscalYearStart();

        if (effectiveModelValue.value !== 0 && effectiveModelValue.value !== undefined) {
            const testFy = FiscalYearStart.fromNumber(effectiveModelValue.value);
            if (testFy) {
                fy = testFy;
            }
        }
        
        const monthDay = fy.toMonthDashDayString();
        return formatMonthDayToLongDate(monthDay);
    });

    const disabledDates = (date: Date) => {
        // Disable February 29 (leap day)
        return date.getMonth() === 1 && date.getDate() === 29; 
    };

    const selectedDate = computed<string>({
        get: () => {
            if (props.modelValue === undefined) {
                return getCurrentFiscalYearStart().toMonthDashDayString();
            }
            return getterModelValue();
        },
        set: (value: string) => {
            if (emit) {
                const numericValue = setterModelValue(value);
                emit('update:modelValue', numericValue);
            }
        }
    });

    return {
        // computed states
        displayName,
        disabledDates,
        effectiveModelValue,
        selectedDate,
        // functions
        getterModelValue,
        setterModelValue,
    }
}
