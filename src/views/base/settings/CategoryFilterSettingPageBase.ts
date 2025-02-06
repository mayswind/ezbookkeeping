import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { CategoryType } from '@/core/category.ts';
import type { TransactionCategory, TransactionCategoriesWithVisibleCount } from '@/models/transaction_category.ts';

import {
    copyObjectTo,
    arrayItemToObjectField
} from '@/lib/common.ts';
import {
    allTransactionCategoriesWithVisibleCount,
    containsAnyAvailableCategory,
    containsAvailableCategory,
    selectAllSubCategories,
    isCategoryOrSubCategoriesAllChecked
} from '@/lib/category.ts';

export function useCategoryFilterSettingPageBase(type?: string, allowCategoryTypesStr?: string) {
    const { tt } = useI18n();

    const settingsStore = useSettingsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionsStore = useTransactionsStore();
    const statisticsStore = useStatisticsStore();

    const allowCategoryTypes: Record<string, boolean> | undefined = allowCategoryTypesStr ? arrayItemToObjectField(allowCategoryTypesStr.split(','), true) : undefined;

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const filterCategoryIds = ref<Record<string, boolean>>({});

    const title = computed<string>(() => {
        if (type === 'statisticsDefault') {
            return 'Default Transaction Category Filter';
        } else {
            return 'Filter Transaction Categories';
        }
    });

    const applyText = computed<string>(() => {
        if (type === 'statisticsDefault') {
            return 'Save';
        } else {
            return 'Apply';
        }
    });

    const allTransactionCategories = computed<Record<number, TransactionCategoriesWithVisibleCount>>(() => allTransactionCategoriesWithVisibleCount(transactionCategoriesStore.allTransactionCategories, allowCategoryTypes));
    const hasAnyAvailableCategory = computed<boolean>(() => containsAnyAvailableCategory(allTransactionCategories.value, true));
    const hasAnyVisibleCategory = computed<boolean>(() => containsAnyAvailableCategory(allTransactionCategories.value, showHidden.value));
    const hasAvailableCategory = computed<Record<number, boolean>>(() => containsAvailableCategory(allTransactionCategories.value, showHidden.value));

    function isCategoryChecked(category: TransactionCategory, filterCategoryIds: Record<string, boolean>): boolean {
        return !filterCategoryIds[category.id];
    }

    function getCategoryTypeName(categoryType: CategoryType): string {
        switch (categoryType) {
            case CategoryType.Income:
                return tt('Income Categories');
            case CategoryType.Expense:
                return tt('Expense Categories');
            case CategoryType.Transfer:
                return tt('Transfer Categories');
            default:
                return tt('Transaction Categories');
        }
    }

    function loadFilterCategoryIds(): boolean {
        const allCategoryIds: Record<string, boolean> = {};

        for (const categoryId in transactionCategoriesStore.allTransactionCategoriesMap) {
            if (!Object.prototype.hasOwnProperty.call(transactionCategoriesStore.allTransactionCategoriesMap, categoryId)) {
                continue;
            }

            const category = transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

            if (allowCategoryTypes && !allowCategoryTypes[category.type]) {
                continue;
            }

            if (type === 'transactionListCurrent' && transactionsStore.allFilterCategoryIdsCount > 0) {
                allCategoryIds[category.id] = true;
            } else {
                allCategoryIds[category.id] = false;
            }
        }

        if (type === 'statisticsDefault') {
            filterCategoryIds.value = copyObjectTo(settingsStore.appSettings.statistics.defaultTransactionCategoryFilter, allCategoryIds) as Record<string, boolean>;
            return true;
        } else if (type === 'statisticsCurrent') {
            filterCategoryIds.value = copyObjectTo(statisticsStore.transactionStatisticsFilter.filterCategoryIds, allCategoryIds) as Record<string, boolean>;
            return true;
        } else if (type === 'transactionListCurrent') {
            for (const categoryId in transactionsStore.allFilterCategoryIds) {
                if (!Object.prototype.hasOwnProperty.call(transactionsStore.allFilterCategoryIds, categoryId)) {
                    continue;
                }

                const category = transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                if (category && (!category.secondaryCategories || !category.secondaryCategories.length)) {
                    allCategoryIds[category.id] = false;
                } else if (category) {
                    selectAllSubCategories(allCategoryIds, category, false);
                }
            }

            filterCategoryIds.value = allCategoryIds;
            return true;
        } else {
            return false;
        }
    }

    function saveFilterCategoryIds(): boolean {
        const filteredCategoryIds: Record<string, boolean> = {};
        let isAllSelected = true;
        let finalCategoryIds = '';
        let changed = true;

        for (const categoryId in filterCategoryIds.value) {
            if (!Object.prototype.hasOwnProperty.call(filterCategoryIds.value, categoryId)) {
                continue;
            }

            const category = transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

            if (!isCategoryOrSubCategoriesAllChecked(category, filterCategoryIds.value)) {
                filteredCategoryIds[categoryId] = true;
                isAllSelected = false;
            } else {
                if (finalCategoryIds.length > 0) {
                    finalCategoryIds += ',';
                }

                finalCategoryIds += categoryId;
            }
        }

        if (type === 'statisticsDefault') {
            settingsStore.setStatisticsDefaultTransactionCategoryFilter(filteredCategoryIds);
        } else if (type === 'statisticsCurrent') {
            changed = statisticsStore.updateTransactionStatisticsFilter({
                filterCategoryIds: filteredCategoryIds
            });
        } else if (type === 'transactionListCurrent') {
            changed = transactionsStore.updateTransactionListFilter({
                categoryIds: isAllSelected ? '' : finalCategoryIds
            });

            if (changed) {
                transactionsStore.updateTransactionListInvalidState(true);
            }
        }

        return changed;
    }

    return {
        // states
        loading,
        showHidden,
        filterCategoryIds,
        // computed states
        title,
        applyText,
        allTransactionCategories,
        hasAnyAvailableCategory,
        hasAnyVisibleCategory,
        hasAvailableCategory,
        // functions
        isCategoryChecked,
        getCategoryTypeName,
        loadFilterCategoryIds,
        saveFilterCategoryIds
    };
}
