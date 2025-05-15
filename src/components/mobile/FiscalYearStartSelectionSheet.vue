<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="fiscal-year-start-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="tt('Clear')" @click="clear"></f7-link>
            </div>
            <div class="center">{{ selectedDisplayName(selectedDate) }}</div>
            <div class="right">
                <f7-link :text="tt('Done')" @click="confirm"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <div class="block block-outline no-margin no-padding">
                <vue-date-picker inline auto-apply hide-offset-dates disable-year-select
                                 month-name-format="long"
                                 model-type="MM-dd"
                                 six-weeks="center"
                                 class="justify-content-center"
                                 :title="selectedDisplayName(selectedDate)"
                                 :enable-time-picker="false"
                                 :clearable="true"
                                 :dark="isDarkMode"
                                 :week-start="firstDayOfWeek"
                                 :day-names="dayNames"
                                 :disabled-dates="disabledDates"
                                 v-model="selectedDate">
                    <template #month="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                    <template #month-overlay-value="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                </vue-date-picker>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';

import {
    type FiscalYearStartSelectionBaseProps,
    type FiscalYearStartSelectionBaseEmits,
    useFiscalYearStartSelectionBase
} from '@/components/base/FiscalYearStartSelectionBase.ts';

interface FiscalYearStartSelectionSheetProps extends FiscalYearStartSelectionBaseProps {
    clearable?: boolean;
    disabled?: boolean;
    label?: string;
    readonly?: boolean;
    show: boolean;
    title?: string;
}

const props = defineProps<FiscalYearStartSelectionSheetProps>();

interface FiscalYearStartSelectionSheetEmits extends FiscalYearStartSelectionBaseEmits {
    (e: 'update:show', value: boolean): void;
    (e: 'update:title', value: string): void;
}

const emit = defineEmits<FiscalYearStartSelectionSheetEmits>();

const { tt, getMonthShortName, getCurrentFiscalYearStartFormatted } = useI18n();

const environmentsStore = useEnvironmentsStore();

const {
    dayNames,
    disabledDates,
    firstDayOfWeek,
    getModelValueToDateString,
    setModelValueFromDateString,
    selectedDisplayName,
} = useFiscalYearStartSelectionBase(props, emit);

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);

const selectedDate = ref<string>(getModelValueToDateString());

function onSheetOpen(): void {
    selectedDate.value = getModelValueToDateString();
}

function clear(): void {
    selectedDate.value = getCurrentFiscalYearStartFormatted();
}

function confirm(): void {
    emit('update:modelValue', setModelValueFromDateString(selectedDate.value));
    emit('update:show', false);
    emit('update:title', selectedDisplayName(selectedDate.value));
}

</script>

<style>
.fiscal-year-start-selection-sheet .dp__menu {
    border: 0;
}
</style>
