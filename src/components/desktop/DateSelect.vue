<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :clearable="modelValue ? clearable : false"
        :label="label"
        :menu-props="{ contentClass: 'date-select-menu' }"
        v-model="displayTime"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayTime }}</span>
        </template>

        <template #no-data>
            <date-time-picker :is-dark-mode="isDarkMode"
                              :enable-time-picker="false"
                              :clearable="true"
                              :show-alternate-dates="true"
                              v-model="dateTime">
            </date-time-picker>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { type TextualYearMonthDay } from '@/core/datetime.ts';
import { ThemeType } from '@/core/theme.ts';

import {
    getLocalDateFromYearDashMonthDashDay,
    getGregorianCalendarYearAndMonthFromLocalDate
} from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue?: TextualYearMonthDay;
    disabled?: boolean;
    readonly?: boolean;
    clearable?: boolean;
    label?: string;
    noDataText?: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: TextualYearMonthDay | ''): void;
}>();

const theme = useTheme();
const { tt, formatGregorianCalendarYearDashMonthDashDayToLongDate } = useI18n();

const dateTime = computed<Date | null>({
    get: () => props.modelValue ? getLocalDateFromYearDashMonthDashDay(props.modelValue) : null,
    set: (value: Date | null) => emit('update:modelValue', value ? getGregorianCalendarYearAndMonthFromLocalDate(value) : '')
});

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const displayTime = computed<string>({
    get: () => {
        if (props.modelValue) {
            return formatGregorianCalendarYearDashMonthDashDayToLongDate(props.modelValue);
        } else if (props.noDataText) {
            return props.noDataText;
        } else {
            return tt('Unspecified');
        }
    },
    set: (value: string) => {
        if (!value) {
            dateTime.value = null;
        }
    }
});
</script>

<style>
.date-select-menu {
    max-height: inherit !important;
}

.date-select-menu .dp__menu {
    border: 0;
}
</style>
