<template>
    <f7-sheet style="height:auto" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">{{ hint }}</p>
                <slot></slot>
                <f7-list no-hairlines inline-labels class="no-margin-top margin-bottom">
                    <f7-list-input
                        :label="$t('Begin Time')"
                        type="datepicker"
                        class="date-range-sheet-time-item"
                        :calendar-params="{
                            timePicker: true,
                            dateFormat: $t('input-format.datetime.long'),
                            firstDay: defaultFirstDayOfWeek,
                            toolbarCloseText: $t('Done'),
                            timePickerPlaceholder: $t('Select Time'),
                            timePickerFormat: $locale.getInputTimeIntlDateTimeFormatOptions(),
                            monthNames: $locale.getAllLongMonthNames(),
                            monthNamesShort: $locale.getAllShortMonthNames(),
                            dayNames: $locale.getAllLongWeekdayNames(),
                            dayNamesShort: $locale.getAllShortWeekdayNames()}"
                        :value="currentMinDate"
                        @calendar:change="currentMinDate = $event"
                    >
                    </f7-list-input>

                    <f7-list-input
                        :label="$t('End Time')"
                        type="datepicker"
                        class="date-range-sheet-time-item"
                        :calendar-params="{
                            timePicker: true,
                            dateFormat: $t('input-format.datetime.long'),
                            firstDay: defaultFirstDayOfWeek,
                            toolbarCloseText: $t('Done'),
                            timePickerPlaceholder: $t('Select Time'),
                            timePickerFormat: $locale.getInputTimeIntlDateTimeFormatOptions(),
                            monthNames: $locale.getAllLongMonthNames(),
                            monthNamesShort: $locale.getAllShortMonthNames(),
                            dayNames: $locale.getAllLongWeekdayNames(),
                            dayNamesShort: $locale.getAllShortWeekdayNames()}"
                        :value="currentMaxDate"
                        @calendar:change="currentMaxDate = $event"
                    >
                    </f7-list-input>
                </f7-list>
                <f7-button large fill
                           :class="{ 'disabled': !currentMinDate || !currentMaxDate }"
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
export default {
    props: [
        'minTime',
        'maxTime',
        'title',
        'hint',
        'show'
    ],
    data() {
        const self = this;
        let minDate = self.$utilities.getTodayFirstUnixTime();
        let maxDate = self.$utilities.getCurrentUnixTime();

        if (self.minTime) {
            minDate = self.minTime;
        }

        if (self.maxTime) {
            maxDate = self.maxTime;
        }

        minDate = self.$utilities.getDummyUnixTimeForLocalUsage(minDate, self.$utilities.getTimezoneOffsetMinutes(), self.$utilities.getBrowserTimezoneOffsetMinutes());
        maxDate = self.$utilities.getDummyUnixTimeForLocalUsage(maxDate, self.$utilities.getTimezoneOffsetMinutes(), self.$utilities.getBrowserTimezoneOffsetMinutes());

        return {
            currentMinDate: [self.$utilities.getLocalDatetimeFromUnixTime(minDate)],
            currentMaxDate: [self.$utilities.getLocalDatetimeFromUnixTime(maxDate)]
        }
    },
    computed: {
        defaultFirstDayOfWeek() {
            return this.$store.getters.currentUserFirstDayOfWeek;
        }
    },
    watch: {
        'currentMinDate': function (newValue) {
            if (!newValue) {
                this.currentMinDate = [this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getCurrentUnixTime())];
            }
        },
        'currentMaxDate': function (newValue) {
            if (!newValue) {
                this.currentMaxDate = [this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getCurrentUnixTime())];
            }
        }
    },
    methods: {
        onSheetOpen() {
            if (this.minTime) {
                const minTime = this.$utilities.getDummyUnixTimeForLocalUsage(this.minTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes());
                this.currentMinDate = [this.$utilities.getLocalDatetimeFromUnixTime(minTime)];
            }

            if (this.maxTime) {
                const maxTime = this.$utilities.getDummyUnixTimeForLocalUsage(this.maxTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes());
                this.currentMaxDate = [this.$utilities.getLocalDatetimeFromUnixTime(maxTime)];
            }
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        confirm() {
            if (!this.currentMinDate || !this.currentMaxDate) {
                return;
            }

            let currentMinDate = this.currentMinDate;

            if (this.$utilities.isArray(this.currentMinDate)) {
                currentMinDate = this.currentMinDate[0];
            }

            let currentMaxDate = this.currentMaxDate;

            if (this.$utilities.isArray(this.currentMaxDate)) {
                currentMaxDate = this.currentMaxDate[0];
            }

            let minUnixTime = this.$utilities.getMinuteFirstUnixTime(currentMinDate);
            let maxUnixTime = this.$utilities.getMinuteLastUnixTime(currentMaxDate);

            if (minUnixTime < 0 || maxUnixTime < 0) {
                this.$toast('Date is too early');
                return;
            }

            minUnixTime = this.$utilities.getActualUnixTimeForStore(minUnixTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes());
            maxUnixTime = this.$utilities.getActualUnixTimeForStore(maxUnixTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes());

            this.$emit('dateRange:change', minUnixTime, maxUnixTime);
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>

<style>
.list .date-range-sheet-time-item > .item-content {
    padding-left: 0;
}
</style>
