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
            <input id="time-picker-input" style="display: none" type="text" readonly="readonly"/>
            <div class="margin-top text-align-center">
                <div class="display-flex padding-horizontal justify-content-space-between">
                    <div class="align-self-center">{{ displayTime }}</div>
                    <f7-button outline :text="switchButtonTitle" @click="switchMode"></f7-button>
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
    getBrowserTimezoneOffsetMinutes,
    getLocalDatetimeFromUnixTime,
    getActualUnixTimeForStore,
    getTimezoneOffsetMinutes,
    getTimeValues,
    getCombinedDateAndTimeValues
} from '@/lib/datetime.js';
import { createInlinePicker } from '@/lib/ui.mobile.js';

export default {
    props: [
        'modelValue',
        'initMode',
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
            mode: self.initMode || 'time',
            yearRange: [
                2000,
                getCurrentYear() + 1
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
        },
        displayTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, getActualUnixTimeForStore(getUnixTime(this.dateTime), getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
        },
        switchButtonTitle() {
            if (this.mode === 'time') {
                return this.$t('Date');
            } else {
                return this.$t('Time');
            }
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
            self.mode = self.initMode || 'time';

            if (self.modelValue) {
                self.dateTime = getLocalDatetimeFromUnixTime(self.modelValue);
            }

            self.timeValues = self.getTimeValues(self.dateTime);

            if (!self.timePickerHolder) {
                self.timePickerHolder = createInlinePicker('#time-picker-container', '#time-picker-input',
                    self.getTimePickerColumns(), self.timeValues, {
                        change(picker, values) {
                            if (self.mode === 'time') {
                                self.timeValues = values;
                                self.dateTime = getCombinedDateAndTimeValues(self.dateTime, self.timeValues, self.is24Hour, self.isMeridiemIndicatorFirst);
                            }
                        }
                    });
            } else {
                self.updateTimePicker(true);
            }

            self.$refs.datetimepicker.switchView('calendar');
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        switchMode() {
            if (this.mode === 'time') {
                this.mode = 'date';
            } else {
                this.mode = 'time';
                this.updateTimePicker(true);
            }
        },
        setCurrentTime() {
            this.dateTime = getLocalDatetimeFromUnixTime(getCurrentUnixTime());
            this.timeValues = this.getTimeValues(this.dateTime);

            if (this.mode === 'time') {
                this.updateTimePicker(false);
            }
        },
        confirm() {
            if (!this.dateTime) {
                return;
            }

            const unixTime = getUnixTime(this.dateTime);

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
        },
        updateTimePicker(lazy) {
            const self = this;

            if (lazy) {
                self.$nextTick(() => {
                    if (self.timePickerHolder) {
                        self.timePickerHolder.setValue(self.timeValues);
                    }
                });
            } else {
                if (self.timePickerHolder) {
                    self.timePickerHolder.setValue(self.timeValues);
                }
            }
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
