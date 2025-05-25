import { computed } from 'vue';

import { FiscalYearStart } from '@/core/fiscalyear.ts';

import { useUserStore } from '@/stores/user.ts';

import { useI18n } from '@/locales/helpers.ts';

import { arrangeArrayWithNewStartIndex } from '@/lib/common';

export interface FiscalYearStartSelectionBaseProps {
    modelValue?: number;
}

export interface FiscalYearStartSelectionBaseEmits {
    (e: 'update:modelValue', value: number): void;
}

export function useFiscalYearStartSelectionBase(props: FiscalYearStartSelectionBaseProps, emit: FiscalYearStartSelectionBaseEmits) {
    const { getAllMinWeekdayNames,getCurrentFiscalYearStart, formatMonthDayToLongDay, getCurrentFiscalYearStartFormatted } = useI18n();

    const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));

    const displayName = computed<string>(() => {
        let fy = getCurrentFiscalYearStart();

        if (selectedFiscalYearStart.value !== 0 && selectedFiscalYearStart.value !== undefined) {
            const testFy = FiscalYearStart.fromNumber(selectedFiscalYearStart.value);
            if (testFy) {
                fy = testFy;
            }
        }
        
        const monthDay = fy.toMonthDashDayString();
        return formatMonthDayToLongDay(monthDay);
    });

    const disabledDates = (date: Date) => {
        // Disable February 29 (leap day)
        return date.getMonth() === 1 && date.getDate() === 29; 
    };

    const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);

    const selectedFiscalYearStart = computed<number>(() => {
        return props.modelValue !== undefined ? props.modelValue : getCurrentFiscalYearStart().value;
    });

    const userStore = useUserStore();

    function selectedDisplayName(dateString: string): string {
        const fyString = FiscalYearStart.fromMonthDashDayString(dateString);
        if (fyString) {
            const monthDay = fyString.toMonthDashDayString();
            return formatMonthDayToLongDay(monthDay);
        }
        return displayName.value;
    }

    function getModelValueToDateString(): string {
        const input = selectedFiscalYearStart.value;
        
        if (input !== 0 && input !== undefined) {
            const fy = FiscalYearStart.fromNumber(input);
            if (fy) {
                return fy.toMonthDashDayString();
            }
        }
        return getCurrentFiscalYearStartFormatted();
    }

    function setModelValueFromDateString(input: string): number {
        const fyString = FiscalYearStart.fromMonthDashDayString(input);
        if (fyString) {
            return fyString.value;
        }
        return getCurrentFiscalYearStart().value;
    }
    
    return {
        // functions
        getModelValueToDateString,
        setModelValueFromDateString,
        selectedDisplayName,
        // computed states
        dayNames,
        displayName,
        disabledDates,
        firstDayOfWeek,
        selectedFiscalYearStart,
    }
}
