<template>
    <f7-page :ptr="!sortable" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :class="{ 'disabled': loading }" :back-link="tt('Back')" v-if="!sortable"></f7-nav-left>
            <f7-nav-left v-else-if="sortable">
                <f7-link icon-f7="xmark" :class="{ 'disabled': displayOrderSaving }" @click="cancelSort"></f7-link>
            </f7-nav-left>
            <f7-nav-title :title="tt(title)"></f7-nav-title>
            <f7-nav-right :class="{ 'navbar-compact-icons': true, 'disabled': loading }">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !categories.length || sortable }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="plus" :href="'/category/add?type=' + categoryType + '&parentId=' + primaryCategoryId + (currentPrimaryCategory ? `&color=${currentPrimaryCategory.color}&icon=${currentPrimaryCategory.icon}` : '')" v-if="!sortable"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': displayOrderSaving || !displayOrderModified }" @click="saveSortResult" v-else-if="sortable"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-top skeleton-text" v-if="loading">
            <f7-list-item title="Category Name"
                          :link="hasSubCategories ? '#' : null"
                          :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                <template #media>
                    <f7-icon f7="app_fill"></f7-icon>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-top" v-if="!loading && noAvailableCategory">
            <f7-list-item :title="tt('No available category')"></f7-list-item>
            <f7-list-button v-if="hasSubCategories && noCategory"
                            :title="tt('Add Default Categories')"
                            :href="'/category/preset?type=' + categoryType"></f7-list-button>
        </f7-list>

        <f7-list strong inset dividers sortable class="margin-top category-list"
                 :sortable-enabled="sortable"
                 v-if="!loading"
                 @sortable:sort="onSort">
            <f7-list-item swipeout
                          :class="{ 'actual-first-child': category.id === firstShowingId, 'actual-last-child': category.id === lastShowingId }"
                          :id="getCategoryDomId(category)"
                          :title="category.name"
                          :footer="category.comment"
                          :link="hasSubCategories ? '/category/list?type=' + categoryType + '&id=' + category.id : null"
                          :key="category.id"
                          v-for="category in categories"
                          v-show="showHidden || !category.hidden"
                          @taphold="setSortable()">
                <template #media>
                    <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color">
                        <f7-badge color="gray" class="right-bottom-icon" v-if="category.hidden">
                            <f7-icon f7="eye_slash_fill"></f7-icon>
                        </f7-badge>
                    </ItemIcon>
                </template>
                <f7-swipeout-actions :left="textDirection === TextDirection.LTR"
                                     :right="textDirection === TextDirection.RTL"
                                     v-if="sortable">
                    <f7-swipeout-button :color="category.hidden ? 'blue' : 'gray'" class="padding-horizontal"
                                        overswipe close @click="hide(category, !category.hidden)">
                        <f7-icon :f7="category.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
                <f7-swipeout-actions :left="textDirection === TextDirection.RTL"
                                     :right="textDirection === TextDirection.LTR"
                                     v-if="!sortable">
                    <f7-swipeout-button color="orange" close :text="tt('Edit')" @click="edit(category)"></f7-swipeout-button>
                    <f7-swipeout-button color="red" class="padding-horizontal" @click="remove(category, false)">
                        <f7-icon f7="trash"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ tt('Sort') }}</f7-actions-button>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ tt('Show Hidden Transaction Categories') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ tt('Hide Hidden Transaction Categories') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ tt('Are you sure you want to delete this category?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(categoryToDelete, true)">{{ tt('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading, onSwipeoutDeleted } from '@/lib/ui/mobile.ts';
import { useCategoryListPageBase } from '@/views/base/categories/CategoryListPageBase.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { TextDirection } from '@/core/text.ts';
import { CategoryType } from '@/core/category.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

import {
    isNoAvailableCategory,
    getFirstShowingId,
    getLastShowingId
} from '@/lib/category.ts';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const { tt, getCurrentLanguageTextDirection } = useI18n();
const { showAlert, showToast, routeBackOnError } = useI18nUIComponents();
const { loading, primaryCategoryId, currentPrimaryCategory } = useCategoryListPageBase();

const transactionCategoriesStore = useTransactionCategoriesStore();

const hasSubCategories = ref<boolean>(false);
const categoryType = ref<CategoryType | 0>(0);
const loadingError = ref<unknown | null>(null);
const showHidden = ref<boolean>(false);
const sortable = ref<boolean>(false);
const categoryToDelete = ref<TransactionCategory | null>(null);
const showMoreActionSheet = ref<boolean>(false);
const showDeleteActionSheet = ref<boolean>(false);
const displayOrderModified = ref<boolean>(false);
const displayOrderSaving = ref<boolean>(false);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());

const categories = computed<TransactionCategory[]>(() => {
    if (!primaryCategoryId.value || primaryCategoryId.value === '' || primaryCategoryId.value === '0') {
        if (!transactionCategoriesStore.allTransactionCategories || !transactionCategoriesStore.allTransactionCategories[categoryType.value]) {
            return [];
        }

        return transactionCategoriesStore.allTransactionCategories[categoryType.value] ?? [];
    } else if (primaryCategoryId.value && primaryCategoryId.value !== '' && primaryCategoryId.value !== '0') {
        if (!transactionCategoriesStore.allTransactionCategoriesMap || !transactionCategoriesStore.allTransactionCategoriesMap[primaryCategoryId.value]) {
            return [];
        }

        return transactionCategoriesStore.allTransactionCategoriesMap[primaryCategoryId.value]?.subCategories ?? [];
    } else {
        return [];
    }
});

const title = computed<string>(() => {
    let title = '';

    switch (categoryType.value) {
        case CategoryType.Income:
            title = 'Income';
            break;
        case CategoryType.Expense:
            title = 'Expense';
            break;
        case CategoryType.Transfer:
            title = 'Transfer';
            break;
        default:
            title = 'Transaction';
            break;
    }

    switch (hasSubCategories.value) {
        case true:
            title += ' Primary';
            break;
        case false:
            title += ' Secondary';
            break;
    }

    return title + ' Categories';
});

const firstShowingId = computed<string | null>(() => getFirstShowingId(categories.value, showHidden.value));
const lastShowingId = computed<string | null>(() => getLastShowingId(categories.value, showHidden.value));
const noAvailableCategory = computed<boolean>(() => isNoAvailableCategory(categories.value, showHidden.value));
const noCategory = computed<boolean>(() => categories.value.length < 1);

function getCategoryDomId(category: TransactionCategory): string {
    return 'category_' + category.id;
}

function parseCategoryIdFromDomId(domId: string): string | null {
    if (!domId || domId.indexOf('category_') !== 0) {
        return null;
    }

    return domId.substring(9); // category_
}

function init(): void {
    const query = props.f7route.query;

    categoryType.value = parseInt(query['type'] || '0');

    if (categoryType.value !== CategoryType.Income &&
        categoryType.value !== CategoryType.Expense &&
        categoryType.value !== CategoryType.Transfer) {
        showToast('Parameter Invalid');
        loadingError.value = 'Parameter Invalid';
        return;
    }

    if (query['id'] && query['id'] !== '0') {
        primaryCategoryId.value = query['id'];
        hasSubCategories.value = false;
    } else {
        primaryCategoryId.value = '0';
        hasSubCategories.value = true;
    }

    loading.value = true;

    transactionCategoriesStore.loadAllCategories({
        force: false
    }).then(() => {
        loading.value = false;
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function reload(done?: () => void): void {
    if (sortable.value) {
        done?.();
        return;
    }

    const force = !!done;

    transactionCategoriesStore.loadAllCategories({
        force: force
    }).then(() => {
        done?.();

        if (force) {
            showToast('Category list has been updated');
        }
    }).catch(error => {
        done?.();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function edit(category: TransactionCategory): void {
    props.f7router.navigate('/category/edit?id=' + category.id);
}

function hide(category: TransactionCategory, hidden: boolean): void {
    showLoading();

    transactionCategoriesStore.hideCategory({
        category: category,
        hidden: hidden
    }).then(() => {
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function remove(category: TransactionCategory | null, confirm: boolean): void {
    if (!category) {
        showAlert('An error occurred');
        return;
    }

    if (!confirm) {
        categoryToDelete.value = category;
        showDeleteActionSheet.value = true;
        return;
    }

    showDeleteActionSheet.value = false;
    categoryToDelete.value = null;
    showLoading();

    transactionCategoriesStore.deleteCategory({
        category: category,
        beforeResolve: (done) => {
            onSwipeoutDeleted(getCategoryDomId(category), done);
        }
    }).then(() => {
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function setSortable(): void {
    if (sortable.value) {
        return;
    }

    showHidden.value = true;
    sortable.value = true;
    displayOrderModified.value = false;
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        showHidden.value = false;
        sortable.value = false;
        return;
    }

    displayOrderSaving.value = true;
    showLoading();

    transactionCategoriesStore.updateCategoryDisplayOrders({
        type: categoryType.value as CategoryType,
        parentId: primaryCategoryId.value
    }).then(() => {
        displayOrderSaving.value = false;
        hideLoading();

        showHidden.value = false;
        sortable.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        displayOrderSaving.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function cancelSort(): void {
    if (!displayOrderModified.value) {
        showHidden.value = false;
        sortable.value = false;
        return;
    }

    displayOrderSaving.value = true;
    showLoading();

    transactionCategoriesStore.loadAllCategories({
        force: false
    }).then(() => {
        displayOrderSaving.value = false;
        hideLoading();

        showHidden.value = false;
        sortable.value = false;
        displayOrderModified.value = false;
    }).catch(error => {
        displayOrderSaving.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function onSort(event: { el: { id: string }; from: number; to: number }): void {
    if (!event || !event.el || !event.el.id) {
        showToast('Unable to move category');
        return;
    }

    const id = parseCategoryIdFromDomId(event.el.id);

    if (!id) {
        showToast('Unable to move category');
        return;
    }

    transactionCategoriesStore.changeCategoryDisplayOrder({
        categoryId: id,
        from: event.from,
        to: event.to
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        showToast(error.message || error);
    });
}

function onPageAfterIn(): void {
    if (transactionCategoriesStore.transactionCategoryListStateInvalid && !loading.value) {
        reload();
    }

    routeBackOnError(props.f7router, loadingError);
}

init();
</script>

<style>
.category-list {
    --f7-list-item-footer-font-size: var(--ebk-large-footer-font-size);
}

.category-list .item-footer {
    padding-top: 4px;
}
</style>
