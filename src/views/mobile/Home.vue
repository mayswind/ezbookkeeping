<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-title :title="$t('global.app.title')"></f7-nav-title>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item link="#" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar_today"></f7-icon>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>Today</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer">
                            <span>MM/DD/YYYY</span>
                        </div>
                         <div slot="after">
                            <div class="text-color-red">
                                <small>0.00 USD</small>
                            </div>
                            <div class="text-color-teal">
                                <small>0.00 USD</small>
                            </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item link="#" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar"></f7-icon>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>This week</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer">
                            <span>MM/DD - MM/DD</span>
                        </div>
                         <div slot="after">
                            <div class="text-color-red">
                                <small>0.00 USD</small>
                            </div>
                            <div class="text-color-teal">
                                <small>0.00 USD</small>
                            </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item link="#" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar"></f7-icon>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>This month</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer">
                            <span>MM/DD - MM/DD</span>
                        </div>
                         <div slot="after">
                            <div class="text-color-red">
                                <small>0.00 USD</small>
                            </div>
                            <div class="text-color-teal">
                                <small>0.00 USD</small>
                            </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item link="#" chevron-center>
                        <div slot="media">
                            <f7-icon f7="square_stack_3d_up"></f7-icon>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>This year</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer">
                            <span>YYYY</span>
                        </div>
                         <div slot="after">
                            <div class="text-color-red">
                                <small>0.00 USD</small>
                            </div>
                            <div class="text-color-teal">
                                <small>0.00 USD</small>
                            </div>
                        </div>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item link="/transaction/list?dateType=1" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar_today"></f7-icon>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>{{ $t('Today' )}}</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer">
                            <span>{{ dateRange.today.startTime | moment($t('format.date.long')) }}</span>
                        </div>
                         <div slot="after">
                             <div class="text-color-red">
                                 <small v-if="transactionOverview.today">{{ transactionOverview.today.incomeAmount | currency(defaultCurrency) }}</small>
                             </div>
                             <div class="text-color-teal">
                                 <small v-if="transactionOverview.today">{{ transactionOverview.today.expenseAmount | currency(defaultCurrency) }}</small>
                             </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item link="/transaction/list?dateType=5" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar"></f7-icon>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>{{ $t('This Week' )}}</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer">
                            <span>{{ dateRange.thisWeek.startTime | moment($t('format.date.monthDay')) }}</span>
                            <span>-</span>
                            <span>{{ dateRange.thisWeek.endTime | moment($t('format.date.monthDay')) }}</span>
                        </div>
                         <div slot="after">
                             <div class="text-color-red">
                                 <small v-if="transactionOverview.thisWeek">{{ transactionOverview.thisWeek.incomeAmount | currency(defaultCurrency) }}</small>
                             </div>
                             <div class="text-color-teal">
                                 <small v-if="transactionOverview.thisWeek">{{ transactionOverview.thisWeek.expenseAmount | currency(defaultCurrency) }}</small>
                             </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item link="/transaction/list?dateType=7" chevron-center>
                        <div slot="media">
                            <f7-icon f7="calendar"></f7-icon>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>{{ $t('This Month' )}}</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer">
                            <span>{{ dateRange.thisMonth.startTime | moment($t('format.date.monthDay')) }}</span>
                            <span>-</span>
                            <span>{{ dateRange.thisMonth.endTime | moment($t('format.date.monthDay')) }}</span>
                        </div>
                         <div slot="after">
                             <div class="text-color-red">
                                 <small v-if="transactionOverview.thisMonth">{{ transactionOverview.thisMonth.incomeAmount | currency(defaultCurrency) }}</small>
                             </div>
                             <div class="text-color-teal">
                                 <small v-if="transactionOverview.thisMonth">{{ transactionOverview.thisMonth.expenseAmount | currency(defaultCurrency) }}</small>
                             </div>
                        </div>
                    </f7-list-item>

                    <f7-list-item link="/transaction/list?dateType=9" chevron-center>
                        <div slot="media">
                            <f7-icon f7="square_stack_3d_up"></f7-icon>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>{{ $t('This Year' )}}</span>
                        </div>
                        <div slot="footer" class="overview-transaction-footer">
                            <span>{{ dateRange.thisYear.startTime | moment($t('format.date.year')) }}</span>
                        </div>
                         <div slot="after">
                            <div class="text-color-red">
                                <small v-if="transactionOverview.thisYear">{{ transactionOverview.thisYear.incomeAmount | currency(defaultCurrency) }}</small>
                            </div>
                            <div class="text-color-teal">
                                <small v-if="transactionOverview.thisYear">{{ transactionOverview.thisYear.expenseAmount | currency(defaultCurrency) }}</small>
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
            <f7-link href="/statistic/overview">
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
            dateRange: self.getCurrentDateRange(),
            loading: true
        };
    },
    computed: {
        transactionOverview() {
            return this.$store.state.transactionOverview;
        },
        defaultCurrency() {
            return this.$store.getters.currentUserDefaultCurrency || this.$t('default.currency');
        }
    },
    created() {
        const self = this;

        self.loading = true;

        self.$store.dispatch('loadTransactionOverview', {
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
            const newDateRange = this.getCurrentDateRange();

            if (newDateRange.today.startTime !== this.dateRange.today.startTime ||
                newDateRange.today.endTime !== this.dateRange.today.endTime) {
                this.dateRange = newDateRange;
            }

            if (this.$store.state.transactionOverviewStateInvalid && !this.loading) {
                this.reload(null);
            }
        },
        reload(done) {
            const self = this;

            self.$store.dispatch('loadTransactionOverview', {
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
        getCurrentDateRange() {
            const self = this;

            return {
                today: {
                    startTime: self.$utilities.getTodayFirstUnixTime(),
                    endTime: self.$utilities.getTodayLastUnixTime()
                },
                thisWeek: {
                    startTime: self.$utilities.getThisWeekFirstUnixTime(),
                    endTime: self.$utilities.getThisWeekLastUnixTime()
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
        }
    }
}
</script>

<style>
.home-overview-card {
    background-color: var(--f7-color-yellow);
    height: 300px;
}

.theme-dark .home-overview-card {
    background-color: var(--f7-theme-color);
}

.overview-transaction-footer {
    padding-top: 2px;
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
