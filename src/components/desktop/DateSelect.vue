<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :clearable="modelValue ? clearable : false"
        :label="label"
        :menu-props="{ 'content-class': 'date-select-menu' }"
        v-model="dateTime"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayTime }}</span>
        </template>

        <template #no-data>
            <vue-date-picker inline vertical auto-apply
                             ref="datepicker"
                             month-name-format="long"
                             model-type="yyyy-MM-dd"
                             :clearable="true"
                             :enable-time-picker="false"
                             :dark="isDarkMode"
                             :week-start="firstDayOfWeek"
                             :year-range="yearRange"
                             :day-names="dayNames"
                             :year-first="isYearFirst"
                             v-model="dateTime">
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
import { ref, computed } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { ThemeType } from '@/core/theme.ts';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import { getCurrentYear } from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue?: string;
    disabled?: boolean;
    readonly?: boolean;
    clearable?: boolean;
    label?: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const theme = useTheme();
const { tt, getAllMinWeekdayNames, getMonthShortName, formatDateToLongDate, isLongDateMonthAfterYear } = useI18n();

const userStore = useUserStore();

const yearRange = ref<number[]>([
    2000,
    getCurrentYear() + 1
]);

const dateTime = computed<string>({
    get: () => props.modelValue ?? '',
    set: (value: string) => emit('update:modelValue', value)
});

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const firstDayOfWeek = computed<number>(() => userStore.currentUserFirstDayOfWeek);
const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));
const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());
const displayTime = computed<string>(() => {
    if (props.modelValue) {
        return formatDateToLongDate(props.modelValue);
    } else {
        return tt('Unspecified');
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
