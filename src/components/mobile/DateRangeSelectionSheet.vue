<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler"></div>
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
                <date-time-picker ref="datetimepicker"
                                  datetime-picker-class="justify-content-center margin-bottom"
                                  :is-dark-mode="isDarkMode"
                                  :enable-time-picker="true"
                                  :preset-dates="presetRanges"
                                  :show-alternate-dates="true"
                                  v-model="dateRange">
                </date-time-picker>
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
import DateTimePicker from '@/components/common/DateTimePicker.vue';
import { computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { type CommonDateRangeSelectionProps, useDateRangeSelectionBase } from '@/components/base/DateRangeSelectionBase.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';

type DateTimePickerType = InstanceType<typeof DateTimePicker>;

const props = defineProps<CommonDateRangeSelectionProps>();
const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'dateRange:change', minUnixTime: number, maxUnixTime: number): void;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();
const {
    dateRange,
    beginDateTime,
    endDateTime,
    presetRanges,
    getLocalDatetimeFromSameDateTimeOfUnixTime,
    getFinalDateRange
} = useDateRangeSelectionBase(props);

const environmentsStore = useEnvironmentsStore();

const datetimepicker = useTemplateRef<DateTimePickerType>('datetimepicker');
const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);

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
        dateRange.value[0] = getLocalDatetimeFromSameDateTimeOfUnixTime(props.minTime);
    }

    if (props.maxTime) {
        dateRange.value[1] = getLocalDatetimeFromSameDateTimeOfUnixTime(props.maxTime);
    }

    window.dispatchEvent(new Event('resize')); // fix vue-datepicker preset max-width
    datetimepicker.value?.switchView('calendar');
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>
