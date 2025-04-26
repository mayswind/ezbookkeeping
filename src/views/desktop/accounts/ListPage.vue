<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <span class="text-subtitle-2">{{ tt('Net assets') }}</span>
                            <p class="account-statistic-item-value text-income text-truncate mt-1 mb-3">
                                <span v-if="!loading || allAccountCount > 0">{{ netAssets }}</span>
                                <span v-else-if="loading && allAccountCount <= 0">
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-1" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <span class="text-subtitle-2">{{ tt('Total liabilities') }}</span>
                            <p class="account-statistic-item-value text-expense text-truncate mt-1 mb-3">
                                <span v-if="!loading || allAccountCount > 0">{{ totalLiabilities }}</span>
                                <span v-else-if="loading && allAccountCount <= 0">
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-1" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <span class="text-subtitle-2">{{ tt('Total assets') }}</span>
                            <p class="account-statistic-item-value mt-1">
                                <span v-if="!loading || allAccountCount > 0">{{ totalAssets }}</span>
                                <span v-else-if="loading && allAccountCount <= 0">
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-1" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="account-category-tabs my-4" direction="vertical"
                                :disabled="loading" v-model="activeAccountCategoryType">
                            <v-tab class="tab-text-truncate" :key="accountCategory.type" :value="accountCategory.type"
                                   v-for="accountCategory in AccountCategory.values()">
                                <ItemIcon icon-type="account" :icon-id="accountCategory.defaultAccountIconId" />
                                <div class="d-flex flex-column text-truncate ml-2">
                                    <small class="text-truncate text-left smaller" v-if="!loading || allAccountCount > 0">{{ accountCategoryTotalBalance(accountCategory) }}</small>
                                    <small class="text-truncate text-left smaller my-1" v-else-if="loading && allAccountCount <= 0">
                                        <v-skeleton-loader class="skeleton-no-margin"
                                                           width="100px" height="16" type="text" :loading="true"></v-skeleton-loader>
                                    </small>
                                    <span class="text-truncate text-left">{{ tt(accountCategory.name) }}</span>
                                </div>
                            </v-tab>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="accountPage">
                                <v-card variant="flat" min-height="780">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="mr-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="mdiMenu" size="24" />
                                            </v-btn>
                                            <span>{{ tt('Account List') }}</span>
                                            <v-btn class="ml-3" color="default" variant="outlined"
                                                   :disabled="loading" @click="add">{{ tt('Add') }}</v-btn>
                                            <v-btn class="ml-3" color="primary" variant="tonal"
                                                   :disabled="loading" @click="saveSortResult"
                                                   v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ml-2" :icon="true" :loading="loading" @click="reload(true)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="mdiRefresh" size="24" />
                                                <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-spacer/>
                                            <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                                                   :disabled="loading" :icon="true">
                                                <v-icon :icon="mdiDotsVertical" />
                                                <v-menu activator="parent">
                                                    <v-list>
                                                        <v-list-item :prepend-icon="mdiEyeOutline"
                                                                     :title="tt('Show Hidden Accounts')"
                                                                     v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                                        <v-list-item :prepend-icon="mdiEyeOffOutline"
                                                                     :title="tt('Hide Hidden Accounts')"
                                                                     v-if="showHidden" @click="showHidden = false"></v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <v-card-text class="accounts-overview-title text-truncate pt-0">
                                        <span class="accounts-overview-subtitle">{{ tt('Balance') }}</span>
                                        <v-skeleton-loader class="skeleton-no-margin ml-3 mb-2" width="120px" type="text" :loading="true" v-if="loading && activeAccountCategory && !hasAccount(activeAccountCategory)"></v-skeleton-loader>
                                        <span class="accounts-overview-amount ml-3" v-else-if="!loading || !activeAccountCategory || hasAccount(activeAccountCategory)">{{ activeAccountCategoryTotalBalance }}</span>
                                        <v-btn class="ml-2" density="compact" color="default" variant="text"
                                               :icon="true" :disabled="loading"
                                               @click="showAccountBalance = !showAccountBalance">
                                            <v-icon :icon="showAccountBalance ? mdiEyeOffOutline : mdiEyeOutline" size="20" />
                                            <v-tooltip activator="parent">{{ showAccountBalance ? tt('Hide Account Balance') : tt('Show Account Balance') }}</v-tooltip>
                                        </v-btn>
                                    </v-card-text>

                                    <v-row class="pl-6 pr-6 pr-md-8" v-if="loading && activeAccountCategory && !hasAccount(activeAccountCategory)">
                                        <v-col cols="12">
                                            <v-card border class="card-title-with-bg account-card mb-8 h-auto">
                                                <template #title>
                                                    <div class="account-title d-flex align-center">
                                                        <v-icon class="disabled mr-0" size="28px" :icon="mdiSquareRounded" />
                                                        <span class="account-name text-truncate ml-2">
                                                            <v-skeleton-loader class="skeleton-no-margin my-1"
                                                                               width="120px" type="text" :loading="true"></v-skeleton-loader>
                                                        </span>
                                                        <v-spacer/>
                                                        <span class="align-self-center">
                                                            <v-icon class="disabled" :icon="mdiDrag"/>
                                                        </span>
                                                    </div>
                                                </template>
                                                <v-divider/>
                                                <v-card-text>
                                                    <div class="d-flex account-toolbar align-center">
                                                        <v-btn class="px-2" density="comfortable" color="default" variant="text"
                                                               :disabled="true" :prepend-icon="mdiListBoxOutline">
                                                            {{ tt('Transaction List') }}
                                                        </v-btn>
                                                        <v-spacer/>
                                                        <span class="account-balance ml-2">
                                                            <v-skeleton-loader class="skeleton-no-margin"
                                                                               width="100px" type="text" :loading="true"></v-skeleton-loader>
                                                        </span>
                                                    </div>
                                                </v-card-text>
                                            </v-card>
                                        </v-col>
                                    </v-row>

                                    <v-row class="pl-5 pr-2 pr-md-4" v-if="!loading && activeAccountCategory && !hasAccount(activeAccountCategory)">
                                        <v-col cols="12">
                                            {{ tt('No available account') }}
                                        </v-col>
                                    </v-row>

                                    <v-row class="pl-6 pr-6 pr-md-8">
                                        <v-col cols="12">
                                            <draggable-list
                                                class="list-group"
                                                item-key="id"
                                                handle=".drag-handle"
                                                ghost-class="dragging-item"
                                                :disabled="activeAccountCategoryVisibleAccountCount <= 1"
                                                :list="allCategorizedAccountsMap[activeAccountCategory.type].accounts"
                                                v-if="activeAccountCategory && allCategorizedAccountsMap[activeAccountCategory.type] && allCategorizedAccountsMap[activeAccountCategory.type].accounts && allCategorizedAccountsMap[activeAccountCategory.type].accounts.length"
                                                @change="onMove"
                                            >
                                                <template #item="{ element }">
                                                    <div class="list-group-item">
                                                        <v-card border class="card-title-with-bg account-card mb-8 h-auto" v-if="showHidden || !element.hidden">
                                                            <template #title>
                                                                <div class="account-title d-flex align-baseline">
                                                                    <ItemIcon size="1.5rem" icon-type="account" :icon-id="element.icon"
                                                                              :color="element.color" :hidden-status="element.hidden" />
                                                                    <span class="account-name text-truncate ml-2">{{ element.name }}</span>
                                                                    <small class="account-currency text-truncate ml-2">
                                                                        {{ accountCurrency(element) }}
                                                                    </small>
                                                                    <v-spacer/>
                                                                    <span class="align-self-center">
                                                                        <v-icon :class="!loading && activeAccountCategoryVisibleAccountCount > 1 ? 'drag-handle' : 'disabled'"
                                                                                :icon="mdiDrag"/>
                                                                        <v-tooltip activator="parent" v-if="!loading && activeAccountCategoryVisibleAccountCount > 1">{{ tt('Drag to Reorder') }}</v-tooltip>
                                                                    </span>
                                                                </div>

                                                                <div class="mt-4" v-if="element.type === AccountType.MultiSubAccounts.type">
                                                                    <v-btn-toggle
                                                                        class="account-subaccounts"
                                                                        variant="outlined"
                                                                        color="primary"
                                                                        density="compact"
                                                                        mandatory="force"
                                                                        divided rounded="xl"
                                                                        :disabled="loading"
                                                                        v-model="activeSubAccount[element.id]"
                                                                    >
                                                                        <v-btn :value="''">
                                                                            <span>{{ tt('All') }}</span>
                                                                        </v-btn>
                                                                        <v-btn :key="subAccount.id" :value="subAccount.id"
                                                                               v-for="subAccount in element.subAccounts"
                                                                               v-show="showHidden || !subAccount.hidden">
                                                                            <ItemIcon size="1.5rem" icon-type="account" :icon-id="subAccount.icon"
                                                                                      :color="subAccount.color" :hidden-status="subAccount.hidden" />
                                                                            <span class="ml-2">{{ subAccount.name }}</span>
                                                                        </v-btn>
                                                                    </v-btn-toggle>
                                                                </div>
                                                            </template>

                                                            <v-divider/>

                                                            <v-card-text v-if="element.getAccountOrSubAccountComment(activeSubAccount[element.id])">
                                                                {{ element.getAccountOrSubAccountComment(activeSubAccount[element.id]) }}
                                                            </v-card-text>

                                                            <v-card-text>
                                                                <div class="d-flex account-toolbar align-center">
                                                                    <v-btn class="px-2" density="comfortable" color="default" variant="text"
                                                                           :disabled="loading" :prepend-icon="mdiListBoxOutline"
                                                                           :to="`/transaction/list?accountIds=${element.getAccountOrSubAccountId(activeSubAccount[element.id])}`">
                                                                        {{ tt('Transaction List') }}
                                                                    </v-btn>
                                                                    <v-btn class="px-2 ml-1" density="comfortable" color="default" variant="text"
                                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                           :disabled="loading"
                                                                           :prepend-icon="element.isAccountOrSubAccountHidden(activeSubAccount[element.id]) ? mdiEyeOutline : mdiEyeOffOutline"
                                                                           v-if="!activeSubAccount[element.id] || element.getSubAccount(activeSubAccount[element.id])"
                                                                           @click="hide(element, element.getAccountOrSubAccount(activeSubAccount[element.id]), !element.isAccountOrSubAccountHidden(activeSubAccount[element.id]))">
                                                                        {{ element.isAccountOrSubAccountHidden(activeSubAccount[element.id]) ? tt('Show') : tt('Hide') }}
                                                                    </v-btn>
                                                                    <v-btn class="px-2 ml-1" density="comfortable" color="default" variant="text"
                                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                           :disabled="loading" :prepend-icon="mdiPencilOutline"
                                                                           v-if="!activeSubAccount[element.id] || element.getSubAccount(activeSubAccount[element.id])"
                                                                           @click="edit(element)">
                                                                        {{ tt('Edit') }}
                                                                    </v-btn>
                                                                    <v-btn class="px-2 ml-1" density="comfortable" color="default" variant="text"
                                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                           :disabled="loading" :prepend-icon="mdiDeleteOutline"
                                                                           v-if="!activeSubAccount[element.id] || element.getSubAccount(activeSubAccount[element.id])"
                                                                           @click="remove(element)">
                                                                        {{ tt('Delete') }}
                                                                    </v-btn>
                                                                    <v-spacer/>
                                                                    <span class="account-balance ml-2">{{ accountBalance(element) }}</span>
                                                                </div>
                                                            </v-card-text>
                                                        </v-card>
                                                    </div>
                                                </template>
                                            </draggable-list>
                                        </v-col>
                                    </v-row>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <edit-dialog ref="editDialog" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import EditDialog from './list/dialogs/EditDialog.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';
import { useDisplay } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useAccountListPageBaseBase } from '@/views/base/accounts/AccountListPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';

import { AccountType, AccountCategory } from '@/core/account.ts';
import type { Account } from '@/models/account.ts';

import { isObject, isString } from '@/lib/common.ts';

import {
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiRefresh,
    mdiSquareRounded,
    mdiMenu,
    mdiPencilOutline,
    mdiDeleteOutline,
    mdiListBoxOutline,
    mdiDrag,
    mdiDotsVertical
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type EditDialogType = InstanceType<typeof EditDialog>;

const display = useDisplay();

const { tt, getCurrencyName, formatAmountWithCurrency, joinMultiText } = useI18n();

const {
    loading,
    showHidden,
    displayOrderModified,
    showAccountBalance,
    allAccounts,
    allCategorizedAccountsMap,
    allAccountCount,
    netAssets,
    totalAssets,
    totalLiabilities,
    accountCategoryTotalBalance
} = useAccountListPageBaseBase();

const accountsStore = useAccountsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const activeAccountCategoryType = ref<number>(AccountCategory.Default.type);
const activeTab = ref<string>('accountPage');
const activeSubAccount = ref<Record<string, string>>({});
const alwaysShowNav = ref<boolean>(display.mdAndUp.value);
const showNav = ref<boolean>(display.mdAndUp.value);

const activeAccountCategory = computed<AccountCategory | undefined>(() => AccountCategory.valueOf(activeAccountCategoryType.value));
const activeAccountCategoryTotalBalance = computed<string>(() => accountCategoryTotalBalance(activeAccountCategory.value));

const activeAccountCategoryVisibleAccountCount = computed<number>(() => {
    if (!activeAccountCategory.value || !allCategorizedAccountsMap.value[activeAccountCategory.value.type] || !allCategorizedAccountsMap.value[activeAccountCategory.value.type].accounts) {
        return 0;
    }

    const accounts = allCategorizedAccountsMap.value[activeAccountCategory.value.type].accounts;

    if (showHidden.value) {
        return accounts.length;
    }

    let visibleCount = 0;

    for (let i = 0; i < accounts.length; i++) {
        if (!accounts[i].hidden) {
            visibleCount++;
        }
    }

    return visibleCount;
});

function reload(force: boolean): void {
    loading.value = true;

    accountsStore.loadAllAccounts({
        force: force
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        if (allAccounts.value) {
            for (let i = 0; i < allAccounts.value.length; i++) {
                const account = allAccounts.value[i];

                if (account.type === AccountType.MultiSubAccounts.type && !activeSubAccount.value[account.id]) {
                    activeSubAccount.value[account.id] = '';
                }
            }
        }

        if (force) {
            snackbar.value?.showMessage('Account list has been updated');
        }
    }).catch(error => {
        loading.value = false;

        if (error && error.isUpToDate) {
            displayOrderModified.value = false;
        }

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function hasAccount(accountCategory: AccountCategory): boolean {
    return accountsStore.hasAccount(accountCategory, !showHidden.value);
}

function accountCurrency(account: Account): string | null {
    if (account.type === AccountType.SingleAccount.type) {
        return getCurrencyName(account.currency);
    } else if (account.type === AccountType.MultiSubAccounts.type) {
        const subAccountCurrencies = account.getSubAccountCurrencies(showHidden.value, activeSubAccount.value[account.id])
            .map(currencyCode => getCurrencyName(currencyCode));
        return joinMultiText(subAccountCurrencies);
    } else {
        return null;
    }
}

function accountBalance(account: Account): string | null {
    if (account.type === AccountType.SingleAccount.type) {
        const balance = accountsStore.getAccountBalance(showAccountBalance.value, account);

        if (!isString(balance)) {
            return '';
        }

        return formatAmountWithCurrency(balance, account.currency);
    } else if (account.type === AccountType.MultiSubAccounts.type) {
        const balanceResult = accountsStore.getAccountSubAccountBalance(showAccountBalance.value, showHidden.value, account, activeSubAccount.value[account.id]);

        if (!isObject(balanceResult)) {
            return '';
        }

        return formatAmountWithCurrency(balanceResult.balance, balanceResult.currency);
    } else {
        return null;
    }
}

function add(): void {
    editDialog.value?.open({
        category: activeAccountCategoryType.value
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function edit(account: Account): void {
    editDialog.value?.open({
        id: account.id,
        currentAccount: account
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        if (accountsStore.accountListStateInvalid && !loading.value) {
            reload(false);
        }
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function hide(account: Account, targetAccount: Account, hidden: boolean): void {
    loading.value = true;

    accountsStore.hideAccount({
        account: targetAccount,
        hidden: hidden
    }).then(() => {
        if (hidden && !showHidden.value && activeSubAccount.value[account.id]) {
            activeSubAccount.value[account.id] = '';
        }

        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function remove(account: Account): void {
    if (activeSubAccount.value[account.id]) {
        const subAccount: Account | null = account.getSubAccount(activeSubAccount.value[account.id]);

        if (!subAccount) {
            snackbar.value?.showMessage('Unable to delete this sub-account');
            return;
        }

        confirmDialog.value?.open('Are you sure you want to delete this sub-account?').then(() => {
            loading.value = true;

            accountsStore.deleteSubAccount({
                subAccount: subAccount
            }).then(() => {
                loading.value = false;
            }).catch(error => {
                loading.value = false;

                if (!error.processed) {
                    snackbar.value?.showError(error);
                }
            });
        });
    } else {
        confirmDialog.value?.open('Are you sure you want to delete this account?').then(() => {
            loading.value = true;

            accountsStore.deleteAccount({
                account: account
            }).then(() => {
                loading.value = false;
            }).catch(error => {
                loading.value = false;

                if (!error.processed) {
                    snackbar.value?.showError(error);
                }
            });
        });
    }
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        return;
    }

    loading.value = true;

    accountsStore.updateAccountDisplayOrders().then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function onMove(event: { moved: { element: { id: string }, oldIndex: number, newIndex: number } }): void {
    if (!event || !event.moved) {
        return;
    }

    const moveEvent = event.moved;

    if (!moveEvent.element || !moveEvent.element.id) {
        snackbar.value?.showMessage('Unable to move account');
        return;
    }

    accountsStore.changeAccountDisplayOrder({
        accountId: moveEvent.element.id,
        from: moveEvent.oldIndex,
        to: moveEvent.newIndex,
        updateListOrder: false,
        updateGlobalListOrder: true
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        snackbar.value?.showError(error);
    });
}

watch(() => display.mdAndUp.value, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});

reload(false);
</script>

<style>
.account-statistic-item-value {
    font-size: 1rem;
}

.account-category-tabs .v-tab.v-tab.v-btn {
    height: calc(var(--v-tabs-height) * 1.5);
}

.accounts-overview-title {
    line-height: 2rem !important;
    min-height: 52px;
    display: flex;
    align-items: flex-end;
}

.accounts-overview-amount {
    font-size: 1.5rem;
    color: rgba(var(--v-theme-on-background), var(--v-high-emphasis-opacity));
    overflow: hidden;
    text-overflow: ellipsis;
}

.accounts-overview-subtitle {
    font-size: 1rem;
    line-height: 1.75rem;
}

.account-card > .v-card-item {
    padding-top: 0.875rem;
    padding-bottom: 0.875rem;
}

.account-card .account-title {
    font-size: 1rem;
    line-height: 1.5rem !important;
}

.account-card .account-title .account-name {
    color: rgba(var(--v-theme-on-background), var(--v-high-emphasis-opacity));
}

.account-card .account-currency {
    font-size: 0.8rem;
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity));
}

.account-card .account-subaccounts {
    overflow-x: auto;
    white-space: nowrap;
}

.account-card .account-subaccounts.v-btn-toggle {
    height: auto !important;
    padding: 0;
    border: none;
}

.account-card .account-subaccounts.v-btn-toggle > .v-btn {
    border-color: rgba(var(--v-border-color), var(--v-border-opacity));
}

.account-card .account-subaccounts.v-btn-toggle > .v-btn:not(:first-child) {
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
    border-left: none;
}

.account-card .account-subaccounts.v-btn-toggle > .v-btn:not(:last-child) {
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
}

.account-card .account-subaccounts.v-btn-toggle > .v-btn {
    border: thin solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.account-card .account-subaccounts.v-btn-toggle button.v-btn {
    width: auto !important;
}

.account-card .account-toolbar {
    overflow-x: auto;
    white-space: nowrap;
}

.account-card .account-toolbar .hover-display {
    display: none;
}

.account-card .account-toolbar:hover .hover-display {
    display: grid;
}

.account-card .account-balance {
    font-size: 1.25rem;
    color: rgba(var(--v-theme-on-background), var(--v-high-emphasis-opacity));
}
</style>
