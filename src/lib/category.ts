import { reversed, keys, values } from '@/core/base.ts';
import { type LocalizedPresetCategory, CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    type TransactionCategoryCreateRequest,
    type TransactionCategoryCreateWithSubCategories,
    TransactionCategory
} from '@/models/transaction_category.ts';

export function transactionTypeToCategoryType(transactionType: TransactionType): CategoryType | null {
    if (transactionType === TransactionType.Income) {
        return CategoryType.Income;
    } else if (transactionType === TransactionType.Expense) {
        return CategoryType.Expense;
    } else if (transactionType === TransactionType.Transfer) {
        return CategoryType.Transfer;
    } else {
        return null;
    }
}

export function categoryTypeToTransactionType(categoryType: CategoryType): TransactionType | null {
    if (categoryType === CategoryType.Income) {
        return TransactionType.Income;
    } else if (categoryType === CategoryType.Expense) {
        return TransactionType.Expense;
    } else if (categoryType === CategoryType.Transfer) {
        return TransactionType.Transfer;
    } else {
        return null;
    }
}

export function localizedPresetCategoryToTransactionCategoryCreateWithSubCategorys(presetCategory: LocalizedPresetCategory): TransactionCategoryCreateWithSubCategories {
    const subCategories: TransactionCategoryCreateRequest[] = [];

    for (const subPresetCategory of presetCategory.subCategories) {
        const subCategory: TransactionCategoryCreateRequest = {
            name: subPresetCategory.name,
            type: subPresetCategory.type,
            parentId: '0',
            icon: subPresetCategory.icon,
            color: subPresetCategory.color,
            comment: '',
            clientSessionId: ''
        };

        subCategories.push(subCategory);
    }

    const categoryWithSubCategories: TransactionCategoryCreateWithSubCategories = {
        name: presetCategory.name,
        type: presetCategory.type,
        icon: presetCategory.icon,
        color: presetCategory.color,
        subCategories: subCategories
    };

    return categoryWithSubCategories;
}

export function localizedPresetCategoriesToTransactionCategoryCreateWithSubCategories(presetCategories: LocalizedPresetCategory[]): TransactionCategoryCreateWithSubCategories[] {
    const categories: TransactionCategoryCreateWithSubCategories[] = [];

    for (const presetCategory of presetCategories) {
        const categoryWithSubCategories = localizedPresetCategoryToTransactionCategoryCreateWithSubCategorys(presetCategory);
        categories.push(categoryWithSubCategories);
    }

    return categories;
}

export function getSecondaryTransactionMapByName(allCategories?: TransactionCategory[]): Record<string, TransactionCategory> {
    const ret: Record<string, TransactionCategory> = {};

    if (!allCategories) {
        return ret;
    }

    for (const category of allCategories) {
        if (category.subCategories) {
            for (const subCategory of category.subCategories) {
                ret[subCategory.name] = subCategory;
            }
        }
    }

    return ret;
}

export function getTransactionPrimaryCategoryName(categoryId: string | null | undefined, allCategories?: TransactionCategory[]): string {
    if (!allCategories) {
        return '';
    }

    for (const category of allCategories) {
        const subCategoryList = category.subCategories;

        if (!subCategoryList) {
            continue;
        }

        for (const subCategory of subCategoryList) {
            if (subCategory.id === categoryId) {
                return category.name;
            }
        }
    }

    return '';
}

export function getTransactionSecondaryCategoryName(categoryId: string | null | undefined, allCategories?: TransactionCategory[]): string {
    if (!allCategories) {
        return '';
    }

    for (const category of allCategories) {
        const subCategoryList = category.subCategories;

        if (!subCategoryList) {
            continue;
        }

        for (const subCategory of subCategoryList) {
            if (subCategory.id === categoryId) {
                return subCategory.name;
            }
        }
    }

    return '';
}

export function filterTransactionCategories(allTransactionCategories: Record<number, TransactionCategory[]>, allowCategoryTypes?: Record<number, boolean>, allowCategoryName?: string, showHidden?: boolean): Record<number, TransactionCategory[]> {
    const ret: Record<string, TransactionCategory[]> = {};
    const hasAllowCategoryTypes = allowCategoryTypes
        && (allowCategoryTypes[CategoryType.Income]
            || allowCategoryTypes[CategoryType.Expense]
            || allowCategoryTypes[CategoryType.Transfer]);

    const allCategoryTypes = [ CategoryType.Income, CategoryType.Expense, CategoryType.Transfer ];
    const lowercaseFilterContent = allowCategoryName ? allowCategoryName.toLowerCase() : '';

    for (const categoryType of allCategoryTypes) {
        const allCategories = allTransactionCategories[categoryType];

        if (!allCategories || allCategories.length < 1) {
            continue;
        }

        if (hasAllowCategoryTypes && !allowCategoryTypes[categoryType]) {
            continue;
        }

        const allFilteredCategories: TransactionCategory[] = [];

        for (const category of allCategories) {
            if (!showHidden && category.hidden) {
                continue;
            }

            const categoryMatchesName = !lowercaseFilterContent || category.name.toLowerCase().includes(lowercaseFilterContent);
            const filteredSubCategories: TransactionCategory[] = [];

            if (category.subCategories) {
                for (const subCategory of category.subCategories) {
                    if (!showHidden && subCategory.hidden) {
                        continue;
                    }

                    if (!categoryMatchesName && lowercaseFilterContent && !subCategory.name.toLowerCase().includes(lowercaseFilterContent)) {
                        continue;
                    }

                    const filteredSubCategory = subCategory.clone();
                    filteredSubCategories.push(filteredSubCategory);
                }
            }

            if (!categoryMatchesName && filteredSubCategories.length < 1) {
                continue;
            }

            const filteredCategory = category.clone();
            filteredCategory.subCategories = filteredSubCategories;
            allFilteredCategories.push(filteredCategory);
        }

        ret[`${categoryType}`] = allFilteredCategories;
    }

    return ret;
}

export function allVisiblePrimaryTransactionCategoriesByType(allTransactionCategories: Record<number, TransactionCategory[]>, categoryType: number): TransactionCategory[] {
    const allCategories = allTransactionCategories[categoryType];
    const visibleCategories: TransactionCategory[] = [];

    if (!allCategories) {
        return visibleCategories;
    }

    for (const category of allCategories) {
        if (category.hidden) {
            continue;
        }

        visibleCategories.push(category);
    }

    return visibleCategories;
}

export function getFinalCategoryIdsByFilteredCategoryIds(allTransactionCategoriesMap: Record<number, TransactionCategory>, filteredCategoryIds: Record<string, boolean>): string {
    let finalCategoryIds = '';

    if (!allTransactionCategoriesMap) {
        return finalCategoryIds;
    }

    for (const category of values(allTransactionCategoriesMap)) {
        if (filteredCategoryIds && !isCategoryOrSubCategoriesAllChecked(category, filteredCategoryIds)) {
            continue;
        }

        if (finalCategoryIds.length > 0) {
            finalCategoryIds += ',';
        }

        finalCategoryIds += category.id;
    }

    return finalCategoryIds;
}

export function isSubCategoryIdAvailable(categories: TransactionCategory[], categoryId: string): boolean {
    if (!categories || !categories.length) {
        return false;
    }

    for (const primaryCategory of categories) {
        if (primaryCategory.hidden) {
            continue;
        }

        const subCategoryList = primaryCategory.subCategories;

        if (!subCategoryList) {
            continue;
        }

        for (const secondaryCategory of subCategoryList) {
            if (secondaryCategory.hidden) {
                continue;
            }

            if (secondaryCategory.id === categoryId) {
                return true;
            }
        }
    }

    return false;
}

export function getFirstVisibleCategoryId(categories?: TransactionCategory[]): string {
    if (!categories || !categories.length) {
        return '';
    }

    for (const primaryCategory of categories) {
        if (primaryCategory.hidden) {
            continue;
        }

        const subCategoryList = primaryCategory.subCategories;

        if (!subCategoryList) {
            continue;
        }

        for (const secondaryCategory of subCategoryList) {
            if (secondaryCategory.hidden) {
                continue;
            }

            return secondaryCategory.id;
        }
    }

    return '';
}

export function getFirstAvailableSubCategoryId(categories: TransactionCategory[], categoryId: string): string {
    if (!categories || !categories.length) {
        return '';
    }

    for (const primaryCategory of categories) {
        if (primaryCategory.hidden || primaryCategory.id !== categoryId) {
            continue;
        }

        const subCategoryList = primaryCategory.subCategories;

        if (!subCategoryList) {
            return '';
        }

        for (const secondaryCategory of subCategoryList) {
            if (secondaryCategory.hidden) {
                continue;
            }

            return secondaryCategory.id;
        }

        return '';
    }

    return '';
}

export function isNoAvailableCategory(categories: TransactionCategory[], showHidden: boolean): boolean {
    for (const category of categories) {
        if (showHidden || !category.hidden) {
            return false;
        }
    }

    return true;
}

export function getAvailableCategoryCount(categories: TransactionCategory[], showHidden: boolean): number {
    let count = 0;

    for (const category of categories) {
        if (showHidden || !category.hidden) {
            count++;
        }
    }

    return count;
}

export function getFirstShowingId(categories: TransactionCategory[], showHidden: boolean): string | null {
    for (const category of categories) {
        if (showHidden || !category.hidden) {
            return category.id;
        }
    }

    return null;
}

export function getLastShowingId(categories: TransactionCategory[], showHidden: boolean): string | null {
    for (const category of reversed(categories)) {
        if (showHidden || !category.hidden) {
            return category.id;
        }
    }

    return null;
}

export function selectAllSubCategories(filterCategoryIds: Record<string, boolean>, value: boolean, category?: TransactionCategory): void {
    if (!category || !category.subCategories || !category.subCategories.length) {
        return;
    }

    for (const subCategory of category.subCategories) {
        filterCategoryIds[subCategory.id] = value;
    }
}

export function selectAllVisible(filterCategoryIds: Record<string, boolean>, allTransactionCategoriesMap: Record<string, TransactionCategory>): void {
    for (const categoryId of keys(filterCategoryIds)) {
        const category = allTransactionCategoriesMap[categoryId];

        if (category) {
            if (category.hidden) {
                continue;
            }

            if (category.parentId && allTransactionCategoriesMap[category.parentId] && allTransactionCategoriesMap[category.parentId]!.hidden) {
                continue;
            }

            filterCategoryIds[category.id] = false;
        }
    }
}

export function selectAll(filterCategoryIds: Record<string, boolean>, allTransactionCategoriesMap: Record<string, TransactionCategory>): void {
    for (const categoryId of keys(filterCategoryIds)) {
        const category = allTransactionCategoriesMap[categoryId];

        if (category) {
            filterCategoryIds[category.id] = false;
        }
    }
}

export function selectNone(filterCategoryIds: Record<string, boolean>, allTransactionCategoriesMap: Record<string, TransactionCategory>): void {
    for (const categoryId of keys(filterCategoryIds)) {
        const category = allTransactionCategoriesMap[categoryId];

        if (category) {
            filterCategoryIds[category.id] = true;
        }
    }
}

export function selectInvert(filterCategoryIds: Record<string, boolean>, allTransactionCategoriesMap: Record<string, TransactionCategory>): void {
    for (const categoryId of keys(filterCategoryIds)) {
        const category = allTransactionCategoriesMap[categoryId];

        if (category) {
            filterCategoryIds[category.id] = !filterCategoryIds[category.id];
        }
    }
}

export function isCategoryOrSubCategoriesAllChecked(category: TransactionCategory, filterCategoryIds: Record<string, boolean>): boolean {
    if (!category.subCategories || category.subCategories.length < 1) {
        return !filterCategoryIds[category.id];
    }

    for (const subCategory of category.subCategories) {
        if (filterCategoryIds[subCategory.id]) {
            return false;
        }
    }

    return true;
}

export function isSubCategoriesAllChecked(category: TransactionCategory, filterCategoryIds: Record<string, boolean>): boolean {
    if (!category.subCategories || category.subCategories.length < 1) {
        return false;
    }

    for (const subCategory of category.subCategories) {
        if (filterCategoryIds[subCategory.id]) {
            return false;
        }
    }

    return true;
}

export function isSubCategoriesHasButNotAllChecked(category: TransactionCategory, filterCategoryIds: Record<string, boolean>): boolean {
    let checkedCount = 0;

    if (!category.subCategories || category.subCategories.length < 1) {
        return false;
    }

    for (const subCategory of category.subCategories) {
        if (!filterCategoryIds[subCategory.id]) {
            checkedCount++;
        }
    }

    return checkedCount > 0 && checkedCount < category.subCategories.length;
}
