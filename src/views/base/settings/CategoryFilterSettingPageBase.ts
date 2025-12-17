import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { keys, keysIfValueEquals, values } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

import {
    arrayItemToObjectField
} from '@/lib/common.ts';
import {
    filterTransactionCategories,
    selectAllSubCategories,
    isCategoryOrSubCategoriesAllChecked
} from '@/lib/category.ts';

export type CategoryFilterType = 'statisticsDefault' | 'statisticsCurrent' | 'homePageOverview' | 'transactionListCurrent' | 'custom';

export function useCategoryFilterSettingPageBase(type?: CategoryFilterType, allowCategoryTypesStr?: string, selectedCategoryIds?: string[]) {
    const { tt } = useI18n();

    const settingsStore = useSettingsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionsStore = useTransactionsStore();
    const statisticsStore = useStatisticsStore();
    const overviewStore = useOverviewStore();

    const allowCategoryTypes: Record<string, boolean> | undefined = allowCategoryTypesStr ? arrayItemToObjectField(allowCategoryTypesStr.split(','), true) : undefined;

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const filterContent = ref<string>('');
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

    const allVisibleTransactionCategories = computed<Record<string, TransactionCategory[]>>(() => filterTransactionCategories(transactionCategoriesStore.allTransactionCategories, allowCategoryTypes, filterContent.value, showHidden.value));
    const allVisibleTransactionCategoryMap = computed<Record<string, TransactionCategory>>(() => {
        const categoryMap: Record<string, TransactionCategory> = {};

        for (const categories of values(allVisibleTransactionCategories.value)) {
            for (const category of categories) {
                categoryMap[category.id] = category;

                if (category.subCategories) {
                    for (const subCategory of category.subCategories) {
                        categoryMap[subCategory.id] = subCategory;
                    }
                }
            }
        }

        return categoryMap;
    });
    const hasAnyAvailableCategory = computed<boolean>(() => transactionCategoriesStore.allAvailablePrimaryCategoriesCount > 0 || transactionCategoriesStore.allAvailableSecondaryCategoriesCount > 0);
    const hasAnyVisibleCategory = computed<boolean>(() => {
        for (const categories of values(allVisibleTransactionCategories.value)) {
            if (categories.length > 0) {
                return true;
            }
        }

        return false;
    });

    function isCategoryChecked(category: TransactionCategory, filterCategoryIds: Record<string, boolean>): boolean {
        return !filterCategoryIds[category.id];
    }

    function getCategoryTypeName(categoryType: number): string {
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

        for (const category of values(transactionCategoriesStore.allTransactionCategoriesMap)) {
            if (allowCategoryTypes && !allowCategoryTypes[category.type]) {
                continue;
            }

            if (type === 'transactionListCurrent' && transactionsStore.allFilterCategoryIdsCount > 0) {
                allCategoryIds[category.id] = true;
            } else if (type === 'custom') {
                allCategoryIds[category.id] = true;
            } else {
                allCategoryIds[category.id] = false;
            }
        }

        if (type === 'statisticsDefault') {
            filterCategoryIds.value = Object.assign(allCategoryIds, settingsStore.appSettings.statistics.defaultTransactionCategoryFilter);
            return true;
        } else if (type === 'statisticsCurrent') {
            filterCategoryIds.value = Object.assign(allCategoryIds, statisticsStore.transactionStatisticsFilter.filterCategoryIds);
            return true;
        } else if (type === 'homePageOverview') {
            filterCategoryIds.value = Object.assign(allCategoryIds, settingsStore.appSettings.overviewTransactionCategoryFilterInHomePage);
            return true;
        } else if (type === 'transactionListCurrent') {
            for (const categoryId of keysIfValueEquals(transactionsStore.allFilterCategoryIds, true)) {
                const category = transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                if (category && (!category.subCategories || !category.subCategories.length)) {
                    allCategoryIds[category.id] = false;
                } else if (category) {
                    selectAllSubCategories(allCategoryIds, false, category);
                }
            }

            filterCategoryIds.value = allCategoryIds;
            return true;
        } else if (type === 'custom') {
            if (selectedCategoryIds) {
                for (const categoryId of selectedCategoryIds) {
                    const category = transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

                    if (category && (!category.subCategories || !category.subCategories.length)) {
                        allCategoryIds[category.id] = false;
                    } else if (category) {
                        selectAllSubCategories(allCategoryIds, false, category);
                    }
                }
            }

            filterCategoryIds.value = allCategoryIds;
            return true;
        } else {
            return false;
        }
    }

    function saveFilterCategoryIds(): [boolean, string[]] {
        const selectedCategoryIds: string[] = [];
        const filteredCategoryIds: Record<string, boolean> = {};
        let isAllSelected = true;
        let finalCategoryIds = '';
        let changed = true;

        for (const categoryId of keys(filterCategoryIds.value)) {
            const category = transactionCategoriesStore.allTransactionCategoriesMap[categoryId];

            if (!category) {
                continue;
            }

            if (!isCategoryOrSubCategoriesAllChecked(category, filterCategoryIds.value)) {
                filteredCategoryIds[categoryId] = true;
                isAllSelected = false;
            } else {
                if (finalCategoryIds.length > 0) {
                    finalCategoryIds += ',';
                }

                finalCategoryIds += categoryId;
                selectedCategoryIds.push(categoryId);
            }
        }

        if (type === 'statisticsDefault') {
            settingsStore.setStatisticsDefaultTransactionCategoryFilter(filteredCategoryIds);
        } else if (type === 'statisticsCurrent') {
            changed = statisticsStore.updateTransactionStatisticsFilter({
                filterCategoryIds: filteredCategoryIds
            });
        } else if (type === 'homePageOverview') {
            settingsStore.setOverviewTransactionCategoryFilterInHomePage(filteredCategoryIds);
            overviewStore.updateTransactionOverviewInvalidState(true);
        } else if (type === 'transactionListCurrent') {
            changed = transactionsStore.updateTransactionListFilter({
                categoryIds: isAllSelected ? '' : finalCategoryIds
            });

            if (changed) {
                transactionsStore.updateTransactionListInvalidState(true);
            }
        }

        return [changed, selectedCategoryIds];
    }

    return {
        // states
        loading,
        showHidden,
        filterContent,
        filterCategoryIds,
        // computed states
        title,
        applyText,
        allVisibleTransactionCategories,
        allVisibleTransactionCategoryMap,
        hasAnyAvailableCategory,
        hasAnyVisibleCategory,
        // functions
        isCategoryChecked,
        getCategoryTypeName,
        loadFilterCategoryIds,
        saveFilterCategoryIds
    };
}
