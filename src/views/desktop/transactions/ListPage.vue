<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="[
                                { name: $t('All Types'), value: 0 },
                                { name: $t('Modify Balance'), value: 1 },
                                { name: $t('Income'), value: 2 },
                                { name: $t('Expense'), value: 3 },
                                { name: $t('Transfer'), value: 4 }
                            ]" v-model="query.type" @update:model-value="changeTypeFilter" />
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <small>{{ $t('Transactions Per Page') }}</small>
                            <v-select class="mt-2" density="compact" :disabled="loading"
                                      :items="[ 5, 10, 15, 20, 25, 30, 50 ]"
                                      v-model="countPerPage"
                            />
                        </div>
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="recentDateRangeType">
                            <v-tab class="tab-text-truncate" :key="idx" :value="idx" v-for="(recentDateRange, idx) in recentMonthDateRanges"
                                   @click="changeDateFilter(recentDateRange)">
                                <span class="text-truncate">{{ recentDateRange.displayName }}</span>
                            </v-tab>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="transactionPage">
                                <v-card variant="flat" min-height="920">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center text-no-wrap">
                                            <v-btn class="mr-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="icons.menu" size="24" />
                                            </v-btn>
                                            <span>{{ $t('Transaction List') }}</span>
                                            <v-btn class="ml-3" color="default" variant="outlined"
                                                   :disabled="loading || !canAddTransaction" @click="add">{{ $t('Add') }}</v-btn>
                                            <v-btn density="compact" color="default" variant="text"
                                                   class="ml-2" :icon="true" :disabled="loading"
                                                   v-if="!loading" @click="reload">
                                                <v-icon :icon="icons.refresh" size="24" />
                                                <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-progress-circular indeterminate size="20" class="ml-3" v-if="loading"></v-progress-circular>
                                            <v-spacer/>
                                            <div class="transaction-keyword-filter ml-2">
                                                <v-text-field density="compact" :disabled="loading"
                                                              :prepend-inner-icon="icons.search"
                                                              :append-inner-icon="searchKeyword !== query.keyword ? icons.check : null"
                                                              :placeholder="$t('Search transaction description')"
                                                              v-model="searchKeyword"
                                                              @click:append-inner="changeKeywordFilter(searchKeyword)"
                                                              @keyup.enter="changeKeywordFilter(searchKeyword)"
                                                />
                                            </div>
                                        </div>
                                    </template>

                                    <v-card-text class="pt-0">
                                        <div class="transaction-list-datetime-range d-flex align-center">
                                            <span class="text-body-1">{{ $t('Date Range') }}</span>
                                            <span class="text-body-1 transaction-list-datetime-range-text ml-2"
                                                  v-if="!query.minTime && !query.maxTime">
                                                <span class="text-sm">{{ $t('All') }}</span>
                                            </span>
                                            <span class="text-body-1 transaction-list-datetime-range-text ml-2"
                                                  v-else-if="query.minTime || query.maxTime">
                                                <v-btn class="mr-1" size="small"
                                                       density="comfortable" color="default" variant="outlined"
                                                       :icon="icons.arrowLeft" :disabled="loading"
                                                       @click="shiftDateRange(query.minTime, query.maxTime, -1)"/>
                                                <span class="text-sm">{{ `${queryMinTime} - ${queryMaxTime}` }}</span>
                                                <v-btn class="ml-1" size="small"
                                                       density="comfortable" color="default" variant="outlined"
                                                       :icon="icons.arrowRight" :disabled="loading"
                                                       @click="shiftDateRange(query.minTime, query.maxTime, 1)"/>
                                            </span>
                                            <v-spacer/>
                                            <div class="skeleton-no-margin d-flex align-center" v-if="showTotalAmountInTransactionListPage && currentMonthTotalAmount">
                                                <span class="ml-2 text-subtitle-1">{{ $t('Total Income') }}</span>
                                                <span class="text-income ml-2" v-if="loading">
                                                    <v-skeleton-loader type="text" style="width: 60px" :loading="true"></v-skeleton-loader>
                                                </span>
                                                <span class="text-income ml-2" v-else-if="!loading">
                                                    {{ currentMonthTotalAmount.income }}
                                                </span>
                                                <span class="text-subtitle-1 ml-3">{{ $t('Total Expense') }}</span>
                                                <span class="text-income ml-2" v-if="loading">
                                                    <v-skeleton-loader type="text" style="width: 60px" :loading="true"></v-skeleton-loader>
                                                </span>
                                                <span class="text-expense ml-2" v-else-if="!loading">
                                                    {{ currentMonthTotalAmount.expense }}
                                                </span>
                                            </div>
                                        </div>
                                    </v-card-text>

                                    <v-table class="transaction-table" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th class="transaction-table-column-time text-uppercase">{{ $t('Time') }}</th>
                                            <th class="transaction-table-column-category text-uppercase">
                                                <v-menu ref="categoryFilterMenu" class="transaction-category-menu"
                                                        eager location="bottom" max-height="500"
                                                        :disabled="query.type === 1"
                                                        :close-on-content-click="false"
                                                        v-model="categoryMenuState"
                                                        @update:model-value="scrollCategoryMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center"
                                                            :class="{ 'readonly': loading, 'cursor-pointer': query.type !== 1, 'text-primary': query.categoryId > 0 }" v-bind="props">
                                                            <span>{{ queryCategoryName }}</span>
                                                            <v-icon :icon="icons.dropdownMenu" v-show="query.type !== 1" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[query.categoryId]">
                                                        <v-list-item key="0" value="0" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.categoryId === '0' }"
                                                                     :append-icon="(query.categoryId === '0' ? icons.check : null)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeCategoryFilter('0')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="icons.all" />
                                                                    <span class="text-sm ml-3">{{ $t('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>

                                                        <template :key="categoryType"
                                                                  v-for="(categories, categoryType) in allPrimaryCategories">
                                                            <v-list-item density="compact">
                                                                <v-list-item-title>
                                                                    <span class="text-sm">{{ getTransactionTypeName(getTransactionTypeFromCategoryType(categoryType), 'Type') }}</span>
                                                                </v-list-item-title>
                                                            </v-list-item>

                                                            <v-list-group :key="category.id" v-for="category in categories">
                                                                <template #activator="{ props }" v-if="!category.hidden">
                                                                    <v-divider />
                                                                    <v-list-item class="text-sm" density="compact"
                                                                                 :class="getCategoryListItemCheckedClass(category, query.categoryId)"
                                                                                 v-bind="props">
                                                                        <v-list-item-title>
                                                                            <div class="d-flex align-center">
                                                                                <ItemIcon icon-type="category" size="24px" :icon-id="category.icon" :color="category.color"></ItemIcon>
                                                                                <span class="text-sm ml-3">{{ category.name }}</span>
                                                                            </div>
                                                                        </v-list-item-title>
                                                                    </v-list-item>
                                                                </template>

                                                                <v-divider />
                                                                <v-list-item class="text-sm" density="compact"
                                                                             :value="category.id"
                                                                             :append-icon="(query.categoryId === category.id ? icons.check : null)">
                                                                    <v-list-item-title class="cursor-pointer"
                                                                                       @click="changeCategoryFilter(category.id)">
                                                                        <div class="d-flex align-center">
                                                                            <v-icon :icon="icons.all" />
                                                                            <span class="text-sm ml-3">{{ $t('All') }}</span>
                                                                        </div>
                                                                    </v-list-item-title>
                                                                </v-list-item>

                                                                <template :key="subCategory.id"
                                                                          v-for="subCategory in category.subCategories">
                                                                    <v-divider v-if="!subCategory.hidden" />
                                                                    <v-list-item class="text-sm" density="compact"
                                                                                 :value="subCategory.id"
                                                                                 :class="{ 'list-item-selected': query.categoryId === subCategory.id }"
                                                                                 :append-icon="(query.categoryId === subCategory.id ? icons.check : null)"
                                                                                 v-if="!subCategory.hidden">
                                                                        <v-list-item-title class="cursor-pointer"
                                                                                           @click="changeCategoryFilter(subCategory.id)">
                                                                            <div class="d-flex align-center">
                                                                                <ItemIcon icon-type="category" size="24px" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                                                                <span class="text-sm ml-3">{{ subCategory.name }}</span>
                                                                            </div>
                                                                        </v-list-item-title>
                                                                    </v-list-item>
                                                                </template>
                                                            </v-list-group>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-amount text-uppercase">{{ $t('Amount') }}</th>
                                            <th class="transaction-table-column-account text-uppercase">
                                                <v-menu ref="accountFilterMenu" class="transaction-account-menu"
                                                        eager location="bottom" max-height="500"
                                                        @update:model-value="scrollAccountMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center cursor-pointer"
                                                             :class="{ 'readonly': loading, 'text-primary': query.accountId > 0 }" v-bind="props">
                                                            <span>{{ queryAccountName }}</span>
                                                            <v-icon :icon="icons.dropdownMenu" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[query.accountId]">
                                                        <v-list-item key="0" value="0" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.accountId === '0' }"
                                                                     :append-icon="(query.accountId === '0' ? icons.check : null)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeAccountFilter('0')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="icons.all" />
                                                                    <span class="text-sm ml-3">{{ $t('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <template :key="account.id"
                                                                  v-for="account in allAccounts">
                                                            <v-divider v-if="!account.hidden" />
                                                            <v-list-item class="text-sm" density="compact"
                                                                         :value="account.id"
                                                                         :class="{ 'list-item-selected': query.accountId === account.id }"
                                                                         :append-icon="(query.accountId === account.id ? icons.check : null)"
                                                                         v-if="!account.hidden">
                                                                <v-list-item-title class="cursor-pointer"
                                                                                   @click="changeAccountFilter(account.id)">
                                                                    <div class="d-flex align-center">
                                                                        <ItemIcon icon-type="account" size="24px" :icon-id="account.icon" :color="account.color"></ItemIcon>
                                                                        <span class="text-sm ml-3">{{ account.name }}</span>
                                                                    </div>
                                                                </v-list-item-title>
                                                            </v-list-item>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-description text-uppercase">{{ $t('Description') }}</th>
                                        </tr>
                                        </thead>

                                        <tbody v-if="loading && (!transactions || !transactions.length || transactions.length < 1)">
                                        <tr :key="itemIdx" v-for="itemIdx in skeletonData">
                                            <td class="px-0" colspan="5">
                                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                            </td>
                                        </tr>
                                        </tbody>

                                        <tbody v-if="!loading && (!transactions || !transactions.length || transactions.length < 1)">
                                        <tr>
                                            <td colspan="5">{{ $t('No transaction data') }}</td>
                                        </tr>
                                        </tbody>

                                        <tbody :key="transaction.id"
                                               :class="{ 'disabled': loading, 'has-bottom-border': idx < transactions.length - 1 }"
                                               v-for="(transaction, idx) in transactions">
                                            <tr class="transaction-list-row-date no-hover text-sm"
                                                v-if="idx === 0 || (idx > 0 && (transaction.date !== transactions[idx - 1].date))">
                                                <td colspan="5" class="font-weight-bold">
                                                    <div class="d-flex align-center">
                                                        <span>{{ getLongDate(transaction) }}</span>
                                                        <v-chip class="ml-1" color="default" size="x-small">
                                                            {{ getWeekdayLongName(transaction) }}
                                                        </v-chip>
                                                    </div>
                                                </td>
                                            </tr>
                                            <tr class="transaction-table-row-data text-sm"
                                                :class="{ 'cursor-pointer': transaction.type !== allTransactionTypes.ModifyBalance }"
                                                @click="show(transaction)">
                                                <td class="transaction-table-column-time">
                                                    <div class="d-flex flex-column">
                                                        <span>{{ getDisplayTime(transaction) }}</span>
                                                        <span class="text-caption" v-if="transaction.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimezone(transaction) }}</span>
                                                        <v-tooltip activator="parent" v-if="transaction.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimeInDefaultTimezone(transaction) }}</v-tooltip>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-category">
                                                    <div class="d-flex align-center">
                                                        <ItemIcon size="24px" icon-type="category"
                                                                  :icon-id="transaction.category.icon"
                                                                  :color="transaction.category.color"
                                                                  v-if="transaction.category && transaction.category.color"></ItemIcon>
                                                        <v-icon size="24" :icon="icons.modifyBalance" v-else-if="!transaction.category || !transaction.category.color" />
                                                        <span class="ml-2" v-if="transaction.type === allTransactionTypes.ModifyBalance">
                                                            {{ $t('Modify Balance') }}
                                                        </span>
                                                        <span class="ml-2" v-else-if="transaction.type !== allTransactionTypes.ModifyBalance && transaction.category">
                                                            {{ transaction.category.name }}
                                                        </span>
                                                        <span class="ml-2" v-else-if="transaction.type !== allTransactionTypes.ModifyBalance && !transaction.category">
                                                            {{ getTransactionTypeName(transaction.type, 'Transaction') }}
                                                        </span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-amount" :class="{ 'text-expense': transaction.type === allTransactionTypes.Expense, 'text-income': transaction.type === allTransactionTypes.Income }">
                                                    <div v-if="transaction.sourceAccount">
                                                        <span v-if="!query.accountId || query.accountId === '0' || (transaction.sourceAccount && (transaction.sourceAccount.id === query.accountId || transaction.sourceAccount.parentId === query.accountId))">{{ getDisplayAmount(transaction.sourceAmount, transaction.sourceAccount.currency, transaction.hideAmount) }}</span>
                                                        <span v-else-if="query.accountId && query.accountId !== '0' && transaction.destinationAccount && (transaction.destinationAccount.id === query.accountId || transaction.destinationAccount.parentId === query.accountId)">{{ getDisplayAmount(transaction.destinationAmount, transaction.destinationAccount.currency, transaction.hideAmount) }}</span>
                                                        <span v-else></span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-account">
                                                    <div class="d-flex align-center">
                                                        <span v-if="transaction.sourceAccount">{{ transaction.sourceAccount.name }}</span>
                                                        <v-icon class="mx-1" size="13" :icon="icons.arrowRight" v-if="transaction.sourceAccount && transaction.type === allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id"></v-icon>
                                                        <span v-if="transaction.sourceAccount && transaction.type === allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">{{ transaction.destinationAccount.name }}</span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-description text-truncate">
                                                    {{ transaction.comment }}
                                                </td>
                                            </tr>
                                        </tbody>
                                    </v-table>

                                    <div class="mt-2 mb-4">
                                        <v-pagination :total-visible="6" :length="totalPageCount"
                                                      v-model="paginationCurrentPage"></v-pagination>
                                    </div>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <date-range-selection-dialog :title="$t('Custom Date Range')"
                                 :min-time="query.minTime"
                                 :max-time="query.maxTime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="changeCustomDateFilter" />
    <edit-dialog ref="editDialog" :persistent="true" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
import EditDialog from './list/dialogs/EditDialog.vue';

import { useDisplay } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionsStore } from '@/stores/transaction.js';

import datetimeConstants from '@/consts/datetime.js';
import currencyConstants from '@/consts/currency.js';
import accountConstants from '@/consts/account.js';
import transactionConstants from '@/consts/transaction.js';
import { getNameByKeyValue } from '@/lib/common.js';
import {
    getCurrentUnixTime,
    parseDateFromUnixTime,
    getYear,
    getMonth,
    getUnixTime,
    getSpecifiedDayFirstUnixTime,
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore,
    getShiftedDateRangeAndDateType,
    getDateRangeByDateType,
    getRecentDateRangeType,
    isDateRangeMatchOneMonth
} from '@/lib/datetime.js';
import {
    categoryTypeToTransactionType,
    transactionTypeToCategoryType
} from '@/lib/category.js';
import { scrollToSelectedItem } from '@/lib/ui.desktop.js';

import {
    mdiMagnify,
    mdiCheck,
    mdiTextBoxCheckOutline,
    mdiRefresh,
    mdiMenu,
    mdiMenuDown,
    mdiPencilBoxOutline,
    mdiArrowLeft,
    mdiArrowRight,
    mdiDotsVertical
} from '@mdi/js';

export default {
    components: {
        EditDialog
    },
    props: [
        'initDateType',
        'initMaxTime',
        'initMinTime',
        'initType',
        'initCategoryId',
        'initAccountId'
    ],
    data() {
        const { mdAndUp } = useDisplay();

        return {
            loading: true,
            updating: false,
            activeTab: 'transactionPage',
            currentPage: 1,
            temporaryCountPerPage: null,
            totalCount: 1,
            searchKeyword: '',
            currentPageTransactions: [],
            categoryMenuState: false,
            alwaysShowNav: mdAndUp.value,
            showNav: mdAndUp.value,
            showCustomDateRangeDialog: false,
            icons: {
                search: mdiMagnify,
                check: mdiCheck,
                all: mdiTextBoxCheckOutline,
                refresh: mdiRefresh,
                menu: mdiMenu,
                dropdownMenu: mdiMenuDown,
                modifyBalance: mdiPencilBoxOutline,
                arrowLeft: mdiArrowLeft,
                arrowRight: mdiArrowRight,
                more: mdiDotsVertical
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useTransactionsStore),
        defaultCurrency() {
            if (this.query.accountId && this.query.accountId !== '0') {
                const account = this.allAccounts[this.query.accountId];

                if (account && account.currency && account.currency !== currencyConstants.parentAccountCurrencyPlaceholder) {
                    return account.currency;
                }
            }

            return this.userStore.currentUserDefaultCurrency;
        },
        canAddTransaction() {
            if (this.query.accountId && this.query.accountId !== '0') {
                const account = this.allAccounts[this.query.accountId];

                if (account && account.type === accountConstants.allAccountTypes.MultiSubAccounts) {
                    return false;
                }
            }

            return true;
        },
        currentTimezoneOffsetMinutes() {
            return getTimezoneOffsetMinutes(this.settingsStore.appSettings.timeZone);
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        recentDateRangeType: {
            get: function () {
                return getRecentDateRangeType(this.recentMonthDateRanges, this.query.dateType, this.query.minTime, this.query.maxTime, this.firstDayOfWeek);
            },
            set: function (value) {
                if (value < 0 || value >= this.recentMonthDateRanges.length) {
                    value = 0;
                }

                this.changeDateFilter(this.recentMonthDateRanges[value]);
            }
        },
        query() {
            return this.transactionsStore.transactionsFilter;
        },
        queryMinTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.minTime);
        },
        queryMaxTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.maxTime);
        },
        queryCategoryName() {
            return getNameByKeyValue(this.allCategories, this.query.categoryId, null, 'name', this.$t('Category'));
        },
        queryAccountName() {
            return getNameByKeyValue(this.allAccounts, this.query.accountId, null, 'name', this.$t('Account'));
        },
        queryMonthlyData() {
            return isDateRangeMatchOneMonth(this.query.minTime, this.query.maxTime);
        },
        countPerPage: {
            get: function () {
                if (this.temporaryCountPerPage) {
                    return this.temporaryCountPerPage;
                }

                return this.settingsStore.appSettings.itemsCountInTransactionListPage;
            },
            set: function(value) {
                const newTotalPageCount = Math.ceil(this.totalCount / value);

                if (this.currentPage > newTotalPageCount) {
                    this.currentPage = newTotalPageCount;
                }

                this.temporaryCountPerPage = value;

                if (!this.queryMonthlyData) {
                    this.reload(false);
                }
            }
        },
        totalPageCount() {
            return Math.ceil(this.totalCount / this.countPerPage);
        },
        paginationCurrentPage: {
            get: function () {
                return this.currentPage;
            },
            set: function (value) {
                this.currentPage = value;

                if (!this.queryMonthlyData) {
                    this.reload(false);
                }
            }
        },
        skeletonData() {
            const data = [];

            for (let i = 0; i < this.countPerPage; i++) {
                data.push(i);
            }

            return data;
        },
        currentMonthTransactionData() {
            const allTransactions = this.transactionsStore.transactions;

            if (!allTransactions || !allTransactions.length) {
                return null;
            }

            const currentMonthMinDate = parseDateFromUnixTime(this.query.minTime);
            const currentYear = getYear(currentMonthMinDate);
            const currentMonth = getMonth(currentMonthMinDate);

            for (let i = 0; i < allTransactions.length; i++) {
                if (allTransactions[i].year === currentYear && allTransactions[i].month === currentMonth) {
                    return allTransactions[i];
                }
            }

            return null;
        },
        transactions() {
            if (this.queryMonthlyData) {
                const transactionData = this.currentMonthTransactionData;

                if (!transactionData || !transactionData.items) {
                    return [];
                }

                const firstIndex = (this.currentPage - 1) * this.countPerPage;
                const lastIndex = this.currentPage * this.countPerPage;

                return transactionData.items.slice(firstIndex, lastIndex);
            } else {
                return this.currentPageTransactions;
            }
        },
        currentMonthTotalAmount() {
            if (this.queryMonthlyData) {
                const transactionData = this.currentMonthTransactionData;

                if (!transactionData) {
                    return null;
                }

                return {
                    income: this.getDisplayMonthTotalAmount(transactionData.totalAmount.income, this.defaultCurrency, '', transactionData.totalAmount.incompleteIncome),
                    expense: this.getDisplayMonthTotalAmount(transactionData.totalAmount.expense, this.defaultCurrency, '', transactionData.totalAmount.incompleteExpense)
                };
            } else {
                return null;
            }
        },
        allTransactionTypes() {
            return transactionConstants.allTransactionTypes;
        },
        allAccounts() {
            return this.accountsStore.allAccountsMap;
        },
        allCategories() {
            return this.transactionCategoriesStore.allTransactionCategoriesMap;
        },
        allPrimaryCategories() {
            const primaryCategories = {};

            for (const categoryType in this.transactionCategoriesStore.allTransactionCategories) {
                if (!Object.prototype.hasOwnProperty.call(this.transactionCategoriesStore.allTransactionCategories, categoryType)) {
                    continue;
                }

                if (this.query.type && this.getTransactionTypeFromCategoryType(categoryType) !== this.query.type) {
                    continue;
                }

                primaryCategories[categoryType] = this.transactionCategoriesStore.allTransactionCategories[categoryType];
            }

            return primaryCategories;
        },
        recentMonthDateRanges() {
            return this.$locale.getAllRecentMonthDateRanges(this.userStore, true, true);
        },
        showTotalAmountInTransactionListPage() {
            return this.settingsStore.appSettings.showTotalAmountInTransactionListPage;
        }
    },
    created() {
        this.init({
            dateType: this.initDateType,
            minTime: this.initMinTime,
            maxTime: this.initMaxTime,
            type: this.initType,
            categoryId: this.initCategoryId,
            accountId: this.initAccountId
        });
    },
    setup() {
        const display = useDisplay();

        return {
            display: display
        };
    },
    watch: {
        'display.mdAndUp.value': function (newValue) {
            this.alwaysShowNav = newValue;

            if (!this.showNav) {
                this.showNav = newValue;
            }
        }
    },
    beforeRouteUpdate(to) {
        if (to.query) {
            this.init({
                dateType: to.query.dateType,
                minTime: to.query.minTime,
                maxTime: to.query.maxTime,
                type: to.query.type,
                categoryId: to.query.categoryId,
                accountId: to.query.accountId
            });
        }
    },
    methods: {
        init(query) {
            let dateRange = getDateRangeByDateType(query.dateType ? parseInt(query.dateType) : undefined, self.firstDayOfWeek);

            if (!dateRange &&
                query.dateType === datetimeConstants.allDateRanges.Custom.type.toString() &&
                parseInt(query.maxTime) > 0 && parseInt(query.minTime) > 0) {
                dateRange = {
                    dateType: parseInt(query.dateType),
                    maxTime: parseInt(query.maxTime),
                    minTime: parseInt(query.minTime)
                };
            }

            this.transactionsStore.initTransactionListFilter({
                dateType: dateRange ? dateRange.dateType : undefined,
                maxTime: dateRange ? dateRange.maxTime : undefined,
                minTime: dateRange ? dateRange.minTime : undefined,
                type: parseInt(query.type) > 0 ? parseInt(query.type) : undefined,
                categoryId: query.categoryId,
                accountId: query.accountId,
                keyword: this.searchKeyword
            });

            this.currentPage = 1;
            this.reload(false);
        },
        reload(force) {
            const self = this;

            self.loading = true;

            const page = self.currentPage;

            Promise.all([
                self.accountsStore.loadAllAccounts({ force: false }),
                self.transactionCategoriesStore.loadAllCategories({ force: false })
            ]).then(() => {
                if (this.queryMonthlyData) {
                    const currentMonthMinDate = parseDateFromUnixTime(this.query.minTime);
                    const currentYear = getYear(currentMonthMinDate);
                    const currentMonth = getMonth(currentMonthMinDate);

                    return self.transactionsStore.loadMonthlyAllTransactions({
                        year: currentYear,
                        month: currentMonth,
                        force: force,
                        autoExpand: true,
                        defaultCurrency: self.defaultCurrency
                    });
                } else {
                    return self.transactionsStore.loadTransactions({
                        reload: true,
                        force: force,
                        count: self.countPerPage,
                        page: page,
                        withCount: page <= 1,
                        autoExpand: true,
                        defaultCurrency: self.defaultCurrency
                    });
                }
            }).then(data => {
                self.loading = false;
                self.currentPageTransactions = data && data.items && data.items.length ? data.items : [];

                if (page <= 1) {
                    self.totalCount = data && data.totalCount ? data.totalCount : 1;
                }

                if (force) {
                    self.$refs.snackbar.showMessage('Data has been updated');
                }
            }).catch(error => {
                self.loading = false;
                self.currentPageTransactions = [];
                self.totalCount = 1;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        shiftDateRange(startTime, endTime, scale) {
            if (this.recentDateRangeType === datetimeConstants.allDateRanges.All.type) {
                return;
            }

            const newDateRange = getShiftedDateRangeAndDateType(startTime, endTime, scale, this.firstDayOfWeek);

            this.transactionsStore.updateTransactionListFilter({
                dateType: newDateRange.dateType,
                maxTime: newDateRange.maxTime,
                minTime: newDateRange.minTime
            });

            this.loading = true;
            this.currentPageTransactions = [];
            this.transactionsStore.clearTransactions();
            this.$router.push(this.getFilterLinkUrl());
        },
        changeDateFilter(recentDateRange) {
            if (recentDateRange.dateType === datetimeConstants.allDateRanges.Custom.type &&
                !recentDateRange.minTime && !recentDateRange.maxTime) { // Custom
                if (!this.query.minTime || !this.query.maxTime) {
                    this.query.maxTime = getActualUnixTimeForStore(getCurrentUnixTime(), this.currentTimezoneOffsetMinutes, getBrowserTimezoneOffsetMinutes());
                    this.query.minTime = getSpecifiedDayFirstUnixTime(this.query.maxTime);
                }

                this.showCustomDateRangeDialog = true;
                return;
            }

            this.transactionsStore.updateTransactionListFilter({
                dateType: recentDateRange.dateType,
                maxTime: recentDateRange.maxTime,
                minTime: recentDateRange.minTime
            });

            this.loading = true;
            this.currentPageTransactions = [];
            this.transactionsStore.clearTransactions();
            this.$router.push(this.getFilterLinkUrl());
        },
        changeCustomDateFilter(minTime, maxTime) {
            if (!minTime || !maxTime) {
                return;
            }

            this.transactionsStore.updateTransactionListFilter({
                dateType: datetimeConstants.allDateRanges.Custom.type,
                maxTime: maxTime,
                minTime: minTime
            });

            this.showCustomDateRangeDialog = false;

            this.loading = true;
            this.currentPageTransactions = [];
            this.transactionsStore.clearTransactions();
            this.$router.push(this.getFilterLinkUrl());
        },
        changeTypeFilter(type) {
            let removeCategoryFilter = false;

            if (type && this.query.categoryId) {
                const category = this.allCategories[this.query.categoryId];

                if (category && category.type !== transactionTypeToCategoryType(type)) {
                    removeCategoryFilter = true;
                }
            }

            this.transactionsStore.updateTransactionListFilter({
                type: type,
                categoryId: removeCategoryFilter ? '0' : undefined
            });

            this.loading = true;
            this.currentPageTransactions = [];
            this.transactionsStore.clearTransactions();
            this.$router.push(this.getFilterLinkUrl());
        },
        changeCategoryFilter(categoryId) {
            this.categoryMenuState = false;

            if (this.query.categoryId === categoryId) {
                return;
            }

            this.transactionsStore.updateTransactionListFilter({
                categoryId: categoryId
            });

            this.loading = true;
            this.currentPageTransactions = [];
            this.transactionsStore.clearTransactions();
            this.$router.push(this.getFilterLinkUrl());
        },
        changeAccountFilter(accountId) {
            if (this.query.accountId === accountId) {
                return;
            }

            this.transactionsStore.updateTransactionListFilter({
                accountId: accountId
            });

            this.loading = true;
            this.currentPageTransactions = [];
            this.transactionsStore.clearTransactions();
            this.$router.push(this.getFilterLinkUrl());
        },
        changeKeywordFilter(keyword) {
            if (this.query.keyword === keyword) {
                return;
            }

            this.transactionsStore.updateTransactionListFilter({
                keyword: keyword
            });

            this.currentPage = 1;
            this.reload(false);
        },
        add() {
            const self = this;

            self.$refs.editDialog.open({
                type: self.query.type,
                categoryId: self.query.categoryId,
                accountId: self.query.accountId
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }

                self.reload(false);
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        show(transaction) {
            const self = this;

            if (transaction.type === self.allTransactionTypes.ModifyBalance) {
                return;
            }

            self.$refs.editDialog.open({
                id: transaction.id,
                currentTransaction: transaction
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }

                self.reload(false);
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        scrollCategoryMenuToSelectedItem(opened) {
            if (opened) {
                this.scrollMenuToSelectedItem(this.$refs.categoryFilterMenu);
            }
        },
        scrollAccountMenuToSelectedItem(opened) {
            if (opened) {
                this.scrollMenuToSelectedItem(this.$refs.accountFilterMenu);
            }
        },
        scrollMenuToSelectedItem(menu) {
            this.$nextTick(() => {
                scrollToSelectedItem(menu.contentEl, 'div.v-list', 'div.v-list-item.list-item-selected');
            });
        },
        getDisplayTime(transaction) {
            return this.$locale.formatUnixTimeToShortTime(this.userStore, transaction.time, transaction.utcOffset, this.currentTimezoneOffsetMinutes);
        },
        getDisplayTimeInDefaultTimezone(transaction) {
            return `${this.$locale.formatUnixTimeToLongDateTime(this.userStore, transaction.time)} (UTC${getTimezoneOffset(this.settingsStore.appSettings.timeZone)})`;
        },
        getDisplayTimezone(transaction) {
            return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
        },
        getDisplayAmount(amount, currency, hideAmount) {
            if (hideAmount) {
                return this.getDisplayCurrency('***', currency);
            }

            return this.getDisplayCurrency(amount, currency);
        },
        getDisplayMonthTotalAmount(amount, currency, symbol, incomplete) {
            const displayAmount = this.getDisplayCurrency(amount, currency);
            return symbol + displayAmount + (incomplete ? '+' : '');
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.getDisplayCurrency(value, currencyCode, {
                currencyDisplayMode: this.settingsStore.appSettings.currencyDisplayMode,
                enableThousandsSeparator: this.settingsStore.appSettings.thousandsSeparator
            });
        },
        getLongDate(transaction) {
            const transactionTime = getUnixTime(parseDateFromUnixTime(transaction.time, transaction.utcOffset, this.currentTimezoneOffsetMinutes));
            return this.$locale.formatUnixTimeToLongDate(this.userStore, transactionTime);
        },
        getWeekdayLongName(transaction) {
            return this.$locale.getWeekdayLongName(transaction.dayOfWeek);
        },
        getTransactionTypeName(type, defaultName) {
            switch (type){
                case this.allTransactionTypes.ModifyBalance:
                    return this.$t('Modify Balance');
                case this.allTransactionTypes.Income:
                    return this.$t('Income');
                case this.allTransactionTypes.Expense:
                    return this.$t('Expense');
                case this.allTransactionTypes.Transfer:
                    return this.$t('Transfer');
                default:
                    return this.$t(defaultName);
            }
        },
        getTransactionTypeFromCategoryType(categoryType) {
            return categoryTypeToTransactionType(parseInt(categoryType));
        },
        getCategoryListItemCheckedClass(category, queryCategoryId) {
            if (category.id === queryCategoryId) {
                return {
                    'list-item-selected': true,
                    'has-children-item-selected': true
                };
            }

            for (let i = 0; i < category.subCategories.length; i++) {
                if (category.subCategories[i].id === queryCategoryId) {
                    return {
                        'list-item-selected': true,
                        'has-children-item-selected': true
                    };
                }
            }

            return [];
        },
        getFilterLinkUrl() {
            return `/transaction/list?${this.transactionsStore.getTransactionListPageParams()}`;
        }
    }
};
</script>

<style>
.transaction-keyword-filter .v-input--density-compact {
    --v-input-control-height: 36px !important;
    --v-input-padding-top: 5px !important;
    --v-input-padding-bottom: 5px !important;
    --v-input-chips-margin-top: 0px !important;
    --v-input-chips-margin-bottom: 0px !important;
    inline-size: 20rem;
}

.transaction-list-datetime-range {
    min-height: 28px;
    flex-wrap: wrap;
    row-gap: 1rem;
}

.transaction-list-datetime-range .transaction-list-datetime-range-text {
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
}

.v-table.transaction-table .transaction-list-row-date > td {
    height: 40px !important;
}

.transaction-table .transaction-table-column-time {
    width: 110px;
    white-space: nowrap;
}

.transaction-table .transaction-table-column-category {
    width: 140px;
    white-space: nowrap;
}

.transaction-table .transaction-table-column-amount {
    width: 120px;
    white-space: nowrap;
}

.transaction-table .transaction-table-column-account {
    width: 160px;
    white-space: nowrap;
}

.transaction-table-column-description {
    max-width: 300px;
}

.transaction-table .transaction-table-column-category .v-btn,
.transaction-table .transaction-table-column-account .v-btn {
    font-size: 0.75rem;
}

.transaction-table .transaction-table-column-category .v-btn .v-btn__append,
.transaction-table .transaction-table-column-account .v-btn .v-btn__append {
    margin-left: 0in;
}

.transaction-category-menu .item-icon,
.transaction-account-menu .item-icon,
.transaction-table .item-icon {
    padding-bottom: 3px;
}

.transaction-category-menu .has-children-item-selected span {
    font-weight: bold;
}
</style>
