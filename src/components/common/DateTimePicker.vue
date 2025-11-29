<template>
    <vue-date-picker ref="datetimepicker"
                     inline auto-apply
                     six-weeks="center"
                     :class="`datetime-picker ${showAlternateDates && alternateCalendarType ? 'datetime-picker-with-alternate-date' : ''} ${datetimePickerClass}`"
                     :config="{ noSwipe: !!noSwipeAndScroll, monthChangeOnScroll: !noSwipeAndScroll }"
                     :time-config="{ enableTimePicker: enableTimePicker, enableSeconds: true, is24: is24Hour }"
                     :input-attrs="{ clearable: !!clearable }"
                     :dark="isDarkMode"
                     :vertical="vertical"
                     :disable-year-select="disableYearSelect"
                     :year-range="yearRange"
                     :day-names="dayNames"
                     :week-start="firstDayOfWeek"
                     :year-first="isYearFirst"
                     :min-date="minDate"
                     :max-date="maxDate"
                     :disabled-dates="disabledDates"
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
            <div class="datetime-picker-display-dates">
                <span>{{ getDisplayDay(date) }}</span>
                <span class="datetime-picker-alternate-date" v-if="showAlternateDates && alternateCalendarType && getAlternateDate(date)">{{ getAlternateDate(date) }}</span>
            </div>
        </template>
        <template #am-pm-button="{ toggle, value }">
            <button class="dp__pm_am_button" tabindex="0" @click="toggle">{{ tt(`datetime.${value}.content`) }}</button>
        </template>
    </vue-date-picker>
</template>

<script setup lang="ts">
import { computed, useTemplateRef } from 'vue';
import { type MenuView, VueDatePicker } from '@vuepic/vue-datepicker';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import type { CalendarType } from '@/core/calendar.ts';
import { NumeralSystem } from '@/core/numeral.ts';
import type { PresetDateRange, WeekDayValue } from '@/core/datetime.ts';
import { isDefined, isArray, arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import { getAllowedYearRange, getYearMonthDayDateTime } from '@/lib/datetime.ts';

type VueDatePickerType = InstanceType<typeof VueDatePicker>;
type SupportedModelValue = Date | Date[] | null;

const props = defineProps<{
    modelValue: SupportedModelValue;
    datetimePickerClass?: string;
    isDarkMode: boolean;
    numeralSystem?: number;
    enableTimePicker: boolean;
    disableYearSelect?: boolean;
    vertical?: boolean;
    noSwipeAndScroll?: boolean;
    clearable?: boolean;
    minDate?: Date;
    maxDate?: Date;
    disabledDates?: (date: Date) => boolean;
    showAlternateDates?: boolean;
    presetRanges?: PresetDateRange[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: SupportedModelValue): void;
}>();

const {
    tt,
    getAllMinWeekdayNames,
    getCurrentCalendarDisplayType,
    getCurrentNumeralSystemType,
    isLongDateMonthAfterYear,
    isLongTime24HourFormat,
    getCalendarDisplayShortYearFromUnixTime,
    getCalendarDisplayShortMonthFromUnixTime,
    getCalendarDisplayDayOfMonthFromUnixTime,
    getCalendarAlternateDate
} = useI18n();

const userStore = useUserStore();

const datetimepicker = useTemplateRef<VueDatePickerType>('datetimepicker');

const yearRange = getAllowedYearRange();

const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));
const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());
const is24Hour = computed<boolean>(() => isLongTime24HourFormat());
const alternateCalendarType = computed<CalendarType | undefined>(() => getCurrentCalendarDisplayType().secondaryCalendarType);

const actualNumeralSystem = computed<NumeralSystem>(() => {
    if (isDefined(props.numeralSystem)) {
        return NumeralSystem.valueOf(props.numeralSystem) ?? NumeralSystem.Default;
    } else {
        return getCurrentNumeralSystemType();
    }
});

const dateTime = computed<SupportedModelValue>({
    get: () => props.modelValue,
    set: (value: SupportedModelValue) => emit('update:modelValue', value)
});

const isDateRange = computed<boolean>(() => isArray(props.modelValue));

function getAlternateDate(date: Date): string | undefined {
    if (!props.showAlternateDates) {
        return undefined;
    }

    return getCalendarAlternateDate({
        year: date.getFullYear(),
        month: date.getMonth() + 1,
        day: date.getDate()
    })?.displayDate;
};

function switchView(viewType: MenuView): void {
    datetimepicker.value?.switchView(viewType);
}

function getDisplayYear(year: number): string {
    return getCalendarDisplayShortYearFromUnixTime(getYearMonthDayDateTime(year, 1, 1).getUnixTime(), actualNumeralSystem.value);
}

function getDisplayMonth(month: number): string {
    if (isArray(dateTime.value)) {
        return getCalendarDisplayShortMonthFromUnixTime(getYearMonthDayDateTime(dateTime.value[0]!.getFullYear(), month + 1, 1).getUnixTime(), actualNumeralSystem.value);
    } else if (dateTime.value) {
        return getCalendarDisplayShortMonthFromUnixTime(getYearMonthDayDateTime(dateTime.value.getFullYear(), month + 1, 1).getUnixTime(), actualNumeralSystem.value);
    } else {
        return getCalendarDisplayShortMonthFromUnixTime(getYearMonthDayDateTime(new Date().getFullYear(), month + 1, 1).getUnixTime(), actualNumeralSystem.value);
    }
}

function getDisplayDay(date: Date): string {
    return getCalendarDisplayDayOfMonthFromUnixTime(getYearMonthDayDateTime(date.getFullYear(), date.getMonth() + 1, date.getDate()).getUnixTime(), actualNumeralSystem.value);
}

defineExpose({
    switchView
});
</script>

<style>
.datetime-picker-display-dates {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.datetime-picker-alternate-date {
    margin-top: -2px;
    opacity: 0.5;
    font-size: 0.8rem;
}

.dp__cell_disabled .datetime-picker-alternate-date,
.dp__cell_offset .datetime-picker-alternate-date {
    opacity: 0.8;
}

.dp__main.datetime-picker .dp__calendar .dp__calendar_row > .dp__calendar_item .datetime-picker-display-dates > span.datetime-picker-alternate-date {
    display: block;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.dp__main.datetime-picker.datetime-picker-with-alternate-date .dp__calendar .dp__calendar_row {
    --dp-cell-size: 45px;
}
</style>
