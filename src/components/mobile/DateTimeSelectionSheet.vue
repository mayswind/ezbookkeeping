<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="date-time-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="$t('Now')" @click="setCurrentTime"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="$t('Done')" @click="confirm"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
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
                             v-model="dateTime">
                <template #month="{ text }">
                    {{ getMonthShortName(text) }}
                </template>
                <template #month-overlay-value="{ text }">
                    {{ getMonthShortName(text) }}
                </template>
            </vue-date-picker>
            <div class="block block-outline no-margin no-padding padding-vertical-half">
                <div id="time-picker-container" class="time-picker-container"></div>
            </div>
            <input id="time-picker-input" style="display: none" type="text" readonly="readonly"/>
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
    getCurrentDateTime,
    getUnixTime,
    getLocalDatetimeFromUnixTime,
    getYear,
    getTimeValues,
    getCombinedDatetimeByDateAndTimeValues
} from '@/lib/datetime.js';
import { createInlinePicker } from '@/lib/ui.mobile.js';

export default {
    props: [
        'modelValue',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        const userStore = useUserStore();
        const self = this;

        self.is24Hour = self.$locale.isLongTime24HourFormat(userStore);
        self.isMeridiemIndicatorFirst = self.$locale.isLongTimeMeridiemIndicatorFirst(userStore);
        self.timePickerHolder = null;

        let value = getCurrentUnixTime();

        if (self.modelValue) {
            value = self.modelValue;
        }

        const datetime = getLocalDatetimeFromUnixTime(value);

        return {
            yearRange: [
                2000,
                getYear(getCurrentDateTime()) + 1
            ],
            dateTime: datetime,
            timeValues: self.getTimeValues(datetime),
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
        }
    },
    beforeUnmount() {
        if (this.timePickerHolder) {
            this.timePickerHolder.destroy();
            this.timePickerHolder = null;
        }
    },
    methods: {
        onSheetOpen() {
            const self = this;

            if (self.modelValue) {
                self.dateTime = getLocalDatetimeFromUnixTime(self.modelValue);
            }

            self.timeValues = self.getTimeValues(self.dateTime);

            if (!self.timePickerHolder) {
                self.timePickerHolder = createInlinePicker('#time-picker-container', '#time-picker-input',
                    self.getTimePickerColumns(), self.timeValues, {
                        change(picker, values) {
                            self.timeValues = values;
                        }
                    });
            } else {
                self.timePickerHolder.setValue(self.timeValues);
            }

            self.$refs.datetimepicker.switchView('calendar');
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        setCurrentTime() {
            this.dateTime = getLocalDatetimeFromUnixTime(getCurrentUnixTime())
            this.timeValues = this.getTimeValues(this.dateTime);

            if (this.timePickerHolder) {
                this.timePickerHolder.setValue(this.timeValues);
            }
        },
        confirm() {
            if (!this.dateTime) {
                return;
            }

            const finalDatetime = getCombinedDatetimeByDateAndTimeValues(this.dateTime, this.timeValues, this.is24Hour, this.isMeridiemIndicatorFirst);
            const unixTime = getUnixTime(finalDatetime);

            if (unixTime < 0) {
                this.$toast('Date is too early');
                return;
            }

            this.$emit('update:modelValue', unixTime);
            this.$emit('update:show', false);
        },
        getMonthShortName(month) {
            return this.$locale.getMonthShortName(month);
        },
        getTimeValues(datetime) {
            return getTimeValues(datetime, this.is24Hour, this.isMeridiemIndicatorFirst);
        },
        getTimePickerColumns() {
            const self = this;
            const ret = [];

            if (!self.is24Hour && this.isMeridiemIndicatorFirst) {
                ret.push({
                    values: datetimeConstants.allMeridiemIndicatorsArray,
                    displayValues: self.$locale.getAllMeridiemIndicatorNames()
                });
            }

            // Hours
            ret.push({
                values: self.generateAllHours()
            });
            // Divider
            ret.push({
                divider: true,
                content: ':',
            });
            // Minutes
            ret.push({
                values: self.generateAllMinutesOrSeconds()
            });
            // Divider
            ret.push({
                divider: true,
                content: ':',
            });
            // Seconds
            ret.push({
                values: self.generateAllMinutesOrSeconds()
            });

            if (!self.is24Hour && !this.isMeridiemIndicatorFirst) {
                ret.push({
                    values: datetimeConstants.allMeridiemIndicatorsArray,
                    displayValues: self.$locale.getAllMeridiemIndicatorNames()
                });
            }

            return ret;
        },
        getDisplayTimeValue(value) {
            if (value < 10) {
                return `0${value}`;
            } else {
                return value.toString();
            }
        },
        generateAllHours() {
            const ret = [];
            const startHour = this.is24Hour ? 0 : 1;
            const endHour = this.is24Hour ? 23 : 11;

            if (!this.is24Hour) {
                ret.push('12');
            }

            for (let i = startHour; i <= endHour; i++) {
                ret.push(this.getDisplayTimeValue(i));
            }

            return ret;
        },
        generateAllMinutesOrSeconds() {
            const ret = [];

            for (let i = 0; i < 60; i++) {
                ret.push(this.getDisplayTimeValue(i));
            }

            return ret;
        }
    }
}
</script>

<style>
.date-time-selection-sheet .dp__menu {
    border: 0;
}

.date-time-selection-sheet .time-picker-container .picker-columns {
    justify-content: space-evenly;
}
</style>
