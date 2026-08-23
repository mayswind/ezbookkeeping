<template>
    <v-card class="h-100" :class="{ disabled: loading }">
        <template #title>
            <span class="text-title-medium">{{ title || tt('Expense Category Ranking') }}</span>
        </template>

        <v-card-text class="py-0">
            <v-list class="py-0" lines="two" v-if="rankingItems.length">
                <template :key="item.id" v-for="(item, idx) in rankingItems">
                    <v-list-item class="px-0 py-3 no-min-height" density="compact">
                        <template #prepend>
                            <router-link class="ranking-list-item" :to="getTransactionItemLinkUrl(item.id)">
                                <ItemIcon size="28px" :icon-type="getCategoryIconType(item.iconType)"
                                          :icon-id="item.icon" :color="item.color" />
                            </router-link>
                        </template>
                        <router-link class="ranking-list-item link-no-color" :to="getTransactionItemLinkUrl(item.id)">
                            <div class="d-flex flex-column ms-2">
                                <div class="d-flex align-center">
                                    <span class="text-truncate">{{ item.name }}</span>
                                    <small class="ranking-percent">{{ formatPercentToLocalizedNumerals(item.percent, 2, '<0.01') }}</small>
                                    <v-spacer />
                                    <span class="ranking-amount ms-2">{{ getDisplayAmount(item.value, rankingData.incomplete) }}</span>
                                </div>
                                <div class="mt-1">
                                    <v-progress-linear :color="item.color ? getCategoryDisplayColor(item.color) : 'primary'"
                                                       :bg-color="isDarkMode ? '#444444' : '#f8f8f8'" :bg-opacity="1"
                                                       :model-value="item.percent" :height="4" />
                                </div>
                            </div>
                        </router-link>
                    </v-list-item>
                    <v-divider v-if="idx !== rankingItems.length - 1" />
                </template>
            </v-list>
            <div v-if="loading && !rankingItems.length">
                <v-skeleton-loader class="skeleton-no-margin py-5 mb-h1" type="text" :key="idx" :loading="true" v-for="idx in props.itemCount"></v-skeleton-loader>
            </div>
            <div class="text-medium-emphasis text-center pt-4" v-if="!loading && !rankingItems.length">{{ tt('No data') }}</div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { CategoryType } from '@/core/category.ts';
import { ThemeType } from '@/core/theme.ts';
import { DISPLAY_HIDDEN_AMOUNT, INCOMPLETE_AMOUNT_SUFFIX } from '@/consts/numeral.ts';

import type { Account } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

import { BIG_DECIMAL_ZERO, parseBigDecimal } from '@/lib/numeral.ts';
import { getCategoryIconType } from '@/lib/icon.ts';
import { getCategoryDisplayColor } from '@/lib/color.ts';

interface RankingItem {
    id: string;
    name: string;
    icon: string;
    iconType: number;
    color: string;
    value: BigDecimal;
    percent: number;
}

const props = defineProps<{
    loading: boolean;
    title?: string;
    dateType: number;
    categoryLevel: string;
    itemCount: number
}>();

const theme = useTheme();

const {
    tt,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatPercentToLocalizedNumerals
} = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const categoriesStore = useTransactionCategoriesStore();
const overviewStore = useOverviewStore();
const exchangeRatesStore = useExchangeRatesStore();

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const rankingItems = computed<RankingItem[]>(() => rankingData.value.items.slice(0, props.itemCount));

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

function getDisplayAmount(amount: BigDecimal, incomplete: boolean): string {
    if (!showAmountInHomePage.value) {
        return formatAmountToLocalizedNumeralsWithCurrency(DISPLAY_HIDDEN_AMOUNT, defaultCurrency.value);
    }

    return formatAmountToLocalizedNumeralsWithCurrency(amount, defaultCurrency.value) + (incomplete ? INCOMPLETE_AMOUNT_SUFFIX : '');
}

function isAccountExcluded(account: Account): boolean {
    const excludedAccounts = settingsStore.appSettings.overviewAccountFilterInHomePage;
    return !!excludedAccounts[account.id] || (account.parentId !== '0' && !!excludedAccounts[account.parentId]);
}

function isCategoryExcluded(category: TransactionCategory): boolean {
    const excludedCategories = settingsStore.appSettings.overviewTransactionCategoryFilterInHomePage;
    return !!excludedCategories[category.id] || (category.parentId !== '0' && !!excludedCategories[category.parentId]);
}

function getTransactionItemLinkUrl(categoryId: string): string {
    return `/transaction/list?${overviewStore.getTransactionListPageParams({ dateType: props.dateType })}&type=3&categoryIds=${categoryId}`;
}
</script>

<style scoped>
.ranking-list-item {
    display: block;
    width: 100%;
    overflow: hidden;
    text-decoration: none;
}

.ranking-percent {
    flex: 0 0 auto;
    font-size: 0.75rem;
    opacity: 0.7;
    margin-inline-start: 6px;
}

.ranking-amount {
    flex: 0 0 auto;
    opacity: 0.8;
}
</style>
