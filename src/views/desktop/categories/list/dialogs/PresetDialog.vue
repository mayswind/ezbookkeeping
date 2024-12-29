<template>
    <v-dialog width="800" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ $t('Default Categories') }}</h4>
                </div>
            </template>
            <v-card-text class="preset-transaction-categories mt-sm-2 mt-md-4 pt-0">
                <template :key="categoryType" v-for="(categories, categoryType) in allPresetCategories">
                    <div class="d-flex align-center mb-1">
                        <h4>{{ getCategoryTypeName(categoryType) }}</h4>
                        <v-spacer/>
                        <v-menu location="bottom">
                            <template #activator="{ props }">
                                <v-btn variant="text" :disabled="submitting"
                                       v-bind="props">{{ currentLanguageName }}</v-btn>
                            </template>
                            <v-list>
                                <v-list-item :key="lang.languageTag" :value="lang.languageTag" v-for="lang in allLanguages">
                                    <v-list-item-title class="cursor-pointer" @click="currentLocale = lang.languageTag">
                                        {{ lang.displayName }}
                                    </v-list-item-title>
                                </v-list-item>
                            </v-list>
                        </v-menu>
                    </div>

                    <v-expansion-panels class="border rounded mb-2" variant="accordion" multiple :disabled="submitting">
                        <v-expansion-panel :key="idx" v-for="(category, idx) in categories">
                            <v-expansion-panel-title class="py-0">
                                <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                                <span class="ml-3">{{ category.name }}</span>
                            </v-expansion-panel-title>
                            <v-expansion-panel-text v-if="category.subCategories.length">
                                <v-list rounded density="comfortable" class="pa-0">
                                    <template :key="subIdx"
                                              v-for="(subCategory, subIdx) in category.subCategories">
                                        <v-list-item>
                                            <template #prepend>
                                                <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                            </template>
                                            <span class="ml-3">{{ subCategory.name }}</span>
                                        </v-list-item>
                                        <v-divider v-if="subIdx !== category.subCategories.length - 1"/>
                                    </template>
                                </v-list>
                            </v-expansion-panel-text>
                        </v-expansion-panel>
                    </v-expansion-panels>
                </template>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn :disabled="submitting" @click="save">
                        {{ $t('Save') }}
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" density="default" variant="tonal"
                           :disabled="submitting" @click="showState = false">{{ $t('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';

import { CategoryType } from '@/core/category.ts';
import { categorizedArrayToPlainArray } from '@/lib/common.ts';

import {
    mdiDotsVertical
} from '@mdi/js';

export default {
    props: [
        'categoryType',
        'persistent',
        'show'
    ],
    emits: [
        'update:show',
        'category:saved'
    ],
    data() {
        const self = this;

        return {
            currentLocale: self.$locale.getCurrentLanguageTag(),
            allCategoryTypes: [],
            submitting: false,
            icons: {
                more: mdiDotsVertical
            }
        };
    },
    computed: {
        ...mapStores(useTransactionCategoriesStore),
        showState: {
            get: function () {
                return this.show;
            },
            set: function (value) {
                this.$emit('update:show', value);
            }
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfoArray(false);
        },
        allPresetCategories() {
            return this.$locale.getAllTransactionDefaultCategories(this.categoryType, this.currentLocale);
        },
        currentLanguageName() {
            const languageInfo = this.$locale.getLanguageInfo(this.currentLocale);

            if (!languageInfo) {
                return '';
            }

            return languageInfo.displayName;
        }
    },
    methods: {
        save() {
            const self = this;

            self.submitting = true;

            const submitCategories = categorizedArrayToPlainArray(self.allPresetCategories);

            self.transactionCategoriesStore.addCategories({
                categories: submitCategories
            }).then(() => {
                self.submitting = false;
                self.showState = false;

                this.$emit('category:saved', {
                    message: 'You have added preset categories'
                });
            }).catch(error => {
                self.submitting = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        getCategoryTypeName(categoryType) {
            switch (categoryType) {
                case CategoryType.Income.toString():
                    return this.$t('Income Categories');
                case CategoryType.Expense.toString():
                    return this.$t('Expense Categories');
                case CategoryType.Transfer.toString():
                    return this.$t('Transfer Categories');
                default:
                    return this.$t('Transaction Categories');
            }
        }
    }
}
</script>

<style>
.preset-transaction-categories .v-expansion-panel-text__wrapper {
    padding: 0 0 0 20px;
}
</style>
