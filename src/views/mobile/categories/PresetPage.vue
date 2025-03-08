<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Default Categories')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" v-if="isPresetHasCategories" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="tt('Save')" :class="{ 'disabled': submitting }" v-if="isPresetHasCategories" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="no-padding no-margin" :key="categoryType" v-for="(categories, categoryType) in allPresetCategories">
            <f7-block-title class="margin-top margin-horizontal">{{ getCategoryTypeName(categoryType) }}</f7-block-title>

            <f7-list strong inset dividers class="margin-top">
                <f7-list-item :title="category.name"
                              :accordion-item="!!category.subCategories.length"
                              :key="idx"
                              v-for="(category, idx) in categories">
                    <template #media>
                        <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                    </template>

                    <f7-accordion-content v-if="category.subCategories.length" class="padding-left">
                        <f7-list>
                            <f7-list-item :title="subCategory.name"
                                          :key="subIdx"
                                          v-for="(subCategory, subIdx) in category.subCategories">
                                <template #media>
                                    <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-content>
                </f7-list-item>
            </f7-list>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="showChangeLocaleSheet = true">{{ tt('Change Language') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <list-item-selection-sheet value-type="item"
                                   key-field="languageTag" value-field="languageTag"
                                   title-field="nativeDisplayName" after-field="displayName"
                                   :items="allLanguages"
                                   v-model:show="showChangeLocaleSheet"
                                   v-model="currentLocale">
        </list-item-selection-sheet>
    </f7-page>
</template>

<script setup lang="ts">
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

import type { PartialRecord } from '@/core/base.ts';
import type { LanguageOption } from '@/locales/index.ts';
import { type LocalizedPresetCategory, CategoryType } from '@/core/category.ts';
import { getObjectOwnFieldCount, categorizedArrayToPlainArray } from '@/lib/common.ts';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const { tt, getCurrentLanguageTag, getAllLanguageOptions, getAllTransactionDefaultCategories } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const transactionCategoriesStore = useTransactionCategoriesStore();

const loadingError = ref<unknown | null>(null);
const currentLocale = ref<string>(getCurrentLanguageTag());
const categoryType = ref<number>(0);
const submitting = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const showChangeLocaleSheet = ref<boolean>(false);

const allLanguages = computed<LanguageOption[]>(() => getAllLanguageOptions(false));
const allPresetCategories = computed<PartialRecord<CategoryType, LocalizedPresetCategory[]>>(() => getAllTransactionDefaultCategories(categoryType.value, currentLocale.value));
const isPresetHasCategories = computed<boolean>(() => getObjectOwnFieldCount(allPresetCategories.value) > 0);

function getCategoryTypeName(categoryType: CategoryType): string {
    switch (categoryType) {
        case CategoryType.Income:
            return tt('Income Categories');
        case CategoryType.Expense:
            return tt('Expense Categories');
        case CategoryType.Transfer:
            return tt('Transfer Categories');
        default:
            return tt('Transaction Categories');
    }
}

function save(): void {
    const router = props.f7router;

    submitting.value = true;
    showLoading(() => submitting.value);

    const submitCategories = categorizedArrayToPlainArray(allPresetCategories.value);

    transactionCategoriesStore.addCategories({
        categories: submitCategories
    }).then(() => {
        submitting.value = false;
        hideLoading();

        showToast('You have added preset categories');
        router.back();
    }).catch(error => {
        submitting.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

const query = props.f7route.query;
categoryType.value = parseInt(query['type'] || '0');

if (categoryType.value !== 0 &&
    categoryType.value !== CategoryType.Income &&
    categoryType.value !== CategoryType.Expense &&
    categoryType.value !== CategoryType.Transfer) {
    showToast('Parameter Invalid');
    loadingError.value = 'Parameter Invalid';
}
</script>
