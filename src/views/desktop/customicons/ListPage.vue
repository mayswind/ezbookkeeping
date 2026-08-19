<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ tt('Custom Icons') }}</span>
                        <v-btn class="ms-3" color="default" variant="outlined"
                               :disabled="loading || updating" @click="add">{{ tt('Add') }}</v-btn>
                        <v-btn class="ms-3" color="primary" variant="tonal"
                               :disabled="loading || updating" @click="saveSortResult"
                               v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text"
                               class="ms-2" :icon="true" :disabled="loading || updating"
                               :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-card-text class="pt-2" v-if="allCustomIcons.length < 1">
                    <div class="custom-icon-grid" v-if="loading">
                        <v-card border class="custom-icon-card card-title-with-bg" variant="flat"
                                :key="idx" v-for="idx in [ 1, 2, 3 ]">
                            <template #title>
                                <div class="d-flex align-center">
                                    <v-btn density="compact" color="default" variant="text"
                                           :icon="true" :disabled="true">
                                        <v-icon :icon="mdiDeleteOutline" size="22" />
                                    </v-btn>
                                    <v-spacer/>
                                    <span class="align-self-center">
                                        <v-icon class="disabled" :icon="mdiDrag"/>
                                    </span>
                                </div>
                            </template>

                            <v-card-text class="d-flex pa-1">
                                <v-skeleton-loader class="w-100 justify-center skeleton-no-margin" type="avatar" :loading="true"></v-skeleton-loader>
                            </v-card-text>
                        </v-card>
                    </div>

                    <div class="text-body-medium" v-else-if="!loading">{{ tt('No available custom icons') }}</div>
                </v-card-text>

                <v-card-text class="pt-2" v-else-if="allCustomIcons.length > 0">
                    <draggable-list class="custom-icon-grid"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    v-model="allCustomIcons"
                                    @change="onMove">
                        <template #item="{ element }">
                            <v-card border class="custom-icon-card card-title-with-bg" variant="flat"
                                    @mouseenter="hoveredIconId = element.id" @mouseleave="hoveredIconId = ''">
                                <template #title>
                                    <div class="d-flex align-center">
                                        <v-btn density="compact" color="default" variant="text"
                                               :icon="true" :disabled="updating"
                                               @click="deleteCustomIcon(element)">
                                            <v-icon :icon="mdiDeleteOutline" size="22" v-if="deletingCustomIconId !== element.id" />
                                            <v-progress-circular indeterminate size="20" v-else-if="deletingCustomIconId === element.id" />
                                            <v-tooltip activator="parent" v-if="hoveredIconId === element.id">{{ tt('Delete') }}</v-tooltip>
                                        </v-btn>
                                        <v-spacer/>
                                        <span class="align-self-center">
                                            <v-icon :class="!loading && !updating && allCustomIcons.length > 1 ? 'drag-handle' : 'disabled'"
                                                    :icon="mdiDrag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && allCustomIcons.length > 1 && hoveredIconId === element.id">{{ tt('Drag to Reorder') }}</v-tooltip>
                                        </span>
                                    </div>
                                </template>

                                <v-card-text class="pa-1">
                                    <v-img class="custom-icon-thumbnail" aspect-ratio="1" cover :src="getCustomIconUrl(element)" />
                                </v-card-text>
                            </v-card>
                        </template>
                    </draggable-list>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <upload-dialog ref="uploadDialog" />

    <confirm-dialog ref="confirmDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import UploadDialog from '@/views/desktop/customicons/dialogs/UploadDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserCustomIconsStore } from '@/stores/userCustomIcon.ts';

import type { UserCustomIconInfoResponse } from '@/models/user_custom_icon.ts';

import {
    mdiRefresh,
    mdiDrag,
    mdiDeleteOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type UploadDialogType = InstanceType<typeof UploadDialog>;

const { tt } = useI18n();

const customIconsStore = useUserCustomIconsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const uploadDialog = useTemplateRef<UploadDialogType>('uploadDialog');

const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const hoveredIconId = ref<string>('');
const displayOrderModified = ref<boolean>(false);
const deletingCustomIconId = ref<string>('');

const allCustomIcons = computed<UserCustomIconInfoResponse[]>(() => customIconsStore.allCustomIcons);

function getCustomIconUrl(customIcon: UserCustomIconInfoResponse): string {
    return customIconsStore.getUserCustomIconUrlWithToken(customIcon);
}

function init(): void {
    loading.value = true;

    customIconsStore.loadAllCustomIcons({ force: false }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function reload(): void {
    loading.value = true;

    customIconsStore.loadAllCustomIcons({ force: true }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function add(): void {
    uploadDialog.value?.open().then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function deleteCustomIcon(customIcon: UserCustomIconInfoResponse): void {
    confirmDialog.value?.open('Are you sure you want to delete this custom icon?').then(() => {
        updating.value = true;
        deletingCustomIconId.value = customIcon.id;

        customIconsStore.deleteCustomIcon({
            customIcon: customIcon
        }).then(() => {
            updating.value = false;
            deletingCustomIconId.value = '';
        }).catch(error => {
            updating.value = false;
            deletingCustomIconId.value = '';

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        return;
    }

    loading.value = true;

    customIconsStore.updateCustomIconDisplayOrders().then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function onMove(event: { moved: { element: { id: string }, oldIndex: number, newIndex: number } }): void {
    if (!event || !event.moved) {
        return;
    }

    const moveEvent = event.moved;

    if (!moveEvent.element || !moveEvent.element.id) {
        snackbar.value?.showMessage('Unable to move custom icon');
        return;
    }

    customIconsStore.changeCustomIconDisplayOrder({
        iconId: moveEvent.element.id,
        from: moveEvent.oldIndex,
        to: moveEvent.newIndex
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        snackbar.value?.showError(error);
    });
}

init();
</script>

<style>
.custom-icon-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
    gap: 8px;

    > .custom-icon-card {
        > .v-card-item {
            padding: 4px 2px;
        }

        > .v-card-text {
            min-height: 78px;
        }
    }
}

.custom-icon-thumbnail {
    border-radius: var(--ebk-radius-md);
    background: rgba(var(--v-theme-on-surface), 0.04);
}
</style>
