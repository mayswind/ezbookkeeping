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

            <f7-list-item class="list-item-with-header-and-title list-item-with-multi-item">
                <template #default>
                    <div class="grid grid-cols-2">
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="category.showIconSelectionSheet = true">
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
                                                  v-model:show="category.showIconSelectionSheet"
                                                  v-model="category.icon"
                            ></icon-selection-sheet>
                        </div>
                        <div class="list-item-subitem no-chevron">
                            <a class="item-link" href="#" @click="category.showColorSelectionSheet = true">
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
                                                   v-model:show="category.showColorSelectionSheet"
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
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';

import categoryConstants from '@/consts/category.js';
import iconConstants from '@/consts/icon.js';
import colorConstants from '@/consts/color.js';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        const self = this;
        const query = self.f7route.query;

        return {
            editCategoryId: null,
            loading: false,
            loadingError: null,
            category: {
                type: parseInt(query.type),
                name: '',
                parentId: query.parentId,
                icon: iconConstants.defaultCategoryIconId,
                color: colorConstants.defaultCategoryColor,
                comment: '',
                visible: true,
                showIconSelectionSheet: false,
                showColorSelectionSheet: false
            },
            submitting: false
        };
    },
    computed: {
        ...mapStores(useTransactionCategoriesStore),
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
            return iconConstants.allCategoryIcons;
        },
        allCategoryColors() {
            return colorConstants.allCategoryColors;
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (!this.category.name) {
                return 'Category name cannot be empty';
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
                self.category.id = category.id;
                self.category.type = category.type;
                self.category.parentId = category.type.parentId;
                self.category.name = category.name;
                self.category.icon = category.icon;
                self.category.color = category.color;
                self.category.comment = category.comment;
                self.category.visible = !category.hidden;

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

            if (categoryType !== categoryConstants.allCategoryTypes.Income &&
                categoryType !== categoryConstants.allCategoryTypes.Expense &&
                categoryType !== categoryConstants.allCategoryTypes.Transfer) {
                self.$toast('Parameter Invalid');
                self.loadingError = 'Parameter Invalid';
                return;
            }

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

            const submitCategory = {
                type: self.category.type,
                name: self.category.name,
                parentId: self.category.parentId,
                icon: self.category.icon,
                color: self.category.color,
                comment: self.category.comment
            };

            if (self.editCategoryId) {
                submitCategory.id = self.category.id;
                submitCategory.hidden = !self.category.visible;
            }

            self.transactionCategoriesStore.saveCategory({
                category: submitCategory
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
        }
    }
}
</script>
