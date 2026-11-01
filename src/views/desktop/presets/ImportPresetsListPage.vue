<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ tt('Import Presets') }}</span>
                        <v-btn class="ms-3" color="default" variant="outlined"
                               :disabled="loading || updating" @click="add">{{ tt('Add') }}</v-btn>
                        <v-btn class="ms-3" color="primary" variant="tonal"
                               :disabled="loading || updating" @click="saveSortResult"
                               v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" class="ms-2"
                               :aria-label="tt('Refresh')" :icon="true" :disabled="loading || updating"
                               :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-table class="import-presets-table table-striped" :hover="!loading">
                    <thead>
                    <tr>
                        <th>
                            <div class="d-flex align-center">
                                <span>{{ tt('Preset Name') }}</span>
                                <v-spacer/>
                                <span>{{ tt('Operation') }}</span>
                            </div>
                        </th>
                    </tr>
                    </thead>

                    <tbody v-if="loading && !presets.length">
                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                        <td class="px-0">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>
                    </tbody>

                    <tbody v-if="!loading && !presets.length">
                    <tr>
                        <td>{{ tt('No available import presets') }}</td>
                    </tr>
                    </tbody>

                    <draggable-list tag="tbody"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    :disabled="!presets.length"
                                    v-model="presets"
                                    @change="onMove">
                        <template #item="{ element }">
                            <tr @mouseenter="hoveredPresetId = element.id" @mouseleave="hoveredPresetId = ''">
                                <td>
                                    <div class="d-flex align-center">
                                        <v-icon size="20" start :icon="mdiFileTableBoxOutline"/>
                                        <span class="import-preset-name">{{ element.name }}</span>
                                        <v-spacer/>

                                        <template v-if="hoveredPresetId === element.id && !loading">
                                            <v-btn class="px-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="mdiPencilOutline"
                                                   :disabled="loading || updating"
                                                   @click="edit(element)">
                                                {{ tt('Edit') }}
                                            </v-btn>
                                            <v-btn class="px-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="mdiDeleteOutline"
                                                   :loading="presetRemoving[element.id]"
                                                   :disabled="loading || updating"
                                                   @click="remove(element)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ tt('Delete') }}
                                            </v-btn>
                                        </template>

                                        <span class="ms-2">
                                            <v-icon :class="!loading && !updating && presets.length > 1 ? 'drag-handle' : 'disabled'"
                                                    :aria-label="tt('Drag to Reorder')" :icon="mdiDrag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && presets.length > 1 && hoveredPresetId === element.id">{{ tt('Drag to Reorder') }}</v-tooltip>
                                        </span>
                                    </div>
                                </td>
                            </tr>
                        </template>
                    </draggable-list>
                </v-table>
            </v-card>
        </v-col>
    </v-row>

    <import-preset-edit-dialog ref="editDialog" />
    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import ImportPresetEditDialog from './dialogs/ImportPresetEditDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useImportPresetsStore } from '@/stores/importPreset.ts';

import { type ImportPresetBasicInfo } from '@/models/import_preset.ts';

import {
    mdiRefresh,
    mdiPencilOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiFileTableBoxOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type EditDialogType = InstanceType<typeof ImportPresetEditDialog>;

const { tt } = useI18n();

const importPresetsStore = useImportPresetsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const hoveredPresetId = ref<string>('');
const presetRemoving = ref<Record<string, boolean>>({});
const displayOrderModified = ref<boolean>(false);

const presets = computed<ImportPresetBasicInfo[]>(() => importPresetsStore.allImportPresetBasicInfos);

function init(): void {
    loading.value = true;

    importPresetsStore.loadAllPresetBasicInfos({
        force: false
    }).then(() => {
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

    importPresetsStore.loadAllPresetBasicInfos({
        force: true
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        snackbar.value?.showMessage('Import preset list has been updated');
    }).catch(error => {
        loading.value = false;

        if (error.isUpToDate) {
            displayOrderModified.value = false;
        }

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function add(): void {
    editDialog.value?.open({}).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function edit(preset: ImportPresetBasicInfo): void {
    editDialog.value?.open({
        id: preset.id
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function remove(preset: ImportPresetBasicInfo): void {
    confirmDialog.value?.open('Are you sure you want to delete this import preset?').then(() => {
        updating.value = true;
        presetRemoving.value[preset.id] = true;

        importPresetsStore.deletePreset({
            preset: preset
        }).then(() => {
            updating.value = false;
            presetRemoving.value[preset.id] = false;
        }).catch(error => {
            updating.value = false;
            presetRemoving.value[preset.id] = false;

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

    importPresetsStore.updatePresetDisplayOrders().then(() => {
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
        snackbar.value?.showMessage('Unable to move import preset');
        return;
    }

    importPresetsStore.changePresetDisplayOrder({
        presetId: moveEvent.element.id,
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
