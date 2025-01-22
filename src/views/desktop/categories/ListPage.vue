<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer ref="navbar" :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="[
                                { name: tt('Expense'), value: CategoryType.Expense },
                                { name: tt('Income'), value: CategoryType.Income },
                                { name: tt('Transfer'), value: CategoryType.Transfer }
                            ]" v-model="activeCategoryType" @update:model-value="switchAllPrimaryCategories" />
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="primaryCategoryId">
                            <v-tab class="tab-text-truncate" value="0" @click="switchAllPrimaryCategories">
                                <span class="text-truncate">{{ tt('Primary Categories') }}</span>
                            </v-tab>
                            <template :key="category.id" v-for="category in primaryCategories">
                                <v-tab class="tab-text-truncate" :value="category.id" v-if="!category.hidden"
                                       @click="switchPrimaryCategory(category)">
                                    <span class="text-truncate">{{ category.name }}</span>
                                </v-tab>
                            </template>
                            <template v-if="loading && (!primaryCategories || primaryCategories.length < 1)">
                                <v-skeleton-loader class="skeleton-no-margin mx-5 mt-4 mb-3" type="text"
                                                   :key="itemIdx" :loading="true" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]"></v-skeleton-loader>
                            </template>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="categoryPage">
                                <v-card variant="flat" :min-height="cardMinHeight">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="mr-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="icons.menu" size="24" />
                                            </v-btn>
                                            <span>{{ tt('Transaction Categories') }}</span>
                                            <v-btn class="ml-3" color="default" variant="outlined"
                                                   :disabled="loading || updating" @click="add">{{ tt('Add') }}</v-btn>
                                            <v-btn class="ml-3" color="primary" variant="tonal"
                                                   :disabled="loading || updating" @click="saveSortResult"
                                                   v-if="displayOrderModified">{{ tt('Save Display Order') }}</v-btn>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ml-2" :icon="true" :loading="loading || updating" @click="reload(true)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="icons.refresh" size="24" />
                                                <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-spacer/>
                                            <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                                                   :disabled="loading || updating" :icon="true">
                                                <v-icon :icon="icons.more" />
                                                <v-menu activator="parent">
                                                    <v-list>
                                                        <v-list-item :prepend-icon="icons.show"
                                                                     :title="tt('Show Hidden Transaction Categories')"
                                                                     v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                                        <v-list-item :prepend-icon="icons.hide"
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
                                                    <v-btn class="ml-3" color="default" variant="outlined"
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
                                                <tr class="transaction-category-table-row text-sm" v-if="showHidden || !element.hidden">
                                                    <td>
                                                        <div class="d-flex align-center">
                                                            <div class="d-flex align-center" :class="{ 'cursor-pointer': isCategorySupportSwitch(element) }"
                                                                 @click="switchPrimaryCategory(element)">
                                                                <ItemIcon icon-type="category"
                                                                          :icon-id="element.icon" :color="element.color"
                                                                          :hidden-status="element.hidden" />
                                                                <div class="d-flex flex-column py-2">
                                                                    <span class="ml-2">{{ element.name }}</span>
                                                                    <span class="transaction-category-comment ml-2">{{ element.comment }}</span>
                                                                </div>
                                                            </div>

                                                            <v-spacer/>

                                                            <v-btn class="px-2 ml-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="element.hidden ? icons.show : icons.hide"
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
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="icons.edit"
                                                                   :disabled="loading || updating"
                                                                   @click="edit(element)">
                                                                {{ tt('Edit') }}
                                                            </v-btn>
                                                            <v-btn class="px-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="icons.remove"
                                                                   :loading="categoryRemoving[element.id]"
                                                                   :disabled="loading || updating"
                                                                   @click="remove(element)">
                                                                <template #loader>
                                                                    <v-progress-circular indeterminate size="20" width="2"/>
                                                                </template>
                                                                {{ tt('Delete') }}
                                                            </v-btn>
                                                            <span class="ml-2">
                                                                <v-icon :class="!loading && !updating && availableCategoryCount > 1 ? 'drag-handle' : 'disabled'"
                                                                        :icon="icons.drag"/>
                                                                <v-tooltip activator="parent" v-if="!loading && !updating && availableCategoryCount > 1">{{ tt('Drag to Reorder') }}</v-tooltip>
                                                            </span>
                                                        </div>
                                                    </td>
                                                </tr>
                                            </template>
                                        </draggable-list>
                                    </v-table>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <edit-dialog ref="editDialog" :persistent="true" />
    <preset-dialog :category-type="activeCategoryType" v-model:show="showPresetDialog"
                            @category:saved="onPresetCategorySaved" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import { VNavigationDrawer } from 'vuetify/components/VNavigationDrawer';

import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import EditDialog from './list/dialogs/EditDialog.vue';
import PresetDialog from './list/dialogs/PresetDialog.vue';

import { ref, computed, useTemplateRef, watch, nextTick } from 'vue';
import { useDisplay } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { CategoryType } from '@/core/category.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

import {
    isNoAvailableCategory,
    getAvailableCategoryCount
} from '@/lib/category.ts';
import { getNavSideBarOuterHeight } from '@/lib/ui/desktop.ts';

import {
    mdiRefresh,
    mdiMenu,
    mdiPencilOutline,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type EditDialogType = InstanceType<typeof EditDialog>;

const display = useDisplay();
const { tt } = useI18n();

const transactionCategoriesStore = useTransactionCategoriesStore();

const icons = {
    refresh: mdiRefresh,
    menu: mdiMenu,
    edit: mdiPencilOutline,
    show: mdiEyeOutline,
    hide: mdiEyeOffOutline,
    remove: mdiDeleteOutline,
    drag: mdiDrag,
    more: mdiDotsVertical
};

const navbar = useTemplateRef<VNavigationDrawer>('navbar');
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<EditDialogType>('editDialog');

const activeCategoryType = ref<CategoryType>(CategoryType.Expense);
const activeTab = ref<string>('categoryPage');
const primaryCategoryId = ref<string>('0');
const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const categoryHiding = ref<Record<string, boolean>>({});
const categoryRemoving = ref<Record<string, boolean>>({});
const displayOrderModified = ref<boolean>(false);
const cardMinHeight = ref<number>(680);
const alwaysShowNav = ref<boolean>(display.mdAndUp.value);
const showNav = ref<boolean>(display.mdAndUp.value);
const showHidden = ref<boolean>(false);
const showPresetDialog = ref<boolean>(false);

const primaryCategories = computed<TransactionCategory[]>(() => {
    if (!transactionCategoriesStore.allTransactionCategories || !transactionCategoriesStore.allTransactionCategories[activeCategoryType.value]) {
        return [];
    }

    return transactionCategoriesStore.allTransactionCategories[activeCategoryType.value];
});

const secondaryCategories = computed<TransactionCategory[]>(() => {
    if (!transactionCategoriesStore.allTransactionCategoriesMap || !transactionCategoriesStore.allTransactionCategoriesMap[primaryCategoryId.value]) {
        return [];
    }

    return transactionCategoriesStore.allTransactionCategoriesMap[primaryCategoryId.value].secondaryCategories || [];
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

function updateCardMinHeight(): void {
    nextTick(() => {
        if (navbar.value && navbar.value.$el && navbar.value.$el.nextElementSibling) {
            const navbarHeight = getNavSideBarOuterHeight(navbar.value.$el.nextElementSibling);
            cardMinHeight.value = Math.max(navbarHeight, 680);
        }
    });
}

function isCategorySupportSwitch(category: TransactionCategory): boolean {
    if (!category || category.hidden) {
        return false;
    }

    return !category.parentId || category.parentId === '' || category.parentId === '0';
}

function switchAllPrimaryCategories(): void {
    primaryCategoryId.value = '0';
    updateCardMinHeight();
}

function switchPrimaryCategory(category: TransactionCategory): void {
    if (!category || category.hidden) {
        return;
    }

    if (!category.parentId || category.parentId === '' || category.parentId === '0') {
        primaryCategoryId.value = category.id;
    }

    updateCardMinHeight();
}

function reload(force: boolean): void {
    loading.value = true;

    transactionCategoriesStore.loadAllCategories({
        force: force
    }).then(() => {
        loading.value = false;
        displayOrderModified.value = false;

        if (force) {
            snackbar.value?.showMessage('Category list has been updated');
        }

        updateCardMinHeight();
    }).catch(error => {
        loading.value = false;

        if (error && error.isUpToDate) {
            displayOrderModified.value = false;
        }

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function add(): void {
    editDialog.value?.open({
        type: activeCategoryType.value,
        parentId: primaryCategoryId.value
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        updateCardMinHeight();
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function edit(category: TransactionCategory): void {
    editDialog.value?.open({
        id: category.id,
        currentCategory: category
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        if (transactionCategoriesStore.transactionCategoryListStateInvalid) {
            reload(true);
        }

        updateCardMinHeight();
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
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

        updateCardMinHeight();
    }).catch(error => {
        updating.value = false;
        categoryHiding.value[category.id] = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
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

            updateCardMinHeight();
        }).catch(error => {
            updating.value = false;
            categoryRemoving.value[category.id] = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
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
            snackbar.value?.showError(error);
        }
    });
}

function onMove(event: { moved: { element: { id: string }, oldIndex: number, newIndex: number } }): void {
    if (!event || !event.moved) {
        return;
    }

    const moveEvent = event.moved;

    if (!moveEvent.element || !moveEvent.element.id) {
        snackbar.value?.showMessage('Unable to move category');
        return;
    }

    transactionCategoriesStore.changeCategoryDisplayOrder({
        categoryId: moveEvent.element.id,
        from: moveEvent.oldIndex,
        to: moveEvent.newIndex
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        snackbar.value?.showError(error);
    });
}

function onPresetCategorySaved(e: { message: string }): void {
    if (e && e.message) {
        snackbar.value?.showMessage(e.message);
        reload(false);
    }
}

watch(() => display.mdAndUp.value, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});

reload(false);
</script>

<style>
.transaction-category-table tr.transaction-category-table-row .hover-display {
    display: none;
}

.transaction-category-table tr.transaction-category-table-row:hover .hover-display {
    display: grid;
}

.transaction-category-table .transaction-category-comment {
    font-size: 0.8rem;
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
}
</style>
