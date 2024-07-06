import categoryConstants from '@/consts/category.js';
import transactionConstants from '@/consts/transaction.js';

export function setCategoryModelByAnotherCategory(category, category2) {
    category.id = category2.id;
    category.type = category2.type;
    category.parentId = category2.parentId;
    category.name = category2.name;
    category.icon = category2.icon;
    category.color = category2.color;
    category.comment = category2.comment;
    category.visible = !category2.hidden;
}

export function transactionTypeToCategoryType(transactionType) {
    if (transactionType === transactionConstants.allTransactionTypes.Income) {
        return categoryConstants.allCategoryTypes.Income;
    } else if (transactionType === transactionConstants.allTransactionTypes.Expense) {
        return categoryConstants.allCategoryTypes.Expense;
    } else if (transactionType === transactionConstants.allTransactionTypes.Transfer) {
        return categoryConstants.allCategoryTypes.Transfer;
    } else {
        return null;
    }
}

export function categoryTypeToTransactionType(categoryType) {
    if (categoryType === categoryConstants.allCategoryTypes.Income) {
        return transactionConstants.allTransactionTypes.Income;
    } else if (categoryType === categoryConstants.allCategoryTypes.Expense) {
        return transactionConstants.allTransactionTypes.Expense;
    } else if (categoryType === categoryConstants.allCategoryTypes.Transfer) {
        return transactionConstants.allTransactionTypes.Transfer;
    } else {
        return null;
    }
}

export function getTransactionPrimaryCategoryName(categoryId, allCategories) {
    for (let i = 0; i < allCategories.length; i++) {
        for (let j = 0; j < allCategories[i].subCategories.length; j++) {
            const subCategory = allCategories[i].subCategories[j];
            if (subCategory.id === categoryId) {
                return allCategories[i].name;
            }
        }
    }

    return '';
}

export function getTransactionSecondaryCategoryName(categoryId, allCategories) {
    for (let i = 0; i < allCategories.length; i++) {
        for (let j = 0; j < allCategories[i].subCategories.length; j++) {
            const subCategory = allCategories[i].subCategories[j];
            if (subCategory.id === categoryId) {
                return subCategory.name;
            }
        }
    }

    return '';
}

export function allVisibleTransactionCategories(allTransactionCategories) {
    const ret = {};

    for (let key in categoryConstants.allCategoryTypes) {
        if (!Object.prototype.hasOwnProperty.call(categoryConstants.allCategoryTypes, key)) {
            continue;
        }

        const categoryType = categoryConstants.allCategoryTypes[key];

        if (!allTransactionCategories[categoryType]) {
            continue;
        }

        const allCategories = allTransactionCategories[categoryType];
        const visibleCategories = [];
        const allVisibleSubCategories = {};

        for (let j = 0; j < allCategories.length; j++) {
            const category = allCategories[j];

            if (category.hidden) {
                continue;
            }

            visibleCategories.push(category);

            if (category.subCategories) {
                const visibleSubCategories = [];

                for (let k = 0; k < category.subCategories.length; k++) {
                    const subCategory = category.subCategories[k];

                    if (!subCategory.hidden) {
                        visibleSubCategories.push(subCategory);
                    }
                }

                if (visibleSubCategories.length > 0) {
                    allVisibleSubCategories[category.id] = visibleSubCategories;
                }
            }
        }

        ret[categoryType.toString()] = {
            type: categoryType.toString(),
            visibleCategories: visibleCategories,
            visibleSubCategories: allVisibleSubCategories
        };
    }

    return ret;
}

export function allVisiblePrimaryTransactionCategoriesByType(allTransactionCategories, type) {
    const allVisibleCategories = allVisibleTransactionCategories(allTransactionCategories);

    if (!allVisibleCategories) {
        return [];
    }

    if (!allVisibleCategories[type.toString()]) {
        return [];
    }

    return allVisibleCategories[type.toString()].visibleCategories;
}

export function isSubCategoryIdAvailable(categories, categoryId) {
    if (!categories || !categories.length) {
        return false;
    }

    for (let i = 0; i < categories.length; i++) {
        for (let j = 0; j < categories[i].subCategories.length; j++) {
            if (categories[i].subCategories[j].id === categoryId) {
                return true;
            }
        }
    }

    return false;
}

export function getFirstAvailableCategoryId(categories) {
    if (!categories || !categories.length) {
        return '';
    }

    for (let i = 0; i < categories.length; i++) {
        for (let j = 0; j < categories[i].subCategories.length; j++) {
            return categories[i].subCategories[j].id;
        }
    }
}

export function getFirstAvaiableSubCategoryId(categories, categoryId) {
    if (!categories || !categories.length) {
        return '';
    }

    for (let i = 0; i < categories.length; i++) {
        if (categories[i].id !== categoryId) {
            continue;
        }

        if (categories[i].subCategories.length <= 0) {
            return '';
        }

        return categories[i].subCategories[0].id;
    }

    return '';
}

export function hasAnyAvailableCategory(allVisibleTransactionCategories) {
    for (let type in allVisibleTransactionCategories) {
        if (!Object.prototype.hasOwnProperty.call(allVisibleTransactionCategories, type)) {
            continue;
        }

        const categoryType = allVisibleTransactionCategories[type];

        if (categoryType.visibleCategories && categoryType.visibleCategories.length > 0) {
            return true;
        }
    }

    return false;
}

export function hasAvailableCategory(allVisibleTransactionCategories) {
    const result = {};

    for (let type in allVisibleTransactionCategories) {
        if (!Object.prototype.hasOwnProperty.call(allVisibleTransactionCategories, type)) {
            continue;
        }

        const categoryType = allVisibleTransactionCategories[type];
        result[type] = categoryType.visibleCategories && categoryType.visibleCategories.length > 0;
    }

    return result;
}

export function selectSubCategories(filterCategoryIds, category, value) {
    if (!category || !category.subCategories || !category.subCategories.length) {
        return;
    }

    for (let i = 0; i < category.subCategories.length; i++) {
        const subCategory = category.subCategories[i];
        filterCategoryIds[subCategory.id] = value;
    }
}

export function selectAll(filterCategoryIds, allTransactionCategoriesMap) {
    for (let categoryId in filterCategoryIds) {
        if (!Object.prototype.hasOwnProperty.call(filterCategoryIds, categoryId)) {
            continue;
        }

        const category = allTransactionCategoriesMap[categoryId];

        if (category) {
            filterCategoryIds[category.id] = false;
        }
    }
}

export function selectNone(filterCategoryIds, allTransactionCategoriesMap) {
    for (let categoryId in filterCategoryIds) {
        if (!Object.prototype.hasOwnProperty.call(filterCategoryIds, categoryId)) {
            continue;
        }

        const category = allTransactionCategoriesMap[categoryId];

        if (category) {
            filterCategoryIds[category.id] = true;
        }
    }
}

export function selectInvert(filterCategoryIds, allTransactionCategoriesMap) {
    for (let categoryId in filterCategoryIds) {
        if (!Object.prototype.hasOwnProperty.call(filterCategoryIds, categoryId)) {
            continue;
        }

        const category = allTransactionCategoriesMap[categoryId];

        if (category) {
            filterCategoryIds[category.id] = !filterCategoryIds[category.id];
        }
    }
}

export function isCategoryOrSubCategoriesAllChecked(category, filterCategoryIds) {
    if (!category.subCategories) {
        return !filterCategoryIds[category.id];
    }

    for (let i = 0; i < category.subCategories.length; i++) {
        const subCategory = category.subCategories[i];
        if (filterCategoryIds[subCategory.id]) {
            return false;
        }
    }

    return true;
}

export function isSubCategoriesAllChecked(category, filterCategoryIds) {
    for (let i = 0; i < category.subCategories.length; i++) {
        const subCategory = category.subCategories[i];
        if (filterCategoryIds[subCategory.id]) {
            return false;
        }
    }

    return true;
}

export function isSubCategoriesHasButNotAllChecked(category, filterCategoryIds) {
    let checkedCount = 0;

    for (let i = 0; i < category.subCategories.length; i++) {
        const subCategory = category.subCategories[i];
        if (!filterCategoryIds[subCategory.id]) {
            checkedCount++;
        }
    }

    return checkedCount > 0 && checkedCount < category.subCategories.length;
}
