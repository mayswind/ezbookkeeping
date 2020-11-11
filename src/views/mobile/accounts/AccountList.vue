<template>
    <f7-page ptr @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Account List')"></f7-nav-title>
            <f7-nav-right>
                <f7-link href="/account/add" icon-f7="plus"></f7-link>
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

        <f7-card v-for="accountCategory in usedAccountCategories" :key="accountCategory.id">
            <f7-card-header>{{ $t(accountCategory.name) }}</f7-card-header>
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item v-for="account in accounts[accountCategory.id]"
                                  :key="account.id" :id="account | accountDomId"
                                  :title="account.name" :after="account.balance | currency(account.currency)"
                                  link="#" swipeout
                    >
                        <f7-swipeout-actions right>
                            <f7-swipeout-button color="orange" :text="$t('Edit')" @click="edit(account)"></f7-swipeout-button>
                            <f7-swipeout-button color="red" :text="$t('Delete')" @click="remove(account)"></f7-swipeout-button>
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
            loading: true
        };
    },
    computed: {
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
        edit() {

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
                        self.$alert({error: error.response.data});
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
