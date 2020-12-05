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
                    <f7-list-item title="Category Name"></f7-list-item>
                    <f7-list-item title="Category Name 2"></f7-list-item>
                    <f7-list-item title="Category Name 3"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list v-if="noAvailableCategory">
                    <f7-list-item :title="$t('No available category')"></f7-list-item>
                    <f7-list-button v-if="hasSubCategories"
                        :title="$t('Add Default Categories')"
                        :href="'/category/default?type=' + categoryType"></f7-list-button>
                </f7-list>

                <f7-list sortable :sortable-enabled="sortable" @sortable:sort="onSort">
                    <f7-list-item v-for="category in categories"
                                  :key="category.id"
                                  :id="category | categoryDomId"
                                  :title="category.name"
                                  :link="hasSubCategories ? '/category/list?type=' + categoryType + '&id=' + category.id : null"
                                  v-show="showHidden || !category.hidden"
                                  swipeout @taphold.native="setSortable()">
                        <f7-icon slot="media" :icon="category.icon | categoryIcon" :style="{ color: '#' + category.color }">
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
                            <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(category)">
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
                <f7-actions-button color="red" @click="remove(categoryToDelete)">{{ $t('Delete') }}</f7-actions-button>
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
            categories: [],
            hasSubCategories: false,
            categoryType: '',
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
        title() {
            let title = '';

            switch (this.categoryType) {
                case '1':
                    title = 'Expense';
                    break;
                case '2':
                    title = 'Income';
                    break;
                case '3':
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

        if (query.type !== '1' && query.type !== '2' && query.type !== '3') {
            self.$toast('Parameter Invalid');
            router.back();
            return;
        }

        self.categoryType = query.type;

        if (query.id && query.id !== '0') {
            self.categoryId = query.id;
            self.hasSubCategories = false;
        } else {
            self.categoryId = '0';
            self.hasSubCategories = true;
        }

        self.loading = true;

        self.$services.getAllTransactionCategories({
            type: self.categoryType,
            parentId: self.categoryId
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                self.$toast('Unable to get category list');
                router.back();
                return;
            }

            if (data.result[self.categoryType]) {
                self.categories = data.result[self.categoryType];
            } else {
                self.categories = [];
            }

            self.loading = false;
        }).catch(error => {
            self.$logger.error('failed to load category list', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                self.$toast({ error: error.response.data });
                router.back();
            } else if (!error.processed) {
                self.$toast('Unable to get category list');
                router.back();
            }
        });
    },
    methods: {
        onPageAfterIn() {
            const self = this;
            const previousRoute = self.$f7router.previousRoute;

            if (previousRoute && (previousRoute.path === '/category/add' || previousRoute.path === '/category/edit' || previousRoute.path === '/category/default') && !self.loading) {
                self.reload(null);
            }
        },
        reload(done) {
            if (this.sortable) {
                done();
                return;
            }

            const self = this;

            self.$services.getAllTransactionCategories({
                type: self.categoryType,
                parentId: self.categoryId
            }).then(response => {
                if (done) {
                    done();
                }

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to get category list');
                    return;
                }

                if (data.result[self.categoryType]) {
                    self.categories = data.result[self.categoryType];
                } else {
                    self.categories = [];
                }
            }).catch(error => {
                self.$logger.error('failed to reload category list', error);

                if (done) {
                    done();
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to get category list');
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
            if (!event || !event.el || !event.el.id || event.el.id.indexOf('category_') !== 0) {
                this.$toast('Unable to move category');
                return;
            }

            const id = event.el.id.substr(9); // category_
            let category = null;

            for (let i = 0; i < this.categories.length; i++) {
                if (this.categories[i].id === id) {
                    category = this.categories[i];
                    break;
                }
            }

            if (!category || !this.categories[event.to]) {
                this.$toast('Unable to move category');
                return;
            }

            this.categories.splice(event.to, 0, this.categories.splice(event.from, 1)[0]);

            this.displayOrderModified = true;
        },
        saveSortResult() {
            const self = this;
            const newDisplayOrders = [];

            if (!self.displayOrderModified) {
                self.showHidden = false;
                self.sortable = false;
                return;
            }

            self.displayOrderSaving = true;

            for (let i = 0; i < self.categories.length; i++) {
                newDisplayOrders.push({
                    id: self.categories[i].id,
                    displayOrder: i + 1
                });
            }

            self.$showLoading();

            self.$services.moveTransactionCategory({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to move category');
                    return;
                }

                self.showHidden = false;
                self.sortable = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.$logger.error('failed to save categories display order', error);

                self.displayOrderSaving = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to move category');
                }
            });
        },
        edit(category) {
            this.$f7router.navigate('/category/edit?id=' + category.id);
        },
        hide(category, hidden) {
            const self = this;

            self.$showLoading();

            self.$services.hideTransactionCategory({
                id: category.id,
                hidden: hidden
            }).then(response => {
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (hidden) {
                        self.$toast('Unable to hide this category');
                    } else {
                        self.$toast('Unable to unhide this category');
                    }

                    return;
                }

                category.hidden = hidden;
            }).catch(error => {
                self.$logger.error('failed to change category visibility', error);

                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    if (hidden) {
                        self.$toast('Unable to hide this category');
                    } else {
                        self.$toast('Unable to unhide this category');
                    }
                }
            });
        },
        remove(category) {
            const self = this;
            const app = self.$f7;
            const $$ = app.$;

            if (!category) {
                self.$alert('An error has occurred');
                return;
            }

            if (!self.showDeleteActionSheet) {
                self.categoryToDelete = category;
                self.showDeleteActionSheet = true;
                return;
            }

            self.showDeleteActionSheet = false;
            self.categoryToDelete = null;
            self.$showLoading();

            self.$services.deleteTransactionCategory({
                id: category.id
            }).then(response => {
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to delete this category');
                    return;
                }

                app.swipeout.delete($$(`#${self.$options.filters.categoryDomId(category)}`), () => {
                    for (let i = 0; i < self.categories.length; i++) {
                        if (self.categories[i].id === category.id) {
                            self.categories.splice(i, 1);
                        }
                    }
                });
            }).catch(error => {
                self.$logger.error('failed to delete category', error);

                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to delete this category');
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
