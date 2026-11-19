<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ tt('Transaction Templates') }}</span>
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
                        <v-spacer/>
                        <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                               :aria-label="tt('More')" :disabled="loading || updating" :icon="true">
                            <v-icon :icon="mdiDotsVertical" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :prepend-icon="mdiEyeOutline"
                                                 :title="tt('Show Hidden Transaction Templates')"
                                                 v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                    <v-list-item :prepend-icon="mdiEyeOffOutline"
                                                 :title="tt('Hide Hidden Transaction Templates')"
                                                 v-if="showHidden" @click="showHidden = false"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </div>
                </template>

                <v-table class="transaction-templates-table table-striped" :hover="!loading">
                    <thead>
                    <tr>
                        <th>
                            <div class="d-flex align-center">
                                <span>{{ tt('Template Name') }}</span>
                                <v-spacer/>
                                <span>{{ tt('Operation') }}</span>
                            </div>
                        </th>
                    </tr>
                    </thead>

                    <tbody v-if="loading && noAvailableTemplate">
                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                        <td class="px-0">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>
                    </tbody>

                    <tbody v-if="!loading && noAvailableTemplate">
                    <tr>
                        <td>{{ tt('No available template. Once you add templates, you can quickly add a new transaction using the dropdown menu of the Add button on the transaction list page') }}</td>
                    </tr>
                    </tbody>

                    <draggable-list tag="tbody"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    :disabled="noAvailableTemplate"
                                    v-model="templates"
                                    @change="onMove">
                        <template #item="{ element }">
                            <tr class="transaction-templates-table-row" v-if="showHidden || !element.hidden"
                                @mouseenter="hoveredTemplateId = element.id" @mouseleave="hoveredTemplateId = ''">
                                <td>
                                    <div class="d-flex align-center">
                                        <div class="d-flex align-center">
                                            <v-badge class="right-bottom-icon" color="secondary"
                                                     location="bottom right" offset-x="8"
                                                     :aria-label="tt('Hidden')" :icon="mdiEyeOffOutline"
                                                     v-if="element.hidden">
                                                <v-icon size="20" start :icon="mdiTextBoxOutline"/>
                                            </v-badge>
                                            <v-icon size="20" start :icon="mdiTextBoxOutline" v-else-if="!element.hidden"/>
                                            <span class="transaction-template-name">{{ element.name }}</span>
                                        </div>

                                        <v-spacer/>

                                        <template v-if="hoveredTemplateId === element.id && !loading">
                                            <v-btn class="px-2 ms-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="element.hidden ? mdiEyeOutline : mdiEyeOffOutline"
                                                   :loading="templateHiding[element.id]"
                                                   :disabled="loading || updating"
                                                   @click="hide(element, !element.hidden)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ element.hidden ? tt('Show') : tt('Hide') }}
                                            </v-btn>
                                            <v-btn class="px-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="mdiPencilOutline"
                                                   :disabled="loading || updating"
                                                   @click="edit(element)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ tt('Edit') }}
                                            </v-btn>
                                            <v-btn class="px-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="mdiDeleteOutline"
                                                   :loading="templateRemoving[element.id]"
                                                   :disabled="loading || updating"
                                                   @click="remove(element)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ tt('Delete') }}
                                            </v-btn>
                                        </template>

                                        <span class="ms-2">
                                            <v-icon :class="!loading && !updating && availableTemplateCount > 1 ? 'drag-handle' : 'disabled'"
                                                    :aria-label="tt('Drag to Reorder')" :icon="mdiDrag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && availableTemplateCount > 1 && hoveredTemplateId === element.id">{{ tt('Drag to Reorder') }}</v-tooltip>
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

    <edit-dialog ref="editDialog" :type="TransactionEditPageType.Template" />

    <confirm-dialog ref="confirmDialog"/>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import EditDialog from '@/views/desktop/transactions/list/dialogs/EditDialog.vue';
import { TransactionEditPageType } from '@/views/base/transactions/TransactionEditPageBase.ts';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import { TemplateType } from '@/core/template.ts';
import { TransactionTemplate } from '@/models/transaction_template.ts';

import {
    isNoAvailableTemplate,
    getAvailableTemplateCount
} from '@/lib/template.ts';

import {
    mdiRefresh,
    mdiPencilOutline,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical,
    mdiTextBoxOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type EditDialogType = InstanceType<typeof EditDialog>;

const { tt } = useI18n();

const transactionTemplatesStore = useTransactionTemplatesStore();
const desktopPageStore = useDesktopPageStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const hoveredTemplateId = ref<string>('');
const templateHiding = ref<Record<string, boolean>>({});
const templateRemoving = ref<Record<string, boolean>>({});
const displayOrderModified = ref<boolean>(false);
const showHidden = ref<boolean>(false);

const templates = computed<TransactionTemplate[]>(() => transactionTemplatesStore.allNormalTransactionTemplates);
const noAvailableTemplate = computed<boolean>(() => isNoAvailableTemplate(templates.value, showHidden.value));
const availableTemplateCount = computed<number>(() => getAvailableTemplateCount(templates.value, showHidden.value));

function init(): void {
    loading.value = true;

    transactionTemplatesStore.loadAllNormalTemplates({
        force: false
    }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function reload(): void {
    loading.value = true;

    transactionTemplatesStore.loadAllNormalTemplates({
        force: true
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        desktopPageStore.addToast('Template list has been updated');
    }).catch(error => {
        loading.value = false;

        if (error && error.isUpToDate) {
            displayOrderModified.value = false;
        }

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function add(): void {
    editDialog.value?.open({
        templateType: TemplateType.Normal.type
    }).then(result => {
        if (result && result.message) {
            desktopPageStore.addToast(result.message);
        }
    }).catch(error => {
        if (error) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function edit(template: TransactionTemplate): void {
    editDialog.value?.open({
        id: template.id,
        currentTemplate: template
    }).then(result => {
        if (result && result.message) {
            desktopPageStore.addToast(result.message);
        }
    }).catch(error => {
        if (error) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function hide(template: TransactionTemplate, hidden: boolean): void {
    updating.value = true;
    templateHiding.value[template.id] = true;

    transactionTemplatesStore.hideTemplate({
        template: template,
        hidden: hidden
    }).then(() => {
        updating.value = false;
        templateHiding.value[template.id] = false;
    }).catch(error => {
        updating.value = false;
        templateHiding.value[template.id] = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function remove(template: TransactionTemplate): void {
    confirmDialog.value?.open('Are you sure you want to delete this template?').then(() => {
        updating.value = true;
        templateRemoving.value[template.id] = true;

        transactionTemplatesStore.deleteTemplate({
            template: template
        }).then(() => {
            updating.value = false;
            templateRemoving.value[template.id] = false;
        }).catch(error => {
            updating.value = false;
            templateRemoving.value[template.id] = false;

            if (!error.processed) {
                desktopPageStore.addErrorToast(error);
            }
        });
    });
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        return;
    }

    loading.value = true;

    transactionTemplatesStore.updateNormalTemplateDisplayOrders().then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function onMove(event: { moved: { element: { id: string }, oldIndex: number, newIndex: number } }): void {
    if (!event || !event.moved) {
        return;
    }

    const moveEvent = event.moved;

    if (!moveEvent.element || !moveEvent.element.id) {
        desktopPageStore.addToast('Unable to move template');
        return;
    }

    transactionTemplatesStore.changeNormalTemplateDisplayOrder({
        templateId: moveEvent.element.id,
        from: moveEvent.oldIndex,
        to: moveEvent.newIndex
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        desktopPageStore.addErrorToast(error);
    });
}

init();
</script>
