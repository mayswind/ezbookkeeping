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

        <f7-list strong inset dividers accordion-list class="margin-top" v-if="!loading && !hasAnyAvailableAccount">
            <f7-list-item :title="$t('No available account')"></f7-list-item>
        </f7-list>

        <f7-block class="combination-list-wrapper margin-vertical"
                  :key="accountCategory.category"
                  v-for="accountCategory in allVisibleCategorizedAccounts"
                  v-else-if="!loading && hasAnyAvailableAccount">
            <f7-accordion-item :opened="collapseStates[accountCategory.category].opened"
                               @accordion:open="collapseStates[accountCategory.category].opened = true"
                               @accordion:close="collapseStates[accountCategory.category].opened = false">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="combination-list-header"
                                 :class="collapseStates[accountCategory.category].opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item>
                                <template #title>
                                    <span>{{ $t(accountCategory.name) }}</span>
                                    <f7-icon class="combination-list-chevron-icon" :f7="collapseStates[accountCategory.category].opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates[accountCategory.category].opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item checkbox
                                      :class="{ 'has-child-list-item': account.type === allAccountTypes.MultiSubAccounts && accountCategory.visibleSubAccounts[account.id] }"
                                      :title="account.name"
                                      :value="account.id"
                                      :checked="isAccountOrSubAccountsAllChecked(account, filterAccountIds)"
                                      :indeterminate="isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds)"
                                      :key="account.id"
                                      v-for="account in accountCategory.visibleAccounts"
                                      @change="selectAccountOrSubAccounts">
                            <template #media>
                                <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"></ItemIcon>
                            </template>

                            <template #root>
                                <ul class="padding-left"
                                    v-if="account.type === allAccountTypes.MultiSubAccounts && accountCategory.visibleSubAccounts[account.id]">
                                    <f7-list-item checkbox
                                                  :title="subAccount.name"
                                                  :value="subAccount.id"
                                                  :checked="isAccountChecked(subAccount, filterAccountIds)"
                                                  :key="subAccount.id"
                                                  v-for="subAccount in accountCategory.visibleSubAccounts[account.id]"
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
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useAccountsStore } from '@/stores/account.js';
import { useStatisticsStore } from '@/stores/statistics.js';

import accountConstants from '@/consts/account.js';
import { copyObjectTo } from '@/lib/common.js';
import {
    getVisibleCategorizedAccounts,
    selectAccountOrSubAccounts,
    selectAll,
    selectNone,
    selectInvert,
    isAccountOrSubAccountsAllChecked,
    isAccountOrSubAccountsHasButNotAllChecked
} from '@/lib/account.js';

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
        ...mapStores(useSettingsStore, useAccountsStore, useStatisticsStore),
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
        allAccountTypes() {
            return accountConstants.allAccountTypes;
        },
        allVisibleCategorizedAccounts() {
            return getVisibleCategorizedAccounts(this.accountsStore.allCategorizedAccounts);
        },
        hasAnyAvailableAccount() {
            return this.accountsStore.allVisibleAccountsCount > 0;
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        self.modifyDefault = !!query.modifyDefault;

        self.accountsStore.loadAllAccounts({
            force: false
        }).then(() => {
            self.loading = false;

            const allAccountIds = {};

            for (let accountId in self.accountsStore.allAccountsMap) {
                if (!Object.prototype.hasOwnProperty.call(self.accountsStore.allAccountsMap, accountId)) {
                    continue;
                }

                const account = self.accountsStore.allAccountsMap[accountId];
                allAccountIds[account.id] = false;
            }

            if (self.modifyDefault) {
                self.filterAccountIds = copyObjectTo(self.settingsStore.appSettings.statistics.defaultAccountFilter, allAccountIds);
            } else {
                self.filterAccountIds = copyObjectTo(self.statisticsStore.transactionStatisticsFilter.filterAccountIds, allAccountIds);
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
                self.settingsStore.setStatisticsDefaultAccountFilter(filteredAccountIds);
            } else {
                self.statisticsStore.updateTransactionStatisticsFilter({
                    filterAccountIds: filteredAccountIds
                });
            }

            router.back();
        },
        selectAccountOrSubAccounts(e) {
            const accountId = e.target.value;
            const account = this.accountsStore.allAccountsMap[accountId];

            if (!account) {
                return;
            }

            selectAccountOrSubAccounts(this.filterAccountIds, account, !e.target.checked);
        },
        selectAccount(e) {
            const accountId = e.target.value;
            const account = this.accountsStore.allAccountsMap[accountId];

            if (!account) {
                return;
            }

            this.filterAccountIds[account.id] = !e.target.checked;
        },
        selectAll() {
            selectAll(this.filterAccountIds, this.accountsStore.allAccountsMap);
        },
        selectNone() {
            selectNone(this.filterAccountIds, this.accountsStore.allAccountsMap);
        },
        selectInvert() {
            selectInvert(this.filterAccountIds, this.accountsStore.allAccountsMap);
        },
        isAccountChecked(account, filterAccountIds) {
            return !filterAccountIds[account.id];
        },
        isAccountOrSubAccountsAllChecked(account, filterAccountIds) {
            return isAccountOrSubAccountsAllChecked(account, filterAccountIds);
        },
        isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds) {
            return isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds);
        },
        getCollapseStates() {
            const collapseStates = {};

            for (let categoryType in accountConstants.allCategories) {
                if (!Object.prototype.hasOwnProperty.call(accountConstants.allCategories, categoryType)) {
                    continue;
                }

                const accountCategory = accountConstants.allCategories[categoryType];

                collapseStates[accountCategory.id] = {
                    opened: true
                };
            }

            return collapseStates;
        }
    }
}
</script>
