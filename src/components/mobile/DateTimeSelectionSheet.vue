<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="date-time-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="$t('Current Time')" @click="setCurrentTime"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="$t('Done')" @click="confirm"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <vue-date-picker inline enable-seconds
                             auto-apply month-name-format="long"
                             class="justify-content-center"
                             :dark="isDarkMode"
                             :week-start="firstDayOfWeek"
                             :year-range="yearRange"
                             :day-names="dayNames"
                             :is24="is24Hour"
                             v-model="dateTime">
                <template #month="{ text }">
                    {{ $t(`datetime.${text}.short`) }}
                </template>
                <template #month-overlay-value="{ text }">
                    {{ $t(`datetime.${text}.short`) }}
                </template>
            </vue-date-picker>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
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
        setCurrentTime() {
            this.dateTime = this.$utilities.getLocalDatetimeFromUnixTime(this.$utilities.getCurrentUnixTime())
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
        }
    }
}
</script>

<style>
.date-time-selection-sheet .dp__menu {
    border: 0;
}
</style>
