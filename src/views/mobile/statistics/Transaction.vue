<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title>
                <f7-link popover-open=".chart-data-type-popover-menu">
                    <span>{{ query.chartDataType | chartDataTypeName(allChartDataTypes) | localized }}</span>
                    <f7-icon size="14px" :f7="showChartDataTypePopover ? 'arrowtriangle_up_fill' : 'arrowtriangle_down_fill'"></f7-icon>
                </f7-link>
            </f7-nav-title>
        </f7-navbar>

        <f7-popover class="chart-data-type-popover-menu" :opened="showChartDataTypePopover"
                    @popover:open="showChartDataTypePopover = true" @popover:close="showChartDataTypePopover = false">
            <f7-list>
                <f7-list-item
                    v-for="dataType in allChartDataTypes" :key="dataType.type"
                    :title="$t(dataType.name)" @click="setChartDataType(dataType.type)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.chartDataType === dataType.type"></f7-icon>
                </f7-list-item>
            </f7-list>
        </f7-popover>

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
            loading: true,
            showChartDataTypePopover: false
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
        allChartDataTypes() {
            return [
                {
                    type: this.$constants.statistics.allChartDataTypes.ExpenseByAccount,
                    name: 'Expense By Account'
                },
                {
                    type: this.$constants.statistics.allChartDataTypes.ExpenseByPrimaryCategory,
                    name: 'Expense By Primary Category'
                },
                {
                    type: this.$constants.statistics.allChartDataTypes.ExpenseBySecondaryCategory,
                    name: 'Expense By Secondary Category'
                },
                {
                    type: this.$constants.statistics.allChartDataTypes.IncomeByAccount,
                    name: 'Income By Account'
                },
                {
                    type: this.$constants.statistics.allChartDataTypes.IncomeByPrimaryCategory,
                    name: 'Income By Primary Category'
                },
                {
                    type: this.$constants.statistics.allChartDataTypes.IncomeBySecondaryCategory,
                    name: 'Income By Secondary Category'
                },
            ];
        },
        statisticsData() {
            const self = this;
            const combinedData = {};

            let allAmount = 0;

            for (let i = 0; i < self.$store.state.transactionStatistics.items.length; i++) {
                const item = self.$store.state.transactionStatistics.items[i];

                if (!item.account || !item.primaryCategory || !item.category) {
                    continue;
                }

                if (self.query.chartDataType === self.$constants.statistics.allChartDataTypes.ExpenseByAccount ||
                    self.query.chartDataType === self.$constants.statistics.allChartDataTypes.ExpenseByPrimaryCategory ||
                    self.query.chartDataType === self.$constants.statistics.allChartDataTypes.ExpenseBySecondaryCategory) {
                    if (item.category.type !== self.$constants.category.allCategoryTypes.Expense) {
                        continue;
                    }
                } else if (self.query.chartDataType === self.$constants.statistics.allChartDataTypes.IncomeByAccount ||
                    self.query.chartDataType === self.$constants.statistics.allChartDataTypes.IncomeByPrimaryCategory ||
                    self.query.chartDataType === self.$constants.statistics.allChartDataTypes.IncomeBySecondaryCategory) {
                    if (item.category.type !== self.$constants.category.allCategoryTypes.Income) {
                        continue;
                    }
                } else {
                    continue;
                }

                if (self.query.chartDataType === self.$constants.statistics.allChartDataTypes.ExpenseByAccount ||
                    self.query.chartDataType === self.$constants.statistics.allChartDataTypes.IncomeByAccount) {
                    if (self.$utilities.isNumber(item.amountInDefaultCurrency)) {
                        let data = combinedData[item.account.id];

                        if (data) {
                            data.totalAmount += item.amountInDefaultCurrency;
                        } else {
                            data = {
                                name: item.account.name,
                                icon: item.account.icon || self.$constants.icons.defaultAccountIcon.icon,
                                color: item.account.color || self.$constants.colors.defaultAccountColor,
                                totalAmount: item.amountInDefaultCurrency
                            }
                        }

                        allAmount += item.amountInDefaultCurrency;
                        combinedData[item.account.id] = data;
                    }
                } else if (self.query.chartDataType === self.$constants.statistics.allChartDataTypes.ExpenseByPrimaryCategory ||
                    self.query.chartDataType === self.$constants.statistics.allChartDataTypes.IncomeByPrimaryCategory) {
                    if (self.$utilities.isNumber(item.amountInDefaultCurrency)) {
                        let data = combinedData[item.primaryCategory.id];

                        if (data) {
                            data.totalAmount += item.amountInDefaultCurrency;
                        } else {
                            data = {
                                name: item.primaryCategory.name,
                                icon: item.account.icon || self.$constants.icons.defaultCategoryIcon.icon,
                                color: item.primaryCategory.color || self.$constants.colors.defaultCategoryColor,
                                totalAmount: item.amountInDefaultCurrency
                            }
                        }

                        allAmount += item.amountInDefaultCurrency;
                        combinedData[item.primaryCategory.id] = data;
                    }
                } else if (self.query.chartDataType === self.$constants.statistics.allChartDataTypes.ExpenseBySecondaryCategory ||
                    self.query.chartDataType === self.$constants.statistics.allChartDataTypes.IncomeBySecondaryCategory) {
                    if (self.$utilities.isNumber(item.amountInDefaultCurrency)) {
                        let data = combinedData[item.category.id];

                        if (data) {
                            data.totalAmount += item.amountInDefaultCurrency;
                        } else {
                            data = {
                                name: item.category.name,
                                icon: item.account.icon || self.$constants.icons.defaultCategoryIcon.icon,
                                color: item.category.color || self.$constants.colors.defaultCategoryColor,
                                totalAmount: item.amountInDefaultCurrency
                            }
                        }

                        allAmount += item.amountInDefaultCurrency;
                        combinedData[item.category.id] = data;
                    }
                }
            }

            const allStatisticsData = [];

            for (let id in combinedData) {
                if (!Object.prototype.hasOwnProperty.call(combinedData, id)) {
                    continue;
                }

                const data = combinedData[id];
                data.percent = data.totalAmount * 100 / allAmount;

                allStatisticsData.push(data);
            }

            allStatisticsData.sort(function (data1, data2) {
                return data1.totalAmount - data2.totalAmount;
            });

            return allStatisticsData;
        },
        chartData() {
            const self = this;

            if (!self.$store.state.transactionStatistics ||
                !self.$store.state.transactionStatistics.items ||
                !self.$store.state.transactionStatistics.items.length) {
                return self.skeletonChart();
            }

            const allData = [];

            for (let i = 0; i < this.statisticsData.length; i++) {
                const data = this.statisticsData[i];

                allData.push({
                    name: data.name,
                    value: data.totalAmount / 100,
                    itemStyle: {
                        color: `#${data.color}`
                    }
                });
            }

            const chartData =  {
                label: {
                    show: true,
                    overflow: 'truncate',
                    align: 'left',
                    formatter: params => {
                        return `${params.name} ${self.$options.filters.currency(params.value * 100, self.defaultCurrency)}`;
                    }
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'shadow'
                    }
                },
                series: [
                    {
                        data: allData,
                    }
                ],
                animation: false
            };

            if (this.query.chartType === this.$constants.statistics.allChartTypes.Bar) {
                return this.$utilities.copyObjectTo({
                    grid: {
                        left: 30,
                        top: 30,
                        right: 50,
                        bottom: 50
                    },
                    xAxis: {
                        type: 'value'
                    },
                    yAxis: {
                        type: 'category'
                    },
                    tooltip: {
                        trigger: 'axis',
                        axisPointer: {
                            type: 'shadow'
                        }
                    },
                    series: [{
                        type: 'bar'
                    }]
                }, chartData);
            } else {
                return this.$utilities.copyObjectTo({
                    tooltip: {
                        trigger: 'item'
                    },
                    series: [{
                        type: 'pie',
                        label: {
                            position: 'inside'
                        },
                        emphasis: {
                            itemStyle: {
                                shadowBlur: 10,
                                shadowOffsetX: 0,
                                shadowColor: 'rgba(0, 0, 0, 0.5)'
                            }
                        }
                    }]
                }, chartData);
            }
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
        },
        setChartDataType(chartDataType) {
            this.$store.dispatch('updateTransactionStatisticsFilter', {
                chartDataType: chartDataType
            });
            this.showChartDataTypePopover = false;
        },
        skeletonChart() {
            const skeletonChartData = {
                series: [{
                    data: [{
                        value: 20,
                        itemStyle: {
                            color: '#7c7c7f'
                        }
                    },{
                        value: 20,
                        itemStyle: {
                            color: '#a5a5aa'
                        }
                    },{
                        value: 60,
                        itemStyle: {
                            color: '#c5c5c9'
                        }
                    }]
                }],
                animation: false
            };

            if (this.query.chartType === this.$constants.statistics.allChartTypes.Bar) {
                return this.$utilities.copyObjectTo({
                    grid: {
                        left: 30,
                        top: 30,
                        right: 50,
                        bottom: 50
                    },
                    xAxis: {
                        type: 'value',
                        axisLabel: {
                            show: false
                        },
                        splitLine: {
                            show: false
                        }
                    },
                    yAxis: {
                        type: 'category'
                    },
                    series: [{
                        type: 'bar'
                    }]
                }, skeletonChartData);
            } else {
                return this.$utilities.copyObjectTo({
                    series: [{
                        type: 'pie',
                        label: {
                            position: 'inside'
                        }
                    }]
                }, skeletonChartData);
            }
        }
    },
    filters: {
        chartDataTypeName(dataType, allChartDataTypes) {
            for (let i = 0; i < allChartDataTypes.length; i++) {
                if (allChartDataTypes[i].type === dataType) {
                    return allChartDataTypes[i].name;
                }
            }

            return 'Statistics';
        },
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
