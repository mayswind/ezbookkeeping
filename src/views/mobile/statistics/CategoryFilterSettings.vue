<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t('Save')" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>
                <small class="card-header-content">
                    <span>Transaction Category</span>
                </small>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item checkbox class="disabled" title="Category Name">
                        <f7-icon slot="media" f7="app_fill"></f7-icon>
                        <ul slot="root" class="padding-left">
                            <f7-list-item checkbox class="disabled" title="Sub Category Name">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </f7-list-item>
                            <f7-list-item checkbox class="disabled" title="Sub Category Name 2">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </f7-list-item>
                            <f7-list-item checkbox class="disabled" title="Sub Category Name 3">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </f7-list-item>
                        </ul>
                    </f7-list-item>
                    <f7-list-item checkbox class="disabled" title="Category Name 2">
                        <f7-icon slot="media" f7="app_fill"></f7-icon>
                        <ul slot="root" class="padding-left">
                            <f7-list-item checkbox class="disabled" title="Sub Category Name">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </f7-list-item>
                            <f7-list-item checkbox class="disabled" title="Sub Category Name 2">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </f7-list-item>
                            <f7-list-item checkbox class="disabled" title="Sub Category Name 3">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </f7-list-item>
                        </ul>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-block class="no-padding no-margin" v-if="!loading">
            <f7-card v-for="(categories, categoryType) in allTransactionCategories" :key="categoryType">
                <f7-card-header>
                    <small class="card-header-content">
                        <span>{{ categoryType | categoryTypeName($constants.category.allCategoryTypes) | localized }}</span>
                    </small>
                </f7-card-header>
                <f7-card-content class="no-safe-areas" :padding="false">
                    <f7-list>
                        <f7-list-item checkbox v-for="category in categories"
                                      v-show="!category.hidden"
                                      :key="category.id"
                                      :title="category.name"
                                      :value="category.id"
                                      :checked="category | subCategoriesAllChecked(filterCategoryIds)"
                                      :indeterminate="category | subCategoriesHasButNotAllChecked(filterCategoryIds)"
                                      @change="selectSubCategories">
                            <f7-icon slot="media"
                                     :icon="category.icon | categoryIcon"
                                     :style="category.color | categoryIconStyle('var(--default-icon-color)')">
                            </f7-icon>

                            <ul slot="root" v-if="category.subCategories.length" class="padding-left">
                                <f7-list-item checkbox v-for="subCategory in category.subCategories"
                                              v-show="!subCategory.hidden"
                                              :key="subCategory.id"
                                              :title="subCategory.name"
                                              :value="subCategory.id"
                                              :checked="subCategory | categoryChecked(filterCategoryIds) "
                                              @change="selectCategory">
                                    <f7-icon slot="media"
                                             :icon="subCategory.icon | categoryIcon"
                                             :style="subCategory.color | categoryIconStyle('var(--default-icon-color)')">
                                    </f7-icon>
                                </f7-list-item>
                            </ul>
                        </f7-list-item>
                    </f7-list>
                </f7-card-content>
            </f7-card>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="selectAll">{{ $t('Select All') }}</f7-actions-button>
                <f7-actions-button @click="selectNone">{{ $t('Select None') }}</f7-actions-button>
                <f7-actions-button @click="selectInvert">{{ $t('Invert Selection') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
export default {
    data: function () {
        return {
            loading: true,
            modifyDefault: false,
            filterCategoryIds: {},
            showMoreActionSheet: false
        }
    },
    computed: {
        title() {
            if (this.modifyDefault) {
                return 'Default Transaction Category Filter';
            } else {
                return 'Filter Transaction Categories';
            }
        },
        allTransactionCategories: function () {
            return this.$store.state.allTransactionCategories;
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;
        const router = self.$f7router;

        self.modifyDefault = !!query.modifyDefault;

        self.$store.dispatch('loadAllCategories', {
            force: false
        }).then(() => {
            self.loading = false;

            const allCategoryIds = {};

            for (let categoryId in self.$store.state.allTransactionCategoriesMap) {
                if (!Object.prototype.hasOwnProperty.call(self.$store.state.allTransactionCategoriesMap, categoryId)) {
                    continue;
                }

                const category = self.$store.state.allTransactionCategoriesMap[categoryId];
                allCategoryIds[category.id] = false;
            }

            if (self.modifyDefault) {
                self.filterCategoryIds = self.$utilities.copyObjectTo(self.$settings.getStatisticsDefaultTransactionCategoryFilter(), allCategoryIds);
            } else {
                self.filterCategoryIds = self.$utilities.copyObjectTo(self.$store.state.transactionStatisticsFilter.filterCategoryIds, allCategoryIds);
            }
        }).catch(error => {
            self.logining = false;

            if (!error.processed) {
                self.$toast(error.message || error);
                router.back();
            }
        });
    },
    methods: {
        save() {
            const self = this;
            const router = self.$f7router;

            const filteredCategoryIds = {};

            for (let categoryId in self.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(self.filterCategoryIds, categoryId)) {
                    continue;
                }

                if (self.filterCategoryIds[categoryId]) {
                    filteredCategoryIds[categoryId] = true;
                }
            }

            if (self.modifyDefault) {
                self.$settings.setStatisticsDefaultTransactionCategoryFilter(filteredCategoryIds);
            } else {
                self.$store.dispatch('updateTransactionStatisticsFilter', {
                    filterCategoryIds: filteredCategoryIds
                });
            }

            router.back();
        },
        selectCategory(e) {
            const categoryId = e.target.value;
            const category = this.$store.state.allTransactionCategoriesMap[categoryId];

            if (!category) {
                return;
            }

            this.filterCategoryIds[category.id] = !e.target.checked;
        },
        selectSubCategories(e) {
            const categoryId = e.target.value;
            const category = this.$store.state.allTransactionCategoriesMap[categoryId];

            if (!category || !category.subCategories || !category.subCategories.length) {
                return;
            }

            for (let i = 0; i < category.subCategories.length; i++) {
                const subCategory = category.subCategories[i];
                this.filterCategoryIds[subCategory.id] = !e.target.checked;
            }
        },
        selectAll() {
            for (let categoryId in this.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterCategoryIds, categoryId)) {
                    continue;
                }

                const category = this.$store.state.allTransactionCategoriesMap[categoryId];

                if (category) {
                    this.filterCategoryIds[category.id] = false;
                }
            }
        },
        selectNone() {
            for (let categoryId in this.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterCategoryIds, categoryId)) {
                    continue;
                }

                const category = this.$store.state.allTransactionCategoriesMap[categoryId];

                if (category) {
                    this.filterCategoryIds[category.id] = true;
                }
            }
        },
        selectInvert() {
            for (let categoryId in this.filterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(this.filterCategoryIds, categoryId)) {
                    continue;
                }

                const category = this.$store.state.allTransactionCategoriesMap[categoryId];

                if (category) {
                    this.filterCategoryIds[category.id] = !this.filterCategoryIds[category.id];
                }
            }
        }
    },
    filters: {
        categoryTypeName(categoryType, allCategoryTypes) {
            switch (categoryType) {
                case allCategoryTypes.Income.toString():
                    return 'Income Categories';
                case allCategoryTypes.Expense.toString():
                    return 'Expense Categories';
                case allCategoryTypes.Transfer.toString():
                    return 'Transfer Categories';
                default:
                    return 'Transaction Categories';
            }
        },
        categoryChecked(category, filterCategoryIds) {
            return !filterCategoryIds[category.id];
        },
        subCategoriesAllChecked(category, filterCategoryIds) {
            for (let i = 0; i < category.subCategories.length; i++) {
                const subCategory = category.subCategories[i];
                if (filterCategoryIds[subCategory.id]) {
                    return false;
                }
            }

            return true;
        },
        subCategoriesHasButNotAllChecked(category, filterCategoryIds) {
            let checkedCount = 0;

            for (let i = 0; i < category.subCategories.length; i++) {
                const subCategory = category.subCategories[i];
                if (!filterCategoryIds[subCategory.id]) {
                    checkedCount++;
                }
            }

            return checkedCount > 0 && checkedCount < category.subCategories.length;
        }
    }
}
</script>
