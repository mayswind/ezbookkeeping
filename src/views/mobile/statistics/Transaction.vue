<template>
    <f7-page>
        <f7-navbar :title="$t('Statistics')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas chart-container" :padding="false">
                <v-chart :options="chartData" v-if="chartData" />
            </f7-card-content>
        </f7-card>

        <f7-toolbar tabbar bottom class="toolbar-item-auto-size">
            <f7-link>
                <f7-icon f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.maxTime > 0 || query.minTime > 0 }">{{ query | dateRange }}</span>
            </f7-link>
            <f7-link>
                <f7-icon f7="arrow_right_square"></f7-icon>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" @click="setChartType($constants.statistics.allChartTypes.Pie)">
                <span :class="{ 'tabbar-item-changed': query.chartType === $constants.statistics.allChartTypes.Pie }">{{ $t('Pie Chart') }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" @click="setChartType($constants.statistics.allChartTypes.Bar)">
                <span :class="{ 'tabbar-item-changed': query.chartType === $constants.statistics.allChartTypes.Bar }">{{ $t('Bar Chart') }}</span>
            </f7-link>
        </f7-toolbar>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            loading: true
        };
    },
    computed: {
        defaultCurrency() {
            if (this.query.accountId && this.query.accountId !== '0') {
                const account = this.allAccounts[this.query.accountId];

                if (account && account.currency && account.currency !== this.$constants.currency.parentAccountCurrencyPlaceholder) {
                    return account.currency;
                }
            }

            return this.$store.getters.currentUserDefaultCurrency || this.$t('default.currency');
        },
        query() {
            return this.$store.state.transactionStatisticsFilter;
        },
        chartData() {
            if (!this.$store.state.transactionStatistics ||
                !this.$store.state.transactionStatistics.items ||
                !this.$store.state.transactionStatistics.items.length) {
                return null;
            }

            const combinedData = {};
            const data = [];

            for (let i = 0; i < this.$store.state.transactionStatistics.items.length; i++) {
                const item = this.$store.state.transactionStatistics.items[i];

                if (!item.account || !item.category) {
                    continue;
                }

                if (item.category.type !== this.$constants.category.allCategoryTypes.Expense) {
                    continue;
                }

                if (this.query.chartLegendType === this.$constants.statistics.allChartLegendTypes.Account) {
                    if (this.$utilities.isNumber(item.amountInDefaultCurrency)) {
                        let totalAmount = combinedData[item.account.name];

                        if (totalAmount) {
                            totalAmount += totalAmount = item.amountInDefaultCurrency;
                        } else {
                            totalAmount = item.amountInDefaultCurrency;
                        }

                        combinedData[item.account.name] = totalAmount;
                    }
                } else if (this.query.chartLegendType === this.$constants.statistics.allChartLegendTypes.SecondaryCategory) {
                    if (this.$utilities.isNumber(item.amountInDefaultCurrency)) {
                        let totalAmount = combinedData[item.category.name];

                        if (totalAmount) {
                            totalAmount += totalAmount = item.amountInDefaultCurrency;
                        } else {
                            totalAmount = item.amountInDefaultCurrency;
                        }

                        combinedData[item.category.name] = totalAmount;
                    }
                }
            }

            let chartType = 'pie';

            if (this.query.chartType === this.$constants.statistics.allChartTypes.Bar) {
                chartType = 'bar';
            }

            for (let legendName in combinedData) {
                if (!Object.prototype.hasOwnProperty.call(combinedData, legendName)) {
                    continue;
                }

                data.push({
                    name: legendName,
                    value: combinedData[legendName] / 100
                });
            }

            return {
                series: [
                    {
                        type: chartType,
                        data: data,
                        label: {
                            position: 'inside'
                        },
                        animation: false,
                    }
                ]
            };
        }
    },
    created() {
        const self = this;
        const router = self.$f7router;

        if (self.query.startTime < 0 || self.query.endTime < 0) {
            const dateParam = self.$utilities.getDateRangeByDateType(self.$constants.datetime.allDateRanges.ThisMonth.type);

            self.$store.dispatch('updateTransactionStatisticsFilter', {
                startTime: dateParam.minTime,
                endTime: dateParam.maxTime
            });
        }

        Promise.all([
            self.$store.dispatch('loadAllAccounts', { force: false }),
            self.$store.dispatch('loadAllCategories', { force: false })
        ]).then(() => {
            return self.$store.dispatch('loadTransactionStatistics', {
                defaultCurrency: self.defaultCurrency
            });
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            self.loading = false;

            if (!error.processed) {
                self.$toast(error.message || error);
                router.back();
            }
        });
    },
    methods: {
        reload(done) {
            const self = this;

            self.$store.dispatch('loadTransactionStatistics', {
                defaultCurrency: self.defaultCurrency
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
        setChartType(chartType) {
            this.$store.dispatch('updateTransactionStatisticsFilter', {
                chartType: chartType
            });
        }
    },
    filters: {
        dateRange() {
            return 'Date';
        }
    }
};
</script>

<style>
.chart-container {
    height: 400px;
}

.chart-container .echarts {
    width: 100%;
    height: 100%;
}
</style>
