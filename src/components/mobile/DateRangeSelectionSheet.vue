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
                        <button class="dp__pm_am_button" tabindex="0" @click="toggle">{{ $t(`datetime.${value}.content`) }}</button>
                    </template>
                </vue-date-picker>
                <f7-button large fill
                           :class="{ 'disabled': !dateRange[0] || !dateRange[1] }"
                           :text="$t('Continue')"
                           @click="confirm">
                </f7-button>
                <div class="margin-top text-align-center">
                    <f7-link @click="cancel" :text="$t('Cancel')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';

import datetimeConstants from '@/consts/datetime.js';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.js';
import {
    getCurrentUnixTime,
    getCurrentYear,
    getUnixTime,
    getLocalDatetimeFromUnixTime,
    getTodayFirstUnixTime,
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
                getCurrentYear() + 1
            ],
            dateRange: [
                getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(minDate, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())),
                getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(maxDate, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
            ]
        }
    },
    computed: {
        ...mapStores(useUserStore),
        isDarkMode() {
            return this.$root.isDarkMode;
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        dayNames() {
            return arrangeArrayWithNewStartIndex(this.$locale.getAllMinWeekdayNames(), this.firstDayOfWeek);
        },
        isYearFirst() {
            return this.$locale.isLongDateMonthAfterYear(this.userStore);
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
                    value: [
                        getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(dateRange.minTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes())),
                        getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(dateRange.maxTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()))
                    ]
                });
            });

            return presetRanges;
        }
    },
    methods: {
        onSheetOpen() {
            if (this.minTime) {
                this.dateRange[0] = getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(this.minTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
            }

            if (this.maxTime) {
                this.dateRange[1] = getLocalDatetimeFromUnixTime(getDummyUnixTimeForLocalUsage(this.maxTime, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
            }

            window.dispatchEvent(new Event('resize')); // fix vue-datepicker preset max-width
            this.$refs.datetimepicker.switchView('calendar');
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
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
        },
        getMonthShortName(month) {
            return this.$locale.getMonthShortName(month);
        }
    }
}
</script>
