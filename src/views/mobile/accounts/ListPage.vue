<template>
    <f7-page :ptr="!sortable" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Account List')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !allAccountCount }" v-if="!sortable" @click="showMoreActionSheet = true"></f7-link>
                <f7-link href="/account/add" icon-f7="plus" v-if="!sortable"></f7-link>
                <f7-link :text="$t('Done')" :class="{ 'disabled': displayOrderSaving }" @click="saveSortResult" v-else-if="sortable"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="account-overview-card" :class="{ 'skeleton-text': loading }">
            <f7-card-header class="display-block" style="padding-top: 120px;">
                <p class="no-margin">
                    <small class="card-header-content" v-if="loading">Net assets</small>
                    <small class="card-header-content" v-else-if="!loading">{{ $t('Net assets') }}</small>
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
                        <span>{{ $t('Total assets') }}</span>
                        <span>{{ totalAssets }}</span>
                        <span>|</span>
                        <span>{{ $t('Total liabilities') }}</span>
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
            <f7-list-item :title="$t('No available account')"></f7-list-item>
        </f7-list>

        <div :key="accountCategory.type"
             v-for="accountCategory in allAccountCategories"
             v-show="(showHidden && hasAccount(accountCategory, false)) || hasAccount(accountCategory, true)">
            <f7-list strong inset dividers sortable class="list-has-group-title account-list margin-vertical"
                     :sortable-enabled="sortable"
                     v-if="allCategorizedAccountsMap[accountCategory.type]"
                     @sortable:sort="onSort">
                <f7-list-item group-title :sortable="false">
                    <small>
                        <span>{{ $t(accountCategory.name) }}</span>
                        <span style="margin-left: 10px">{{ accountCategoryTotalBalance(accountCategory) }}</span>
                    </small>
                </f7-list-item>
                <f7-list-item swipeout
                              class="nested-list-item"
                              :id="getAccountDomId(account)"
                              :class="{ 'has-child-list-item': account.type === allAccountTypes.MultiSubAccounts.type && hasVisibleSubAccount(account), 'actual-first-child': account.id === firstShowingIds.accounts[accountCategory.type], 'actual-last-child': account.id === lastShowingIds.accounts[accountCategory.type] }"
                              :after="accountBalance(account)"
                              :link="!sortable ? '/transaction/list?accountIds=' + account.id : null"
                              :key="account.id"
                              v-for="account in allCategorizedAccountsMap[accountCategory.type].accounts"
                              v-show="showHidden || !account.hidden"
                              @taphold="setSortable()"
                >
                    <template #media v-if="account.type !== allAccountTypes.MultiSubAccounts.type || !hasVisibleSubAccount(account)">
                        <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color">
                            <f7-badge color="gray" class="right-bottom-icon" v-if="account.hidden">
                                <f7-icon f7="eye_slash_fill"></f7-icon>
                            </f7-badge>
                        </ItemIcon>
                    </template>

                    <template #title>
                        <div class="display-flex padding-top-half padding-bottom-half">
                            <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"
                                      v-if="account.type === allAccountTypes.MultiSubAccounts.type && hasVisibleSubAccount(account)">
                                <f7-badge color="gray" class="right-bottom-icon" v-if="account.hidden">
                                    <f7-icon f7="eye_slash_fill"></f7-icon>
                                </f7-badge>
                            </ItemIcon>
                            <div class="nested-list-item-title">
                                <span>{{ account.name }}</span>
                                <div class="item-footer" v-if="account.comment">{{ account.comment }}</div>
                            </div>
                        </div>
                        <li v-if="account.type === allAccountTypes.MultiSubAccounts.type">
                            <ul class="no-padding">
                                <f7-list-item class="no-sortable nested-list-item-child"
                                              :class="{ 'actual-first-child': subAccount.id === firstShowingIds.subAccounts[account.id], 'actual-last-child': subAccount.id === lastShowingIds.subAccounts[account.id] }"
                                              :id="getAccountDomId(subAccount)"
                                              :title="subAccount.name" :footer="subAccount.comment" :after="accountBalance(subAccount)"
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
                        <f7-swipeout-button color="orange" close :text="$t('Edit')" @click="edit(account)"></f7-swipeout-button>
                        <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(account, false)">
                            <f7-icon f7="trash"></f7-icon>
                        </f7-swipeout-button>
                    </f7-swipeout-actions>
                </f7-list-item>
            </f7-list>
        </div>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ $t('Sort') }}</f7-actions-button>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ $t('Show Hidden Accounts') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ $t('Hide Hidden Accounts') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ $t('Are you sure you want to delete this account?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(accountToDelete, true)">{{ $t('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { AccountType, AccountCategory } from '@/core/account.ts';
import { onSwipeoutDeleted } from '@/lib/ui/mobile.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            loading: true,
            loadingError: null,
            showHidden: false,
            sortable: false,
            accountToDelete: null,
            showMoreActionSheet: false,
            showDeleteActionSheet: false,
            displayOrderModified: false,
            displayOrderSaving: false
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useExchangeRatesStore),
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        allAccountTypes() {
            return AccountType.all();
        },
        allAccountCategories() {
            return AccountCategory.values();
        },
        allCategorizedAccountsMap() {
            return this.accountsStore.allCategorizedAccountsMap;
        },
        allAccountCount() {
            return this.accountsStore.allAvailableAccountsCount;
        },
        firstShowingIds() {
            return this.accountsStore.getFirstShowingIds(this.showHidden);
        },
        lastShowingIds() {
            return this.accountsStore.getLastShowingIds(this.showHidden);
        },
        noAvailableAccount() {
            if (this.showHidden) {
                return this.accountsStore.allAvailableAccountsCount < 1;
            } else {
                return this.accountsStore.allVisibleAccountsCount < 1;
            }
        },
        netAssets() {
            const netAssets = this.accountsStore.getNetAssets(this.showAccountBalance);
            return this.getDisplayCurrency(netAssets, this.defaultCurrency);
        },
        totalAssets() {
            const totalAssets = this.accountsStore.getTotalAssets(this.showAccountBalance);
            return this.getDisplayCurrency(totalAssets, this.defaultCurrency);
        },
        totalLiabilities() {
            const totalLiabilities = this.accountsStore.getTotalLiabilities(this.showAccountBalance);
            return this.getDisplayCurrency(totalLiabilities, this.defaultCurrency);
        },
        showAccountBalance: {
            get: function () {
                return this.settingsStore.appSettings.showAccountBalance;
            },
            set: function (value) {
                this.settingsStore.setShowAccountBalance(value);
            }
        }
    },
    created() {
        const self = this;

        self.loading = true;

        self.accountsStore.loadAllAccounts({
            force: false
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            if (this.accountsStore.accountListStateInvalid && !this.loading) {
                this.reload(null);
            }

            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            if (this.sortable) {
                done();
                return;
            }

            const self = this;
            const force = !!done;

            self.accountsStore.loadAllAccounts({
                force: force
            }).then(() => {
                if (done) {
                    done();
                }

                if (force) {
                    self.$toast('Account list has been updated');
                }
            }).catch(error => {
                if (done) {
                    done();
                }

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        hasAccount(accountCategory, visibleOnly) {
            return this.accountsStore.hasAccount(accountCategory, visibleOnly);
        },
        hasVisibleSubAccount(account) {
            return this.accountsStore.hasVisibleSubAccount(this.showHidden, account);
        },
        accountBalance(account) {
            const balance = this.accountsStore.getAccountBalance(this.showAccountBalance, account);
            return this.getDisplayCurrency(balance, account.currency);
        },
        accountCategoryTotalBalance(accountCategory) {
            const totalBalance = this.accountsStore.getAccountCategoryTotalBalance(this.showAccountBalance, accountCategory);
            return this.getDisplayCurrency(totalBalance, this.defaultCurrency);
        },
        setSortable() {
            if (this.sortable) {
                return;
            }

            this.showHidden = true;
            this.sortable = true;
            this.displayOrderModified = false;
        },
        onSort(event) {
            const self = this;

            if (!event || !event.el || !event.el.id) {
                self.$toast('Unable to move account');
                return;
            }

            const id = self.parseAccountIdFromDomId(event.el.id);

            if (!id) {
                self.$toast('Unable to move account');
                return;
            }

            self.accountsStore.changeAccountDisplayOrder({
                accountId: id,
                from: event.from - 1, // first item in the list is title, so the index need minus one
                to: event.to - 1,
                updateListOrder: true,
                updateGlobalListOrder: true
            }).then(() => {
                self.displayOrderModified = true;
            }).catch(error => {
                self.$toast(error.message || error);
            });
        },
        saveSortResult() {
            const self = this;

            if (!self.displayOrderModified) {
                self.showHidden = false;
                self.sortable = false;
                return;
            }

            self.displayOrderSaving = true;
            self.$showLoading();

            self.accountsStore.updateAccountDisplayOrders().then(() => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                self.showHidden = false;
                self.sortable = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        edit(account) {
            this.f7router.navigate('/account/edit?id=' + account.id);
        },
        hide(account, hidden) {
            const self = this;

            self.$showLoading();

            self.accountsStore.hideAccount({
                account: account,
                hidden: hidden
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        remove(account, confirm) {
            const self = this;

            if (!account) {
                self.$alert('An error occurred');
                return;
            }

            if (!confirm) {
                self.accountToDelete = account;
                self.showDeleteActionSheet = true;
                return;
            }

            self.showDeleteActionSheet = false;
            self.accountToDelete = null;
            self.$showLoading();

            self.accountsStore.deleteAccount({
                account: account,
                beforeResolve: (done) => {
                    onSwipeoutDeleted(self.getAccountDomId(account), done);
                }
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        },
        getAccountDomId(account) {
            return 'account_' + account.id;
        },
        parseAccountIdFromDomId(domId) {
            if (!domId || domId.indexOf('account_') !== 0) {
                return null;
            }

            return domId.substring(8); // account_
        }
    }
};
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
    --f7-list-group-title-height: var(--ebk-account-list-group-title-height);
    --f7-list-item-footer-font-size: var(--ebk-large-footer-font-size);
}

.account-list .item-footer {
    padding-top: 4px;
}
</style>
