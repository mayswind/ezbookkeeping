<template>
    <v-card :class="{ 'pa-2 pa-sm-4 pa-md-8': dialogMode }">
        <template #title>
            <div class="d-flex align-center justify-center" v-if="dialogMode">
                <div class="w-100 text-center">
                    <h4 class="text-h4">{{ $t(title) }}</h4>
                </div>
                <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                       :disabled="loading || !hasAnyAvailableTag" :icon="true">
                    <v-icon :icon="icons.more" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="icons.selectAll"
                                         :title="$t('Select All')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectInvert"></v-list-item>
                            <v-divider class="my-2"/>
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
            <div class="d-flex align-center" v-else-if="!dialogMode">
                <span>{{ $t(title) }}</span>
                <v-spacer/>
                <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                       :disabled="loading" :icon="true">
                    <v-icon :icon="icons.more" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="icons.selectAll"
                                         :title="$t('Select All')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         :disabled="!hasAnyVisibleTag"
                                         @click="selectInvert"></v-list-item>
                            <v-divider class="my-2"/>
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

        <div v-if="loading">
            <v-skeleton-loader type="paragraph" :loading="loading"
                               :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
        </div>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-if="!loading && !hasAnyVisibleTag">
            <span class="text-body-1">{{ $t('No available tag') }}</span>
        </v-card-text>

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-else-if="!loading && hasAnyVisibleTag">
            <v-expansion-panels class="tag-categories" multiple v-model="expandTagCategories">
                <v-expansion-panel class="border" key="default" value="default">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ml-3">{{ $t('Tags') }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                        <v-list rounded density="comfortable" class="pa-0">
                            <template :key="transactionTag.id"
                                      v-for="transactionTag in allTags">
                                <v-list-item v-if="showHidden || !transactionTag.hidden">
                                    <template #prepend>
                                        <v-checkbox :model-value="!filterTagIds[transactionTag.id]"
                                                    @update:model-value="selectTransactionTag(transactionTag, $event)">
                                            <template #label>
                                                <v-badge class="right-bottom-icon" color="secondary"
                                                         location="bottom right" offset-x="2" offset-y="2" :icon="icons.hide"
                                                         v-if="transactionTag.hidden">
                                                    <v-icon size="24" :icon="icons.tag"/>
                                                </v-badge>
                                                <v-icon size="24" :icon="icons.tag" v-else-if="!transactionTag.hidden"/>
                                                <span class="ml-3">{{ transactionTag.name }}</span>
                                            </template>
                                        </v-checkbox>
                                    </template>
                                </v-list-item>
                            </template>
                        </v-list>
                    </v-expansion-panel-text>
                </v-expansion-panel>
            </v-expansion-panels>
        </v-card-text>

        <v-card-text class="overflow-y-visible" v-if="dialogMode">
            <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                <v-btn :disabled="!hasAnyVisibleTag" @click="save">{{ $t(applyText) }}</v-btn>
                <v-btn color="secondary" variant="tonal" @click="cancel">{{ $t('Cancel') }}</v-btn>
            </div>
        </v-card-text>
    </v-card>
    
    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';
import { useTransactionsStore } from '@/stores/transaction.js';

import {
    selectAll,
    selectNone,
    selectInvert
} from '@/lib/common.js';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiDotsVertical,
    mdiPound
} from '@mdi/js';

export default {
    props: [
        'dialogMode',
        'type',
        'autoSave'
    ],
    emits: [
        'settings:change'
    ],
    data: function () {
        return {
            loading: true,
            expandTagCategories: [ 'default' ],
            filterTagIds: {},
            showHidden: false,
            icons: {
                selectAll: mdiSelectAll,
                selectNone: mdiSelect,
                selectInverse: mdiSelectInverse,
                show: mdiEyeOutline,
                hide: mdiEyeOffOutline,
                more: mdiDotsVertical,
                tag: mdiPound
            }
        }
    },
    computed: {
        ...mapStores(useTransactionTagsStore, useTransactionsStore),
        title() {
            return 'Filter Transaction Tags';
        },
        applyText() {
            return 'Apply';
        },
        allTags() {
            return this.transactionTagsStore.allTransactionTags;
        },
        hasAnyAvailableTag() {
            return this.transactionTagsStore.allAvailableTagsCount > 0;
        },
        hasAnyVisibleTag() {
            if (this.showHidden) {
                return this.transactionTagsStore.allAvailableTagsCount > 0;
            } else {
                return this.transactionTagsStore.allVisibleTagsCount > 0;
            }
        }
    },
    created() {
        const self = this;

        self.transactionTagsStore.loadAllTags({
            force: false
        }).then(() => {
            self.loading = false;

            const allTransactionTagIds = {};

            for (let transactionTagId in self.transactionTagsStore.allTransactionTagsMap) {
                if (!Object.prototype.hasOwnProperty.call(self.transactionTagsStore.allTransactionTagsMap, transactionTagId)) {
                    continue;
                }

                const transactionTag = self.transactionTagsStore.allTransactionTagsMap[transactionTagId];

                if (self.type === 'transactionListCurrent' && self.transactionsStore.allFilterTagIdsCount > 0) {
                    allTransactionTagIds[transactionTag.id] = true;
                } else {
                    allTransactionTagIds[transactionTag.id] = false;
                }
            }

            if (self.type === 'transactionListCurrent') {
                for (let transactionTagId in self.transactionsStore.allFilterTagIds) {
                    if (!Object.prototype.hasOwnProperty.call(self.transactionsStore.allFilterTagIds, transactionTagId)) {
                        continue;
                    }

                    const transactionTag = self.transactionTagsStore.allTransactionTagsMap[transactionTagId];

                    if (transactionTag) {
                        allTransactionTagIds[transactionTag.id] = false;
                    }
                }
                self.filterTagIds = allTransactionTagIds;
            } else {
                self.$refs.snackbar.showError('Parameter Invalid');
            }
        }).catch(error => {
            self.loading = false;

            if (!error.processed) {
                self.$refs.snackbar.showError(error);
            }
        });
    },
    methods: {
        save() {
            const self = this;

            const filteredTagIds = {};
            let finalTagIds = '';
            let changed = true;

            for (let transactionTagId in self.filterTagIds) {
                if (!Object.prototype.hasOwnProperty.call(self.filterTagIds, transactionTagId)) {
                    continue;
                }

                const transactionTag = self.transactionTagsStore.allTransactionTagsMap[transactionTagId];

                if (self.filterTagIds[transactionTag.id]) {
                    filteredTagIds[transactionTag.id] = true;
                } else {
                    if (finalTagIds.length > 0) {
                        finalTagIds += ',';
                    }

                    finalTagIds += transactionTag.id;
                }
            }

            if (this.type === 'transactionListCurrent') {
                changed = self.transactionsStore.updateTransactionListFilter({
                    tagIds: finalTagIds
                });

                if (changed) {
                    self.transactionsStore.updateTransactionListInvalidState(true);
                }
            }

            self.$emit('settings:change', changed);
        },
        cancel() {
            this.$emit('settings:change', false);
        },
        selectTransactionTag(transactionTag, value) {
            this.filterTagIds[transactionTag.id] = !value;

            if (this.autoSave) {
                this.save();
            }
        },
        selectAll() {
            selectAll(this.filterTagIds, this.transactionTagsStore.allTransactionTagsMap);

            if (this.autoSave) {
                this.save();
            }
        },
        selectNone() {
            selectNone(this.filterTagIds, this.transactionTagsStore.allTransactionTagsMap);

            if (this.autoSave) {
                this.save();
            }
        },
        selectInvert() {
            selectInvert(this.filterTagIds, this.transactionTagsStore.allTransactionTagsMap);

            if (this.autoSave) {
                this.save();
            }
        }
    }
}
</script>

<style>
.tag-categories .v-expansion-panel-text__wrapper {
    padding: 0 0 0 20px;
}

.tag-categories .v-expansion-panel--active:not(:first-child),
.tag-categories .v-expansion-panel--active + .v-expansion-panel {
    margin-top: 1rem;
}
</style>
