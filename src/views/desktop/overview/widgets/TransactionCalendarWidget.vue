<template>
    <v-card class="overview-widget transaction-calendar-container" :class="{ disabled: loading }">
        <template #title>
            <overview-widget-header :title="title || tt('Transaction Calendar')" :icon="mdiCalendarMonthOutline" />
        </template>

        <v-card-text class="pa-0">
            <transaction-calendar calendar-class="auto-height-calendar justify-content-center"
                                  day-has-transaction-class="font-weight-bold"
                                  week-day-name-type="short"
                                  :readonly="loading || editing" :is-dark-mode="isDarkMode"
                                  :default-currency="defaultCurrency"
                                  :min-date="transactionCalendarMinDate" :max-date="transactionCalendarMaxDate"
                                  :daily-total-amounts="dailyTotalAmounts"
                                  :show-amount="showAmountInHomePage"
                                  :show-income-amount="showIncome"
                                  :show-expense-amount="showExpense"
                                  :show-alternate-date="showAlternateDate"
                                  :model-value="currentCalendarDate"
                                  @update:model-value="selectDate" />
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import OverviewWidgetHeader from './OverviewWidgetHeader.vue';

import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useTransactionCalendarWidgetBase } from '@/views/base/overview/TransactionCalendarWidgetBase.ts';

import { ThemeType } from '@/core/theme.ts';
import type { TextualYearMonthDay } from '@/core/datetime.ts';

import {
    mdiCalendarMonthOutline
} from '@mdi/js';

const props = defineProps<{
    loading: boolean;
    editing?: boolean;
    title?: string;
    transactionTypes: number[];
    showAlternateDate: boolean;
}>();

const router = useRouter();
const theme = useTheme();

const { tt } = useI18n();

const {
    currentCalendarDate,
    showAmountInHomePage,
    showIncome,
    showExpense,
    defaultCurrency,
    dailyTotalAmounts,
    transactionCalendarMinDate,
    transactionCalendarMaxDate,
    getTransactionListUrl
} = useTransactionCalendarWidgetBase(props);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

function selectDate(date: TextualYearMonthDay): void {
    currentCalendarDate.value = date;
    const url = getTransactionListUrl(date);

    if (url) {
        router.push(url);
    }
}
</script>

<style>
.overview-widget {
    &.transaction-calendar-container {
        .v-card-text {
            min-height: 0;

            .dp--main {
                .dp--menu {
                    --dp-background-color: transparent;
                    --dp-menu-border-color: transparent;
                }

                .dp--calendar {
                    .dp--calendar-row {
                        > .dp--calendar-item {
                            .transaction-calendar-daily-amounts {
                                > span.transaction-calendar-alternate-date,
                                > span.transaction-calendar-daily-amount {
                                    font-size: 0.8125rem;
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
</style>
