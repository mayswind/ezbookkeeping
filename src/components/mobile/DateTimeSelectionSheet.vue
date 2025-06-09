<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="date-time-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="tt('Now')" @click="setCurrentTime"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="tt('Done')" @click="confirm"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content class="padding-bottom">
            <div class="block block-outline no-margin no-padding">
                <vue-date-picker inline enable-seconds auto-apply
                                 ref="datetimepicker"
                                 month-name-format="long"
                                 six-weeks="center"
                                 class="justify-content-center"
                                 :enable-time-picker="false"
                                 :clearable="false"
                                 :dark="isDarkMode"
                                 :week-start="firstDayOfWeek"
                                 :year-range="yearRange"
                                 :day-names="dayNames"
                                 :year-first="isYearFirst"
                                 v-model="dateTime"
                                 v-show="mode === 'date'">
                    <template #month="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                    <template #month-overlay-value="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                </vue-date-picker>
            </div>
            <div class="block block-outline no-margin no-padding padding-vertical-half" v-show="mode === 'time'">
                <div id="time-picker-container" class="time-picker-container"></div>
            </div>
            <input id="time-picker-input" style="display: none" type="text" :readonly="true"/>
            <div class="margin-top text-align-center">
                <div class="display-flex padding-horizontal justify-content-space-between">
                    <div class="align-self-center">{{ displayTime }}</div>
                    <f7-button outline :text="tt(switchButtonTitle)" @click="switchMode"></f7-button>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, useTemplateRef, onBeforeUnmount } from 'vue';
import VueDatePicker from '@vuepic/vue-datepicker';
import type { Picker } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, createInlinePicker } from '@/lib/ui/mobile.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';
import { useUserStore } from '@/stores/user.ts';

import { type WeekDayValue } from '@/core/datetime.ts';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import {
    getCurrentUnixTime,
    getCurrentYear,
    getUnixTime,
    getBrowserTimezoneOffsetMinutes,
    getLocalDatetimeFromUnixTime,
    getActualUnixTimeForStore,
    getTimezoneOffsetMinutes,
    getTimeValues,
    getCombinedDateAndTimeValues
} from '@/lib/datetime.ts';

type VueDatePickerType = InstanceType<typeof VueDatePicker>;

const props = defineProps<{
    modelValue: number;
    initMode?: string;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, getAllMinWeekdayNames, getAllMeridiemIndicators, getMonthShortName, formatUnixTimeToLongDateTime, isLongDateMonthAfterYear, isLongTime24HourFormat, isLongTimeMeridiemIndicatorFirst } = useI18n();
const { showToast } = useI18nUIComponents();

const environmentsStore = useEnvironmentsStore();
const userStore = useUserStore();

const is24Hour: boolean = isLongTime24HourFormat();
const isMeridiemIndicatorFirst: boolean = isLongTimeMeridiemIndicatorFirst() || false;
let timePickerHolder: Picker.Picker | null = null;

const mode = ref<string>(props.initMode || 'time');
const yearRange = ref<number[]>([
    2000,
    getCurrentYear() + 1
]);
const dateTime = ref<Date>(getLocalDatetimeFromUnixTime(props.modelValue || getCurrentUnixTime()));
const timeValues = ref<string[]>(getTimeValues(dateTime.value, is24Hour, isMeridiemIndicatorFirst));

const datetimepicker = useTemplateRef<VueDatePickerType>('datetimepicker');

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);
const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));
const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());
const displayTime = computed<string>(() => formatUnixTimeToLongDateTime(getActualUnixTimeForStore(getUnixTime(dateTime.value), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())));
const switchButtonTitle = computed<string>(() => mode.value === 'time' ? 'Date' : 'Time');

function switchMode(): void {
    if (mode.value === 'time') {
        mode.value = 'date';
    } else {
        mode.value = 'time';
        updateTimePicker(true);
    }
}

function setCurrentTime(): void {
    dateTime.value = getLocalDatetimeFromUnixTime(getCurrentUnixTime());
    timeValues.value = getTimeValues(dateTime.value, is24Hour, isMeridiemIndicatorFirst);

    if (mode.value === 'time') {
        updateTimePicker(false);
    }
}

function confirm(): void {
    if (!dateTime.value) {
        return;
    }

    const unixTime = getUnixTime(dateTime.value);

    if (unixTime < 0) {
        showToast('Date is too early');
        return;
    }

    emit('update:modelValue', unixTime);
    emit('update:show', false);
}

function getTimePickerColumns(): Picker.ColumnParameters[] {
    const ret: Picker.ColumnParameters[] = [];

    if (!is24Hour && isMeridiemIndicatorFirst) {
        ret.push(getAllMeridiemIndicators());
    }

    // Hours
    ret.push({
        values: generateAllHours()
    });
    // Divider
    ret.push({
        divider: true,
        content: ':',
    });
    // Minutes
    ret.push({
        values: generateAllMinutesOrSeconds()
    });
    // Divider
    ret.push({
        divider: true,
        content: ':',
    });
    // Seconds
    ret.push({
        values: generateAllMinutesOrSeconds()
    });

    if (!is24Hour && !isMeridiemIndicatorFirst) {
        ret.push(getAllMeridiemIndicators());
    }

    return ret;
}

function getDisplayTimeValue(value: number): string {
    if (value < 10) {
        return `0${value}`;
    } else {
        return value.toString();
    }
}

function generateAllHours(): string[] {
    const ret: string[] = [];
    const startHour = is24Hour ? 0 : 1;
    const endHour = is24Hour ? 23 : 11;

    if (!is24Hour) {
        ret.push('12');
    }

    for (let i = startHour; i <= endHour; i++) {
        ret.push(getDisplayTimeValue(i));
    }

    return ret;
}

function generateAllMinutesOrSeconds(): string[] {
    const ret: string[] = [];

    for (let i = 0; i < 60; i++) {
        ret.push(getDisplayTimeValue(i));
    }

    return ret;
}

function updateTimePicker(lazy: boolean): void {
    if (lazy) {
        nextTick(() => {
            if (timePickerHolder) {
                timePickerHolder.setValue(timeValues.value);
            }
        });
    } else {
        if (timePickerHolder) {
            timePickerHolder.setValue(timeValues.value);
        }
    }
}

function onSheetOpen(): void {
    mode.value = props.initMode || 'time';

    if (props.modelValue) {
        dateTime.value = getLocalDatetimeFromUnixTime(props.modelValue);
    }

    timeValues.value = getTimeValues(dateTime.value, is24Hour, isMeridiemIndicatorFirst);

    if (!timePickerHolder) {
        timePickerHolder = createInlinePicker('#time-picker-container', '#time-picker-input',
            getTimePickerColumns(), timeValues.value, {
                change(picker, values) {
                    if (mode.value === 'time' && values) {
                        timeValues.value = values as string[];
                        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, timeValues.value, is24Hour, isMeridiemIndicatorFirst);
                    }
                }
            });
    } else {
        updateTimePicker(true);
    }

    if (datetimepicker.value) {
        datetimepicker.value.switchView('calendar');
    }
}

function onSheetClosed(): void {
    emit('update:show', false);
}

onBeforeUnmount(() => {
    if (timePickerHolder) {
        timePickerHolder.destroy();
        timePickerHolder = null;
    }
});
</script>

<style>
.date-time-selection-sheet .dp__menu {
    border: 0;
}

.date-time-selection-sheet .time-picker-container .picker-columns {
    justify-content: space-evenly;
}
</style>
