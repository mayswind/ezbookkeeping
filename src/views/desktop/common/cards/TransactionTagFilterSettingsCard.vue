<template>
    <v-card :class="{ 'pa-sm-1 pa-md-2': dialogMode }">
        <template #title>
            <v-row>
                <v-col cols="6">
                    <div :class="{ 'text-h4': dialogMode, 'text-wrap': true }">
                        {{ tt(title) }}
                    </div>
                </v-col>
                <v-col cols="6" class="d-flex align-center">
                    <v-spacer v-if="!dialogMode"/>
                    <v-text-field density="compact" :disabled="loading || !hasAnyAvailableTag"
                                  :prepend-inner-icon="mdiMagnify"
                                  :placeholder="tt('Find tag')"
                                  v-model="filterContent"
                                  v-if="dialogMode"></v-text-field>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :disabled="loading || !hasAnyAvailableTag" :icon="true">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiSelectAll"
                                             :title="tt('Set All to Included')"
                                             :disabled="!hasAnyVisibleTag"
                                             @click="setAllTagsState(TransactionTagFilterState.Include)"></v-list-item>
                                <v-list-item :prepend-icon="mdiSelectAll"
                                             :title="tt('Set All to Default')"
                                             :disabled="!hasAnyVisibleTag"
                                             @click="setAllTagsState(TransactionTagFilterState.Default)"></v-list-item>
                                <v-list-item :prepend-icon="mdiSelectAll"
                                             :title="tt('Set All to Excluded')"
                                             :disabled="!hasAnyVisibleTag"
                                             @click="setAllTagsState(TransactionTagFilterState.Exclude)"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-item :prepend-icon="mdiEyeOutline"
                                             :title="tt('Show Hidden Transaction Tags')"
                                             v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                <v-list-item :prepend-icon="mdiEyeOffOutline"
                                             :title="tt('Hide Hidden Transaction Tags')"
                                             v-if="showHidden" @click="showHidden = false"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </v-col>
            </v-row>
        </template>

        <div v-if="loading">
            <v-skeleton-loader type="paragraph" :loading="loading"
                               :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
        </div>

        <v-card-text v-if="!loading && !hasAnyVisibleTag">
            <span class="text-body-1">{{ tt('No available tag') }}</span>
        </v-card-text>

        <v-card-text :class="{ 'flex-grow-1 overflow-y-auto': dialogMode }" v-else-if="!loading && hasAnyVisibleTag">
            <v-expansion-panels class="tag-categories" multiple v-model="expandTagGroups">
                <template :key="tagGroup.id" v-for="tagGroup in allTagGroupsWithDefault">
                    <v-expansion-panel class="border" :value="tagGroup.id" v-if="allVisibleTags[tagGroup.id] && allVisibleTags[tagGroup.id]!.length > 0">
                        <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                            <span class="ms-3 text-truncate">{{ tagGroup.name }}</span>
                            <v-spacer/>
                            <div class="d-flex me-3" v-if="groupTagFilterTypesMap[tagGroup.id] && groupTagFilterStateCountMap[tagGroup.id]">
                                <v-btn color="secondary" density="compact" variant="outlined"
                                       v-if="groupTagFilterStateCountMap[tagGroup.id]![TransactionTagFilterState.Include] && groupTagFilterStateCountMap[tagGroup.id]![TransactionTagFilterState.Include] > 1">
                                    {{ groupTagFilterTypesMap[tagGroup.id]!.includeType === TransactionTagFilterType.HasAll.type ? tt(TransactionTagFilterType.HasAll.name) : tt(TransactionTagFilterType.HasAny.name) }}
                                    <v-menu activator="parent">
                                        <v-list>
                                            <v-list-item :key="filterType.type" :title="tt(filterType.name)"
                                                         :append-icon="groupTagFilterTypesMap[tagGroup.id]!.includeType === filterType.type ? mdiCheck : undefined"
                                                         v-for="filterType in [TransactionTagFilterType.HasAny, TransactionTagFilterType.HasAll]"
                                                         @click="updateTransactionTagGroupIncludeType(tagGroup, filterType)"></v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                                <v-btn class="ms-2" color="secondary" density="compact" variant="outlined"
                                       v-if="groupTagFilterStateCountMap[tagGroup.id]![TransactionTagFilterState.Exclude] && groupTagFilterStateCountMap[tagGroup.id]![TransactionTagFilterState.Exclude] > 1">
                                    {{ groupTagFilterTypesMap[tagGroup.id]!.excludeType === TransactionTagFilterType.NotHasAll.type ? tt(TransactionTagFilterType.NotHasAll.name) : tt(TransactionTagFilterType.NotHasAny.name) }}
                                    <v-menu activator="parent">
                                        <v-list>
                                            <v-list-item :key="filterType.type" :title="tt(filterType.name)"
                                                         :append-icon="groupTagFilterTypesMap[tagGroup.id]!.excludeType === filterType.type ? mdiCheck : undefined"
                                                         v-for="filterType in [TransactionTagFilterType.NotHasAny, TransactionTagFilterType.NotHasAll]"
                                                         @click="updateTransactionTagGroupExcludeType(tagGroup, filterType)"></v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                            </div>
                        </v-expansion-panel-title>
                        <v-expansion-panel-text>
                            <v-list rounded density="comfortable" class="pa-0">
                                <template :key="transactionTag.id"
                                          v-for="transactionTag in allVisibleTags[tagGroup.id]">
                                    <v-list-item class="ps-2">
                                        <template #prepend>
                                            <v-badge class="right-bottom-icon" color="secondary"
                                                     location="bottom right" offset-x="2" offset-y="2" :icon="mdiEyeOffOutline"
                                                     v-if="transactionTag.hidden">
                                                <v-icon size="24" :icon="mdiPound"/>
                                            </v-badge>
                                            <v-icon size="24" :icon="mdiPound" v-else-if="!transactionTag.hidden"/>
                                            <span class="ms-3">{{ transactionTag.name }}</span>
                                        </template>
                                        <template #append>
                                            <v-btn-toggle class="toggle-buttons" density="compact" variant="outlined"
                                                          mandatory="force" divided
                                                          :model-value="tagFilterStateMap[transactionTag.id]"
                                                          @update:model-value="updateTransactionTagState(transactionTag, $event)">
                                                <v-btn :value="TransactionTagFilterState.Include">{{ tt('Included') }}</v-btn>
                                                <v-btn :value="TransactionTagFilterState.Default">{{ tt('Default') }}</v-btn>
                                                <v-btn :value="TransactionTagFilterState.Exclude">{{ tt('Excluded') }}</v-btn>
                                            </v-btn-toggle>
                                        </template>
                                    </v-list-item>
                                </template>
                            </v-list>
                        </v-expansion-panel-text>
                    </v-expansion-panel>
                </template>
            </v-expansion-panels>
        </v-card-text>

        <v-card-text class="overflow-y-visible" v-if="dialogMode">
            <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                <v-btn :disabled="!hasAnyAvailableTag" @click="save">{{ tt(applyText) }}</v-btn>
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
    useTransactionTagFilterSettingPageBase,
    TransactionTagFilterState
} from '@/views/base/settings/TransactionTagFilterSettingPageBase.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { values } from '@/core/base.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';

import type { TransactionTagGroup } from '@/models/transaction_tag_group.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';

import {
    mdiMagnify,
    mdiCheck,
    mdiSelectAll,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiDotsVertical,
    mdiPound
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    type: string;
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
    filterContent,
    tagFilterStateMap,
    groupTagFilterTypesMap,
    title,
    applyText,
    groupTagFilterStateCountMap,
    allTagGroupsWithDefault,
    allVisibleTags,
    allVisibleTagGroupIds,
    hasAnyAvailableTag,
    hasAnyVisibleTag,
    loadFilterTagIds,
    saveFilterTagIds
} = useTransactionTagFilterSettingPageBase(props.type);

const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const expandTagGroups = ref<string[]>(allVisibleTagGroupIds.value);

function init(): void {
    transactionTagsStore.loadAllTags({
        force: false
    }).then(() => {
        loading.value = false;
        expandTagGroups.value = allVisibleTagGroupIds.value;

        if (!loadFilterTagIds()) {
            snackbar.value?.showError('Parameter Invalid');
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function updateTransactionTagState(transactionTag: TransactionTag, value: TransactionTagFilterState): void {
    tagFilterStateMap.value[transactionTag.id] = value;

    if (props.autoSave) {
        save();
    }
}

function updateTransactionTagGroupIncludeType(tagGroup: TransactionTagGroup, filterType: TransactionTagFilterType): void {
    const tagFilterTypes = groupTagFilterTypesMap.value[tagGroup.id];

    if (!tagFilterTypes) {
        return;
    }

    tagFilterTypes.includeType = filterType.type;

    if (props.autoSave) {
        save();
    }
}

function updateTransactionTagGroupExcludeType(tagGroup: TransactionTagGroup, filterType: TransactionTagFilterType): void {
    const tagFilterTypes = groupTagFilterTypesMap.value[tagGroup.id];

    if (!tagFilterTypes) {
        return;
    }

    tagFilterTypes.excludeType = filterType.type;

    if (props.autoSave) {
        save();
    }
}

function setAllTagsState(value: TransactionTagFilterState): void {
    for (const tags of values(allVisibleTags.value)) {
        for (const tag of tags) {
            tagFilterStateMap.value[tag.id] = value;
        }
    }

    if (props.autoSave) {
        save();
    }
}

function save(): void {
    const changed = saveFilterTagIds();
    emit('settings:change', changed);
}

function cancel(): void {
    emit('settings:change', false);
}

init();
</script>

<style>
.tag-categories .tag-filter-state-toggle {
    overflow-x: auto;
    white-space: nowrap;
}

.tag-categories .v-expansion-panel-text__wrapper {
    padding: 0 0 0 0;
    padding-inline-start: 20px;
}

.tag-categories .v-expansion-panel--active:not(:first-child),
.tag-categories .v-expansion-panel--active + .v-expansion-panel {
    margin-top: 1rem;
}
</style>
