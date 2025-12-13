<template>
    <v-card :class="{ 'pa-sm-1 pa-md-2': dialogMode }">
        <template #title>
            <div class="d-flex align-center justify-center" v-if="dialogMode">
                <div class="w-100 text-center">
                    <h4 class="text-h4">{{ tt(title) }}</h4>
                </div>
                <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                       :disabled="loading || !hasAnyAvailableAccount" :icon="true">
                    <v-icon :icon="mdiDotsVertical" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectAllAccounts"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelect"
                                         :title="tt('Select None')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectNoneAccounts"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelectInverse"
                                         :title="tt('Invert Selection')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectInvertAccounts"></v-list-item>
                            <v-divider class="my-2" v-if="allowHiddenAccount"/>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All Visible')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectAllVisibleAccounts"
                                         v-if="allowHiddenAccount"></v-list-item>
                            <v-divider class="my-2" v-if="allowHiddenAccount"/>
                            <v-list-item :prepend-icon="mdiEyeOutline"
                                         :title="tt('Show Hidden Accounts')"
                                         v-if="allowHiddenAccount && !showHidden" @click="showHidden = true"></v-list-item>
                            <v-list-item :prepend-icon="mdiEyeOffOutline"
                                         :title="tt('Hide Hidden Accounts')"
                                         v-if="allowHiddenAccount && showHidden" @click="showHidden = false"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </div>
            <div class="d-flex align-center" v-else-if="!dialogMode">
                <span>{{ tt(title) }}</span>
                <v-spacer/>
                <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                       :disabled="loading" :icon="true">
                    <v-icon :icon="mdiDotsVertical" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectAllAccounts"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelect"
                                         :title="tt('Select None')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectNoneAccounts"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelectInverse"
                                         :title="tt('Invert Selection')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectInvertAccounts"></v-list-item>
                            <v-divider class="my-2" v-if="allowHiddenAccount"/>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All Visible')"
                                         :disabled="!hasAnyVisibleAccount"
                                         @click="selectAllVisibleAccounts"
                                         v-if="allowHiddenAccount"></v-list-item>
                            <v-divider class="my-2" v-if="allowHiddenAccount"/>
                            <v-list-item :prepend-icon="mdiEyeOutline"
                                         :title="tt('Show Hidden Accounts')"
                                         v-if="allowHiddenAccount && !showHidden" @click="showHidden = true"></v-list-item>
                            <v-list-item :prepend-icon="mdiEyeOffOutline"
                                         :title="tt('Hide Hidden Accounts')"
                                         v-if="allowHiddenAccount && showHidden" @click="showHidden = false"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </div>
        </template>

        <div v-if="loading">
            <v-skeleton-loader type="paragraph" :loading="loading"
                               :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
        </div>

        <v-card-text v-if="!loading && !hasAnyVisibleAccount">
            <span class="text-body-1">{{ tt('No available account') }}</span>
        </v-card-text>

        <v-card-text :class="{ 'flex-grow-1 overflow-y-auto': dialogMode }" v-else-if="!loading && hasAnyVisibleAccount">
            <v-expansion-panels class="account-categories" multiple v-model="expandAccountCategories">
                <v-expansion-panel :key="accountCategory.category"
                                   :value="accountCategory.category"
                                   class="border"
                                   v-for="accountCategory in allCategorizedAccounts"
                                   v-show="showHidden || accountCategory.allVisibleAccountCount > 0">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ms-3">{{ tt(accountCategory.name) }}</span>
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
                                                    @update:model-value="updateAccountOrSubAccountsSelected(account, $event)">
                                            <template #label>
                                                <ItemIcon class="d-flex" icon-type="account" :icon-id="account.icon"
                                                          :color="account.color" :hidden-status="account.hidden"></ItemIcon>
                                                <span class="ms-3">{{ account.name }}</span>
                                            </template>
                                        </v-checkbox>
                                    </template>
                                </v-list-item>

                                <v-divider v-if="(showHidden || !account.hidden) && account.type === AccountType.MultiSubAccounts.type && ((showHidden && accountCategory.allSubAccounts[account.id]) || accountCategory.allVisibleSubAccountCounts[account.id])"/>

                                <v-list rounded density="comfortable" class="pa-0 ms-4"
                                        v-if="(showHidden || !account.hidden) && account.type === AccountType.MultiSubAccounts.type && ((showHidden && accountCategory.allSubAccounts[account.id]) || accountCategory.allVisibleSubAccountCounts[account.id])">
                                    <template :key="subAccount.id"
                                              v-for="(subAccount, subIdx) in accountCategory.allSubAccounts[account.id]">
                                        <v-divider v-if="showHidden ? subIdx > 0 : (!subAccount.hidden ? subIdx > (accountCategory.allFirstVisibleSubAccountIndexes[account.id] as number) : false)"/>

                                        <v-list-item v-if="showHidden || !subAccount.hidden">
                                            <template #prepend>
                                                <v-checkbox :model-value="isAccountChecked(subAccount, filterAccountIds)"
                                                            @update:model-value="updateAccountSelected(subAccount, $event)">
                                                    <template #label>
                                                        <ItemIcon class="d-flex" icon-type="account" :icon-id="subAccount.icon"
                                                                  :color="subAccount.color" :hidden-status="subAccount.hidden"></ItemIcon>
                                                        <span class="ms-3">{{ subAccount.name }}</span>
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
            <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                <v-btn :disabled="!hasAnyVisibleAccount" @click="save">{{ tt(applyText) }}</v-btn>
                <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
            </div>
        </v-card-text>
    </v-card>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type AccountFilterType,
    useAccountFilterSettingPageBase
} from '@/views/base/settings/AccountFilterSettingPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';

import { AccountType, AccountCategory } from '@/core/account.ts';
import type { Account } from '@/models/account.ts';

import {
    selectAccountOrSubAccounts,
    selectAllVisible,
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

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    type: AccountFilterType;
    dialogMode?: boolean;
    autoSave?: boolean;
}>();

const emit = defineEmits<{
    (e: 'settings:change', changed: boolean): void;
}>();

const { tt } = useI18n();

const {
    loading,
    showHidden,
    filterAccountIds,
    title,
    applyText,
    allowHiddenAccount,
    allCategorizedAccounts,
    hasAnyAvailableAccount,
    hasAnyVisibleAccount,
    isAccountChecked,
    loadFilterAccountIds,
    saveFilterAccountIds
} = useAccountFilterSettingPageBase(props.type);

const accountsStore = useAccountsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const expandAccountCategories = ref<number[]>(AccountCategory.values().map(category => category.type));

function init(): void {
    accountsStore.loadAllAccounts({
        force: false
    }).then(() => {
        loading.value = false;

        if (!loadFilterAccountIds()) {
            snackbar.value?.showError('Parameter Invalid');
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function updateAccountOrSubAccountsSelected(account: Account, value: boolean | null): void {
    selectAccountOrSubAccounts(filterAccountIds.value, account, !value);

    if (props.autoSave) {
        save();
    }
}

function updateAccountSelected(account: Account, value: boolean | null): void {
    filterAccountIds.value[account.id] = !value;

    if (props.autoSave) {
        save();
    }
}

function selectAllAccounts(): void {
    selectAll(filterAccountIds.value, accountsStore.allAccountsMap, !allowHiddenAccount.value);

    if (props.autoSave) {
        save();
    }
}

function selectNoneAccounts(): void {
    selectNone(filterAccountIds.value, accountsStore.allAccountsMap, !allowHiddenAccount.value);

    if (props.autoSave) {
        save();
    }
}

function selectInvertAccounts(): void {
    selectInvert(filterAccountIds.value, accountsStore.allAccountsMap, !allowHiddenAccount.value);

    if (props.autoSave) {
        save();
    }
}

function selectAllVisibleAccounts(): void {
    selectAllVisible(filterAccountIds.value, accountsStore.allAccountsMap);

    if (props.autoSave) {
        save();
    }
}

function save(): void {
    const changed = saveFilterAccountIds();
    emit('settings:change', changed);
}

function cancel(): void {
    emit('settings:change', false);
}

init();
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
