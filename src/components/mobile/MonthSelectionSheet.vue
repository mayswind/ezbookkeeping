<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="month-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top" v-if="hint">{{ hint }}</p>
                <slot></slot>
                <month-picker month-picker-class="justify-content-center margin-bottom"
                              :is-dark-mode="isDarkMode" v-model="monthValue"></month-picker>
                <f7-button large fill
                           :class="{ 'disabled': !monthValue }"
                           :text="tt('Continue')"
                           @click="confirm">
                </f7-button>
                <div class="margin-top text-align-center">
                    <f7-link @click="cancel" :text="tt('Cancel')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';

import type { Year0BasedMonth } from '@/core/datetime.ts';
import { useEnvironmentsStore } from '@/stores/environment.ts';

import { getYear0BasedMonthObjectFromUnixTime, getThisMonthFirstUnixTime } from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue?: Year0BasedMonth;
    title?: string;
    hint?: string;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Year0BasedMonth): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const environmentsStore = useEnvironmentsStore();

const monthValue = ref<Year0BasedMonth>(getYear0BasedMonthObjectFromUnixTime(getThisMonthFirstUnixTime()));
const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);

function confirm(): void {
    if (monthValue.value.year <= 0 || monthValue.value.month0base < 0) {
        showToast('Date is too early');
        return;
    }

    emit('update:modelValue', monthValue.value);
}

function cancel(): void {
    emit('update:show', false);
}

function onSheetOpen(): void {
    if (props.modelValue) {
        monthValue.value = props.modelValue;
    }
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>

<style>
.month-selection-sheet .dp__main .dp__instance_calendar .dp__overlay.dp--overlay-relative {
    width: 100% !important;
}

.month-selection-sheet .dp__main .dp__instance_calendar .dp__overlay.dp--overlay-relative .dp__selection_grid_header .dp--year-mode-picker .dp--arrow-btn-nav {
    display: flex;
}

.month-selection-sheet .dp__main .dp__instance_calendar .dp__overlay.dp--overlay-relative .dp__selection_grid_header .dp--year-mode-picker .dp--year-select+.dp--arrow-btn-nav {
    justify-content: end;
}
</style>
