<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="date-time-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="tt('Now')" @click="setCurrentTime"></f7-link>
            </div>
            <div class="right">
                <f7-link :icon-f7="mode === 'time' ? 'calendar' : 'clock'" @click="switchMode"></f7-link>
                <f7-button round fill icon-f7="checkmark_alt" @click="confirm"></f7-button>
            </div>
        </f7-toolbar>
        <f7-page-content class="padding-bottom">
            <div class="block no-margin no-padding">
                <date-time-picker ref="datetimepicker"
                                  datetime-picker-class="justify-content-center"
                                  :is-dark-mode="isDarkMode"
                                  :enable-time-picker="false"
                                  :show-alternate-dates="true"
                                  v-model="dateTime"
                                  v-show="mode === 'date'">
                </date-time-picker>
            </div>
            <div class="block no-margin no-padding padding-vertical-half" v-show="mode === 'time'">
                <div class="time-picker-container" ref="timePickerContainer">
                    <div class="picker picker-inline picker-3d">
                        <div class="picker-columns">
                            <div class="picker-column" v-if="!is24Hour && isMeridiemIndicatorFirst">
                                <div class="picker-items picker-items-meridiem-indicator-first"
                                     @scroll="onPickerColumnScroll('picker-items-meridiem-indicator-first', 'picker-meridiem-indicator', false)">
                                    <div :class="{ 'picker-item': true, 'picker-meridiem-indicator': true, 'picker-item-selected': currentMeridiemIndicator === item.value }"
                                         :key="item.value" :data-value="item.value"
                                         @click="currentMeridiemIndicator = item.value; scrollToSelectedItem('picker-items-meridiem-indicator-first', 'picker-meridiem-indicator', item.value)"
                                         v-for="item in meridiemItems">
                                        <span>{{ item.name }}</span>
                                    </div>
                                </div>
                            </div>
                            <div class="picker-column">
                                <div class="picker-items picker-items-hour"
                                     @scroll="onPickerColumnScroll('picker-items-hour', 'picker-hour', false)"
                                     @scrollend="onPickerColumnScroll('picker-items-hour', 'picker-hour', true)">
                                    <div :class="{ 'picker-item': true, 'picker-hour': true, 'picker-item-selected': currentHour === item.value }"
                                         :key="`${item.itemsIndex}_${item.value}`" :data-items-index="item.itemsIndex" :data-value="item.value"
                                         @click="currentHour = item.value; scrollToSelectedItem('picker-items-hour', 'picker-hour', item.value)"
                                         v-for="item in hourItems">
                                        <span :style="getTimerPickerItemStyle(item.value, currentHour, item.itemsIndex, hourItems)">{{ item.value }}</span>
                                    </div>
                                </div>
                            </div>
                            <div class="picker-column picker-column-divider">:</div>
                            <div class="picker-column">
                                <div class="picker-items picker-items-minute"
                                     @scroll="onPickerColumnScroll('picker-items-minute', 'picker-minute', false)"
                                     @scrollend="onPickerColumnScroll('picker-items-minute', 'picker-minute', true)">
                                    <div :class="{ 'picker-item': true, 'picker-minute': true, 'picker-item-selected': currentMinute === item.value }"
                                         :key="`${item.itemsIndex}_${item.value}`" :data-items-index="item.itemsIndex" :data-value="item.value"
                                         @click="currentMinute = item.value; scrollToSelectedItem('picker-items-minute', 'picker-minute', item.value)"
                                         v-for="item in minuteItems">
                                        <span :style="getTimerPickerItemStyle(item.value, currentMinute, item.itemsIndex, minuteItems)">{{ item.value }}</span>
                                    </div>
                                </div>
                            </div>
                            <div class="picker-column picker-column-divider">:</div>
                            <div class="picker-column">
                                <div class="picker-items picker-items-second"
                                     @scroll="onPickerColumnScroll('picker-items-second', 'picker-second', false)"
                                     @scrollend="onPickerColumnScroll('picker-items-second', 'picker-second', true)">
                                    <div :class="{ 'picker-item': true, 'picker-second': true, 'picker-item-selected': currentSecond === item.value }"
                                         :key="`${item.itemsIndex}_${item.value}`" :data-items-index="item.itemsIndex" :data-value="item.value"
                                         @click="currentSecond = item.value; scrollToSelectedItem('picker-items-second', 'picker-second', item.value)"
                                         v-for="item in secondItems">
                                        <span :style="getTimerPickerItemStyle(item.value, currentSecond, item.itemsIndex, secondItems)">{{ item.value }}</span>
                                    </div>
                                </div>
                            </div>
                            <div class="picker-column" v-if="!is24Hour && !isMeridiemIndicatorFirst">
                                <div class="picker-items picker-items-meridiem-indicator-last"
                                     @scroll="onPickerColumnScroll('picker-items-meridiem-indicator-last', 'picker-meridiem-indicator', false)">
                                    <div :class="{ 'picker-item': true, 'picker-meridiem-indicator': true, 'picker-item-selected': currentMeridiemIndicator === item.value }"
                                         :key="item.value" :data-value="item.value"
                                         @click="currentMeridiemIndicator = item.value; scrollToSelectedItem('picker-items-meridiem-indicator-last', 'picker-meridiem-indicator', item.value)"
                                         v-for="item in meridiemItems">
                                        <span>{{ item.name }}</span>
                                    </div>
                                </div>
                            </div>
                            <div class="picker-center-highlight"></div>
                        </div>
                    </div>
                </div>
            </div>
            <input id="time-picker-input" style="display: none" type="text" :readonly="true"/>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import DateTimePicker from '@/components/common/DateTimePicker.vue';
import { ref, computed, nextTick, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { type TimePickerValue, useDateTimeSelectionBase } from '@/components/base/DateTimeSelectionBase.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';

import { NumeralSystem } from '@/core/numeral.ts';

import { isDefined } from '@/lib/common.ts';
import {
    getHourIn12HourFormat,
    getLocalDatetimeFromUnixTime,
    getCurrentUnixTime,
    getUnixTimeFromLocalDatetime,
    getAMOrPM,
    getCombinedDateAndTimeValues
} from '@/lib/datetime.ts';

type DateTimePickerType = InstanceType<typeof DateTimePicker>;

const props = defineProps<{
    modelValue: number;
    initMode?: string;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
    (e: 'update:show', value: boolean): void;
}>();

const {
    tt,
    getCurrentNumeralSystemType
} = useI18n();
const { showToast } = useI18nUIComponents();

const {
    is24Hour,
    isHourTwoDigits,
    isMinuteTwoDigits,
    isSecondTwoDigits,
    isMeridiemIndicatorFirst,
    meridiemItems,
    getDisplayTimeValue,
    generateAllHours,
    generateAllMinutesOrSeconds
} = useDateTimeSelectionBase();

const environmentsStore = useEnvironmentsStore();

const datetimepicker = useTemplateRef<DateTimePickerType>('datetimepicker');
const timePickerContainer = useTemplateRef<HTMLDivElement>('timePickerContainer');

let resetTimePickerItemPositionItemsClass: string | undefined = undefined;
let resetTimePickerItemPositionItemClass: string | undefined = undefined;
let resetTimePickerItemPositionItemsLastOffsetTop: number | undefined = undefined;
let resetTimePickerItemPositionCheckedFrames: number | undefined = undefined;

const mode = ref<string>(props.initMode || 'time');
const dateTime = ref<Date>(getLocalDatetimeFromUnixTime(props.modelValue || getCurrentUnixTime()));
const timePickerContainerHeight = ref<number | undefined>(undefined);
const timePickerItemHeight = ref<number | undefined>(undefined);

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);
const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

const hourItems = computed<TimePickerValue[]>(() => generateAllHours(3, isHourTwoDigits.value));
const minuteItems = computed<TimePickerValue[]>(() => generateAllMinutesOrSeconds(3, isMinuteTwoDigits.value));
const secondItems = computed<TimePickerValue[]>(() => generateAllMinutesOrSeconds(3, isSecondTwoDigits.value));

const currentMeridiemIndicator = computed<string>({
    get: () => {
        return getAMOrPM(dateTime.value.getHours())
    },
    set: (value: string) => {
        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, numeralSystem.value, currentHour.value, currentMinute.value, currentSecond.value, value, is24Hour.value);
    }
});
const currentHour = computed<string>({
    get: () => {
        return getDisplayTimeValue(is24Hour.value ? dateTime.value.getHours() : getHourIn12HourFormat(dateTime.value.getHours()), isHourTwoDigits.value);
    },
    set: (value: string) => {
        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, numeralSystem.value, value, currentMinute.value, currentSecond.value, currentMeridiemIndicator.value, is24Hour.value);
    }
});
const currentMinute = computed<string>({
    get: () => {
        return getDisplayTimeValue(dateTime.value.getMinutes(), isMinuteTwoDigits.value);
    },
    set: (value: string) => {
        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, numeralSystem.value, currentHour.value, value, currentSecond.value, currentMeridiemIndicator.value, is24Hour.value);
    }
});
const currentSecond = computed<string>({
    get: () => {
        return getDisplayTimeValue(dateTime.value.getSeconds(), isSecondTwoDigits.value);
    },
    set: (value: string) => {
        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, numeralSystem.value, currentHour.value, currentMinute.value, value, currentMeridiemIndicator.value, is24Hour.value);
    }
});

function switchMode(): void {
    if (mode.value === 'time') {
        mode.value = 'date';
    } else {
        mode.value = 'time';
    }
}

function setCurrentTime(): void {
    dateTime.value = getLocalDatetimeFromUnixTime(getCurrentUnixTime());

    if (mode.value === 'time') {
        scrollAllTimeSelectedItems();
    }
}

function confirm(): void {
    if (!dateTime.value) {
        return;
    }

    const unixTime = getUnixTimeFromLocalDatetime(dateTime.value);

    if (unixTime < 0) {
        showToast('Date is too early');
        return;
    }

    emit('update:modelValue', unixTime);
    emit('update:show', false);
}

function getTimerPickerItemStyle(textualValue: string, textualCurrentValue: string, itemsIndex: number, values: TimePickerValue[]): string {
    if (!timePickerContainerHeight.value || !timePickerItemHeight.value) {
        return '';
    }

    const minValue = parseInt(values[0]!.value);
    const maxValue = parseInt(values[values.length - 1]!.value);
    const value = parseInt(textualValue, 10);
    const currentValue = parseInt(textualCurrentValue, 10);
    let valueDiff = value - currentValue;

    if (Math.abs(valueDiff) >= 5) {
        if (itemsIndex === 0 && maxValue - 5 < value && currentValue < minValue + 5) {
            valueDiff = value - (maxValue + currentValue + 1);
        } else if (itemsIndex === 2 && maxValue - 5 < currentValue && value < minValue + 5) {
            valueDiff = (maxValue + value + 1) - currentValue;
        }
    }

    const angle = -24 * valueDiff;

    if (angle > 180) {
        return '';
    }
    if (angle < -180) {
        return '';
    }

    return `transform: translate3d(0, ${-valueDiff * timePickerItemHeight.value}px, -100px) rotateX(${angle}deg)`;
}

function initTimePickerStyle(): void {
    const pickerItems = timePickerContainer.value?.querySelectorAll('.picker-item');
    const firstPickerItem = pickerItems ? pickerItems[0] : null;

    if (timePickerContainer.value) {
        timePickerContainerHeight.value = timePickerContainer.value.offsetHeight as number;
    }

    if (firstPickerItem && 'offsetHeight' in firstPickerItem) {
        timePickerItemHeight.value = firstPickerItem.offsetHeight as number;
    }

    if (timePickerContainer.value && firstPickerItem && 'offsetHeight' in firstPickerItem) {
        timePickerContainer.value.style.setProperty('--f7-picker-scroll-padding', `${(timePickerContainer.value.offsetHeight - (firstPickerItem.offsetHeight as number)) / 2}px`);
    }
}

function scrollAllTimeSelectedItems(): void {
    scrollToSelectedItem('picker-items-hour', 'picker-hour', currentHour.value);
    scrollToSelectedItem('picker-items-minute', 'picker-minute', currentMinute.value);
    scrollToSelectedItem('picker-items-second', 'picker-second', currentSecond.value);
    scrollToSelectedItem('picker-items-meridiem-indicator-first', 'picker-meridiem-indicator', currentMeridiemIndicator.value);
    scrollToSelectedItem('picker-items-meridiem-indicator-last', 'picker-meridiem-indicator', currentMeridiemIndicator.value);
}

function scrollTimeSelectedItems(itemsClass: string, itemClass: string): void {
    switch (resetTimePickerItemPositionItemClass) {
        case 'picker-hour':
            scrollToSelectedItem(itemsClass, itemClass, currentHour.value);
            break;
        case 'picker-minute':
            scrollToSelectedItem(itemsClass, itemClass, currentMinute.value);
            break;
        case 'picker-second':
            scrollToSelectedItem(itemsClass, itemClass, currentSecond.value);
            break;
    }
}

function scrollToSelectedItem(itemsClass: string, itemClass: string, value: string): void {
    const itemsElement = timePickerContainer.value?.querySelector(`.${itemsClass}`);
    const itemElements = itemsElement?.querySelectorAll(`.${itemClass}`);

    if (!itemsElement || !itemElements || !itemElements.length) {
        return;
    }

    for (let i = 0; i < itemElements.length; i++) {
        const itemElement = itemElements[i] as HTMLElement;

        if ('offsetHeight' in itemsElement && 'offsetTop' in itemElement && 'offsetHeight' in itemElement
            && (!itemElement.hasAttribute('data-items-index') || itemElement.getAttribute('data-items-index') === '1')
            && itemElement.getAttribute('data-value') === value) {
            itemsElement.scrollTop = (itemElement.offsetTop as number) - ((itemsElement.offsetHeight as number) / 2) + ((itemElement.offsetHeight as number) / 2);
            break;
        }
    }
}

function onPickerColumnScroll(itemsClass: string, itemClass: string, scrollEnd: boolean): void {
    const itemsElement = timePickerContainer.value?.querySelector(`.${itemsClass}`);
    const itemElements = itemsElement?.querySelectorAll(`.${itemClass}`);
    const firstPickerElement = itemElements ? itemElements[0] : null;

    if (!itemsElement || !itemElements || !itemElements.length || !firstPickerElement || !('offsetHeight' in firstPickerElement)) {
        return;
    }

    const itemHeight = firstPickerElement.offsetHeight as number;
    const scrollTop = itemsElement?.scrollTop || 0;
    const index = Math.round(scrollTop / itemHeight);
    const selectedItem = itemElements[index];

    if (selectedItem) {
        const value = selectedItem.getAttribute('data-value');
        const itemsIndex = selectedItem.getAttribute('data-items-index');

        if (value) {
            switch (itemClass) {
                case 'picker-hour':
                    currentHour.value = value;
                    break;
                case 'picker-minute':
                    currentMinute.value = value;
                    break;
                case 'picker-second':
                    currentSecond.value = value;
                    break;
                case 'picker-meridiem-indicator':
                    currentMeridiemIndicator.value = value;
                    break;
            }

            if (itemsIndex === '0' || itemsIndex === '2') {
                if (scrollEnd) {
                    scrollToSelectedItem(itemsClass, itemClass, value);
                } else {
                    if (resetTimePickerItemPositionItemsClass && resetTimePickerItemPositionItemClass
                        && resetTimePickerItemPositionItemsClass !== itemsClass && resetTimePickerItemPositionItemClass !== itemClass) {
                        scrollTimeSelectedItems(resetTimePickerItemPositionItemsClass, resetTimePickerItemPositionItemClass);
                        resetTimePickerItemPositionItemsClass = undefined;
                        resetTimePickerItemPositionItemClass = undefined;
                        resetTimePickerItemPositionItemsLastOffsetTop = undefined;
                        resetTimePickerItemPositionCheckedFrames = undefined;
                    }

                    if (!resetTimePickerItemPositionCheckedFrames && window.requestAnimationFrame) {
                        resetTimePickerItemPositionItemsClass = itemsClass;
                        resetTimePickerItemPositionItemClass = itemClass;
                        resetTimePickerItemPositionItemsLastOffsetTop = itemsElement.scrollTop;
                        resetTimePickerItemPositionCheckedFrames = 1;
                        window.requestAnimationFrame(delayCheckAndResetTimePickerItemPosition);
                    }
                }
            }
        }
    }
}

function delayCheckAndResetTimePickerItemPosition(): void {
    if (!resetTimePickerItemPositionItemsClass || !resetTimePickerItemPositionItemClass || !isDefined(resetTimePickerItemPositionItemsLastOffsetTop) || !isDefined(resetTimePickerItemPositionCheckedFrames)) {
        return;
    }

    const itemsElement = timePickerContainer.value?.querySelector(`.${resetTimePickerItemPositionItemsClass}`);

    if (!itemsElement) {
        return;
    }

    if (itemsElement.scrollTop === resetTimePickerItemPositionItemsLastOffsetTop) {
        resetTimePickerItemPositionCheckedFrames++;
    } else {
        resetTimePickerItemPositionItemsLastOffsetTop = itemsElement.scrollTop;
        resetTimePickerItemPositionCheckedFrames = 0;
    }

    if (resetTimePickerItemPositionCheckedFrames > 3) {
        scrollTimeSelectedItems(resetTimePickerItemPositionItemsClass, resetTimePickerItemPositionItemClass);
        resetTimePickerItemPositionItemsClass = undefined;
        resetTimePickerItemPositionItemClass = undefined;
        resetTimePickerItemPositionItemsLastOffsetTop = undefined;
        resetTimePickerItemPositionCheckedFrames = undefined;
        return;
    }

    window.requestAnimationFrame(delayCheckAndResetTimePickerItemPosition);
}

function onSheetOpen(): void {
    mode.value = props.initMode || 'time';

    if (props.modelValue) {
        dateTime.value = getLocalDatetimeFromUnixTime(props.modelValue);
    }

    if (mode.value === 'time') {
        nextTick(() => {
            initTimePickerStyle();
            scrollAllTimeSelectedItems();
        });
    }

    datetimepicker.value?.switchView('calendar');
}

function onSheetClosed(): void {
    emit('update:show', false);
}

watch(mode, (newValue) => {
    if (newValue === 'date') {
        datetimepicker.value?.switchView('calendar');
    } else if (newValue === 'time') {
        nextTick(() => {
            initTimePickerStyle();
            scrollAllTimeSelectedItems();
        });
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

.picker-hour,
.picker-minute,
.picker-second {
    font-variant-numeric: tabular-nums;
}
</style>
