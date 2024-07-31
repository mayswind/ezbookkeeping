<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ $t('Transaction Templates') }}</span>
                        <v-btn class="ml-3" color="default" variant="outlined"
                               :disabled="loading || updating" @click="add">{{ $t('Add') }}</v-btn>
                        <v-btn class="ml-3" color="primary" variant="tonal"
                               :disabled="loading || updating" @click="saveSortResult"
                               v-if="displayOrderModified">{{ $t('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ml-2" :icon="true" :disabled="loading || updating"
                               :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="icons.refresh" size="24" />
                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-spacer/>
                        <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                               :disabled="loading || updating" :icon="true">
                            <v-icon :icon="icons.more" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :prepend-icon="icons.show"
                                                 :title="$t('Show Hidden Templates')"
                                                 v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                    <v-list-item :prepend-icon="icons.hide"
                                                 :title="$t('Hide Hidden Templates')"
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
                                <span>{{ $t('Template Name') }}</span>
                                <v-spacer/>
                                <span>{{ $t('Operation') }}</span>
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
                        <td>{{ $t('No available template') }}</td>
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
                                                <v-icon size="20" start :icon="icons.text"/>
                                            </v-badge>
                                            <v-icon size="20" start :icon="icons.text" v-else-if="!element.hidden"/>
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
                                            {{ element.hidden ? $t('Show') : $t('Hide') }}
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
                                            {{ $t('Edit') }}
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
                                            {{ $t('Delete') }}
                                        </v-btn>
                                        <span class="ml-2">
                                            <v-icon :class="!loading && !updating && availableTemplateCount > 1 ? 'drag-handle' : 'disabled'"
                                                    :icon="icons.drag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && availableTemplateCount > 1">{{ $t('Drag to Reorder') }}</v-tooltip>
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

<script>
import EditDialog from '@/views/desktop/transactions/list/dialogs/EditDialog.vue';

import { mapStores } from 'pinia';
import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.js';

import templateConstants from '@/consts/template.js';

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
    mdiTextBoxOutline
} from '@mdi/js';

export default {
    components: {
        EditDialog
    },
    data() {
        return {
            templateType: templateConstants.allTemplateTypes.Normal,
            loading: true,
            updating: false,
            templateHiding: {},
            templateRemoving: {},
            displayOrderModified: false,
            showHidden: false,
            icons: {
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
                text: mdiTextBoxOutline
            }
        };
    },
    computed: {
        ...mapStores(useTransactionTemplatesStore),
        templates() {
            return this.transactionTemplatesStore.allTransactionTemplates[this.templateType] || [];
        },
        noAvailableTemplate() {
            for (let i = 0; i < this.templates.length; i++) {
                if (this.showHidden || !this.templates[i].hidden) {
                    return false;
                }
            }

            return true;
        },
        availableTemplateCount() {
            let count = 0;

            for (let i = 0; i < this.templates.length; i++) {
                if (this.showHidden || !this.templates[i].hidden) {
                    count++;
                }
            }

            return count;
        }
    },
    created() {
        const self = this;

        self.loading = true;

        self.transactionTemplatesStore.loadAllTemplates({
            templateType: self.templateType,
            force: false
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            self.loading = false;

            if (!error.processed) {
                self.$refs.snackbar.showError(error);
            }
        });
    },
    methods: {
        reload() {
            const self = this;
            self.loading = true;

            self.transactionTemplatesStore.loadAllTemplates({
                templateType: self.templateType,
                force: true
            }).then(() => {
                self.loading = false;
                self.displayOrderModified = false;

                self.$refs.snackbar.showMessage('Template list has been updated');
            }).catch(error => {
                self.loading = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        onMove(event) {
            if (!event || !event.moved) {
                return;
            }

            const self = this;
            const moveEvent = event.moved;

            if (!moveEvent.element || !moveEvent.element.id) {
                self.$refs.snackbar.showMessage('Unable to move template');
                return;
            }

            self.transactionTemplatesStore.changeTemplateDisplayOrder({
                templateType: self.templateType,
                templateId: moveEvent.element.id,
                from: moveEvent.oldIndex,
                to: moveEvent.newIndex
            }).then(() => {
                self.displayOrderModified = true;
            }).catch(error => {
                self.$refs.snackbar.showError(error);
            });
        },
        saveSortResult() {
            const self = this;

            if (!self.displayOrderModified) {
                return;
            }

            self.loading = true;

            self.transactionTemplatesStore.updateTemplateDisplayOrders({
                templateType: self.templateType
            }).then(() => {
                self.loading = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.loading = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        add() {
            const self = this;

            self.$refs.editDialog.open({
                templateType: self.templateType
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        edit(template) {
            const self = this;

            self.$refs.editDialog.open({
                id: template.id,
                currentTemplate: {
                    name: template.name,
                    type: template.type,
                    categoryId: template.categoryId,
                    sourceAccountId: template.sourceAccountId,
                    destinationAccountId: template.destinationAccountId,
                    sourceAmount: template.sourceAmount,
                    destinationAmount: template.destinationAmount,
                    hideAmount: template.hideAmount,
                    tagIds: template.tagIds,
                    comment: template.comment
                }
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        hide(template, hidden) {
            const self = this;

            self.updating = true;
            self.templateHiding[template.id] = true;

            self.transactionTemplatesStore.hideTemplate({
                template: template,
                hidden: hidden
            }).then(() => {
                self.updating = false;
                self.templateHiding[template.id] = false;
            }).catch(error => {
                self.updating = false;
                self.templateHiding[template.id] = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        remove(template) {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to delete this template?').then(() => {
                self.updating = true;
                self.templateRemoving[template.id] = true;

                self.transactionTemplatesStore.deleteTemplate({
                    template: template
                }).then(() => {
                    self.updating = false;
                    self.templateRemoving[template.id] = false;
                }).catch(error => {
                    self.updating = false;
                    self.templateRemoving[template.id] = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        }
    }
}
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
