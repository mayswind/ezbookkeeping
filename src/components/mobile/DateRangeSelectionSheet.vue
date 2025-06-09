<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top" v-if="hint">{{ hint }}</p>
                <p class="no-margin-top margin-bottom" v-if="beginDateTime && endDateTime">
                    <span>{{ beginDateTime }}</span>
                    <span> - </span>
                    <span>{{ endDateTime }}</span>
                </p>
                <slot></slot>
                <vue-date-picker inline enable-seconds auto-apply
                                 ref="datetimepicker"
                                 month-name-format="long"
                                 six-weeks="center"
                                 class="justify-content-center margin-bottom"
                                 :clearable="false"
                                 :dark="isDarkMode"
                                 :week-start="firstDayOfWeek"
                                 :year-range="yearRange"
                                 :day-names="dayNames"
                                 :year-first="isYearFirst"
                                 :is24="is24Hour"
                                 :range="{ partialRange: false }"
                                 :preset-dates="presetRanges"
                                 v-model="dateRange">
                    <template #month="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                    <template #month-overlay-value="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                    <template #am-pm-button="{ toggle, value }">
                        <button class="dp__pm_am_button" tabindex="0" @click="toggle">{{ tt(`datetime.${value}.content`) }}</button>
                    </template>
                </vue-date-picker>
                <f7-button large fill
                           :class="{ 'disabled': !dateRange[0] || !dateRange[1] }"
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
import { computed, useTemplateRef } from 'vue';
import VueDatePicker from '@vuepic/vue-datepicker';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { type CommonDateRangeSelectionProps, useDateRangeSelectionBase } from '@/components/base/DateRangeSelectionBase.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';
import { useUserStore } from '@/stores/user.ts';

import { type WeekDayValue } from '@/core/datetime.ts';

import {
    getLocalDatetimeFromUnixTime,
    getDummyUnixTimeForLocalUsage,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes
} from '@/lib/datetime.ts';

type VueDatePickerType = InstanceType<typeof VueDatePicker>;

const props = defineProps<CommonDateRangeSelectionProps>();
const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'dateRange:change', minUnixTime: number, maxUnixTime: number): void;
}>();

const { tt, getMonthShortName } = useI18n();
const { showToast } = useI18nUIComponents();
const { yearRange, dateRange, dayNames, isYearFirst, is24Hour, beginDateTime, endDateTime, presetRanges, getFinalDateRange } = useDateRangeSelectionBase(props);

const environmentsStore = useEnvironmentsStore();
const userStore = useUserStore();

const datetimepicker = useTemplateRef<VueDatePickerType>('datetimepicker');
const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);
const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);

function confirm(): void {
    try {
        const finalDateRange = getFinalDateRange();

        if (!finalDateRange) {
            return;
        }

        emit('dateRange:change', finalDateRange.minUnixTime, finalDateRange.maxUnixTime);
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
    if (props.minTime) {
        dateRange.value[0] = getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(props.minTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
    }

    if (props.maxTime) {
        dateRange.value[1] = getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(props.maxTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
    }

    window.dispatchEvent(new Event('resize')); // fix vue-datepicker preset max-width

    if (datetimepicker.value) {
        datetimepicker.value.switchView('calendar');
    }
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>
