<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Default Categories')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" v-if="allCategories && allCategories.length" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t('Save')" :class="{ 'disabled': submitting }" v-if="allCategories && allCategories.length" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="no-padding no-margin" :key="categoryInfo.type" v-for="categoryInfo in allCategories">
            <f7-block-title class="margin-top margin-horizontal">{{ getCategoryTypeName(categoryInfo.type) }}</f7-block-title>

            <f7-list strong inset dividers class="margin-top">
                <f7-list-item v-for="(category, idx) in categoryInfo.categories"
                              :key="idx"
                              :accordion-item="!!category.subCategories.length"
                              :title="$t('category.' + category.name, currentLocale)">
                    <template #media>
                        <ItemIcon icon-type="category" :icon-id="category.categoryIconId" :color="category.color"></ItemIcon>
                    </template>

                    <f7-accordion-content v-if="category.subCategories.length" class="padding-left">
                        <f7-list>
                            <f7-list-item v-for="(subCategory, subIdx) in category.subCategories"
                                          :key="subIdx"
                                          :title="$t('category.' + subCategory.name, currentLocale)">

                                <template #media>
                                    <ItemIcon icon-type="category" :icon-id="subCategory.categoryIconId" :color="subCategory.color"></ItemIcon>
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
export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        const self = this;

        return {
            loadingError: null,
            currentLocale: self.$i18n.locale,
            categoryType: 0,
            allCategories: [],
            submitting: false,
            showMoreActionSheet: false,
            showChangeLocaleSheet: false
        };
    },
    computed: {
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        self.categoryType = parseInt(query.type);

        if (self.categoryType !== 0 &&
            self.categoryType !== this.$constants.category.allCategoryTypes.Income &&
            self.categoryType !== this.$constants.category.allCategoryTypes.Expense &&
            self.categoryType !== this.$constants.category.allCategoryTypes.Transfer) {
            self.$toast('Parameter Invalid');
            self.loadingError = 'Parameter Invalid';
            return;
        }

        if (self.categoryType === 0) {
            for (let i = 1; i <= 3; i++) {
                self.allCategories.push({
                    type: i,
                    categories: self.$utilities.copyArrayTo(self.getDefaultCategories(i), [])
                });
            }
        } else {
            self.allCategories.push({
                type: self.categoryType,
                categories: self.$utilities.copyArrayTo(self.getDefaultCategories(self.categoryType), [])
            });
        }
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        getDefaultCategories(categoryType) {
            switch (categoryType) {
                case this.$constants.category.allCategoryTypes.Income:
                    return this.$constants.category.defaultIncomeCategories;
                case this.$constants.category.allCategoryTypes.Expense:
                    return this.$constants.category.defaultExpenseCategories;
                case this.$constants.category.allCategoryTypes.Transfer:
                    return this.$constants.category.defaultTransferCategories;
                default:
                    return [];
            }
        },
        save() {
            const self = this;
            const router = self.f7router;

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            const categories = [];

            for (let i = 0; i < self.allCategories.length; i++) {
                const categoryInfo = self.allCategories[i];

                for (let j = 0; j < categoryInfo.categories.length; j++) {
                    const category = categoryInfo.categories[j];
                    const submitCategory = {
                        name: self.$t('category.' + category.name, self.currentLocale),
                        type: categoryInfo.type,
                        icon: category.categoryIconId,
                        color: category.color,
                        subCategories: []
                    }

                    for (let k = 0; k < category.subCategories.length; k++) {
                        const subCategory = category.subCategories[k];
                        submitCategory.subCategories.push({
                            name: self.$t('category.' + subCategory.name, self.currentLocale),
                            type: categoryInfo.type,
                            icon: subCategory.categoryIconId,
                            color: subCategory.color
                        });
                    }

                    categories.push(submitCategory);
                }
            }

            self.$store.dispatch('addCategories', {
                categories: categories
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
                case this.$constants.category.allCategoryTypes.Income:
                    return this.$t('Income Categories');
                case this.$constants.category.allCategoryTypes.Expense:
                    return this.$t('Expense Categories');
                case this.$constants.category.allCategoryTypes.Transfer:
                    return this.$t('Transfer Categories');
                default:
                    return this.$t('Transaction Categories');
            }
        }
    }
};
</script>
