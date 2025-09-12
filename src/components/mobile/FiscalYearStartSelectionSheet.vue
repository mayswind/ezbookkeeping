<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="fiscal-year-start-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="tt('Reset')" @click="reset"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="tt('Done')" @click="confirm"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <div class="block block-outline no-margin no-padding">
                <date-time-picker datetime-picker-class="justify-content-center"
                                  :is-dark-mode="isDarkMode"
                                  :numeral-system="numeralSystem"
                                  :enable-time-picker="false"
                                  :disable-year-select="true"
                                  :no-swipe-and-scroll="true"
                                  :min-date="allowedMinDate"
                                  :max-date="allowedMaxDate"
                                  :disabled-dates="disabledDates"
                                  v-model="selectedFiscalYearStartValue">
                </date-time-picker>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';
import { useUserStore } from '@/stores/user.ts';

import {
    type CommonFiscalYearStartSelectionProps,
    type CommonFiscalYearStartSelectionEmits,
    useFiscalYearStartSelectionBase
} from '@/components/base/FiscalYearStartSelectionBase.ts';

import { FiscalYearStart } from '@/core/fiscalyear.ts';

interface MobileFiscalYearStartSelectionSheetProps extends CommonFiscalYearStartSelectionProps {
    show: boolean;
}

interface MobileFiscalYearStartSelectionSheetEmits extends CommonFiscalYearStartSelectionEmits {
    (e: 'update:show', value: boolean): void;
}

const props = defineProps<MobileFiscalYearStartSelectionSheetProps>();
const emit = defineEmits<MobileFiscalYearStartSelectionSheetEmits>();

const { tt } = useI18n();

const environmentsStore = useEnvironmentsStore();
const userStore = useUserStore();

const {
    disabledDates,
    selectedFiscalYearStart,
    selectedFiscalYearStartValue,
    allowedMinDate,
    allowedMaxDate,
} = useFiscalYearStartSelectionBase(props);

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);

function confirm(): void {
    emit('update:modelValue', selectedFiscalYearStart.value);
    emit('update:show', false);
}

function reset(): void {
    selectedFiscalYearStart.value = userStore.currentUserFiscalYearStart;
}

function onSheetOpen(): void {
    if (props.modelValue) {
        const fiscalYearStart = FiscalYearStart.valueOf(props.modelValue);

        if (fiscalYearStart) {
            selectedFiscalYearStart.value = fiscalYearStart.value;
        }
    }
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>

<style>
.fiscal-year-start-selection-sheet .dp__menu {
    border: 0;
}
</style>
