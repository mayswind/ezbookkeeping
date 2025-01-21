<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ 'content-class': 'date-time-select-menu' }"
        v-model="dateTime"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayTime }}</span>
        </template>

        <template #no-data>
            <vue-date-picker inline vertical time-picker-inline enable-seconds auto-apply
                             ref="datepicker"
                             month-name-format="long"
                             :clearable="false"
                             :dark="isDarkMode"
                             :week-start="firstDayOfWeek"
                             :year-range="yearRange"
                             :day-names="dayNames"
                             :year-first="isYearFirst"
                             :is24="is24Hour"
                             v-model="dateTime">
                <template #month="{ text }">
                    {{ getMonthShortName(text) }}
                </template>
                <template #month-overlay-value="{ text }">
                    {{ getMonthShortName(text) }}
                </template>
                <template #am-pm-button="{ toggle, value }">
                    <button class="dp__pm_am_button" tabindex="0" @click="toggle">{{ tt(`datetime.${value}.content`) }}</button>
                </template>
            </vue-date-picker>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { ThemeType } from '@/core/theme.ts';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import {
    getCurrentYear,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getLocalDatetimeFromUnixTime,
    getActualUnixTimeForStore,
    getUnixTime
} from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue: number;
    disabled?: boolean;
    readonly?: boolean;
    label?: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
    (e: 'error', message: string): void;
}>();

const theme = useTheme();
const { tt, getAllMinWeekdayNames, getMonthShortName, formatUnixTimeToLongDateTime, isLongDateMonthAfterYear, isLongTime24HourFormat } = useI18n();

const userStore = useUserStore();

const yearRange = ref<number[]>([
    2000,
    getCurrentYear() + 1
]);

const dateTime = computed<Date>({
    get: () => {
        return getLocalDatetimeFromUnixTime(props.modelValue);
    },
    set: (value: Date) => {
        const unixTime = getUnixTime(value);

        if (unixTime < 0) {
            emit('error', 'Date is too early');
            return;
        }

        emit('update:modelValue', unixTime);
    }
});

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);
const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));
const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());
const is24Hour = computed<boolean>(() => isLongTime24HourFormat());
const displayTime = computed<string>(() => formatUnixTimeToLongDateTime(getActualUnixTimeForStore(getUnixTime(dateTime.value), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())));
</script>

<style>
.date-time-select-menu {
    max-height: inherit !important;
}

.date-time-select-menu .dp__menu {
    border: 0;
}
</style>
