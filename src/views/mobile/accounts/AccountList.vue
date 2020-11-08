<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Account List')" :back-link="$t('Back')"></f7-nav-title>
            <f7-nav-right>
                <f7-link href="/account/add" icon-f7="plus"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="skeleton-text" v-if="loading">
            <f7-block-title>Account Category</f7-block-title>
            <f7-list media-list>
                <f7-list-item title="Account Name" after="0.00 USD"></f7-list-item>
            </f7-list>
        </f7-block>

        <f7-block v-for="accountCategory in usedAccountCategories" :key="accountCategory.id">
            <f7-block-title>{{ $t(accountCategory.name) }}</f7-block-title>
            <f7-list media-list>
                <f7-list-item v-for="account in accounts[accountCategory.id]" :key="account.id"
                    :title="account.name" :after="account.balance | currency(account.currency)"></f7-list-item>
            </f7-list>
        </f7-block>
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

                if (this.$utils.isArray(this.accounts[accountCategory.id]) && this.accounts[accountCategory.id].length) {
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

            self.accounts = {};

            for (let i = 0; i < data.result.length; i++) {
                const account = data.result[i];

                if (!self.accounts[account.category]) {
                    self.accounts[account.category] = [];
                }

                const accountList = self.accounts[account.category];
                accountList.push(account);
            }
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

    }
};
</script>
