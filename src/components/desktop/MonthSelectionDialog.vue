<template>
    <v-dialog class="month-selection-dialog" width="640" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ title }}</h4>
                </div>
            </template>
            <template #subtitle>
                <div class="text-body-1 text-center text-wrap mt-6">
                    <p v-if="hint">{{ hint }}</p>
                    <slot></slot>
                </div>
            </template>
            <v-card-text class="mb-md-4 w-100 d-flex justify-center">
                <v-row class="match-height">
                    <v-col>
                        <month-picker :is-dark-mode="isDarkMode" v-model="monthValue"></month-picker>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="!monthValue" @click="confirm">{{ tt('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { ThemeType } from '@/core/theme.ts';
import type { Year0BasedMonth } from '@/core/datetime.ts';

import { getYear0BasedMonthObjectFromUnixTime, getThisMonthFirstUnixTime } from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue?: Year0BasedMonth;
    title?: string;
    hint?: string;
    show: boolean;
    persistent?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Year0BasedMonth): void;
    (e: 'update:show', value: boolean): void;
    (e: 'error', message: string): void;
}>();

const theme = useTheme();

const { tt } = useI18n();

const monthValue = ref<Year0BasedMonth>(getYear0BasedMonthObjectFromUnixTime(getThisMonthFirstUnixTime()));

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const showState = computed<boolean>({
    get: () => props.show || false,
    set: (value) => emit('update:show', value)
});

function confirm(): void {
    if (monthValue.value.year <= 0 || monthValue.value.month0base < 0) {
        emit('error', 'Date is too early');
        return;
    }

    emit('update:modelValue', monthValue.value);
}

function cancel(): void {
    emit('update:show', false);
}

watch(() => props.modelValue, (newValue) => {
    if (newValue) {
        monthValue.value = newValue;
    }
});

watch(() => props.show, (newValue) => {
    if (newValue && props.modelValue) {
        monthValue.value = props.modelValue;
    }
});
</script>

<style>
.month-selection-dialog .dp__preset_ranges {
    white-space: nowrap !important;
}

.month-selection-dialog .dp__overlay {
    width: 100% !important;
    height: 100% !important;
}
</style>
