<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ contentClass: 'fiscal-year-start-select-menu' }"
        v-model="selectedFiscalYearStartValue"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayFiscalYearStartDate }}</span>
        </template>

        <template #no-data>
            <date-time-picker :is-dark-mode="isDarkMode"
                              :numeral-system="numeralSystem"
                              :vertical="true"
                              :enable-time-picker="false"
                              :disable-year-select="true"
                              :no-swipe-and-scroll="true"
                              :min-date="allowedMinDate"
                              :max-date="allowedMaxDate"
                              :disabled-dates="disabledDates"
                              v-model="selectedFiscalYearStartValue">
            </date-time-picker>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { useTheme } from 'vuetify';

import {
    type CommonFiscalYearStartSelectionProps,
    type CommonFiscalYearStartSelectionEmits,
    useFiscalYearStartSelectionBase
} from '@/components/base/FiscalYearStartSelectionBase.ts';

import { ThemeType } from '@/core/theme.ts';

const props = defineProps<CommonFiscalYearStartSelectionProps>();
const emit = defineEmits<CommonFiscalYearStartSelectionEmits>();

const theme = useTheme();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const {
    disabledDates,
    selectedFiscalYearStart,
    selectedFiscalYearStartValue,
    displayFiscalYearStartDate,
    allowedMinDate,
    allowedMaxDate,
} = useFiscalYearStartSelectionBase(props);

watch(() => props.modelValue, (newValue) => {
    if (newValue) {
        selectedFiscalYearStart.value = newValue;
    }
});

watch(selectedFiscalYearStart, (newValue) => {
    emit('update:modelValue', newValue);
});
</script>

<style>
.fiscal-year-start-select-menu {
    max-height: inherit !important;
}

.fiscal-year-start-select-menu .dp__menu {
    border: 0;
}
</style>
