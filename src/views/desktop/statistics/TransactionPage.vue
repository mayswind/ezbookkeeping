<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="[
                                { name: $t('Pie Chart'), value: allChartTypes.Pie },
                                { name: $t('Bar Chart'), value: allChartTypes.Bar }
                            ]" v-model="query.chartType" @update:modelValue="setChartType" />
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="query.chartDataType">
                            <v-tab class="tab-text-truncate" :key="dataType.type" :value="dataType.type"
                                   v-for="dataType in allChartDataTypes">
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
                                            <span>{{ $t('Statistics Data') }}</span>
                                            <v-btn-group class="ml-4" color="default" density="comfortable" variant="outlined" divided>
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
                                                                <div class="statistics-custom-datetime-range" v-if="dateRange.type === allDateRanges.Custom.type && query.dateType === allDateRanges.Custom.type && query.startTime && query.endTime">
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
                                                        <v-list-item to="/app/settings?tab=statisticsSetting"
                                                                     :prepend-icon="icons.filterSettings"
                                                                     :title="$t('Settings')"></v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <div v-if="initing">
                                        <v-skeleton-loader type="paragraph" :loading="initing"
                                                           :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4 ]"></v-skeleton-loader>
                                    </div>

                                    <v-card-text class="statistics-overview-title pt-0" :class="{ 'disabled': loading }"
                                                 v-if="!initing && statisticsData && statisticsData.items && statisticsData.items.length">
                                        <span class="text-subtitle-1">{{ totalAmountName }}</span>
                                        <span class="statistics-overview-amount ml-3" :class="statisticsTextColor">
                                            {{ getDisplayAmount(statisticsData.totalAmount, defaultCurrency) }}
                                        </span>
                                    </v-card-text>

                                    <v-card-text class="statistics-overview-title pt-0"
                                                 v-else-if="!initing && (!statisticsData || !statisticsData.items || !statisticsData.items.length)">
                                        <span class="text-subtitle-1 statistics-overview-empty-tip">{{ $t('No transaction data') }}</span>
                                    </v-card-text>

                                    <v-card-text :class="{ 'readonly': loading }" v-if="!initing && query.chartType === allChartTypes.Pie">
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

                                    <v-card-text :class="{ 'readonly': loading }" v-if="!initing && query.chartType === allChartTypes.Bar">
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
                    </v-main>
                </v-layout>
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
import { useDisplay } from 'vuetify';

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
    getShiftedDateRangeAndDateType,
    getDateRangeByDateType
} from '@/lib/datetime.js';

import {
    mdiCheck,
    mdiArrowLeft,
    mdiArrowRight,
    mdiSort,
    mdiRefresh,
    mdiMenu,
    mdiFilterOutline,
    mdiFilterCogOutline,
    mdiPencilOutline,
    mdiDotsVertical,
} from '@mdi/js';

import AccountFilterSettingsCard from './settings/cards/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from './settings/cards/CategoryFilterSettingsCard.vue';

export default {
    components: {
        AccountFilterSettingsCard,
        CategoryFilterSettingsCard
    },
    data() {
        const { mdAndUp } = useDisplay();

        return {
            activeTab: 'statisticsPage',
            initing: true,
            loading: true,
            alwaysShowNav: mdAndUp.value,
            showNav: mdAndUp.value,
            showCustomDateRangeDialog: false,
            showFilterAccountDialog: false,
            showFilterCategoryDialog: false,
            icons: {
                check: mdiCheck,
                left: mdiArrowLeft,
                right: mdiArrowRight,
                sort: mdiSort,
                refresh: mdiRefresh,
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
        },
        'display.mdAndUp.value': function (newValue) {
            this.alwaysShowNav = newValue;

            if (!this.showNav) {
                this.showNav = newValue;
            }
        }
    },
    created() {
        const self = this;

        self.statisticsStore.initTransactionStatisticsFilter();

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
    setup() {
        const display = useDisplay();

        return {
            display: display
        };
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

            const newDateRange = getShiftedDateRangeAndDateType(startTime, endTime, scale, this.firstDayOfWeek);

            this.statisticsStore.updateTransactionStatisticsFilter({
                dateType: newDateRange.dateType,
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

            return this.$locale.getDateRangeDisplayName(this.userStore, query.dateType, query.startTime, query.endTime);
        },
        clickPieChartItem(item) {
            this.$router.push(this.getItemLinkUrl(item));
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
            return `/transaction/list?${this.statisticsStore.getTransactionListPageParams(item)}`;
        }
    }
}
</script>

<style>
.statistics-custom-datetime-range {
    font-size: 0.7rem;
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
    font-size: 0.7rem;
    opacity: 0.7;
    margin-left: 6px;
}

.statistics-list-item .statistics-amount {
    opacity: 0.8;
}
</style>
