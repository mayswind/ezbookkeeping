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
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         @click="selectInvert"></v-list-item>
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
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         @click="selectInvert"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </div>
        </template>

        <div v-if="loading">
            <v-skeleton-loader type="paragraph" :loading="loading"
                               :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
        </div>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-if="!loading && !hasAnyAvailableAccount">
            <span class="text-body-1">{{ $t('No available account') }}</span>
        </v-card-text>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-else-if="!loading && hasAnyAvailableAccount">
            <v-expansion-panels class="account-categories" multiple v-model="expandAccountCategories">
                <v-expansion-panel :key="accountCategory.category"
                                   :value="accountCategory.category"
                                   class="border"
                                   v-for="accountCategory in allVisibleCategorizedAccounts">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ml-3">{{ $t(accountCategory.name) }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                        <v-list rounded density="comfortable" class="pa-0">
                            <template :key="account.id"
                                      v-for="(account, idx) in accountCategory.visibleAccounts">
                                <v-list-item>
                                    <template #prepend>
                                        <v-checkbox :model-value="isAccountOrSubAccountsAllChecked(account, filterAccountIds)"
                                                    :indeterminate="isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds)"
                                                    @update:model-value="selectAccountOrSubAccounts(account, $event)">
                                            <template #label>
                                                <ItemIcon class="d-flex" icon-type="account"
                                                          :icon-id="account.icon" :color="account.color"></ItemIcon>
                                                <span class="ml-3">{{ account.name }}</span>
                                            </template>
                                        </v-checkbox>
                                    </template>
                                </v-list-item>

                                <v-divider v-if="account.type === allAccountTypes.MultiSubAccounts && accountCategory.visibleSubAccounts[account.id]"/>

                                <v-list rounded density="comfortable" class="pa-0 ml-4"
                                        v-if="account.type === allAccountTypes.MultiSubAccounts && accountCategory.visibleSubAccounts[account.id]">
                                    <template :key="subAccount.id"
                                              v-for="(subAccount, subIdx) in accountCategory.visibleSubAccounts[account.id]">
                                        <v-list-item>
                                            <template #prepend>
                                                <v-checkbox :model-value="isAccountChecked(subAccount, filterAccountIds)"
                                                            @update:model-value="selectAccount(subAccount, $event)">
                                                    <template #label>
                                                        <ItemIcon class="d-flex" icon-type="account"
                                                                  :icon-id="subAccount.icon" :color="subAccount.color"></ItemIcon>
                                                        <span class="ml-3">{{ subAccount.name }}</span>
                                                    </template>
                                                </v-checkbox>
                                            </template>
                                        </v-list-item>
                                        <v-divider v-if="subIdx !== accountCategory.visibleSubAccounts[account.id].length - 1"/>
                                    </template>
                                </v-list>

                                <v-divider v-if="idx !== accountCategory.visibleAccounts.length - 1"/>
                            </template>
                        </v-list>
                    </v-expansion-panel-text>
                </v-expansion-panel>
            </v-expansion-panels>
        </v-card-text>

        <v-card-text class="overflow-y-visible" v-if="dialogMode">
            <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                <v-btn :disabled="!hasAnyAvailableAccount" @click="save">{{ $t(applyText) }}</v-btn>
                <v-btn color="secondary" variant="tonal" @click="cancel">{{ $t('Cancel') }}</v-btn>
            </div>
        </v-card-text>
    </v-card>
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

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiDotsVertical
} from '@mdi/js';

export default {
    props: [
        'dialogMode',
        'modifyDefault',
        'autoSave'
    ],
    emits: [
        'settings:change'
    ],
    data: function () {
        return {
            loading: true,
            expandAccountCategories: accountConstants.allCategories.map(category => category.id),
            filterAccountIds: {},
            icons: {
                selectAll: mdiSelectAll,
                selectNone: mdiSelect,
                selectInverse: mdiSelectInverse,
                more: mdiDotsVertical
            }
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

            this.$emit('settings:change', true);
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
    padding: 0 0 0 20px;
}

.account-categories .v-expansion-panel--active:not(:first-child),
.account-categories .v-expansion-panel--active + .v-expansion-panel {
    margin-top: 1rem;
}
</style>
