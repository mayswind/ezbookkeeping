<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar :class="{ 'disabled': loading }" :title="tt('Budgets')">
            <f7-nav-right>
                <f7-link :class="{ 'disabled': loading }" icon-f7="plus" @click="goToEdit"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <!-- Month picker header -->
        <f7-list strong inset dividers class="budget-month-picker">
            <f7-list-item>
                <div class="month-picker-row">
                    <f7-button icon-f7="chevron_left" @click="changeMonth(-1)"></f7-button>
                    <span class="month-label" @click="resetToCurrentMonth">{{ formatMonth(currentMonthStartDate) }}</span>
                    <f7-button icon-f7="chevron_right" @click="changeMonth(1)"></f7-button>
                </div>
            </f7-list-item>
        </f7-list>

        <!-- Summary card -->
        <f7-block v-if="overview && overview.items.length > 0" class="budget-summary">
            <div class="summary-row">
                <span class="summary-label">{{ tt('Total Budget') }}</span>
                <span class="summary-value">{{ displayAmount(overview.totalLimit, defaultCurrency) }}</span>
            </div>
            <div class="summary-row">
                <span class="summary-label">{{ tt('Total Spent') }}</span>
                <span class="summary-value spent">{{ displayAmount(overview.totalActualExpenseAmount, defaultCurrency) }}</span>
            </div>
            <div class="summary-row total">
                <span class="summary-label">{{ tt('Total Available') }}</span>
                <span class="summary-value" :class="{ 'overspent': overview.totalAvailableAmount < 0 }">
                    {{ displayAmount(overview.totalAvailableAmount, defaultCurrency) }}
                </span>
            </div>
        </f7-block>

        <!-- Loading skeleton -->
        <f7-list strong inset dividers class="margin-top-half skeleton-text" v-if="loading">
            <f7-list-item title="Category Name"></f7-list-item>
            <f7-list-item title="Category Name"></f7-list-item>
            <f7-list-item title="Category Name"></f7-list-item>
        </f7-list>

        <!-- Empty state -->
        <f7-block v-else-if="!overview || overview.items.length === 0" class="text-align-center">
            <f7-icon f7="money_pound_circle" size="48" style="color: var(--f7-list-item-footer-text-color)"></f7-icon>
            <p>{{ tt('No budgets for this month. Tap + to set one.') }}</p>
        </f7-block>

        <!-- Budget list -->
        <f7-list strong inset dividers v-else>
            <f7-list-item
                v-for="item in overview.items"
                :key="item.id"
                :title="item.categoryName || tt('(Unknown Category)')"
                @click="editBudget(item)">
                <div class="budget-item-content">
                    <f7-progressbar
                        :progress="progressPercent(item)"
                        :class="{ 'progress-overspent': item.availableAmount < 0 }">
                    </f7-progressbar>
                    <div class="budget-amounts">
                        <span class="spent">{{ displayAmount(item.actualExpenseAmount, item.currency) }}</span>
                        <span class="separator">/</span>
                        <span class="limit">{{ displayAmount(item.amount, item.currency) }}</span>
                    </div>
                </div>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useUserStore } from '@/stores/user.ts';
import { useCategoryBudgetStore } from '@/stores/categoryBudget.ts';
import { parseBigDecimal } from '@/lib/numeral.ts';

import type { CategoryBudgetOverviewItem } from '@/models/category_budget_limit.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, formatAmountToLocalizedNumeralsWithCurrency } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const userStore = useUserStore();
const budgetStore = useCategoryBudgetStore();

const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);

// First-of-month unix timestamp for the currently selected month (client timezone).
function firstOfMonth(date: Date): number {
    return new Date(date.getFullYear(), date.getMonth(), 1).getTime() / 1000;
}
function addMonthsToDate(unixSeconds: number, months: number): number {
    const d = new Date(unixSeconds * 1000);
    return new Date(d.getFullYear(), d.getMonth() + months, 1).getTime() / 1000;
}

const currentMonthStartDate = ref<number>(firstOfMonth(new Date()));

const overview = computed(() => budgetStore.currentOverview);
const defaultCurrency = computed(() => userStore.currentUserDefaultCurrency || 'USD');

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
    loadCurrentOverview();
}

function loadCurrentOverview(): void {
    loading.value = true;
    budgetStore.loadOverview({ startDate: currentMonthStartDate.value }).catch(error => {
        loadingError.value = error;

        if (!error.processed) {
            showToast(error.message || error);
        }
    }).finally(() => {
        loading.value = false;
    });
}

function reload(done?: () => void): void {
    budgetStore.loadOverview({ startDate: currentMonthStartDate.value }).then(() => {
        done?.();
        showToast('Budget overview has been updated');
    }).catch(error => {
        done?.();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function changeMonth(delta: number): void {
    currentMonthStartDate.value = addMonthsToDate(currentMonthStartDate.value, delta);
    loadCurrentOverview();
}

function resetToCurrentMonth(): void {
    currentMonthStartDate.value = firstOfMonth(new Date());
    loadCurrentOverview();
}

function progressPercent(item: CategoryBudgetOverviewItem): number {
    if (item.amount <= 0) {
        return 0;
    }
    const pct = (item.actualExpenseAmount / item.amount) * 100;
    return Math.min(100, Math.max(0, pct));
}

function formatMonth(unixSeconds: number): string {
    const d = new Date(unixSeconds * 1000);
    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
    return `${tt(`datetime.${monthNames[d.getMonth()]}.long`)} ${d.getFullYear()}`;
}

function editBudget(item: CategoryBudgetOverviewItem): void {
    props.f7router.navigate('/budget/edit', {
        props: {
            editBudgetId: item.id,
            editCategoryId: item.categoryId,
            editAmount: item.amount,
            editCurrency: item.currency,
            monthStartDate: currentMonthStartDate.value
        }
    });
}

function goToEdit(): void {
    props.f7router.navigate('/budget/edit', {
        props: {
            monthStartDate: currentMonthStartDate.value
        }
    });
}

function displayAmount(cents: number, currency: string): string {
    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(cents), currency);
}
</script>

<style scoped>
.budget-month-picker .month-picker-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
}
.budget-month-picker .month-label {
    font-weight: 600;
    font-size: var(--f7-list-item-title-font-size);
}
.budget-summary {
    background: var(--f7-list-bg-color);
    border-radius: var(--f7-list-inset-border-radius);
    padding: 12px 16px;
    margin-bottom: 8px;
}
.budget-summary .summary-row {
    display: flex;
    justify-content: space-between;
    padding: 4px 0;
}
.budget-summary .summary-row.total {
    border-top: 1px solid var(--f7-list-item-divider-border-color);
    margin-top: 4px;
    padding-top: 8px;
    font-weight: 600;
}
.budget-summary .summary-value.spent {
    color: var(--f7-color-red);
}
.budget-summary .summary-value.overspent {
    color: var(--f7-color-red);
    font-weight: 700;
}
.budget-item-content {
    width: 100%;
    padding-left: 8px;
}
.budget-item-content .budget-amounts {
    display: flex;
    gap: 4px;
    align-items: baseline;
    margin-top: 4px;
    font-size: var(--f7-list-item-after-font-size);
}
.budget-item-content .budget-amounts .spent {
    color: var(--f7-color-red);
    font-weight: 600;
}
.budget-item-content .budget-amounts .separator {
    color: var(--f7-list-item-after-text-color);
}
.budget-item-content .budget-amounts .limit {
    color: var(--f7-list-item-after-text-color);
}
.progress-overspent :deep(.progressbar-fill) {
    background: var(--f7-color-red) !important;
}
</style>
