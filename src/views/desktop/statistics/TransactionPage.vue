<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="[
                                { name: $t('Categorical Analysis'), value: allAnalysisTypes.CategoricalAnalysis },
                                { name: $t('Trend Analysis'), value: allAnalysisTypes.TrendAnalysis }
                            ]" v-model="queryAnalysisType" />
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ $t('Chart Type') }}</span>
                            <v-select
                                item-title="displayName"
                                item-value="type"
                                class="mt-2"
                                density="compact"
                                :disabled="loading"
                                :items="allChartTypes"
                                v-model="queryChartType"
                            />
                        </div>
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ $t('Sort Order') }}</span>
                            <v-select
                                item-title="displayName"
                                item-value="type"
                                class="mt-2"
                                density="compact"
                                :disabled="loading"
                                :items="allSortingTypes"
                                v-model="querySortingType"
                            />
                        </div>
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="queryChartDataType">
                            <v-tab class="tab-text-truncate" :key="dataType.type" :value="dataType.type"
                                   v-for="dataType in allChartDataTypes" v-show="dataType.isAvailableAnalysisType(queryAnalysisType)">
                                <span class="text-truncate">{{ $t(dataType.name) }}</span>
                                <v-tooltip activator="parent" location="right">{{ $t(dataType.name) }}</v-tooltip>
                            </v-tab>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="statisticsPage">
                                <v-card variant="flat" min-height="680">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="mr-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="icons.menu" size="24" />
                                            </v-btn>
                                            <span>{{ $t('Statistics & Analysis') }}</span>
                                            <v-btn-group class="ml-4" color="default" density="comfortable" variant="outlined" divided>
                                                <v-btn :icon="icons.left"
                                                       :disabled="loading || !canShiftDateRange()"
                                                       @click="shiftDateRange(-1)"/>
                                                <v-menu location="bottom">
                                                    <template #activator="{ props }">
                                                        <v-btn :disabled="loading || queryChartDataType === allChartDataTypes.AccountTotalAssets.type || queryChartDataType === allChartDataTypes.AccountTotalLiabilities.type"
                                                               v-bind="props">{{ dateRangeName() }}</v-btn>
                                                    </template>
                                                    <v-list :selected="[queryDateType]">
                                                        <v-list-item :key="dateRange.type" :value="dateRange.type"
                                                                     :append-icon="(queryDateType === dateRange.type ? icons.check : null)"
                                                                     v-for="dateRange in allDateRangesArray">
                                                            <v-list-item-title
                                                                class="cursor-pointer"
                                                                @click="setDateFilter(dateRange.type)">
                                                                {{ dateRange.displayName }}
                                                                <div class="statistics-custom-datetime-range" v-if="dateRange.type === allDateRanges.Custom.type && showCustomDateRange()">
                                                                    <span>{{ queryStartTime }}</span>
                                                                    <span>&nbsp;-&nbsp;</span>
                                                                    <br/>
                                                                    <span>{{ queryEndTime }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                    </v-list>
                                                </v-menu>
                                                <v-btn :icon="icons.right"
                                                       :disabled="loading || !canShiftDateRange()"
                                                       @click="shiftDateRange(1)"/>
                                            </v-btn-group>

                                            <v-menu location="bottom" v-if="queryAnalysisType === allAnalysisTypes.TrendAnalysis">
                                                <template #activator="{ props }">
                                                    <v-btn class="ml-3" color="default" variant="outlined"
                                                           :prepend-icon="icons.dateAggregation" :disabled="loading"
                                                           v-bind="props">{{ queryTrendDateAggregationTypeName }}</v-btn>
                                                </template>
                                                <v-list>
                                                    <v-list-item class="cursor-pointer" :key="aggregationType.type" :value="aggregationType.type"
                                                                 :append-icon="(trendDateAggregationType === aggregationType.type ? icons.check : null)"
                                                                 :title="aggregationType.displayName"
                                                                 v-for="aggregationType in allDateAggregationTypes"
                                                                 @click="setTrendDateAggregationType(aggregationType.type)">
                                                    </v-list-item>
                                                </v-list>
                                            </v-menu>

                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ml-2" :icon="true" :loading="loading" @click="reload">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="icons.refresh" size="24" />
                                                <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-spacer/>
                                            <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                                                   :disabled="loading" :icon="true">
                                                <v-icon :icon="icons.more" />
                                                <v-menu activator="parent">
                                                    <v-list>
                                                        <v-list-item :prepend-icon="icons.filter"
                                                                     :title="$t('Filter Accounts')"
                                                                     @click="showFilterAccountDialog = true"></v-list-item>
                                                        <v-list-item :prepend-icon="icons.filter"
                                                                     :title="$t('Filter Transaction Categories')"
                                                                     @click="showFilterCategoryDialog = true"></v-list-item>
                                                        <v-list-item :prepend-icon="icons.filter"
                                                                     :title="$t('Filter Transaction Tags')"
                                                                     @click="showFilterTagDialog = true"></v-list-item>
                                                        <v-divider class="my-2"/>
                                                        <v-list-item to="/app/settings?tab=statisticsSetting"
                                                                     :prepend-icon="icons.filterSettings"
                                                                     :title="$t('Settings')"></v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <v-card-text class="statistics-overview-title pt-0" :class="{ 'disabled': loading }"
                                                 v-if="queryAnalysisType === allAnalysisTypes.CategoricalAnalysis && (initing || (categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length))">
                                        <span class="statistics-subtitle">{{ totalAmountName }}</span>
                                        <span class="statistics-overview-amount ml-3"
                                              :class="statisticsTextColor"
                                              v-if="!initing && categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                                            {{ getDisplayAmount(categoricalAnalysisData.totalAmount, defaultCurrency) }}
                                        </span>
                                        <v-skeleton-loader class="skeleton-no-margin ml-3 mb-2"
                                                           width="120px" type="text" :loading="true"
                                                           v-else-if="initing"></v-skeleton-loader>
                                    </v-card-text>

                                    <v-card-text class="statistics-overview-title pt-0"
                                                 v-else-if="!initing && ((queryAnalysisType === allAnalysisTypes.CategoricalAnalysis && (!categoricalAnalysisData || !categoricalAnalysisData.items || !categoricalAnalysisData.items.length))
                                                  || (queryAnalysisType === allAnalysisTypes.TrendAnalysis && (!trendsAnalysisData || !trendsAnalysisData.items || !trendsAnalysisData.items.length)))">
                                        <span class="statistics-subtitle statistics-overview-empty-tip">{{ $t('No transaction data') }}</span>
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === allAnalysisTypes.CategoricalAnalysis && query.categoricalChartType === allCategoricalChartTypes.Pie.type">
                                        <pie-chart
                                            :items="[
                                                {id: '1', name: '---', value: 60, color: '7c7c7f'},
                                                {id: '2', name: '---', value: 20, color: 'a5a5aa'},
                                                {id: '3', name: '---', value: 20, color: 'c5c5c9'}
                                            ]"
                                            :skeleton="true"
                                            id-field="id"
                                            name-field="name"
                                            value-field="value"
                                            color-field="color"
                                            v-if="initing"
                                        ></pie-chart>
                                        <pie-chart
                                            :items="categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length ? categoricalAnalysisData.items : []"
                                            :min-valid-percent="0.0001"
                                            :show-value="showAmountInChart"
                                            :enable-click-item="true"
                                            :default-currency="defaultCurrency"
                                            id-field="id"
                                            name-field="name"
                                            value-field="totalAmount"
                                            percent-field="percent"
                                            hidden-field="hidden"
                                            v-else-if="!initing"
                                            @click="clickPieChartItem"
                                        />
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === allAnalysisTypes.CategoricalAnalysis && query.categoricalChartType === allCategoricalChartTypes.Bar.type">
                                        <v-list rounded lines="two" v-if="initing">
                                            <template :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                                                <v-list-item class="pl-0">
                                                    <template #prepend>
                                                        <div>
                                                            <v-icon class="disabled mr-0" size="34" :icon="icons.square" />
                                                        </div>
                                                    </template>
                                                    <div class="d-flex flex-column ml-2">
                                                        <div class="d-flex">
                                                            <v-skeleton-loader class="skeleton-no-margin my-2"
                                                                               width="120px" type="text" :loading="true"></v-skeleton-loader>
                                                        </div>
                                                        <div>
                                                            <v-progress-linear :model-value="0" :height="4"></v-progress-linear>
                                                        </div>
                                                    </div>
                                                </v-list-item>
                                                <v-divider v-if="itemIdx < 3"/>
                                            </template>
                                        </v-list>
                                        <v-list class="py-0" rounded lines="two" v-else-if="!initing && categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                                            <template :key="idx"
                                                      v-for="(item, idx) in categoricalAnalysisData.items">
                                                <v-list-item class="pl-0" v-if="!item.hidden">
                                                    <template #prepend>
                                                        <router-link class="statistics-list-item" :to="getTransactionItemLinkUrl(item.id)">
                                                            <ItemIcon :icon-type="queryChartDataCategory" size="34px"
                                                                      :icon-id="item.icon"
                                                                      :color="item.color"></ItemIcon>
                                                        </router-link>
                                                    </template>
                                                    <router-link class="statistics-list-item" :to="getTransactionItemLinkUrl(item.id)">
                                                        <div class="d-flex flex-column ml-2">
                                                            <div class="d-flex">
                                                                <span>{{ item.name }}</span>
                                                                <small class="statistics-percent" v-if="item.percent >= 0">{{ getDisplayPercent(item.percent, 2, '&lt;0.01') }}</small>
                                                                <v-spacer/>
                                                                <span class="statistics-amount">{{ getDisplayAmount(item.totalAmount, (item.currency || defaultCurrency)) }}</span>
                                                            </div>
                                                            <div>
                                                                <v-progress-linear :color="item.color ? '#' + item.color : 'primary'"
                                                                                   :bg-color="isDarkMode ? '#161616' : '#f8f8f8'" :bg-opacity="1"
                                                                                   :model-value="item.percent >= 0 ? item.percent : 0"
                                                                                   :height="4"></v-progress-linear>
                                                            </div>
                                                        </div>
                                                    </router-link>
                                                </v-list-item>
                                                <v-divider v-if="!item.hidden && idx !== categoricalAnalysisData.items.length - 1"/>
                                            </template>
                                        </v-list>
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === allAnalysisTypes.TrendAnalysis">
                                        <trends-chart
                                            :type="queryChartType"
                                            :start-year-month="query.trendChartStartYearMonth"
                                            :end-year-month="query.trendChartEndYearMonth"
                                            :items="[]"
                                            :skeleton="true"
                                            id-field="id"
                                            name-field="name"
                                            value-field="value"
                                            color-field="color"
                                            v-if="initing"
                                        ></trends-chart>
                                        <trends-chart
                                            :type="queryChartType"
                                            :start-year-month="query.trendChartStartYearMonth"
                                            :end-year-month="query.trendChartEndYearMonth"
                                            :sorting-type="querySortingType"
                                            :date-aggregation-type="trendDateAggregationType"
                                            :items="trendsAnalysisData && trendsAnalysisData.items && trendsAnalysisData.items.length ? trendsAnalysisData.items : []"
                                            :translate-name="translateNameInTrendsChart"
                                            :show-value="showAmountInChart"
                                            :enable-click-item="true"
                                            :default-currency="defaultCurrency"
                                            :show-total-amount-in-tooltip="showTotalAmountInTrendsChart"
                                            id-field="id"
                                            name-field="name"
                                            value-field="totalAmount"
                                            hidden-field="hidden"
                                            display-orders-field="displayOrders"
                                            v-else-if="!initing && trendsAnalysisData && trendsAnalysisData.items && trendsAnalysisData.items.length"
                                            @click="clickTrendChartItem"
                                        />
                                    </v-card-text>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <date-range-selection-dialog :title="$t('Custom Date Range')"
                                  :min-time="query.categoricalChartStartTime"
                                  :max-time="query.categoricalChartEndTime"
                                  v-model:show="showCustomDateRangeDialog"
                                  @dateRange:change="setCustomDateFilter"
                                  @error="showError" />

    <month-range-selection-dialog :title="$t('Custom Date Range')"
                                  :min-time="query.trendChartStartYearMonth"
                                  :max-time="query.trendChartEndYearMonth"
                                  v-model:show="showCustomMonthRangeDialog"
                                  @dateRange:change="setCustomDateFilter"
                                  @error="showError" />

    <v-dialog width="800" v-model="showFilterAccountDialog">
        <account-filter-settings-card type="statisticsCurrent" :dialog-mode="true"
            @settings:change="setAccountFilter" />
    </v-dialog>

    <v-dialog width="800" v-model="showFilterCategoryDialog">
        <category-filter-settings-card type="statisticsCurrent" :dialog-mode="true"
            @settings:change="setCategoryFilter" />
    </v-dialog>

    <v-dialog width="800" v-model="showFilterTagDialog">
        <transaction-tag-filter-settings-card type="statisticsCurrent" :dialog-mode="true"
                                              @settings:change="setTagFilter" />
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script>
import { useDisplay, useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useStatisticsStore } from '@/stores/statistics.js';

import { DateRangeScene, DateRange } from '@/core/datetime.ts';
import { ThemeType } from '@/core/theme.ts';
import {
    StatisticsAnalysisType,
    CategoricalChartType,
    ChartDataType,
    ChartSortingType,
    ChartDateAggregationType
} from '@/core/statistics.ts';
import {
    isDefined,
    limitText,
    getNameByKeyValue,
    arrayItemToObjectField
} from '@/lib/common.ts'
import { formatPercent } from '@/lib/numeral.ts';
import {
    getYearAndMonthFromUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getShiftedDateRangeAndDateType,
    getDateTypeByDateRange,
    getDateRangeByDateType
} from '@/lib/datetime.ts';

import {
    mdiCheck,
    mdiArrowLeft,
    mdiArrowRight,
    mdiCalendarRangeOutline,
    mdiRefresh,
    mdiSquareRounded,
    mdiMenu,
    mdiFilterOutline,
    mdiFilterCogOutline,
    mdiPencilOutline,
    mdiDotsVertical
} from '@mdi/js';

import AccountFilterSettingsCard from '@/views/desktop/common/cards/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from '@/views/desktop/common/cards/CategoryFilterSettingsCard.vue';
import TransactionTagFilterSettingsCard from '@/views/desktop/common/cards/TransactionTagFilterSettingsCard.vue';

export default {
    components: {
        AccountFilterSettingsCard,
        CategoryFilterSettingsCard,
        TransactionTagFilterSettingsCard
    },
    props: [
        'initAnalysisType',
        'initChartDataType',
        'initChartType',
        'initChartDateType',
        'initStartTime',
        'initEndTime',
        'initFilterAccountIds',
        'initFilterCategoryIds',
        'initTagIds',
        'initTagFilterType',
        'initSortingType',
        'initTrendDateAggregationType'
    ],
    data() {
        const { mdAndUp } = useDisplay();

        return {
            activeTab: 'statisticsPage',
            initing: true,
            loading: true,
            alwaysShowNav: mdAndUp.value,
            showNav: mdAndUp.value,
            analysisType: StatisticsAnalysisType.CategoricalAnalysis,
            trendDateAggregationType: ChartDateAggregationType.Default.type,
            showCustomDateRangeDialog: false,
            showCustomMonthRangeDialog: false,
            showFilterAccountDialog: false,
            showFilterCategoryDialog: false,
            showFilterTagDialog: false,
            icons: {
                check: mdiCheck,
                left: mdiArrowLeft,
                right: mdiArrowRight,
                dateAggregation: mdiCalendarRangeOutline,
                refresh: mdiRefresh,
                square: mdiSquareRounded,
                menu: mdiMenu,
                filter: mdiFilterOutline,
                filterSettings: mdiFilterCogOutline,
                pencil: mdiPencilOutline,
                more: mdiDotsVertical
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useStatisticsStore),
        isDarkMode() {
            return this.globalTheme.global.name.value === ThemeType.Dark;
        },
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        queryAnalysisType: {
            get: function () {
                return this.analysisType;
            },
            set: function(value) {
                this.setAnalysisType(value);
            }
        },
        query() {
            return this.statisticsStore.transactionStatisticsFilter;
        },
        queryChartDataCategory() {
            return this.statisticsStore.categoricalAnalysisChartDataCategory;
        },
        queryChartType: {
            get: function () {
                if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                    return this.query.categoricalChartType;
                } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                    return this.query.trendChartType;
                } else {
                    return null;
                }
            },
            set: function(value) {
                this.setChartType(value);
            }
        },
        queryChartDataType: {
            get: function () {
                return this.query.chartDataType;
            },
            set: function(value) {
                this.setChartDataType(value);
            }
        },
        querySortingType: {
            get: function () {
                return this.query.sortingType;
            },
            set: function(value) {
                this.setSortingType(value);
            }
        },
        queryDateType() {
            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.query.categoricalChartDateType;
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.query.trendChartDateType;
            } else {
                return null;
            }
        },
        queryTrendDateAggregationTypeName() {
            return getNameByKeyValue(this.allDateAggregationTypes, this.trendDateAggregationType, 'type', 'displayName', '');
        },
        queryStartTime() {
            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.categoricalChartStartTime);
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, getYearMonthFirstUnixTime(this.query.trendChartStartYearMonth));
            } else {
                return '';
            }
        },
        queryEndTime() {
            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.categoricalChartEndTime);
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, getYearMonthLastUnixTime(this.query.trendChartEndYearMonth));
            } else {
                return '';
            }
        },
        allAnalysisTypes() {
            return StatisticsAnalysisType;
        },
        allChartTypes() {
            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.getAllCategoricalChartTypes();
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.getAllTrendChartTypes();
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
            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.getAllDateRanges(DateRangeScene.Normal, true);
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.getAllDateRanges(DateRangeScene.TrendAnalysis, true);
            } else {
                return [];
            }
        },
        showAccountBalance() {
            return this.settingsStore.appSettings.showAccountBalance;
        },
        totalAmountName() {
            if (this.queryChartDataType === ChartDataType.IncomeByAccount.type
                || this.queryChartDataType === ChartDataType.IncomeByPrimaryCategory.type
                || this.queryChartDataType === ChartDataType.IncomeBySecondaryCategory.type) {
                return this.$t('Total Income');
            } else if (this.queryChartDataType === ChartDataType.ExpenseByAccount.type
                || this.queryChartDataType === ChartDataType.ExpenseByPrimaryCategory.type
                || this.queryChartDataType === ChartDataType.ExpenseBySecondaryCategory.type) {
                return this.$t('Total Expense');
            } else if (this.queryChartDataType === ChartDataType.AccountTotalAssets.type) {
                return this.$t('Total Assets');
            } else if (this.queryChartDataType === ChartDataType.AccountTotalLiabilities.type) {
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
            return this.queryChartDataType === ChartDataType.TotalExpense.type ||
                this.queryChartDataType === ChartDataType.TotalIncome.type ||
                this.queryChartDataType === ChartDataType.TotalBalance.type;
        },
        showTotalAmountInTrendsChart() {
            return this.queryChartDataType !== ChartDataType.TotalExpense.type &&
                this.queryChartDataType !== ChartDataType.TotalIncome.type &&
                this.queryChartDataType !== ChartDataType.TotalBalance.type;
        },
        statisticsTextColor() {
            if (this.queryChartDataType === ChartDataType.ExpenseByAccount.type ||
                this.queryChartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
                this.queryChartDataType === ChartDataType.ExpenseBySecondaryCategory.type) {
                return 'text-expense';
            } else if (this.queryChartDataType === ChartDataType.IncomeByAccount.type ||
                this.queryChartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
                this.queryChartDataType === ChartDataType.IncomeBySecondaryCategory.type) {
                return 'text-income';
            } else {
                return 'text-default';
            }
        },
        showAmountInChart() {
            if (!this.showAccountBalance
                && (this.queryChartDataType === ChartDataType.AccountTotalAssets.type || this.queryChartDataType === ChartDataType.AccountTotalLiabilities.type)) {
                return false;
            }

            return true;
        }
    },
    watch: {
        'display.mdAndUp.value': function (newValue) {
            this.alwaysShowNav = newValue;

            if (!this.showNav) {
                this.showNav = newValue;
            }
        }
    },
    created() {
        this.init({
            analysisType: this.initAnalysisType,
            chartDataType: this.initChartDataType,
            chartType: this.initChartType,
            chartDateType: this.initChartDateType,
            startTime: this.initStartTime,
            endTime: this.initEndTime,
            filterAccountIds: this.initFilterAccountIds,
            filterCategoryIds: this.initFilterCategoryIds,
            tagIds: this.initTagIds,
            tagFilterType: this.initTagFilterType,
            sortingType: this.initSortingType,
            trendDateAggregationType: this.initTrendDateAggregationType,
        });
    },
    setup() {
        const display = useDisplay();
        const theme = useTheme();

        return {
            display: display,
            globalTheme: theme
        };
    },
    beforeRouteUpdate(to) {
        if (to.query) {
            this.init({
                analysisType: to.query.analysisType,
                chartDataType: to.query.chartDataType,
                chartType: to.query.chartType,
                chartDateType: to.query.chartDateType,
                startTime: to.query.startTime,
                endTime: to.query.endTime,
                filterAccountIds: to.query.filterAccountIds,
                filterCategoryIds: to.query.filterCategoryIds,
                tagIds: to.query.tagIds,
                tagFilterType: to.query.tagFilterType,
                sortingType: to.query.sortingType,
                trendDateAggregationType: to.query.trendDateAggregationType
            });
        } else {
            this.init({});
        }
    },
    methods: {
        init(query) {
            const self = this;
            let needReload = !isDefined(query.analysisType);

            let filter = {
                chartDataType: query.chartDataType ? parseInt(query.chartDataType) : undefined,
                filterAccountIds: query.filterAccountIds ? arrayItemToObjectField(query.filterAccountIds.split(','), true) : {},
                filterCategoryIds: query.filterCategoryIds ? arrayItemToObjectField(query.filterCategoryIds.split(','), true) : {},
                tagIds: query.tagIds,
                tagFilterType: query.tagFilterType && parseInt(query.tagFilterType) >= 0 ? parseInt(query.tagFilterType) : undefined,
                sortingType: query.sortingType ? parseInt(query.sortingType) : undefined
            };

            if (query.analysisType === StatisticsAnalysisType.CategoricalAnalysis.toString()) {
                filter.categoricalChartType = query.chartType ? parseInt(query.chartType) : undefined;
                filter.categoricalChartDateType = query.chartDateType ? parseInt(query.chartDateType) : undefined;
                filter.categoricalChartStartTime = query.startTime ? parseInt(query.startTime) : undefined;
                filter.categoricalChartEndTime = query.endTime ? parseInt(query.endTime) : undefined;

                if (filter.categoricalChartDateType !== self.query.categoricalChartDateType) {
                    needReload = true;
                } else if (filter.categoricalChartDateType === DateRange.Custom.type) {
                    if (filter.categoricalChartStartTime !== self.query.categoricalChartStartTime
                        || filter.categoricalChartEndTime !== self.query.categoricalChartEndTime) {
                        needReload = true;
                    }
                }

                if (query.analysisType !== self.queryAnalysisType.toString()) {
                    self.analysisType = StatisticsAnalysisType.CategoricalAnalysis;
                    needReload = true;
                }
            } else if (query.analysisType === StatisticsAnalysisType.TrendAnalysis.toString()) {
                filter.trendChartType = query.chartType ? parseInt(query.chartType) : undefined;
                filter.trendChartDateType = query.chartDateType ? parseInt(query.chartDateType) : undefined;
                filter.trendChartStartYearMonth = query.startTime;
                filter.trendChartEndYearMonth = query.endTime;

                if (filter.trendChartDateType !== self.query.trendChartDateType) {
                    needReload = true;
                } else if (filter.trendChartDateType === DateRange.Custom.type) {
                    if (filter.trendChartStartYearMonth !== self.query.trendChartStartYearMonth
                        || filter.trendChartEndYearMonth !== self.query.trendChartEndYearMonth) {
                        needReload = true;
                    }
                }

                if (query.analysisType !== self.queryAnalysisType.toString()) {
                    self.analysisType = StatisticsAnalysisType.TrendAnalysis;
                    needReload = true;
                }

                if (query.trendDateAggregationType) {
                    self.trendDateAggregationType = parseInt(query.trendDateAggregationType);
                }
            }

            if (!isDefined(query.analysisType)) {
                self.analysisType = StatisticsAnalysisType.CategoricalAnalysis;
                filter = null;
            }

            self.statisticsStore.initTransactionStatisticsFilter(self.queryAnalysisType, filter);

            if (!needReload && !self.statisticsStore.transactionStatisticsStateInvalid) {
                self.loading = false;
                self.initing = false;
                return;
            }

            Promise.all([
                self.accountsStore.loadAllAccounts({force: false}),
                self.transactionCategoriesStore.loadAllCategories({force: false})
            ]).then(() => {
                if (self.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                    return self.statisticsStore.loadCategoricalAnalysis({
                        force: false
                    });
                } else if (self.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                    return self.statisticsStore.loadTrendAnalysis({
                        force: false
                    });
                }
            }).then(() => {
                self.loading = false;
                self.initing = false;
            }).catch(error => {
                self.loading = false;
                self.initing = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        reload(force) {
            const self = this;
            let dispatchPromise = null;

            self.loading = true;

            if (self.queryChartDataType === ChartDataType.ExpenseByAccount.type ||
                self.queryChartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
                self.queryChartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
                self.queryChartDataType === ChartDataType.IncomeByAccount.type ||
                self.queryChartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
                self.queryChartDataType === ChartDataType.IncomeBySecondaryCategory.type ||
                self.queryChartDataType === ChartDataType.TotalExpense.type ||
                self.queryChartDataType === ChartDataType.TotalIncome.type ||
                self.queryChartDataType === ChartDataType.TotalBalance.type) {
                if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                    dispatchPromise = self.statisticsStore.loadCategoricalAnalysis({
                        force: force
                    });
                } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                    dispatchPromise = self.statisticsStore.loadTrendAnalysis({
                        force: force
                    });
                }
            } else if (self.queryChartDataType === ChartDataType.AccountTotalAssets.type ||
                self.queryChartDataType === ChartDataType.AccountTotalLiabilities.type) {
                dispatchPromise = self.accountsStore.loadAllAccounts({
                    force: force
                });
            }

            if (dispatchPromise) {
                dispatchPromise.then(() => {
                    self.loading = false;

                    if (force) {
                        self.$refs.snackbar.showMessage('Data has been updated');
                    }
                }).catch(error => {
                    self.loading = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            }

            return dispatchPromise;
        },
        setAnalysisType(analysisType) {
            if (this.analysisType === analysisType) {
                return;
            }

            if (!ChartDataType.isAvailableForAnalysisType(this.queryChartDataType, analysisType)) {
                this.statisticsStore.updateTransactionStatisticsFilter({
                    chartDataType: ChartDataType.Default.type
                });
            }

            if (this.analysisType !== StatisticsAnalysisType.TrendAnalysis) {
                this.trendDateAggregationType = ChartDateAggregationType.Month.type;
            }

            this.analysisType = analysisType;
            this.loading = true;
            this.statisticsStore.updateTransactionStatisticsInvalidState(true);
            this.$router.push(this.getFilterLinkUrl());
        },
        setChartType(chartType) {
            let changed = false;

            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartType: chartType
                });
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    trendChartType: chartType
                });
            }

            if (changed) {
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        setChartDataType(chartDataType) {
            const changed = this.statisticsStore.updateTransactionStatisticsFilter({
                chartDataType: chartDataType
            });

            if (changed) {
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        setSortingType(sortingType) {
            if (sortingType < ChartSortingType.Amount.type || sortingType > ChartSortingType.Name.type) {
                return;
            }

            const changed = this.statisticsStore.updateTransactionStatisticsFilter({
                sortingType: sortingType
            });

            if (changed) {
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        setTrendDateAggregationType(aggregationType) {
            const changed = this.trendDateAggregationType !== aggregationType;
            this.trendDateAggregationType = aggregationType;

            if (changed) {
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        setDateFilter(dateType) {
            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                if (dateType === this.allDateRanges.Custom.type) { // Custom
                    this.showCustomDateRangeDialog = true;
                    return;
                } else if (this.query.categoricalChartDateType === dateType) {
                    return;
                }
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                if (dateType === this.allDateRanges.Custom.type) { // Custom
                    this.showCustomMonthRangeDialog = true;
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

            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: dateRange.dateType,
                    categoricalChartStartTime: dateRange.minTime,
                    categoricalChartEndTime: dateRange.maxTime
                });
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    trendChartDateType: dateRange.dateType,
                    trendChartStartYearMonth: getYearAndMonthFromUnixTime(dateRange.minTime),
                    trendChartEndYearMonth: getYearAndMonthFromUnixTime(dateRange.maxTime)
                });
            }

            if (changed) {
                this.loading = true;
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        setCustomDateFilter(startTime, endTime) {
            if (!startTime || !endTime) {
                return;
            }

            let changed = false;

            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                const chartDateType = getDateTypeByDateRange(startTime, endTime, this.firstDayOfWeek, DateRangeScene.Normal);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: chartDateType,
                    categoricalChartStartTime: startTime,
                    categoricalChartEndTime: endTime
                });

                this.showCustomDateRangeDialog = false;
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                const chartDateType = getDateTypeByDateRange(getYearMonthFirstUnixTime(startTime), getYearMonthLastUnixTime(endTime), this.firstDayOfWeek, DateRangeScene.TrendAnalysis);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    trendChartDateType: chartDateType,
                    trendChartStartYearMonth: startTime,
                    trendChartEndYearMonth: endTime
                });

                this.showCustomMonthRangeDialog = false;
            }

            if (changed) {
                this.loading = true;
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        showCustomDateRange() {
            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.query.categoricalChartDateType === this.allDateRanges.Custom.type && this.query.categoricalChartStartTime && this.query.categoricalChartEndTime;
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.query.trendChartDateType === this.allDateRanges.Custom.type && this.query.trendChartStartYearMonth && this.query.trendChartEndYearMonth;
            } else {
                return false;
            }
        },
        canShiftDateRange() {
            if (this.queryChartDataType === ChartDataType.AccountTotalAssets.type || this.queryChartDataType === ChartDataType.AccountTotalLiabilities.type) {
                return false;
            }

            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.query.categoricalChartDateType !== this.allDateRanges.All.type;
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.query.trendChartDateType !== this.allDateRanges.All.type;
            } else {
                return false;
            }
        },
        shiftDateRange(scale) {
            if (this.query.categoricalChartDateType === this.allDateRanges.All.type) {
                return;
            }

            let changed = false;

            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                const newDateRange = getShiftedDateRangeAndDateType(this.query.categoricalChartStartTime, this.query.categoricalChartEndTime, scale, this.firstDayOfWeek, DateRangeScene.Normal);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    categoricalChartDateType: newDateRange.dateType,
                    categoricalChartStartTime: newDateRange.minTime,
                    categoricalChartEndTime: newDateRange.maxTime
                });
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                const newDateRange = getShiftedDateRangeAndDateType(getYearMonthFirstUnixTime(this.query.trendChartStartYearMonth), getYearMonthLastUnixTime(this.query.trendChartEndYearMonth), scale, this.firstDayOfWeek, DateRangeScene.TrendAnalysis);

                changed = this.statisticsStore.updateTransactionStatisticsFilter({
                    trendChartDateType: newDateRange.dateType,
                    trendChartStartYearMonth: getYearAndMonthFromUnixTime(newDateRange.minTime),
                    trendChartEndYearMonth: getYearAndMonthFromUnixTime(newDateRange.maxTime)
                });
            }

            if (changed) {
                this.loading = true;
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        dateRangeName() {
            if (this.queryChartDataType === ChartDataType.AccountTotalAssets.type ||
                this.queryChartDataType === ChartDataType.AccountTotalLiabilities.type) {
                return this.$t(this.allDateRanges.All.name);
            }

            if (this.queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis) {
                return this.$locale.getDateRangeDisplayName(this.userStore, this.query.categoricalChartDateType, this.query.categoricalChartStartTime, this.query.categoricalChartEndTime);
            } else if (this.queryAnalysisType === StatisticsAnalysisType.TrendAnalysis) {
                return this.$locale.getDateRangeDisplayName(this.userStore, this.query.trendChartDateType, getYearMonthFirstUnixTime(this.query.trendChartStartYearMonth), getYearMonthLastUnixTime(this.query.trendChartEndYearMonth));
            } else {
                return '';
            }
        },
        setAccountFilter(changed) {
            this.showFilterAccountDialog = false;

            if (changed) {
                this.loading = true;
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        setCategoryFilter(changed) {
            this.showFilterCategoryDialog = false;

            if (changed) {
                this.loading = true;
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        setTagFilter(changed) {
            this.showFilterTagDialog = false;

            if (changed) {
                this.loading = true;
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        clickPieChartItem(item) {
            this.$router.push(this.getTransactionItemLinkUrl(item.id));
        },
        clickTrendChartItem(item) {
            this.$router.push(this.getTransactionItemLinkUrl(item.itemId, item.dateRange));
        },
        showError(message) {
            this.$refs.snackbar.showError(message);
        },
        getDisplayAmount(amount, currency, textLimit) {
            amount = this.getDisplayCurrency(amount, currency);

            if (!this.showAccountBalance
                && (this.queryChartDataType === ChartDataType.AccountTotalAssets.type
                    || this.queryChartDataType === ChartDataType.AccountTotalLiabilities.type)
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
        getFilterLinkUrl() {
            return `/statistics/transaction?${this.statisticsStore.getTransactionStatisticsPageParams(this.queryAnalysisType, this.trendDateAggregationType)}`;
        },
        getTransactionItemLinkUrl(itemId, dateRange) {
            return `/transaction/list?${this.statisticsStore.getTransactionListPageParams(this.queryAnalysisType, itemId, dateRange)}`;
        }
    }
}
</script>

<style>
.statistics-custom-datetime-range {
    font-size: 0.7rem;
    line-height: 1rem;
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
}

.statistics-overview-title {
    line-height: 2rem !important;
    height: 46px;
    display: flex;
    align-items: flex-end;
}

.statistics-overview-amount {
    font-size: 1.5rem;
    overflow: hidden;
    text-overflow: ellipsis;
}

.statistics-subtitle {
    font-size: 1rem;
    line-height: 1.75rem
}

.statistics-overview-empty-tip {
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
}

.statistics-list-item {
    color: var(--v-theme-on-default);
    font-size: 1rem !important;
    line-height: 1.75rem;
    overflow: hidden;
    text-overflow: ellipsis;
}

.statistics-list-item .statistics-percent {
    font-size: 0.75rem;
    opacity: 0.7;
    margin-left: 6px;
}

.statistics-list-item .statistics-amount {
    opacity: 0.8;
}
</style>
