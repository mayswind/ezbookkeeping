<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ templateType === TemplateType.Schedule.type ? tt('Scheduled Transactions') : tt('Transaction Templates') }}</span>
                        <v-btn class="ml-3" color="default" variant="outlined"
                               :disabled="loading || updating" @click="add">{{ tt('Add') }}</v-btn>
                        <v-btn class="ml-3" color="primary" variant="tonal"
                               :disabled="loading || updating" @click="saveSortResult"
                               v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ml-2" :icon="true" :disabled="loading || updating"
                               :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="icons.refresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-spacer/>
                        <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                               :disabled="loading || updating" :icon="true">
                            <v-icon :icon="icons.more" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :prepend-icon="icons.show"
                                                 :title="tt('Show Hidden Transaction Templates')"
                                                 v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                    <v-list-item :prepend-icon="icons.hide"
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
                        <td v-if="templateType === TemplateType.Normal.type">{{ tt('No available template. Once you add templates, you can quickly add a new transaction using the dropdown menu of the Add button on the transaction list page') }}</td>
                        <td v-else-if="templateType === TemplateType.Schedule.type">{{ tt('No available scheduled transactions') }}</td>
                        <td v-else>{{ tt('No available template') }}</td>
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
                            <tr class="transaction-templates-table-row text-sm" v-if="showHidden || !element.hidden">
                                <td>
                                    <div class="d-flex align-center">
                                        <div class="d-flex align-center">
                                            <v-badge class="right-bottom-icon" color="secondary"
                                                     location="bottom right" offset-x="8" :icon="icons.hide"
                                                     v-if="element.hidden">
                                                <v-icon size="20" start :icon="templateType === TemplateType.Schedule.type ? icons.clock : icons.text"/>
                                            </v-badge>
                                            <v-icon size="20" start :icon="templateType === TemplateType.Schedule.type ? icons.clock : icons.text" v-else-if="!element.hidden"/>
                                            <span class="transaction-template-name">{{ element.name }}</span>
                                        </div>

                                        <v-spacer/>

                                        <v-btn class="px-2 ml-2" color="default"
                                               density="comfortable" variant="text"
                                               :class="{ 'd-none': loading, 'hover-display': !loading }"
                                               :prepend-icon="element.hidden ? icons.show : icons.hide"
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
                                               :class="{ 'd-none': loading, 'hover-display': !loading }"
                                               :prepend-icon="icons.edit"
                                               :disabled="loading || updating"
                                               @click="edit(element)">
                                            <template #loader>
                                                <v-progress-circular indeterminate size="20" width="2"/>
                                            </template>
                                            {{ tt('Edit') }}
                                        </v-btn>
                                        <v-btn class="px-2" color="default"
                                               density="comfortable" variant="text"
                                               :class="{ 'd-none': loading, 'hover-display': !loading }"
                                               :prepend-icon="icons.remove"
                                               :loading="templateRemoving[element.id]"
                                               :disabled="loading || updating"
                                               @click="remove(element)">
                                            <template #loader>
                                                <v-progress-circular indeterminate size="20" width="2"/>
                                            </template>
                                            {{ tt('Delete') }}
                                        </v-btn>
                                        <span class="ml-2">
                                            <v-icon :class="!loading && !updating && availableTemplateCount > 1 ? 'drag-handle' : 'disabled'"
                                                    :icon="icons.drag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && availableTemplateCount > 1">{{ tt('Drag to Reorder') }}</v-tooltip>
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

    <edit-dialog ref="editDialog" type="template" :persistent="true" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import EditDialog from '@/views/desktop/transactions/list/dialogs/EditDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';

import { TemplateType } from '@/core/template.ts';
import { TransactionTemplate } from '@/models/transaction_template.ts';

import {
    isNoAvailableTemplate,
    getAvailableTemplateCount
} from '@/lib/template.ts';

import {
    mdiRefresh,
    mdiPlus,
    mdiPencilOutline,
    mdiCheck,
    mdiClose,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical,
    mdiTextBoxOutline,
    mdiClockTimeNineOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type EditDialogType = InstanceType<typeof EditDialog>;

const props = defineProps<{
    initType: number;
}>();

const { tt } = useI18n();

const transactionTemplatesStore = useTransactionTemplatesStore();

const icons = {
    refresh: mdiRefresh,
    add: mdiPlus,
    edit: mdiPencilOutline,
    confirm: mdiCheck,
    cancel: mdiClose,
    show: mdiEyeOutline,
    hide: mdiEyeOffOutline,
    remove: mdiDeleteOutline,
    drag: mdiDrag,
    more: mdiDotsVertical,
    text: mdiTextBoxOutline,
    clock: mdiClockTimeNineOutline
};

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const templateType = ref<number>(TemplateType.Normal.type);
const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const templateHiding = ref<Record<string, boolean>>({});
const templateRemoving = ref<Record<string, boolean>>({});
const displayOrderModified = ref<boolean>(false);
const showHidden = ref<boolean>(false);

const templates = computed<TransactionTemplate[]>(() => transactionTemplatesStore.allTransactionTemplates[templateType.value] || []);
const noAvailableTemplate = computed<boolean>(() => isNoAvailableTemplate(templates.value, showHidden.value));
const availableTemplateCount = computed<number>(() => getAvailableTemplateCount(templates.value, showHidden.value));

function init(): void {
    templateType.value = props.initType;
    loading.value = true;

    transactionTemplatesStore.loadAllTemplates({
        templateType: templateType.value,
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

    transactionTemplatesStore.loadAllTemplates({
        templateType: templateType.value,
        force: true
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        snackbar.value?.showMessage('Template list has been updated');
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
    editDialog.value?.open({
        templateType: templateType.value
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

function edit(template: TransactionTemplate): void {
    editDialog.value?.open({
        id: template.id,
        currentTemplate: template
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
            snackbar.value?.showError(error);
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

    transactionTemplatesStore.updateTemplateDisplayOrders({
        templateType: templateType.value
    }).then(() => {
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
        snackbar.value?.showMessage('Unable to move template');
        return;
    }

    transactionTemplatesStore.changeTemplateDisplayOrder({
        templateType: templateType.value,
        templateId: moveEvent.element.id,
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
.transaction-templates-table tr.transaction-templates-table-row .hover-display {
    display: none;
}

.transaction-templates-table tr.transaction-templates-table-row:hover .hover-display {
    display: inline-grid;
}

.transaction-templates-table tr:not(:last-child) > td > div {
    padding-bottom: 1px;
}

.transaction-templates-table .has-bottom-border tr:last-child > td > div {
    padding-bottom: 1px;
}

.transaction-templates-table tr.transaction-templates-table-row .right-bottom-icon .v-badge__badge {
    padding-bottom: 1px;
}

.transaction-templates-table .v-text-field .v-input__prepend {
    margin-right: 0;
    color: rgba(var(--v-theme-on-surface));
}

.transaction-templates-table .v-text-field .v-input__prepend .v-badge > .v-badge__wrapper > .v-icon {
    opacity: var(--v-medium-emphasis-opacity);
}

.transaction-templates-table .v-text-field.v-input--plain-underlined .v-input__prepend {
    padding-top: 10px;
}

.transaction-templates-table .v-text-field .v-field__input {
    font-size: 0.875rem;
    padding-top: 0;
    color: rgba(var(--v-theme-on-surface));
}

.transaction-templates-table .transaction-template-name {
    font-size: 0.875rem;
}

.transaction-templates-table tr .v-text-field .v-field__input {
    padding-bottom: 1px;
}
</style>
