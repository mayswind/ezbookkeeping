<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <span class="text-subtitle-2">{{ $t('Net assets') }}</span>
                            <p class="account-statistic-item-value text-income text-truncate mt-1 mb-3">
                                <span v-if="!loading || allAccountCount > 0">{{ netAssets }}</span>
                                <span v-else-if="loading && allAccountCount <= 0">
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-1" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <span class="text-subtitle-2">{{ $t('Total liabilities') }}</span>
                            <p class="account-statistic-item-value text-expense text-truncate mt-1 mb-3">
                                <span v-if="!loading || allAccountCount > 0">{{ totalLiabilities }}</span>
                                <span v-else-if="loading && allAccountCount <= 0">
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-1" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <span class="text-subtitle-2">{{ $t('Total assets') }}</span>
                            <p class="account-statistic-item-value mt-1">
                                <span v-if="!loading || allAccountCount > 0">{{ totalAssets }}</span>
                                <span v-else-if="loading && allAccountCount <= 0">
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-1" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="account-category-tabs my-4" direction="vertical"
                                :disabled="loading" v-model="activeAccountCategoryId">
                            <v-tab class="tab-text-truncate" :key="accountCategory.id" :value="accountCategory.id"
                                   v-for="accountCategory in allAccountCategories">
                                <ItemIcon icon-type="account" :icon-id="accountCategory.defaultAccountIconId" />
                                <div class="d-flex flex-column text-truncate ml-2">
                                    <small class="text-truncate text-left smaller" v-if="!loading || allAccountCount > 0">{{ accountCategoryTotalBalance(accountCategory) }}</small>
                                    <small class="text-truncate text-left smaller my-1" v-else-if="loading && allAccountCount <= 0">
                                        <v-skeleton-loader class="skeleton-no-margin"
                                                           width="100px" height="16" type="text" :loading="true"></v-skeleton-loader>
                                    </small>
                                    <span class="text-truncate text-left">{{ $t(accountCategory.name) }}</span>
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
                                                <v-icon :icon="icons.menu" size="24" />
                                            </v-btn>
                                            <span>{{ $t('Account List') }}</span>
                                            <v-btn class="ml-3" color="default" variant="outlined"
                                                   :disabled="loading" @click="add">{{ $t('Add') }}</v-btn>
                                            <v-btn class="ml-3" color="primary" variant="tonal"
                                                   :disabled="loading" @click="saveSortResult"
                                                   v-if="displayOrderModified">{{ $t('Save Display Order') }}</v-btn>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ml-2" :icon="true" :loading="loading" @click="reload">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="icons.refresh" size="24" />
                                                <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-spacer/>
                                            <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                                                   :disabled="loading" :icon="true">
                                                <v-icon :icon="icons.more" />
                                                <v-menu activator="parent">
                                                    <v-list>
                                                        <v-list-item :prepend-icon="icons.show"
                                                                     :title="$t('Show Hidden Accounts')"
                                                                     v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                                        <v-list-item :prepend-icon="icons.hide"
                                                                     :title="$t('Hide Hidden Accounts')"
                                                                     v-if="showHidden" @click="showHidden = false"></v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <v-card-text class="accounts-overview-title text-truncate pt-0">
                                        <span class="accounts-overview-subtitle">{{ $t('Balance') }}</span>
                                        <v-skeleton-loader class="skeleton-no-margin ml-3 mb-2" width="120px" type="text" :loading="true" v-if="loading && !hasAccount(activeAccountCategory)"></v-skeleton-loader>
                                        <span class="accounts-overview-amount ml-3" v-else-if="!loading || hasAccount(activeAccountCategory)">{{ activeAccountCategoryTotalBalance }}</span>
                                        <v-btn class="ml-2" density="compact" color="default" variant="text"
                                               :icon="true" :disabled="loading"
                                               @click="showAccountBalance = !showAccountBalance">
                                            <v-icon :icon="showAccountBalance ? icons.eyeSlash : icons.eye" size="20" />
                                            <v-tooltip activator="parent">{{ showAccountBalance ? $t('Hide Account Balance') : $t('Show Account Balance') }}</v-tooltip>
                                        </v-btn>
                                    </v-card-text>

                                    <v-row class="pl-6 pr-6 pr-md-8" v-if="loading && !hasAccount(activeAccountCategory)">
                                        <v-col cols="12">
                                            <v-card border class="card-title-with-bg account-card mb-8 h-auto">
                                                <template #title>
                                                    <div class="account-title d-flex align-center">
                                                        <v-icon class="disabled mr-0" size="28px" :icon="icons.square" />
                                                        <span class="account-name text-truncate ml-2">
                                                            <v-skeleton-loader class="skeleton-no-margin my-1"
                                                                               width="120px" type="text" :loading="true"></v-skeleton-loader>
                                                        </span>
                                                        <v-spacer/>
                                                        <span class="align-self-center">
                                                            <v-icon class="disabled" :icon="icons.drag"/>
                                                        </span>
                                                    </div>
                                                </template>
                                                <v-divider/>
                                                <v-card-text>
                                                    <div class="d-flex account-toolbar align-center">
                                                        <v-btn class="px-2" density="comfortable" color="default" variant="text"
                                                               :disabled="true" :prepend-icon="icons.transactions">
                                                            {{ $t('Transaction List') }}
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

                                    <v-row class="pl-5 pr-2 pr-md-4" v-if="!loading && !hasAccount(activeAccountCategory)">
                                        <v-col cols="12">
                                            {{ $t('No available account') }}
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
                                                :list="allCategorizedAccountsMap[activeAccountCategory.id].accounts"
                                                v-if="allCategorizedAccountsMap[activeAccountCategory.id] && allCategorizedAccountsMap[activeAccountCategory.id].accounts && allCategorizedAccountsMap[activeAccountCategory.id].accounts.length"
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
                                                                                :icon="icons.drag"/>
                                                                        <v-tooltip activator="parent" v-if="!loading && activeAccountCategoryVisibleAccountCount > 1">{{ $t('Drag to Reorder') }}</v-tooltip>
                                                                    </span>
                                                                </div>

                                                                <div class="mt-4" v-if="element.type === allAccountTypes.MultiSubAccounts">
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
                                                                        <v-btn :value="undefined">
                                                                            <span>{{ $t('All') }}</span>
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

                                                            <v-card-text v-if="accountComment(element)">
                                                                {{ accountComment(element) }}
                                                            </v-card-text>

                                                            <v-card-text>
                                                                <div class="d-flex account-toolbar align-center">
                                                                    <v-btn class="px-2" density="comfortable" color="default" variant="text"
                                                                           :disabled="loading" :prepend-icon="icons.transactions"
                                                                           :to="`/transaction/list?accountIds=${accountOrSubAccountId(element)}`">
                                                                        {{ $t('Transaction List') }}
                                                                    </v-btn>
                                                                    <v-btn class="px-2 ml-1" density="comfortable" color="default" variant="text"
                                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                           :disabled="loading"
                                                                           :prepend-icon="element.hidden ? icons.show : icons.hide"
                                                                           v-if="!activeSubAccount[element.id]"
                                                                           @click="hide(element, !element.hidden)">
                                                                        {{ element.hidden ? $t('Show') : $t('Hide') }}
                                                                    </v-btn>
                                                                    <v-btn class="px-2 ml-1" density="comfortable" color="default" variant="text"
                                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                           :disabled="loading" :prepend-icon="icons.edit"
                                                                           v-if="!activeSubAccount[element.id]"
                                                                           @click="edit(element)">
                                                                        {{ $t('Edit') }}
                                                                    </v-btn>
                                                                    <v-btn class="px-2 ml-1" density="comfortable" color="default" variant="text"
                                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                           :disabled="loading" :prepend-icon="icons.remove"
                                                                           v-if="!activeSubAccount[element.id]"
                                                                           @click="remove(element)">
                                                                        {{ $t('Delete') }}
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

    <edit-dialog ref="editDialog" :persistent="true" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
import EditDialog from './list/dialogs/EditDialog.vue';

import { useDisplay } from 'vuetify';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import accountConstants from '@/consts/account.js';
import { isObject } from '@/lib/common.js';
import {
    getAccountCategoryInfo,
    getSubAccountCurrencies,
    getAccountOrSubAccountId,
    getAccountOrSubAccountComment
} from '@/lib/account.js';

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

export default {
    components: {
        EditDialog
    },
    data() {
        const { mdAndUp } = useDisplay();

        return {
            activeAccountCategoryId: accountConstants.allCategories[0].id,
            activeTab: 'accountPage',
            activeSubAccount: {},
            loading: true,
            displayOrderModified: false,
            alwaysShowNav: mdAndUp.value,
            showNav: mdAndUp.value,
            showHidden: false,
            icons: {
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline,
                refresh: mdiRefresh,
                square: mdiSquareRounded,
                menu: mdiMenu,
                edit: mdiPencilOutline,
                show: mdiEyeOutline,
                hide: mdiEyeOffOutline,
                remove: mdiDeleteOutline,
                transactions: mdiListBoxOutline,
                drag: mdiDrag,
                more: mdiDotsVertical
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useExchangeRatesStore),
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        allAccountTypes() {
            return accountConstants.allAccountTypes;
        },
        allAccountCategories() {
            return accountConstants.allCategories;
        },
        allCategorizedAccountsMap() {
            return this.accountsStore.allCategorizedAccountsMap;
        },
        allAccountCount() {
            return this.accountsStore.allAvailableAccountsCount;
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
        activeAccountCategory() {
            return getAccountCategoryInfo(this.activeAccountCategoryId);
        },
        activeAccountCategoryTotalBalance() {
            return this.accountCategoryTotalBalance(this.activeAccountCategory);
        },
        activeAccountCategoryVisibleAccountCount() {
            if (!this.allCategorizedAccountsMap[this.activeAccountCategory.id] || !this.allCategorizedAccountsMap[this.activeAccountCategory.id].accounts) {
                return 0;
            }

            const accounts = this.allCategorizedAccountsMap[this.activeAccountCategory.id].accounts;

            if (this.showHidden) {
                return accounts.length;
            }

            let visibleCount = 0;

            for (let i = 0; i < accounts.length; i++) {
                if (!accounts[i].hidden) {
                    visibleCount++;
                }
            }

            return visibleCount;
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
        this.reload(false);
    },
    setup() {
        const display = useDisplay();

        return {
            display: display
        };
    },
    watch: {
        'display.mdAndUp.value': function (newValue) {
            this.alwaysShowNav = newValue;

            if (!this.showNav) {
                this.showNav = newValue;
            }
        }
    },
    methods: {
        reload(force) {
            const self = this;

            self.loading = true;

            self.accountsStore.loadAllAccounts({
                force: force
            }).then(() => {
                self.loading = false;
                self.displayOrderModified = false;

                if (force) {
                    self.$refs.snackbar.showMessage('Account list has been updated');
                }
            }).catch(error => {
                self.loading = false;

                if (error && error.message === 'Account list is up to date') {
                    self.displayOrderModified = false;
                }

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        hasAccount(accountCategory) {
            if (this.showHidden) {
                return this.accountsStore.hasAccount(accountCategory, false);
            } else {
                return this.accountsStore.hasAccount(accountCategory, true);
            }
        },
        accountOrSubAccountId(account) {
            return getAccountOrSubAccountId(account, this.activeSubAccount[account.id]);
        },
        accountComment(account) {
            return getAccountOrSubAccountComment(account, this.activeSubAccount[account.id]);
        },
        accountCurrency(account) {
            const self = this;

            if (account.type === self.allAccountTypes.SingleAccount) {
                return self.$locale.getCurrencyName(account.currency);
            } else if (account.type === self.allAccountTypes.MultiSubAccounts) {
                const subAccountCurrencies = getSubAccountCurrencies(account, self.showHidden, self.activeSubAccount[account.id])
                    .map(currencyCode => self.$locale.getCurrencyName(currencyCode));
                return self.$locale.joinMultiText(subAccountCurrencies);
            } else {
                return null;
            }
        },
        accountBalance(account) {
            if (account.type === this.allAccountTypes.SingleAccount) {
                const balance = this.accountsStore.getAccountBalance(this.showAccountBalance, account);
                return this.getDisplayCurrency(balance, account.currency);
            } else if (account.type === this.allAccountTypes.MultiSubAccounts) {
                const balanceResult = this.accountsStore.getAccountSubAccountBalance(this.showAccountBalance, this.showHidden, account, this.activeSubAccount[account.id]);

                if (!isObject(balanceResult)) {
                    return this.getDisplayCurrency(balanceResult, this.defaultCurrency);
                }

                return this.getDisplayCurrency(balanceResult.balance, balanceResult.currency);
            } else {
                return null;
            }
        },
        accountCategoryTotalBalance(accountCategory) {
            const totalBalance = this.accountsStore.getAccountCategoryTotalBalance(this.showAccountBalance, accountCategory);
            return this.getDisplayCurrency(totalBalance, this.defaultCurrency);
        },
        onMove(event) {
            if (!event || !event.moved) {
                return;
            }

            const self = this;
            const moveEvent = event.moved;

            if (!moveEvent.element || !moveEvent.element.id) {
                self.$refs.snackbar.showMessage('Unable to move account');
                return;
            }

            self.accountsStore.changeAccountDisplayOrder({
                accountId: moveEvent.element.id,
                from: moveEvent.oldIndex,
                to: moveEvent.newIndex,
                updateListOrder: false,
                updateGlobalListOrder: true
            }).then(() => {
                self.displayOrderModified = true;
            }).catch(error => {
                self.$refs.snackbar.showError(error);
            });
        },
        saveSortResult() {
            const self = this;

            if (!self.displayOrderModified) {
                return;
            }

            self.loading = true;

            self.accountsStore.updateAccountDisplayOrders().then(() => {
                self.loading = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.loading = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        add() {
            const self = this;

            self.$refs.editDialog.open({
                category: self.activeAccountCategoryId
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        edit(account) {
            const self = this;

            self.$refs.editDialog.open({
                id: account.id,
                currentAccount: account
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }

                if (self.accountsStore.accountListStateInvalid && !self.loading) {
                    self.reload(false);
                }
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        hide(account, hidden) {
            const self = this;

            self.loading = true;

            self.accountsStore.hideAccount({
                account: account,
                hidden: hidden
            }).then(() => {
                self.loading = false;
            }).catch(error => {
                self.loading = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        remove(account) {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to delete this account?').then(() => {
                self.loading = true;

                self.accountsStore.deleteAccount({
                    account: account
                }).then(() => {
                    self.loading = false;
                }).catch(error => {
                    self.loading = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        }
    }
}
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
