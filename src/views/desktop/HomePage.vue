<template>
    <v-row class="match-height">
        <v-col cols="12" lg="4" md="12">
            <v-card :class="{ 'disabled': loadingOverview }">
                <template #title>
                    <span class="text-2xl font-weight-bold">{{ displayDateRange.thisMonth.displayTime }}</span>
                    <span>·</span>
                    <small>{{ $t('Expense') }}</small>
                    <v-btn density="compact" color="default" variant="text"
                           class="ml-2" :icon="true"
                           v-if="!loadingOverview" @click="reload(true)">
                        <v-icon :icon="icons.refresh" size="24" />
                        <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                    </v-btn>
                    <v-progress-circular indeterminate size="24" class="ml-2" v-if="loadingOverview"></v-progress-circular>
                </template>

                <v-card-text>
                    <h5 class="text-2xl font-weight-medium text-primary">
                        {{ transactionOverview && transactionOverview.thisMonth ? transactionOverview.thisMonth.expenseAmount : '-' }}
                        <v-btn density="compact" color="default" variant="text"
                               :icon="true" @click="showAmountInHomePage = !showAmountInHomePage">
                            <v-icon :icon="showAmountInHomePage ? icons.eyeSlash : icons.eye" size="20" />
                        </v-btn>
                    </h5>
                    <p>
                        <span class="mr-2">{{ $t('Monthly income') }}</span>
                        <span>{{ transactionOverview && transactionOverview.thisMonth ? transactionOverview.thisMonth.incomeAmount : '-' }}</span>
                    </p>
                    <v-btn size="small" to="/transactions">{{ $t('View Details') }}</v-btn>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" lg="2" md="6">
            <income-expense-overview-card
                :disabled="loadingOverview" :icon="icons.calendarToday"
                :title="$t('Today')"
                :expense-amount="transactionOverview.today && transactionOverview.today.valid ? transactionOverview.today.expenseAmount : '-'"
                :income-amount="transactionOverview.today && transactionOverview.today.valid ? transactionOverview.today.incomeAmount : '-'"
                :datetime="displayDateRange.today.displayTime"
            >
                <template #menus>
                    <v-list-item :to="'/transactions?dateType=' + allDateRanges.Today.type">
                        <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                    </v-list-item>
                </template>
            </income-expense-overview-card>
        </v-col>

        <v-col cols="12" lg="2" md="6">
            <income-expense-overview-card
                :disabled="loadingOverview" :icon="icons.calendarWeek"
                :title="$t('This Week')"
                :expense-amount="transactionOverview.thisWeek && transactionOverview.thisWeek.valid ? transactionOverview.thisWeek.expenseAmount : '-'"
                :income-amount="transactionOverview.thisWeek && transactionOverview.thisWeek.valid ? transactionOverview.thisWeek.incomeAmount : '-'"
                :datetime="displayDateRange.thisWeek.startTime + '-' + displayDateRange.thisWeek.endTime"
            >
                <template #menus>
                    <v-list-item :to="'/transactions?dateType=' + allDateRanges.ThisWeek.type">
                        <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                    </v-list-item>
                </template>
            </income-expense-overview-card>
        </v-col>

        <v-col cols="12" lg="2" md="6">
            <income-expense-overview-card
                :disabled="loadingOverview" :icon="icons.calendarMonth"
                :title="$t('This Month')"
                :expense-amount="transactionOverview.thisMonth && transactionOverview.thisMonth.valid ? transactionOverview.thisMonth.expenseAmount : '-'"
                :income-amount="transactionOverview.thisMonth && transactionOverview.thisMonth.valid ? transactionOverview.thisMonth.incomeAmount : '-'"
                :datetime="displayDateRange.thisMonth.startTime + '-' + displayDateRange.thisMonth.endTime"
            >
                <template #menus>
                    <v-list-item :to="'/transactions?dateType=' + allDateRanges.ThisMonth.type">
                        <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                    </v-list-item>
                </template>
            </income-expense-overview-card>
        </v-col>

        <v-col cols="12" lg="2" md="6">
            <income-expense-overview-card
                :disabled="loadingOverview" :icon="icons.calendarYear"
                :title="$t('This Year')"
                :expense-amount="transactionOverview.thisYear && transactionOverview.thisYear.valid ? transactionOverview.thisYear.expenseAmount : '-'"
                :income-amount="transactionOverview.thisYear && transactionOverview.thisYear.valid ? transactionOverview.thisYear.incomeAmount : '-'"
                :datetime="displayDateRange.thisYear.displayTime"
            >
                <template #menus>
                    <v-list-item :to="'/transactions?dateType=' + allDateRanges.ThisYear.type">
                        <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                    </v-list-item>
                </template>
            </income-expense-overview-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script>
import IncomeExpenseOverviewCard from './overview/IncomeExpenseOverviewCard.vue';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useOverviewStore } from '@/stores/overview.js';

import datetimeConstants from '@/consts/datetime.js';
import {
    formatUnixTime,
    getTodayFirstUnixTime,
    getTodayLastUnixTime,
    getThisWeekFirstUnixTime,
    getThisWeekLastUnixTime,
    getThisMonthFirstUnixTime,
    getThisMonthLastUnixTime,
    getThisYearFirstUnixTime,
    getThisYearLastUnixTime
} from '@/lib/datetime.js';

import {
    mdiRefresh,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiCalendarTodayOutline,
    mdiCalendarWeekOutline,
    mdiCalendarMonthOutline,
    mdiLayersTripleOutline,
    mdiDotsVertical
} from '@mdi/js';

export default {
    components: {
        IncomeExpenseOverviewCard
    },
    data() {
        return {
            loadingOverview: true,
            todayFirstUnixTime: getTodayFirstUnixTime(),
            todayLastUnixTime: getTodayLastUnixTime(),
            icons: {
                refresh: mdiRefresh,
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline,
                calendarToday: mdiCalendarTodayOutline,
                calendarWeek: mdiCalendarWeekOutline,
                calendarMonth: mdiCalendarMonthOutline,
                calendarYear: mdiLayersTripleOutline,
                more: mdiDotsVertical
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useOverviewStore),
        showAmountInHomePage: {
            get: function() {
                return this.settingsStore.appSettings.showAmountInHomePage;
            },
            set: function(value) {
                this.settingsStore.setShowAmountInHomePage(value);
            }
        },
        isEnableThousandsSeparator() {
            return this.settingsStore.appSettings.thousandsSeparator;
        },
        currencyDisplayMode() {
            return this.settingsStore.appSettings.currencyDisplayMode;
        },
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        allDateRanges() {
            return datetimeConstants.allDateRanges;
        },
        dateRange() {
            const self = this;

            return {
                today: {
                    startTime: self.todayFirstUnixTime,
                    endTime: self.todayLastUnixTime
                },
                thisWeek: {
                    startTime: getThisWeekFirstUnixTime(self.firstDayOfWeek),
                    endTime: getThisWeekLastUnixTime(self.firstDayOfWeek)
                },
                thisMonth: {
                    startTime: getThisMonthFirstUnixTime(),
                    endTime: getThisMonthLastUnixTime()
                },
                thisYear: {
                    startTime: getThisYearFirstUnixTime(),
                    endTime: getThisYearLastUnixTime()
                }
            };
        },
        displayDateRange() {
            const self = this;

            return {
                today: {
                    displayTime: self.$locale.formatUnixTimeToLongDate(self.userStore, self.dateRange.today.startTime),
                },
                thisWeek: {
                    startTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.dateRange.thisWeek.startTime),
                    endTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.dateRange.thisWeek.endTime)
                },
                thisMonth: {
                    displayTime: formatUnixTime(self.dateRange.thisMonth.startTime, 'MMMM'),
                    startTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.dateRange.thisMonth.startTime),
                    endTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.dateRange.thisMonth.endTime)
                },
                thisYear: {
                    displayTime: self.$locale.formatUnixTimeToLongYear(self.userStore, self.dateRange.thisYear.startTime)
                }
            };
        },
        transactionOverview() {
            if (!this.overviewStore.transactionOverview || !this.overviewStore.transactionOverview.thisMonth) {
                return {
                    thisMonth: {
                        valid: false,
                        incomeAmount: this.getDisplayAmount(0, false),
                        expenseAmount: this.getDisplayAmount(0, false)
                    }
                };
            }

            const originalOverview = this.overviewStore.transactionOverview;
            const displayOverview = {};

            [ 'today', 'thisWeek', 'thisMonth', 'thisYear' ].forEach(key => {
                if (!originalOverview[key]) {
                    return;
                }

                const item = originalOverview[key];

                displayOverview[key] = {
                    valid: true,
                    incomeAmount: this.getDisplayAmount(item.incomeAmount, item.incompleteIncomeAmount),
                    expenseAmount: this.getDisplayAmount(item.expenseAmount, item.incompleteExpenseAmount)
                };
            });

            return displayOverview;
        }
    },
    created() {
        if (this.$user.isUserLogined() && this.$user.isUserUnlocked()) {
            this.reload(false);
        }
    },
    methods: {
        reload(force) {
            const self = this;

            self.loadingOverview = true;

            self.overviewStore.loadTransactionOverview({
                defaultCurrency: self.defaultCurrency,
                dateRange: self.dateRange,
                force: force
            }).then(() => {
                self.loadingOverview = false;

                if (force) {
                    self.$refs.snackbar.showMessage('Data has been updated');
                }
            }).catch(error => {
                self.loadingOverview = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        getDisplayAmount(amount, incomplete) {
            if (!this.showAmountInHomePage) {
                return this.$locale.getDisplayCurrency('***', this.defaultCurrency);
            }

            return this.$locale.getDisplayCurrency(amount, this.defaultCurrency) + (incomplete ? '+' : '');
        }
    }
}
</script>
