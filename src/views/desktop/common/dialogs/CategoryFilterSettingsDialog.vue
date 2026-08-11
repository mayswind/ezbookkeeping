<template>
    <v-dialog width="800" v-model="showState">
        <one-column-dialog-layout content-class="pt-4"
                                  :title="tt(title)" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-text-field class="mx-2" density="compact"
                              :disabled="loading || !hasAnyAvailableCategory"
                              :prepend-inner-icon="mdiMagnify"
                              :placeholder="tt('Find category')"
                              v-model="filterContent"></v-text-field>

                <v-btn density="comfortable" variant="outlined"
                       :disabled="!hasAnyAvailableCategory" @click="save">{{ tt(applyText) }}</v-btn>

                <v-btn density="compact" color="default" variant="text" class="ms-2"
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
                            <v-list-item :prepend-icon="mdiEyeOutline"
                                         :title="tt('Show Hidden Transaction Categories')"
                                         v-if="!showHidden" @click="showHidden = true"></v-list-item>
                            <v-list-item :prepend-icon="mdiEyeOffOutline"
                                         :title="tt('Hide Hidden Transaction Categories')"
                                         v-if="showHidden" @click="showHidden = false"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </template>

            <template #content>
                <div v-if="loading">
                    <v-skeleton-loader type="paragraph" :loading="loading"
                                       :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]"></v-skeleton-loader>
                </div>

                <div v-else-if="!loading">
                    <v-expansion-panels class="category-types" multiple v-model="expandCategoryTypes">
                        <v-expansion-panel :key="categoryType"
                                           :value="parseInt(categoryType) as CategoryType"
                                           class="border"
                                           v-for="(categories, categoryType) in allVisibleTransactionCategories">
                            <v-expansion-panel-title class="expand-panel-title-with-bg py-0">
                                <span class="text-body-large">{{ getCategoryTypeName(parseInt(categoryType)) }}</span>
                            </v-expansion-panel-title>
                            <v-expansion-panel-text>
                                <v-list rounded density="comfortable" class="pa-0">
                                    <div class="ms-5 py-3" v-if="!categories || !categories.length">{{ tt('No available category') }}</div>

                                    <template :key="category.id"
                                              v-for="(category, idx) in categories">
                                        <v-divider v-if="idx > 0"/>

                                        <v-list-item>
                                            <template #prepend>
                                                <v-checkbox :model-value="isSubCategoriesAllChecked(category, filterCategoryIds)"
                                                            :indeterminate="isSubCategoriesHasButNotAllChecked(category, filterCategoryIds)"
                                                            @update:model-value="updateAllSubCategoriesSelected(category, $event)">
                                                    <template #label>
                                                        <ItemIcon class="d-flex ms-1" icon-type="category" :icon-id="category.icon"
                                                                  :color="category.color" :hidden-status="category.hidden"></ItemIcon>
                                                        <span class="text-body-medium ms-2">{{ category.name }}</span>
                                                    </template>
                                                </v-checkbox>
                                            </template>
                                        </v-list-item>

                                        <v-divider v-if="category.subCategories && category.subCategories.length"/>

                                        <v-list rounded density="comfortable" class="pa-0 ms-4"
                                                v-if="category.subCategories && category.subCategories.length">
                                            <template :key="subCategory.id"
                                                      v-for="(subCategory, subIdx) in category.subCategories">
                                                <v-divider v-if="subIdx > 0"/>

                                                <v-list-item>
                                                    <template #prepend>
                                                        <v-checkbox :model-value="isCategoryChecked(subCategory, filterCategoryIds)"
                                                                    @update:model-value="updateCategorySelected(subCategory, $event)">
                                                            <template #label>
                                                                <ItemIcon class="d-flex ms-1" icon-type="category" :icon-id="subCategory.icon"
                                                                          :color="subCategory.color" :hidden-status="subCategory.hidden"></ItemIcon>
                                                                <span class="text-body-medium ms-2">{{ subCategory.name }}</span>
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
                </div>
            </template>
        </one-column-dialog-layout>

        <snack-bar ref="snackbar" />
    </v-dialog>
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

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
    selectAll,
    selectNone,
    selectInvert,
    isSubCategoriesAllChecked,
    isSubCategoriesHasButNotAllChecked
} from '@/lib/category.ts';

import {
    mdiMagnify,
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
    selectedCategoryIds?: string[];
    autoSave?: boolean;
    categoryTypes?: string;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'settings:change', changed: boolean, selectedCategoryIds?: string[]): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const {
    loading,
    showHidden,
    filterContent,
    filterCategoryIds,
    title,
    applyText,
    allVisibleTransactionCategories,
    allVisibleTransactionCategoryMap,
    hasAnyAvailableCategory,
    hasAnyVisibleCategory,
    isCategoryChecked,
    getCategoryTypeName,
    loadFilterCategoryIds,
    saveFilterCategoryIds
} = useCategoryFilterSettingPageBase(props.type, props.categoryTypes, props.selectedCategoryIds);

const transactionCategoriesStore = useTransactionCategoriesStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const expandCategoryTypes = ref<CategoryType[]>([
    CategoryType.Income,
    CategoryType.Expense,
    CategoryType.Transfer
]);

const showState = computed<boolean>({
    get: () => props.show || false,
    set: (value) => emit('update:show', value)
});

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
    selectAll(filterCategoryIds.value, allVisibleTransactionCategoryMap.value);

    if (props.autoSave) {
        save();
    }
}

function selectNoneCategories(): void {
    selectNone(filterCategoryIds.value, allVisibleTransactionCategoryMap.value);

    if (props.autoSave) {
        save();
    }
}

function selectInvertCategories(): void {
    selectInvert(filterCategoryIds.value, allVisibleTransactionCategoryMap.value);

    if (props.autoSave) {
        save();
    }
}

function save(): void {
    const [changed, selectedCategoryIds] = saveFilterCategoryIds();
    emit('settings:change', changed, selectedCategoryIds);
}

function cancel(): void {
    emit('settings:change', false);
}

watch(() => props.show, (newValue) => {
    if (newValue) {
        loadFilterCategoryIds();
        showHidden.value = false;
        filterContent.value = '';
    }
});

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

