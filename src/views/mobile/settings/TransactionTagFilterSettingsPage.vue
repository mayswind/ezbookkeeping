<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableTag }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t(applyText)" :class="{ 'disabled': !hasAnyVisibleTag }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="combination-list-wrapper margin-vertical skeleton-text" v-if="loading">
            <f7-accordion-item>
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="combination-list-header combination-list-opened">
                            <f7-list-item>
                                <template #title>
                                    <span>Tags</span>
                                    <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content style="height: auto">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item checkbox class="disabled" title="Tag Name"
                                      :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                            <template #media>
                                <f7-icon f7="app_fill"></f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-list strong inset dividers accordion-list class="margin-top" v-if="!loading && !hasAnyVisibleTag">
            <f7-list-item :title="$t('No available tag')"></f7-list-item>
        </f7-list>

        <f7-block class="combination-list-wrapper margin-vertical" key="default" v-show="!loading && hasAnyVisibleTag">
            <f7-list class="margin-top-half margin-bottom" strong inset dividers v-if="type === 'statisticsCurrent'">
                <f7-list-item radio
                              :title="filterType.displayName"
                              :value="filterType.type"
                              :checked="tagFilterType === filterType.type"
                              :key="filterType.type"
                              v-for="filterType in allTagFilterTypes"
                              @change="tagFilterType = filterType.type">
                </f7-list-item>
            </f7-list>

            <f7-accordion-item :opened="collapseStates['default'].opened"
                               @accordion:open="collapseStates['default'].opened = true"
                               @accordion:close="collapseStates['default'].opened = false">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="combination-list-header"
                                 :class="collapseStates['default'].opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item>
                                <template #title>
                                    <span>{{ $t('Tags') }}</span>
                                    <f7-icon class="combination-list-chevron-icon" :f7="collapseStates['default'].opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates['default'].opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content">
                        <f7-list-item checkbox
                                      :title="transactionTag.name"
                                      :value="transactionTag.id"
                                      :checked="!filterTagIds[transactionTag.id]"
                                      :key="transactionTag.id"
                                      v-for="transactionTag in allTags"
                                      v-show="showHidden || !transactionTag.hidden"
                                      @change="selectTransactionTag">
                            <template #media>
                                <f7-icon class="transaction-tag-icon" f7="number">
                                    <f7-badge color="gray" class="right-bottom-icon" v-if="transactionTag.hidden">
                                        <f7-icon f7="eye_slash_fill"></f7-icon>
                                    </f7-badge>
                                </f7-icon>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="selectAll">{{ $t('Select All') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="selectNone">{{ $t('Select None') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleTag }" @click="selectInvert">{{ $t('Invert Selection') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ $t('Show Hidden Transaction Tags') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ $t('Hide Hidden Transaction Tags') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { TransactionTagFilterType } from '@/core/transaction.ts';

import {
    selectAll,
    selectNone,
    selectInvert
} from '@/lib/common.ts';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data: function () {
        return {
            loading: true,
            loadingError: null,
            type: null,
            filterTagIds: {},
            tagFilterType: TransactionTagFilterType.Default.type,
            showHidden: false,
            collapseStates: {
                'default': {
                    opened: true
                }
            },
            showMoreActionSheet: false
        }
    },
    computed: {
        ...mapStores(useTransactionTagsStore, useTransactionsStore, useStatisticsStore),
        title() {
            return 'Filter Transaction Tags';
        },
        applyText() {
            return 'Apply';
        },
        allTags() {
            return this.transactionTagsStore.allTransactionTags;
        },
        allTagFilterTypes() {
            return this.$locale.getAllTransactionTagFilterTypes();
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
        const query = self.f7route.query;

        self.type = query.type;

        self.transactionTagsStore.loadAllTags({
            force: false
        }).then(() => {
            self.loading = false;

            const allTransactionTagIds = {};

            for (const transactionTagId in self.transactionTagsStore.allTransactionTagsMap) {
                if (!Object.prototype.hasOwnProperty.call(self.transactionTagsStore.allTransactionTagsMap, transactionTagId)) {
                    continue;
                }

                const transactionTag = self.transactionTagsStore.allTransactionTagsMap[transactionTagId];
                allTransactionTagIds[transactionTag.id] = true;
            }

            if (self.type === 'statisticsCurrent') {
                const transactionTagIds = self.statisticsStore.transactionStatisticsFilter.tagIds ? self.statisticsStore.transactionStatisticsFilter.tagIds.split(',') : [];

                for (let i = 0; i < transactionTagIds.length; i++) {
                    const transactionTagId = transactionTagIds[i];
                    const transactionTag = self.transactionTagsStore.allTransactionTagsMap[transactionTagId];

                    if (transactionTag) {
                        allTransactionTagIds[transactionTag.id] = false;
                    }
                }
                self.filterTagIds = allTransactionTagIds;
                self.tagFilterType = self.statisticsStore.transactionStatisticsFilter.tagFilterType;
            } else if (self.type === 'transactionListCurrent') {
                for (const transactionTagId in self.transactionsStore.allFilterTagIds) {
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
                self.$toast('Parameter Invalid');
                self.loadingError = 'Parameter Invalid';
            }
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        save() {
            const self = this;
            const router = self.f7router;

            const filteredTagIds = {};
            let finalTagIds = '';
            let changed = true;

            for (const transactionTagId in self.filterTagIds) {
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

            if (this.type === 'statisticsCurrent') {
                changed = self.statisticsStore.updateTransactionStatisticsFilter({
                    tagIds: finalTagIds,
                    tagFilterType: self.tagFilterType
                });

                if (changed) {
                    self.statisticsStore.updateTransactionStatisticsInvalidState(true);
                }
            } else if (this.type === 'transactionListCurrent') {
                changed = self.transactionsStore.updateTransactionListFilter({
                    tagIds: finalTagIds
                });

                if (changed) {
                    self.transactionsStore.updateTransactionListInvalidState(true);
                }
            }

            router.back();
        },
        selectTransactionTag(e) {
            const transactionTagId = e.target.value;
            const transactionTag = this.transactionTagsStore.allTransactionTagsMap[transactionTagId];

            if (!transactionTag) {
                return;
            }

            this.filterTagIds[transactionTag.id] = !e.target.checked;
        },
        selectAll() {
            selectAll(this.filterTagIds, this.transactionTagsStore.allTransactionTagsMap);
        },
        selectNone() {
            selectNone(this.filterTagIds, this.transactionTagsStore.allTransactionTagsMap);
        },
        selectInvert() {
            selectInvert(this.filterTagIds, this.transactionTagsStore.allTransactionTagsMap);
        }
    }
}
</script>
