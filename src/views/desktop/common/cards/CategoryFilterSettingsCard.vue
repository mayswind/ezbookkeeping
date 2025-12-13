<template>
    <v-card :class="{ 'pa-sm-1 pa-md-2': dialogMode }">
        <template #title>
            <div class="d-flex align-center justify-center" v-if="dialogMode">
                <div class="w-100 text-center">
                    <h4 class="text-h4">{{ tt(title) }}</h4>
                </div>
                <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                       :disabled="loading || !hasAnyAvailableCategory" :icon="true">
                    <v-icon :icon="mdiDotsVertical" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectAllCategories"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelect"
                                         :title="tt('Select None')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectNoneCategories"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelectInverse"
                                         :title="tt('Invert Selection')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectInvertCategories"></v-list-item>
                            <v-divider class="my-2"/>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All Visible')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectAllVisibleCategories"></v-list-item>
                            <v-divider class="my-2"/>
                            <v-list-item :prepend-icon="mdiEyeOutline"
                                         :title="tt('Show Hidden Transaction Categories')"
                                         v-if="!showHidden" @click="showHidden = true"></v-list-item>
                            <v-list-item :prepend-icon="mdiEyeOffOutline"
                                         :title="tt('Hide Hidden Transaction Categories')"
                                         v-if="showHidden" @click="showHidden = false"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </div>
            <div class="d-flex align-center" v-else-if="!dialogMode">
                <span>{{ tt(title) }}</span>
                <v-spacer/>
                <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                       :disabled="loading" :icon="true">
                    <v-icon :icon="mdiDotsVertical" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectAllCategories"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelect"
                                         :title="tt('Select None')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectNoneCategories"></v-list-item>
                            <v-list-item :prepend-icon="mdiSelectInverse"
                                         :title="tt('Invert Selection')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectInvertCategories"></v-list-item>
                            <v-divider class="my-2"/>
                            <v-list-item :prepend-icon="mdiSelectAll"
                                         :title="tt('Select All Visible')"
                                         :disabled="!hasAnyVisibleCategory"
                                         @click="selectAllVisibleCategories"></v-list-item>
                            <v-divider class="my-2"/>
                            <v-list-item :prepend-icon="mdiEyeOutline"
                                         :title="tt('Show Hidden Transaction Categories')"
                                         v-if="!showHidden" @click="showHidden = true"></v-list-item>
                            <v-list-item :prepend-icon="mdiEyeOffOutline"
                                         :title="tt('Hide Hidden Transaction Categories')"
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

        <v-card-text :class="{ 'flex-grow-1 overflow-y-auto': dialogMode }" v-else-if="!loading">
            <v-expansion-panels class="category-types" multiple v-model="expandCategoryTypes">
                <v-expansion-panel :key="transactionType.type"
                                   :value="transactionType.type"
                                   class="border"
                                   v-for="transactionType in allTransactionCategories">
                    <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                        <span class="ms-3">{{ getCategoryTypeName(transactionType.type) }}</span>
                    </v-expansion-panel-title>
                    <v-expansion-panel-text>
                        <v-list rounded density="comfortable" class="pa-0">
                            <div class="ms-5 py-3" v-if="!hasAvailableCategory[transactionType.type]">{{ tt('No available category') }}</div>

                            <template :key="category.id"
                                      v-for="(category, idx) in transactionType.allCategories">
                                <v-divider v-if="showHidden ? idx > 0 : (!category.hidden ? idx > transactionType.firstVisibleCategoryIndex : false)"/>

                                <v-list-item v-if="showHidden || !category.hidden">
                                    <template #prepend>
                                        <v-checkbox :model-value="isSubCategoriesAllChecked(category, filterCategoryIds)"
                                                    :indeterminate="isSubCategoriesHasButNotAllChecked(category, filterCategoryIds)"
                                                    @update:model-value="updateAllSubCategoriesSelected(category, $event)">
                                            <template #label>
                                                <ItemIcon class="d-flex" icon-type="category" :icon-id="category.icon"
                                                          :color="category.color" :hidden-status="category.hidden"></ItemIcon>
                                                <span class="ms-3">{{ category.name }}</span>
                                            </template>
                                        </v-checkbox>
                                    </template>
                                </v-list-item>

                                <v-divider v-if="(showHidden || !category.hidden) && ((showHidden && transactionType.allSubCategories[category.id]) || transactionType.allVisibleSubCategoryCounts[category.id])"/>

                                <v-list rounded density="comfortable" class="pa-0 ms-4"
                                        v-if="(showHidden || !category.hidden) && ((showHidden && transactionType.allSubCategories[category.id]) || transactionType.allVisibleSubCategoryCounts[category.id])">
                                    <template :key="subCategory.id"
                                              v-for="(subCategory, subIdx) in transactionType.allSubCategories[category.id]">
                                        <v-divider v-if="showHidden ? subIdx > 0 : (!subCategory.hidden ? subIdx > (transactionType.allFirstVisibleSubCategoryIndexes[category.id] as number) : false)"/>

                                        <v-list-item v-if="showHidden || !subCategory.hidden">
                                            <template #prepend>
                                                <v-checkbox :model-value="isCategoryChecked(subCategory, filterCategoryIds)"
                                                            @update:model-value="updateCategorySelected(subCategory, $event)">
                                                    <template #label>
                                                        <ItemIcon class="d-flex" icon-type="category" :icon-id="subCategory.icon"
                                                                  :color="subCategory.color" :hidden-status="subCategory.hidden"></ItemIcon>
                                                        <span class="ms-3">{{ subCategory.name }}</span>
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
            <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                <v-btn :disabled="!hasAnyVisibleCategory" @click="save">{{ tt(applyText) }}</v-btn>
                <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
            </div>
        </v-card-text>
    </v-card>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type CategoryFilterType,
    useCategoryFilterSettingPageBase
} from '@/views/base/settings/CategoryFilterSettingPageBase.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { CategoryType } from '@/core/category.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

import {
    selectAllSubCategories,
    selectAllVisible,
    selectAll,
    selectNone,
    selectInvert,
    isSubCategoriesAllChecked,
    isSubCategoriesHasButNotAllChecked
} from '@/lib/category.ts';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiDotsVertical
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    type: CategoryFilterType;
    dialogMode?: boolean;
    autoSave?: boolean;
    categoryTypes?: string;
}>();

const emit = defineEmits<{
    (e: 'settings:change', changed: boolean): void;
}>();

const { tt } = useI18n();

const {
    loading,
    showHidden,
    filterCategoryIds,
    title,
    applyText,
    allTransactionCategories,
    hasAnyAvailableCategory,
    hasAnyVisibleCategory,
    hasAvailableCategory,
    isCategoryChecked,
    getCategoryTypeName,
    loadFilterCategoryIds,
    saveFilterCategoryIds
} = useCategoryFilterSettingPageBase(props.type, props.categoryTypes);

const transactionCategoriesStore = useTransactionCategoriesStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const expandCategoryTypes = ref<CategoryType[]>([
    CategoryType.Income,
    CategoryType.Expense,
    CategoryType.Transfer
]);

function init(): void {
    transactionCategoriesStore.loadAllCategories({
        force: false
    }).then(() => {
        loading.value = false;

        if (!loadFilterCategoryIds()) {
            snackbar.value?.showError('Parameter Invalid');
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function updateCategorySelected(category: TransactionCategory, value: boolean | null): void {
    if (!category) {
        return;
    }

    filterCategoryIds.value[category.id] = !value;

    if (props.autoSave) {
        save();
    }
}

function updateAllSubCategoriesSelected(category: TransactionCategory, value: boolean | null): void {
    selectAllSubCategories(filterCategoryIds.value, !value, category);

    if (props.autoSave) {
        save();
    }
}

function selectAllCategories(): void {
    selectAll(filterCategoryIds.value, transactionCategoriesStore.allTransactionCategoriesMap);

    if (props.autoSave) {
        save();
    }
}

function selectNoneCategories(): void {
    selectNone(filterCategoryIds.value, transactionCategoriesStore.allTransactionCategoriesMap);

    if (props.autoSave) {
        save();
    }
}

function selectInvertCategories(): void {
    selectInvert(filterCategoryIds.value, transactionCategoriesStore.allTransactionCategoriesMap);

    if (props.autoSave) {
        save();
    }
}

function selectAllVisibleCategories(): void {
    selectAllVisible(filterCategoryIds.value, transactionCategoriesStore.allTransactionCategoriesMap);

    if (props.autoSave) {
        save();
    }
}

function save(): void {
    const changed = saveFilterCategoryIds();
    emit('settings:change', changed);
}

function cancel(): void {
    emit('settings:change', false);
}

init();
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

