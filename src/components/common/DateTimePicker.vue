<template>
    <vue-date-picker ref="datetimepicker"
                     inline auto-apply
                     enable-seconds
                     six-weeks="center"
                     :class="datetimePickerClass"
                     :config="noSwipeAndScroll ? { noSwipe: true } : undefined"
                     :dark="isDarkMode"
                     :vertical="vertical"
                     :enable-time-picker="enableTimePicker"
                     :disable-year-select="disableYearSelect"
                     :clearable="!!clearable"
                     :year-range="yearRange"
                     :day-names="dayNames"
                     :week-start="firstDayOfWeek"
                     :year-first="isYearFirst"
                     :is24="is24Hour"
                     :min-date="minDate"
                     :max-date="maxDate"
                     :disabled-dates="disabledDates"
                     :month-change-on-scroll="!noSwipeAndScroll"
                     :range="isDateRange ? { partialRange: false } : undefined"
                     :preset-dates="presetRanges"
                     v-model="dateTime">
        <template #year="{ value }">
            {{ getDisplayYear(value) }}
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
        <template #day="{ date }">
            {{ getDisplayDay(date) }}
        </template>
        <template #am-pm-button="{ toggle, value }">
            <button class="dp__pm_am_button" tabindex="0" @click="toggle">{{ tt(`datetime.${value}.content`) }}</button>
        </template>
    </vue-date-picker>
</template>

<script setup lang="ts">
import { computed, useTemplateRef } from 'vue';
import VueDatePicker, { type MenuView } from '@vuepic/vue-datepicker';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { type PresetDateRange, type WeekDayValue } from '@/core/datetime.ts';
import { isArray, arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import { getAllowedYearRange, getYearMonthDayDateTime } from '@/lib/datetime.ts';

type VueDatePickerType = InstanceType<typeof VueDatePicker>;
type SupportedModelValue = Date | Date[] | null;

const props = defineProps<{
    modelValue: SupportedModelValue;
    datetimePickerClass?: string;
    isDarkMode: boolean;
    enableTimePicker: boolean;
    disableYearSelect?: boolean;
    vertical?: boolean;
    noSwipeAndScroll?: boolean;
    clearable?: boolean;
    minDate?: Date;
    maxDate?: Date;
    disabledDates?: (date: Date) => boolean;
    presetRanges?: PresetDateRange[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: SupportedModelValue): void;
}>();

const {
    tt,
    getAllMinWeekdayNames,
    isLongDateMonthAfterYear,
    isLongTime24HourFormat,
    getCalendarShortYearFromUnixTime,
    getCalendarShortMonthFromUnixTime,
    getCalendarDayOfMonthFromUnixTime
} = useI18n();

const userStore = useUserStore();

const datetimepicker = useTemplateRef<VueDatePickerType>('datetimepicker');

const yearRange = getAllowedYearRange();

const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));
const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());
const is24Hour = computed<boolean>(() => isLongTime24HourFormat());

const dateTime = computed<SupportedModelValue>({
    get: () => props.modelValue,
    set: (value: SupportedModelValue) => emit('update:modelValue', value)
});

const isDateRange = computed<boolean>(() => isArray(props.modelValue));

function switchView(viewType: MenuView): void {
    datetimepicker.value?.switchView(viewType);
}

function getDisplayYear(year: number): string {
    return getCalendarShortYearFromUnixTime(getYearMonthDayDateTime(year, 1, 1).getUnixTime());
}

function getDisplayMonth(month: number): string {
    if (isArray(dateTime.value)) {
        return getCalendarShortMonthFromUnixTime(getYearMonthDayDateTime(dateTime.value[0].getFullYear(), month + 1, 1).getUnixTime());
    } else if (dateTime.value) {
        return getCalendarShortMonthFromUnixTime(getYearMonthDayDateTime(dateTime.value.getFullYear(), month + 1, 1).getUnixTime());
    } else {
        return getCalendarShortMonthFromUnixTime(getYearMonthDayDateTime(new Date().getFullYear(), month + 1, 1).getUnixTime());
    }
}

function getDisplayDay(date: Date): string {
    return getCalendarDayOfMonthFromUnixTime(getYearMonthDayDateTime(date.getFullYear(), date.getMonth() + 1, date.getDate()).getUnixTime());
}

defineExpose({
    switchView
});
</script>
