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
                <vue-date-picker inline month-picker auto-apply
                                 month-name-format="long"
                                 class="justify-content-center margin-bottom"
                                 :clearable="false"
                                 :dark="isDarkMode"
                                 :year-range="yearRange"
                                 :year-first="isYearFirst"
                                 v-model="monthValue">
                    <template #month="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                    <template #month-overlay-value="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                </vue-date-picker>
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
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { type CommonMonthSelectionProps, useMonthSelectionBase } from '@/components/base/MonthSelectionBase.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';

import { getYearMonthObjectFromString } from '@/lib/datetime.ts';

const props = defineProps<CommonMonthSelectionProps>();
const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, getMonthShortName } = useI18n();
const { showToast } = useI18nUIComponents();
const { yearRange, monthValue, isYearFirst, getTextualYearMonth } = useMonthSelectionBase(props);

const environmentsStore = useEnvironmentsStore();

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);

function confirm(): void {
    try {
        const finalMonthRange = getTextualYearMonth();

        if (!finalMonthRange) {
            return;
        }

        emit('update:modelValue', finalMonthRange);
    } catch (ex: unknown) {
        if (ex instanceof Error) {
            showToast(ex.message);
        }
    }
}

function cancel(): void {
    emit('update:show', false);
}

function onSheetOpen(): void {
    if (props.modelValue) {
        const yearMonth = getYearMonthObjectFromString(props.modelValue);

        if (yearMonth) {
            monthValue.value = yearMonth;
        }
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
