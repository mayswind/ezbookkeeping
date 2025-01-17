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

        <f7-popover class="chart-data-type-popover-menu"
                    v-model:opened="showChartDataTypePopover"
                    @popover:open="scrollPopoverToSelectedItem">
            <f7-list dividers>
                <f7-list-group>
                    <f7-list-item group-title :title="$t('Categorical Analysis')" />
                    <f7-list-item :title="$t(dataType.name)"
                                  :class="{ 'list-item-selected': analysisType === allAnalysisTypes.CategoricalAnalysis && query.chartDataType === dataType.type }"
                                  :key="dataType.type"
                                  v-for="dataType in allChartDataTypes"
                                  v-show="dataType.isAvailableAnalysisType(allAnalysisTypes.CategoricalAnalysis)"
                                  @click="setChartDataType(allAnalysisTypes.CategoricalAnalysis, dataType.type)">
                        <template #after>
                            <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="analysisType === allAnalysisTypes.CategoricalAnalysis && query.chartDataType === dataType.type"></f7-icon>
                        </template>
                    </f7-list-item>
                </f7-list-group>
                <f7-list-group>
                    <f7-list-item group-title :title="$t('Trend Analysis')" />
                    <f7-list-item :title="$t(dataType.name)"
                                  :class="{ 'list-item-selected': analysisType === allAnalysisTypes.TrendAnalysis && query.chartDataType === dataType.type }"
                                  :key="dataType.type"
                                  v-for="dataType in allChartDataTypes"
                                  v-show="dataType.isAvailableAnalysisType(allAnalysisTypes.TrendAnalysis)"
                                  @click="setChartDataType(allAnalysisTypes.TrendAnalysis, dataType.type)">
                        <template #after>
                            <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="analysisType === allAnalysisTypes.TrendAnalysis && query.chartDataType === dataType.type"></f7-icon>
                        </template>
                    </f7-list-item>
                </f7-list-group>
            </f7-list>
        </f7-popover>

        <f7-card v-if="analysisType === allAnalysisTypes.CategoricalAnalysis && query.categoricalChartType === allCategoricalChartTypes.Pie.type">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header full-line text-align-right">
                    <span style="margin-right: 4px;">{{ $t('Sort by') }}</span>
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
                    :items="categoricalAnalysisData.items"
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
                    hidden-field="hidden"
                    v-else-if="!loading"
                    @click="clickPieChartItem"
                >
                    <text class="statistics-pie-chart-total-amount-title" v-if="categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                        {{ totalAmountName }}
                    </text>
                    <text class="statistics-pie-chart-total-amount-value" v-if="categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                        {{ getDisplayAmount(categoricalAnalysisData.totalAmount, defaultCurrency, 16) }}
                    </text>
                    <text class="statistics-pie-chart-total-no-data" cy="50%" v-if="!categoricalAnalysisData.items || !categoricalAnalysisData.items.length">
                        {{ $t('No data') }}
                    </text>
                </pie-chart>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="analysisType === allAnalysisTypes.CategoricalAnalysis && query.categoricalChartType === allCategoricalChartTypes.Bar.type">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header display-flex full-line justify-content-space-between">
                    <div>
                        {{ totalAmountName }}
                    </div>
                    <div class="align-self-flex-end">
                        <span style="margin-right: 4px;">{{ $t('Sort by') }}</span>
                        <f7-link href="#" popover-open=".sorting-type-popover-menu">{{ querySortingTypeName }}</f7-link>
                    </div>
                </div>
                <div class="display-flex full-line">
                    <div :class="{ 'statistics-list-item-overview-amount': true, 'text-expense': query.chartDataType === allChartDataTypes.ExpenseByAccount.type || query.chartDataType === allChartDataTypes.ExpenseByPrimaryCategory.type || query.chartDataType === allChartDataTypes.ExpenseBySecondaryCategory.type, 'text-income': query.chartDataType === allChartDataTypes.IncomeByAccount.type || query.chartDataType === allChartDataTypes.IncomeByPrimaryCategory.type || query.chartDataType === allChartDataTypes.IncomeBySecondaryCategory.type }">
                        <span v-if="!loading && categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                            {{ getDisplayAmount(categoricalAnalysisData.totalAmount, defaultCurrency) }}
                        </span>
                        <span :class="{ 'skeleton-text': loading }" v-else-if="loading || !categoricalAnalysisData || !categoricalAnalysisData.items || !categoricalAnalysisData.items.length">
                            {{ loading ? '***.**' : '---' }}
                        </span>
                    </div>
                </div>
            </f7-card-header>
            <f7-card-content style="margin-top: -14px" :padding="false">
                <f7-list class="statistics-list-item skeleton-text" v-if="loading">
                    <f7-list-item link="#" :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                        <template #media>
                            <div class="display-flex no-padding-horizontal">
                                <div class="display-flex align-items-center statistics-icon">
                                    <f7-icon f7="app_fill"></f7-icon>
                                </div>
                            </div>
                        </template>
                        <template #title>
                            <div class="statistics-list-item-text">
                                <span>Category Name</span>
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

                <f7-list v-else-if="!loading && (!categoricalAnalysisData || !categoricalAnalysisData.items || !categoricalAnalysisData.items.length)">
                    <f7-list-item :title="$t('No transaction data')"></f7-list-item>
                </f7-list>

                <f7-list v-else-if="!loading && categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                    <f7-list-item class="statistics-list-item"
                                  :link="getTransactionItemLinkUrl(item.id)"
                                  :key="idx"
                                  v-for="(item, idx) in categoricalAnalysisData.items"
                                  v-show="!item.hidden"
                    >
                        <template #media>
                            <div class="display-flex no-padding-horizontal">
                                <div class="display-flex align-items-center statistics-icon">
                                    <ItemIcon :icon-type="queryChartDataCategory" :icon-id="item.icon" :color="item.color" v-if="item.icon"></ItemIcon>
                                    <f7-icon f7="pencil_ellipsis_rectangle" v-else-if="!item.icon"></f7-icon>
                                </div>
                            </div>
                        </template>

                        <template #title>
                            <div class="statistics-list-item-text">
                                <span>{{ item.name }}</span>
                                <small class="statistics-percent" v-if="item.percent >= 0">{{ getDisplayPercent(item.percent, 2, '&lt;0.01') }}</small>
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

        <f7-card v-else-if="analysisType === allAnalysisTypes.TrendAnalysis">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header display-flex full-line justify-content-space-between">
                    <div></div>
                    <div class="align-self-flex-end">
                        <span style="margin-right: 4px;">{{ $t('Sort by') }}</span>
                        <f7-link href="#" popover-open=".sorting-type-popover-menu">{{ querySortingTypeName }}</f7-link>
                    </div>
                </div>
            </f7-card-header>
            <f7-card-content style="margin-top: -14px" :padding="false">
                <trends-bar-chart
                    :loading="loading || reloading"
                    :start-year-month="query.trendChartStartYearMonth"
                    :end-year-month="query.trendChartEndYearMonth"
                    :sorting-type="query.sortingType"
                    :date-aggregation-type="trendDateAggregationType"
                    :items="trendsAnalysisData && trendsAnalysisData.items && trendsAnalysisData.items.length ? trendsAnalysisData.items : []"
                    :translate-name="translateNameInTrendsChart"
                    :default-currency="defaultCurrency"
                    id-field="id"
                    name-field="name"
                    value-field="totalAmount"
                    hidden-field="hidden"
                    display-orders-field="displayOrders"
                    @click="clickTrendChartItem"
                />
            </f7-card-content>
        </f7-card>

        <f7-popover class="sorting-type-popover-menu"
                    v-model:opened="showSortingTypePopover">
            <f7-list dividers>
                <f7-list-item :title="sortingType.displayName"
                              :class="{ 'list-item-selected': query.sortingType === sortingType.type }"
                              :key="sortingType.type"
                              v-for="sortingType in allSortingTypes"
                              @click="setSortingType(sortingType.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.sortingType === sortingType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-toolbar tabbar bottom class="toolbar-item-auto-size">
            <f7-link :class="{ 'disabled': reloading || !canShiftDateRange(query) }" @click="shiftDateRange(query, -1)">
                <f7-icon f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': reloading || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type }" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.maxTime > 0 || query.minTime > 0 }">{{ dateRangeName(query) }}</span>
            </f7-link>
            <f7-link :class="{ 'disabled': reloading || !canShiftDateRange(query) }" @click="shiftDateRange(query, 1)">
                <f7-icon f7="arrow_right_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': reloading }" popover-open=".date-aggregation-popover-menu"
                     v-if="analysisType === allAnalysisTypes.TrendAnalysis">
                <span :class="{ 'tabbar-item-changed': trendDateAggregationType !== defaultTrendDateAggregationType }">{{ queryTrendDateAggregationTypeName }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" :key="chartType.type"
                     v-for="chartType in allChartTypes" @click="setChartType(chartType.type)">
                <span :class="{ 'tabbar-item-changed': queryChartType === chartType.type }">{{ chartType.displayName }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-popover class="date-popover-menu"
                    v-model:opened="showDatePopover"
                    @popover:open="scrollPopoverToSelectedItem">
            <f7-list dividers>
                <f7-list-item :title="dateRange.displayName"
                              :class="{ 'list-item-selected': queryDateType === dateRange.type }"
                              :key="dateRange.type"
                              v-for="dateRange in allDateRangesArray"
                              @click="setDateFilter(dateRange.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="queryDateType === dateRange.type"></f7-icon>
                    </template>
                    <template #footer>
                        <div v-if="dateRange.type === allDateRanges.Custom.type && showCustomDateRange()">
                            <span>{{ queryStartTime }}</span>
                            <span>&nbsp;-&nbsp;</span>
                            <br/>
                            <span>{{ queryEndTime }}</span>
                        </div>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="date-aggregation-popover-menu"
                    v-model:opened="showDateAggregationPopover"
                    @popover:open="scrollPopoverToSelectedItem">
            <f7-list dividers>
                <f7-list-item :title="aggregationType.displayName"
                              :class="{ 'list-item-selected': trendDateAggregationType === aggregationType.type }"
                              :key="aggregationType.type"
                              v-for="aggregationType in allDateAggregationTypes"
                              @click="setTrendDateAggregationType(aggregationType.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="trendDateAggregationType === aggregationType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <date-range-selection-sheet :title="$t('Custom Date Range')"
                                    :min-time="query.categoricalChartStartTime"
                                    :max-time="query.categoricalChartEndTime"
                                    v-model:show="showCustomDateRangeSheet"
                                    @dateRange:change="setCustomDateFilter">
        </date-range-selection-sheet>

        <month-range-selection-sheet :title="$t('Custom Date Range')"
                                     :min-time="query.trendChartStartYearMonth"
                                     :max-time="query.trendChartEndYearMonth"
                                     v-model:show="showCustomMonthRangeSheet"
                                     @dateRange:change="setCustomDateFilter">
        </month-range-selection-sheet>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="filterAccounts">{{ $t('Filter Accounts') }}</f7-actions-button>
                <f7-actions-button @click="filterCategories">{{ $t('Filter Transaction Categories') }}</f7-actions-button>
                <f7-actions-button @click="filterTags">{{ $t('Filter Transaction Tags') }}</f7-actions-button>
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
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useStatisticsStore } from '@/stores/statistics.js';

import { DateRangeScene, DateRange } from '@/core/datetime.ts';
import {
    StatisticsAnalysisType,
    CategoricalChartType,
    ChartDataType,
    ChartSortingType,
    ChartDateAggregationType
} from '@/core/statistics.ts';
import { getNameByKeyValue, limitText } from '@/lib/common.ts';
import { formatPercent } from '@/lib/numeral.ts';
import {
    getYearAndMonthFromUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getShiftedDateRangeAndDateType,
    getDateTypeByDateRange,
    getDateRangeByDateType
} from '@/lib/datetime.ts';
import { scrollToSelectedItem } from '@/lib/ui/mobile.ts';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            loading: true,
            loadingError: null,
            reloading: false,
            analysisType: StatisticsAnalysisType.CategoricalAnalysis,
            trendDateAggregationType: ChartDateAggregationType.Default.type,
            showChartDataTypePopover: false,
            showSortingTypePopover: false,
            showDatePopover: false,
            showDateAggregationPopover: false,
            showCustomDateRangeSheet: false,
            showCustomMonthRangeSheet: false,
            showMoreActionSheet: false
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useStatisticsStore),
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        defaultTrendDateAggregationType() {
            return ChartDateAggregationType.Default.type;
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        query() {
            return this.statisticsStore.transactionStatisticsFilter;
        },
        queryChartDataCategory() {
            return this.statisticsStore.categoricalAnalysisChartDataCategory;
        },
        queryChartType: {
            get: function () {
                if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                    return this.query.categoricalChartType;
                } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                    return this.query.trendChartType;
                } else {
                    return null;
                }
            },
            set: function(value) {
                this.setChartType(value);
            }
        },
        queryChartDataTypeName() {
            const queryChartDataTypeName = getNameByKeyValue(ChartDataType.values(), this.query.chartDataType, 'type', 'name', 'Statistics');
            return this.$t(queryChartDataTypeName);
        },
        querySortingTypeName() {
            const querySortingTypeName = getNameByKeyValue(ChartSortingType.values(), this.query.sortingType, 'type', 'name', 'System Default');
            return this.$t(querySortingTypeName);
        },
        queryDateType() {
            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.query.categoricalChartDateType;
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.query.trendChartDateType;
            } else {
                return null;
            }
        },
        queryTrendDateAggregationTypeName() {
            return getNameByKeyValue(this.allDateAggregationTypes, this.trendDateAggregationType, 'type', 'displayName', '');
        },
        queryStartTime() {
            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.categoricalChartStartTime);
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, getYearMonthFirstUnixTime(this.query.trendChartStartYearMonth));
            } else {
                return '';
            }
        },
        queryEndTime() {
            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.categoricalChartEndTime);
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, getYearMonthLastUnixTime(this.query.trendChartEndYearMonth));
            } else {
                return '';
            }
        },
        allAnalysisTypes() {
            return StatisticsAnalysisType;
        },
        allChartTypes() {
            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.getAllCategoricalChartTypes();
            } else {
                return [];
            }
        },
        allCategoricalChartTypes() {
            return CategoricalChartType.all();
        },
        allChartDataTypes() {
            return ChartDataType.all();
        },
        allSortingTypes() {
            return this.$locale.getAllStatisticsSortingTypes();
        },
        allDateAggregationTypes() {
            return this.$locale.getAllStatisticsDateAggregationTypes();
        },
        allDateRanges() {
            return DateRange.all();
        },
        allDateRangesArray() {
            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.getAllDateRanges(DateRangeScene.Normal, true);
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.getAllDateRanges(DateRangeScene.TrendAnalysis, true);
            } else {
                return [];
            }
        },
        showAccountBalance() {
            return this.settingsStore.appSettings.showAccountBalance;
        },
        totalAmountName() {
            if (this.query.chartDataType === ChartDataType.IncomeByAccount.type
                || this.query.chartDataType === ChartDataType.IncomeByPrimaryCategory.type
                || this.query.chartDataType === ChartDataType.IncomeBySecondaryCategory.type) {
                return this.$t('Total Income');
            } else if (this.query.chartDataType === ChartDataType.ExpenseByAccount.type
                || this.query.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type
                || this.query.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type) {
                return this.$t('Total Expense');
            } else if (this.query.chartDataType === ChartDataType.AccountTotalAssets.type) {
                return this.$t('Total Assets');
            } else if (this.query.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
                return this.$t('Total Liabilities');
            }

            return this.$t('Total Amount');
        },
        categoricalAnalysisData() {
            return this.statisticsStore.categoricalAnalysisData;
        },
        trendsAnalysisData() {
            return this.statisticsStore.trendsAnalysisData;
        },
        translateNameInTrendsChart() {
            return this.query.chartDataType === ChartDataType.TotalExpense.type ||
                this.query.chartDataType === ChartDataType.TotalIncome.type ||
                this.query.chartDataType === ChartDataType.TotalBalance.type;
        },
        showAmountInChart() {
            if (!this.showAccountBalance
                && (this.query.chartDataType === ChartDataType.AccountTotalAssets.type || this.query.chartDataType === ChartDataType.AccountTotalLiabilities.type)) {
                return false;
            }

            return true;
        }
    },
    created() {
        const self = this;

        self.statisticsStore.initTransactionStatisticsFilter(self.analysisType);

        Promise.all([
            self.accountsStore.loadAllAccounts({ force: false }),
            self.transactionCategoriesStore.loadAllCategories({ force: false })
        ]).then(() => {
            if (self.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return self.statisticsStore.loadCategoricalAnalysis({
                    force: false
                });
            } else if (self.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                return self.statisticsStore.loadTrendAnalysis({
                    force: false
                });
            }
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
            if (this.statisticsStore.transactionStatisticsStateInvalid && !this.loading) {
                this.reload(null);
            }

            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            const self = this;
            const force = !!done;
            let dispatchPromise = null;

            self.reloading = true;

            if (self.query.chartDataType === ChartDataType.ExpenseByAccount.type ||
                self.query.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
                self.query.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
                self.query.chartDataType === ChartDataType.IncomeByAccount.type ||
                self.query.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
                self.query.chartDataType === ChartDataType.IncomeBySecondaryCategory.type ||
                self.query.chartDataType === ChartDataType.TotalExpense.type ||
                self.query.chartDataType === ChartDataType.TotalIncome.type ||
                self.query.chartDataType === ChartDataType.TotalBalance.type) {
                if (self.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                    dispatchPromise = self.statisticsStore.loadCategoricalAnalysis({
                        force: force
                    });
                } else if (self.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                    dispatchPromise = self.statisticsStore.loadTrendAnalysis({
                        force: force
                    });
                }
            } else if (self.query.chartDataType === ChartDataType.AccountTotalAssets.type ||
                self.query.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
                dispatchPromise = self.accountsStore.loadAllAccounts({
                    force: force
                });
            }

            if (dispatchPromise) {
                dispatchPromise.then(() => {
                    self.reloading = false;

                    if (done) {
                        done();
                    }

                    if (force) {
                        self.$toast('Data has been updated');
                    }
                }).catch(error => {
                    self.reloading = false;

                    if (done) {
                        done();
                    }

                    if (!error.processed) {
                        self.$toast(error.message || error);
                    }
                });
            } else {
                self.reloading = false;
            }
        },
        setChartType(chartType) {
            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartType: chartType
                });
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                this.statisticsStore.updateTransactionStatisticsFilter({
                    trendChartType: chartType
                });
            }
        },
        setChartDataType(analysisType, chartDataType) {
            let analysisTypeChanged = false;

            if (this.analysisType !== analysisType) {
                if (!ChartDataType.isAvailableForAnalysisType(this.queryChartDataType, analysisType)) {
                    this.statisticsStore.updateTransactionStatisticsFilter({
                        chartDataType: ChartDataType.Default.type
                    });
                }

                this.analysisType = analysisType;
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
                analysisTypeChanged = true;
            }

            this.statisticsStore.updateTransactionStatisticsFilter({
                chartDataType: chartDataType
            });

            this.showChartDataTypePopover = false;

            if (analysisTypeChanged) {
                this.reload(null);
            }
        },
        setSortingType(sortingType) {
            if (sortingType < ChartSortingType.Amount.type || sortingType > ChartSortingType.Name.type) {
                this.showSortingTypePopover = false;
                return;
            }

            this.statisticsStore.updateTransactionStatisticsFilter({
                sortingType: sortingType
            });

            this.showSortingTypePopover = false;
        },
        setTrendDateAggregationType(aggregationType) {
            this.trendDateAggregationType = aggregationType;
            this.showDateAggregationPopover = false;
        },
        setDateFilter(dateType) {
            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                if (dateType === this.allDateRanges.Custom.type) { // Custom
                    this.showCustomDateRangeSheet = true;
                    this.showDatePopover = false;
                    return;
                } else if (this.query.categoricalChartDateType === dateType) {
                    return;
                }
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                if (dateType === this.allDateRanges.Custom.type) { // Custom
                    this.showCustomMonthRangeSheet = true;
                    this.showDatePopover = false;
                    return;
                } else if (this.query.trendChartDateType === dateType) {
                    return;
                }
            }

            const dateRange = getDateRangeByDateType(dateType, this.firstDayOfWeek);

            if (!dateRange) {
                return;
            }

            let changed = false;

            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: dateRange.dateType,
                    categoricalChartStartTime: dateRange.minTime,
                    categoricalChartEndTime: dateRange.maxTime
                });
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    trendChartDateType: dateRange.dateType,
                    trendChartStartYearMonth: getYearAndMonthFromUnixTime(dateRange.minTime),
                    trendChartEndYearMonth: getYearAndMonthFromUnixTime(dateRange.maxTime)
                });
            }

            this.showDatePopover = false;

            if (changed) {
                this.reload(null);
            }
        },
        setCustomDateFilter(startTime, endTime) {
            if (!startTime || !endTime) {
                return;
            }

            let changed = false;

            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                const chartDateType = getDateTypeByDateRange(startTime, endTime, this.firstDayOfWeek, DateRangeScene.Normal);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: chartDateType,
                    categoricalChartStartTime: startTime,
                    categoricalChartEndTime: endTime
                });

                this.showCustomDateRangeSheet = false;
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                const chartDateType = getDateTypeByDateRange(getYearMonthFirstUnixTime(startTime), getYearMonthLastUnixTime(endTime), this.firstDayOfWeek, DateRangeScene.TrendAnalysis);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    trendChartDateType: chartDateType,
                    trendChartStartYearMonth: startTime,
                    trendChartEndYearMonth: endTime
                });

                this.showCustomMonthRangeSheet = false;
            }

            if (changed) {
                this.reload(null);
            }
        },
        showCustomDateRange() {
            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.query.categoricalChartDateType === this.allDateRanges.Custom.type && this.query.categoricalChartStartTime && this.query.categoricalChartEndTime;
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.query.trendChartDateType === this.allDateRanges.Custom.type && this.query.trendChartStartYearMonth && this.query.trendChartEndYearMonth;
            } else {
                return false;
            }
        },
        canShiftDateRange(query) {
            if (query.chartDataType === ChartDataType.AccountTotalAssets.type || query.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
                return false;
            }

            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return query.categoricalChartDateType !== this.allDateRanges.All.type;
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                return query.trendChartDateType !== this.allDateRanges.All.type;
            } else {
                return false;
            }
        },
        shiftDateRange(query, scale) {
            if (this.query.categoricalChartDateType === this.allDateRanges.All.type) {
                return;
            }

            let changed = false;

            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                const newDateRange = getShiftedDateRangeAndDateType(query.categoricalChartStartTime, query.categoricalChartEndTime, scale, this.firstDayOfWeek, DateRangeScene.Normal);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: newDateRange.dateType,
                    categoricalChartStartTime: newDateRange.minTime,
                    categoricalChartEndTime: newDateRange.maxTime
                });
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                const newDateRange = getShiftedDateRangeAndDateType(getYearMonthFirstUnixTime(query.trendChartStartYearMonth), getYearMonthLastUnixTime(query.trendChartEndYearMonth), scale, this.firstDayOfWeek, DateRangeScene.TrendAnalysis);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    trendChartDateType: newDateRange.dateType,
                    trendChartStartYearMonth: getYearAndMonthFromUnixTime(newDateRange.minTime),
                    trendChartEndYearMonth: getYearAndMonthFromUnixTime(newDateRange.maxTime)
                });
            }

            if (changed) {
                this.reload(null);
            }
        },
        dateRangeName(query) {
            if (this.query.chartDataType === ChartDataType.AccountTotalAssets.type ||
                this.query.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
                return this.$t(this.allDateRanges.All.name);
            }

            if (this.analysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.getDateRangeDisplayName(this.userStore, query.categoricalChartDateType, query.categoricalChartStartTime, query.categoricalChartEndTime);
            } else if (this.analysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.getDateRangeDisplayName(this.userStore, query.trendChartDateType, getYearMonthFirstUnixTime(query.trendChartStartYearMonth), getYearMonthLastUnixTime(query.trendChartEndYearMonth));
            } else {
                return '';
            }
        },
        clickPieChartItem(item) {
            this.f7router.navigate(this.getTransactionItemLinkUrl(item.id));
        },
        clickTrendChartItem(item) {
            this.f7router.navigate(this.getTransactionItemLinkUrl(item.itemId, item.dateRange));
        },
        filterAccounts() {
            this.f7router.navigate('/settings/filter/account?type=statisticsCurrent');
        },
        filterCategories() {
            this.f7router.navigate('/settings/filter/category?type=statisticsCurrent');
        },
        filterTags() {
            this.f7router.navigate('/settings/filter/tag?type=statisticsCurrent');
        },
        settings() {
            this.f7router.navigate('/statistic/settings');
        },
        scrollPopoverToSelectedItem(event) {
            scrollToSelectedItem(event.$el, '.popover-inner', 'li.list-item-selected');
        },
        getDisplayAmount(amount, currency, textLimit) {
            amount = this.getDisplayCurrency(amount, currency);

            if (!this.showAccountBalance
                && (this.query.chartDataType === ChartDataType.AccountTotalAssets.type
                    || this.query.chartDataType === ChartDataType.AccountTotalLiabilities.type)
            ) {
                return '***';
            }

            if (textLimit) {
                return limitText(amount, textLimit);
            }

            return amount;
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        },
        getDisplayPercent(value, precision, lowPrecisionValue) {
            return formatPercent(value, precision, lowPrecisionValue);
        },
        getTransactionItemLinkUrl(itemId, dateRange) {
            return `/transaction/list?${this.statisticsStore.getTransactionListPageParams(this.analysisType, itemId, dateRange)}`;
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

.dark .statistics-percent-line .progressbar {
    --f7-progressbar-bg-color: #161616;
}

.chart-data-type-popover-menu .popover-inner{
    max-height: 440px;
    overflow-y: auto;
}
</style>
