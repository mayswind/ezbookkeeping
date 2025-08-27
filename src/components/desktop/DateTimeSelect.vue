<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ contentClass: 'date-time-select-menu' }"
        v-model="dateTime"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayTime }}</span>
        </template>

        <template #no-data>
            <date-time-picker :is-dark-mode="isDarkMode"
                              :enable-time-picker="false"
                              :vertical="true"
                              v-model="dateTime">
            </date-time-picker>
            <div class="date-time-select-time-picker-container">
                <v-btn class="px-3" color="primary" variant="flat"
                       v-if="!is24Hour && isMeridiemIndicatorFirst"
                       @click="toggleMeridiemIndicator">
                    {{ tt(`datetime.${currentMeridiemIndicator}.content`) }}
                </v-btn>
                <v-autocomplete eager ref="hourInput"
                                density="compact"
                                max-width="70px"
                                item-title="value"
                                item-value="value"
                                auto-select-first="exact"
                                :items="hourItems"
                                :hide-no-data="true"
                                v-model="currentHour"
                                @update:focused="onFocused(hourInput, $event)"
                                @click="onFocused(hourInput, true)"
                                @keydown="onKeyDown('hour', $event)"
                />
                <span>:</span>
                <v-autocomplete eager ref="minuteInput"
                                density="compact"
                                max-width="70px"
                                item-title="value"
                                item-value="value"
                                auto-select-first="exact"
                                :items="minuteItems"
                                :hide-no-data="true"
                                v-model="currentMinute"
                                @update:focused="onFocused(minuteInput, $event)"
                                @click="onFocused(minuteInput, true)"
                                @keydown="onKeyDown('minute', $event)"
                />
                <span>:</span>
                <v-autocomplete eager ref="secondInput"
                                density="compact"
                                max-width="70px"
                                item-title="value"
                                item-value="value"
                                auto-select-first="exact"
                                :items="secondItems"
                                :hide-no-data="true"
                                v-model="currentSecond"
                                @update:focused="onFocused(secondInput, $event)"
                                @click="onFocused(secondInput, true)"
                                @keydown="onKeyDown('second', $event)"
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
import { VAutocomplete } from 'vuetify/components/VAutocomplete';

import { computed, useTemplateRef, nextTick } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { type TimePickerValue, useDateTimeSelectionBase } from '@/components/base/DateTimeSelectionBase.ts';

import { ThemeType } from '@/core/theme.ts';
import { MeridiemIndicator } from '@/core/datetime.ts';
import {
    getHourIn12HourFormat,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getLocalDatetimeFromUnixTime,
    getActualUnixTimeForStore,
    getUnixTimeFromLocalDatetime,
    getAMOrPM,
    getCombinedDateAndTimeValues
} from '@/lib/datetime.ts';
import { setChildInputFocus } from '@/lib/ui/desktop.ts';

const props = defineProps<{
    modelValue: number;
    disabled?: boolean;
    readonly?: boolean;
    label?: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
    (e: 'error', message: string): void;
}>();

const theme = useTheme();
const { tt, formatUnixTimeToLongDateTime } = useI18n();

const {
    is24Hour,
    isHourTwoDigits,
    isMinuteTwoDigits,
    isSecondTwoDigits,
    isMeridiemIndicatorFirst,
    getDisplayTimeValue,
    generateAllHours,
    generateAllMinutesOrSeconds
} = useDateTimeSelectionBase();

const hourInput = useTemplateRef<VAutocomplete>('hourInput');
const minuteInput = useTemplateRef<VAutocomplete>('minuteInput');
const secondInput = useTemplateRef<VAutocomplete>('secondInput');

const dateTime = computed<Date>({
    get: () => {
        return getLocalDatetimeFromUnixTime(props.modelValue);
    },
    set: (value: Date) => {
        const unixTime = getUnixTimeFromLocalDatetime(value);

        if (unixTime < 0) {
            emit('error', 'Date is too early');
            return;
        }

        emit('update:modelValue', unixTime);
    }
});

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

        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, currentHour.value, currentMinute.value, currentSecond.value, value, is24Hour.value);
    }
});
const currentHour = computed<string>({
    get: () => {
        return getDisplayTimeValue(is24Hour.value ? dateTime.value.getHours() : getHourIn12HourFormat(dateTime.value.getHours()), isHourTwoDigits.value);
    },
    set: (value: string) => {
        const hour = parseInt(value);

        if (isNaN(hour) || hour < 0 || (is24Hour.value ? hour > 23 : hour > 12)) {
            return;
        }

        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, value, currentMinute.value, currentSecond.value, currentMeridiemIndicator.value, is24Hour.value);
    }
});
const currentMinute = computed<string>({
    get: () => {
        return getDisplayTimeValue(dateTime.value.getMinutes(), isMinuteTwoDigits.value);
    },
    set: (value: string) => {
        const minute = parseInt(value);

        if (isNaN(minute) || minute < 0 || minute > 59) {
            return;
        }

        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, currentHour.value, value, currentSecond.value, currentMeridiemIndicator.value, is24Hour.value);
    }
});
const currentSecond = computed<string>({
    get: () => {
        return getDisplayTimeValue(dateTime.value.getSeconds(), isSecondTwoDigits.value);
    },
    set: (value: string) => {
        const second = parseInt(value);

        if (isNaN(second) || second < 0 || second > 59) {
            return;
        }

        dateTime.value = getCombinedDateAndTimeValues(dateTime.value, currentHour.value, currentMinute.value, value, currentMeridiemIndicator.value, is24Hour.value);
    }
});

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const displayTime = computed<string>(() => formatUnixTimeToLongDateTime(getActualUnixTimeForStore(getUnixTimeFromLocalDatetime(dateTime.value), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())));

function toggleMeridiemIndicator(): void {
    if (currentMeridiemIndicator.value === MeridiemIndicator.AM.name) {
        currentMeridiemIndicator.value = MeridiemIndicator.PM.name;
    } else {
        currentMeridiemIndicator.value = MeridiemIndicator.AM.name;
    }
}

function onFocused(input: VAutocomplete | null | undefined, focused: boolean): void {
    if (input && focused) {
        nextTick(() => {
            setChildInputFocus(input?.$el, 'input');
        });
    }
}

function onKeyDown(type: string, e: KeyboardEvent): void {
    if (e.altKey || e.ctrlKey || e.metaKey || (e.key.indexOf('F') === 0 && (e.key.length === 2 || e.key.length === 3))
        || e.key === 'ArrowLeft' || e.key === 'ArrowRight'
        || e.key === 'Home' || e.key === 'End'
        || e.key === 'Backspace' || e.key === 'Delete' || e.key === 'Del') {
        return;
    }

    if (e.key.length === 1 && '0' <= e.key && e.key <= '9') {
        return;
    }

    let value = '';

    if (e.target instanceof HTMLInputElement) {
        const input = e.target as HTMLInputElement;
        value = input.value;
    }

    if (!value) {
        return;
    }

    if (e.key === 'Tab' || e.key === 'Enter') {
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
            nextTick(() => {
                setTimeout(() => {
                    setChildInputFocus(hourInput.value?.$el, 'input');
                }, 50);
            });
        } else if (type === 'second') {
            nextTick(() => {
                setTimeout(() => {
                    setChildInputFocus(minuteInput.value?.$el, 'input');
                }, 50);
            });
        }

        e.preventDefault();
        e.stopPropagation();
        return;
    }

    if (!e.shiftKey && (e.key === 'Tab' || e.key === 'Enter')) {
        if (type === 'hour') {
            nextTick(() => {
                setTimeout(() => {
                    setChildInputFocus(minuteInput.value?.$el, 'input');
                }, 50);
            });
        } else if (type === 'minute') {
            nextTick(() => {
                setTimeout(() => {
                    setChildInputFocus(secondInput.value?.$el, 'input');
                }, 50);
            });
        }
    }

    e.preventDefault();
}
</script>

<style>
.date-time-select-menu {
    max-height: inherit !important;
}

.date-time-select-menu .dp__menu {
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
