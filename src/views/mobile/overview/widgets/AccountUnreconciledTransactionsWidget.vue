<template>
    <f7-list strong inset dividers class="overview-widget-list no-margin-top margin-bottom" :class="{ 'skeleton-text': loading }">
        <f7-list-item group-title v-if="showTitle">
            <small>{{ title || tt('Account Unreconciled Transactions') }}</small>
        </f7-list-item>
        <f7-list-item :key="item.account.id" :title="item.account.name"
                      :after="formatNumberToLocalizedNumerals(item.count)"
                      :link="getReconciliationStatementUrl(item.account.id)"
                      v-for="item in displayAccounts">
            <template #media>
                <ItemIcon :icon-type="getAccountIconType(item.account.iconType)" :icon-id="item.account.icon" :color="item.account.color" />
            </template>
        </f7-list-item>
        <template v-if="loading && !displayAccounts.length">
            <f7-list-item link="#" title="Account" after="0" :key="idx" v-for="idx in itemCount">
                <template #media>
                    <f7-icon f7="app_fill"></f7-icon>
                </template>
            </f7-list-item>
        </template>
        <f7-list-item :title="noDataText" v-else-if="!loading && !displayAccounts.length"></f7-list-item>
    </f7-list>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';

import {
    type CommonAccountUnreconciledTransactionsWidgetProps,
    useAccountUnreconciledTransactionsWidgetBase
} from '@/views/base/overview/AccountUnreconciledTransactionsWidgetBase.ts';

import { DateRange } from '@/core/datetime.ts';

import { getAccountIconType } from '@/lib/icon.ts';

interface MobileAccountUnreconciledTransactionsWidgetProps extends CommonAccountUnreconciledTransactionsWidgetProps {
    showTitle: boolean;
}

const props = defineProps<MobileAccountUnreconciledTransactionsWidgetProps>();

const { tt, formatNumberToLocalizedNumerals } = useI18n();

const { noDataText, displayAccounts, getTransactionListDateRange } = useAccountUnreconciledTransactionsWidgetBase(props);

function getReconciliationStatementUrl(accountId: string): string {
    const dateRange = getTransactionListDateRange(accountId);

    if (!dateRange) {
        return `/account/reconciliation_statements?accountId=${accountId}&dateType=${DateRange.All.type}`;
    }

    return `/account/reconciliation_statements?accountId=${accountId}&dateType=${dateRange.dateType}`;
}
</script>
