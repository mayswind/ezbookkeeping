<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableCategory }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t(applyText)" :class="{ 'disabled': !hasAnyVisibleCategory }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <div class="skeleton-text" v-if="loading">
            <f7-block class="combination-list-wrapper margin-vertical"
                      :key="blockIdx" v-for="blockIdx in [ 1, 2 ]">
                <f7-accordion-item>
                    <f7-block-title>
                        <f7-accordion-toggle>
                            <f7-list strong inset dividers media-list
                                     class="combination-list-header combination-list-opened">
                                <f7-list-item>
                                    <template #title>
                                        <span>Transaction Category</span>
                                        <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
                                    </template>
                                </f7-list-item>
                            </f7-list>
                        </f7-accordion-toggle>
                    </f7-block-title>
                    <f7-accordion-content style="height: auto">
                        <f7-list strong inset dividers accordion-list class="combination-list-content">
                            <f7-list-item checkbox class="disabled" title="Category Name"
                                          :key="itemIdx" v-for="itemIdx in [ 1, 2 ]">
                                <template #media>
                                    <f7-icon f7="app_fill"></f7-icon>
                                </template>
                                <template #root>
                                    <ul class="padding-left">
                                        <f7-list-item checkbox class="disabled" title="Sub Category Name"
                                                      :key="subItemIdx" v-for="subItemIdx in [ 1, 2, 3 ]">
                                            <template #media>
                                                <f7-icon f7="app_fill"></f7-icon>
                                            </template>
                                        </f7-list-item>
                                    </ul>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-content>
                </f7-accordion-item>
            </f7-block>
        </div>

        <f7-block class="combination-list-wrapper margin-vertical"
                  :key="transactionType.type"
                  v-for="transactionType in allTransactionCategories"
                  v-else-if="!loading">
            <f7-accordion-item :opened="collapseStates[transactionType.type].opened"
                               @accordion:open="collapseStates[transactionType.type].opened = true"
                               @accordion:close="collapseStates[transactionType.type].opened = false">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="combination-list-header"
                                 :class="collapseStates[transactionType.type].opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item>
                                <template #title>
                                    <span>{{ getCategoryTypeName(transactionType.type) }}</span>
                                    <f7-icon class="combination-list-chevron-icon" :f7="collapseStates[transactionType.type].opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates[transactionType.type].opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content" v-if="!hasAvailableCategory[transactionType.type]">
                        <f7-list-item :title="$t('No available category')"></f7-list-item>
                    </f7-list>
                    <f7-list strong inset dividers accordion-list class="combination-list-content" v-else-if="hasAvailableCategory[transactionType.type]">
                        <f7-list-item checkbox
                                      :class="{ 'has-child-list-item': (showHidden && transactionType.allSubCategories[category.id]) || transactionType.allVisibleSubCategoryCounts[category.id] }"
                                      :title="category.name"
                                      :value="category.id"
                                      :checked="isSubCategoriesAllChecked(category, filterCategoryIds)"
                                      :indeterminate="isSubCategoriesHasButNotAllChecked(category, filterCategoryIds)"
                                      :key="category.id"
                                      v-for="category in transactionType.allCategories"
                                      v-show="showHidden || !category.hidden"
                                      @change="selectSubCategories">
                            <template #media>
                                <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color">
                                    <f7-badge color="gray" class="right-bottom-icon" v-if="category.hidden">
                                        <f7-icon f7="eye_slash_fill"></f7-icon>
                                    </f7-badge>
                                </ItemIcon>
                            </template>

                            <template #root>
                                <ul class="padding-left"
                                    v-if="(showHidden && transactionType.allSubCategories[category.id]) || transactionType.allVisibleSubCategoryCounts[category.id]">
                                    <f7-list-item checkbox
                                                  :title="subCategory.name"
                                                  :value="subCategory.id"
                                                  :checked="isCategoryChecked(subCategory, filterCategoryIds)"
                                                  :key="subCategory.id"
                                                  v-for="subCategory in transactionType.allSubCategories[category.id]"
                                                  v-show="showHidden || !subCategory.hidden"
                                                  @change="selectCategory">
                                        <template #media>
                                            <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color">
                                                <f7-badge color="gray" class="right-bottom-icon" v-if="subCategory.hidden">
                                                    <f7-icon f7="eye_slash_fill"></f7-icon>
                                                </f7-badge>
                                            </ItemIcon>
                                        </template>
                                    </f7-list-item>
                                </ul>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleCategory }" @click="selectAll">{{ $t('Select All') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleCategory }" @click="selectNone">{{ $t('Select None') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleCategory }" @click="selectInvert">{{ $t('Invert Selection') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ $t('Show Hidden Transaction Categories') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ $t('Hide Hidden Transaction Categories') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { CategoryType } from '@/core/category.ts';
import { copyObjectTo, arrayItemToObjectField } from '@/lib/common.ts';
import {
    allTransactionCategoriesWithVisibleCount,
    hasAnyAvailableCategory,
    hasAvailableCategory,
    selectSubCategories,
    selectAll,
    selectNone,
    selectInvert,
    isCategoryOrSubCategoriesAllChecked,
    isSubCategoriesAllChecked,
    isSubCategoriesHasButNotAllChecked
} from '@/lib/category.ts';

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
            allowCategoryTypes: null,
            filterCategoryIds: {},
            showHidden: false,
            collapseStates: {
                [CategoryType.Income]: {
                    opened: true
                },
                [CategoryType.Expense]: {
                    opened: true
                },
                [CategoryType.Transfer]: {
                    opened: true
                }
            },
            showMoreActionSheet: false
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useTransactionCategoriesStore, useTransactionsStore, useStatisticsStore),
        title() {
            if (this.type === 'statisticsDefault') {
                return 'Default Transaction Category Filter';
            } else {
                return 'Filter Transaction Categories';
            }
        },
        applyText() {
            if (this.type === 'statisticsDefault') {
                return 'Save';
            } else {
                return 'Apply';
            }
        },
        allTransactionCategories() {
            return allTransactionCategoriesWithVisibleCount(this.transactionCategoriesStore.allTransactionCategories, this.allowCategoryTypes);
        },
        hasAnyAvailableCategory() {
            return hasAnyAvailableCategory(this.allTransactionCategories, true);
        },
        hasAnyVisibleCategory() {
            return hasAnyAvailableCategory(this.allTransactionCategories, this.showHidden);
        },
        hasAvailableCategory() {
            return hasAvailableCategory(this.allTransactionCategories, this.showHidden);
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        self.type = query.type;
        self.allowCategoryTypes = query.allowCategoryTypes ? arrayItemToObjectField(query.allowCategoryTypes.split(','), true) : null;

        self.transactionCategoriesStore.loadAllCategories({
            force: false
        }).then(() => {
            self.loading = false;

            const allCategoryIds = {};

            for (const categoryId in self.transactionCategoriesStore.allTransactionCategoriesMap) {
                if (!Object.prototype.hasOwnProperty.call(self.transactionCategoriesStore.allTransactionCategoriesMap, categoryId)) {
                    continue;
                }

                const category = self.transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                if (self.allowCategoryTypes && !self.allowCategoryTypes[category.type]) {
                    continue;
                }

                if (self.type === 'transactionListCurrent' && self.transactionsStore.allFilterCategoryIdsCount > 0) {
                    allCategoryIds[category.id] = true;
                } else {
                    allCategoryIds[category.id] = false;
                }
            }

            if (self.type === 'statisticsDefault') {
                self.filterCategoryIds = copyObjectTo(self.settingsStore.appSettings.statistics.defaultTransactionCategoryFilter, allCategoryIds);
            } else if (self.type === 'statisticsCurrent') {
                self.filterCategoryIds = copyObjectTo(self.statisticsStore.transactionStatisticsFilter.filterCategoryIds, allCategoryIds);
            } else if (self.type === 'transactionListCurrent') {
                for (const categoryId in self.transactionsStore.allFilterCategoryIds) {
                    if (!Object.prototype.hasOwnProperty.call(self.transactionsStore.allFilterCategoryIds, categoryId)) {
                        continue;
                    }

                    const category = self.transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                    if (category && (!category.subCategories || !category.subCategories.length)) {
                        allCategoryIds[category.id] = false;
                    } else if (category) {
                        selectSubCategories(allCategoryIds, category, false);
                    }
                }

                self.filterCategoryIds = allCategoryIds;
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

            const filteredCategoryIds = {};
            let isAllSelected = true;
            let finalCategoryIds = '';

            for (const categoryId in self.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(self.filterCategoryIds, categoryId)) {
                    continue;
                }

                const category = self.transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                if (!isCategoryOrSubCategoriesAllChecked(category, self.filterCategoryIds)) {
                    filteredCategoryIds[categoryId] = true;
                    isAllSelected = false;
                } else {
                    if (finalCategoryIds.length > 0) {
                        finalCategoryIds += ',';
                    }

                    finalCategoryIds += categoryId;
                }
            }

            if (this.type === 'statisticsDefault') {
                self.settingsStore.setStatisticsDefaultTransactionCategoryFilter(filteredCategoryIds);
            } else if (this.type === 'statisticsCurrent') {
                self.statisticsStore.updateTransactionStatisticsFilter({
                    filterCategoryIds: filteredCategoryIds
                });
            } else if (this.type === 'transactionListCurrent') {
                const changed = self.transactionsStore.updateTransactionListFilter({
                    categoryIds: isAllSelected ? '' : finalCategoryIds
                });

                if (changed) {
                    self.transactionsStore.updateTransactionListInvalidState(true);
                }
            }

            router.back();
        },
        selectCategory(e) {
            const categoryId = e.target.value;
            const category = this.transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

            if (!category) {
                return;
            }

            this.filterCategoryIds[category.id] = !e.target.checked;
        },
        selectSubCategories(e) {
            const categoryId = e.target.value;
            const category = this.transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

            selectSubCategories(this.filterCategoryIds, category, !e.target.checked);
        },
        selectAll() {
            selectAll(this.filterCategoryIds, this.transactionCategoriesStore.allTransactionCategoriesMap);
        },
        selectNone() {
            selectNone(this.filterCategoryIds, this.transactionCategoriesStore.allTransactionCategoriesMap);
        },
        selectInvert() {
            selectInvert(this.filterCategoryIds, this.transactionCategoriesStore.allTransactionCategoriesMap);
        },
        getCategoryTypeName(categoryType) {
            switch (categoryType) {
                case CategoryType.Income:
                    return this.$t('Income Categories');
                case CategoryType.Expense:
                    return this.$t('Expense Categories');
                case CategoryType.Transfer:
                    return this.$t('Transfer Categories');
                default:
                    return this.$t('Transaction Categories');
            }
        },
        isCategoryChecked(category, filterCategoryIds) {
            return !filterCategoryIds[category.id];
        },
        isSubCategoriesAllChecked(category, filterCategoryIds) {
            return isSubCategoriesAllChecked(category, filterCategoryIds);
        },
        isSubCategoriesHasButNotAllChecked(category, filterCategoryIds) {
            return isSubCategoriesHasButNotAllChecked(category, filterCategoryIds);
        }
    }
}
</script>
