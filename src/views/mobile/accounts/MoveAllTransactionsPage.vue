<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Move All Transactions')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': !fromAccount || !toAccountId || fromAccount?.id === toAccountId || !toAccountName || !isToAccountNameValid || moving }" :text="tt('Confirm')" @click="confirm"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card
            :title="tt('Are you sure you want to move all transactions?')"
            :content="tt('format.misc.moveTransactionsInAccountTip', { fromAccount: fromAccount?.name, toAccount: displayToAccountName })">
        </f7-card>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Target Account" title="Unspecified"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Confirm Target Account Name" title="Unspecified"></f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'disabled': moving || !allVisibleAccounts.length }"
                :header="tt('Target Account')"
                :title="Account.findAccountNameById(allAccounts, toAccountId, tt('Unspecified'))"
                @click="showAccountSheet = true"
            >
                <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                      primary-title-field="name" primary-footer-field="displayBalance"
                                                      primary-icon-field="icon" primary-icon-type="account"
                                                      primary-sub-items-field="accounts"
                                                      :primary-title-i18n="true"
                                                      secondary-key-field="id" secondary-value-field="id"
                                                      secondary-title-field="name" secondary-footer-field="displayBalance"
                                                      secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                      :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                      :items="allVisibleCategorizedAccounts"
                                                      v-model:show="showAccountSheet"
                                                      v-model="toAccountId">
                </two-column-list-item-selection-sheet>
            </f7-list-item>

            <f7-list-input
                type="text"
                clear-button
                :label="tt('Confirm Target Account Name')"
                :placeholder="tt('Please re-enter the target account name to confirm')"
                v-model:value="toAccountName"
            ></f7-list-input>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref} from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useMoveAllTransactionsPageBase } from '@/views/base/accounts/MoveAllTransactionsPageBase.ts'

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import { Account } from '@/models/account.ts';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const { tt } = useI18n();

const { showToast, routeBackOnError } = useI18nUIComponents();

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

const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);
const showAccountSheet = ref<boolean>(false);

function init(): void {
    const query = props.f7route.query;
    const fromAccountId = query['fromAccountId'];

    if (!fromAccountId) {
        showToast('Parameter Invalid');
        loadingError.value = 'Parameter Invalid';
        return;
    }

    accountsStore.loadAllAccounts({
        force: false
    }).then(() => {
        loading.value = false;

        const account = accountsStore.allAccountsMap[fromAccountId];

        if (!account) {
            showToast('Parameter Invalid');
            loadingError.value = 'Parameter Invalid';
            return;
        }

        fromAccount.value = account;
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function confirm(): void {
    const router = props.f7router;

    if (!fromAccount.value || !toAccountId.value || fromAccount.value?.id === toAccountId.value || !toAccountName.value || !isToAccountNameValid.value) {
        return;
    }

    moving.value = true;
    showLoading(() => moving.value);

    transactionsStore.moveAllTransactionsBetweenAccounts({
        fromAccountId: fromAccount.value.id,
        toAccountId: toAccountId.value
    }).then(() => {
        moving.value = false;
        hideLoading();

        showToast('All transactions in this account has been moved.');
        router.back();
    }).catch(error => {
        moving.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

init();
</script>
