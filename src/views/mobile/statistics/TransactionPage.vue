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
                    v-model:opened="showChartDataTypePopover">
            <f7-list dividers>
                <f7-list-item group-title :title="$t('Categorical Analysis')" />
                <f7-list-item :title="$t(dataType.name)"
                              :key="dataType.type"
                              v-for="dataType in allChartDataTypes"
                              v-show="dataType.availableAnalysisTypes[analysisType]"
                              @click="setChartDataType(dataType.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.chartDataType === dataType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-card v-if="query.categoricalChartType === allCategoricalChartTypes.Pie">
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

        <f7-card v-else-if="query.categoricalChartType === allCategoricalChartTypes.Bar">
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
                                  :link="getTransactionItemLinkUrl(item)"
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

        <f7-popover class="sorting-type-popover-menu"
                    v-model:opened="showSortingTypePopover">
            <f7-list dividers>
                <f7-list-item :title="$t(sortingType.name)"
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
            <f7-link :class="{ 'disabled': reloading || query.categoricalChartDateType === allDateRanges.All.type || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type }" @click="shiftDateRange(-1)">
                <f7-icon f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': reloading || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type }" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.maxTime > 0 || query.minTime > 0 }">{{ dateRangeName() }}</span>
            </f7-link>
            <f7-link :class="{ 'disabled': reloading || query.categoricalChartDateType === allDateRanges.All.type || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type }" @click="shiftDateRange(1)">
                <f7-icon f7="arrow_right_square"></f7-icon>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" :key="chartType.type"
                     v-for="chartType in allChartTypes" @click="setChartType(chartType.type)">
                <span :class="{ 'tabbar-item-changed': query.categoricalChartType === chartType.type }">{{ chartType.displayName }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-popover class="date-popover-menu"
                    v-model:opened="showDatePopover"
                    @popover:open="scrollPopoverToSelectedItem">
            <f7-list dividers>
                <f7-list-item :title="dateRange.displayName"
                              :class="{ 'list-item-selected': query.categoricalChartDateType === dateRange.type }"
                              :key="dateRange.type"
                              v-for="dateRange in allDateRangesArray"
                              @click="setDateFilter(dateRange.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.categoricalChartDateType === dateRange.type"></f7-icon>
                    </template>
                    <template #footer>
                        <div v-if="dateRange.type === allDateRanges.Custom.type && query.categoricalChartDateType === allDateRanges.Custom.type && query.categoricalChartStartTime && query.categoricalChartEndTime">
                            <span>{{ queryStartTime }}</span>
                            <span>&nbsp;-&nbsp;</span>
                            <br/>
                            <span>{{ queryEndTime }}</span>
                        </div>
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
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useStatisticsStore } from '@/stores/statistics.js';

import datetimeConstants from '@/consts/datetime.js';
import statisticsConstants from '@/consts/statistics.js';
import { getNameByKeyValue, limitText } from '@/lib/common.js';
import { formatPercent } from '@/lib/numeral.js';
import {
    getShiftedDateRangeAndDateType,
    getDateTypeByDateRange,
    getDateRangeByDateType
} from '@/lib/datetime.js';
import { scrollToSelectedItem } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            loading: true,
            loadingError: null,
            reloading: false,
            analysisType: statisticsConstants.allAnalysisTypes.CategoricalAnalysis,
            showChartDataTypePopover: false,
            showSortingTypePopover: false,
            showDatePopover: false,
            showCustomDateRangeSheet: false,
            showMoreActionSheet: false
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useStatisticsStore),
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
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
        queryChartDataTypeName() {
            const queryChartDataTypeName = getNameByKeyValue(this.allChartDataTypes, this.query.chartDataType, 'type', 'name', 'Statistics');
            return this.$t(queryChartDataTypeName);
        },
        querySortingTypeName() {
            const querySortingTypeName = getNameByKeyValue(this.allSortingTypes, this.query.sortingType, 'type', 'name', 'System Default');
            return this.$t(querySortingTypeName);
        },
        queryStartTime() {
            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.categoricalChartStartTime);
            } else {
                return '';
            }
        },
        queryEndTime() {
            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.categoricalChartEndTime);
            } else {
                return '';
            }
        },
        allChartTypes() {
            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                return this.$locale.getAllCategoricalChartTypes();
            } else {
                return [];
            }
        },
        allCategoricalChartTypes() {
            return statisticsConstants.allCategoricalChartTypes;
        },
        allChartDataTypes() {
            return statisticsConstants.allChartDataTypes;
        },
        allSortingTypes() {
            return statisticsConstants.allSortingTypes;
        },
        allDateRanges() {
            return datetimeConstants.allDateRanges;
        },
        allDateRangesArray() {
            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                return this.$locale.getAllDateRanges(datetimeConstants.allDateRangeScenes.Normal, true);
            } else {
                return [];
            }
        },
        showAccountBalance() {
            return this.settingsStore.appSettings.showAccountBalance;
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
        categoricalAnalysisData() {
            return this.statisticsStore.categoricalAnalysisData;
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

        self.statisticsStore.initTransactionStatisticsFilter(self.analysisType);

        Promise.all([
            self.accountsStore.loadAllAccounts({ force: false }),
            self.transactionCategoriesStore.loadAllCategories({ force: false })
        ]).then(() => {
            return self.statisticsStore.loadCategoricalAnalysis({
                force: false
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

            if (self.query.chartDataType === self.allChartDataTypes.ExpenseByAccount.type ||
                self.query.chartDataType === self.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.ExpenseBySecondaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeByAccount.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeByPrimaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeBySecondaryCategory.type) {
                dispatchPromise = self.statisticsStore.loadCategoricalAnalysis({
                    force: force
                });
            } else if (self.query.chartDataType === self.allChartDataTypes.AccountTotalAssets.type ||
                self.query.chartDataType === self.allChartDataTypes.AccountTotalLiabilities.type) {
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
            this.statisticsStore.updateTransactionStatisticsFilter({
                categoricalChartType: chartType
            });
        },
        setChartDataType(chartDataType) {
            this.statisticsStore.updateTransactionStatisticsFilter({
                chartDataType: chartDataType
            });
            this.showChartDataTypePopover = false;
        },
        setSortingType(sortingType) {
            if (sortingType < this.allSortingTypes.Amount.type || sortingType > this.allSortingTypes.Name.type) {
                this.showSortingTypePopover = false;
                return;
            }

            this.statisticsStore.updateTransactionStatisticsFilter({
                sortingType: sortingType
            });

            this.showSortingTypePopover = false;
        },
        setDateFilter(dateType) {
            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                if (dateType === this.allDateRanges.Custom.type) { // Custom
                    this.showCustomDateRangeSheet = true;
                    this.showDatePopover = false;
                    return;
                } else if (this.query.categoricalChartDateType === dateType) {
                    return;
                }
            }

            const dateRange = getDateRangeByDateType(dateType, this.firstDayOfWeek);

            if (!dateRange) {
                return;
            }

            let changed = false;

            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: dateRange.dateType,
                    categoricalChartStartTime: dateRange.minTime,
                    categoricalChartEndTime: dateRange.maxTime
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

            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                const chartDateType = getDateTypeByDateRange(startTime, endTime, this.firstDayOfWeek, datetimeConstants.allDateRangeScenes.Normal);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: chartDateType,
                    categoricalChartStartTime: startTime,
                    categoricalChartEndTime: endTime
                });

                this.showCustomDateRangeSheet = false;
            }

            if (changed) {
                this.reload(null);
            }
        },
        shiftDateRange(scale) {
            if (this.query.categoricalChartDateType === this.allDateRanges.All.type) {
                return;
            }

            let changed = false;

            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                const newDateRange = getShiftedDateRangeAndDateType(this.query.categoricalChartStartTime, this.query.categoricalChartEndTime, scale, this.firstDayOfWeek, datetimeConstants.allDateRangeScenes.Normal);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: newDateRange.dateType,
                    categoricalChartStartTime: newDateRange.minTime,
                    categoricalChartEndTime: newDateRange.maxTime
                });
            }

            if (changed) {
                this.reload(null);
            }
        },
        dateRangeName() {
            if (this.query.chartDataType === this.allChartDataTypes.AccountTotalAssets.type ||
                this.query.chartDataType === this.allChartDataTypes.AccountTotalLiabilities.type) {
                return this.$t(this.allDateRanges.All.name);
            }

            if (this.analysisType === statisticsConstants.allAnalysisTypes.CategoricalAnalysis) {
                return this.$locale.getDateRangeDisplayName(this.userStore, this.query.categoricalChartDateType, this.query.categoricalChartStartTime, this.query.categoricalChartEndTime);
            } else {
                return '';
            }
        },
        clickPieChartItem(item) {
            this.f7router.navigate(this.getTransactionItemLinkUrl(item));
        },
        filterAccounts() {
            this.f7router.navigate('/settings/filter/account?type=statisticsCurrent');
        },
        filterCategories() {
            this.f7router.navigate('/settings/filter/category?type=statisticsCurrent');
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
                && (this.query.chartDataType === this.allChartDataTypes.AccountTotalAssets.type
                    || this.query.chartDataType === this.allChartDataTypes.AccountTotalLiabilities.type)
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
        getTransactionItemLinkUrl(item) {
            return `/transaction/list?${this.statisticsStore.getTransactionListPageParams(statisticsConstants.allAnalysisTypes.CategoricalAnalysis, item)}`;
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
