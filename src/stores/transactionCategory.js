import { defineStore } from 'pinia';

import categoryConstants from '@/consts/category.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';

function loadTransactionCategoryList(state, allCategories) {
    state.allTransactionCategories = allCategories;
    state.allTransactionCategoriesMap = {};

    for (let categoryType in allCategories) {
        if (!Object.prototype.hasOwnProperty.call(allCategories, categoryType)) {
            continue;
        }

        const categories = allCategories[categoryType];

        for (let i = 0; i < categories.length; i++) {
            const category = categories[i];
            state.allTransactionCategoriesMap[category.id] = category;

            for (let j = 0; j < category.subCategories.length; j++) {
                const subCategory = category.subCategories[j];
                state.allTransactionCategoriesMap[subCategory.id] = subCategory;
            }
        }
    }
}

function addCategoryToTransactionCategoryList(state, category) {
    let categoryList = null;

    if (!category.parentId || category.parentId === '0') {
        categoryList = state.allTransactionCategories[category.type];
    } else if (state.allTransactionCategoriesMap[category.parentId]) {
        categoryList = state.allTransactionCategoriesMap[category.parentId].subCategories;
    }

    if (categoryList) {
        categoryList.push(category);
    }

    state.allTransactionCategoriesMap[category.id] = category;
}

function updateCategoryInTransactionCategoryList(state, category) {
    let categoryList = null;

    if (!category.parentId || category.parentId === '0') {
        categoryList = state.allTransactionCategories[category.type];
    } else if (state.allTransactionCategoriesMap[category.parentId]) {
        categoryList = state.allTransactionCategoriesMap[category.parentId].subCategories;
    }

    if (categoryList) {
        for (let i = 0; i < categoryList.length; i++) {
            if (categoryList[i].id === category.id) {
                categoryList.splice(i, 1, category);
                break;
            }
        }
    }

    state.allTransactionCategoriesMap[category.id] = category;
}

function updateCategoryDisplayOrderInCategoryList(state, { category, from, to }) {
    let categoryList = null;

    if (!category.parentId || category.parentId === '0') {
        categoryList = state.allTransactionCategories[category.type];
    } else if (state.allTransactionCategoriesMap[category.parentId]) {
        categoryList = state.allTransactionCategoriesMap[category.parentId].subCategories;
    }

    if (categoryList) {
        categoryList.splice(to, 0, categoryList.splice(from, 1)[0]);
    }
}

function updateCategoryVisibilityInTransactionCategoryList(state, { category, hidden }) {
    if (state.allTransactionCategoriesMap[category.id]) {
        state.allTransactionCategoriesMap[category.id].hidden = hidden;
    }
}

function removeCategoryFromTransactionCategoryList(state, category) {
    let categoryList = null;

    if (!category.parentId || category.parentId === '0') {
        categoryList = state.allTransactionCategories[category.type];
    } else if (state.allTransactionCategoriesMap[category.parentId]) {
        categoryList = state.allTransactionCategoriesMap[category.parentId].subCategories;
    }

    if (categoryList) {
        for (let i = 0; i < categoryList.length; i++) {
            if (categoryList[i].id === category.id) {
                categoryList.splice(i, 1);
                break;
            }
        }
    }

    if (state.allTransactionCategoriesMap[category.id] && state.allTransactionCategoriesMap[category.id].subCategories) {
        const subCategories = state.allTransactionCategoriesMap[category.id].subCategories;

        for (let i = 0; i < subCategories.length; i++) {
            const subCategory = subCategories[i];
            if (state.allTransactionCategoriesMap[subCategory.id]) {
                delete state.allTransactionCategoriesMap[subCategory.id];
            }
        }
    }

    if (state.allTransactionCategoriesMap[category.id]) {
        delete state.allTransactionCategoriesMap[category.id];
    }
}

export const useTransactionCategoriesStore = defineStore('transactionCategories', {
    state: () => ({
        allTransactionCategories: {},
        allTransactionCategoriesMap: {},
        transactionCategoryListStateInvalid: true,
    }),
    actions: {
        updateTransactionCategoryListInvalidState(invalidState) {
            this.transactionCategoryListStateInvalid = invalidState;
        },
        resetTransactionCategories() {
            this.allTransactionCategories = {};
            this.allTransactionCategoriesMap = {};
            this.transactionCategoryListStateInvalid = true;
        },
        loadAllCategories({ force }) {
            const self = this;

            if (!force && !self.transactionCategoryListStateInvalid) {
                return new Promise((resolve) => {
                    resolve(self.allTransactionCategories);
                });
            }

            return new Promise((resolve, reject) => {
                services.getAllTransactionCategories().then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to get category list' });
                        return;
                    }

                    if (!data.result[categoryConstants.allCategoryTypes.Income]) {
                        data.result[categoryConstants.allCategoryTypes.Income] = [];
                    }

                    if (!data.result[categoryConstants.allCategoryTypes.Expense]) {
                        data.result[categoryConstants.allCategoryTypes.Expense] = [];
                    }

                    if (!data.result[categoryConstants.allCategoryTypes.Transfer]) {
                        data.result[categoryConstants.allCategoryTypes.Transfer] = [];
                    }

                    for (let categoryType in data.result) {
                        if (!Object.prototype.hasOwnProperty.call(data.result, categoryType)) {
                            continue;
                        }

                        const categories = data.result[categoryType];

                        for (let i = 0; i < categories.length; i++) {
                            const category = categories[i];

                            if (!category.subCategories) {
                                category.subCategories = [];
                            }
                        }
                    }

                    loadTransactionCategoryList(self, data.result);

                    if (self.transactionCategoryListStateInvalid) {
                        self.updateTransactionCategoryListInvalidState(false);
                    }

                    resolve(data.result);
                }).catch(error => {
                    if (force) {
                        logger.error('failed to force load category list', error);
                    } else {
                        logger.error('failed to load category list', error);
                    }

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to get category list' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        getCategory({ categoryId }) {
            return new Promise((resolve, reject) => {
                services.getTransactionCategory({
                    id: categoryId
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to get category' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to load category info', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to get category' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        saveCategory({ category }) {
            const self = this;

            return new Promise((resolve, reject) => {
                let promise = null;

                if (!category.id) {
                    promise = services.addTransactionCategory(category);
                } else {
                    promise = services.modifyTransactionCategory(category);
                }

                promise.then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (!category.id) {
                            reject({ message: 'Unable to add category' });
                        } else {
                            reject({ message: 'Unable to save category' });
                        }
                        return;
                    }

                    if (!data.result.subCategories) {
                        data.result.subCategories = [];
                    }

                    if (!category.id) {
                        addCategoryToTransactionCategoryList(self, data.result);
                    } else {
                        updateCategoryInTransactionCategoryList(self, data.result);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to save category', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        if (!category.id) {
                            reject({ message: 'Unable to add category' });
                        } else {
                            reject({ message: 'Unable to save category' });
                        }
                    } else {
                        reject(error);
                    }
                });
            });
        },
        addCategories({ categories }) {
            const self = this;

            return new Promise((resolve, reject) => {
                services.addTransactionCategoryBatch({
                    categories: categories
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to add preset categories' });
                        return;
                    }

                    if (!self.transactionCategoryListStateInvalid) {
                        self.updateTransactionCategoryListInvalidState(true);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to add preset categories', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to add preset categories' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        changeCategoryDisplayOrder({ categoryId, from, to }) {
            const self = this;
            const category = self.allTransactionCategoriesMap[categoryId];

            return new Promise((resolve, reject) => {
                if (!category) {
                    reject({ message: 'Unable to move category' });
                    return;
                }

                if (!category.parentId || category.parentId === '0') {
                    if (!self.allTransactionCategories[category.type] ||
                        !self.allTransactionCategories[category.type][to]) {
                        reject({ message: 'Unable to move category' });
                        return;
                    }
                } else {
                    if (!self.allTransactionCategoriesMap[category.parentId].subCategories ||
                        !self.allTransactionCategoriesMap[category.parentId].subCategories[to]) {
                        reject({ message: 'Unable to move category' });
                        return;
                    }
                }

                if (!self.transactionCategoryListStateInvalid) {
                    self.updateTransactionCategoryListInvalidState(true);
                }

                updateCategoryDisplayOrderInCategoryList(self, {
                    category: category,
                    from: from,
                    to: to
                });

                resolve();
            });
        },
        updateCategoryDisplayOrders({ type, parentId }) {
            const self = this;
            const newDisplayOrders = [];

            let categoryList = null;

            if (!parentId || parentId === '0') {
                categoryList = self.allTransactionCategories[type];
            } else if (self.allTransactionCategoriesMap[parentId]) {
                categoryList = self.allTransactionCategoriesMap[parentId].subCategories;
            }

            if (categoryList) {
                for (let i = 0; i < categoryList.length; i++) {
                    newDisplayOrders.push({
                        id: categoryList[i].id,
                        displayOrder: i + 1
                    });
                }
            }

            return new Promise((resolve, reject) => {
                services.moveTransactionCategory({
                    newDisplayOrders: newDisplayOrders
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to move category' });
                        return;
                    }

                    if (self.transactionCategoryListStateInvalid) {
                        self.updateTransactionCategoryListInvalidState(false);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to save categories display order', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to move category' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        hideCategory({ category, hidden }) {
            const self = this;

            return new Promise((resolve, reject) => {
                services.hideTransactionCategory({
                    id: category.id,
                    hidden: hidden
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (hidden) {
                            reject({ message: 'Unable to hide this category' });
                        } else {
                            reject({ message: 'Unable to unhide this category' });
                        }

                        return;
                    }

                    updateCategoryVisibilityInTransactionCategoryList(self, {
                        category: category,
                        hidden: hidden
                    });

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to change category visibility', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        if (hidden) {
                            reject({ message: 'Unable to hide this category' });
                        } else {
                            reject({ message: 'Unable to unhide this category' });
                        }
                    } else {
                        reject(error);
                    }
                });
            });
        },
        deleteCategory({ category, beforeResolve }) {
            const self = this;

            return new Promise((resolve, reject) => {
                services.deleteTransactionCategory({
                    id: category.id
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to delete this category' });
                        return;
                    }

                    if (beforeResolve) {
                        beforeResolve(() => {
                            removeCategoryFromTransactionCategoryList(self, category);
                        });
                    } else {
                        removeCategoryFromTransactionCategoryList(self, category);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to delete category', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to delete this category' });
                    } else {
                        reject(error);
                    }
                });
            });
        }
    }
});
