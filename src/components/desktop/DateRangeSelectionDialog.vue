<template>
    <v-dialog class="date-range-selection-dialog" width="640" :persistent="!!persistent" v-model="showState">
        <one-column-dialog-layout :title="title" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       :disabled="!dateRange[0] || !dateRange[1]" @click="confirm">{{ tt('OK') }}</v-btn>
            </template>

            <template #content>
                <div class="text-body-1 mt-3" v-if="beginDateTime && endDateTime">
                    <span>{{ beginDateTime }}</span>
                    <span> - </span>
                    <span>{{ endDateTime }}</span>
                </div>

                <div class="text-body-1 text-wrap mt-2" v-if="hint">
                    <span>{{ hint }}</span>
                    <slot></slot>
                </div>

                <v-row class="mt-1">
                    <v-col>
                        <date-time-picker :is-dark-mode="isDarkMode"
                                          :enable-time-picker="true"
                                          :vertical="true"
                                          :preset-dates="presetRanges"
                                          :show-alternate-dates="true"
                                          v-model="dateRange">
                        </date-time-picker>
                    </v-col>
                </v-row>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonDateRangeSelectionProps, useDateRangeSelectionBase } from '@/components/base/DateRangeSelectionBase.ts';

import { ThemeType } from '@/core/theme.ts';

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
const {
    dateRange,
    beginDateTime,
    endDateTime,
    presetRanges,
    getLocalDatetimeFromSameDateTimeOfUnixTime,
    getFinalDateRange
} = useDateRangeSelectionBase(props);

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
        dateRange.value[0] = getLocalDatetimeFromSameDateTimeOfUnixTime(newValue);
    }
});

watch(() => props.maxTime, (newValue) => {
    if (newValue) {
        dateRange.value[1] = getLocalDatetimeFromSameDateTimeOfUnixTime(newValue);
    }
});
</script>

<style>
.date-range-selection-dialog .dp--main .dp--instance-calendar .dp--preset-ranges {
    white-space: nowrap !important;
}
</style>
