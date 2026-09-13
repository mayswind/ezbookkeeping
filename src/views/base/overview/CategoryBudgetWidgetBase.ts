import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { DateRange } from '@/core/datetime.ts';
import { CategoryType } from '@/core/category.ts';

import { DISPLAY_HIDDEN_AMOUNT, INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';

import { BIG_DECIMAL_ZERO } from '@/lib/numeral.ts';
import { getCurrentDateTime } from '@/lib/datetime.ts';

import {
    type CategoryBudgetProgress,
    useCategoryBudgetProgressBase
} from '@/views/base/categories/CategoryBudgetProgressBase.ts';

export interface CommonCategoryBudgetWidgetProps {
    loading: boolean;
    title?: string;
    itemCount: number;
}

export function useCategoryBudgetWidgetBase(props: CommonCategoryBudgetWidgetProps) {
    const { formatAmountToLocalizedNumeralsWithCurrency } = useI18n();
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const overviewStore = useOverviewStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
    const statistics = computed(() => overviewStore.transactionCategoryStatisticsData[DateRange.ThisMonth.type] ?? null);
    const { budgetProgressMap, getBudgetAmountText } = useCategoryBudgetProgressBase(statistics, true);

    const budgetItems = computed<CategoryBudgetProgress[]>(() => Object.values(budgetProgressMap.value)
        .sort((a, b) => {
            const percentComparison = b.percent - a.percent;
            return percentComparison !== 0 ? percentComparison : b.spent.compareTo(a.spent);
        })
        .slice(0, props.itemCount));

    const budgetSummary = computed<{ spent: BigDecimal; limit: BigDecimal; percent: number; incomplete: boolean }>(() => {
        let spent = BIG_DECIMAL_ZERO;
        let limit = BIG_DECIMAL_ZERO;
        let incomplete = false;

        for (const progress of Object.values(budgetProgressMap.value)) {
            const parentHasBudget = progress.category.parentId !== '0' && !!budgetProgressMap.value[progress.category.parentId];

            if (parentHasBudget) {
                continue;
            }

            const convertedSpent = convertToDefaultCurrency(progress.spent, progress.category.budgetCurrency);
            const convertedLimit = convertToDefaultCurrency(progress.limit, progress.category.budgetCurrency);

            if (!convertedSpent || !convertedLimit) {
                incomplete = true;
                continue;
            }

            spent = spent.add(convertedSpent);
            limit = limit.add(convertedLimit);
            incomplete = incomplete || progress.incomplete;
        }

        return {
            spent,
            limit,
            percent: limit.isPositive() ? spent.divide(limit).multiply(100).toDoubleNumber() : 0,
            incomplete
        };
    });

    function convertToDefaultCurrency(amount: BigDecimal, currency: string): BigDecimal | null {
        if (currency === defaultCurrency.value) {
            return amount;
        }

        return exchangeRatesStore.getExchangedAmount(amount, currency, defaultCurrency.value)?.truncate() ?? null;
    }

    function getSummaryAmountText(amount: BigDecimal): string {
        if (!settingsStore.appSettings.showAmountInHomePage) {
            return formatAmountToLocalizedNumeralsWithCurrency(DISPLAY_HIDDEN_AMOUNT, defaultCurrency.value);
        }

        if (!defaultCurrency.value) {
            return '';
        }

        return formatAmountToLocalizedNumeralsWithCurrency(amount, defaultCurrency.value) + (budgetSummary.value.incomplete ? INCOMPLETE_AMOUNT_SUFFIX : '');
    }

    function getItemAmountText(item: CategoryBudgetProgress): string {
        if (!settingsStore.appSettings.showAmountInHomePage) {
            return formatAmountToLocalizedNumeralsWithCurrency(DISPLAY_HIDDEN_AMOUNT, item.category.budgetCurrency);
        }

        return `${getBudgetAmountText(item.spent, item.category.budgetCurrency, item.incomplete)} / ${getBudgetAmountText(item.limit, item.category.budgetCurrency)}`;
    }

    function getCategoryPageLink(item: CategoryBudgetProgress, mobile: boolean): string {
        const now = getCurrentDateTime();
        const month = `${now.getGregorianCalendarYear()}-${now.getGregorianCalendarMonth().toString().padStart(2, '0')}`;
        const primaryCategoryId = item.category.parentId !== '0' ? item.category.parentId : item.category.id;

        if (mobile) {
            return `/category/list?type=${CategoryType.Expense}&id=${primaryCategoryId}&month=${month}`;
        }

        return `/category/list?id=${primaryCategoryId}&month=${month}`;
    }

    return {
        budgetItems,
        budgetSummary,
        getSummaryAmountText,
        getItemAmountText,
        getCategoryPageLink
    };
}
