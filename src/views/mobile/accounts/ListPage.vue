<template>
    <f7-page :ptr="!sortable" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Account List')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !allAccountCount }" v-if="!sortable" @click="showMoreActionSheet = true"></f7-link>
                <f7-link href="/account/add" icon-f7="plus" v-if="!sortable"></f7-link>
                <f7-link :text="tt('Done')" :class="{ 'disabled': displayOrderSaving }" @click="saveSortResult" v-else-if="sortable"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="account-overview-card" :class="{ 'skeleton-text': loading }">
            <f7-card-header class="display-block" style="padding-top: 120px;">
                <p class="no-margin">
                    <small class="card-header-content" v-if="loading">Net assets</small>
                    <small class="card-header-content" v-else-if="!loading">{{ tt('Net assets') }}</small>
                </p>
                <p class="no-margin">
                    <span class="net-assets" v-if="loading">0.00 USD</span>
                    <span class="net-assets" v-else-if="!loading">{{ netAssets }}</span>
                    <f7-link class="margin-left-half" @click="showAccountBalance = !showAccountBalance">
                        <f7-icon class="ebk-hide-icon" :f7="showAccountBalance ? 'eye_slash_fill' : 'eye_fill'"></f7-icon>
                    </f7-link>
                </p>
                <p class="no-margin">
                    <small class="account-overview-info" v-if="loading">
                        <span>Total assets | Total liabilities</span>
                    </small>
                    <small class="account-overview-info" v-else-if="!loading">
                        <span>{{ tt('Total assets') }}</span>
                        <span>{{ totalAssets }}</span>
                        <span>|</span>
                        <span>{{ tt('Total liabilities') }}</span>
                        <span>{{ totalLiabilities }}</span>
                    </small>
                </p>
            </f7-card-header>
        </f7-card>

        <div class="skeleton-text" v-if="loading">
            <f7-list strong inset dividers sortable class="list-has-group-title account-list margin-vertical"
                     :key="listIdx" v-for="listIdx in [ 1, 2, 3 ]">
                <f7-list-item group-title :sortable="false">
                    <small>
                        <span>Account Category</span>
                        <span style="margin-left: 10px">0.00 USD</span>
                    </small>
                </f7-list-item>
                <f7-list-item class="nested-list-item" after="0.00 USD" link="#"
                              :key="itemIdx" v-for="itemIdx in (listIdx === 1 ? [ 1 ] : [ 1, 2 ])">
                    <template #media>
                        <f7-icon f7="app_fill"></f7-icon>
                    </template>
                    <template #title>
                        <div class="display-flex padding-top-half padding-bottom-half">
                            <div class="nested-list-item-title">
                                <span>Account Name</span>
                            </div>
                        </div>
                    </template>
                </f7-list-item>
            </f7-list>
        </div>

        <f7-list strong inset dividers class="margin-vertical" v-if="!loading && noAvailableAccount">
            <f7-list-item :title="tt('No available account')"></f7-list-item>
        </f7-list>

        <div :key="accountCategory.type"
             v-for="accountCategory in AccountCategory.values()"
             v-show="!loading && ((showHidden && hasAccount(accountCategory, false)) || hasAccount(accountCategory, true))">
            <f7-list strong inset dividers sortable class="list-has-group-title account-list margin-vertical"
                     :sortable-enabled="sortable"
                     v-if="allCategorizedAccountsMap[accountCategory.type]"
                     @sortable:sort="onSort">
                <f7-list-item group-title :sortable="false">
                    <small>
                        <span>{{ tt(accountCategory.name) }}</span>
                        <span style="margin-left: 10px">{{ accountCategoryTotalBalance(accountCategory) }}</span>
                    </small>
                </f7-list-item>
                <f7-list-item swipeout
                              class="nested-list-item"
                              :id="getAccountDomId(account)"
                              :class="{ 'has-child-list-item': account.type === AccountType.MultiSubAccounts.type && hasVisibleSubAccount(account), 'actual-first-child': account.id === firstShowingIds.accounts[accountCategory.type], 'actual-last-child': account.id === lastShowingIds.accounts[accountCategory.type] }"
                              :after="account.type === AccountType.SingleAccount.type ? accountBalance(account) : ''"
                              :link="!sortable ? '/transaction/list?accountIds=' + account.id : null"
                              :key="account.id"
                              v-for="account in allCategorizedAccountsMap[accountCategory.type].accounts"
                              v-show="showHidden || !account.hidden"
                              @taphold="setSortable()"
                >
                    <template #media v-if="account.type !== AccountType.MultiSubAccounts.type || !hasVisibleSubAccount(account)">
                        <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color">
                            <f7-badge color="gray" class="right-bottom-icon" v-if="account.hidden">
                                <f7-icon f7="eye_slash_fill"></f7-icon>
                            </f7-badge>
                        </ItemIcon>
                    </template>

                    <template #title>
                        <div class="nested-list-item-inner display-flex padding-top-half padding-bottom-half">
                            <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"
                                      v-if="account.type === AccountType.MultiSubAccounts.type && hasVisibleSubAccount(account)">
                                <f7-badge color="gray" class="right-bottom-icon" v-if="account.hidden">
                                    <f7-icon f7="eye_slash_fill"></f7-icon>
                                </f7-badge>
                            </ItemIcon>
                            <div class="nested-list-item-title">
                                <span>{{ account.name }}</span>
                                <div class="item-footer" v-if="account.comment">{{ account.comment }}</div>
                            </div>
                            <div class="nested-list-item-after" v-if="account.type === AccountType.MultiSubAccounts.type">
                                <span>{{ accountBalance(account) }}</span>
                            </div>
                        </div>
                        <li v-if="account.type === AccountType.MultiSubAccounts.type">
                            <ul class="no-padding">
                                <f7-list-item class="no-sortable nested-list-item-child"
                                              :class="{ 'actual-first-child': subAccount.id === firstShowingIds.subAccounts[account.id], 'actual-last-child': subAccount.id === lastShowingIds.subAccounts[account.id] }"
                                              :id="getAccountDomId(subAccount)"
                                              :title="subAccount.name" :footer="subAccount.comment" :after="accountBalance(account, subAccount.id)"
                                              :link="!sortable ? '/transaction/list?accountIds=' + subAccount.id : null"
                                              :key="subAccount.id"
                                              v-for="subAccount in account.subAccounts"
                                              v-show="showHidden || !subAccount.hidden"
                                >
                                    <template #media>
                                        <ItemIcon icon-type="account" :icon-id="subAccount.icon" :color="subAccount.color">
                                            <f7-badge color="gray" class="right-bottom-icon" v-if="subAccount.hidden">
                                                <f7-icon f7="eye_slash_fill"></f7-icon>
                                            </f7-badge>
                                        </ItemIcon>
                                    </template>
                                </f7-list-item>
                            </ul>
                        </li>
                    </template>
                    <f7-swipeout-actions left v-if="sortable">
                        <f7-swipeout-button :color="account.hidden ? 'blue' : 'gray'" class="padding-left padding-right"
                                            overswipe close @click="hide(account, !account.hidden)">
                            <f7-icon :f7="account.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                        </f7-swipeout-button>
                    </f7-swipeout-actions>
                    <f7-swipeout-actions right v-if="!sortable">
                        <f7-swipeout-button color="orange" close :text="tt('Edit')" @click="edit(account)"></f7-swipeout-button>
                        <f7-swipeout-button color="primary" close :text="tt('More')" @click="showMoreActionSheetForAccount(account)"></f7-swipeout-button>
                        <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(account, false)">
                            <f7-icon f7="trash"></f7-icon>
                        </f7-swipeout-button>
                    </f7-swipeout-actions>
                </f7-list-item>
            </f7-list>
        </div>

        <f7-actions close-by-outside-click close-on-escape :opened="showAccountMoreActionSheet" @actions:closed="showAccountMoreActionSheet = false">
            <f7-actions-group v-if="accountForMoreActionSheet && accountForMoreActionSheet.type === AccountType.SingleAccount.type">
                <f7-actions-button @click="showReconciliationStatement(accountForMoreActionSheet)">{{ tt('Reconciliation Statement') }}</f7-actions-button>
            </f7-actions-group>
            <template v-if="accountForMoreActionSheet && accountForMoreActionSheet.type === AccountType.MultiSubAccounts.type">
                <f7-actions-group :key="subAccount.id"
                                  v-for="subAccount in accountForMoreActionSheet.subAccounts"
                                  v-show="showHidden || !subAccount.hidden">
                    <f7-actions-label>{{ subAccount.name }}</f7-actions-label>
                    <f7-actions-button @click="showReconciliationStatement(subAccount)">{{ tt('Reconciliation Statement') }}</f7-actions-button>
                </f7-actions-group>
            </template>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ tt('Sort') }}</f7-actions-button>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ tt('Show Hidden Accounts') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ tt('Hide Hidden Accounts') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group v-if="hasAnyVisibleAccount">
                <f7-actions-button @click="setAccountsIncludedInTotal()">{{ tt('Set Accounts Included in Total') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ tt('Are you sure you want to delete this account?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(accountToDelete, true)">{{ tt('Delete') }}</f7-actions-button>
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
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useAccountListPageBaseBase } from '@/views/base/accounts/AccountListPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';

import { AccountType, AccountCategory } from '@/core/account.ts';
import type { Account, AccountShowingIds } from '@/models/account.ts';

import { onSwipeoutDeleted } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showAlert, showToast, routeBackOnError } = useI18nUIComponents();

const {
    loading,
    showHidden,
    displayOrderModified,
    showAccountBalance,
    allCategorizedAccountsMap,
    allAccountCount,
    netAssets,
    totalAssets,
    totalLiabilities,
    accountCategoryTotalBalance,
    accountBalance
} = useAccountListPageBaseBase();

const accountsStore = useAccountsStore();

const loadingError = ref<unknown | null>(null);
const sortable = ref<boolean>(false);
const accountForMoreActionSheet = ref<Account | null>(null);
const accountToDelete = ref<Account | null>(null);
const showAccountMoreActionSheet = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const showDeleteActionSheet = ref<boolean>(false);
const displayOrderSaving = ref<boolean>(false);

const firstShowingIds = computed<AccountShowingIds>(() => accountsStore.getFirstShowingIds(showHidden.value));
const lastShowingIds = computed<AccountShowingIds>(() => accountsStore.getLastShowingIds(showHidden.value));
const hasAnyVisibleAccount = computed<boolean>(() => accountsStore.allVisibleAccountsCount > 0);
const noAvailableAccount = computed<boolean>(() => {
    if (showHidden.value) {
        return accountsStore.allAvailableAccountsCount < 1;
    } else {
        return accountsStore.allVisibleAccountsCount < 1;
    }
});

function hasAccount(accountCategory: AccountCategory, visibleOnly: boolean): boolean {
    return accountsStore.hasAccount(accountCategory, visibleOnly);
}

function hasVisibleSubAccount(account: Account): boolean {
    return accountsStore.hasVisibleSubAccount(showHidden.value, account);
}

function getAccountDomId(account: Account): string {
    return 'account_' + account.id;
}

function parseAccountIdFromDomId(domId: string): string | null {
    if (!domId || domId.indexOf('account_') !== 0) {
        return null;
    }

    return domId.substring(8); // account_
}

function init(): void {
    loading.value = true;

    accountsStore.loadAllAccounts({
        force: false
    }).then(() => {
        loading.value = false;
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function reload(done?: () => void): void {
    if (sortable.value) {
        done?.();
        return;
    }

    const force = !!done;

    accountsStore.loadAllAccounts({
        force: force
    }).then(() => {
        done?.();

        if (force) {
            showToast('Account list has been updated');
        }
    }).catch(error => {
        done?.();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function edit(account: Account): void {
    props.f7router.navigate('/account/edit?id=' + account.id);
}

function showMoreActionSheetForAccount(account: Account): void {
    accountForMoreActionSheet.value = account;
    showAccountMoreActionSheet.value = true;
}

function showReconciliationStatement(account: Account | null): void {
    if (!account) {
        showAlert('An error occurred');
        return;
    }

    props.f7router.navigate('/account/reconciliation_statements?accountId=' + account.id);
    showAccountMoreActionSheet.value = false;
    accountForMoreActionSheet.value = null;
}

function hide(account: Account, hidden: boolean): void {
    showLoading();

    accountsStore.hideAccount({
        account: account,
        hidden: hidden
    }).then(() => {
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function remove(account: Account | null, confirm: boolean): void {
    if (!account) {
        showAlert('An error occurred');
        return;
    }

    if (!confirm) {
        accountToDelete.value = account;
        showDeleteActionSheet.value = true;
        return;
    }

    showDeleteActionSheet.value = false;
    accountToDelete.value = null;
    showLoading();

    accountsStore.deleteAccount({
        account: account,
        beforeResolve: (done) => {
            onSwipeoutDeleted(getAccountDomId(account), done);
        }
    }).then(() => {
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function setSortable(): void {
    if (sortable.value) {
        return;
    }

    showHidden.value = true;
    sortable.value = true;
    displayOrderModified.value = false;
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        showHidden.value = false;
        sortable.value = false;
        return;
    }

    displayOrderSaving.value = true;
    showLoading();

    accountsStore.updateAccountDisplayOrders().then(() => {
        displayOrderSaving.value = false;
        hideLoading();

        showHidden.value = false;
        sortable.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        displayOrderSaving.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function setAccountsIncludedInTotal(): void {
    props.f7router.navigate('/settings/filter/account?type=accountListTotalAmount');
}

function onSort(event: { el: { id: string }; from: number; to: number }): void {
    if (!event || !event.el || !event.el.id) {
        showToast('Unable to move account');
        return;
    }

    const id = parseAccountIdFromDomId(event.el.id);

    if (!id) {
        showToast('Unable to move account');
        return;
    }

    accountsStore.changeAccountDisplayOrder({
        accountId: id,
        from: event.from - 1, // first item in the list is title, so the index need minus one
        to: event.to - 1,
        updateListOrder: true,
        updateGlobalListOrder: true
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        showToast(error.message || error);
    });
}

function onPageAfterIn(): void {
    if (accountsStore.accountListStateInvalid && !loading.value) {
        reload();
    }

    routeBackOnError(props.f7router, loadingError);
}

init();
</script>

<style>
.account-overview-card {
    background-color: var(--f7-color-yellow);
}

.dark .account-overview-card {
    background-color: var(--f7-theme-color);
}

.dark .account-overview-card a {
    color: var(--f7-text-color);
    opacity: 0.6;
}

.net-assets {
    font-size: 1.5em;
}

.account-overview-info {
    opacity: 0.6;
}

.account-overview-info > span {
    margin-right: 4px;
}

.account-overview-info > span:last-child {
    margin-right: 0;
}

.account-list {
    --f7-list-item-footer-font-size: var(--ebk-large-footer-font-size);
}

.account-list .item-footer {
    padding-top: 4px;
}
</style>
