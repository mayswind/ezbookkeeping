<template>
    <f7-card class="monthly-expense-widget no-margin-top margin-bottom" :class="{ 'skeleton-text': loading }">
        <f7-card-content class="padding-horizontal padding-vertical">
            <div class="monthly-expense-widget__header display-flex align-items-baseline">
                <span class="monthly-expense-widget__month font-weight-bold" v-if="loading">Month</span>
                <span class="monthly-expense-widget__month font-weight-bold" v-else>{{ displayDateRange?.thisMonth?.displayTime }}</span>
                <span class="monthly-expense-widget__dot">·</span>
                <span class="monthly-expense-widget__type font-weight-bold">{{ tt('Expense') }}</span>
            </div>

            <div class="monthly-expense-widget__content">
                <div class="monthly-expense-widget__amount-container margin-top-half display-flex align-items-center">
                    <span class="monthly-expense-widget__amount text-color-primary font-weight-bold">
                        <span v-if="!loading">{{ transactionOverview?.thisMonth ? getDisplayExpenseAmount(transactionOverview.thisMonth) : '-' }}</span>
                        <span v-else>0.00 USD</span>
                    </span>
                    <f7-link class="margin-inline-start-half"
                           :aria-label="showAmountInHomePage ? tt('Hide Amount') : tt('Show Amount')"
                           @click="showAmountInHomePage = !showAmountInHomePage">
                        <f7-icon :f7="showAmountInHomePage ? 'eye_slash' : 'eye'" size="20"></f7-icon>
                    </f7-link>
                </div>

                <div class="monthly-expense-widget__income margin-top-half">
                    <span class="monthly-expense-widget__income-title text-color-gray">{{ tt('Monthly income') }}</span>
                    <span class="monthly-expense-widget__income-amount font-weight-bold" v-if="!loading">{{ transactionOverview?.thisMonth ? getDisplayIncomeAmount(transactionOverview.thisMonth) : '-' }}</span>
                    <span class="monthly-expense-widget__income-amount font-weight-bold" v-else>0.00 USD</span>
                </div>
            </div>

            <div class="monthly-expense-widget__illustration img-with-direction" aria-hidden="true">
                <svg viewBox="0 0 168 132" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <circle class="monthly-expense-widget__illustration-halo" cx="112" cy="73" r="52" />
                    <circle class="monthly-expense-widget__illustration-dot" cx="151" cy="25" r="4" />
                    <circle class="monthly-expense-widget__illustration-dot monthly-expense-widget__illustration-dot--muted" cx="18" cy="108" r="3" />

                    <g class="monthly-expense-widget__illustration-paper">
                        <rect x="28" y="16" width="78" height="101" rx="12" />
                        <path d="M85 16H94C100.627 16 106 21.3726 106 28V38L85 16Z" />
                        <rect class="monthly-expense-widget__illustration-line" x="43" y="39" width="34" height="6" rx="3" />
                        <rect class="monthly-expense-widget__illustration-line monthly-expense-widget__illustration-line--short" x="43" y="53" width="49" height="4" rx="2" />
                        <rect class="monthly-expense-widget__illustration-line monthly-expense-widget__illustration-line--short" x="43" y="64" width="39" height="4" rx="2" />
                        <rect class="monthly-expense-widget__illustration-total" x="42" y="84" width="49" height="18" rx="6" />
                        <path class="monthly-expense-widget__illustration-total-line" d="M50 93H68" />
                    </g>

                    <g class="monthly-expense-widget__illustration-chart">
                        <circle cx="117" cy="82" r="29" />
                        <circle class="monthly-expense-widget__illustration-chart-track" cx="117" cy="82" r="18" />
                        <circle class="monthly-expense-widget__illustration-chart-value" cx="117" cy="82" r="18" pathLength="100" />
                        <circle class="monthly-expense-widget__illustration-chart-center" cx="117" cy="82" r="8" />
                    </g>

                    <g class="monthly-expense-widget__illustration-coins">
                        <ellipse cx="139" cy="111" rx="15" ry="5" />
                        <path d="M124 104V111C124 113.761 130.716 116 139 116C147.284 116 154 113.761 154 111V104" />
                        <ellipse cx="139" cy="104" rx="15" ry="5" />
                    </g>
                </svg>
            </div>
        </f7-card-content>
    </f7-card>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { usePeriodStatisticsWidgetBase } from '@/views/base/overview/PeriodStatisticsWidgetBase.ts';

defineProps<{
    loading: boolean;
}>();

const { tt } = useI18n();
const {
    showAmountInHomePage,
    displayDateRange,
    transactionOverview,
    getDisplayIncomeAmount,
    getDisplayExpenseAmount
} = usePeriodStatisticsWidgetBase({});
</script>

<style scoped>
.monthly-expense-widget {
    background:
        radial-gradient(circle at 100% 100%, rgba(var(--f7-theme-color-rgb), 0.075), transparent 11rem),
        linear-gradient(145deg, rgba(var(--f7-theme-color-rgb), 0.035), transparent 56%),
        var(--f7-card-bg-color);
    position: relative;
    overflow: hidden;
}

.monthly-expense-widget::after {
    position: absolute;
    z-index: 0;
    width: 9rem;
    height: 9rem;
    border-radius: 50%;
    background: rgba(var(--f7-theme-color-rgb), 0.055);
    content: '';
    filter: blur(1.5rem);
    bottom: -4.5rem;
    right: -3.5rem;
    pointer-events: none;
}

html[dir="rtl"] .monthly-expense-widget::after {
    right: auto;
    left: -3.5rem;
}

.monthly-expense-widget__header {
    position: relative;
    z-index: 2;
}

.monthly-expense-widget__month {
    font-size: 1.25rem;
    font-weight: bold;
}

.monthly-expense-widget__dot {
    margin: 0 4px;
    font-size: 1.25rem;
    color: var(--f7-text-color);
    opacity: 0.6;
}

.monthly-expense-widget__type {
    font-size: 1rem;
    font-weight: bold;
    color: var(--f7-text-color);
    opacity: 0.8;
}

.monthly-expense-widget__content {
    position: relative;
    z-index: 2;
    display: flex;
    width: calc(100% - 7rem);
    min-width: 10rem;
    flex-direction: column;
    align-items: flex-start;
    padding-bottom: 2rem;
}

.monthly-expense-widget__amount {
    font-size: 1.75rem;
}

.monthly-expense-widget__income {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.1rem;
    margin-top: 1rem;
}

.monthly-expense-widget__income-title {
    font-size: 0.85rem;
    color: var(--f7-text-color);
    opacity: 0.6;
    font-weight: normal;
}

.monthly-expense-widget__income-amount {
    font-size: 1.15rem;
}

.monthly-expense-widget__illustration {
    position: absolute;
    z-index: 1;
    width: clamp(7.25rem, 28%, 9.25rem);
    color: var(--f7-theme-color);
    bottom: 0.3rem;
    right: 1rem;
    opacity: 0.9;
}

html[dir="rtl"] .monthly-expense-widget__illustration {
    right: auto;
    left: 1rem;
}

.monthly-expense-widget__illustration svg {
    display: block;
    width: 100%;
    height: auto;
    overflow: visible;
    filter: drop-shadow(0 0.55rem 0.8rem rgba(0, 0, 0, 0.08));
}

.monthly-expense-widget__illustration-halo {
    fill: currentColor;
    opacity: 0.055;
}

.monthly-expense-widget__illustration-dot {
    fill: currentColor;
    opacity: 0.32;
}

.monthly-expense-widget__illustration-dot--muted {
    opacity: 0.16;
}

.monthly-expense-widget__illustration-paper {
    fill: var(--f7-card-bg-color);
    stroke: currentColor;
    stroke-width: 2;
    stroke-linejoin: round;
}

.monthly-expense-widget__illustration-paper > path {
    fill: currentColor;
    opacity: 0.1;
}

.monthly-expense-widget__illustration-line {
    fill: currentColor;
    stroke: none;
    opacity: 0.34;
}

.monthly-expense-widget__illustration-line--short {
    opacity: 0.15;
}

.monthly-expense-widget__illustration-total {
    fill: currentColor;
    stroke: none;
    opacity: 0.1;
}

.monthly-expense-widget__illustration-total-line {
    fill: none !important;
    stroke: currentColor;
    stroke-linecap: round;
    stroke-width: 3;
    opacity: 0.5 !important;
}

.monthly-expense-widget__illustration-chart {
    fill: var(--f7-card-bg-color);
    stroke: currentColor;
    stroke-width: 2;
}

.monthly-expense-widget__illustration-chart-track,
.monthly-expense-widget__illustration-chart-value {
    fill: none;
    stroke-width: 6;
}

.monthly-expense-widget__illustration-chart-track {
    opacity: 0.12;
}

.monthly-expense-widget__illustration-chart-value {
    stroke-dasharray: 58 42;
    stroke-linecap: round;
    transform: rotate(-90deg);
    transform-origin: 117px 82px;
    opacity: 0.65;
}

.monthly-expense-widget__illustration-chart-center {
    fill: currentColor;
    stroke: none;
    opacity: 0.12;
}

.monthly-expense-widget__illustration-coins {
    fill: var(--f7-card-bg-color);
    stroke: currentColor;
    stroke-width: 2;
}

@media (max-width: 350px) {
    .monthly-expense-widget__illustration {
        display: none !important;
    }
}
</style>
