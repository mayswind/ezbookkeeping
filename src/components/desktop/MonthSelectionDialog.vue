<template>
    <v-dialog class="month-selection-dialog" width="640" :persistent="!!persistent" v-model="showState">
        <one-column-dialog-layout :title="title" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       :disabled="!monthValue" @click="confirm">{{ tt('OK') }}</v-btn>
            </template>

            <template #content>
                <div class="text-body-1 text-wrap mt-2" v-if="hint">
                    <span>{{ hint }}</span>
                    <slot></slot>
                </div>
                <v-row class="mt-2">
                    <v-col>
                        <month-picker :is-dark-mode="isDarkMode" v-model="monthValue"></month-picker>
                    </v-col>
                </v-row>
            </template>
        </one-column-dialog-layout>
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
.month-selection-dialog .dp--main .dp--instance-calendar .dp--overlay.dp--overlay-relative {
    width: 100% !important;
}
</style>
