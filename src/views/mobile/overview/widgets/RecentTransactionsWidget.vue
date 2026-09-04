<template>
    <f7-list strong inset dividers class="overview-widget-list transaction-info-list no-margin-top margin-bottom"
             :class="{ 'skeleton-text': loading }" :media-list="loading || !!transactions.length">
        <f7-list-item group-title v-if="showTitle">
            <small>{{ title || tt('Recent Transactions') }}</small>
        </f7-list-item>

        <template v-if="loading && !transactions.length">
            <f7-list-item link="#" chevron-center class="transaction-info" :key="idx" v-for="idx in itemCount">
                <template #media>
                    <div class="display-flex flex-direction-column transaction-date">
                        <span class="transaction-day width-100 flex-direction-column">DD</span>
                        <span class="transaction-day-of-week width-100 flex-direction-column">Sun</span>
                    </div>
                </template>
                <template #inner>
                    <div class="display-flex no-padding-horizontal">
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
                                <div class="transaction-footer">
                                    <span>HH:mm</span>
                                    <span>·</span>
                                    <span>Source Account</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </template>
            </f7-list-item>
        </template>

        <f7-list-item :title="tt('No data')" v-else-if="!loading && !transactions.length"></f7-list-item>

        <f7-list-item chevron-center class="transaction-info"
                      :link="`/transaction/detail?id=${transaction.id}&type=${transaction.type}`"
                      :key="transaction.id"
                      v-for="(transaction, idx) in transactions"
                      v-else-if="transactions.length">
            <template #media>
                <div class="display-flex flex-direction-column transaction-date" :style="getTransactionDateStyle(transaction, idx > 0 ? transactions[idx - 1] : undefined)">
                    <span class="transaction-day width-100 flex-direction-column">
                        {{ getDisplayDayOfMonth(transaction) }}
                    </span>
                    <span class="transaction-day-of-week width-100 flex-direction-column">
                        {{ getDisplayDayOfWeek(transaction) }}
                    </span>
                </div>
            </template>
            <template #inner>
                <div class="display-flex no-padding-horizontal">
                    <div class="item-media">
                        <div class="transaction-icon display-flex align-items-center">
                            <ItemIcon :icon-type="getCategoryIconType(getTransactionCategory(transaction)?.iconType)"
                                      :icon-id="getTransactionCategory(transaction)?.icon ?? ''"
                                      :color="getTransactionCategory(transaction)?.color"
                                      v-if="getTransactionCategory(transaction) && getTransactionCategory(transaction)?.color"></ItemIcon>
                            <f7-icon v-else-if="!getTransactionCategory(transaction) || !getTransactionCategory(transaction)?.color"
                                     f7="pencil_ellipsis_rectangle">
                            </f7-icon>
                        </div>
                    </div>
                    <div class="actual-item-inner">
                        <div class="item-title-row">
                            <div class="item-title">
                                <div class="transaction-category-name no-padding">
                                    <span>{{ getTransactionCategoryName(transaction) }}</span>
                                </div>
                            </div>
                            <div class="item-after">
                                <div class="transaction-amount"
                                     :class="{ 'text-expense': transaction.type === TransactionType.Expense, 'text-income': transaction.type === TransactionType.Income }">
                                    <span>{{ getDisplayAmount(transaction) }}</span>
                                </div>
                            </div>
                        </div>
                        <div class="item-text">
                            <div class="transaction-description" v-if="transaction.comment">
                                <span>{{ transaction.comment }}</span>
                            </div>
                        </div>
                        <div class="item-footer">
                            <div class="transaction-footer">
                                <span>{{ getDisplayTime(transaction) }}</span>
                                <template v-if="getSourceAccount(transaction)">
                                    <span>·</span>
                                    <span>{{ getSourceAccount(transaction)?.name }}</span>
                                </template>
                                <template v-if="getSourceAccount(transaction) && transaction.type === TransactionType.Transfer && getDestinationAccount(transaction) && getSourceAccount(transaction)?.id !== getDestinationAccount(transaction)?.id">
                                    <f7-icon class="transaction-account-arrow icon-with-direction" f7="arrow_right"></f7-icon>
                                    <span>{{ getDestinationAccount(transaction)?.name }}</span>
                                </template>
                            </div>
                        </div>
                    </div>
                </div>
            </template>
        </f7-list-item>
    </f7-list>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { TransactionType } from '@/core/transaction.ts';
import type { OverviewRecentTransactionsQuery } from '@/core/overview_layout.ts';
import { DISPLAY_HIDDEN_AMOUNT } from '@/consts/numeral.ts';

import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionInfoResponse } from '@/models/transaction.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTimeWithTimezoneOffset } from '@/lib/datetime.ts';
import { getCategoryIconType } from '@/lib/icon.ts';
import { getOverviewRecentTransactionsQuery } from '@/lib/overview_layout.ts';

const props = defineProps<{
    loading: boolean;
    title?: string;
    showTitle: boolean;
    editing?: boolean;
    itemCount: number;
    accountIds?: string[];
    categoryIds?: string[];
    tagFilter?: string;
    amountFilter?: string;
    keyword?: string;
}>();

const {
    tt,
    getWeekdayShortName,
    getCalendarDisplayDayOfMonthFromDateTime,
    formatDateTimeToShortTime,
    formatAmountToLocalizedNumeralsWithCurrency
} = useI18n();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const overviewStore = useOverviewStore();

const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
const recentTransactionsQuery = computed<OverviewRecentTransactionsQuery>(() => getOverviewRecentTransactionsQuery({
    itemCount: props.itemCount,
    accountIds: props.accountIds ?? [],
    categoryIds: props.categoryIds ?? [],
    tagFilter: props.tagFilter ?? '',
    amountFilter: props.amountFilter ?? '',
    keyword: props.keyword ?? ''
}));
const transactions = computed<TransactionInfoResponse[]>(() => overviewStore.getRecentTransactions(recentTransactionsQuery.value).slice(0, props.itemCount));

function getTransactionDateStyle(transaction: TransactionInfoResponse, previousTransaction: TransactionInfoResponse | undefined): Record<string, string> {
    if (!previousTransaction) {
        return {};
    }

    const transactionTime = getTransactionDateTime(transaction);
    const previousTransactionTime = getTransactionDateTime(previousTransaction);

    if (transactionTime.getGregorianCalendarDay() !== previousTransactionTime.getGregorianCalendarDay()) {
        return {};
    }

    return {
        color: 'transparent'
    };
}

function getTransactionDateTime(transaction: TransactionInfoResponse) {
    return parseDateTimeFromUnixTimeWithTimezoneOffset(transaction.time, transaction.utcOffset);
}

function getDisplayDayOfMonth(transaction: TransactionInfoResponse): string {
    return getCalendarDisplayDayOfMonthFromDateTime(getTransactionDateTime(transaction));
}

function getDisplayDayOfWeek(transaction: TransactionInfoResponse): string {
    return getWeekdayShortName(getTransactionDateTime(transaction).getWeekDay());
}

function getDisplayTime(transaction: TransactionInfoResponse): string {
    return formatDateTimeToShortTime(getTransactionDateTime(transaction));
}

function getDisplayAmount(transaction: TransactionInfoResponse): string {
    const account = getSourceAccount(transaction);
    const amount = !showAmountInHomePage.value || transaction.hideAmount ? DISPLAY_HIDDEN_AMOUNT : parseBigDecimal(transaction.sourceAmount);
    return formatAmountToLocalizedNumeralsWithCurrency(amount, account?.currency ?? '');
}

function getSourceAccount(transaction: TransactionInfoResponse): Account | undefined {
    return accountsStore.allAccountsMap[transaction.sourceAccountId];
}

function getDestinationAccount(transaction: TransactionInfoResponse): Account | undefined {
    return accountsStore.allAccountsMap[transaction.destinationAccountId];
}

function getTransactionCategory(transaction: TransactionInfoResponse): TransactionCategory | undefined {
    if (transaction.type === TransactionType.ModifyBalance) {
        return undefined;
    }

    return transactionCategoriesStore.allTransactionCategoriesMap[transaction.categoryId];
}

function getTransactionCategoryName(transaction: TransactionInfoResponse): string {
    if (transaction.type === TransactionType.ModifyBalance) {
        return tt('Modify Balance');
    }

    const category = getTransactionCategory(transaction);

    if (category) {
        return category.name;
    } else if (transaction.type === TransactionType.Income) {
        return tt('Income');
    } else if (transaction.type === TransactionType.Expense) {
        return tt('Expense');
    } else if (transaction.type === TransactionType.Transfer) {
        return tt('Transfer');
    } else {
        return tt('Transaction');
    }
}
</script>
