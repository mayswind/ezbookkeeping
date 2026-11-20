<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <template #title>
                    <div class="title-and-toolbar d-flex align-center">
                        <span>{{ tt('Transaction Categories') }}</span>
                        <v-btn class="ms-3" color="default" variant="outlined"
                               :disabled="loading || updating" @click="add">{{ tt('Add') }}</v-btn>
                        <v-btn class="ms-3" color="primary" variant="tonal"
                               :disabled="loading || updating" @click="saveSortResult"
                               v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                        <v-btn density="compact" color="default" variant="text" class="ms-2"
                               :aria-label="tt('Refresh')" :icon="true" :loading="loading || updating" @click="reload(true)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="mdiRefresh" size="24" />
                            <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                        </v-btn>
                        <v-spacer/>
                        <v-select class="ms-2 flex-0-1" density="compact" variant="outlined"
                                  item-title="name" item-value="value"
                                  :disabled="loading"
                                  :items="availableRootCategories"
                                  v-model="primaryCategoryId">
                        </v-select>
                        <v-btn-toggle class="toggle-buttons ms-3" style="min-height: 36px"
                                      density="compact" color="default" variant="outlined"
                                      mandatory="force"
                                      :disabled="loading"
                                      v-model="activeCategoryType"
                                      @update:model-value="switchAllPrimaryCategories">
                            <v-btn :value="CategoryType.Expense">{{ tt('Expense') }}</v-btn>
                            <v-btn :value="CategoryType.Income">{{ tt('Income') }}</v-btn>
                            <v-btn :value="CategoryType.Transfer">{{ tt('Transfer') }}</v-btn>
                        </v-btn-toggle>
                        <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                               :aria-label="tt('More')" :disabled="loading || updating" :icon="true">
                            <v-icon :icon="mdiDotsVertical" />
                            <v-menu activator="parent">
                                <v-list>
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

                <v-table class="transaction-category-table table-striped" :hover="!loading">
                    <thead>
                    <tr>
                        <th>
                            <div class="d-flex align-center">
                                <span>{{ tt('Category Name') }}</span>
                                <v-spacer/>
                                <span>{{ tt('Operation') }}</span>
                            </div>
                        </th>
                    </tr>
                    </thead>

                    <tbody v-if="loading && noAvailableCategory">
                    <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]">
                        <td class="px-0">
                            <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                        </td>
                    </tr>
                    </tbody>

                    <tbody v-if="!loading && noAvailableCategory">
                    <tr>
                        <td>
                            <div class="d-flex align-center">
                                <span>{{ tt('No available category') }}</span>
                                <v-btn class="ms-3" color="default" variant="outlined"
                                       @click="showPresetDialog = true"
                                       v-if="hasSubCategories && noCategory">
                                    {{ tt('Add Default Categories') }}
                                </v-btn>
                            </div>
                        </td>
                    </tr>
                    </tbody>

                    <draggable-list tag="tbody"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    :disabled="noAvailableCategory"
                                    v-model="categories"
                                    @change="onMove">
                        <template #item="{ element }">
                            <tr class="transaction-category-table-row" v-if="showHidden || !element.hidden"
                                @mouseenter="hoveredCategoryId = element.id" @mouseleave="hoveredCategoryId = ''">
                                <td>
                                    <div class="d-flex align-center">
                                        <div class="d-flex align-center" :class="{ 'cursor-pointer': isCategorySupportSwitch(element) }"
                                             @click="switchPrimaryCategory(element)">
                                            <ItemIcon :icon-type="getCategoryIconType(element.iconType)"
                                                      :icon-id="element.icon" :color="element.color"
                                                      :hidden-status="element.hidden" />
                                            <div class="d-flex flex-column py-2">
                                                <span class="ms-2">{{ element.name }}</span>
                                                <span class="transaction-category-comment ms-2">{{ element.comment }}</span>
                                            </div>
                                        </div>

                                        <v-spacer/>

                                        <template v-if="hoveredCategoryId === element.id && !loading">
                                            <v-btn class="px-2 ms-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="element.hidden ? mdiEyeOutline : mdiEyeOffOutline"
                                                   :loading="categoryHiding[element.id]"
                                                   :disabled="loading || updating"
                                                   @click="hide(element, !element.hidden)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ element.hidden ? tt('Show') : tt('Hide') }}
                                            </v-btn>
                                            <v-btn class="px-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="mdiPencilOutline"
                                                   :disabled="loading || updating"
                                                   @click="edit(element)">
                                                {{ tt('Edit') }}
                                            </v-btn>
                                            <v-btn class="px-2" color="default"
                                                   density="comfortable" variant="text"
                                                   :prepend-icon="mdiDeleteOutline"
                                                   :loading="categoryRemoving[element.id]"
                                                   :disabled="loading || updating"
                                                   @click="remove(element)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                </template>
                                                {{ tt('Delete') }}
                                            </v-btn>
                                        </template>

                                        <span class="ms-2">
                                            <v-icon :class="!loading && !updating && availableCategoryCount > 1 ? 'drag-handle' : 'disabled'"
                                                    :aria-label="tt('Drag to Reorder')" :icon="mdiDrag"/>
                                            <v-tooltip activator="parent" v-if="!loading && !updating && availableCategoryCount > 1 && hoveredCategoryId === element.id">{{ tt('Drag to Reorder') }}</v-tooltip>
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

    <edit-dialog ref="editDialog" />
    <preset-dialog :category-type="activeCategoryType" v-model:show="showPresetDialog"
                   @category:saved="onPresetCategorySaved" />

    <confirm-dialog ref="confirmDialog"/>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import EditDialog from './list/dialogs/EditDialog.vue';
import PresetDialog from './list/dialogs/PresetDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useCategoryListPageBase } from '@/views/base/categories/CategoryListPageBase.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import type { NameValue } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

import { getCategoryIconType } from '@/lib/icon.ts';
import {
    isNoAvailableCategory,
    getAvailableCategoryCount
} from '@/lib/category.ts';

import {
    mdiRefresh,
    mdiPencilOutline,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type EditDialogType = InstanceType<typeof EditDialog>;

const { tt } = useI18n();
const { loading, primaryCategoryId, currentPrimaryCategory } = useCategoryListPageBase();

const transactionCategoriesStore = useTransactionCategoriesStore();
const desktopPageStore = useDesktopPageStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const activeCategoryType = ref<CategoryType>(CategoryType.Expense);
const updating = ref<boolean>(false);
const hoveredCategoryId = ref<string>('');
const categoryHiding = ref<Record<string, boolean>>({});
const categoryRemoving = ref<Record<string, boolean>>({});
const displayOrderModified = ref<boolean>(false);
const showHidden = ref<boolean>(false);
const showPresetDialog = ref<boolean>(false);

const availableRootCategories = computed<NameValue[]>(() => {
    const categories: NameValue[] = [];

    categories.push({
        name: tt('Primary Categories'),
        value: '0'
    });

    categories.push(...primaryCategories.value.filter(category => !category.hidden).map(category => ({
        name: category.name,
        value: category.id
    })));

    return categories;
});

const primaryCategories = computed<TransactionCategory[]>(() => {
    if (!transactionCategoriesStore.allTransactionCategories || !transactionCategoriesStore.allTransactionCategories[activeCategoryType.value]) {
        return [];
    }

    return transactionCategoriesStore.allTransactionCategories[activeCategoryType.value] ?? [];
});

const secondaryCategories = computed<TransactionCategory[]>(() => {
    if (!transactionCategoriesStore.allTransactionCategoriesMap || !transactionCategoriesStore.allTransactionCategoriesMap[primaryCategoryId.value]) {
        return [];
    }

    return transactionCategoriesStore.allTransactionCategoriesMap[primaryCategoryId.value]?.subCategories ?? [];
});

const hasSubCategories = computed<boolean>(() => {
    return !primaryCategoryId.value || primaryCategoryId.value === '' || primaryCategoryId.value === '0';
});

const categories = computed<TransactionCategory[]>(() => {
    if (hasSubCategories.value) {
        return primaryCategories.value;
    } else {
        return secondaryCategories.value;
    }
});

const noAvailableCategory = computed<boolean>(() => isNoAvailableCategory(categories.value, showHidden.value));
const noCategory = computed<boolean>(() => categories.value.length < 1);
const availableCategoryCount = computed<number>(() => getAvailableCategoryCount(categories.value, showHidden.value));

function isCategorySupportSwitch(category: TransactionCategory): boolean {
    if (!category || category.hidden) {
        return false;
    }

    return !category.parentId || category.parentId === '' || category.parentId === '0';
}

function switchAllPrimaryCategories(): void {
    primaryCategoryId.value = '0';
}

function switchPrimaryCategory(category: TransactionCategory): void {
    if (!category || category.hidden) {
        return;
    }

    if (!category.parentId || category.parentId === '' || category.parentId === '0') {
        primaryCategoryId.value = category.id;
    }
}

function reload(force: boolean): void {
    loading.value = true;

    transactionCategoriesStore.loadAllCategories({
        force: force
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        if (force) {
            desktopPageStore.addToast('Category list has been updated');
        }
    }).catch(error => {
        loading.value = false;

        if (error && error.isUpToDate) {
            displayOrderModified.value = false;
        }

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function add(): void {
    editDialog.value?.open({
        type: activeCategoryType.value,
        parentId: primaryCategoryId.value,
        color: currentPrimaryCategory.value?.color,
        icon: currentPrimaryCategory.value?.icon
    }).then(result => {
        if (result && result.message) {
            desktopPageStore.addToast(result.message);
        }
    }).catch(error => {
        if (error) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function edit(category: TransactionCategory): void {
    editDialog.value?.open({
        id: category.id,
        currentCategory: category
    }).then(result => {
        if (result && result.message) {
            desktopPageStore.addToast(result.message);
        }

        if (transactionCategoriesStore.transactionCategoryListStateInvalid) {
            reload(true);
        }
    }).catch(error => {
        if (error) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function hide(category: TransactionCategory, hidden: boolean): void {
    updating.value = true;
    categoryHiding.value[category.id] = true;

    transactionCategoriesStore.hideCategory({
        category: category,
        hidden: hidden
    }).then(() => {
        updating.value = false;
        categoryHiding.value[category.id] = false;
    }).catch(error => {
        updating.value = false;
        categoryHiding.value[category.id] = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function remove(category: TransactionCategory): void {
    confirmDialog.value?.open('Are you sure you want to delete this category?').then(() => {
        updating.value = true;
        categoryRemoving.value[category.id] = true;

        transactionCategoriesStore.deleteCategory({
            category: category
        }).then(() => {
            updating.value = false;
            categoryRemoving.value[category.id] = false;
        }).catch(error => {
            updating.value = false;
            categoryRemoving.value[category.id] = false;

            if (!error.processed) {
                desktopPageStore.addErrorToast(error);
            }
        });
    });
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        return;
    }

    loading.value = true;

    transactionCategoriesStore.updateCategoryDisplayOrders({
        type: activeCategoryType.value,
        parentId: primaryCategoryId.value
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function onMove(event: { moved: { element: { id: string }, oldIndex: number, newIndex: number } }): void {
    if (!event || !event.moved) {
        return;
    }

    const moveEvent = event.moved;

    if (!moveEvent.element || !moveEvent.element.id) {
        desktopPageStore.addToast('Unable to move category');
        return;
    }

    transactionCategoriesStore.changeCategoryDisplayOrder({
        categoryId: moveEvent.element.id,
        from: moveEvent.oldIndex,
        to: moveEvent.newIndex
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        desktopPageStore.addErrorToast(error);
    });
}

function onPresetCategorySaved(e: { message: string }): void {
    if (e && e.message) {
        desktopPageStore.addToast(e.message);
        reload(false);
    }
}

reload(false);
</script>

<style>
.transaction-category-table .transaction-category-comment {
    font-size: 0.8rem;
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
}
</style>
