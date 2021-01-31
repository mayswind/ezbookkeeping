<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-title :title="$t('global.app.title')"></f7-nav-title>
        </f7-navbar>

        <f7-card class="home-summary-card" :class="{ 'skeleton-text': loading }">
            <f7-card-header class="display-block" style="padding-top: 120px;">
                <p class="no-margin">
                    <span :style="{ opacity: 0.6 }" v-if="loading">
                        <span class="home-summary-month">MM</span>
                        <span>·</span>
                        <small>Expense</small>
                    </span>
                    <span :style="{ opacity: 0.6 }" v-else-if="!loading">
                        <span class="home-summary-month">{{ dateRange.thisMonth.startTime | moment('MMMM') | monthNameLocalizedKey | localized }}</span>
                        <span>·</span>
                        <small>{{ $t('Expense') }}</small>
                    </span>
                </p>
                <p class="no-margin">
                    <span class="month-expense" v-if="loading">0.00 USD</span>
                    <span class="month-expense" v-else-if="!loading">{{ thisMonthAmount.expenseAmount | amount(thisMonthAmount.incompleteExpenseAmount, showAmountInHomePage) | currency(defaultCurrency) }}</span>
                    <f7-link class="margin-left-half" @click="toggleShowAmountInHomePage()">
                        <f7-icon :f7="showAmountInHomePage ? 'eye_slash_fill' : 'eye_fill'" size="18px"></f7-icon>
                    </f7-link>
                </p>
                <p class="no-margin">
                    <small class="home-summary-misc" v-if="loading">Income of this month 0.00 USD</small>
                    <small class="home-summary-misc" v-else-if="!loading">
                        <span>{{ $t('Income of this month') }}</span>
                        <span>{{ thisMonthAmount.incomeAmount | amount(thisMonthAmount.incompleteIncomeAmount, showAmountInHomePage) | currency(defaultCurrency) }}</span>
                    </small>
                </p>
            </f7-card-header>
        </f7-card>

        <f7-card :class="{ 'skeleton-text': loading }">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :link="'/transaction/list?dateType=' + $constants.datetime.allDateRanges.Today.type" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar_today"></f7-icon>
                        </div>
                        <div slot="title" class="padding-top-half">
                            <span v-if="loading">Today</span>
                            <span v-else-if="!loading">{{ $t('Today' )}}</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer padding-bottom-half">
                            <span v-if="loading">MM/DD/YYYY</span>
                            <span v-else-if="!loading">{{ dateRange.today.startTime | moment($t('format.date.long')) }}</span>
                        </div>
                         <div slot="after">
                             <div class="text-color-red">
                                 <small v-if="loading">0.00 USD</small>
                                 <small v-else-if="!loading && transactionOverview.today">{{ transactionOverview.today.incomeAmount | amount(transactionOverview.today.incompleteIncomeAmount, showAmountInHomePage) | currency(defaultCurrency) }}</small>
                             </div>
                             <div class="text-color-teal">
                                 <small v-if="loading">0.00 USD</small>
                                 <small v-else-if="!loading && transactionOverview.today">{{ transactionOverview.today.expenseAmount | amount(transactionOverview.today.incompleteExpenseAmount, showAmountInHomePage) | currency(defaultCurrency) }}</small>
                             </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item :link="'/transaction/list?dateType=' + $constants.datetime.allDateRanges.ThisWeek.type" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar"></f7-icon>
                        </div>
                        <div slot="title" class="padding-top-half">
                            <span v-if="loading">This Week</span>
                            <span v-else-if="!loading">{{ $t('This Week' )}}</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer padding-bottom-half">
                            <span v-if="loading">MM/DD</span>
                            <span v-else-if="!loading">{{ dateRange.thisWeek.startTime | moment($t('format.monthDay.long')) }}</span>
                            <span>-</span>
                            <span v-if="loading">MM/DD</span>
                            <span v-else-if="!loading">{{ dateRange.thisWeek.endTime | moment($t('format.monthDay.long')) }}</span>
                        </div>
                         <div slot="after">
                             <div class="text-color-red">
                                 <small v-if="loading">0.00 USD</small>
                                 <small v-else-if="!loading && transactionOverview.thisWeek">{{ transactionOverview.thisWeek.incomeAmount | amount(transactionOverview.thisWeek.incompleteIncomeAmount, showAmountInHomePage) | currency(defaultCurrency) }}</small>
                             </div>
                             <div class="text-color-teal">
                                 <small v-if="loading">0.00 USD</small>
                                 <small v-else-if="!loading && transactionOverview.thisWeek">{{ transactionOverview.thisWeek.expenseAmount | amount(transactionOverview.thisWeek.incompleteExpenseAmount, showAmountInHomePage) | currency(defaultCurrency) }}</small>
                             </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item :link="'/transaction/list?dateType=' + $constants.datetime.allDateRanges.ThisMonth.type" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar"></f7-icon>
                        </div>
                        <div slot="title" class="padding-top-half">
                            <span v-if="loading">This Month</span>
                            <span v-else-if="!loading">{{ $t('This Month' )}}</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer padding-bottom-half">
                            <span v-if="loading">MM/DD</span>
                            <span v-else-if="!loading">{{ dateRange.thisMonth.startTime | moment($t('format.monthDay.long')) }}</span>
                            <span>-</span>
                            <span v-if="loading">MM/DD</span>
                            <span v-else-if="!loading">{{ dateRange.thisMonth.endTime | moment($t('format.monthDay.long')) }}</span>
                        </div>
                         <div slot="after">
                             <div class="text-color-red">
                                 <small v-if="loading">0.00 USD</small>
                                 <small v-else-if="!loading && transactionOverview.thisMonth">{{ transactionOverview.thisMonth.incomeAmount | amount(transactionOverview.thisMonth.incompleteIncomeAmount, showAmountInHomePage) | currency(defaultCurrency) }}</small>
                             </div>
                             <div class="text-color-teal">
                                 <small v-if="loading">0.00 USD</small>
                                 <small v-else-if="!loading && transactionOverview.thisMonth">{{ transactionOverview.thisMonth.expenseAmount | amount(transactionOverview.thisMonth.incompleteExpenseAmount, showAmountInHomePage) | currency(defaultCurrency) }}</small>
                             </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item :link="'/transaction/list?dateType=' + $constants.datetime.allDateRanges.ThisYear.type" chevron-center>
                        <div slot="media">
                            <f7-icon f7="square_stack_3d_up"></f7-icon>
                        </div>
                        <div slot="title" class="padding-top-half">
                            <span v-if="loading">This Year</span>
                            <span v-else-if="!loading">{{ $t('This Year' )}}</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer padding-bottom-half">
                            <span v-if="loading">YYYY</span>
                            <span v-else-if="!loading">{{ dateRange.thisYear.startTime | moment($t('format.year.long')) }}</span>
                        </div>
                         <div slot="after">
                            <div class="text-color-red">
                                <small v-if="loading">0.00 USD</small>
                                <small v-else-if="!loading && transactionOverview.thisYear">{{ transactionOverview.thisYear.incomeAmount | amount(transactionOverview.thisYear.incompleteIncomeAmount, showAmountInHomePage) | currency(defaultCurrency) }}</small>
                            </div>
                            <div class="text-color-teal">
                                <small v-if="loading">0.00 USD</small>
                                <small v-else-if="!loading && transactionOverview.thisYear">{{ transactionOverview.thisYear.expenseAmount | amount(transactionOverview.thisYear.incompleteExpenseAmount, showAmountInHomePage) | currency(defaultCurrency) }}</small>
                            </div>
                        </div>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-toolbar tabbar labels bottom>
            <f7-link href="/transaction/list">
                <f7-icon f7="square_list"></f7-icon>
                <span class="tabbar-label">{{ $t('Details') }}</span>
            </f7-link>
            <f7-link href="/account/list">
                <f7-icon f7="creditcard"></f7-icon>
                <span class="tabbar-label">{{ $t('Accounts') }}</span>
            </f7-link>
            <f7-link href="/transaction/add">
                <f7-icon f7="plus_square" class="lab-tarbar-big-icon"></f7-icon>
            </f7-link>
            <f7-link href="/statistic/transaction">
                <f7-icon f7="chart_pie"></f7-icon>
                <span class="tabbar-label">{{ $t('Statistics') }}</span>
            </f7-link>
            <f7-link href="/settings">
                <f7-icon f7="gear_alt"></f7-icon>
                <span class="tabbar-label">{{ $t('Settings') }}</span>
            </f7-link>
        </f7-toolbar>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            loading: true,
            todayFirstUnixTime: self.$utilities.getTodayFirstUnixTime(),
            todayLastUnixTime: self.$utilities.getTodayLastUnixTime(),
            showAmountInHomePage: self.$settings.isShowAmountInHomePage()
        };
    },
    computed: {
        transactionOverview() {
            return this.$store.state.transactionOverview;
        },
        defaultCurrency() {
            return this.$store.getters.currentUserDefaultCurrency || this.$t('default.currency');
        },
        firstDayOfWeek() {
            if (this.$utilities.isNumber(this.$store.getters.currentUserFirstDayOfWeek)) {
                return this.$store.getters.currentUserFirstDayOfWeek;
            }

            if (this.$constants.datetime.allWeekDays[this.$t('default.firstDayOfWeek')]) {
                return this.$constants.datetime.allWeekDays[this.$t('default.firstDayOfWeek')].type;
            }

            return 0;
        },
        dateRange() {
            const self = this;

            return {
                today: {
                    startTime: self.todayFirstUnixTime,
                    endTime: self.todayLastUnixTime
                },
                thisWeek: {
                    startTime: self.$utilities.getThisWeekFirstUnixTime(self.firstDayOfWeek),
                    endTime: self.$utilities.getThisWeekLastUnixTime(self.firstDayOfWeek)
                },
                thisMonth: {
                    startTime: self.$utilities.getThisMonthFirstUnixTime(),
                    endTime: self.$utilities.getThisMonthLastUnixTime()
                },
                thisYear: {
                    startTime: self.$utilities.getThisYearFirstUnixTime(),
                    endTime: self.$utilities.getThisYearLastUnixTime()
                }
            };
        },
        thisMonthAmount() {
            if (!this.$store.state.transactionOverview || !this.$store.state.transactionOverview.thisMonth) {
                return {
                    incomeAmount : 0,
                    expenseAmount : 0,
                    incompleteIncomeAmount: false,
                    incompleteExpenseAmount : false
                };
            }

            return this.$store.state.transactionOverview.thisMonth;
        }
    },
    created() {
        const self = this;

        self.loading = true;

        self.$store.dispatch('loadTransactionOverview', {
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
    },
    methods: {
        onPageAfterIn() {
            this.showAmountInHomePage = this.$settings.isShowAmountInHomePage();

            let dateChanged = false;

            if (this.todayFirstUnixTime !== this.$utilities.getTodayFirstUnixTime()) {
                dateChanged = true;

                this.todayFirstUnixTime = this.$utilities.getTodayFirstUnixTime();
                this.todayLastUnixTime = this.$utilities.getTodayLastUnixTime();
            }

            if ((dateChanged || this.$store.state.transactionOverviewStateInvalid) && !this.loading) {
                this.reload(null);
            }
        },
        reload(done) {
            const self = this;

            self.$store.dispatch('loadTransactionOverview', {
                defaultCurrency: self.defaultCurrency,
                dateRange: self.dateRange,
                force: true
            }).then(() => {
                if (done) {
                    done();
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
        toggleShowAmountInHomePage() {
            this.showAmountInHomePage = !this.showAmountInHomePage;
            this.$settings.setShowAmountInHomePage(this.showAmountInHomePage);
        }
    },
    filters: {
        amount(amount, incomplete, showAmount) {
            if (!showAmount) {
                return '***';
            }

            return amount + (incomplete ? '+' : '');
        },
        monthNameLocalizedKey(monthName) {
            return `datetime.${monthName}.long`;
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

.theme-dark .home-summary-card {
    background-color: var(--f7-theme-color);
}

.theme-dark .home-summary-card a {
    color: var(--f7-text-color);
    opacity: 0.6;
}

.overview-transaction-footer {
    padding-top: 6px;
    font-size: 13px;
}

.overview-transaction-footer > span {
    margin-right: 4px;
}

.tabbar-labels i.lab-tarbar-big-icon {
    font-size: 42px;
    width: 42px;
    height: 42px;
    line-height: 42px;
}
</style>
