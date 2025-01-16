<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t(saveButtonTitle)" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-top skeleton-text" v-if="loading">
            <f7-list-input label="Category Name" placeholder="Your category name"></f7-list-input>
            <f7-list-item class="list-item-with-header-and-title" header="Primary Category" title="Primary Category"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                <template #default>
                    <div class="grid grid-cols-2">
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>Category Icon</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <f7-icon f7="app_fill"></f7-icon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>
                        </div>
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>Category Color</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <f7-icon f7="app_fill"></f7-icon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>
                        </div>
                    </div>
                </template>
            </f7-list-item>
            <f7-list-item class="list-item-toggle" header="Visible" after="True"></f7-list-item>
            <f7-list-input label="Description" type="textarea" placeholder="Your category description (optional)"></f7-list-input>
        </f7-list>

        <f7-list form strong inset dividers class="margin-top" v-else-if="!loading">
            <f7-list-input
                type="text"
                clear-button
                :label="$t('Category Name')"
                :placeholder="$t('Your category name')"
                v-model:value="category.name"
            ></f7-list-input>

            <f7-list-item
                link="#" no-chevron
                class="list-item-with-header-and-title"
                :header="$t('Primary Category')"
                :title="getPrimaryCategoryName(category.parentId)"
                @click="showPrimaryCategorySheet = true"
                v-if="editCategoryId && category.parentId && category.parentId !== '0'"
            >
                <list-item-selection-sheet value-type="item"
                                           key-field="id" value-field="id" title-field="name"
                                           icon-field="icon" icon-type="category" color-field="color"
                                           :items="allAvailableCategories"
                                           v-model:show="showPrimaryCategorySheet"
                                           v-model="category.parentId">
                </list-item-selection-sheet>
            </f7-list-item>

            <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                <template #default>
                    <div class="grid grid-cols-2">
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="showIconSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ $t('Category Icon') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <icon-selection-sheet :all-icon-infos="allCategoryIcons"
                                                  :color="category.color"
                                                  v-model:show="showIconSelectionSheet"
                                                  v-model="category.icon"
                            ></icon-selection-sheet>
                        </div>
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="showColorSelectionSheet = true">
                                <div class="item-content">
                                    <div class="item-inner">
                                        <div class="item-header">
                                            <span>{{ $t('Category Color') }}</span>
                                        </div>
                                        <div class="item-title">
                                            <div class="list-item-custom-title no-padding">
                                                <ItemIcon icon-type="fixed-f7" icon-id="app_fill" :color="category.color"></ItemIcon>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </a>

                            <color-selection-sheet :all-color-infos="allCategoryColors"
                                                   v-model:show="showColorSelectionSheet"
                                                   v-model="category.color"
                            ></color-selection-sheet>
                        </div>
                    </div>
                </template>
            </f7-list-item>

            <f7-list-item :title="$t('Visible')" v-if="editCategoryId">
                <f7-toggle :checked="category.visible" @toggle:change="category.visible = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-input
                type="textarea"
                style="height: auto"
                :label="$t('Description')"
                :placeholder="$t('Your category description (optional)')"
                v-textarea-auto-size
                v-model:value="category.comment"
            ></f7-list-input>
        </f7-list>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { CategoryType } from '@/core/category.ts';
import { ALL_CATEGORY_ICONS } from '@/consts/icon.ts';
import { ALL_CATEGORY_COLORS } from '@/consts/color.ts';
import { TransactionCategory } from '@/models/transaction_category.ts';

import { getNameByKeyValue } from '@/lib/common.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import { allVisiblePrimaryTransactionCategoriesByType } from '@/lib/category.ts';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        const query = this.f7route.query;

        return {
            editCategoryId: null,
            clientSessionId: '',
            loading: false,
            loadingError: null,
            category: TransactionCategory.createNewCategory(parseInt(query.type), query.parentId),
            showPrimaryCategorySheet: false,
            showIconSelectionSheet: false,
            showColorSelectionSheet: false,
            submitting: false
        };
    },
    computed: {
        ...mapStores(useTransactionCategoriesStore),
        allAvailableCategories() {
            return allVisiblePrimaryTransactionCategoriesByType(this.transactionCategoriesStore.allTransactionCategories, this.category.type);
        },
        title() {
            if (!this.editCategoryId) {
                if (this.category.parentId === '0') {
                    return 'Add Primary Category';
                } else {
                    return 'Add Secondary Category';
                }
            } else {
                return 'Edit Category';
            }
        },
        saveButtonTitle() {
            if (!this.editCategoryId) {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        allCategoryIcons() {
            return ALL_CATEGORY_ICONS;
        },
        allCategoryColors() {
            return ALL_CATEGORY_COLORS;
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (!this.category.name) {
                return 'Category name cannot be blank';
            } else {
                return null;
            }
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        if (!query.id && !query.parentId) {
            self.$toast('Parameter Invalid');
            self.loadingError = 'Parameter Invalid';
            return;
        }

        if (query.id) {
            self.loading = true;

            self.editCategoryId = query.id;
            self.transactionCategoriesStore.getCategory({
                categoryId: self.editCategoryId
            }).then(category => {
                self.category.from(category);
                self.loading = false;
            }).catch(error => {
                if (error.processed) {
                    self.loading = false;
                } else {
                    self.loadingError = error;
                    self.$toast(error.message || error);
                }
            });
        } else if (query.parentId) {
            const categoryType = parseInt(query.type);

            if (categoryType !== CategoryType.Income &&
                categoryType !== CategoryType.Expense &&
                categoryType !== CategoryType.Transfer) {
                self.$toast('Parameter Invalid');
                self.loadingError = 'Parameter Invalid';
                return;
            }

            self.clientSessionId = generateRandomUUID();
            self.loading = false;
        }
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        save() {
            const self = this;
            const router = self.f7router;

            const problemMessage = self.inputEmptyProblemMessage;

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            self.transactionCategoriesStore.saveCategory({
                category: self.category,
                isEdit: !!self.editCategoryId,
                clientSessionId: self.clientSessionId
            }).then(() => {
                self.submitting = false;
                self.$hideLoading();

                if (!self.editCategoryId) {
                    self.$toast('You have added a new category');
                } else {
                    self.$toast('You have saved this category');
                }

                router.back();
            }).catch(error => {
                self.submitting = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        getPrimaryCategoryName(parentId) {
            return getNameByKeyValue(this.allAvailableCategories, parentId, 'id', 'name');
        }
    }
}
</script>
