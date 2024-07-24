<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ $t('Transaction Tags') }}</span>
                        <v-btn class="ml-3" color="default" variant="outlined"
                               :disabled="loading || updating || hasEditingTag" @click="add">{{ $t('Add') }}</v-btn>
                        <v-btn class="ml-3" color="primary" variant="tonal"
                               :disabled="loading || updating || hasEditingTag" @click="saveSortResult"
                               v-if="displayOrderModified">{{ $t('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ml-2" :icon="true" :disabled="loading || updating || hasEditingTag"
                               :loading="loading" @click="reload">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="icons.refresh" size="24" />
                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-spacer/>
                        <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                               :disabled="loading || updating || hasEditingTag" :icon="true">
                            <v-icon :icon="icons.more" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :prepend-icon="icons.show"
                                                 :title="$t('Show Hidden Transaction Tags')"
                                                 v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                    <v-list-item :prepend-icon="icons.hide"
                                                 :title="$t('Hide Hidden Transaction Tags')"
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
                                <span>{{ $t('Tag Title') }}</span>
                                <v-spacer/>
                                <span>{{ $t('Operation') }}</span>
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
                        <td>{{ $t('No available tag') }}</td>
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
                                                     location="bottom right" offset-x="8" :icon="icons.hide"
                                                     v-if="element.hidden">
                                                <v-icon size="20" start :icon="icons.tag"/>
                                            </v-badge>
                                            <v-icon size="20" start :icon="icons.tag" v-else-if="!element.hidden"/>
                                            <span class="transaction-tag-name">{{ element.name }}</span>
                                        </div>

                                        <v-text-field class="w-100 mr-2" type="text"
                                            density="compact" variant="underlined"
                                            :disabled="loading || updating"
                                            :placeholder="$t('Tag Title')"
                                            v-model="editingTag.name"
                                            v-else-if="editingTag.id === element.id"
                                            @keyup.enter="save(editingTag)"
                                        >
                                            <template #prepend>
                                                <v-badge class="right-bottom-icon" color="secondary"
                                                         location="bottom right" offset-x="8" :icon="icons.hide"
                                                         v-if="element.hidden">
                                                    <v-icon size="20" start :icon="icons.tag"/>
                                                </v-badge>
                                                <v-icon size="20" start :icon="icons.tag" v-else-if="!element.hidden"/>
                                            </template>
                                        </v-text-field>

                                        <v-spacer/>

                                        <v-btn class="px-2 ml-2" color="default"
                                               density="comfortable" variant="text"
                                               :class="{ 'd-none': loading, 'hover-display': !loading }"
                                               :prepend-icon="element.hidden ? icons.show : icons.hide"
                                               :loading="tagHiding[element.id]"
                                               :disabled="loading || updating"
                                               v-if="editingTag.id !== element.id"
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
                                               :loading="tagUpdating[element.id]"
                                               :disabled="loading || updating"
                                               v-if="editingTag.id !== element.id"
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
                                               :loading="tagRemoving[element.id]"
                                               :disabled="loading || updating"
                                               v-if="editingTag.id !== element.id"
                                               @click="remove(element)">
                                            <template #loader>
                                                <v-progress-circular indeterminate size="20" width="2"/>
                                            </template>
                                            {{ $t('Delete') }}
                                        </v-btn>
                                        <v-btn class="px-2"
                                               density="comfortable" variant="text"
                                               :prepend-icon="icons.confirm"
                                               :loading="tagUpdating[element.id]"
                                               :disabled="loading || updating || !isTagModified(element)"
                                               v-if="editingTag.id === element.id" @click="save(editingTag)">
                                            <template #loader>
                                                <v-progress-circular indeterminate size="20" width="2"/>
                                            </template>
                                            {{ $t('Save') }}
                                        </v-btn>
                                        <v-btn class="px-2" color="default"
                                               density="comfortable" variant="text"
                                               :prepend-icon="icons.cancel"
                                               :disabled="loading || updating"
                                               v-if="editingTag.id === element.id" @click="cancelSave(editingTag)">
                                            {{ $t('Cancel') }}
                                        </v-btn>
                                        <span class="ml-2">
                                            <v-icon :class="!loading && !updating && availableTagCount > 1 ? 'drag-handle' : 'disabled'"
                                                    :icon="icons.drag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && availableTagCount > 1">{{ $t('Drag to Reorder') }}</v-tooltip>
                                        </span>
                                    </div>
                                </td>
                            </tr>
                        </template>
                    </draggable-list>

                    <tbody v-if="newTag">
                    <tr class="text-sm" :class="{ 'even-row': availableTagCount & 1 === 1}">
                        <td>
                            <div class="d-flex align-center">
                                <v-text-field class="w-100 mr-2" type="text" color="primary"
                                              density="compact" variant="underlined"
                                              :disabled="loading || updating" :placeholder="$t('Tag Title')"
                                              v-model="newTag.name" @keyup.enter="save(newTag)">
                                    <template #prepend>
                                        <v-icon size="20" start :icon="icons.tag"/>
                                    </template>
                                </v-text-field>

                                <v-spacer/>

                                <v-btn class="px-2" density="comfortable" variant="text"
                                       :prepend-icon="icons.confirm"
                                       :loading="tagUpdating[null]"
                                       :disabled="loading || updating || !isTagModified(newTag)"
                                       @click="save(newTag)">
                                    <template #loader>
                                        <v-progress-circular indeterminate size="20" width="2"/>
                                    </template>
                                    {{ $t('Save') }}
                                </v-btn>
                                <v-btn class="px-2" color="default"
                                       density="comfortable" variant="text"
                                       :prepend-icon="icons.cancel"
                                       :disabled="loading || updating"
                                       @click="cancelSave(newTag)">
                                    {{ $t('Cancel') }}
                                </v-btn>
                                <span class="ml-2">
                                    <v-icon class="disabled" :icon="icons.drag"/>
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

<script>
import { mapStores } from 'pinia';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';

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
    mdiPound
} from '@mdi/js';

export default {
    data() {
        const transactionTagsStore = useTransactionTagsStore();
        const newTransactionTag = transactionTagsStore.generateNewTransactionTagModel();

        return {
            newTag: null,
            editingTag: newTransactionTag,
            loading: true,
            updating: false,
            tagUpdating: {},
            tagHiding: {},
            tagRemoving: {},
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
                tag: mdiPound
            }
        };
    },
    computed: {
        ...mapStores(useTransactionTagsStore),
        tags() {
            return this.transactionTagsStore.allTransactionTags;
        },
        noAvailableTag() {
            for (let i = 0; i < this.tags.length; i++) {
                if (this.showHidden || !this.tags[i].hidden) {
                    return false;
                }
            }

            return true;
        },
        availableTagCount() {
            let count = 0;

            for (let i = 0; i < this.tags.length; i++) {
                if (this.showHidden || !this.tags[i].hidden) {
                    count++;
                }
            }

            return count;
        },
        hasEditingTag() {
            return !!(this.newTag || (this.editingTag.id && this.editingTag.id !== ''));
        }
    },
    created() {
        const self = this;

        self.loading = true;

        self.transactionTagsStore.loadAllTags({
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
            if (this.hasEditingTag) {
                return;
            }

            const self = this;
            self.loading = true;

            self.transactionTagsStore.loadAllTags({
                force: true
            }).then(() => {
                self.loading = false;
                self.displayOrderModified = false;

                self.$refs.snackbar.showMessage('Tag list has been updated');
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
                self.$refs.snackbar.showMessage('Unable to move tag');
                return;
            }

            self.transactionTagsStore.changeTagDisplayOrder({
                tagId: moveEvent.element.id,
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

            self.transactionTagsStore.updateTagDisplayOrders().then(() => {
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
            this.newTag = {
                name: ''
            };
        },
        edit(tag) {
            this.editingTag.id = tag.id;
            this.editingTag.name = tag.name;
        },
        save(tag) {
            const self = this;

            self.updating = true;
            self.tagUpdating[tag.id || null] = true;

            self.transactionTagsStore.saveTag({
                tag: tag
            }).then(() => {
                self.updating = false;
                self.tagUpdating[tag.id || null] = false;

                if (tag.id) {
                    self.editingTag.id = '';
                    self.editingTag.name = '';
                } else {
                    self.newTag = null;
                }
            }).catch(error => {
                self.updating = false;
                self.tagUpdating[tag.id || null] = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        cancelSave(tag) {
            if (tag.id) {
                this.editingTag.id = '';
                this.editingTag.name = '';
            } else {
                this.newTag = null;
            }
        },
        isTagModified(tag) {
            if (tag.id) {
                return this.editingTag.name !== '' && this.editingTag.name !== tag.name;
            } else {
                return tag.name !== '';
            }
        },
        hide(tag, hidden) {
            const self = this;

            self.updating = true;
            self.tagHiding[tag.id] = true;

            self.transactionTagsStore.hideTag({
                tag: tag,
                hidden: hidden
            }).then(() => {
                self.updating = false;
                self.tagHiding[tag.id] = false;
            }).catch(error => {
                self.updating = false;
                self.tagHiding[tag.id] = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        remove(tag) {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to delete this tag?').then(() => {
                self.updating = true;
                self.tagRemoving[tag.id] = true;

                self.transactionTagsStore.deleteTag({
                    tag: tag
                }).then(() => {
                    self.updating = false;
                    self.tagRemoving[tag.id] = false;
                }).catch(error => {
                    self.updating = false;
                    self.tagRemoving[tag.id] = false;

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
    margin-right: 0;
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
