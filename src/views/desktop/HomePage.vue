<template>
    <v-row class="match-height">
        <v-col cols="12" lg="4" md="12">
            <v-card :class="{ 'disabled': loadingOverview }">
                <template #title>
                    <div class="d-flex align-center">
                        <div class="d-flex align-baseline">
                            <span class="text-2xl font-weight-bold">{{ displayDateRange.thisMonth.displayTime }}</span>
                            <span>·</span>
                            <span style="font-size: 1rem">{{ $t('Expense') }}</span>
                        </div>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ml-2" :icon="true" :loading="loadingOverview" @click="reload(true)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="icons.refresh" size="24" />
                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-card-text>
                    <h4 class="text-2xl font-weight-medium text-primary">
                        <span v-if="!loadingOverview || (transactionOverview && transactionOverview.thisMonth && transactionOverview.thisMonth.valid)">{{ transactionOverview && transactionOverview.thisMonth ? getDisplayExpenseAmount(transactionOverview.thisMonth) : '-' }}</span>
                        <v-skeleton-loader class="d-inline-block skeleton-no-margin mt-3 pb-1" width="120px" type="text" :loading="true" v-else-if="loadingOverview && (!transactionOverview || !transactionOverview.thisMonth || !transactionOverview.thisMonth.valid)"></v-skeleton-loader>
                        <v-btn class="ml-1" density="compact" color="default" variant="text"
                               :icon="true" @click="showAmountInHomePage = !showAmountInHomePage">
                            <v-icon :icon="showAmountInHomePage ? icons.eyeSlash : icons.eye" size="20" />
                        </v-btn>
                    </h4>
                    <div class="mt-1 mb-3">
                        <span class="mr-2">{{ $t('Monthly income') }}</span>
                        <span v-if="!loadingOverview || (transactionOverview && transactionOverview.thisMonth && transactionOverview.thisMonth.valid)">{{ transactionOverview && transactionOverview.thisMonth ? getDisplayIncomeAmount(transactionOverview.thisMonth) : '-' }}</span>
                        <v-skeleton-loader class="d-inline-block skeleton-no-margin mt-2" width="120px" type="text" :loading="true" v-else-if="loadingOverview && (!transactionOverview || !transactionOverview.thisMonth || !transactionOverview.thisMonth.valid)"></v-skeleton-loader>
                    </div>
                    <v-btn size="small" to="/transaction/list?dateType=7">{{ $t('View Details') }}</v-btn>
                    <v-img class="overview-card-background" src="img/desktop/card-background.png"/>
                    <v-img class="overview-card-background-image" width="116px" src="img/desktop/document.svg"/>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" lg="8" md="12">
            <v-card :class="{ 'disabled': loadingOverview }">
                <template #title>
                    <span>{{ $t('Asset Summary') }}</span>
                </template>

                <v-card-text>
                    <div class="mb-8">
                        <span class="text-body-1" v-if="!loadingOverview || (allAccounts && allAccounts.length)">{{ $t('format.misc.youHaveAccounts', { count: allAccounts.length }) }}</span>
                        <v-skeleton-loader class="skeleton-no-margin mt-1 mb-2 pb-1" width="200px" type="text" :loading="true" v-else-if="loadingOverview && (!allAccounts || !allAccounts.length)"></v-skeleton-loader>
                    </div>

                    <v-row>
                        <v-col cols="12" md="4">
                            <div class="d-flex align-center">
                                <div class="me-3">
                                    <v-avatar rounded color="secondary" size="42" class="elevation-1">
                                        <v-icon size="24" :icon="icons.totalAssets"/>
                                    </v-avatar>
                                </div>

                                <div class="d-flex flex-column">
                                    <span class="text-caption">{{ $t('Total assets') }}</span>
                                    <span class="text-h5" v-if="!loadingOverview || (allAccounts && allAccounts.length)">{{ totalAssets }}</span>
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-2" width="120px" type="text" :loading="true" v-else-if="loadingOverview && (!allAccounts || !allAccounts.length)"></v-skeleton-loader>
                                </div>
                            </div>
                        </v-col>

                        <v-col cols="12" md="4">
                            <div class="d-flex align-center">
                                <div class="me-3">
                                    <v-avatar rounded color="expense" size="42" class="elevation-1">
                                        <v-icon size="24" :icon="icons.totalLiabilities"/>
                                    </v-avatar>
                                </div>

                                <div class="d-flex flex-column">
                                    <span class="text-caption">{{ $t('Total liabilities') }}</span>
                                    <span class="text-h5" v-if="!loadingOverview || (allAccounts && allAccounts.length)">{{ totalLiabilities }}</span>
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-2" width="120px" type="text" :loading="true" v-else-if="loadingOverview && (!allAccounts || !allAccounts.length)"></v-skeleton-loader>
                                </div>
                            </div>
                        </v-col>

                        <v-col cols="12" md="4">
                            <div class="d-flex align-center">
                                <div class="me-3">
                                    <v-avatar rounded color="primary" size="42" class="elevation-1">
                                        <v-icon size="24" :icon="icons.netAssets"/>
                                    </v-avatar>
                                </div>

                                <div class="d-flex flex-column">
                                    <span class="text-caption">{{ $t('Net assets') }}</span>
                                    <span class="text-h5" v-if="!loadingOverview || (allAccounts && allAccounts.length)">{{ netAssets }}</span>
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-2" width="120px" type="text" :loading="true" v-else-if="loadingOverview && (!allAccounts || !allAccounts.length)"></v-skeleton-loader>
                                </div>
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" md="6">
            <v-row>
                <v-col cols="12" md="6">
                    <income-expense-overview-card
                        :loading="loadingOverview" :disabled="loadingOverview" :icon="icons.calendarToday"
                        :title="$t('Today')"
                        :expense-amount="transactionOverview.today && transactionOverview.today.valid ? getDisplayExpenseAmount(transactionOverview.today) : ''"
                        :income-amount="transactionOverview.today && transactionOverview.today.valid ? getDisplayIncomeAmount(transactionOverview.today) : ''"
                        :datetime="displayDateRange.today.displayTime"
                    >
                        <template #menus>
                            <v-list-item :prepend-icon="icons.viewDetails" :to="'/transaction/list?dateType=' + allDateRanges.Today.type">
                                <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                            </v-list-item>
                        </template>
                    </income-expense-overview-card>
                </v-col>

                <v-col cols="12" md="6">
                    <income-expense-overview-card
                        :loading="loadingOverview" :disabled="loadingOverview" :icon="icons.calendarWeek"
                        :title="$t('This Week')"
                        :expense-amount="transactionOverview.thisWeek && transactionOverview.thisWeek.valid ? getDisplayExpenseAmount(transactionOverview.thisWeek) : ''"
                        :income-amount="transactionOverview.thisWeek && transactionOverview.thisWeek.valid ? getDisplayIncomeAmount(transactionOverview.thisWeek) : ''"
                        :datetime="displayDateRange.thisWeek.startTime + '-' + displayDateRange.thisWeek.endTime"
                    >
                        <template #menus>
                            <v-list-item :prepend-icon="icons.viewDetails" :to="'/transaction/list?dateType=' + allDateRanges.ThisWeek.type">
                                <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                            </v-list-item>
                        </template>
                    </income-expense-overview-card>
                </v-col>

                <v-col cols="12" md="6">
                    <income-expense-overview-card
                        :loading="loadingOverview" :disabled="loadingOverview" :icon="icons.calendarMonth"
                        :title="$t('This Month')"
                        :expense-amount="transactionOverview.thisMonth && transactionOverview.thisMonth.valid ? getDisplayExpenseAmount(transactionOverview.thisMonth) : ''"
                        :income-amount="transactionOverview.thisMonth && transactionOverview.thisMonth.valid ? getDisplayIncomeAmount(transactionOverview.thisMonth) : ''"
                        :datetime="displayDateRange.thisMonth.startTime + '-' + displayDateRange.thisMonth.endTime"
                    >
                        <template #menus>
                            <v-list-item :prepend-icon="icons.viewDetails" :to="'/transaction/list?dateType=' + allDateRanges.ThisMonth.type">
                                <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                            </v-list-item>
                        </template>
                    </income-expense-overview-card>
                </v-col>

                <v-col cols="12" md="6">
                    <income-expense-overview-card
                        :loading="loadingOverview" :disabled="loadingOverview" :icon="icons.calendarYear"
                        :title="$t('This Year')"
                        :expense-amount="transactionOverview.thisYear && transactionOverview.thisYear.valid ? getDisplayExpenseAmount(transactionOverview.thisYear) : ''"
                        :income-amount="transactionOverview.thisYear && transactionOverview.thisYear.valid ? getDisplayIncomeAmount(transactionOverview.thisYear) : ''"
                        :datetime="displayDateRange.thisYear.displayTime"
                    >
                        <template #menus>
                            <v-list-item :prepend-icon="icons.viewDetails" :to="'/transaction/list?dateType=' + allDateRanges.ThisYear.type">
                                <v-list-item-title>{{ $t('View Details') }}</v-list-item-title>
                            </v-list-item>
                        </template>
                    </income-expense-overview-card>
                </v-col>
            </v-row>
        </v-col>

        <v-col cols="12" md="6">
            <monthly-income-and-expense-card :data="monthlyIncomeAndExpenseData" :is-dark-mode="isDarkMode"
                                             :loading="loadingOverview" :disabled="loadingOverview"
                                             :enable-click-item="true" @click="clickMonthlyIncomeOrExpense" />
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script>
import { useTheme } from 'vuetify';

import IncomeExpenseOverviewCard from './overview/cards/IncomeExpenseOverviewCard.vue';
import MonthlyIncomeAndExpenseCard from './overview/cards/MonthlyIncomeAndExpenseCard.vue';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useOverviewStore } from '@/stores/overview.js';

import { DateRange } from '@/core/datetime.ts';
import { ThemeType } from '@/core/theme.ts';
import {
    formatUnixTime,
    getUnixTimeBeforeUnixTime,
    getUnixTimeAfterUnixTime
} from '@/lib/datetime.js';

import {
    mdiRefresh,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiBankOutline,
    mdiCreditCardOutline,
    mdiPiggyBankOutline,
    mdiCalendarTodayOutline,
    mdiCalendarWeekOutline,
    mdiCalendarMonthOutline,
    mdiLayersTripleOutline,
    mdiListBoxOutline,
    mdiDotsVertical
} from '@mdi/js';

export default {
    components: {
        IncomeExpenseOverviewCard,
        MonthlyIncomeAndExpenseCard
    },
    data() {
        return {
            loadingOverview: true,
            icons: {
                refresh: mdiRefresh,
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline,
                totalAssets: mdiBankOutline,
                totalLiabilities: mdiCreditCardOutline,
                netAssets: mdiPiggyBankOutline,
                calendarToday: mdiCalendarTodayOutline,
                calendarWeek: mdiCalendarWeekOutline,
                calendarMonth: mdiCalendarMonthOutline,
                calendarYear: mdiLayersTripleOutline,
                viewDetails: mdiListBoxOutline,
                more: mdiDotsVertical
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useOverviewStore),
        isDarkMode() {
            return this.globalTheme.global.name.value === ThemeType.Dark;
        },
        showAmountInHomePage: {
            get: function() {
                return this.settingsStore.appSettings.showAmountInHomePage;
            },
            set: function(value) {
                this.settingsStore.setShowAmountInHomePage(value);
            }
        },
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        allDateRanges() {
            return DateRange.all();
        },
        allAccounts() {
            return this.accountsStore.allAccounts;
        },
        netAssets() {
            const netAssets = this.accountsStore.getNetAssets(this.showAmountInHomePage);
            return this.getDisplayCurrency(netAssets, this.defaultCurrency);
        },
        totalAssets() {
            const totalAssets = this.accountsStore.getTotalAssets(this.showAmountInHomePage);
            return this.getDisplayCurrency(totalAssets, this.defaultCurrency);
        },
        totalLiabilities() {
            const totalLiabilities = this.accountsStore.getTotalLiabilities(this.showAmountInHomePage);
            return this.getDisplayCurrency(totalLiabilities, this.defaultCurrency);
        },
        displayDateRange() {
            const self = this;

            return {
                today: {
                    displayTime: self.$locale.formatUnixTimeToLongDate(self.userStore, self.overviewStore.transactionDataRange.today.startTime),
                },
                thisWeek: {
                    startTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.overviewStore.transactionDataRange.thisWeek.startTime),
                    endTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.overviewStore.transactionDataRange.thisWeek.endTime)
                },
                thisMonth: {
                    displayTime: formatUnixTime(self.overviewStore.transactionDataRange.thisMonth.startTime, 'MMMM'),
                    startTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.overviewStore.transactionDataRange.thisMonth.startTime),
                    endTime: self.$locale.formatUnixTimeToLongMonthDay(self.userStore, self.overviewStore.transactionDataRange.thisMonth.endTime)
                },
                thisYear: {
                    displayTime: self.$locale.formatUnixTimeToLongYear(self.userStore, self.overviewStore.transactionDataRange.thisYear.startTime)
                }
            };
        },
        transactionOverview() {
            return this.overviewStore.transactionOverview;
        },
        monthlyIncomeAndExpenseData() {
            const self = this;
            const data = [];

            if (!self.transactionOverview || !self.transactionOverview.thisMonth || !self.transactionOverview.thisMonth.valid) {
                return data;
            }

            [
                'monthBeforeLast10Months',
                'monthBeforeLast9Months',
                'monthBeforeLast8Months',
                'monthBeforeLast7Months',
                'monthBeforeLast6Months',
                'monthBeforeLast5Months',
                'monthBeforeLast4Months',
                'monthBeforeLast3Months',
                'monthBeforeLast2Months',
                'monthBeforeLastMonth',
                'lastMonth',
                'thisMonth'
            ].forEach(fieldName => {
                if (!Object.prototype.hasOwnProperty.call(self.transactionOverview, fieldName)) {
                    return;
                }

                const dateRange = self.overviewStore.transactionDataRange[fieldName];

                if (!dateRange) {
                    return;
                }

                const item = self.overviewStore.transactionOverview[fieldName];

                data.push({
                    monthStartTime: dateRange.startTime,
                    incomeAmount: item ? item.incomeAmount : 0,
                    expenseAmount: item ? item.expenseAmount : 0,
                    incompleteIncomeAmount: item ? item.incompleteIncomeAmount : true,
                    incompleteExpenseAmount: item ? item.incompleteExpenseAmount : true
                });
            });

            return data;
        }
    },
    created() {
        if (this.$user.isUserLogined() && this.$user.isUserUnlocked()) {
            this.reload(false);
        }
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    methods: {
        reload(force) {
            const self = this;

            self.loadingOverview = true;

            const promises = [
                self.accountsStore.loadAllAccounts({ force: false }),
                self.overviewStore.loadTransactionOverview({ force: force, loadLast11Months: true })
            ];

            Promise.all(promises).then(() => {
                self.loadingOverview = false;

                if (force) {
                    self.$refs.snackbar.showMessage('Data has been updated');
                }
            }).catch(error => {
                self.loadingOverview = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        clickMonthlyIncomeOrExpense(e) {
            const minTime = e.monthStartTime;
            const maxTime = getUnixTimeBeforeUnixTime(getUnixTimeAfterUnixTime(minTime, 1, 'months'), 1, 'seconds');
            const type = e.transactionType;

            this.$router.push(`/transaction/list?type=${type}&dateType=${DateRange.Custom.type}&maxTime=${maxTime}&minTime=${minTime}`);
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        },
        getDisplayAmount(amount, incomplete) {
            if (!this.showAmountInHomePage) {
                return this.getDisplayCurrency('***', this.defaultCurrency);
            }

            return this.getDisplayCurrency(amount, this.defaultCurrency) + (incomplete ? '+' : '');
        },
        getDisplayIncomeAmount(category) {
            return this.getDisplayAmount(category.incomeAmount, category.incompleteIncomeAmount);
        },
        getDisplayExpenseAmount(category) {
            return this.getDisplayAmount(category.expenseAmount, category.incompleteExpenseAmount);
        }
    }
}
</script>

<style>
.overview-card-background {
    position: absolute;
    inline-size: 9rem;
    inset-block-end: 0;
    inset-inline-end: 0;
}

.overview-card-background-image {
    position: absolute;
    inline-size: 5rem;
    inset-block-end: 0.5rem;
    inset-inline-end: 1rem;
}
</style>
