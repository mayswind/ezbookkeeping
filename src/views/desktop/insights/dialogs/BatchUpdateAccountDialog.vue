<template>
    <v-dialog width="600" :persistent="true" v-model="showState">
        <one-column-dialog-layout content-class="pt-0" :disabled="loading || submitting"
                                  :title="isDestinationAccount ? tt('Update Destination Accounts for Transactions') : tt('Update Accounts for Transactions')"
                                  :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #after-title>
                <v-btn density="compact" color="default" variant="text"
                       class="ms-2" :icon="true" :disabled="loading || submitting"
                       :loading="loading" @click="reload">
                    <template #loader>
                        <v-progress-circular indeterminate size="20"/>
                    </template>
                    <v-icon :icon="mdiRefresh" size="22" />
                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                </v-btn>
            </template>

            <template #content>
                <div class="mt-5">
                    <two-column-select primary-key-field="id" primary-value-field="category"
                                       primary-title-field="name" primary-footer-field="displayBalance"
                                       primary-icon-field="icon" primary-icon-type="account"
                                       primary-sub-items-field="accounts"
                                       :primary-title-i18n="true"
                                       secondary-key-field="id" secondary-value-field="id"
                                       secondary-title-field="name" secondary-footer-field="displayBalance"
                                       secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                       :disabled="loading || !allVisibleAccounts.length"
                                       :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                       :custom-selection-primary-text="getAccountDisplayName(accountId)"
                                       :label="!isDestinationAccount ? tt('Account') : tt('Destination Account')"
                                       :placeholder="!isDestinationAccount ? tt('Account') : tt('Destination Account')"
                                       :items="allVisibleCategorizedAccounts"
                                       v-model="accountId">
                    </two-column-select>
                </div>
            </template>

            <template #footer>
                <v-btn color="secondary" variant="tonal" :disabled="loading || submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                <v-spacer/>
                <v-btn :disabled="loading || submitting || updateIds.length < 1 || !accountId" @click="confirm">
                    {{ tt('OK') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                </v-btn>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';

import {
    mdiRefresh
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    getCategorizedAccountsWithDisplayBalance
} = useI18n();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionsStore = useTransactionsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const submitting = ref<boolean>(false);
const isDestinationAccount = ref<boolean>(false);
const updateIds = ref<string[]>([]);
const accountId = ref<string>('');

let resolveFunc: ((response: number) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
const customAccountCategoryOrder = computed<string>(() => settingsStore.appSettings.accountCategoryOrders);

const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value, customAccountCategoryOrder.value));

function getAccountDisplayName(accountId?: string): string {
    if (accountId) {
        return Account.findAccountNameById(allAccounts.value, accountId) || '';
    } else {
        return tt('None');
    }
}

function open(options: { isDestinationAccount: boolean; updateIds: string[] }): Promise<number> {
    isDestinationAccount.value = options.isDestinationAccount;
    updateIds.value = options.updateIds;
    accountId.value = '';
    submitting.value = false;
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function reload(): void {
    loading.value = true;

    accountsStore.loadAllAccounts({ force: true }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function confirm(): void {
    submitting.value = true;

    transactionsStore.batchUpdateTransactionAccounts({
        transactionIds: updateIds.value,
        accountId: accountId.value,
        isDestinationAccount: isDestinationAccount.value
    }).then(() => {
        submitting.value = false;
        showState.value = false;
        resolveFunc?.(updateIds.value.length);
    }).catch(error => {
        submitting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
