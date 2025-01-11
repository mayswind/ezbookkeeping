<template>
    <v-dialog class="date-range-selection-dialog" width="640" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ title }}</h4>
                </div>
            </template>
            <template #subtitle>
                <div class="text-body-1 text-center text-wrap mt-6">
                    <p v-if="hint">{{ hint }}</p>
                    <span v-if="beginDateTime && endDateTime">
                        <span>{{ beginDateTime }}</span>
                        <span> - </span>
                        <span>{{ endDateTime }}</span>
                    </span>
                    <slot></slot>
                </div>
            </template>
            <v-card-text class="mb-md-4 w-100 d-flex justify-center">
                <vue-date-picker inline enable-seconds auto-apply
                                 month-name-format="long"
                                 six-weeks="center"
                                 :clearable="false"
                                 :dark="isDarkMode"
                                 :week-start="firstDayOfWeek"
                                 :year-range="yearRange"
                                 :day-names="dayNames"
                                 :year-first="isYearFirst"
                                 :is24="is24Hour"
                                 :range="{ partialRange: false }"
                                 :preset-dates="presetRanges"
                                 v-model="dateRange">
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
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="!dateRange[0] || !dateRange[1]" @click="confirm">{{ tt('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { useTheme } from 'vuetify';

import { type CommonDateRangeSelectionProps, useDateRangeSelectionBase } from '@/components/base/DateRangeSelectionBase.ts';

import { useI18n } from '@/locales/helpers.ts';
import { useUserStore } from '@/stores/user.ts';

import { ThemeType } from '@/core/theme.ts';

import {
    getLocalDatetimeFromUnixTime,
    getDummyUnixTimeForLocalUsage,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes
} from '@/lib/datetime.ts';

interface DesktopDateRangeSelectionProps extends CommonDateRangeSelectionProps {
    persistent?: boolean;
}

const props = defineProps<DesktopDateRangeSelectionProps>();
const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'dateRange:change', minUnixTime: number, maxUnixTime: number): void;
    (e: 'error', message: string): void;
}>();

const theme = useTheme();
const { tt, getMonthShortName } = useI18n();

const userStore = useUserStore();

const { yearRange, dateRange, dayNames, isYearFirst, is24Hour, beginDateTime, endDateTime, presetRanges, getFinalDateRange } = useDateRangeSelectionBase(props);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);
const showState = computed<boolean>({
    get: () => props.show || false,
    set: (value) => emit('update:show', value)
});

function confirm(): void {
    try {
        const finalDateRange = getFinalDateRange();

        if (!finalDateRange) {
            return;
        }

        emit('dateRange:change', finalDateRange.minUnixTime, finalDateRange.maxUnixTime);
    } catch (ex: unknown) {
        if (ex instanceof Error) {
            emit('error', ex.message);
        }
    }
}

function cancel(): void {
    emit('update:show', false);
}

watch(() => props.minTime, (newValue) => {
    if (newValue) {
        dateRange.value[0] = getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(newValue, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
    }
});

watch(() => props.maxTime, (newValue) => {
    if (newValue) {
        dateRange.value[1] = getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(newValue, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
    }
});
</script>

<style>
.date-range-selection-dialog .dp__preset_ranges {
    white-space: nowrap !important;
}
</style>
