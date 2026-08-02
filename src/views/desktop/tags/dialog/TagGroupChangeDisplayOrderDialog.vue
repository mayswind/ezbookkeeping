<template>
    <v-dialog width="800" :persistent="displayOrderModified" v-model="showState">
        <one-column-dialog-layout content-class="pa-0" :disabled="loading || updating"
                                  :title="tt('Change Group Display Order')" :cancel-button-title="tt('Close')"
                                  @cancel="close">
            <template #after-title>
                <v-btn density="compact" color="default" variant="text" size="22"
                       class="ms-2" :icon="true" :disabled="loading || updating"
                       :loading="loading" @click="reload">
                    <template #loader>
                        <v-progress-circular indeterminate size="20"/>
                    </template>
                    <v-icon :icon="mdiRefresh" size="22" />
                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                </v-btn>
                <v-btn density="compact" color="primary" variant="text" class="ms-1" :icon="true"
                       :disabled="loading || updating || !displayOrderModified" @click="saveDisplayOrder">
                    <v-icon :icon="mdiCheck" size="22" />
                    <v-tooltip activator="parent">{{ tt('Save Display Order') }}</v-tooltip>
                </v-btn>
            </template>

            <template #content>
                <v-table hover density="comfortable" class="w-100 table-striped">
                    <tbody v-if="loading && (!allTagGroups || allTagGroups.length < 1)">
                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6 ]">
                        <td class="px-0">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>
                    </tbody>

                    <tbody v-if="!loading && (!allTagGroups || allTagGroups.length < 1)">
                    <tr>
                        <td>{{ tt('No available tag group') }}</td>
                    </tr>
                    </tbody>

                    <draggable-list tag="tbody"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    v-model="allTagGroups"
                                    @change="onMove">
                        <template #item="{ element }">
                            <tr>
                                <td>
                                    <div class="d-flex align-center">
                                        <div class="d-flex align-center">
                                            <span>{{ element.name }}</span>
                                        </div>

                                        <v-spacer/>

                                        <span class="ms-2">
                                            <v-icon :class="!loading && !updating && allTagGroups && allTagGroups.length > 0 ? 'drag-handle' : 'disabled'"
                                                    :icon="mdiDrag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && allTagGroups && allTagGroups.length > 0">{{ tt('Drag to Reorder') }}</v-tooltip>
                                        </span>
                                    </div>
                                </td>
                            </tr>
                        </template>
                    </draggable-list>
                </v-table>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { type TransactionTagGroup } from '@/models/transaction_tag_group.ts';

import {
    mdiRefresh,
    mdiCheck,
    mdiDrag
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const transactionTagsStore = useTransactionTagsStore();

let resolveFunc: (() => void) | null = null;

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const displayOrderModified = ref<boolean>(false);

const allTagGroups = computed<TransactionTagGroup[]>(() => transactionTagsStore.allTransactionTagGroups);

function open(): Promise<void> {
    showState.value = true;
    loading.value = true;

    transactionTagsStore.loadAllTagGroups({
        force: false
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });

    return new Promise<void>((resolve) => {
        resolveFunc = resolve;
    });
}

function reload(): void {
    loading.value = true;

    transactionTagsStore.loadAllTagGroups({
        force: true
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        snackbar.value?.showMessage('Tag group list has been updated');
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

function saveDisplayOrder(): void {
    if (!displayOrderModified.value) {
        return;
    }

    loading.value = true;

    transactionTagsStore.updateTagGroupDisplayOrders().then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function close(): void {
    if (loading.value || updating.value) {
        return;
    }

    resolveFunc?.();
    showState.value = false;
}

function onMove(event: { moved: { element: { id: string }; oldIndex: number; newIndex: number } }): void {
    if (!event || !event.moved) {
        return;
    }

    const moveEvent = event.moved;

    if (!moveEvent.element || !moveEvent.element.id) {
        snackbar.value?.showMessage('Unable to move tag group');
        return;
    }

    transactionTagsStore.changeTagGroupDisplayOrder({
        tagGroupId: moveEvent.element.id,
        from: moveEvent.oldIndex,
        to: moveEvent.newIndex
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        snackbar.value?.showError(error);
    });
}

defineExpose({
    open
});
</script>
