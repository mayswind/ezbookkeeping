<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Reconciliation Statement')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="" v-if="!finishQuery"></f7-link>
                <f7-link :class="{ 'disabled': !validQuery }" :text="tt('Next')" @click="reload(false)" v-if="!finishQuery"></f7-link>
                <f7-link style="color: transparent" :text="tt('Next')" v-if="finishQuery"></f7-link>
                <f7-link :class="{ 'disabled': loading }" icon-f7="ellipsis" v-if="finishQuery" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

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
                    <div v-if="((dateRange.isBillingCycle || dateRange.type === DateRange.Custom.type) && queryDateRangeType === dateRange.type) && startTime && endTime">
                        <span>{{ displayStartTime }}</span>
                        <span>&nbsp;-&nbsp;</span>
                        <br/>
                        <span>{{ displayEndTime }}</span>
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
            <f7-list-item :title="tt('Total Transactions')" :after="reconciliationStatements.length || '0'"></f7-list-item>
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
                 v-if="finishQuery && loading">
            <ul>
                <f7-list-item chevron-center media-item
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
                                        <span style="margin-left: 4px">0.00 USD</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </f7-list-item>
            </ul>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical"
                 v-if="finishQuery && !loading && (!allReconciliationStatementVirtualListItems || !allReconciliationStatementVirtualListItems.length)">
            <f7-list-item :title="tt('No transaction data')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers media-list virtual-list
                 class="margin-vertical transaction-info-list reconciliation-statement-list"
                 :virtual-list-params="{ items: allReconciliationStatementVirtualListItems, renderExternal, height: 'auto' }"
                 v-if="finishQuery && !loading && allReconciliationStatementVirtualListItems && allReconciliationStatementVirtualListItems.length">
            <ul>
                <f7-list-item
                    chevron-center
                    media-item
                    :key="item.index"
                    :class="{ 'transaction-info': item.type == 'transaction', 'last-transaction-of-day': allReconciliationStatementVirtualListItems[item.index + 1] && allReconciliationStatementVirtualListItems[item.index + 1].type === 'date', 'reconciliation-statement-transaction-date': item.type == 'date' }"
                    :style="`top: ${virtualDataItems.topPosition}px`"
                    :virtual-list-index="item.index"
                    :link="item.type == 'transaction' && item.transaction && item.transaction.type !== TransactionType.ModifyBalance ? `/transaction/detail?id=${item.transaction?.id}&type=${item.transaction.type}` : null"
                    v-for="item in virtualDataItems.items"
                >
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
                                          :icon-id="allCategoriesMap[item.transaction.categoryId]?.icon"
                                          :color="allCategoriesMap[item.transaction.categoryId]?.color"
                                          v-if="allCategoriesMap[item.transaction.categoryId] && allCategoriesMap[item.transaction.categoryId]?.color"></ItemIcon>
                                <f7-icon v-else-if="!allCategoriesMap[item.transaction.categoryId] || !allCategoriesMap[item.transaction.categoryId]?.color"
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
                                        <span v-else-if="item.transaction.type !== TransactionType.ModifyBalance && allCategoriesMap[item.transaction.categoryId]">
                                            {{ allCategoriesMap[item.transaction.categoryId].name }}
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
                                        <span v-if="item.transaction.utcOffset !== currentTimezoneOffsetMinutes">{{ `(${getDisplayTimezone(item.transaction)})` }}</span>
                                    </div>
                                    <div class="account-balance flex-shrink-1">
                                        <span>{{ isCurrentLiabilityAccount ? tt('Outstanding Balance') : tt('Balance') }}</span>
                                        <span style="margin-left: 4px">{{ getDisplayAccountBalance(item.transaction) }}</span>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </f7-list-item>
            </ul>
        </f7-list>

        <date-range-selection-sheet :title="tt('Custom Date Range')"
                                    :min-time="startTime"
                                    :max-time="endTime"
                                    v-model:show="showCustomDateRangeSheet"
                                    @dateRange:change="changeCustomDateFilter">
        </date-range-selection-sheet>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading }" @click="addTransaction()">{{ tt('Add Transaction') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': loading }" @click="reload(true)">{{ tt('Refresh') }}</f7-actions-button>
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
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useReconciliationStatementPageBase } from '@/views/base/transactions/ReconciliationStatementPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import { type TimeRangeAndDateType, DateRange, DateRangeScene } from '@/core/datetime.ts';
import { AccountType } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import { type TransactionReconciliationStatementResponseItem } from '@/models/transaction.ts';

import {
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
    transaction?: TransactionReconciliationStatementResponseItem;
}

type ReconciliationStatementVirtualListItemType = 'transaction' | 'date';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const { tt, getAllDateRanges, formatUnixTimeToLongDateTime } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const {
    accountId,
    startTime,
    endTime,
    reconciliationStatements,
    openingBalance,
    closingBalance,
    firstDayOfWeek,
    fiscalYearStart,
    currentTimezoneOffsetMinutes,
    isCurrentLiabilityAccount,
    allCategoriesMap,
    currentAccount,
    displayStartDateTime,
    displayEndDateTime,
    displayTotalOutflows,
    displayTotalInflows,
    displayTotalBalance,
    displayOpeningBalance,
    displayClosingBalance,
    getDisplayDate,
    getDisplayTime,
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
const showCustomDateRangeSheet = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const virtualDataItems = ref<ReconciliationStatementVirtualListData>({
    items: [],
    topPosition: 0
});

const validQuery = computed(() => currentAccount.value && currentAccount.value.type === AccountType.SingleAccount.type);
const allAvailableDateRanges = computed(() => getAllDateRanges(DateRangeScene.Normal, true, !!accountsStore.getAccountStatementDate(accountId.value)));
const displayStartTime = computed<string>(() => formatUnixTimeToLongDateTime(startTime.value));
const displayEndTime = computed<string>(() => formatUnixTimeToLongDateTime(endTime.value));

const allReconciliationStatementVirtualListItems = computed<ReconciliationStatementVirtualListItem[]>(() => {
    const ret: ReconciliationStatementVirtualListItem[] = [];

    if (!reconciliationStatements.value || reconciliationStatements.value.length < 1) {
        return ret;
    }

    let index = 0;
    let lastDisplayDate: string | null = null;

    for (let i = 0; i < reconciliationStatements.value.length; i++) {
        const transaction = reconciliationStatements.value[i];
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

function init(): void {
    const query = props.f7route.query;
    const defaultDateRange = getDateRangeByDateType(queryDateRangeType.value, firstDayOfWeek.value, fiscalYearStart.value);

    finishQuery.value = false;
    loading.value = false;
    accountId.value = query['accountId'] || '';
    startTime.value = defaultDateRange?.minTime || 0;
    endTime.value = defaultDateRange?.maxTime || 0;
    reconciliationStatements.value = [];

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
            showToast('Data has been updated');
        }

        loading.value = false;
        reconciliationStatements.value = result.transactions;
        openingBalance.value = result.openingBalance;
        closingBalance.value = result.closingBalance;
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

function renderExternal(vl: unknown, vlData: ReconciliationStatementVirtualListData): void {
    virtualDataItems.value = vlData;
}

function onPageAfterIn(): void {
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
    padding-left: 0 !important;
}

.list.reconciliation-statement-list li.reconciliation-statement-transaction-date > .item-content > .item-inner {
    padding-left: calc(var(--f7-list-item-padding-horizontal) + var(--f7-safe-area-left));
}

.list.reconciliation-statement-list li.reconciliation-statement-transaction-date > .item-content > .item-inner:after {
    background-color: inherit;
}

.list.reconciliation-statement-list li.transaction-info.last-transaction-of-day > .item-link > .item-content > .item-inner:after {
    background-color: inherit;
}

.list.reconciliation-statement-list li.transaction-info .account-balance {
    margin-left: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
</style>
