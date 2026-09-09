<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :clearable="!emptyValue ? clearable : false"
        :label="label"
        :menu-props="{ contentClass: 'date-time-select-menu' }"
        v-model="dateTime"
        @paste="onPaste"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayTime }}</span>
        </template>

        <template #no-data>
            <date-time-picker :is-dark-mode="isDarkMode"
                              :enable-time-picker="false"
                              :vertical="true"
                              :show-alternate-dates="true"
                              v-model="dateTime">
            </date-time-picker>
            <div class="date-time-select-time-picker-container"
                 @focusin="onFocused"
                 @click="onFocused"
                 @keydown="onKeyDown">
                <v-btn class="px-3" color="primary" variant="flat"
                       v-if="!is24Hour && isMeridiemIndicatorFirst"
                       @click="toggleMeridiemIndicator">
                    {{ tt(`datetime.${currentMeridiemIndicator}.content`) }}
                </v-btn>
                <v-autocomplete eager
                                density="compact"
                                max-width="70px"
                                item-title="value"
                                item-value="value"
                                auto-select-first="exact"
                                :items="hourItems"
                                :hide-no-data="true"
                                v-model="currentHour"
                />
                <span>:</span>
                <v-autocomplete eager
                                density="compact"
                                max-width="70px"
                                item-title="value"
                                item-value="value"
                                auto-select-first="exact"
                                :items="minuteItems"
                                :hide-no-data="true"
                                v-model="currentMinute"
                />
                <span>:</span>
                <v-autocomplete eager
                                density="compact"
                                max-width="70px"
                                item-title="value"
                                item-value="value"
                                auto-select-first="exact"
                                :items="secondItems"
                                :hide-no-data="true"
                                v-model="currentSecond"
                />
                <v-btn class="px-3" color="primary" variant="flat"
                       v-if="!is24Hour && !isMeridiemIndicatorFirst"
                       @click="toggleMeridiemIndicator">
                    {{ tt(`datetime.${currentMeridiemIndicator}.content`) }}
                </v-btn>
            </div>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { computed, nextTick } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { type TimePickerValue, useDateTimeSelectionBase } from '@/components/base/DateTimeSelectionBase.ts';

import { ThemeType } from '@/core/theme.ts';
import { NumeralSystem } from '@/core/numeral.ts';
import {
    type DateTime,
    type DateFormatOrder,
    MeridiemIndicator,
    KnownDateTimeFormat
} from '@/core/datetime.ts';
import {
    getHourIn12HourFormat,
    getLocalDatetimeFromUnixTime,
    getSameDateTimeWithBrowserTimezone,
    parseDateTimeFromUnixTimeWithTimezoneOffset,
    parseDateTimeFromKnownDateTimeFormat,
    getAMOrPM,
    getCombinedDateAndTimeValues
} from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue: number;
    timezoneUtcOffset: number;
    emptyValue?: boolean;
    disabled?: boolean;
    readonly?: boolean;
    clearable?: boolean;
    label?: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
    (e: 'clear:modelValue'): void;
    (e: 'error', message: string): void;
}>();

const theme = useTheme();
const {
    tt,
    getCurrentNumeralSystemType,
    getLongDateFormatOrder,
    getShortDateFormatOrder,
    parseDateTimeFromLongDateTime,
    parseDateTimeFromShortDateTime,
    formatDateTimeToLongDateTime
} = useI18n();

const {
    is24Hour,
    isHourTwoDigits,
    isMinuteTwoDigits,
    isSecondTwoDigits,
    isMeridiemIndicatorFirst,
    getLocalDatetimeFromSameDateTimeOfUnixTime,
    getUnixTimeFromSameDateTimeOfLocalDatetime,
    getDisplayTimeValue,
    generateAllHours,
    generateAllMinutesOrSeconds
} = useDateTimeSelectionBase();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const longDateFormatOrder = computed<DateFormatOrder>(() => getLongDateFormatOrder());
const shortDateFormatOrder = computed<DateFormatOrder>(() => getShortDateFormatOrder());

const dateTime = computed<Date>({
    get: () => {
        return getLocalDatetimeFromSameDateTimeOfUnixTime(props.modelValue, props.timezoneUtcOffset);
    },
    set: (value: Date | null) => {
        if (!value) {
            emit('clear:modelValue');
            return;
        }

        const unixTime = getUnixTimeFromSameDateTimeOfLocalDatetime(value, props.timezoneUtcOffset);

        if (unixTime < 0) {
            emit('error', 'Date is too early');
            return;
        }

        emit('update:modelValue', unixTime);
    }
});

const displayTime = computed<string>(() => props.emptyValue ? tt('None') : formatDateTimeToLongDateTime(parseDateTimeFromUnixTimeWithTimezoneOffset(props.modelValue, props.timezoneUtcOffset)));

const hourItems = computed<TimePickerValue[]>(() => generateAllHours(1, isHourTwoDigits.value));
const minuteItems = computed<TimePickerValue[]>(() => generateAllMinutesOrSeconds(1, isMinuteTwoDigits.value));
const secondItems = computed<TimePickerValue[]>(() => generateAllMinutesOrSeconds(1, isSecondTwoDigits.value));

const currentMeridiemIndicator = computed<string>({
    get: () => {
        return getAMOrPM(dateTime.value.getHours())
    },
    set: (value: string) => {
        if (value !== MeridiemIndicator.AM.name && value !== MeridiemIndicator.PM.name) {
            return;
        }

        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, numeralSystem.value, currentHour.value, currentMinute.value, currentSecond.value, value, is24Hour.value);
    }
});
const currentHour = computed<string>({
    get: () => {
        return getDisplayTimeValue(is24Hour.value ? dateTime.value.getHours() : getHourIn12HourFormat(dateTime.value.getHours()), isHourTwoDigits.value);
    },
    set: (value: string) => {
        const hour = numeralSystem.value.parseInt(value);

        if (isNaN(hour) || hour < 0 || (is24Hour.value ? hour > 23 : hour > 12)) {
            return;
        }

        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, numeralSystem.value, value, currentMinute.value, currentSecond.value, currentMeridiemIndicator.value, is24Hour.value);
    }
});
const currentMinute = computed<string>({
    get: () => {
        return getDisplayTimeValue(dateTime.value.getMinutes(), isMinuteTwoDigits.value);
    },
    set: (value: string) => {
        const minute = numeralSystem.value.parseInt(value);

        if (isNaN(minute) || minute < 0 || minute > 59) {
            return;
        }

        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, numeralSystem.value, currentHour.value, value, currentSecond.value, currentMeridiemIndicator.value, is24Hour.value);
    }
});
const currentSecond = computed<string>({
    get: () => {
        return getDisplayTimeValue(dateTime.value.getSeconds(), isSecondTwoDigits.value);
    },
    set: (value: string) => {
        const second = numeralSystem.value.parseInt(value);

        if (isNaN(second) || second < 0 || second > 59) {
            return;
        }

        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, numeralSystem.value, currentHour.value, currentMinute.value, value, currentMeridiemIndicator.value, is24Hour.value);
    }
});

function setTimeInputFocus(container: HTMLElement, inputIndex: number): void {
    nextTick(() => {
        setTimeout(() => {
            const input = container.querySelectorAll<HTMLInputElement>('input[role="combobox"]')[inputIndex];

            input?.focus();
            input?.select();
        }, 50);
    });
}

function toggleMeridiemIndicator(): void {
    if (currentMeridiemIndicator.value === MeridiemIndicator.AM.name) {
        currentMeridiemIndicator.value = MeridiemIndicator.PM.name;
    } else {
        currentMeridiemIndicator.value = MeridiemIndicator.AM.name;
    }
}

function onPaste(event: ClipboardEvent): void {
    if (!event.clipboardData || props.readonly || props.disabled) {
        event.preventDefault();
        return;
    }

    let text = event.clipboardData.getData('Text');

    if (!text) {
        event.preventDefault();
        return;
    }

    text = text.trim();

    const formats = KnownDateTimeFormat.detect(text, longDateFormatOrder.value, shortDateFormatOrder.value);
    let dt: DateTime | undefined = undefined;

    if (formats && (formats.length === 1 || (formats.length > 1 && formats[0]!.type === longDateFormatOrder.value && formats[0]!.type === shortDateFormatOrder.value))) {
        dt = parseDateTimeFromKnownDateTimeFormat(text, formats[0] as KnownDateTimeFormat);

        if (dt) {
            dateTime.value = getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(dt).getUnixTime());
            return;
        }
    }

    dt = parseDateTimeFromLongDateTime(text);

    if (dt) {
        dateTime.value = getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(dt).getUnixTime());
        return;
    }

    dt = parseDateTimeFromShortDateTime(text);

    if (dt) {
        dateTime.value = getLocalDatetimeFromUnixTime(getSameDateTimeWithBrowserTimezone(dt).getUnixTime());
        return;
    }

    event.preventDefault();
}

function onFocused(e: Event): void {
    if (e.target instanceof HTMLInputElement && e.target.role === 'combobox') {
        const input = e.target;

        nextTick(() => {
            input.focus();
            input.select();
        });
    }
}

function onKeyDown(e: KeyboardEvent): void {
    if (!(e.currentTarget instanceof HTMLElement) || !(e.target instanceof HTMLInputElement)) {
        return;
    }

    const container = e.currentTarget;
    const inputs = Array.from(container.querySelectorAll<HTMLInputElement>('input[role="combobox"]'));
    const inputIndex = inputs.indexOf(e.target);
    const type = ['hour', 'minute', 'second'][inputIndex];

    if (!type) {
        return;
    }

    if (e.altKey || e.ctrlKey || e.metaKey || (e.key.indexOf('F') === 0 && (e.key.length === 2 || e.key.length === 3))
        || e.key === 'ArrowLeft' || e.key === 'ArrowRight'
        || e.key === 'Home' || e.key === 'End'
        || e.key === 'Backspace' || e.key === 'Delete' || e.key === 'Del') {
        return;
    }

    if (e.key.length === 1 && numeralSystem.value.isDigit(e.key)) {
        return;
    }

    let value = '';

    if (e.target instanceof HTMLInputElement) {
        const input = e.target as HTMLInputElement;
        value = input.value;
    }

    if (value && (e.key === 'Tab' || e.key === 'Enter')) {
        if (type === 'hour') {
            currentHour.value = value;
        } else if (type === 'minute') {
            currentMinute.value = value;
        } else if (type === 'second') {
            currentSecond.value = value;
        }
    }

    if (e.shiftKey && e.key === 'Tab') {
        if (type === 'minute') {
            setTimeInputFocus(container, 0);
        } else if (type === 'second') {
            setTimeInputFocus(container, 1);
        }

        e.preventDefault();
        e.stopPropagation();
        return;
    }

    if (!e.shiftKey && (e.key === 'Tab' || e.key === 'Enter')) {
        if (type === 'hour') {
            setTimeInputFocus(container, 1);

            e.preventDefault();
            e.stopPropagation();
            return;
        } else if (type === 'minute') {
            setTimeInputFocus(container, 2);

            e.preventDefault();
            e.stopPropagation();
            return;
        }
    }

    e.preventDefault();
}
</script>

<style>
.date-time-select-menu {
    max-height: inherit !important;
}

.date-time-select-menu .dp--menu {
    border: 0;
}

.date-time-select-time-picker-container {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--dp-menu-padding);
    padding-bottom: 0;
    column-gap: 8px;
}

.date-time-select-time-picker-container .v-autocomplete.v-input--density-compact {
    --v-input-control-height: 38px;
    --v-field-input-padding-top: 4px;
    --v-field-input-padding-bottom: 4px;
}

.date-time-select-time-picker-container .v-autocomplete.v-input--density-compact .v-field {
    --v-field-padding-start: 12px;
    --v-field-padding-end: 0;
}

.date-time-select-time-picker-container .v-autocomplete.v-input--density-compact .v-field__input {
    min-height: 38px !important;
}

.date-time-select-time-picker-container .v-autocomplete.v-input--density-compact .v-field__append-inner .v-autocomplete__menu-icon {
    margin-inline-start: 0;
}

.date-time-select-time-picker-container .v-autocomplete .v-field--appended {
    padding-inline-end: 8px;
}
</style>
