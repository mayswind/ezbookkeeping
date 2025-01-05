<template>
    <v-dialog class="month-range-selection-dialog" width="640" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ title }}</h4>
                </div>
            </template>
            <template #subtitle>
                <div class="text-body-1 text-center text-wrap mt-6">
                    <p v-if="hint">{{ hint }}</p>
                    <span v-if="beginDateTime && endDateTime">
                        <span>{{ beginDateTime }}</span>
                        <span> - </span>
                        <span>{{ endDateTime }}</span>
                    </span>
                    <slot></slot>
                </div>
            </template>
            <v-card-text class="mb-md-4 w-100 d-flex justify-center">
                <v-row class="match-height">
                    <v-col cols="12" md="6">
                        <vue-date-picker inline month-picker auto-apply
                                         month-name-format="long"
                                         :clearable="false"
                                         :dark="isDarkMode"
                                         :year-range="yearRange"
                                         :year-first="isYearFirst"
                                         v-model="startTime">
                            <template #month="{ text }">
                                {{ getMonthShortName(text) }}
                            </template>
                            <template #month-overlay-value="{ text }">
                                {{ getMonthShortName(text) }}
                            </template>
                        </vue-date-picker>
                    </v-col>
                    <v-col cols="12" md="6">
                        <vue-date-picker inline month-picker auto-apply
                                         month-name-format="long"
                                         :clearable="false"
                                         :dark="isDarkMode"
                                         :year-range="yearRange"
                                         :year-first="isYearFirst"
                                         v-model="endTime">
                            <template #month="{ text }">
                                {{ getMonthShortName(text) }}
                            </template>
                            <template #month-overlay-value="{ text }">
                                {{ getMonthShortName(text) }}
                            </template>
                        </vue-date-picker>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="!startTime || !endTime" @click="confirm">{{ $t('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ $t('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.ts';

import { ThemeType } from '@/core/theme.ts';
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
        'persistent',
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
            startTime: minDate,
            endTime: maxDate
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
            return this.globalTheme.global.name.value === ThemeType.Dark;
        },
        isYearFirst() {
            return this.$locale.isLongDateMonthAfterYear(this.userStore);
        },
        beginDateTime() {
            return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, getYearMonthFirstUnixTime(this.startTime));
        },
        endDateTime() {
            return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, getYearMonthLastUnixTime(this.endTime));
        }
    },
    watch: {
        'minTime': function (newValue) {
            if (newValue) {
                this.startTime = getYearMonthObjectFromString(newValue);
            }
        },
        'maxTime': function (newValue) {
            if (newValue) {
                this.endTime = getYearMonthObjectFromString(newValue);
            }
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
            if (!this.startTime || !this.endTime) {
                return;
            }

            if (this.startTime.year <= 0 || this.startTime.month < 0 || this.endTime.year <= 0 || this.endTime.month < 0) {
                this.$toast('Date is too early');
                return;
            }

            const minYearMonth = getYearMonthStringFromObject(this.startTime);
            const maxYearMonth = getYearMonthStringFromObject(this.endTime);

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
.month-range-selection-dialog .dp__preset_ranges {
    white-space: nowrap !important;
}

.month-range-selection-dialog .dp__overlay {
    width: 100% !important;
    height: 100% !important;
}
</style>
