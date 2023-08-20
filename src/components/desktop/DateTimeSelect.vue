<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ 'content-class': 'date-time-select-menu' }"
        v-model="dateTime"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayTime }}</span>
        </template>

        <template #no-data>
            <vue-date-picker inline vertical time-picker-inline enable-seconds auto-apply
                             ref="datepicker"
                             month-name-format="long"
                             :clearable="false"
                             :dark="isDarkMode"
                             :week-start="firstDayOfWeek"
                             :year-range="yearRange"
                             :day-names="dayNames"
                             :is24="is24Hour"
                             v-model="dateTime">
                <template #month="{ text }">
                    {{ getMonthShortName(text) }}
                </template>
                <template #month-overlay-value="{ text }">
                    {{ getMonthShortName(text) }}
                </template>
                <template #am-pm-button="{ toggle, value }">
                    <button class="dp__pm_am_button" tabindex="0" @click="toggle">{{ $t(`datetime.${value}.content`) }}</button>
                </template>
            </vue-date-picker>
        </template>
    </v-select>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';

import { arrangeArrayWithNewStartIndex } from '@/lib/common.js';
import {
    getCurrentUnixTime,
    getCurrentDateTime,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getLocalDatetimeFromUnixTime,
    getActualUnixTimeForStore,
    getYear, getUnixTime
} from '@/lib/datetime.js';

export default {
    props: [
        'modelValue',
        'disabled',
        'readonly',
        'label'
    ],
    emits: [
        'update:modelValue',
        'error'
    ],
    data() {
        return {
            yearRange: [
                2000,
                getYear(getCurrentDateTime()) + 1
            ]
        }
    },
    computed: {
        ...mapStores(useUserStore),
        dateTime: {
            get: function () {
                return getLocalDatetimeFromUnixTime(this.modelValue);
            },
            set: function (value) {
                const unixTime = getUnixTime(value);

                if (unixTime < 0) {
                    this.$emit('error', 'Date is too early');
                    return;
                }

                this.$emit('update:modelValue', unixTime);
            }
        },
        isDarkMode() {
            return this.globalTheme.global.name.value === 'dark';
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        dayNames() {
            return arrangeArrayWithNewStartIndex(this.$locale.getAllMinWeekdayNames(), this.firstDayOfWeek);
        },
        is24Hour() {
            return this.$locale.isLongTime24HourFormat(this.userStore);
        },
        displayTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, getActualUnixTimeForStore(this.modelValue, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
        }
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    methods: {
        setCurrentTime() {
            this.dateTime = getLocalDatetimeFromUnixTime(getCurrentUnixTime())
        },
        getMonthShortName(month) {
            return this.$locale.getMonthShortName(month);
        }
    }
}
</script>

<style>
.date-time-select-menu {
    max-height: inherit !important;
}

.date-time-select-menu .dp__menu {
    border: 0;
}
</style>
