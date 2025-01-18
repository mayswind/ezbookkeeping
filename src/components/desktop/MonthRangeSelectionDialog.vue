<template>
    <v-dialog class="month-range-selection-dialog" width="640" :persistent="!!persistent" v-model="showState">
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
                <v-row class="match-height">
                    <v-col cols="12" md="6">
                        <vue-date-picker inline month-picker auto-apply
                                         month-name-format="long"
                                         :clearable="false"
                                         :dark="isDarkMode"
                                         :year-range="yearRange"
                                         :year-first="isYearFirst"
                                         v-model="dateRange[0]">
                            <template #month="{ text }">
                                {{ getMonthShortName(text) }}
                            </template>
                            <template #month-overlay-value="{ text }">
                                {{ getMonthShortName(text) }}
                            </template>
                        </vue-date-picker>
                    </v-col>
                    <v-col cols="12" md="6">
                        <vue-date-picker inline month-picker auto-apply
                                         month-name-format="long"
                                         :clearable="false"
                                         :dark="isDarkMode"
                                         :year-range="yearRange"
                                         :year-first="isYearFirst"
                                         v-model="dateRange[1]">
                            <template #month="{ text }">
                                {{ getMonthShortName(text) }}
                            </template>
                            <template #month-overlay-value="{ text }">
                                {{ getMonthShortName(text) }}
                            </template>
                        </vue-date-picker>
                    </v-col>
                </v-row>
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

import { useI18n } from '@/locales/helpers.ts';
import { type CommonMonthRangeSelectionProps, useMonthRangeSelectionBase } from '@/components/base/MonthRangeSelectionBase.ts';

import { ThemeType } from '@/core/theme.ts';
import { getYearMonthObjectFromString } from '@/lib/datetime.ts';

interface DesktopMonthRangeSelectionProps extends CommonMonthRangeSelectionProps {
    persistent?: boolean;
}

const props = defineProps<DesktopMonthRangeSelectionProps>();
const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'dateRange:change', minYearMonth: string, maxYearMonth: string): void;
    (e: 'error', message: string): void;
}>();

const theme = useTheme();

const { tt, getMonthShortName } = useI18n();
const { yearRange, dateRange, isYearFirst, beginDateTime, endDateTime, getFinalMonthRange } = useMonthRangeSelectionBase(props);

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
        const yearMonth = getYearMonthObjectFromString(newValue);

        if (yearMonth) {
            dateRange.value[0] = yearMonth;
        }
    }
});

watch(() => props.maxTime, (newValue) => {
    if (newValue) {
        const yearMonth = getYearMonthObjectFromString(newValue);

        if (yearMonth) {
            dateRange.value[1] = yearMonth;
        }
    }
});
</script>

<style>
.month-range-selection-dialog .dp__preset_ranges {
    white-space: nowrap !important;
}

.month-range-selection-dialog .dp__overlay {
    width: 100% !important;
    height: 100% !important;
}
</style>
