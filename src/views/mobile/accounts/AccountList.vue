<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Account List')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" v-if="!sortable" @click="showMoreActionSheet = true"></f7-link>
                <f7-link href="/account/add" icon-f7="plus" v-if="!sortable"></f7-link>
                <f7-link :text="$t('Done')" :class="{ 'disabled': displayOrderSaving }" @click="saveSortResult" v-else-if="sortable"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card :class="{ 'bg-color-yellow': true, 'skeleton-text': loading }">
            <f7-card-header class="display-block" style="padding-top: 100px;">
                <small :style="{ opacity: 0.6 }">{{ loading ? 'Net assets' : $t('Net assets') }}</small><br />
                <span class="net-assets">{{ netAssets | currency(defaultCurrency) }}</span>
                <f7-link class="margin-left-half" @click="toggleShowAccountBalance()">
                    <f7-icon :f7="showAccountBalance ? 'eye_slash_fill' : 'eye_fill'" size="18px"></f7-icon>
                </f7-link>
                <br />
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
            </f7-card-header>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>Account Category</f7-card-header>
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item title="Account Name" after="0.00 USD"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>Account Category 2</f7-card-header>
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item title="Account Name" after="0.00 USD"></f7-list-item>
                    <f7-list-item title="Account Name 2" after="0.00 USD"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>Account Category 3</f7-card-header>
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item title="Account Name" after="0.00 USD"></f7-list-item>
                    <f7-list-item title="Account Name 2" after="0.00 USD"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-if="!loading && noAvailableAccount">
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('No available account')"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-for="accountCategory in usedAccountCategories" :key="accountCategory.id" v-show="showHidden || hasShownAccount(accountCategory)">
            <f7-card-header>
                <small :style="{ opacity: 0.6 }">
                    <span>{{ $t(accountCategory.name) }}</span>
                    <span style="margin-left: 10px">{{ accountCategoryTotalBalance(accountCategory) | currency(defaultCurrency) }}</span>
                </small>
            </f7-card-header>
            <f7-card-content :padding="false">
                <f7-list sortable :sortable-enabled="sortable" @sortable:sort="onSort">
                    <f7-list-item v-for="account in accounts[accountCategory.id]" v-show="showHidden || !account.hidden"
                                  :key="account.id" :id="account | accountDomId"
                                  :class="{ 'nested-list-item': true, 'has-child-list-item': account.type === $constants.account.allAccountTypes.MultiSubAccounts }"
                                  :after="accountBalance(account) | currency(account.currency)"
                                  :link="account.type === $constants.account.allAccountTypes.SingleAccount ? '#' : null"
                                  swipeout @taphold.native="setSortable()"
                    >
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex">
                                <f7-icon slot="media" :f7="account.icon | accountIcon" :style="{ color: '#' + account.color }">
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
                                                  link="#"
                                    >
                                        <f7-icon slot="media" :f7="subAccount.icon | accountIcon" :style="{ color: '#' + subAccount.color }">
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
                                                overswipe @click="hide(account, !account.hidden)">
                                <f7-icon :f7="account.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                            </f7-swipeout-button>
                        </f7-swipeout-actions>
                        <f7-swipeout-actions right v-if="!sortable">
                            <f7-swipeout-button color="orange" :text="$t('Edit')" @click="edit(account)"></f7-swipeout-button>
                            <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(account)">
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
            accounts: {},
            loading: true,
            showHidden: false,
            sortable: false,
            showAccountBalance: this.$settings.isShowAccountBalance(),
            showMoreActionSheet: false,
            displayOrderModified: false,
            displayOrderSaving: false
        };
    },
    computed: {
        defaultCurrency() {
            return this.$user.getUserInfo() ? this.$user.getUserInfo().defaultCurrency : this.$t('default.currency');
        },
        noAvailableAccount() {
            let allAccountCount = 0;
            let shownAccountCount = 0;

            for (let category in this.accounts) {
                if (!Object.prototype.hasOwnProperty.call(this.accounts, category)) {
                    continue;
                }

                const accountList = this.accounts[category];

                for (let i = 0; i < accountList.length; i++) {
                    if (!accountList[i].hidden) {
                        shownAccountCount++;
                    }

                    allAccountCount++;
                }
            }

            if (this.showHidden) {
                return allAccountCount < 1;
            } else {
                return shownAccountCount < 1;
            }
        },
        usedAccountCategories() {
            const allAccountCategories = this.$constants.account.allCategories;
            const usedAccountCategories = [];

            for (let i = 0; i < allAccountCategories.length; i++) {
                const accountCategory = allAccountCategories[i];

                if (this.$utilities.isArray(this.accounts[accountCategory.id]) && this.accounts[accountCategory.id].length) {
                    usedAccountCategories.push(accountCategory);
                }
            }

            return usedAccountCategories;
        },
        netAssets() {
            if (!this.showAccountBalance) {
                return '***';
            }

            const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(this.accounts, () => true);
            let netAssets = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === this.defaultCurrency) {
                    netAssets += accountsBalance[i].balance;
                } else {
                    const balance = this.$exchangeRates.getOtherCurrencyAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

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

            const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(this.accounts, category => category.isAsset);
            let totalAssets = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === this.defaultCurrency) {
                    totalAssets += accountsBalance[i].balance;
                } else {
                    const balance = this.$exchangeRates.getOtherCurrencyAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

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

            const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(this.accounts, category => category.isLiability);
            let totalLiabilities = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === this.defaultCurrency) {
                    totalLiabilities += accountsBalance[i].balance;
                } else {
                    const balance = this.$exchangeRates.getOtherCurrencyAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

                    if (!this.$utilities.isNumber(balance)) {
                        hasUnCalculatedAmount = true;
                        continue;
                    }

                    totalLiabilities += Math.floor(balance);
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

        self.$services.getAllAccounts().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                self.$alert('Unable to get account list', () => {
                    router.back();
                });
                return;
            }

            self.accounts = self.$utilities.getCategorizedAccounts(data.result);
            self.loading = false;
        }).catch(error => {
            self.$logger.error('failed to load account list', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                self.$alert({ error: error.response.data }, () => {
                    router.back();
                });
            } else if (!error.processed) {
                self.$alert('Unable to get account list', () => {
                    router.back();
                });
            }
        });
    },
    methods: {
        reload(done) {
            const self = this;

            self.$services.getAllAccounts().then(response => {
                done();

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to get account list');
                    return;
                }

                self.accounts = self.$utilities.getCategorizedAccounts(data.result);
            }).catch(error => {
                self.$logger.error('failed to reload account list', error);

                done();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to get account list');
                }
            });
        },
        hasShownAccount(accountCategory) {
            if (!this.accounts[accountCategory.id].length) {
                return false;
            }

            let shownCount = 0;

            for (let i = 0; i < this.accounts[accountCategory.id].length; i++) {
                const account = this.accounts[accountCategory.id][i];

                if (!account.hidden) {
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
                return account.balance;
            } else {
                return '***';
            }
        },
        accountCategoryTotalBalance(accountCategory) {
            if (!this.showAccountBalance) {
                return '***';
            }

            const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(this.accounts, category => category.id === accountCategory.id);
            let totalBalance = 0;
            let hasUnCalculatedAmount = false;

            for (let i = 0; i < accountsBalance.length; i++) {
                if (accountsBalance[i].currency === this.defaultCurrency) {
                    totalBalance += accountsBalance[i].balance;
                } else {
                    const balance = this.$exchangeRates.getOtherCurrencyAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

                    if (!this.$utilities.isNumber(balance)) {
                        hasUnCalculatedAmount = true;
                        continue;
                    }

                    totalBalance += Math.floor(balance);
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
            if (!event || !event.el || !event.el.id || event.el.id.indexOf('account_') !== 0) {
                this.$toast('Unable to move account');
                return;
            }

            const id = event.el.id.substr(8);
            const account = this.$utilities.getAccountByAccountId(this.accounts, id);

            if (!account || !this.accounts[account.category] || !this.accounts[account.category][event.to]) {
                this.$toast('Unable to move account');
                return;
            }

            const accountList = this.accounts[account.category];
            accountList.splice(event.to, 0, accountList.splice(event.from, 1)[0]);

            this.displayOrderModified = true;
        },
        saveSortResult() {
            const self = this;
            const newDisplayOrders = [];

            if (!self.displayOrderModified) {
                self.showHidden = false;
                self.sortable = false;
                return;
            }

            self.displayOrderSaving = true;

            for (let category in self.accounts) {
                if (!Object.prototype.hasOwnProperty.call(self.accounts, category)) {
                    continue;
                }

                const accountList = self.accounts[category];

                for (let i = 0; i < accountList.length; i++) {
                    newDisplayOrders.push({
                        id: accountList[i].id,
                        displayOrder: i + 1
                    });
                }
            }

            self.$showLoading();

            self.$services.moveAccount({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$alert('Unable to move account');
                    return;
                }

                self.showHidden = false;
                self.sortable = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.$logger.error('failed to save accounts display order', error);

                self.displayOrderSaving = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert({ error: error.response.data });
                } else if (!error.processed) {
                    self.$alert('Unable to move account');
                }
            });
        },
        edit(account) {
            this.$f7router.navigate('/account/edit?id=' + account.id);
        },
        hide(account, hidden) {
            const self = this;

            self.$showLoading();

            self.$services.hideAccount({
                id: account.id,
                hidden: hidden
            }).then(response => {
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (hidden) {
                        self.$toast('Unable to hide this account');
                    } else {
                        self.$toast('Unable to unhide this account');
                    }

                    return;
                }

                account.hidden = hidden;
            }).catch(error => {
                self.$logger.error('failed to change account visibility', error);

                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    if (hidden) {
                        self.$toast('Unable to hide this account');
                    } else {
                        self.$toast('Unable to unhide this account');
                    }
                }
            });
        },
        remove(account) {
            const self = this;
            const app = self.$f7;
            const $$ = app.$;

            self.$confirm('Are you sure you want to delete this account?', () => {
                self.$showLoading();

                self.$services.deleteAccount({
                    id: account.id
                }).then(response => {
                    self.$hideLoading();
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        self.$alert('Unable to delete this account');
                        return;
                    }

                    app.swipeout.delete($$(`#${self.$options.filters.accountDomId(account)}`), () => {
                        const accountList = self.accounts[account.category];
                        for (let i = 0; i < accountList.length; i++) {
                            if (accountList[i].id === account.id) {
                                accountList.splice(i, 1);
                            }
                        }
                    });
                }).catch(error => {
                    self.$logger.error('failed to delete account', error);

                    self.$hideLoading();

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        self.$alert({ error: error.response.data });
                    } else if (!error.processed) {
                        self.$alert('Unable to delete this account');
                    }
                });
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
