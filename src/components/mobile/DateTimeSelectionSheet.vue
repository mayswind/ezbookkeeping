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
                <p class="no-margin-top margin-bottom" v-if="displayDateTime">
                    <span>{{ displayDateTime }}</span>
                </p>
                <slot></slot>
                <VueDatePicker inline enable-seconds
                               auto-apply month-name-format="long"
                               class="margin-bottom"
                               :dark="isDarkMode"
                               :week-start="firstDayOfWeek"
                               :year-range="yearRange"
                               :day-names="dayNames"
                               :is24="is24Hour"
                               v-model="dateTime">
                    <template #month="{ text, value }">
                        {{ $t(`datetime.${text}.short`) }}
                    </template>
                    <template #month-overlay-value="{ text, value }">
                        {{ $t(`datetime.${text}.short`) }}
                    </template>
                </VueDatePicker>
                <f7-button large fill
                           :class="{ 'disabled': !dateTime }"
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
        'modelValue',
        'title',
        'hint',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        const self = this;
        let value = self.$utilities.getCurrentUnixTime();

        if (self.modelValue) {
            value = self.modelValue;
        }

        return {
            yearRange: [
                2000,
                this.$utilities.getYear(this.$utilities.getCurrentDateTime()) + 1
            ],
            dateTime: this.$utilities.getLocalDatetimeFromUnixTime(value),
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
        displayDateTime() {
            const unixTime = this.$utilities.getUnixTime(this.dateTime);
            const actualUnixTime = this.$utilities.getActualUnixTimeForStore(unixTime, this.$utilities.getTimezoneOffsetMinutes(), this.$utilities.getBrowserTimezoneOffsetMinutes());
            return this.$utilities.formatUnixTime(actualUnixTime, this.$t('format.datetime.long'));
        }
    },
    methods: {
        onSheetOpen() {
            if (this.modelValue) {
                this.dateTime = this.$utilities.getLocalDatetimeFromUnixTime(this.modelValue)
            }
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        confirm() {
            if (!this.dateTime) {
                return;
            }

            const unixTime = this.$utilities.getUnixTime(this.dateTime);

            if (unixTime < 0) {
                this.$toast('Date is too early');
                return;
            }

            this.$emit('update:modelValue', unixTime);
            this.$emit('update:show', false);
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>
