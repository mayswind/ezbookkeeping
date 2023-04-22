<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableAccount }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t(applyText)" :class="{ 'disabled': !hasAnyAvailableAccount }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="combination-list-wrapper margin-vertical skeleton-text"
                  :key="blockIdx" v-for="blockIdx in [ 1, 2, 3 ]" v-if="loading">
            <f7-accordion-item>
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="combination-list-header combination-list-opened">
                            <f7-list-item>
                                <template #title>
                                    <span>Account Category</span>
                                    <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content style="height: auto">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item checkbox class="disabled" title="Account Name"
                                      :key="itemIdx" v-for="itemIdx in (blockIdx === 1 ? [ 1 ] : [ 1, 2 ])">
                            <template #media>
                                <f7-icon f7="app_fill"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-list strong inset dividers accordion-list class="margin-top" v-if="!hasAnyAvailableAccount">
            <f7-list-item :title="$t('No available account')"></f7-list-item>
        </f7-list>

        <f7-block class="combination-list-wrapper margin-vertical"
                  :key="accountCategory.id"
                  v-for="accountCategory in allAccountCategories"
                  v-else-if="!loading && hasAnyAvailableAccount">
            <f7-accordion-item :opened="collapseStates[accountCategory.id].opened"
                               v-show="hasShownAccount(accountCategory)"
                               @accordion:open="collapseStates[accountCategory.id].opened = true"
                               @accordion:close="collapseStates[accountCategory.id].opened = false">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="combination-list-header"
                                 :class="collapseStates[accountCategory.id].opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item>
                                <template #title>
                                    <span>{{ $t(accountCategory.name) }}</span>
                                    <f7-icon class="combination-list-chevron-icon" :f7="collapseStates[accountCategory.id].opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates[accountCategory.id].opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content"
                             v-if="categorizedAccounts[accountCategory.id]">
                        <f7-list-item checkbox
                                      :title="account.name"
                                      :value="account.id"
                                      :checked="isAccountOrSubAccountsAllChecked(account, filterAccountIds)"
                                      :indeterminate="isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds)"
                                      :key="account.id"
                                      v-for="account in categorizedAccounts[accountCategory.id].accounts"
                                      v-show="!account.hidden"
                                      @change="selectAccountOrSubAccounts">
                            <template #media>
                                <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"></ItemIcon>
                            </template>

                            <template #root>
                                <ul v-if="account.type === $constants.account.allAccountTypes.MultiSubAccounts" class="padding-left">
                                    <f7-list-item checkbox
                                                  :title="subAccount.name"
                                                  :value="subAccount.id"
                                                  :checked="isAccountChecked(subAccount, filterAccountIds)"
                                                  :key="subAccount.id"
                                                  v-for="subAccount in account.subAccounts"
                                                  v-show="!subAccount.hidden"
                                                  @change="selectAccount">
                                        <template #media>
                                            <ItemIcon icon-type="account" :icon-id="subAccount.icon" :color="subAccount.color"></ItemIcon>
                                        </template>
                                    </f7-list-item>
                                </ul>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
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
    props: [
        'f7route',
        'f7router'
    ],
    data: function () {
        const self = this;

        return {
            loading: true,
            loadingError: null,
            modifyDefault: false,
            filterAccountIds: {},
            collapseStates: self.getCollapseStates(),
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
        applyText() {
            if (this.modifyDefault) {
                return 'Save';
            } else {
                return 'Apply';
            }
        },
        allAccountCategories() {
            return this.$constants.account.allCategories;
        },
        categorizedAccounts() {
            return this.$store.state.allCategorizedAccounts;
        },
        hasAnyAvailableAccount() {
            return this.$store.getters.allVisibleAccountsCount > 0;
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

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
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        save() {
            const self = this;
            const router = self.f7router;

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
        isAccountChecked(account, filterAccountIds) {
            return !filterAccountIds[account.id];
        },
        isAccountOrSubAccountsAllChecked(account, filterAccountIds) {
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
        isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds) {
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
        },
        getCollapseStates() {
            const collapseStates = {};

            for (let categoryType in this.$constants.account.allCategories) {
                if (!Object.prototype.hasOwnProperty.call(this.$constants.account.allCategories, categoryType)) {
                    continue;
                }

                const accountCategory = this.$constants.account.allCategories[categoryType];

                collapseStates[accountCategory.id] = {
                    opened: true
                };
            }

            return collapseStates;
        }
    }
}
</script>
