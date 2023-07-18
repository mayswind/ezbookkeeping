<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ $t('Transaction Tags') }}</span>
                        <v-btn density="compact" color="default" variant="text"
                               class="ml-2" :icon="true" :disabled="loading || updating || hasEditingTag"
                               v-if="!loading" @click="reload">
                            <v-icon :icon="icons.refresh" size="24" />
                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
                        <v-btn density="compact" color="default" variant="text"
                               class="ml-2" :icon="true" :disabled="loading || updating || hasEditingTag"
                               @click="add">
                            <v-icon :icon="icons.add" size="24" />
                            <v-tooltip activator="parent" v-if="!loading && !updating && !hasEditingTag">{{ $t('Add') }}</v-tooltip>
                        </v-btn>
                        <v-spacer/>
                        <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                               :disabled="loading || updating || hasEditingTag" :icon="true">
                            <v-icon :icon="icons.more" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :prepend-icon="icons.show"
                                                 :title="$t('Show Hidden Transaction Tag')"
                                                 v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                    <v-list-item :prepend-icon="icons.hide"
                                                 :title="$t('Hide Hidden Transaction Tag')"
                                                 v-if="showHidden" @click="showHidden = false"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </div>
                </template>

                <v-table class="transaction-tags-table">
                    <thead>
                    <tr>
                        <th class="text-uppercase" style="width: 50%">{{ $t('Tag Title') }}</th>
                        <th class="text-uppercase text-right">{{ $t('Operation') }}</th>
                    </tr>
                    </thead>

                    <tbody>
                    <tr :key="itemIdx" v-for="itemIdx in (loading && noAvailableTag ? [ 1, 2, 3 ] : [])">
                        <td class="px-0" colspan="2">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>

                    <tr v-if="!loading && noAvailableTag && !newTag">
                        <td colspan="2">{{ $t('No available tag') }}</td>
                    </tr>

                    <template :key="tag.id" v-for="tag in tags">
                        <tr v-if="showHidden || !tag.hidden">
                            <td>
                                <div class="d-flex align-center" v-if="editingTag.id !== tag.id">
                                    <v-badge class="right-bottom-icon" color="secondary"
                                             location="bottom right" offset-x="8" :icon="icons.hide"
                                             v-if="tag.hidden">
                                        <v-icon size="24" start :icon="icons.tag"/>
                                    </v-badge>
                                    <v-icon size="24" start :icon="icons.tag" v-else-if="!tag.hidden"/>
                                    {{ tag.name }}
                                </div>
                                <v-text-field
                                    type="text"
                                    clearable
                                    density="compact"
                                    variant="underlined"
                                    :disabled="loading || updating"
                                    :placeholder="$t('Tag Title')"
                                    v-model="editingTag.name"
                                    v-else-if="editingTag.id === tag.id"
                                    @keyup.enter="save(newTag)"
                                >
                                    <template #prepend>
                                        <v-badge class="right-bottom-icon" color="secondary"
                                                 location="bottom right" offset-x="8" :icon="icons.hide"
                                                 v-if="tag.hidden">
                                            <v-icon size="24" start :icon="icons.tag"/>
                                        </v-badge>
                                        <v-icon size="24" start :icon="icons.tag" v-else-if="!tag.hidden"/>
                                    </template>
                                </v-text-field>
                            </td>
                            <td class="text-uppercase text-right">
                                <v-btn class="px-2" color="default"
                                       density="comfortable" variant="text"
                                       :prepend-icon="icons.edit"
                                       :loading="tagUpdating[tag.id]"
                                       :disabled="loading || updating"
                                       v-if="editingTag.id !== tag.id"
                                       @click="edit(tag)">
                                    {{ $t('Edit') }}
                                </v-btn>
                                <v-btn class="px-2 ml-2" color="default"
                                       density="comfortable" variant="text"
                                       :prepend-icon="tag.hidden ? icons.show : icons.hide"
                                       :loading="tagHiding[tag.id]"
                                       :disabled="loading || updating"
                                       v-if="editingTag.id !== tag.id"
                                       @click="hide(tag, !tag.hidden)">
                                    {{ tag.hidden ? $t('Show') : $t('Hide') }}
                                </v-btn>
                                <v-btn class="px-2 ml-2" color="default"
                                       density="comfortable" variant="text"
                                       :prepend-icon="icons.remove"
                                       :loading="tagRemoving[tag.id]"
                                       :disabled="loading || updating"
                                       v-if="editingTag.id !== tag.id"
                                       @click="remove(tag)">
                                    {{ $t('Delete') }}
                                </v-btn>
                                <v-btn class="px-2"
                                       density="comfortable" variant="text"
                                       :prepend-icon="icons.confirm"
                                       :loading="tagUpdating[tag.id]"
                                       :disabled="loading || updating || !isTagModified(tag)"
                                       v-if="editingTag.id === tag.id" @click="save(editingTag)">
                                    {{ $t('Save') }}
                                </v-btn>
                                <v-btn class="px-2 ml-2" color="default"
                                       density="comfortable" variant="text"
                                       :prepend-icon="icons.cancel"
                                       :disabled="loading || updating"
                                       v-if="editingTag.id === tag.id" @click="cancelSave(editingTag)">
                                    {{ $t('Cancel') }}
                                </v-btn>
                            </td>
                        </tr>
                    </template>

                    <tr v-if="newTag">
                        <td>
                            <v-text-field type="text" color="primary" clearable
                                          density="compact" variant="underlined"
                                          :disabled="loading || updating" :placeholder="$t('Tag Title')"
                                          v-model="newTag.name" @keyup.enter="save(newTag)">
                                <template #prepend>
                                    <v-icon size="24" start :icon="icons.tag"/>
                                </template>
                            </v-text-field>
                        </td>
                        <td class="text-uppercase text-right">
                            <v-btn class="px-2" density="comfortable" variant="text"
                                   :prepend-icon="icons.confirm"
                                   :loading="tagUpdating[null]"
                                   :disabled="loading || updating || !isTagModified(newTag)"
                                   @click="save(newTag)">
                                {{ $t('Save') }}
                            </v-btn>
                            <v-btn class="px-2 ml-2" color="default"
                                   density="comfortable" variant="text"
                                   :prepend-icon="icons.cancel"
                                   :disabled="loading || updating"
                                   @click="cancelSave(newTag)">
                                {{ $t('Cancel') }}
                            </v-btn>
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
    mdiDotsVertical,
    mdiPound
} from '@mdi/js';

export default {
    data() {
        return {
            newTag: null,
            editingTag: {
                id: '',
                name: ''
            },
            loading: true,
            updating: false,
            tagUpdating: {},
            tagHiding: {},
            tagRemoving: {},
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
                self.$refs.snackbar.showMessage('Tag list has been updated');
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
.transaction-tags-table .v-text-field .v-input__prepend {
    margin-right: 0;
    color: rgba(var(--v-theme-on-surface));
}

.transaction-tags-table .v-text-field .v-input__prepend .v-badge > .v-badge__wrapper > .v-icon {
    opacity: var(--v-medium-emphasis-opacity);
}

.transaction-tags-table .v-text-field.v-text-field--plain-underlined .v-input__prepend {
    padding-top: 10px;
}

.transaction-tags-table tr:last-child .v-text-field.v-text-field--plain-underlined .v-input__prepend {
    padding-top: 9px;
}

.transaction-tags-table .v-text-field .v-field__input {
    padding-top: 2px;
    color: rgba(var(--v-theme-on-surface));
}
</style>
