<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !hasAnyAvailableCategory }" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t(applyText)" :class="{ 'disabled': !hasAnyAvailableCategory }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="combination-list-wrapper margin-vertical skeleton-text"
                  :key="blockIdx" v-for="blockIdx in [ 1, 2 ]" v-if="loading">
            <f7-accordion-item>
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="combination-list-header combination-list-opened">
                            <f7-list-item>
                                <template #title>
                                    <span>Transaction Category</span>
                                    <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
                                </template>
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
                                <ul class="padding-left">
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

        <f7-block class="combination-list-wrapper margin-vertical"
                  :key="transactionType.type"
                  v-for="transactionType in allVisibleTransactionCategories"
                  v-else-if="!loading">
            <f7-accordion-item :opened="collapseStates[transactionType.type].opened"
                               @accordion:open="collapseStates[transactionType.type].opened = true"
                               @accordion:close="collapseStates[transactionType.type].opened = false">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="combination-list-header"
                                 :class="collapseStates[transactionType.type].opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item>
                                <template #title>
                                    <span>{{ getCategoryTypeName(transactionType.type) }}</span>
                                    <f7-icon class="combination-list-chevron-icon" :f7="collapseStates[transactionType.type].opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: collapseStates[transactionType.type].opened ? 'auto' : '' }">
                    <f7-list strong inset dividers accordion-list class="combination-list-content" v-if="!hasAvailableCategory[transactionType.type]">
                        <f7-list-item :title="$t('No available category')"></f7-list-item>
                    </f7-list>
                    <f7-list strong inset dividers accordion-list class="combination-list-content" v-else-if="hasAvailableCategory[transactionType.type]">
                        <f7-list-item checkbox
                                      :class="{ 'has-child-list-item': transactionType.visibleSubCategories[category.id] }"
                                      :title="category.name"
                                      :value="category.id"
                                      :checked="isSubCategoriesAllChecked(category, filterCategoryIds)"
                                      :indeterminate="isSubCategoriesHasButNotAllChecked(category, filterCategoryIds)"
                                      :key="category.id"
                                      v-for="category in transactionType.visibleCategories"
                                      @change="selectSubCategories">
                            <template #media>
                                <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                            </template>

                            <template #root>
                                <ul class="padding-left"
                                    v-if="transactionType.visibleSubCategories[category.id]">
                                    <f7-list-item checkbox
                                                  :title="subCategory.name"
                                                  :value="subCategory.id"
                                                  :checked="isCategoryChecked(subCategory, filterCategoryIds)"
                                                  :key="subCategory.id"
                                                  v-for="subCategory in transactionType.visibleSubCategories[category.id]"
                                                  @change="selectCategory">
                                        <template #media>
                                            <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
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
    props: [
        'f7route',
        'f7router'
    ],
    data: function () {
        const self = this;

        return {
            loading: true,
            loadingError: null,
            modifyDefault: false,
            filterCategoryIds: {},
            collapseStates: self.getCollapseStates(),
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
        applyText() {
            if (this.modifyDefault) {
                return 'Save';
            } else {
                return 'Apply';
            }
        },
        allVisibleTransactionCategories() {
            return this.$utilities.allVisibleTransactionCategories(this.$store.state.allTransactionCategories);
        },
        hasAnyAvailableCategory() {
            for (let type in this.allVisibleTransactionCategories) {
                if (!Object.prototype.hasOwnProperty.call(this.allVisibleTransactionCategories, type)) {
                    continue;
                }

                const categoryType = this.allVisibleTransactionCategories[type];

                if (categoryType.visibleCategories && categoryType.visibleCategories.length > 0) {
                    return true;
                }
            }

            return false;
        },
        hasAvailableCategory() {
            const result = {};

            for (let type in this.allVisibleTransactionCategories) {
                if (!Object.prototype.hasOwnProperty.call(this.allVisibleTransactionCategories, type)) {
                    continue;
                }

                const categoryType = this.allVisibleTransactionCategories[type];
                result[type] = categoryType.visibleCategories && categoryType.visibleCategories.length > 0;
            }

            return result;
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

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
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        save() {
            const self = this;
            const router = self.f7router;

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
        },
        getCategoryTypeName(categoryType) {
            switch (categoryType) {
                case this.$constants.category.allCategoryTypes.Income.toString():
                    return this.$t('Income Categories');
                case this.$constants.category.allCategoryTypes.Expense.toString():
                    return this.$t('Expense Categories');
                case this.$constants.category.allCategoryTypes.Transfer.toString():
                    return this.$t('Transfer Categories');
                default:
                    return this.$t('Transaction Categories');
            }
        },
        isCategoryChecked(category, filterCategoryIds) {
            return !filterCategoryIds[category.id];
        },
        isSubCategoriesAllChecked(category, filterCategoryIds) {
            for (let i = 0; i < category.subCategories.length; i++) {
                const subCategory = category.subCategories[i];
                if (filterCategoryIds[subCategory.id]) {
                    return false;
                }
            }

            return true;
        },
        isSubCategoriesHasButNotAllChecked(category, filterCategoryIds) {
            let checkedCount = 0;

            for (let i = 0; i < category.subCategories.length; i++) {
                const subCategory = category.subCategories[i];
                if (!filterCategoryIds[subCategory.id]) {
                    checkedCount++;
                }
            }

            return checkedCount > 0 && checkedCount < category.subCategories.length;
        },
        getCollapseStates() {
            const collapseStates = {};

            for (let categoryTypeField in this.$constants.category.allCategoryTypes) {
                if (!Object.prototype.hasOwnProperty.call(this.$constants.category.allCategoryTypes, categoryTypeField)) {
                    continue;
                }

                const categoryType = this.$constants.category.allCategoryTypes[categoryTypeField];

                collapseStates[categoryType] = {
                    opened: true
                };
            }

            return collapseStates;
        }
    }
}
</script>
