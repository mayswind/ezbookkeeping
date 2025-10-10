<template>
    <v-dialog width="640" :persistent="true" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4 text-wrap">{{ tt('Are you sure you want to move all transactions?') }}</h4>
                </div>
            </template>
            <v-card-text>{{ tt('format.misc.moveTransactionsInAccountTip', { fromAccount: fromAccount?.name, toAccount: displayToAccountName }) }}</v-card-text>
            <v-card-text class="mb-md-4 w-100 d-flex justify-center">
                <v-row>
                    <v-col cols="12" md="12">
                        <two-column-select primary-key-field="id" primary-value-field="category"
                                           primary-title-field="name" primary-footer-field="displayBalance"
                                           primary-icon-field="icon" primary-icon-type="account"
                                           primary-sub-items-field="accounts"
                                           :primary-title-i18n="true"
                                           secondary-key-field="id" secondary-value-field="id"
                                           secondary-title-field="name" secondary-footer-field="displayBalance"
                                           secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                           :disabled="loading || moving || !allVisibleAccounts.length"
                                           :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                           :label="tt('Target Account')"
                                           :placeholder="tt('Target Account')"
                                           :items="allVisibleCategorizedAccounts"
                                           :no-item-text="Account.findAccountNameById(allAccounts, toAccountId, tt('Unspecified'))"
                                           v-model="toAccountId">
                        </two-column-select>
                    </v-col>

                    <v-col cols="12" md="12">
                        <v-text-field type="text"
                                      persistent-placeholder
                                      :disabled="moving"
                                      :label="tt('Confirm Target Account Name')"
                                      :placeholder="tt('Please re-enter the target account name to confirm')"
                                      v-model="toAccountName"
                        />
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="!fromAccount || !toAccountId || fromAccount?.id === toAccountId || !toAccountName || !isToAccountNameValid || moving" @click="confirm">
                        {{ tt('Confirm') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="moving"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="moving" @click="cancel">
                        {{ tt('Cancel') }}
                    </v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useMoveAllTransactionsPageBase } from '@/views/base/accounts/MoveAllTransactionsPageBase.ts'

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import { Account } from '@/models/account.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const {
    moving,
    fromAccount,
    toAccountId,
    toAccountName,
    allAccounts,
    allVisibleAccounts,
    allVisibleCategorizedAccounts,
    displayToAccountName,
    isToAccountNameValid
} = useMoveAllTransactionsPageBase();

const accountsStore = useAccountsStore();
const transactionsStore = useTransactionsStore();

let resolveFunc: (() => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);

function init(): void {
    accountsStore.loadAllAccounts({
        force: false
    }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function open(account: Account): Promise<void> {
    showState.value = true;
    moving.value = false;
    fromAccount.value = account;
    toAccountId.value = '';
    toAccountName.value = '';

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (!fromAccount.value || !toAccountId.value || fromAccount.value?.id === toAccountId.value || !toAccountName.value || !isToAccountNameValid.value) {
        return;
    }

    moving.value = true;

    transactionsStore.moveAllTransactionsBetweenAccounts({
        fromAccountId: fromAccount.value.id,
        toAccountId: toAccountId.value
    }).then(() => {
        moving.value = false;

        resolveFunc?.();
        showState.value = false;
    }).catch(error => {
        moving.value = false;

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

init();
</script>
