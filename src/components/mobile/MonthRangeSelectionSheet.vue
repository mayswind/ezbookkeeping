<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="month-range-selection-sheet" style="height:auto"
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
                <vue-date-picker inline month-picker auto-apply
                                 month-name-format="long"
                                 class="justify-content-center margin-bottom"
                                 :clearable="false"
                                 :dark="isDarkMode"
                                 :year-range="yearRange"
                                 :year-first="isYearFirst"
                                 :range="{ partialRange: false }"
                                 v-model="dateRange">
                    <template #month="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                    <template #month-overlay-value="{ text }">
                        {{ getMonthShortName(text) }}
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
import { useUserStore } from '@/stores/user.ts';

import {
    getYearMonthObjectFromString,
    getYearMonthStringFromObject,
    getCurrentUnixTime,
    getCurrentYear,
    getThisYearFirstUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime
} from '@/lib/datetime.ts';

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
        let minDate = getThisYearFirstUnixTime();
        let maxDate = getCurrentUnixTime();

        if (self.minTime) {
            minDate = getYearMonthObjectFromString(self.minTime);
        }

        if (self.maxTime) {
            maxDate = getYearMonthObjectFromString(self.maxTime);
        }

        return {
            yearRange: [
                2000,
                getCurrentYear() + 1
            ],
            dateRange: [
                minDate,
                maxDate
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
        isYearFirst() {
            return this.$locale.isLongDateMonthAfterYear(this.userStore);
        },
        is24Hour() {
            return this.$locale.isLongTime24HourFormat(this.userStore);
        },
        beginDateTime() {
            return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, getYearMonthFirstUnixTime(this.dateRange[0]));
        },
        endDateTime() {
            return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, getYearMonthLastUnixTime(this.dateRange[1]));
        }
    },
    methods: {
        onSheetOpen() {
            if (this.minTime) {
                this.dateRange[0] = getYearMonthObjectFromString(this.minTime);
            }

            if (this.maxTime) {
                this.dateRange[1] = getYearMonthObjectFromString(this.maxTime);
            }
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        confirm() {
            if (!this.dateRange[0] || !this.dateRange[1]) {
                return;
            }

            if (this.dateRange[0].year <= 0 || this.dateRange[0].month < 0 || this.dateRange[1].year <= 0 || this.dateRange[1].month < 0) {
                this.$toast('Date is too early');
                return;
            }

            const minYearMonth = getYearMonthStringFromObject(this.dateRange[0]);
            const maxYearMonth = getYearMonthStringFromObject(this.dateRange[1]);

            this.$emit('dateRange:change', minYearMonth, maxYearMonth);
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

<style>
.month-range-selection-sheet .dp__main .dp__instance_calendar .dp__overlay.dp--overlay-relative {
    width: 100% !important;
}

.month-range-selection-sheet .dp__main .dp__instance_calendar .dp__overlay.dp--overlay-relative .dp__selection_grid_header .dp--year-mode-picker .dp--arrow-btn-nav {
    display: flex;
}

.month-range-selection-sheet .dp__main .dp__instance_calendar .dp__overlay.dp--overlay-relative .dp__selection_grid_header .dp--year-mode-picker .dp--year-select+.dp--arrow-btn-nav {
    justify-content: end;
}
</style>
