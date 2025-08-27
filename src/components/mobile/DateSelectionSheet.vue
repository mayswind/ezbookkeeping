<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="date-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="tt('Clear')" @click="clear"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="tt('Done')" @click="confirm"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <div class="block block-outline no-margin no-padding">
                <date-time-picker datetime-picker-class="justify-content-center"
                                  :is-dark-mode="isDarkMode"
                                  :enable-time-picker="false"
                                  :clearable="true"
                                  v-model="dateTime">
                </date-time-picker>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';

import { type TextualYearMonthDay } from '@/core/datetime.ts';

import {
    getLocalDateFromYearDashMonthDashDay,
    getGregorianCalendarYearAndMonthFromLocalDate
} from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue?: TextualYearMonthDay;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: TextualYearMonthDay | ''): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const environmentsStore = useEnvironmentsStore();

const dateTime = ref<Date | null>(null);

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);

function clear(): void {
    dateTime.value = null;
    confirm();
}

function confirm(): void {
    emit('update:modelValue', dateTime.value ? getGregorianCalendarYearAndMonthFromLocalDate(dateTime.value) : '');
    emit('update:show', false);
}

function onSheetOpen(): void {
    dateTime.value = props.modelValue ? getLocalDateFromYearDashMonthDashDay(props.modelValue) : null;
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>

<style>
.date-selection-sheet .dp__menu {
    border: 0;
}
</style>
