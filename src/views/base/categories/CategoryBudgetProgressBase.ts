import { computed, type Ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { CategoryType } from '@/core/category.ts';

import type { TransactionStatisticResponse } from '@/models/transaction.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

import { BIG_DECIMAL_ZERO, parseBigDecimal } from '@/lib/numeral.ts';

export interface CategoryBudgetProgress {
    readonly category: TransactionCategory;
    spent: BigDecimal;
    readonly limit: BigDecimal;
    percent: number;
    overrun: BigDecimal;
    incomplete: boolean;
}

export function useCategoryBudgetProgressBase(statistics: Ref<TransactionStatisticResponse | null>, applyOverviewFilters = false) {
    const { formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

    const settingsStore = useSettingsStore();
    const accountsStore = useAccountsStore();
    const categoriesStore = useTransactionCategoriesStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const budgetProgressMap = computed<Record<string, CategoryBudgetProgress>>(() => {
        const result: Record<string, CategoryBudgetProgress> = {};

        for (const category of Object.values(categoriesStore.allTransactionCategoriesMap)) {
            if (category.type !== CategoryType.Expense || category.budgetAmount <= 0 || !category.budgetCurrency || category.hidden) {
                continue;
            }

            result[category.id] = {
                category,
                spent: BIG_DECIMAL_ZERO,
                limit: parseBigDecimal(category.budgetAmount),
                percent: 0,
                overrun: BIG_DECIMAL_ZERO,
                incomplete: false
            };
        }

        for (const item of statistics.value?.items ?? []) {
            const account = accountsStore.allAccountsMap[item.accountId];
            const category = categoriesStore.allTransactionCategoriesMap[item.categoryId];

            if (!account || !category || category.type !== CategoryType.Expense || isAccountExcluded(account.id, account.parentId) || isCategoryExcluded(category)) {
                continue;
            }

            addAmount(result[category.id], item.amount, account.currency);

            if (category.parentId && category.parentId !== '0') {
                const parentCategory = categoriesStore.allTransactionCategoriesMap[category.parentId];

                if (parentCategory && !isCategoryExcluded(parentCategory)) {
                    addAmount(result[parentCategory.id], item.amount, account.currency);
                }
            }
        }

        for (const progress of Object.values(result)) {
            progress.percent = progress.limit.isPositive()
                ? progress.spent.divide(progress.limit).multiply(100).toDoubleNumber()
                : 0;
            progress.overrun = progress.spent.compareTo(progress.limit) > 0
                ? progress.spent.subtract(progress.limit)
                : BIG_DECIMAL_ZERO;
        }

        return result;
    });

    function addAmount(progress: CategoryBudgetProgress | undefined, textualAmount: string, sourceCurrency: string): void {
        if (!progress) {
            return;
        }

        let amount = parseBigDecimal(textualAmount);

        if (sourceCurrency !== progress.category.budgetCurrency) {
            const exchangedAmount = exchangeRatesStore.getExchangedAmount(amount, sourceCurrency, progress.category.budgetCurrency);

            if (!exchangedAmount) {
                progress.incomplete = true;
                return;
            }

            amount = exchangedAmount.truncate();
        }

        progress.spent = progress.spent.add(amount);
    }

    function isAccountExcluded(accountId: string, parentAccountId: string): boolean {
        if (!applyOverviewFilters) {
            return false;
        }

        const excludedAccounts = settingsStore.appSettings.overviewAccountFilterInHomePage;
        return !!excludedAccounts[accountId] || (parentAccountId !== '0' && !!excludedAccounts[parentAccountId]);
    }

    function isCategoryExcluded(category: TransactionCategory): boolean {
        if (!applyOverviewFilters) {
            return false;
        }

        const excludedCategories = settingsStore.appSettings.overviewTransactionCategoryFilterInHomePage;
        return !!excludedCategories[category.id] || (category.parentId !== '0' && !!excludedCategories[category.parentId]);
    }

    function getBudgetProgress(categoryId: string): CategoryBudgetProgress | undefined {
        return budgetProgressMap.value[categoryId];
    }

    function getBudgetAmountText(amount: BigDecimal, currency: string, incomplete = false): string {
        return formatAmountToLocalizedNumeralsWithCurrency(amount, currency) + (incomplete ? '*' : '');
    }

    return {
        budgetProgressMap,
        getBudgetProgress,
        getBudgetAmountText
    };
}
