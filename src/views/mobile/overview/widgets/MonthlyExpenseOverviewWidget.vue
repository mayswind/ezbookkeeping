<template>
    <f7-card class="home-summary-card no-margin-top margin-bottom" :class="{ 'skeleton-text': loading }" :style="cardStyle">
        <f7-card-header class="display-block" :style="cardHeaderStyle">
            <p class="no-margin">
                <span class="card-header-content" v-if="loading">
                    <span class="home-summary-month">Month</span>
                    <span>·</span>
                    <small>Expense</small>
                </span>
                <span class="card-header-content" v-else>
                    <span class="home-summary-month">{{ displayDateRange?.thisMonth?.displayTime }}</span>
                    <span>·</span>
                    <small>{{ tt('Expense') }}</small>
                </span>
            </p>
            <p class="no-margin">
                <span class="month-expense" v-if="loading">0.00 USD</span>
                <span class="month-expense" v-else>{{ transactionOverview?.thisMonth ? getDisplayExpenseAmount(transactionOverview.thisMonth) : '-' }}</span>
                <f7-link class="display-inline-flex margin-inline-start-half" :style="iconStyle"
                         @click="showAmountInHomePage = !showAmountInHomePage">
                    <f7-icon class="ebk-hide-icon" :f7="showAmountInHomePage ? 'eye_slash_fill' : 'eye_fill'"></f7-icon>
                </f7-link>
            </p>
            <p class="no-margin">
                <small class="home-summary-misc" v-if="loading">Monthly income 0.00 USD</small>
                <small class="home-summary-misc" v-else>
                    <span>{{ tt('Monthly income') }}</span>
                    <span>{{ transactionOverview?.thisMonth ? getDisplayIncomeAmount(transactionOverview.thisMonth) : '-' }}</span>
                </small>
            </p>
        </f7-card-header>
    </f7-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { usePeriodStatisticsWidgetBase } from '@/views/base/overview/PeriodStatisticsWidgetBase.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';

import type { ColorValue } from '@/core/color.ts';
import {
    DEFAULT_MOBILE_OVERVIEW_WIDGET_LIGHT_BACKGROUND_COLOR,
    DEFAULT_MOBILE_OVERVIEW_WIDGET_DARK_BACKGROUND_COLOR
} from '@/consts/color.ts';

import { getDisplayColor, getContrastTextColor, getContrastIconColor } from '@/lib/color.ts';

const props = defineProps<{
    loading: boolean;
    height: number;
    lightBackgroundColor?: ColorValue;
    darkBackgroundColor?: ColorValue;
}>();

const environmentsStore = useEnvironmentsStore();

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);
const backgroundColor = computed<ColorValue>(() => isDarkMode.value ?
    props.darkBackgroundColor ?? DEFAULT_MOBILE_OVERVIEW_WIDGET_DARK_BACKGROUND_COLOR :
    props.lightBackgroundColor ?? DEFAULT_MOBILE_OVERVIEW_WIDGET_LIGHT_BACKGROUND_COLOR);
const foregroundColor = computed<ColorValue>(() => getContrastTextColor(backgroundColor.value));
const iconColor = computed<ColorValue>(() => getContrastIconColor(backgroundColor.value));

const cardStyle = computed<Record<string, string>>(() => ({
    'background-color': getDisplayColor(backgroundColor.value),
    'color': getDisplayColor(foregroundColor.value)
}));

const cardHeaderStyle = computed<Record<string, string>>(() => {
    const finalStyle: Record<string, string> = {
        color: getDisplayColor(foregroundColor.value)
    };

    if (props.height === 1) {
        finalStyle['padding-top'] = '10px';
    } else if (props.height === 2) {
        finalStyle['padding-top'] = '60px';
    } else {
        finalStyle['padding-top'] = '120px';
    }

    return finalStyle;
});
const iconStyle = computed<Record<string, string>>(() => ({
    color: getDisplayColor(iconColor.value),
    opacity: '1'
}));

const { tt } = useI18n();
const {
    showAmountInHomePage,
    displayDateRange,
    transactionOverview,
    getDisplayIncomeAmount,
    getDisplayExpenseAmount
} = usePeriodStatisticsWidgetBase({});
</script>
