import { ref, computed } from 'vue';

import { type WeekDayValue } from '@/core/datetime.ts';
import { FiscalYearStart } from '@/core/fiscalyear.ts';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

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
    const {
        getAllMinWeekdayNames,
        formatMonthDayToLongDay
    } = useI18n();

    const userStore = useUserStore();

    const disabledDates = (date: Date) => {
        // Disable February 29 (leap day)
        return date.getMonth() === 1 && date.getDate() === 29;
    };

    const selectedFiscalYearStart = ref<number>(getFiscalYearStartFromProps(props));

    const selectedFiscalYearStartValue = computed<string>({
        get: () => {
            const fiscalYearStart = FiscalYearStart.valueOf(selectedFiscalYearStart.value);

            if (fiscalYearStart) {
                return fiscalYearStart.toMonthDashDayString();
            } else {
                return FiscalYearStart.Default.toMonthDashDayString();
            }
        },
        set: (value: string) => {
            const fiscalYearStart = FiscalYearStart.parse(value);

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

        return formatMonthDayToLongDay(fiscalYearStart.toMonthDashDayString());
    });

    const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
    const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));

    return {
        // constants
        disabledDates,
        // states,
        selectedFiscalYearStart,
        // computed states
        selectedFiscalYearStartValue,
        displayFiscalYearStartDate,
        firstDayOfWeek,
        dayNames
    };
}
