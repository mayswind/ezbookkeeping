<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title>
                <f7-link popover-open=".chart-data-type-popover-menu">
                    <span style="color: var(--f7-text-color)">{{ queryChartDataTypeName }}</span>
                    <f7-icon class="page-title-bar-icon" color="gray" style="opacity: 0.5" f7="chevron_down_circle_fill"></f7-icon>
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
                    <f7-list-item group-title :title="tt('Categorical Analysis')" />
                    <f7-list-item :title="tt(dataType.name)"
                                  :class="{ 'list-item-selected': analysisType === StatisticsAnalysisType.CategoricalAnalysis && query.chartDataType === dataType.type }"
                                  :key="dataType.type"
                                  v-for="dataType in ChartDataType.values(StatisticsAnalysisType.CategoricalAnalysis)"
                                  @click="setChartDataType(StatisticsAnalysisType.CategoricalAnalysis, dataType.type)">
                        <template #after>
                            <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="analysisType === StatisticsAnalysisType.CategoricalAnalysis && query.chartDataType === dataType.type"></f7-icon>
                        </template>
                    </f7-list-item>
                </f7-list-group>
                <f7-list-group>
                    <f7-list-item group-title :title="tt('Trend Analysis')" />
                    <f7-list-item :title="tt(dataType.name)"
                                  :class="{ 'list-item-selected': analysisType === StatisticsAnalysisType.TrendAnalysis && query.chartDataType === dataType.type }"
                                  :key="dataType.type"
                                  v-for="dataType in ChartDataType.values(StatisticsAnalysisType.TrendAnalysis)"
                                  @click="setChartDataType(StatisticsAnalysisType.TrendAnalysis, dataType.type)">
                        <template #after>
                            <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="analysisType === StatisticsAnalysisType.TrendAnalysis && query.chartDataType === dataType.type"></f7-icon>
                        </template>
                    </f7-list-item>
                </f7-list-group>
            </f7-list>
        </f7-popover>

        <f7-card v-if="analysisType === StatisticsAnalysisType.CategoricalAnalysis && query.categoricalChartType === CategoricalChartType.Pie.type">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header full-line text-align-right">
                    <span style="margin-right: 4px;">{{ tt('Sort by') }}</span>
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
                    name-field="name"
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
                    @click="onClickPieChartItem"
                >
                    <text class="statistics-pie-chart-total-amount-title" v-if="categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                        {{ totalAmountName }}
                    </text>
                    <text class="statistics-pie-chart-total-amount-value" v-if="categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                        {{ getDisplayAmount(categoricalAnalysisData.totalAmount, defaultCurrency, 16) }}
                    </text>
                    <text class="statistics-pie-chart-total-no-data" cy="50%" v-if="!categoricalAnalysisData.items || !categoricalAnalysisData.items.length">
                        {{ tt('No data') }}
                    </text>
                </pie-chart>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="analysisType === StatisticsAnalysisType.CategoricalAnalysis && query.categoricalChartType === CategoricalChartType.Bar.type">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header display-flex full-line justify-content-space-between">
                    <div>
                        {{ totalAmountName }}
                    </div>
                    <div class="align-self-flex-end">
                        <span style="margin-right: 4px;">{{ tt('Sort by') }}</span>
                        <f7-link href="#" popover-open=".sorting-type-popover-menu">{{ querySortingTypeName }}</f7-link>
                    </div>
                </div>
                <div class="display-flex full-line">
                    <div :class="{ 'statistics-list-item-overview-amount': true, 'text-expense': query.chartDataType === ChartDataType.ExpenseByAccount.type || query.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type || query.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type, 'text-income': query.chartDataType === ChartDataType.IncomeByAccount.type || query.chartDataType === ChartDataType.IncomeByPrimaryCategory.type || query.chartDataType === ChartDataType.IncomeBySecondaryCategory.type }">
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
                    <f7-list-item :title="tt('No transaction data')"></f7-list-item>
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
                                <small class="statistics-percent" v-if="item.percent >= 0">{{ formatPercent(item.percent, 2, '&lt;0.01') }}</small>
                            </div>
                        </template>

                        <template #after>
                            <span>{{ getDisplayAmount(item.totalAmount, defaultCurrency) }}</span>
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

        <f7-card v-else-if="analysisType === StatisticsAnalysisType.TrendAnalysis">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header display-flex full-line justify-content-space-between">
                    <div></div>
                    <div class="align-self-flex-end">
                        <span style="margin-right: 4px;">{{ tt('Sort by') }}</span>
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
                    :fiscal-year-start="fiscalYearStart"
                    :items="trendsAnalysisData && trendsAnalysisData.items && trendsAnalysisData.items.length ? trendsAnalysisData.items : []"
                    :translate-name="translateNameInTrendsChart"
                    :default-currency="defaultCurrency"
                    id-field="id"
                    name-field="name"
                    value-field="totalAmount"
                    hidden-field="hidden"
                    display-orders-field="displayOrders"
                    @click="onClickTrendChartItem"
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
            <f7-link :class="{ 'disabled': reloading || !canShiftDateRange }" @click="shiftDateRange(-1)">
                <f7-icon f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': reloading || query.chartDataType === ChartDataType.AccountTotalAssets.type || query.chartDataType === ChartDataType.AccountTotalLiabilities.type }" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': isQueryDateRangeChanged }">{{ queryDateRangeName }}</span>
            </f7-link>
            <f7-link :class="{ 'disabled': reloading || !canShiftDateRange }" @click="shiftDateRange(1)">
                <f7-icon f7="arrow_right_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': reloading }" popover-open=".date-aggregation-popover-menu"
                     v-if="analysisType === StatisticsAnalysisType.TrendAnalysis">
                <span :class="{ 'tabbar-item-changed': trendDateAggregationType !== ChartDateAggregationType.Default.type }">{{ queryTrendDateAggregationTypeName }}</span>
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
                              v-for="dateRange in allDateRanges"
                              @click="setDateFilter(dateRange.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="queryDateType === dateRange.type"></f7-icon>
                    </template>
                    <template #footer>
                        <div v-if="dateRange.type === DateRange.Custom.type && showCustomDateRange">
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

        <date-range-selection-sheet :title="tt('Custom Date Range')"
                                    :min-time="query.categoricalChartStartTime"
                                    :max-time="query.categoricalChartEndTime"
                                    v-model:show="showCustomDateRangeSheet"
                                    @dateRange:change="setCustomDateFilter">
        </date-range-selection-sheet>

        <month-range-selection-sheet :title="tt('Custom Date Range')"
                                     :min-time="query.trendChartStartYearMonth"
                                     :max-time="query.trendChartEndYearMonth"
                                     v-model:show="showCustomMonthRangeSheet"
                                     @dateRange:change="setCustomDateFilter">
        </month-range-selection-sheet>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="filterAccounts">{{ tt('Filter Accounts') }}</f7-actions-button>
                <f7-actions-button @click="filterCategories">{{ tt('Filter Transaction Categories') }}</f7-actions-button>
                <f7-actions-button @click="filterTags">{{ tt('Filter Transaction Tags') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button @click="settings">{{ tt('Settings') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useStatisticsTransactionPageBase } from '@/views/base/statistics/StatisticsTransactionPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { type TimeRangeAndDateType, DateRangeScene, DateRange } from '@/core/datetime.ts';
import {
    StatisticsAnalysisType,
    CategoricalChartType,
    ChartDataType,
    ChartSortingType,
    ChartDateAggregationType
} from '@/core/statistics.ts';

import { isString, isNumber } from '@/lib/common.ts';
import {
    getYearAndMonthFromUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getShiftedDateRangeAndDateType,
    getDateTypeByDateRange,
    getDateRangeByDateType
} from '@/lib/datetime.ts';
import { type Framework7Dom, useI18nUIComponents, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, getAllCategoricalChartTypes, formatPercent } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const {
    loading,
    analysisType,
    trendDateAggregationType,
    defaultCurrency,
    firstDayOfWeek,
    fiscalYearStart,
    allDateRanges,
    allSortingTypes,
    allDateAggregationTypes,
    query,
    queryChartDataCategory,
    queryDateType,
    queryStartTime,
    queryEndTime,
    queryDateRangeName,
    queryChartDataTypeName,
    querySortingTypeName,
    queryTrendDateAggregationTypeName,
    isQueryDateRangeChanged,
    canShiftDateRange,
    showCustomDateRange,
    showAmountInChart,
    totalAmountName,
    translateNameInTrendsChart,
    categoricalAnalysisData,
    trendsAnalysisData,
    getDisplayAmount
} = useStatisticsTransactionPageBase();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const statisticsStore = useStatisticsStore();

const loadingError = ref<unknown | null>(null);
const reloading = ref<boolean>(false);
const showChartDataTypePopover = ref<boolean>(false);
const showSortingTypePopover = ref<boolean>(false);
const showDatePopover = ref<boolean>(false);
const showDateAggregationPopover = ref<boolean>(false);
const showCustomDateRangeSheet = ref<boolean>(false);
const showCustomMonthRangeSheet = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);

const allChartTypes = computed<TypeAndDisplayName[]>(() => {
    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        return getAllCategoricalChartTypes();
    } else {
        return [];
    }
});

const queryChartType = computed<number | undefined>({
    get: () => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return query.value.categoricalChartType;
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return query.value.trendChartType;
        } else {
            return undefined;
        }
    },
    set: (value: number | undefined) => {
        setChartType(value);
    }
});

function getTransactionItemLinkUrl(itemId: string, dateRange?: TimeRangeAndDateType): string {
    return `/transaction/list?${statisticsStore.getTransactionListPageParams(analysisType.value, itemId, dateRange)}`;
}

function init(): void {
    statisticsStore.initTransactionStatisticsFilter(analysisType.value);

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false })
    ]).then(() => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return statisticsStore.loadCategoricalAnalysis({
                force: false
            }) as Promise<unknown>;
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return statisticsStore.loadTrendAnalysis({
                force: false
            }) as Promise<unknown>;
        } else {
            return Promise.reject('An error occurred');
        }
    }).then(() => {
        loading.value = false;
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function reload(done?: () => void): void {
    const force = !!done;
    let dispatchPromise: Promise<unknown> | null = null;

    reloading.value = true;

    if (query.value.chartDataType === ChartDataType.ExpenseByAccount.type ||
        query.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
        query.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
        query.value.chartDataType === ChartDataType.IncomeByAccount.type ||
        query.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
        query.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type ||
        query.value.chartDataType === ChartDataType.TotalExpense.type ||
        query.value.chartDataType === ChartDataType.TotalIncome.type ||
        query.value.chartDataType === ChartDataType.TotalBalance.type) {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            dispatchPromise = statisticsStore.loadCategoricalAnalysis({
                force: force
            });
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            dispatchPromise = statisticsStore.loadTrendAnalysis({
                force: force
            });
        }
    } else if (query.value.chartDataType === ChartDataType.AccountTotalAssets.type ||
        query.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
        dispatchPromise = accountsStore.loadAllAccounts({
            force: force
        });
    }

    if (dispatchPromise) {
        dispatchPromise.then(() => {
            reloading.value = false;

            done?.();

            if (force) {
                showToast('Data has been updated');
            }
        }).catch(error => {
            reloading.value = false;

            done?.();

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    } else {
        reloading.value = false;
    }
}

function setChartType(type?: number): void {
    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        statisticsStore.updateTransactionStatisticsFilter({
            categoricalChartType: type
        });
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        statisticsStore.updateTransactionStatisticsFilter({
            trendChartType: type
        });
    }
}

function setChartDataType(type: number, chartDataType: number): void {
    let analysisTypeChanged = false;

    if (analysisType.value !== type) {
        if (!ChartDataType.isAvailableForAnalysisType(query.value.chartDataType, type)) {
            statisticsStore.updateTransactionStatisticsFilter({
                chartDataType: ChartDataType.Default.type
            });
        }

        analysisType.value = type;
        statisticsStore.updateTransactionStatisticsInvalidState(true);
        analysisTypeChanged = true;
    }

    statisticsStore.updateTransactionStatisticsFilter({
        chartDataType: chartDataType
    });

    showChartDataTypePopover.value = false;

    if (analysisTypeChanged) {
        reload();
    }
}

function setSortingType(type: number): void {
    if (type < ChartSortingType.Amount.type || type > ChartSortingType.Name.type) {
        showSortingTypePopover.value = false;
        return;
    }

    statisticsStore.updateTransactionStatisticsFilter({
        sortingType: type
    });

    showSortingTypePopover.value = false;
}

function setTrendDateAggregationType(type: number): void {
    trendDateAggregationType.value = type;
    showDateAggregationPopover.value = false;
}

function setDateFilter(dateType: number): void {
    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        if (dateType === DateRange.Custom.type) { // Custom
            showCustomDateRangeSheet.value = true;
            showDatePopover.value = false;
            return;
        } else if (query.value.categoricalChartDateType === dateType) {
            return;
        }
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        if (dateType === DateRange.Custom.type) { // Custom
            showCustomMonthRangeSheet.value = true;
            showDatePopover.value = false;
            return;
        } else if (query.value.trendChartDateType === dateType) {
            return;
        }
    }

    const dateRange = getDateRangeByDateType(dateType, firstDayOfWeek.value, fiscalYearStart.value);

    if (!dateRange) {
        return;
    }

    let changed = false;

    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        changed = statisticsStore.updateTransactionStatisticsFilter({
            categoricalChartDateType: dateRange.dateType,
            categoricalChartStartTime: dateRange.minTime,
            categoricalChartEndTime: dateRange.maxTime
        });
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        changed = statisticsStore.updateTransactionStatisticsFilter({
            trendChartDateType: dateRange.dateType,
            trendChartStartYearMonth: getYearAndMonthFromUnixTime(dateRange.minTime),
            trendChartEndYearMonth: getYearAndMonthFromUnixTime(dateRange.maxTime)
        });
    }

    showDatePopover.value = false;

    if (changed) {
        reload();
    }
}

function setCustomDateFilter(startTime: number | string, endTime: number | string): void {
    if (!startTime || !endTime) {
        return;
    }

    let changed = false;

    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis && isNumber(startTime) && isNumber(endTime)) {
        const chartDateType = getDateTypeByDateRange(startTime, endTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);

        changed = statisticsStore.updateTransactionStatisticsFilter({
            categoricalChartDateType: chartDateType,
            categoricalChartStartTime: startTime,
            categoricalChartEndTime: endTime
        });

        showCustomDateRangeSheet.value = false;
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis && isString(startTime) && isString(endTime)) {
        const chartDateType = getDateTypeByDateRange(getYearMonthFirstUnixTime(startTime), getYearMonthLastUnixTime(endTime), firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.TrendAnalysis);

        changed = statisticsStore.updateTransactionStatisticsFilter({
            trendChartDateType: chartDateType,
            trendChartStartYearMonth: startTime,
            trendChartEndYearMonth: endTime
        });

        showCustomMonthRangeSheet.value = false;
    }

    if (changed) {
        reload();
    }
}

function shiftDateRange(scale: number): void {
    if (query.value.categoricalChartDateType === DateRange.All.type) {
        return;
    }

    let changed = false;

    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        const newDateRange = getShiftedDateRangeAndDateType(query.value.categoricalChartStartTime, query.value.categoricalChartEndTime, scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);

        changed = statisticsStore.updateTransactionStatisticsFilter({
            categoricalChartDateType: newDateRange.dateType,
            categoricalChartStartTime: newDateRange.minTime,
            categoricalChartEndTime: newDateRange.maxTime
        });
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        const newDateRange = getShiftedDateRangeAndDateType(getYearMonthFirstUnixTime(query.value.trendChartStartYearMonth), getYearMonthLastUnixTime(query.value.trendChartEndYearMonth), scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.TrendAnalysis);

        changed = statisticsStore.updateTransactionStatisticsFilter({
            trendChartDateType: newDateRange.dateType,
            trendChartStartYearMonth: getYearAndMonthFromUnixTime(newDateRange.minTime),
            trendChartEndYearMonth: getYearAndMonthFromUnixTime(newDateRange.maxTime)
        });
    }

    if (changed) {
        reload();
    }
}

function filterAccounts(): void {
    props.f7router.navigate('/settings/filter/account?type=statisticsCurrent');
}

function filterCategories(): void {
    props.f7router.navigate('/settings/filter/category?type=statisticsCurrent');
}

function filterTags(): void {
    props.f7router.navigate('/settings/filter/tag?type=statisticsCurrent');
}

function settings(): void {
    props.f7router.navigate('/statistic/settings');
}

function scrollPopoverToSelectedItem(event: { $el: Framework7Dom }): void {
    scrollToSelectedItem(event.$el, '.popover-inner', 'li.list-item-selected');
}

function onClickPieChartItem(item: Record<string, unknown>): void {
    props.f7router.navigate(getTransactionItemLinkUrl(item['id'] as string));
}

function onClickTrendChartItem(item: { itemId: string, dateRange: TimeRangeAndDateType }): void {
    props.f7router.navigate(getTransactionItemLinkUrl(item.itemId, item.dateRange));
}

function onPageAfterIn(): void {
    if (statisticsStore.transactionStatisticsStateInvalid && !loading.value) {
        reload();
    }

    routeBackOnError(props.f7router, loadingError);
}

init();
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
