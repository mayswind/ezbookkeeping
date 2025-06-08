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
            <vue-date-picker inline auto-apply disable-year-select
                             month-name-format="long"
                             model-type="MM-dd"
                             six-weeks="center"
                             :clearable="false"
                             :enable-time-picker="false"
                             :dark="isDarkMode"
                             :week-start="firstDayOfWeek"
                             :day-names="dayNames"
                             :disabled-dates="disabledDates"
                             v-model="selectedFiscalYearStartValue"
                             >
                <template #month="{ text }">
                    {{ getMonthShortName(text) }}
                </template>
                <template #month-overlay-value="{ text }">
                    {{ getMonthShortName(text) }}
                </template>
            </vue-date-picker>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import {
    type CommonFiscalYearStartSelectionProps,
    type CommonFiscalYearStartSelectionEmits,
    useFiscalYearStartSelectionBase
} from '@/components/base/FiscalYearStartSelectionBase.ts';

import { ThemeType } from '@/core/theme.ts';

const props = defineProps<CommonFiscalYearStartSelectionProps>();
const emit = defineEmits<CommonFiscalYearStartSelectionEmits>();

const { getMonthShortName } = useI18n();

const theme = useTheme();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const {
    disabledDates,
    selectedFiscalYearStart,
    selectedFiscalYearStartValue,
    displayFiscalYearStartDate,
    firstDayOfWeek,
    dayNames
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
