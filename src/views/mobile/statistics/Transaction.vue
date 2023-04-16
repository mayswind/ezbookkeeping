<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title>
                <f7-link popover-open=".chart-data-type-popover-menu">
                    <span>{{ queryChartDataTypeName }}</span>
                    <f7-icon size="14px" :f7="showChartDataTypePopover ? 'arrowtriangle_up_fill' : 'arrowtriangle_down_fill'"></f7-icon>
                </f7-link>
            </f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-popover class="chart-data-type-popover-menu" :opened="showChartDataTypePopover"
                    @popover:open="showChartDataTypePopover = true" @popover:close="showChartDataTypePopover = false">
            <f7-list dividers>
                <f7-list-item
                    v-for="dataType in allChartDataTypes" :key="dataType.type"
                    :title="$t(dataType.name)" @click="setChartDataType(dataType.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.chartDataType === dataType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-card v-if="query.chartType === $constants.statistics.allChartTypes.Pie">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header full-line text-align-right">
                    <span style="margin-right: 4px;">{{ $t('Sort By') }}</span>
                    <f7-link href="#" popover-open=".sorting-type-popover-menu">{{ querySortingTypeName }}</f7-link>
                </div>
            </f7-card-header>
            <f7-card-content class="pie-chart-container" style="margin-top: -6px" :padding="false">
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
                        {{ totalAmountName }}
                    </text>
                    <text class="statistics-pie-chart-total-amount-value" v-if="statisticsData.items && statisticsData.items.length">
                        {{ getDisplayAmount(statisticsData.totalAmount, defaultCurrency, 16) }}
                    </text>
                    <text class="statistics-pie-chart-total-no-data" cy="50%" v-if="!statisticsData.items || !statisticsData.items.length">
                        {{ $t('No data') }}
                    </text>
                </pie-chart>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="query.chartType === $constants.statistics.allChartTypes.Bar">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header display-flex full-line justify-content-space-between">
                    <div>
                        {{ totalAmountName }}
                    </div>
                    <div class="align-self-flex-end">
                        <span style="margin-right: 4px;">{{ $t('Sort By') }}</span>
                        <f7-link href="#" popover-open=".sorting-type-popover-menu">{{ querySortingTypeName }}</f7-link>
                    </div>
                </div>
                <div class="display-flex full-line">
                    <div :class="{ 'statistics-list-item-overview-amount': true, 'text-color-teal': query.chartDataType === allChartDataTypes.ExpenseByAccount.type || query.chartDataType === allChartDataTypes.ExpenseByPrimaryCategory.type || query.chartDataType === allChartDataTypes.ExpenseBySecondaryCategory.type, 'text-color-red': query.chartDataType === allChartDataTypes.IncomeByAccount.type || query.chartDataType === allChartDataTypes.IncomeByPrimaryCategory.type || query.chartDataType === allChartDataTypes.IncomeBySecondaryCategory.type }">
                        <span v-if="!loading && statisticsData && statisticsData.items && statisticsData.items.length">
                            {{ getDisplayAmount(statisticsData.totalAmount, defaultCurrency) }}
                        </span>
                        <span :class="{ 'skeleton-text': loading }" v-else-if="loading || !statisticsData || !statisticsData.items || !statisticsData.items.length">
                            ***.**
                        </span>
                    </div>
                </div>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" style="margin-top: -14px" :padding="false">
                <f7-list dividers class="statistics-list-item skeleton-text" v-if="loading">
                    <f7-list-item link="#">
                        <template #media>
                            <div class="display-flex no-padding-horizontal">
                                <div class="display-flex align-items-center statistics-icon">
                                    <f7-icon f7="app_fill"></f7-icon>
                                </div>
                            </div>
                        </template>

                        <template #title>
                            <div class="statistics-list-item-text">
                                <span>Category Name 1</span>
                                <small class="statistics-percent">33.33</small>
                            </div>
                        </template>

                        <template #after>
                            <span>0.00 USD</span>
                        </template>

                        <template #inner-end>
                            <div class="statistics-item-end">
                                <div class="statistics-percent-line">
                                    <f7-progressbar></f7-progressbar>
                                </div>
                            </div>
                        </template>
                    </f7-list-item>
                    <f7-list-item link="#">
                        <template #media>
                            <div class="display-flex no-padding-horizontal">
                                <div class="display-flex align-items-center statistics-icon">
                                    <f7-icon f7="app_fill"></f7-icon>
                                </div>
                            </div>
                        </template>

                        <template #title>
                            <div class="statistics-list-item-text">
                                <span>Category Name 2</span>
                                <small class="statistics-percent">33.33</small>
                            </div>
                        </template>

                        <template #after>
                            <span>0.00 USD</span>
                        </template>

                        <template #inner-end>
                            <div class="statistics-item-end">
                                <div class="statistics-percent-line">
                                    <f7-progressbar></f7-progressbar>
                                </div>
                            </div>
                        </template>
                    </f7-list-item>
                    <f7-list-item link="#">
                        <template #media>
                            <div class="display-flex no-padding-horizontal">
                                <div class="display-flex align-items-center statistics-icon">
                                    <f7-icon f7="app_fill"></f7-icon>
                                </div>
                            </div>
                        </template>

                        <template #title>
                            <div class="statistics-list-item-text">
                                <span>Category Name 3</span>
                                <small class="statistics-percent">33.33</small>
                            </div>
                        </template>

                        <template #after>
                            <span>0.00 USD</span>
                        </template>

                        <template #inner-end>
                            <div class="statistics-item-end">
                                <div class="statistics-percent-line">
                                    <f7-progressbar></f7-progressbar>
                                </div>
                            </div>
                        </template>
                    </f7-list-item>
                </f7-list>
                <f7-list dividers v-else-if="!loading && (!statisticsData || !statisticsData.items || !statisticsData.items.length)">
                    <f7-list-item :title="$t('No transaction data')"></f7-list-item>
                </f7-list>
                <f7-list dividers v-else-if="!loading && statisticsData && statisticsData.items && statisticsData.items.length">
                    <f7-list-item v-for="(item, idx) in statisticsData.items" :key="idx"
                                  class="statistics-list-item"
                                  :link="getItemLinkUrl(item)"
                                  v-show="!item.hidden"
                    >
                        <template #media>
                            <div class="display-flex no-padding-horizontal">
                                <div class="display-flex align-items-center statistics-icon">
                                    <ItemIcon icon-type="category" :icon-id="item.icon" :color="item.color" v-if="item.icon"></ItemIcon>
                                    <f7-icon f7="pencil_ellipsis_rectangle" v-else-if="!item.icon"></f7-icon>
                                </div>
                            </div>
                        </template>

                        <template #title>
                            <div class="statistics-list-item-text">
                                <span>{{ item.name }}</span>
                                <small class="statistics-percent" v-if="item.percent >= 0">{{ $utilities.formatPercent(item.percent, 2, '&lt;0.01') }}</small>
                            </div>
                        </template>

                        <template #after>
                            <span>{{ getDisplayAmount(item.totalAmount, (item.currency || defaultCurrency)) }}</span>
                        </template>

                        <template #inner-end>
                            <div class="statistics-item-end">
                                <div class="statistics-percent-line">
                                    <f7-progressbar :progress="item.percent >= 0 ? item.percent : 0" :style="{ '--f7-progressbar-progress-color': (item.color ? '#' + item.color : '') } "></f7-progressbar>
                                </div>
                            </div>
                        </template>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-popover class="sorting-type-popover-menu" :opened="showSortingTypePopover"
                    @popover:open="scrollPopoverToSelectedItem"
                    @popover:opened="showSortingTypePopover = true" @popover:closed="showSortingTypePopover = false">
            <f7-list dividers>
                <f7-list-item v-for="sortingType in allSortingTypes"
                              :key="sortingType.type"
                              :class="{ 'list-item-selected': query.sortingType === sortingType.type }"
                              :title="$t(sortingType.name)"
                              @click="setSortingType(sortingType.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.sortingType === sortingType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

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
            <f7-list dividers>
                <f7-list-item v-for="dateRange in allDateRanges"
                              :key="dateRange.type"
                              :class="{ 'list-item-selected': query.dateType === dateRange.type }"
                              :title="$t(dateRange.name)"
                              @click="setDateFilter(dateRange.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.dateType === dateRange.type"></f7-icon>
                    </template>
                    <template #footer>
                        <div v-if="dateRange.type === allDateRanges.Custom.type && query.dateType === allDateRanges.Custom.type && query.startTime && query.endTime">
                            <span>{{ $utilities.formatUnixTime(query.startTime, $t('format.datetime.long-without-second')) }}</span>
                            <span>&nbsp;-&nbsp;</span>
                            <br/>
                            <span>{{ $utilities.formatUnixTime(query.endTime, $t('format.datetime.long-without-second')) }}</span>
                        </div>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <date-range-selection-sheet :title="$t('Custom Date Range')"
                                    :min-time="query.startTime"
                                    :max-time="query.endTime"
                                    v-model:show="showCustomDateRangeSheet"
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
    props: [
        'f7router'
    ],
    data() {
        const self = this;

        return {
            loading: true,
            loadingError: null,
            showAccountBalance: self.$settings.isShowAccountBalance(),
            showChartDataTypePopover: false,
            showSortingTypePopover: false,
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
        queryChartDataTypeName() {
            const queryChartDataTypeName = this.$utilities.getNameByKeyValue(this.allChartDataTypes, this.query.chartDataType, 'type', 'name', 'Statistics');
            return this.$t(queryChartDataTypeName);
        },
        querySortingTypeName() {
            const querySortingTypeName = this.$utilities.getNameByKeyValue(this.allSortingTypes, this.query.sortingType, 'type', 'name', 'System Default');
            return this.$t(querySortingTypeName);
        },
        allChartDataTypes() {
            return this.$constants.statistics.allChartDataTypes;
        },
        allSortingTypes() {
            return this.$constants.statistics.allSortingTypes;
        },
        allDateRanges() {
            return this.$constants.datetime.allDateRanges;
        },
        totalAmountName() {
            if (this.query.chartDataType === this.allChartDataTypes.IncomeByAccount.type
                || this.query.chartDataType === this.allChartDataTypes.IncomeByPrimaryCategory.type
                || this.query.chartDataType === this.allChartDataTypes.IncomeBySecondaryCategory.type) {
                return this.$t('Total Income');
            } else if (this.query.chartDataType === this.allChartDataTypes.ExpenseByAccount.type
                || this.query.chartDataType === this.allChartDataTypes.ExpenseByPrimaryCategory.type
                || this.query.chartDataType === this.allChartDataTypes.ExpenseBySecondaryCategory.type) {
                return this.$t('Total Expense');
            } else if (this.query.chartDataType === this.allChartDataTypes.AccountTotalAssets.type) {
                return this.$t('Total Assets');
            } else if (this.query.chartDataType === this.allChartDataTypes.AccountTotalLiabilities.type) {
                return this.$t('Total Liabilities');
            }

            return this.$t('Total Amount');
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

            if (self.query.sortingType === this.allSortingTypes.DisplayOrder.type) {
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
            } else if (self.query.sortingType === this.allSortingTypes.Name.type) {
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

        let defaultSortType = self.$settings.getStatisticsSortingType();

        if (defaultSortType < self.allSortingTypes.Amount.type || defaultSortType > self.allSortingTypes.Name.type) {
            defaultSortType = self.$constants.statistics.defaultSortingType;
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
            sortingType: defaultSortType,
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
            if (this.$store.state.transactionStatisticsStateInvalid && !this.loading) {
                this.reload(null);
            }

            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            const self = this;
            let dispatchPromise = null;

            if (self.query.chartDataType === self.allChartDataTypes.ExpenseByAccount.type ||
                self.query.chartDataType === self.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.ExpenseBySecondaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeByAccount.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeByPrimaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeBySecondaryCategory.type) {
                dispatchPromise = self.$store.dispatch('loadTransactionStatistics', {
                    defaultCurrency: self.defaultCurrency
                });
            } else if (self.query.chartDataType === self.allChartDataTypes.AccountTotalAssets.type ||
                self.query.chartDataType === self.allChartDataTypes.AccountTotalLiabilities.type) {
                dispatchPromise = self.$store.dispatch('loadAllAccounts', {
                    force: true
                });
            }

            if (dispatchPromise) {
                dispatchPromise.then(() => {
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
            }
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
        setSortingType(sortingType) {
            if (sortingType < this.allSortingTypes.Amount.type || sortingType > this.allSortingTypes.Name.type) {
                this.showSortingTypePopover = false;
                return;
            }

            this.$store.dispatch('updateTransactionStatisticsFilter', {
                sortingType: sortingType
            });

            this.showSortingTypePopover = false;
            this.reload(null);
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
            this.f7router.navigate(this.getItemLinkUrl(item));
        },
        filterAccounts() {
            this.f7router.navigate('/statistic/filter/account');
        },
        filterCategories() {
            this.f7router.navigate('/statistic/filter/category');
        },
        settings() {
            this.f7router.navigate('/statistic/settings');
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
        },
        getDisplayAmount(amount, currency, textLimit) {
            amount = this.$locale.getDisplayCurrency(amount, currency);

            if (!this.showAccountBalance
                && (this.query.chartDataType === this.allChartDataTypes.AccountTotalAssets.type
                    || this.query.chartDataType === this.allChartDataTypes.AccountTotalLiabilities.type)
            ) {
                return '***';
            }

            if (textLimit) {
                this.$utilities.limitText(amount, textLimit);
            }

            return amount;
        },
        getItemLinkUrl(item) {
            const querys = [];

            if (this.query.chartDataType === this.allChartDataTypes.IncomeByAccount.type
                || this.query.chartDataType === this.allChartDataTypes.IncomeByPrimaryCategory.type
                || this.query.chartDataType === this.allChartDataTypes.IncomeBySecondaryCategory.type) {
                querys.push('type=2');
            } else if (this.query.chartDataType === this.allChartDataTypes.ExpenseByAccount.type
                || this.query.chartDataType === this.allChartDataTypes.ExpenseByPrimaryCategory.type
                || this.query.chartDataType === this.allChartDataTypes.ExpenseBySecondaryCategory.type) {
                querys.push('type=3');
            }

            if (this.query.chartDataType === this.allChartDataTypes.IncomeByAccount.type
                || this.query.chartDataType === this.allChartDataTypes.ExpenseByAccount.type
                || this.query.chartDataType === this.allChartDataTypes.AccountTotalAssets.type
                || this.query.chartDataType === this.allChartDataTypes.AccountTotalLiabilities.type) {
                querys.push('accountId=' + item.id);
            } else if (this.query.chartDataType === this.allChartDataTypes.IncomeByPrimaryCategory.type
                || this.query.chartDataType === this.allChartDataTypes.IncomeBySecondaryCategory.type
                || this.query.chartDataType === this.allChartDataTypes.ExpenseByPrimaryCategory.type
                || this.query.chartDataType === this.allChartDataTypes.ExpenseBySecondaryCategory.type) {
                querys.push('categoryId=' + item.id);
            }

            if (this.query.chartDataType !== this.allChartDataTypes.AccountTotalAssets.type
                && this.query.chartDataType !== this.allChartDataTypes.AccountTotalLiabilities.type) {
                querys.push('dateType=' + this.query.dateType);
                querys.push('minTime=' + this.query.startTime);
                querys.push('maxTime=' + this.query.endTime);
            }

            return '/transaction/list?' + querys.join('&');
        }
    }
};
</script>

<style>
.card-header.no-border:after {
    display: none;
}

.statistics-chart-header {
    font-size: var(--f7-list-item-header-font-size);
}

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

.statistics-list-item-overview-amount {
    margin-top: 2px;
    font-size: 1.5em;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-bottom: 6px;
}

.statistics-list-item-text {
    overflow: hidden;
    text-overflow: ellipsis;
}

.statistics-list-item .item-content {
    margin-top: 8px;
    margin-bottom: 12px;
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

.dark .statistics-percent-line .progressbar {
    --f7-progressbar-bg-color: #161616;
}
</style>
