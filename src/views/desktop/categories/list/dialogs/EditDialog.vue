<template>
    <v-dialog width="600" :persistent="!!persistent" v-model="showState">
        <v-card>
            <v-toolbar color="primary">
                <v-toolbar-title>
                    <span>{{ $t(title) }}</span>
                    <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                </v-toolbar-title>
            </v-toolbar>
            <v-card-text class="px-5 pt-7 pb-5">
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
                            :disabled="loading || submitting"
                            :label="$t('Description')"
                            :placeholder="$t('Your category description (optional)')"
                            v-model="category.comment"
                        />
                    </v-col>
                    <v-col class="pt-0" cols="12" md="12" v-if="editCategoryId">
                        <v-switch inset :disabled="loading || submitting"
                                  :label="$t('Visible')" v-model="category.visible"/>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="gray" :disabled="loading || submitting" @click="cancel">{{ $t('Cancel') }}</v-btn>
                <v-btn :disabled="inputIsEmpty || loading || submitting" @click="save">
                    {{ $t(saveButtonTitle) }}
                    <v-progress-circular indeterminate size="24" class="ml-2" v-if="submitting"></v-progress-circular>
                </v-btn>
            </v-card-actions>
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

export default {
    props: [
        'persistent',
        'show'
    ],
    expose: [
        'open'
    ],
    data() {
        return {
            showState: false,
            editCategoryId: null,
            loading: false,
            category: {
                type: categoryConstants.allCategoryTypes.Income,
                name: '',
                parentId: '0',
                icon: iconConstants.defaultCategoryIconId,
                color: colorConstants.defaultCategoryColor,
                comment: '',
                visible: true
            },
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

            self.category.id = null;
            self.category.type = categoryConstants.allCategoryTypes.Income;
            self.category.parentId = '0';
            self.category.name = '';
            self.category.icon = iconConstants.defaultCategoryIconId;
            self.category.color = colorConstants.defaultCategoryColor;
            self.category.comment = '';
            self.category.visible = true;

            if (options.id) {
                if (options.currentCategory) {
                    self.setCategory(options.currentCategory);
                }

                self.editCategoryId = options.id;
                self.transactionCategoriesStore.getCategory({
                    categoryId: self.editCategoryId
                }).then(category => {
                    self.setCategory(category);
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

            const submitCategory = {
                type: self.category.type,
                name: self.category.name,
                parentId: self.category.parentId,
                icon: self.category.icon,
                color: self.category.color,
                comment: self.category.comment
            };

            if (self.editCategoryId) {
                submitCategory.id = self.category.id;
                submitCategory.hidden = !self.category.visible;
            }

            self.transactionCategoriesStore.saveCategory({
                category: submitCategory
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
        },
        setCategory(category) {
            this.category.id = category.id;
            this.category.type = category.type;
            this.category.parentId = category.parentId;
            this.category.name = category.name;
            this.category.icon = category.icon;
            this.category.color = category.color;
            this.category.comment = category.comment;
            this.category.visible = !category.hidden;
        }
    }
}
</script>
