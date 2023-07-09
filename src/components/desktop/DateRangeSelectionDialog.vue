<template>
    <v-dialog width="460" v-model="showState">
        <v-card>
            <v-toolbar color="primary">
                <v-toolbar-title>{{ title }}</v-toolbar-title>
            </v-toolbar>
            <v-card-text class="pa-4">
                <p v-if="hint">{{ hint }}</p>
                <span v-if="beginDateTime && endDateTime">
                    <span>{{ beginDateTime }}</span>
                    <span> - </span>
                    <span>{{ endDateTime }}</span>
                </span>
                <slot></slot>
            </v-card-text>
            <v-card-text class="pa-4 w-100 d-flex justify-center">
                <vue-date-picker range inline enable-seconds auto-apply
                                 ref="datetimepicker"
                                 month-name-format="long"
                                 six-weeks="center"
                                 :clearable="false"
                                 :dark="isDarkMode"
                                 :week-start="firstDayOfWeek"
                                 :year-range="yearRange"
                                 :day-names="dayNames"
                                 :is24="is24Hour"
                                 :partial-range="false"
                                 :preset-ranges="presetRanges"
                                 v-model="dateRange">
                    <template #month="{ text }">
                        {{ $t(`datetime.${text}.short`) }}
                    </template>
                    <template #month-overlay-value="{ text }">
                        {{ $t(`datetime.${text}.short`) }}
                    </template>
                </vue-date-picker>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="gray" @click="cancel">{{ $t('Cancel') }}</v-btn>
                <v-btn :disabled="!dateRange[0] || !dateRange[1]" @click="confirm">{{ $t('OK') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';

import datetimeConstants from '@/consts/datetime.js';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.js';
import {
    getCurrentUnixTime,
    getCurrentDateTime,
    getUnixTime,
    getLocalDatetimeFromUnixTime,
    getTodayFirstUnixTime,
    getYear,
    getDummyUnixTimeForLocalUsage,
    getActualUnixTimeForStore,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getDateRangeByDateType
} from '@/lib/datetime.js';

export default {
    props: [
        'minTime',
        'maxTime',
        'title',
        'hint',
        'show'
    ],
    emits: [
        'update:show',
        'dateRange:change'
    ],
    data() {
        const self = this;
        let minDate = getTodayFirstUnixTime();
        let maxDate = getCurrentUnixTime();

        if (self.minTime) {
            minDate = self.minTime;
        }

        if (self.maxTime) {
            maxDate = self.maxTime;
        }

        return {
            yearRange: [
                2000,
                getYear(getCurrentDateTime()) + 1
            ],
            dateRange: [
                getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(minDate, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())),
                getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(maxDate, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
            ]
        }
    },
    computed: {
        ...mapStores(useUserStore),
        showState: {
            get: function () {
                return this.show;
            },
            set: function (value) {
                this.$emit('update:show', value);
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
        beginDateTime() {
            const actualBeginUnixTime = getActualUnixTimeForStore(getUnixTime(this.dateRange[0]), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes());
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, actualBeginUnixTime);
        },
        endDateTime() {
            const actualEndUnixTime = getActualUnixTimeForStore(getUnixTime(this.dateRange[1]), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes());
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, actualEndUnixTime);
        },
        presetRanges() {
            const presetRanges = [];

            [
                datetimeConstants.allDateRanges.Today,
                datetimeConstants.allDateRanges.LastSevenDays,
                datetimeConstants.allDateRanges.LastThirtyDays,
                datetimeConstants.allDateRanges.ThisWeek,
                datetimeConstants.allDateRanges.ThisMonth,
                datetimeConstants.allDateRanges.ThisYear
            ].forEach(dateRangeType => {
                const dateRange = getDateRangeByDateType(dateRangeType.type, this.firstDayOfWeek);

                presetRanges.push({
                    label: this.$t(dateRangeType.name),
                    range: [
                        getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(dateRange.minTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())),
                        getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(dateRange.maxTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
                    ]
                });
            });

            return presetRanges;
        }
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    methods: {
        confirm() {
            if (!this.dateRange[0] || !this.dateRange[1]) {
                return;
            }

            const currentMinDate = this.dateRange[0];
            const currentMaxDate = this.dateRange[1];

            let minUnixTime = getUnixTime(currentMinDate);
            let maxUnixTime = getUnixTime(currentMaxDate);

            if (minUnixTime < 0 || maxUnixTime < 0) {
                this.$toast('Date is too early');
                return;
            }

            minUnixTime = getActualUnixTimeForStore(minUnixTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes());
            maxUnixTime = getActualUnixTimeForStore(maxUnixTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes());

            this.$emit('dateRange:change', minUnixTime, maxUnixTime);
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>
