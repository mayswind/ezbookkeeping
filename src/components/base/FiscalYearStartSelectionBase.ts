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
    const { getAllMinWeekdayNames, formatMonthDayToLongDay, getCurrentFiscalYearStartFormatted } = useI18n();

    const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));

    const displayName = computed<string>(() => {
        let fy = FiscalYearStart.fromNumber(selectedFiscalYearStart.value);

        if ( fy ) {
            return formatMonthDayToLongDay(fy.toMonthDashDayString())
        }

        return formatMonthDayToLongDay(FiscalYearStart.strictFromNumber(userStore.currentUserFiscalYearStart).toMonthDashDayString());

    });

    const disabledDates = (date: Date) => {
        // Disable February 29 (leap day)
        return date.getMonth() === 1 && date.getDate() === 29; 
    };

    const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);

    const selectedFiscalYearStart = computed<number>(() => {
        return props.modelValue !== undefined ? props.modelValue : userStore.currentUserFiscalYearStart;
    });

    const userStore = useUserStore();

    function selectedDisplayName(dateString: string): string {
        let fy = FiscalYearStart.fromMonthDashDayString(dateString);
        if ( fy ) {
            return formatMonthDayToLongDay(fy.toMonthDashDayString());
        }
        return displayName.value;
    }

    function getModelValueToDateString(): string {
        const input = selectedFiscalYearStart.value;
        
        let fy = FiscalYearStart.fromNumber(input);

        if ( fy ) {
            return fy.toMonthDashDayString();
        }

        return getCurrentFiscalYearStartFormatted();
    }

    function getDateStringToModelValue(input: string): number {
        const fyString = FiscalYearStart.fromMonthDashDayString(input);
        if (fyString) {
            return fyString.value;
        }
        return userStore.currentUserFiscalYearStart;
    }
    
    return {
        // functions
        getDateStringToModelValue,
        getModelValueToDateString,
        selectedDisplayName,
        // computed states
        dayNames,
        displayName,
        disabledDates,
        firstDayOfWeek,
        selectedFiscalYearStart,
    }
}
