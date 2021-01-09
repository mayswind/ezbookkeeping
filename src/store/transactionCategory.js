import categoryContants from '../consts/category.js';
import services from '../lib/services.js';
import logger from '../lib/logger.js';

import {
    LOAD_TRANSACTION_CATEGORY_LIST,
    ADD_CATEGORY_TO_TRANSACTION_CATEGORY_LIST,
    SAVE_CATEGORY_IN_TRANSACTION_CATEGORY_LIST,
    CHANGE_CATEGORY_DISPLAY_ORDER_IN_CATEGORY_LIST,
    UPDATE_CATEGORY_VISIBILITY_IN_TRANSACTION_CATEGORY_LIST,
    REMOVE_CATEGORY_FROM_TRANSACTION_CATEGORYLIST,
    UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE,
} from './mutations.js';

export function loadAllCategories(context, { force }) {
    if (!force && !context.state.transactionCategoryListStateInvalid) {
        return new Promise((resolve) => {
            resolve(context.state.allTransactionCategories);
        });
    }

    return new Promise((resolve, reject) => {
        services.getAllTransactionCategories().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get category list' });
                return;
            }

            if (!data.result[categoryContants.allCategoryTypes.Income]) {
                data.result[categoryContants.allCategoryTypes.Income] = [];
            }

            if (!data.result[categoryContants.allCategoryTypes.Expense]) {
                data.result[categoryContants.allCategoryTypes.Expense] = [];
            }

            if (!data.result[categoryContants.allCategoryTypes.Transfer]) {
                data.result[categoryContants.allCategoryTypes.Transfer] = [];
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

            context.commit(LOAD_TRANSACTION_CATEGORY_LIST, data.result);
            context.commit(UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE, false);

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
}

export function getCategory(context, { categoryId }) {
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
}

export function saveCategory(context, { category }) {
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
                context.commit(ADD_CATEGORY_TO_TRANSACTION_CATEGORY_LIST, data.result);
            } else {
                context.commit(SAVE_CATEGORY_IN_TRANSACTION_CATEGORY_LIST, data.result);
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
}

export function addCategories(context, { categories }) {
    return new Promise((resolve, reject) => {
        services.addTransactionCategoryBatch({
            categories: categories
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to add preset categories' });
                return;
            }

            context.commit(UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE, true);

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
}

export function changeCategoryDisplayOrder(context, { categoryId, from, to }) {
    const category = context.state.allTransactionCategoriesMap[categoryId];

    return new Promise((resolve, reject) => {
        if (!category) {
            reject({ message: 'Unable to move category' });
            return;
        }

        if (!category.parentId || category.parentId === '0') {
            if (!context.state.allTransactionCategories[category.type] ||
                !context.state.allTransactionCategories[category.type][to]) {
                reject({ message: 'Unable to move category' });
                return;
            }
        } else {
            if (!context.state.allTransactionCategoriesMap[category.parentId].subCategories ||
                !context.state.allTransactionCategoriesMap[category.parentId].subCategories[to]) {
                reject({ message: 'Unable to move category' });
                return;
            }
        }

        context.commit(UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE, true);
        context.commit(CHANGE_CATEGORY_DISPLAY_ORDER_IN_CATEGORY_LIST, {
            category: category,
            from: from,
            to: to
        });

        resolve();
    });
}

export function updateCategoryDisplayOrders(context, { type, parentId }) {
    const newDisplayOrders = [];

    let categoryList = null;

    if (!parentId || parentId === '0') {
        categoryList = context.state.allTransactionCategories[type];
    } else if (context.state.allTransactionCategoriesMap[parentId]) {
        categoryList = context.state.allTransactionCategoriesMap[parentId].subCategories;
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

            context.commit(UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE, false);

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
}

export function hideCategory(context, { category, hidden }) {
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

            context.commit(UPDATE_CATEGORY_VISIBILITY_IN_TRANSACTION_CATEGORY_LIST, {
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
}

export function deleteCategory(context, { category, beforeResolve }) {
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
                    context.commit(REMOVE_CATEGORY_FROM_TRANSACTION_CATEGORYLIST, category);
                });
            } else {
                context.commit(REMOVE_CATEGORY_FROM_TRANSACTION_CATEGORYLIST, category);
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
