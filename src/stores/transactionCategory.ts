import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import type { BeforeResolveFunction } from '@/core/base.ts';

import { itemAndIndex, values } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';

import {
    type TransactionCategoryInfoResponse,
    type TransactionCategoryCreateBatchRequest,
    type TransactionCategoryNewDisplayOrderRequest,
    TransactionCategory,
} from '@/models/transaction_category.ts';

import { isEquals } from '@/lib/common.ts';
import { getFirstVisibleCategoryId } from '@/lib/category.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useTransactionCategoriesStore = defineStore('transactionCategories', () =>{
    const allTransactionCategories = ref<Record<number, TransactionCategory[]>>({});
    const allTransactionCategoriesMap = ref<Record<string, TransactionCategory>>({});
    const transactionCategoryListStateInvalid = ref<boolean>(true);

    const allAvailablePrimaryCategoriesCount = computed<number>(() => {
        let count = 0;

        for (const categories of values(allTransactionCategories.value)) {
            count += categories.length;
        }

        return count;
    });

    const allAvailableSecondaryCategoriesCount = computed<number>(() => {
        let count = 0;

        for (const categories of values(allTransactionCategories.value)) {
            for (const category of categories) {
                if (category.subCategories) {
                    count += category.subCategories.length;
                }
            }
        }

        return count;
    });

    const hasVisibleExpenseCategories = computed<boolean>(() => {
        if (!allTransactionCategories.value || !allTransactionCategories.value[CategoryType.Expense] || !allTransactionCategories.value[CategoryType.Expense].length) {
            return false;
        }

        const firstVisibleCategoryId = getFirstVisibleCategoryId(allTransactionCategories.value[CategoryType.Expense]);
        return firstVisibleCategoryId !== '';
    });

    const hasVisibleIncomeCategories = computed<boolean>(() => {
        if (!allTransactionCategories.value || !allTransactionCategories.value[CategoryType.Income] || !allTransactionCategories.value[CategoryType.Income].length) {
            return false;
        }

        const firstVisibleCategoryId = getFirstVisibleCategoryId(allTransactionCategories.value[CategoryType.Income]);
        return firstVisibleCategoryId !== '';
    });

    const hasVisibleTransferCategories = computed<boolean>(() => {
        if (!allTransactionCategories.value || !allTransactionCategories.value[CategoryType.Transfer] || !allTransactionCategories.value[CategoryType.Transfer].length) {
            return false;
        }

        const firstVisibleCategoryId = getFirstVisibleCategoryId(allTransactionCategories.value[CategoryType.Transfer]);
        return firstVisibleCategoryId !== '';
    });

    function loadTransactionCategoryList(allCategories: Record<number, TransactionCategory[]>): void {
        allTransactionCategories.value = allCategories;
        allTransactionCategoriesMap.value = {};

        for (const categories of values(allCategories)) {
            for (const category of categories) {
                allTransactionCategoriesMap.value[category.id] = category;

                if (!category.subCategories) {
                    continue;
                }

                for (const subCategory of category.subCategories) {
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
            categoryList = allTransactionCategoriesMap.value[category.parentId]!.subCategories;
        }

        if (categoryList) {
            categoryList.push(category);
        }

        allTransactionCategoriesMap.value[category.id] = category;
    }

    function updateCategoryInTransactionCategoryList(currentCategory: TransactionCategory, oldCategory?: TransactionCategory): boolean {
        if (oldCategory && currentCategory.parentId !== oldCategory.parentId) {
            return false;
        }

        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!currentCategory.parentId || currentCategory.parentId === '0') {
            categoryList = allTransactionCategories.value[currentCategory.type];
        } else if (allTransactionCategoriesMap.value[currentCategory.parentId]) {
            categoryList = allTransactionCategoriesMap.value[currentCategory.parentId]!.subCategories;
        }

        if (categoryList) {
            for (const [category, index] of itemAndIndex(categoryList)) {
                if (category.id === currentCategory.id) {
                    if (!currentCategory.parentId || currentCategory.parentId === '0') {
                        currentCategory.subCategories = category.subCategories;
                    }

                    categoryList.splice(index, 1, currentCategory);
                    break;
                }
            }
        }

        allTransactionCategoriesMap.value[currentCategory.id] = currentCategory;
        return true;
    }

    function updateCategoryDisplayOrderInCategoryList({ category, from, to }: { category: TransactionCategory, from: number, to: number }): void {
        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!category.parentId || category.parentId === '0') {
            categoryList = allTransactionCategories.value[category.type];
        } else if (allTransactionCategoriesMap.value[category.parentId]) {
            categoryList = allTransactionCategoriesMap.value[category.parentId]!.subCategories;
        }

        if (categoryList) {
            categoryList.splice(to, 0, categoryList.splice(from, 1)[0] as TransactionCategory);
        }
    }

    function updateCategoryVisibilityInTransactionCategoryList({ category, hidden }: { category: TransactionCategory, hidden: boolean }): void {
        if (allTransactionCategoriesMap.value[category.id]) {
            allTransactionCategoriesMap.value[category.id]!.visible = !hidden;
        }
    }

    function removeCategoryFromTransactionCategoryList(currentCategory: TransactionCategory): void {
        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!currentCategory.parentId || currentCategory.parentId === '0') {
            categoryList = allTransactionCategories.value[currentCategory.type];
        } else if (allTransactionCategoriesMap.value[currentCategory.parentId]) {
            categoryList = allTransactionCategoriesMap.value[currentCategory.parentId]!.subCategories;
        }

        if (categoryList) {
            for (const [category, index] of itemAndIndex(categoryList)) {
                if (category.id === currentCategory.id) {
                    categoryList.splice(index, 1);
                    break;
                }
            }
        }

        if (allTransactionCategoriesMap.value[currentCategory.id] && allTransactionCategoriesMap.value[currentCategory.id]!.subCategories) {
            const subCategoryList = allTransactionCategoriesMap.value[currentCategory.id]!.subCategories;

            if (subCategoryList) {
                for (const subCategory of subCategoryList) {
                    if (allTransactionCategoriesMap.value[subCategory.id]) {
                        delete allTransactionCategoriesMap.value[subCategory.id];
                    }
                }
            }
        }

        if (allTransactionCategoriesMap.value[currentCategory.id]) {
            delete allTransactionCategoriesMap.value[currentCategory.id];
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

    function loadAllCategories({ force }: { force?: boolean }): Promise<Record<number, TransactionCategory[]>> {
        if (!force && !transactionCategoryListStateInvalid.value) {
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

                if (force && data.result && isEquals(allTransactionCategories.value, transactionCategories)) {
                    reject({ message: 'Category list is up to date', isUpToDate: true });
                    return;
                }

                loadTransactionCategoryList(transactionCategories);

                resolve(transactionCategories);
            }).catch(error => {
                if (force) {
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

    function getCategory({ categoryId }: { categoryId: string }): Promise<TransactionCategory> {
        return new Promise((resolve, reject) => {
            services.getTransactionCategory({
                id: categoryId
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

    function saveCategory({ category, isEdit, clientSessionId }: { category: TransactionCategory, isEdit: boolean, clientSessionId: string }): Promise<TransactionCategory> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<TransactionCategoryInfoResponse>;

            if (!isEdit) {
                promise = services.addTransactionCategory(category.toCreateRequest(clientSessionId));
            } else {
                promise = services.modifyTransactionCategory(category.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add category' });
                    } else {
                        reject({ message: 'Unable to save category' });
                    }
                    return;
                }

                const transactionCategory = TransactionCategory.of(data.result);

                if (!isEdit) {
                    addCategoryToTransactionCategoryList(transactionCategory);
                } else {
                    const result = updateCategoryInTransactionCategoryList(transactionCategory, allTransactionCategoriesMap.value[category.id]);

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
                    if (!isEdit) {
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
                    reject({ message: 'Unable to add category' });
                    return;
                }

                if (!transactionCategoryListStateInvalid.value) {
                    updateTransactionCategoryListInvalidState(true);
                }

                const transactionCategories = TransactionCategory.ofMap(data.result);

                resolve(transactionCategories);
            }).catch(error => {
                logger.error('failed to add categories', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to add category' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function addPresetCategories(req: TransactionCategoryCreateBatchRequest): Promise<Record<number, TransactionCategory[]>> {
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

    function changeCategoryDisplayOrder({ categoryId, from, to }: { categoryId: string, from: number, to: number }): Promise<void> {
        const category = allTransactionCategoriesMap.value[categoryId];

        return new Promise((resolve, reject) => {
            if (!category) {
                reject({ message: 'Unable to move category' });
                return;
            }

            if (!category.parentId || category.parentId === '0') {
                if (!allTransactionCategories.value[category.type] ||
                    !allTransactionCategories.value[category.type]![to]) {
                    reject({ message: 'Unable to move category' });
                    return;
                }
            } else {
                const subCategoryList = allTransactionCategoriesMap.value[category.parentId]?.subCategories;

                if (!subCategoryList || !subCategoryList[to]) {
                    reject({ message: 'Unable to move category' });
                    return;
                }
            }

            if (!transactionCategoryListStateInvalid.value) {
                updateTransactionCategoryListInvalidState(true);
            }

            updateCategoryDisplayOrderInCategoryList({ category, from, to });

            resolve();
        });
    }

    function updateCategoryDisplayOrders({ type, parentId }: { type: CategoryType, parentId: string }): Promise<boolean> {
        const newDisplayOrders: TransactionCategoryNewDisplayOrderRequest[] = [];

        let categoryList: TransactionCategory[] | undefined = undefined;

        if (!parentId || parentId === '0') {
            categoryList = allTransactionCategories.value[type];
        } else if (allTransactionCategoriesMap.value[parentId]) {
            categoryList = allTransactionCategoriesMap.value[parentId].subCategories;
        }

        if (categoryList) {
            for (const [category, index] of itemAndIndex(categoryList)) {
                newDisplayOrders.push({
                    id: category.id,
                    displayOrder: index + 1
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

                loadAllCategories({ force: false }).finally(() => {
                    if (transactionCategoryListStateInvalid.value) {
                        updateTransactionCategoryListInvalidState(false);
                    }

                    resolve(data.result);
                });
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

    function hideCategory({ category, hidden }: { category: TransactionCategory, hidden: boolean }): Promise<boolean> {
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

                updateCategoryVisibilityInTransactionCategoryList({ category, hidden });

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

    function deleteCategory({ category, beforeResolve }: { category: TransactionCategory, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
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
                        removeCategoryFromTransactionCategoryList(category);
                    });
                } else {
                    removeCategoryFromTransactionCategoryList(category);
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
        // computed states
        allAvailablePrimaryCategoriesCount,
        allAvailableSecondaryCategoriesCount,
        hasVisibleExpenseCategories,
        hasVisibleIncomeCategories,
        hasVisibleTransferCategories,
        // functions
        updateTransactionCategoryListInvalidState,
        resetTransactionCategories,
        loadAllCategories,
        getCategory,
        saveCategory,
        addCategories,
        addPresetCategories,
        changeCategoryDisplayOrder,
        updateCategoryDisplayOrders,
        hideCategory,
        deleteCategory
    };
});
