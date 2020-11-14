<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Account List')"></f7-nav-title>
            <f7-nav-right>
                <f7-link href="/account/add" icon-f7="plus" v-if="!sortable"></f7-link>
                <f7-link :text="$t('Done')" :class="{ 'disabled': displayOrderSaving }" @click="saveSortResult" v-else-if="sortable"></f7-link>
            </f7-nav-right>
        </f7-navbar>

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

        <f7-card v-if="noAvailableAccount">
            <f7-card-content :padding="false">
                <f7-list sortable sortable-tap-hold :sortable-enabled="sortable" @sortable:sort="onSort">
                    <f7-list-item :title="$t('No available account')" @taphold.native="setSortable()"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-for="accountCategory in usedAccountCategories" :key="accountCategory.id" v-show="showHidden || hasShownAccount(accountCategory)">
            <f7-card-header>{{ $t(accountCategory.name) }}</f7-card-header>
            <f7-card-content :padding="false">
                <f7-list sortable sortable-tap-hold :sortable-enabled="sortable" @sortable:sort="onSort">
                    <f7-list-item v-for="account in accounts[accountCategory.id]" v-show="showHidden || !account.hidden"
                                  :key="account.id" :id="account | accountDomId"
                                  :title="account.name" :after="accountBalance(account) | currency(account.currency)"
                                  link="#" swipeout @taphold.native="setSortable()"
                    >
                        <f7-icon slot="media" :f7="account.icon | accountIcon"></f7-icon>
                        <f7-swipeout-actions left v-if="sortable">
                            <f7-swipeout-button :color="account.hidden ? 'blue' : 'gray'" class="padding-left padding-right" @click="hide(account, !account.hidden)">
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
            displayOrderModified: false,
            displayOrderSaving: false
        };
    },
    computed: {
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
        }
    },
    created() {
        const self = this;
        const router = self.$f7router;

        self.loading = true;

        self.$services.getAllAccounts().then(response => {
            self.loading = false;
            const data = response.data;

            if (!data || !data.success || !data.result) {
                self.$alert('Unable to get account list', () => {
                    router.back();
                });
                return;
            }

            self.accounts = self.$utilities.getCategorizedAccounts(data.result);
        }).catch(error => {
            self.loading = false;

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
        accountBalance(account) {
            if (this.$settings.isShowAccountBalance()) {
                return account.balance;
            } else {
                return '---';
            }
        },
        setSortable() {
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
                self.displayOrderSaving = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$alert({ error: error.response.data });
                } else if (!error.processed) {
                    self.$alert('Unable to move account');
                }
            });
        },
        edit() {

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
