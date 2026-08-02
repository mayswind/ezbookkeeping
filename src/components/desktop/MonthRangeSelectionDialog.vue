<template>
    <v-dialog class="month-range-selection-dialog" width="640" :persistent="!!persistent" v-model="showState">
        <one-column-dialog-layout :title="title" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="me-2" density="comfortable" variant="outlined"
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
                    <v-col cols="12" md="6">
                        <month-picker :is-dark-mode="isDarkMode" v-model="dateRange[0]"></month-picker>
                    </v-col>
                    <v-col cols="12" md="6">
                        <month-picker :is-dark-mode="isDarkMode" v-model="dateRange[1]"></month-picker>
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
import { type CommonMonthRangeSelectionProps, useMonthRangeSelectionBase } from '@/components/base/MonthRangeSelectionBase.ts';

import { ThemeType } from '@/core/theme.ts';
import { type TextualYearMonth } from '@/core/datetime.ts';

import { getYear0BasedMonthObjectFromString } from '@/lib/datetime.ts';

interface DesktopMonthRangeSelectionProps extends CommonMonthRangeSelectionProps {
    persistent?: boolean;
}

const props = defineProps<DesktopMonthRangeSelectionProps>();
const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'dateRange:change', minYearMonth: TextualYearMonth | '', maxYearMonth: TextualYearMonth | ''): void;
    (e: 'error', message: string): void;
}>();

const theme = useTheme();

const { tt } = useI18n();
const { dateRange, beginDateTime, endDateTime, getFinalMonthRange } = useMonthRangeSelectionBase(props);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const showState = computed<boolean>({
    get: () => props.show || false,
    set: (value) => emit('update:show', value)
});

function confirm(): void {
    try {
        const finalMonthRange = getFinalMonthRange();

        if (!finalMonthRange) {
            return;
        }

        emit('dateRange:change', finalMonthRange.minYearMonth, finalMonthRange.maxYearMonth);
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
        const yearMonth = getYear0BasedMonthObjectFromString(newValue);

        if (yearMonth) {
            dateRange.value[0] = yearMonth;
        }
    }
});

watch(() => props.maxTime, (newValue) => {
    if (newValue) {
        const yearMonth = getYear0BasedMonthObjectFromString(newValue);

        if (yearMonth) {
            dateRange.value[1] = yearMonth;
        }
    }
});
</script>

<style>
.month-range-selection-dialog .dp--main .dp--instance-calendar .dp--overlay.dp--overlay-relative {
    width: 100% !important;
}
</style>
