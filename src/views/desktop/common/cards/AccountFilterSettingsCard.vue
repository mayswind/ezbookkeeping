<template>
    <v-card :class="{ 'pa-2 pa-sm-4 pa-md-8': dialogMode }">
        <template #title>
            <div class="d-flex align-center justify-center" v-if="dialogMode">
                <div class="w-100 text-center">
                    <h4 class="text-h4">{{ $t(title) }}</h4>
                </div>
                <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                       :disabled="loading || !hasAnyAvailableAccount" :icon="true">
                    <v-icon :icon="icons.more" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="icons.selectAll"
                                         :title="$t('Select All')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectInvert"></v-list-item>
                            <v-divider class="my-2"/>
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
            <div class="d-flex align-center" v-else-if="!dialogMode">
                <span>{{ $t(title) }}</span>
                <v-spacer/>
                <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                       :disabled="loading" :icon="true">
                    <v-icon :icon="icons.more" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="icons.selectAll"
                                         :title="$t('Select All')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectInvert"></v-list-item>
                            <v-divider class="my-2"/>
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

        <div v-if="loading">
            <v-skeleton-loader type="paragraph" :loading="loading"
                               :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
        </div>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-if="!loading && !hasAnyVisibleAccount">
            <span class="text-body-1">{{ $t('No available account') }}</span>
        </v-card-text>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-else-if="!loading && hasAnyVisibleAccount">
            <v-expansion-panels class="account-categories" multiple v-model="expandAccountCategories">
                <v-expansion-panel :key="accountCategory.category"
                                   :value="accountCategory.category"
                                   class="border"
                                   v-for="accountCategory in allCategorizedAccounts"
                                   v-show="showHidden || accountCategory.allVisibleAccountCount > 0">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ml-3">{{ $t(accountCategory.name) }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                        <v-list rounded density="comfortable" class="pa-0">
                            <template :key="account.id"
                                      v-for="(account, idx) in accountCategory.allAccounts">
                                <v-divider v-if="showHidden ? idx > 0 : (!account.hidden ? idx > accountCategory.firstVisibleAccountIndex : false)"/>

                                <v-list-item v-if="showHidden || !account.hidden">
                                    <template #prepend>
                                        <v-checkbox :model-value="isAccountOrSubAccountsAllChecked(account, filterAccountIds)"
                                                    :indeterminate="isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds)"
                                                    @update:model-value="selectAccountOrSubAccounts(account, $event)">
                                            <template #label>
                                                <ItemIcon class="d-flex" icon-type="account" :icon-id="account.icon"
                                                          :color="account.color" :hidden-status="account.hidden"></ItemIcon>
                                                <span class="ml-3">{{ account.name }}</span>
                                            </template>
                                        </v-checkbox>
                                    </template>
                                </v-list-item>

                                <v-divider v-if="(showHidden || !account.hidden) && account.type === allAccountTypes.MultiSubAccounts.type && ((showHidden && accountCategory.allSubAccounts[account.id]) || accountCategory.allVisibleSubAccountCounts[account.id])"/>

                                <v-list rounded density="comfortable" class="pa-0 ml-4"
                                        v-if="(showHidden || !account.hidden) && account.type === allAccountTypes.MultiSubAccounts.type && ((showHidden && accountCategory.allSubAccounts[account.id]) || accountCategory.allVisibleSubAccountCounts[account.id])">
                                    <template :key="subAccount.id"
                                              v-for="(subAccount, subIdx) in accountCategory.allSubAccounts[account.id]">
                                        <v-divider v-if="showHidden ? subIdx > 0 : (!subAccount.hidden ? subIdx > accountCategory.allFirstVisibleSubAccountIndexes[account.id] : false)"/>

                                        <v-list-item v-if="showHidden || !subAccount.hidden">
                                            <template #prepend>
                                                <v-checkbox :model-value="isAccountChecked(subAccount, filterAccountIds)"
                                                            @update:model-value="selectAccount(subAccount, $event)">
                                                    <template #label>
                                                        <ItemIcon class="d-flex" icon-type="account" :icon-id="subAccount.icon"
                                                                  :color="subAccount.color" :hidden-status="subAccount.hidden"></ItemIcon>
                                                        <span class="ml-3">{{ subAccount.name }}</span>
                                                    </template>
                                                </v-checkbox>
                                            </template>
                                        </v-list-item>
                                    </template>
                                </v-list>
                            </template>
                        </v-list>
                    </v-expansion-panel-text>
                </v-expansion-panel>
            </v-expansion-panels>
        </v-card-text>

        <v-card-text class="overflow-y-visible" v-if="dialogMode">
            <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                <v-btn :disabled="!hasAnyVisibleAccount" @click="save">{{ $t(applyText) }}</v-btn>
                <v-btn color="secondary" variant="tonal" @click="cancel">{{ $t('Cancel') }}</v-btn>
            </div>
        </v-card-text>
    </v-card>

    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { AccountType, AccountCategory } from '@/core/account.ts';
import { copyObjectTo } from '@/lib/common.ts';
import {
    getCategorizedAccountsWithVisibleCount,
    selectAccountOrSubAccounts,
    selectAll,
    selectNone,
    selectInvert,
    isAccountOrSubAccountsAllChecked,
    isAccountOrSubAccountsHasButNotAllChecked
} from '@/lib/account.ts';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiDotsVertical
} from '@mdi/js';

export default {
    props: [
        'dialogMode',
        'type',
        'autoSave'
    ],
    emits: [
        'settings:change'
    ],
    data: function () {
        return {
            loading: true,
            expandAccountCategories: AccountCategory.values().map(category => category.type),
            filterAccountIds: {},
            showHidden: false,
            icons: {
                selectAll: mdiSelectAll,
                selectNone: mdiSelect,
                selectInverse: mdiSelectInverse,
                show: mdiEyeOutline,
                hide: mdiEyeOffOutline,
                more: mdiDotsVertical
            }
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useAccountsStore, useTransactionsStore, useStatisticsStore),
        title() {
            if (this.type === 'statisticsDefault') {
                return 'Default Account Filter';
            } else {
                return 'Filter Accounts';
            }
        },
        applyText() {
            if (this.type === 'statisticsDefault') {
                return 'Save';
            } else {
                return 'Apply';
            }
        },
        allAccountTypes() {
            return AccountType.all();
        },
        allCategorizedAccounts() {
            return getCategorizedAccountsWithVisibleCount(this.accountsStore.allCategorizedAccountsMap);
        },
        hasAnyAvailableAccount() {
            return this.accountsStore.allAvailableAccountsCount > 0;
        },
        hasAnyVisibleAccount() {
            if (this.showHidden) {
                return this.accountsStore.allAvailableAccountsCount > 0;
            } else {
                return this.accountsStore.allVisibleAccountsCount > 0;
            }
        }
    },
    created() {
        const self = this;

        self.accountsStore.loadAllAccounts({
            force: false
        }).then(() => {
            self.loading = false;

            const allAccountIds = {};

            for (const accountId in self.accountsStore.allAccountsMap) {
                if (!Object.prototype.hasOwnProperty.call(self.accountsStore.allAccountsMap, accountId)) {
                    continue;
                }

                const account = self.accountsStore.allAccountsMap[accountId];

                if (self.type === 'transactionListCurrent' && self.transactionsStore.allFilterAccountIdsCount > 0) {
                    allAccountIds[account.id] = true;
                } else {
                    allAccountIds[account.id] = false;
                }
            }

            if (self.type === 'statisticsDefault') {
                self.filterAccountIds = copyObjectTo(self.settingsStore.appSettings.statistics.defaultAccountFilter, allAccountIds);
            } else if (self.type === 'statisticsCurrent') {
                self.filterAccountIds = copyObjectTo(self.statisticsStore.transactionStatisticsFilter.filterAccountIds, allAccountIds);
            } else if (self.type === 'transactionListCurrent') {
                for (const accountId in self.transactionsStore.allFilterAccountIds) {
                    if (!Object.prototype.hasOwnProperty.call(self.transactionsStore.allFilterAccountIds, accountId)) {
                        continue;
                    }

                    const account = self.accountsStore.allAccountsMap[accountId];

                    if (account) {
                        selectAccountOrSubAccounts(allAccountIds, account, false);
                    }
                }
                self.filterAccountIds = allAccountIds;
            } else {
                self.$refs.snackbar.showError('Parameter Invalid');
            }
        }).catch(error => {
            self.loading = false;

            if (!error.processed) {
                self.$refs.snackbar.showError(error);
            }
        });
    },
    methods: {
        save() {
            const self = this;

            const filteredAccountIds = {};
            let isAllSelected = true;
            let finalAccountIds = '';
            let changed = true;

            for (const accountId in self.filterAccountIds) {
                if (!Object.prototype.hasOwnProperty.call(self.filterAccountIds, accountId)) {
                    continue;
                }

                const account = self.accountsStore.allAccountsMap[accountId];

                if (!isAccountOrSubAccountsAllChecked(account, self.filterAccountIds)) {
                    filteredAccountIds[accountId] = true;
                    isAllSelected = false;
                } else {
                    if (finalAccountIds.length > 0) {
                        finalAccountIds += ',';
                    }

                    finalAccountIds += accountId;
                }
            }

            if (this.type === 'statisticsDefault') {
                self.settingsStore.setStatisticsDefaultAccountFilter(filteredAccountIds);
            } else if (this.type === 'statisticsCurrent') {
                changed = self.statisticsStore.updateTransactionStatisticsFilter({
                    filterAccountIds: filteredAccountIds
                });

                if (changed) {
                    self.statisticsStore.updateTransactionStatisticsInvalidState(true);
                }
            } else if (this.type === 'transactionListCurrent') {
                changed = self.transactionsStore.updateTransactionListFilter({
                    accountIds: isAllSelected ? '' : finalAccountIds
                });

                if (changed) {
                    self.transactionsStore.updateTransactionListInvalidState(true);
                }
            }

            self.$emit('settings:change', changed);
        },
        cancel() {
            this.$emit('settings:change', false);
        },
        selectAccountOrSubAccounts(account, value) {
            selectAccountOrSubAccounts(this.filterAccountIds, account, !value);

            if (this.autoSave) {
                this.save();
            }
        },
        selectAccount(account, value) {
            this.filterAccountIds[account.id] = !value;

            if (this.autoSave) {
                this.save();
            }
        },
        selectAll() {
            selectAll(this.filterAccountIds, this.accountsStore.allAccountsMap);

            if (this.autoSave) {
                this.save();
            }
        },
        selectNone() {
            selectNone(this.filterAccountIds, this.accountsStore.allAccountsMap);

            if (this.autoSave) {
                this.save();
            }
        },
        selectInvert() {
            selectInvert(this.filterAccountIds, this.accountsStore.allAccountsMap);

            if (this.autoSave) {
                this.save();
            }
        },
        isAccountChecked(account, filterAccountIds) {
            return !filterAccountIds[account.id];
        },
        isAccountOrSubAccountsAllChecked(account, filterAccountIds) {
            return isAccountOrSubAccountsAllChecked(account, filterAccountIds);
        },
        isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds) {
            return isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds);
        }
    }
}
</script>

<style>
.account-categories .v-expansion-panel-text__wrapper {
    padding: 0 0 0 0;
}

.account-categories .v-expansion-panel--active:not(:first-child),
.account-categories .v-expansion-panel--active + .v-expansion-panel {
    margin-top: 1rem;
}
</style>
