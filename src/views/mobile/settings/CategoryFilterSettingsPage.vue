<template>
    <f7-page with-subnavbar @page:beforein="onPageBeforeIn" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt(title)"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableCategory }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': !hasAnyAvailableCategory }" @click="save"></f7-link>
            </f7-nav-right>

            <f7-subnavbar :inner="false">
                <f7-searchbar
                    custom-searchs
                    :value="filterContent"
                    :placeholder="tt('Find category')"
                    :disable-button-text="tt('Cancel')"
                    @change="filterContent = $event.target.value"
                ></f7-searchbar>
            </f7-subnavbar>
        </f7-navbar>

        <div class="skeleton-text" v-if="loading">
            <f7-block class="combination-list-wrapper margin-vertical"
                      :key="blockIdx" v-for="blockIdx in [ 1, 2 ]">
                <f7-accordion-item>
                    <f7-block-title>
                        <f7-accordion-toggle>
                            <f7-list strong inset dividers
                                     class="combination-list-header combination-list-opened">
                                <f7-list-item group-title>
                                    <small>Transaction Category</small>
                                    <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
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
                                    <ul class="padding-inline-start">
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
                  :key="categoryType"
                  v-for="(categories, categoryType) in allVisibleTransactionCategories"
                  v-else-if="!loading">
            <f7-accordion-item :opened="collapseStates[categoryType]!.opened"
                               @accordion:open="collapseStates[categoryType]!.opened = true"
                               @accordion:close="collapseStates[categoryType]!.opened = false">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers
                                 class="combination-list-header"
                                 :class="collapseStates[categoryType]!.opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item group-title>
                                <small>{{ getCategoryTypeName(parseInt(categoryType)) }}</small>
                                <f7-icon class="combination-list-chevron-icon" :f7="collapseStates[categoryType]!.opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates[categoryType]!.opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content" v-if="!categories || !categories.length">
                        <f7-list-item :title="tt('No available category')"></f7-list-item>
                    </f7-list>
                    <f7-list strong inset dividers accordion-list class="combination-list-content" v-else-if="categories && categories.length">
                        <f7-list-item checkbox
                                      :class="{ 'has-child-list-item': category.subCategories && category.subCategories.length > 0 }"
                                      :title="category.name"
                                      :value="category.id"
                                      :checked="isSubCategoriesAllChecked(category, filterCategoryIds)"
                                      :indeterminate="isSubCategoriesHasButNotAllChecked(category, filterCategoryIds)"
                                      :key="category.id"
                                      v-for="category in categories"
                                      @change="updateAllSubCategoriesSelected">
                            <template #media>
                                <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color">
                                    <f7-badge color="gray" class="right-bottom-icon" v-if="category.hidden">
                                        <f7-icon f7="eye_slash_fill"></f7-icon>
                                    </f7-badge>
                                </ItemIcon>
                            </template>

                            <template #root>
                                <ul class="padding-inline-start"
                                    v-if="category.subCategories && category.subCategories.length > 0">
                                    <f7-list-item checkbox
                                                  :title="subCategory.name"
                                                  :value="subCategory.id"
                                                  :checked="isCategoryChecked(subCategory, filterCategoryIds)"
                                                  :key="subCategory.id"
                                                  v-for="subCategory in category.subCategories"
                                                  @change="updateCategorySelected">
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
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleCategory }" @click="selectAllCategories">{{ tt('Select All') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleCategory }" @click="selectNoneCategories">{{ tt('Select None') }}</f7-actions-button>
                <f7-actions-button :class="{ 'disabled': !hasAnyVisibleCategory }" @click="selectInvertCategories">{{ tt('Invert Selection') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ tt('Show Hidden Transaction Categories') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ tt('Hide Hidden Transaction Categories') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import {
    type CategoryFilterType,
    useCategoryFilterSettingPageBase
} from '@/views/base/settings/CategoryFilterSettingPageBase.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { CategoryType } from '@/core/category.ts';

import {
    selectAllSubCategories,
    selectAll,
    selectNone,
    selectInvert,
    isSubCategoriesAllChecked,
    isSubCategoriesHasButNotAllChecked
} from '@/lib/category.ts';

interface CollapseState {
    opened: boolean;
}

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const query = props.f7route.query;

const { tt } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const {
    loading,
    showHidden,
    filterContent,
    filterCategoryIds,
    title,
    allVisibleTransactionCategories,
    allVisibleTransactionCategoryMap,
    hasAnyAvailableCategory,
    hasAnyVisibleCategory,
    isCategoryChecked,
    getCategoryTypeName,
    loadFilterCategoryIds,
    saveFilterCategoryIds
} = useCategoryFilterSettingPageBase(query['type'] as CategoryFilterType, query['allowCategoryTypes']);

const transactionCategoriesStore = useTransactionCategoriesStore();

const loadingError = ref<unknown | null>(null);
const showMoreActionSheet = ref<boolean>(false);

const collapseStates = ref<Record<string, CollapseState>>({
    [CategoryType.Income.toString()]: {
        opened: true
    },
    [CategoryType.Expense.toString()]: {
        opened: true
    },
    [CategoryType.Transfer.toString()]: {
        opened: true
    }
});

function init(): void {
    transactionCategoriesStore.loadAllCategories({
        force: false
    }).then(() => {
        loading.value = false;

        if (!loadFilterCategoryIds()) {
            showToast('Parameter Invalid');
            loadingError.value = 'Parameter Invalid';
        }
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function updateCategorySelected(e: Event): void {
    const target = e.target as HTMLInputElement;
    const categoryId = target.value;
    const category = allVisibleTransactionCategoryMap.value[categoryId];

    if (!category) {
        return;
    }

    filterCategoryIds.value[category.id] = !target.checked;
}

function updateAllSubCategoriesSelected(e: Event): void {
    const target = e.target as HTMLInputElement;
    const categoryId = target.value;
    const category = allVisibleTransactionCategoryMap.value[categoryId];

    selectAllSubCategories(filterCategoryIds.value, !target.checked, category);
}

function selectAllCategories(): void {
    selectAll(filterCategoryIds.value, allVisibleTransactionCategoryMap.value);
}

function selectNoneCategories(): void {
    selectNone(filterCategoryIds.value, allVisibleTransactionCategoryMap.value);
}

function selectInvertCategories(): void {
    selectInvert(filterCategoryIds.value, allVisibleTransactionCategoryMap.value);
}

function save(): void {
    saveFilterCategoryIds();
    props.f7router.back();
}

function onPageBeforeIn(): void {
    filterContent.value = '';
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

init();
</script>
