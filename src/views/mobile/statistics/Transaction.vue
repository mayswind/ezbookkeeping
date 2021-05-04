<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title>
                <f7-link popover-open=".chart-data-type-popover-menu">
                    <span>{{ query.chartDataType | optionName(allChartDataTypes, 'type', 'name', 'Statistics') | localized }}</span>
                    <f7-icon size="14px" :f7="showChartDataTypePopover ? 'arrowtriangle_up_fill' : 'arrowtriangle_down_fill'"></f7-icon>
                </f7-link>
            </f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-popover class="chart-data-type-popover-menu" :opened="showChartDataTypePopover"
                    @popover:open="showChartDataTypePopover = true" @popover:close="showChartDataTypePopover = false">
            <f7-list>
                <f7-list-item
                    v-for="dataType in allChartDataTypes" :key="dataType.type"
                    :title="$t(dataType.name)" @click="setChartDataType(dataType.type)">
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.chartDataType === dataType.type"></f7-icon>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-card v-if="query.chartType === $constants.statistics.allChartTypes.Pie">
            <f7-card-content class="pie-chart-container" :padding="false">
                <pie-chart
                    :items="[{value: 60, color: '7c7c7f'}, {value: 20, color: 'a5a5aa'}, {value: 20, color: 'c5c5c9'}]"
                    :skeleton="true"
                    :show-center-text="true"
                    :show-selected-item-info="true"
                    class="statistics-pie-chart"
                    value-field="value"
                    color-field="color"
                    center-text-background="#cccccc"
                    v-if="loading"
                ></pie-chart>
                <pie-chart
                    :items="statisticsData.items"
                    :min-valid-percent="0.0001"
                    :show-value="showAmountInChart"
                    :show-center-text="true"
                    :show-selected-item-info="true"
                    :enable-click-item="true"
                    :default-currency="defaultCurrency"
                    class="statistics-pie-chart"
                    name-field="name"
                    value-field="totalAmount"
                    percent-field="percent"
                    currency-field="currency"
                    hidden-field="hidden"
                    v-else-if="!loading"
                    @click="clickPieChartItem"
                >
                    <text class="statistics-pie-chart-total-amount-title" v-if="statisticsData.items && statisticsData.items.length">
                        {{ query.chartDataType | totalAmountName(allChartDataTypes) | localized }}
                    </text>
                    <text class="statistics-pie-chart-total-amount-value" v-if="statisticsData.items && statisticsData.items.length">
                        {{ statisticsData.totalAmount | currency(defaultCurrency) | finalAmount(showAccountBalance, query.chartDataType, allChartDataTypes) | textLimit(16) }}
                    </text>
                    <text class="statistics-pie-chart-total-no-data" cy="50%" v-if="!statisticsData.items || !statisticsData.items.length">
                        {{ $t('No data') }}
                    </text>
                </pie-chart>
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

                        <div slot="title" class="statistics-list-item-text">
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

                        <div slot="title" class="statistics-list-item-text">
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

                        <div slot="title" class="statistics-list-item-text">
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
                        <div slot="header">
                            {{ query.chartDataType | totalAmountName(allChartDataTypes) | localized }}
                        </div>
                        <div slot="title"
                             :class="{ 'statistics-list-item-overview-amount': true, 'text-color-teal': query.chartDataType === allChartDataTypes.ExpenseByAccount.type || query.chartDataType === allChartDataTypes.ExpenseByPrimaryCategory.type || query.chartDataType === allChartDataTypes.ExpenseBySecondaryCategory.type, 'text-color-red': query.chartDataType === allChartDataTypes.IncomeByAccount.type || query.chartDataType === allChartDataTypes.IncomeByPrimaryCategory.type || query.chartDataType === allChartDataTypes.IncomeBySecondaryCategory.type }">
                            {{ statisticsData.totalAmount | currency(defaultCurrency) | finalAmount(showAccountBalance, query.chartDataType, allChartDataTypes) }}
                        </div>
                    </f7-list-item>
                    <f7-list-item v-for="(item, idx) in statisticsData.items" :key="idx"
                                  class="statistics-list-item"
                                  :link="item | itemLinkUrl(query, allChartDataTypes)"
                                  v-show="!item.hidden"
                    >
                        <div slot="media" class="display-flex no-padding-horizontal">
                            <div class="display-flex align-items-center statistics-icon">
                                <f7-icon v-if="item.icon"
                                         :icon="item.icon | icon(item.type)"
                                         :style="item.color | iconStyle(item.type, 'var(--category-icon-color)')">
                                </f7-icon>
                                <f7-icon v-else-if="!item.icon"
                                         f7="pencil_ellipsis_rectangle">
                                </f7-icon>
                            </div>
                        </div>

                        <div slot="title" class="statistics-list-item-text">
                            <span>{{ item.name }}</span>
                            <small class="statistics-percent" v-if="item.percent >= 0">{{ item.percent | percent(2, '&lt;0.01') }}</small>
                        </div>

                        <div slot="after">
                            <span>{{ item.totalAmount | currency(item.currency || defaultCurrency) | finalAmount(showAccountBalance, query.chartDataType, allChartDataTypes) }}</span>
                        </div>

                        <div slot="inner-end" class="statistics-item-end">
                            <div class="statistics-percent-line">
                                <f7-progressbar :progress="item.percent >= 0 ? item.percent : 0" :style="{ '--f7-progressbar-progress-color': (item.color ? '#' + item.color : '') } "></f7-progressbar>
                            </div>
                        </div>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-toolbar tabbar bottom class="toolbar-item-auto-size">
            <f7-link :class="{ 'disabled': query.dateType === allDateRanges.All.type || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type }" @click="shiftDateRange(query.startTime, query.endTime, -1)">
                <f7-icon f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type }" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.maxTime > 0 || query.minTime > 0 }">{{ dateRangeName(query) }}</span>
            </f7-link>
            <f7-link :class="{ 'disabled': query.dateType === allDateRanges.All.type || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type }" @click="shiftDateRange(query.startTime, query.endTime, 1)">
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
                    @popover:open="scrollPopoverToSelectedItem"
                    @popover:opened="showDatePopover = true" @popover:closed="showDatePopover = false">
            <f7-list>
                <f7-list-item v-for="dateRange in allDateRanges"
                              :key="dateRange.type"
                              :class="{ 'list-item-selected': query.dateType === dateRange.type }"
                              :title="dateRange.name | localized"
                              @click="setDateFilter(dateRange.type)">
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.dateType === dateRange.type"></f7-icon>
                    <div slot="footer"
                         v-if="dateRange.type === allDateRanges.Custom.type && query.dateType === allDateRanges.Custom.type && query.startTime && query.endTime">
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

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="filterAccounts">{{ $t('Filter Accounts') }}</f7-actions-button>
                <f7-actions-button @click="filterCategories">{{ $t('Filter Transaction Categories') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button @click="settings">{{ $t('Settings') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            loading: true,
            loadingError: null,
            sortBy: self.$settings.getStatisticsSortingType(),
            showAccountBalance: self.$settings.isShowAccountBalance(),
            showChartDataTypePopover: false,
            showDatePopover: false,
            showCustomDateRangeSheet: false,
            showMoreActionSheet: false
        };
    },
    computed: {
        defaultCurrency() {
            return this.$store.getters.currentUserDefaultCurrency;
        },
        firstDayOfWeek() {
            return this.$store.getters.currentUserFirstDayOfWeek;
        },
        query() {
            return this.$store.state.transactionStatisticsFilter;
        },
        allChartDataTypes() {
            return this.$constants.statistics.allChartDataTypes;
        },
        allDateRanges() {
            return this.$constants.datetime.allDateRanges;
        },
        statisticsData() {
            const self = this;
            let combinedData = {
                items: [],
                totalAmount: 0
            };

            if (self.query.chartDataType === self.allChartDataTypes.ExpenseByAccount.type ||
                self.query.chartDataType === self.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.ExpenseBySecondaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeByAccount.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeByPrimaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeBySecondaryCategory.type) {
                combinedData = this.$store.getters.statisticsItemsByTransactionStatisticsData;
            } else if (self.query.chartDataType === self.allChartDataTypes.AccountTotalAssets.type ||
                self.query.chartDataType === self.allChartDataTypes.AccountTotalLiabilities.type) {
                combinedData = this.$store.getters.statisticsItemsByAccountsData;
            }

            const allStatisticsItems = [];

            for (let id in combinedData.items) {
                if (!Object.prototype.hasOwnProperty.call(combinedData.items, id)) {
                    continue;
                }

                const data = combinedData.items[id];

                if (data.totalAmount > 0) {
                    data.percent = data.totalAmount * 100 / combinedData.totalNonNegativeAmount;
                } else {
                    data.percent = 0;
                }

                if (data.percent < 0) {
                    data.percent = 0;
                }

                allStatisticsItems.push(data);
            }

            if (this.sortBy === this.$constants.statistics.allSortingTypes.ByDisplayOrder) {
                allStatisticsItems.sort(function (data1, data2) {
                    for (let i = 0; i < Math.min(data1.displayOrders.length, data2.displayOrders.length); i++) {
                        if (data1.displayOrders[i] !== data2.displayOrders[i]) {
                            return data1.displayOrders[i] - data2.displayOrders[i]; // asc
                        }
                    }

                    return data1.name.localeCompare(data2.name, undefined, { // asc
                        numeric: true,
                        sensitivity: 'base'
                    });
                });
            } else if (this.sortBy === this.$constants.statistics.allSortingTypes.ByName) {
                allStatisticsItems.sort(function (data1, data2) {
                    return data1.name.localeCompare(data2.name, undefined, { // asc
                        numeric: true,
                        sensitivity: 'base'
                    });
                });
            } else {
                allStatisticsItems.sort(function (data1, data2) {
                    if (data1.totalAmount !== data2.totalAmount) {
                        return data2.totalAmount - data1.totalAmount; // desc
                    }

                    return data1.name.localeCompare(data2.name, undefined, { // asc
                        numeric: true,
                        sensitivity: 'base'
                    });
                });
            }

            return {
                totalAmount: combinedData.totalAmount,
                items: allStatisticsItems
            };
        },
        showAmountInChart() {
            if (!this.showAccountBalance
                && (this.query.chartDataType === this.allChartDataTypes.AccountTotalAssets.type || this.query.chartDataType === this.allChartDataTypes.AccountTotalLiabilities.type)) {
                return false;
            }

            return true;
        }
    },
    created() {
        const self = this;

        let defaultChartType = self.$settings.getStatisticsDefaultChartType();

        if (defaultChartType !== self.$constants.statistics.allChartTypes.Pie && defaultChartType !== self.$constants.statistics.allChartTypes.Bar) {
            defaultChartType = self.$constants.statistics.defaultChartType;
        }

        let defaultChartDataType = self.$settings.getStatisticsDefaultChartDataType();

        if (defaultChartDataType < self.allChartDataTypes.ExpenseByAccount.type || defaultChartDataType > self.allChartDataTypes.AccountTotalLiabilities.type) {
            defaultChartDataType = self.$constants.statistics.defaultChartDataType;
        }

        let defaultDateRange = self.$settings.getStatisticsDefaultDateRange();

        if (defaultDateRange < self.allDateRanges.All.type || defaultDateRange >= self.allDateRanges.Custom.type) {
            defaultDateRange = self.$constants.statistics.defaultDataRangeType;
        }

        const dateRange = self.$utilities.getDateRangeByDateType(defaultDateRange, self.firstDayOfWeek);

        self.$store.dispatch('initTransactionStatisticsFilter', {
            dateType: dateRange ? dateRange.dateType : undefined,
            startTime: dateRange ? dateRange.minTime : undefined,
            endTime: dateRange ? dateRange.maxTime : undefined,
            chartType: defaultChartType,
            chartDataType: defaultChartDataType,
            filterAccountIds: self.$settings.getStatisticsDefaultAccountFilter() || {},
            filterCategoryIds: self.$settings.getStatisticsDefaultTransactionCategoryFilter() || {},
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
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            if (this.sortBy !== this.$settings.getStatisticsSortingType()) {
                this.sortBy = this.$settings.getStatisticsSortingType();
            }

            if (this.$store.state.transactionStatisticsStateInvalid && !this.loading) {
                this.reload(null);
            }

            this.$routeBackOnError('loadingError');
        },
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
            if (dateType === this.allDateRanges.Custom.type) { // Custom
                this.showCustomDateRangeSheet = true;
                this.showDatePopover = false;
                return;
            } else if (this.query.dateType === dateType) {
                return;
            }

            const dateRange = this.$utilities.getDateRangeByDateType(dateType, this.firstDayOfWeek);

            if (!dateRange) {
                return;
            }

            this.$store.dispatch('updateTransactionStatisticsFilter', {
                dateType: dateRange.dateType,
                startTime: dateRange.minTime,
                endTime: dateRange.maxTime
            });

            this.showDatePopover = false;
            this.reload(null);
        },
        setCustomDateFilter(startTime, endTime) {
            if (!startTime || !endTime) {
                return;
            }

            this.$store.dispatch('updateTransactionStatisticsFilter', {
                dateType: this.allDateRanges.Custom.type,
                startTime: startTime,
                endTime: endTime
            });

            this.showCustomDateRangeSheet = false;

            this.reload(null);
        },
        shiftDateRange(startTime, endTime, scale) {
            if (this.query.dateType === this.allDateRanges.All.type) {
                return;
            }

            const newDateRange = this.$utilities.getShiftedDateRange(startTime, endTime, scale);
            let newDateType = this.allDateRanges.Custom.type;

            for (let dateRangeField in this.allDateRanges) {
                if (!Object.prototype.hasOwnProperty.call(this.allDateRanges, dateRangeField)) {
                    continue;
                }

                const dateRangeType = this.allDateRanges[dateRangeField];
                const dateRange = this.$utilities.getDateRangeByDateType(dateRangeType.type, this.firstDayOfWeek);

                if (dateRange && dateRange.minTime === newDateRange.minTime && dateRange.maxTime === newDateRange.maxTime) {
                    newDateType = dateRangeType.type;
                    break;
                }
            }

            this.$store.dispatch('updateTransactionStatisticsFilter', {
                dateType: newDateType,
                startTime: newDateRange.minTime,
                endTime: newDateRange.maxTime
            });

            this.reload(null);
        },
        dateRangeName(query) {
            if (query.chartDataType === this.allChartDataTypes.AccountTotalAssets.type ||
                query.chartDataType === this.allChartDataTypes.AccountTotalLiabilities.type) {
                return this.$t(this.allDateRanges.All.name);
            }

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

            if (this.$utilities.isDateRangeMatchFullYears(query.startTime, query.endTime)) {
                const displayStartTime = this.$utilities.formatUnixTime(query.startTime, this.$t('format.year.short'));
                const displayEndTime = this.$utilities.formatUnixTime(query.endTime, this.$t('format.year.short'));

                return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
            }

            if (this.$utilities.isDateRangeMatchFullMonths(query.startTime, query.endTime)) {
                const displayStartTime = this.$utilities.formatUnixTime(query.startTime, this.$t('format.yearMonth.short'));
                const displayEndTime = this.$utilities.formatUnixTime(query.endTime, this.$t('format.yearMonth.short'));

                return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
            }

            const startTimeYear = this.$utilities.getYear(this.$utilities.parseDateFromUnixTime(query.startTime));
            const endTimeYear = this.$utilities.getYear(this.$utilities.parseDateFromUnixTime(query.endTime));

            const displayStartTime = this.$utilities.formatUnixTime(query.startTime, this.$t('format.date.short'));
            const displayEndTime = this.$utilities.formatUnixTime(query.endTime, this.$t('format.date.short'));

            if (displayStartTime === displayEndTime) {
                return displayStartTime;
            } else if (startTimeYear === endTimeYear) {
                const displayShortEndTime = this.$utilities.formatUnixTime(query.endTime, this.$t('format.monthDay.short'));
                return `${displayStartTime} ~ ${displayShortEndTime}`;
            }

            return `${displayStartTime} ~ ${displayEndTime}`;
        },
        clickPieChartItem(item) {
            this.$f7router.navigate(this.$options.filters.itemLinkUrl(item, this.query, this.allChartDataTypes));
        },
        filterAccounts() {
            this.$f7router.navigate('/statistic/filter/account');
        },
        filterCategories() {
            this.$f7router.navigate('/statistic/filter/category');
        },
        settings() {
            this.$f7router.navigate('/statistic/settings');
        },
        scrollPopoverToSelectedItem(event) {
            if (!event || !event.$el || !event.$el.length) {
                return;
            }

            const container = event.$el.find('.popover-inner');
            const selectedItem = event.$el.find('li.list-item-selected');

            if (!container.length || !selectedItem.length) {
                return;
            }

            let targetPos = selectedItem.offset().top - container.offset().top - parseInt(container.css('padding-top'), 10)
                - (container.outerHeight() - selectedItem.outerHeight()) / 2;

            if (targetPos <= 0) {
                return;
            }

            container.scrollTop(targetPos);
        }
    },
    filters: {
        finalAmount(amount, isShowAccountBalance, dataType, allChartDataTypes) {
            if (!isShowAccountBalance && (dataType === allChartDataTypes.AccountTotalAssets.type || dataType === allChartDataTypes.AccountTotalLiabilities.type)) {
                return '***';
            }

            return amount;
        },
        totalAmountName(dataType, allChartDataTypes) {
            if (dataType === allChartDataTypes.IncomeByAccount.type || dataType === allChartDataTypes.IncomeByPrimaryCategory.type || dataType === allChartDataTypes.IncomeBySecondaryCategory.type) {
                return 'Total Income';
            } else if (dataType === allChartDataTypes.ExpenseByAccount.type || dataType === allChartDataTypes.ExpenseByPrimaryCategory.type || dataType === allChartDataTypes.ExpenseBySecondaryCategory.type) {
                return 'Total Expense';
            } else if (dataType === allChartDataTypes.AccountTotalAssets.type) {
                return 'Total Assets';
            } else if (dataType === allChartDataTypes.AccountTotalLiabilities.type) {
                return 'Total Liabilities';
            }

            return 'Total Amount';
        },
        itemLinkUrl(item, query, allChartDataTypes) {
            const querys = [];

            if (query.chartDataType === allChartDataTypes.IncomeByAccount.type || query.chartDataType === allChartDataTypes.IncomeByPrimaryCategory.type || query.chartDataType === allChartDataTypes.IncomeBySecondaryCategory.type) {
                querys.push('type=2');
            } else if (query.chartDataType === allChartDataTypes.ExpenseByAccount.type || query.chartDataType === allChartDataTypes.ExpenseByPrimaryCategory.type || query.chartDataType === allChartDataTypes.ExpenseBySecondaryCategory.type) {
                querys.push('type=3');
            }

            if (query.chartDataType === allChartDataTypes.IncomeByAccount.type || query.chartDataType === allChartDataTypes.ExpenseByAccount.type || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type) {
                querys.push('accountId=' + item.id);
            } else if (query.chartDataType === allChartDataTypes.IncomeByPrimaryCategory.type || query.chartDataType === allChartDataTypes.IncomeBySecondaryCategory.type || query.chartDataType === allChartDataTypes.ExpenseByPrimaryCategory.type || query.chartDataType === allChartDataTypes.ExpenseBySecondaryCategory.type) {
                querys.push('categoryId=' + item.id);
            }

            if (query.chartDataType !== allChartDataTypes.AccountTotalAssets.type && query.chartDataType !== allChartDataTypes.AccountTotalLiabilities.type) {
                querys.push('dateType=' + query.dateType);
                querys.push('minTime=' + query.startTime);
                querys.push('maxTime=' + query.endTime);
            }

            return '/transaction/list?' + querys.join('&');
        }
    }
};
</script>

<style>
.statistics-pie-chart .pie-chart-text-group {
    fill: #fff;
    text-anchor: middle;
}

.statistics-pie-chart-total-amount-title {
    -moz-transform: translateY(0.5em);
    -ms-transform: translateY(0.5em);
    -webkit-transform: translateY(0.5em);
    transform: translateY(0.5em);
}

.statistics-pie-chart-total-amount-value {
    -moz-transform: translateY(2em);
    -ms-transform: translateY(2em);
    -webkit-transform: translateY(2em);
    transform: translateY(2em);
}

.statistics-pie-chart-total-no-data {
    -moz-transform: translateY(1.5em);
    -ms-transform: translateY(1.5em);
    -webkit-transform: translateY(1.5em);
    transform: translateY(1.5em);
}

.statistics-list-item-overview {
    padding-top: 12px;
    margin-bottom: 6px;
}

.statistics-list-item-overview-amount {
    margin-top: 2px;
    font-size: 1.5em;
    overflow: hidden;
    text-overflow: ellipsis;
}

.statistics-list-item-text {
    overflow: hidden;
    text-overflow: ellipsis;
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
</style>
