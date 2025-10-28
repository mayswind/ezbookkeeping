<template>
    <v-card :class="{ 'pa-2 pa-sm-4 pa-md-8': dialogMode }">
        <template #title>
            <div class="d-flex align-center justify-center" v-if="dialogMode">
                <div class="w-100 text-center">
                    <h4 class="text-h4">{{ tt(title) }}</h4>
                </div>
                <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                       :disabled="loading || !hasAnyAvailableTag" :icon="true">
                    <v-icon :icon="mdiDotsVertical" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectAllTransactionTags"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelect"
                                         :title="tt('Select None')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectNoneTransactionTags"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelectInverse"
                                         :title="tt('Invert Selection')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectInvertTransactionTags"></v-list-item>
                            <v-divider class="my-2"/>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All Visible')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectAllVisibleTransactionTags"></v-list-item>
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
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectAllTransactionTags"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelect"
                                         :title="tt('Select None')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectNoneTransactionTags"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelectInverse"
                                         :title="tt('Invert Selection')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectInvertTransactionTags"></v-list-item>
                            <v-divider class="my-2"/>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All Visible')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectAllVisibleTransactionTags"></v-list-item>
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
            </div>
        </template>

        <div v-if="loading">
            <v-skeleton-loader type="paragraph" :loading="loading"
                               :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
        </div>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-if="!loading && !hasAnyVisibleTag">
            <span class="text-body-1">{{ tt('No available tag') }}</span>
        </v-card-text>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-else-if="!loading && hasAnyVisibleTag">
            <div class="tag-filter-types d-flex flex-column mb-4" v-if="type === 'statisticsCurrent'">
                <v-btn border class="justify-start" :key="filterType.type"
                       :color="tagFilterType === filterType.type ? 'primary' : 'default'"
                       :variant="tagFilterType === filterType.type ? 'tonal' : 'outlined'"
                       :append-icon="(tagFilterType === filterType.type ? mdiCheck : undefined)"
                       v-for="filterType in allTagFilterTypes"
                       @click="tagFilterType = filterType.type">
                    {{ filterType.displayName }}
                </v-btn>
            </div>

            <v-expansion-panels class="tag-categories" multiple v-model="expandTagCategories">
                <v-expansion-panel class="border" key="default" value="default">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ms-3">{{ tt('Tags') }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                        <v-list rounded density="comfortable" class="pa-0">
                            <template :key="transactionTag.id"
                                      v-for="transactionTag in allTags">
                                <v-list-item v-if="showHidden || !transactionTag.hidden">
                                    <template #prepend>
                                        <v-checkbox :model-value="!filterTagIds[transactionTag.id]"
                                                    @update:model-value="updateTransactionTagSelected(transactionTag, $event)">
                                            <template #label>
                                                <v-badge class="right-bottom-icon" color="secondary"
                                                         location="bottom right" offset-x="2" offset-y="2" :icon="mdiEyeOffOutline"
                                                         v-if="transactionTag.hidden">
                                                    <v-icon size="24" :icon="mdiPound"/>
                                                </v-badge>
                                                <v-icon size="24" :icon="mdiPound" v-else-if="!transactionTag.hidden"/>
                                                <span class="ms-3">{{ transactionTag.name }}</span>
                                            </template>
                                        </v-checkbox>
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
                <v-btn :disabled="!hasAnyVisibleTag" @click="save">{{ tt(applyText) }}</v-btn>
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
import { useTransactionTagFilterSettingPageBase } from '@/views/base/settings/TransactionTagFilterSettingPageBase.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import type { TransactionTag } from '@/models/transaction_tag.ts';

import {
    selectAllVisible,
    selectAll,
    selectNone,
    selectInvert
} from '@/lib/common.ts';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiDotsVertical,
    mdiCheck,
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
    filterTagIds,
    tagFilterType,
    title,
    applyText,
    allTags,
    allTagFilterTypes,
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

function updateTransactionTagSelected(transactionTag: TransactionTag, value: boolean | null): void {
    filterTagIds.value[transactionTag.id] = !value;

    if (props.autoSave) {
        save();
    }
}

function selectAllTransactionTags(): void {
    selectAll(filterTagIds.value, transactionTagsStore.allTransactionTagsMap);

    if (props.autoSave) {
        save();
    }
}

function selectNoneTransactionTags(): void {
    selectNone(filterTagIds.value, transactionTagsStore.allTransactionTagsMap);

    if (props.autoSave) {
        save();
    }
}

function selectInvertTransactionTags(): void {
    selectInvert(filterTagIds.value, transactionTagsStore.allTransactionTagsMap);

    if (props.autoSave) {
        save();
    }
}

function selectAllVisibleTransactionTags(): void {
    selectAllVisible(filterTagIds.value, transactionTagsStore.allTransactionTagsMap);

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
.tag-filter-types .v-btn:not(:first-child) {
    border-top-left-radius: inherit;
    border-top-right-radius: inherit;
}

.tag-filter-types .v-btn:not(:last-child) {
    border-bottom: 0;
    border-bottom-left-radius: inherit;
    border-bottom-right-radius: inherit;
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
