<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler"></div>
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top" v-if="hint">{{ hint }}</p>
                <p class="no-margin-top margin-bottom" v-if="beginDateTime && endDateTime">
                    <span>{{ beginDateTime }}</span>
                    <span> - </span>
                    <span>{{ endDateTime }}</span>
                </p>
                <slot></slot>
                <vue-date-picker range inline enable-seconds six-weeks
                                 auto-apply month-name-format="long"
                                 class="justify-content-center margin-bottom"
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
        let minDate = self.$utilities.getTodayFirstUnixTime();
        let maxDate = self.$utilities.getCurrentUnixTime();

        if (self.minTime) {
            minDate = self.minTime;
        }

        if (self.maxTime) {
            maxDate = self.maxTime;
        }

        return {
            yearRange: [
                2000,
                this.$utilities.getYear(this.$utilities.getCurrentDateTime()) + 1
            ],
            dateRange: [
                this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getDummyUnixTimeForLocalUsage(minDate, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes())),
                this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getDummyUnixTimeForLocalUsage(maxDate, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes()))
            ]
        }
    },
    computed: {
        isDarkMode() {
            return this.$root.isDarkMode;
        },
        firstDayOfWeek() {
            return this.$store.getters.currentUserFirstDayOfWeek;
        },
        dayNames() {
            return this.$locale.getAllMinWeekdayNames();
        },
        is24Hour() {
            const datetimeFormat = this.$t('format.datetime.long');
            return this.$utilities.is24HourFormat(datetimeFormat);
        },
        beginDateTime() {
            const actualBeginUnixTime = this.$utilities.getActualUnixTimeForStore(this.$utilities.getUnixTime(this.dateRange[0]), this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes());
            return this.$utilities.formatUnixTime(actualBeginUnixTime, this.$t('format.datetime.long'));
        },
        endDateTime() {
            const actualEndUnixTime = this.$utilities.getActualUnixTimeForStore(this.$utilities.getUnixTime(this.dateRange[1]), this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes());
            return this.$utilities.formatUnixTime(actualEndUnixTime, this.$t('format.datetime.long'));
        },
        presetRanges() {
            const presetRanges = [];

            [
                this.$constants.datetime.allDateRanges.Today,
                this.$constants.datetime.allDateRanges.LastSevenDays,
                this.$constants.datetime.allDateRanges.LastThirtyDays,
                this.$constants.datetime.allDateRanges.ThisWeek,
                this.$constants.datetime.allDateRanges.ThisMonth,
                this.$constants.datetime.allDateRanges.ThisYear
            ].forEach(dateRangeType => {
                const dateRange = this.$utilities.getDateRangeByDateType(dateRangeType.type, this.firstDayOfWeek);

                presetRanges.push({
                    label: this.$t(dateRangeType.name),
                    range: [
                        this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getDummyUnixTimeForLocalUsage(dateRange.minTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes())),
                        this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getDummyUnixTimeForLocalUsage(dateRange.maxTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes()))
                    ]
                });
            });

            return presetRanges;
        }
    },
    methods: {
        onSheetOpen() {
            if (this.minTime) {
                this.dateRange[0] = this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getDummyUnixTimeForLocalUsage(this.minTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes()));
            }

            if (this.maxTime) {
                this.dateRange[1] = this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getDummyUnixTimeForLocalUsage(this.maxTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes()));
            }
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

            let minUnixTime = this.$utilities.getUnixTime(currentMinDate);
            let maxUnixTime = this.$utilities.getUnixTime(currentMaxDate);

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
