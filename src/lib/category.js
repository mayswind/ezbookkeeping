import categoryConstants from '@/consts/category.js';
import transactionConstants from '@/consts/transaction.js';

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
