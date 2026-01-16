<template>
    <f7-page with-subnavbar @page:beforein="onPageBeforeIn" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt(title)"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableTag }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': !hasAnyAvailableTag }" @click="save"></f7-link>
            </f7-nav-right>

            <f7-subnavbar :inner="false">
                <f7-searchbar
                    custom-searchs
                    :value="filterContent"
                    :placeholder="tt('Find tag')"
                    :disable-button-text="tt('Cancel')"
                    @change="filterContent = $event.target.value"
                ></f7-searchbar>
            </f7-subnavbar>
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

        <f7-block class="combination-list-wrapper margin-vertical"
                  :key="tagGroup.id" v-for="tagGroup in allTagGroupsWithDefault"
                  v-show="!loading && hasAnyVisibleTag">
            <f7-accordion-item :opened="collapseStates[tagGroup.id]?.opened ?? true"
                               @accordion:open="collapseStates[tagGroup.id]!.opened = true"
                               @accordion:close="collapseStates[tagGroup.id]!.opened = false"
                               v-if="allVisibleTags[tagGroup.id] && allVisibleTags[tagGroup.id]!.length > 0">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers
                                 class="combination-list-header"
                                 :class="collapseStates[tagGroup.id]?.opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item group-title>
                                <small class="tag-group-title">{{ tagGroup.name }}</small>
                                <f7-icon class="combination-list-chevron-icon" :f7="collapseStates[tagGroup.id]?.opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates[tagGroup.id]?.opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item link="#"
                                      popover-open=".tag-filter-include-type-popover-menu"
                                      :title="tt(TransactionTagFilterType.parse(groupTagFilterTypesMap[tagGroup.id]?.includeType as number)?.name as string)"
                                      @click="currentTransactionTagGroupId = tagGroup.id"
                                      v-if="groupTagFilterTypesMap[tagGroup.id] && groupTagFilterStateCountMap[tagGroup.id] && groupTagFilterStateCountMap[tagGroup.id]![TransactionTagFilterState.Include] && groupTagFilterStateCountMap[tagGroup.id]![TransactionTagFilterState.Include] > 1 && TransactionTagFilterType.parse(groupTagFilterTypesMap[tagGroup.id]?.includeType as number)">
                        </f7-list-item>
                        <f7-list-item link="#"
                                      popover-open=".tag-filter-exclude-type-popover-menu"
                                      :title="tt(TransactionTagFilterType.parse(groupTagFilterTypesMap[tagGroup.id]?.excludeType as number)?.name as string)"
                                      @click="currentTransactionTagGroupId = tagGroup.id"
                                      v-if="groupTagFilterTypesMap[tagGroup.id] && groupTagFilterStateCountMap[tagGroup.id] && groupTagFilterStateCountMap[tagGroup.id]![TransactionTagFilterState.Exclude] && groupTagFilterStateCountMap[tagGroup.id]![TransactionTagFilterState.Exclude] > 1 && TransactionTagFilterType.parse(groupTagFilterTypesMap[tagGroup.id]?.excludeType as number)">
                        </f7-list-item>
                        <f7-list-item link="#"
                                      popover-open=".tag-filter-state-popover-menu"
                                      :title="transactionTag.name"
                                      :value="transactionTag.id"
                                      :key="transactionTag.id"
                                      :after="tt(tagFilterStateMap[transactionTag.id] === TransactionTagFilterState.Include ? 'Included' : tagFilterStateMap[transactionTag.id] === TransactionTagFilterState.Exclude ? 'Excluded' : 'Default')"
                                      v-for="transactionTag in allVisibleTags[tagGroup.id]"
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

        <f7-popover class="tag-filter-include-type-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt(filterType.name)"
                              :class="{ 'list-item-selected': groupTagFilterTypesMap[currentTransactionTagGroupId]?.includeType === filterType.type }"
                              :key="filterType.type"
                              v-for="filterType in [TransactionTagFilterType.HasAny, TransactionTagFilterType.HasAll]"
                              @click="updateTransactionTagGroupIncludeType(filterType)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="groupTagFilterTypesMap[currentTransactionTagGroupId]?.includeType === filterType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="tag-filter-exclude-type-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="tt(filterType.name)"
                              :class="{ 'list-item-selected': groupTagFilterTypesMap[currentTransactionTagGroupId]?.excludeType === filterType.type }"
                              :key="filterType.type"
                              v-for="filterType in [TransactionTagFilterType.NotHasAny, TransactionTagFilterType.NotHasAll]"
                              @click="updateTransactionTagGroupExcludeType(filterType)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="groupTagFilterTypesMap[currentTransactionTagGroupId]?.excludeType === filterType.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="tag-filter-state-popover-menu">
            <f7-list dividers>
                <f7-list-item link="#" no-chevron popover-close
                              :title="state.displayName"
                              :class="{ 'list-item-selected': tagFilterStateMap[currentTransactionTagId] === state.type }"
                              :key="state.type"
                              v-for="state in [
                                  { type: TransactionTagFilterState.Include, displayName: tt('Included') },
                                  { type: TransactionTagFilterState.Default, displayName: tt('Default') },
                                  { type: TransactionTagFilterState.Exclude, displayName: tt('Excluded') }
                              ]"
                              @click="updateCurrentTransactionTagState(state.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="tagFilterStateMap[currentTransactionTagId] === state.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(TransactionTagFilterState.Include)">{{ tt('Set All to Included') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(TransactionTagFilterState.Default)">{{ tt('Set All to Default') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="setAllTagsState(TransactionTagFilterState.Exclude)">{{ tt('Set All to Excluded') }}</f7-actions-button>
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

import { values } from '@/core/base.ts';
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
    filterContent,
    tagFilterStateMap,
    groupTagFilterTypesMap,
    title,
    groupTagFilterStateCountMap,
    allTagGroupsWithDefault,
    allVisibleTags,
    allVisibleTagGroupIds,
    hasAnyAvailableTag,
    hasAnyVisibleTag,
    loadFilterTagIds,
    saveFilterTagIds
} = useTransactionTagFilterSettingPageBase(query['type']);

const transactionTagsStore = useTransactionTagsStore();

const loadingError = ref<unknown | null>(null);
const currentTransactionTagGroupId = ref<string>('');
const currentTransactionTagId = ref<string>('');
const showMoreActionSheet = ref<boolean>(false);

const collapseStates = ref<Record<string, CollapseState>>(getInitCollapseState(allVisibleTagGroupIds.value));

function getInitCollapseState(tagGroupIds: string[]): Record<string, CollapseState> {
    const states: Record<string, CollapseState> = {};

    for (const tagGroupId of tagGroupIds) {
        states[tagGroupId] = {
            opened: true
        };
    }

    return states;
}

function init(): void {
    transactionTagsStore.loadAllTags({
        force: false
    }).then(() => {
        loading.value = false;
        collapseStates.value = getInitCollapseState(allVisibleTagGroupIds.value);

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

function updateTransactionTagGroupIncludeType(filterType: TransactionTagFilterType): void {
    const tagFilterTypes = groupTagFilterTypesMap.value[currentTransactionTagGroupId.value];

    if (!tagFilterTypes) {
        return;
    }

    tagFilterTypes.includeType = filterType.type;
}

function updateTransactionTagGroupExcludeType(filterType: TransactionTagFilterType): void {
    const tagFilterTypes = groupTagFilterTypesMap.value[currentTransactionTagGroupId.value];

    if (!tagFilterTypes) {
        return;
    }

    tagFilterTypes.excludeType = filterType.type;
}

function updateCurrentTransactionTagState(state: number): void {
    tagFilterStateMap.value[currentTransactionTagId.value] = state;
    currentTransactionTagId.value = '';
}

function setAllTagsState(value: TransactionTagFilterState): void {
    for (const tags of values(allVisibleTags.value)) {
        for (const tag of tags) {
            tagFilterStateMap.value[tag.id] = value;
        }
    }
}

function save(): void {
    saveFilterTagIds();
    props.f7router.back();
}

function onPageBeforeIn(): void {
    filterContent.value = '';
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

init();
</script>

<style>
.tag-group-title {
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
