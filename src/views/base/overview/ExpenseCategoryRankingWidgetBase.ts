import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import { DISPLAY_HIDDEN_AMOUNT, INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';

import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

import { BIG_DECIMAL_ZERO, parseBigDecimal } from '@/lib/numeral.ts';

export interface RankingItem {
    id: string;
    name: string;
    icon: string;
    iconType: number;
    color: string;
    value: BigDecimal;
    percent: number;
}

export interface CommonExpenseCategoryRankingWidgetProps {
    loading: boolean;
    title?: string;
    dateType: number;
    categoryLevel: string;
    itemCount: number;
}

export function useExpenseCategoryRankingWidgetBase(props: CommonExpenseCategoryRankingWidgetProps) {
    const { formatAmountToLocalizedNumeralsWithCurrency } = useI18n();

    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const categoriesStore = useTransactionCategoriesStore();
    const overviewStore = useOverviewStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
    const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);


    const rankingData = computed<{ total: BigDecimal; incomplete: boolean; items: RankingItem[] }>(() => {
        const response = overviewStore.transactionCategoryStatisticsData[props.dateType];
        const values: Record<string, RankingItem> = {};
        let total: BigDecimal = BIG_DECIMAL_ZERO;
        let hasIncompleteAmount: boolean = false;

        for (const responseItem of response?.items ?? []) {
            const account = accountsStore.allAccountsMap[responseItem.accountId];
            const category = categoriesStore.allTransactionCategoriesMap[responseItem.categoryId];

            if (!account || !category || category.type !== CategoryType.Expense || isAccountExcluded(account) || isCategoryExcluded(category)) {
                continue;
            }

            let finalCategory: TransactionCategory = category;

            if (props.categoryLevel === 'primary' && category.parentId !== '0') {
                finalCategory = categoriesStore.allTransactionCategoriesMap[category.parentId] ?? category;
            }

            if (isCategoryExcluded(finalCategory)) {
                continue;
            }

            let amount = parseBigDecimal(responseItem.amount);

            if (account.currency !== defaultCurrency.value) {
                const exchangedAmount = exchangeRatesStore.getExchangedAmount(amount, account.currency, defaultCurrency.value);

                if (!exchangedAmount) {
                    hasIncompleteAmount = true;
                    continue;
                }

                amount = exchangedAmount.truncate();
            }

            const existing = values[finalCategory.id];

            if (existing) {
                existing.value = existing.value.add(amount);
            } else {
                values[finalCategory.id] = {
                    id: finalCategory.id,
                    name: finalCategory.name,
                    icon: finalCategory.icon,
                    iconType: finalCategory.iconType,
                    color: finalCategory.color,
                    value: amount,
                    percent: 0
                };
            }

            total = total.add(amount);
        }

        const items = Object.values(values).sort((a, b) => b.value.compareTo(a.value));

        for (const item of items) {
            item.percent = total.isPositive() ? item.value.divide(total).multiply(100).toDoubleNumber() : 0;
        }

        return {
            total: total,
            incomplete: hasIncompleteAmount,
            items: items
        };
    });

    const rankingItems = computed<RankingItem[]>(() => rankingData.value.items.slice(0, props.itemCount));

    function isAccountExcluded(account: Account): boolean {
        const excludedAccounts = settingsStore.appSettings.overviewAccountFilterInHomePage;
        return !!excludedAccounts[account.id] || (account.parentId !== '0' && !!excludedAccounts[account.parentId]);
    }

    function isCategoryExcluded(category: TransactionCategory): boolean {
        const excludedCategories = settingsStore.appSettings.overviewTransactionCategoryFilterInHomePage;
        return !!excludedCategories[category.id] || (category.parentId !== '0' && !!excludedCategories[category.parentId]);
    }

    function getDisplayAmount(amount: BigDecimal, incomplete: boolean): string {
        if (!showAmountInHomePage.value) {
            return formatAmountToLocalizedNumeralsWithCurrency(DISPLAY_HIDDEN_AMOUNT, defaultCurrency.value);
        }

        return formatAmountToLocalizedNumeralsWithCurrency(amount, defaultCurrency.value) + (incomplete ? INCOMPLETE_AMOUNT_SUFFIX : '');
    }

    function getTransactionItemLinkUrl(categoryId: string): string {
        return `/transaction/list?${overviewStore.getTransactionListPageParams({ type: TransactionType.Expense, dateType: props.dateType, categoryIds: [categoryId] })}`;
    }

    return {
        // computed states
        rankingData,
        rankingItems,
        // functions
        getDisplayAmount,
        getTransactionItemLinkUrl
    };
}
