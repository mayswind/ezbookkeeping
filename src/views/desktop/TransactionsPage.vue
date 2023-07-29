<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <div class="d-flex flex-column flex-md-row">
                    <div>
                        <div class="mx-6 my-4">
                            <div class="transaction-type-buttons d-flex flex-column">
                                <v-btn border :color="query.type === 0 ? 'primary' : 'default'"
                                       :variant="query.type === 0 ? 'tonal' : 'outlined'" :disabled="loading"
                                       @click="changeTypeFilter(0)">
                                    {{ $t('All Types') }}
                                </v-btn>
                                <v-btn border :color="query.type === 1 ? 'primary' : 'default'"
                                       :variant="query.type === 1 ? 'tonal' : 'outlined'" :disabled="loading"
                                       @click="changeTypeFilter(1)">
                                    {{ $t('Modify Balance') }}
                                </v-btn>
                                <v-btn border :color="query.type === 2 ? 'primary' : 'default'"
                                       :variant="query.type === 2 ? 'tonal' : 'outlined'" :disabled="loading"
                                       @click="changeTypeFilter(2)">
                                    {{ $t('Income') }}
                                </v-btn>
                                <v-btn border :color="query.type === 3 ? 'primary' : 'default'"
                                       :variant="query.type === 3 ? 'tonal' : 'outlined'" :disabled="loading"
                                       @click="changeTypeFilter(3)">
                                    {{ $t('Expense') }}
                                </v-btn>
                                <v-btn border :color="query.type === 4 ? 'primary' : 'default'"
                                       :variant="query.type === 4 ? 'tonal' : 'outlined'" :disabled="loading"
                                       @click="changeTypeFilter(4)">
                                    {{ $t('Transfer') }}
                                </v-btn>
                            </div>
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="recentDateRangeType">
                            <v-tab :key="idx" :value="idx" v-for="(recentDateRange, idx) in recentMonthDateRanges"
                                   @click="changeDateFilter(recentDateRange)">
                                {{ recentDateRange.displayName }}
                            </v-tab>
                        </v-tabs>
                    </div>
                    <v-window class="d-flex flex-grow-1 ml-md-5 disable-tab-transition w-100-window-container" v-model="activeTab">
                        <v-window-item value="transactionPage">
                            <v-card variant="flat">
                                <template #title>
                                    <div class="transaction-list-title d-flex align-center text-no-wrap">
                                        <span>{{ $t('Transaction List') }}</span>
                                        <v-btn class="ml-3" color="default" variant="outlined"
                                               :disabled="loading || !canAddTransaction" @click="add">{{ $t('Add') }}</v-btn>
                                        <v-btn density="compact" color="default" variant="text"
                                               class="ml-2" :icon="true" :disabled="loading"
                                               v-if="!loading" @click="reload">
                                            <v-icon :icon="icons.refresh" size="24" />
                                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                        </v-btn>
                                        <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
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
                                        <span class="text-body-1 transaction-list-datetime-range-text ml-2">
                                            <span v-if="!this.query.minTime && !this.query.maxTime">{{ $t('All') }}</span>
                                            <span v-else-if="this.query.minTime || this.query.maxTime">{{ `${queryMinTime} - ${queryMaxTime}` }}</span>
                                        </span>
                                        <v-spacer/>
                                        <div class="transaction-list-total-amount-text d-flex align-center" v-if="showTotalAmountInTransactionListPage && monthlyDataTotalAmount">
                                            <span class="ml-2 text-subtitle-1">{{ $t('Total Income') }}</span>
                                            <span class="text-income ml-2" v-if="loading">
                                                <v-skeleton-loader type="text" style="width: 60px" :loading="true"></v-skeleton-loader>
                                            </span>
                                            <span class="text-income ml-2" v-else-if="!loading">
                                                {{ monthlyDataTotalAmount.income }}
                                            </span>
                                            <span class="text-subtitle-1 ml-3">{{ $t('Total Expense') }}</span>
                                            <span class="text-income ml-2" v-if="loading">
                                                <v-skeleton-loader type="text" style="width: 60px" :loading="true"></v-skeleton-loader>
                                            </span>
                                            <span class="text-expense ml-2" v-else-if="!loading">
                                                {{ monthlyDataTotalAmount.expense }}
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

                                    <tbody v-if="loading">
                                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]">
                                        <td class="px-0" colspan="5">
                                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                        </td>
                                    </tr>
                                    </tbody>

                                    <tbody v-if="!loading && noTransaction">
                                    <tr>
                                        <td colspan="5">{{ $t('No transaction data') }}</td>
                                    </tr>
                                    </tbody>

                                    <template :key="transactionMonthList.yearMonth"
                                              v-for="(transactionMonthList, monthIdx) in transactions">
                                        <tbody :class="{ 'has-bottom-border': monthIdx < transactions.length - 1 }" v-if="shouldShowMonthlyData(transactionMonthList)">
                                        <template :key="transaction.id" v-for="(transaction, idx) in transactionMonthList.items">
                                            <template v-if="monthlyDatePageFirstIndex <= idx && idx < monthlyDatePageLastIndex">
                                                <tr class="transaction-list-row-date no-hover text-sm"
                                                    v-if="idx === 0 || monthlyDatePageFirstIndex === idx || (idx > 0 && (transaction.day !== transactionMonthList.items[idx - 1].day))">
                                                    <td colspan="5" class="font-weight-bold">
                                                        <div class="d-flex align-center">
                                                            <span>{{ getLongDate(transaction) }}</span>
                                                            <v-chip class="ml-1" color="default" size="x-small">
                                                                {{ getWeekdayLongName(transaction) }}
                                                            </v-chip>
                                                        </div>
                                                    </td>
                                                </tr>
                                                <tr class="transaction-table-row-data text-sm cursor-pointer"
                                                    @click="show(transaction)">
                                                    <td class="transaction-table-column-time">
                                                        <div class="d-flex flex-column">
                                                            <span>{{ getDisplayTime(transaction) }}</span>
                                                            <span class="text-caption" v-if="transaction.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimezone(transaction) }}</span>
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
                                                        <span v-if="transaction.sourceAccount">{{ transaction.sourceAccount.name }}</span>
                                                        <v-icon :icon="icons.arrowRight" v-if="transaction.sourceAccount && transaction.type === allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id"></v-icon>
                                                        <span v-if="transaction.sourceAccount && transaction.type === allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">{{ transaction.destinationAccount.name }}</span>
                                                    </td>
                                                    <td class="transaction-table-column-description text-truncate">
                                                        {{ transaction.comment }}
                                                    </td>
                                                </tr>
                                            </template>
                                        </template>
                                        </tbody>
                                    </template>
                                </v-table>

                                <div class="mt-2 mb-4" v-if="isShowMonthlyData">
                                    <v-pagination :total-visible="5" :length="monthlyDatePageCount"
                                                  v-model="currentPage"></v-pagination>
                                </div>
                            </v-card>
                        </v-window-item>
                    </v-window>
                </div>
            </v-card>
        </v-col>
    </v-row>

    <date-range-selection-dialog :title="$t('Custom Date Range')" :persistent="true"
                                 :min-time="query.minTime"
                                 :max-time="query.maxTime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="changeCustomDateFilter" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
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
    getUnixTime,
    getSpecifiedDayFirstUnixTime,
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore,
    getDateRangeByDateType,
    getRecentDateRangeType
} from '@/lib/datetime.js';
import {
    categoryTypeToTransactionType,
    transactionTypeToCategoryType
} from '@/lib/category.js';
import { scrollToMenuListItem } from '@/lib/ui.desktop.js';

import {
    mdiMagnify,
    mdiCheck,
    mdiTextBoxCheckOutline,
    mdiRefresh,
    mdiMenuDown,
    mdiPencilBoxOutline,
    mdiArrowRightThin,
    mdiDeleteOutline,
    mdiDotsVertical
} from '@mdi/js';

export default {
    props: [
        'initDateType',
        'initMaxTime',
        'initMinTime',
        'initType',
        'initCategoryId',
        'initAccountId'
    ],
    data() {
        return {
            loading: true,
            updating: false,
            activeTab: 'transactionPage',
            currentPage: 1,
            countPerPage: 15,
            searchKeyword: '',
            showCustomDateRangeDialog: false,
            transactionRemoving: {},
            icons: {
                search: mdiMagnify,
                check: mdiCheck,
                all: mdiTextBoxCheckOutline,
                refresh: mdiRefresh,
                dropdownMenu: mdiMenuDown,
                modifyBalance: mdiPencilBoxOutline,
                arrowRight: mdiArrowRightThin,
                remove: mdiDeleteOutline,
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
        transactions() {
            if (this.loading) {
                return [];
            }

            return this.transactionsStore.transactions;
        },
        noTransaction() {
            return this.transactionsStore.noTransaction;
        },
        hasMoreTransaction() {
            return this.transactionsStore.hasMoreTransaction;
        },
        isShowMonthlyData() {
            const recentDateRange = this.recentMonthDateRanges[this.recentDateRangeType];
            return recentDateRange.year && recentDateRange.month;
        },
        monthlyDatePageCount() {
            const recentDateRange = this.recentMonthDateRanges[this.recentDateRangeType];

            if (!recentDateRange || !recentDateRange.year || !recentDateRange.month) {
                return 1;
            }

            if (!this.transactions || !this.transactions.length) {
                return 1;
            }

            for (let i = 0; i < this.transactions.length; i++) {
                if (this.transactions[i].year === recentDateRange.year &&
                    this.transactions[i].month === recentDateRange.month) {
                    return Math.ceil(this.transactions[i].items.length / this.countPerPage);
                }
            }

            return 1;
        },
        monthlyDatePageFirstIndex() {
            const currentPage = this.currentPage >= 1 ? this.currentPage : 1;

            if (this.isShowMonthlyData) {
                return (currentPage - 1) * this.countPerPage;
            } else {
                return 0;
            }
        },
        monthlyDatePageLastIndex() {
            const currentPage = this.currentPage >= 1 ? this.currentPage : 1;

            if (this.isShowMonthlyData) {
                return currentPage * this.countPerPage;
            } else {
                return Number.MAX_SAFE_INTEGER;
            }
        },
        monthlyDataTotalAmount() {
            const recentDateRange = this.recentMonthDateRanges[this.recentDateRangeType];

            if (!recentDateRange || !recentDateRange.year || !recentDateRange.month) {
                return null;
            }

            if (!this.transactions || !this.transactions.length) {
                return null;
            }

            for (let i = 0; i < this.transactions.length; i++) {
                if (this.transactions[i].year === recentDateRange.year &&
                    this.transactions[i].month === recentDateRange.month) {
                    return {
                        income: this.getDisplayMonthTotalAmount(this.transactions[i].totalAmount.income, this.defaultCurrency, '', this.transactions[i].totalAmount.incompleteIncome),
                        expense: this.getDisplayMonthTotalAmount(this.transactions[i].totalAmount.expense, this.defaultCurrency, '', this.transactions[i].totalAmount.incompleteExpense)
                    };
                }
            }

            return null;
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
            return this.$locale.getAllRecentMonthDateRanges(this.userStore, true);
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

            this.reload(false);
        },
        reload(force) {
            const self = this;

            self.loading = true;

            Promise.all([
                self.accountsStore.loadAllAccounts({ force: false }),
                self.transactionCategoriesStore.loadAllCategories({ force: false })
            ]).then(() => {
                const recentDateRange = this.recentMonthDateRanges[this.recentDateRangeType];
                let yearMonth = null;

                if (recentDateRange.year && recentDateRange.month) {
                    yearMonth = {
                        year: recentDateRange.year,
                        month: recentDateRange.month
                    };
                }

                return self.transactionsStore.loadTransactions({
                    reload: true,
                    yearMonth: yearMonth,
                    force: force,
                    autoExpand: true,
                    defaultCurrency: self.defaultCurrency
                });
            }).then(() => {
                self.loading = false;
                self.currentPage = 1;

                if (force) {
                    self.$refs.snackbar.showMessage('Data has been updated');
                }
            }).catch(error => {
                self.loading = false;
                self.currentPage = 1;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
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

            this.$router.push(this.getFilterLinkUrl());
        },
        changeTypeFilter(type) {
            if (this.query.type === type) {
                return;
            }

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

            this.$router.push(this.getFilterLinkUrl());
        },
        changeCategoryFilter(categoryId) {
            if (this.query.categoryId === categoryId) {
                return;
            }

            this.transactionsStore.updateTransactionListFilter({
                categoryId: categoryId
            });

            this.showCategoryPopover = false;
            this.$router.push(this.getFilterLinkUrl());
        },
        changeAccountFilter(accountId) {
            if (this.query.accountId === accountId) {
                return;
            }

            this.transactionsStore.updateTransactionListFilter({
                accountId: accountId
            });

            this.$router.push(this.getFilterLinkUrl());
        },
        changeKeywordFilter(keyword) {
            if (this.query.keyword === keyword) {
                return;
            }

            this.transactionsStore.updateTransactionListFilter({
                keyword: keyword
            });

            this.reload(false);
        },
        add() {

        },
        duplicate() {

        },
        show() {

        },
        edit() {

        },
        remove(transaction) {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to delete this transaction?').then(() => {
                self.updating = true;
                self.transactionRemoving[transaction.id] = true;

                self.transactionsStore.deleteTransaction({
                    transaction: transaction,
                    defaultCurrency: self.defaultCurrency
                }).then(() => {
                    self.updating = false;
                    self.transactionRemoving[transaction.id] = false;
                }).catch(error => {
                    self.updating = false;
                    self.transactionRemoving[transaction.id] = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        shouldShowMonthlyData(transactionMonthList) {
            const recentDateRange = this.recentMonthDateRanges[this.recentDateRangeType];

            if (!recentDateRange || !recentDateRange.year || !recentDateRange.month) {
                return true;
            }

            return transactionMonthList.year === recentDateRange.year && transactionMonthList.month === recentDateRange.month;
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
                scrollToMenuListItem(menu.contentEl);
            });
        },
        getDisplayTime(transaction) {
            return this.$locale.formatUnixTimeToShortTime(this.userStore, transaction.time, transaction.utcOffset, this.currentTimezoneOffsetMinutes);
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
            return `/transactions?${this.transactionsStore.getTransactionListPageParams()}`;
        }
    }
};
</script>

<style>
.transaction-keyword-filter .v-input--density-compact {
    --v-input-control-height: 36px;
    --v-input-padding-top: 5px;
    --v-input-padding-bottom: 5px;
    inline-size: 20rem;
}

.transaction-type-buttons .v-btn:not(:first-child) {
    border-top-left-radius: inherit;
    border-top-right-radius: inherit;
}

.transaction-type-buttons .v-btn:not(:last-child) {
    border-bottom: 0;
    border-bottom-left-radius: inherit;
    border-bottom-right-radius: inherit;
}

.transaction-list-title {
    overflow-x: auto;
    white-space: nowrap;
}

.transaction-list-datetime-range {
    height: 28px;
    overflow-x: auto;
    white-space: nowrap;
}

.transaction-list-datetime-range .transaction-list-datetime-range-text {
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
}

.transaction-list-total-amount-text .v-skeleton-loader__text {
    margin: 0;
}

.v-table.transaction-table .transaction-list-row-date > td {
    height: 40px !important;
}

.transaction-table tr.transaction-table-row-data .hover-display {
    display: none;
}

.transaction-table tr.transaction-table-row-data:hover .hover-display {
    display: grid;
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
