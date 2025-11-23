<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="[
                                { name: tt('Categorical Analysis'), value: StatisticsAnalysisType.CategoricalAnalysis },
                                { name: tt('Trend Analysis'), value: StatisticsAnalysisType.TrendAnalysis },
                                { name: tt('Asset Trends'), value: StatisticsAnalysisType.AssetTrends }
                            ]" v-model="queryAnalysisType" />
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ tt('Chart Type') }}</span>
                            <v-select
                                item-title="displayName"
                                item-value="type"
                                class="mt-2"
                                density="compact"
                                :disabled="loading"
                                :items="allChartTypes"
                                v-model="queryChartType"
                                v-show="!isQuerySpecialChartType"
                            />
                            <v-select
                                item-title="displayName"
                                item-value="type"
                                class="mt-2"
                                density="compact"
                                :disabled="true"
                                :items="[{ displayName: tt('Sankey Chart'), type: 0 }]"
                                :model-value="0"
                                v-show="isQuerySpecialChartType && queryChartDataType === ChartDataType.Overview.type"
                            />
                        </div>
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ tt('Sort Order') }}</span>
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
                                   v-for="dataType in ChartDataType.values(undefined, true)"
                                   v-show="dataType.isAvailableAnalysisType(queryAnalysisType)">
                                <span class="text-truncate">{{ tt(dataType.name) }}</span>
                                <v-tooltip activator="parent" location="right">{{ tt(dataType.name) }}</v-tooltip>
                            </v-tab>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="statisticsPage">
                                <v-card variant="flat" :min-height="queryAnalysisType === StatisticsAnalysisType.TrendAnalysis || queryAnalysisType === StatisticsAnalysisType.AssetTrends ? '900' : '800'">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="me-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="mdiMenu" size="24" />
                                            </v-btn>
                                            <span>{{ tt('Statistics & Analysis') }}</span>
                                            <v-btn-group class="ms-4" color="default" density="comfortable" variant="outlined" divided>
                                                <v-btn class="button-icon-with-direction" :icon="mdiArrowLeft"
                                                       :disabled="loading || !canShiftDateRange"
                                                       @click="shiftDateRange(-1)"/>
                                                <v-menu location="bottom" max-height="500">
                                                    <template #activator="{ props }">
                                                        <v-btn :disabled="loading || !canChangeDateRange"
                                                               v-bind="props">{{ queryDateRangeName }}</v-btn>
                                                    </template>
                                                    <v-list :selected="[queryDateType]">
                                                        <v-list-item :key="dateRange.type" :value="dateRange.type"
                                                                     :append-icon="(queryDateType === dateRange.type ? mdiCheck : undefined)"
                                                                     v-for="dateRange in allDateRanges">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="setDateFilter(dateRange.type)">
                                                                <div class="d-flex align-center">
                                                                    <span>{{ dateRange.displayName }}</span>
                                                                </div>
                                                                <div class="statistics-custom-datetime-range smaller" v-if="dateRange.isUserCustomRange && canShowCustomDateRange(dateRange.type)">
                                                                    <span>{{ queryStartTime }}</span>
                                                                    <span>&nbsp;-&nbsp;</span>
                                                                    <br/>
                                                                    <span>{{ queryEndTime }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                    </v-list>
                                                </v-menu>
                                                <v-btn class="button-icon-with-direction" :icon="mdiArrowRight"
                                                       :disabled="loading || !canShiftDateRange"
                                                       @click="shiftDateRange(1)"/>
                                            </v-btn-group>

                                            <v-menu location="bottom" max-height="500" v-if="queryAnalysisType === StatisticsAnalysisType.TrendAnalysis">
                                                <template #activator="{ props }">
                                                    <v-btn class="ms-3" color="default" variant="outlined"
                                                           :prepend-icon="mdiCalendarRangeOutline" :disabled="loading"
                                                           v-bind="props">{{ queryTrendDateAggregationTypeName }}</v-btn>
                                                </template>
                                                <v-list>
                                                    <v-list-item class="cursor-pointer" :key="aggregationType.type" :value="aggregationType.type"
                                                                 :append-icon="(trendDateAggregationType === aggregationType.type ? mdiCheck : undefined)"
                                                                 :title="aggregationType.displayName"
                                                                 v-for="aggregationType in allTrendAnalysisDateAggregationTypes"
                                                                 @click="setTrendDateAggregationType(aggregationType.type)">
                                                    </v-list-item>
                                                </v-list>
                                            </v-menu>

                                            <v-menu location="bottom" max-height="500" v-if="queryAnalysisType === StatisticsAnalysisType.AssetTrends">
                                                <template #activator="{ props }">
                                                    <v-btn class="ms-3" color="default" variant="outlined"
                                                           :prepend-icon="mdiCalendarRangeOutline" :disabled="loading"
                                                           v-bind="props">{{ queryAssetTrendsDateAggregationTypeName }}</v-btn>
                                                </template>
                                                <v-list>
                                                    <v-list-item class="cursor-pointer" :key="aggregationType.type" :value="aggregationType.type"
                                                                 :append-icon="(assetTrendsDateAggregationType === aggregationType.type ? mdiCheck : undefined)"
                                                                 :title="aggregationType.displayName"
                                                                 v-for="aggregationType in allAssetTrendsDateAggregationTypes"
                                                                 @click="setAssetTrendsDateAggregationType(aggregationType.type)">
                                                    </v-list-item>
                                                </v-list>
                                            </v-menu>

                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ms-2" :icon="true" :loading="loading" @click="reload(true)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="mdiRefresh" size="24" />
                                                <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-spacer/>
                                            <div class="transaction-keyword-filter ms-2">
                                                <v-text-field density="compact" :disabled="loading"
                                                              :prepend-inner-icon="mdiMagnify"
                                                              :append-inner-icon="filterKeyword !== query.keyword ? mdiCheck : undefined"
                                                              :placeholder="tt('Filter transaction description')"
                                                              v-model="filterKeyword"
                                                              v-if="canUseKeywordFilter"
                                                              @click:append-inner="setKeywordFilter(filterKeyword)"
                                                              @keyup.enter="setKeywordFilter(filterKeyword)"
                                                />
                                            </div>
                                            <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                                                   :disabled="loading" :icon="true">
                                                <v-icon :icon="mdiDotsVertical" />
                                                <v-menu activator="parent">
                                                    <v-list>
                                                        <v-list-item :disabled="loading"
                                                                     :prepend-icon="mdiFilterOutline"
                                                                     :title="tt('Filter Accounts')"
                                                                     @click="showFilterAccountDialog = true"></v-list-item>
                                                        <v-list-item :disabled="loading"
                                                                     :prepend-icon="mdiFilterOutline"
                                                                     :title="tt('Filter Transaction Categories')"
                                                                     @click="showFilterCategoryDialog = true"
                                                                     v-if="canUseCategoryFilter"></v-list-item>
                                                        <v-list-item :disabled="loading"
                                                                     :prepend-icon="mdiFilterOutline"
                                                                     :title="tt('Filter Transaction Tags')"
                                                                     @click="showFilterTagDialog = true"
                                                                     v-if="canUseTagFilter"></v-list-item>
                                                        <v-divider class="my-2" v-if="!isQuerySpecialChartType" />
                                                        <v-list-item :prepend-icon="mdiExport"
                                                                     :title="tt('Export Results')"
                                                                     :disabled="!statisticsDataHasData"
                                                                     @click="exportResults"
                                                                     v-if="!isQuerySpecialChartType"></v-list-item>
                                                        <v-divider class="my-2"/>
                                                        <v-list-item to="/app/settings?tab=statisticsSetting"
                                                                     :prepend-icon="mdiFilterCogOutline"
                                                                     :title="tt('Settings')"></v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <v-card-text class="statistics-overview-title pt-0" :class="{ 'disabled': loading }"
                                                 v-if="queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis && isQuerySpecialChartType && queryChartDataType === ChartDataType.Overview.type && (initing || categoricalOverviewAnalysisData && categoricalOverviewAnalysisData.items && categoricalOverviewAnalysisData.items.length)">
                                        <span class="statistics-subtitle">{{ tt('Total Income') }}</span>
                                        <span class="statistics-overview-amount ms-3 text-income"
                                              v-if="!initing && categoricalOverviewAnalysisData && categoricalOverviewAnalysisData.items && categoricalOverviewAnalysisData.items.length">
                                            {{ getDisplayAmount(categoricalOverviewAnalysisData.totalIncome, defaultCurrency) }}
                                        </span>
                                        <v-skeleton-loader class="skeleton-no-margin ms-3 mb-2"
                                                           width="120px" type="text" :loading="true"
                                                           v-else-if="initing"></v-skeleton-loader>
                                        <span class="statistics-subtitle ms-3">{{ tt('Total Expense') }}</span>
                                        <span class="statistics-overview-amount ms-3 text-expense"
                                              v-if="!initing && categoricalOverviewAnalysisData && categoricalOverviewAnalysisData.items && categoricalOverviewAnalysisData.items.length">
                                            {{ getDisplayAmount(categoricalOverviewAnalysisData.totalExpense, defaultCurrency) }}
                                        </span>
                                        <v-skeleton-loader class="skeleton-no-margin ms-3 mb-2"
                                                           width="120px" type="text" :loading="true"
                                                           v-else-if="initing"></v-skeleton-loader>
                                    </v-card-text>

                                    <v-card-text class="statistics-overview-title pt-0" :class="{ 'disabled': loading }"
                                                 v-else-if="queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis && !isQuerySpecialChartType && (initing || (categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length))">
                                        <span class="statistics-subtitle">{{ totalAmountName }}</span>
                                        <span class="statistics-overview-amount ms-3"
                                              :class="statisticsTextColor"
                                              v-if="!initing && categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length">
                                            {{ getDisplayAmount(categoricalAnalysisData.totalAmount, defaultCurrency) }}
                                        </span>
                                        <v-skeleton-loader class="skeleton-no-margin ms-3 mb-2"
                                                           width="120px" type="text" :loading="true"
                                                           v-else-if="initing"></v-skeleton-loader>
                                    </v-card-text>

                                    <v-card-text class="statistics-overview-title pt-0"
                                                 v-else-if="!loading && (
                                                     (queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis && isQuerySpecialChartType && queryChartDataType === ChartDataType.Overview.type && (!categoricalOverviewAnalysisData || !categoricalOverviewAnalysisData.items || !categoricalOverviewAnalysisData.items.length))
                                                  || (queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis && !isQuerySpecialChartType && (!categoricalAnalysisData || !categoricalAnalysisData.items || !categoricalAnalysisData.items.length))
                                                  || (queryAnalysisType === StatisticsAnalysisType.TrendAnalysis && (!trendsAnalysisData || !trendsAnalysisData.items || !trendsAnalysisData.items.length))
                                                  || (queryAnalysisType === StatisticsAnalysisType.AssetTrends && (!assetTrendsData || !assetTrendsData.items || !assetTrendsData.items.length))
                                                  )">
                                        <span class="statistics-subtitle statistics-overview-empty-tip">{{ tt('No transaction data') }}</span>
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis && queryChartDataType === ChartDataType.Overview.type">
                                        <account-and-category-sankey-chart
                                            :items="[]"
                                            :sorting-type="querySortingType"
                                            :skeleton="true"
                                            v-if="initing"
                                        />
                                        <account-and-category-sankey-chart
                                            :items="categoricalOverviewAnalysisData && categoricalOverviewAnalysisData.items && categoricalOverviewAnalysisData.items.length ? categoricalOverviewAnalysisData.items : []"
                                            :enable-click-item="true"
                                            :default-currency="defaultCurrency"
                                            v-else-if="!initing"
                                            @click="onClickSankeyChartItem"
                                        />
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis && !isQuerySpecialChartType && query.categoricalChartType === CategoricalChartType.Pie.type">
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
                                        />
                                        <pie-chart
                                            :items="categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length ? categoricalAnalysisData.items : []"
                                            :min-valid-percent="0.0001"
                                            :show-value="showAmountInChart"
                                            :show-percent="showPercentInCategoricalChart"
                                            :enable-click-item="true"
                                            :default-currency="defaultCurrency"
                                            id-field="id"
                                            name-field="name"
                                            value-field="totalAmount"
                                            percent-field="percent"
                                            hidden-field="hidden"
                                            v-else-if="!initing"
                                            @click="onClickPieChartItem"
                                        />
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis && !isQuerySpecialChartType && query.categoricalChartType === CategoricalChartType.Bar.type">
                                        <v-list rounded lines="two" v-if="initing">
                                            <template :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                                                <v-list-item class="ps-0">
                                                    <template #prepend>
                                                        <div>
                                                            <v-icon class="disabled me-0" size="34" :icon="mdiSquareRounded" />
                                                        </div>
                                                    </template>
                                                    <div class="d-flex flex-column ms-2">
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
                                                <v-list-item class="ps-0" v-if="!item.hidden">
                                                    <template #prepend>
                                                        <router-link class="statistics-list-item" :to="getTransactionItemLinkUrl(item.id)">
                                                            <ItemIcon :icon-type="queryChartDataCategory" size="34px"
                                                                      :icon-id="item.icon"
                                                                      :color="item.color"></ItemIcon>
                                                        </router-link>
                                                    </template>
                                                    <router-link class="statistics-list-item" :to="getTransactionItemLinkUrl(item.id)">
                                                        <div class="d-flex flex-column ms-2">
                                                            <div class="d-flex">
                                                                <span>{{ item.name }}</span>
                                                                <small class="statistics-percent" v-if="showPercentInCategoricalChart && item.percent >= 0">{{ formatPercentToLocalizedNumerals(item.percent, 2, '&lt;0.01') }}</small>
                                                                <v-spacer/>
                                                                <span class="statistics-amount">{{ getDisplayAmount(item.totalAmount, defaultCurrency) }}</span>
                                                            </div>
                                                            <div>
                                                                <v-progress-linear :color="item.color ? getTransactionCategoricalAnalysisDataItemDisplayColor(item) : 'primary'"
                                                                                   :bg-color="isDarkMode ? '#444444' : '#f8f8f8'" :bg-opacity="1"
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

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === StatisticsAnalysisType.CategoricalAnalysis && !isQuerySpecialChartType && query.categoricalChartType === CategoricalChartType.Radar.type">
                                        <radar-chart
                                            :items="[
                                                {name: '---', value: 10},
                                                {name: '---', value: 10},
                                                {name: '---', value: 10},
                                                {name: '---', value: 10},
                                                {name: '---', value: 10},
                                                {name: '---', value: 10}
                                            ]"
                                            :skeleton="true"
                                            name-field="name"
                                            value-field="value"
                                            v-if="initing"
                                        />
                                        <radar-chart
                                            :items="categoricalAnalysisData && categoricalAnalysisData.items && categoricalAnalysisData.items.length ? categoricalAnalysisData.items : []"
                                            :min-valid-percent="0.0001"
                                            :show-value="showAmountInChart"
                                            :show-percent="showPercentInCategoricalChart"
                                            :default-currency="defaultCurrency"
                                            name-field="name"
                                            value-field="totalAmount"
                                            percent-field="percent"
                                            hidden-field="hidden"
                                            v-else-if="!initing"
                                        />
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === StatisticsAnalysisType.TrendAnalysis">
                                        <trends-chart
                                            chart-mode="monthly"
                                            :type="queryChartType"
                                            :start-time="undefined"
                                            :end-time="undefined"
                                            :start-year-month="query.trendChartStartYearMonth"
                                            :end-year-month="query.trendChartEndYearMonth"
                                            :sorting-type="querySortingType"
                                            :data-aggregation-type="ChartDataAggregationType.Sum"
                                            :date-aggregation-type="trendDateAggregationType"
                                            :fiscal-year-start="fiscalYearStart"
                                            :items="[]"
                                            :skeleton="true"
                                            id-field="id"
                                            name-field="name"
                                            value-field="value"
                                            color-field="color"
                                            v-if="initing"
                                        />
                                        <trends-chart
                                            chart-mode="monthly"
                                            :type="queryChartType"
                                            :start-time="undefined"
                                            :end-time="undefined"
                                            :start-year-month="query.trendChartStartYearMonth"
                                            :end-year-month="query.trendChartEndYearMonth"
                                            :sorting-type="querySortingType"
                                            :data-aggregation-type="ChartDataAggregationType.Sum"
                                            :date-aggregation-type="trendDateAggregationType"
                                            :fiscal-year-start="fiscalYearStart"
                                            :items="trendsAnalysisData && trendsAnalysisData.items && trendsAnalysisData.items.length ? trendsAnalysisData.items : []"
                                            :translate-name="translateNameInTrendsChart"
                                            :show-value="showAmountInChart"
                                            :enable-click-item="true"
                                            :default-currency="defaultCurrency"
                                            :stacked="showStackedInTrendsChart"
                                            :show-total-amount-in-tooltip="showTotalAmountInTrendsChart"
                                            ref="monthlyTrendsChart"
                                            id-field="id"
                                            name-field="name"
                                            value-field="totalAmount"
                                            hidden-field="hidden"
                                            display-orders-field="displayOrders"
                                            v-else-if="!initing && trendsAnalysisData && trendsAnalysisData.items && trendsAnalysisData.items.length"
                                            @click="onClickTrendChartItem"
                                        />
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="queryAnalysisType === StatisticsAnalysisType.AssetTrends">
                                        <trends-chart
                                            chart-mode="daily"
                                            :type="queryChartType"
                                            :start-time="query.assetTrendsChartStartTime"
                                            :end-time="query.assetTrendsChartEndTime"
                                            :start-year-month="undefined"
                                            :end-year-month="undefined"
                                            :sorting-type="querySortingType"
                                            :data-aggregation-type="ChartDataAggregationType.Last"
                                            :date-aggregation-type="assetTrendsDateAggregationType"
                                            :fiscal-year-start="fiscalYearStart"
                                            :items="[]"
                                            :skeleton="true"
                                            id-field="id"
                                            name-field="name"
                                            value-field="value"
                                            color-field="color"
                                            v-if="initing"
                                        />
                                        <trends-chart
                                            chart-mode="daily"
                                            :type="queryChartType"
                                            :start-time="query.assetTrendsChartStartTime"
                                            :end-time="query.assetTrendsChartEndTime"
                                            :start-year-month="undefined"
                                            :end-year-month="undefined"
                                            :sorting-type="querySortingType"
                                            :data-aggregation-type="ChartDataAggregationType.Last"
                                            :date-aggregation-type="assetTrendsDateAggregationType"
                                            :fiscal-year-start="fiscalYearStart"
                                            :items="assetTrendsData && assetTrendsData.items && assetTrendsData.items.length ? assetTrendsData.items : []"
                                            :translate-name="translateNameInTrendsChart"
                                            :show-value="showAmountInChart"
                                            :enable-click-item="true"
                                            :default-currency="defaultCurrency"
                                            :stacked="showStackedInTrendsChart"
                                            :show-total-amount-in-tooltip="showTotalAmountInTrendsChart"
                                            ref="dailyTrendsChart"
                                            id-field="id"
                                            name-field="name"
                                            value-field="totalAmount"
                                            hidden-field="hidden"
                                            display-orders-field="displayOrders"
                                            v-else-if="!initing && assetTrendsData && assetTrendsData.items && assetTrendsData.items.length"
                                            @click="onClickTrendChartItem"
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

    <date-range-selection-dialog :title="tt('Custom Date Range')"
                                  :min-time="query.categoricalChartStartTime"
                                  :max-time="query.categoricalChartEndTime"
                                  v-model:show="showCustomDateRangeDialog"
                                  @dateRange:change="setCustomDateFilter"
                                  @error="onShowDateRangeError" />

    <month-range-selection-dialog :title="tt('Custom Date Range')"
                                  :min-time="query.trendChartStartYearMonth"
                                  :max-time="query.trendChartEndYearMonth"
                                  v-model:show="showCustomMonthRangeDialog"
                                  @dateRange:change="setCustomDateFilter"
                                  @error="onShowDateRangeError" />

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

    <export-dialog ref="exportDialog" />

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import TrendsChart from '@/components/desktop/TrendsChart.vue';
import AccountFilterSettingsCard from '@/views/desktop/common/cards/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from '@/views/desktop/common/cards/CategoryFilterSettingsCard.vue';
import TransactionTagFilterSettingsCard from '@/views/desktop/common/cards/TransactionTagFilterSettingsCard.vue';
import ExportDialog from '@/views/desktop/statistics/transaction/dialogs/ExportDialog.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';
import { useRouter, onBeforeRouteUpdate } from 'vue-router';
import { useDisplay, useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useStatisticsTransactionPageBase } from '@/views/base/statistics/StatisticsTransactionPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { type TransactionStatisticsPartialFilter, useStatisticsStore } from '@/stores/statistics.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { type TextualYearMonth, type TimeRangeAndDateType, DateRangeScene, DateRange } from '@/core/datetime.ts';
import { ThemeType } from '@/core/theme.ts';
import {
    ChartDataAggregationType,
    StatisticsAnalysisType,
    CategoricalChartType,
    ChartDataType,
    ChartSortingType,
    ChartDateAggregationType
} from '@/core/statistics.ts';

import {
    isDefined,
    isString,
    isNumber,
    arrayItemToObjectField
} from '@/lib/common.ts';
import {
    getGregorianCalendarYearAndMonthFromUnixTime,
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
    mdiMagnify,
    mdiMenu,
    mdiFilterOutline,
    mdiFilterCogOutline,
    mdiExport,
    mdiDotsVertical
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;
type TrendsChartType = InstanceType<typeof TrendsChart>;
type ExportDialogType = InstanceType<typeof ExportDialog>;

interface TransactionStatisticsProps {
    initAnalysisType?: string,
    initChartDataType?: string,
    initChartType?: string,
    initChartDateType?: string,
    initStartTime?: TextualYearMonth | '',
    initEndTime?: TextualYearMonth | '',
    initFilterAccountIds?: string,
    initFilterCategoryIds?: string,
    initTagFilter?: string,
    initKeyword?: string;
    initSortingType?: string,
    initTrendDateAggregationType?: string
    initAssetTrendsDateAggregationType?: string
}

const props = defineProps<TransactionStatisticsProps>();

const router = useRouter();
const display = useDisplay();
const theme = useTheme();

const {
    tt,
    getAllCategoricalChartTypes,
    getAllTrendChartTypes,
    formatAmountToWesternArabicNumeralsWithoutDigitGrouping,
    formatPercentToLocalizedNumerals
} = useI18n();

const {
    loading,
    analysisType,
    trendDateAggregationType,
    assetTrendsDateAggregationType,
    defaultCurrency,
    firstDayOfWeek,
    fiscalYearStart,
    allDateRanges,
    allSortingTypes,
    allTrendAnalysisDateAggregationTypes,
    allAssetTrendsDateAggregationTypes,
    query,
    queryChartDataCategory,
    queryDateType,
    queryStartTime,
    queryEndTime,
    queryDateRangeName,
    queryTrendDateAggregationTypeName,
    queryAssetTrendsDateAggregationTypeName,
    canChangeDateRange,
    canShiftDateRange,
    canUseCategoryFilter,
    canUseTagFilter,
    canUseKeywordFilter,
    showAmountInChart,
    totalAmountName,
    showPercentInCategoricalChart,
    showTotalAmountInTrendsChart,
    showStackedInTrendsChart,
    translateNameInTrendsChart,
    categoricalOverviewAnalysisData,
    categoricalAnalysisData,
    trendsAnalysisData,
    assetTrendsData,
    canShowCustomDateRange,
    getTransactionCategoricalAnalysisDataItemDisplayColor,
    getDisplayAmount
} = useStatisticsTransactionPageBase();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const statisticsStore = useStatisticsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const monthlyTrendsChart = useTemplateRef<TrendsChartType>('monthlyTrendsChart');
const dailyTrendsChart = useTemplateRef<TrendsChartType>('dailyTrendsChart');
const exportDialog = useTemplateRef<ExportDialogType>('exportDialog');

const activeTab = ref<string>('statisticsPage');
const initing = ref<boolean>(true);
const filterKeyword = ref<string>('');
const alwaysShowNav = ref<boolean>(display.mdAndUp.value);
const showNav = ref<boolean>(display.mdAndUp.value);
const showCustomDateRangeDialog = ref<boolean>(false);
const showCustomMonthRangeDialog = ref<boolean>(false);
const showFilterAccountDialog = ref<boolean>(false);
const showFilterCategoryDialog = ref<boolean>(false);
const showFilterTagDialog = ref<boolean>(false);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

const statisticsDataHasData = computed<boolean>(() => {
    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        return !!categoricalAnalysisData.value && !!categoricalAnalysisData.value.items && categoricalAnalysisData.value.items.length > 0;
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        return !!trendsAnalysisData.value && !!trendsAnalysisData.value.items && trendsAnalysisData.value.items.length > 0 && !!monthlyTrendsChart.value;
    } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
        return !!assetTrendsData.value && !!assetTrendsData.value.items && assetTrendsData.value.items.length > 0 && !!dailyTrendsChart.value;
    }

    return false;
});

const allChartTypes = computed<TypeAndDisplayName[]>(() => {
    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        return getAllCategoricalChartTypes(true);
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        return getAllTrendChartTypes();
    } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
        return getAllTrendChartTypes();
    } else {
        return [];
    }
});

const queryAnalysisType = computed<StatisticsAnalysisType>({
    get: () => analysisType.value,
    set: (value: number) => {
        setAnalysisType(value);
    }
});

const queryChartType = computed<number | undefined>({
    get: () => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return query.value.categoricalChartType;
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return query.value.trendChartType;
        } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
            return query.value.assetTrendsChartType;
        } else {
            return undefined;
        }
    },
    set: (value: number | undefined) => {
        setChartType(value);
    }
});

const queryChartDataType = computed<number>({
    get: () => query.value.chartDataType,
    set: (value: number) => {
        setChartDataType(value);
    }
});

const querySortingType = computed<number>({
    get: () => query.value.sortingType,
    set: (value: number) => {
        setSortingType(value);
    }
});

const isQuerySpecialChartType = computed<boolean>(() => {
    return ChartDataType.valueOf(queryChartDataType.value)?.specialChart ?? false;
});

const statisticsTextColor = computed<string>(() => {
    if (query.value.chartDataType === ChartDataType.OutflowsByAccount.type ||
        query.value.chartDataType === ChartDataType.ExpenseByAccount.type ||
        query.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
        query.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type) {
        return 'text-expense';
    } else if (query.value.chartDataType === ChartDataType.InflowsByAccount.type ||
        query.value.chartDataType === ChartDataType.IncomeByAccount.type ||
        query.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
        query.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type) {
        return 'text-income';
    } else {
        return 'text-default';
    }
});

function getFilterLinkUrl(): string {
    return `/statistics/transaction?${statisticsStore.getTransactionStatisticsPageParams(analysisType.value, trendDateAggregationType.value, assetTrendsDateAggregationType.value)}`;
}

function getTransactionItemLinkUrl(itemId: string, dateRange?: TimeRangeAndDateType): string {
    return `/transaction/list?${statisticsStore.getTransactionListPageParams(analysisType.value, itemId, dateRange)}`;
}

function init(initProps: TransactionStatisticsProps): void {
    let needReload = !isDefined(initProps.initAnalysisType);

    const filter: TransactionStatisticsPartialFilter = {
        chartDataType: initProps.initChartDataType ? parseInt(initProps.initChartDataType) : undefined,
        filterAccountIds: initProps.initFilterAccountIds ? arrayItemToObjectField(initProps.initFilterAccountIds.split(','), true) : {},
        filterCategoryIds: initProps.initFilterCategoryIds ? arrayItemToObjectField(initProps.initFilterCategoryIds.split(','), true) : {},
        tagFilter: initProps.initTagFilter,
        keyword: initProps.initKeyword,
        sortingType: initProps.initSortingType ? parseInt(initProps.initSortingType) : undefined
    };

    filterKeyword.value = filter.keyword || '';

    if (initProps.initAnalysisType === StatisticsAnalysisType.CategoricalAnalysis.toString()) {
        filter.categoricalChartType = initProps.initChartType ? parseInt(initProps.initChartType) : undefined;
        filter.categoricalChartDateType = initProps.initChartDateType ? parseInt(initProps.initChartDateType) : undefined;
        filter.categoricalChartStartTime = initProps.initStartTime ? parseInt(initProps.initStartTime) : undefined;
        filter.categoricalChartEndTime = initProps.initEndTime ? parseInt(initProps.initEndTime) : undefined;

        if (filter.categoricalChartDateType !== query.value.categoricalChartDateType) {
            needReload = true;
        } else if (filter.categoricalChartDateType === DateRange.Custom.type) {
            if (filter.categoricalChartStartTime !== query.value.categoricalChartStartTime
                || filter.categoricalChartEndTime !== query.value.categoricalChartEndTime) {
                needReload = true;
            }
        }

        if (initProps.initAnalysisType !== analysisType.value.toString()) {
            analysisType.value = StatisticsAnalysisType.CategoricalAnalysis;
            needReload = true;
        }
    } else if (initProps.initAnalysisType === StatisticsAnalysisType.TrendAnalysis.toString()) {
        filter.trendChartType = initProps.initChartType ? parseInt(initProps.initChartType) : undefined;
        filter.trendChartDateType = initProps.initChartDateType ? parseInt(initProps.initChartDateType) : undefined;
        filter.trendChartStartYearMonth = initProps.initStartTime;
        filter.trendChartEndYearMonth = initProps.initEndTime;

        if (filter.trendChartDateType !== query.value.trendChartDateType) {
            needReload = true;
        } else if (filter.trendChartDateType === DateRange.Custom.type) {
            if (filter.trendChartStartYearMonth !== query.value.trendChartStartYearMonth
                || filter.trendChartEndYearMonth !== query.value.trendChartEndYearMonth) {
                needReload = true;
            }
        }

        if (initProps.initAnalysisType !== analysisType.value.toString()) {
            analysisType.value = StatisticsAnalysisType.TrendAnalysis;
            needReload = true;
        }

        if (initProps.initTrendDateAggregationType) {
            trendDateAggregationType.value = parseInt(initProps.initTrendDateAggregationType);
        }
    } else if (initProps.initAnalysisType === StatisticsAnalysisType.AssetTrends.toString()) {
        filter.assetTrendsChartType = initProps.initChartType ? parseInt(initProps.initChartType) : undefined;
        filter.assetTrendsChartDateType = initProps.initChartDateType ? parseInt(initProps.initChartDateType) : undefined;
        filter.assetTrendsChartStartTime = initProps.initStartTime ? parseInt(initProps.initStartTime) : undefined;
        filter.assetTrendsChartEndTime = initProps.initEndTime ? parseInt(initProps.initEndTime) : undefined;

        if (filter.assetTrendsChartDateType !== query.value.assetTrendsChartDateType) {
            needReload = true;
        } else if (filter.assetTrendsChartDateType === DateRange.Custom.type) {
            if (filter.assetTrendsChartStartTime !== query.value.assetTrendsChartStartTime
                || filter.assetTrendsChartEndTime !== query.value.assetTrendsChartEndTime) {
                needReload = true;
            }
        }

        if (initProps.initAnalysisType !== analysisType.value.toString()) {
            analysisType.value = StatisticsAnalysisType.AssetTrends;
            needReload = true;
        }

        if (initProps.initAssetTrendsDateAggregationType) {
            assetTrendsDateAggregationType.value = parseInt(initProps.initAssetTrendsDateAggregationType);
        }
    }

    if (!isDefined(initProps.initAnalysisType)) {
        analysisType.value = StatisticsAnalysisType.CategoricalAnalysis;
        statisticsStore.initTransactionStatisticsFilter(analysisType.value);
    } else {
        statisticsStore.initTransactionStatisticsFilter(analysisType.value, filter);
    }

    if (!needReload && !statisticsStore.transactionStatisticsStateInvalid) {
        loading.value = false;
        initing.value = false;
        return;
    }

    Promise.all([
        accountsStore.loadAllAccounts({force: false}),
        transactionCategoriesStore.loadAllCategories({force: false})
    ]).then(() => {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            return statisticsStore.loadCategoricalAnalysis({
                force: false
            }) as Promise<unknown>;
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            return statisticsStore.loadTrendAnalysis({
                force: false
            }) as Promise<unknown>;
        } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
            return statisticsStore.loadAssetTrends({
                force: false
            }) as Promise<unknown>;
        } else {
            return Promise.reject('An error occurred');
        }
    }).then(() => {
        loading.value = false;
        initing.value = false;
    }).catch(error => {
        loading.value = false;
        initing.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function reload(force: boolean): Promise<unknown> | null {
    let dispatchPromise: Promise<unknown> | null = null;

    loading.value = true;

    if (query.value.chartDataType === ChartDataType.Overview.type ||
        query.value.chartDataType === ChartDataType.OutflowsByAccount.type ||
        query.value.chartDataType === ChartDataType.ExpenseByAccount.type ||
        query.value.chartDataType === ChartDataType.ExpenseByPrimaryCategory.type ||
        query.value.chartDataType === ChartDataType.ExpenseBySecondaryCategory.type ||
        query.value.chartDataType === ChartDataType.InflowsByAccount.type ||
        query.value.chartDataType === ChartDataType.IncomeByAccount.type ||
        query.value.chartDataType === ChartDataType.IncomeByPrimaryCategory.type ||
        query.value.chartDataType === ChartDataType.IncomeBySecondaryCategory.type ||
        query.value.chartDataType === ChartDataType.TotalOutflows.type ||
        query.value.chartDataType === ChartDataType.TotalExpense.type ||
        query.value.chartDataType === ChartDataType.TotalInflows.type ||
        query.value.chartDataType === ChartDataType.TotalIncome.type ||
        query.value.chartDataType === ChartDataType.NetCashFlow.type ||
        query.value.chartDataType === ChartDataType.NetIncome.type ||
        query.value.chartDataType === ChartDataType.NetWorth.type) {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            dispatchPromise = statisticsStore.loadCategoricalAnalysis({
                force: force
            });
        } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
            dispatchPromise = statisticsStore.loadTrendAnalysis({
                force: force
            });
        } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
            dispatchPromise = statisticsStore.loadAssetTrends({
                force: force
            });
        }
    } else if (query.value.chartDataType === ChartDataType.AccountTotalAssets.type ||
        query.value.chartDataType === ChartDataType.AccountTotalLiabilities.type) {
        if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
            dispatchPromise = accountsStore.loadAllAccounts({
                force: force
            });
        } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
            dispatchPromise = statisticsStore.loadAssetTrends({
                force: force
            });
        }
    }

    if (dispatchPromise) {
        dispatchPromise.then(() => {
            loading.value = false;

            if (force) {
                snackbar.value?.showMessage('Data has been updated');
            }
        }).catch(error => {
            loading.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }

    return dispatchPromise;
}

function setAnalysisType(type: StatisticsAnalysisType): void {
    if (analysisType.value === type) {
        return;
    }

    if (!ChartDataType.isAvailableForAnalysisType(query.value.chartDataType, type)) {
        let defaultChartDataType: ChartDataType = ChartDataType.Default;

        if (type === StatisticsAnalysisType.AssetTrends) {
            defaultChartDataType = ChartDataType.DefaultForAssetTrends;
        }

        statisticsStore.updateTransactionStatisticsFilter({
            chartDataType: defaultChartDataType.type
        });
    }

    if (analysisType.value !== StatisticsAnalysisType.TrendAnalysis && type === StatisticsAnalysisType.TrendAnalysis) {
        trendDateAggregationType.value = ChartDateAggregationType.Default.type;
    } else if (analysisType.value !== StatisticsAnalysisType.AssetTrends && type === StatisticsAnalysisType.AssetTrends) {
        assetTrendsDateAggregationType.value = ChartDateAggregationType.Default.type;
    }

    analysisType.value = type;
    loading.value = true;
    statisticsStore.updateTransactionStatisticsInvalidState(true);
    router.push(getFilterLinkUrl());
}

function setChartType(type?: number): void {
    let changed = false;

    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        changed = statisticsStore.updateTransactionStatisticsFilter({
            categoricalChartType: type
        });
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        changed = statisticsStore.updateTransactionStatisticsFilter({
            trendChartType: type
        });
    } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
        changed = statisticsStore.updateTransactionStatisticsFilter({
            assetTrendsChartType: type
        });
    }

    if (changed) {
        router.push(getFilterLinkUrl());
    }
}

function setChartDataType(type: number): void {
    const changed = statisticsStore.updateTransactionStatisticsFilter({
        chartDataType: type
    });

    if (changed) {
        router.push(getFilterLinkUrl());
    }
}

function setSortingType(type: number): void {
    if (type < ChartSortingType.Amount.type || type > ChartSortingType.Name.type) {
        return;
    }

    const changed = statisticsStore.updateTransactionStatisticsFilter({
        sortingType: type
    });

    if (changed) {
        router.push(getFilterLinkUrl());
    }
}

function setTrendDateAggregationType(type: number): void {
    const changed = trendDateAggregationType.value !== type;
    trendDateAggregationType.value = type;

    if (changed) {
        router.push(getFilterLinkUrl());
    }
}

function setAssetTrendsDateAggregationType(type: number): void {
    const changed = assetTrendsDateAggregationType.value !== type;
    assetTrendsDateAggregationType.value = type;

    if (changed) {
        router.push(getFilterLinkUrl());
    }
}

function setDateFilter(dateType: number): void {
    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        if (dateType === DateRange.Custom.type) { // Custom
            showCustomDateRangeDialog.value = true;
            return;
        } else if (query.value.categoricalChartDateType === dateType) {
            return;
        }
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        if (dateType === DateRange.Custom.type) { // Custom
            showCustomMonthRangeDialog.value = true;
            return;
        } else if (query.value.trendChartDateType === dateType) {
            return;
        }
    } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
        if (dateType === DateRange.Custom.type) { // Custom
            showCustomDateRangeDialog.value = true;
            return;
        } else if (query.value.assetTrendsChartDateType === dateType) {
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
            trendChartStartYearMonth: getGregorianCalendarYearAndMonthFromUnixTime(dateRange.minTime),
            trendChartEndYearMonth: getGregorianCalendarYearAndMonthFromUnixTime(dateRange.maxTime)
        });
    } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
        changed = statisticsStore.updateTransactionStatisticsFilter({
            assetTrendsChartDateType: dateRange.dateType,
            assetTrendsChartStartTime: dateRange.minTime,
            assetTrendsChartEndTime: dateRange.maxTime
        });
    }

    if (changed) {
        loading.value = true;
        statisticsStore.updateTransactionStatisticsInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function setCustomDateFilter(startTime: number | TextualYearMonth, endTime: number | TextualYearMonth): void {
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

        showCustomDateRangeDialog.value = false;
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis && isString(startTime) && isString(endTime)) {
        const chartDateType = getDateTypeByDateRange(getYearMonthFirstUnixTime(startTime), getYearMonthLastUnixTime(endTime), firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.TrendAnalysis);

        changed = statisticsStore.updateTransactionStatisticsFilter({
            trendChartDateType: chartDateType,
            trendChartStartYearMonth: startTime,
            trendChartEndYearMonth: endTime
        });

        showCustomMonthRangeDialog.value = false;
    } else if (analysisType.value === StatisticsAnalysisType.AssetTrends && isNumber(startTime) && isNumber(endTime)) {
        const chartDateType = getDateTypeByDateRange(startTime, endTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.AssetTrends);

        changed = statisticsStore.updateTransactionStatisticsFilter({
            assetTrendsChartDateType: chartDateType,
            assetTrendsChartStartTime: startTime,
            assetTrendsChartEndTime: endTime
        });

        showCustomDateRangeDialog.value = false;
    }

    if (changed) {
        loading.value = true;
        statisticsStore.updateTransactionStatisticsInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function shiftDateRange(scale: number): void {
    let changed = false;

    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        if (query.value.categoricalChartDateType === DateRange.All.type) {
            return;
        }

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
            trendChartStartYearMonth: getGregorianCalendarYearAndMonthFromUnixTime(newDateRange.minTime),
            trendChartEndYearMonth: getGregorianCalendarYearAndMonthFromUnixTime(newDateRange.maxTime)
        });
    } else if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
        if (query.value.assetTrendsChartDateType === DateRange.All.type) {
            return;
        }

        const newDateRange = getShiftedDateRangeAndDateType(query.value.assetTrendsChartStartTime, query.value.assetTrendsChartEndTime, scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.AssetTrends);

        changed = statisticsStore.updateTransactionStatisticsFilter({
            assetTrendsChartDateType: newDateRange.dateType,
            assetTrendsChartStartTime: newDateRange.minTime,
            assetTrendsChartEndTime: newDateRange.maxTime
        });
    }

    if (changed) {
        loading.value = true;
        statisticsStore.updateTransactionStatisticsInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function setAccountFilter(changed: boolean): void {
    showFilterAccountDialog.value = false;

    if (changed) {
        loading.value = true;
        statisticsStore.updateTransactionStatisticsInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function setCategoryFilter(changed: boolean): void {
    showFilterCategoryDialog.value = false;

    if (changed) {
        loading.value = true;
        statisticsStore.updateTransactionStatisticsInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function setTagFilter(changed: boolean): void {
    showFilterTagDialog.value = false;

    if (changed) {
        loading.value = true;
        statisticsStore.updateTransactionStatisticsInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function setKeywordFilter(keyword: string): void {
    if (analysisType.value === StatisticsAnalysisType.AssetTrends) {
        return;
    }

    if (query.value.keyword === keyword) {
        return;
    }

    let changed = false;

    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis) {
        changed = statisticsStore.updateTransactionStatisticsFilter({
            keyword: keyword
        });
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis) {
        changed = statisticsStore.updateTransactionStatisticsFilter({
            keyword: keyword
        });
    }

    if (changed) {
        loading.value = true;
        statisticsStore.updateTransactionStatisticsInvalidState(true);
        router.push(getFilterLinkUrl());
    }
}

function exportResults(): void {
    if (analysisType.value === StatisticsAnalysisType.CategoricalAnalysis && categoricalAnalysisData.value && categoricalAnalysisData.value.items) {
        exportDialog.value?.open({
            headers: [
                tt('Name'),
                tt('Amount') + ` (${defaultCurrency.value})`,
                tt('Proportion (%)')
            ],
            data: categoricalAnalysisData.value.items
                .filter(item => !item.hidden)
                .map(item => [
                    item.name,
                    formatAmountToWesternArabicNumeralsWithoutDigitGrouping(item.totalAmount),
                    item.percent.toFixed(4)
                ])
        });
    } else if (analysisType.value === StatisticsAnalysisType.TrendAnalysis && trendsAnalysisData.value && trendsAnalysisData.value.items && monthlyTrendsChart.value) {
        const exportData = monthlyTrendsChart.value.exportData();
        exportDialog.value?.open({
            headers: exportData.headers || [],
            data: exportData.data || []
        });
    } else if (analysisType.value === StatisticsAnalysisType.AssetTrends && assetTrendsData.value && assetTrendsData.value.items && dailyTrendsChart.value) {
        const exportData = dailyTrendsChart.value.exportData();
        exportDialog.value?.open({
            headers: exportData.headers || [],
            data: exportData.data || []
        });
    }
}

function onClickSankeyChartItem(sourceItemType: 'account' | 'category', sourceItemId: string, targetItemType?: 'account' | 'category', targetItemId?: string): void {
    if (sourceItemType === 'category' && targetItemType === 'category' && sourceItemId && targetItemId) {
        const sourceCategory = transactionCategoriesStore.allTransactionCategoriesMap[sourceItemId];
        const targetCategory = transactionCategoriesStore.allTransactionCategoriesMap[targetItemId];

        if (sourceCategory?.parentId === targetCategory?.id) {
            router.push(getTransactionItemLinkUrl(`${sourceItemType}:${sourceItemId}`));
            return;
        } else if (targetCategory?.parentId === sourceCategory?.id) {
            router.push(getTransactionItemLinkUrl(`${targetItemType}:${targetItemId}`));
            return;
        }
    }

    router.push(getTransactionItemLinkUrl(`${sourceItemType}:${sourceItemId}` + (targetItemType && targetItemId ? `-${targetItemType}:${targetItemId}` : '')));
}

function onClickPieChartItem(item: Record<string, unknown>): void {
    router.push(getTransactionItemLinkUrl(item['id'] as string));
}

function onClickTrendChartItem(item: { itemId: string, dateRange: TimeRangeAndDateType }): void {
    router.push(getTransactionItemLinkUrl(item.itemId, item.dateRange));
}

function onShowDateRangeError(message: string): void {
    snackbar.value?.showError(message);
}

onBeforeRouteUpdate((to) => {
    if (to.query) {
        init({
            initAnalysisType: (to.query['analysisType'] as string | null) || undefined,
            initChartDataType: (to.query['chartDataType'] as string | null) || undefined,
            initChartType: (to.query['chartType'] as string | null) || undefined,
            initChartDateType: (to.query['chartDateType'] as string | null) || undefined,
            initStartTime: (to.query['startTime'] as TextualYearMonth | null) || undefined,
            initEndTime: (to.query['endTime'] as TextualYearMonth | null) || undefined,
            initFilterAccountIds: (to.query['filterAccountIds'] as string | null) || undefined,
            initFilterCategoryIds: (to.query['filterCategoryIds'] as string | null) || undefined,
            initTagFilter: (to.query['tagFilter'] as string | null) || undefined,
            initKeyword: (to.query['keyword'] as string | null) || undefined,
            initSortingType: (to.query['sortingType'] as string | null) || undefined,
            initTrendDateAggregationType: (to.query['trendDateAggregationType'] as string | null) || undefined,
            initAssetTrendsDateAggregationType: (to.query['assetTrendsDateAggregationType'] as string | null) || undefined
        });
    } else {
        init({});
    }
});

watch(() => display.mdAndUp.value, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});

init(props);
</script>

<style>
.statistics-custom-datetime-range {
    line-height: 1rem;
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
    margin-inline-start: 6px;
}

.statistics-list-item .statistics-amount {
    opacity: 0.8;
}
</style>
