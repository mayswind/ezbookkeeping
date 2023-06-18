<template>
    <f7-page :ptr="!sortable" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" v-if="!sortable && this.categories.length" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :href="'/category/add?type=' + categoryType + '&parentId=' + categoryId" icon-f7="plus" v-if="!sortable"></f7-link>
                <f7-link :text="$t('Done')" :class="{ 'disabled': displayOrderSaving }" @click="saveSortResult" v-else-if="sortable"></f7-link>
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
            <f7-list-item :title="$t('No available category')"></f7-list-item>
            <f7-list-button v-if="hasSubCategories"
                            :title="$t('Add Default Categories')"
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
                <f7-swipeout-actions left v-if="sortable">
                    <f7-swipeout-button :color="category.hidden ? 'blue' : 'gray'" class="padding-left padding-right"
                                        overswipe close @click="hide(category, !category.hidden)">
                        <f7-icon :f7="category.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
                <f7-swipeout-actions right v-if="!sortable">
                    <f7-swipeout-button color="orange" close :text="$t('Edit')" @click="edit(category)"></f7-swipeout-button>
                    <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(category, false)">
                        <f7-icon f7="trash"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ $t('Sort') }}</f7-actions-button>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ $t('Show Hidden Transaction Category') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ $t('Hide Hidden Transaction Category') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ $t('Are you sure you want to delete this category?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(categoryToDelete, true)">{{ $t('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';

import categoryConstants from '@/consts/category.js';
import { onSwipeoutDeleted } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        return {
            hasSubCategories: false,
            categoryType: 0,
            categoryId: '',
            loading: true,
            loadingError: null,
            showHidden: false,
            sortable: false,
            categoryToDelete: null,
            showMoreActionSheet: false,
            showDeleteActionSheet: false,
            displayOrderModified: false,
            displayOrderSaving: false
        };
    },
    computed: {
        ...mapStores(useTransactionCategoriesStore),
        categories() {
            if (!this.categoryId || this.categoryId === '' || this.categoryId === '0') {
                if (!this.transactionCategoriesStore.allTransactionCategories || !this.transactionCategoriesStore.allTransactionCategories[this.categoryType]) {
                    return [];
                }

                return this.transactionCategoriesStore.allTransactionCategories[this.categoryType];
            } else if (this.categoryId && this.categoryId !== '' && this.categoryId !== '0') {
                if (!this.transactionCategoriesStore.allTransactionCategoriesMap || !this.transactionCategoriesStore.allTransactionCategoriesMap[this.categoryId]) {
                    return [];
                }

                return this.transactionCategoriesStore.allTransactionCategoriesMap[this.categoryId].subCategories;
            } else {
                return [];
            }
        },
        title() {
            let title = '';

            switch (this.categoryType) {
                case categoryConstants.allCategoryTypes.Income:
                    title = 'Income';
                    break;
                case categoryConstants.allCategoryTypes.Expense:
                    title = 'Expense';
                    break;
                case categoryConstants.allCategoryTypes.Transfer:
                    title = 'Transfer';
                    break;
                default:
                    title = 'Transaction';
                    break;
            }

            switch (this.hasSubCategories) {
                case true:
                    title += ' Primary';
                    break;
                case false:
                    title += ' Secondary';
                    break;
            }

            return title + ' Categories';
        },
        firstShowingId() {
            for (let i = 0; i < this.categories.length; i++) {
                if (this.showHidden || !this.categories[i].hidden) {
                    return this.categories[i].id;
                }
            }

            return null;
        },
        lastShowingId() {
            for (let i = this.categories.length - 1; i >= 0; i--) {
                if (this.showHidden || !this.categories[i].hidden) {
                    return this.categories[i].id;
                }
            }

            return null;
        },
        noAvailableCategory() {
            for (let i = 0; i < this.categories.length; i++) {
                if (this.showHidden || !this.categories[i].hidden) {
                    return false;
                }
            }

            return true;
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        self.categoryType = parseInt(query.type);

        if (self.categoryType !== categoryConstants.allCategoryTypes.Income &&
            self.categoryType !== categoryConstants.allCategoryTypes.Expense &&
            self.categoryType !== categoryConstants.allCategoryTypes.Transfer) {
            self.$toast('Parameter Invalid');
            self.loadingError = 'Parameter Invalid';
            return;
        }

        if (query.id && query.id !== '0') {
            self.categoryId = query.id;
            self.hasSubCategories = false;
        } else {
            self.categoryId = '0';
            self.hasSubCategories = true;
        }

        self.loading = true;

        self.transactionCategoriesStore.loadAllCategories({
            force: false
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            if (this.transactionCategoriesStore.transactionCategoryListStateInvalid && !this.loading) {
                this.reload(null);
            }

            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            if (this.sortable) {
                done();
                return;
            }

            const self = this;
            const force = !!done;

            self.transactionCategoriesStore.loadAllCategories({
                force: force
            }).then(() => {
                if (done) {
                    done();
                }

                if (force) {
                    self.$toast('Category list has been updated');
                }
            }).catch(error => {
                if (done) {
                    done();
                }

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        setSortable() {
            if (this.sortable) {
                return;
            }

            this.showHidden = true;
            this.sortable = true;
            this.displayOrderModified = false;
        },
        onSort(event) {
            const self = this;

            if (!event || !event.el || !event.el.id) {
                self.$toast('Unable to move category');
                return;
            }

            const id = self.parseCategoryIdFromDomId(event.el.id);

            if (!id) {
                self.$toast('Unable to move category');
                return;
            }

            self.transactionCategoriesStore.changeCategoryDisplayOrder({
                categoryId: id,
                from: event.from,
                to: event.to
            }).then(() => {
                self.displayOrderModified = true;
            }).catch(error => {
                self.$toast(error.message || error);
            });
        },
        saveSortResult() {
            const self = this;

            if (!self.displayOrderModified) {
                self.showHidden = false;
                self.sortable = false;
                return;
            }

            self.displayOrderSaving = true;
            self.$showLoading();

            self.transactionCategoriesStore.updateCategoryDisplayOrders({
                type: self.categoryType,
                parentId: self.categoryId,
            }).then(() => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                self.showHidden = false;
                self.sortable = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        edit(category) {
            this.f7router.navigate('/category/edit?id=' + category.id);
        },
        hide(category, hidden) {
            const self = this;

            self.$showLoading();

            self.transactionCategoriesStore.hideCategory({
                category: category,
                hidden: hidden
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        remove(category, confirm) {
            const self = this;

            if (!category) {
                self.$alert('An error has occurred');
                return;
            }

            if (!confirm) {
                self.categoryToDelete = category;
                self.showDeleteActionSheet = true;
                return;
            }

            self.showDeleteActionSheet = false;
            self.categoryToDelete = null;
            self.$showLoading();

            self.transactionCategoriesStore.deleteCategory({
                category: category,
                beforeResolve: (done) => {
                    onSwipeoutDeleted(self.getCategoryDomId(category), done);
                }
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        getCategoryDomId(category) {
            return 'category_' + category.id;
        },
        parseCategoryIdFromDomId(domId) {
            if (!domId || domId.indexOf('category_') !== 0) {
                return null;
            }

            return domId.substring(9); // category_
        }
    }
};
</script>

<style>
.category-list {
    --f7-list-item-footer-font-size: var(--ebk-large-footer-font-size);
}

.category-list .item-footer {
    padding-top: 4px;
}
</style>
