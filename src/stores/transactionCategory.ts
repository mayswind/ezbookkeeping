import { ref } from 'vue';
import { defineStore } from 'pinia';

import type { BeforeResolveFunction } from '@/core/base.ts';

import { CategoryType } from '@/core/category.ts';

import {
    type TransactionCategoryInfoResponse,
    type TransactionCategoryCreateBatchRequest,
    type TransactionCategoryNewDisplayOrderRequest,
    TransactionCategory,
} from '@/models/transaction_category.ts';

import { isEquals } from '@/lib/common.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useTransactionCategoriesStore = defineStore('transactionCategories', () =>{
    const allTransactionCategories = ref<Record<number, TransactionCategory[]>>({});
    const allTransactionCategoriesMap = ref<Record<string, TransactionCategory>>({});
    const transactionCategoryListStateInvalid = ref<boolean>(true);

    function loadTransactionCategoryList(allCategories: Record<number, TransactionCategory[]>): void {
        allTransactionCategories.value = allCategories;
        allTransactionCategoriesMap.value = {};

        for (const categoryType in allCategories) {
            if (!Object.prototype.hasOwnProperty.call(allCategories, categoryType)) {
                continue;
            }

            const categories = allCategories[categoryType];

            for (let i = 0; i < categories.length; i++) {
                const category = categories[i];
                allTransactionCategoriesMap.value[category.id] = category;

                if (!category.secondaryCategories) {
                    continue;
                }

                for (let j = 0; j < category.secondaryCategories.length; j++) {
                    const subCategory = category.secondaryCategories[j];
                    allTransactionCategoriesMap.value[subCategory.id] = subCategory;
                }
            }
        }
    }

    function addCategoryToTransactionCategoryList(category: TransactionCategory): void {
        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!category.parentId || category.parentId === '0') {
            categoryList = allTransactionCategories.value[category.type];
        } else if (allTransactionCategoriesMap.value[category.parentId]) {
            categoryList = allTransactionCategoriesMap.value[category.parentId].secondaryCategories;
        }

        if (categoryList) {
            categoryList.push(category);
        }

        allTransactionCategoriesMap.value[category.id] = category;
    }

    function updateCategoryInTransactionCategoryList(category: TransactionCategory, oldCategory: TransactionCategory): boolean {
        if (oldCategory && category.parentId !== oldCategory.parentId) {
            return false;
        }

        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!category.parentId || category.parentId === '0') {
            categoryList = allTransactionCategories.value[category.type];
        } else if (allTransactionCategoriesMap.value[category.parentId]) {
            categoryList = allTransactionCategoriesMap.value[category.parentId].secondaryCategories;
        }

        if (categoryList) {
            for (let i = 0; i < categoryList.length; i++) {
                if (categoryList[i].id === category.id) {
                    if (!category.parentId || category.parentId === '0') {
                        category.secondaryCategories = categoryList[i].secondaryCategories;
                    }

                    categoryList.splice(i, 1, category);
                    break;
                }
            }
        }

        allTransactionCategoriesMap.value[category.id] = category;
        return true;
    }

    function updateCategoryDisplayOrderInCategoryList(params: { category: TransactionCategory, from: number, to: number }): void {
        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!params.category.parentId || params.category.parentId === '0') {
            categoryList = allTransactionCategories.value[params.category.type];
        } else if (allTransactionCategoriesMap.value[params.category.parentId]) {
            categoryList = allTransactionCategoriesMap.value[params.category.parentId].secondaryCategories;
        }

        if (categoryList) {
            categoryList.splice(params.to, 0, categoryList.splice(params.from, 1)[0]);
        }
    }

    function updateCategoryVisibilityInTransactionCategoryList(params: { category: TransactionCategory, hidden: boolean }): void {
        if (allTransactionCategoriesMap.value[params.category.id]) {
            allTransactionCategoriesMap.value[params.category.id].visible = !params.hidden;
        }
    }

    function removeCategoryFromTransactionCategoryList(category: TransactionCategory): void {
        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!category.parentId || category.parentId === '0') {
            categoryList = allTransactionCategories.value[category.type];
        } else if (allTransactionCategoriesMap.value[category.parentId]) {
            categoryList = allTransactionCategoriesMap.value[category.parentId].secondaryCategories;
        }

        if (categoryList) {
            for (let i = 0; i < categoryList.length; i++) {
                if (categoryList[i].id === category.id) {
                    categoryList.splice(i, 1);
                    break;
                }
            }
        }

        if (allTransactionCategoriesMap.value[category.id] && allTransactionCategoriesMap.value[category.id].secondaryCategories) {
            const subCategoryList = allTransactionCategoriesMap.value[category.id].secondaryCategories;

            if (subCategoryList) {
                for (let i = 0; i < subCategoryList.length; i++) {
                    const subCategory = subCategoryList[i];
                    if (allTransactionCategoriesMap.value[subCategory.id]) {
                        delete allTransactionCategoriesMap.value[subCategory.id];
                    }
                }
            }
        }

        if (allTransactionCategoriesMap.value[category.id]) {
            delete allTransactionCategoriesMap.value[category.id];
        }
    }

    function updateTransactionCategoryListInvalidState(invalidState: boolean): void {
        transactionCategoryListStateInvalid.value = invalidState;
    }

    function resetTransactionCategories(): void {
        allTransactionCategories.value = {};
        allTransactionCategoriesMap.value = {};
        transactionCategoryListStateInvalid.value = true;
    }

    function loadAllCategories(params: { force?: boolean }): Promise<Record<number, TransactionCategory[]>> {
        if (!params.force && !transactionCategoryListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allTransactionCategories.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllTransactionCategories().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve category list' });
                    return;
                }

                if (!data.result[CategoryType.Income]) {
                    data.result[CategoryType.Income] = [];
                }

                if (!data.result[CategoryType.Expense]) {
                    data.result[CategoryType.Expense] = [];
                }

                if (!data.result[CategoryType.Transfer]) {
                    data.result[CategoryType.Transfer] = [];
                }

                if (transactionCategoryListStateInvalid.value) {
                    updateTransactionCategoryListInvalidState(false);
                }

                const transactionCategories = TransactionCategory.ofMap(data.result);

                if (params.force && data.result && isEquals(allTransactionCategories.value, transactionCategories)) {
                    reject({ message: 'Category list is up to date' });
                    return;
                }

                loadTransactionCategoryList(transactionCategories);

                resolve(transactionCategories);
            }).catch(error => {
                if (params.force) {
                    logger.error('failed to force load category list', error);
                } else {
                    logger.error('failed to load category list', error);
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve category list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getCategory(params: { categoryId: string }): Promise<TransactionCategory> {
        return new Promise((resolve, reject) => {
            services.getTransactionCategory({
                id: params.categoryId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve category' });
                    return;
                }

                const transactionCategory = TransactionCategory.of(data.result);

                resolve(transactionCategory);
            }).catch(error => {
                logger.error('failed to load category info', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve category' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function saveCategory(params: { category: TransactionCategory, isEdit: boolean, clientSessionId: string }): Promise<TransactionCategory> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<TransactionCategoryInfoResponse>;

            if (!params.isEdit) {
                promise = services.addTransactionCategory(params.category.toCreateRequest(params.clientSessionId));
            } else {
                promise = services.modifyTransactionCategory(params.category.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!params.isEdit) {
                        reject({ message: 'Unable to add category' });
                    } else {
                        reject({ message: 'Unable to save category' });
                    }
                    return;
                }

                const transactionCategory = TransactionCategory.of(data.result);

                if (!params.isEdit) {
                    addCategoryToTransactionCategoryList(transactionCategory);
                } else {
                    const result = updateCategoryInTransactionCategoryList(transactionCategory, allTransactionCategoriesMap.value[params.category.id]);

                    if (!result && !transactionCategoryListStateInvalid.value) {
                        updateTransactionCategoryListInvalidState(true);
                    }
                }

                resolve(transactionCategory);
            }).catch(error => {
                logger.error('failed to save category', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!params.isEdit) {
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

    function addCategories(req: TransactionCategoryCreateBatchRequest): Promise<Record<number, TransactionCategory[]>> {
        return new Promise((resolve, reject) => {
            services.addTransactionCategoryBatch(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to add preset categories' });
                    return;
                }

                if (!transactionCategoryListStateInvalid.value) {
                    updateTransactionCategoryListInvalidState(true);
                }

                const transactionCategories = TransactionCategory.ofMap(data.result);

                resolve(transactionCategories);
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

    function changeCategoryDisplayOrder(params: { categoryId: string, from: number, to: number }): Promise<void> {
        const category = allTransactionCategoriesMap.value[params.categoryId];

        return new Promise((resolve, reject) => {
            if (!category) {
                reject({ message: 'Unable to move category' });
                return;
            }

            if (!category.parentId || category.parentId === '0') {
                if (!allTransactionCategories.value[category.type] ||
                    !allTransactionCategories.value[category.type][params.to]) {
                    reject({ message: 'Unable to move category' });
                    return;
                }
            } else {
                const subCategoryList = allTransactionCategoriesMap.value[category.parentId].secondaryCategories;

                if (!subCategoryList || !subCategoryList[params.to]) {
                    reject({ message: 'Unable to move category' });
                    return;
                }
            }

            if (!transactionCategoryListStateInvalid.value) {
                updateTransactionCategoryListInvalidState(true);
            }

            updateCategoryDisplayOrderInCategoryList({
                category: category,
                from: params.from,
                to: params.to
            });

            resolve();
        });
    }

    function updateCategoryDisplayOrders(params: { type: CategoryType, parentId: string }): Promise<boolean> {
        const newDisplayOrders: TransactionCategoryNewDisplayOrderRequest[] = [];

        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!params.parentId || params.parentId === '0') {
            categoryList = allTransactionCategories.value[params.type];
        } else if (allTransactionCategoriesMap.value[params.parentId]) {
            categoryList = allTransactionCategoriesMap.value[params.parentId].secondaryCategories;
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

                if (transactionCategoryListStateInvalid.value) {
                    updateTransactionCategoryListInvalidState(false);
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
    }

    function hideCategory(params: { category: TransactionCategory, hidden: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.hideTransactionCategory({
                id: params.category.id,
                hidden: params.hidden
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (params.hidden) {
                        reject({ message: 'Unable to hide this category' });
                    } else {
                        reject({ message: 'Unable to unhide this category' });
                    }

                    return;
                }

                updateCategoryVisibilityInTransactionCategoryList({
                    category: params.category,
                    hidden: params.hidden
                });

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to change category visibility', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (params.hidden) {
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

    function deleteCategory(params: { category: TransactionCategory, beforeResolve: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteTransactionCategory({
                id: params.category.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this category' });
                    return;
                }

                if (params.beforeResolve) {
                    params.beforeResolve(() => {
                        removeCategoryFromTransactionCategoryList(params.category);
                    });
                } else {
                    removeCategoryFromTransactionCategoryList(params.category);
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

    return {
        // states
        allTransactionCategories,
        allTransactionCategoriesMap,
        transactionCategoryListStateInvalid,
        // functions
        updateTransactionCategoryListInvalidState,
        resetTransactionCategories,
        loadAllCategories,
        getCategory,
        saveCategory,
        addCategories,
        changeCategoryDisplayOrder,
        updateCategoryDisplayOrders,
        hideCategory,
        deleteCategory
    };
});
