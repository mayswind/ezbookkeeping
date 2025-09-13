<template>
    <vue-date-picker inline auto-apply
                     month-picker
                     :class="monthPickerClass"
                     :dark="isDarkMode"
                     :clearable="!!clearable"
                     :year-range="yearRange"
                     :year-first="isYearFirst"
                     :range="isDateRange ? { partialRange: false } : undefined"
                     v-model="dateTime">
        <!-- @vue-expect-error It seems to be a bug in vue-date-picker (https://github.com/Vuepic/vue-datepicker/issues/1154), when using the month picker, it does not provide the value and text props in the slot, but provides the year. -->
        <template #year="{ year }">
            {{ getDisplayYear(year) }}
        </template>
        <template #year-overlay-value="{ value }">
            {{ getDisplayYear(value) }}
        </template>
        <template #month="{ value }">
            {{ getDisplayMonth(value) }}
        </template>
        <template #month-overlay-value="{ value }">
            {{ getDisplayMonth(value) }}
        </template>
    </vue-date-picker>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import VueDatePicker from '@vuepic/vue-datepicker';

import { useI18n } from '@/locales/helpers.ts';

import type { Year0BasedMonth } from '@/core/datetime.ts';

import { isArray } from '@/lib/common.ts';
import { getAllowedYearRange, getYearMonthDayDateTime } from '@/lib/datetime.ts';

export interface MonthSelectionValue {
    year: number;
    month: number; // 0-based month (0 = January, 11 = December)
}

type SupportedModelValue = Year0BasedMonth | Year0BasedMonth[];

const props = defineProps<{
    modelValue: SupportedModelValue;
    monthPickerClass?: string;
    isDarkMode: boolean;
    clearable?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: SupportedModelValue): void;
}>();

const {
    isLongDateMonthAfterYear,
    getCalendarDisplayShortYearFromUnixTime,
    getCalendarDisplayShortMonthFromUnixTime
} = useI18n();

const yearRange = getAllowedYearRange();

const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());

const dateTime = computed<MonthSelectionValue | MonthSelectionValue[]>({
    get: () => {
        if (isArray(props.modelValue)) {
            return props.modelValue.map(item => getMonthSelectionValueFromYear0BasedMonth(item));
        } else {
            return getMonthSelectionValueFromYear0BasedMonth(props.modelValue);
        }
    },
    set: (value: MonthSelectionValue | MonthSelectionValue[]) => {
        if (isArray(value)) {
            emit('update:modelValue', value.map(item => getYear0BasedMonthFromMonthSelectionValue(item)));
        } else {
            emit('update:modelValue', getYear0BasedMonthFromMonthSelectionValue(value));
        }
    }
});

const isDateRange = computed<boolean>(() => isArray(props.modelValue));

function getMonthSelectionValueFromYear0BasedMonth(value: Year0BasedMonth): MonthSelectionValue {
    return {
        year: value.year,
        month: value.month0base
    };
}

function getYear0BasedMonthFromMonthSelectionValue(value: MonthSelectionValue): Year0BasedMonth {
    return {
        year: value.year,
        month0base: value.month
    };
}

function getDisplayYear(year: number): string {
    return getCalendarDisplayShortYearFromUnixTime(getYearMonthDayDateTime(year, 1, 1).getUnixTime());
}

function getDisplayMonth(month: number): string {
    if (isArray(dateTime.value)) {
        return getCalendarDisplayShortMonthFromUnixTime(getYearMonthDayDateTime(dateTime.value[0]!.year, month + 1, 1).getUnixTime());
    } else {
        return getCalendarDisplayShortMonthFromUnixTime(getYearMonthDayDateTime(dateTime.value.year, month + 1, 1).getUnixTime());
    }
}
</script>
