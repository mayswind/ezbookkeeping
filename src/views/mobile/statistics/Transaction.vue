<template>
    <f7-page ptr @ptr:refresh="reload">
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

        <f7-card v-if="query.chartType === $constants.statistics.allChartTypes.Pie">
            <f7-card-content class="no-safe-areas chart-container" :padding="false">
                <v-chart :options="chartData" v-if="chartData" />
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="query.chartType === $constants.statistics.allChartTypes.Bar">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list class="statistics-list-item skeleton-text" v-if="loading">
                    <f7-list-item link="#">
                        <div slot="media" class="display-flex no-padding-horizontal">
                            <div class="display-flex align-items-center statistics-icon">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </div>
                        </div>

                        <div slot="title">
                            <span>Category Name 1</span>
                            <small class="statistics-percent">33.33</small>
                        </div>

                        <div slot="after">
                            <span>0.00 USD</span>
                        </div>

                        <div slot="inner-end" class="statistics-item-end">
                            <div class="statistics-percent-line">
                                <f7-progressbar></f7-progressbar>
                            </div>
                        </div>
                    </f7-list-item>
                    <f7-list-item link="#">
                        <div slot="media" class="display-flex no-padding-horizontal">
                            <div class="display-flex align-items-center statistics-icon">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </div>
                        </div>

                        <div slot="title">
                            <span>Category Name 2</span>
                            <small class="statistics-percent">33.33</small>
                        </div>

                        <div slot="after">
                            <span>0.00 USD</span>
                        </div>

                        <div slot="inner-end" class="statistics-item-end">
                            <div class="statistics-percent-line">
                                <f7-progressbar></f7-progressbar>
                            </div>
                        </div>
                    </f7-list-item>
                    <f7-list-item link="#">
                        <div slot="media" class="display-flex no-padding-horizontal">
                            <div class="display-flex align-items-center statistics-icon">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </div>
                        </div>

                        <div slot="title">
                            <span>Category Name 3</span>
                            <small class="statistics-percent">33.33</small>
                        </div>

                        <div slot="after">
                            <span>0.00 USD</span>
                        </div>

                        <div slot="inner-end" class="statistics-item-end">
                            <div class="statistics-percent-line">
                                <f7-progressbar></f7-progressbar>
                            </div>
                        </div>
                    </f7-list-item>
                </f7-list>
                <f7-list v-else-if="!loading && (!statisticsData || !statisticsData.items || !statisticsData.items.length)">
                    <f7-list-item :title="$t('No transaction data')"></f7-list-item>
                </f7-list>
                <f7-list v-else-if="!loading && statisticsData && statisticsData.items && statisticsData.items.length">
                    <f7-list-item class="statistics-list-item-overview">
                        <div slot="header" v-if="query.chartDataType === $constants.statistics.allChartDataTypes.ExpenseByAccount || query.chartDataType === $constants.statistics.allChartDataTypes.ExpenseByPrimaryCategory || query.chartDataType === $constants.statistics.allChartDataTypes.ExpenseBySecondaryCategory">{{ $t('Total Expense') }}</div>
                        <div slot="header" v-if="query.chartDataType === $constants.statistics.allChartDataTypes.IncomeByAccount || query.chartDataType === $constants.statistics.allChartDataTypes.IncomeByPrimaryCategory || query.chartDataType === $constants.statistics.allChartDataTypes.IncomeBySecondaryCategory">{{ $t('Total Income') }}</div>
                        <div slot="title"
                             :class="{ 'statistics-list-item-overview-amount': true, 'text-color-teal': query.chartDataType === $constants.statistics.allChartDataTypes.ExpenseByAccount || query.chartDataType === $constants.statistics.allChartDataTypes.ExpenseByPrimaryCategory || query.chartDataType === $constants.statistics.allChartDataTypes.ExpenseBySecondaryCategory, 'text-color-red': query.chartDataType === $constants.statistics.allChartDataTypes.IncomeByAccount || query.chartDataType === $constants.statistics.allChartDataTypes.IncomeByPrimaryCategory || query.chartDataType === $constants.statistics.allChartDataTypes.IncomeBySecondaryCategory }">
                            {{ statisticsData.totalAmount | currency(defaultCurrency) }}
                        </div>
                    </f7-list-item>
                    <f7-list-item v-for="(data, idx) in statisticsData.items" :key="idx"
                                  class="statistics-list-item"
                                  :link="data | itemLinkUrl(query, $constants.statistics.allChartDataTypes)">
                        <div slot="media" class="display-flex no-padding-horizontal">
                            <div class="display-flex align-items-center statistics-icon">
                                <f7-icon v-if="data.icon"
                                         :icon="data.icon | icon(data.type)"
                                         :style="data.color | iconStyle(data.type, 'var(--category-icon-color)')">
                                </f7-icon>
                                <f7-icon v-else-if="!data.icon"
                                         f7="pencil_ellipsis_rectangle">
                                </f7-icon>
                            </div>
                        </div>

                        <div slot="title">
                            <span>{{ data.name }}</span>
                            <small class="statistics-percent">{{ data.percent | percent(2, '&lt;0.01') }}</small>
                        </div>

                        <div slot="after">
                            <span>{{ data.totalAmount | currency(defaultCurrency) }}</span>
                        </div>

                        <div slot="inner-end" class="statistics-item-end">
                            <div class="statistics-percent-line">
                                <f7-progressbar :progress="data.percent" :style="{ '--f7-progressbar-progress-color': (data.color ? '#' + data.color : '') } "></f7-progressbar>
                            </div>
                        </div>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-toolbar tabbar bottom class="toolbar-item-auto-size">
            <f7-link @click="backwardDateRange(query.startTime, query.endTime)">
                <f7-icon f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.maxTime > 0 || query.minTime > 0 }">{{ dateRangeName(query) }}</span>
            </f7-link>
            <f7-link @click="forwardDateRange(query.startTime, query.endTime)">
                <f7-icon f7="arrow_right_square"></f7-icon>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" @click="setChartType($constants.statistics.allChartTypes.Pie)">
                <span :class="{ 'tabbar-item-changed': query.chartType === $constants.statistics.allChartTypes.Pie }">{{ $t('Pie Chart') }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" @click="setChartType($constants.statistics.allChartTypes.Bar)">
                <span :class="{ 'tabbar-item-changed': query.chartType === $constants.statistics.allChartTypes.Bar }">{{ $t('Bar Chart') }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-popover class="date-popover-menu" :opened="showDatePopover"
                    @popover:opened="showDatePopover = true" @popover:closed="showDatePopover = false">
            <f7-list>
                <f7-list-item v-for="dateRange in allDateRanges"
                              :key="dateRange.type"
                              :title="dateRange.name | localized"
                              @click="setDateFilter(dateRange.type)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === dateRange.type"></f7-icon>
                    <div slot="footer"
                         v-if="dateRange.type === $constants.datetime.allDateRanges.Custom.type && query.dateType === $constants.datetime.allDateRanges.Custom.type && query.startTime && query.endTime">
                        <span>{{ query.startTime | moment($t('format.datetime.long-without-second')) }}</span>
                        <span>&nbsp;-&nbsp;</span>
                        <br/>
                        <span>{{ query.endTime | moment($t('format.datetime.long-without-second')) }}</span>
                    </div>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <date-range-selection-sheet :title="$t('Custom Date Range')"
                                    :show.sync="showCustomDateRangeSheet"
                                    :min-time="query.startTime"
                                    :max-time="query.endTime"
                                    @dateRange:change="setCustomDateFilter">
        </date-range-selection-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            loading: true,
            showChartDataTypePopover: false,
            showDatePopover: false,
            showCustomDateRangeSheet: false
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
        allDateRanges() {
            return this.$constants.datetime.allDateRanges;
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
                                type: 'account',
                                id: item.account.id,
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
                                type: 'category',
                                id: item.primaryCategory.id,
                                icon: item.primaryCategory.icon || self.$constants.icons.defaultCategoryIcon.icon,
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
                                type: 'category',
                                id: item.category.id,
                                icon: item.category.icon || self.$constants.icons.defaultCategoryIcon.icon,
                                color: item.category.color || self.$constants.colors.defaultCategoryColor,
                                totalAmount: item.amountInDefaultCurrency
                            }
                        }

                        allAmount += item.amountInDefaultCurrency;
                        combinedData[item.category.id] = data;
                    }
                }
            }

            const allStatisticsItems = [];

            for (let id in combinedData) {
                if (!Object.prototype.hasOwnProperty.call(combinedData, id)) {
                    continue;
                }

                const data = combinedData[id];
                data.percent = data.totalAmount * 100 / allAmount;

                allStatisticsItems.push(data);
            }

            allStatisticsItems.sort(function (data1, data2) {
                return data2.totalAmount - data1.totalAmount;
            });

            return {
                totalAmount: allAmount,
                items: allStatisticsItems
            };
        },
        chartData() {
            const self = this;

            if (!self.$store.state.transactionStatistics ||
                !self.$store.state.transactionStatistics.items ||
                !self.$store.state.transactionStatistics.items.length) {
                return self.skeletonChart();
            }

            const allData = [];

            for (let i = 0; i < this.statisticsData.items.length; i++) {
                const data = this.statisticsData.items[i];

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

        const dateParam = self.$utilities.getDateRangeByDateType(self.query.dateType);

        self.$store.dispatch('initTransactionStatisticsFilter', {
            dateType: dateParam ? dateParam.dateType : undefined,
            startTime: dateParam ? dateParam.minTime : undefined,
            endTime: dateParam ? dateParam.maxTime : undefined,
        });

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
        setDateFilter(dateType) {
            if (dateType === this.$constants.datetime.allDateRanges.Custom.type) { // Custom
                this.showCustomDateRangeSheet = true;
                this.showDatePopover = false;
                return;
            } else if (this.query.dateType === dateType) {
                return;
            }

            const dateParam = this.$utilities.getDateRangeByDateType(dateType);

            if (!dateParam) {
                return;
            }

            this.$store.dispatch('updateTransactionStatisticsFilter', {
                dateType: dateParam.dateType,
                startTime: dateParam.minTime,
                endTime: dateParam.maxTime
            });

            this.showDatePopover = false;
            this.reload(null);
        },
        setCustomDateFilter(startTime, endTime) {
            if (!startTime || !endTime) {
                return;
            }

            this.$store.dispatch('updateTransactionStatisticsFilter', {
                dateType: this.$constants.datetime.allDateRanges.Custom.type,
                startTime: startTime,
                endTime: endTime
            });

            this.showCustomDateRangeSheet = false;

            this.reload(null);
        },
        backwardDateRange(startTime, endTime) {
            this.$store.dispatch('updateTransactionStatisticsFilter', {
                dateType: this.$constants.datetime.allDateRanges.Custom.type,
                startTime: startTime - (endTime - startTime),
                endTime: startTime - 1
            });

            this.reload(null);
        },
        forwardDateRange(startTime, endTime) {
            this.$store.dispatch('updateTransactionStatisticsFilter', {
                dateType: this.$constants.datetime.allDateRanges.Custom.type,
                startTime: endTime + 1,
                endTime: endTime + (endTime - startTime)
            });

            this.reload(null);
        },
        dateRangeName(query) {
            if (query.dateType === this.allDateRanges.All.type) {
                return this.$t(this.allDateRanges.All.name);
            }

            for (let dateRangeField in this.allDateRanges) {
                if (!Object.prototype.hasOwnProperty.call(this.allDateRanges, dateRangeField)) {
                    continue;
                }

                const dateRange = this.allDateRanges[dateRangeField];

                if (dateRange && dateRange.type !== this.allDateRanges.Custom.type && dateRange.type === query.dateType && dateRange.name) {
                    return this.$t(dateRange.name);
                }
            }

            const startTimeYear = this.$utilities.getYear(this.$utilities.parseDateFromUnixTime(query.startTime));
            const endTimeYear = this.$utilities.getYear(this.$utilities.parseDateFromUnixTime(query.endTime));

            const displayStartTime = this.$utilities.formatUnixTime(query.startTime, this.$t('format.date.short'));
            const displayEndTime = this.$utilities.formatUnixTime(query.endTime, this.$t(startTimeYear !== endTimeYear ? 'format.date.short' : 'format.date.shortMonthDay'));

            return `${displayStartTime} ~ ${displayEndTime}`;
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
        itemLinkUrl(item, query, allChartDataTypes) {
            const querys = [];

            if (query.chartDataType === allChartDataTypes.IncomeByAccount || query.chartDataType === allChartDataTypes.IncomeByPrimaryCategory || query.chartDataType === allChartDataTypes.IncomeBySecondaryCategory) {
                querys.push('type=2');
            } else if (query.chartDataType === allChartDataTypes.ExpenseByAccount || query.chartDataType === allChartDataTypes.ExpenseByPrimaryCategory || query.chartDataType === allChartDataTypes.ExpenseBySecondaryCategory) {
                querys.push('type=3');
            }

            if (query.chartDataType === allChartDataTypes.IncomeByAccount || query.chartDataType === allChartDataTypes.ExpenseByAccount) {
                querys.push('accountId=' + item.id);
            } else if (query.chartDataType === allChartDataTypes.IncomeByPrimaryCategory || query.chartDataType === allChartDataTypes.IncomeBySecondaryCategory || query.chartDataType === allChartDataTypes.ExpenseByPrimaryCategory || query.chartDataType === allChartDataTypes.ExpenseBySecondaryCategory) {
                querys.push('categoryId=' + item.id);
            }

            querys.push('dateType=' + query.dateType);
            querys.push('minTime=' + query.startTime);
            querys.push('maxTime=' + query.endTime);

            return '/transaction/list?' + querys.join('&');
        }
    }
};
</script>

<style>
.statistics-list-item-overview {
    padding-top: 12px;
    margin-bottom: 6px;
}

.statistics-list-item-overview-amount {
    margin-top: 2px;
    font-size: 1.5em;
}

.statistics-list-item .item-content {
    margin-top: 8px;
    margin-bottom: 12px;
}

.statistics-list-item-overview .item-inner:after, .statistics-list-item .item-inner:after {
    background-color: transparent;
}

.statistics-icon {
    margin-bottom: -2px;
}

.statistics-percent {
    font-size: 0.7em;
    opacity: 0.6;
    margin-left: 6px;
}

.statistics-item-end {
    position: absolute;
    bottom: 0;
    width: 100%;
}

.statistics-percent-line {
    margin-right: calc(var(--f7-list-chevron-icon-area) + var(--f7-list-item-padding-horizontal) + var(--f7-safe-area-right));
}

.statistics-percent-line .progressbar {
    height: 4px;
    --f7-progressbar-bg-color: #f8f8f8;
}

.theme-dark .statistics-percent-line .progressbar {
    --f7-progressbar-bg-color: #161616;
}

.chart-container {
    height: 400px;
}

.chart-container .echarts {
    width: 100%;
    height: 100%;
}
</style>
