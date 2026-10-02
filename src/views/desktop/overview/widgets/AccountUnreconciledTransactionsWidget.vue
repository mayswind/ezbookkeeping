<template>
    <v-card class="overview-widget overview-widget--list overview-widget--accounts h-100" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || tt('Account Unreconciled Transactions')" :icon="mdiCreditCardOutline" />
        </template>

        <v-card-text class="overview-widget__body overview-widget__list-body">
            <v-list class="py-0" density="compact" v-if="displayAccounts.length">
                <v-list-item class="overview-widget__list-item" :key="item.account.id" :title="item.account.name"
                             @click="showReconciliationStatementDialog(item.account.id)"
                             v-for="item in displayAccounts">
                    <template #prepend>
                        <span class="overview-widget__item-icon">
                            <ItemIcon size="24px" :icon-type="getAccountIconType(item.account.iconType)" :icon-id="item.account.icon" :color="item.account.color" />
                        </span>
                    </template>
                    <template #append>
                        <span class="overview-widget__list-amount overview-widget__amount">{{ formatNumberToLocalizedNumerals(item.count) }}</span>
                    </template>
                </v-list-item>
            </v-list>
            <div v-if="loading && !displayAccounts.length">
                <v-skeleton-loader class="skeleton-no-margin mx-2 py-3" style="margin-bottom: 1px" type="text" :key="idx" :loading="true" v-for="idx in props.itemCount"></v-skeleton-loader>
            </div>
            <div class="overview-widget__empty" v-else-if="!loading && !displayAccounts.length">
                <v-icon :icon="mdiCreditCardOutline" size="32" />
                <span>{{ noDataText }}</span>
            </div>
        </v-card-text>

        <reconciliation-statement-dialog ref="reconciliationStatementDialog"
                                         @update:last-reconciled-time="onUpdateLastReconciledTime"
                                         @error="onShowReconciliationStatementError" />
        <snack-bar ref="snackbar" />
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import ReconciliationStatementDialog from '@/views/desktop/accounts/list/dialogs/ReconciliationStatementDialog.vue';

import { useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import {
    type CommonAccountUnreconciledTransactionsWidgetProps,
    useAccountUnreconciledTransactionsWidgetBase
} from '@/views/base/overview/AccountUnreconciledTransactionsWidgetBase.ts';

import { getAccountIconType } from '@/lib/icon.ts';

import {
    mdiCreditCardOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;
type ReconciliationStatementDialogType = InstanceType<typeof ReconciliationStatementDialog>;

const emit = defineEmits<{
    (e: 'refresh'): void
}>();

interface DesktopAccountUnreconciledTransactionsWidgetProps extends CommonAccountUnreconciledTransactionsWidgetProps {
    editing?: boolean
}

const props = defineProps<DesktopAccountUnreconciledTransactionsWidgetProps>();

const { tt, formatNumberToLocalizedNumerals } = useI18n();

const { noDataText, displayAccounts, getTransactionListDateRange } = useAccountUnreconciledTransactionsWidgetBase(props);

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const reconciliationStatementDialog = useTemplateRef<ReconciliationStatementDialogType>('reconciliationStatementDialog');

function showReconciliationStatementDialog(accountId: string): void {
    if (props.editing || props.loading) {
        return;
    }

    const dateRange = getTransactionListDateRange(accountId);

    reconciliationStatementDialog.value?.open({
        accountId: accountId,
        startTime: dateRange?.minTime ?? 0,
        endTime: dateRange?.maxTime ?? 0
    });
}

function onUpdateLastReconciledTime(): void {
    emit('refresh');
}

function onShowReconciliationStatementError(message: string): void {
    snackbar.value?.showError(message);
}
</script>
