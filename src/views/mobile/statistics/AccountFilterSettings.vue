<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t('Save')" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>
                <small :style="{ opacity: 0.6 }">
                    <span>Account Category</span>
                </small>
            </f7-card-header>

            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item checkbox class="disabled" title="Account Name">
                        <f7-icon slot="media" f7="app_fill"></f7-icon>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>
                <small :style="{ opacity: 0.6 }">
                    <span>Account Category 2</span>
                </small>
            </f7-card-header>

            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item checkbox class="disabled" title="Account Name">
                        <f7-icon slot="media" f7="app_fill"></f7-icon>
                    </f7-list-item>
                    <f7-list-item checkbox class="disabled" title="Account Name 2">
                        <f7-icon slot="media" f7="app_fill"></f7-icon>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-block class="no-padding no-margin" v-if="!loading">
            <f7-card v-for="accountCategory in allAccountCategories" :key="accountCategory.id" v-show="hasShownAccount(accountCategory)">
                <f7-card-header>
                    <small :style="{ opacity: 0.6 }">
                        <span>{{ $t(accountCategory.name) }}</span>
                    </small>
                </f7-card-header>
                <f7-card-content class="no-safe-areas" :padding="false">
                    <f7-list v-if="categorizedAccounts[accountCategory.id]">
                        <f7-list-item checkbox v-for="account in categorizedAccounts[accountCategory.id].accounts"
                                      v-show="!account.hidden"
                                      :key="account.id"
                                      :title="account.name"
                                      :value="account.id"
                                      :checked="account | accountOrSubAccountsAllChecked(filterAccountIds)"
                                      :indeterminate="account | accountOrSubAccountsHasButNotAllChecked(filterAccountIds)"
                                      @change="selectAccountOrSubAccounts">
                            <f7-icon slot="media"
                                     :icon="account.icon | accountIcon"
                                     :style="account.color | accountIconStyle('var(--default-icon-color)')">
                            </f7-icon>

                            <ul slot="root" v-if="account.type === $constants.account.allAccountTypes.MultiSubAccounts" class="padding-left">
                                <f7-list-item checkbox v-for="subAccount in account.subAccounts"
                                              v-show="!subAccount.hidden"
                                              :key="subAccount.id"
                                              :title="subAccount.name"
                                              :value="subAccount.id"
                                              :checked="subAccount | accountChecked(filterAccountIds) "
                                              @change="selectAccount">
                                    <f7-icon slot="media"
                                             :icon="subAccount.icon | accountIcon"
                                             :style="subAccount.color | accountIconStyle('var(--default-icon-color)')">
                                    </f7-icon>
                                </f7-list-item>
                            </ul>
                        </f7-list-item>
                    </f7-list>
                </f7-card-content>
            </f7-card>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="selectAll">{{ $t('Select All') }}</f7-actions-button>
                <f7-actions-button @click="selectNone">{{ $t('Select None') }}</f7-actions-button>
                <f7-actions-button @click="selectInvert">{{ $t('Invert Selection') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
export default {
    data: function () {
        return {
            loading: true,
            modifyDefault: false,
            filterAccountIds: {},
            showMoreActionSheet: false
        }
    },
    computed: {
        title() {
            if (this.modifyDefault) {
                return 'Default Account Filter';
            } else {
                return 'Filter Accounts';
            }
        },
        allAccountCategories() {
            return this.$constants.account.allCategories;
        },
        categorizedAccounts() {
            return this.$store.state.allCategorizedAccounts;
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;
        const router = self.$f7router;

        self.modifyDefault = !!query.modifyDefault;

        self.$store.dispatch('loadAllAccounts', {
            force: false
        }).then(() => {
            self.loading = false;

            const allAccountIds = {};

            for (let accountId in self.$store.state.allAccountsMap) {
                if (!Object.prototype.hasOwnProperty.call(self.$store.state.allAccountsMap, accountId)) {
                    continue;
                }

                const account = self.$store.state.allAccountsMap[accountId];
                allAccountIds[account.id] = false;
            }

            if (self.modifyDefault) {
                self.filterAccountIds = self.$utilities.copyObjectTo(self.$settings.getStatisticsDefaultAccountFilter(), allAccountIds);
            } else {
                self.filterAccountIds = self.$utilities.copyObjectTo(self.$store.state.transactionStatisticsFilter.filterAccountIds, allAccountIds);
            }
        }).catch(error => {
            self.logining = false;

            if (!error.processed) {
                self.$toast(error.message || error);
                router.back();
            }
        });
    },
    methods: {
        save() {
            const self = this;
            const router = self.$f7router;

            const filteredAccountIds = {};

            for (let accountId in self.filterAccountIds) {
                if (!Object.prototype.hasOwnProperty.call(self.filterAccountIds, accountId)) {
                    continue;
                }

                if (self.filterAccountIds[accountId]) {
                    filteredAccountIds[accountId] = true;
                }
            }

            if (self.modifyDefault) {
                self.$settings.setStatisticsDefaultAccountFilter(filteredAccountIds);
            } else {
                self.$store.dispatch('updateTransactionStatisticsFilter', {
                    filterAccountIds: filteredAccountIds
                });
            }

            router.back();
        },
        selectAccountOrSubAccounts(e) {
            const accountId = e.target.value;
            const account = this.$store.state.allAccountsMap[accountId];

            if (!account) {
                return;
            }

            if (account.type === this.$constants.account.allAccountTypes.SingleAccount) {
                this.filterAccountIds[account.id] = !e.target.checked;
            } else if (account.type === this.$constants.account.allAccountTypes.MultiSubAccounts) {
                if (!account.subAccounts || !account.subAccounts.length) {
                    return;
                }

                for (let i = 0; i < account.subAccounts.length; i++) {
                    const subAccount = account.subAccounts[i];
                    this.filterAccountIds[subAccount.id] = !e.target.checked;
                }
            }
        },
        selectAccount(e) {
            const accountId = e.target.value;
            const account = this.$store.state.allAccountsMap[accountId];

            if (!account) {
                return;
            }

            this.filterAccountIds[account.id] = !e.target.checked;
        },
        selectAll() {
            for (let accountId in this.filterAccountIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterAccountIds, accountId)) {
                    continue;
                }

                const account = this.$store.state.allAccountsMap[accountId];

                if (account && account.type === this.$constants.account.allAccountTypes.SingleAccount) {
                    this.filterAccountIds[account.id] = false;
                }
            }
        },
        selectNone() {
            for (let accountId in this.filterAccountIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterAccountIds, accountId)) {
                    continue;
                }

                const account = this.$store.state.allAccountsMap[accountId];

                if (account && account.type === this.$constants.account.allAccountTypes.SingleAccount) {
                    this.filterAccountIds[account.id] = true;
                }
            }
        },
        selectInvert() {
            for (let accountId in this.filterAccountIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterAccountIds, accountId)) {
                    continue;
                }

                const account = this.$store.state.allAccountsMap[accountId];

                if (account && account.type === this.$constants.account.allAccountTypes.SingleAccount) {
                    this.filterAccountIds[account.id] = !this.filterAccountIds[account.id];
                }
            }
        },
        hasShownAccount(accountCategory) {
            if (!this.categorizedAccounts[accountCategory.id] ||
                !this.categorizedAccounts[accountCategory.id].accounts ||
                !this.categorizedAccounts[accountCategory.id].accounts.length) {
                return false;
            }

            for (let i = 0; i < this.categorizedAccounts[accountCategory.id].accounts.length; i++) {
                const account = this.categorizedAccounts[accountCategory.id].accounts[i];

                if (!account.hidden) {
                    return true;
                }
            }

            return false;
        }
    },
    filters: {
        accountChecked(account, filterAccountIds) {
            return !filterAccountIds[account.id];
        },
        accountOrSubAccountsAllChecked(account, filterAccountIds) {
            if (!account.subAccounts) {
                return !filterAccountIds[account.id];
            }

            for (let i = 0; i < account.subAccounts.length; i++) {
                const subAccount = account.subAccounts[i];
                if (filterAccountIds[subAccount.id]) {
                    return false;
                }
            }

            return true;
        },
        accountOrSubAccountsHasButNotAllChecked(account, filterAccountIds) {
            if (!account.subAccounts) {
                return false;
            }

            let checkedCount = 0;

            for (let i = 0; i < account.subAccounts.length; i++) {
                const subAccount = account.subAccounts[i];
                if (!filterAccountIds[subAccount.id]) {
                    checkedCount++;
                }
            }

            return checkedCount > 0 && checkedCount < account.subAccounts.length;
        }
    }
}
</script>
