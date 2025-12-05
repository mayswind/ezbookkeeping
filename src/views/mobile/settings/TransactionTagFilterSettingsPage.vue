<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt(title)"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableTag }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': !hasAnyVisibleTag }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="combination-list-wrapper margin-vertical skeleton-text" v-if="loading">
            <f7-accordion-item>
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers
                                 class="combination-list-header combination-list-opened">
                            <f7-list-item group-title>
                                <small>Tags</small>
                                <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content style="height: auto">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item checkbox class="disabled" title="Tag Name"
                                      :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                            <template #media>
                                <f7-icon f7="app_fill"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-list strong inset dividers accordion-list class="margin-top" v-if="!loading && !hasAnyVisibleTag">
            <f7-list-item :title="tt('No available tag')"></f7-list-item>
        </f7-list>

        <f7-block class="combination-list-wrapper margin-vertical" key="default" v-show="!loading && hasAnyVisibleTag">
            <f7-list class="margin-top-half margin-bottom" strong inset dividers v-if="includeTagsCount > 1">
                <f7-list-item radio
                              :title="tt(filterType.name)"
                              :key="filterType.type"
                              :value="filterType.type"
                              :checked="includeTagFilterType === filterType.type"
                              v-for="filterType in [TransactionTagFilterType.HasAny, TransactionTagFilterType.HasAll]"
                              @change="includeTagFilterType = filterType.type">
                </f7-list-item>
            </f7-list>

            <f7-list class="margin-top-half margin-bottom" strong inset dividers v-if="excludeTagsCount > 1">
                <f7-list-item radio
                              :title="tt(filterType.name)"
                              :key="filterType.type"
                              :value="filterType.type"
                              :checked="excludeTagFilterType === filterType.type"
                              v-for="filterType in [TransactionTagFilterType.NotHasAny, TransactionTagFilterType.NotHasAll]"
                              @change="excludeTagFilterType = filterType.type">
                </f7-list-item>
            </f7-list>

            <f7-accordion-item :opened="collapseStates['default']!.opened"
                               @accordion:open="collapseStates['default']!.opened = true"
                               @accordion:close="collapseStates['default']!.opened = false">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers
                                 class="combination-list-header"
                                 :class="collapseStates['default']!.opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item group-title>
                                <small>{{ tt('Tags') }}</small>
                                <f7-icon class="combination-list-chevron-icon" :f7="collapseStates['default']!.opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates['default']!.opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item link="#"
                                      popover-open=".tag-filter-state-popover-menu"
                                      :title="transactionTag.name"
                                      :value="transactionTag.id"
                                      :key="transactionTag.id"
                                      :after="tt(filterTagIds[transactionTag.id] === TransactionTagFilterState.Include ? 'Included' : filterTagIds[transactionTag.id] === TransactionTagFilterState.Exclude ? 'Excluded' : 'Default')"
                                      v-for="transactionTag in allTags"
                                      v-show="showHidden || !transactionTag.hidden"
                                      @click="currentTransactionTagId = transactionTag.id">
                            <template #media>
                                <f7-icon class="transaction-tag-icon" f7="number">
                                    <f7-badge color="gray" class="right-bottom-icon" v-if="transactionTag.hidden">
                                        <f7-icon f7="eye_slash_fill"></f7-icon>
                                    </f7-badge>
                                </f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-popover class="tag-filter-state-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="state.displayName"
                              :class="{ 'list-item-selected': filterTagIds[currentTransactionTagId] === state.type }"
                              :key="state.type"
                              v-for="state in [
                                  { type: TransactionTagFilterState.Include, displayName: tt('Included') },
                                  { type: TransactionTagFilterState.Default, displayName: tt('Default') },
                                  { type: TransactionTagFilterState.Exclude, displayName: tt('Excluded') }
                              ]"
                              @click="updateCurrentTransactionTagState(state.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="filterTagIds[currentTransactionTagId] === state.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(false, TransactionTagFilterState.Include)">{{ tt('Set All to Included') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(false, TransactionTagFilterState.Default)">{{ tt('Set All to Default') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(false, TransactionTagFilterState.Exclude)">{{ tt('Set All to Excluded') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(true, TransactionTagFilterState.Include)">{{ tt('Set All Visible Items to Included') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(true, TransactionTagFilterState.Default)">{{ tt('Set All Visible Items to Default') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(true, TransactionTagFilterState.Exclude)">{{ tt('Set All Visible Items to Excluded') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ tt('Show Hidden Transaction Tags') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ tt('Hide Hidden Transaction Tags') }}</f7-actions-button>
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
    useTransactionTagFilterSettingPageBase,
    TransactionTagFilterState
} from '@/views/base/settings/TransactionTagFilterSettingPageBase.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { TransactionTagFilterType } from '@/core/transaction.ts';

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
    filterTagIds,
    includeTagFilterType,
    excludeTagFilterType,
    includeTagsCount,
    excludeTagsCount,
    title,
    allTags,
    hasAnyAvailableTag,
    hasAnyVisibleTag,
    loadFilterTagIds,
    saveFilterTagIds
} = useTransactionTagFilterSettingPageBase(query['type']);

const transactionTagsStore = useTransactionTagsStore();

const loadingError = ref<unknown | null>(null);
const currentTransactionTagId = ref<string>('');
const showMoreActionSheet = ref<boolean>(false);

const collapseStates = ref<Record<string, CollapseState>>({
    default: {
        opened: true
    }
});

function init(): void {
    transactionTagsStore.loadAllTags({
        force: false
    }).then(() => {
        loading.value = false;

        if (!loadFilterTagIds()) {
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

function updateCurrentTransactionTagState(state: number): void {
    filterTagIds.value[currentTransactionTagId.value] = state;
    currentTransactionTagId.value = '';
}

function setAllTagsState(onlyVisible: boolean, value: TransactionTagFilterState): void {
    for (const tag of allTags.value) {
        if (onlyVisible && !showHidden.value && tag.hidden) {
            continue;
        }

        filterTagIds.value[tag.id] = value;
    }
}

function save(): void {
    saveFilterTagIds();
    props.f7router.back();
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

init();
</script>
