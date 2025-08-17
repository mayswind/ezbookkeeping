<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
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
                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
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
        </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { TransactionTag } from '@/models/transaction_tag.ts';

import {
    isNoAvailableTag,
    getAvailableTagCount
} from '@/lib/tag.ts';

import {
    mdiRefresh,
    mdiPencilOutline,
    mdiCheck,
    mdiClose,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical,
    mdiPound
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const transactionTagsStore = useTransactionTagsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const newTag = ref<TransactionTag | null>(null);
const editingTag = ref<TransactionTag>(TransactionTag.createNewTag());
const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const tagUpdating = ref<Record<string, boolean>>({});
const tagHiding = ref<Record<string, boolean>>({});
const tagRemoving = ref<Record<string, boolean>>({});
const displayOrderModified = ref<boolean>(false);
const showHidden = ref<boolean>(false);

const tags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);
const noAvailableTag = computed<boolean>(() => isNoAvailableTag(tags.value, showHidden.value));
const availableTagCount = computed<number>(() => getAvailableTagCount(tags.value, showHidden.value));
const hasEditingTag = computed<boolean>(() => !!(newTag.value || (editingTag.value.id && editingTag.value.id !== '')));

function isTagModified(tag: TransactionTag): boolean {
    if (tag.id) {
        return editingTag.value.name !== '' && editingTag.value.name !== tag.name;
    } else {
        return tag.name !== '';
    }
}

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

function add(): void {
    newTag.value = TransactionTag.createNewTag();
}

function edit(tag: TransactionTag): void {
    editingTag.value.id = tag.id;
    editingTag.value.name = tag.name;
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

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        return;
    }

    loading.value = true;

    transactionTagsStore.updateTagDisplayOrders().then(() => {
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
</script>

<style>
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
