<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-title :title="$t('global.app.title')"></f7-nav-title>
        </f7-navbar>

        <f7-card class="home-summary-card" :class="{ 'skeleton-text': loading }">
            <f7-card-header class="display-block" style="padding-top: 120px;">
                <p class="no-margin">
                    <span class="card-header-content" v-if="loading">
                        <span class="home-summary-month">MM</span>
                        <span>·</span>
                        <small>Expense</small>
                    </span>
                    <span class="card-header-content" v-else-if="!loading">
                        <span class="home-summary-month">{{ displayDateRange.thisMonth.displayTime }}</span>
                        <span>·</span>
                        <small>{{ $t('Expense') }}</small>
                    </span>
                </p>
                <p class="no-margin">
                    <span class="month-expense" v-if="loading">0.00 USD</span>
                    <span class="month-expense" v-else-if="!loading">{{ transactionOverview.thisMonth.expenseAmount }}</span>
                    <f7-link class="margin-left-half" @click="showAmountInHomePage = !showAmountInHomePage">
                        <f7-icon class="ebk-hide-icon" :f7="showAmountInHomePage ? 'eye_slash_fill' : 'eye_fill'"></f7-icon>
                    </f7-link>
                </p>
                <p class="no-margin">
                    <small class="home-summary-misc" v-if="loading">Monthly income 0.00 USD</small>
                    <small class="home-summary-misc" v-else-if="!loading">
                        <span>{{ $t('Monthly income') }}</span>
                        <span>{{ transactionOverview.thisMonth.incomeAmount }}</span>
                    </small>
                </p>
            </f7-card-header>
        </f7-card>

        <f7-list strong inset dividers class="margin-top overview-transaction-list" :class="{ 'skeleton-text': loading }">
            <f7-list-item :link="'/transaction/list?dateType=' + allDateRanges.Today.type" chevron-center>
                <template #media>
                    <f7-icon f7="calendar_today"></f7-icon>
                </template>
                <template #title>
                    <div class="padding-top-half">
                        <span v-if="loading">Today</span>
                        <span v-else-if="!loading">{{ $t('Today') }}</span>
                    </div>
                </template>
                <template #footer>
                    <div class="overview-transaction-footer padding-bottom-half">
                        <span v-if="loading">MM/DD/YYYY</span>
                        <span v-else-if="!loading">{{ displayDateRange.today.displayTime }}</span>
                    </div>
                </template>
                <template #after>
                    <div class="overview-transaction-amount">
                        <div class="text-color-red text-align-right">
                            <small v-if="loading">0.00 USD</small>
                            <small v-else-if="!loading && transactionOverview.today && transactionOverview.today.valid">{{ transactionOverview.today.incomeAmount }}</small>
                        </div>
                        <div class="text-color-teal text-align-right">
                            <small v-if="loading">0.00 USD</small>
                            <small v-else-if="!loading && transactionOverview.today && transactionOverview.today.valid">{{ transactionOverview.today.expenseAmount }}</small>
                        </div>
                    </div>
                </template>
            </f7-list-item>

            <f7-list-item :link="'/transaction/list?dateType=' + allDateRanges.ThisWeek.type" chevron-center>
                <template #media>
                    <f7-icon f7="calendar"></f7-icon>
                </template>
                <template #title>
                    <div class="padding-top-half">
                        <span v-if="loading">This Week</span>
                        <span v-else-if="!loading">{{ $t('This Week') }}</span>
                    </div>
                </template>
                <template #footer>
                    <div class="overview-transaction-footer padding-bottom-half">
                        <span v-if="loading">MM/DD</span>
                        <span v-else-if="!loading">{{ displayDateRange.thisWeek.startTime }}</span>
                        <span>-</span>
                        <span v-if="loading">MM/DD</span>
                        <span v-else-if="!loading">{{ displayDateRange.thisWeek.endTime }}</span>
                    </div>
                </template>
                <template #after>
                    <div class="overview-transaction-amount">
                        <div class="text-color-red text-align-right">
                            <small v-if="loading">0.00 USD</small>
                            <small v-else-if="!loading && transactionOverview.thisWeek && transactionOverview.thisWeek.valid">{{ transactionOverview.thisWeek.incomeAmount }}</small>
                        </div>
                        <div class="text-color-teal text-align-right">
                            <small v-if="loading">0.00 USD</small>
                            <small v-else-if="!loading && transactionOverview.thisWeek && transactionOverview.thisWeek.valid">{{ transactionOverview.thisWeek.expenseAmount }}</small>
                        </div>
                    </div>
                </template>
            </f7-list-item>

            <f7-list-item :link="'/transaction/list?dateType=' + allDateRanges.ThisMonth.type" chevron-center>
                <template #media>
                    <f7-icon f7="calendar"></f7-icon>
                </template>
                <template #title>
                    <div class="padding-top-half">
                        <span v-if="loading">This Month</span>
                        <span v-else-if="!loading">{{ $t('This Month') }}</span>
                    </div>
                </template>
                <template #footer>
                    <div class="overview-transaction-footer padding-bottom-half">
                        <span v-if="loading">MM/DD</span>
                        <span v-else-if="!loading">{{ displayDateRange.thisMonth.startTime }}</span>
                        <span>-</span>
                        <span v-if="loading">MM/DD</span>
                        <span v-else-if="!loading">{{ displayDateRange.thisMonth.endTime }}</span>
                    </div>
                </template>
                <template #after>
                    <div class="overview-transaction-amount">
                        <div class="text-color-red text-align-right">
                            <small v-if="loading">0.00 USD</small>
                            <small v-else-if="!loading && transactionOverview.thisMonth && transactionOverview.thisMonth.valid">{{ transactionOverview.thisMonth.incomeAmount }}</small>
                        </div>
                        <div class="text-color-teal text-align-right">
                            <small v-if="loading">0.00 USD</small>
                            <small v-else-if="!loading && transactionOverview.thisMonth && transactionOverview.thisMonth.valid">{{ transactionOverview.thisMonth.expenseAmount }}</small>
                        </div>
                    </div>
                </template>
            </f7-list-item>

            <f7-list-item :link="'/transaction/list?dateType=' + allDateRanges.ThisYear.type" chevron-center>
                <template #media>
                    <f7-icon f7="square_stack_3d_up"></f7-icon>
                </template>
                <template #title>
                    <div class="padding-top-half">
                        <span v-if="loading">This Year</span>
                        <span v-else-if="!loading">{{ $t('This Year') }}</span>
                    </div>
                </template>
                <template #footer>
                    <div class="overview-transaction-footer padding-bottom-half">
                        <span v-if="loading">YYYY</span>
                        <span v-else-if="!loading">{{ displayDateRange.thisYear.displayTime }}</span>
                    </div>
                </template>
                <template #after>
                    <div class="overview-transaction-amount">
                        <div class="text-color-red text-align-right">
                            <small v-if="loading">0.00 USD</small>
                            <small v-else-if="!loading && transactionOverview.thisYear && transactionOverview.thisYear.valid">{{ transactionOverview.thisYear.incomeAmount }}</small>
                        </div>
                        <div class="text-color-teal text-align-right">
                            <small v-if="loading">0.00 USD</small>
                            <small v-else-if="!loading && transactionOverview.thisYear && transactionOverview.thisYear.valid">{{ transactionOverview.thisYear.expenseAmount }}</small>
                        </div>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-toolbar tabbar icons bottom class="main-tabbar">
            <f7-link class="link" href="/transaction/list">
                <f7-icon f7="square_list"></f7-icon>
                <span class="tabbar-label">{{ $t('Details') }}</span>
            </f7-link>
            <f7-link class="link" href="/account/list">
                <f7-icon f7="creditcard"></f7-icon>
                <span class="tabbar-label">{{ $t('Accounts') }}</span>
            </f7-link>
            <f7-link class="link" href="/transaction/add">
                <f7-icon f7="plus_square" class="ebk-tarbar-big-icon"></f7-icon>
            </f7-link>
            <f7-link class="link" href="/statistic/transaction">
                <f7-icon f7="chart_pie"></f7-icon>
                <span class="tabbar-label">{{ $t('Statistics') }}</span>
            </f7-link>
            <f7-link class="link" href="/settings">
                <f7-icon f7="gear_alt"></f7-icon>
                <span class="tabbar-label">{{ $t('Settings') }}</span>
            </f7-link>
        </f7-toolbar>
    </f7-page>
</template>

<script>
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

export default {
    data() {
        const self = this;

        return {
            loading: true,
            todayFirstUnixTime: getTodayFirstUnixTime(),
            todayLastUnixTime: getTodayLastUnixTime()
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
            // make sure this computed property refers these property, so these property can trigger this computed property to update
            const isEnableThousandsSeparator = this.isEnableThousandsSeparator; // eslint-disable-line
            const currencyDisplayMode = this.currencyDisplayMode; // eslint-disable-line

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
        const self = this;

        if (self.$user.isUserLogined() && self.$user.isUserUnlocked()) {
            self.loading = true;

            self.overviewStore.loadTransactionOverview({
                defaultCurrency: self.defaultCurrency,
                dateRange: self.dateRange,
                force: false
            }).then(() => {
                self.loading = false;
            }).catch(error => {
                self.loading = false;

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        }
    },
    methods: {
        onPageAfterIn() {
            let dateChanged = false;

            if (this.todayFirstUnixTime !== getTodayFirstUnixTime()) {
                dateChanged = true;

                this.todayFirstUnixTime = getTodayFirstUnixTime();
                this.todayLastUnixTime = getTodayLastUnixTime();
            }

            if ((dateChanged || this.overviewStore.transactionOverviewStateInvalid) && !this.loading) {
                this.reload(null);
            }
        },
        reload(done) {
            const self = this;
            const force = !!done;

            self.overviewStore.loadTransactionOverview({
                defaultCurrency: self.defaultCurrency,
                dateRange: self.dateRange,
                force: force
            }).then(() => {
                if (done) {
                    done();
                }

                if (force) {
                    self.$toast('Data has been updated');
                }
            }).catch(error => {
                if (done) {
                    done();
                }

                if (!error.processed) {
                    self.$toast(error.message || error);
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

<style>
.home-summary-card {
    background-color: var(--f7-color-yellow);
}

.home-summary-card .home-summary-month {
    font-size: 1.3em;
}

.home-summary-card .month-expense {
    font-size: 1.5em;
}

.home-summary-card .home-summary-misc {
    opacity: 0.6;
}

.home-summary-misc > span {
    margin-right: 4px;
}

.home-summary-misc > span:last-child {
    margin-right: 0;
}

.dark .home-summary-card {
    background-color: var(--f7-theme-color);
}

.dark .home-summary-card a {
    color: var(--f7-text-color);
    opacity: 0.6;
}

.overview-transaction-list .item-title > div {
    overflow: hidden;
    text-overflow: ellipsis;
}

.overview-transaction-list .item-after {
    max-width: 100%;
}

.overview-transaction-list .overview-transaction-footer {
    padding-top: 6px;
    font-size: var(--ebk-large-footer-font-size);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.overview-transaction-list .overview-transaction-footer > span {
    margin-right: 4px;
}

.overview-transaction-list .overview-transaction-amount {
    max-width: 100%;
}

.overview-transaction-list .overview-transaction-amount > div {
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}

.tabbar.main-tabbar .link i + span.tabbar-label {
    margin-top: var(--ebk-icon-text-margin);
}

.tabbar.main-tabbar .link i.ebk-tarbar-big-icon {
    font-size: var(--ebk-big-icon-button-size);
    width: var(--ebk-big-icon-button-size);
    height: var(--ebk-big-icon-button-size);
    line-height: var(--ebk-big-icon-button-size);
}
</style>
