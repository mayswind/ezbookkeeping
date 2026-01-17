<template>
    <v-dialog width="800" :persistent="displayOrderModified" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex align-center">
                        <h4 class="text-h4">{{ tt('Change Explorer Display Order') }}</h4>
                        <v-btn class="ms-3" color="primary" variant="tonal"
                               :disabled="loading || updating" @click="saveDisplayOrder"
                               v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" :disabled="loading || updating"
                               :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                    <v-spacer/>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :disabled="loading || updating" :icon="true">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiEyeOutline"
                                             :title="tt('Show Hidden Explorers')"
                                             v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                <v-list-item :prepend-icon="mdiEyeOffOutline"
                                             :title="tt('Hide Hidden Explorers')"
                                             v-if="showHidden" @click="showHidden = false"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>

            <v-card-text class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto">
                <v-table hover density="comfortable" class="explorers-table w-100 table-striped">
                    <tbody v-if="loading && noAvailableExplorer">
                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6 ]">
                        <td class="px-0">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>
                    </tbody>

                    <tbody v-if="!loading && noAvailableExplorer">
                    <tr>
                        <td>{{ tt('No available explorer') }}</td>
                    </tr>
                    </tbody>

                    <draggable-list tag="tbody"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    v-model="allExplorers"
                                    @change="onMove">
                        <template #item="{ element }">
                            <tr class="explorers-table-row text-sm" v-if="showHidden || !element.hidden">
                                <td>
                                    <div class="d-flex align-center">
                                        <div class="d-flex align-center">
                                            <span>{{ element.name }}</span>
                                        </div>

                                        <v-spacer/>

                                        <v-btn class="px-2 ms-2" color="default"
                                               density="compact" variant="text"
                                               :class="{ 'd-none': loading, 'hover-display': !loading }"
                                               :prepend-icon="element.hidden ? mdiEyeOutline : mdiEyeOffOutline"
                                               :loading="explorerHiding[element.id]"
                                               :disabled="loading || updating"
                                               @click="hide(element, !element.hidden)">
                                            <template #loader>
                                                <v-progress-circular indeterminate size="20" width="2"/>
                                            </template>
                                            {{ element.hidden ? tt('Show') : tt('Hide') }}
                                        </v-btn>
                                        <span class="ms-2">
                                            <v-icon :class="!loading && !updating && !noAvailableExplorer ? 'drag-handle' : 'disabled'"
                                                    :icon="mdiDrag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && !noAvailableExplorer">{{ tt('Drag to Reorder') }}</v-tooltip>
                                        </span>
                                    </div>
                                </td>
                            </tr>
                        </template>
                    </draggable-list>
                </v-table>
            </v-card-text>

            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn color="secondary" variant="tonal"
                           :disabled="loading || updating" @click="close">{{ tt('Close') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useExplorersStore } from '@/stores/explorer.ts';

import { type InsightsExplorerBasicInfo } from '@/models/explorer.ts';

import {
    mdiDotsVertical,
    mdiRefresh,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiDrag
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const explorersStore = useExplorersStore();

let resolveFunc: (() => void) | null = null;

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const explorerHiding = ref<Record<string, boolean>>({});
const displayOrderModified = ref<boolean>(false);
const showHidden = ref<boolean>(false);

const allExplorers = computed<InsightsExplorerBasicInfo[]>(() => explorersStore.allInsightsExplorerBasicInfos);

const noAvailableExplorer = computed<boolean>(() => {
    for (const explorer of allExplorers.value) {
        if (showHidden.value || !explorer.hidden) {
            return false;
        }
    }

    return true;
});

function open(): Promise<void> {
    showHidden.value = false;
    showState.value = true;
    loading.value = true;

    explorersStore.loadAllInsightsExplorerBasicInfos({
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

    explorersStore.loadAllInsightsExplorerBasicInfos({
        force: true
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        snackbar.value?.showMessage('Explorer list has been updated');
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

function hide(explorer: InsightsExplorerBasicInfo, hidden: boolean): void {
    updating.value = true;
    explorerHiding.value[explorer.id] = true;

    explorersStore.hideInsightsExplorer({
        explorer: explorer,
        hidden: hidden
    }).then(() => {
        updating.value = false;
        explorerHiding.value[explorer.id] = false;
    }).catch(error => {
        updating.value = false;
        explorerHiding.value[explorer.id] = false;

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

    explorersStore.updateInsightsExplorerDisplayOrders().then(() => {
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
        snackbar.value?.showMessage('Unable to move explorer');
        return;
    }

    explorersStore.changeInsightsExplorerDisplayOrder({
        explorerId: moveEvent.element.id,
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

<style>
.explorers-table tr.explorers-table-row .hover-display {
    display: none;
}

.explorers-table tr.explorers-table-row:hover .hover-display {
    display: inline-grid;
}
</style>
