<template>
    <v-dialog width="800" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h5 class="text-h5">{{ $t(title) }}</h5>
                    <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                </div>
            </template>
            <v-card-text>
                <v-form class="mt-2 mt-md-6">
                    <v-row>
                        <v-col cols="12" md="12">
                            <v-text-field
                                type="text"
                                clearable
                                persistent-placeholder
                                :disabled="loading || submitting"
                                :label="$t('Category Name')"
                                :placeholder="$t('Category Name')"
                                v-model="category.name"
                            />
                        </v-col>
                        <v-col cols="12" md="6">
                            <icon-select icon-type="category"
                                         :all-icon-infos="allCategoryIcons"
                                          :label="$t('Category Icon')"
                                          :color="category.color"
                                          :disabled="loading || submitting"
                                          v-model="category.icon" />
                        </v-col>
                        <v-col cols="12" md="6">
                            <color-select :all-color-infos="allCategoryColors"
                                         :label="$t('Category Color')"
                                         :disabled="loading || submitting"
                                         v-model="category.color" />
                        </v-col>
                        <v-col cols="12" md="12">
                            <v-textarea
                                type="text"
                                persistent-placeholder
                                rows="3"
                                :disabled="loading || submitting"
                                :label="$t('Description')"
                                :placeholder="$t('Your category description (optional)')"
                                v-model="category.comment"
                            />
                        </v-col>
                        <v-col class="py-0" cols="12" md="12" v-if="editCategoryId">
                            <v-switch inset :disabled="loading || submitting"
                                      :label="$t('Visible')" v-model="category.visible"/>
                        </v-col>
                    </v-row>
                </v-form>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn :disabled="inputIsEmpty || loading || submitting" @click="save">
                        {{ $t(saveButtonTitle) }}
                        <v-progress-circular indeterminate size="24" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal"
                           :disabled="loading || submitting" @click="cancel">{{ $t('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';

import categoryConstants from '@/consts/category.js';
import iconConstants from '@/consts/icon.js';
import colorConstants from '@/consts/color.js';
import { setCategoryModelByAnotherCategory } from '@/lib/category.js';

export default {
    props: [
        'persistent',
        'show'
    ],
    expose: [
        'open'
    ],
    data() {
        const transactionCategoriesStore = useTransactionCategoriesStore();
        const newTransactionCategory = transactionCategoriesStore.generateNewTransactionCategoryModel();

        return {
            showState: false,
            editCategoryId: null,
            loading: false,
            category: newTransactionCategory,
            submitting: false,
            resolve: null,
            reject: null
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
    methods: {
        open(options) {
            const self = this;
            self.showState = true;
            self.loading = true;
            self.submitting = false;

            const newTransactionCategory = self.transactionCategoriesStore.generateNewTransactionCategoryModel();
            setCategoryModelByAnotherCategory(self.category, newTransactionCategory);

            if (options.id) {
                if (options.currentCategory) {
                    setCategoryModelByAnotherCategory(self.category, options.currentCategory);
                }

                self.editCategoryId = options.id;
                self.transactionCategoriesStore.getCategory({
                    categoryId: self.editCategoryId
                }).then(category => {
                    setCategoryModelByAnotherCategory(self.category, category);
                    self.loading = false;
                }).catch(error => {
                    self.loading = false;
                    self.showState = false;

                    if (!error.processed) {
                        if (self.reject) {
                            self.reject(error);
                        }
                    }
                });
            } else if (options.parentId) {
                self.editCategoryId = null;

                const categoryType = parseInt(options.type);

                if (categoryType !== categoryConstants.allCategoryTypes.Income &&
                    categoryType !== categoryConstants.allCategoryTypes.Expense &&
                    categoryType !== categoryConstants.allCategoryTypes.Transfer) {
                    self.loading = false;
                    self.showState = false;

                    if (self.reject) {
                        self.reject('Parameter Invalid');
                    }

                    return;
                }

                self.category.type = categoryType;
                self.category.parentId = options.parentId;

                self.loading = false;
            }

            return new Promise((resolve, reject) => {
                self.resolve = resolve;
                self.reject = reject;
            });
        },
        save() {
            const self = this;

            const problemMessage = self.inputEmptyProblemMessage;

            if (problemMessage) {
                self.$refs.snackbar.showMessage(problemMessage);
                return;
            }

            self.submitting = true;

            self.transactionCategoriesStore.saveCategory({
                category: self.category,
                isEdit: !!self.editCategoryId
            }).then(() => {
                self.submitting = false;

                let message = 'You have saved this category';

                if (!self.editCategoryId) {
                    message = 'You have added a new category';
                }

                if (self.resolve) {
                    self.resolve({
                        message: message
                    });
                }

                self.showState = false;
            }).catch(error => {
                self.submitting = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        cancel() {
            if (this.reject) {
                this.reject();
            }

            this.showState = false;
        }
    }
}
</script>
