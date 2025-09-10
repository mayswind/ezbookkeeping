<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableAccount }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="tt(applyText)" :class="{ 'disabled': !hasAnyVisibleAccount }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <div class="skeleton-text" v-if="loading">
            <f7-block class="combination-list-wrapper margin-vertical"
                      :key="blockIdx" v-for="blockIdx in [ 1, 2, 3 ]">
                <f7-accordion-item>
                    <f7-block-title>
                        <f7-accordion-toggle>
                            <f7-list strong inset dividers
                                     class="combination-list-header combination-list-opened">
                                <f7-list-item group-title>
                                    <small>Account Category</small>
                                    <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
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
        </div>

        <f7-list strong inset dividers accordion-list class="margin-top" v-if="!loading && !hasAnyVisibleAccount">
            <f7-list-item :title="tt('No available account')"></f7-list-item>
        </f7-list>

        <f7-block class="no-margin no-padding" v-show="!loading && hasAnyVisibleAccount">
            <f7-block class="combination-list-wrapper margin-vertical"
                      :key="accountCategory.category"
                      v-for="accountCategory in allCategorizedAccounts"
                      v-show="showHidden || accountCategory.allVisibleAccountCount > 0">
                <f7-accordion-item :opened="collapseStates[accountCategory.category]!.opened"
                                   @accordion:open="collapseStates[accountCategory.category]!.opened = true"
                                   @accordion:close="collapseStates[accountCategory.category]!.opened = false">
                    <f7-block-title>
                        <f7-accordion-toggle>
                            <f7-list strong inset dividers
                                     class="combination-list-header"
                                     :class="collapseStates[accountCategory.category]!.opened ? 'combination-list-opened' : 'combination-list-closed'">
                                <f7-list-item group-title>
                                    <small>{{ tt(accountCategory.name) }}</small>
                                    <f7-icon class="combination-list-chevron-icon" :f7="collapseStates[accountCategory.category]!.opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                                </f7-list-item>
                            </f7-list>
                        </f7-accordion-toggle>
                    </f7-block-title>
                    <f7-accordion-content :style="{ height: collapseStates[accountCategory.category]!.opened ? 'auto' : '' }">
                        <f7-list strong inset dividers accordion-list class="combination-list-content">
                            <f7-list-item checkbox
                                          :class="{ 'has-child-list-item': account.type === AccountType.MultiSubAccounts.type && ((showHidden && accountCategory.allSubAccounts[account.id]) || accountCategory.allVisibleSubAccountCounts[account.id]) }"
                                          :title="account.name"
                                          :value="account.id"
                                          :checked="isAccountOrSubAccountsAllChecked(account, filterAccountIds)"
                                          :indeterminate="isAccountOrSubAccountsHasButNotAllChecked(account, filterAccountIds)"
                                          :key="account.id"
                                          v-for="account in accountCategory.allAccounts"
                                          v-show="showHidden || !account.hidden"
                                          @change="updateAccountOrSubAccountsSelected">
                                <template #media>
                                    <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color">
                                        <f7-badge color="gray" class="right-bottom-icon" v-if="account.hidden">
                                            <f7-icon f7="eye_slash_fill"></f7-icon>
                                        </f7-badge>
                                    </ItemIcon>
                                </template>

                                <template #root>
                                    <ul class="padding-inline-start"
                                        v-if="account.type === AccountType.MultiSubAccounts.type && ((showHidden && accountCategory.allSubAccounts[account.id]) || accountCategory.allVisibleSubAccountCounts[account.id])">
                                        <f7-list-item checkbox
                                                      :title="subAccount.name"
                                                      :value="subAccount.id"
                                                      :checked="isAccountChecked(subAccount, filterAccountIds)"
                                                      :key="subAccount.id"
                                                      v-for="subAccount in accountCategory.allSubAccounts[account.id]"
                                                      v-show="showHidden || !subAccount.hidden"
                                                      @change="updateAccountSelected">
                                            <template #media>
                                                <ItemIcon icon-type="account" :icon-id="subAccount.icon" :color="subAccount.color">
                                                    <f7-badge color="gray" class="right-bottom-icon" v-if="subAccount.hidden">
                                                        <f7-icon f7="eye_slash_fill"></f7-icon>
                                                    </f7-badge>
                                                </ItemIcon>
                                            </template>
                                        </f7-list-item>
                                    </ul>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-content>
                </f7-accordion-item>
            </f7-block>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleAccount }" @click="selectAllAccounts">{{ tt('Select All') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleAccount }" @click="selectNoneAccounts">{{ tt('Select None') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleAccount }" @click="selectInvertAccounts">{{ tt('Invert Selection') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group v-if="allowHiddenAccount">
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ tt('Show Hidden Accounts') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ tt('Hide Hidden Accounts') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import {
    type AccountFilterType,
    useAccountFilterSettingPageBase
} from '@/views/base/settings/AccountFilterSettingPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';

import { AccountType, AccountCategory } from '@/core/account.ts';
import {
    selectAccountOrSubAccounts,
    selectAll,
    selectNone,
    selectInvert,
    isAccountOrSubAccountsAllChecked,
    isAccountOrSubAccountsHasButNotAllChecked
} from '@/lib/account.ts';

interface CollapseState {
    opened: boolean;
}

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const query = props.f7route.query;

const { tt } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

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
} = useAccountFilterSettingPageBase(query['type'] as AccountFilterType);

const accountsStore = useAccountsStore();

const collapseStates = ref<Record<number, CollapseState>>(getCollapseStates());
const loadingError = ref<unknown | null>(null);
const showMoreActionSheet = ref<boolean>(false);

function getCollapseStates(): Record<number, CollapseState> {
    const collapseStates: Record<number, CollapseState> = {};
    const allCategories = AccountCategory.values();

    for (const accountCategory of allCategories) {
        collapseStates[accountCategory.type] = {
            opened: true
        };
    }

    return collapseStates;
}

function init(): void {
    accountsStore.loadAllAccounts({
        force: false
    }).then(() => {
        loading.value = false;

        if (!loadFilterAccountIds()) {
            showToast('Parameter Invalid');
            loadingError.value = 'Parameter Invalid';
        }
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function updateAccountOrSubAccountsSelected(e: Event): void {
    const target = e.target as HTMLInputElement;
    const accountId = target.value;
    const account = accountsStore.allAccountsMap[accountId];

    if (!account) {
        return;
    }

    selectAccountOrSubAccounts(filterAccountIds.value, account, !target.checked);
}

function updateAccountSelected(e: Event): void {
    const target = e.target as HTMLInputElement;
    const accountId = target.value;
    const account = accountsStore.allAccountsMap[accountId];

    if (!account) {
        return;
    }

    filterAccountIds.value[account.id] = !target.checked;
}

function selectAllAccounts(): void {
    selectAll(filterAccountIds.value, accountsStore.allAccountsMap, !allowHiddenAccount.value);
}

function selectNoneAccounts(): void {
    selectNone(filterAccountIds.value, accountsStore.allAccountsMap, !allowHiddenAccount.value);
}

function selectInvertAccounts(): void {
    selectInvert(filterAccountIds.value, accountsStore.allAccountsMap, !allowHiddenAccount.value);
}

function save(): void {
    saveFilterAccountIds();
    props.f7router.back();
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

init();
</script>
