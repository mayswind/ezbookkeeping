<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <div class="d-flex flex-column flex-md-row">
                    <div>
                        <v-tabs show-arrows direction="vertical"
                                class="text-uppercase my-4" v-model="query.chartDataType">
                            <v-tab :key="dataType.type" :value="dataType.type"
                                   v-for="dataType in allChartDataTypes">
                                {{ $t(dataType.name) }}
                            </v-tab>
                        </v-tabs>
                    </div>
                    <v-window class="d-flex flex-grow-1 ml-md-5 disable-tab-transition statistics-container" v-model="activeTab">
                        <v-window-item value="statisticsPage">
                            <v-card variant="flat">
                                <template #title>
                                    <div class="d-flex align-center">
                                        <div class="statistics-toolbar">
                                            <v-btn-toggle
                                                variant="outlined"
                                                color="primary"
                                                density="comfortable"
                                                mandatory="force"
                                                divided
                                                :disabled="loading"
                                                v-model="query.chartType"
                                            >
                                                <v-btn :value="allChartTypes.Pie" @click="setChartType(allChartTypes.Pie)">
                                                    {{ $t('Pie Chart') }}
                                                </v-btn>
                                                <v-btn :value="allChartTypes.Bar" @click="setChartType(allChartTypes.Bar)">
                                                    {{ $t('Bar Chart') }}
                                                </v-btn>
                                            </v-btn-toggle>

                                            <v-btn-group class="ml-3" color="default"
                                                         density="comfortable" variant="outlined"
                                                         divided>
                                                <v-btn :icon="icons.left"
                                                       :disabled="loading || query.dateType === allDateRanges.All.type || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type"
                                                       @click="shiftDateRange(query.startTime, query.endTime, -1)"/>
                                                <v-menu location="bottom">
                                                    <template #activator="{ props }">
                                                        <v-btn :disabled="loading || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type"
                                                               v-bind="props">{{ dateRangeName(query) }}</v-btn>
                                                    </template>
                                                    <v-list>
                                                        <v-list-item :key="dateRange.type" :value="dateRange.type"
                                                                     :append-icon="(query.dateType === dateRange.type ? icons.check : null)"
                                                                     v-for="dateRange in allDateRanges">
                                                            <v-list-item-title
                                                                class="cursor-pointer"
                                                                @click="setDateFilter(dateRange.type)">
                                                                {{ $t(dateRange.name) }}
                                                                <div class="text-body-2" v-if="dateRange.type === allDateRanges.Custom.type && query.dateType === allDateRanges.Custom.type && query.startTime && query.endTime">
                                                                    <small>
                                                                        <span>{{ queryStartTime }}</span>
                                                                        <span>&nbsp;-&nbsp;</span>
                                                                        <br/>
                                                                        <span>{{ queryEndTime }}</span>
                                                                    </small>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                    </v-list>
                                                </v-menu>
                                                <v-btn :icon="icons.right"
                                                       :disabled="loading || query.dateType === allDateRanges.All.type || query.chartDataType === allChartDataTypes.AccountTotalAssets.type || query.chartDataType === allChartDataTypes.AccountTotalLiabilities.type"
                                                       @click="shiftDateRange(query.startTime, query.endTime, 1)"/>
                                            </v-btn-group>

                                            <v-menu location="bottom">
                                                <template #activator="{ props }">
                                                    <v-btn class="ml-3" color="default" variant="outlined"
                                                           :prepend-icon="icons.sort" :disabled="loading"
                                                           v-bind="props">{{ querySortingTypeName }}</v-btn>
                                                </template>
                                                <v-list>
                                                    <v-list-item :key="sortingType.type" :value="sortingType.type"
                                                                 :append-icon="(query.sortingType === sortingType.type ? icons.check : null)"
                                                                 v-for="sortingType in allSortingTypes">
                                                        <v-list-item-title
                                                            class="cursor-pointer"
                                                            @click="setSortingType(sortingType.type)">
                                                            {{ $t(sortingType.fullName) }}
                                                        </v-list-item-title>
                                                    </v-list-item>
                                                </v-list>
                                            </v-menu>
                                        </div>
                                        <v-btn density="compact" color="default" variant="text"
                                               class="ml-2" :icon="true" :disabled="loading"
                                               v-if="!loading" @click="reload">
                                            <v-icon :icon="icons.refresh" size="24" />
                                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                        </v-btn>
                                        <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
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
                                                    <v-divider class="my-2"/>
                                                    <v-list-item :prepend-icon="icons.filterSettings"
                                                                 :title="$t('Settings')"
                                                                 @click="settings"></v-list-item>
                                                </v-list>
                                            </v-menu>
                                        </v-btn>
                                    </div>
                                </template>

                                <v-card-text v-if="initing">
                                    <v-skeleton-loader type="paragraph" :loading="initing"
                                                       :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
                                </v-card-text>

                                <v-card-title class="statistics-overview-title pt-0" v-if="!initing">
                                    <div>{{ totalAmountName }}</div>
                                    <div class="statistics-overview-amount ml-3" :class="statisticsTextColor"
                                          v-if="statisticsData && statisticsData.items && statisticsData.items.length">
                                        {{ getDisplayAmount(statisticsData.totalAmount, defaultCurrency) }}
                                    </div>
                                    <div class="text-subtitle-1 ml-3"
                                          v-else-if="!statisticsData || !statisticsData.items || !statisticsData.items.length">
                                        {{ $t('No transaction data') }}
                                    </div>
                                </v-card-title>

                                <v-card-text v-if="!initing && query.chartType === allChartTypes.Pie">
                                    <pie-chart
                                        :items="statisticsData && statisticsData.items && statisticsData.items.length ? statisticsData.items : []"
                                        :min-valid-percent="0.0001"
                                        :show-value="showAmountInChart"
                                        :enable-click-item="true"
                                        :default-currency="defaultCurrency"
                                        id-field="id"
                                        name-field="name"
                                        value-field="totalAmount"
                                        percent-field="percent"
                                        currency-field="currency"
                                        hidden-field="hidden"
                                        @click="clickPieChartItem"
                                    />
                                </v-card-text>

                                <v-card-text v-if="!initing && query.chartType === allChartTypes.Bar">
                                    <v-list rounded lines="two"
                                            v-if="statisticsData && statisticsData.items && statisticsData.items.length">
                                        <template :key="idx"
                                                  v-for="(item, idx) in statisticsData.items">
                                            <v-list-item class="pl-0" v-if="!item.hidden">
                                                <template #prepend>
                                                    <router-link class="statistics-list-item" :to="getItemLinkUrl(item)">
                                                        <ItemIcon :icon-type="queryChartDataCategory" size="34px"
                                                                  :icon-id="item.icon"
                                                                  :color="item.color"></ItemIcon>
                                                    </router-link>
                                                </template>
                                                <router-link class="statistics-list-item" :to="getItemLinkUrl(item)">
                                                    <div class="d-flex flex-column ml-2">
                                                        <div class="d-flex">
                                                            <span>{{ item.name }}</span>
                                                            <small class="statistics-percent" v-if="item.percent >= 0">{{ getDisplayPercent(item.percent, 2, '&lt;0.01') }}</small>
                                                            <v-spacer/>
                                                            <span class="statistics-amount">{{ getDisplayAmount(item.totalAmount, (item.currency || defaultCurrency)) }}</span>
                                                        </div>
                                                        <div>
                                                            <v-progress-linear :color="item.color ? '#' + item.color : 'primary'"
                                                                               :model-value="item.percent >= 0 ? item.percent : 0"
                                                                               :height="4"></v-progress-linear>
                                                        </div>
                                                    </div>
                                                </router-link>
                                            </v-list-item>
                                            <v-divider v-if="!item.hidden && idx !== statisticsData.items.length - 1"/>
                                        </template>
                                    </v-list>
                                </v-card-text>
                            </v-card>
                        </v-window-item>
                    </v-window>
                </div>
            </v-card>
        </v-col>
    </v-row>

    <date-range-selection-dialog :title="$t('Custom Date Range')"
                                  :min-time="query.startTime"
                                  :max-time="query.endTime"
                                  v-model:show="showCustomDateRangeDialog"
                                  @dateRange:change="setCustomDateFilter" />

    <v-dialog scrollable max-width="600" max-height="600" v-model="showFilterAccountDialog">
        <account-filter-settings-card
            :dialog-mode="true" :modify-default="false"
            @settings:change="showFilterAccountDialog = false" />
    </v-dialog>

    <v-dialog scrollable max-width="600" max-height="600" v-model="showFilterCategoryDialog">
        <category-filter-settings-card
            :dialog-mode="true" :modify-default="false"
            @settings:change="showFilterCategoryDialog = false" />
    </v-dialog>

    <snack-bar ref="snackbar" />
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
import { getNameByKeyValue, limitText, formatPercent } from '@/lib/common.js'
import {
    parseDateFromUnixTime,
    getYear,
    getShiftedDateRange,
    getDateRangeByDateType,
    isDateRangeMatchFullYears,
    isDateRangeMatchFullMonths
} from '@/lib/datetime.js';

import {
    mdiCheck,
    mdiArrowLeft,
    mdiArrowRight,
    mdiSort,
    mdiRefresh,
    mdiFilterOutline,
    mdiFilterCogOutline,
    mdiPencilOutline,
    mdiDotsVertical,
} from '@mdi/js';

import AccountFilterSettingsCard from '@/views/desktop/statistics/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from '@/views/desktop/statistics/CategoryFilterSettingsCard.vue';

export default {
    components: {
        AccountFilterSettingsCard,
        CategoryFilterSettingsCard
    },
    data() {
        return {
            activeTab: 'statisticsPage',
            initing: true,
            loading: true,
            showCustomDateRangeDialog: false,
            showFilterAccountDialog: false,
            showFilterCategoryDialog: false,
            icons: {
                check: mdiCheck,
                left: mdiArrowLeft,
                right: mdiArrowRight,
                sort: mdiSort,
                refresh: mdiRefresh,
                filter: mdiFilterOutline,
                filterSettings: mdiFilterCogOutline,
                pencil: mdiPencilOutline,
                more: mdiDotsVertical
            }
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
            return this.statisticsStore.transactionStatisticsChartDataCategory;
        },
        querySortingTypeName() {
            const querySortingTypeName = getNameByKeyValue(this.allSortingTypes, this.query.sortingType, 'type', 'fullName', 'System Default');
            return this.$t(querySortingTypeName);
        },
        queryStartTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.startTime);
        },
        queryEndTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.endTime);
        },
        allChartTypes() {
            return statisticsConstants.allChartTypes;
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
        statisticsData() {
            return this.statisticsStore.statisticsData;
        },
        statisticsTextColor() {
            if (this.query.chartDataType === this.allChartDataTypes.ExpenseByAccount.type ||
                this.query.chartDataType === this.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                this.query.chartDataType === this.allChartDataTypes.ExpenseBySecondaryCategory.type) {
                return 'text-expense';
            } else if (this.query.chartDataType === this.allChartDataTypes.IncomeByAccount.type ||
                this.query.chartDataType === this.allChartDataTypes.IncomeByPrimaryCategory.type ||
                this.query.chartDataType === this.allChartDataTypes.IncomeBySecondaryCategory.type) {
                return 'text-income';
            } else {
                return 'text-default';
            }
        },
        showAmountInChart() {
            if (!this.showAccountBalance
                && (this.query.chartDataType === this.allChartDataTypes.AccountTotalAssets.type || this.query.chartDataType === this.allChartDataTypes.AccountTotalLiabilities.type)) {
                return false;
            }

            return true;
        }
    },
    watch: {
        'query.chartDataType': function (newValue) {
            this.statisticsStore.updateTransactionStatisticsFilter({
                chartDataType: newValue
            });
        }
    },
    created() {
        const self = this;

        let defaultChartType = self.settingsStore.appSettings.statistics.defaultChartType;

        if (defaultChartType !== self.allChartTypes.Pie && defaultChartType !== self.allChartTypes.Bar) {
            defaultChartType = statisticsConstants.defaultChartType;
        }

        let defaultChartDataType = self.settingsStore.appSettings.statistics.defaultChartDataType;

        if (defaultChartDataType < self.allChartDataTypes.ExpenseByAccount.type || defaultChartDataType > self.allChartDataTypes.AccountTotalLiabilities.type) {
            defaultChartDataType = statisticsConstants.defaultChartDataType;
        }

        let defaultDateRange = self.settingsStore.appSettings.statistics.defaultDataRangeType;

        if (defaultDateRange < self.allDateRanges.All.type || defaultDateRange >= self.allDateRanges.Custom.type) {
            defaultDateRange = statisticsConstants.defaultDataRangeType;
        }

        let defaultSortType = self.settingsStore.appSettings.statistics.defaultSortingType;

        if (defaultSortType < self.allSortingTypes.Amount.type || defaultSortType > self.allSortingTypes.Name.type) {
            defaultSortType = statisticsConstants.defaultSortingType;
        }

        const dateRange = getDateRangeByDateType(defaultDateRange, self.firstDayOfWeek);

        self.statisticsStore.initTransactionStatisticsFilter({
            dateType: dateRange ? dateRange.dateType : undefined,
            startTime: dateRange ? dateRange.minTime : undefined,
            endTime: dateRange ? dateRange.maxTime : undefined,
            chartType: defaultChartType,
            chartDataType: defaultChartDataType,
            filterAccountIds: self.settingsStore.appSettings.statistics.defaultAccountFilter || {},
            filterCategoryIds: self.settingsStore.appSettings.statistics.defaultTransactionCategoryFilter || {},
            sortingType: defaultSortType,
        });

        Promise.all([
            self.accountsStore.loadAllAccounts({ force: false }),
            self.transactionCategoriesStore.loadAllCategories({ force: false })
        ]).then(() => {
            return self.statisticsStore.loadTransactionStatistics({
                force: false
            });
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
    methods: {
        reload(force) {
            const self = this;
            let dispatchPromise = null;

            self.loading = true;

            if (self.query.chartDataType === self.allChartDataTypes.ExpenseByAccount.type ||
                self.query.chartDataType === self.allChartDataTypes.ExpenseByPrimaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.ExpenseBySecondaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeByAccount.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeByPrimaryCategory.type ||
                self.query.chartDataType === self.allChartDataTypes.IncomeBySecondaryCategory.type) {
                dispatchPromise = self.statisticsStore.loadTransactionStatistics({
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
        },
        setChartType(chartType) {
            this.statisticsStore.updateTransactionStatisticsFilter({
                chartType: chartType
            });
        },
        setSortingType(sortingType) {
            if (sortingType < this.allSortingTypes.Amount.type || sortingType > this.allSortingTypes.Name.type) {
                return;
            }

            this.statisticsStore.updateTransactionStatisticsFilter({
                sortingType: sortingType
            });

            this.reload(null);
        },
        setDateFilter(dateType) {
            if (dateType === this.allDateRanges.Custom.type) { // Custom
                this.showCustomDateRangeDialog = true;
                return;
            } else if (this.query.dateType === dateType) {
                return;
            }

            const dateRange = getDateRangeByDateType(dateType, this.firstDayOfWeek);

            if (!dateRange) {
                return;
            }

            this.statisticsStore.updateTransactionStatisticsFilter({
                dateType: dateRange.dateType,
                startTime: dateRange.minTime,
                endTime: dateRange.maxTime
            });

            this.reload(null);
        },
        setCustomDateFilter(startTime, endTime) {
            if (!startTime || !endTime) {
                return;
            }

            this.statisticsStore.updateTransactionStatisticsFilter({
                dateType: this.allDateRanges.Custom.type,
                startTime: startTime,
                endTime: endTime
            });

            this.showCustomDateRangeDialog = false;

            this.reload(null);
        },
        shiftDateRange(startTime, endTime, scale) {
            if (this.query.dateType === this.allDateRanges.All.type) {
                return;
            }

            const newDateRange = getShiftedDateRange(startTime, endTime, scale);
            let newDateType = this.allDateRanges.Custom.type;

            for (let dateRangeField in this.allDateRanges) {
                if (!Object.prototype.hasOwnProperty.call(this.allDateRanges, dateRangeField)) {
                    continue;
                }

                const dateRangeType = this.allDateRanges[dateRangeField];
                const dateRange = getDateRangeByDateType(dateRangeType.type, this.firstDayOfWeek);

                if (dateRange && dateRange.minTime === newDateRange.minTime && dateRange.maxTime === newDateRange.maxTime) {
                    newDateType = dateRangeType.type;
                    break;
                }
            }

            this.statisticsStore.updateTransactionStatisticsFilter({
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

            if (isDateRangeMatchFullYears(query.startTime, query.endTime)) {
                const displayStartTime = this.$locale.formatUnixTimeToShortYear(this.userStore, query.startTime);
                const displayEndTime = this.$locale.formatUnixTimeToShortYear(this.userStore, query.endTime);

                return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
            }

            if (isDateRangeMatchFullMonths(query.startTime, query.endTime)) {
                const displayStartTime = this.$locale.formatUnixTimeToShortYearMonth(this.userStore, query.startTime);
                const displayEndTime = this.$locale.formatUnixTimeToShortYearMonth(this.userStore, query.endTime);

                return displayStartTime !== displayEndTime ? `${displayStartTime} ~ ${displayEndTime}` : displayStartTime;
            }

            const startTimeYear = getYear(parseDateFromUnixTime(query.startTime));
            const endTimeYear = getYear(parseDateFromUnixTime(query.endTime));

            const displayStartTime = this.$locale.formatUnixTimeToShortDate(this.userStore, query.startTime);
            const displayEndTime = this.$locale.formatUnixTimeToShortDate(this.userStore, query.endTime);

            if (displayStartTime === displayEndTime) {
                return displayStartTime;
            } else if (startTimeYear === endTimeYear) {
                const displayShortEndTime = this.$locale.formatUnixTimeToShortMonthDay(this.userStore, query.endTime);
                return `${displayStartTime} ~ ${displayShortEndTime}`;
            }

            return `${displayStartTime} ~ ${displayEndTime}`;
        },
        clickPieChartItem(item) {
            this.$router.push(this.getItemLinkUrl(item));
        },
        settings() {
            this.$router.push('/app/settings?tab=statisticsSetting');
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
            return this.$locale.getDisplayCurrency(value, currencyCode, {
                currencyDisplayMode: this.settingsStore.appSettings.currencyDisplayMode,
                enableThousandsSeparator: this.settingsStore.appSettings.thousandsSeparator
            });
        },
        getDisplayPercent(value, precision, lowPrecisionValue) {
            return formatPercent(value, precision, lowPrecisionValue);
        },
        getItemLinkUrl(item) {
            return `/transactions?${this.statisticsStore.getTransactionListPageParams(item)}`;
        }
    }
}
</script>

<style>
.statistics-container.v-window > .v-window__container {
    width: 100%;
}

.statistics-toolbar {
    overflow-x: auto;
    white-space: nowrap;
}

.statistics-overview-title {
    line-height: 2rem !important;
    display: flex;
    align-items: flex-end;
}

.statistics-overview-amount {
    font-size: 1.5rem;
    overflow: hidden;
    text-overflow: ellipsis;
}

.statistics-list-item {
    color: var(--v-theme-on-default);
    font-size: 1rem !important;
    line-height: 1.75rem;
    overflow: hidden;
    text-overflow: ellipsis;
}

.statistics-list-item .statistics-percent {
    font-size: 0.7rem;
    opacity: 0.7;
    margin-left: 6px;
}

.statistics-list-item .statistics-amount {
    opacity: 0.8;
}
</style>
