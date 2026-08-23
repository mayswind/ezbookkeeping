<template>
    <v-card class="transaction-calendar-heatmap-card h-100" :class="{ disabled: loading }">
        <template #title>
            <span class="text-title-medium">{{ title || tt('Transaction Calendar Heatmap') }}</span>
        </template>

        <v-card-text class="transaction-calendar-heatmap-body pa-0 overflow-hidden">
            <date-range-calendar-heat-map-chart :start-time="startTime" :end-time="endTime"
                                                :items="items" :value-type="ChartValueType.Amount"
                                                :value-type-name="tt(transactionTypeName)" :default-currency="defaultCurrency"
                                                :show-value="showAmountInHomePage" :enable-click-item="!editing"
                                                @click="clickDate" />
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { BigDecimal } from '@/core/numeral.ts';
import { DateRange, KnownDateTimeFormat } from '@/core/datetime.ts';
import { ChartValueType, type CalendarChartSourceDataItem } from '@/core/chart.ts';
import { TransactionType } from '@/core/transaction.ts';

import { BIG_DECIMAL_ZERO, parseBigDecimal } from '@/lib/numeral.ts';
import {
    parseDateTimeFromKnownDateTimeFormat,
    getUnixTimeBeforeUnixTime,
    getUnixTimeAfterUnixTime
} from '@/lib/datetime.ts';

const props = defineProps<{
    loading: boolean;
    months: number;
    transactionType: TransactionType;
    editing?: boolean;
    title?: string
}>();

const router = useRouter();

const { tt } = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const overviewStore = useOverviewStore();
const exchangeRatesStore = useExchangeRatesStore();

const showAmountInHomePage = computed<boolean>(() => settingsStore.appSettings.showAmountInHomePage);
const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);
const startTime = computed<number>(() => getUnixTimeAfterUnixTime(getUnixTimeBeforeUnixTime(overviewStore.transactionDataRange.today.startTime, props.months, 'months'), 1, 'days'));
const endTime = computed<number>(() => overviewStore.transactionDataRange.today.endTime);
const transactionTypeName = computed<string>(() => props.transactionType === TransactionType.Income ? 'Income' : 'Expense');

const items = computed<CalendarChartSourceDataItem[]>(() => {
    const result: CalendarChartSourceDataItem[] = [];

    for (const dailyItem of overviewStore.transactionDailyAmountsData) {
        const dateTime = parseDateTimeFromKnownDateTimeFormat(dailyItem.date, KnownDateTimeFormat.DefaultDate);

        if (!dateTime || dateTime.getUnixTime() < startTime.value) {
            continue;
        }

        let total: BigDecimal = BIG_DECIMAL_ZERO;

        for (const amountItem of dailyItem.amounts) {
            let amount: BigDecimal = parseBigDecimal(props.transactionType === TransactionType.Income ? amountItem.incomeAmount : amountItem.expenseAmount);

            if (amountItem.currency !== defaultCurrency.value) {
                const exchangedAmount = exchangeRatesStore.getExchangedAmount(amount, amountItem.currency, defaultCurrency.value);

                if (!exchangedAmount) {
                    continue;
                }

                amount = exchangedAmount.truncate();
            }

            total = total.add(amount);
        }

        if (!total.isZero()) {
            result.push({
                id: dailyItem.date,
                value: total
            });
        }
    }

    return result;
});

function clickDate(date: string): void {
    const dateTime = parseDateTimeFromKnownDateTimeFormat(date, KnownDateTimeFormat.DefaultDate);

    if (!dateTime) {
        return;
    }

    const minTime = dateTime.getUnixTime();
    const maxTime = getUnixTimeBeforeUnixTime(getUnixTimeAfterUnixTime(minTime, 1, 'days'), 1, 'seconds');
    router.push(`/transaction/list?${overviewStore.getTransactionListPageParams({ type: props.transactionType, dateType: DateRange.Custom.type, minTime, maxTime })}`);
}
</script>

<style scoped>
.transaction-calendar-heatmap-card {
    display: flex;
    flex-direction: column;
}

.transaction-calendar-heatmap-body {
    flex: 1 1 0;
    min-height: 0;
}
</style>
