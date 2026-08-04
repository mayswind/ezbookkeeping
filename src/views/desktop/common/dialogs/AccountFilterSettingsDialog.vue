<template>
    <v-dialog width="800" v-model="showState">
        <one-column-dialog-layout content-class="pt-4"
                                  :title="tt(title)" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-text-field class="mx-2" density="compact"
                              :disabled="loading || !hasAnyAvailableAccount"
                              :prepend-inner-icon="mdiMagnify"
                              :placeholder="tt('Find account')"
                              v-model="filterContent"></v-text-field>

                <v-btn class="ms-2" density="comfortable" variant="outlined"
                       :disabled="!hasAnyAvailableAccount" @click="save">{{ tt(applyText) }}</v-btn>

                <v-btn density="compact" color="default" variant="text" class="ms-2"
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
                            <v-list-item :prepend-icon="mdiEyeOutline"
                                         :title="tt('Show Hidden Accounts')"
                                         v-if="allowHiddenAccount && !showHidden" @click="showHidden = true"></v-list-item>
                            <v-list-item :prepend-icon="mdiEyeOffOutline"
                                         :title="tt('Hide Hidden Accounts')"
                                         v-if="allowHiddenAccount && showHidden" @click="showHidden = false"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </template>

            <template #content>
                <div v-if="loading">
                    <v-skeleton-loader type="paragraph" :loading="loading"
                                       :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
                </div>

                <div v-if="!loading && !hasAnyVisibleAccount">
                    <span class="text-body-1">{{ tt('No available account') }}</span>
                </div>

                <div v-else-if="!loading && hasAnyVisibleAccount">
                    <v-expansion-panels class="account-categories" multiple v-model="expandAccountCategories">
                        <v-expansion-panel :key="accountCategory.category"
                                           :value="accountCategory.category"
                                           class="border"
                                           v-for="accountCategory in allCategorizedAccounts">
                            <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                                <span class="ms-3">{{ tt(accountCategory.name) }}</span>
                            </v-expansion-panel-title>
                            <v-expansion-panel-text>
                                <v-list rounded density="comfortable" class="pa-0">
                                    <template :key="account.id"
                                              v-for="(account, idx) in accountCategory.accounts">
                                        <v-divider v-if="idx > 0"/>

                                        <v-list-item>
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

                                        <v-divider v-if="account.type === AccountType.MultiSubAccounts.type && account.subAccounts && account.subAccounts.length > 0"/>

                                        <v-list rounded density="comfortable" class="pa-0 ms-4"
                                                v-if="account.type === AccountType.MultiSubAccounts.type && account.subAccounts && account.subAccounts.length > 0">
                                            <template :key="subAccount.id"
                                                      v-for="(subAccount, subIdx) in account.subAccounts">
                                                <v-divider v-if="subIdx > 0"/>

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
                </div>
            </template>
        </one-column-dialog-layout>

        <snack-bar ref="snackbar" />
    </v-dialog>
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

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
    selectAll,
    selectNone,
    selectInvert,
    isAccountOrSubAccountsAllChecked,
    isAccountOrSubAccountsHasButNotAllChecked
} from '@/lib/account.ts';

import {
    mdiMagnify,
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
    selectedAccountIds?: string[];
    autoSave?: boolean;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'settings:change', changed: boolean, selectedAccountIds?: string[]): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const {
    loading,
    showHidden,
    filterContent,
    filterAccountIds,
    title,
    applyText,
    allowHiddenAccount,
    customAccountCategoryOrder,
    allCategorizedAccounts,
    allVisibleAccountMap,
    hasAnyAvailableAccount,
    hasAnyVisibleAccount,
    isAccountChecked,
    loadFilterAccountIds,
    saveFilterAccountIds
} = useAccountFilterSettingPageBase(props.type, props.selectedAccountIds);

const accountsStore = useAccountsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const expandAccountCategories = ref<number[]>(AccountCategory.values(customAccountCategoryOrder.value).map(category => category.type));

const showState = computed<boolean>({
    get: () => props.show || false,
    set: (value) => emit('update:show', value)
});

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
    selectAll(filterAccountIds.value, allVisibleAccountMap.value);

    if (props.autoSave) {
        save();
    }
}

function selectNoneAccounts(): void {
    selectNone(filterAccountIds.value, allVisibleAccountMap.value);

    if (props.autoSave) {
        save();
    }
}

function selectInvertAccounts(): void {
    selectInvert(filterAccountIds.value, allVisibleAccountMap.value);

    if (props.autoSave) {
        save();
    }
}

function save(): void {
    const [changed, selectedAccountIds] = saveFilterAccountIds();
    emit('settings:change', changed, selectedAccountIds);
}

function cancel(): void {
    emit('settings:change', false);
}

watch(() => props.show, (newValue) => {
    if (newValue) {
        loadFilterAccountIds();
        showHidden.value = false;
        filterContent.value = '';
    }
});

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
