<template>
    <v-row class="match-height">
        <v-col cols="12" lg="4" md="12">
            <v-card :class="{ 'disabled': loadingOverview }">
                <template #title>
                    <div class="d-flex align-center">
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
                    </div>
                </template>

                <v-card-text>
                    <h5 class="text-2xl font-weight-medium text-primary">
                        {{ transactionOverview && transactionOverview.thisMonth ? getDisplayExpenseAmount(transactionOverview.thisMonth) : '-' }}
                        <v-btn density="compact" color="default" variant="text"
                               :icon="true" @click="showAmountInHomePage = !showAmountInHomePage">
                            <v-icon :icon="showAmountInHomePage ? icons.eyeSlash : icons.eye" size="20" />
                        </v-btn>
                    </h5>
                    <p>
                        <span class="mr-2">{{ $t('Monthly income') }}</span>
                        <span>{{ transactionOverview && transactionOverview.thisMonth ? getDisplayIncomeAmount(transactionOverview.thisMonth) : '-' }}</span>
                    </p>
                    <v-btn size="small" to="/transactions">{{ $t('View Details') }}</v-btn>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" lg="2" md="6">
            <income-expense-overview-card
                :disabled="loadingOverview" :icon="icons.calendarToday"
                :title="$t('Today')"
                :expense-amount="transactionOverview.today && transactionOverview.today.valid ? getDisplayExpenseAmount(transactionOverview.today) : '-'"
                :income-amount="transactionOverview.today && transactionOverview.today.valid ? getDisplayIncomeAmount(transactionOverview.today) : '-'"
                :datetime="displayDateRange.today.displayTime"
            >
                <template #menus>
                    <v-list-item :prepend-icon="icons.viewDetails" :to="'/transactions?dateType=' + allDateRanges.Today.type">
                        <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                    </v-list-item>
                </template>
            </income-expense-overview-card>
        </v-col>

        <v-col cols="12" lg="2" md="6">
            <income-expense-overview-card
                :disabled="loadingOverview" :icon="icons.calendarWeek"
                :title="$t('This Week')"
                :expense-amount="transactionOverview.thisWeek && transactionOverview.thisWeek.valid ? getDisplayExpenseAmount(transactionOverview.thisWeek) : '-'"
                :income-amount="transactionOverview.thisWeek && transactionOverview.thisWeek.valid ? getDisplayIncomeAmount(transactionOverview.thisWeek) : '-'"
                :datetime="displayDateRange.thisWeek.startTime + '-' + displayDateRange.thisWeek.endTime"
            >
                <template #menus>
                    <v-list-item :prepend-icon="icons.viewDetails" :to="'/transactions?dateType=' + allDateRanges.ThisWeek.type">
                        <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                    </v-list-item>
                </template>
            </income-expense-overview-card>
        </v-col>

        <v-col cols="12" lg="2" md="6">
            <income-expense-overview-card
                :disabled="loadingOverview" :icon="icons.calendarMonth"
                :title="$t('This Month')"
                :expense-amount="transactionOverview.thisMonth && transactionOverview.thisMonth.valid ? getDisplayExpenseAmount(transactionOverview.thisMonth) : '-'"
                :income-amount="transactionOverview.thisMonth && transactionOverview.thisMonth.valid ? getDisplayIncomeAmount(transactionOverview.thisMonth) : '-'"
                :datetime="displayDateRange.thisMonth.startTime + '-' + displayDateRange.thisMonth.endTime"
            >
                <template #menus>
                    <v-list-item :prepend-icon="icons.viewDetails" :to="'/transactions?dateType=' + allDateRanges.ThisMonth.type">
                        <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                    </v-list-item>
                </template>
            </income-expense-overview-card>
        </v-col>

        <v-col cols="12" lg="2" md="6">
            <income-expense-overview-card
                :disabled="loadingOverview" :icon="icons.calendarYear"
                :title="$t('This Year')"
                :expense-amount="transactionOverview.thisYear && transactionOverview.thisYear.valid ? getDisplayExpenseAmount(transactionOverview.thisYear) : '-'"
                :income-amount="transactionOverview.thisYear && transactionOverview.thisYear.valid ? getDisplayIncomeAmount(transactionOverview.thisYear) : '-'"
                :datetime="displayDateRange.thisYear.displayTime"
            >
                <template #menus>
                    <v-list-item :prepend-icon="icons.viewDetails" :to="'/transactions?dateType=' + allDateRanges.ThisYear.type">
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
import { formatUnixTime } from '@/lib/datetime.js';

import {
    mdiRefresh,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiCalendarTodayOutline,
    mdiCalendarWeekOutline,
    mdiCalendarMonthOutline,
    mdiLayersTripleOutline,
    mdiListBoxOutline,
    mdiDotsVertical
} from '@mdi/js';

export default {
    components: {
        IncomeExpenseOverviewCard
    },
    data() {
        return {
            loadingOverview: true,
            icons: {
                refresh: mdiRefresh,
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline,
                calendarToday: mdiCalendarTodayOutline,
                calendarWeek: mdiCalendarWeekOutline,
                calendarMonth: mdiCalendarMonthOutline,
                calendarYear: mdiLayersTripleOutline,
                viewDetails: mdiListBoxOutline,
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
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        allDateRanges() {
            return datetimeConstants.allDateRanges;
        },
        displayDateRange() {
            const self = this;

            return {
                today: {
                    displayTime: self.$locale.formatUnixTimeToLongDate(self.userStore, self.overviewStore.transactionDataRange.today.startTime),
                },
                thisWeek: {
                    startTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.overviewStore.transactionDataRange.thisWeek.startTime),
                    endTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.overviewStore.transactionDataRange.thisWeek.endTime)
                },
                thisMonth: {
                    displayTime: formatUnixTime(self.overviewStore.transactionDataRange.thisMonth.startTime, 'MMMM'),
                    startTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.overviewStore.transactionDataRange.thisMonth.startTime),
                    endTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.overviewStore.transactionDataRange.thisMonth.endTime)
                },
                thisYear: {
                    displayTime: self.$locale.formatUnixTimeToLongYear(self.userStore, self.overviewStore.transactionDataRange.thisYear.startTime)
                }
            };
        },
        transactionOverview() {
            return this.overviewStore.transactionOverview;
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
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.getDisplayCurrency(value, currencyCode, {
                currencyDisplayMode: this.settingsStore.appSettings.currencyDisplayMode,
                enableThousandsSeparator: this.settingsStore.appSettings.thousandsSeparator
            });
        },
        getDisplayAmount(amount, incomplete) {
            if (!this.showAmountInHomePage) {
                return this.getDisplayCurrency('***', this.defaultCurrency);
            }

            return this.getDisplayCurrency(amount, this.defaultCurrency) + (incomplete ? '+' : '');
        },
        getDisplayIncomeAmount(category) {
            return this.getDisplayAmount(category.incomeAmount, category.incompleteIncomeAmount);
        },
        getDisplayExpenseAmount(category) {
            return this.getDisplayAmount(category.expenseAmount, category.incompleteExpenseAmount);
        }
    }
}
</script>
