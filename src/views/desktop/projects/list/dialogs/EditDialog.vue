<template>
    <v-dialog width="800" :persistent="true" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ tt(title) }}</h4>
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="loading"></v-progress-circular>
                </div>
            </template>
            <v-card-text class="pt-0">
                <v-form class="mt-md-6">
                    <v-row>
                        <v-col cols="12" md="12">
                            <v-text-field
                                type="text"
                                persistent-placeholder
                                :disabled="loading || submitting"
                                :label="tt('Project Name')"
                                :placeholder="tt('Project Name')"
                                v-model="project.name"
                            />
                        </v-col>
                        <v-col cols="12" md="12">
                            <color-select :all-color-infos="ALL_CATEGORY_COLORS"
                                         :label="tt('Project Color')"
                                         :disabled="loading || submitting"
                                         v-model="project.color" />
                        </v-col>
                        <v-col cols="12" md="12">
                            <v-textarea
                                type="text"
                                persistent-placeholder
                                rows="3"
                                :disabled="loading || submitting"
                                :label="tt('Description')"
                                :placeholder="tt('Your project description (optional)')"
                                v-model="project.comment"
                            />
                        </v-col>
                        <v-col class="py-0" cols="12" md="12" v-if="editProjectId">
                            <v-switch :disabled="loading || submitting"
                                      :label="tt('Visible')" v-model="visible"/>
                        </v-col>
                    </v-row>
                </v-form>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-tooltip :disabled="!inputIsEmpty" :text="inputEmptyProblemMessage ? tt(inputEmptyProblemMessage) : ''">
                        <template v-slot:activator="{ props }">
                            <div v-bind="props" class="d-inline-block">
                                <v-btn :disabled="inputIsEmpty || loading || submitting" @click="save">
                                    {{ tt(saveButtonTitle) }}
                                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                                </v-btn>
                            </div>
                        </template>
                    </v-tooltip>
                    <v-btn color="secondary" variant="tonal"
                           :disabled="loading || submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
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
import { useProjectsStore } from '@/stores/project.ts';
import { Project } from '@/models/project.ts';
import { ALL_CATEGORY_COLORS } from '@/consts/color.ts';
import { generateRandomUUID } from '@/lib/misc.ts';

interface ProjectEditResponse {
    message: string;
}

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();
const projectsStore = useProjectsStore();
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const submitting = ref<boolean>(false);
const editProjectId = ref<string | null>(null);
const project = ref<Project>(new Project('', '', '', '', 0, false));
const clientSessionId = ref<string>('');

let resolveFunc: ((value: ProjectEditResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const title = computed<string>(() => editProjectId.value ? 'Edit Project' : 'Add Project');
const saveButtonTitle = computed<string>(() => editProjectId.value ? 'Save' : 'Add');
const inputIsEmpty = computed<boolean>(() => !project.value.name);
const inputEmptyProblemMessage = computed<string | null>(() => !project.value.name ? 'Project name cannot be blank' : null);

const visible = computed<boolean>({
    get: () => !project.value.hidden,
    set: (val) => project.value.hidden = !val
});

function open(options: { id?: string; currentProject?: Project }): Promise<ProjectEditResponse> {
    showState.value = true;
    loading.value = true;
    submitting.value = false;

    // Reset
    project.value = new Project('', '', '', '', 0, false);
    project.value.color = ALL_CATEGORY_COLORS[0] || ''; // Default color

    if (options.id) {
        editProjectId.value = options.id;
        if (options.currentProject) {
            project.value = new Project(
                options.currentProject.id,
                options.currentProject.name,
                options.currentProject.color,
                options.currentProject.comment,
                options.currentProject.displayOrder,
                options.currentProject.hidden
            );
            loading.value = false;
        } else {
             const p = projectsStore.getProject(options.id);
             if (p) {
                 project.value = new Project(p.id, p.name, p.color, p.comment, p.displayOrder, p.hidden);
                 loading.value = false;
             } else {
                 // Fallback or error
                 loading.value = false;
             }
        }
    } else {
        editProjectId.value = null;
        clientSessionId.value = generateRandomUUID();
        loading.value = false;
    }

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function save(): void {
    if (inputIsEmpty.value) return;

    submitting.value = true;

    const req = editProjectId.value ? {
        id: editProjectId.value,
        name: project.value.name,
        color: project.value.color,
        comment: project.value.comment,
        hidden: project.value.hidden
    } : {
        name: project.value.name,
        color: project.value.color,
        comment: project.value.comment,
        clientSessionId: clientSessionId.value
    };

    projectsStore.saveProject({
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        project: req as any
    }).then(() => {
        submitting.value = false;
        resolveFunc?.({ message: editProjectId.value ? 'Project saved' : 'Project added' });
        showState.value = false;
    }).catch(error => {
        submitting.value = false;
        snackbar.value?.showError(error);
    });
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
