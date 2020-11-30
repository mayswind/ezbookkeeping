<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Default Categories')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :text="$t('Save')" :class="{ 'disabled': saving }" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card v-for="categoryInfo in allCategories" :key="categoryInfo.type">
            <f7-card-header>
                <small :style="{ opacity: 0.6 }">
                    <span>{{ categoryInfo.type | categoryTypeName | t }}</span>
                </small>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item v-for="(category, idx) in categoryInfo.categories"
                                  :key="idx"
                                  :accordion-item="!!category.subCategories.length"
                                  :title="$t('category.' + category.name)">
                        <f7-icon slot="media"
                                 :icon="category.categoryIconId | categoryIcon"
                                 :style="{ color: '#' + category.color }">
                        </f7-icon>

                        <f7-accordion-content v-if="category.subCategories.length" class="padding-left">
                            <f7-list>
                                <f7-list-item v-for="(subCategory, subIdx) in category.subCategories"
                                              :key="subIdx"
                                              :title="$t('category.' + subCategory.name)">
                                    <f7-icon slot="media"
                                             :icon="subCategory.categoryIconId | categoryIcon"
                                             :style="{ color: '#' + subCategory.color }">
                                    </f7-icon>
                                </f7-list-item>
                            </f7-list>
                        </f7-accordion-content>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            categoryType: '',
            allCategories: [],
            saving: false
        };
    },
    created() {
        const self = this;
        const query = self.$f7route.query;
        const router = self.$f7router;

        if (query.type !== '0' && query.type !== '1' && query.type !== '2' && query.type !== '3') {
            self.$toast('Parameter Invalid');
            router.back();
            return;
        }

        self.categoryType = query.type;

        if (query.type === '0') {
            for (let i = 1; i <= 3; i++) {
                self.allCategories.push({
                    type: i.toString(),
                    categories: self.$utilities.copyArrayTo(self.getDefaultCategories(i.toString()), [])
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
                case '1':
                    return this.$constants.category.defaultExpenseCategories;
                case '2':
                    return this.$constants.category.defaultIncomeCategories;
                case '3':
                    return this.$constants.category.defaultTransferCategories;
                default:
                    return [];
            }
        },
        save() {

        }
    },
    filters: {
        categoryTypeName(categoryType) {
            switch (categoryType) {
                case '1':
                    return 'Expense Categories';
                case '2':
                    return 'Income Categories';
                case '3':
                    return 'Transfer Categories';
                default:
                    return 'Transaction Categories';
            }
        }
    }
};
</script>
