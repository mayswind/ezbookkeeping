<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title>
                <span style="color: var(--f7-text-color)" v-if="!finishQuery">{{ tt('Reconciliation Statement') }}</span>
                <f7-link popover-open=".display-mode-popover-menu" v-if="finishQuery">
                    <span style="color: var(--f7-text-color)">{{ tt('Reconciliation Statement') }}</span>
                    <f7-icon class="page-title-bar-icon" style="opacity: 0.5"
                             color="gray" f7="chevron_down_circle_fill"></f7-icon>
                </f7-link>
            </f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': !validQuery }" @click="reload(false)" v-if="!finishQuery"></f7-link>
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': loading }" v-if="finishQuery" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-popover class="display-mode-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt('Transaction List')"
                              :class="{ 'list-item-selected': !showAccountBalanceTrendsCharts }"
                              @click="showAccountBalanceTrendsCharts = false">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!showAccountBalanceTrendsCharts"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt('Account Balance Trends')"
                              :class="{ 'list-item-selected': showAccountBalanceTrendsCharts }"
                              @click="showAccountBalanceTrendsCharts = true">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="showAccountBalanceTrendsCharts"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!finishQuery">
            <f7-list-item group-title>
                <small>{{ tt('Date Range') }}</small>
            </f7-list-item>
            <f7-list-item :key="dateRange.type"
                          :title="dateRange.displayName"
                          :disabled="!validQuery"
                          v-for="dateRange in allAvailableDateRanges"
                          @click="changeDateFilter(dateRange.type)">
                <template #after>
                    <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="queryDateRangeType === dateRange.type"></f7-icon>
                </template>
                <template #footer>
                    <div v-if="dateRange.isUserCustomRange && queryDateRangeType === dateRange.type && startTime && endTime">
                        <span>{{ displayStartDateTime }}</span>
                        <span>&nbsp;-&nbsp;</span>
                        <br/>
                        <span>{{ displayEndDateTime }}</span>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="finishQuery && !startTime && !endTime">
            <f7-list-item :title="tt('Date Range')" :after="tt('All')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="finishQuery && (startTime || endTime)">
            <f7-list-item :title="tt('Start Time')" :after="displayStartDateTime"></f7-list-item>
            <f7-list-item :title="tt('End Time')" :after="displayEndDateTime"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="finishQuery && loading">
            <f7-list-item :title="tt('Total Transactions')" after="Count"></f7-list-item>
            <f7-list-item :title="tt('Total Inflows')" after="Count"></f7-list-item>
            <f7-list-item :title="tt('Total Outflows')" after="Count"></f7-list-item>
            <f7-list-item :title="tt('Net Cash Flow')" after="Count"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="finishQuery && !loading">
            <f7-list-item :title="tt('Total Transactions')"
                          :after="formatNumberToLocalizedNumerals(reconciliationStatements.transactions.length)"
                          v-if="reconciliationStatements && reconciliationStatements.transactions"></f7-list-item>
            <f7-list-item :title="tt('Total Inflows')" :after="displayTotalInflows"></f7-list-item>
            <f7-list-item :title="tt('Total Outflows')" :after="displayTotalOutflows"></f7-list-item>
            <f7-list-item :title="tt('Net Cash Flow')" :after="displayTotalBalance"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="finishQuery && loading">
            <f7-list-item :title="tt('Opening Balance')" after="Count"></f7-list-item>
            <f7-list-item :title="tt('Closing Balance')" after="Count"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-if="finishQuery && !loading">
            <f7-list-item :title="tt('Opening Balance')" :after="displayOpeningBalance"></f7-list-item>
            <f7-list-item :title="tt('Closing Balance')" :after="displayClosingBalance"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers media-list
                 class="skeleton-text margin-vertical transaction-info-list reconciliation-statement-list"
                 v-if="finishQuery && !showAccountBalanceTrendsCharts && loading">
            <ul>
                <f7-list-item chevron-center
                    :key="index"
                    :class="{ 'transaction-info': type === 't', 'last-transaction-of-day': index === 2, 'reconciliation-statement-transaction-date': type === 'd' }"
                    :link="type === 't' ? '#' : null"
                    v-for="(type, index) in [ 'd', 't', 't', 'd', 't', 't', 't' ]"
                >
                    <div class="display-flex no-padding-horizontal" v-if="type === 'd'">
                        <div class="actual-item-inner">
                            <div class="item-title-row">
                                <div class="item-title">
                                    <small>yyyy-MM-dd</small>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="display-flex no-padding-horizontal" v-if="type === 't'">
                        <div class="item-media">
                            <div class="transaction-icon display-flex align-items-center">
                                <f7-icon f7="app_fill"></f7-icon>
                            </div>
                        </div>
                        <div class="actual-item-inner">
                            <div class="item-title-row">
                                <div class="item-title">
                                    <div class="transaction-category-name no-padding">
                                        <span>Category</span>
                                    </div>
                                </div>
                                <div class="item-after">
                                    <div class="transaction-amount">
                                        <span>0.00 USD</span>
                                    </div>
                                </div>
                            </div>
                            <div class="item-text">
                                <div class="transaction-description">
                                    <span>Transaction Description</span>
                                </div>
                            </div>
                            <div class="item-footer">
                                <div class="transaction-footer display-flex justify-content-space-between">
                                    <div class="flex-shrink-0">
                                        <span>HH:mm</span>
                                    </div>
                                    <div class="account-balance flex-shrink-1">
                                        <span>Balance</span>
                                        <span style="margin-inline-start: 4px">0.00 USD</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </f7-list-item>
            </ul>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical"
                 v-if="finishQuery && !showAccountBalanceTrendsCharts && !loading && (!allReconciliationStatementVirtualListItems || !allReconciliationStatementVirtualListItems.length)">
            <f7-list-item :title="tt('No transaction data')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers media-list virtual-list
                 class="margin-vertical transaction-info-list reconciliation-statement-list"
                 :virtual-list-params="{ items: allReconciliationStatementVirtualListItems, renderExternal, height: 'auto' }"
                 v-if="finishQuery && !showAccountBalanceTrendsCharts && !loading && allReconciliationStatementVirtualListItems && allReconciliationStatementVirtualListItems.length">
            <ul>
                <f7-list-item chevron-center
                              :key="item.index"
                              :id="item.transaction ? getTransactionDomId(item.transaction) : undefined"
                              :class="{ 'transaction-info': item.type == 'transaction', 'last-transaction-of-day': allReconciliationStatementVirtualListItems[item.index + 1] && allReconciliationStatementVirtualListItems[item.index + 1]!.type === 'date', 'reconciliation-statement-transaction-date': item.type == 'date' }"
                              :style="`top: ${virtualDataItems.topPosition}px`"
                              :virtual-list-index="item.index"
                              :swipeout="item.type === 'transaction' && !!item.transaction"
                              :accordion-item="item.type === 'transaction' && !!item.transaction"
                              :link="item.type == 'transaction' && item.transaction && item.transaction.type !== TransactionType.ModifyBalance ? `/transaction/detail?id=${item.transaction?.id}&type=${item.transaction.type}` : null"
                              v-for="item in virtualDataItems.items"
                >
                    <template #inner>
                        <div class="display-flex no-padding-horizontal" v-if="item.type == 'date' && item.displayDate">
                            <div class="actual-item-inner">
                                <div class="item-title-row">
                                    <div class="item-title">
                                        <small>{{ item.displayDate }}</small>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="display-flex no-padding-horizontal" v-if="item.type == 'transaction' && item.transaction">
                            <div class="item-media">
                                <div class="transaction-icon display-flex align-items-center">
                                    <ItemIcon icon-type="category"
                                              :icon-id="item.transaction.category?.icon"
                                              :color="item.transaction.category?.color"
                                              v-if="item.transaction.category && item.transaction.category?.color"></ItemIcon>
                                    <f7-icon v-else-if="!item.transaction.category || !item.transaction.category?.color"
                                             f7="pencil_ellipsis_rectangle">
                                    </f7-icon>
                                </div>
                            </div>
                            <div class="actual-item-inner">
                                <div class="item-title-row">
                                    <div class="item-title">
                                        <div class="transaction-category-name no-padding">
                                        <span v-if="item.transaction.type === TransactionType.ModifyBalance">
                                            {{ tt('Modify Balance') }}
                                        </span>
                                            <span v-else-if="item.transaction.type !== TransactionType.ModifyBalance && item.transaction.category">
                                            {{ item.transaction.category.name }}
                                        </span>
                                        </div>
                                    </div>
                                    <div class="item-after">
                                        <div class="transaction-amount"
                                             :class="{ 'text-expense': item.transaction.type === TransactionType.Expense, 'text-income': item.transaction.type === TransactionType.Income }">
                                            <span v-if="item.transaction.type === TransactionType.Transfer && item.transaction.destinationAccountId === accountId">{{ getDisplayDestinationAmount(item.transaction) }}</span>
                                            <span v-else-if="item.transaction.type !== TransactionType.Transfer || item.transaction.destinationAccountId !== accountId">{{ getDisplaySourceAmount(item.transaction) }}</span>
                                        </div>
                                    </div>
                                </div>
                                <div class="item-text">
                                    <div class="transaction-description" v-if="item.transaction.comment">
                                        <span>{{ item.transaction.comment }}</span>
                                    </div>
                                </div>
                                <div class="item-footer">
                                    <div class="transaction-footer display-flex justify-content-space-between">
                                        <div class="flex-shrink-0">
                                            <span>{{ getDisplayTime(item.transaction) }}</span>
                                            <span style="margin-inline-start: 4px" v-if="!isSameAsDefaultTimezoneOffsetMinutes(item.transaction)">{{ `(${getDisplayTimezone(item.transaction)})` }}</span>
                                        </div>
                                        <div class="account-balance flex-shrink-1">
                                            <span>{{ isCurrentLiabilityAccount ? tt('Outstanding Balance') : tt('Balance') }}</span>
                                            <span style="margin-inline-start: 4px">{{ getDisplayAccountBalance(item.transaction) }}</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </template>
                    <f7-swipeout-actions :left="textDirection === TextDirection.RTL"
                                         :right="textDirection === TextDirection.LTR"
                                         v-if="item.type == 'transaction' && item.transaction">
                        <f7-swipeout-button color="primary" close
                                            :text="tt('Duplicate')"
                                            v-if="item.transaction.type !== TransactionType.ModifyBalance"
                                            @click="duplicateTransaction(item.transaction)"></f7-swipeout-button>
                        <f7-swipeout-button color="orange" close
                                            :text="tt('Edit')"
                                            v-if="item.transaction.editable && item.transaction.type !== TransactionType.ModifyBalance"
                                            @click="editTransaction(item.transaction)"></f7-swipeout-button>
                        <f7-swipeout-button color="red" class="padding-horizontal"
                                            v-if="item.transaction.editable"
                                            @click="removeTransaction(item.transaction, false)">
                            <f7-icon f7="trash"></f7-icon>
                        </f7-swipeout-button>
                    </f7-swipeout-actions>
                </f7-list-item>
            </ul>
        </f7-list>

        <f7-card v-if="finishQuery && showAccountBalanceTrendsCharts">
            <f7-card-header class="no-border display-block">
                <div class="statistics-chart-header display-flex full-line justify-content-space-between">
                    <div></div>
                    <div class="align-self-flex-end">
                        <span style="margin-inline-end: 4px;">{{ tt('Time Granularity') }}</span>
                        <f7-link :class="{ 'disabled': loading }" href="#" popover-open=".chart-data-date-aggregation-type-popover-menu">{{ chartDataDateAggregationTypeDisplayName }}</f7-link>
                    </div>
                </div>
            </f7-card-header>
            <f7-card-content style="margin-top: -14px" :padding="false">
                <account-balance-trends-bar-chart
                    :loading="loading"
                    :date-aggregation-type="chartDataDateAggregationType"
                    :fiscal-year-start="fiscalYearStart"
                    :items="reconciliationStatements?.transactions"
                    :account="currentAccount"
                />
            </f7-card-content>
        </f7-card>

        <f7-popover class="chart-data-date-aggregation-type-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="dateAggregationType.displayName"
                              :class="{ 'list-item-selected': chartDataDateAggregationType === dateAggregationType.type }"
                              :key="dateAggregationType.type"
                              v-for="dateAggregationType in allDateAggregationTypes"
                              @click="setChartDataDateAggregationType(dateAggregationType.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="chartDataDateAggregationType === dateAggregationType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <date-range-selection-sheet :title="tt('Custom Date Range')"
                                    :min-time="startTime"
                                    :max-time="endTime"
                                    v-model:show="showCustomDateRangeSheet"
                                    @dateRange:change="changeCustomDateFilter">
        </date-range-selection-sheet>

        <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                          :max-value="TRANSACTION_MAX_AMOUNT"
                          :currency="currentAccountCurrency"
                          :hint="tt('Please enter the new closing balance for the account')"
                          v-model:show="showNewClosingBalanceSheet"
                          v-model="newClosingBalance"
                          @update:model-value="updateClosingBalance"
        ></number-pad-sheet>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading }" @click="addTransaction()">{{ tt('Add Transaction') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': loading }" @click="updateClosingBalance(undefined)">{{ tt('Update Closing Balance') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading }" @click="reload(true)">{{ tt('Refresh') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ tt('Are you sure you want to delete this transaction?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="removeTransaction(transactionToDelete, true)">{{ tt('Delete') }}</f7-actions-button>
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
import { useI18nUIComponents, showLoading, hideLoading, onSwipeoutDeleted } from '@/lib/ui/mobile.ts';
import { useReconciliationStatementPageBase } from '@/views/base/accounts/ReconciliationStatementPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import { TextDirection } from '@/core/text.ts';
import { type TimeRangeAndDateType, DateRange, DateRangeScene } from '@/core/datetime.ts';
import { AccountType } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import { ChartDateAggregationType } from '@/core/statistics.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '@/consts/transaction.ts';
import { type TransactionReconciliationStatementResponseItemWithInfo } from '@/models/transaction.ts';

import { isDefined, isEquals, findDisplayNameByType } from '@/lib/common.ts';
import {
    getCurrentUnixTime,
    getDateTypeByDateRange,
    getDateTypeByBillingCycleDateRange,
    getDateRangeByDateType,
    getDateRangeByBillingCycleDateType
} from '@/lib/datetime.ts';

interface ReconciliationStatementVirtualListData {
    items: ReconciliationStatementVirtualListItem[],
    topPosition: number
}

interface ReconciliationStatementVirtualListItem {
    index: number;
    type: ReconciliationStatementVirtualListItemType;
    displayDate?: string;
    transaction?: TransactionReconciliationStatementResponseItemWithInfo;
}

type ReconciliationStatementVirtualListItemType = 'transaction' | 'date';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const {
    tt,
    getCurrentLanguageTextDirection,
    getAllDateRanges,
    formatNumberToLocalizedNumerals
} = useI18n();

const { showAlert, showToast, routeBackOnError } = useI18nUIComponents();

const {
    accountId,
    startTime,
    endTime,
    reconciliationStatements,
    firstDayOfWeek,
    fiscalYearStart,
    allDateAggregationTypes,
    isCurrentLiabilityAccount,
    currentAccount,
    currentAccountCurrency,
    displayStartDateTime,
    displayEndDateTime,
    displayTotalInflows,
    displayTotalOutflows,
    displayTotalBalance,
    displayOpeningBalance,
    displayClosingBalance,
    setReconciliationStatements,
    getDisplayDate,
    getDisplayTime,
    isSameAsDefaultTimezoneOffsetMinutes,
    getDisplayTimezone,
    getDisplaySourceAmount,
    getDisplayDestinationAmount,
    getDisplayAccountBalance
} = useReconciliationStatementPageBase();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionsStore = useTransactionsStore();

const finishQuery = ref<boolean>(false);
const loading = ref<boolean>(false);
const loadingError = ref<unknown | null>(null);
const queryDateRangeType = ref<number>(DateRange.ThisMonth.type);
const showAccountBalanceTrendsCharts = ref<boolean>(false);
const chartDataDateAggregationType = ref<number>(ChartDateAggregationType.Day.type);
const transactionToDelete = ref<TransactionReconciliationStatementResponseItemWithInfo | null>(null);
const newClosingBalance = ref<number>(0);
const showCustomDateRangeSheet = ref<boolean>(false);
const showNewClosingBalanceSheet = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const showDeleteActionSheet = ref<boolean>(false);
const virtualDataItems = ref<ReconciliationStatementVirtualListData>({
    items: [],
    topPosition: 0
});

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const validQuery = computed(() => currentAccount.value && currentAccount.value.type === AccountType.SingleAccount.type);
const allAvailableDateRanges = computed(() => getAllDateRanges(DateRangeScene.Normal, true, !!accountsStore.getAccountStatementDate(accountId.value)));

const allReconciliationStatementVirtualListItems = computed<ReconciliationStatementVirtualListItem[]>(() => {
    const ret: ReconciliationStatementVirtualListItem[] = [];

    if (!reconciliationStatements.value || !reconciliationStatements.value.transactions || reconciliationStatements.value.transactions.length < 1) {
        return ret;
    }

    let index = 0;
    let lastDisplayDate: string | null = null;

    for (const transaction of reconciliationStatements.value.transactions) {
        const displayDate = getDisplayDate(transaction);

        if (lastDisplayDate !== displayDate) {
            lastDisplayDate = displayDate;
            ret.push({
                index: index++,
                type: 'date',
                displayDate: displayDate
            });
        }

        ret.push({
            index: index++,
            type: 'transaction',
            transaction: transaction
        });
    }

    return ret;
});

const chartDataDateAggregationTypeDisplayName = computed<string>(() => {
    return findDisplayNameByType(allDateAggregationTypes.value, chartDataDateAggregationType.value) || tt('Unknown');
});

function getTransactionDomId(transaction: TransactionReconciliationStatementResponseItemWithInfo): string {
    return 'transaction_' + transaction.id;
}

function init(): void {
    const query = props.f7route.query;
    const defaultDateRange = getDateRangeByDateType(queryDateRangeType.value, firstDayOfWeek.value, fiscalYearStart.value);

    finishQuery.value = false;
    loading.value = false;
    accountId.value = query['accountId'] || '';
    startTime.value = defaultDateRange?.minTime || 0;
    endTime.value = defaultDateRange?.maxTime || 0;
    reconciliationStatements.value = undefined;

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false })
    ]).catch(error => {
        loadingError.value = error;
        showToast(error.message || error);
    });
}

function changeDateFilter(dateRangeType: number): void {
    if (dateRangeType === DateRange.Custom.type) {
        showCustomDateRangeSheet.value = true;
        return;
    }

    let dateRange: TimeRangeAndDateType | null = null;

    if (DateRange.isBillingCycle(dateRangeType)) {
        dateRange = getDateRangeByBillingCycleDateType(dateRangeType, firstDayOfWeek.value, fiscalYearStart.value, accountsStore.getAccountStatementDate(accountId.value));
    } else {
        dateRange = getDateRangeByDateType(dateRangeType, firstDayOfWeek.value, fiscalYearStart.value);
    }

    if (!dateRange) {
        return;
    }

    queryDateRangeType.value = dateRange.dateType;
    startTime.value = dateRange.minTime;
    endTime.value = dateRange.maxTime;
}

function changeCustomDateFilter(minTime: number, maxTime: number): void {
    if (!minTime || !maxTime) {
        return;
    }

    let dateType: number | null = getDateTypeByBillingCycleDateRange(minTime, maxTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal, accountsStore.getAccountStatementDate(accountId.value));

    if (!dateType) {
        dateType = getDateTypeByDateRange(minTime, maxTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);
    }

    queryDateRangeType.value = dateType;
    startTime.value = minTime;
    endTime.value = maxTime;
    showCustomDateRangeSheet.value = false;
}

function reload(force: boolean): void {
    finishQuery.value = true;
    loading.value = true;

    transactionsStore.getReconciliationStatements({
        accountId: accountId.value,
        startTime: startTime.value,
        endTime: endTime.value
    }).then(result => {
        if (force) {
            if (isEquals(reconciliationStatements.value, result)) {
                showToast('Data is up to date');
            } else {
                showToast('Data has been updated');
            }
        }

        loading.value = false;
        setReconciliationStatements(result);
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function addTransaction(): void {
    props.f7router.navigate(`/transaction/add?accountId=${accountId.value}`);
}

function duplicateTransaction(transaction: TransactionReconciliationStatementResponseItemWithInfo): void {
    props.f7router.navigate(`/transaction/add?id=${transaction.id}&type=${transaction.type}`);
}

function editTransaction(transaction: TransactionReconciliationStatementResponseItemWithInfo): void {
    props.f7router.navigate(`/transaction/edit?id=${transaction.id}&type=${transaction.type}`);
}

function updateClosingBalance(balance?: number): void {
    let currentClosingBalance = reconciliationStatements.value?.closingBalance ?? 0;

    if (isCurrentLiabilityAccount.value) {
        currentClosingBalance = -currentClosingBalance;
    }

    if (!isDefined(balance)) {
        newClosingBalance.value = currentClosingBalance;
        showNewClosingBalanceSheet.value = true;
        return;
    }

    const currentUnixTime = getCurrentUnixTime();
    let setTransactionTime = false;
    let newTransactionTime: number | undefined = undefined;

    if (endTime.value < currentUnixTime) {
        setTransactionTime = true;
        newTransactionTime = endTime.value;
    } else if (currentUnixTime < startTime.value) {
        setTransactionTime = true;
        newTransactionTime = startTime.value;
    }

    let newTransactionType: TransactionType = isCurrentLiabilityAccount.value ? TransactionType.Expense : TransactionType.Income;
    let newTransactionAmount: number = balance - currentClosingBalance;

    if (newTransactionAmount < 0) {
        newTransactionType = isCurrentLiabilityAccount.value ? TransactionType.Income : TransactionType.Expense;
        newTransactionAmount = -newTransactionAmount;
    }

    const params: string[] = [];

    if (setTransactionTime) {
        params.push(`time=${newTransactionTime}`);
    }

    params.push(`type=${newTransactionType}`);
    params.push(`amount=${newTransactionAmount}`);
    params.push(`accountId=${accountId.value}`);
    params.push(`noTransactionDraft=true`);

    props.f7router.navigate(`/transaction/add?${params.join('&')}`);
}

function removeTransaction(transaction: TransactionReconciliationStatementResponseItemWithInfo | null, confirm: boolean): void {
    if (!transaction) {
        showAlert('An error occurred');
        return;
    }

    if (!confirm) {
        transactionToDelete.value = transaction;
        showDeleteActionSheet.value = true;
        return;
    }

    showDeleteActionSheet.value = false;
    transactionToDelete.value = null;
    showLoading();

    transactionsStore.deleteTransaction({
        transaction: transaction,
        defaultCurrency: currentAccountCurrency.value,
        beforeResolve: (done) => {
            onSwipeoutDeleted(getTransactionDomId(transaction), done);
        }
    }).then(() => {
        hideLoading();
        reload(false);
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function setChartDataDateAggregationType(type: number): void {
    chartDataDateAggregationType.value = type;
}

function renderExternal(vl: unknown, vlData: ReconciliationStatementVirtualListData): void {
    virtualDataItems.value = vlData;
}

function onPageAfterIn(): void {
    if (finishQuery.value && transactionsStore.transactionReconciliationStatementStateInvalid) {
        reload(false);
    }

    routeBackOnError(props.f7router, loadingError);
}

init();
</script>

<style>
.list.reconciliation-statement-list li.reconciliation-statement-transaction-date:first-child {
    border-radius: var(--f7-list-inset-border-radius) var(--f7-list-inset-border-radius) 0 0;
}

.list.reconciliation-statement-list li.reconciliation-statement-transaction-date {
    display: flex;
    align-items: center;
    align-content: center;
    padding-top: 0;
    padding-bottom: 0;
    height: var(--f7-list-group-title-height);
    line-height: var(--f7-list-group-title-height);
    background-color: var(--f7-list-group-title-bg-color);
}

.list.reconciliation-statement-list li.reconciliation-statement-transaction-date > .item-content {
    padding-inline-start: 0 !important;
}

.list.reconciliation-statement-list li.reconciliation-statement-transaction-date > .item-content > .item-inner {
    padding-inline-start: calc(var(--f7-list-item-padding-horizontal) + var(--f7-safe-area-left));
}

.list.reconciliation-statement-list li.reconciliation-statement-transaction-date > .item-content > .item-inner:after {
    background-color: inherit;
}

.list.reconciliation-statement-list li.transaction-info.last-transaction-of-day > .item-link > .item-content > .item-inner:after,
.list.reconciliation-statement-list li.transaction-info.last-transaction-of-day > .swipeout-content > .item-link > .item-content > .item-inner:after {
    background-color: inherit;
}

.list.reconciliation-statement-list li.transaction-info .account-balance {
    margin-inline-start: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
</style>
