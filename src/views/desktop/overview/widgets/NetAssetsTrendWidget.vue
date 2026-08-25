<template>
    <v-card class="net-assets-trend-card h-100" :class="{ disabled: loading }">
        <template #title>
            <span class="text-title-medium">{{ title || tt('Net Assets Trends') }}</span>
        </template>

        <trends-chart hide-y-axis-labels hide-horizontal-grid-lines
                      class="mb-2" chart-mode="daily" :type="TrendChartType.Area.type"
                      :start-time="startTime" :end-time="endTime"
                      :start-year-month="undefined" :end-year-month="undefined"
                      :sorting-type="ChartSortingType.Amount.type"
                      :data-aggregation-type="ChartDataAggregationType.Last"
                      :date-aggregation-type="ChartDateAggregationType.Month.type"
                      :fiscal-year-start="fiscalYearStart" :items="items"
                      :value-type="ChartValueType.Amount" :show-value="showAmountInHomePage"
                      :default-currency="defaultCurrency"
                      :hide-legend="!showLegend" legend-position="bottom"
                      :hide-x-axis-labels="!showXAxisLabels" />
    </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { entries } from '@/core/base.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import { ChartValueType } from '@/core/chart.ts';
import { ChartSortingType, ChartDataAggregationType, ChartDateAggregationType, TrendChartType } from '@/core/statistics.ts';
import type { TransactionAssetTrendsAnalysisDataItem, TransactionAssetTrendsAnalysisDataAmount } from '@/models/transaction.ts';

import { BIG_DECIMAL_ZERO, parseBigDecimal } from '@/lib/numeral.ts';
import { getUnixTimeBeforeUnixTime } from '@/lib/datetime.ts';

const props = defineProps<{
    loading: boolean;
    title?: string;
    months: number;
    showLegend: boolean;
    showXAxisLabels: boolean;
}>();

const { tt } = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const overviewStore = useOverviewStore();
const exchangeRatesStore = useExchangeRatesStore();

const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const startTime = computed<number>(() => getUnixTimeBeforeUnixTime(overviewStore.transactionDataRange.thisMonth.startTime, props.months - 1, 'months'));
const endTime = computed<number>(() => overviewStore.transactionDataRange.thisMonth.endTime);
const fiscalYearStart = computed<number>(() => userStore.currentUserFiscalYearStart);

const items = computed<TransactionAssetTrendsAnalysisDataItem[]>(() => {
    const lastBalances: Record<string, string> = {};
    const amounts: TransactionAssetTrendsAnalysisDataAmount[] = [];

    for (const dailyItem of overviewStore.transactionAssetTrendsData) {
        for (const accountItem of dailyItem.items) {
            lastBalances[accountItem.accountId] = accountItem.accountClosingBalance;
        }

        let total: BigDecimal = BIG_DECIMAL_ZERO;

        for (const [accountId, balanceText] of entries(lastBalances)) {
            const account = accountsStore.allAccountsMap[accountId];

            if (!account || settingsStore.appSettings.overviewAccountFilterInHomePage[account.id]) {
                continue;
            }

            let balance: BigDecimal = parseBigDecimal(balanceText);

            if (account.currency !== defaultCurrency.value) {
                const exchangedBalance = exchangeRatesStore.getExchangedAmount(balance, account.currency, defaultCurrency.value);

                if (!exchangedBalance) {
                    continue;
                }

                balance = exchangedBalance.truncate();
            }

            if (account.isAsset) {
                total = total.add(balance);
            } else if (account.isLiability) {
                total = total.add(balance);
            }
        }

        const dailyAmount: TransactionAssetTrendsAnalysisDataAmount = {
            year: dailyItem.year,
            month: dailyItem.month,
            day: dailyItem.day,
            value: total
        };

        amounts.push(dailyAmount);
    }

    return [{
        name: tt('Net assets'),
        type: 'total',
        id: 'total',
        icon: '',
        iconType: 0,
        color: '',
        hidden: false,
        displayOrders: [1],
        value: amounts.at(-1)?.value ?? BIG_DECIMAL_ZERO,
        items: amounts
    }];
});
</script>

<style scoped>
.net-assets-trend-card {
    display: flex;
    flex-direction: column;
}

.trends-chart-container {
    flex: 1 1 0;
    min-height: 0;
    height: auto;
}
</style>
