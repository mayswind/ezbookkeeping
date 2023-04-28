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
                    <f7-link class="margin-left-half" @click="toggleShowAmountInHomePage()">
                        <f7-icon :f7="showAmountInHomePage ? 'eye_slash_fill' : 'eye_fill'" size="18px"></f7-icon>
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

        <f7-list strong inset dividers class="margin-top" :class="{ 'skeleton-text': loading }">
            <f7-list-item :link="'/transaction/list?dateType=' + $constants.datetime.allDateRanges.Today.type" chevron-center>
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
                    <div>
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

            <f7-list-item :link="'/transaction/list?dateType=' + $constants.datetime.allDateRanges.ThisWeek.type" chevron-center>
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
                    <div>
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

            <f7-list-item :link="'/transaction/list?dateType=' + $constants.datetime.allDateRanges.ThisMonth.type" chevron-center>
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
                    <div>
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

            <f7-list-item :link="'/transaction/list?dateType=' + $constants.datetime.allDateRanges.ThisYear.type" chevron-center>
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
                    <div>
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
export default {
    data() {
        const self = this;

        return {
            loading: true,
            todayFirstUnixTime: self.$utilities.getTodayFirstUnixTime(),
            todayLastUnixTime: self.$utilities.getTodayLastUnixTime(),
            showAmountInHomePage: self.$settings.isShowAmountInHomePage(),
            isEnableThousandsSeparator: self.$settings.isEnableThousandsSeparator(),
            currencyDisplayMode: self.$settings.getCurrencyDisplayMode()
        };
    },
    computed: {
        defaultCurrency() {
            return this.$store.getters.currentUserDefaultCurrency;
        },
        firstDayOfWeek() {
            return this.$store.getters.currentUserFirstDayOfWeek;
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
        displayDateRange() {
            const self = this;

            return {
                today: {
                    displayTime: self.$utilities.formatUnixTime(self.dateRange.today.startTime, self.$t('format.date.long')),
                },
                thisWeek: {
                    startTime: self.$utilities.formatUnixTime(self.dateRange.thisWeek.startTime, self.$t('format.monthDay.long')),
                    endTime: self.$utilities.formatUnixTime(self.dateRange.thisWeek.endTime, self.$t('format.monthDay.long'))
                },
                thisMonth: {
                    displayTime: self.$utilities.formatUnixTime(self.dateRange.thisMonth.startTime, 'MMMM'),
                    startTime: self.$utilities.formatUnixTime(self.dateRange.thisMonth.startTime, self.$t('format.monthDay.long')),
                    endTime: self.$utilities.formatUnixTime(self.dateRange.thisMonth.endTime, self.$t('format.monthDay.long'))
                },
                thisYear: {
                    displayTime: self.$utilities.formatUnixTime(self.dateRange.thisYear.startTime, self.$t('format.year.long'))
                }
            };
        },
        transactionOverview() {
            // make sure this computed property refers these property, so these property can trigger this computed property to update
            const isEnableThousandsSeparator = this.isEnableThousandsSeparator; // eslint-disable-line
            const currencyDisplayMode = this.currencyDisplayMode; // eslint-disable-line

            if (!this.$store.state.transactionOverview || !this.$store.state.transactionOverview.thisMonth) {
                return {
                    thisMonth: {
                        valid: false,
                        incomeAmount: this.getDisplayAmount(0, false),
                        expenseAmount: this.getDisplayAmount(0, false)
                    }
                };
            }

            const originalOverview = this.$store.state.transactionOverview;
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
        }
    },
    methods: {
        onPageAfterIn() {
            this.showAmountInHomePage = this.$settings.isShowAmountInHomePage();

            if (this.isEnableThousandsSeparator !== this.$settings.isEnableThousandsSeparator() || this.currencyDisplayMode !== this.$settings.getCurrencyDisplayMode()) {
                this.isEnableThousandsSeparator = this.$settings.isEnableThousandsSeparator();
                this.currencyDisplayMode = this.$settings.getCurrencyDisplayMode();
                this.$forceUpdate();
            }

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

.overview-transaction-footer {
    padding-top: 6px;
    font-size: 13px;
}

.overview-transaction-footer > span {
    margin-right: 4px;
}

.tabbar.main-tabbar .link i + span.tabbar-label {
    margin-top: 2px;
}

.tabbar.main-tabbar .link i.ebk-tarbar-big-icon {
    font-size: 42px;
    width: 42px;
    height: 42px;
    line-height: 42px;
}
</style>
