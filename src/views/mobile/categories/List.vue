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

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item title="Category Name">
                        <f7-icon slot="media" f7="app_fill"></f7-icon>
                    </f7-list-item>
                    <f7-list-item title="Category Name 2">
                        <f7-icon slot="media" f7="app_fill"></f7-icon>
                    </f7-list-item>
                    <f7-list-item title="Category Name 3">
                        <f7-icon slot="media" f7="app_fill"></f7-icon>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list v-if="noAvailableCategory">
                    <f7-list-item :title="$t('No available category')"></f7-list-item>
                    <f7-list-button v-if="hasSubCategories"
                        :title="$t('Add Default Categories')"
                        :href="'/category/preset?type=' + categoryType"></f7-list-button>
                </f7-list>

                <f7-list sortable :sortable-enabled="sortable" @sortable:sort="onSort">
                    <f7-list-item v-for="category in categories"
                                  :key="category.id"
                                  :id="category | categoryDomId"
                                  :title="category.name"
                                  :link="hasSubCategories ? '/category/list?type=' + categoryType + '&id=' + category.id : null"
                                  v-show="showHidden || !category.hidden"
                                  swipeout @taphold.native="setSortable()">
                        <f7-icon slot="media"
                                 :icon="category.icon | categoryIcon"
                                 :style="category.color | categoryIconStyle('var(--default-icon-color)')">
                            <f7-badge color="gray" class="right-bottom-icon" v-if="category.hidden">
                                <f7-icon f7="eye_slash_fill"></f7-icon>
                            </f7-badge>
                        </f7-icon>
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
            </f7-card-content>
        </f7-card>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ $t('Sort') }}</f7-actions-button>
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
export default {
    data() {
        return {
            hasSubCategories: false,
            categoryType: 0,
            categoryId: '',
            loading: true,
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
        categories() {
            if (!this.categoryId || this.categoryId === '' || this.categoryId === '0') {
                if (!this.$store.state.allTransactionCategories || !this.$store.state.allTransactionCategories[this.categoryType]) {
                    return [];
                }

                return this.$store.state.allTransactionCategories[this.categoryType];
            } else if (this.categoryId && this.categoryId !== '' && this.categoryId !== '0') {
                if (!this.$store.state.allTransactionCategoriesMap || !this.$store.state.allTransactionCategoriesMap[this.categoryId]) {
                    return [];
                }

                return this.$store.state.allTransactionCategoriesMap[this.categoryId].subCategories;
            } else {
                return [];
            }
        },
        title() {
            let title = '';

            switch (this.categoryType) {
                case this.$constants.category.allCategoryTypes.Income:
                    title = 'Income';
                    break;
                case this.$constants.category.allCategoryTypes.Expense:
                    title = 'Expense';
                    break;
                case this.$constants.category.allCategoryTypes.Transfer:
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
        const query = self.$f7route.query;
        const router = self.$f7router;

        self.categoryType = parseInt(query.type);

        if (self.categoryType !== this.$constants.category.allCategoryTypes.Income &&
            self.categoryType !== this.$constants.category.allCategoryTypes.Expense &&
            self.categoryType !== this.$constants.category.allCategoryTypes.Transfer) {
            self.$toast('Parameter Invalid');
            router.back();
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

        self.$store.dispatch('loadAllCategories', {
            force: false
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            self.logining = false;

            if (!error.processed) {
                self.$toast(error.message || error);
                router.back();
            }
        });
    },
    methods: {
        onPageAfterIn() {
            if (this.$store.state.transactionCategoryListStateInvalid && !this.loading) {
                this.reload(null);
            }
        },
        reload(done) {
            if (this.sortable) {
                done();
                return;
            }

            const self = this;

            self.$store.dispatch('loadAllCategories', {
                force: true
            }).then(() => {
                if (done) {
                    done();
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

            if (!event || !event.el || !event.el.id || event.el.id.indexOf('category_') !== 0) {
                this.$toast('Unable to move category');
                return;
            }

            const id = event.el.id.substr(9); // category_

            self.$store.dispatch('changeCategoryDisplayOrder', {
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

            self.$store.dispatch('updateCategoryDisplayOrders', {
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
            this.$f7router.navigate('/category/edit?id=' + category.id);
        },
        hide(category, hidden) {
            const self = this;

            self.$showLoading();

            self.$store.dispatch('hideCategory', {
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
            const app = self.$f7;
            const $$ = app.$;

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

            self.$store.dispatch('deleteCategory', {
                category: category,
                beforeResolve: (done) => {
                    app.swipeout.delete($$(`#${self.$options.filters.categoryDomId(category)}`), () => {
                        done();
                    });
                }
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        }
    },
    filters: {
        categoryDomId(category) {
            return 'category_' + category.id;
        }
    }
};
</script>
