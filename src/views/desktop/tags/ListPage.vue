<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <span class="text-subtitle-2">{{ tt('Total tags') }}</span>
                            <p class="transaction-tags-statistic-item-value mt-1">
                                <span v-if="!loading || totalAvailableTagsCount > 0">{{ displayTotalAvailableTagsCount }}</span>
                                <span v-else-if="loading && totalAvailableTagsCount <= 0">
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-1" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                        <v-divider />
                        <v-tabs show-arrows
                                class="scrollable-vertical-tabs"
                                style="max-height: calc(100% - 88px)"
                                direction="vertical"
                                :prev-icon="mdiMenuUp" :next-icon="mdiMenuDown"
                                :disabled="loading || updating" v-model="activeTagGroupId">
                            <v-tab class="tab-text-truncate" :disabled="loading || updating || displayOrderModified || hasEditingTag"
                                   :key="tagGroup.id" :value="tagGroup.id"
                                   v-for="tagGroup in allTagGroupsWithDefault"
                                   @click="switchTagGroup(tagGroup.id)">
                                <span class="text-truncate">{{ tagGroup.name }}</span>
                            </v-tab>
                            <template v-if="loading && (!allTagGroupsWithDefault || allTagGroupsWithDefault.length < 2)">
                                <v-skeleton-loader class="skeleton-no-margin mx-5 mt-4 mb-3" type="text"
                                                   :key="itemIdx" :loading="true" v-for="itemIdx in [ 1, 2, 3, 4, 5 ]"></v-skeleton-loader>
                            </template>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="tagListPage">
                                <v-card variant="flat" min-height="780">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="me-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="mdiMenu" size="24" />
                                            </v-btn>
                                            <span>{{ tt('Transaction Tags') }}</span>
                                            <v-btn class="ms-3" color="default" variant="outlined"
                                                   :disabled="loading || updating || hasEditingTag" @click="add">{{ tt('Add') }}</v-btn>
                                            <v-btn class="ms-3" color="primary" variant="tonal"
                                                   :disabled="loading || updating || hasEditingTag" @click="saveSortResult"
                                                   v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ms-2" :icon="true" :disabled="loading || updating || hasEditingTag"
                                                   :loading="loading" @click="reload">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="mdiRefresh" size="24" />
                                                <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-spacer/>
                                            <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                                                   :disabled="loading || updating || hasEditingTag" :icon="true">
                                                <v-icon :icon="mdiDotsVertical" />
                                                <v-menu activator="parent">
                                                    <v-list>
                                                        <v-list-item :prepend-icon="mdiPlus" @click="addTagGroup">
                                                            <v-list-item-title>{{ tt('Add Tag Group') }}</v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item :prepend-icon="mdiPencilOutline"
                                                                     @click="renameTagGroup"
                                                                     v-if="activeTagGroupId && activeTagGroupId !== DEFAULT_TAG_GROUP_ID">
                                                            <v-list-item-title>{{ tt('Rename Tag Group') }}</v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item :prepend-icon="mdiDeleteOutline"
                                                                     :disabled="tags && tags.length > 0"
                                                                     @click="removeTagGroup"
                                                                     v-if="activeTagGroupId && activeTagGroupId !== DEFAULT_TAG_GROUP_ID">
                                                            <v-list-item-title>{{ tt('Delete Tag Group') }}</v-list-item-title>
                                                        </v-list-item>
                                                        <v-divider class="my-2" v-if="allTagGroupsWithDefault.length >= 2"/>
                                                        <v-list-item :prepend-icon="mdiSort"
                                                                     :disabled="!allTagGroupsWithDefault || allTagGroupsWithDefault.length < 2"
                                                                     :title="tt('Change Group Display Order')"
                                                                     v-if="allTagGroupsWithDefault.length >= 2"
                                                                     @click="showChangeGroupDisplayOrderDialog"></v-list-item>
                                                        <v-divider class="my-2"/>
                                                        <v-list-item :prepend-icon="mdiSortAlphabeticalAscending"
                                                                     :disabled="!tags || tags.length < 2"
                                                                     :title="tt('Sort by Name (A to Z)')"
                                                                     @click="sortByName(false)"></v-list-item>
                                                        <v-list-item :prepend-icon="mdiSortAlphabeticalDescending"
                                                                     :disabled="!tags || tags.length < 2"
                                                                     :title="tt('Sort by Name (Z to A)')"
                                                                     @click="sortByName(true)"></v-list-item>
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

                                    <v-table class="transaction-tags-table table-striped" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th>
                                                <div class="d-flex align-center">
                                                    <span>{{ tt('Tag Title') }}</span>
                                                    <v-spacer/>
                                                    <span>{{ tt('Operation') }}</span>
                                                </div>
                                            </th>
                                        </tr>
                                        </thead>

                                        <tbody v-if="loading && noAvailableTag && !newTag">
                                        <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5 ]">
                                            <td class="px-0">
                                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                            </td>
                                        </tr>
                                        </tbody>

                                        <tbody v-if="!loading && noAvailableTag && !newTag">
                                        <tr>
                                            <td>{{ tt('No available tag') }}</td>
                                        </tr>
                                        </tbody>

                                        <draggable-list tag="tbody"
                                                        item-key="id"
                                                        handle=".drag-handle"
                                                        ghost-class="dragging-item"
                                                        :class="{ 'has-bottom-border': newTag }"
                                                        :disabled="noAvailableTag"
                                                        v-model="tags"
                                                        @change="onMove">
                                            <template #item="{ element }">
                                                <tr class="transaction-tags-table-row-tag text-sm" v-if="showHidden || !element.hidden">
                                                    <td>
                                                        <div class="d-flex align-center">
                                                            <div class="d-flex align-center" v-if="editingTag.id !== element.id">
                                                                <v-badge class="right-bottom-icon" color="secondary"
                                                                         location="bottom right" offset-x="8" :icon="mdiEyeOffOutline"
                                                                         v-if="element.hidden">
                                                                    <v-icon size="20" start :icon="mdiPound"/>
                                                                </v-badge>
                                                                <v-icon size="20" start :icon="mdiPound" v-else-if="!element.hidden"/>
                                                                <span class="transaction-tag-name">{{ element.name }}</span>
                                                            </div>

                                                            <v-text-field class="w-100 me-2" type="text"
                                                                          density="compact" variant="underlined"
                                                                          :disabled="loading || updating"
                                                                          :placeholder="tt('Tag Title')"
                                                                          v-model="editingTag.name"
                                                                          v-else-if="editingTag.id === element.id"
                                                                          @keyup.enter="save(editingTag)"
                                                            >
                                                                <template #prepend>
                                                                    <v-badge class="right-bottom-icon" color="secondary"
                                                                             location="bottom right" offset-x="8" :icon="mdiEyeOffOutline"
                                                                             v-if="element.hidden">
                                                                        <v-icon size="20" start :icon="mdiPound"/>
                                                                    </v-badge>
                                                                    <v-icon size="20" start :icon="mdiPound" v-else-if="!element.hidden"/>
                                                                </template>
                                                            </v-text-field>

                                                            <v-spacer/>

                                                            <v-btn class="px-2 ms-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="element.hidden ? mdiEyeOutline : mdiEyeOffOutline"
                                                                   :loading="tagHiding[element.id]"
                                                                   :disabled="loading || updating"
                                                                   v-if="editingTag.id !== element.id"
                                                                   @click="hide(element, !element.hidden)">
                                                                <template #loader>
                                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                                </template>
                                                                {{ element.hidden ? tt('Show') : tt('Hide') }}
                                                            </v-btn>
                                                            <v-btn class="px-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="mdiFolderMoveOutline"
                                                                   :loading="tagMoving[element.id]"
                                                                   :disabled="loading || updating || allTagGroupsWithDefault.length < 2"
                                                                   v-if="editingTag.id !== element.id">
                                                                <template #loader>
                                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                                </template>
                                                                {{ tt('Move') }}
                                                                <v-menu activator="parent" max-height="500">
                                                                    <v-list>
                                                                        <v-list-subheader :title="tt('Move to...')"/>
                                                                        <template :key="tagGroup.id" v-for="tagGroup in allTagGroupsWithDefault">
                                                                            <v-list-item class="text-sm" density="compact"
                                                                                         :value="tagGroup.id" v-if="activeTagGroupId !== tagGroup.id">
                                                                                <v-list-item-title class="cursor-pointer"
                                                                                                   @click="moveTagToGroup(element, tagGroup.id)">
                                                                                    <div class="d-flex align-center">
                                                                                        <span class="text-sm ms-3">{{ tagGroup.name }}</span>
                                                                                    </div>
                                                                                </v-list-item-title>
                                                                            </v-list-item>
                                                                        </template>
                                                                    </v-list>
                                                                </v-menu>
                                                            </v-btn>
                                                            <v-btn class="px-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="mdiPencilOutline"
                                                                   :loading="tagUpdating[element.id]"
                                                                   :disabled="loading || updating"
                                                                   v-if="editingTag.id !== element.id"
                                                                   @click="edit(element)">
                                                                <template #loader>
                                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                                </template>
                                                                {{ tt('Edit') }}
                                                            </v-btn>
                                                            <v-btn class="px-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="mdiDeleteOutline"
                                                                   :loading="tagRemoving[element.id]"
                                                                   :disabled="loading || updating"
                                                                   v-if="editingTag.id !== element.id"
                                                                   @click="remove(element)">
                                                                <template #loader>
                                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                                </template>
                                                                {{ tt('Delete') }}
                                                            </v-btn>
                                                            <v-btn class="px-2"
                                                                   density="comfortable" variant="text"
                                                                   :prepend-icon="mdiCheck"
                                                                   :loading="tagUpdating[element.id]"
                                                                   :disabled="loading || updating || !isTagModified(element)"
                                                                   v-if="editingTag.id === element.id" @click="save(editingTag)">
                                                                <template #loader>
                                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                                </template>
                                                                {{ tt('Save') }}
                                                            </v-btn>
                                                            <v-btn class="px-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :prepend-icon="mdiClose"
                                                                   :disabled="loading || updating"
                                                                   v-if="editingTag.id === element.id" @click="cancelSave(editingTag)">
                                                                {{ tt('Cancel') }}
                                                            </v-btn>
                                                            <span class="ms-2">
                                                                <v-icon :class="!loading && !updating && !hasEditingTag && availableTagCount > 1 ? 'drag-handle' : 'disabled'"
                                                                        :icon="mdiDrag"/>
                                                                <v-tooltip activator="parent" v-if="!loading && !updating && !hasEditingTag && availableTagCount > 1">{{ tt('Drag to Reorder') }}</v-tooltip>
                                                            </span>
                                                        </div>
                                                    </td>
                                                </tr>
                                            </template>
                                        </draggable-list>

                                        <tbody v-if="newTag">
                                        <tr class="text-sm" :class="{ 'even-row': (availableTagCount & 1) === 1}">
                                            <td>
                                                <div class="d-flex align-center">
                                                    <v-text-field class="w-100 me-2" type="text" color="primary"
                                                                  density="compact" variant="underlined"
                                                                  :disabled="loading || updating" :placeholder="tt('Tag Title')"
                                                                  v-model="newTag.name" @keyup.enter="save(newTag)">
                                                        <template #prepend>
                                                            <v-icon size="20" start :icon="mdiPound"/>
                                                        </template>
                                                    </v-text-field>

                                                    <v-spacer/>

                                                    <v-btn class="px-2" density="comfortable" variant="text"
                                                           :prepend-icon="mdiCheck"
                                                           :loading="tagUpdating['']"
                                                           :disabled="loading || updating || !isTagModified(newTag)"
                                                           @click="save(newTag)">
                                                        <template #loader>
                                                            <v-progress-circular indeterminate size="20" width="2"/>
                                                        </template>
                                                        {{ tt('Save') }}
                                                    </v-btn>
                                                    <v-btn class="px-2" color="default"
                                                           density="comfortable" variant="text"
                                                           :prepend-icon="mdiClose"
                                                           :disabled="loading || updating"
                                                           @click="cancelSave(newTag)">
                                                        {{ tt('Cancel') }}
                                                    </v-btn>
                                                    <span class="ms-2">
                                                        <v-icon class="disabled" :icon="mdiDrag"/>
                                                    </span>
                                                </div>
                                            </td>
                                        </tr>
                                        </tbody>
                                    </v-table>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <tag-group-change-display-order-dialog ref="tagGroupChangeDisplayOrderDialog" />

    <rename-dialog ref="renameDialog"
                   :default-title="tt('Rename Tag Group')"
                   :label="tt('Tag Group Name')" :placeholder="tt('Tag Group Name')" />
    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import TagGroupChangeDisplayOrderDialog from './dialog/TagGroupChangeDisplayOrderDialog.vue';
import RenameDialog from '@/components/desktop/RenameDialog.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';
import { useDisplay } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useTagListPageBase } from '@/views/base/tags/TagListPageBase.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { DEFAULT_TAG_GROUP_ID } from '@/consts/tag.ts';

import { TransactionTagGroup } from '@/models/transaction_tag_group.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';

import { getAvailableTagCount } from '@/lib/tag.ts';

import {
    mdiRefresh,
    mdiMenuUp,
    mdiMenuDown,
    mdiPencilOutline,
    mdiCheck,
    mdiClose,
    mdiSortAlphabeticalAscending,
    mdiSortAlphabeticalDescending,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiSort,
    mdiMenu,
    mdiPlus,
    mdiFolderMoveOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical,
    mdiPound
} from '@mdi/js';

type TagGroupChangeDisplayOrderDialogType = InstanceType<typeof TagGroupChangeDisplayOrderDialog>;
type RenameDialogType = InstanceType<typeof RenameDialog>;
type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const display = useDisplay();

const { tt, formatNumberToLocalizedNumerals } = useI18n();

const {
    activeTagGroupId,
    newTag,
    editingTag,
    loading,
    showHidden,
    displayOrderModified,
    allTagGroupsWithDefault,
    tags,
    noAvailableTag,
    hasEditingTag,
    isTagModified,
    switchTagGroup,
    add,
    edit
} = useTagListPageBase();

const transactionTagsStore = useTransactionTagsStore();

const tagGroupChangeDisplayOrderDialog = useTemplateRef<TagGroupChangeDisplayOrderDialogType>('tagGroupChangeDisplayOrderDialog');
const renameDialog = useTemplateRef<RenameDialogType>('renameDialog');
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const updating = ref<boolean>(false);
const activeTab = ref<string>('tagListPage');
const alwaysShowNav = ref<boolean>(display.mdAndUp.value);
const showNav = ref<boolean>(display.mdAndUp.value);
const tagUpdating = ref<Record<string, boolean>>({});
const tagHiding = ref<Record<string, boolean>>({});
const tagMoving = ref<Record<string, boolean>>({});
const tagRemoving = ref<Record<string, boolean>>({});

const totalAvailableTagsCount = computed<number>(() => transactionTagsStore.allAvailableTagsCount);
const displayTotalAvailableTagsCount = computed<string>(() => formatNumberToLocalizedNumerals(transactionTagsStore.allAvailableTagsCount));
const availableTagCount = computed<number>(() => getAvailableTagCount(tags.value, showHidden.value));

function reload(): void {
    if (hasEditingTag.value) {
        return;
    }

    loading.value = true;

    transactionTagsStore.loadAllTags({
        force: true
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        snackbar.value?.showMessage('Tag list has been updated');
    }).catch(error => {
        loading.value = false;

        if (error && error.isUpToDate) {
            displayOrderModified.value = false;
        }

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function addTagGroup(): void {
    renameDialog.value?.open('', tt('New Tag Group Name')).then((newName: string) => {
        updating.value = true;

        transactionTagsStore.saveTagGroup({
            tagGroup: TransactionTagGroup.createNewTagGroup(newName)
        }).then(tagGroup => {
            updating.value = false;
            activeTagGroupId.value = tagGroup.id;
        }).catch(error => {
            updating.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function renameTagGroup(): void {
    const tagGroup = transactionTagsStore.allTransactionTagGroupsMap[activeTagGroupId.value];

    if (!tagGroup) {
        snackbar.value?.showMessage('Unable to rename this tag group');
        return;
    }

    renameDialog.value?.open(tagGroup.name || '').then((newName: string) => {
        updating.value = true;

        const newTagGroup = tagGroup.clone();
        newTagGroup.name = newName;

        transactionTagsStore.saveTagGroup({
            tagGroup: newTagGroup
        }).then(() => {
            updating.value = false;
        }).catch(error => {
            updating.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function showChangeGroupDisplayOrderDialog(): void {
    tagGroupChangeDisplayOrderDialog.value?.open().then(() => {
        if (transactionTagsStore.transactionTagGroupListStateInvalid) {
            loading.value = true;

            transactionTagsStore.loadAllTagGroups({
                force: false
            }).then(() => {
                loading.value = false;
            }).catch(() => {
                loading.value = false;
            });
        }
    });
}

function removeTagGroup(): void {
    const tagGroup = transactionTagsStore.allTransactionTagGroupsMap[activeTagGroupId.value];

    if (!tagGroup) {
        snackbar.value?.showMessage('Unable to delete this tag group');
        return;
    }

    const currentTagGroupIndex = allTagGroupsWithDefault.value.findIndex(group => group.id === tagGroup.id);

    confirmDialog.value?.open('Are you sure you want to delete this tag group?').then(() => {
        updating.value = true;

        transactionTagsStore.deleteTagGroup({
            tagGroup: tagGroup
        }).then(() => {
            updating.value = false;

            if (allTagGroupsWithDefault.value[currentTagGroupIndex]) {
                const newActiveTagGroup = allTagGroupsWithDefault.value[currentTagGroupIndex];
                activeTagGroupId.value = newActiveTagGroup ? newActiveTagGroup.id : DEFAULT_TAG_GROUP_ID;
            } else if (allTagGroupsWithDefault.value[currentTagGroupIndex - 1]) {
                const newActiveTagGroup = allTagGroupsWithDefault.value[currentTagGroupIndex - 1];
                activeTagGroupId.value = newActiveTagGroup ? newActiveTagGroup.id : DEFAULT_TAG_GROUP_ID;
            } else {
                activeTagGroupId.value = DEFAULT_TAG_GROUP_ID;
            }
        }).catch(error => {
            updating.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function moveTagToGroup(tag: TransactionTag, targetTagGroupId: string): void {
    updating.value = true;
    tagMoving.value[tag.id] = true;

    const newTag = tag.clone();
    newTag.groupId = targetTagGroupId;

    transactionTagsStore.saveTag({
        tag: newTag
    }).then(() => {
        updating.value = false;
        tagMoving.value[tag.id] = false;
    }).catch(error => {
        updating.value = false;
        tagMoving.value[tag.id] = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function save(tag: TransactionTag): void {
    updating.value = true;
    tagUpdating.value[tag.id || ''] = true;

    transactionTagsStore.saveTag({
        tag: tag
    }).then(() => {
        updating.value = false;
        tagUpdating.value[tag.id || ''] = false;

        if (tag.id) {
            editingTag.value.id = '';
            editingTag.value.name = '';
        } else {
            newTag.value = null;
        }
    }).catch(error => {
        updating.value = false;
        tagUpdating.value[tag.id || ''] = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function cancelSave(tag: TransactionTag): void {
    if (tag.id) {
        editingTag.value.id = '';
        editingTag.value.name = '';
    } else {
        newTag.value = null;
    }
}

function sortByName(desc: boolean): void {
    const changed = transactionTagsStore.sortTagDisplayOrderByTagName(activeTagGroupId.value, desc);

    if (changed) {
        displayOrderModified.value = true;
    }
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        return;
    }

    loading.value = true;

    transactionTagsStore.updateTagDisplayOrders(activeTagGroupId.value).then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function hide(tag: TransactionTag, hidden: boolean): void {
    updating.value = true;
    tagHiding.value[tag.id] = true;

    transactionTagsStore.hideTag({
        tag: tag,
        hidden: hidden
    }).then(() => {
        updating.value = false;
        tagHiding.value[tag.id] = false;
    }).catch(error => {
        updating.value = false;
        tagHiding.value[tag.id] = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function remove(tag: TransactionTag): void {
    confirmDialog.value?.open('Are you sure you want to delete this tag?').then(() => {
        updating.value = true;
        tagRemoving.value[tag.id] = true;

        transactionTagsStore.deleteTag({
            tag: tag
        }).then(() => {
            updating.value = false;
            tagRemoving.value[tag.id] = false;
        }).catch(error => {
            updating.value = false;
            tagRemoving.value[tag.id] = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function onMove(event: { moved: { element: { id: string }; oldIndex: number; newIndex: number } }): void {
    if (!event || !event.moved) {
        return;
    }

    const moveEvent = event.moved;

    if (!moveEvent.element || !moveEvent.element.id) {
        snackbar.value?.showMessage('Unable to move tag');
        return;
    }

    transactionTagsStore.changeTagDisplayOrder({
        tagId: moveEvent.element.id,
        from: moveEvent.oldIndex,
        to: moveEvent.newIndex
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        snackbar.value?.showError(error);
    });
}

transactionTagsStore.loadAllTags({
    force: false
}).then(() => {
    loading.value = false;
}).catch(error => {
    loading.value = false;

    if (!error.processed) {
        snackbar.value?.showError(error);
    }
});

watch(() => display.mdAndUp.value, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});
</script>

<style>
.transaction-tags-statistic-item-value {
    font-size: 1rem;
}

.transaction-tags-table tr.transaction-tags-table-row-tag .hover-display {
    display: none;
}

.transaction-tags-table tr.transaction-tags-table-row-tag:hover .hover-display {
    display: inline-grid;
}

.transaction-tags-table tr:not(:last-child) > td > div {
    padding-bottom: 1px;
}

.transaction-tags-table .has-bottom-border tr:last-child > td > div {
    padding-bottom: 1px;
}

.transaction-tags-table tr.transaction-tags-table-row-tag .right-bottom-icon .v-badge__badge {
    padding-bottom: 1px;
}

.transaction-tags-table .v-text-field .v-input__prepend {
    margin-inline-end: 0;
    color: rgba(var(--v-theme-on-surface));
}

.transaction-tags-table .v-text-field .v-input__prepend .v-badge > .v-badge__wrapper > .v-icon {
    opacity: var(--v-medium-emphasis-opacity);
}

.transaction-tags-table .v-text-field.v-input--plain-underlined .v-input__prepend {
    padding-top: 10px;
}

.transaction-tags-table .v-text-field .v-field__input {
    font-size: 0.875rem;
    padding-top: 0;
    color: rgba(var(--v-theme-on-surface));
}

.transaction-tags-table .transaction-tag-name {
    font-size: 0.875rem;
}

.transaction-tags-table tr .v-text-field .v-field__input {
    padding-bottom: 1px;
}
</style>
