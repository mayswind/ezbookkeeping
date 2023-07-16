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
            </vue-date-picker>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';

import { arrangeArrayWithNewStartIndex } from '@/lib/common.js';
import {
    getCurrentUnixTime,
    getCurrentDateTime,
    getUnixTime,
    getLocalDatetimeFromUnixTime,
    getYear
} from '@/lib/datetime.js';

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
        let value = getCurrentUnixTime();

        if (self.modelValue) {
            value = self.modelValue;
        }

        return {
            yearRange: [
                2000,
                getYear(getCurrentDateTime()) + 1
            ],
            dateTime: getLocalDatetimeFromUnixTime(value),
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
        is24Hour() {
            return this.$locale.isLongTime24HourFormat(this.userStore);
        }
    },
    methods: {
        onSheetOpen() {
            if (this.modelValue) {
                this.dateTime = getLocalDatetimeFromUnixTime(this.modelValue)
            }

            this.$refs.datetimepicker.switchView('calendar');
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        setCurrentTime() {
            this.dateTime = getLocalDatetimeFromUnixTime(getCurrentUnixTime())
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
        }
    }
}
</script>

<style>
.date-time-selection-sheet .dp__menu {
    border: 0;
}
</style>
