<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Default Categories')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" v-if="isPresetHasCategories" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t('Save')" :class="{ 'disabled': submitting }" v-if="isPresetHasCategories" @click="save"></f7-link>
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
                <f7-actions-button @click="showChangeLocaleSheet = true">{{ $t('Change Language') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <list-item-selection-sheet value-type="index"
                                   title-field="displayName"
                                   :items="allLanguages"
                                   v-model:show="showChangeLocaleSheet"
                                   v-model="currentLocale">
        </list-item-selection-sheet>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';

import categoryConstants from '@/consts/category.js';
import { getObjectOwnFieldCount, categorizedArrayToPlainArray } from '@/lib/common.js';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        const self = this;

        return {
            loadingError: null,
            currentLocale: self.$locale.getCurrentLanguageTag(),
            categoryType: 0,
            submitting: false,
            showMoreActionSheet: false,
            showChangeLocaleSheet: false
        };
    },
    computed: {
        ...mapStores(useTransactionCategoriesStore),
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        },
        allPresetCategories() {
            return this.$locale.getAllTransactionDefaultCategories(this.categoryType, this.currentLocale);
        },
        isPresetHasCategories() {
            return getObjectOwnFieldCount(this.allPresetCategories);
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        self.categoryType = parseInt(query.type);

        if (self.categoryType !== 0 &&
            self.categoryType !== categoryConstants.allCategoryTypes.Income &&
            self.categoryType !== categoryConstants.allCategoryTypes.Expense &&
            self.categoryType !== categoryConstants.allCategoryTypes.Transfer) {
            self.$toast('Parameter Invalid');
            self.loadingError = 'Parameter Invalid';
        }
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        save() {
            const self = this;
            const router = self.f7router;

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            const submitCategories = categorizedArrayToPlainArray(self.allPresetCategories);

            self.transactionCategoriesStore.addCategories({
                categories: submitCategories
            }).then(() => {
                self.submitting = false;
                self.$hideLoading();

                self.$toast('You have added preset categories');
                router.back();
            }).catch(error => {
                self.submitting = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
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
        }
    }
};
</script>
