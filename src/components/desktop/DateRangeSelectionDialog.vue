<template>
    <v-dialog class="date-range-selection-dialog" width="640" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <h4 class="text-h4">{{ title }}</h4>
            </template>
            <template #subtitle>
                <div class="text-body-1 text-wrap mt-2">
                    <p v-if="hint">{{ hint }}</p>
                    <span v-if="beginDateTime && endDateTime">
                        <span>{{ beginDateTime }}</span>
                        <span> - </span>
                        <span>{{ endDateTime }}</span>
                    </span>
                    <slot></slot>
                </div>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
                <date-time-picker :is-dark-mode="isDarkMode"
                                  :enable-time-picker="true"
                                  :vertical="true"
                                  :preset-dates="presetRanges"
                                  :show-alternate-dates="true"
                                  v-model="dateRange">
                </date-time-picker>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
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

import { useI18n } from '@/locales/helpers.ts';
import { type CommonDateRangeSelectionProps, useDateRangeSelectionBase } from '@/components/base/DateRangeSelectionBase.ts';

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

const { tt } = useI18n();
const { dateRange, beginDateTime, endDateTime, presetRanges, getFinalDateRange } = useDateRangeSelectionBase(props);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
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
