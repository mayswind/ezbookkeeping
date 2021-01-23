<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Default Categories')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :text="$t('Save')" :class="{ 'disabled': submitting }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card v-for="categoryInfo in allCategories" :key="categoryInfo.type">
            <f7-card-header>
                <small :style="{ opacity: 0.6 }">
                    <span>{{ categoryInfo.type | categoryTypeName($constants.category.allCategoryTypes) | localized }}</span>
                </small>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item v-for="(category, idx) in categoryInfo.categories"
                                  :key="idx"
                                  :accordion-item="!!category.subCategories.length"
                                  :title="$t('category.' + category.name, currentLocale)">
                        <f7-icon slot="media"
                                 :icon="category.categoryIconId | categoryIcon"
                                 :style="category.color | categoryIconStyle('var(--default-icon-color)')">
                        </f7-icon>

                        <f7-accordion-content v-if="category.subCategories.length" class="padding-left">
                            <f7-list>
                                <f7-list-item v-for="(subCategory, subIdx) in category.subCategories"
                                              :key="subIdx"
                                              :title="$t('category.' + subCategory.name, currentLocale)">
                                    <f7-icon slot="media"
                                             :icon="subCategory.categoryIconId | categoryIcon"
                                             :style="subCategory.color | categoryIconStyle('var(--default-icon-color)')">
                                    </f7-icon>
                                </f7-list-item>
                            </f7-list>
                        </f7-accordion-content>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

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
                                   :show.sync="showChangeLocaleSheet"
                                   v-model="currentLocale">
        </list-item-selection-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
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
            return this.$locale.getAllLanguages();
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;
        const router = self.$f7router;

        self.categoryType = parseInt(query.type);

        if (self.categoryType !== 0 &&
            self.categoryType !== this.$constants.category.allCategoryTypes.Income &&
            self.categoryType !== this.$constants.category.allCategoryTypes.Expense &&
            self.categoryType !== this.$constants.category.allCategoryTypes.Transfer) {
            self.$toast('Parameter Invalid');
            router.back();
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
            const router = self.$f7router;

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
        }
    },
    filters: {
        categoryTypeName(categoryType, allCategoryTypes) {
            switch (categoryType) {
                case allCategoryTypes.Income:
                    return 'Income Categories';
                case allCategoryTypes.Expense:
                    return 'Expense Categories';
                case allCategoryTypes.Transfer:
                    return 'Transfer Categories';
                default:
                    return 'Transaction Categories';
            }
        }
    }
};
</script>
