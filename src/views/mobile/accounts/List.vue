<template>
    <f7-page :ptr="!sortable" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Account List')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" v-if="!sortable && allAccountCount" @click="showMoreActionSheet = true"></f7-link>
                <f7-link href="/account/add" icon-f7="plus" v-if="!sortable"></f7-link>
                <f7-link :text="$t('Done')" :class="{ 'disabled': displayOrderSaving }" @click="saveSortResult" v-else-if="sortable"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="account-overview-card" :class="{ 'skeleton-text': loading }">
            <f7-card-header class="display-block" style="padding-top: 120px;">
                <p class="no-margin">
                    <small :style="{ opacity: 0.6 }" v-if="loading">Net assets</small>
                    <small :style="{ opacity: 0.6 }" v-else-if="!loading">{{ $t('Net assets') }}</small>
                </p>
                <p class="no-margin">
                    <span class="net-assets" v-if="loading">0.00 USD</span>
                    <span class="net-assets" v-else-if="!loading">{{ netAssets | currency(defaultCurrency) }}</span>
                    <f7-link class="margin-left-half" @click="toggleShowAccountBalance()">
                        <f7-icon :f7="showAccountBalance ? 'eye_slash_fill' : 'eye_fill'" size="18px"></f7-icon>
                    </f7-link>
                </p>
                <p class="no-margin">
                    <small class="account-overview-info" :style="{ opacity: 0.6 }" v-if="loading">
                        <span>Total assets | Total liabilities</span>
                    </small>
                    <small class="account-overview-info" :style="{ opacity: 0.6 }" v-else-if="!loading">
                        <span>{{ $t('Total assets') }}</span>
                        <span>{{ totalAssets | currency(defaultCurrency) }}</span>
                        <span>|</span>
                        <span>{{ $t('Total liabilities') }}</span>
                        <span>{{ totalLiabilities | currency(defaultCurrency) }}</span>
                    </small>
                </p>
            </f7-card-header>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>
                <small style="opacity: 0.6;">Account Category</small>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item class="nested-list-item" after="0.00 USD" link="#">
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex padding-top-half padding-bottom-half">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                                <div class="nested-list-item-title">Account Name</div>
                            </div>
                        </f7-block>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>
                <small style="opacity: 0.6;">Account Category 2</small>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item class="nested-list-item" after="0.00 USD" link="#">
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex padding-top-half padding-bottom-half">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                                <div class="nested-list-item-title">Account Name</div>
                            </div>
                        </f7-block>
                    </f7-list-item>
                    <f7-list-item class="nested-list-item" after="0.00 USD" link="#">
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex padding-top-half padding-bottom-half">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                                <div class="nested-list-item-title">Account Name 2</div>
                            </div>
                        </f7-block>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>
                <small style="opacity: 0.6;">Account Category 3</small>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item class="nested-list-item" after="0.00 USD" link="#">
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex padding-top-half padding-bottom-half">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                                <div class="nested-list-item-title">Account Name</div>
                            </div>
                        </f7-block>
                    </f7-list-item>
                    <f7-list-item class="nested-list-item" after="0.00 USD" link="#">
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex padding-top-half padding-bottom-half">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                                <div class="nested-list-item-title">Account Name 2</div>
                            </div>
                        </f7-block>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-if="!loading && noAvailableAccount">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('No available account')"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-for="accountCategory in allAccountCategories" :key="accountCategory.id" v-show="(showHidden && hasAccount(accountCategory, false)) || hasAccount(accountCategory, true)">
            <f7-card-header>
                <small :style="{ opacity: 0.6 }">
                    <span>{{ $t(accountCategory.name) }}</span>
                    <span style="margin-left: 10px">{{ accountCategoryTotalBalance(accountCategory) | currency(defaultCurrency) }}</span>
                </small>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list sortable :sortable-enabled="sortable" @sortable:sort="onSort" v-if="categorizedAccounts[accountCategory.id]">
                    <f7-list-item v-for="account in categorizedAccounts[accountCategory.id].accounts" v-show="showHidden || !account.hidden"
                                  :key="account.id" :id="account | accountDomId"
                                  :class="{ 'nested-list-item': true, 'has-child-list-item': account.type === $constants.account.allAccountTypes.MultiSubAccounts }"
                                  :after="accountBalance(account) | currency(account.currency)"
                                  :link="!sortable && account.type === $constants.account.allAccountTypes.SingleAccount ? '/transaction/list?accountId=' + account.id : null"
                                  swipeout @taphold.native="setSortable()"
                    >
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex padding-top-half padding-bottom-half">
                                <f7-icon slot="media" :icon="account.icon | accountIcon"
                                         :style="account.color | accountIconStyle('var(--default-icon-color)')">
                                    <f7-badge color="gray" class="right-bottom-icon" v-if="account.hidden">
                                        <f7-icon f7="eye_slash_fill"></f7-icon>
                                    </f7-badge>
                                </f7-icon>
                                <div class="nested-list-item-title">{{ account.name }}</div>
                            </div>
                            <li v-if="account.type === $constants.account.allAccountTypes.MultiSubAccounts">
                                <ul class="no-padding">
                                    <f7-list-item class="no-sortable nested-list-item-child" v-for="subAccount in account.subAccounts" v-show="showHidden || !subAccount.hidden"
                                                  :key="subAccount.id" :id="subAccount | accountDomId"
                                                  :title="subAccount.name" :after="accountBalance(subAccount) | currency(subAccount.currency)"
                                                  :link="!sortable ? '/transaction/list?accountId=' + subAccount.id : null"
                                    >
                                        <f7-icon slot="media" :icon="subAccount.icon | accountIcon"
                                                 :style="subAccount.color | accountIconStyle('var(--default-icon-color)')">
                                            <f7-badge color="gray" class="right-bottom-icon" v-if="subAccount.hidden">
                                                <f7-icon f7="eye_slash_fill"></f7-icon>
                                            </f7-badge>
                                        </f7-icon>
                                    </f7-list-item>
                                </ul>
                            </li>
                        </f7-block>
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
            </f7-card-content>
        </f7-card>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ $t('Sort') }}</f7-actions-button>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ $t('Show Hidden Account') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ $t('Hide Hidden Account') }}</f7-actions-button>
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
export default {
    data() {
        return {
            loading: true,
            showHidden: false,
            sortable: false,
            accountToDelete: null,
            showAccountBalance: this.$settings.isShowAccountBalance(),
            showMoreActionSheet: false,
            showDeleteActionSheet: false,
            displayOrderModified: false,
            displayOrderSaving: false
        };
    },
    computed: {
        defaultCurrency() {
            return this.$store.getters.currentUserDefaultCurrency || this.$t('default.currency');
        },
        allAccountCategories() {
            return this.$constants.account.allCategories;
        },
        categorizedAccounts() {
            return this.$store.state.allCategorizedAccounts;
        },
        allAccountCount() {
            return this.$store.getters.allAvailableAccountsCount;
        },
        noAvailableAccount() {
            if (this.showHidden) {
                return this.$store.getters.allAvailableAccountsCount < 1;
            } else {
                return this.$store.getters.allVisibleAccountsCount < 1;
            }
        },
        netAssets() {
            if (!this.showAccountBalance) {
                return '***';
            }

            const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(this.categorizedAccounts, () => true);
            let netAssets = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === this.defaultCurrency) {
                    netAssets += accountsBalance[i].balance;
                } else {
                    const balance = this.$store.getters.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

                    if (!this.$utilities.isNumber(balance)) {
                        hasUnCalculatedAmount = true;
                        continue;
                    }

                    netAssets += Math.floor(balance);
                }
            }

            if (hasUnCalculatedAmount) {
                return netAssets + '+';
            } else {
                return netAssets;
            }
        },
        totalAssets() {
            if (!this.showAccountBalance) {
                return '***';
            }

            const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(this.categorizedAccounts, account => account.isAsset);
            let totalAssets = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === this.defaultCurrency) {
                    totalAssets += accountsBalance[i].balance;
                } else {
                    const balance = this.$store.getters.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

                    if (!this.$utilities.isNumber(balance)) {
                        hasUnCalculatedAmount = true;
                        continue;
                    }

                    totalAssets += Math.floor(balance);
                }
            }

            if (hasUnCalculatedAmount) {
                return totalAssets + '+';
            } else {
                return totalAssets;
            }
        },
        totalLiabilities() {
            if (!this.showAccountBalance) {
                return '***';
            }

            const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(this.categorizedAccounts, account => account.isLiability);
            let totalLiabilities = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === this.defaultCurrency) {
                    totalLiabilities -= accountsBalance[i].balance;
                } else {
                    const balance = this.$store.getters.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

                    if (!this.$utilities.isNumber(balance)) {
                        hasUnCalculatedAmount = true;
                        continue;
                    }

                    totalLiabilities -= Math.floor(balance);
                }
            }

            if (hasUnCalculatedAmount) {
                return totalLiabilities + '+';
            } else {
                return totalLiabilities;
            }
        }
    },
    created() {
        const self = this;
        const router = self.$f7router;

        self.loading = true;

        self.$store.dispatch('loadAllAccounts', {
            force: false
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            self.loading = false;

            if (!error.processed) {
                self.$toast(error.message || error);
                router.back();
            }
        });
    },
    methods: {
        onPageAfterIn() {
            if (this.$store.state.accountListStateInvalid && !this.loading) {
                this.reload(null);
            }
        },
        reload(done) {
            if (this.sortable) {
                done();
                return;
            }

            const self = this;

            self.$store.dispatch('loadAllAccounts', {
                force: true
            }).then(() => {
                if (done) {
                    done();
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
            if (!this.categorizedAccounts[accountCategory.id] ||
                !this.categorizedAccounts[accountCategory.id].accounts ||
                !this.categorizedAccounts[accountCategory.id].accounts.length) {
                return false;
            }

            let shownCount = 0;

            for (let i = 0; i < this.categorizedAccounts[accountCategory.id].accounts.length; i++) {
                const account = this.categorizedAccounts[accountCategory.id].accounts[i];

                if (!visibleOnly || !account.hidden) {
                    shownCount++;
                }
            }

            return shownCount > 0;
        },
        toggleShowAccountBalance() {
            this.showAccountBalance = !this.showAccountBalance;
            this.$settings.setShowAccountBalance(this.showAccountBalance);
        },
        accountBalance(account) {
            if (account.type !== this.$constants.account.allAccountTypes.SingleAccount) {
                return null;
            }

            if (this.showAccountBalance) {
                if (account.isAsset) {
                    return account.balance;
                } else if (account.isLiability) {
                    return -account.balance;
                } else {
                    return account.balance;
                }
            } else {
                return '***';
            }
        },
        accountCategoryTotalBalance(accountCategory) {
            if (!this.showAccountBalance) {
                return '***';
            }

            const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(this.categorizedAccounts, account => account.category === accountCategory.id);
            let totalBalance = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === this.defaultCurrency) {
                    if (accountsBalance[i].isAsset) {
                        totalBalance += accountsBalance[i].balance;
                    } else if (accountsBalance[i].isLiability) {
                        totalBalance -= accountsBalance[i].balance;
                    } else {
                        totalBalance += accountsBalance[i].balance;
                    }
                } else {
                    const balance = this.$store.getters.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

                    if (!this.$utilities.isNumber(balance)) {
                        hasUnCalculatedAmount = true;
                        continue;
                    }

                    if (accountsBalance[i].isAsset) {
                        totalBalance += Math.floor(balance);
                    } else if (accountsBalance[i].isLiability) {
                        totalBalance -= Math.floor(balance);
                    } else {
                        totalBalance += Math.floor(balance);
                    }
                }
            }

            if (hasUnCalculatedAmount) {
                return totalBalance + '+';
            } else {
                return totalBalance;
            }
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

            if (!event || !event.el || !event.el.id || event.el.id.indexOf('account_') !== 0) {
                self.$toast('Unable to move account');
                return;
            }

            const id = event.el.id.substr(8); // account_

            self.$store.dispatch('changeAccountDisplayOrder', {
                accountId: id,
                from: event.from,
                to: event.to
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

            self.$store.dispatch('updateAccountDisplayOrders').then(() => {
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
            this.$f7router.navigate('/account/edit?id=' + account.id);
        },
        hide(account, hidden) {
            const self = this;

            self.$showLoading();

            self.$store.dispatch('hideAccount', {
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
            const app = self.$f7;
            const $$ = app.$;

            if (!account) {
                self.$alert('An error has occurred');
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

            self.$store.dispatch('deleteAccount', {
                account: account,
                beforeResolve: (done) => {
                    app.swipeout.delete($$(`#${self.$options.filters.accountDomId(account)}`), () => {
                        done();
                    });
                }
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        }
    },
    filters: {
        accountDomId(account) {
            return 'account_' + account.id;
        }
    }
};
</script>

<style>
.account-overview-card {
    background-color: var(--f7-color-yellow);
}

.theme-dark .account-overview-card {
    background-color: var(--f7-theme-color);
}

.theme-dark .account-overview-card a {
    color: var(--f7-text-color);
    opacity: 0.6;
}

.net-assets {
    font-size: 1.5em;
}

.account-overview-info > span {
    margin-right: 4px;
}

.account-overview-info > span:last-child {
    margin-right: 0;
}
</style>
