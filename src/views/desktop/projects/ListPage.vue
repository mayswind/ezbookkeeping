<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ tt('Projects') }}</span>
                        <v-btn class="ms-3" color="default" variant="outlined"
                               :disabled="loading || updating" @click="add">{{ tt('Add') }}</v-btn>
                        <v-btn class="ms-3" color="primary" variant="tonal"
                               :disabled="loading || updating" @click="saveSortResult"
                               v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ms-2" :icon="true" :disabled="loading || updating"
                               :loading="loading" @click="reload(true)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-spacer/>
                        <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                               :disabled="loading || updating" :icon="true">
                            <v-icon :icon="mdiDotsVertical" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :prepend-icon="mdiEyeOutline"
                                                 :title="tt('Show Hidden Projects')"
                                                 v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                    <v-list-item :prepend-icon="mdiEyeOffOutline"
                                                 :title="tt('Hide Hidden Projects')"
                                                 v-if="showHidden" @click="showHidden = false"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </div>
                </template>

                <v-table class="project-table table-striped" :hover="!loading">
                    <thead>
                    <tr>
                        <th>
                            <div class="d-flex align-center">
                                <span>{{ tt('Project Name') }}</span>
                                <v-spacer/>
                                <span>{{ tt('Operation') }}</span>
                            </div>
                        </th>
                    </tr>
                    </thead>

                    <tbody v-if="loading && noAvailableProject">
                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                        <td class="px-0">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>
                    </tbody>

                    <tbody v-if="!loading && noAvailableProject">
                    <tr>
                        <td>{{ tt('No available project') }}</td>
                    </tr>
                    </tbody>

                    <draggable-list tag="tbody"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    :disabled="noAvailableProject"
                                    v-model="projects"
                                    @change="onMove">
                        <template #item="{ element }">
                            <tr class="project-table-row text-sm" v-if="showHidden || !element.hidden">
                                <td>
                                    <div class="d-flex align-center">
                                        <div class="d-flex align-center">
                                            <v-icon size="20" start :icon="mdiFolder" :color="element.color" v-if="!element.hidden"/>
                                            <v-badge class="right-bottom-icon" color="secondary"
                                                     location="bottom right" offset-x="8" :icon="mdiEyeOffOutline"
                                                     v-if="element.hidden">
                                                <v-icon size="20" start :icon="mdiFolder" :color="element.color"/>
                                            </v-badge>
                                            <span class="ms-2">{{ element.name }}</span>
                                            <span class="text-caption ms-2 text-medium-emphasis">{{ element.comment }}</span>
                                        </div>

                                        <v-spacer/>

                                        <v-btn class="px-2 ms-2" color="default"
                                               density="comfortable" variant="text"
                                               :class="{ 'd-none': loading, 'hover-display': !loading }"
                                               :prepend-icon="element.hidden ? mdiEyeOutline : mdiEyeOffOutline"
                                               :loading="projectHiding[element.id]"
                                               :disabled="loading || updating"
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
                                               :disabled="loading || updating"
                                               @click="edit(element)">
                                            {{ tt('Edit') }}
                                        </v-btn>
                                        <v-btn class="px-2" color="default"
                                               density="comfortable" variant="text"
                                               :class="{ 'd-none': loading, 'hover-display': !loading }"
                                               :prepend-icon="mdiDeleteOutline"
                                               :loading="projectRemoving[element.id]"
                                               :disabled="loading || updating"
                                               @click="remove(element)">
                                            <template #loader>
                                                <v-progress-circular indeterminate size="20" width="2"/>
                                            </template>
                                            {{ tt('Delete') }}
                                        </v-btn>
                                        <span class="ms-2">
                                            <v-icon :class="!loading && !updating && availableProjectCount > 1 ? 'drag-handle' : 'disabled'"
                                                    :icon="mdiDrag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && availableProjectCount > 1">{{ tt('Drag to Reorder') }}</v-tooltip>
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

    <edit-dialog ref="editDialog" />
    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import EditDialog from './list/dialogs/EditDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useProjectsStore } from '@/stores/project.ts';
import { Project } from '@/models/project.ts';

import {
    mdiRefresh,
    mdiPencilOutline,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical,
    mdiFolder
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type EditDialogType = InstanceType<typeof EditDialog>;

const { tt } = useI18n();
const projectsStore = useProjectsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const projectHiding = ref<Record<string, boolean>>({});
const projectRemoving = ref<Record<string, boolean>>({});
const displayOrderModified = ref<boolean>(false);
const showHidden = ref<boolean>(false);

const projects = computed<Project[]>({
    get: () => projectsStore.allProjects,
    set: (value) => {
        // draggable-list updates v-model
    }
});

const noAvailableProject = computed<boolean>(() => {
    if (!projects.value || projects.value.length < 1) return true;
    if (showHidden.value) return false;
    return !projects.value.some(p => !p.hidden);
});

const availableProjectCount = computed<number>(() => {
    if (!projects.value) return 0;
    if (showHidden.value) return projects.value.length;
    return projects.value.filter(p => !p.hidden).length;
});

function reload(force: boolean): void {
    loading.value = true;
    projectsStore.loadAllProjects({ force: force }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;
        if (force) snackbar.value?.showMessage('Project list has been updated');
    }).catch(error => {
        loading.value = false;
        if (!error.processed) snackbar.value?.showError(error);
    });
}

function add(): void {
    editDialog.value?.open({}).then(result => {
        if (result && result.message) snackbar.value?.showMessage(result.message);
    }).catch(error => {
        if (error) snackbar.value?.showError(error);
    });
}

function edit(project: Project): void {
    editDialog.value?.open({ id: project.id, currentProject: project }).then(result => {
        if (result && result.message) snackbar.value?.showMessage(result.message);
    }).catch(error => {
        if (error) snackbar.value?.showError(error);
    });
}

function hide(project: Project, hidden: boolean): void {
    updating.value = true;
    projectHiding.value[project.id] = true;

    projectsStore.hideProject({
        project: { id: project.id, hidden: hidden },
        hidden: hidden
    }).then(() => {
        updating.value = false;
        projectHiding.value[project.id] = false;
    }).catch(error => {
        updating.value = false;
        projectHiding.value[project.id] = false;
        if (!error.processed) snackbar.value?.showError(error);
    });
}

function remove(project: Project): void {
    confirmDialog.value?.open('Are you sure you want to delete this project?').then(() => {
        updating.value = true;
        projectRemoving.value[project.id] = true;

        projectsStore.deleteProject({
            project: { id: project.id }
        }).then(() => {
            updating.value = false;
            projectRemoving.value[project.id] = false;
        }).catch(error => {
            updating.value = false;
            projectRemoving.value[project.id] = false;
            if (!error.processed) snackbar.value?.showError(error);
        });
    });
}

function saveSortResult(): void {
    if (!displayOrderModified.value) return;

    loading.value = true;
    const newDisplayOrders = projects.value.map((p, index) => ({ id: p.id, displayOrder: index + 1 }));

    projectsStore.moveProject({
        project: { newDisplayOrders }
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;
        if (!error.processed) snackbar.value?.showError(error);
    });
}

function onMove(event: { moved: { element: { id: string }, oldIndex: number, newIndex: number } }): void {
    if (event && event.moved) {
        displayOrderModified.value = true;
    }
}

reload(false);
</script>

<style>
.project-table tr.project-table-row .hover-display {
    display: none;
}

.project-table tr.project-table-row:hover .hover-display {
    display: inline-grid;
}

.project-table tr:not(:last-child) > td > div {
    padding-bottom: 1px;
}
</style>
