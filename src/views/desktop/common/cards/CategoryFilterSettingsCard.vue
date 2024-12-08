<template>
    <v-card :class="{ 'pa-2 pa-sm-4 pa-md-8': dialogMode }">
        <template #title>
            <div class="d-flex align-center justify-center" v-if="dialogMode">
                <div class="w-100 text-center">
                    <h4 class="text-h4">{{ $t(title) }}</h4>
                </div>
                <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                       :disabled="loading || !hasAnyAvailableCategory" :icon="true">
                    <v-icon :icon="icons.more" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="icons.selectAll"
                                         :title="$t('Select All')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectInvert"></v-list-item>
                            <v-divider class="my-2"/>
                            <v-list-item :prepend-icon="icons.show"
                                         :title="$t('Show Hidden Transaction Categories')"
                                         v-if="!showHidden" @click="showHidden = true"></v-list-item>
                            <v-list-item :prepend-icon="icons.hide"
                                         :title="$t('Hide Hidden Transaction Categories')"
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
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectInvert"></v-list-item>
                            <v-divider class="my-2"/>
                            <v-list-item :prepend-icon="icons.show"
                                         :title="$t('Show Hidden Transaction Categories')"
                                         v-if="!showHidden" @click="showHidden = true"></v-list-item>
                            <v-list-item :prepend-icon="icons.hide"
                                         :title="$t('Hide Hidden Transaction Categories')"
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

        <v-card-text :class="{ 'mt-0 mt-sm-2 mt-md-4': dialogMode }" v-else-if="!loading">
            <v-expansion-panels class="category-types" multiple v-model="expandCategoryTypes">
                <v-expansion-panel :key="transactionType.type"
                                   :value="transactionType.type"
                                   class="border"
                                   v-for="transactionType in allTransactionCategories">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ml-3">{{ getCategoryTypeName(transactionType.type) }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                        <v-list rounded density="comfortable" class="pa-0">
                            <div class="py-3" v-if="!hasAvailableCategory[transactionType.type]">{{ $t('No available category') }}</div>

                            <template :key="category.id"
                                      v-for="(category, idx) in transactionType.allCategories">
                                <v-divider v-if="showHidden ? idx > 0 : (!category.hidden ? idx > transactionType.firstVisibleCategoryIndex : false)"/>

                                <v-list-item v-if="showHidden || !category.hidden">
                                    <template #prepend>
                                        <v-checkbox :model-value="isSubCategoriesAllChecked(category, filterCategoryIds)"
                                                    :indeterminate="isSubCategoriesHasButNotAllChecked(category, filterCategoryIds)"
                                                    @update:model-value="selectSubCategories(category, $event)">
                                            <template #label>
                                                <ItemIcon class="d-flex" icon-type="category" :icon-id="category.icon"
                                                          :color="category.color" :hidden-status="category.hidden"></ItemIcon>
                                                <span class="ml-3">{{ category.name }}</span>
                                            </template>
                                        </v-checkbox>
                                    </template>
                                </v-list-item>

                                <v-divider v-if="(showHidden || !category.hidden) && ((showHidden && transactionType.allSubCategories[category.id]) || transactionType.allVisibleSubCategoryCounts[category.id])"/>

                                <v-list rounded density="comfortable" class="pa-0 ml-4"
                                        v-if="(showHidden || !category.hidden) && ((showHidden && transactionType.allSubCategories[category.id]) || transactionType.allVisibleSubCategoryCounts[category.id])">
                                    <template :key="subCategory.id"
                                              v-for="(subCategory, subIdx) in transactionType.allSubCategories[category.id]">
                                        <v-divider v-if="showHidden ? subIdx > 0 : (!subCategory.hidden ? subIdx > transactionType.allFirstVisibleSubCategoryIndexes[category.id] : false)"/>

                                        <v-list-item v-if="showHidden || !subCategory.hidden">
                                            <template #prepend>
                                                <v-checkbox :model-value="isCategoryChecked(subCategory, filterCategoryIds)"
                                                            @update:model-value="selectCategory(subCategory, $event)">
                                                    <template #label>
                                                        <ItemIcon class="d-flex" icon-type="category" :icon-id="subCategory.icon"
                                                                  :color="subCategory.color" :hidden-status="subCategory.hidden"></ItemIcon>
                                                        <span class="ml-3">{{ subCategory.name }}</span>
                                                    </template>
                                                </v-checkbox>
                                            </template>
                                        </v-list-item>
                                    </template>
                                </v-list>
                            </template>
                        </v-list>
                    </v-expansion-panel-text>
                </v-expansion-panel>
            </v-expansion-panels>
        </v-card-text>

        <v-card-text class="overflow-y-visible" v-if="dialogMode">
            <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                <v-btn :disabled="!hasAnyVisibleCategory" @click="save">{{ $t(applyText) }}</v-btn>
                <v-btn color="secondary" variant="tonal" @click="cancel">{{ $t('Cancel') }}</v-btn>
            </div>
        </v-card-text>
    </v-card>

    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useStatisticsStore } from '@/stores/statistics.js';

import categoryConstants from '@/consts/category.js';
import { copyObjectTo, arrayItemToObjectField } from '@/lib/common.js';
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
} from '@/lib/category.js';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiDotsVertical
} from '@mdi/js';

export default {
    props: [
        'dialogMode',
        'categoryTypes',
        'type',
        'autoSave'
    ],
    emits: [
        'settings:change'
    ],
    data: function () {
        return {
            loading: true,
            expandCategoryTypes: [
                categoryConstants.allCategoryTypes.Income.toString(),
                categoryConstants.allCategoryTypes.Expense.toString(),
                categoryConstants.allCategoryTypes.Transfer.toString()
            ],
            filterCategoryIds: {},
            showHidden: false,
            icons: {
                selectAll: mdiSelectAll,
                selectNone: mdiSelect,
                selectInverse: mdiSelectInverse,
                show: mdiEyeOutline,
                hide: mdiEyeOffOutline,
                more: mdiDotsVertical
            }
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
        allowCategoryTypes() {
            return this.categoryTypes ? arrayItemToObjectField(this.categoryTypes.split(','), true) : null;
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

        self.transactionCategoriesStore.loadAllCategories({
            force: false
        }).then(() => {
            self.loading = false;

            const allCategoryIds = {};

            for (let categoryId in self.transactionCategoriesStore.allTransactionCategoriesMap) {
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
                for (let categoryId in self.transactionsStore.allFilterCategoryIds) {
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

            const filteredCategoryIds = {};
            let isAllSelected = true;
            let finalCategoryIds = '';
            let changed = true;

            for (let categoryId in self.filterCategoryIds) {
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
                changed = self.statisticsStore.updateTransactionStatisticsFilter({
                    filterCategoryIds: filteredCategoryIds
                });

                if (changed) {
                    self.statisticsStore.updateTransactionStatisticsInvalidState(true);
                }
            } else if (this.type === 'transactionListCurrent') {
                changed = self.transactionsStore.updateTransactionListFilter({
                    categoryIds: isAllSelected ? '' : finalCategoryIds
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
        selectCategory(category, value) {
            if (!category) {
                return;
            }

            this.filterCategoryIds[category.id] = !value;

            if (this.autoSave) {
                this.save();
            }
        },
        selectSubCategories(category, value) {
            selectSubCategories(this.filterCategoryIds, category, !value);

            if (this.autoSave) {
                this.save();
            }
        },
        selectAll() {
            selectAll(this.filterCategoryIds, this.transactionCategoriesStore.allTransactionCategoriesMap);

            if (this.autoSave) {
                this.save();
            }
        },
        selectNone() {
            selectNone(this.filterCategoryIds, this.transactionCategoriesStore.allTransactionCategoriesMap);

            if (this.autoSave) {
                this.save();
            }
        },
        selectInvert() {
            selectInvert(this.filterCategoryIds, this.transactionCategoriesStore.allTransactionCategoriesMap);

            if (this.autoSave) {
                this.save();
            }
        },
        getCategoryTypeName(categoryType) {
            switch (categoryType) {
                case categoryConstants.allCategoryTypes.Income.toString():
                    return this.$t('Income Categories');
                case categoryConstants.allCategoryTypes.Expense.toString():
                    return this.$t('Expense Categories');
                case categoryConstants.allCategoryTypes.Transfer.toString():
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

<style>
.category-types .v-expansion-panel-text__wrapper {
    padding: 0 0 0 0;
}

.category-types .v-expansion-panel--active:not(:first-child),
.category-types .v-expansion-panel--active + .v-expansion-panel {
    margin-top: 1rem;
}
</style>

