<template>
    <v-card>
        <v-toolbar color="primary" v-if="dialogMode">
            <v-toolbar-title>{{ $t('Default Transaction Category Filter') }}</v-toolbar-title>
            <v-spacer/>
            <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                   :disabled="loading || !hasAnyAvailableCategory" :icon="true">
                <v-icon :icon="icons.more" />
                <v-menu activator="parent">
                    <v-list>
                        <v-list-item :prepend-icon="icons.selectAll"
                                     :title="$t('Select All')"
                                     @click="selectAll"></v-list-item>
                        <v-list-item :prepend-icon="icons.selectNone"
                                     :title="$t('Select None')"
                                     @click="selectNone"></v-list-item>
                        <v-list-item :prepend-icon="icons.selectInverse"
                                     :title="$t('Invert Selection')"
                                     @click="selectInvert"></v-list-item>
                    </v-list>
                </v-menu>
            </v-btn>
        </v-toolbar>

        <template #title v-if="!dialogMode">
            <div class="d-flex align-center">
                <span>{{ $t('Default Transaction Category Filter') }}</span>
                <v-spacer/>
                <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                       :disabled="loading" :icon="true">
                    <v-icon :icon="icons.more" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="icons.selectAll"
                                         :title="$t('Select All')"
                                         @click="selectAll"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectNone"
                                         :title="$t('Select None')"
                                         @click="selectNone"></v-list-item>
                            <v-list-item :prepend-icon="icons.selectInverse"
                                         :title="$t('Invert Selection')"
                                         @click="selectInvert"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </div>
        </template>

        <div v-if="loading">
            <v-skeleton-loader type="paragraph" :loading="loading"
                               :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
        </div>

        <v-card-text v-else-if="!loading">
            <v-expansion-panels class="category-types" multiple v-model="expandCategoryTypes">
                <v-expansion-panel :key="transactionType.type"
                                   :value="transactionType.type"
                                   class="border"
                                   v-for="transactionType in allVisibleTransactionCategories">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ml-3">{{ getCategoryTypeName(transactionType.type) }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                        <v-list rounded density="comfortable" class="pa-0">
                            <div class="py-3" v-if="!hasAvailableCategory[transactionType.type]">{{ $t('No available category') }}</div>

                            <template :key="category.id"
                                      v-for="(category, idx) in transactionType.visibleCategories">
                                <v-list-item>
                                    <template #prepend>
                                        <v-checkbox :model-value="isSubCategoriesAllChecked(category, filterCategoryIds)"
                                                    :indeterminate="isSubCategoriesHasButNotAllChecked(category, filterCategoryIds)"
                                                    @update:model-value="selectSubCategories(category, $event)">
                                            <template #label>
                                                <ItemIcon class="d-flex" icon-type="category"
                                                          :icon-id="category.icon" :color="category.color"></ItemIcon>
                                                <span class="ml-3">{{ category.name }}</span>
                                            </template>
                                        </v-checkbox>
                                    </template>
                                </v-list-item>

                                <v-divider v-if="transactionType.visibleSubCategories[category.id]"/>

                                <v-list rounded density="comfortable" class="pa-0 ml-4"
                                        v-if="transactionType.visibleSubCategories[category.id]">
                                    <template :key="subCategory.id"
                                              v-for="(subCategory, subIdx) in transactionType.visibleSubCategories[category.id]">
                                        <v-list-item>
                                            <template #prepend>
                                                <v-checkbox :model-value="isCategoryChecked(subCategory, filterCategoryIds)"
                                                            @update:model-value="selectCategory(subCategory, $event)">
                                                    <template #label>
                                                        <ItemIcon class="d-flex" icon-type="category"
                                                                  :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                                        <span class="ml-3">{{ subCategory.name }}</span>
                                                    </template>
                                                </v-checkbox>
                                            </template>
                                        </v-list-item>
                                        <v-divider v-if="subIdx !== transactionType.visibleSubCategories[category.id].length - 1"/>
                                    </template>
                                </v-list>

                                <v-divider v-if="idx !== transactionType.visibleCategories - 1"/>
                            </template>
                        </v-list>
                    </v-expansion-panel-text>
                </v-expansion-panel>
            </v-expansion-panels>
        </v-card-text>

        <v-card-actions class="mt-3" v-if="dialogMode">
            <v-spacer></v-spacer>
            <v-btn color="gray" @click="cancel">{{ $t('Cancel') }}</v-btn>
            <v-btn color="primary" :disabled="!hasAnyAvailableCategory" @click="save">{{ $t('OK') }}</v-btn>
        </v-card-actions>
    </v-card>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useStatisticsStore } from '@/stores/statistics.js';

import categoryConstants from '@/consts/category.js';
import { copyObjectTo } from '@/lib/common.js';
import { allVisibleTransactionCategories } from '@/lib/category.js';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiDotsVertical
} from '@mdi/js';

export default {
    props: [
        'dialogMode',
        'modifyDefault',
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
            icons: {
                selectAll: mdiSelectAll,
                selectNone: mdiSelect,
                selectInverse: mdiSelectInverse,
                more: mdiDotsVertical
            }
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useTransactionCategoriesStore, useStatisticsStore),
        title() {
            if (this.modifyDefault) {
                return 'Default Transaction Category Filter';
            } else {
                return 'Filter Transaction Categories';
            }
        },
        applyText() {
            if (this.modifyDefault) {
                return 'Save';
            } else {
                return 'Apply';
            }
        },
        allVisibleTransactionCategories() {
            return allVisibleTransactionCategories(this.transactionCategoriesStore.allTransactionCategories);
        },
        hasAnyAvailableCategory() {
            for (let type in this.allVisibleTransactionCategories) {
                if (!Object.prototype.hasOwnProperty.call(this.allVisibleTransactionCategories, type)) {
                    continue;
                }

                const categoryType = this.allVisibleTransactionCategories[type];

                if (categoryType.visibleCategories && categoryType.visibleCategories.length > 0) {
                    return true;
                }
            }

            return false;
        },
        hasAvailableCategory() {
            const result = {};

            for (let type in this.allVisibleTransactionCategories) {
                if (!Object.prototype.hasOwnProperty.call(this.allVisibleTransactionCategories, type)) {
                    continue;
                }

                const categoryType = this.allVisibleTransactionCategories[type];
                result[type] = categoryType.visibleCategories && categoryType.visibleCategories.length > 0;
            }

            return result;
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
                allCategoryIds[category.id] = false;
            }

            if (self.modifyDefault) {
                self.filterCategoryIds = copyObjectTo(self.settingsStore.appSettings.statistics.defaultTransactionCategoryFilter, allCategoryIds);
            } else {
                self.filterCategoryIds = copyObjectTo(self.statisticsStore.transactionStatisticsFilter.filterCategoryIds, allCategoryIds);
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

            for (let categoryId in self.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(self.filterCategoryIds, categoryId)) {
                    continue;
                }

                if (self.filterCategoryIds[categoryId]) {
                    filteredCategoryIds[categoryId] = true;
                }
            }

            if (self.modifyDefault) {
                self.settingsStore.setStatisticsDefaultTransactionCategoryFilter(filteredCategoryIds);
            } else {
                self.statisticsStore.updateTransactionStatisticsFilter({
                    filterCategoryIds: filteredCategoryIds
                });
            }

            this.$emit('settings:change', true);
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
            if (!category || !category.subCategories || !category.subCategories.length) {
                return;
            }

            for (let i = 0; i < category.subCategories.length; i++) {
                const subCategory = category.subCategories[i];
                this.filterCategoryIds[subCategory.id] = !value;
            }

            if (this.autoSave) {
                this.save();
            }
        },
        selectAll() {
            for (let categoryId in this.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterCategoryIds, categoryId)) {
                    continue;
                }

                const category = this.transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                if (category) {
                    this.filterCategoryIds[category.id] = false;
                }
            }

            if (this.autoSave) {
                this.save();
            }
        },
        selectNone() {
            for (let categoryId in this.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterCategoryIds, categoryId)) {
                    continue;
                }

                const category = this.transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                if (category) {
                    this.filterCategoryIds[category.id] = true;
                }
            }

            if (this.autoSave) {
                this.save();
            }
        },
        selectInvert() {
            for (let categoryId in this.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterCategoryIds, categoryId)) {
                    continue;
                }

                const category = this.transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                if (category) {
                    this.filterCategoryIds[category.id] = !this.filterCategoryIds[category.id];
                }
            }

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
            for (let i = 0; i < category.subCategories.length; i++) {
                const subCategory = category.subCategories[i];
                if (filterCategoryIds[subCategory.id]) {
                    return false;
                }
            }

            return true;
        },
        isSubCategoriesHasButNotAllChecked(category, filterCategoryIds) {
            let checkedCount = 0;

            for (let i = 0; i < category.subCategories.length; i++) {
                const subCategory = category.subCategories[i];
                if (!filterCategoryIds[subCategory.id]) {
                    checkedCount++;
                }
            }

            return checkedCount > 0 && checkedCount < category.subCategories.length;
        }
    }
}
</script>

<style>
.category-types .v-expansion-panel-text__wrapper {
    padding: 0 0 0 20px;
}
</style>

