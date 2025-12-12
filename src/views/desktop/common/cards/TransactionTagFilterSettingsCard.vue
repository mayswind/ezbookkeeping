<template>
    <v-card :class="{ 'pa-2 pa-sm-4 pa-md-8': dialogMode }">
        <template #title>
            <v-row>
                <v-col cols="6">
                    <div :class="{ 'text-h4': dialogMode }">
                        {{ tt(title) }}
                    </div>
                </v-col>
                <v-col cols="6" class="d-flex align-center">
                    <v-text-field eager ref="filterInput" density="compact"
                                  :prepend-inner-icon="mdiMagnify"
                                  :placeholder="tt('Find tag')"
                                  v-model="filterContent"></v-text-field>
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

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-if="!loading && !hasAnyVisibleTag">
            <span class="text-body-1">{{ tt('No available tag') }}</span>
        </v-card-text>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-else-if="!loading && hasAnyVisibleTag">
            <div class="mb-4" v-if="includeTagsCount > 1">
                <v-btn-toggle class="toggle-buttons" density="compact" variant="outlined"
                              mandatory="force" divided
                              :model-value="includeTagFilterType"
                              @update:model-value="updateTransactionTagIncludeType($event)">
                    <v-btn :value="TransactionTagFilterType.HasAny.type">{{ tt(TransactionTagFilterType.HasAny.name) }}</v-btn>
                    <v-btn :value="TransactionTagFilterType.HasAll.type">{{ tt(TransactionTagFilterType.HasAll.name) }}</v-btn>
                </v-btn-toggle>
            </div>

            <div class="mb-4" v-if="excludeTagsCount > 1">
                <v-btn-toggle class="toggle-buttons" density="compact" variant="outlined"
                              mandatory="force" divided
                              :model-value="excludeTagFilterType"
                              @update:model-value="updateTransactionTagExcludeType($event)">
                    <v-btn :value="TransactionTagFilterType.NotHasAny.type">{{ tt(TransactionTagFilterType.NotHasAny.name) }}</v-btn>
                    <v-btn :value="TransactionTagFilterType.NotHasAll.type">{{ tt(TransactionTagFilterType.NotHasAll.name) }}</v-btn>
                </v-btn-toggle>
            </div>

            <v-expansion-panels class="tag-categories" multiple v-model="expandTagCategories">
                <v-expansion-panel class="border" key="default" value="default">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ms-3">{{ tt('Tags') }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                        <v-list rounded density="comfortable" class="pa-0">
                            <template :key="transactionTag.id"
                                      v-for="transactionTag in allVisibleTags">
                                <v-list-item v-if="showHidden || !transactionTag.hidden">
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
                                                      :model-value="filterTagIds[transactionTag.id]"
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
            </v-expansion-panels>
        </v-card-text>

        <v-card-text class="overflow-y-visible" v-if="dialogMode">
            <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
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

import { TransactionTagFilterType } from '@/core/transaction.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';

import {
    mdiMagnify,
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
    filterTagIds,
    includeTagFilterType,
    excludeTagFilterType,
    includeTagsCount,
    excludeTagsCount,
    title,
    applyText,
    allVisibleTags,
    hasAnyAvailableTag,
    hasAnyVisibleTag,
    loadFilterTagIds,
    saveFilterTagIds
} = useTransactionTagFilterSettingPageBase(props.type);

const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const expandTagCategories = ref<string[]>([ 'default' ]);

function init(): void {
    transactionTagsStore.loadAllTags({
        force: false
    }).then(() => {
        loading.value = false;

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
    filterTagIds.value[transactionTag.id] = value;

    if (props.autoSave) {
        save();
    }
}

function updateTransactionTagIncludeType(value: number): void {
    includeTagFilterType.value = value;

    if (props.autoSave) {
        save();
    }
}

function updateTransactionTagExcludeType(value: number): void {
    excludeTagFilterType.value = value;

    if (props.autoSave) {
        save();
    }
}

function setAllTagsState(value: TransactionTagFilterState): void {
    for (const tag of allVisibleTags.value) {
        filterTagIds.value[tag.id] = value;
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
