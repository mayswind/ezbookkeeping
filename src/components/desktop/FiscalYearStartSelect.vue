<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :clearable="modelValue ? clearable : false"
        :label="label"
        :menu-props="{ contentClass: 'fiscal-year-start-select-menu' }"
        v-model="selectedDate"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayName }}</span>
        </template>

        <template #no-data>
            <vue-date-picker inline vertical auto-apply hide-offset-dates disable-year-select
                             ref="datepicker"
                             month-name-format="long"
                             model-type="MM-dd"
                             :clearable="false"
                             :enable-time-picker="false"
                             :dark="isDarkMode"
                             :week-start="firstDayOfWeek"
                             :day-names="dayNames"
                             :disabled-dates="disabledDates"
                             :start-date="selectedDate"
                             v-model="selectedDate"
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
import { computed } from 'vue';
import { useTheme } from 'vuetify';
import { useUserStore } from '@/stores/user.ts';
import { ThemeType } from '@/core/theme.ts';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import {
    type FiscalYearStartSelectionBaseProps,
    type FiscalYearStartSelectionBaseEmits,
    useFiscalYearStartSelectionBase
} from '@/components/base/FiscalYearStartSelectionBase.ts';
import { useI18n } from '@/locales/helpers.ts';

interface FiscalYearStartSelectProps extends FiscalYearStartSelectionBaseProps {
    disabled?: boolean;
    readonly?: boolean;
    clearable?: boolean;
    label?: string;
}

const props = defineProps<FiscalYearStartSelectProps>();
const emit = defineEmits<FiscalYearStartSelectionBaseEmits>();

const { getMonthShortName } = useI18n();
const theme = useTheme();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

// Get all base functionality
const {
    dayNames,
    displayName,
    disabledDates,
    firstDayOfWeek,
    getModelValueToDateString,
    getDateStringToModelValue,
} = useFiscalYearStartSelectionBase(props, emit);

const selectedDate = computed<string>({
    get: () => {
        return getModelValueToDateString();
    },
    set: (value: string) => {
        emit('update:modelValue', getDateStringToModelValue(value));
    }
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
