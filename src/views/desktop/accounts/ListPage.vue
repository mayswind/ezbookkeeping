<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <div class="d-flex flex-column flex-md-row">
                    <div>
                        <div class="mx-6 my-4">
                            <small>{{ $t('Net assets') }}</small>
                            <p class="text-body-1 text-income mt-1 mb-3">
                                <span v-if="!loading">{{ netAssets }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="accounts-summary-skeleton mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <small>{{ $t('Total liabilities') }}</small>
                            <p class="text-body-1 text-expense mt-1 mb-3">
                                <span v-if="!loading">{{ totalLiabilities }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="accounts-summary-skeleton mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <small>{{ $t('Total assets') }}</small>
                            <p class="text-body-1 mt-1">
                                <span v-if="!loading">{{ totalAssets }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="accounts-summary-skeleton mt-3 mb-5" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="account-category-tabs my-4" direction="vertical"
                                :disabled="loading" v-model="activeAccountCategoryId">
                            <v-tab :key="accountCategory.id" :value="accountCategory.id"
                                   v-for="accountCategory in allAccountCategories">
                                <ItemIcon icon-type="account" :icon-id="accountCategory.defaultAccountIconId" />
                                <div class="ml-2 d-flex flex-column">
                                    <small class="text-left smaller">{{ accountCategoryTotalBalance(accountCategory) }}</small>
                                    <div class="text-left">{{ $t(accountCategory.name) }}</div>
                                </div>
                            </v-tab>
                        </v-tabs>
                    </div>
                    <v-window class="d-flex flex-grow-1 ml-md-5 disable-tab-transition w-100-window-container" v-model="activeTab">
                        <v-window-item value="accountPage">
                            <v-card variant="flat">
                                <template #title>
                                    <div class="d-flex align-center">
                                        <span>{{ $t('Account List') }}</span>
                                        <v-btn class="ml-3" color="default" variant="outlined"
                                               :disabled="loading" @click="add">{{ $t('Add') }}</v-btn>
                                        <v-btn class="ml-3" color="primary" variant="tonal"
                                               :disabled="loading" @click="saveSortResult"
                                               v-if="displayOrderModified">{{ $t('Save Display Order') }}</v-btn>
                                        <v-btn density="compact" color="default" variant="text"
                                               class="ml-2" :icon="true" :disabled="loading"
                                               v-if="!loading" @click="reload">
                                            <v-icon :icon="icons.refresh" size="24" />
                                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                        </v-btn>
                                        <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
                                        <v-spacer/>
                                        <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                                               :disabled="loading" :icon="true">
                                            <v-icon :icon="icons.more" />
                                            <v-menu activator="parent">
                                                <v-list>
                                                    <v-list-item :prepend-icon="icons.show"
                                                                 :title="$t('Show Hidden Account')"
                                                                 v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                                    <v-list-item :prepend-icon="icons.hide"
                                                                 :title="$t('Hide Hidden Account')"
                                                                 v-if="showHidden" @click="showHidden = false"></v-list-item>
                                                </v-list>
                                            </v-menu>
                                        </v-btn>
                                    </div>
                                </template>

                                <v-card-text class="accounts-overview-title pt-0">
                                    <span class="text-subtitle-1">{{ $t('Balance') }}</span>
                                    <span class="accounts-overview-amount ml-3">
                                        {{ activeAccountCategoryTotalBalance }}
                                    </span>
                                    <v-btn class="ml-2" density="compact" color="default" variant="text"
                                           :icon="true" :disabled="loading"
                                           @click="showAccountBalance = !showAccountBalance">
                                        <v-icon :icon="showAccountBalance ? icons.eyeSlash : icons.eye" size="20" />
                                        <v-tooltip activator="parent">{{ showAccountBalance ? $t('Hide Account Balance') : $t('Show Account Balance') }}</v-tooltip>
                                    </v-btn>
                                </v-card-text>

                                <div v-if="loading && !hasAccount(activeAccountCategory)">
                                    <v-skeleton-loader type="paragraph" :loading="loading"
                                                       :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
                                </div>

                                <v-row class="pl-4 pr-4 pr-md-8" v-if="!loading && !hasAccount(activeAccountCategory)">
                                    <v-col cols="12">
                                        {{ $t('No available account') }}
                                    </v-col>
                                </v-row>

                                <v-row class="pl-4 pr-4 pr-md-8">
                                    <v-col cols="12">
                                        <draggable-list
                                            class="list-group"
                                            item-key="id"
                                            handle=".drag-handle"
                                            ghost-class="dragging-item"
                                            :disabled="activeAccountCategoryVisibleAccountCount <= 1"
                                            :list="categorizedAccounts[activeAccountCategory.id].accounts"
                                            v-if="categorizedAccounts[activeAccountCategory.id] && categorizedAccounts[activeAccountCategory.id].accounts && categorizedAccounts[activeAccountCategory.id].accounts.length"
                                            @change="onMove"
                                        >
                                            <template #item="{ element }">
                                                <div class="list-group-item">
                                                    <v-card border class="card-title-with-bg account-card mb-8 h-auto" v-if="showHidden || !element.hidden">
                                                        <template #title>
                                                            <div class="account-title d-flex align-baseline">
                                                                <ItemIcon size="1.5rem" icon-type="account" :icon-id="element.icon"
                                                                          :color="element.color" :hidden-status="element.hidden" />
                                                                <span class="account-name ml-2">{{ element.name }}</span>
                                                                <small class="account-currency ml-2">
                                                                    {{ accountCurrency(element) }}
                                                                </small>
                                                                <v-spacer/>
                                                                <span class="align-self-center">
                                                                    <v-icon :class="!loading && activeAccountCategoryVisibleAccountCount > 1 ? 'drag-handle' : 'disabled'"
                                                                            :icon="icons.drag"/>
                                                                    <v-tooltip activator="parent" v-if="!loading && activeAccountCategoryVisibleAccountCount > 1">{{ $t('Drag and Drop to Change Order') }}</v-tooltip>
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
                                                                       :to="`/transaction/list?accountId=${accountOrSubAccountId(element)}`">
                                                                    {{ $t('Transaction List') }}
                                                                </v-btn>
                                                                <v-btn class="hover-display px-2 ml-2" density="comfortable" color="default" variant="text"
                                                                       :disabled="loading" :prepend-icon="icons.edit"
                                                                       @click="edit(element)">
                                                                    {{ $t('Edit') }}
                                                                </v-btn>
                                                                <v-btn class="hover-display px-2 ml-2" density="comfortable" color="default" variant="text"
                                                                       :disabled="loading" :prepend-icon="icons.hide"
                                                                       v-if="!element.hidden"
                                                                       @click="hide(element, true)">
                                                                    {{ $t('Hide') }}
                                                                </v-btn>
                                                                <v-btn class="hover-display px-2 ml-2" density="comfortable" color="default" variant="text"
                                                                       :disabled="loading" :prepend-icon="icons.show"
                                                                       v-if="element.hidden"
                                                                       @click="hide(element, false)">
                                                                    {{ $t('Show') }}
                                                                </v-btn>
                                                                <v-btn class="hover-display px-2 ml-2" density="comfortable" color="default" variant="text"
                                                                       :disabled="loading" :prepend-icon="icons.remove"
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
                </div>
            </v-card>
        </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
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
    mdiPencilOutline,
    mdiDeleteOutline,
    mdiListBoxOutline,
    mdiDrag,
    mdiDotsVertical,
} from '@mdi/js';

export default {
    data() {
        return {
            activeAccountCategoryId: accountConstants.allCategories[0].id,
            activeTab: 'accountPage',
            activeSubAccount: {},
            loading: true,
            displayOrderModified: false,
            showHidden: false,
            icons: {
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline,
                refresh: mdiRefresh,
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
        categorizedAccounts() {
            return this.accountsStore.allCategorizedAccounts;
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
            if (!this.categorizedAccounts[this.activeAccountCategory.id] || !this.categorizedAccounts[this.activeAccountCategory.id].accounts) {
                return 0;
            }

            const accounts = this.categorizedAccounts[this.activeAccountCategory.id].accounts;

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

        },
        edit() {

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
            return this.$locale.getDisplayCurrency(value, currencyCode, {
                currencyDisplayMode: this.settingsStore.appSettings.currencyDisplayMode,
                enableThousandsSeparator: this.settingsStore.appSettings.thousandsSeparator
            });
        }
    }
}
</script>

<style>
.accounts-summary-skeleton .v-skeleton-loader__text {
    margin: 0;
}

.account-category-tabs .v-tab.v-tab {
    --v-btn-height: calc(var(--v-tabs-height) * 1.5);
}

.accounts-overview-title {
    line-height: 2rem !important;
    height: 46px;
    display: flex;
    align-items: flex-end;
}

.accounts-overview-amount {
    font-size: 1.5rem;
    color: rgba(var(--v-theme-on-background), var(--v-high-emphasis-opacity));
    overflow: hidden;
    text-overflow: ellipsis;
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
